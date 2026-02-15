# Frontend Client

React-based dashboard with ECharts visualization.

## 📁 Structure

```
client/
├── src/
│   ├── api/
│   │   └── analytics.js          # API client
│   ├── components/
│   │   ├── RevenueChart.jsx      # Revenue forecast chart
│   │   ├── OrdersChart.jsx       # Orders bar chart
│   │   ├── SummaryCards.jsx      # Metric cards
│   │   └── Insights.jsx          # AI insights
│   ├── App.jsx                   # Main component
│   ├── App.css                   # Styles
│   └── main.jsx                  # Entry point
├── public/
├── index.html
├── package.json
└── vite.config.js
```

## 🚀 Quick Start

### Development

```bash
# Install dependencies
npm install

# Start dev server
npm run dev

# Open http://localhost:3000
```

### Production

```bash
# Build
npm run build

# Preview build
npm run preview
```

## 🎨 Components

### RevenueChart

Line chart showing actual, predicted, and target revenue.

**Props:**
```javascript
data: {
  labels: string[],      // ["2026-01", "2026-02", ...]
  actual: number[],      // Actual revenue values
  predicted: number[],   // Predicted values
  target: number[]       // Target values
}
```

**Usage:**
```jsx
<RevenueChart data={dashboard.revenue_chart} />
```

### OrdersChart

Bar chart showing order counts over time.

**Props:**
```javascript
data: {
  labels: string[],      // ["2026-01", "2026-02", ...]
  order_counts: number[] // Order counts
}
```

### SummaryCards

Display key metrics in card format.

**Props:**
```javascript
summary: {
  current_month_revenue: number,
  month_over_month_growth: number,
  year_over_year_growth: number,
  total_orders: number,
  avg_order_value: number,
  forecast_next_month: number,
  forecast_confidence: number
}
```

### Insights

Display AI-generated insights.

**Props:**
```javascript
insights: Array<{
  type: "success" | "warning" | "info",
  title: string,
  description: string,
  value?: string
}>
```

## 🔌 API Integration

### API Client

```javascript
// src/api/analytics.js
import axios from 'axios';

const api = axios.create({
  baseURL: 'http://localhost:8080/api/v1',
  timeout: 30000
});

export const analyticsAPI = {
  getForecast: (restaurantId, months) => 
    api.get('/analytics/forecast', { params: { restaurant_id: restaurantId, months }}),
  
  getDashboard: (restaurantId) =>
    api.get('/analytics/dashboard', { params: { restaurant_id: restaurantId }})
};
```

### Usage in Component

```jsx
import { analyticsAPI } from './api/analytics';

const [dashboard, setDashboard] = useState(null);

useEffect(() => {
  const fetchData = async () => {
    const data = await analyticsAPI.getDashboard(restaurantId);
    setDashboard(data);
  };
  fetchData();
}, [restaurantId]);
```

## 🎨 Styling

### CSS Variables

```css
:root {
  --color-primary: #667eea;
  --color-success: #10b981;
  --color-warning: #f59e0b;
  --color-error: #ef4444;
  --color-background: #f5f7fa;
  --color-text: #333;
}
```

### Responsive Design

```css
/* Mobile first */
.dashboard {
  display: grid;
  grid-template-columns: 1fr;
  gap: 20px;
}

/* Tablet and up */
@media (min-width: 768px) {
  .dashboard {
    grid-template-columns: repeat(2, 1fr);
  }
}

/* Desktop */
@media (min-width: 1024px) {
  .dashboard {
    grid-template-columns: repeat(4, 1fr);
  }
}
```

## 📊 ECharts Configuration

### Customizing Charts

```javascript
// In RevenueChart.jsx
const option = {
  // Colors
  color: ['#3b82f6', '#10b981', '#f59e0b'],
  
  // Title
  title: {
    text: 'Dự báo Doanh thu',
    textStyle: {
      fontSize: 20,
      fontWeight: 600
    }
  },
  
  // Tooltip
  tooltip: {
    trigger: 'axis',
    formatter: (params) => {
      // Custom formatting
    }
  },
  
  // Legend
  legend: {
    data: ['Thực tế', 'Dự đoán', 'Mục tiêu']
  },
  
  // Axis
  xAxis: {
    type: 'category',
    data: labels
  },
  yAxis: {
    type: 'value'
  },
  
  // Series
  series: [{
    name: 'Thực tế',
    type: 'line',
    data: actualData,
    smooth: true
  }]
};
```

## 🧪 Testing

```bash
# Run tests
npm test

# With coverage
npm test -- --coverage

# Watch mode
npm test -- --watch
```

## 🔧 Configuration

### Environment Variables

Create `.env` file:

```env
VITE_API_URL=http://localhost:8080/api/v1
```

### Vite Config

```javascript
// vite.config.js
export default defineConfig({
  plugins: [react()],
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  }
});
```

## 📦 Dependencies

Main dependencies:
- **react**: UI library
- **react-dom**: React DOM rendering
- **axios**: HTTP client
- **echarts**: Charting library
- **echarts-for-react**: React wrapper for ECharts

## 🚀 Deployment

### Build for Production

```bash
npm run build
# Output: dist/
```

### Deploy to Nginx

```nginx
server {
    listen 80;
    server_name yourdomain.com;
    root /var/www/dist;
    
    location / {
        try_files $uri $uri/ /index.html;
    }
    
    location /api {
        proxy_pass http://backend:8080;
    }
}
```

### Deploy to Vercel

```bash
npm install -g vercel
vercel deploy
```

## 🎯 Performance Optimization

### Code Splitting

```javascript
// Lazy load components
const RevenueChart = lazy(() => import('./components/RevenueChart'));

<Suspense fallback={<Loading />}>
  <RevenueChart data={data} />
</Suspense>
```

### Memoization

```javascript
// Memoize expensive calculations
const chartOption = useMemo(() => ({
  // Chart config
}), [data]);
```

### Debouncing

```javascript
// Debounce API calls
const fetchData = useCallback(
  debounce((id) => {
    analyticsAPI.getDashboard(id);
  }, 500),
  []
);
```

## 🐛 Debugging

### React DevTools

Install React DevTools browser extension.

### Console Logging

```javascript
console.log('Dashboard data:', dashboard);
console.table(dashboard.revenue_chart);
```

## 📚 Adding New Features

### 1. Create Component

```jsx
// src/components/NewChart.jsx
const NewChart = ({ data }) => {
  return (
    <div className="chart-container">
      <ReactECharts option={option} />
    </div>
  );
};

export default NewChart;
```

### 2. Import in App

```jsx
import NewChart from './components/NewChart';

// Use in App.jsx
<NewChart data={dashboard.new_data} />
```

### 3. Add Styling

```css
/* App.css */
.new-chart-container {
  background: white;
  border-radius: 12px;
  padding: 24px;
}
```
