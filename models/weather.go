package models

import "time"

// WeatherData represents weather information for a location
type WeatherData struct {
	ID            int       `json:"id"`
	City          string    `json:"city"`
	Country       string    `json:"country"`
	State         string    `json:"state"`
	Temperature   float64   `json:"temperature"`
	Description   string    `json:"description"`
	Humidity      int       `json:"humidity"`
	Icon          string    `json:"icon"`
	ConditionCode int       `json:"condition_code"`
	Timestamp     time.Time `json:"timestamp"`
}

// WeatherAPISearchResult represents a search result from WeatherAPI
type WeatherAPISearchResult struct {
	Name    string  `json:"name"`
	Region  string  `json:"region"`
	Country string  `json:"country"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
}

// WeatherAPICurrentResult represents current weather data from WeatherAPI
type WeatherAPICurrentResult struct {
	Location struct {
		Name      string  `json:"name"`
		Region    string  `json:"region"`
		Country   string  `json:"country"`
		Lat       float64 `json:"lat"`
		Lon       float64 `json:"lon"`
		Localtime string  `json:"localtime"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
			Icon string `json:"icon"`
			Code int    `json:"code"`
		} `json:"condition"`
		Humidity    int    `json:"humidity"`
		LastUpdated string `json:"last_updated"`
	} `json:"current"`
}

// APIError represents an API error response
type APIError struct {
	Error string `json:"error"`
}

// Constants for weather condition codes
const (
	HistoryLimit = 3
)

// WeatherConditionCodes maps condition codes to descriptions
var WeatherConditionCodes = map[int]string{
	1000: "Clear",
	1003: "Partly cloudy",
	1006: "Cloudy",
	1009: "Overcast",
	1030: "Mist",
	1063: "Patchy rain possible",
	1066: "Patchy snow possible",
	1069: "Patchy sleet possible",
	1087: "Thundery outbreaks possible",
	1114: "Blowing snow",
	1117: "Blizzard",
	1135: "Fog",
	1147: "Freezing fog",
	1150: "Patchy light drizzle",
	1153: "Light drizzle",
	1168: "Freezing drizzle",
	1171: "Heavy freezing drizzle",
	1180: "Slight rain shower",
	1183: "Light rain",
	1186: "Moderate rain at times",
	1189: "Moderate rain",
	1192: "Heavy rain at times",
	1195: "Heavy rain",
	1198: "Light freezing rain",
	1201: "Moderate or heavy freezing rain",
	1204: "Light sleet",
	1207: "Moderate or heavy sleet",
	1210: "Patchy light snow",
	1213: "Light snow",
	1216: "Patchy moderate snow",
	1219: "Moderate snow",
	1222: "Patchy heavy snow",
	1225: "Heavy snow",
	1237: "Ice pellets",
	1240: "Light rain shower",
	1243: "Moderate or heavy rain shower",
	1246: "Torrential rain shower",
	1249: "Light sleet showers",
	1252: "Moderate or heavy sleet showers",
	1255: "Light snow showers",
	1258: "Moderate or heavy snow showers",
	1261: "Light showers of ice pellets",
	1264: "Moderate or heavy showers of ice pellets",
	1273: "Patchy light rain with thunder",
	1276: "Moderate or heavy rain with thunder",
	1279: "Patchy light snow with thunder",
	1282: "Moderate or heavy snow with thunder",
}

// GetWeatherConditionDescription returns a human-readable description based on condition code
func GetWeatherConditionDescription(code int) string {
	if description, exists := WeatherConditionCodes[code]; exists {
		return description
	}
	return "Unknown"
}
