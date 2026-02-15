"""
Database connection and utilities
"""
from pymongo import MongoClient
from typing import List, Dict, Any
import pandas as pd
from loguru import logger
from config import config


class MongoDB:
    def __init__(self):
        self.client = MongoClient(config.MONGODB_URI)
        self.db = self.client[config.MONGODB_DATABASE]
        logger.info(f"Connected to MongoDB: {config.MONGODB_DATABASE}")
    
    def get_collection(self, collection_name: str):
        return self.db[collection_name]
    
    def fetch_features(self, restaurant_id: str) -> pd.DataFrame:
        """
        Fetch feature_revenue_monthly data for a restaurant
        """
        collection = self.get_collection("feature_revenue_monthly")
        
        query = {"restaurant_id": restaurant_id}
        cursor = collection.find(query).sort("month", 1)
        
        data = list(cursor)
        if not data:
            logger.warning(f"No features found for restaurant: {restaurant_id}")
            return pd.DataFrame()
        
        df = pd.DataFrame(data)
        df['month'] = pd.to_datetime(df['month'])
        df = df.sort_values('month')
        
        logger.info(f"Fetched {len(df)} records for restaurant {restaurant_id}")
        return df
    
    def get_all_restaurants(self) -> List[str]:
        """
        Get list of all active restaurant IDs
        """
        collection = self.get_collection("restaurants")
        cursor = collection.find({"status": "active"}, {"restaurant_id": 1})
        
        restaurants = [doc["restaurant_id"] for doc in cursor]
        logger.info(f"Found {len(restaurants)} active restaurants")
        return restaurants
    
    def save_model_metadata(self, metadata: Dict[str, Any]):
        """
        Save model metadata to ml_models collection
        """
        collection = self.get_collection("ml_models")
        
        # Check if model version exists
        existing = collection.find_one({
            "model_name": metadata["model_name"],
            "version": metadata["version"]
        })
        
        if existing:
            # Update existing
            collection.update_one(
                {"_id": existing["_id"]},
                {"$set": metadata}
            )
            logger.info(f"Updated model metadata: {metadata['model_name']} {metadata['version']}")
        else:
            # Insert new
            collection.insert_one(metadata)
            logger.info(f"Saved new model metadata: {metadata['model_name']} {metadata['version']}")
    
    def get_latest_model_version(self, model_name: str) -> str:
        """
        Get latest model version for a model name
        """
        collection = self.get_collection("ml_models")
        
        model = collection.find_one(
            {"model_name": model_name, "is_production": True},
            sort=[("trained_at", -1)]
        )
        
        if model:
            return model["version"]
        return None
    
    def close(self):
        self.client.close()
        logger.info("MongoDB connection closed")


# Singleton instance
db = MongoDB()
