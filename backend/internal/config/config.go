package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	// Server
	APIPort     string
	APIHost     string
	Environment string

	// Database
	MongoURI string
	MongoDB  string

	// Redis
	RedisHost     string
	RedisPort     string
	RedisPassword string

	// JWT
	JWTSecret string

	// Rate Limiting
	RateLimit int

	// Logging
	LogLevel string
}

func Load() *Config {
	// Load .env file
	if err := godotenv.Load("../../.env"); err != nil {
		log.Warn("No .env file found")
	}

	rateLimit, _ := strconv.Atoi(getEnv("API_RATE_LIMIT", "100"))

	return &Config{
		APIPort:       getEnv("API_PORT", "8080"),
		APIHost:       getEnv("API_HOST", "0.0.0.0"),
		Environment:   getEnv("ENVIRONMENT", "development"),
		MongoURI:      getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		MongoDB:       getEnv("MONGODB_DATABASE", "ai_analytics"),
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		RateLimit:     rateLimit,
		LogLevel:      getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
