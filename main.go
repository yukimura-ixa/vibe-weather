package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"weather-dashboard/config"
	"weather-dashboard/handlers"
	"weather-dashboard/services"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database service
	dbService, err := services.NewDatabaseService(cfg.Database.Path)
	if err != nil {
		log.Fatalf("Failed to initialize database service: %v", err)
	}
	defer dbService.Close()

	// Initialize weather service
	weatherService := services.NewWeatherService(&cfg.Weather)

	// Initialize handlers
	weatherHandler := handlers.NewWeatherHandler(weatherService, dbService)

	// Setup Gin router
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	// Setup routes
	setupRoutes(r, weatherHandler)

	// Start server
	log.Printf("Server starting on %s", cfg.GetServerAddress())
	if err := r.Run(cfg.GetServerAddress()); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// setupRoutes configures all application routes
func setupRoutes(r *gin.Engine, weatherHandler *handlers.WeatherHandler) {
	// Main page
	r.GET("/", weatherHandler.ServeIndex)

	// API routes
	api := r.Group("/api")
	{
		api.GET("/weather/:city", weatherHandler.GetWeatherByCity)
		api.GET("/weather/coordinates/:lat/:lon", weatherHandler.GetWeatherByCoordinates)
		api.GET("/history", weatherHandler.GetWeatherHistory)
	}
}
