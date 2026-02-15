# Hướng dẫn Triển khai (Deployment Guide)

## 🚀 Phương pháp Triển khai

### Phương pháp 1: Docker Compose (Khuyến nghị cho Development/Staging)

#### Bước 1: Chuẩn bị môi trường

```bash
# Clone repository
git clone <repository-url>
cd AI_analysis

# Copy environment file
cp .env.example .env

# Chỉnh sửa .env với thông tin thực tế
nano .env
```

#### Bước 2: Build và khởi động services

```bash
# Build tất cả images
docker-compose build

# Khởi động tất cả services
docker-compose up -d

# Xem logs
docker-compose logs -f

# Kiểm tra status
docker-compose ps
```

#### Bước 3: Khởi tạo Database

```bash
# Kết nối MongoDB
docker exec -it ai_analytics_mongodb mongosh

# Chạy trong mongosh
use ai_analytics

# Load schemas và indexes
load("/docker-entrypoint-initdb.d/schemas.js")
load("/docker-entrypoint-initdb.d/indexes.js")
load("/docker-entrypoint-initdb.d/sample_data.js")
```

#### Bước 4: Training Model lần đầu

```bash
# Chạy training một lần
docker-compose run --rm ml_training

# Kiểm tra models đã được tạo
docker-compose exec ml_prediction ls -la /app/models
```

#### Bước 5: Verify deployment

```bash
# Health check Backend
curl http://localhost:8080/health

# Test Analytics API
curl "http://localhost:8080/api/v1/analytics/dashboard?restaurant_id=REST001"

# Mở browser
# Frontend: http://localhost:3000
# MongoDB Express (optional): http://localhost:8081
```

---

### Phương pháp 2: Manual Deployment (Development)

#### 1. MongoDB Setup

```bash
# Cài đặt MongoDB
# Windows: Download từ mongodb.com
# Linux: sudo apt-get install mongodb-org
# Mac: brew install mongodb-community

# Khởi động MongoDB
mongod --dbpath /data/db

# Import sample data
mongoimport --db ai_analytics --collection restaurants --file database/mongodb/sample_data.js
```

#### 2. Redis Setup

```bash
# Windows: Download từ redis.io
# Linux: sudo apt-get install redis-server
# Mac: brew install redis

# Khởi động Redis
redis-server
```

#### 3. ETL Worker

```bash
cd etl

# Install dependencies
go mod download

# Copy environment
cp ../.env .env

# Run ETL worker
go run cmd/etl_worker/main.go
```

#### 4. ML Services

```bash
cd ml

# Create virtual environment
python -m venv venv
source venv/bin/activate  # Windows: venv\Scripts\activate

# Install dependencies
pip install -r requirements.txt

# Copy environment
cp ../.env .env

# Run training
python training/train_revenue_forecast.py

# Run prediction
python prediction/batch_predict.py
```

#### 5. Backend API

```bash
cd backend

# Install dependencies
go mod download

# Copy environment
cp ../.env .env

# Run API server
go run cmd/api/main.go
```

#### 6. Frontend Client

```bash
cd client

# Install dependencies
npm install

# Copy environment
echo "VITE_API_URL=http://localhost:8080/api/v1" > .env

# Run development server
npm run dev

# Build for production
npm run build
```

---

### Phương pháp 3: Kubernetes Deployment (Production)

#### Bước 1: Tạo Kubernetes manifests

**mongodb-deployment.yaml**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mongodb
spec:
  replicas: 1
  selector:
    matchLabels:
      app: mongodb
  template:
    metadata:
      labels:
        app: mongodb
    spec:
      containers:
      - name: mongodb
        image: mongo:7.0
        ports:
        - containerPort: 27017
        env:
        - name: MONGO_INITDB_ROOT_USERNAME
          valueFrom:
            secretKeyRef:
              name: mongodb-secret
              key: username
        - name: MONGO_INITDB_ROOT_PASSWORD
          valueFrom:
            secretKeyRef:
              name: mongodb-secret
              key: password
        volumeMounts:
        - name: mongodb-storage
          mountPath: /data/db
      volumes:
      - name: mongodb-storage
        persistentVolumeClaim:
          claimName: mongodb-pvc
---
apiVersion: v1
kind: Service
metadata:
  name: mongodb
spec:
  selector:
    app: mongodb
  ports:
  - port: 27017
    targetPort: 27017
```

**backend-deployment.yaml**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: backend-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: backend-api
  template:
    metadata:
      labels:
        app: backend-api
    spec:
      containers:
      - name: api
        image: ai-analytics-backend:latest
        ports:
        - containerPort: 8080
        env:
        - name: MONGODB_URI
          valueFrom:
            secretKeyRef:
              name: mongodb-secret
              key: uri
        - name: REDIS_HOST
          value: redis
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: backend-api
spec:
  type: LoadBalancer
  selector:
    app: backend-api
  ports:
  - port: 80
    targetPort: 8080
```

#### Bước 2: Deploy to Kubernetes

```bash
# Create namespace
kubectl create namespace ai-analytics

# Create secrets
kubectl create secret generic mongodb-secret \
  --from-literal=username=admin \
  --from-literal=password=SecurePassword123 \
  --from-literal=uri=mongodb://admin:SecurePassword123@mongodb:27017 \
  -n ai-analytics

# Apply manifests
kubectl apply -f k8s/ -n ai-analytics

# Check status
kubectl get pods -n ai-analytics
kubectl get services -n ai-analytics

# View logs
kubectl logs -f deployment/backend-api -n ai-analytics
```

---

## 🔧 Configuration Management

### Environment Variables by Service

#### ETL Worker
```bash
MONGODB_URI=mongodb://admin:password@mongodb:27017
MONGODB_DATABASE=ai_analytics
ETL_CRON_SCHEDULE="0 2 * * *"
ETL_BATCH_SIZE=1000
ENVIRONMENT=production
LOG_LEVEL=info
```

#### ML Services
```bash
MONGODB_URI=mongodb://admin:password@mongodb:27017
MONGODB_DATABASE=ai_analytics
ML_MODEL_PATH=/app/models
MLFLOW_TRACKING_URI=http://mlflow:5000
PREDICTION_MONTHS_AHEAD=12
ENVIRONMENT=production
```

#### Backend API
```bash
API_PORT=8080
API_HOST=0.0.0.0
MONGODB_URI=mongodb://admin:password@mongodb:27017
MONGODB_DATABASE=ai_analytics
REDIS_HOST=redis
REDIS_PORT=6379
JWT_SECRET=your-super-secret-jwt-key-change-in-production
API_RATE_LIMIT=100
ENVIRONMENT=production
LOG_LEVEL=info
```

#### Client
```bash
VITE_API_URL=https://api.yourdomain.com/api/v1
```

---

## 🔐 Security Checklist

### Pre-deployment Security

- [ ] Change default MongoDB credentials
- [ ] Change JWT secret key
- [ ] Enable MongoDB authentication
- [ ] Set up Redis password
- [ ] Configure firewall rules
- [ ] Enable HTTPS/TLS
- [ ] Set up rate limiting
- [ ] Configure CORS properly
- [ ] Enable audit logging
- [ ] Set up backup strategy

### SSL/TLS Setup

```bash
# Generate self-signed certificate (Development)
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout server.key -out server.crt

# Use Let's Encrypt (Production)
certbot certonly --standalone -d yourdomain.com
```

---

## 📊 Monitoring Setup

### Prometheus + Grafana (Khuyến nghị)

**prometheus.yml**
```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'backend-api'
    static_configs:
      - targets: ['backend:8080']
  
  - job_name: 'mongodb'
    static_configs:
      - targets: ['mongodb-exporter:9216']
  
  - job_name: 'redis'
    static_configs:
      - targets: ['redis-exporter:9121']
```

### Logging with ELK Stack

```bash
# Filebeat config for Backend logs
filebeat.inputs:
- type: log
  enabled: true
  paths:
    - /var/log/backend/*.log
  json.keys_under_root: true

output.elasticsearch:
  hosts: ["elasticsearch:9200"]
```

---

## 🔄 Backup & Recovery

### MongoDB Backup

```bash
# Daily backup script
#!/bin/bash
DATE=$(date +%Y%m%d)
BACKUP_DIR="/backups/mongodb"

mongodump --uri="mongodb://admin:password@mongodb:27017" \
  --out="$BACKUP_DIR/$DATE"

# Compress
tar -czf "$BACKUP_DIR/$DATE.tar.gz" "$BACKUP_DIR/$DATE"
rm -rf "$BACKUP_DIR/$DATE"

# Keep only last 7 days
find $BACKUP_DIR -name "*.tar.gz" -mtime +7 -delete
```

### Model Backup

```bash
# Sync models to S3/Cloud Storage
aws s3 sync /app/models s3://your-bucket/ml-models/$(date +%Y%m%d)/
```

---

## 🚨 Troubleshooting

### Common Issues

#### 1. MongoDB Connection Failed
```bash
# Check MongoDB is running
docker-compose ps mongodb

# Check logs
docker-compose logs mongodb

# Test connection
mongosh "mongodb://admin:password@localhost:27017"
```

#### 2. Backend API not responding
```bash
# Check backend logs
docker-compose logs backend

# Check if port is occupied
netstat -tulpn | grep 8080

# Test endpoint
curl http://localhost:8080/health
```

#### 3. ML Training failed
```bash
# Check Python dependencies
docker-compose exec ml_training pip list

# Check data availability
docker-compose exec ml_training python -c "from database import db; print(len(db.get_all_restaurants()))"

# Check logs
docker-compose logs ml_training
```

#### 4. Frontend can't connect to API
```bash
# Check CORS configuration in backend
# Check API URL in client .env
# Check browser console for errors
```

---

## 📈 Performance Tuning

### MongoDB Optimization

```javascript
// Create indexes
db.orders.createIndex({ restaurant_id: 1, created_at: -1 })
db.revenue_predictions.createIndex({ restaurant_id: 1, month: 1 })

// Enable profiling
db.setProfilingLevel(1, { slowms: 100 })
```

### Redis Optimization

```bash
# Adjust maxmemory
redis-cli CONFIG SET maxmemory 2gb
redis-cli CONFIG SET maxmemory-policy allkeys-lru
```

### Backend Optimization

```go
// Connection pool settings
clientOptions := options.Client().
    ApplyURI(mongoURI).
    SetMaxPoolSize(100).
    SetMinPoolSize(10)
```

---

## ✅ Health Check Endpoints

```bash
# Backend Health
curl http://localhost:8080/health

# MongoDB Health
docker-compose exec mongodb mongosh --eval "db.adminCommand('ping')"

# Redis Health
docker-compose exec redis redis-cli ping
```

---

## 📞 Support & Maintenance

### Daily Checks
- [ ] Check ETL job completion
- [ ] Verify ML predictions updated
- [ ] Monitor API response times
- [ ] Check error logs
- [ ] Verify backups completed

### Weekly Checks
- [ ] Review model performance metrics
- [ ] Check disk space usage
- [ ] Review security logs
- [ ] Update dependencies

### Monthly Checks
- [ ] Retrain ML models
- [ ] Performance optimization review
- [ ] Security audit
- [ ] Capacity planning
