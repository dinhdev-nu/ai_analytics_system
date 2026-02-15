package services

import (
	"context"
	"fmt"
	"time"

	"ai-analytics/backend/internal/database"
	"ai-analytics/backend/internal/models"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

type AnalyticsService struct {
	db *database.MongoDB
}

func NewAnalyticsService(db *database.MongoDB) *AnalyticsService {
	return &AnalyticsService{db: db}
}

// GetRevenueForecast returns revenue data with predictions
func (s *AnalyticsService) GetRevenueForecast(ctx context.Context, restaurantID string, months int) (*models.ForecastResponse, error) {
	log.Infof("Fetching revenue forecast for restaurant: %s, months: %d", restaurantID, months)

	// Get predictions
	predictions, err := s.getPredictions(ctx, restaurantID, months)
	if err != nil {
		return nil, err
	}

	// Get actual data
	actuals, err := s.getActualRevenue(ctx, restaurantID, months)
	if err != nil {
		return nil, err
	}

	// Get features (for targets)
	features, err := s.getFeatures(ctx, restaurantID, months)
	if err != nil {
		return nil, err
	}

	// Merge data
	response := s.mergeForecastData(predictions, actuals, features)
	response.RestaurantID = restaurantID

	// Get model info
	modelInfo, err := s.getModelInfo(ctx, restaurantID)
	if err == nil {
		response.ModelInfo = modelInfo
	}

	return response, nil
}

// GetDashboard returns complete dashboard data
func (s *AnalyticsService) GetDashboard(ctx context.Context, restaurantID string) (*models.DashboardResponse, error) {
	log.Infof("Fetching dashboard for restaurant: %s", restaurantID)

	response := &models.DashboardResponse{
		RestaurantID: restaurantID,
	}

	// Get summary
	summary, err := s.getSummary(ctx, restaurantID)
	if err != nil {
		log.Errorf("Failed to get summary: %v", err)
	} else {
		response.Summary = summary
	}

	// Get revenue chart data (last 12 months + 6 months forecast)
	forecast, err := s.GetRevenueForecast(ctx, restaurantID, 18)
	if err != nil {
		log.Errorf("Failed to get forecast: %v", err)
	} else {
		response.RevenueChart = models.RevenueChartData{
			Labels:    forecast.Timestamps,
			Actual:    forecast.Actual,
			Predicted: forecast.Predicted,
			Target:    forecast.Target,
		}
	}

	// Get orders chart
	ordersChart, err := s.getOrdersChart(ctx, restaurantID, 12)
	if err != nil {
		log.Errorf("Failed to get orders chart: %v", err)
	} else {
		response.OrdersChart = ordersChart
	}

	// Generate insights
	response.Insights = s.generateInsights(summary, forecast)

	return response, nil
}

func (s *AnalyticsService) getPredictions(ctx context.Context, restaurantID string, months int) ([]bson.M, error) {
	collection := s.db.GetCollection("revenue_predictions")

	// Get latest predictions
	cursor, err := collection.Find(ctx, bson.M{
		"restaurant_id": restaurantID,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var predictions []bson.M
	if err := cursor.All(ctx, &predictions); err != nil {
		return nil, err
	}

	return predictions, nil
}

func (s *AnalyticsService) getActualRevenue(ctx context.Context, restaurantID string, months int) ([]bson.M, error) {
	collection := s.db.GetCollection("feature_revenue_monthly")

	// Get last N months
	cursor, err := collection.Find(ctx, bson.M{
		"restaurant_id": restaurantID,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var actuals []bson.M
	if err := cursor.All(ctx, &actuals); err != nil {
		return nil, err
	}

	return actuals, nil
}

func (s *AnalyticsService) getFeatures(ctx context.Context, restaurantID string, months int) ([]bson.M, error) {
	collection := s.db.GetCollection("feature_revenue_monthly")

	cursor, err := collection.Find(ctx, bson.M{
		"restaurant_id": restaurantID,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var features []bson.M
	if err := cursor.All(ctx, &features); err != nil {
		return nil, err
	}

	return features, nil
}

func (s *AnalyticsService) mergeForecastData(predictions, actuals, features []bson.M) *models.ForecastResponse {
	response := &models.ForecastResponse{
		Timestamps:     []string{},
		Actual:         []float64{},
		Predicted:      []float64{},
		Target:         []float64{},
		ConfidenceData: []models.ConfidenceRange{},
	}

	// Create maps for quick lookup
	actualsMap := make(map[string]float64)
	for _, a := range actuals {
		if month, ok := a["month"].(string); ok {
			if revenue, ok := a["revenue"].(float64); ok {
				actualsMap[month] = revenue
			}
		}
	}

	targetsMap := make(map[string]float64)
	for _, f := range features {
		if month, ok := f["month"].(string); ok {
			if target, ok := f["target"].(float64); ok {
				targetsMap[month] = target
			}
		}
	}

	// Merge predictions with actuals
	for _, p := range predictions {
		month, _ := p["month"].(string)
		predicted, _ := p["predicted"].(float64)
		lowerCI, _ := p["lower_ci"].(float64)
		upperCI, _ := p["upper_ci"].(float64)

		response.Timestamps = append(response.Timestamps, month)
		response.Predicted = append(response.Predicted, predicted)

		// Add actual if available
		if actual, exists := actualsMap[month]; exists {
			response.Actual = append(response.Actual, actual)
		} else {
			response.Actual = append(response.Actual, 0)
		}

		// Add target if available
		if target, exists := targetsMap[month]; exists {
			response.Target = append(response.Target, target)
		} else {
			response.Target = append(response.Target, 0)
		}

		// Add confidence interval
		response.ConfidenceData = append(response.ConfidenceData, models.ConfidenceRange{
			Month: month,
			Lower: lowerCI,
			Upper: upperCI,
		})
	}

	return response
}

func (s *AnalyticsService) getModelInfo(ctx context.Context, restaurantID string) (models.ModelInfo, error) {
	collection := s.db.GetCollection("ml_models")

	var model bson.M
	err := collection.FindOne(ctx, bson.M{
		"model_name":    "revenue_forecast_prophet_" + restaurantID,
		"is_production": true,
	}).Decode(&model)

	if err != nil {
		return models.ModelInfo{}, err
	}

	modelInfo := models.ModelInfo{
		ModelName:    model["model_name"].(string),
		ModelVersion: model["version"].(string),
		LastUpdated:  model["trained_at"].(time.Time),
	}

	// Extract metrics
	if metrics, ok := model["metrics"].(bson.M); ok {
		modelInfo.Metrics = models.Metrics{
			MAPE:    metrics["mape"].(float64),
			RMSE:    metrics["rmse"].(float64),
			MAE:     metrics["mae"].(float64),
			R2Score: metrics["r2_score"].(float64),
		}
	}

	return modelInfo, nil
}

func (s *AnalyticsService) getSummary(ctx context.Context, restaurantID string) (models.Summary, error) {
	collection := s.db.GetCollection("feature_revenue_monthly")

	// Get last 2 months data
	now := time.Now()
	currentMonth := now.Format("2006-01")
	previousMonth := now.AddDate(0, -1, 0).Format("2006-01")
	lastYearMonth := now.AddDate(-1, 0, 0).Format("2006-01")

	var current, previous, lastYear bson.M

	collection.FindOne(ctx, bson.M{
		"restaurant_id": restaurantID,
		"month":         currentMonth,
	}).Decode(&current)

	collection.FindOne(ctx, bson.M{
		"restaurant_id": restaurantID,
		"month":         previousMonth,
	}).Decode(&previous)

	collection.FindOne(ctx, bson.M{
		"restaurant_id": restaurantID,
		"month":         lastYearMonth,
	}).Decode(&lastYear)

	summary := models.Summary{}

	if revenue, ok := current["revenue"].(float64); ok {
		summary.CurrentMonthRevenue = revenue
	}

	if revenue, ok := previous["revenue"].(float64); ok {
		summary.PreviousMonthRevenue = revenue
	}

	// Calculate growth
	if summary.PreviousMonthRevenue > 0 {
		summary.MonthOverMonthGrowth = ((summary.CurrentMonthRevenue - summary.PreviousMonthRevenue) / summary.PreviousMonthRevenue) * 100
	}

	if lastYearRevenue, ok := lastYear["revenue"].(float64); ok && lastYearRevenue > 0 {
		summary.YearOverYearGrowth = ((summary.CurrentMonthRevenue - lastYearRevenue) / lastYearRevenue) * 100
	}

	if orderCount, ok := current["order_count"].(int32); ok {
		summary.TotalOrders = int(orderCount)
	}

	if avgOrderValue, ok := current["avg_order_value"].(float64); ok {
		summary.AvgOrderValue = avgOrderValue
	}

	// Get next month forecast
	predCollection := s.db.GetCollection("revenue_predictions")
	nextMonth := now.AddDate(0, 1, 0).Format("2006-01")

	var prediction bson.M
	predCollection.FindOne(ctx, bson.M{
		"restaurant_id": restaurantID,
		"month":         nextMonth,
	}).Decode(&prediction)

	if forecast, ok := prediction["predicted"].(float64); ok {
		summary.ForecastNextMonth = forecast
	}

	if confidence, ok := prediction["confidence_score"].(float64); ok {
		summary.ForecastConfidence = confidence
	}

	return summary, nil
}

func (s *AnalyticsService) getOrdersChart(ctx context.Context, restaurantID string, months int) (models.OrdersChartData, error) {
	collection := s.db.GetCollection("feature_revenue_monthly")

	cursor, err := collection.Find(ctx, bson.M{
		"restaurant_id": restaurantID,
	})
	if err != nil {
		return models.OrdersChartData{}, err
	}
	defer cursor.Close(ctx)

	var data []bson.M
	cursor.All(ctx, &data)

	ordersChart := models.OrdersChartData{
		Labels:      []string{},
		OrderCounts: []int{},
	}

	for _, d := range data {
		if month, ok := d["month"].(string); ok {
			ordersChart.Labels = append(ordersChart.Labels, month)
		}
		if orderCount, ok := d["order_count"].(int32); ok {
			ordersChart.OrderCounts = append(ordersChart.OrderCounts, int(orderCount))
		}
	}

	return ordersChart, nil
}

func (s *AnalyticsService) generateInsights(summary models.Summary, forecast *models.ForecastResponse) []models.Insight {
	insights := []models.Insight{}

	// Growth insight
	if summary.MonthOverMonthGrowth > 10 {
		insights = append(insights, models.Insight{
			Type:        "success",
			Title:       "Tăng trưởng mạnh",
			Description: "Doanh thu tháng này tăng mạnh so với tháng trước",
			Value:       formatPercent(summary.MonthOverMonthGrowth),
		})
	} else if summary.MonthOverMonthGrowth < -5 {
		insights = append(insights, models.Insight{
			Type:        "warning",
			Title:       "Doanh thu giảm",
			Description: "Doanh thu tháng này giảm so với tháng trước",
			Value:       formatPercent(summary.MonthOverMonthGrowth),
		})
	}

	// Forecast insight
	if summary.ForecastNextMonth > summary.CurrentMonthRevenue {
		insights = append(insights, models.Insight{
			Type:        "info",
			Title:       "Dự báo tích cực",
			Description: "Doanh thu tháng tới dự kiến tăng",
			Value:       formatMoney(summary.ForecastNextMonth),
		})
	}

	return insights
}

func formatPercent(value float64) string {
	return fmt.Sprintf("%.1f%%", value)
}

func formatMoney(value float64) string {
	return fmt.Sprintf("%.0f VND", value)
}
