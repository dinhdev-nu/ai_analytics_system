"""
Configuration for ML training and prediction
"""
import os
from dotenv import load_dotenv

load_dotenv()

class Config:
    # MongoDB
    MONGODB_URI = os.getenv("MONGODB_URI", "mongodb://localhost:27017")
    MONGODB_DATABASE = os.getenv("MONGODB_DATABASE", "ai_analytics")
    
    # MLflow
    MLFLOW_TRACKING_URI = os.getenv("MLFLOW_TRACKING_URI", "http://localhost:5000")
    MLFLOW_EXPERIMENT_NAME = "revenue_forecasting"
    
    # Model paths
    MODEL_PATH = os.getenv("ML_MODEL_PATH", "./models")
    MODEL_VERSION = "v1.0.0"
    
    # Training parameters
    TRAIN_TEST_SPLIT = 0.8
    MIN_TRAINING_MONTHS = 12  # Minimum data required for training
    
    # Prediction parameters
    PREDICTION_MONTHS_AHEAD = int(os.getenv("PREDICTION_MONTHS_AHEAD", "12"))
    
    # Environment
    ENVIRONMENT = os.getenv("ENVIRONMENT", "development")
    LOG_LEVEL = os.getenv("LOG_LEVEL", "INFO")


config = Config()
