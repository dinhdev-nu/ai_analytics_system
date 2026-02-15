# AI Analytics System - Production Ready

## 🎯 Mục tiêu
Hệ thống phân tích dữ liệu kinh doanh với khả năng dự đoán (forecasting) dựa trên dữ liệu lịch sử. Backend xử lý toàn bộ logic AI/ML, client chỉ nhận dữ liệu đã xử lý để render dashboard.

## 🏗️ Kiến trúc tổng quan

```
Database (MongoDB) 
    ↓
ETL / Feature Engineering (Go)
    ↓
Model Training (Python - Offline)
    ↓
Batch Prediction (Python - Cron)
    ↓
Backend API (Go + Gin)
    ↓
Client (React + ECharts)
```

## 📁 Cấu trúc dự án

```
AI_analysis/
├── backend/               # Go API Server (port 8080)
├── etl/                   # Go ETL Workers
├── ml/                    # Python ML Training & Prediction
├── client/                # React Dashboard
├── database/              # Database schemas
├── docker/                # Docker configurations
└── docs/                  # Documentation
```

## 🚀 Tech Stack

| Tầng | Công nghệ |
|------|-----------|
| Database | MongoDB |
| ETL | Go + Cron |
| Feature Store | MongoDB Collections |
| Training | Python, XGBoost, Prophet |
| Prediction | Batch Job (Python) |
| API | Go + Gin Framework |
| Client | React + ECharts |
| Cache | Redis |
| Container | Docker + Docker Compose |

## 🔧 Quick Start

### Sử dụng Management Script (Đơn giản nhất)

```bash
# Linux/Mac
chmod +x manage.sh
./manage.sh setup      # Setup lần đầu (tự động làm mọi thứ)
./manage.sh start      # Khởi động services
./manage.sh health     # Kiểm tra health

# Windows
manage.bat setup       # Setup lần đầu
manage.bat start       # Khởi động services
manage.bat health      # Kiểm tra health
```

### Cách 1: Docker Compose (Khuyến nghị - 5 phút)

```bash
# Clone repository
git clone <your-repo-url>
cd AI_analysis

# Copy và chỉnh sửa environment
cp .env.example .env

# Khởi động tất cả services
docker-compose up -d

# Import sample data
docker exec -it ai_analytics_mongodb mongosh
use ai_analytics
load("/docker-entrypoint-initdb.d/sample_data.js")
exit

# Chạy ETL và Training
docker-compose run --rm etl_worker
docker-compose run --rm ml_training
docker-compose run --rm ml_prediction

# Mở dashboard
open http://localhost:3000
```

➡️ **Xem chi tiết**: [Quick Start Guide](docs/quickstart.md)

### Cách 2: Manual Setup (Development)

```bash
# 1. MongoDB
mongod --dbpath /data/db

# 2. Redis
redis-server

# 3. ETL Worker
cd etl && go run cmd/etl_worker/main.go

# 4. ML Training
cd ml && python training/train_revenue_forecast.py

# 5. Backend API
cd backend && go run cmd/api/main.go

# 6. Frontend
cd client && npm install && npm run dev
```

➡️ **Xem chi tiết**: [Deployment Guide](docs/deployment.md)

## 📊 Luồng dữ liệu chính

### 1. ETL Pipeline (Daily)
- Đọc raw data từ MongoDB (orders, payments)
- Tính toán features: rolling average, growth rate, seasonality
- Lưu vào Feature Store

### 2. Model Training (Weekly/Monthly)
- Load features từ Feature Store
- Train model (Prophet/XGBoost)
- Evaluate và lưu model artifact
- Track version với MLflow

### 3. Batch Prediction (Daily)
- Load model đã train
- Predict 3-12 tháng tiếp theo
- Lưu predictions vào MongoDB

### 4. API Serving (Real-time)
- Query actual data + predictions
- Merge và format dữ liệu
- Return JSON render-ready

### 5. Client Rendering
- Gọi API endpoints
- Render charts với ECharts
- Interactive dashboard

## 🔑 API Endpoints

```
GET  /api/v1/analytics/revenue?restaurant_id=xxx
GET  /api/v1/analytics/orders?restaurant_id=xxx
GET  /api/v1/analytics/forecast?restaurant_id=xxx&months=6
GET  /api/v1/analytics/dashboard?restaurant_id=xxx
```

## 📈 Database Collections

### Raw Data
- `orders` - Đơn hàng
- `payments` - Thanh toán
- `restaurants` - Nhà hàng

### Feature Store
- `feature_revenue_monthly` - Features cho revenue forecast
- `feature_order_patterns` - Features cho order prediction

### Predictions
- `revenue_predictions` - Kết quả dự đoán doanh thu
- `order_predictions` - Kết quả dự đoán đơn hàng

## 🧪 Testing

```bash
# Test Backend
cd backend && go test ./...

# Test ETL
cd etl && go test ./...
### 📖 Getting Started
- [⚡ Quick Start (5 phút)](docs/quickstart.md) - Setup nhanh nhất
- [🚀 Deployment Guide](docs/deployment.md) - Hướng dẫn triển khai đầy đủ
- [📘 Usage Guide](docs/usage.md) - Hướng dẫn sử dụng từng service

### 🏗️ Technical Documentation
- [🏛️ Architecture Details](docs/architecture.md) - Kiến trúc chi tiết
- [🔌 API Documentation](docs/api.md) - API endpoints và examples

### 📊 Features by Component
- **ETL Service**: Feature engineering, data transformation
- **ML Training**: Prophet/XGBoost models, evaluation metrics
- **Prediction Service**: Batch forecasting, confidence intervals
- **Backend API**: REST endpoints, caching, rate limiting
- **Frontend**: React dashboard, ECharts visualization
## 📦 Deployment

```bash
# Build all services
docker-compose build

# Deploy
docker-compose up -d
```

## 📚 Documentation

- [Architecture Details](docs/architecture.md)
- [ETL Pipeline](docs/etl.md)
- [ML Models](docs/ml_models.md)
- [API Documentation](docs/api.md)
- [Frontend Guide](docs/frontend.md)

## ⚡ Performance

- API Response: < 100ms (cached)
- Batch Prediction: ~10 phút / 1000 restaurants
- ETL Processing: ~30 phút / ngày
- Model Training: ~2 giờ / lần

## 🔒 Security

- JWT Authentication
- RBAC (Role-Based Access Control)
- API Rate Limiting
- Data Encryption at rest

## 📝 License

MIT License
