# Project Structure

Complete directory structure and file organization of AI Analytics system.

---

## 📁 Root Directory

```
AI_analysis/
├── backend/                    # Go backend API service
├── client/                     # React frontend application
├── database/                   # Database schemas and scripts
├── docs/                       # Documentation
├── etl/                        # ETL worker service (Go)
├── ml/                         # Machine learning services (Python)
├── .env.example                # Environment variables template
├── .gitignore                  # Git ignore rules
├── CHANGELOG.md                # Version history and changes
├── CODE_OF_CONDUCT.md          # Community guidelines
├── CONTRIBUTING.md             # Contribution guidelines
├── docker-compose.yml          # Docker orchestration
├── FAQ.md                      # Frequently asked questions
├── LICENSE                     # MIT License
├── Makefile                    # Build and deployment automation
├── manage.bat                  # Management script (Windows)
├── manage.sh                   # Management script (Linux/Mac)
├── PROJECT_STRUCTURE.md        # This file
├── README.md                   # Project overview
├── ROADMAP.md                  # Future development plans
├── SECURITY.md                 # Security policy
└── TROUBLESHOOTING.md          # Problem-solving guide
```

---

## 🔧 Backend Service (Go)

```
backend/
├── cmd/
│   └── api/
│       └── main.go             # API server entry point
├── internal/
│   ├── config/
│   │   └── config.go           # Configuration management
│   ├── database/
│   │   └── mongodb.go          # MongoDB connection handler
│   ├── handlers/
│   │   └── analytics_handler.go # HTTP request handlers
│   ├── models/
│   │   └── response.go         # API response models
│   └── services/
│       └── analytics_service.go # Business logic layer
├── Dockerfile                  # Multi-stage Docker build
├── go.mod                      # Go module dependencies
├── go.sum                      # Dependency checksums
└── README.md                   # Backend documentation
```

### Key Files

| File | Purpose | Lines |
|------|---------|-------|
| `cmd/api/main.go` | API server, routing, middleware setup | ~200 |
| `internal/handlers/analytics_handler.go` | HTTP handlers for endpoints | ~150 |
| `internal/services/analytics_service.go` | Core business logic | ~500 |
| `internal/config/config.go` | Config loading from env vars | ~80 |
| `internal/database/mongodb.go` | Database connection pooling | ~100 |

---

## 🎨 Frontend Client (React)

```
client/
├── public/                     # Static assets
├── src/
│   ├── api/
│   │   └── analytics.js        # API client (Axios)
│   ├── components/
│   │   ├── Insights.jsx        # AI insights component
│   │   ├── OrdersChart.jsx     # Bar chart for orders
│   │   ├── RevenueChart.jsx    # Revenue forecast chart
│   │   └── SummaryCards.jsx    # Metric cards
│   ├── App.css                 # Application styles
│   ├── App.jsx                 # Main application component
│   └── main.jsx                # React entry point
├── Dockerfile                  # Multi-stage build with Nginx
├── index.html                  # HTML template
├── nginx.conf                  # Nginx configuration
├── package.json                # Node dependencies
├── README.md                   # Frontend documentation
└── vite.config.js              # Vite build configuration
```

### Component Hierarchy

```
App.jsx
├── Restaurant Selector (dropdown)
├── SummaryCards.jsx
│   ├── Current Month Revenue
│   ├── MoM Growth
│   ├── YoY Growth
│   ├── Total Orders
│   ├── Average Order Value
│   ├── Next Month Forecast
│   └── Forecast Confidence
├── RevenueChart.jsx
│   ├── Actual Revenue (line)
│   ├── Predicted Revenue (dashed line)
│   └── Target Revenue (dotted line)
├── OrdersChart.jsx
│   └── Order Count (bars)
└── Insights.jsx
    ├── Success Insights (green)
    ├── Warning Insights (yellow)
    └── Info Insights (blue)
```

---

## 🗄️ Database Layer

```
database/
├── mongodb/
│   ├── indexes.js              # Performance indexes
│   ├── sample_data.js          # Sample data generator
│   └── schemas.js              # Collection schemas
└── README.md                   # Database documentation
```

### Collections

| Collection | Documents | Purpose |
|------------|-----------|---------|
| `restaurants` | Restaurant master data | Store restaurant info |
| `orders` | Raw order transactions | Historical order data |
| `payments` | Payment transactions | Payment tracking |
| `feature_revenue_monthly` | Engineered features | ML model inputs |
| `revenue_predictions` | ML predictions | Forecast results |
| `ml_models` | Model metadata | Model versioning |
| `etl_logs` | ETL job logs | Job tracking |
| `api_logs` | API request logs | Monitoring (future) |

---

## ⚙️ ETL Service (Go)

```
etl/
├── cmd/
│   └── etl_worker/
│       └── main.go             # ETL worker with cron
├── internal/
│   ├── config/
│   │   └── config.go           # ETL configuration
│   ├── database/
│   │   └── mongodb.go          # Database connection
│   ├── etl/
│   │   └── revenue_features.go # Feature engineering logic
│   └── models/
│       └── models.go           # Data models
├── Dockerfile                  # Docker build
├── go.mod                      # Go dependencies
├── go.sum                      # Checksums
└── README.md                   # ETL documentation
```

### Feature Engineering Pipeline

```
Raw Orders (MongoDB)
    ↓
Extract & Aggregate
    ↓ (monthly grouping)
Base Metrics
    - revenue
    - order_count
    - avg_order_value
    - unique_customers
    ↓
Calculate Rolling Averages
    - rolling_avg_3m
    - rolling_avg_6m
    - rolling_avg_12m
    ↓
Calculate Growth Rates
    - mom_growth (Month-over-Month %)
    - yoy_growth (Year-over-Year %)
    ↓
Detect Seasonality
    - is_holiday_season (boolean)
    - seasonality_index (float)
    ↓
Store Features
    ↓
feature_revenue_monthly (MongoDB)
```

---

## 🤖 ML Services (Python)

```
ml/
├── models/                     # Saved models (generated)
│   └── *.pkl                   # Pickled Prophet models
├── notebooks/                  # Jupyter notebooks (optional)
├── prediction/
│   └── batch_predict.py        # Batch prediction service
├── training/
│   └── train_revenue_forecast.py # Model training service
├── config.py                   # ML configuration
├── database.py                 # MongoDB utilities
├── Dockerfile.prediction       # Prediction service Docker
├── Dockerfile.training         # Training service Docker
├── README.md                   # ML documentation
└── requirements.txt            # Python dependencies
```

### ML Pipeline

```
feature_revenue_monthly
    ↓
Load Data (training/train_revenue_forecast.py)
    ↓
Prepare Prophet Format
    ds (date), y (revenue), regressors
    ↓
Train/Test Split (80/20)
    ↓
Train Prophet Model
    - yearly_seasonality: True
    - changepoint_prior_scale: 0.05
    - seasonality_mode: multiplicative
    ↓
Evaluate on Test Set
    - MAPE, RMSE, MAE, R² Score
    ↓
Save Model (.pkl file)
Store Metadata (ml_models collection)
    ↓
Generate Predictions (prediction/batch_predict.py)
    ↓
12-Month Forecast
    - predicted value
    - confidence interval
    - confidence score
    ↓
Store Predictions (revenue_predictions collection)
```

---

## 📚 Documentation

```
docs/
├── api.md                      # API reference & examples
├── architecture.md             # System architecture
├── deployment.md               # Deployment guides
├── quickstart.md               # 5-minute setup guide
└── usage.md                    # Usage instructions
```

### Documentation Coverage

| Document | Content | Pages |
|----------|---------|-------|
| `architecture.md` | System design, data flow, components | 15 |
| `deployment.md` | Docker, Kubernetes, manual setup | 12 |
| `api.md` | Endpoint specs, request/response examples | 10 |
| `usage.md` | Service-by-service usage guide | 8 |
| `quickstart.md` | Quick start in 5 minutes | 3 |

---

## 🐳 Docker Configuration

```
docker-compose.yml              # Orchestrate 6 services
backend/Dockerfile              # Backend API image
client/Dockerfile               # Frontend + Nginx image
etl/Dockerfile                  # ETL worker image
ml/Dockerfile.training          # ML training image
ml/Dockerfile.prediction        # ML prediction image
```

### Service Architecture

```
┌─────────────────────────────────────────────────┐
│                  Host Machine                    │
│  ┌───────────────────────────────────────────┐  │
│  │         Docker Compose Network            │  │
│  │                                           │  │
│  │  ┌──────────┐    ┌──────────┐           │  │
│  │  │ MongoDB  │    │  Redis   │           │  │
│  │  │  :27017  │    │  :6379   │           │  │
│  │  └────┬─────┘    └────┬─────┘           │  │
│  │       │               │                  │  │
│  │  ┌────┴───────────────┴─────────┐       │  │
│  │  │                               │       │  │
│  │  │  ┌────────────┐  ┌─────────┐ │       │  │
│  │  │  │ ETL Worker │  │Backend  │ │       │  │
│  │  │  │  (cron)    │  │API:8080 │ │       │  │
│  │  │  └────────────┘  └────┬────┘ │       │  │
│  │  │                        │      │       │  │
│  │  │  ┌────────────┐  ┌────┴────┐ │       │  │
│  │  │  │ML Training │  │   ML    │ │       │  │
│  │  │  │  (weekly)  │  │Predict  │ │       │  │
│  │  │  └────────────┘  └─────────┘ │       │  │
│  │  │                               │       │  │
│  │  └───────────────────────────────┘       │  │
│  │                    │                     │  │
│  │          ┌─────────┴─────────┐          │  │
│  │          │  Client (Nginx)   │          │  │
│  │          │      :3000        │          │  │
│  │          └───────────────────┘          │  │
│  └───────────────────────────────────────────┘  │
│                     │                            │
└─────────────────────┼────────────────────────────┘
                      │
                  Browser
               http://localhost:3000
```

---

## 🛠️ Build & Deployment

```
Makefile                        # 30+ make targets
manage.sh                       # Bash management script
manage.bat                      # Windows batch script
```

### Makefile Targets

```
Setup & Management:
  make setup          # Initialize system
  make start          # Start all services
  make stop           # Stop all services
  make restart        # Restart services
  make status         # Check service status
  make logs           # View all logs

ML Operations:
  make etl            # Run ETL job
  make train          # Train models
  make predict        # Generate predictions

Testing:
  make test           # Run all tests
  make test-backend   # Test backend only
  make test-ml        # Test ML services
  make test-client    # Test frontend

Building:
  make build          # Build all services
  make build-backend  # Build backend image
  make build-etl      # Build ETL image
  make build-client   # Build frontend image

Development:
  make dev-backend    # Run backend in dev mode
  make dev-client     # Run frontend dev server
  make dev-ml         # ML development mode

Database:
  make db-shell       # MongoDB shell
  make db-backup      # Backup database
  make db-restore     # Restore database

Monitoring:
  make health         # System health check
  make metrics        # View metrics

Cleanup:
  make clean          # Clean all
  make clean-cache    # Clean cache only
  make clean-logs     # Clean logs

Deployment:
  make deploy-staging # Deploy to staging
  make deploy-prod    # Deploy to production

Utilities:
  make fmt            # Format code
  make lint           # Run linters
  make deps           # Update dependencies
```

---

## 📦 Dependencies

### Backend (Go)

```go
// Key dependencies
github.com/gin-gonic/gin           // Web framework
go.mongodb.org/mongo-driver        // MongoDB driver
github.com/redis/go-redis/v9       // Redis client
github.com/joho/godotenv           // Environment variables
github.com/robfig/cron/v3          // Cron scheduler
```

### Frontend (JavaScript)

```json
{
  "react": "^18.2.0",              // UI library
  "react-dom": "^18.2.0",          // React DOM
  "axios": "^1.6.0",               // HTTP client
  "echarts": "^5.4.0",             // Charts library
  "echarts-for-react": "^3.0.2"    // React wrapper
}
```

### ML (Python)

```
prophet==1.1.5                     # Time-series forecasting
xgboost==2.0.3                     # Gradient boosting (future)
scikit-learn==1.3.2                # ML utilities
pandas==2.1.4                      # Data manipulation
numpy==1.26.2                      # Numerical computing
pymongo==4.6.1                     # MongoDB driver
```

---

## 📊 File Statistics

### Lines of Code

| Language | Files | Code | Comments | Blank | Total |
|----------|-------|------|----------|-------|-------|
| Go       | 15    | 2,500| 400      | 600   | 3,500 |
| Python   | 8     | 1,800| 300      | 400   | 2,500 |
| JavaScript | 12  | 1,500| 200      | 300   | 2,000 |
| Markdown | 18    | 8,000| -        | 1,500 | 9,500 |
| YAML     | 2     | 200  | 50       | 30    | 280   |
| **Total** | **55** | **14,000** | **950** | **2,830** | **17,780** |

### Documentation

- **Total docs:** 18 markdown files
- **Total words:** ~45,000
- **Total pages:** ~120 (estimated)
- **Coverage:** System design, API, deployment, usage, troubleshooting

---

## 🔗 File Relationships

### Data Flow

```
Raw Data (orders.csv)
    ↓ mongoimport
MongoDB (orders collection)
    ↓ ETL reads
ETL Service (Go)
    ↓ Feature engineering
MongoDB (feature_revenue_monthly)
    ↓ ML reads
ML Training (Python)
    ↓ Model training
Saved Models (*.pkl)
    ↓ Prediction uses
ML Prediction (Python)
    ↓ Generate forecasts
MongoDB (revenue_predictions)
    ↓ API reads
Backend API (Go)
    ↓ REST endpoints
Frontend Client (React)
    ↓ Visualize
User's Browser
```

### Code Dependencies

```
Backend API
├─→ internal/config (config loading)
├─→ internal/database (MongoDB connection)
├─→ internal/services
│   ├─→ internal/models (data models)
│   └─→ external: MongoDB, Redis
└─→ internal/handlers
    └─→ internal/services (business logic)

Frontend Client
├─→ src/api (API client)
│   └─→ external: Backend API
├─→ src/components (UI components)
│   ├─→ external: ECharts
│   └─→ src/api (data fetching)
└─→ src/App.jsx
    └─→ src/components

ETL Worker
├─→ internal/config
├─→ internal/database
├─→ internal/models
└─→ internal/etl
    ├─→ internal/models
    └─→ external: MongoDB

ML Services
├─→ config.py
├─→ database.py (MongoDB utilities)
│   └─→ external: MongoDB
└─→ training/prediction
    ├─→ config.py
    ├─→ database.py
    └─→ external: Prophet, scikit-learn
```

---

## 🎯 Quick Navigation

### I want to...

**Add a new API endpoint:**
→ `backend/internal/handlers/` + `backend/cmd/api/main.go`

**Modify the dashboard:**
→ `client/src/components/` + `client/src/App.jsx`

**Add a new feature:**
→ `etl/internal/etl/revenue_features.go`

**Train a new model:**
→ `ml/training/` (create new trainer)

**Change database schema:**
→ `database/mongodb/schemas.js`

**Update documentation:**
→ `docs/` or root `*.md` files

**Fix a bug:**
→ See [TROUBLESHOOTING.md](TROUBLESHOOTING.md)

**Deploy to production:**
→ See [docs/deployment.md](docs/deployment.md)

---

## 📝 File Naming Conventions

### Backend (Go)
- **Files:** `snake_case.go`
- **Packages:** `lowercase`
- **Types:** `PascalCase`
- **Functions:** `PascalCase` (exported), `camelCase` (private)

### Frontend (JavaScript)
- **Files:** `PascalCase.jsx` (components), `camelCase.js` (utilities)
- **Components:** `PascalCase`
- **Functions:** `camelCase`
- **Constants:** `UPPER_SNAKE_CASE`

### ML (Python)
- **Files:** `snake_case.py`
- **Classes:** `PascalCase`
- **Functions:** `snake_case`
- **Constants:** `UPPER_SNAKE_CASE`

### Documentation
- **Files:** `UPPER_CASE.md` (root), `lowercase.md` (docs/)

---

## 🔍 Search Tips

Find files by purpose:

```bash
# Configuration files
find . -name "config.*" -o -name "*.config.*"

# Docker files
find . -name "Dockerfile*" -o -name "docker-compose.yml"

# Test files
find . -name "*_test.*" -o -name "*.test.*"

# Documentation
find . -name "*.md"
```

---

**Last Updated:** February 15, 2026
**Total Files:** 55+
**Total Size:** ~25 MB (with dependencies)
