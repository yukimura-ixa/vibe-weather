package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidCityName(t *testing.T) {
	tests := []struct {
		name     string
		city     string
		expected bool
	}{
		{
			name:     "valid city name with letters only",
			city:     "London",
			expected: true,
		},
		{
			name:     "valid city name with spaces",
			city:     "New York",
			expected: true,
		},
		{
			name:     "valid city name with hyphen",
			city:     "Saint-Jean",
			expected: true,
		},
		{
			name:     "valid city name with apostrophe",
			city:     "King's Lynn",
			expected: true,
		},
		{
			name:     "empty string",
			city:     "",
			expected: false,
		},
		{
			name:     "whitespace only",
			city:     "   ",
			expected: false,
		},
		{
			name:     "single character",
			city:     "A",
			expected: false,
		},
		{
			name:     "contains numbers",
			city:     "London123",
			expected: false,
		},
		{
			name:     "contains special characters",
			city:     "London@",
			expected: false,
		},
		{
			name:     "contains dots",
			city:     "St. Petersburg",
			expected: false,
		},
		{
			name:     "contains commas",
			city:     "London, UK",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidCityName(tt.city)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSanitizeCityName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normal case",
			input:    "london",
			expected: "London",
		},
		{
			name:     "mixed case",
			input:    "LoNdOn",
			expected: "London",
		},
		{
			name:     "with extra spaces",
			input:    "  new  york  ",
			expected: "New York",
		},
		{
			name:     "multiple spaces",
			input:    "new    york",
			expected: "New York",
		},
		{
			name:     "with tabs",
			input:    "new\tyork",
			expected: "New York",
		},
		{
			name:     "with newlines",
			input:    "new\nyork",
			expected: "New York",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "whitespace only",
			input:    "   ",
			expected: "",
		},
		{
			name:     "single word",
			input:    "paris",
			expected: "Paris",
		},
		{
			name:     "already title case",
			input:    "London",
			expected: "London",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeCityName(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsValidCoordinate(t *testing.T) {
	tests := []struct {
		name     string
		lat      string
		lon      string
		expected bool
	}{
		{
			name:     "valid positive coordinates",
			lat:      "51.5074",
			lon:      "0.1278",
			expected: true,
		},
		{
			name:     "valid negative coordinates",
			lat:      "-33.8688",
			lon:      "151.2093",
			expected: true,
		},
		{
			name:     "valid integer coordinates",
			lat:      "51",
			lon:      "0",
			expected: true,
		},
		{
			name:     "valid coordinates with many decimal places",
			lat:      "51.5074000",
			lon:      "0.1278000",
			expected: true,
		},
		{
			name:     "empty latitude",
			lat:      "",
			lon:      "0.1278",
			expected: false,
		},
		{
			name:     "empty longitude",
			lat:      "51.5074",
			lon:      "",
			expected: false,
		},
		{
			name:     "both empty",
			lat:      "",
			lon:      "",
			expected: false,
		},
		{
			name:     "invalid latitude with letters",
			lat:      "51.5074a",
			lon:      "0.1278",
			expected: false,
		},
		{
			name:     "invalid longitude with letters",
			lat:      "51.5074",
			lon:      "0.1278b",
			expected: false,
		},
		{
			name:     "latitude with special characters",
			lat:      "51.5074@",
			lon:      "0.1278",
			expected: false,
		},
		{
			name:     "longitude with special characters",
			lat:      "51.5074",
			lon:      "0.1278#",
			expected: false,
		},
		{
			name:     "latitude out of range (too high)",
			lat:      "91.0",
			lon:      "0.1278",
			expected: true, // Current validation only checks format, not range
		},
		{
			name:     "longitude out of range (too high)",
			lat:      "51.5074",
			lon:      "181.0",
			expected: true, // Current validation only checks format, not range
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidCoordinate(tt.lat, tt.lon)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEdgeCases(t *testing.T) {
	t.Run("city name with only spaces and special characters", func(t *testing.T) {
		assert.False(t, IsValidCityName("   @#$%   "))
	})

	t.Run("coordinate with leading zeros", func(t *testing.T) {
		assert.True(t, IsValidCoordinate("051.5074", "00.1278"))
	})

	t.Run("coordinate with scientific notation", func(t *testing.T) {
		// Current validation doesn't support scientific notation
		assert.False(t, IsValidCoordinate("5.15074e1", "1.278e-1"))
	})

	t.Run("sanitize with special characters", func(t *testing.T) {
		// Sanitize should handle special characters gracefully
		result := SanitizeCityName("  london@#$%  ")
		assert.Equal(t, "London@#$%", result)
	})
}
