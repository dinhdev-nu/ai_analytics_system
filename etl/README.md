# ETL Worker Service

Go-based ETL service for feature engineering.

## 📁 Structure

```
etl/
├── cmd/
│   └── etl_worker/
│       └── main.go               # Entry point with cron
├── internal/
│   ├── config/
│   │   └── config.go             # Configuration
│   ├── database/
│   │   └── mongodb.go            # MongoDB connection
│   ├── models/
│   │   └── models.go             # Data models
│   └── etl/
│       └── revenue_features.go   # Feature engineering logic
├── go.mod
└── Dockerfile
```

## 🚀 Quick Start

### Development

```bash
# Install dependencies
go mod download

# Run once
go run cmd/etl_worker/main.go

# Run with cron
ETL_MODE=cron go run cmd/etl_worker/main.go
```

### Production

```bash
# Build
go build -o bin/etl_worker cmd/etl_worker/main.go

# Run in cron mode
ETL_MODE=cron ./bin/etl_worker
```

## ⚙️ Configuration

### Environment Variables

```env
# MongoDB
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=ai_analytics

# ETL Settings
ETL_MODE=cron              # cron or once
ETL_SCHEDULE=0 2 * * *     # Daily at 2 AM
ETL_BATCH_SIZE=1000
ETL_WORKERS=4

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

### Cron Schedule Format

```
┌───────── minute (0 - 59)
│ ┌─────── hour (0 - 23)
│ │ ┌───── day of month (1 - 31)
│ │ │ ┌─── month (1 - 12)
│ │ │ │ ┌─ day of week (0 - 6, Sunday=0)
│ │ │ │ │
* * * * *

Examples:
0 2 * * *     # Daily at 2:00 AM
0 */6 * * *   # Every 6 hours
0 0 * * 0     # Weekly on Sunday
0 0 1 * *     # Monthly on the 1st
```

## 🔧 Feature Engineering

### Monthly Revenue Features

Aggregated from raw orders data:

```go
type FeatureRevenueMonthly struct {
    RestaurantID      string    `bson:"restaurant_id"`
    Month             string    `bson:"month"` // YYYY-MM
    Revenue           float64   `bson:"revenue"`
    OrderCount        int       `bson:"order_count"`
    AvgOrderValue     float64   `bson:"avg_order_value"`
    UniqueCustomers   int       `bson:"unique_customers"`
    
    // Rolling averages
    RollingAvg3M      float64   `bson:"rolling_avg_3m"`
    RollingAvg6M      float64   `bson:"rolling_avg_6m"`
    RollingAvg12M     float64   `bson:"rolling_avg_12m"`
    
    // Growth rates
    MoMGrowth         float64   `bson:"mom_growth"` // Month-over-Month
    YoYGrowth         float64   `bson:"yoy_growth"` // Year-over-Year
    
    // Seasonality
    IsHolidaySeason   bool      `bson:"is_holiday_season"`
    SeasonalityIndex  float64   `bson:"seasonality_index"`
    
    UpdatedAt         time.Time `bson:"updated_at"`
}
```

### Calculation Logic

#### 1. Base Aggregation

```go
// Aggregate from orders
pipeline := []bson.M{
    {"$match": bson.M{
        "restaurant_id": restaurantID,
        "status": "completed",
        "order_date": bson.M{
            "$gte": startDate,
            "$lt": endDate,
        },
    }},
    {"$group": bson.M{
        "_id": "$restaurant_id",
        "revenue": bson.M{"$sum": "$total_amount"},
        "order_count": bson.M{"$sum": 1},
        "unique_customers": bson.M{"$addToSet": "$customer_id"},
    }},
}
```

#### 2. Rolling Averages

```go
func calculateRollingAverage(values []float64, window int) float64 {
    if len(values) < window {
        window = len(values)
    }
    
    sum := 0.0
    for i := len(values) - window; i < len(values); i++ {
        sum += values[i]
    }
    
    return sum / float64(window)
}
```

#### 3. Growth Rates

```go
// Month-over-Month Growth
func calculateMoMGrowth(current, previous float64) float64 {
    if previous == 0 {
        return 0
    }
    return ((current - previous) / previous) * 100
}

// Year-over-Year Growth
func calculateYoYGrowth(current, yearAgo float64) float64 {
    if yearAgo == 0 {
        return 0
    }
    return ((current - yearAgo) / yearAgo) * 100
}
```

#### 4. Seasonality Detection

```go
func detectSeasonality(month int) (bool, float64) {
    // Holiday months in Vietnam
    holidayMonths := map[int]bool{
        1: true,  // Lunar New Year
        4: true,  // Reunification Day
        9: true,  // National Day
        12: true, // Christmas
    }
    
    isHoliday := holidayMonths[month]
    seasonalityIndex := calculateSeasonalityIndex(month)
    
    return isHoliday, seasonalityIndex
}
```

## 📊 Data Flow

```
┌─────────────┐
│ Raw Orders  │
│ Collection  │
└──────┬──────┘
       │
       │ ETL reads
       ▼
┌──────────────┐
│ Aggregation  │
│ & Transform  │
└──────┬───────┘
       │
       │ Calculate features
       ▼
┌────────────────────┐
│ feature_revenue_   │
│ monthly Collection │
└────────────────────┘
       │
       │ ML reads
       ▼
┌────────────────┐
│ Model Training │
└────────────────┘
```

## 🔄 Execution Modes

### Single Run (Testing)

```bash
ETL_MODE=once go run cmd/etl_worker/main.go
```

### Cron Mode (Production)

```bash
ETL_MODE=cron ETL_SCHEDULE="0 2 * * *" go run cmd/etl_worker/main.go
```

### Manual Trigger via API (Future)

```bash
curl -X POST http://localhost:8080/api/v1/etl/trigger \
  -H "Content-Type: application/json" \
  -d '{"restaurant_id": "REST001", "start_date": "2026-01-01"}'
```

## 📝 Logging

ETL logs are stored in MongoDB:

```javascript
{
  job_id: "etl_revenue_20260215_020000",
  job_type: "revenue_features",
  status: "success",
  started_at: ISODate("2026-02-15T02:00:00Z"),
  completed_at: ISODate("2026-02-15T02:05:30Z"),
  duration_seconds: 330,
  records_processed: 50000,
  records_updated: 36,
  errors: [],
  metadata: {
    restaurants_processed: ["REST001", "REST002", "REST003"],
    date_range: {
      start: "2024-01-01",
      end: "2026-02-01"
    }
  }
}
```

### Query Logs

```javascript
// Get recent jobs
db.etl_logs.find().sort({started_at: -1}).limit(10)

// Get failed jobs
db.etl_logs.find({status: "failed"})

// Get job statistics
db.etl_logs.aggregate([
  {$group: {
    _id: "$status",
    count: {$sum: 1},
    avg_duration: {$avg: "$duration_seconds"}
  }}
])
```

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Test specific package
go test ./internal/etl/

# With coverage
go test -cover ./...

# Integration tests
go test -tags=integration ./...
```

### Test Example

```go
func TestCalculateRollingAverage(t *testing.T) {
    values := []float64{100, 200, 300, 400, 500}
    
    avg3 := calculateRollingAverage(values, 3)
    expected := (300 + 400 + 500) / 3.0
    
    if avg3 != expected {
        t.Errorf("Expected %f, got %f", expected, avg3)
    }
}
```

## 🐛 Debugging

### Enable Debug Logging

```bash
LOG_LEVEL=debug go run cmd/etl_worker/main.go
```

### Dry Run Mode

```go
// Add flag to skip database writes
if os.Getenv("ETL_DRY_RUN") == "true" {
    log.Println("DRY RUN: Would insert", feature)
    return nil
}
```

## 📊 Performance Optimization

### Batch Processing

```go
// Process in batches
batchSize := 1000
for i := 0; i < len(orders); i += batchSize {
    end := i + batchSize
    if end > len(orders) {
        end = len(orders)
    }
    
    batch := orders[i:end]
    processBatch(batch)
}
```

### Parallel Workers

```go
// Use worker pool
workerCount := 4
jobs := make(chan Restaurant, 100)
results := make(chan Result, 100)

// Start workers
for w := 0; w < workerCount; w++ {
    go worker(jobs, results)
}

// Send jobs
for _, restaurant := range restaurants {
    jobs <- restaurant
}
close(jobs)
```

### Indexing

Ensure proper indexes exist:

```javascript
db.orders.createIndex({ restaurant_id: 1, order_date: 1 })
db.orders.createIndex({ status: 1 })
db.feature_revenue_monthly.createIndex({ restaurant_id: 1, month: 1 }, { unique: true })
```

## 🚨 Error Handling

### Retry Logic

```go
func withRetry(fn func() error, maxRetries int) error {
    var err error
    for i := 0; i < maxRetries; i++ {
        err = fn()
        if err == nil {
            return nil
        }
        
        log.Printf("Retry %d/%d: %v", i+1, maxRetries, err)
        time.Sleep(time.Second * time.Duration(i+1))
    }
    return err
}
```

### Alerting

```go
func sendAlert(err error) {
    // Send to monitoring system
    // Slack, PagerDuty, etc.
}
```

## 📚 Adding New Features

### 1. Define Model

```go
// internal/models/models.go
type NewFeature struct {
    RestaurantID string  `bson:"restaurant_id"`
    Value        float64 `bson:"value"`
    UpdatedAt    time.Time `bson:"updated_at"`
}
```

### 2. Implement Calculation

```go
// internal/etl/new_feature.go
func CalculateNewFeatures(db *mongo.Database) error {
    // Implementation
}
```

### 3. Add to ETL Job

```go
// cmd/etl_worker/main.go
func runETLJob() {
    calculateRevenueFeatures()
    calculateNewFeatures() // Add here
}
```

## 🔐 Security

- Use environment variables for credentials
- Never commit secrets to Git
- Validate all inputs
- Sanitize MongoDB queries
- Use connection pooling with limits
