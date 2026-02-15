.PHONY: help setup start stop restart status logs health clean test build deploy

# Default target
.DEFAULT_GOAL := help

help: ## Show this help message
	@echo "AI Analytics - Makefile Commands"
	@echo "================================="
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Setup & Management

setup: ## Initial setup (run once)
	@echo "Setting up AI Analytics System..."
	@cp -n .env.example .env || true
	@docker-compose build
	@docker-compose up -d mongodb redis
	@sleep 10
	@docker exec -i ai_analytics_mongodb mongosh < database/mongodb/sample_data.js
	@$(MAKE) etl
	@$(MAKE) train
	@$(MAKE) predict
	@docker-compose up -d
	@echo "Setup completed! Frontend: http://localhost:3000"

start: ## Start all services
	@echo "Starting all services..."
	@docker-compose up -d
	@echo "Frontend: http://localhost:3000"
	@echo "Backend API: http://localhost:8080"

stop: ## Stop all services
	@echo "Stopping all services..."
	@docker-compose down

restart: stop start ## Restart all services

status: ## Show service status
	@docker-compose ps

logs: ## Show logs (usage: make logs SERVICE=backend)
ifdef SERVICE
	@docker-compose logs -f $(SERVICE)
else
	@docker-compose logs -f
endif

##@ ML Operations

etl: ## Run ETL job
	@echo "Running ETL job..."
	@docker-compose run --rm etl_worker

train: ## Train ML models
	@echo "Training ML models..."
	@docker-compose run --rm ml_training

predict: ## Generate predictions
	@echo "Generating predictions..."
	@docker-compose run --rm ml_prediction

##@ Testing

test: test-backend test-ml ## Run all tests

test-backend: ## Test backend (Go)
	@echo "Testing backend..."
	@cd backend && go test ./...

test-ml: ## Test ML services (Python)
	@echo "Testing ML services..."
	@cd ml && python -m pytest tests/

test-client: ## Test frontend (React)
	@echo "Testing frontend..."
	@cd client && npm test

##@ Building

build: ## Build all Docker images
	@echo "Building Docker images..."
	@docker-compose build

build-backend: ## Build backend only
	@cd backend && go build -o bin/api cmd/api/main.go

build-etl: ## Build ETL worker only
	@cd etl && go build -o bin/etl_worker cmd/etl_worker/main.go

build-client: ## Build frontend only
	@cd client && npm run build

##@ Development

dev-backend: ## Run backend in development mode
	@cd backend && go run cmd/api/main.go

dev-client: ## Run frontend in development mode
	@cd client && npm run dev

dev-ml: ## Run ML training locally
	@cd ml && python training/train_revenue_forecast.py

##@ Database

db-shell: ## Open MongoDB shell
	@docker exec -it ai_analytics_mongodb mongosh

db-backup: ## Backup MongoDB
	@mkdir -p backups
	@docker exec ai_analytics_mongodb mongodump \
		--uri="mongodb://admin:password123@localhost:27017" \
		--out=/tmp/backup
	@docker cp ai_analytics_mongodb:/tmp/backup backups/mongodb-$(shell date +%Y%m%d-%H%M%S)
	@echo "Backup completed"

db-restore: ## Restore MongoDB (usage: make db-restore BACKUP=backups/mongodb-20260215-120000)
	@docker exec -i ai_analytics_mongodb mongorestore --uri="mongodb://admin:password123@localhost:27017" $(BACKUP)

##@ Monitoring

health: ## Check system health
	@echo "Checking system health..."
	@echo -n "Backend API: "
	@curl -sf http://localhost:8080/health > /dev/null && echo "✓ OK" || echo "✗ FAILED"
	@echo -n "MongoDB: "
	@docker exec ai_analytics_mongodb mongosh --eval "db.adminCommand('ping')" > /dev/null 2>&1 && echo "✓ OK" || echo "✗ FAILED"
	@echo -n "Redis: "
	@docker exec ai_analytics_redis redis-cli ping > /dev/null 2>&1 && echo "✓ OK" || echo "✗ FAILED"

metrics: ## Show system metrics
	@echo "System Metrics:"
	@echo "=============="
	@docker stats --no-stream --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}"

##@ Cleanup

clean: ## Remove all containers and data
	@echo "⚠️  This will remove all containers, volumes, and data!"
	@read -p "Are you sure? (yes/no): " confirm && \
	if [ "$$confirm" = "yes" ]; then \
		docker-compose down -v; \
		rm -rf ml/models/*.pkl; \
		echo "Cleanup completed"; \
	else \
		echo "Cleanup cancelled"; \
	fi

clean-cache: ## Clear Redis cache
	@docker exec ai_analytics_redis redis-cli FLUSHALL
	@echo "Cache cleared"

clean-logs: ## Remove old logs
	@docker-compose logs --tail=0 > /dev/null
	@echo "Logs cleared"

##@ Documentation

docs: ## Generate API documentation
	@echo "Generating API documentation..."
	@cd backend && swag init -g cmd/api/main.go

docs-serve: ## Serve documentation locally
	@cd docs && python -m http.server 8000

##@ Deployment

deploy-staging: ## Deploy to staging
	@echo "Deploying to staging..."
	@docker-compose -f docker-compose.staging.yml up -d

deploy-prod: ## Deploy to production
	@echo "Deploying to production..."
	@docker-compose -f docker-compose.prod.yml up -d

##@ Utilities

fmt: ## Format code
	@echo "Formatting Go code..."
	@cd backend && go fmt ./...
	@cd etl && go fmt ./...
	@echo "Formatting Python code..."
	@cd ml && black .
	@echo "Formatting done"

lint: ## Lint code
	@echo "Linting Go code..."
	@cd backend && golangci-lint run
	@cd etl && golangci-lint run
	@echo "Linting Python code..."
	@cd ml && pylint training/ prediction/

deps: ## Update dependencies
	@echo "Updating Go dependencies..."
	@cd backend && go get -u ./... && go mod tidy
	@cd etl && go get -u ./... && go mod tidy
	@echo "Updating Python dependencies..."
	@cd ml && pip install --upgrade -r requirements.txt
	@echo "Updating Node dependencies..."
	@cd client && npm update
