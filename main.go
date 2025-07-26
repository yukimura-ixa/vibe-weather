package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

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

type WeatherAPISearchResult struct {
	Name    string  `json:"name"`
	Region  string  `json:"region"`
	Country string  `json:"country"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
}

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

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("sqlite", "./weather.db")
	if err != nil {
		log.Fatal(err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS weather_data (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		city TEXT NOT NULL,
		country TEXT,
		state TEXT,
		temperature REAL NOT NULL,
		description TEXT NOT NULL,
		humidity INTEGER NOT NULL,
		icon TEXT,
		condition_code INTEGER,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}
}

func searchCityWeatherAPI(city string) ([]WeatherAPISearchResult, error) {
	apiKey := os.Getenv("WEATHERAPI_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("API key not found")
	}
	url := fmt.Sprintf("http://api.weatherapi.com/v1/search.json?key=%s&q=%s", apiKey, city)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var results []WeatherAPISearchResult
	err = json.Unmarshal(body, &results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

// getWeatherConditionDescription returns a human-readable description based on condition code
func getWeatherConditionDescription(code int) string {
	conditions := map[int]string{
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

	if description, exists := conditions[code]; exists {
		return description
	}
	return "Unknown"
}

func fetchWeatherByCoordinates(lat, lon string) (*WeatherData, error) {
	apiKey := os.Getenv("WEATHERAPI_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("API key not found")
	}
	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s,%s", apiKey, lat, lon)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result WeatherAPICurrentResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	weatherData := &WeatherData{
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
	return weatherData, nil
}

func fetchWeatherData(city string) (*WeatherData, error) {
	results, err := searchCityWeatherAPI(city)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("city not found: %s", city)
	}
	// Use the first result
	lat := fmt.Sprintf("%f", results[0].Lat)
	lon := fmt.Sprintf("%f", results[0].Lon)
	return fetchWeatherByCoordinates(lat, lon)
}

func saveWeatherData(data *WeatherData) error {
	_, err := db.Exec("INSERT INTO weather_data (city, country, state, temperature, description, humidity, icon, condition_code, timestamp) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		data.City, data.Country, data.State, data.Temperature, data.Description, data.Humidity, data.Icon, data.ConditionCode, data.Timestamp)
	return err
}

func getWeatherHistory() ([]WeatherData, error) {
	rows, err := db.Query("SELECT id, city, country, state, temperature, description, humidity, icon, condition_code, timestamp FROM weather_data ORDER BY timestamp DESC LIMIT 3")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []WeatherData
	for rows.Next() {
		var data WeatherData
		err := rows.Scan(&data.ID, &data.City, &data.Country, &data.State, &data.Temperature, &data.Description, &data.Humidity, &data.Icon, &data.ConditionCode, &data.Timestamp)
		if err != nil {
			return nil, err
		}
		history = append(history, data)
	}
	return history, nil
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Warning: .env file not found, using system environment variables")
	}

	initDB()
	defer db.Close()

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/api/weather/:city", func(c *gin.Context) {
		city := c.Param("city")
		weatherData, err := fetchWeatherData(city)
		if err != nil {
			log.Printf("Error fetching weather for %s: %v", city, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		err = saveWeatherData(weatherData)
		if err != nil {
			log.Printf("Error saving weather data: %v", err)
		}
		c.JSON(http.StatusOK, weatherData)
	})

	r.GET("/api/weather/coordinates/:lat/:lon", func(c *gin.Context) {
		lat := c.Param("lat")
		lon := c.Param("lon")
		weatherData, err := fetchWeatherByCoordinates(lat, lon)
		if err != nil {
			log.Printf("Error fetching weather for coordinates %s,%s: %v", lat, lon, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		err = saveWeatherData(weatherData)
		if err != nil {
			log.Printf("Error saving weather data: %v", err)
		}
		c.JSON(http.StatusOK, weatherData)
	})

	r.GET("/api/history", func(c *gin.Context) {
		history, err := getWeatherHistory()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, history)
	})

	log.Println("Server starting on :8080")
	r.Run(":8080")
}
