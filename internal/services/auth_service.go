package services

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"
	"go-lastfm/internal/config" // Adjusted import path
)

type AuthService struct {
	sessionKey string
}

// NewAuthService initializes a new AuthService.
func NewAuthService() *AuthService {
	return &AuthService{}
}

// Authenticate exchanges the token from Last.fm for a session key
func (s *AuthService) Authenticate(token string) error {
	apiKey := config.AppConfig.APIKey
	apiSecret := config.AppConfig.SharedSecret

	// Build params
	params := url.Values{}
	params.Set("method", "auth.getSession")
	params.Set("api_key", apiKey)
	params.Set("token", token)

	// Generate the api_sig by concatenating the parameters in the correct order
	apiSig := generateApiSig(map[string]string{
		"method":  "auth.getSession",
		"api_key": apiKey,
		"token":   token,
	}, apiSecret)
	params.Set("api_sig", apiSig)
	params.Set("format", "json")

	// Create request URL
	requestURL := fmt.Sprintf("https://ws.audioscrobbler.com/2.0/?%s", params.Encode())

	// Make the request
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Println("Error making request:", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to get session key, status: %s, body: %s", resp.Status, string(body))
	}

	// Parse response
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return err
	}

	// Get session key from response
	session, ok := result["session"].(map[string]interface{})
	if !ok {
		return errors.New("failed to parse session response")
	}
	s.sessionKey, ok = session["key"].(string)
	if !ok {
		return errors.New("no session key found in response")
	}

	log.Println("Successfully authenticated. Session Key:", s.sessionKey)
	return nil
}

// GetSessionKey returns the session key
func (s *AuthService) GetSessionKey() string {
	return s.sessionKey
}

// Helper function to generate the API signature according to Last.fm rules
func generateApiSig(params map[string]string, apiSecret string) string {
	// Sort the keys of the parameters
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build the concatenated string of keys and values
	var signString string
	for _, key := range keys {
		signString += key + params[key]
	}

	// Append the API secret
	signString += apiSecret

	// Generate the MD5 hash
	hash := md5.Sum([]byte(signString))
	return hex.EncodeToString(hash[:])
}
