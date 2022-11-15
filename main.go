package main

import (
	"os"
	"github.com/gin-gonic/gin"
	routes "stronka3d-backend/routes"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())
	routes.ProductRoutes(router)
	router.Use(middleware.Authentication())

	router.run(":" + port)
}