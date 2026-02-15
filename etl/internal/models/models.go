package models

import (
	"time"
)

// Order represents raw order data
type Order struct {
	OrderID      string    `bson:"order_id"`
	RestaurantID string    `bson:"restaurant_id"`
	TotalPrice   float64   `bson:"total_price"`
	Status       string    `bson:"status"`
	CreatedAt    time.Time `bson:"created_at"`
	CompletedAt  time.Time `bson:"completed_at"`
}

// Payment represents raw payment data
type Payment struct {
	PaymentID string    `bson:"payment_id"`
	OrderID   string    `bson:"order_id"`
	Amount    float64   `bson:"amount"`
	Method    string    `bson:"method"`
	Status    string    `bson:"status"`
	PaidAt    time.Time `bson:"paid_at"`
}

// FeatureRevenueMonthly represents aggregated monthly revenue features
type FeatureRevenueMonthly struct {
	RestaurantID string `bson:"restaurant_id"`
	Month        string `bson:"month"` // Format: "2026-01"
	Year         int    `bson:"year"`
	MonthNum     int    `bson:"month_num"`

	// Raw metrics
	Revenue       float64 `bson:"revenue"`
	OrderCount    int     `bson:"order_count"`
	AvgOrderValue float64 `bson:"avg_order_value"`

	// Rolling features
	RollingAvg3M  float64 `bson:"rolling_avg_3m"`
	RollingAvg6M  float64 `bson:"rolling_avg_6m"`
	RollingAvg12M float64 `bson:"rolling_avg_12m"`

	// Growth features
	MoMGrowth float64 `bson:"mom_growth"`
	YoYGrowth float64 `bson:"yoy_growth"`

	// Seasonality
	Season       string  `bson:"season"`
	IsHoliday    bool    `bson:"is_holiday"`
	DayOfWeekAvg float64 `bson:"day_of_week_avg"`

	// Target
	Target float64 `bson:"target"`

	// Metadata
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
	Version   string    `bson:"version"`
}

// ETLLog represents ETL job execution log
type ETLLog struct {
	JobName          string                 `bson:"job_name"`
	JobType          string                 `bson:"job_type"`
	Status           string                 `bson:"status"`
	StartedAt        time.Time              `bson:"started_at"`
	CompletedAt      time.Time              `bson:"completed_at"`
	Duration         int64                  `bson:"duration"` // seconds
	RecordsProcessed int                    `bson:"records_processed"`
	RecordsFailed    int                    `bson:"records_failed"`
	ErrorMessage     string                 `bson:"error_message,omitempty"`
	StackTrace       string                 `bson:"stack_trace,omitempty"`
	Metadata         map[string]interface{} `bson:"metadata,omitempty"`
}
