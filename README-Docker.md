# 🐳 Weather Dashboard Docker Deployment Guide

A production-ready Docker setup for the Weather Dashboard application with comprehensive security, monitoring, and scalability features.

## 🚀 Quick Start

### Prerequisites
- Docker and Docker Compose installed
- WeatherAPI key from [weatherapi.com](https://weatherapi.com)

### 1. Environment Setup
```bash
# Clone the repository
git clone <your-repo-url>
cd weather-dashboard

# Create environment file
cp .env.example .env
# Edit .env and add your WeatherAPI key
```

### 2. Development Deployment
```bash
# Deploy in development mode
./deploy.sh dev

# Or manually
docker-compose up -d
```

### 3. Production Deployment
```bash
# Deploy in production mode
./deploy.sh prod

# Or manually
docker-compose -f docker-compose.prod.yml up -d
```

## 📁 Docker Files Overview

### Core Files
- **`Dockerfile`** - Multi-stage production build
- **`docker-compose.yml`** - Development environment
- **`docker-compose.prod.yml`** - Production environment with Nginx
- **`.dockerignore`** - Optimized build context
- **`deploy.sh`** - Automated deployment script

### Production Files
- **`nginx.conf`** - Reverse proxy configuration
- **`README-Docker.md`** - This deployment guide

## 🔧 Configuration

### Environment Variables
```bash
# Required
WEATHERAPI_KEY=your_api_key_here

# Optional (with defaults)
PORT=8080
HOST=0.0.0.0
DB_PATH=/app/data/weather.db
GIN_MODE=release
```

### Volume Mounts
- **`weather_data:/app/data`** - Persistent database storage
- **`./templates:/app/templates:ro`** - Read-only templates
- **`./static:/app/static:ro`** - Read-only static files

## 🛡️ Security Features

### Container Security
- ✅ Non-root user execution
- ✅ Read-only filesystem where possible
- ✅ No new privileges
- ✅ Resource limits
- ✅ Health checks

### Network Security
- ✅ Isolated Docker networks
- ✅ Rate limiting (API: 10r/s, General: 30r/s)
- ✅ Security headers
- ✅ HTTPS enforcement (production)

### Application Security
- ✅ Input validation
- ✅ SQL injection protection
- ✅ XSS protection headers
- ✅ CSRF protection

## 📊 Monitoring & Health Checks

### Health Check Endpoints
- **Application**: `http://localhost:8080/`
- **Nginx**: `http://localhost/health`

### Logging
```bash
# View application logs
docker-compose logs -f weather-dashboard

# View nginx logs (production)
docker-compose -f docker-compose.prod.yml logs -f nginx

# View all logs
docker-compose logs -f
```

### Metrics
- Container resource usage
- Application response times
- Error rates
- Database performance

## 🚀 Production Deployment

### 1. SSL Certificate Setup
```bash
# Create SSL directory
mkdir -p ssl

# Add your certificates
cp your-cert.pem ssl/cert.pem
cp your-key.pem ssl/key.pem
```

### 2. Domain Configuration
Edit `nginx.conf` and `docker-compose.prod.yml`:
```nginx
server_name your-domain.com;
```

### 3. Deploy with Nginx
```bash
# Deploy with nginx profile
docker-compose -f docker-compose.prod.yml --profile nginx up -d
```

### 4. Verify Deployment
```bash
# Check container status
docker-compose -f docker-compose.prod.yml ps

# Test health endpoint
curl https://your-domain.com/health

# Check SSL
curl -I https://your-domain.com/
```

## 🔄 Scaling & Updates

### Rolling Updates
```bash
# Update application
git pull
docker-compose -f docker-compose.prod.yml build
docker-compose -f docker-compose.prod.yml up -d --no-deps weather-dashboard
```

### Database Backups
```bash
# Backup database
docker exec weather-dashboard-prod sqlite3 /app/data/weather.db ".backup /app/data/backup_$(date +%Y%m%d_%H%M%S).db"

# Restore database
docker exec -i weather-dashboard-prod sqlite3 /app/data/weather.db < backup_file.db
```

## 🐛 Troubleshooting

### Common Issues

#### 1. Container Won't Start
```bash
# Check logs
docker-compose logs weather-dashboard

# Check environment variables
docker-compose config

# Verify API key
echo $WEATHERAPI_KEY
```

#### 2. Database Issues
```bash
# Check database file
docker exec weather-dashboard ls -la /app/data/

# Reset database (WARNING: data loss)
docker exec weather-dashboard rm /app/data/weather.db
```

#### 3. Network Issues
```bash
# Check network connectivity
docker network ls
docker network inspect weather-dashboard_weather-network

# Test internal communication
docker exec weather-dashboard wget -qO- http://localhost:8080/
```

#### 4. Performance Issues
```bash
# Check resource usage
docker stats weather-dashboard

# Check nginx logs (production)
docker-compose -f docker-compose.prod.yml logs nginx
```

### Debug Mode
```bash
# Run in debug mode
GIN_MODE=debug docker-compose up

# Access debug endpoints
curl http://localhost:8080/debug/vars
```

## 📈 Performance Optimization

### Resource Limits
- **CPU**: 0.5 cores (limit), 0.25 cores (reservation)
- **Memory**: 512MB (limit), 256MB (reservation)

### Caching
- Static files cached for 1 year
- Gzip compression enabled
- Browser caching headers

### Database Optimization
- SQLite with WAL mode
- Connection pooling
- Prepared statements

## 🔍 Monitoring Commands

```bash
# Container status
docker-compose ps

# Resource usage
docker stats

# Network connectivity
docker network inspect weather-dashboard_weather-network

# Volume usage
docker volume ls
docker volume inspect weather-dashboard_weather_data

# Health check
curl -f http://localhost:8080/ || echo "Health check failed"
```

## 🎯 Best Practices

### Security
- ✅ Always use HTTPS in production
- ✅ Regularly update base images
- ✅ Monitor security advisories
- ✅ Use secrets management for API keys

### Performance
- ✅ Monitor resource usage
- ✅ Set appropriate limits
- ✅ Use health checks
- ✅ Implement logging

### Maintenance
- ✅ Regular backups
- ✅ Update dependencies
- ✅ Monitor logs
- ✅ Test deployments

## 📞 Support

For issues and questions:
1. Check the troubleshooting section
2. Review application logs
3. Verify environment configuration
4. Test with minimal setup

---

**Happy Weather Monitoring! 🌤️** 