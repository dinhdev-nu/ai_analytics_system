# Backend API Service

Go-based REST API service using Gin framework.

## 📁 Structure

```
backend/
├── cmd/
│   └── api/
│       └── main.go           # Entry point
├── internal/
│   ├── config/
│   │   └── config.go         # Configuration management
│   ├── database/
│   │   └── mongodb.go        # MongoDB connection
│   ├── handlers/
│   │   └── analytics_handler.go  # HTTP handlers
│   ├── models/
│   │   └── response.go       # Response models
│   └── services/
│       └── analytics_service.go  # Business logic
├── go.mod
└── Dockerfile
```

## 🚀 Quick Start

### Development

```bash
# Install dependencies
go mod download

# Run server
go run cmd/api/main.go

# Or with hot reload (install air first)
air
```

### Production

```bash
# Build
go build -o bin/api cmd/api/main.go

# Run
./bin/api
```

## 🔌 API Endpoints

### Health Check
```bash
GET /health
```

### Get Revenue Forecast
```bash
GET /api/v1/analytics/forecast?restaurant_id=REST001&months=12
```

### Get Dashboard
```bash
GET /api/v1/analytics/dashboard?restaurant_id=REST001
```

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 📝 Configuration

Edit `.env` file:

```env
API_PORT=8080
API_HOST=0.0.0.0
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=ai_analytics
REDIS_HOST=localhost
REDIS_PORT=6379
JWT_SECRET=your-secret-key
LOG_LEVEL=info
```

## 🔐 Authentication

Currently supports JWT tokens (to be implemented).

```go
// Add to headers
Authorization: Bearer <jwt_token>
```

## 📊 Monitoring

Logs are output in JSON format:

```json
{
  "level": "info",
  "time": "2026-02-15T10:30:00Z",
  "method": "GET",
  "path": "/api/v1/analytics/dashboard",
  "status": 200,
  "latency_ms": 45
}
```

## 🐛 Debugging

```bash
# Install Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug
dlv debug cmd/api/main.go
```

## 📚 Adding New Endpoints

1. Add handler in `internal/handlers/`
2. Add service logic in `internal/services/`
3. Register route in `cmd/api/main.go`

Example:

```go
// internal/handlers/new_handler.go
func (h *Handler) NewEndpoint(c *gin.Context) {
    // Implementation
}

// cmd/api/main.go
router.GET("/api/v1/new-endpoint", handler.NewEndpoint)
```
