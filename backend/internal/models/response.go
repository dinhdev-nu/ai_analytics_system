package models

import "time"

// RevenueAnalytics response
type RevenueAnalytics struct {
	RestaurantID string        `json:"restaurant_id"`
	Period       string        `json:"period"` // "monthly", "yearly"
	Data         []RevenueData `json:"data"`
}

type RevenueData struct {
	Month     string  `json:"month"` // "2026-01"
	Actual    float64 `json:"actual"`
	Predicted float64 `json:"predicted,omitempty"`
	Target    float64 `json:"target,omitempty"`
	LowerCI   float64 `json:"lower_ci,omitempty"`
	UpperCI   float64 `json:"upper_ci,omitempty"`
}

// ForecastResponse for forecast-specific endpoint
type ForecastResponse struct {
	RestaurantID   string            `json:"restaurant_id"`
	Timestamps     []string          `json:"timestamps"`
	Actual         []float64         `json:"actual"`
	Predicted      []float64         `json:"predicted"`
	Target         []float64         `json:"target"`
	ConfidenceData []ConfidenceRange `json:"confidence_data"`
	ModelInfo      ModelInfo         `json:"model_info"`
}

type ConfidenceRange struct {
	Month string  `json:"month"`
	Lower float64 `json:"lower"`
	Upper float64 `json:"upper"`
}

type ModelInfo struct {
	ModelName    string    `json:"model_name"`
	ModelVersion string    `json:"model_version"`
	LastUpdated  time.Time `json:"last_updated"`
	Metrics      Metrics   `json:"metrics"`
}

type Metrics struct {
	MAPE    float64 `json:"mape"`
	RMSE    float64 `json:"rmse"`
	MAE     float64 `json:"mae"`
	R2Score float64 `json:"r2_score"`
}

// Dashboard summary response
type DashboardResponse struct {
	RestaurantID string           `json:"restaurant_id"`
	Summary      Summary          `json:"summary"`
	RevenueChart RevenueChartData `json:"revenue_chart"`
	OrdersChart  OrdersChartData  `json:"orders_chart"`
	Insights     []Insight        `json:"insights"`
}

type Summary struct {
	CurrentMonthRevenue  float64 `json:"current_month_revenue"`
	PreviousMonthRevenue float64 `json:"previous_month_revenue"`
	MonthOverMonthGrowth float64 `json:"month_over_month_growth"`
	YearOverYearGrowth   float64 `json:"year_over_year_growth"`
	TotalOrders          int     `json:"total_orders"`
	AvgOrderValue        float64 `json:"avg_order_value"`
	ForecastNextMonth    float64 `json:"forecast_next_month"`
	ForecastConfidence   float64 `json:"forecast_confidence"`
}

type RevenueChartData struct {
	Labels    []string  `json:"labels"`
	Actual    []float64 `json:"actual"`
	Predicted []float64 `json:"predicted"`
	Target    []float64 `json:"target"`
}

type OrdersChartData struct {
	Labels      []string `json:"labels"`
	OrderCounts []int    `json:"order_counts"`
}

type Insight struct {
	Type        string `json:"type"` // "warning", "success", "info"
	Title       string `json:"title"`
	Description string `json:"description"`
	Value       string `json:"value,omitempty"`
}

// Error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// Success response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
