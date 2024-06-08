package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	server.GET("/snooze", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"message": "Snooooze"})
	})

	server.Run(":8080")
}
