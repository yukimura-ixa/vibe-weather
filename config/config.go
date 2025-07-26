package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Weather  WeatherConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port string
	Host string
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Path string
}

// WeatherConfig holds weather API-related configuration
type WeatherConfig struct {
	APIKey     string
	BaseURL    string
	SearchURL  string
	CurrentURL string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Warning: .env file not found, using system environment variables")
	} else {
		log.Println("Loaded .env file")
	}

	config := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Host: getEnv("HOST", "localhost"),
		},
		Database: DatabaseConfig{
			Path: getEnv("DB_PATH", "./weather.db"),
		},
		Weather: WeatherConfig{
			APIKey:     getEnv("WEATHERAPI_KEY", ""),
			BaseURL:    "http://api.weatherapi.com/v1",
			SearchURL:  "http://api.weatherapi.com/v1/search.json",
			CurrentURL: "http://api.weatherapi.com/v1/current.json",
		},
	}

	// Validate required configuration
	if config.Weather.APIKey == "" {
		return nil, fmt.Errorf("WEATHERAPI_KEY is required")
	}

	return config, nil
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// GetServerAddress returns the full server address
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)
}
