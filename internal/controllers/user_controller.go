package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-lastfm/internal/services" // Adjust this import path
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{userService: userService}
}

// GetUserInfoHandler handles the GET request for user info
func (uc *UserController) GetUserInfoHandler(c *gin.Context) {
	userInfo, err := uc.userService.GetUserInfo() // Get user info using session key
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userInfo)
}
