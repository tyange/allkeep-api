package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	work := server.Group("/work")
	work.POST("/save", workSave)

	auth := server.Group("/auth")
	auth.POST("/signup", signup)
	auth.POST("/login", login)
}
