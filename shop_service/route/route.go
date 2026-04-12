package routes

import (
	"shop_service/controllers"
	"shop_service/middleware"
	"shop_service/services"
)

func RegisterRoutes(router *Router) {

	service := services.NewProductService()
	controller := controllers.NewProductController(service)

	router.GET("/products", middleware.AuthRequired(controller.GetAll))
	router.GET("/products/:id", middleware.AuthRequired(controller.GetByID))
	router.POST("/products", middleware.AuthRequired(controller.Create))
	router.PUT("/products/:id", middleware.AuthRequired(controller.Update))
	router.DELETE("/products/:id", middleware.AuthRequired(controller.Delete))
}
