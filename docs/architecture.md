# Chi tiết Kiến trúc Hệ thống AI Analytics

## 📐 Tổng quan Kiến trúc

Hệ thống được thiết kế theo mô hình **microservices** với các tầng xử lý riêng biệt:

```
┌─────────────────────────────────────────────────────────────┐
│                         CLIENT LAYER                         │
│                    (React + ECharts)                         │
└────────────────────────┬────────────────────────────────────┘
                         │ HTTP/REST
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                      API GATEWAY LAYER                       │
│                      (Go + Gin)                              │
│  - Authentication & Authorization                            │
│  - Rate Limiting                                             │
│  - Response Caching (Redis)                                  │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                     SERVICE LAYER                            │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │  Analytics   │  │     ETL      │  │      ML      │      │
│  │   Service    │  │   Worker     │  │   Service    │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                      DATA LAYER                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   MongoDB    │  │    Redis     │  │    Model     │      │
│  │  (Primary)   │  │   (Cache)    │  │   Storage    │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
```

## 🔄 Luồng Dữ liệu Chi tiết

### 1. ETL Pipeline (Extract-Transform-Load)

**Tần suất**: Daily (2 AM)  
**Công nghệ**: Go Worker + Cron  
**Input**: Raw data từ MongoDB collections (orders, payments)  
**Output**: Feature Store (feature_revenue_monthly)

#### Các bước xử lý:

1. **Extract (Trích xuất)**
   ```
   - Đọc orders, payments từ MongoDB
   - Lọc theo trạng thái (status = "completed")
   - Aggregate theo restaurant_id và tháng
   ```

2. **Transform (Chuyển đổi)**
   ```
   - Tính toán raw metrics:
     * revenue (sum của total_price)
     * order_count
     * avg_order_value
   
   - Tính toán rolling features:
     * rolling_avg_3m
     * rolling_avg_6m
     * rolling_avg_12m
   
   - Tính toán growth features:
     * mom_growth (Month-over-Month)
     * yoy_growth (Year-over-Year)
   
   - Tính toán seasonality:
     * season (Q1, Q2, Q3, Q4)
     * is_holiday (boolean)
     * day_of_week_avg
   ```

3. **Load (Tải dữ liệu)**
   ```
   - Upsert vào feature_revenue_monthly collection
   - Log execution vào etl_logs collection
   ```

#### Code flow:
```go
func (e *RevenueFeatureETL) Run(ctx context.Context) error {
    // 1. Get list of active restaurants
    restaurants := e.getRestaurants(ctx)
    
    // 2. Process each restaurant
    for _, restaurantID := range restaurants {
        // 2.1. Calculate monthly features for last 12 months
        for month := startDate; month < now; month = month.AddMonth(1) {
            feature := e.calculateMonthlyFeatures(ctx, restaurantID, month)
            e.saveFeature(ctx, feature)
        }
    }
    
    // 3. Log execution
    e.saveETLLog(ctx, etlLog)
}
```

---

### 2. Model Training Pipeline

**Tần suất**: Weekly/On-demand  
**Công nghệ**: Python + Prophet/XGBoost  
**Input**: Feature Store  
**Output**: Model artifacts (.pkl files)

#### Training Process:

1. **Data Preparation**
   ```python
   # Load features from MongoDB
   df = db.fetch_features(restaurant_id)
   
   # Prophet format transformation
   prophet_df = pd.DataFrame({
       'ds': df['month'],
       'y': df['revenue']
   })
   
   # Add regressors
   prophet_df['rolling_avg_3m'] = df['rolling_avg_3m']
   prophet_df['is_holiday'] = df['is_holiday']
   ```

2. **Model Training**
   ```python
   # Initialize Prophet
   model = Prophet(
       yearly_seasonality=True,
       seasonality_mode='multiplicative',
       changepoint_prior_scale=0.05
   )
   
   # Add custom regressors
   model.add_regressor('rolling_avg_3m')
   model.add_regressor('is_holiday')
   
   # Fit model
   model.fit(prophet_df)
   ```

3. **Model Evaluation**
   ```python
   # Train-test split (80-20)
   split_idx = int(len(df) * 0.8)
   df_train = df[:split_idx]
   df_test = df[split_idx:]
   
   # Predict on test set
   forecast = model.predict(future)
   
   # Calculate metrics
   mape = mean_absolute_percentage_error(y_true, y_pred)
   rmse = sqrt(mean_squared_error(y_true, y_pred))
   mae = mean_absolute_error(y_true, y_pred)
   r2 = r2_score(y_true, y_pred)
   ```

4. **Model Persistence**
   ```python
   # Save model file
   joblib.dump(model, 'revenue_forecast_prophet_REST001_v1.0.0.pkl')
   
   # Save metadata to MongoDB
   db.save_model_metadata({
       'model_name': 'revenue_forecast_prophet_REST001',
       'version': 'v1.0.0',
       'metrics': {'mape': 5.2, 'rmse': 8500000, ...},
       'trained_at': datetime.now(),
       'is_production': True
   })
   ```

---

### 3. Batch Prediction Pipeline

**Tần suất**: Daily (3 AM)  
**Công nghệ**: Python + Trained Models  
**Input**: Trained models + Latest features  
**Output**: Predictions (revenue_predictions collection)

#### Prediction Process:

1. **Load Model**
   ```python
   model = joblib.load('revenue_forecast_prophet_REST001_v1.0.0.pkl')
   ```

2. **Generate Predictions**
   ```python
   # Create future dataframe (12 months ahead)
   future = model.make_future_dataframe(periods=12, freq='MS')
   
   # Add regressors for future dates
   future['is_holiday'] = future['ds'].dt.month.isin([1,2,4,5,9,12])
   future['rolling_avg_3m'] = last_known_value
   
   # Predict
   forecast = model.predict(future)
   ```

3. **Save Predictions**
   ```python
   for _, row in forecast.iterrows():
       prediction = {
           'restaurant_id': 'REST001',
           'month': row['ds'].strftime('%Y-%m'),
           'predicted': row['yhat'],
           'lower_ci': row['yhat_lower'],
           'upper_ci': row['yhat_upper'],
           'model_version': 'v1.0.0',
           'confidence_score': calculate_confidence(row),
           'predicted_at': datetime.now()
       }
       db.save_prediction(prediction)
   ```

---

### 4. API Serving Layer

**Technology**: Go + Gin Framework  
**Port**: 8080  
**Caching**: Redis

#### Key Endpoints:

##### 1. `/api/v1/analytics/forecast`

**Purpose**: Get revenue forecast with predictions

**Query Params**:
- `restaurant_id` (required): Restaurant identifier
- `months` (optional): Number of months to return (default: 12)

**Response Format**:
```json
{
  "restaurant_id": "REST001",
  "timestamps": ["2026-01", "2026-02", "2026-03"],
  "actual": [120000000, 135000000, 0],
  "predicted": [125000000, 135000000, 142000000],
  "target": [130000000, 130000000, 140000000],
  "confidence_data": [
    {
      "month": "2026-01",
      "lower": 115000000,
      "upper": 135000000
    }
  ],
  "model_info": {
    "model_name": "revenue_forecast_prophet_REST001",
    "model_version": "v1.0.0",
    "last_updated": "2026-02-15T02:00:00Z",
    "metrics": {
      "mape": 5.2,
      "rmse": 8500000,
      "mae": 6200000,
      "r2_score": 0.92
    }
  }
}
```

**Implementation Logic**:
```go
func (s *AnalyticsService) GetRevenueForecast(ctx, restaurantID, months) {
    // 1. Check cache
    cacheKey := fmt.Sprintf("forecast:%s:%d", restaurantID, months)
    if cached := redis.Get(cacheKey); cached != nil {
        return cached
    }
    
    // 2. Get predictions from MongoDB
    predictions := db.GetPredictions(restaurantID, months)
    
    // 3. Get actual data
    actuals := db.GetActuals(restaurantID, months)
    
    // 4. Merge data
    response := mergeForecastData(predictions, actuals)
    
    // 5. Cache result (TTL: 1 hour)
    redis.Set(cacheKey, response, 3600)
    
    return response
}
```

##### 2. `/api/v1/analytics/dashboard`

**Purpose**: Get complete dashboard data

**Response Format**:
```json
{
  "restaurant_id": "REST001",
  "summary": {
    "current_month_revenue": 135000000,
    "previous_month_revenue": 120000000,
    "month_over_month_growth": 12.5,
    "year_over_year_growth": 15.2,
    "total_orders": 2800,
    "avg_order_value": 48214,
    "forecast_next_month": 142000000,
    "forecast_confidence": 0.85
  },
  "revenue_chart": {
    "labels": ["2025-03", "2025-04", ...],
    "actual": [130M, 135M, ...],
    "predicted": [0, 0, 142M, 148M, ...],
    "target": [130M, 140M, ...]
  },
  "orders_chart": {
    "labels": ["2025-03", "2025-04", ...],
    "order_counts": [2500, 2800, ...]
  },
  "insights": [
    {
      "type": "success",
      "title": "Tăng trưởng mạnh",
      "description": "Doanh thu tháng này tăng mạnh so với tháng trước",
      "value": "+12.5%"
    }
  ]
}
```

---

### 5. Frontend Client Layer

**Technology**: React 18 + Vite + ECharts

#### Component Architecture:

```
App.jsx (Root)
  ├── SummaryCards.jsx (4 metric cards)
  ├── RevenueChart.jsx (Line chart with predictions)
  ├── OrdersChart.jsx (Bar chart)
  └── Insights.jsx (AI-generated insights)
```

#### Data Flow:

1. **Component Mount**
   ```jsx
   useEffect(() => {
       fetchDashboard();
   }, [restaurantId]);
   ```

2. **API Call**
   ```javascript
   const fetchDashboard = async () => {
       const data = await analyticsAPI.getDashboard(restaurantId);
       setDashboard(data);
   };
   ```

3. **Chart Rendering**
   ```jsx
   <ReactECharts 
       option={chartOption} 
       style={{ height: '400px' }}
   />
   ```

#### Chart Configuration Example:

```javascript
const option = {
  title: { text: 'Dự báo Doanh thu' },
  tooltip: {
    trigger: 'axis',
    formatter: (params) => {
      // Custom tooltip formatting
    }
  },
  legend: {
    data: ['Thực tế', 'Dự đoán', 'Mục tiêu']
  },
  xAxis: {
    type: 'category',
    data: timestamps
  },
  yAxis: {
    type: 'value',
    axisLabel: {
      formatter: (value) => (value / 1000000) + 'M'
    }
  },
  series: [
    {
      name: 'Thực tế',
      type: 'line',
      data: actualData,
      smooth: true
    },
    {
      name: 'Dự đoán',
      type: 'line',
      data: predictedData,
      lineStyle: { type: 'dashed' }
    }
  ]
};
```

---

## 🔐 Bảo mật & Hiệu suất

### Authentication Flow (Dự kiến mở rộng)

```
Client → Login API → JWT Token → Stored in localStorage
                                 ↓
                    Subsequent requests with Bearer token
                                 ↓
                    API validates token + RBAC check
```

### Caching Strategy

```
Layer 1: Browser Cache (static assets)
Layer 2: Redis Cache (API responses, TTL: 1 hour)
Layer 3: MongoDB (persistent storage)
```

### Performance Optimization

- **Backend**: Connection pooling, query optimization, indexes
- **Frontend**: Code splitting, lazy loading, image optimization
- **Database**: Compound indexes, TTL indexes for logs
- **Caching**: Redis for frequently accessed data

---

## 📊 Monitoring & Logging

### ETL Logs
```javascript
{
  job_name: "revenue_feature_engineering",
  status: "success",
  started_at: ISODate,
  completed_at: ISODate,
  duration: 120, // seconds
  records_processed: 1000,
  records_failed: 5
}
```

### API Logs (JSON format)
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

---

## 🚀 Scaling Strategy

### Horizontal Scaling

- **API Server**: Multiple instances behind Load Balancer
- **ETL Worker**: Distributed processing with message queue
- **ML Training**: Separate training cluster with GPU support

### Vertical Scaling

- **Database**: MongoDB sharding for large datasets
- **Cache**: Redis Cluster for high availability
- **Model Storage**: Cloud object storage (S3/GCS)

---

## 🔄 CI/CD Pipeline (Khuyến nghị)

```
Git Push → GitHub Actions
           ↓
    1. Run Tests
    2. Build Docker Images
    3. Push to Registry
    4. Deploy to Staging
    5. Run E2E Tests
    6. Deploy to Production
```

---

## 📈 Metrics & KPIs

### System Metrics
- API Response Time: < 100ms (p95)
- ETL Processing Time: < 30 minutes
- Model Training Time: < 2 hours
- Prediction Accuracy: MAPE < 10%

### Business Metrics
- Revenue Forecast Accuracy
- Order Volume Trends
- Growth Rate Analysis
- Seasonal Patterns
