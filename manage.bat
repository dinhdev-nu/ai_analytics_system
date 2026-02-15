@echo off
REM AI Analytics Management Script for Windows
REM Usage: manage.bat [command]

setlocal enabledelayedexpansion

if "%1"=="" goto help
if "%1"=="setup" goto setup
if "%1"=="start" goto start
if "%1"=="stop" goto stop
if "%1"=="restart" goto restart
if "%1"=="status" goto status
if "%1"=="logs" goto logs
if "%1"=="etl" goto etl
if "%1"=="train" goto train
if "%1"=="predict" goto predict
if "%1"=="backup" goto backup
if "%1"=="health" goto health
if "%1"=="clean" goto clean
if "%1"=="help" goto help
goto help

:setup
echo [INFO] Setting up AI Analytics System...

REM Copy environment file
if not exist .env (
    copy .env.example .env
    echo [INFO] Created .env file from .env.example
    echo [WARN] Please edit .env with your configuration
) else (
    echo [WARN] .env already exists, skipping...
)

REM Build images
echo [INFO] Building Docker images...
docker-compose build

REM Start services
echo [INFO] Starting services...
docker-compose up -d mongodb redis

REM Wait for MongoDB
echo [INFO] Waiting for MongoDB to be ready...
timeout /t 10

REM Import sample data
echo [INFO] Importing sample data...
docker exec -i ai_analytics_mongodb mongosh --eval "load('/docker-entrypoint-initdb.d/sample_data.js')"

REM Run ETL
echo [INFO] Running ETL job...
docker-compose run --rm etl_worker

REM Train models
echo [INFO] Training ML models...
docker-compose run --rm ml_training

REM Generate predictions
echo [INFO] Generating predictions...
docker-compose run --rm ml_prediction

REM Start all services
echo [INFO] Starting all services...
docker-compose up -d

echo [INFO] Setup completed successfully!
echo [INFO] Frontend: http://localhost:3000
echo [INFO] Backend API: http://localhost:8080
goto end

:start
echo [INFO] Starting all services...
docker-compose up -d
echo [INFO] All services started
echo [INFO] Frontend: http://localhost:3000
echo [INFO] Backend API: http://localhost:8080
goto end

:stop
echo [INFO] Stopping all services...
docker-compose down
echo [INFO] All services stopped
goto end

:restart
echo [INFO] Restarting all services...
call :stop
call :start
goto end

:status
echo [INFO] Service status:
docker-compose ps
goto end

:logs
if "%2"=="" (
    docker-compose logs -f
) else (
    docker-compose logs -f %2
)
goto end

:etl
echo [INFO] Running ETL job...
docker-compose run --rm etl_worker
goto end

:train
echo [INFO] Training ML models...
docker-compose run --rm ml_training
goto end

:predict
echo [INFO] Generating predictions...
docker-compose run --rm ml_prediction
goto end

:backup
set BACKUP_DIR=backups\%date:~10,4%%date:~4,2%%date:~7,2%_%time:~0,2%%time:~3,2%%time:~6,2%
mkdir %BACKUP_DIR%

echo [INFO] Backing up MongoDB to %BACKUP_DIR%...
docker exec ai_analytics_mongodb mongodump --uri="mongodb://admin:password123@localhost:27017" --out=/tmp/backup
docker cp ai_analytics_mongodb:/tmp/backup %BACKUP_DIR%\mongodb

echo [INFO] Backup completed: %BACKUP_DIR%
goto end

:health
echo [INFO] Checking system health...

REM Check backend
echo Backend API:
curl -sf http://localhost:8080/health >nul 2>&1
if %errorlevel%==0 (
    echo [OK] Backend is running
) else (
    echo [FAILED] Backend is not responding
)

REM Check MongoDB
echo MongoDB:
docker exec ai_analytics_mongodb mongosh --eval "db.adminCommand('ping')" >nul 2>&1
if %errorlevel%==0 (
    echo [OK] MongoDB is running
) else (
    echo [FAILED] MongoDB is not responding
)

REM Check Redis
echo Redis:
docker exec ai_analytics_redis redis-cli ping >nul 2>&1
if %errorlevel%==0 (
    echo [OK] Redis is running
) else (
    echo [FAILED] Redis is not responding
)

REM Check Frontend
echo Frontend:
curl -sf http://localhost:3000 >nul 2>&1
if %errorlevel%==0 (
    echo [OK] Frontend is running
) else (
    echo [FAILED] Frontend is not responding
)
goto end

:clean
echo [WARN] This will remove all containers, volumes, and data!
set /p confirm="Are you sure? (yes/no): "
if "%confirm%"=="yes" (
    echo [INFO] Cleaning up...
    docker-compose down -v
    del /q ml\models\*.pkl 2>nul
    echo [INFO] Cleanup completed
) else (
    echo [INFO] Cleanup cancelled
)
goto end

:help
echo AI Analytics Management Script
echo.
echo Usage: manage.bat [command]
echo.
echo Commands:
echo   setup         Initial setup (run once)
echo   start         Start all services
echo   stop          Stop all services
echo   restart       Restart all services
echo   status        Show service status
echo   logs [svc]    Show logs (optional: specific service)
echo.
echo   etl           Run ETL job
echo   train         Train ML models
echo   predict       Generate predictions
echo.
echo   backup        Backup MongoDB data
echo   health        Check system health
echo   clean         Remove all containers and data
echo.
echo   help          Show this help message
echo.
echo Examples:
echo   manage.bat setup                  # Initial setup
echo   manage.bat start                  # Start all services
echo   manage.bat logs backend           # Show backend logs
echo   manage.bat etl                    # Run ETL manually
echo   manage.bat health                 # Check health
goto end

:end
endlocal
