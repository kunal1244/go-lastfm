package config

import (
	"encoding/json"
	"log"
	"os"
)

// Config struct to hold the API key, secret, and callback URL
type Config struct {
	APIKey       string `json:"lastfm_api_key"`
	SharedSecret string `json:"lastfm_shared_secret"`
	CallbackURL  string `json:"callback_url"`
}

// AppConfig holds the loaded configuration
var AppConfig Config

// LoadConfig loads configuration from a JSON file
func LoadConfig() {
	// Read the config.json file
	file, err := os.ReadFile("config/config.json")
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// Parse the JSON file into the config struct
	err = json.Unmarshal(file, &AppConfig)
	if err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}

	// Ensure API key and secret are loaded
	if AppConfig.APIKey == "" || AppConfig.SharedSecret == "" {
		log.Fatal("Missing API key or shared secret in config file")
	}
}
