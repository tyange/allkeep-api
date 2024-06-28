package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tyange/white-shadow-api/db"
	"github.com/tyange/white-shadow-api/routes"
)

func main() {
	db.InitDB()
	server := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"}
	config.AllowHeaders = []string{
		"Authorization",
		"Content-Type",
	}
	config.AllowCredentials = true

	server.Use(cors.New(config))

	routes.RegisterRoutes(server)

	server.Run(":8080")
}
