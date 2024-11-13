package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/files"

	_ "restapi/docs"

	_ "github.com/mattn/go-sqlite3"

)


func RegisterRoutes(server *gin.Engine, database *sql.DB) {

	server.POST("/signup", func(context *gin.Context) {
		Signup(context, database)
	})

	server.POST("/login", func(context *gin.Context) {
		Login(context, database)
	})

	server.POST("/authorize-token", func(context *gin.Context) {
		AuthorizeToken(context)
	})
	server.POST("/revoke-token", func(context *gin.Context) {
		RevokeToken(context)
	})
	
	server.POST("/refresh-token", func(context *gin.Context) {
		RefreshToken(context)
	})

	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}