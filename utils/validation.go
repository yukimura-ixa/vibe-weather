package utils

import (
	"regexp"
	"strings"
)

// IsValidCityName checks if a city name is valid
func IsValidCityName(city string) bool {
	if strings.TrimSpace(city) == "" {
		return false
	}

	// Check for minimum length
	if len(strings.TrimSpace(city)) < 2 {
		return false
	}

	// Check for valid characters (letters, spaces, hyphens, apostrophes)
	validPattern := regexp.MustCompile(`^[a-zA-Z\s\-']+$`)
	return validPattern.MatchString(city)
}

// SanitizeCityName cleans and normalizes city name input
func SanitizeCityName(city string) string {
	// Trim whitespace
	cleaned := strings.TrimSpace(city)

	// Convert to title case
	cleaned = strings.Title(strings.ToLower(cleaned))

	// Remove extra spaces
	spacePattern := regexp.MustCompile(`\s+`)
	cleaned = spacePattern.ReplaceAllString(cleaned, " ")

	return cleaned
}

// IsValidCoordinate checks if latitude and longitude are valid
func IsValidCoordinate(lat, lon string) bool {
	// Basic validation - could be enhanced with proper coordinate validation
	if lat == "" || lon == "" {
		return false
	}

	// Check if they can be parsed as floats
	latPattern := regexp.MustCompile(`^-?\d+(\.\d+)?$`)
	lonPattern := regexp.MustCompile(`^-?\d+(\.\d+)?$`)

	return latPattern.MatchString(lat) && lonPattern.MatchString(lon)
}
