package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tyange/triplework-backend/db"
	"github.com/tyange/triplework-backend/routes"
)

func main() {
	db.InitDB()
	server := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true

	server.Use(cors.New(config))

	routes.RegisterRoutes(server)

	server.Run(":8080")
}
