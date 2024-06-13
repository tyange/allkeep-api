package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tyange/white-shadow-api/middlewares"
)

func RegisterRoutes(server *gin.Engine) {
	work := server.Group("/work")
	work.Use(middlewares.Authenticate)
	work.POST("/create", createWork)

	auth := server.Group("/auth")
	auth.POST("/signup", signup)
	auth.POST("/login", login)

	company := server.Group("/company")
	company.Use(middlewares.Authenticate)
	company.POST("/create", createCompany)
}
