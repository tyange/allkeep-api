package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tyange/white-shadow-api/middlewares"
)

func RegisterRoutes(server *gin.Engine) {
	work := server.Group("/work")
	work.Use(middlewares.Authenticate)
	work.POST("/save", workSave)

	auth := server.Group("/auth")
	auth.POST("/signup", signup)
	auth.POST("/login", login)
}
