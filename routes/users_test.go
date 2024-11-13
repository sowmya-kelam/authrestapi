package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"restapi/models"
	"restapi/utils"
	"testing"
	"fmt"
	//"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)


func mockSave(user *models.User) error {
	fmt.Println("user saved",user)
	return nil
}

func mockValidate(user *models.User) error {
	
	if user.Email == "test@example.com" && user.Password == "password" {
		return nil
	}
	
	return nil
}

func TestSignup(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.POST("/signup", func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		err := mockSave(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "user created"})
	})

	user := models.User{
		Email:    "test@example.com",
		Password: "password",
	}

	
	body, _ := json.Marshal(user)

	req, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusCreated, res.Code)
	assert.Contains(t, res.Body.String(), "user created")
}

func TestLogin(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Mock the handler for Login without any actual DB logic
	router.POST("/login", func(c *gin.Context) {
		// Mock the validation logic
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		err := mockValidate(&user)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}

		// Mock the successful response with fake tokens
		accesstoken := "fake-access-token"
		refreshtoken := "fake-refresh-token"
		c.JSON(http.StatusOK, gin.H{
			"message":    "Login successful",
			"accesstoken": accesstoken,
			"refreshtoken": refreshtoken,
		})
	})

	// Create a user for login testing
	user := models.User{
		Email:    "test@example.com",
		Password: "password",
	}

	// Marshal the user to JSON
	body, _ := json.Marshal(user)

	// Create a new POST request for the /login route
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the response
	rec := httptest.NewRecorder()

	// Perform the HTTP request
	router.ServeHTTP(rec, req)

	// Assert that the response is what we expect
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Login successful")
}

func TestAuthorizeToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/authorize", AuthorizeToken)

	token, _, _ := utils.GenerateToken("test@example.com", 1)
	tokenReq := models.TokenRequest{Token: "Bearer " + token}

	body, _ := json.Marshal(tokenReq)
	req, _ := http.NewRequest(http.MethodPost, "/authorize", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Vaild Token")
}

func TestRevokeToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/revoke", RevokeToken)

	token, _, _ := utils.GenerateToken("test@example.com", 1)
	tokenReq := models.TokenRequest{Token: "Bearer " + token}

	body, _ := json.Marshal(tokenReq)
	req, _ := http.NewRequest(http.MethodPost, "/revoke", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Token is Revoked")
}

func TestRefreshToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/refresh", RefreshToken)

	_, refreshToken, _ := utils.GenerateToken("test@example.com", 1)
	refreshTokenReq := models.RefreshTokenRequest{RefreshToken: refreshToken}

	body, _ := json.Marshal(refreshTokenReq)
	req, _ := http.NewRequest(http.MethodPost, "/refresh", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "New Token Generated")
}
