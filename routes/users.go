package routes

import (
	"restapi/models"
	"restapi/utils"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/golang-jwt/jwt/v5"

	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var revokedTokens = make(map[string]bool)
var mu sync.Mutex
const secretkey = "supersecret"

// Signup User Handler godoc
// @Summary Signingup User
// @Description This endpoint allows user to sign up using email and password .
// @Tags Auth Rest Api's
// @Accept  json
// @Produce json
// @Param user body models.User  true "User Details"
// @Success 201 {object} map[string]string "message: user created"
// @Failure 400 {object} map[string]string "Error: Bad Request"
// @Failure 500 {object} map[string]string "Error: Internal Server Error"
// @Router /signup [post]
func Signup(context *gin.Context, database *sql.DB) {
	var user models.User
	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	fmt.Println(user)
	err = user.Save(database)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	context.JSON(http.StatusCreated, gin.H{"message": "user created "})
}

// Login User Handler godoc
// @Summary Login User
// @Description This endpoint allows user to Login using given email and password .
// @Tags Auth Rest Api's
// @Accept  json
// @Produce json
// @Param user body models.User  true "User Details"
// @Success 200 {object} map[string]string "message: Login successful"
// @Failure 400 {object} map[string]string "Error: Bad Request"
// @Failure 401 {object} map[string]string "Error: Unauthorized"
// @Failure 500 {object} map[string]string "Error: Internal Server Error"
// @Router /login [post]
func Login(context *gin.Context, database *sql.DB) {
	var user models.User
	err := context.ShouldBindJSON(&user)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = user.Validate(database)

	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	accesstoken,refreshtoken, err := utils.GenerateToken(user.Email, user.ID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Login successful", "accesstoken": accesstoken,"refreshtoken":refreshtoken})
}

// Authorize Handler godoc
// @Summary Authorize Token
// @Description This endpoint allows user to Authorize access token .
// @Tags Auth Rest Api's
// @Accept  json
// @Produce json
// @Param token body models.TokenRequest  true "Access Token"
// @Success 200 {object} map[string]string "message: valid token"
// @Failure 400 {object} map[string]string "Error: Bad Request"
// @Failure 401 {object} map[string]string "Error: unauthorized"
// @Router /authorize-token [post]
func AuthorizeToken(context *gin.Context) {
	var token models.TokenRequest
	err := context.ShouldBindJSON(&token)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if token.Token == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "token is empty"})
		return
	}

	_, err = utils.VerifyToken(token.Token)

	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	if IsRevokedToken(token.Token) {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Revoked Token"})
		return
	}
	
	

	context.JSON(http.StatusOK, gin.H{"message": " Vaild Token "})
}

// Revoke Token Handler godoc
// @Summary Revoke Token 
// @Description This endpoint allows user to revoke the access token .
// @Tags Auth Rest Api's
// @Accept  json
// @Produce json
// @Param token body models.TokenRequest  true "Access Token"
// @Success 200 {object} map[string]string "message: Token is Revoked "
// @Failure 400 {object} map[string]string "Error: Bad Request"
// @Router /revoke-token [post]
func RevokeToken(context *gin.Context) {
	var token models.TokenRequest
	err := context.ShouldBindJSON(&token)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if token.Token == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "token is empty"})
		return
	}

	if IsRevokedToken(token.Token) {
		context.JSON(http.StatusBadRequest, gin.H{"message": "token is already revoked"})
		return
	}
	
	mu.Lock()
    revokedTokens[token.Token] = true 
    mu.Unlock()

	context.JSON(http.StatusOK, gin.H{"message": " Token is Revoked "})
}


func IsRevokedToken(token string) bool {
	mu.Lock()
    defer mu.Unlock()
    _, exists := revokedTokens[token]
    return exists
}

// Refresh Token Handler godoc
// @Summary  Refresh Token 
// @Description This endpoint allows user to Get New Accesstoken using refresh token .
// @Tags Auth Rest Api's
// @Accept  json
// @Produce json
// @Param refreshtoken body models.RefreshTokenRequest  true "Refresh token request"
// @Success 200 {object} map[string]string "message: New Token Generated"
// @Failure 400 {object} map[string]string "Error: Bad Request"
// @Failure 500 {object} map[string]string "Error: Internal Server Error"
// @Router /refresh-token [post]
func RefreshToken(context *gin.Context) {

	var refreshtoken models.RefreshTokenRequest
	err := context.ShouldBindJSON(&refreshtoken)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if refreshtoken.RefreshToken == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "token is empty"})
		return
	}

	if IsRevokedToken(refreshtoken.RefreshToken) {
		context.JSON(http.StatusBadRequest, gin.H{"message": "token is revoked, login again to get new access token"})
		return
	}

	parsedToken, err := jwt.Parse(refreshtoken.RefreshToken, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, err
        }
        return []byte(secretkey), nil
    })

    if err != nil || !parsedToken.Valid {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
    }

    claims, ok := parsedToken.Claims.(jwt.MapClaims)
    if !ok || !parsedToken.Valid {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
    }

    userId, ok := claims["userId"].(float64)
    if !ok {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
    }

    email, ok := claims["email"].(string)
    if !ok {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
    }
	
	newAccessToken,_, err := utils.GenerateToken(string(email),int64(userId))

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message":"New Token Generated","Token": newAccessToken})
}
