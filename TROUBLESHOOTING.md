# Troubleshooting Guide

Solutions to common problems when running AI Analytics system.

---

## 🔍 Quick Diagnosis

Run this command first:
```bash
./manage.sh health
```

This checks:
- ✅ MongoDB connection
- ✅ Redis connection
- ✅ Backend API status
- ✅ Frontend availability
- ✅ All services running

---

## 📋 Table of Contents

1. [Installation Issues](#installation-issues)
2. [Docker & Containers](#docker--containers)
3. [Database Problems](#database-problems)
4. [API Errors](#api-errors)
5. [Frontend Issues](#frontend-issues)
6. [ML & Predictions](#ml--predictions)
7. [Performance Problems](#performance-problems)
8. [Deployment Issues](#deployment-issues)

---

## 🚀 Installation Issues

### Problem: Docker Compose fails to start

**Symptoms:**
```
ERROR: Version in "./docker-compose.yml" is unsupported
```

**Solution:**
Update Docker Compose to v2.0+:
```bash
# Linux
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# macOS
brew upgrade docker-compose

# Windows
# Update Docker Desktop
```

---

### Problem: Port already in use

**Symptoms:**
```
Error starting userland proxy: listen tcp 0.0.0.0:8080: bind: address already in use
```

**Solution:**

Check what's using the port:
```bash
# Linux/Mac
lsof -i :8080
sudo kill <PID>

# Windows
netstat -ano | findstr :8080
taskkill /PID <PID> /F
```

Or change the port in `.env`:
```env
API_PORT=8081
```

---

### Problem: Permission denied errors

**Symptoms:**
```
mkdir: cannot create directory '/data/db': Permission denied
```

**Solution:**
```bash
# Fix permissions
sudo chown -R $USER:$USER .

# Or run with sudo (not recommended)
sudo docker-compose up
```

---

## 🐳 Docker & Containers

### Problem: Container exits immediately

**Diagnosis:**
```bash
docker-compose ps
docker-compose logs <service-name>
```

**Common causes:**

1. **Missing environment variables**
   ```bash
   # Create .env from template
   cp .env.example .env
   ```

2. **Configuration error**
   ```bash
   # Check service logs
   docker-compose logs backend
   docker-compose logs ml_training
   ```

3. **Dependency not ready**
   ```bash
   # Restart in order
   docker-compose up -d mongodb
   sleep 5
   docker-compose up -d redis
   docker-compose up -d backend
   ```

---

### Problem: Out of memory errors

**Symptoms:**
```
MongoDB killed by Docker (exit code 137)
```

**Solution:**

Increase Docker memory:
- Docker Desktop → Settings → Resources → Memory → 4GB+

Or edit `docker-compose.yml`:
```yaml
services:
  mongodb:
    mem_limit: 2g
  backend:
    mem_limit: 1g
```

---

### Problem: Slow container startup

**Solution:**

1. **Use Docker image cache:**
   ```bash
   docker-compose build --parallel
   ```

2. **Prune unused resources:**
   ```bash
   docker system prune -a
   ```

3. **Check disk space:**
   ```bash
   df -h
   docker system df
   ```

---

## 💾 Database Problems

### Problem: MongoDB connection refused

**Symptoms:**
```
Error: connect ECONNREFUSED 127.0.0.1:27017
```

**Diagnosis:**
```bash
# Check MongoDB is running
docker-compose ps mongodb

# Test connection
mongosh mongodb://localhost:27017
```

**Solutions:**

1. **Start MongoDB:**
   ```bash
   docker-compose up -d mongodb
   ```

2. **Check MongoDB logs:**
   ```bash
   docker-compose logs mongodb
   ```

3. **Restart MongoDB:**
   ```bash
   docker-compose restart mongodb
   ```

4. **Verify port mapping:**
   ```bash
   docker port ai_analysis_mongodb
   ```

---

### Problem: Authentication failed

**Symptoms:**
```
MongoServerError: Authentication failed
```

**Solution:**

Update `.env` with correct credentials:
```env
MONGODB_URI=mongodb://admin:password@localhost:27017
MONGODB_DATABASE=ai_analytics
```

Or disable auth in development:
```yaml
# docker-compose.yml
mongodb:
  environment:
    - MONGO_INITDB_ROOT_USERNAME=  # Remove
    - MONGO_INITDB_ROOT_PASSWORD=  # Remove
```

---

### Problem: Collection not found

**Symptoms:**
```
MongoServerError: ns not found
```

**Solution:**

Initialize database:
```bash
# Load schemas
mongosh mongodb://localhost:27017/ai_analytics < database/mongodb/schemas.js

# Create indexes
mongosh mongodb://localhost:27017/ai_analytics < database/mongodb/indexes.js

# Load sample data
mongosh mongodb://localhost:27017/ai_analytics < database/mongodb/sample_data.js
```

---

### Problem: Slow queries

**Diagnosis:**
```javascript
// Enable profiling
db.setProfilingLevel(2)

// Check slow queries
db.system.profile.find({millis: {$gt: 100}}).sort({ts: -1})
```

**Solutions:**

1. **Add missing indexes:**
   ```javascript
   db.orders.createIndex({ restaurant_id: 1, order_date: 1 })
   db.orders.createIndex({ status: 1 })
   ```

2. **Check index usage:**
   ```javascript
   db.orders.find({restaurant_id: "REST001"}).explain("executionStats")
   ```

3. **Optimize query:**
   ```javascript
   // Bad (fetches all fields)
   db.orders.find({restaurant_id: "REST001"})
   
   // Good (projection)
   db.orders.find(
     {restaurant_id: "REST001"},
     {total_amount: 1, order_date: 1}
   )
   ```

---

## 🔌 API Errors

### Problem: API returns 404 Not Found

**Symptoms:**
```
GET http://localhost:8080/api/v1/analytics/dashboard 404
```

**Diagnosis:**
```bash
# Check backend logs
docker-compose logs backend

# Test health endpoint
curl http://localhost:8080/health
```

**Solutions:**

1. **Check backend is running:**
   ```bash
   docker-compose ps backend
   ```

2. **Verify route exists:**
   ```bash
   # Check backend/cmd/api/main.go
   grep "dashboard" backend/cmd/api/main.go
   ```

3. **Check API base URL:**
   ```javascript
   // client/src/api/analytics.js
   const api = axios.create({
     baseURL: 'http://localhost:8080/api/v1'  // Correct
   });
   ```

---

### Problem: API returns 500 Internal Server Error

**Symptoms:**
```json
{
  "error": "Internal server error",
  "message": "Database connection failed"
}
```

**Diagnosis:**
```bash
# Check backend logs
docker-compose logs --tail=100 backend

# Check dependencies
docker-compose ps mongodb redis
```

**Common causes:**

1. **MongoDB not connected:**
   ```bash
   # Check connection string in .env
   MONGODB_URI=mongodb://localhost:27017
   
   # Test connection
   mongosh $MONGODB_URI
   ```

2. **Redis not connected:**
   ```bash
   # Test Redis
   docker-compose exec redis redis-cli ping
   # Should return: PONG
   ```

3. **Missing data:**
   ```bash
   # Verify restaurant exists
   mongosh mongodb://localhost:27017/ai_analytics
   db.restaurants.findOne({restaurant_id: "REST001"})
   ```

---

### Problem: CORS errors

**Symptoms:**
```
Access to XMLHttpRequest blocked by CORS policy
```

**Solution:**

Update backend CORS configuration:
```go
// backend/cmd/api/main.go
router.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders:     []string{"Content-Type", "Authorization"},
    AllowCredentials: true,
}))
```

Rebuild backend:
```bash
docker-compose up -d --build backend
```

---

### Problem: Timeout errors

**Symptoms:**
```
Error: timeout of 30000ms exceeded
```

**Solutions:**

1. **Increase timeout:**
   ```javascript
   // client/src/api/analytics.js
   const api = axios.create({
     baseURL: 'http://localhost:8080/api/v1',
     timeout: 60000  // 60 seconds
   });
   ```

2. **Check query performance:**
   ```bash
   # Enable slow query log
   docker-compose logs backend | grep "slow query"
   ```

3. **Add database indexes** (see Database Problems section)

---

## 🎨 Frontend Issues

### Problem: Frontend shows blank page

**Diagnosis:**
```bash
# Check browser console (F12)
# Look for errors

# Check frontend logs
docker-compose logs client
```

**Solutions:**

1. **Build issues:**
   ```bash
   cd client
   npm install
   npm run build
   ```

2. **API connection:**
   ```javascript
   // Check baseURL in client/src/api/analytics.js
   baseURL: 'http://localhost:8080/api/v1'
   ```

3. **Nginx configuration:**
   ```bash
   # Check client/nginx.conf
   docker-compose restart client
   ```

---

### Problem: Charts not rendering

**Symptoms:**
Empty chart containers or "undefined" errors.

**Common causes:**

1. **No data:**
   ```bash
   # Check API response
   curl http://localhost:8080/api/v1/analytics/dashboard?restaurant_id=REST001
   ```

2. **Data format mismatch:**
   ```javascript
   // RevenueChart expects:
   {
     labels: ["2026-01", "2026-02"],
     actual: [100000, 150000],
     predicted: [105000, 155000]
   }
   ```

3. **ECharts not loaded:**
   ```bash
   cd client
   npm install echarts echarts-for-react
   ```

---

### Problem: Styles not applied

**Solutions:**

1. **CSS not imported:**
   ```javascript
   // client/src/main.jsx
   import './App.css'
   ```

2. **Cache issue:**
   ```bash
   # Clear browser cache (Ctrl+Shift+Delete)
   # Or hard reload (Ctrl+Shift+R)
   ```

3. **Build CSS:**
   ```bash
   cd client
   npm run build
   ```

---

## 🤖 ML & Predictions

### Problem: No predictions generated

**Symptoms:**
```javascript
db.revenue_predictions.countDocuments()  // Returns 0
```

**Diagnosis:**
```bash
# Check prediction logs
docker-compose logs ml_prediction
```

**Solution:**

Run full pipeline:
```bash
# 1. Generate features
./manage.sh etl

# 2. Train model
./manage.sh train

# 3. Generate predictions
./manage.sh predict

# 4. Verify
mongosh mongodb://localhost:27017/ai_analytics
db.revenue_predictions.find()
```

---

### Problem: Model training fails

**Symptoms:**
```
ModuleNotFoundError: No module named 'prophet'
```

**Solution:**

1. **Install dependencies:**
   ```bash
   cd ml
   pip install -r requirements.txt
   ```

2. **Check Python version:**
   ```bash
   python --version  # Should be 3.8+
   ```

3. **Rebuild ML container:**
   ```bash
   docker-compose build ml_training
   docker-compose up -d ml_training
   ```

---

### Problem: Predictions are inaccurate

**Symptoms:**
MAPE > 20% or predictions don't match reality.

**Solutions:**

1. **Check data quality:**
   ```javascript
   // Look for missing data
   db.orders.countDocuments({
     restaurant_id: "REST001",
     status: "completed"
   })
   ```

2. **Need more historical data:**
   - Minimum: 12 months
   - Recommended: 24+ months

3. **Tune hyperparameters:**
   ```python
   # ml/config.py
   PROPHET_PARAMS = {
       'changepoint_prior_scale': 0.01,  # Lower = less flexible
       'seasonality_prior_scale': 5.0,
       'seasonality_mode': 'multiplicative'
   }
   ```

4. **Retrain model:**
   ```bash
   ./manage.sh train
   ./manage.sh predict
   ```

---

### Problem: ETL job hangs

**Symptoms:**
ETL runs for hours without completing.

**Diagnosis:**
```bash
# Check ETL logs
docker-compose logs etl_worker

# Check MongoDB performance
mongosh
db.currentOp()
```

**Solutions:**

1. **Add indexes:**
   ```javascript
   db.orders.createIndex({ restaurant_id: 1, order_date: 1 })
   ```

2. **Reduce batch size:**
   ```go
   // etl/internal/config/config.go
   ETLBatchSize: 500  // Default 1000
   ```

3. **Kill and restart:**
   ```bash
   docker-compose restart etl_worker
   ```

---

## ⚡ Performance Problems

### Problem: Slow API responses

**Diagnosis:**
```bash
# Measure response time
time curl http://localhost:8080/api/v1/analytics/dashboard?restaurant_id=REST001
```

**Solutions:**

1. **Enable caching:**
   ```bash
   # Check Redis is running
   docker-compose ps redis
   
   # Test cache
   docker-compose exec redis redis-cli
   redis> KEYS *
   ```

2. **Add database indexes:**
   ```javascript
   db.feature_revenue_monthly.createIndex({restaurant_id: 1, month: 1})
   db.revenue_predictions.createIndex({restaurant_id: 1, month: 1})
   ```

3. **Optimize queries:**
   ```go
   // Use projection
   projection := bson.M{"revenue": 1, "month": 1}
   cursor, err := collection.Find(ctx, filter, options.Find().SetProjection(projection))
   ```

4. **Increase cache TTL:**
   ```go
   // backend/internal/services/analytics_service.go
   cacheTTL := 30 * time.Minute  // Increase from default
   ```

---

### Problem: High memory usage

**Diagnosis:**
```bash
# Check container memory
docker stats

# Check MongoDB memory
docker-compose exec mongodb mongo --eval "db.serverStatus().mem"
```

**Solutions:**

1. **Limit MongoDB memory:**
   ```yaml
   # docker-compose.yml
   mongodb:
     command: --wiredTigerCacheSizeGB 1
   ```

2. **Optimize queries** (use projection and limits)

3. **Archive old data:**
   ```javascript
   // Archive orders older than 2 years
   db.orders.deleteMany({
     order_date: {$lt: new Date("2024-01-01")}
   })
   ```

---

### Problem: Disk space running out

**Diagnosis:**
```bash
df -h
docker system df
```

**Solution:**
```bash
# Remove unused Docker resources
docker system prune -a

# Compact MongoDB
docker-compose exec mongodb mongo ai_analytics --eval "db.runCommand({compact: 'orders'})"

# Archive old logs
find logs/ -name "*.log" -mtime +30 -delete
```

---

## 🌐 Deployment Issues

### Problem: Can't connect from external IP

**Symptoms:**
Works on localhost but not from other machines.

**Solution:**

1. **Bind to 0.0.0.0:**
   ```env
   # .env
   API_HOST=0.0.0.0  # Not 127.0.0.1
   ```

2. **Open firewall ports:**
   ```bash
   # Linux (ufw)
   sudo ufw allow 8080/tcp
   
   # AWS Security Group
   # Add inbound rule: TCP 8080 from 0.0.0.0/0
   ```

3. **Update CORS:**
   ```go
   AllowOrigins: []string{"*"}  // Or specific domain
   ```

---

### Problem: SSL/HTTPS issues

**Solution:**

Use Let's Encrypt with Nginx:

1. **Install Certbot:**
   ```bash
   sudo apt install certbot python3-certbot-nginx
   ```

2. **Get certificate:**
   ```bash
   sudo certbot --nginx -d yourdomain.com
   ```

3. **Auto-renewal:**
   ```bash
   sudo crontab -e
   # Add:
   0 0 * * * certbot renew --quiet
   ```

---

### Problem: Environment variable not working

**Symptoms:**
Changes to `.env` file have no effect.

**Solution:**

1. **Rebuild containers:**
   ```bash
   docker-compose down
   docker-compose up -d --build
   ```

2. **Verify variables:**
   ```bash
   docker-compose exec backend env | grep MONGODB
   ```

3. **Check .env location:**
   ```bash
   # Must be in same directory as docker-compose.yml
   ls -la .env
   ```

---

## 🆘 Getting More Help

### Enable Debug Logging

**Backend:**
```env
LOG_LEVEL=debug
```

**ETL:**
```env
ETL_LOG_LEVEL=debug
```

**ML:**
```python
# ml/config.py
LOG_LEVEL = "DEBUG"
```

### Collect Diagnostic Info

```bash
# Save all logs
./manage.sh logs > debug.log

# System info
docker-compose ps > system_status.txt
docker stats --no-stream >> system_status.txt
df -h >> system_status.txt

# Database stats
mongosh --eval "db.serverStatus()" > db_stats.txt
```

### Report an Issue

Include:
1. Error message and stack trace
2. Steps to reproduce
3. System info (OS, Docker version)
4. Relevant logs
5. Configuration (sanitize secrets!)

Submit at: [GitHub Issues](#)

---

## 📚 Additional Resources

- [FAQ](FAQ.md) - Common questions
- [Documentation](docs/) - Full documentation
- [Architecture](docs/architecture.md) - System design
- [API Reference](docs/api.md) - API details

---

**Last Updated:** February 15, 2026
