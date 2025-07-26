@echo off
REM Weather Dashboard Deployment Script for Windows
REM Usage: deploy.bat [dev|prod]

setlocal enabledelayedexpansion

REM Set colors for output
set "BLUE=[94m"
set "GREEN=[92m"
set "YELLOW=[93m"
set "RED=[91m"
set "NC=[0m"

REM Function to print colored output
:print_status
echo %BLUE%[INFO]%NC% %~1
goto :eof

:print_success
echo %GREEN%[SUCCESS]%NC% %~1
goto :eof

:print_warning
echo %YELLOW%[WARNING]%NC% %~1
goto :eof

:print_error
echo %RED%[ERROR]%NC% %~1
goto :eof

REM Check if environment is provided
set "ENVIRONMENT=%1"
if "%ENVIRONMENT%"=="" set "ENVIRONMENT=dev"
call :print_status "Deploying Weather Dashboard in %ENVIRONMENT% mode"

REM Check if Docker is running
docker info >nul 2>&1
if errorlevel 1 (
    call :print_error "Docker is not running. Please start Docker and try again."
    exit /b 1
)

REM Check if .env file exists
if not exist .env (
    call :print_warning ".env file not found. Creating template..."
    (
        echo # Weather Dashboard Environment Variables
        echo WEATHERAPI_KEY=your_weather_api_key_here
        echo PORT=8080
        echo HOST=0.0.0.0
        echo DB_PATH=/app/data/weather.db
        echo GIN_MODE=release
    ) > .env
    call :print_error "Please update .env file with your WeatherAPI key and run the script again."
    exit /b 1
)

REM Load environment variables (simplified for Windows)
for /f "tokens=1,2 delims==" %%a in (.env) do (
    if not "%%a"=="" if not "%%a:~0,1%"=="#" (
        set "%%a=%%b"
    )
)

REM Check if API key is set
if "%WEATHERAPI_KEY%"=="your_weather_api_key_here" (
    call :print_error "Please set your WeatherAPI key in the .env file"
    exit /b 1
)

call :print_status "Environment variables loaded successfully"

REM Stop existing containers
call :print_status "Stopping existing containers..."
docker-compose down --remove-orphans 2>nul

REM Build and start containers
if "%ENVIRONMENT%"=="prod" (
    call :print_status "Building production image..."
    docker-compose -f docker-compose.prod.yml build --no-cache
    
    call :print_status "Starting production services..."
    docker-compose -f docker-compose.prod.yml up -d
    
    REM Wait for services to be ready
    call :print_status "Waiting for services to be ready..."
    timeout /t 10 /nobreak >nul
) else (
    call :print_status "Building development image..."
    docker-compose build --no-cache
    
    call :print_status "Starting development services..."
    docker-compose up -d
)

REM Wait for application to be ready
call :print_status "Waiting for application to be ready..."
for /l %%i in (1,1,30) do (
    curl -f http://localhost:8080/ >nul 2>&1
    if not errorlevel 1 (
        call :print_success "Application is ready!"
        goto :health_check
    )
    timeout /t 1 /nobreak >nul
)
call :print_error "Application failed to start within 30 seconds"
docker-compose logs weather-dashboard
exit /b 1

:health_check
REM Health check
call :print_status "Performing health check..."
curl -f http://localhost:8080/ >nul 2>&1
if errorlevel 1 (
    call :print_error "Health check failed!"
    exit /b 1
)
call :print_success "Health check passed!"

REM Show container status
call :print_status "Container status:"
docker-compose ps

REM Show logs
call :print_status "Recent logs:"
docker-compose logs --tail=20 weather-dashboard

call :print_success "Weather Dashboard deployed successfully!"
call :print_status "Access the application at: http://localhost:8080"
call :print_status "API endpoints available at: http://localhost:8080/api/"

if "%ENVIRONMENT%"=="prod" (
    call :print_status "Production mode enabled with Nginx reverse proxy"
    call :print_status "Make sure to configure SSL certificates for production use"
)

call :print_status "To view logs: docker-compose logs -f weather-dashboard"
call :print_status "To stop: docker-compose down"

endlocal 