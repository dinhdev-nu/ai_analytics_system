# Database Layer

MongoDB schemas, indexes, and sample data for AI Analytics system.

## 📁 Structure

```
database/
└── mongodb/
    ├── schemas.js         # Collection schemas
    ├── indexes.js         # Performance indexes
    └── sample_data.js     # Sample data generator
```

## 🗄️ Collections

### 1. restaurants

Restaurant master data.

```javascript
{
  _id: ObjectId("..."),
  restaurant_id: "REST001",        // Unique identifier
  name: "Nhà hàng ABC",
  address: "123 Nguyễn Huệ, Q1, TP.HCM",
  phone: "0901234567",
  email: "contact@restaurant.com",
  owner_id: "USR001",
  status: "active",                // active | inactive | suspended
  created_at: ISODate("2024-01-01T00:00:00Z"),
  updated_at: ISODate("2026-02-15T00:00:00Z")
}
```

**Indexes:**
```javascript
db.restaurants.createIndex({ restaurant_id: 1 }, { unique: true })
db.restaurants.createIndex({ status: 1 })
db.restaurants.createIndex({ owner_id: 1 })
```

### 2. orders

Raw order transactions.

```javascript
{
  _id: ObjectId("..."),
  order_id: "ORD20260215001",
  restaurant_id: "REST001",
  customer_id: "CUST001",
  order_date: ISODate("2026-02-15T12:30:00Z"),
  items: [
    {
      item_id: "ITEM001",
      name: "Phở bò",
      quantity: 2,
      unit_price: 50000,
      subtotal: 100000
    }
  ],
  subtotal: 100000,              // Before tax/discount
  tax: 10000,                    // 10% VAT
  discount: 5000,
  total_amount: 105000,           // Final amount
  payment_method: "cash",         // cash | card | momo | zalopay
  status: "completed",            // pending | confirmed | preparing | completed | cancelled
  table_number: "A05",
  notes: "Không hành",
  created_at: ISODate("2026-02-15T12:30:00Z"),
  updated_at: ISODate("2026-02-15T13:00:00Z")
}
```

**Indexes:**
```javascript
db.orders.createIndex({ order_id: 1 }, { unique: true })
db.orders.createIndex({ restaurant_id: 1, order_date: 1 })
db.orders.createIndex({ customer_id: 1 })
db.orders.createIndex({ status: 1 })
db.orders.createIndex({ order_date: 1 })
```

### 3. payments

Payment transactions.

```javascript
{
  _id: ObjectId("..."),
  payment_id: "PAY20260215001",
  order_id: "ORD20260215001",
  restaurant_id: "REST001",
  amount: 105000,
  payment_method: "momo",
  status: "success",              // pending | success | failed | refunded
  transaction_id: "MOMO123456",   // External payment provider ID
  payment_date: ISODate("2026-02-15T13:00:00Z"),
  metadata: {
    provider_response: {...},
    fee: 2100                     // Payment provider fee
  },
  created_at: ISODate("2026-02-15T13:00:00Z"),
  updated_at: ISODate("2026-02-15T13:00:05Z")
}
```

**Indexes:**
```javascript
db.payments.createIndex({ payment_id: 1 }, { unique: true })
db.payments.createIndex({ order_id: 1 })
db.payments.createIndex({ restaurant_id: 1, payment_date: 1 })
db.payments.createIndex({ status: 1 })
```

### 4. feature_revenue_monthly

Engineered features for ML models.

```javascript
{
  _id: ObjectId("..."),
  restaurant_id: "REST001",
  month: "2026-02",               // YYYY-MM format
  
  // Base metrics
  revenue: 150000000,
  order_count: 1500,
  avg_order_value: 100000,
  unique_customers: 800,
  
  // Rolling averages
  rolling_avg_3m: 145000000,
  rolling_avg_6m: 140000000,
  rolling_avg_12m: 135000000,
  
  // Growth rates
  mom_growth: 5.5,                // Month-over-Month %
  yoy_growth: 12.3,               // Year-over-Year %
  
  // Seasonality
  is_holiday_season: false,
  seasonality_index: 1.05,
  
  // Additional features
  weekend_revenue_ratio: 0.35,
  peak_hours_revenue_ratio: 0.55,
  avg_items_per_order: 2.5,
  
  updated_at: ISODate("2026-02-15T02:00:00Z")
}
```

**Indexes:**
```javascript
db.feature_revenue_monthly.createIndex(
  { restaurant_id: 1, month: 1 },
  { unique: true }
)
db.feature_revenue_monthly.createIndex({ month: 1 })
db.feature_revenue_monthly.createIndex({ updated_at: 1 })
```

### 5. revenue_predictions

ML model predictions.

```javascript
{
  _id: ObjectId("..."),
  prediction_id: "PRED20260215_REST001_202603",
  restaurant_id: "REST001",
  month: "2026-03",               // Prediction for this month
  
  // Predictions
  predicted: 158000000,
  confidence: 0.85,               // 0-1 scale
  lower_bound: 148000000,         // 95% confidence interval
  upper_bound: 168000000,
  
  // Actual value (filled after month ends)
  actual: null,
  error: null,                    // Calculated when actual is filled
  
  // Model metadata
  model_name: "revenue_forecast_prophet_REST001",
  model_version: "v1.0.0",
  predicted_at: ISODate("2026-02-15T03:00:00Z"),
  
  // Additional predictions
  predicted_orders: 1550,
  predicted_avg_order_value: 102000,
  
  updated_at: ISODate("2026-02-15T03:00:00Z")
}
```

**Indexes:**
```javascript
db.revenue_predictions.createIndex(
  { restaurant_id: 1, month: 1 },
  { unique: true }
)
db.revenue_predictions.createIndex({ month: 1 })
db.revenue_predictions.createIndex({ predicted_at: 1 })
```

### 6. ml_models

Model metadata and versioning.

```javascript
{
  _id: ObjectId("..."),
  model_name: "revenue_forecast_prophet_REST001",
  model_type: "prophet",          // prophet | xgboost | lstm
  restaurant_id: "REST001",
  version: "v1.0.0",
  
  // Training metadata
  trained_at: ISODate("2026-02-15T02:00:00Z"),
  training_duration_seconds: 120,
  data_points: 24,                // Number of months used
  train_start_date: "2024-01-01",
  train_end_date: "2026-01-31",
  
  // Performance metrics
  metrics: {
    mape: 5.2,                    // Mean Absolute Percentage Error
    rmse: 8500000,                // Root Mean Square Error
    mae: 6200000,                 // Mean Absolute Error
    r2_score: 0.92                // R² Score
  },
  
  // Hyperparameters
  hyperparameters: {
    changepoint_prior_scale: 0.05,
    seasonality_mode: "multiplicative",
    yearly_seasonality: true
  },
  
  // Model file info
  model_path: "./models/revenue_forecast_prophet_REST001_v1.0.0.pkl",
  model_size_mb: 2.5,
  
  // Status
  is_production: true,
  status: "active",               // active | archived | failed
  
  created_at: ISODate("2026-02-15T02:00:00Z"),
  updated_at: ISODate("2026-02-15T02:00:00Z")
}
```

**Indexes:**
```javascript
db.ml_models.createIndex(
  { model_name: 1, version: 1 },
  { unique: true }
)
db.ml_models.createIndex({ restaurant_id: 1 })
db.ml_models.createIndex({ is_production: 1, status: 1 })
db.ml_models.createIndex({ trained_at: 1 })
```

### 7. etl_logs

ETL job execution logs.

```javascript
{
  _id: ObjectId("..."),
  job_id: "etl_revenue_20260215_020000",
  job_type: "revenue_features",   // revenue_features | payment_analysis
  status: "success",              // running | success | failed
  
  // Timing
  started_at: ISODate("2026-02-15T02:00:00Z"),
  completed_at: ISODate("2026-02-15T02:05:30Z"),
  duration_seconds: 330,
  
  // Statistics
  records_processed: 50000,
  records_updated: 36,
  records_failed: 0,
  
  // Error details
  errors: [],
  
  // Metadata
  metadata: {
    restaurants_processed: ["REST001", "REST002", "REST003"],
    date_range: {
      start: "2024-01-01",
      end: "2026-02-01"
    },
    trigger: "cron"               // cron | manual | api
  }
}
```

**Indexes:**
```javascript
db.etl_logs.createIndex({ job_id: 1 }, { unique: true })
db.etl_logs.createIndex({ started_at: 1 })
db.etl_logs.createIndex({ status: 1 })
```

### 8. api_logs (Future)

API request logs for monitoring.

```javascript
{
  _id: ObjectId("..."),
  request_id: "REQ20260215120000",
  method: "GET",
  path: "/api/v1/analytics/dashboard",
  query_params: {
    restaurant_id: "REST001"
  },
  status_code: 200,
  response_time_ms: 45,
  timestamp: ISODate("2026-02-15T12:00:00Z"),
  user_agent: "Mozilla/5.0...",
  ip_address: "192.168.1.100"
}
```

## 🔧 Setup Instructions

### 1. Initialize Database

```bash
# Connect to MongoDB
mongosh mongodb://localhost:27017

# Create database
use ai_analytics

# Load schemas
load('database/mongodb/schemas.js')

# Create indexes
load('database/mongodb/indexes.js')
```

### 2. Load Sample Data

```bash
# Generate and insert sample data
load('database/mongodb/sample_data.js')
```

### 3. Verify Setup

```javascript
// Check collections
show collections

// Count documents
db.restaurants.countDocuments()
db.orders.countDocuments()
db.payments.countDocuments()

// Sample query
db.orders.findOne()
```

## 📊 Useful Queries

### Revenue by Month

```javascript
db.orders.aggregate([
  {
    $match: {
      restaurant_id: "REST001",
      status: "completed"
    }
  },
  {
    $group: {
      _id: {
        $dateToString: { format: "%Y-%m", date: "$order_date" }
      },
      revenue: { $sum: "$total_amount" },
      orders: { $sum: 1 }
    }
  },
  {
    $sort: { _id: 1 }
  }
])
```

### Top Customers

```javascript
db.orders.aggregate([
  {
    $match: {
      restaurant_id: "REST001",
      status: "completed"
    }
  },
  {
    $group: {
      _id: "$customer_id",
      total_spent: { $sum: "$total_amount" },
      order_count: { $sum: 1 }
    }
  },
  {
    $sort: { total_spent: -1 }
  },
  {
    $limit: 10
  }
])
```

### Model Performance Over Time

```javascript
db.revenue_predictions.aggregate([
  {
    $match: {
      actual: { $ne: null }
    }
  },
  {
    $addFields: {
      error_pct: {
        $multiply: [
          { $divide: [
            { $abs: { $subtract: ["$predicted", "$actual"] }},
            "$actual"
          ]},
          100
        ]
      }
    }
  },
  {
    $group: {
      _id: null,
      avg_error: { $avg: "$error_pct" },
      max_error: { $max: "$error_pct" },
      min_error: { $min: "$error_pct" }
    }
  }
])
```

## 🔒 Security Best Practices

1. **Authentication**: Enable MongoDB authentication
2. **Network**: Bind to localhost in development
3. **Encryption**: Enable encryption at rest
4. **Backups**: Regular automated backups
5. **Monitoring**: Track slow queries

## 📈 Performance Tips

1. **Use indexes** for all query filters
2. **Limit projection** to needed fields only
3. **Use aggregation pipeline** for complex queries
4. **Enable profiling** to find slow queries
5. **Monitor index usage** and remove unused ones

## 🔄 Backup & Restore

### Backup

```bash
# Full backup
mongodump --db=ai_analytics --out=/backup/$(date +%Y%m%d)

# Specific collection
mongodump --db=ai_analytics --collection=orders --out=/backup
```

### Restore

```bash
# Full restore
mongorestore --db=ai_analytics /backup/20260215/ai_analytics

# Specific collection
mongorestore --db=ai_analytics --collection=orders /backup/ai_analytics/orders.bson
```
