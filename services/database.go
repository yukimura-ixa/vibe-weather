package services

import (
	"database/sql"
	"fmt"

	"weather-dashboard/handlers"
	"weather-dashboard/models"

	_ "modernc.org/sqlite"
)

// DatabaseService handles all database operations
type DatabaseService struct {
	db *sql.DB
}

// NewDatabaseService creates a new database service
func NewDatabaseService(dbPath string) (*DatabaseService, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	service := &DatabaseService{db: db}
	if err := service.initDB(); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return service, nil
}

// Close closes the database connection
func (s *DatabaseService) Close() error {
	return s.db.Close()
}

// initDB initializes the database schema
func (s *DatabaseService) initDB() error {
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

	_, err := s.db.Exec(createTable)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
}

// SaveWeatherData saves weather data to the database
func (s *DatabaseService) SaveWeatherData(data *models.WeatherData) error {
	query := `
		INSERT INTO weather_data 
		(city, country, state, temperature, description, humidity, icon, condition_code, timestamp) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := s.db.Exec(query,
		data.City, data.Country, data.State, data.Temperature,
		data.Description, data.Humidity, data.Icon, data.ConditionCode, data.Timestamp)

	if err != nil {
		return fmt.Errorf("failed to save weather data: %w", err)
	}

	return nil
}

// GetWeatherHistory retrieves recent weather history
func (s *DatabaseService) GetWeatherHistory(limit int) ([]models.WeatherData, error) {
	query := `
		SELECT id, city, country, state, temperature, description, humidity, icon, condition_code, timestamp 
		FROM weather_data 
		ORDER BY timestamp DESC 
		LIMIT ?`

	rows, err := s.db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query weather history: %w", err)
	}
	defer rows.Close()

	var history []models.WeatherData
	for rows.Next() {
		var data models.WeatherData
		err := rows.Scan(
			&data.ID, &data.City, &data.Country, &data.State,
			&data.Temperature, &data.Description, &data.Humidity,
			&data.Icon, &data.ConditionCode, &data.Timestamp)
		if err != nil {
			return nil, fmt.Errorf("failed to scan weather data: %w", err)
		}
		history = append(history, data)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return history, nil
}

// GetWeatherHistoryDefault retrieves weather history with default limit
func (s *DatabaseService) GetWeatherHistoryDefault() ([]models.WeatherData, error) {
	return s.GetWeatherHistory(models.HistoryLimit)
}

// Ensure DatabaseService implements handlers.DatabaseServiceInterface
var _ handlers.DatabaseServiceInterface = (*DatabaseService)(nil)
