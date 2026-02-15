// Sample data for development and testing
// Run: mongoimport --db ai_analytics --collection restaurants --file sample_data.json

// ==========================================
// RESTAURANTS
// ==========================================

db.restaurants.insertMany([
  {
    restaurant_id: "REST001",
    name: "Nhà hàng Phương Nam",
    location: "Hà Nội",
    created_at: new Date("2024-01-01"),
    status: "active",
    metadata: {
      cuisine_type: "Vietnamese",
      capacity: 100
    }
  },
  {
    restaurant_id: "REST002",
    name: "Quán Ăn Sài Gòn",
    location: "TP. Hồ Chí Minh",
    created_at: new Date("2024-02-01"),
    status: "active",
    metadata: {
      cuisine_type: "Vietnamese",
      capacity: 50
    }
  },
  {
    restaurant_id: "REST003",
    name: "BBQ House",
    location: "Đà Nẵng",
    created_at: new Date("2024-03-01"),
    status: "active",
    metadata: {
      cuisine_type: "Korean BBQ",
      capacity: 80
    }
  }
]);

// ==========================================
// ORDERS (Sample for last 12 months)
// ==========================================

// Generate sample orders
const restaurants = ["REST001", "REST002", "REST003"];
const startDate = new Date("2025-02-01");
const endDate = new Date("2026-02-15");

let orders = [];
let orderId = 1;

// Generate orders for each restaurant
restaurants.forEach(restaurantId => {
  let currentDate = new Date(startDate);
  
  while (currentDate <= endDate) {
    // Random 50-150 orders per day
    const ordersPerDay = Math.floor(Math.random() * 100) + 50;
    
    for (let i = 0; i < ordersPerDay; i++) {
      const orderDate = new Date(currentDate);
      orderDate.setHours(Math.floor(Math.random() * 14) + 8); // 8AM - 10PM
      
      orders.push({
        order_id: `ORD${String(orderId).padStart(8, '0')}`,
        restaurant_id: restaurantId,
        total_price: Math.floor(Math.random() * 400000) + 100000, // 100k - 500k VND
        status: "completed",
        items: [
          {
            item_id: `ITEM${Math.floor(Math.random() * 50) + 1}`,
            quantity: Math.floor(Math.random() * 3) + 1,
            price: Math.floor(Math.random() * 150000) + 50000
          }
        ],
        created_at: orderDate,
        completed_at: new Date(orderDate.getTime() + 1800000) // +30 minutes
      });
      
      orderId++;
    }
    
    // Next day
    currentDate.setDate(currentDate.getDate() + 1);
  }
});

db.orders.insertMany(orders.slice(0, 5000)); // Insert first 5000 orders as sample

// ==========================================
// PAYMENTS
// ==========================================

// Generate payments for orders
let payments = [];
orders.slice(0, 5000).forEach((order, idx) => {
  payments.push({
    payment_id: `PAY${String(idx + 1).padStart(8, '0')}`,
    order_id: order.order_id,
    amount: order.total_price,
    method: ["cash", "card", "ewallet"][Math.floor(Math.random() * 3)],
    status: "success",
    paid_at: order.completed_at,
    created_at: order.completed_at
  });
});

db.payments.insertMany(payments);

// ==========================================
// FEATURE_REVENUE_MONTHLY (Sample)
// ==========================================

db.feature_revenue_monthly.insertMany([
  {
    restaurant_id: "REST001",
    month: "2025-12",
    year: 2025,
    month_num: 12,
    revenue: 150000000,
    order_count: 3200,
    avg_order_value: 46875,
    rolling_avg_3m: 145000000,
    rolling_avg_6m: 140000000,
    rolling_avg_12m: 135000000,
    mom_growth: 5.5,
    yoy_growth: 15.2,
    season: "Q4",
    is_holiday: true,
    day_of_week_avg: 450000,
    target: 160000000,
    created_at: new Date(),
    updated_at: new Date(),
    version: "v1.0"
  },
  {
    restaurant_id: "REST001",
    month: "2026-01",
    year: 2026,
    month_num: 1,
    revenue: 120000000,
    order_count: 2800,
    avg_order_value: 42857,
    rolling_avg_3m: 138000000,
    rolling_avg_6m: 140000000,
    rolling_avg_12m: 137000000,
    mom_growth: -20.0,
    yoy_growth: 10.5,
    season: "Q1",
    is_holiday: true,
    day_of_week_avg: 380000,
    target: 130000000,
    created_at: new Date(),
    updated_at: new Date(),
    version: "v1.0"
  }
]);

// ==========================================
// REVENUE_PREDICTIONS (Sample)
// ==========================================

db.revenue_predictions.insertMany([
  {
    restaurant_id: "REST001",
    month: "2026-02",
    predicted: 135000000,
    lower_ci: 125000000,
    upper_ci: 145000000,
    actual: null,
    model_name: "prophet",
    model_version: "v1.0.0",
    confidence_score: 0.85,
    mape: 5.2,
    rmse: 8500000,
    predicted_at: new Date(),
    created_at: new Date()
  },
  {
    restaurant_id: "REST001",
    month: "2026-03",
    predicted: 142000000,
    lower_ci: 132000000,
    upper_ci: 152000000,
    actual: null,
    model_name: "prophet",
    model_version: "v1.0.0",
    confidence_score: 0.82,
    mape: 5.2,
    rmse: 8500000,
    predicted_at: new Date(),
    created_at: new Date()
  },
  {
    restaurant_id: "REST001",
    month: "2026-04",
    predicted: 148000000,
    lower_ci: 138000000,
    upper_ci: 158000000,
    actual: null,
    model_name: "prophet",
    model_version: "v1.0.0",
    confidence_score: 0.80,
    mape: 5.2,
    rmse: 8500000,
    predicted_at: new Date(),
    created_at: new Date()
  }
]);

print("Sample data inserted successfully!");
print(`Inserted ${orders.length} orders, ${payments.length} payments`);
