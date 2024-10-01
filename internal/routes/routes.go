package routes

import (
	"github.com/gin-gonic/gin"
	"go-lastfm/internal/controllers" // Adjusted import path
	"go-lastfm/internal/services"    // Adjusted import path
)

func SetupRoutes(router *gin.Engine, authService *services.AuthService) {
	userService := services.NewUserService(authService)
	authController := controllers.NewAuthController(authService)
	userController := controllers.NewUserController(userService)

	router.GET("/auth", authController.AuthHandler)           // Endpoint to redirect user to Last.fm auth
	router.GET("/auth/callback", authController.CallbackHandler) // Last.fm callback handler
	router.GET("/user/get-info", userController.GetUserInfoHandler)
}
