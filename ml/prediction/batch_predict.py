"""
Batch Prediction Service
Loads trained models and generates predictions for all restaurants
Saves predictions to MongoDB
"""
import os
import sys
import pandas as pd
import numpy as np
import joblib
from datetime import datetime
from loguru import logger

# Add parent directory to path
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from config import config
from database import db


class BatchPredictor:
    def __init__(self):
        self.model = None
        self.restaurant_id = None
        
    def load_model(self, restaurant_id: str) -> bool:
        """
        Load trained model for a restaurant
        """
        self.restaurant_id = restaurant_id
        
        # Find model file
        model_pattern = f"revenue_forecast_prophet_{restaurant_id}_{config.MODEL_VERSION}.pkl"
        model_path = os.path.join(config.MODEL_PATH, model_pattern)
        
        if not os.path.exists(model_path):
            logger.warning(f"Model not found for {restaurant_id}: {model_path}")
            return False
        
        try:
            self.model = joblib.load(model_path)
            logger.info(f"Loaded model for {restaurant_id}")
            return True
        except Exception as e:
            logger.error(f"Failed to load model for {restaurant_id}: {str(e)}")
            return False
    
    def predict(self, periods: int = 12) -> pd.DataFrame:
        """
        Generate predictions for future periods
        """
        if self.model is None:
            raise ValueError("Model not loaded")
        
        # Create future dataframe
        future = self.model.make_future_dataframe(periods=periods, freq='MS')
        
        # Add regressors if needed
        if 'is_holiday' in self.model.extra_regressors:
            future['is_holiday'] = future['ds'].dt.month.isin([1, 2, 4, 5, 9, 12]).astype(int)
        
        if 'rolling_avg_3m' in self.model.extra_regressors:
            future['rolling_avg_3m'] = 0  # Simplified for prediction
        
        # Generate forecast
        forecast = self.model.predict(future)
        
        # Get only future predictions
        future_forecast = forecast.tail(periods)
        
        return future_forecast[['ds', 'yhat', 'yhat_lower', 'yhat_upper']]
    
    def save_predictions(self, predictions: pd.DataFrame):
        """
        Save predictions to MongoDB
        """
        collection = db.get_collection("revenue_predictions")
        
        records = []
        for _, row in predictions.iterrows():
            record = {
                "restaurant_id": self.restaurant_id,
                "month": row['ds'].strftime("%Y-%m"),
                "predicted": float(row['yhat']),
                "lower_ci": float(row['yhat_lower']),
                "upper_ci": float(row['yhat_upper']),
                "actual": None,  # Will be filled later when actual data is available
                "model_name": "prophet",
                "model_version": config.MODEL_VERSION,
                "confidence_score": self.calculate_confidence(row),
                "predicted_at": datetime.now(),
                "created_at": datetime.now()
            }
            records.append(record)
        
        # Bulk upsert
        for record in records:
            collection.update_one(
                {
                    "restaurant_id": record["restaurant_id"],
                    "month": record["month"],
                    "model_version": record["model_version"]
                },
                {"$set": record},
                upsert=True
            )
        
        logger.info(f"Saved {len(records)} predictions for {self.restaurant_id}")
    
    def calculate_confidence(self, row) -> float:
        """
        Calculate confidence score based on prediction interval width
        """
        interval_width = row['yhat_upper'] - row['yhat_lower']
        predicted_value = row['yhat']
        
        if predicted_value == 0:
            return 0.5
        
        # Narrower interval = higher confidence
        relative_width = interval_width / predicted_value
        confidence = max(0.0, min(1.0, 1.0 - relative_width / 2))
        
        return float(confidence)


def run_batch_prediction():
    """
    Run batch prediction for all restaurants
    """
    logger.info("=" * 60)
    logger.info("Starting Batch Prediction")
    logger.info("=" * 60)
    
    # Get all restaurants
    restaurants = db.get_all_restaurants()
    logger.info(f"Found {len(restaurants)} restaurants")
    
    predictor = BatchPredictor()
    success_count = 0
    failed_count = 0
    
    for restaurant_id in restaurants:
        try:
            # Load model
            if not predictor.load_model(restaurant_id):
                failed_count += 1
                continue
            
            # Generate predictions
            predictions = predictor.predict(periods=config.PREDICTION_MONTHS_AHEAD)
            
            # Save to database
            predictor.save_predictions(predictions)
            
            success_count += 1
            logger.info(f"✓ Completed predictions for {restaurant_id}")
            
        except Exception as e:
            logger.error(f"✗ Failed to predict for {restaurant_id}: {str(e)}")
            failed_count += 1
    
    logger.info("=" * 60)
    logger.info(f"Batch Prediction Summary:")
    logger.info(f"  Total: {len(restaurants)}")
    logger.info(f"  Success: {success_count}")
    logger.info(f"  Failed: {failed_count}")
    logger.info("=" * 60)


def update_actuals():
    """
    Update actual values in predictions table when real data becomes available
    """
    logger.info("Updating actual values in predictions...")
    
    collection = db.get_collection("revenue_predictions")
    features_collection = db.get_collection("feature_revenue_monthly")
    
    # Get predictions that don't have actual values yet
    predictions = collection.find({"actual": None})
    
    updated_count = 0
    
    for pred in predictions:
        # Check if actual data exists
        feature = features_collection.find_one({
            "restaurant_id": pred["restaurant_id"],
            "month": pred["month"]
        })
        
        if feature and "revenue" in feature:
            # Update prediction with actual value
            collection.update_one(
                {"_id": pred["_id"]},
                {
                    "$set": {
                        "actual": float(feature["revenue"]),
                        "updated_at": datetime.now()
                    }
                }
            )
            updated_count += 1
    
    logger.info(f"Updated {updated_count} predictions with actual values")


def main():
    """
    Main entry point
    """
    # Run batch predictions
    run_batch_prediction()
    
    # Update actuals
    update_actuals()
    
    # Close database connection
    db.close()


if __name__ == "__main__":
    main()
