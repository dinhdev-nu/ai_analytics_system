#!/bin/bash

# AI Analytics Management Script
# Usage: ./manage.sh [command]

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

function print_info {
    echo -e "${GREEN}[INFO]${NC} $1"
}

function print_warn {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

function print_error {
    echo -e "${RED}[ERROR]${NC} $1"
}

function check_docker {
    if ! command -v docker &> /dev/null; then
        print_error "Docker is not installed"
        exit 1
    fi
    if ! command -v docker-compose &> /dev/null; then
        print_error "Docker Compose is not installed"
        exit 1
    fi
}

function setup {
    print_info "Setting up AI Analytics System..."
    
    # Check prerequisites
    check_docker
    
    # Copy environment file
    if [ ! -f .env ]; then
        cp .env.example .env
        print_info "Created .env file from .env.example"
        print_warn "Please edit .env with your configuration"
    else
        print_warn ".env already exists, skipping..."
    fi
    
    # Build images
    print_info "Building Docker images..."
    docker-compose build
    
    # Start services
    print_info "Starting services..."
    docker-compose up -d mongodb redis
    
    # Wait for MongoDB to be ready
    print_info "Waiting for MongoDB to be ready..."
    sleep 10
    
    # Import sample data
    print_info "Importing sample data..."
    docker exec -i ai_analytics_mongodb mongosh <<EOF
use ai_analytics
load("/docker-entrypoint-initdb.d/schemas.js")
load("/docker-entrypoint-initdb.d/indexes.js")
load("/docker-entrypoint-initdb.d/sample_data.js")
EOF
    
    # Run ETL
    print_info "Running ETL job..."
    docker-compose run --rm etl_worker
    
    # Train models
    print_info "Training ML models..."
    docker-compose run --rm ml_training
    
    # Generate predictions
    print_info "Generating predictions..."
    docker-compose run --rm ml_prediction
    
    # Start all services
    print_info "Starting all services..."
    docker-compose up -d
    
    print_info "Setup completed successfully!"
    print_info "Frontend: http://localhost:3000"
    print_info "Backend API: http://localhost:8080"
}

function start {
    print_info "Starting all services..."
    check_docker
    docker-compose up -d
    print_info "All services started"
    print_info "Frontend: http://localhost:3000"
    print_info "Backend API: http://localhost:8080"
}

function stop {
    print_info "Stopping all services..."
    check_docker
    docker-compose down
    print_info "All services stopped"
}

function restart {
    print_info "Restarting all services..."
    stop
    start
}

function status {
    print_info "Service status:"
    docker-compose ps
}

function logs {
    SERVICE=${1:-}
    if [ -z "$SERVICE" ]; then
        docker-compose logs -f
    else
        docker-compose logs -f $SERVICE
    fi
}

function etl {
    print_info "Running ETL job..."
    docker-compose run --rm etl_worker
}

function train {
    print_info "Training ML models..."
    docker-compose run --rm ml_training
}

function predict {
    print_info "Generating predictions..."
    docker-compose run --rm ml_prediction
}

function backup {
    BACKUP_DIR="backups/$(date +%Y%m%d_%H%M%S)"
    mkdir -p $BACKUP_DIR
    
    print_info "Backing up MongoDB to $BACKUP_DIR..."
    docker exec ai_analytics_mongodb mongodump \
        --uri="mongodb://admin:password123@localhost:27017" \
        --out=/tmp/backup
    
    docker cp ai_analytics_mongodb:/tmp/backup $BACKUP_DIR/mongodb
    
    print_info "Backup completed: $BACKUP_DIR"
}

function health {
    print_info "Checking system health..."
    
    # Check backend
    echo -n "Backend API: "
    if curl -sf http://localhost:8080/health > /dev/null; then
        echo -e "${GREEN}✓ OK${NC}"
    else
        echo -e "${RED}✗ FAILED${NC}"
    fi
    
    # Check MongoDB
    echo -n "MongoDB: "
    if docker exec ai_analytics_mongodb mongosh --eval "db.adminCommand('ping')" > /dev/null 2>&1; then
        echo -e "${GREEN}✓ OK${NC}"
    else
        echo -e "${RED}✗ FAILED${NC}"
    fi
    
    # Check Redis
    echo -n "Redis: "
    if docker exec ai_analytics_redis redis-cli ping > /dev/null 2>&1; then
        echo -e "${GREEN}✓ OK${NC}"
    else
        echo -e "${RED}✗ FAILED${NC}"
    fi
    
    # Check Frontend
    echo -n "Frontend: "
    if curl -sf http://localhost:3000 > /dev/null; then
        echo -e "${GREEN}✓ OK${NC}"
    else
        echo -e "${RED}✗ FAILED${NC}"
    fi
}

function clean {
    print_warn "This will remove all containers, volumes, and data!"
    read -p "Are you sure? (yes/no): " confirm
    if [ "$confirm" == "yes" ]; then
        print_info "Cleaning up..."
        docker-compose down -v
        rm -rf ml/models/*.pkl
        print_info "Cleanup completed"
    else
        print_info "Cleanup cancelled"
    fi
}

function help {
    cat << EOF
AI Analytics Management Script

Usage: ./manage.sh [command]

Commands:
  setup         Initial setup (run once)
  start         Start all services
  stop          Stop all services
  restart       Restart all services
  status        Show service status
  logs [svc]    Show logs (optional: specific service)
  
  etl           Run ETL job
  train         Train ML models
  predict       Generate predictions
  
  backup        Backup MongoDB data
  health        Check system health
  clean         Remove all containers and data
  
  help          Show this help message

Examples:
  ./manage.sh setup                  # Initial setup
  ./manage.sh start                  # Start all services
  ./manage.sh logs backend           # Show backend logs
  ./manage.sh etl                    # Run ETL manually
  ./manage.sh health                 # Check health

EOF
}

# Main script
COMMAND=${1:-help}

case $COMMAND in
    setup)
        setup
        ;;
    start)
        start
        ;;
    stop)
        stop
        ;;
    restart)
        restart
        ;;
    status)
        status
        ;;
    logs)
        logs ${2:-}
        ;;
    etl)
        etl
        ;;
    train)
        train
        ;;
    predict)
        predict
        ;;
    backup)
        backup
        ;;
    health)
        health
        ;;
    clean)
        clean
        ;;
    help|*)
        help
        ;;
esac
