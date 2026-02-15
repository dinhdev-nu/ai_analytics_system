# ML Training & Prediction Services

Python-based machine learning services for revenue forecasting.

## 📁 Structure

```
ml/
├── training/
│   └── train_revenue_forecast.py    # Model training
├── prediction/
│   └── batch_predict.py             # Batch prediction
├── models/                          # Saved models (generated)
├── notebooks/                       # Jupyter notebooks
├── config.py                        # Configuration
├── database.py                      # MongoDB utilities
└── requirements.txt                 # Python dependencies
```

## 🚀 Quick Start

### Setup Environment

```bash
# Create virtual environment
python -m venv venv

# Activate
source venv/bin/activate  # Linux/Mac
venv\Scripts\activate     # Windows

# Install dependencies
pip install -r requirements.txt
```

### Training Models

```bash
# Train all restaurants
python training/train_revenue_forecast.py

# Output: models/*.pkl files
```

### Generate Predictions

```bash
# Generate predictions for all restaurants
python prediction/batch_predict.py
```

## 🤖 Models

### Prophet

Time-series forecasting with seasonality.

**Parameters:**
```python
model = Prophet(
    yearly_seasonality=True,
    weekly_seasonality=False,
    daily_seasonality=False,
    seasonality_mode='multiplicative',
    changepoint_prior_scale=0.05
)
```

**Regressors:**
- `rolling_avg_3m`: 3-month rolling average
- `is_holiday`: Holiday indicator

### XGBoost (Future)

Gradient boosting for tabular data.

## 📊 Model Evaluation

Metrics tracked:
- **MAPE**: Mean Absolute Percentage Error
- **RMSE**: Root Mean Square Error
- **MAE**: Mean Absolute Error
- **R² Score**: Coefficient of determination

Target performance:
- MAPE < 10%
- R² Score > 0.85

## 🧪 Testing

```bash
# Run tests
python -m pytest tests/

# With coverage
python -m pytest --cov=. tests/

# Specific test
python -m pytest tests/test_training.py
```

## 📈 Monitoring

### View Model Performance

```python
from database import db

# Get latest model
model = db.get_collection("ml_models").find_one(
    {"model_name": "revenue_forecast_prophet_REST001"},
    sort=[("trained_at", -1)]
)

print(model["metrics"])
```

### Compare Predictions vs Actuals

```python
predictions = db.get_collection("revenue_predictions").find({
    "restaurant_id": "REST001",
    "actual": {"$ne": None}
})

for pred in predictions:
    error = abs(pred["predicted"] - pred["actual"]) / pred["actual"]
    print(f"{pred['month']}: {error:.2%} error")
```

## 🔧 Configuration

Edit `config.py`:

```python
MONGODB_URI = "mongodb://localhost:27017"
MONGODB_DATABASE = "ai_analytics"
ML_MODEL_PATH = "./models"
MODEL_VERSION = "v1.0.0"
PREDICTION_MONTHS_AHEAD = 12
```

## 📚 Adding New Features

### 1. Add feature in ETL

Edit `etl/internal/etl/revenue_features.go`

### 2. Update model training

```python
# In train_revenue_forecast.py
prophet_df['new_feature'] = df['new_feature']
model.add_regressor('new_feature')
```

### 3. Update prediction

```python
# In batch_predict.py
future['new_feature'] = calculate_future_feature()
```

## 🐛 Debugging

### IPython Interactive Shell

```bash
ipython

>>> from database import db
>>> from training.train_revenue_forecast import RevenueForecastingModel
>>> model = RevenueForecastingModel()
>>> # Test interactively
```

### Jupyter Notebooks

```bash
jupyter notebook notebooks/

# Create new notebook for experiments
```

## 📊 Model Versioning

Models are versioned automatically:

```
revenue_forecast_prophet_REST001_v1.0.0.pkl
revenue_forecast_prophet_REST001_v1.0.1.pkl
```

Metadata stored in MongoDB:

```javascript
{
  model_name: "revenue_forecast_prophet_REST001",
  version: "v1.0.0",
  trained_at: ISODate("2026-02-15T02:00:00Z"),
  metrics: {
    mape: 5.2,
    rmse: 8500000,
    mae: 6200000,
    r2_score: 0.92
  },
  is_production: true
}
```

## 🔄 Retraining Strategy

### Weekly Retrain

```bash
# Add to crontab
0 3 * * 0 cd /app/ml && python training/train_revenue_forecast.py
```

### On-Demand Retrain

```bash
python training/train_revenue_forecast.py --restaurant REST001 --force
```

## 📝 Best Practices

1. **Always validate data** before training
2. **Save model metadata** for versioning
3. **Track performance metrics** over time
4. **A/B test new models** before production
5. **Monitor prediction accuracy** regularly
