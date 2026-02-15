import React, { useState, useEffect } from 'react';
import './App.css';
import { analyticsAPI } from './api/analytics';
import SummaryCards from './components/SummaryCards';
import RevenueChart from './components/RevenueChart';
import OrdersChart from './components/OrdersChart';
import Insights from './components/Insights';

function App() {
    const [restaurantId, setRestaurantId] = useState('REST001');
    const [dashboard, setDashboard] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        fetchDashboard();
    }, [restaurantId]);

    const fetchDashboard = async () => {
        try {
            setLoading(true);
            setError(null);
            const data = await analyticsAPI.getDashboard(restaurantId);
            setDashboard(data);
        } catch (err) {
            setError(err.message || 'Failed to fetch dashboard data');
            console.error('Error fetching dashboard:', err);
        } finally {
            setLoading(false);
        }
    };

    if (loading) {
        return (
            <div className="loading">
                Đang tải dữ liệu...
            </div>
        );
    }

    if (error) {
        return (
            <div className="container">
                <div className="error">
                    <h2>Lỗi</h2>
                    <p>{error}</p>
                    <button onClick={fetchDashboard}>Thử lại</button>
                </div>
            </div>
        );
    }

    return (
        <div>
            <div className="header">
                <div className="container">
                    <h1>🚀 AI Analytics Dashboard</h1>
                    <p>Phân tích doanh thu & Dự đoán thông minh với Machine Learning</p>
                </div>
            </div>

            <div className="container">
                {/* Restaurant Selector */}
                <div className="restaurant-selector">
                    <label htmlFor="restaurant">Chọn nhà hàng: </label>
                    <select
                        id="restaurant"
                        value={restaurantId}
                        onChange={(e) => setRestaurantId(e.target.value)}
                    >
                        <option value="REST001">Nhà hàng Phương Nam</option>
                        <option value="REST002">Quán Ăn Sài Gòn</option>
                        <option value="REST003">BBQ House</option>
                    </select>
                </div>

                {/* Summary Cards */}
                {dashboard?.summary && <SummaryCards summary={dashboard.summary} />}

                {/* Revenue Chart */}
                {dashboard?.revenue_chart && (
                    <div className="chart-container">
                        <RevenueChart data={dashboard.revenue_chart} />
                    </div>
                )}

                {/* Orders Chart */}
                {dashboard?.orders_chart && (
                    <div className="chart-container">
                        <OrdersChart data={dashboard.orders_chart} />
                    </div>
                )}

                {/* Insights */}
                {dashboard?.insights && <Insights insights={dashboard.insights} />}
            </div>
        </div>
    );
}

export default App;
