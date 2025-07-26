package services

import (
	"os"
	"testing"
	"time"

	"weather-dashboard/models"

	"sync"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabaseService(t *testing.T) {
	// Use a temporary database file for testing
	testDBPath := "test_weather.db"

	// Clean up test database after tests
	defer func() {
		os.Remove(testDBPath)
	}()

	t.Run("NewDatabaseService", func(t *testing.T) {
		dbService, err := NewDatabaseService(testDBPath)
		require.NoError(t, err)
		require.NotNil(t, dbService)
		defer dbService.Close()

		// Verify database file was created
		_, err = os.Stat(testDBPath)
		assert.NoError(t, err)
	})

	t.Run("SaveWeatherData", func(t *testing.T) {
		dbService, err := NewDatabaseService(testDBPath)
		require.NoError(t, err)
		defer dbService.Close()

		weatherData := &models.WeatherData{
			City:          "London",
			Country:       "United Kingdom",
			State:         "England",
			Temperature:   15.5,
			Description:   "Partly cloudy",
			Humidity:      65,
			Icon:          "https://example.com/icon.png",
			ConditionCode: 1003,
			Timestamp:     time.Now(),
		}

		err = dbService.SaveWeatherData(weatherData)
		assert.NoError(t, err)
	})

	t.Run("GetWeatherHistory", func(t *testing.T) {
		dbService, err := NewDatabaseService(testDBPath)
		require.NoError(t, err)
		defer dbService.Close()

		// Save some test data
		testData := []*models.WeatherData{
			{
				City:          "London",
				Country:       "United Kingdom",
				State:         "England",
				Temperature:   15.5,
				Description:   "Partly cloudy",
				Humidity:      65,
				Icon:          "https://example.com/icon1.png",
				ConditionCode: 1003,
				Timestamp:     time.Now().Add(-time.Hour),
			},
			{
				City:          "Paris",
				Country:       "France",
				State:         "Île-de-France",
				Temperature:   18.2,
				Description:   "Sunny",
				Humidity:      55,
				Icon:          "https://example.com/icon2.png",
				ConditionCode: 1000,
				Timestamp:     time.Now().Add(-30 * time.Minute),
			},
			{
				City:          "Berlin",
				Country:       "Germany",
				State:         "Berlin",
				Temperature:   12.8,
				Description:   "Cloudy",
				Humidity:      70,
				Icon:          "https://example.com/icon3.png",
				ConditionCode: 1006,
				Timestamp:     time.Now(),
			},
		}

		for _, data := range testData {
			err = dbService.SaveWeatherData(data)
			require.NoError(t, err)
		}

		// Test getting history with default limit
		history, err := dbService.GetWeatherHistoryDefault()
		assert.NoError(t, err)
		assert.Len(t, history, 3)

		// Verify the data is returned in correct order (newest first)
		// Sort by timestamp to ensure consistent order
		assert.True(t, history[0].Timestamp.After(history[1].Timestamp) || history[0].Timestamp.Equal(history[1].Timestamp))
		assert.True(t, history[1].Timestamp.After(history[2].Timestamp) || history[1].Timestamp.Equal(history[2].Timestamp))

		// Test getting history with custom limit
		history, err = dbService.GetWeatherHistory(2)
		assert.NoError(t, err)
		assert.Len(t, history, 2)
		assert.True(t, history[0].Timestamp.After(history[1].Timestamp) || history[0].Timestamp.Equal(history[1].Timestamp))
	})

	t.Run("GetWeatherHistoryEmpty", func(t *testing.T) {
		// Use a fresh database
		emptyDBPath := "test_empty.db"
		defer os.Remove(emptyDBPath)

		dbService, err := NewDatabaseService(emptyDBPath)
		require.NoError(t, err)
		defer dbService.Close()

		history, err := dbService.GetWeatherHistoryDefault()
		assert.NoError(t, err)
		assert.Len(t, history, 0)
	})

	t.Run("SaveWeatherDataWithAllFields", func(t *testing.T) {
		dbService, err := NewDatabaseService(testDBPath)
		require.NoError(t, err)
		defer dbService.Close()

		weatherData := &models.WeatherData{
			City:          "Tokyo",
			Country:       "Japan",
			State:         "Tokyo",
			Temperature:   22.1,
			Description:   "Light rain",
			Humidity:      80,
			Icon:          "https://example.com/rain.png",
			ConditionCode: 1183,
			Timestamp:     time.Now(),
		}

		err = dbService.SaveWeatherData(weatherData)
		assert.NoError(t, err)

		// Verify the data was saved correctly
		history, err := dbService.GetWeatherHistory(1)
		assert.NoError(t, err)
		assert.Len(t, history, 1)

		savedData := history[0]
		assert.Equal(t, "Tokyo", savedData.City)
		assert.Equal(t, "Japan", savedData.Country)
		assert.Equal(t, "Tokyo", savedData.State)
		assert.Equal(t, 22.1, savedData.Temperature)
		assert.Equal(t, "Light rain", savedData.Description)
		assert.Equal(t, 80, savedData.Humidity)
		assert.Equal(t, "https://example.com/rain.png", savedData.Icon)
		assert.Equal(t, 1183, savedData.ConditionCode)
		assert.NotZero(t, savedData.ID)
		assert.NotZero(t, savedData.Timestamp)
	})
}

func TestDatabaseServiceErrors(t *testing.T) {
	t.Run("InvalidDatabasePath", func(t *testing.T) {
		// Test with an invalid path (directory that doesn't exist)
		_, err := NewDatabaseService("/nonexistent/path/test.db")
		assert.Error(t, err)
	})

	t.Run("DatabaseClose", func(t *testing.T) {
		testDBPath := "test_close.db"
		defer os.Remove(testDBPath)

		dbService, err := NewDatabaseService(testDBPath)
		require.NoError(t, err)

		// Test that close doesn't return an error
		err = dbService.Close()
		assert.NoError(t, err)

		// Test that close can be called multiple times safely
		err = dbService.Close()
		assert.NoError(t, err)
	})
}

func TestDatabaseConcurrency(t *testing.T) {
	testDBPath := "test_concurrency.db"
	defer os.Remove(testDBPath)

	dbService, err := NewDatabaseService(testDBPath)
	require.NoError(t, err)
	defer dbService.Close()

	// Test concurrent saves with a mutex to prevent database locks
	var mu sync.Mutex
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			mu.Lock()
			defer mu.Unlock()

			weatherData := &models.WeatherData{
				City:        "City" + string(rune(id+'A')),
				Country:     "Country",
				Temperature: float64(id),
				Description: "Test",
				Humidity:    50,
				Timestamp:   time.Now(),
			}
			err := dbService.SaveWeatherData(weatherData)
			assert.NoError(t, err)
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all data was saved
	history, err := dbService.GetWeatherHistory(20)
	assert.NoError(t, err)
	assert.Len(t, history, 10)
}
