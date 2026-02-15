# Quick Start Guide - 5 Phút Setup

## ⚡ Yêu cầu hệ thống

- Docker & Docker Compose
- Git
- Ports: 3000, 8080, 27017, 6379 (available)

## 🚀 Setup trong 5 bước

### Bước 1: Clone & Setup Environment (1 phút)

```bash
# Clone repository
git clone <your-repo-url>
cd AI_analysis

# Copy environment file
cp .env.example .env

# (Optional) Chỉnh sửa .env nếu cần
```

### Bước 2: Khởi động Services (2 phút)

```bash
# Build và start tất cả services
docker-compose up -d

# Xem logs để verify
docker-compose logs -f
```

Output mong đợi:
```
✓ mongodb running
✓ redis running
✓ backend running
✓ client running
```

### Bước 3: Import Sample Data (30 giây)

```bash
# Exec vào MongoDB container
docker exec -it ai_analytics_mongodb mongosh

# Trong mongosh, chạy:
use ai_analytics
load("/docker-entrypoint-initdb.d/schemas.js")
load("/docker-entrypoint-initdb.d/indexes.js")
load("/docker-entrypoint-initdb.d/sample_data.js")
exit
```

### Bước 4: Run ETL & Training (1 phút)

```bash
# Chạy ETL để tạo features
docker-compose run --rm etl_worker

# Train models
docker-compose run --rm ml_training

# Generate predictions
docker-compose run --rm ml_prediction
```

### Bước 5: Mở Dashboard (10 giây)

```bash
# Mở browser
http://localhost:3000

# Test API
curl http://localhost:8080/api/v1/analytics/dashboard?restaurant_id=REST001
```

## ✅ Verify Installation

### Test Backend API
```bash
curl http://localhost:8080/health
# Response: {"success":true,"message":"API is running"}
```

### Test MongoDB
```bash
docker exec -it ai_analytics_mongodb mongosh --eval "db.adminCommand('ping')"
# Response: { ok: 1 }
```

### Test Redis
```bash
docker exec -it ai_analytics_redis redis-cli ping
# Response: PONG
```

### Check Frontend
```
Browser: http://localhost:3000
Should see: AI Analytics Dashboard with charts
```

## 🎯 Next Steps

1. **Explore Dashboard**: Thử select các restaurant khác nhau
2. **View API Response**: Check Network tab trong browser
3. **Review Code**: Bắt đầu từ `client/src/App.jsx`
4. **Read Docs**: Xem [Architecture](./architecture.md) để hiểu hệ thống

## 🐛 Common Issues

### Port Already in Use
```bash
# Find process using port
netstat -ano | findstr :8080  # Windows
lsof -i :8080                # Mac/Linux

# Kill process or change port in .env
```

### Docker Build Failed
```bash
# Clear Docker cache
docker system prune -a

# Rebuild
docker-compose build --no-cache
docker-compose up -d
```

### MongoDB Connection Failed
```bash
# Check MongoDB is running
docker-compose ps mongodb

# Restart MongoDB
docker-compose restart mongodb

# Check logs
docker-compose logs mongodb
```

### No Data in Dashboard
```bash
# Verify sample data loaded
docker exec -it ai_analytics_mongodb mongosh
use ai_analytics
db.orders.countDocuments()  # Should return > 0

# Re-run ETL and prediction
docker-compose run --rm etl_worker
docker-compose run --rm ml_prediction
```

## 📞 Need Help?

1. Check [Deployment Guide](./deployment.md)
2. Check [Usage Guide](./usage.md)
3. Review logs: `docker-compose logs [service-name]`
4. Open an issue on GitHub

## 🎉 You're Ready!

Hệ thống đã sẵn sàng. Bắt đầu khám phá và phát triển!

### Development Mode
```bash
# Backend
cd backend
go run cmd/api/main.go

# Frontend
cd client
npm run dev

# ML
cd ml
python training/train_revenue_forecast.py
```

### Production Mode
```bash
# Deploy with Docker Compose
docker-compose -f docker-compose.prod.yml up -d

# Or use Kubernetes
kubectl apply -f k8s/
```

Happy Coding! 🚀
