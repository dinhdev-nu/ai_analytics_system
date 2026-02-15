"""
Revenue Forecasting Model Training using Prophet
Prophet is ideal for time-series with strong seasonality
"""
import os
import pandas as pd
import numpy as np
from prophet import Prophet
from sklearn.metrics import mean_absolute_percentage_error, mean_squared_error
import joblib
from datetime import datetime
from loguru import logger
import sys

# Add parent directory to path
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from config import config
from database import db


class RevenueForecastingModel:
    def __init__(self, model_name="prophet"):
        self.model_name = model_name
        self.model = None
        self.metrics = {}
        self.version = config.MODEL_VERSION
        
    def prepare_data(self, df: pd.DataFrame) -> pd.DataFrame:
        """
        Prepare data for Prophet
        Prophet requires columns: ds (date), y (target)
        """
        if df.empty:
            raise ValueError("Empty dataframe provided")
        
        # Prophet format
        prophet_df = pd.DataFrame({
            'ds': df['month'],
            'y': df['revenue']
        })
        
        # Add regressors (optional)
        if 'rolling_avg_3m' in df.columns:
            prophet_df['rolling_avg_3m'] = df['rolling_avg_3m']
        if 'is_holiday' in df.columns:
            prophet_df['is_holiday'] = df['is_holiday'].astype(int)
        
        return prophet_df
    
    def train(self, df: pd.DataFrame):
        """
        Train Prophet model
        """
        logger.info(f"Training Prophet model with {len(df)} data points")
        
        # Initialize Prophet with Vietnamese holidays considerations
        self.model = Prophet(
            yearly_seasonality=True,
            weekly_seasonality=False,
            daily_seasonality=False,
            seasonality_mode='multiplicative',
            changepoint_prior_scale=0.05,
            seasonality_prior_scale=10.0
        )
        
        # Add custom regressors
        if 'rolling_avg_3m' in df.columns:
            self.model.add_regressor('rolling_avg_3m')
        if 'is_holiday' in df.columns:
            self.model.add_regressor('is_holiday')
        
        # Fit model
        self.model.fit(df)
        logger.info("Prophet model trained successfully")
    
    def predict(self, periods: int = 12) -> pd.DataFrame:
        """
        Make predictions for future periods
        """
        if self.model is None:
            raise ValueError("Model not trained yet")
        
        # Create future dataframe
        future = self.model.make_future_dataframe(periods=periods, freq='MS')
        
        # Add regressors for future dates (use last known values)
        if 'rolling_avg_3m' in self.model.extra_regressors:
            # For simplicity, use last known value
            future['rolling_avg_3m'] = future['rolling_avg_3m'].fillna(method='ffill').fillna(0)
        if 'is_holiday' in self.model.extra_regressors:
            # Mark common holiday months
            future['is_holiday'] = future['ds'].dt.month.isin([1, 2, 4, 5, 9, 12]).astype(int)
        
        # Predict
        forecast = self.model.predict(future)
        
        return forecast[['ds', 'yhat', 'yhat_lower', 'yhat_upper']]
    
    def evaluate(self, df_train: pd.DataFrame, df_test: pd.DataFrame):
        """
        Evaluate model performance
        """
        # Predict on test set
        future = self.model.make_future_dataframe(periods=len(df_test), freq='MS')
        
        # Add regressors
        if 'rolling_avg_3m' in self.model.extra_regressors:
            future['rolling_avg_3m'] = 0
        if 'is_holiday' in self.model.extra_regressors:
            future['is_holiday'] = 0
        
        forecast = self.model.predict(future)
        
        # Get predictions for test period
        test_start_idx = len(df_train)
        y_pred = forecast['yhat'].values[test_start_idx:test_start_idx + len(df_test)]
        y_true = df_test['y'].values
        
        # Calculate metrics
        mape = mean_absolute_percentage_error(y_true, y_pred) * 100
        rmse = np.sqrt(mean_squared_error(y_true, y_pred))
        mae = np.mean(np.abs(y_true - y_pred))
        
        # R² score
        ss_res = np.sum((y_true - y_pred) ** 2)
        ss_tot = np.sum((y_true - np.mean(y_true)) ** 2)
        r2_score = 1 - (ss_res / ss_tot) if ss_tot != 0 else 0
        
        self.metrics = {
            "mape": float(mape),
            "rmse": float(rmse),
            "mae": float(mae),
            "r2_score": float(r2_score)
        }
        
        logger.info(f"Model Evaluation Metrics:")
        logger.info(f"  MAPE: {mape:.2f}%")
        logger.info(f"  RMSE: {rmse:,.0f}")
        logger.info(f"  MAE: {mae:,.0f}")
        logger.info(f"  R²: {r2_score:.4f}")
        
        return self.metrics
    
    def save_model(self, restaurant_id: str):
        """
        Save trained model to disk
        """
        if self.model is None:
            raise ValueError("No model to save")
        
        # Create models directory
        os.makedirs(config.MODEL_PATH, exist_ok=True)
        
        # Save model
        model_filename = f"revenue_forecast_prophet_{restaurant_id}_{self.version}.pkl"
        model_path = os.path.join(config.MODEL_PATH, model_filename)
        
        joblib.dump(self.model, model_path)
        logger.info(f"Model saved to: {model_path}")
        
        # Save metadata to MongoDB
        model_metadata = {
            "model_name": f"revenue_forecast_prophet_{restaurant_id}",
            "model_type": "prophet",
            "version": self.version,
            "file_path": model_path,
            "file_size": os.path.getsize(model_path),
            "trained_at": datetime.now(),
            "metrics": self.metrics,
            "status": "active",
            "is_production": True
        }
        
        db.save_model_metadata(model_metadata)
        
        return model_path
    
    def load_model(self, model_path: str):
        """
        Load trained model from disk
        """
        if not os.path.exists(model_path):
            raise FileNotFoundError(f"Model not found: {model_path}")
        
        self.model = joblib.load(model_path)
        logger.info(f"Model loaded from: {model_path}")


def train_restaurant_model(restaurant_id: str):
    """
    Train model for a specific restaurant
    """
    logger.info(f"Starting training for restaurant: {restaurant_id}")
    
    # Fetch features
    df = db.fetch_features(restaurant_id)
    
    if len(df) < config.MIN_TRAINING_MONTHS:
        logger.warning(f"Insufficient data for {restaurant_id}: {len(df)} months < {config.MIN_TRAINING_MONTHS} required")
        return None
    
    # Initialize model
    model = RevenueForecastingModel()
    
    # Prepare data
    prophet_df = model.prepare_data(df)
    
    # Train-test split
    split_idx = int(len(prophet_df) * config.TRAIN_TEST_SPLIT)
    df_train = prophet_df[:split_idx]
    df_test = prophet_df[split_idx:]
    
    # Train
    start_time = datetime.now()
    model.train(df_train)
    training_duration = (datetime.now() - start_time).total_seconds()
    
    # Evaluate
    if len(df_test) > 0:
        model.evaluate(df_train, df_test)
    
    # Save model
    model_path = model.save_model(restaurant_id)
    
    logger.info(f"Training completed for {restaurant_id} in {training_duration:.2f}s")
    
    return model


def main():
    """
    Main training script
    """
    logger.info("=" * 60)
    logger.info("Starting Revenue Forecasting Model Training")
    logger.info("=" * 60)
    
    # Get all restaurants
    restaurants = db.get_all_restaurants()
    logger.info(f"Found {len(restaurants)} restaurants to train")
    
    trained_count = 0
    failed_count = 0
    
    for restaurant_id in restaurants:
        try:
            model = train_restaurant_model(restaurant_id)
            if model:
                trained_count += 1
            else:
                failed_count += 1
        except Exception as e:
            logger.error(f"Failed to train model for {restaurant_id}: {str(e)}")
            failed_count += 1
    
    logger.info("=" * 60)
    logger.info(f"Training Summary:")
    logger.info(f"  Total: {len(restaurants)}")
    logger.info(f"  Trained: {trained_count}")
    logger.info(f"  Failed: {failed_count}")
    logger.info("=" * 60)
    
    db.close()


if __name__ == "__main__":
    main()
