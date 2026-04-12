package routes

import (
	"auth_service/controllers"
	"auth_service/services"
)

func RegisterRoutes(router *Router) {

	authService := services.NewAuthService()
	authController := controllers.NewAuthController(authService)

	router.POST("/register", authController.Register)
	router.POST("/login", authController.Login)
	router.POST("/switch-role", authController.SwitchRole)
	router.POST("/internal/permissions", authController.SyncPermissions)
	router.POST("/internal/check", authController.CheckPermission)
}
