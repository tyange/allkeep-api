package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	work := server.Group("/work")
	work.POST("/start", workStart)
}
