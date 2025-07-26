package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		expectError bool
		checkConfig func(*testing.T, *Config)
	}{
		{
			name: "valid configuration with all required fields",
			envVars: map[string]string{
				"WEATHERAPI_KEY": "test-api-key",
				"PORT":           "9090",
				"HOST":           "0.0.0.0",
				"DB_PATH":        "/tmp/test.db",
			},
			expectError: false,
			checkConfig: func(t *testing.T, cfg *Config) {
				assert.Equal(t, "test-api-key", cfg.Weather.APIKey)
				assert.Equal(t, "9090", cfg.Server.Port)
				assert.Equal(t, "0.0.0.0", cfg.Server.Host)
				assert.Equal(t, "/tmp/test.db", cfg.Database.Path)
				assert.Equal(t, "http://api.weatherapi.com/v1", cfg.Weather.BaseURL)
			},
		},
		{
			name: "missing API key should return error",
			envVars: map[string]string{
				"PORT":    "8080",
				"HOST":    "localhost",
				"DB_PATH": "./weather.db",
			},
			expectError: true,
		},
		{
			name: "default values when env vars not set",
			envVars: map[string]string{
				"WEATHERAPI_KEY": "test-api-key",
			},
			expectError: false,
			checkConfig: func(t *testing.T, cfg *Config) {
				assert.Equal(t, "8080", cfg.Server.Port)
				assert.Equal(t, "localhost", cfg.Server.Host)
				assert.Equal(t, "./weather.db", cfg.Database.Path)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables for test
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			// Clean up environment variables after test
			defer func() {
				for key := range tt.envVars {
					os.Unsetenv(key)
				}
			}()

			cfg, err := Load()

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, cfg)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, cfg)
				if tt.checkConfig != nil {
					tt.checkConfig(t, cfg)
				}
			}
		})
	}
}

func TestGetServerAddress(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Host: "localhost",
			Port: "8080",
		},
	}

	address := cfg.GetServerAddress()
	assert.Equal(t, "localhost:8080", address)
}

func TestGetEnv(t *testing.T) {
	// Test with existing environment variable
	os.Setenv("TEST_VAR", "test-value")
	defer os.Unsetenv("TEST_VAR")

	// This is an internal function, but we can test it indirectly through Load
	// or we can make it public for testing
}

func TestWeatherConfig(t *testing.T) {
	cfg := &Config{
		Weather: WeatherConfig{
			APIKey:     "test-key",
			BaseURL:    "http://api.weatherapi.com/v1",
			SearchURL:  "http://api.weatherapi.com/v1/search.json",
			CurrentURL: "http://api.weatherapi.com/v1/current.json",
		},
	}

	assert.Equal(t, "test-key", cfg.Weather.APIKey)
	assert.Equal(t, "http://api.weatherapi.com/v1", cfg.Weather.BaseURL)
	assert.Equal(t, "http://api.weatherapi.com/v1/search.json", cfg.Weather.SearchURL)
	assert.Equal(t, "http://api.weatherapi.com/v1/current.json", cfg.Weather.CurrentURL)
}
