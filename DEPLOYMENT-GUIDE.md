# 🚀 Deployment Guide - Free Hosting Platforms

This guide will help you deploy your Weather Dashboard to various free hosting platforms.

## 📋 Prerequisites

1. **GitHub Repository**: Your code must be pushed to GitHub
2. **WeatherAPI Key**: Get a free API key from [weatherapi.com](https://weatherapi.com)
3. **Account**: Sign up for your chosen hosting platform

## 🎯 Platform Options

### 1. Railway (Recommended) ⭐

**Free Tier**: $5/month credit (sufficient for small apps)

#### Quick Deploy Steps:

1. **Visit [Railway.app](https://railway.app)**
2. **Sign up** with GitHub
3. **Click "New Project"** → **"Deploy from GitHub repo"**
4. **Select your repository**: `weather-dashboard`
5. **Add Environment Variables**:
   ```
   WEATHERAPI_KEY=your_api_key_here
   PORT=8080
   HOST=0.0.0.0
   DB_PATH=/app/data/weather.db
   GIN_MODE=release
   ```
6. **Deploy** - Railway will automatically build and deploy your app

#### Railway Features:
- ✅ **Automatic HTTPS**
- ✅ **Custom domains**
- ✅ **Auto-deploy from GitHub**
- ✅ **Environment variables**
- ✅ **Logs and monitoring**

---

### 2. Render

**Free Tier**: 750 hours/month

#### Quick Deploy Steps:

1. **Visit [Render.com](https://render.com)**
2. **Sign up** with GitHub
3. **Click "New"** → **"Web Service"**
4. **Connect your repository**: `weather-dashboard`
5. **Configure service**:
   - **Name**: `weather-dashboard`
   - **Environment**: `Docker`
   - **Branch**: `master`
   - **Root Directory**: `/` (leave empty)
6. **Add Environment Variables**:
   ```
   WEATHERAPI_KEY=your_api_key_here
   ```
7. **Deploy** - Render will use the `render.yaml` configuration

#### Render Features:
- ✅ **Automatic HTTPS**
- ✅ **Custom domains**
- ✅ **Auto-deploy from GitHub**
- ✅ **Free SSL certificates**

---

### 3. Fly.io

**Free Tier**: 3 shared-cpu VMs, 3GB persistent volume

#### Quick Deploy Steps:

1. **Install Fly CLI**:
   ```bash
   # Windows (PowerShell)
   iwr https://fly.io/install.ps1 -useb | iex
   
   # macOS
   curl -L https://fly.io/install.sh | sh
   ```

2. **Sign up** at [Fly.io](https://fly.io)

3. **Login and deploy**:
   ```bash
   fly auth login
   fly launch
   ```

4. **Set environment variables**:
   ```bash
   fly secrets set WEATHERAPI_KEY=your_api_key_here
   ```

5. **Deploy**:
   ```bash
   fly deploy
   ```

#### Fly.io Features:
- ✅ **Global deployment** (multiple regions)
- ✅ **Automatic HTTPS**
- ✅ **Custom domains**
- ✅ **Persistent volumes**

---

### 4. Heroku

**Free Tier**: Discontinued, but has affordable plans

#### Quick Deploy Steps:

1. **Install Heroku CLI**:
   ```bash
   # Windows
   winget install --id=Heroku.HerokuCLI
   ```

2. **Sign up** at [Heroku.com](https://heroku.com)

3. **Login and create app**:
   ```bash
   heroku login
   heroku create your-weather-dashboard
   ```

4. **Set environment variables**:
   ```bash
   heroku config:set WEATHERAPI_KEY=your_api_key_here
   heroku config:set PORT=8080
   heroku config:set HOST=0.0.0.0
   ```

5. **Deploy**:
   ```bash
   git push heroku master
   ```

---

## 🔧 Environment Variables

All platforms require these environment variables:

| Variable | Value | Description |
|----------|-------|-------------|
| `WEATHERAPI_KEY` | `your_api_key_here` | **Required** - Your WeatherAPI.com key |
| `PORT` | `8080` | Port for the application |
| `HOST` | `0.0.0.0` | Host binding for cloud deployment |
| `DB_PATH` | `/app/data/weather.db` | SQLite database path |
| `GIN_MODE` | `release` | Production mode for Gin |

## 🌐 Custom Domain Setup

### Railway
1. Go to your project settings
2. Click "Domains"
3. Add your custom domain
4. Update DNS records as instructed

### Render
1. Go to your service settings
2. Click "Custom Domains"
3. Add your domain
4. Update DNS records

### Fly.io
```bash
fly certs add yourdomain.com
```

## 📊 Monitoring and Logs

### Railway
- **Logs**: Available in the Railway dashboard
- **Metrics**: CPU, memory, and network usage
- **Alerts**: Automatic notifications for issues

### Render
- **Logs**: Available in the service dashboard
- **Metrics**: Request count and response times
- **Health checks**: Automatic monitoring

### Fly.io
```bash
# View logs
fly logs

# Monitor metrics
fly status
```

## 🔒 Security Best Practices

1. **Environment Variables**: Never commit API keys to Git
2. **HTTPS**: All platforms provide automatic SSL
3. **Input Validation**: Already implemented in the app
4. **Rate Limiting**: Consider adding for production use

## 🚨 Troubleshooting

### Common Issues:

1. **Build Failures**:
   - Check Dockerfile syntax
   - Verify all files are committed to GitHub
   - Check platform-specific requirements

2. **Runtime Errors**:
   - Verify environment variables are set
   - Check application logs
   - Ensure API key is valid

3. **Database Issues**:
   - Verify database path is writable
   - Check disk space on platform
   - Ensure proper permissions

### Debug Commands:

```bash
# Check application logs
docker logs <container_id>

# Test API endpoint
curl https://your-app-url.com/api/weather/london

# Check environment variables
echo $WEATHERAPI_KEY
```

## 📈 Performance Optimization

1. **Database**: SQLite is sufficient for small apps
2. **Caching**: Consider adding Redis for high traffic
3. **CDN**: Use platform CDN for static assets
4. **Monitoring**: Set up alerts for performance issues

## 🎯 Recommended Deployment Flow

1. **Choose Railway** (easiest for beginners)
2. **Push code to GitHub**
3. **Connect repository to Railway**
4. **Set environment variables**
5. **Deploy and test**
6. **Set up custom domain** (optional)
7. **Monitor and maintain**

## 📞 Support

- **Railway**: [Discord](https://discord.gg/railway)
- **Render**: [Documentation](https://render.com/docs)
- **Fly.io**: [Community](https://community.fly.io)
- **Heroku**: [Support](https://help.heroku.com)

---

**Happy Deploying! 🚀**

Your Weather Dashboard will be live and accessible from anywhere in the world! 