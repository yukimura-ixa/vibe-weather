package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"weather-dashboard/models"
)

// Define interfaces for dependency injection

type WeatherServiceInterface interface {
	SearchCity(city string) ([]models.WeatherAPISearchResult, error)
	GetWeatherByCoordinates(lat, lon string) (*models.WeatherData, error)
	GetWeatherByCity(city string) (*models.WeatherData, error)
}

type DatabaseServiceInterface interface {
	SaveWeatherData(data *models.WeatherData) error
	GetWeatherHistory(limit int) ([]models.WeatherData, error)
	GetWeatherHistoryDefault() ([]models.WeatherData, error)
	Close() error
}

// WeatherHandler handles weather-related HTTP requests
type WeatherHandler struct {
	weatherService WeatherServiceInterface
	dbService      DatabaseServiceInterface
}

// NewWeatherHandler creates a new weather handler
func NewWeatherHandler(weatherService WeatherServiceInterface, dbService DatabaseServiceInterface) *WeatherHandler {
	return &WeatherHandler{
		weatherService: weatherService,
		dbService:      dbService,
	}
}

// GetWeatherByCity handles GET /api/weather/:city
func (h *WeatherHandler) GetWeatherByCity(c *gin.Context) {
	city := c.Param("city")
	if city == "" {
		c.JSON(http.StatusBadRequest, models.APIError{Error: "city parameter is required"})
		return
	}

	weatherData, err := h.weatherService.GetWeatherByCity(city)
	if err != nil {
		log.Printf("Error fetching weather for %s: %v", city, err)
		c.JSON(http.StatusInternalServerError, models.APIError{Error: err.Error()})
		return
	}

	// Save to database
	if err := h.dbService.SaveWeatherData(weatherData); err != nil {
		log.Printf("Error saving weather data: %v", err)
		// Don't return error to client, just log it
	}

	c.JSON(http.StatusOK, weatherData)
}

// GetWeatherByCoordinates handles GET /api/weather/coordinates/:lat/:lon
func (h *WeatherHandler) GetWeatherByCoordinates(c *gin.Context) {
	lat := c.Param("lat")
	lon := c.Param("lon")

	if lat == "" || lon == "" {
		c.JSON(http.StatusBadRequest, models.APIError{Error: "latitude and longitude parameters are required"})
		return
	}

	weatherData, err := h.weatherService.GetWeatherByCoordinates(lat, lon)
	if err != nil {
		log.Printf("Error fetching weather for coordinates %s,%s: %v", lat, lon, err)
		c.JSON(http.StatusInternalServerError, models.APIError{Error: err.Error()})
		return
	}

	// Save to database
	if err := h.dbService.SaveWeatherData(weatherData); err != nil {
		log.Printf("Error saving weather data: %v", err)
		// Don't return error to client, just log it
	}

	c.JSON(http.StatusOK, weatherData)
}

// GetWeatherHistory handles GET /api/history
func (h *WeatherHandler) GetWeatherHistory(c *gin.Context) {
	history, err := h.dbService.GetWeatherHistoryDefault()
	if err != nil {
		log.Printf("Error fetching weather history: %v", err)
		c.JSON(http.StatusInternalServerError, models.APIError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}

// ServeIndex handles GET /
func (h *WeatherHandler) ServeIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}
