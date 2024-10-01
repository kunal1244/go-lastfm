package main

import (
	"go-lastfm/internal/config"   // Adjusted import path
	"go-lastfm/internal/routes"   // Adjusted import path
	"go-lastfm/internal/services" // Adjusted import path
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration (API key, secret, etc.)
	config.LoadConfig()

	// Initialize Gin router
	router := gin.Default()

	// Initialize services
	authService := services.NewAuthService()

	// Setup routes
	routes.SetupRoutes(router, authService)

	// Start the server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
