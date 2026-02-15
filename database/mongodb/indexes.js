// MongoDB Indexes
// Run this file to create indexes for optimal query performance

// ==========================================
// RAW DATA INDEXES
// ==========================================

// restaurants
db.restaurants.createIndex({ "restaurant_id": 1 }, { unique: true });
db.restaurants.createIndex({ "status": 1 });
db.restaurants.createIndex({ "created_at": -1 });

// orders
db.orders.createIndex({ "order_id": 1 }, { unique: true });
db.orders.createIndex({ "restaurant_id": 1, "created_at": -1 });
db.orders.createIndex({ "status": 1 });
db.orders.createIndex({ "created_at": -1 });
db.orders.createIndex({ "completed_at": -1 });

// Compound index for analytics queries
db.orders.createIndex({ 
  "restaurant_id": 1, 
  "status": 1, 
  "created_at": -1 
});

// payments
db.payments.createIndex({ "payment_id": 1 }, { unique: true });
db.payments.createIndex({ "order_id": 1 });
db.payments.createIndex({ "restaurant_id": 1, "paid_at": -1 });
db.payments.createIndex({ "status": 1 });
db.payments.createIndex({ "paid_at": -1 });

// ==========================================
// FEATURE STORE INDEXES
// ==========================================

// feature_revenue_monthly
db.feature_revenue_monthly.createIndex({ 
  "restaurant_id": 1, 
  "month": 1 
}, { unique: true });
db.feature_revenue_monthly.createIndex({ "month": 1 });
db.feature_revenue_monthly.createIndex({ "year": 1, "month_num": 1 });
db.feature_revenue_monthly.createIndex({ "created_at": -1 });

// Compound for time-series queries
db.feature_revenue_monthly.createIndex({
  "restaurant_id": 1,
  "year": 1,
  "month_num": 1
});

// feature_order_patterns
db.feature_order_patterns.createIndex({ 
  "restaurant_id": 1, 
  "date": 1 
}, { unique: true });
db.feature_order_patterns.createIndex({ "date": -1 });
db.feature_order_patterns.createIndex({ "day_of_week": 1, "hour": 1 });

// ==========================================
// PREDICTION INDEXES
// ==========================================

// revenue_predictions
db.revenue_predictions.createIndex({ 
  "restaurant_id": 1, 
  "month": 1,
  "model_version": 1
}, { unique: true });
db.revenue_predictions.createIndex({ "month": 1 });
db.revenue_predictions.createIndex({ "model_version": 1 });
db.revenue_predictions.createIndex({ "predicted_at": -1 });

// Compound for API queries
db.revenue_predictions.createIndex({
  "restaurant_id": 1,
  "month": 1,
  "model_version": 1
});

// order_predictions
db.order_predictions.createIndex({ 
  "restaurant_id": 1, 
  "date": 1 
}, { unique: true });
db.order_predictions.createIndex({ "date": -1 });
db.order_predictions.createIndex({ "predicted_at": -1 });

// ==========================================
// MODEL METADATA INDEXES
// ==========================================

// ml_models
db.ml_models.createIndex({ "model_name": 1, "version": 1 }, { unique: true });
db.ml_models.createIndex({ "status": 1 });
db.ml_models.createIndex({ "is_production": 1 });
db.ml_models.createIndex({ "trained_at": -1 });

// ==========================================
// ETL LOG INDEXES
// ==========================================

// etl_logs
db.etl_logs.createIndex({ "job_name": 1, "started_at": -1 });
db.etl_logs.createIndex({ "status": 1 });
db.etl_logs.createIndex({ "started_at": -1 });
db.etl_logs.createIndex({ "completed_at": -1 });

// TTL index - auto delete logs older than 90 days
db.etl_logs.createIndex({ "created_at": 1 }, { expireAfterSeconds: 7776000 });

print("All indexes created successfully!");
