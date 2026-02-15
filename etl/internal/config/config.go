package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	MongoURI     string
	MongoDB      string
	ETLBatchSize int
	Environment  string
	LogLevel     string
}

func Load() *Config {
	// Load .env file if exists
	if err := godotenv.Load("../../.env"); err != nil {
		log.Warn("No .env file found, using environment variables")
	}

	batchSize, _ := strconv.Atoi(getEnv("ETL_BATCH_SIZE", "1000"))

	return &Config{
		MongoURI:     getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		MongoDB:      getEnv("MONGODB_DATABASE", "ai_analytics"),
		ETLBatchSize: batchSize,
		Environment:  getEnv("ENVIRONMENT", "development"),
		LogLevel:     getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
