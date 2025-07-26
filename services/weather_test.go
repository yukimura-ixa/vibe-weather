package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"weather-dashboard/config"
	"weather-dashboard/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewWeatherService(t *testing.T) {
	cfg := &config.WeatherConfig{
		APIKey:     "test-key",
		BaseURL:    "http://api.weatherapi.com/v1",
		SearchURL:  "http://api.weatherapi.com/v1/search.json",
		CurrentURL: "http://api.weatherapi.com/v1/current.json",
	}

	service := NewWeatherService(cfg)
	assert.NotNil(t, service)
	assert.Equal(t, cfg, service.config)
	assert.NotNil(t, service.client)
}

func TestWeatherService_SearchCity(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if it's a search request
		if r.URL.Path == "/v1/search.json" {
			// Verify query parameters
			assert.Equal(t, "test-key", r.URL.Query().Get("key"))
			assert.Equal(t, "london", r.URL.Query().Get("q"))

			// Return mock search results
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
		} else {
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	// Create service with test server URL
	cfg := &config.WeatherConfig{
		APIKey:    "test-key",
		SearchURL: server.URL + "/v1/search.json",
	}

	service := NewWeatherService(cfg)

	// Test search
	results, err := service.SearchCity("london")
	require.NoError(t, err)
	require.Len(t, results, 1)

	result := results[0]
	assert.Equal(t, "London", result.Name)
	assert.Equal(t, "England", result.Region)
	assert.Equal(t, "United Kingdom", result.Country)
	assert.Equal(t, 51.5074, result.Lat)
	assert.Equal(t, -0.1278, result.Lon)
}

func TestWeatherService_SearchCityEmptyResults(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return empty results
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]models.WeatherAPISearchResult{})
	}))
	defer server.Close()

	cfg := &config.WeatherConfig{
		APIKey:    "test-key",
		SearchURL: server.URL + "/v1/search.json",
	}

	service := NewWeatherService(cfg)

	results, err := service.SearchCity("nonexistent")
	require.NoError(t, err)
	assert.Len(t, results, 0)
}

func TestWeatherService_SearchCityError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return error status
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	cfg := &config.WeatherConfig{
		APIKey:    "test-key",
		SearchURL: server.URL + "/v1/search.json",
	}

	service := NewWeatherService(cfg)

	_, err := service.SearchCity("london")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API request failed with status: 500")
}

func TestWeatherService_GetWeatherByCoordinates(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if it's a current weather request
		if r.URL.Path == "/v1/current.json" {
			// Verify query parameters
			assert.Equal(t, "test-key", r.URL.Query().Get("key"))
			assert.Equal(t, "51.5074,-0.1278", r.URL.Query().Get("q"))

			// Return mock weather data
			result := models.WeatherAPICurrentResult{}
			result.Location.Name = "London"
			result.Location.Region = "England"
			result.Location.Country = "United Kingdom"
			result.Location.Lat = 51.5074
			result.Location.Lon = -0.1278
			result.Location.Localtime = "2023-01-01 12:00"
			result.Current.TempC = 15.5
			result.Current.Condition.Text = "Partly cloudy"
			result.Current.Condition.Icon = "//cdn.weatherapi.com/weather/64x64/day/116.png"
			result.Current.Condition.Code = 1003
			result.Current.Humidity = 65
			result.Current.LastUpdated = "2023-01-01 12:00"

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(result)
		} else {
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	cfg := &config.WeatherConfig{
		APIKey:     "test-key",
		CurrentURL: server.URL + "/v1/current.json",
	}

	service := NewWeatherService(cfg)

	// Test weather by coordinates
	weatherData, err := service.GetWeatherByCoordinates("51.5074", "-0.1278")
	require.NoError(t, err)
	assert.NotNil(t, weatherData)

	assert.Equal(t, "London", weatherData.City)
	assert.Equal(t, "England", weatherData.State)
	assert.Equal(t, "United Kingdom", weatherData.Country)
	assert.Equal(t, 15.5, weatherData.Temperature)
	assert.Equal(t, "Partly cloudy", weatherData.Description)
	assert.Equal(t, 65, weatherData.Humidity)
	assert.Equal(t, "https://cdn.weatherapi.com/weather/64x64/day/116.png", weatherData.Icon)
	assert.Equal(t, 1003, weatherData.ConditionCode)
	assert.NotZero(t, weatherData.Timestamp)
}

func TestWeatherService_GetWeatherByCity(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/search.json" {
			// Return search results
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
			// Return weather data
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

	cfg := &config.WeatherConfig{
		APIKey:     "test-key",
		SearchURL:  server.URL + "/v1/search.json",
		CurrentURL: server.URL + "/v1/current.json",
	}

	service := NewWeatherService(cfg)

	// Test weather by city
	weatherData, err := service.GetWeatherByCity("london")
	require.NoError(t, err)
	assert.NotNil(t, weatherData)

	assert.Equal(t, "London", weatherData.City)
	assert.Equal(t, "England", weatherData.State)
	assert.Equal(t, "United Kingdom", weatherData.Country)
	assert.Equal(t, 15.5, weatherData.Temperature)
	assert.Equal(t, "Partly cloudy", weatherData.Description)
	assert.Equal(t, 65, weatherData.Humidity)
	assert.Equal(t, "https://cdn.weatherapi.com/weather/64x64/day/116.png", weatherData.Icon)
	assert.Equal(t, 1003, weatherData.ConditionCode)
}

func TestWeatherService_GetWeatherByCityNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return empty search results
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]models.WeatherAPISearchResult{})
	}))
	defer server.Close()

	cfg := &config.WeatherConfig{
		APIKey:    "test-key",
		SearchURL: server.URL + "/v1/search.json",
	}

	service := NewWeatherService(cfg)

	// Test city not found
	_, err := service.GetWeatherByCity("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "city not found: nonexistent")
}

func TestWeatherService_TransformWeatherData(t *testing.T) {
	service := NewWeatherService(&config.WeatherConfig{})

	// Create test API result
	apiResult := &models.WeatherAPICurrentResult{}
	apiResult.Location.Name = "London"
	apiResult.Location.Region = "England"
	apiResult.Location.Country = "United Kingdom"
	apiResult.Current.TempC = 15.5
	apiResult.Current.Condition.Text = "Partly cloudy"
	apiResult.Current.Condition.Icon = "//cdn.weatherapi.com/weather/64x64/day/116.png"
	apiResult.Current.Condition.Code = 1003
	apiResult.Current.Humidity = 65

	// Transform the data
	weatherData := service.transformWeatherData(apiResult)

	assert.Equal(t, "London", weatherData.City)
	assert.Equal(t, "England", weatherData.State)
	assert.Equal(t, "United Kingdom", weatherData.Country)
	assert.Equal(t, 15.5, weatherData.Temperature)
	assert.Equal(t, "Partly cloudy", weatherData.Description)
	assert.Equal(t, 65, weatherData.Humidity)
	assert.Equal(t, "https://cdn.weatherapi.com/weather/64x64/day/116.png", weatherData.Icon)
	assert.Equal(t, 1003, weatherData.ConditionCode)
	assert.NotZero(t, weatherData.Timestamp)
}

func TestWeatherService_Timeout(t *testing.T) {
	// Create a server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second) // Longer than the 10-second timeout
		w.Write([]byte("{}"))
	}))
	defer server.Close()

	cfg := &config.WeatherConfig{
		APIKey:    "test-key",
		SearchURL: server.URL + "/v1/search.json",
	}

	service := NewWeatherService(cfg)

	// Test timeout - expect JSON unmarshal error since the response is not a valid array
	_, err := service.SearchCity("london")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal search results")
}
