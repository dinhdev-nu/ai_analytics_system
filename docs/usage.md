# Hướng dẫn Sử dụng từng Service

## 📊 ETL Service (Feature Engineering)

### Chức năng chính
ETL Service thực hiện:
1. Trích xuất dữ liệu raw từ MongoDB
2. Tính toán các features cho ML
3. Lưu vào Feature Store

### Chạy ETL thủ công

```bash
# Với Docker
docker-compose run --rm etl_worker

# Manual
cd etl
go run cmd/etl_worker/main.go
```

### Cấu hình Cron Schedule

Mặc định: Daily lúc 2 AM

```bash
# Chỉnh sửa trong .env
ETL_CRON_SCHEDULE="0 2 * * *"

# Hoặc chạy ngay lập tức (development)
ETL_CRON_SCHEDULE="@every 5m"
```

### Monitoring ETL Jobs

```javascript
// Xem logs ETL trong MongoDB
db.etl_logs.find().sort({started_at: -1}).limit(10)

// Xem job gần nhất
db.etl_logs.findOne({}, {sort: {started_at: -1}})

// Tìm failed jobs
db.etl_logs.find({status: "failed"})
```

### Troubleshooting

**Problem**: ETL job không chạy
```bash
# Check logs
docker-compose logs etl_worker

# Verify cron schedule
docker-compose exec etl_worker cat /proc/1/environ | grep ETL_CRON
```

**Problem**: Không có dữ liệu trong feature store
```bash
# Kiểm tra raw data có tồn tại
db.orders.countDocuments()
db.payments.countDocuments()

# Chạy ETL manually
docker-compose run --rm etl_worker
```

---

## 🤖 ML Training Service

### Model Training Workflow

```
Raw Data → Features → Train-Test Split → Model Training → Evaluation → Save Model
```

### Chạy Training

```bash
# Docker
docker-compose run --rm ml_training

# Manual
cd ml
python training/train_revenue_forecast.py
```

### Model Configuration

**Prophet Parameters** (trong `train_revenue_forecast.py`):
```python
model = Prophet(
    yearly_seasonality=True,      # Detect yearly patterns
    weekly_seasonality=False,     # No weekly (monthly data)
    daily_seasonality=False,      # No daily (monthly data)
    seasonality_mode='multiplicative',  # Multiplicative seasonality
    changepoint_prior_scale=0.05  # Flexibility of trend changes
)
```

### Thêm Custom Regressors

```python
# Trong train_revenue_forecast.py
model.add_regressor('rolling_avg_3m')
model.add_regressor('is_holiday')
model.add_regressor('marketing_spend')  # Thêm mới
```

### Model Evaluation

Metrics được tính tự động:
- **MAPE** (Mean Absolute Percentage Error): < 10% là tốt
- **RMSE** (Root Mean Square Error): So sánh giữa các models
- **MAE** (Mean Absolute Error): Dễ hiểu hơn RMSE
- **R² Score**: 0-1, càng gần 1 càng tốt

### Retrain Models

```bash
# Retrain tất cả restaurants
docker-compose run --rm ml_training

# Retrain specific restaurant (future feature)
docker-compose run --rm ml_training python training/train_revenue_forecast.py --restaurant REST001
```

### Model Versioning

Models được version tự động:
```
revenue_forecast_prophet_REST001_v1.0.0.pkl
revenue_forecast_prophet_REST001_v1.0.1.pkl
```

Version được lưu trong MongoDB:
```javascript
db.ml_models.find({model_name: "revenue_forecast_prophet_REST001"})
```

---

## 🔮 ML Prediction Service

### Batch Prediction Workflow

```
Load Model → Generate Future Dates → Predict → Calculate Confidence → Save to DB
```

### Chạy Prediction

```bash
# Docker
docker-compose run --rm ml_prediction

# Manual
cd ml
python prediction/batch_predict.py
```

### Prediction Configuration

```python
# Số tháng dự đoán
PREDICTION_MONTHS_AHEAD = 12  # Dự đoán 12 tháng

# Confidence calculation
confidence = 1.0 - (interval_width / predicted_value) / 2
```

### Xem Predictions

```javascript
// Xem predictions của 1 restaurant
db.revenue_predictions.find({
    restaurant_id: "REST001"
}).sort({month: 1})

// So sánh predicted vs actual
db.revenue_predictions.aggregate([
    {$match: {actual: {$ne: null}}},
    {$project: {
        month: 1,
        predicted: 1,
        actual: 1,
        error: {$abs: {$subtract: ["$predicted", "$actual"]}},
        error_pct: {
            $multiply: [
                {$divide: [
                    {$abs: {$subtract: ["$predicted", "$actual"]}},
                    "$actual"
                ]},
                100
            ]
        }
    }}
])
```

### Update Actuals

Predictions được update với actual values khi dữ liệu thực tế có sẵn:

```bash
# Chạy update actuals
python prediction/batch_predict.py --update-actuals
```

---

## 🚀 Backend API Service

### Khởi động API

```bash
# Docker
docker-compose up -d backend

# Manual
cd backend
go run cmd/api/main.go
```

### API Endpoints

#### 1. Health Check
```bash
curl http://localhost:8080/health
```

#### 2. Get Forecast
```bash
curl "http://localhost:8080/api/v1/analytics/forecast?restaurant_id=REST001&months=12"
```

#### 3. Get Dashboard
```bash
curl "http://localhost:8080/api/v1/analytics/dashboard?restaurant_id=REST001"
```

### Caching với Redis

API sử dụng Redis để cache responses:

```go
// Cache key format
forecast:{restaurant_id}:{months}
dashboard:{restaurant_id}

// TTL
forecast: 1 hour
dashboard: 30 minutes
```

Clear cache:
```bash
# Clear all cache
docker-compose exec redis redis-cli FLUSHALL

# Clear specific key
docker-compose exec redis redis-cli DEL "forecast:REST001:12"
```

### Rate Limiting

Mặc định: 100 requests/phút per IP

Chỉnh sửa trong `.env`:
```bash
API_RATE_LIMIT=100
```

### Logging

Logs được xuất ra theo format JSON:

```json
{
  "level": "info",
  "time": "2026-02-15T10:30:00Z",
  "method": "GET",
  "path": "/api/v1/analytics/dashboard",
  "status": 200,
  "latency_ms": 45,
  "restaurant_id": "REST001"
}
```

Xem logs:
```bash
docker-compose logs -f backend
```

### Performance Monitoring

```bash
# Check response time
time curl "http://localhost:8080/api/v1/analytics/dashboard?restaurant_id=REST001"

# Monitor with Apache Bench
ab -n 1000 -c 10 "http://localhost:8080/api/v1/analytics/dashboard?restaurant_id=REST001"
```

---

## 💻 Frontend Client

### Khởi động Client

```bash
# Development
cd client
npm install
npm run dev

# Production
npm run build
npm run preview
```

### Component Structure

```
src/
├── api/
│   └── analytics.js          # API client
├── components/
│   ├── RevenueChart.jsx      # Main forecast chart
│   ├── OrdersChart.jsx       # Orders bar chart
│   ├── SummaryCards.jsx      # Metric cards
│   └── Insights.jsx          # AI insights
├── App.jsx                   # Main app component
└── main.jsx                  # Entry point
```

### Chart Customization

**Revenue Chart (ECharts)**

```javascript
// In src/components/RevenueChart.jsx
const option = {
  // Customize colors
  color: ['#3b82f6', '#10b981', '#f59e0b'],
  
  // Customize title
  title: {
    text: 'Dự báo Doanh thu',
    textStyle: {
      fontSize: 20,
      fontWeight: 600,
      color: '#333'
    }
  },
  
  // Add gradient
  series: [{
    areaStyle: {
      color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
        { offset: 0, color: 'rgba(59, 130, 246, 0.3)' },
        { offset: 1, color: 'rgba(59, 130, 246, 0)' }
      ])
    }
  }]
};
```

### Thêm Restaurant mới

```javascript
// In src/App.jsx
<select value={restaurantId} onChange={(e) => setRestaurantId(e.target.value)}>
  <option value="REST001">Nhà hàng Phương Nam</option>
  <option value="REST002">Quán Ăn Sài Gòn</option>
  <option value="REST003">BBQ House</option>
  <option value="REST004">Pizza Italia</option>  {/* Thêm mới */}
</select>
```

### State Management

Currently using React Hooks (useState, useEffect).

For larger apps, consider:
- **Redux**: Global state management
- **React Query**: API state management
- **Zustand**: Lightweight state management

### Error Handling

```javascript
// In src/App.jsx
try {
  const data = await analyticsAPI.getDashboard(restaurantId);
  setDashboard(data);
} catch (err) {
  setError(err.message || 'Failed to fetch dashboard data');
  
  // Show toast notification
  toast.error('Không thể tải dữ liệu dashboard');
}
```

---

## 🗄️ Database Management

### MongoDB Commands

```javascript
// Connect
mongosh "mongodb://localhost:27017"
use ai_analytics

// View collections
show collections

// Count documents
db.orders.countDocuments()
db.feature_revenue_monthly.countDocuments()
db.revenue_predictions.countDocuments()

// Latest feature data
db.feature_revenue_monthly.find().sort({month: -1}).limit(10)

// Check restaurant data
db.orders.aggregate([
    {$match: {restaurant_id: "REST001"}},
    {$group: {
        _id: null,
        total_revenue: {$sum: "$total_price"},
        total_orders: {$sum: 1}
    }}
])
```

### Backup & Restore

```bash
# Backup
mongodump --uri="mongodb://localhost:27017" --db=ai_analytics --out=backup/

# Restore
mongorestore --uri="mongodb://localhost:27017" backup/ai_analytics/

# Export to JSON
mongoexport --uri="mongodb://localhost:27017" --db=ai_analytics --collection=orders --out=orders.json

# Import from JSON
mongoimport --uri="mongodb://localhost:27017" --db=ai_analytics --collection=orders --file=orders.json
```

### Data Cleanup

```javascript
// Delete old logs (older than 90 days)
db.etl_logs.deleteMany({
    created_at: {$lt: new Date(Date.now() - 90*24*60*60*1000)}
})

// Delete old predictions
db.revenue_predictions.deleteMany({
    predicted_at: {$lt: new Date("2025-01-01")}
})
```

---

## 🎯 Best Practices

### ETL
- Chạy ETL sau khi có đủ data mới (sau 12 AM)
- Monitor ETL logs thường xuyên
- Validate data quality trước khi training

### ML Training
- Retrain models khi có đủ data mới (weekly/monthly)
- Track model performance metrics
- A/B test models trước khi deploy production

### Prediction
- Chạy prediction sau khi training xong
- Verify prediction quality
- Update actuals để track accuracy

### API
- Enable caching cho performance
- Monitor response times
- Set proper rate limits

### Frontend
- Implement error boundaries
- Add loading states
- Cache API responses (React Query)

---

## 🔧 Maintenance Tasks

### Daily
```bash
# Check ETL completion
docker-compose logs etl_worker | grep "completed"

# Verify predictions updated
curl "http://localhost:8080/api/v1/analytics/forecast?restaurant_id=REST001" | jq '.model_info.last_updated'
```

### Weekly
```bash
# Retrain models
docker-compose run --rm ml_training

# Review logs
docker-compose logs --tail=100 backend
```

### Monthly
```bash
# Full backup
./scripts/backup.sh

# Review model performance
python ml/scripts/evaluate_models.py

# Update dependencies
cd backend && go get -u ./...
cd client && npm update
```

---

## 📚 Additional Resources

- [Prophet Documentation](https://facebook.github.io/prophet/)
- [ECharts Examples](https://echarts.apache.org/examples/)
- [Go Gin Framework](https://gin-gonic.com/)
- [MongoDB Aggregation](https://docs.mongodb.com/manual/aggregation/)
