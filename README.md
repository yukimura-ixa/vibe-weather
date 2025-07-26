# Weather Dashboard

A minimal weather dashboard built with Go, featuring a beautiful glassmorphism web interface, SQLite database, and external weather API integration.

## Features

- 🌤️ Real-time weather data from OpenWeatherMap API
- 🎨 Beautiful glassmorphism UI design
- 💾 SQLite database for storing search history
- 🐳 Docker containerization
- 📱 Responsive design for mobile devices
- 🔍 Search history with recent queries

## Tech Stack

- **Backend**: Go with Gin framework
- **Database**: SQLite (minimal and file-based)
- **Frontend**: HTML, CSS, JavaScript
- **Design**: Glassmorphism with CSS backdrop-filter
- **Containerization**: Docker & Docker Compose
- **External API**: OpenWeatherMap

## Quick Start

### Using Docker Compose (Recommended)

1. **Clone and navigate to the project:**
   ```bash
   cd weather-dashboard
   ```

2. **Build and run with Docker Compose:**
   ```bash
   docker-compose up --build
   ```

3. **Open your browser and visit:**
   ```
   http://localhost:8080
   ```

### Manual Setup

1. **Install Go dependencies:**
   ```bash
   go mod download
   ```

2. **Run the application:**
   ```bash
   go run main.go
   ```

3. **Open your browser and visit:**
   ```
   http://localhost:8080
   ```

## API Key Setup

The application uses OpenWeatherMap API for weather data. You can get a free API key from [OpenWeatherMap](https://openweathermap.org/api):

1. Sign up at OpenWeatherMap
2. Get your API key
3. Create a `.env` file in the project root:
   ```
   OPENWEATHER_API_KEY=your_api_key_here
   ```

**Note**: The `.env` file is already created with your API key and is included in `.gitignore` for security.

## API Endpoints

- `GET /` - Main dashboard page
- `GET /api/weather/:city` - Get weather data for a specific city
- `GET /api/history` - Get recent search history

## Project Structure

```
weather-dashboard/
├── main.go              # Go application entry point
├── go.mod               # Go module file
├── Dockerfile           # Docker configuration
├── docker-compose.yml   # Docker Compose configuration
├── README.md           # This file
├── templates/
│   └── index.html      # Main HTML template
└── static/
    ├── style.css       # Glassmorphism CSS styles
    └── script.js       # Frontend JavaScript
```

## Features in Detail

### Glassmorphism Design
- Semi-transparent backgrounds with backdrop blur
- Subtle borders and shadows
- Smooth hover animations
- Gradient backgrounds

### Weather Data
- Current temperature in Celsius
- Weather description
- Humidity percentage
- Timestamp of last update

### Search History
- Stores last 10 searches in SQLite database
- Displays city name, temperature, description, and timestamp
- Persistent across application restarts

## Development

### Prerequisites
- Go 1.21 or later
- Docker (optional)

### Local Development
```bash
# Install dependencies
go mod download

# Run the application
go run main.go

# The app will be available at http://localhost:8080
```

### Docker Development
```bash
# Build the image
docker build -t weather-dashboard .

# Run the container
docker run -p 8080:8080 weather-dashboard
```

## License

This project is open source and available under the MIT License. 