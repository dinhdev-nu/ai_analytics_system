package etl

import (
	"context"
	"fmt"
	"time"

	"ai-analytics/etl/internal/database"
	"ai-analytics/etl/internal/models"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RevenueFeatureETL struct {
	db *database.MongoDB
}

func NewRevenueFeatureETL(db *database.MongoDB) *RevenueFeatureETL {
	return &RevenueFeatureETL{db: db}
}

// Run executes the complete ETL pipeline for revenue features
func (e *RevenueFeatureETL) Run(ctx context.Context) error {
	startTime := time.Now()
	log.Info("Starting Revenue Feature ETL job")

	// Create ETL log
	etlLog := &models.ETLLog{
		JobName:   "revenue_feature_engineering",
		JobType:   "etl",
		Status:    "running",
		StartedAt: startTime,
	}

	// Get list of restaurants
	restaurants, err := e.getRestaurants(ctx)
	if err != nil {
		e.saveETLLog(ctx, etlLog, "failed", err)
		return err
	}

	recordsProcessed := 0
	recordsFailed := 0

	// Process each restaurant
	for _, restaurantID := range restaurants {
		log.Infof("Processing restaurant: %s", restaurantID)

		if err := e.processRestaurant(ctx, restaurantID); err != nil {
			log.Errorf("Failed to process restaurant %s: %v", restaurantID, err)
			recordsFailed++
			continue
		}
		recordsProcessed++
	}

	etlLog.RecordsProcessed = recordsProcessed
	etlLog.RecordsFailed = recordsFailed
	etlLog.CompletedAt = time.Now()
	etlLog.Duration = int64(time.Since(startTime).Seconds())
	etlLog.Status = "success"

	e.saveETLLog(ctx, etlLog, "success", nil)

	log.Infof("Revenue Feature ETL completed. Processed: %d, Failed: %d",
		recordsProcessed, recordsFailed)
	return nil
}

func (e *RevenueFeatureETL) getRestaurants(ctx context.Context) ([]string, error) {
	collection := e.db.GetCollection("restaurants")

	cursor, err := collection.Find(ctx, bson.M{"status": "active"})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var restaurants []string
	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			continue
		}
		if id, ok := doc["restaurant_id"].(string); ok {
			restaurants = append(restaurants, id)
		}
	}

	return restaurants, nil
}

func (e *RevenueFeatureETL) processRestaurant(ctx context.Context, restaurantID string) error {
	// Get monthly revenue data for the last 12 months
	now := time.Now()
	startDate := now.AddDate(0, -12, 0)

	for currentDate := startDate; currentDate.Before(now); currentDate = currentDate.AddDate(0, 1, 0) {
		monthStr := currentDate.Format("2006-01")

		feature, err := e.calculateMonthlyFeatures(ctx, restaurantID, currentDate)
		if err != nil {
			log.Errorf("Failed to calculate features for %s - %s: %v",
				restaurantID, monthStr, err)
			continue
		}

		// Upsert feature to database
		if err := e.saveFeature(ctx, feature); err != nil {
			log.Errorf("Failed to save feature for %s - %s: %v",
				restaurantID, monthStr, err)
			continue
		}

		log.Debugf("Saved features for %s - %s", restaurantID, monthStr)
	}

	return nil
}

func (e *RevenueFeatureETL) calculateMonthlyFeatures(ctx context.Context, restaurantID string, month time.Time) (*models.FeatureRevenueMonthly, error) {
	monthStr := month.Format("2006-01")
	startOfMonth := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	// Calculate raw metrics
	revenue, orderCount, err := e.calculateRawMetrics(ctx, restaurantID, startOfMonth, endOfMonth)
	if err != nil {
		return nil, err
	}

	avgOrderValue := float64(0)
	if orderCount > 0 {
		avgOrderValue = revenue / float64(orderCount)
	}

	// Calculate rolling averages
	rollingAvg3M := e.calculateRollingAverage(ctx, restaurantID, month, 3)
	rollingAvg6M := e.calculateRollingAverage(ctx, restaurantID, month, 6)
	rollingAvg12M := e.calculateRollingAverage(ctx, restaurantID, month, 12)

	// Calculate growth rates
	momGrowth := e.calculateMoMGrowth(ctx, restaurantID, month)
	yoyGrowth := e.calculateYoYGrowth(ctx, restaurantID, month)

	// Determine season
	season := e.getSeason(month.Month())
	isHoliday := e.isHolidayMonth(month.Month())

	// Calculate target (simple: 10% growth from current)
	target := revenue * 1.1

	feature := &models.FeatureRevenueMonthly{
		RestaurantID:  restaurantID,
		Month:         monthStr,
		Year:          month.Year(),
		MonthNum:      int(month.Month()),
		Revenue:       revenue,
		OrderCount:    orderCount,
		AvgOrderValue: avgOrderValue,
		RollingAvg3M:  rollingAvg3M,
		RollingAvg6M:  rollingAvg6M,
		RollingAvg12M: rollingAvg12M,
		MoMGrowth:     momGrowth,
		YoYGrowth:     yoyGrowth,
		Season:        season,
		IsHoliday:     isHoliday,
		Target:        target,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Version:       "v1.0",
	}

	return feature, nil
}

func (e *RevenueFeatureETL) calculateRawMetrics(ctx context.Context, restaurantID string, startDate, endDate time.Time) (float64, int, error) {
	collection := e.db.GetCollection("orders")

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{
			{Key: "restaurant_id", Value: restaurantID},
			{Key: "status", Value: "completed"},
			{Key: "created_at", Value: bson.D{
				{Key: "$gte", Value: startDate},
				{Key: "$lt", Value: endDate},
			}},
		}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "total_revenue", Value: bson.D{{Key: "$sum", Value: "$total_price"}}},
			{Key: "order_count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, 0, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		var result struct {
			TotalRevenue float64 `bson:"total_revenue"`
			OrderCount   int     `bson:"order_count"`
		}
		if err := cursor.Decode(&result); err != nil {
			return 0, 0, err
		}
		return result.TotalRevenue, result.OrderCount, nil
	}

	return 0, 0, nil
}

func (e *RevenueFeatureETL) calculateRollingAverage(ctx context.Context, restaurantID string, month time.Time, months int) float64 {
	collection := e.db.GetCollection("feature_revenue_monthly")

	// Get previous N months
	startDate := month.AddDate(0, -months, 0)
	endDate := month

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{
			{Key: "restaurant_id", Value: restaurantID},
			{Key: "year", Value: bson.D{
				{Key: "$gte", Value: startDate.Year()},
				{Key: "$lte", Value: endDate.Year()},
			}},
		}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "avg_revenue", Value: bson.D{{Key: "$avg", Value: "$revenue"}}},
		}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		var result struct {
			AvgRevenue float64 `bson:"avg_revenue"`
		}
		if err := cursor.Decode(&result); err != nil {
			return 0
		}
		return result.AvgRevenue
	}

	return 0
}

func (e *RevenueFeatureETL) calculateMoMGrowth(ctx context.Context, restaurantID string, month time.Time) float64 {
	currentRevenue := e.getRevenueForMonth(ctx, restaurantID, month)
	previousRevenue := e.getRevenueForMonth(ctx, restaurantID, month.AddDate(0, -1, 0))

	if previousRevenue == 0 {
		return 0
	}

	return ((currentRevenue - previousRevenue) / previousRevenue) * 100
}

func (e *RevenueFeatureETL) calculateYoYGrowth(ctx context.Context, restaurantID string, month time.Time) float64 {
	currentRevenue := e.getRevenueForMonth(ctx, restaurantID, month)
	previousYearRevenue := e.getRevenueForMonth(ctx, restaurantID, month.AddDate(-1, 0, 0))

	if previousYearRevenue == 0 {
		return 0
	}

	return ((currentRevenue - previousYearRevenue) / previousYearRevenue) * 100
}

func (e *RevenueFeatureETL) getRevenueForMonth(ctx context.Context, restaurantID string, month time.Time) float64 {
	collection := e.db.GetCollection("feature_revenue_monthly")
	monthStr := month.Format("2006-01")

	var feature models.FeatureRevenueMonthly
	err := collection.FindOne(ctx, bson.M{
		"restaurant_id": restaurantID,
		"month":         monthStr,
	}).Decode(&feature)

	if err != nil {
		return 0
	}

	return feature.Revenue
}

func (e *RevenueFeatureETL) getSeason(month time.Month) string {
	switch {
	case month >= 1 && month <= 3:
		return "Q1"
	case month >= 4 && month <= 6:
		return "Q2"
	case month >= 7 && month <= 9:
		return "Q3"
	default:
		return "Q4"
	}
}

func (e *RevenueFeatureETL) isHolidayMonth(month time.Month) bool {
	// Vietnamese holidays: Tết (Jan/Feb), 30/4, 1/5, 2/9
	holidays := []time.Month{1, 2, 4, 5, 9, 12}
	for _, h := range holidays {
		if month == h {
			return true
		}
	}
	return false
}

func (e *RevenueFeatureETL) saveFeature(ctx context.Context, feature *models.FeatureRevenueMonthly) error {
	collection := e.db.GetCollection("feature_revenue_monthly")

	filter := bson.M{
		"restaurant_id": feature.RestaurantID,
		"month":         feature.Month,
	}

	update := bson.M{
		"$set": feature,
	}

	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (e *RevenueFeatureETL) saveETLLog(ctx context.Context, etlLog *models.ETLLog, status string, err error) {
	etlLog.Status = status
	etlLog.CompletedAt = time.Now()
	etlLog.Duration = int64(time.Since(etlLog.StartedAt).Seconds())

	if err != nil {
		etlLog.ErrorMessage = err.Error()
		etlLog.StackTrace = fmt.Sprintf("%+v", err)
	}

	collection := e.db.GetCollection("etl_logs")
	collection.InsertOne(ctx, etlLog)
}
