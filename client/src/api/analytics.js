import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor
api.interceptors.request.use(
  (config) => {
    // Add auth token if available
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor
api.interceptors.response.use(
  (response) => response.data,
  (error) => {
    if (error.response) {
      console.error('API Error:', error.response.data);
      return Promise.reject(error.response.data);
    }
    return Promise.reject(error);
  }
);

export const analyticsAPI = {
  // Get revenue forecast
  getForecast: (restaurantId, months = 12) => {
    return api.get('/analytics/forecast', {
      params: { restaurant_id: restaurantId, months },
    });
  },

  // Get dashboard data
  getDashboard: (restaurantId) => {
    return api.get('/analytics/dashboard', {
      params: { restaurant_id: restaurantId },
    });
  },

  // Health check
  healthCheck: () => {
    return api.get('/health');
  },
};

export default api;
