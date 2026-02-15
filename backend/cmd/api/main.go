package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ai-analytics/backend/internal/config"
	"ai-analytics/backend/internal/database"
	"ai-analytics/backend/internal/handlers"
	"ai-analytics/backend/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Setup logger
	setupLogger(cfg.LogLevel)

	log.Info("Starting AI Analytics API Server...")
	log.Infof("Environment: %s", cfg.Environment)

	// Connect to MongoDB
	mongoDB, err := database.NewMongoDB(cfg.MongoURI, cfg.MongoDB)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoDB.Close()

	// Initialize services
	analyticsService := services.NewAnalyticsService(mongoDB)

	// Initialize handlers
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService)

	// Setup Gin
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Routes
	setupRoutes(router, analyticsHandler)

	// Start server
	addr := fmt.Sprintf("%s:%s", cfg.APIHost, cfg.APIPort)
	log.Infof("Server starting on %s", addr)

	// Graceful shutdown
	go func() {
		if err := router.Run(addr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")
	log.Info("Server stopped gracefully")
}

func setupRoutes(router *gin.Engine, analyticsHandler *handlers.AnalyticsHandler) {
	// Health check
	router.GET("/health", analyticsHandler.HealthCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		analytics := v1.Group("/analytics")
		{
			analytics.GET("/forecast", analyticsHandler.GetRevenueForecast)
			analytics.GET("/dashboard", analyticsHandler.GetDashboard)
		}
	}

	log.Info("Routes registered:")
	log.Info("  GET  /health")
	log.Info("  GET  /api/v1/analytics/forecast")
	log.Info("  GET  /api/v1/analytics/dashboard")
}

func setupLogger(level string) {
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
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

	log.Infof("Log level set to: %s", level)
}
