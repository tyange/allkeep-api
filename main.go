package main

import (
	"github.com/gin-gonic/gin"
	"github.com/tyange/triplework-backend/db"
	"github.com/tyange/triplework-backend/routes"
)

func main() {
	db.InitDB()
	server := gin.Default()

	routes.RegisterRoutes(server)

	server.Run(":8080")
}
