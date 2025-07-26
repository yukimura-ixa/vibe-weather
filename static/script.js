// Load weather history on page load
document.addEventListener('DOMContentLoaded', function() {
    loadHistory();
});

// Handle Enter key press in search input
document.getElementById('cityInput').addEventListener('keypress', function(e) {
    if (e.key === 'Enter') {
        getWeather();
    }
});

async function getWeather() {
    const cityInput = document.getElementById('cityInput');
    const city = cityInput.value.trim();
    
    if (!city || city.length === 0) {
        showError('Please enter a city name');
        return;
    }

    // Show loading state
    const weatherCard = document.getElementById('currentWeather');
    weatherCard.style.display = 'block';
    weatherCard.innerHTML = '<div class="loading">Loading weather data...</div>';

    try {
        const response = await fetch(`/api/weather/${encodeURIComponent(city)}`);
        
        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || 'City not found or API error');
        }

        const weatherData = await response.json();
        displayWeather(weatherData);
        loadHistory(); // Refresh history after new search
        
    } catch (error) {
        showError('Error fetching weather data: ' + error.message);
        weatherCard.style.display = 'none';
    }
}

function displayWeather(data) {
    const weatherCard = document.getElementById('currentWeather');
    
    // Build location details string
    let locationDetails = '';
    if (data.country && data.state) {
        locationDetails = `${data.state}, ${data.country}`;
    } else if (data.country) {
        locationDetails = data.country;
    } else if (data.state) {
        locationDetails = data.state;
    }
    
    weatherCard.innerHTML = `
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
    
    weatherCard.style.display = 'block';
}

async function getMyLocation() {
    if (!navigator.geolocation) {
        showError('Geolocation is not supported by this browser');
        return;
    }

    // Show loading state
    const weatherCard = document.getElementById('currentWeather');
    weatherCard.style.display = 'block';
    weatherCard.innerHTML = '<div class="loading">Getting your location...</div>';

    try {
        const position = await new Promise((resolve, reject) => {
            navigator.geolocation.getCurrentPosition(resolve, reject, {
                enableHighAccuracy: true,
                timeout: 10000,
                maximumAge: 60000
            });
        });

        const { latitude, longitude } = position.coords;
        
        // Fetch weather for coordinates
        const response = await fetch(`/api/weather/coordinates/${latitude}/${longitude}`);
        
        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || 'Error fetching weather for your location');
        }

        const weatherData = await response.json();
        displayWeather(weatherData);
        loadHistory(); // Refresh history after new search
        
    } catch (error) {
        if (error.code === 1) {
            showError('Location access denied. Please allow location access or search for a city manually.');
        } else if (error.code === 2) {
            showError('Location unavailable. Please search for a city manually.');
        } else if (error.code === 3) {
            showError('Location request timed out. Please search for a city manually.');
        } else {
            showError('Error getting your location: ' + error.message);
        }
        weatherCard.style.display = 'none';
    }
}

async function loadHistory() {
    try {
        const response = await fetch('/api/history');
        const history = await response.json();

        // Defensive: If not an array, show no data
        if (!Array.isArray(history) || history.length === 0) {
            document.getElementById('historyList').innerHTML = '<p class="no-data">No recent searches</p>';
            return;
        }

        document.getElementById('historyList').innerHTML = history.map(item => {
            let locationDetails = '';
            if (item.country && item.state) {
                locationDetails = `${item.state}, ${item.country}`;
            } else if (item.country) {
                locationDetails = item.country;
            } else if (item.state) {
                locationDetails = item.state;
            }
            
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
        }).join('');
        
    } catch (error) {
        console.error('Error loading history:', error);
        document.getElementById('historyList').innerHTML = '<p class="no-data">Error loading history</p>';
    }
}

function showError(message) {
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
} 

window.getMyLocation = getMyLocation; 