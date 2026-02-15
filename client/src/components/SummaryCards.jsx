import React from 'react';

const SummaryCards = ({ summary }) => {
    if (!summary) return null;

    const formatMoney = (value) => {
        return new Intl.NumberFormat('vi-VN', {
            style: 'currency',
            currency: 'VND',
        }).format(value);
    };

    const formatPercent = (value) => {
        const sign = value >= 0 ? '+' : '';
        return `${sign}${value.toFixed(1)}%`;
    };

    return (
        <div className="dashboard">
            <div className="card">
                <div className="card-title">Doanh thu tháng này</div>
                <div className="card-value">{formatMoney(summary.current_month_revenue)}</div>
                <div className={`card-change ${summary.month_over_month_growth >= 0 ? 'positive' : 'negative'}`}>
                    <span>{formatPercent(summary.month_over_month_growth)}</span>
                    <span>so với tháng trước</span>
                </div>
            </div>

            <div className="card">
                <div className="card-title">Tăng trưởng YoY</div>
                <div className="card-value">{formatPercent(summary.year_over_year_growth)}</div>
                <div className="card-change">
                    <span>So với cùng kỳ năm trước</span>
                </div>
            </div>

            <div className="card">
                <div className="card-title">Tổng đơn hàng</div>
                <div className="card-value">{summary.total_orders.toLocaleString()}</div>
                <div className="card-change">
                    <span>Giá trị TB: {formatMoney(summary.avg_order_value)}</span>
                </div>
            </div>

            <div className="card">
                <div className="card-title">Dự báo tháng tới</div>
                <div className="card-value">{formatMoney(summary.forecast_next_month)}</div>
                <div className="card-change">
                    <span>Độ tin cậy: {(summary.forecast_confidence * 100).toFixed(0)}%</span>
                </div>
            </div>
        </div>
    );
};

export default SummaryCards;
