# API Documentation

## Base URL
```
Development: http://localhost:8080/api/v1
Production: https://api.yourdomain.com/api/v1
```

## Authentication
All API requests require authentication using JWT tokens (to be implemented).

```http
Authorization: Bearer <your_jwt_token>
```

---

## Endpoints

### 1. Health Check

**GET** `/health`

Check if API server is running.

**Response**
```json
{
  "success": true,
  "message": "API is running"
}
```

---

### 2. Get Revenue Forecast

**GET** `/analytics/forecast`

Get revenue forecast with predictions for a restaurant.

**Query Parameters**

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| restaurant_id | string | Yes | - | Restaurant identifier |
| months | integer | No | 12 | Number of months to return (1-24) |

**Example Request**
```bash
curl "http://localhost:8080/api/v1/analytics/forecast?restaurant_id=REST001&months=12"
```

**Response 200 OK**
```json
{
  "restaurant_id": "REST001",
  "timestamps": [
    "2026-01",
    "2026-02",
    "2026-03",
    "2026-04",
    "2026-05",
    "2026-06"
  ],
  "actual": [
    120000000,
    135000000,
    0,
    0,
    0,
    0
  ],
  "predicted": [
    125000000,
    135000000,
    142000000,
    148000000,
    155000000,
    162000000
  ],
  "target": [
    130000000,
    140000000,
    145000000,
    150000000,
    160000000,
    165000000
  ],
  "confidence_data": [
    {
      "month": "2026-01",
      "lower": 115000000,
      "upper": 135000000
    },
    {
      "month": "2026-02",
      "lower": 125000000,
      "upper": 145000000
    }
  ],
  "model_info": {
    "model_name": "revenue_forecast_prophet_REST001",
    "model_version": "v1.0.0",
    "last_updated": "2026-02-15T02:00:00Z",
    "metrics": {
      "mape": 5.2,
      "rmse": 8500000,
      "mae": 6200000,
      "r2_score": 0.92
    }
  }
}
```

**Response 400 Bad Request**
```json
{
  "error": "bad_request",
  "message": "restaurant_id is required",
  "code": 400
}
```

**Response 500 Internal Server Error**
```json
{
  "error": "internal_error",
  "message": "Failed to fetch forecast data",
  "code": 500
}
```

---

### 3. Get Dashboard Data

**GET** `/analytics/dashboard`

Get complete dashboard data including summary, charts, and insights.

**Query Parameters**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| restaurant_id | string | Yes | Restaurant identifier |

**Example Request**
```bash
curl "http://localhost:8080/api/v1/analytics/dashboard?restaurant_id=REST001"
```

**Response 200 OK**
```json
{
  "restaurant_id": "REST001",
  "summary": {
    "current_month_revenue": 135000000,
    "previous_month_revenue": 120000000,
    "month_over_month_growth": 12.5,
    "year_over_year_growth": 15.2,
    "total_orders": 2800,
    "avg_order_value": 48214,
    "forecast_next_month": 142000000,
    "forecast_confidence": 0.85
  },
  "revenue_chart": {
    "labels": [
      "2025-03",
      "2025-04",
      "2025-05",
      "2025-06",
      "2025-07",
      "2025-08",
      "2025-09",
      "2025-10",
      "2025-11",
      "2025-12",
      "2026-01",
      "2026-02",
      "2026-03",
      "2026-04",
      "2026-05",
      "2026-06"
    ],
    "actual": [
      110000000,
      115000000,
      118000000,
      122000000,
      125000000,
      128000000,
      130000000,
      133000000,
      138000000,
      150000000,
      120000000,
      135000000,
      0,
      0,
      0,
      0
    ],
    "predicted": [
      0,
      0,
      0,
      0,
      0,
      0,
      0,
      0,
      0,
      0,
      0,
      0,
      142000000,
      148000000,
      155000000,
      162000000
    ],
    "target": [
      110000000,
      115000000,
      120000000,
      125000000,
      130000000,
      135000000,
      135000000,
      140000000,
      145000000,
      160000000,
      130000000,
      140000000,
      145000000,
      150000000,
      160000000,
      165000000
    ]
  },
  "orders_chart": {
    "labels": [
      "2025-03",
      "2025-04",
      "2025-05",
      "2025-06",
      "2025-07",
      "2025-08",
      "2025-09",
      "2025-10",
      "2025-11",
      "2025-12",
      "2026-01",
      "2026-02"
    ],
    "order_counts": [
      2400,
      2500,
      2550,
      2600,
      2650,
      2700,
      2750,
      2800,
      2900,
      3200,
      2800,
      3000
    ]
  },
  "insights": [
    {
      "type": "success",
      "title": "Tăng trưởng mạnh",
      "description": "Doanh thu tháng này tăng mạnh so với tháng trước",
      "value": "+12.5%"
    },
    {
      "type": "info",
      "title": "Dự báo tích cực",
      "description": "Doanh thu tháng tới dự kiến tăng",
      "value": "142000000 VND"
    }
  ]
}
```

---

## Data Models

### Summary Object
```typescript
{
  current_month_revenue: number,      // Doanh thu tháng hiện tại (VND)
  previous_month_revenue: number,     // Doanh thu tháng trước (VND)
  month_over_month_growth: number,    // Tăng trưởng MoM (%)
  year_over_year_growth: number,      // Tăng trưởng YoY (%)
  total_orders: number,               // Tổng số đơn hàng
  avg_order_value: number,            // Giá trị đơn hàng trung bình (VND)
  forecast_next_month: number,        // Dự báo tháng tới (VND)
  forecast_confidence: number         // Độ tin cậy (0-1)
}
```

### Revenue Chart Data
```typescript
{
  labels: string[],        // Array of month strings ("2026-01")
  actual: number[],        // Actual revenue values
  predicted: number[],     // Predicted revenue values (0 for past months)
  target: number[]         // Target revenue values
}
```

### Confidence Range
```typescript
{
  month: string,           // Month identifier
  lower: number,           // Lower confidence interval
  upper: number            // Upper confidence interval
}
```

### Model Info
```typescript
{
  model_name: string,
  model_version: string,
  last_updated: string,    // ISO 8601 timestamp
  metrics: {
    mape: number,          // Mean Absolute Percentage Error (%)
    rmse: number,          // Root Mean Square Error
    mae: number,           // Mean Absolute Error
    r2_score: number       // R² Score (0-1)
  }
}
```

### Insight Object
```typescript
{
  type: "success" | "warning" | "info",
  title: string,
  description: string,
  value?: string           // Optional value
}
```

---

## Error Handling

All error responses follow this format:

```json
{
  "error": "error_code",
  "message": "Human readable error message",
  "code": 400
}
```

### Common Error Codes

| HTTP Status | Error Code | Description |
|-------------|------------|-------------|
| 400 | bad_request | Invalid request parameters |
| 401 | unauthorized | Missing or invalid authentication |
| 403 | forbidden | Insufficient permissions |
| 404 | not_found | Resource not found |
| 429 | rate_limit_exceeded | Too many requests |
| 500 | internal_error | Server error |
| 503 | service_unavailable | Service temporarily unavailable |

---

## Rate Limiting

API requests are rate limited per IP address:

```
- Development: 1000 requests per hour
- Production: 100 requests per minute
```

Rate limit headers are included in responses:

```http
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1234567890
```

---

## Caching

Responses are cached with the following TTL:

| Endpoint | Cache TTL |
|----------|-----------|
| /analytics/forecast | 1 hour |
| /analytics/dashboard | 30 minutes |

Cache can be bypassed using header:

```http
Cache-Control: no-cache
```

---

## Examples

### JavaScript/TypeScript (Axios)

```typescript
import axios from 'axios';

const api = axios.create({
  baseURL: 'http://localhost:8080/api/v1',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  }
});

// Get dashboard
const getDashboard = async (restaurantId: string) => {
  try {
    const response = await api.get('/analytics/dashboard', {
      params: { restaurant_id: restaurantId }
    });
    return response.data;
  } catch (error) {
    console.error('Error fetching dashboard:', error);
    throw error;
  }
};

// Get forecast
const getForecast = async (restaurantId: string, months: number = 12) => {
  try {
    const response = await api.get('/analytics/forecast', {
      params: { restaurant_id: restaurantId, months }
    });
    return response.data;
  } catch (error) {
    console.error('Error fetching forecast:', error);
    throw error;
  }
};
```

### Python (Requests)

```python
import requests

BASE_URL = "http://localhost:8080/api/v1"
headers = {
    "Authorization": f"Bearer {token}",
    "Content-Type": "application/json"
}

def get_dashboard(restaurant_id):
    response = requests.get(
        f"{BASE_URL}/analytics/dashboard",
        params={"restaurant_id": restaurant_id},
        headers=headers
    )
    response.raise_for_status()
    return response.json()

def get_forecast(restaurant_id, months=12):
    response = requests.get(
        f"{BASE_URL}/analytics/forecast",
        params={"restaurant_id": restaurant_id, "months": months},
        headers=headers
    )
    response.raise_for_status()
    return response.json()
```

### cURL

```bash
# Get dashboard
curl -X GET \
  "http://localhost:8080/api/v1/analytics/dashboard?restaurant_id=REST001" \
  -H "Authorization: Bearer your_token_here" \
  -H "Content-Type: application/json"

# Get forecast
curl -X GET \
  "http://localhost:8080/api/v1/analytics/forecast?restaurant_id=REST001&months=12" \
  -H "Authorization: Bearer your_token_here" \
  -H "Content-Type: application/json"
```

---

## Webhooks (Future)

Subscribe to events:

```json
POST /api/v1/webhooks/subscribe
{
  "event": "prediction.completed",
  "url": "https://your-server.com/webhook",
  "secret": "your_webhook_secret"
}
```

Event types:
- `prediction.completed` - Batch prediction completed
- `training.completed` - Model training completed
- `etl.completed` - ETL job completed
- `anomaly.detected` - Anomaly detected in data
