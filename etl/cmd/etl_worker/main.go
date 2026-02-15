package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ai-analytics/etl/internal/config"
	"ai-analytics/etl/internal/database"
	"ai-analytics/etl/internal/etl"

	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Setup logger
	setupLogger(cfg.LogLevel)

	log.Info("Starting ETL Worker...")
	log.Infof("Environment: %s", cfg.Environment)

	// Connect to MongoDB
	mongoDB, err := database.NewMongoDB(cfg.MongoURI, cfg.MongoDB)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoDB.Close()

	// Create ETL instance
	revenueETL := etl.NewRevenueFeatureETL(mongoDB)

	// Run once immediately on startup
	log.Info("Running initial ETL job...")
	ctx := context.Background()
	if err := revenueETL.Run(ctx); err != nil {
		log.Errorf("Initial ETL job failed: %v", err)
	}

	// Setup cron scheduler
	c := cron.New()

	// Schedule ETL job (daily at 2 AM)
	c.AddFunc("0 2 * * *", func() {
		log.Info("Running scheduled ETL job...")
		ctx := context.Background()
		if err := revenueETL.Run(ctx); err != nil {
			log.Errorf("Scheduled ETL job failed: %v", err)
		}
	})

	c.Start()
	log.Info("ETL Worker started. Cron schedule: daily at 2 AM")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Info("Shutting down ETL Worker...")

	// Stop cron scheduler
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	c.Stop()
	log.Info("ETL Worker stopped gracefully")
}

func setupLogger(level string) {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	switch level {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}
