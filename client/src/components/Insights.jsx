import React from 'react';

const Insights = ({ insights }) => {
    if (!insights || insights.length === 0) {
        return null;
    }

    return (
        <div className="insights">
            <h2 className="chart-title">Nhận định & Insights</h2>
            {insights.map((insight, index) => (
                <div key={index} className={`insight-card ${insight.type}`}>
                    <div className="insight-title">{insight.title}</div>
                    <div className="insight-description">
                        {insight.description}
                        {insight.value && <strong> ({insight.value})</strong>}
                    </div>
                </div>
            ))}
        </div>
    );
};

export default Insights;
