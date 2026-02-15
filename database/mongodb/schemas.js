// MongoDB Collections Schema Definition
// This file documents the schema structure for all collections

// ==========================================
// RAW DATA COLLECTIONS
// ==========================================

// restaurants: Thông tin nhà hàng
{
  "_id": ObjectId,
  "restaurant_id": String,      // Unique identifier
  "name": String,
  "location": String,
  "created_at": ISODate,
  "status": String,             // "active", "inactive"
  "metadata": Object
}

// orders: Đơn hàng
{
  "_id": ObjectId,
  "order_id": String,           // Unique order identifier
  "restaurant_id": String,
  "total_price": Number,        // Giá trị đơn hàng
  "status": String,             // "pending", "completed", "cancelled"
  "items": [                    // Chi tiết món
    {
      "item_id": String,
      "quantity": Number,
      "price": Number
    }
  ],
  "created_at": ISODate,
  "completed_at": ISODate
}

// payments: Thanh toán
{
  "_id": ObjectId,
  "payment_id": String,
  "order_id": String,
  "amount": Number,
  "method": String,             // "cash", "card", "ewallet"
  "status": String,             // "pending", "success", "failed"
  "paid_at": ISODate,
  "created_at": ISODate
}

// ==========================================
// FEATURE STORE COLLECTIONS
// ==========================================

// feature_revenue_monthly: Features cho dự đoán doanh thu
{
  "_id": ObjectId,
  "restaurant_id": String,
  "month": String,              // "2026-01"
  "year": Number,
  "month_num": Number,
  
  // Raw metrics
  "revenue": Number,            // Doanh thu thực tế
  "order_count": Number,        // Số đơn hàng
  "avg_order_value": Number,    // Giá trị đơn hàng trung bình
  
  // Rolling features
  "rolling_avg_3m": Number,     // Trung bình 3 tháng
  "rolling_avg_6m": Number,     // Trung bình 6 tháng
  "rolling_avg_12m": Number,    // Trung bình 12 tháng
  
  // Growth features
  "mom_growth": Number,         // Month-over-month growth %
  "yoy_growth": Number,         // Year-over-year growth %
  
  // Seasonality features
  "season": String,             // "Q1", "Q2", "Q3", "Q4"
  "is_holiday": Boolean,        // Có phải tháng lễ không
  "day_of_week_avg": Number,    // Trung bình theo ngày trong tuần
  
  // Target
  "target": Number,             // Revenue target
  
  // Metadata
  "created_at": ISODate,
  "updated_at": ISODate,
  "version": String             // Feature version
}

// feature_order_patterns: Features cho phân tích đơn hàng
{
  "_id": ObjectId,
  "restaurant_id": String,
  "date": ISODate,
  "day_of_week": Number,        // 0-6
  "hour": Number,               // 0-23
  
  "order_count": Number,
  "total_revenue": Number,
  "avg_order_value": Number,
  "unique_customers": Number,
  
  "peak_hour": Number,          // Giờ cao điểm
  "conversion_rate": Number,
  
  "created_at": ISODate
}

// ==========================================
// PREDICTION COLLECTIONS
// ==========================================

// revenue_predictions: Kết quả dự đoán doanh thu
{
  "_id": ObjectId,
  "restaurant_id": String,
  "month": String,              // "2026-03"
  
  // Prediction values
  "predicted": Number,          // Giá trị dự đoán
  "lower_ci": Number,           // Confidence interval dưới
  "upper_ci": Number,           // Confidence interval trên
  
  // Actual (if available)
  "actual": Number,             // Doanh thu thực tế (null nếu chưa có)
  
  // Model info
  "model_name": String,         // "prophet", "xgboost", "lstm"
  "model_version": String,      // "v1.2.3"
  "confidence_score": Number,   // 0-1
  
  // Metrics
  "mape": Number,               // Mean Absolute Percentage Error
  "rmse": Number,               // Root Mean Square Error
  
  // Metadata
  "predicted_at": ISODate,
  "created_at": ISODate
}

// order_predictions: Dự đoán đơn hàng
{
  "_id": ObjectId,
  "restaurant_id": String,
  "date": ISODate,
  
  "predicted_orders": Number,
  "predicted_revenue": Number,
  "confidence_score": Number,
  
  "model_version": String,
  "predicted_at": ISODate
}

// ==========================================
// MODEL METADATA COLLECTION
// ==========================================

// ml_models: Thông tin về models
{
  "_id": ObjectId,
  "model_name": String,         // "revenue_forecast_prophet"
  "model_type": String,         // "prophet", "xgboost", "lstm"
  "version": String,            // "v1.2.3"
  
  "file_path": String,          // Path to model file
  "file_size": Number,          // Model size in bytes
  
  // Training info
  "trained_at": ISODate,
  "training_duration": Number,  // seconds
  "training_data_size": Number, // số records
  "training_date_range": {
    "start": ISODate,
    "end": ISODate
  },
  
  // Performance metrics
  "metrics": {
    "mape": Number,
    "rmse": Number,
    "mae": Number,
    "r2_score": Number
  },
  
  // Hyperparameters
  "hyperparameters": Object,
  
  // Status
  "status": String,             // "training", "active", "deprecated"
  "is_production": Boolean,
  
  "created_at": ISODate,
  "updated_at": ISODate
}

// ==========================================
// ETL JOB LOGS
// ==========================================

// etl_logs: Logs của ETL jobs
{
  "_id": ObjectId,
  "job_name": String,           // "feature_engineering_revenue"
  "job_type": String,           // "etl", "training", "prediction"
  "status": String,             // "running", "success", "failed"
  
  "started_at": ISODate,
  "completed_at": ISODate,
  "duration": Number,           // seconds
  
  "records_processed": Number,
  "records_failed": Number,
  
  "error_message": String,
  "stack_trace": String,
  
  "metadata": Object
}
