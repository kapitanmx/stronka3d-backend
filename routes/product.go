package routes

import (
	"stronka3d-backend/controllers"
	"github.com/gin-gonic/gin"
)

func ProductRoutes(router *gin.Engine) {
	router.POST("/products/", controllers.CreateProduct())
	router.GET("/products/:_id", controllers.GetProduct())
	router.GET("/products/", controllers.GetAllProducts())
	router.DELETE("/products/:_id", controller.DeleteProduct())
}