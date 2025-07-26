package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"weather-dashboard/config"
	"weather-dashboard/handlers"
	"weather-dashboard/models"
)

// WeatherService handles weather API operations
type WeatherService struct {
	config *config.WeatherConfig
	client *http.Client
}

// NewWeatherService creates a new weather service
func NewWeatherService(cfg *config.WeatherConfig) *WeatherService {
	return &WeatherService{
		config: cfg,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SearchCity searches for a city using WeatherAPI
func (s *WeatherService) SearchCity(city string) ([]models.WeatherAPISearchResult, error) {
	baseURL := s.config.SearchURL
	params := url.Values{}
	params.Add("key", s.config.APIKey)
	params.Add("q", city)

	requestURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	resp, err := s.client.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to search city: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var results []models.WeatherAPISearchResult
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search results: %w", err)
	}

	return results, nil
}

// GetWeatherByCoordinates fetches weather data for given coordinates
func (s *WeatherService) GetWeatherByCoordinates(lat, lon string) (*models.WeatherData, error) {
	baseURL := s.config.CurrentURL
	params := url.Values{}
	params.Add("key", s.config.APIKey)
	params.Add("q", fmt.Sprintf("%s,%s", lat, lon))

	requestURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	resp, err := s.client.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get weather by coordinates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result models.WeatherAPICurrentResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal weather data: %w", err)
	}

	return s.transformWeatherData(&result), nil
}

// GetWeatherByCity fetches weather data for a city
func (s *WeatherService) GetWeatherByCity(city string) (*models.WeatherData, error) {
	// First search for the city to get coordinates
	results, err := s.SearchCity(city)
	if err != nil {
		return nil, fmt.Errorf("failed to search city: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("city not found: %s", city)
	}

	// Use the first result to get weather data
	lat := fmt.Sprintf("%f", results[0].Lat)
	lon := fmt.Sprintf("%f", results[0].Lon)

	return s.GetWeatherByCoordinates(lat, lon)
}

// transformWeatherData transforms API response to our model
func (s *WeatherService) transformWeatherData(result *models.WeatherAPICurrentResult) *models.WeatherData {
	return &models.WeatherData{
		City:          result.Location.Name,
		Country:       result.Location.Country,
		State:         result.Location.Region,
		Temperature:   result.Current.TempC,
		Description:   result.Current.Condition.Text,
		Humidity:      result.Current.Humidity,
		Icon:          "https:" + result.Current.Condition.Icon,
		ConditionCode: result.Current.Condition.Code,
		Timestamp:     time.Now(),
	}
}

// Ensure WeatherService implements handlers.WeatherServiceInterface
var _ handlers.WeatherServiceInterface = (*WeatherService)(nil)
