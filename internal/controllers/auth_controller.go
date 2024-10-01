package controllers

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"go-lastfm/internal/config"          // Adjusted import path
	"go-lastfm/internal/services" // Adjusted import path
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// Redirects to Last.fm for authentication
func (ac *AuthController) AuthHandler(c *gin.Context) {
	apiKey := config.AppConfig.APIKey
	callbackURL := "http://localhost:8080/auth/callback" // Replace with your app's actual callback URL

	authURL := fmt.Sprintf("https://www.last.fm/api/auth/?api_key=%s&cb=%s", apiKey, callbackURL)
	c.Redirect(http.StatusFound, authURL)
}

// CallbackHandler simulates the Last.fm callback where the user gets authenticated
func (ac *AuthController) CallbackHandler(c *gin.Context) {
	token := c.Query("token") // Get token from query parameters

	// Authenticate the user using the token
	if err := ac.authService.Authenticate(token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Authentication successful"})
}
