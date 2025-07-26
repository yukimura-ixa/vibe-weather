package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"weather-dashboard/config"
	"weather-dashboard/handlers"
	"weather-dashboard/models"
	"weather-dashboard/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_WeatherFlow(t *testing.T) {
	// Setup test environment
	testDBPath := "test_integration.db"
	defer os.Remove(testDBPath)

	// Load test configuration
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port: "8080",
			Host: "localhost",
		},
		Database: config.DatabaseConfig{
			Path: testDBPath,
		},
		Weather: config.WeatherConfig{
			APIKey:     "test-key",
			BaseURL:    "http://api.weatherapi.com/v1",
			SearchURL:  "http://api.weatherapi.com/v1/search.json",
			CurrentURL: "http://api.weatherapi.com/v1/current.json",
		},
	}

	// Create test server for weather API
	weatherServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/search.json" {
			// Mock search response
			results := []models.WeatherAPISearchResult{
				{
					Name:    "London",
					Region:  "England",
					Country: "United Kingdom",
					Lat:     51.5074,
					Lon:     -0.1278,
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(results)
		} else if r.URL.Path == "/v1/current.json" {
			// Mock current weather response
			result := models.WeatherAPICurrentResult{}
			result.Location.Name = "London"
			result.Location.Region = "England"
			result.Location.Country = "United Kingdom"
			result.Location.Lat = 51.5074
			result.Location.Lon = -0.1278
			result.Current.TempC = 15.5
			result.Current.Condition.Text = "Partly cloudy"
			result.Current.Condition.Icon = "//cdn.weatherapi.com/weather/64x64/day/116.png"
			result.Current.Condition.Code = 1003
			result.Current.Humidity = 65

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(result)
		}
	}))
	defer weatherServer.Close()

	// Update config to use test server
	cfg.Weather.SearchURL = weatherServer.URL + "/v1/search.json"
	cfg.Weather.CurrentURL = weatherServer.URL + "/v1/current.json"

	// Initialize services
	dbService, err := services.NewDatabaseService(cfg.Database.Path)
	require.NoError(t, err)
	defer dbService.Close()

	weatherService := services.NewWeatherService(&cfg.Weather)

	// Initialize handler
	handler := handlers.NewWeatherHandler(weatherService, dbService)

	// Setup router
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/api/weather/:city", handler.GetWeatherByCity)
	r.GET("/api/weather/coordinates/:lat/:lon", handler.GetWeatherByCoordinates)
	r.GET("/api/history", handler.GetWeatherHistory)

	t.Run("CompleteWeatherFlow", func(t *testing.T) {
		// Step 1: Get weather for a city
		req, err := http.NewRequest("GET", "/api/weather/london", nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Assert weather response
		assert.Equal(t, http.StatusOK, w.Code)

		var weatherResponse models.WeatherData
		err = json.Unmarshal(w.Body.Bytes(), &weatherResponse)
		require.NoError(t, err)

		assert.Equal(t, "London", weatherResponse.City)
		assert.Equal(t, "England", weatherResponse.State)
		assert.Equal(t, "United Kingdom", weatherResponse.Country)
		assert.Equal(t, 15.5, weatherResponse.Temperature)
		assert.Equal(t, "Partly cloudy", weatherResponse.Description)
		assert.Equal(t, 65, weatherResponse.Humidity)
		assert.Equal(t, "https://cdn.weatherapi.com/weather/64x64/day/116.png", weatherResponse.Icon)
		assert.Equal(t, 1003, weatherResponse.ConditionCode)

		// Step 2: Check that data was saved to database
		req, err = http.NewRequest("GET", "/api/history", nil)
		require.NoError(t, err)

		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Assert history response
		assert.Equal(t, http.StatusOK, w.Code)

		var historyResponse []models.WeatherData
		err = json.Unmarshal(w.Body.Bytes(), &historyResponse)
		require.NoError(t, err)

		assert.Len(t, historyResponse, 1)
		assert.Equal(t, "London", historyResponse[0].City)
		assert.Equal(t, "England", historyResponse[0].State)
		assert.Equal(t, "United Kingdom", historyResponse[0].Country)
		assert.Equal(t, 15.5, historyResponse[0].Temperature)
		assert.Equal(t, "Partly cloudy", historyResponse[0].Description)
		assert.Equal(t, 65, historyResponse[0].Humidity)
		assert.Equal(t, "https://cdn.weatherapi.com/weather/64x64/day/116.png", historyResponse[0].Icon)
		assert.Equal(t, 1003, historyResponse[0].ConditionCode)
		assert.NotZero(t, historyResponse[0].ID)
		assert.NotZero(t, historyResponse[0].Timestamp)
	})

	t.Run("MultipleWeatherRequests", func(t *testing.T) {
		// Make multiple weather requests
		cities := []string{"london", "paris", "berlin"}

		for _, city := range cities {
			req, err := http.NewRequest("GET", "/api/weather/"+city, nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
		}

		// Check history has multiple entries
		req, err := http.NewRequest("GET", "/api/history", nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var historyResponse []models.WeatherData
		err = json.Unmarshal(w.Body.Bytes(), &historyResponse)
		require.NoError(t, err)

		// Should have at least 3 entries (from previous test + this test)
		assert.GreaterOrEqual(t, len(historyResponse), 3)
	})

	t.Run("WeatherByCoordinates", func(t *testing.T) {
		// Test weather by coordinates
		req, err := http.NewRequest("GET", "/api/weather/coordinates/51.5074/-0.1278", nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var weatherResponse models.WeatherData
		err = json.Unmarshal(w.Body.Bytes(), &weatherResponse)
		require.NoError(t, err)

		assert.Equal(t, "London", weatherResponse.City)
		assert.Equal(t, "England", weatherResponse.State)
		assert.Equal(t, "United Kingdom", weatherResponse.Country)
		assert.Equal(t, 15.5, weatherResponse.Temperature)
	})
}

func TestIntegration_ErrorHandling(t *testing.T) {
	// Setup test environment
	testDBPath := "test_error_integration.db"
	defer os.Remove(testDBPath)

	// Create test server that returns errors
	errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/search.json" {
			// Return empty results for city not found
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode([]models.WeatherAPISearchResult{})
		} else if r.URL.Path == "/v1/current.json" {
			// Return error status
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("API Error"))
		}
	}))
	defer errorServer.Close()

	// Load test configuration
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Path: testDBPath,
		},
		Weather: config.WeatherConfig{
			APIKey:     "test-key",
			SearchURL:  errorServer.URL + "/v1/search.json",
			CurrentURL: errorServer.URL + "/v1/current.json",
		},
	}

	// Initialize services
	dbService, err := services.NewDatabaseService(cfg.Database.Path)
	require.NoError(t, err)
	defer dbService.Close()

	weatherService := services.NewWeatherService(&cfg.Weather)

	// Initialize handler
	handler := handlers.NewWeatherHandler(weatherService, dbService)

	// Setup router
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/api/weather/:city", handler.GetWeatherByCity)
	r.GET("/api/weather/coordinates/:lat/:lon", handler.GetWeatherByCoordinates)

	t.Run("CityNotFound", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/weather/nonexistent", nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var errorResponse models.APIError
		err = json.Unmarshal(w.Body.Bytes(), &errorResponse)
		require.NoError(t, err)

		assert.Contains(t, errorResponse.Error, "city not found")
	})

	t.Run("InvalidCoordinates", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/weather/coordinates//", nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Gin returns 301 for missing parameters in URL
		assert.Equal(t, http.StatusMovedPermanently, w.Code)

		// Don't try to parse response body for redirect
	})
}

func TestIntegration_DatabasePersistence(t *testing.T) {
	// Test that data persists across service restarts
	testDBPath := "test_persistence.db"
	defer os.Remove(testDBPath)

	// First, create a service and save some data
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Path: testDBPath,
		},
		Weather: config.WeatherConfig{
			APIKey: "test-key",
		},
	}

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/search.json" {
			results := []models.WeatherAPISearchResult{
				{
					Name:    "London",
					Region:  "England",
					Country: "United Kingdom",
					Lat:     51.5074,
					Lon:     -0.1278,
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(results)
		} else if r.URL.Path == "/v1/current.json" {
			result := models.WeatherAPICurrentResult{}
			result.Location.Name = "London"
			result.Location.Region = "England"
			result.Location.Country = "United Kingdom"
			result.Current.TempC = 15.5
			result.Current.Condition.Text = "Partly cloudy"
			result.Current.Condition.Icon = "//cdn.weatherapi.com/weather/64x64/day/116.png"
			result.Current.Condition.Code = 1003
			result.Current.Humidity = 65

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(result)
		}
	}))
	defer server.Close()

	cfg.Weather.SearchURL = server.URL + "/v1/search.json"
	cfg.Weather.CurrentURL = server.URL + "/v1/current.json"

	// First service instance
	dbService1, err := services.NewDatabaseService(cfg.Database.Path)
	require.NoError(t, err)

	weatherService1 := services.NewWeatherService(&cfg.Weather)
	handler1 := handlers.NewWeatherHandler(weatherService1, dbService1)

	// Setup router
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/api/weather/:city", handler1.GetWeatherByCity)
	r.GET("/api/history", handler1.GetWeatherHistory)

	// Save some data
	req, err := http.NewRequest("GET", "/api/weather/london", nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Close first service
	dbService1.Close()

	// Create second service instance with same database
	dbService2, err := services.NewDatabaseService(cfg.Database.Path)
	require.NoError(t, err)
	defer dbService2.Close()

	weatherService2 := services.NewWeatherService(&cfg.Weather)
	handler2 := handlers.NewWeatherHandler(weatherService2, dbService2)

	// Setup new router
	r2 := gin.New()
	r2.GET("/api/history", handler2.GetWeatherHistory)

	// Check that data persists
	req, err = http.NewRequest("GET", "/api/history", nil)
	require.NoError(t, err)

	w = httptest.NewRecorder()
	r2.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var historyResponse []models.WeatherData
	err = json.Unmarshal(w.Body.Bytes(), &historyResponse)
	require.NoError(t, err)

	assert.Len(t, historyResponse, 1)
	assert.Equal(t, "London", historyResponse[0].City)
	assert.Equal(t, "England", historyResponse[0].State)
	assert.Equal(t, "United Kingdom", historyResponse[0].Country)
	assert.Equal(t, 15.5, historyResponse[0].Temperature)
}
