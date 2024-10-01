package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"go-lastfm/internal/config" // Adjust this import path
)

type UserInfo struct {
	User struct {
		Name         string `json:"name"`
		Playcount    string `json:"playcount"`
		PlaycountInt int    // Store integer representation of playcount
	} `json:"user"`
}

type UserService struct {
	authService *AuthService // Use AuthService to retrieve session key
}

// NewUserService creates a new instance of UserService
func NewUserService(authService *AuthService) *UserService {
	return &UserService{authService: authService}
}

// GetUserInfo fetches user information using the session key
func (us *UserService) GetUserInfo() (UserInfo, error) {
	sessionKey := us.authService.GetSessionKey() // Get stored session key
	if sessionKey == "" {
		return UserInfo{}, fmt.Errorf("session key not found")
	}

	apiKey := config.AppConfig.APIKey
	baseURL := "https://ws.audioscrobbler.com/2.0/"
	params := url.Values{}
	params.Add("method", "user.getInfo")
	params.Add("api_key", apiKey)
	params.Add("sk", sessionKey)
	params.Add("format", "json")

	resp, err := http.Get(fmt.Sprintf("%s?%s", baseURL, params.Encode()))
	if err != nil {
		return UserInfo{}, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return UserInfo{}, fmt.Errorf("failed to get user info, status: %s, response: %s", resp.Status, body)
	}

	var userInfo UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return UserInfo{}, fmt.Errorf("failed to decode user info: %w", err)
	}

	// Convert playcount from string to int
	playcount, err := strconv.Atoi(userInfo.User.Playcount)
	if err != nil {
		return UserInfo{}, fmt.Errorf("failed to convert playcount to int: %w", err)
	}
	userInfo.User.PlaycountInt = playcount // Store the integer value

	return userInfo, nil
}
