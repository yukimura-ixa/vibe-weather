// Utility functions
const WeatherUtils = {
    // Format location details consistently
    formatLocationDetails(country, state) {
        if (country && state) return `${state}, ${country}`;
        if (country) return country;
        if (state) return state;
        return '';
    },

    // Show loading state
    showLoading(weatherCard, message = 'Loading weather data...') {
        weatherCard.style.display = 'block';
        weatherCard.innerHTML = `<div class="loading">${message}</div>`;
    },

    // Hide weather card
    hideWeatherCard(weatherCard) {
        weatherCard.style.display = 'none';
    },

    // Validate city input
    validateCityInput(city) {
        return city && city.trim().length > 0;
    },

    // Create weather display HTML
    createWeatherHTML(data) {
        const locationDetails = this.formatLocationDetails(data.country, data.state);
        
        return `
            <h2>Current Weather</h2>
            <div class="weather-info">
                <div class="city-name">${data.city}</div>
                ${locationDetails ? `<div class="location-details">${locationDetails}</div>` : ''}
                ${data.icon ? `<img src="${data.icon}" alt="weather icon" class="weather-icon">` : ''}
                <div class="temperature">${Math.round(data.temperature)}°C</div>
                <div class="description">${data.description}</div>
                <div class="humidity">Humidity: ${data.humidity}%</div>
                <div class="timestamp">Updated: ${new Date(data.timestamp).toLocaleString()}</div>
            </div>
        `;
    },

    // Create history item HTML
    createHistoryItemHTML(item) {
        const locationDetails = this.formatLocationDetails(item.country, item.state);
        
        return `
            <div class="history-item">
                <h3>${item.city}</h3>
                ${locationDetails ? `<p>${locationDetails}</p>` : ''}
                ${item.icon ? `<img src="${item.icon}" alt="weather icon" class="weather-icon">` : ''}
                <p>Temperature: ${Math.round(item.temperature)}°C</p>
                <p>${item.description}</p>
                <p>Humidity: ${item.humidity}%</p>
                <p>${new Date(item.timestamp).toLocaleString()}</p>
            </div>
        `;
    }
};

// Error handling utility
const ErrorHandler = {
    showError(message) {
        // Remove any existing error messages
        const existingError = document.querySelector('.error');
        if (existingError) {
            existingError.remove();
        }
        
        // Create and display error message
        const errorDiv = document.createElement('div');
        errorDiv.className = 'error';
        errorDiv.textContent = message;
        
        const searchSection = document.querySelector('.search-section');
        searchSection.appendChild(errorDiv);
        
        // Remove error after 5 seconds
        setTimeout(() => {
            if (errorDiv.parentNode) {
                errorDiv.remove();
            }
        }, 5000);
    },

    handleGeolocationError(error) {
        let message = 'Error getting your location';
        
        if (error.code === 1) {
            message = 'Location access denied. Please allow location access or search for a city manually.';
        } else if (error.code === 2) {
            message = 'Location unavailable. Please search for a city manually.';
        } else if (error.code === 3) {
            message = 'Location request timed out. Please search for a city manually.';
        } else {
            message = `Error getting your location: ${error.message}`;
        }
        
        this.showError(message);
    }
};

// API service
const WeatherAPI = {
    async fetchWeather(city) {
        const response = await fetch(`/api/weather/${encodeURIComponent(city)}`);
        
        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || 'City not found or API error');
        }

        return await response.json();
    },

    async fetchWeatherByCoordinates(latitude, longitude) {
        const response = await fetch(`/api/weather/coordinates/${latitude}/${longitude}`);
        
        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || 'Error fetching weather for your location');
        }

        return await response.json();
    },

    async fetchHistory() {
        const response = await fetch('/api/history');
        return await response.json();
    }
};

// Main application logic
const WeatherApp = {
    weatherCard: null,

    init() {
        this.weatherCard = document.getElementById('currentWeather');
        this.setupEventListeners();
        this.loadHistory();
    },

    setupEventListeners() {
        // Handle Enter key press in search input
        document.getElementById('cityInput').addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                this.getWeather();
            }
        });
    },

    async getWeather() {
        const cityInput = document.getElementById('cityInput');
        const city = cityInput.value.trim();
        
        if (!WeatherUtils.validateCityInput(city)) {
            ErrorHandler.showError('Please enter a city name');
            return;
        }

        WeatherUtils.showLoading(this.weatherCard);

        try {
            const weatherData = await WeatherAPI.fetchWeather(city);
            this.displayWeather(weatherData);
            this.loadHistory();
        } catch (error) {
            ErrorHandler.showError(`Error fetching weather data: ${error.message}`);
            WeatherUtils.hideWeatherCard(this.weatherCard);
        }
    },

    displayWeather(data) {
        this.weatherCard.innerHTML = WeatherUtils.createWeatherHTML(data);
        this.weatherCard.style.display = 'block';
    },

    async getMyLocation() {
        if (!navigator.geolocation) {
            ErrorHandler.showError('Geolocation is not supported by this browser');
            return;
        }

        WeatherUtils.showLoading(this.weatherCard, 'Getting your location...');

        try {
            const position = await new Promise((resolve, reject) => {
                navigator.geolocation.getCurrentPosition(resolve, reject, {
                    enableHighAccuracy: true,
                    timeout: 10000,
                    maximumAge: 60000
                });
            });

            const { latitude, longitude } = position.coords;
            const weatherData = await WeatherAPI.fetchWeatherByCoordinates(latitude, longitude);
            
            this.displayWeather(weatherData);
            this.loadHistory();
        } catch (error) {
            ErrorHandler.handleGeolocationError(error);
            WeatherUtils.hideWeatherCard(this.weatherCard);
        }
    },

    async loadHistory() {
        try {
            const history = await WeatherAPI.fetchHistory();

            // Defensive: If not an array, show no data
            if (!Array.isArray(history) || history.length === 0) {
                document.getElementById('historyList').innerHTML = '<p class="no-data">No recent searches</p>';
                return;
            }

            const historyHTML = history.map(item => WeatherUtils.createHistoryItemHTML(item)).join('');
            document.getElementById('historyList').innerHTML = historyHTML;
        } catch (error) {
            console.error('Error loading history:', error);
            document.getElementById('historyList').innerHTML = '<p class="no-data">Error loading history</p>';
        }
    }
};

// Initialize app when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    WeatherApp.init();
});

// Make getMyLocation globally accessible for HTML onclick
window.getMyLocation = () => WeatherApp.getMyLocation(); 