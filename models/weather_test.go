package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetWeatherConditionDescription(t *testing.T) {
	tests := []struct {
		name     string
		code     int
		expected string
	}{
		{
			name:     "clear weather",
			code:     1000,
			expected: "Clear",
		},
		{
			name:     "partly cloudy",
			code:     1003,
			expected: "Partly cloudy",
		},
		{
			name:     "cloudy",
			code:     1006,
			expected: "Cloudy",
		},
		{
			name:     "overcast",
			code:     1009,
			expected: "Overcast",
		},
		{
			name:     "mist",
			code:     1030,
			expected: "Mist",
		},
		{
			name:     "light rain",
			code:     1183,
			expected: "Light rain",
		},
		{
			name:     "heavy rain",
			code:     1195,
			expected: "Heavy rain",
		},
		{
			name:     "light snow",
			code:     1213,
			expected: "Light snow",
		},
		{
			name:     "heavy snow",
			code:     1225,
			expected: "Heavy snow",
		},
		{
			name:     "thunder",
			code:     1276,
			expected: "Moderate or heavy rain with thunder",
		},
		{
			name:     "unknown code",
			code:     9999,
			expected: "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetWeatherConditionDescription(tt.code)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestWeatherDataStruct(t *testing.T) {
	now := time.Now()
	weatherData := WeatherData{
		ID:            1,
		City:          "London",
		Country:       "United Kingdom",
		State:         "England",
		Temperature:   15.5,
		Description:   "Partly cloudy",
		Humidity:      65,
		Icon:          "https://example.com/icon.png",
		ConditionCode: 1003,
		Timestamp:     now,
	}

	assert.Equal(t, 1, weatherData.ID)
	assert.Equal(t, "London", weatherData.City)
	assert.Equal(t, "United Kingdom", weatherData.Country)
	assert.Equal(t, "England", weatherData.State)
	assert.Equal(t, 15.5, weatherData.Temperature)
	assert.Equal(t, "Partly cloudy", weatherData.Description)
	assert.Equal(t, 65, weatherData.Humidity)
	assert.Equal(t, "https://example.com/icon.png", weatherData.Icon)
	assert.Equal(t, 1003, weatherData.ConditionCode)
	assert.Equal(t, now, weatherData.Timestamp)
}

func TestWeatherAPISearchResultStruct(t *testing.T) {
	searchResult := WeatherAPISearchResult{
		Name:    "London",
		Region:  "England",
		Country: "United Kingdom",
		Lat:     51.5074,
		Lon:     -0.1278,
	}

	assert.Equal(t, "London", searchResult.Name)
	assert.Equal(t, "England", searchResult.Region)
	assert.Equal(t, "United Kingdom", searchResult.Country)
	assert.Equal(t, 51.5074, searchResult.Lat)
	assert.Equal(t, -0.1278, searchResult.Lon)
}

func TestWeatherAPICurrentResultStruct(t *testing.T) {
	currentResult := WeatherAPICurrentResult{}

	// Test Location struct
	currentResult.Location.Name = "London"
	currentResult.Location.Region = "England"
	currentResult.Location.Country = "United Kingdom"
	currentResult.Location.Lat = 51.5074
	currentResult.Location.Lon = -0.1278
	currentResult.Location.Localtime = "2023-01-01 12:00"

	// Test Current struct
	currentResult.Current.TempC = 15.5
	currentResult.Current.Condition.Text = "Partly cloudy"
	currentResult.Current.Condition.Icon = "//cdn.weatherapi.com/weather/64x64/day/116.png"
	currentResult.Current.Condition.Code = 1003
	currentResult.Current.Humidity = 65
	currentResult.Current.LastUpdated = "2023-01-01 12:00"

	assert.Equal(t, "London", currentResult.Location.Name)
	assert.Equal(t, "England", currentResult.Location.Region)
	assert.Equal(t, "United Kingdom", currentResult.Location.Country)
	assert.Equal(t, 51.5074, currentResult.Location.Lat)
	assert.Equal(t, -0.1278, currentResult.Location.Lon)
	assert.Equal(t, "2023-01-01 12:00", currentResult.Location.Localtime)

	assert.Equal(t, 15.5, currentResult.Current.TempC)
	assert.Equal(t, "Partly cloudy", currentResult.Current.Condition.Text)
	assert.Equal(t, "//cdn.weatherapi.com/weather/64x64/day/116.png", currentResult.Current.Condition.Icon)
	assert.Equal(t, 1003, currentResult.Current.Condition.Code)
	assert.Equal(t, 65, currentResult.Current.Humidity)
	assert.Equal(t, "2023-01-01 12:00", currentResult.Current.LastUpdated)
}

func TestAPIErrorStruct(t *testing.T) {
	apiError := APIError{
		Error: "City not found",
	}

	assert.Equal(t, "City not found", apiError.Error)
}

func TestWeatherConditionCodesMap(t *testing.T) {
	// Test that the map contains expected codes
	assert.Contains(t, WeatherConditionCodes, 1000)
	assert.Contains(t, WeatherConditionCodes, 1003)
	assert.Contains(t, WeatherConditionCodes, 1006)
	assert.Contains(t, WeatherConditionCodes, 1009)
	assert.Contains(t, WeatherConditionCodes, 1030)
	assert.Contains(t, WeatherConditionCodes, 1183)
	assert.Contains(t, WeatherConditionCodes, 1195)
	assert.Contains(t, WeatherConditionCodes, 1213)
	assert.Contains(t, WeatherConditionCodes, 1225)
	assert.Contains(t, WeatherConditionCodes, 1276)

	// Test that the map doesn't contain invalid codes
	assert.NotContains(t, WeatherConditionCodes, 9999)
	assert.NotContains(t, WeatherConditionCodes, -1)

	// Test some specific values
	assert.Equal(t, "Clear", WeatherConditionCodes[1000])
	assert.Equal(t, "Partly cloudy", WeatherConditionCodes[1003])
	assert.Equal(t, "Heavy rain", WeatherConditionCodes[1195])
	assert.Equal(t, "Heavy snow", WeatherConditionCodes[1225])
}

func TestConstants(t *testing.T) {
	assert.Equal(t, 3, HistoryLimit)
}
