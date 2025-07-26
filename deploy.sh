#!/bin/bash

# Weather Dashboard Deployment Script
# Usage: ./deploy.sh [dev|prod]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if environment is provided
ENVIRONMENT=${1:-dev}
print_status "Deploying Weather Dashboard in $ENVIRONMENT mode"

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    print_error "Docker is not running. Please start Docker and try again."
    exit 1
fi

# Check if .env file exists
if [ ! -f .env ]; then
    print_warning ".env file not found. Creating template..."
    cat > .env << EOF
# Weather Dashboard Environment Variables
WEATHERAPI_KEY=your_weather_api_key_here
PORT=8080
HOST=0.0.0.0
DB_PATH=/app/data/weather.db
GIN_MODE=release
EOF
    print_error "Please update .env file with your WeatherAPI key and run the script again."
    exit 1
fi

# Load environment variables
source .env

# Check if API key is set
if [ "$WEATHERAPI_KEY" = "your_weather_api_key_here" ] || [ -z "$WEATHERAPI_KEY" ]; then
    print_error "Please set your WeatherAPI key in the .env file"
    exit 1
fi

print_status "Environment variables loaded successfully"

# Stop existing containers
print_status "Stopping existing containers..."
docker-compose down --remove-orphans 2>/dev/null || true

# Build and start containers
if [ "$ENVIRONMENT" = "prod" ]; then
    print_status "Building production image..."
    docker-compose -f docker-compose.prod.yml build --no-cache
    
    print_status "Starting production services..."
    docker-compose -f docker-compose.prod.yml up -d
    
    # Wait for services to be ready
    print_status "Waiting for services to be ready..."
    sleep 10
    
    # Check if nginx profile is enabled
    if docker-compose -f docker-compose.prod.yml ps | grep -q nginx; then
        print_status "Nginx reverse proxy is running"
    fi
else
    print_status "Building development image..."
    docker-compose build --no-cache
    
    print_status "Starting development services..."
    docker-compose up -d
fi

# Wait for application to be ready
print_status "Waiting for application to be ready..."
for i in {1..30}; do
    if curl -f http://localhost:8080/ > /dev/null 2>&1; then
        print_success "Application is ready!"
        break
    fi
    if [ $i -eq 30 ]; then
        print_error "Application failed to start within 30 seconds"
        docker-compose logs weather-dashboard
        exit 1
    fi
    sleep 1
done

# Health check
print_status "Performing health check..."
if curl -f http://localhost:8080/ > /dev/null 2>&1; then
    print_success "Health check passed!"
else
    print_error "Health check failed!"
    exit 1
fi

# Show container status
print_status "Container status:"
docker-compose ps

# Show logs
print_status "Recent logs:"
docker-compose logs --tail=20 weather-dashboard

print_success "Weather Dashboard deployed successfully!"
print_status "Access the application at: http://localhost:8080"
print_status "API endpoints available at: http://localhost:8080/api/"

if [ "$ENVIRONMENT" = "prod" ]; then
    print_status "Production mode enabled with Nginx reverse proxy"
    print_status "Make sure to configure SSL certificates for production use"
fi

print_status "To view logs: docker-compose logs -f weather-dashboard"
print_status "To stop: docker-compose down" 