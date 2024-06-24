package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tyange/white-shadow-api/middlewares"
)

func RegisterRoutes(server *gin.Engine) {
	// TODO: get work
	// TODO: edit work
	// TODO: delete work
	work := server.Group("/work")
	work.Use(middlewares.Authenticate)
	work.POST("/create", createWork)

	// TODO: get user info
	// TODO: edit user info
	// TODO: delete user
	auth := server.Group("/auth")
	auth.POST("/signup", signup)
	auth.POST("/login", login)
	auth.POST("/google", googleLoginCallBack)

	// TODO: get company
	// TODO: edit company
	// TODO: delete company
	company := server.Group("/company")
	company.Use(middlewares.Authenticate)
	company.POST("/create", createCompany)
}
