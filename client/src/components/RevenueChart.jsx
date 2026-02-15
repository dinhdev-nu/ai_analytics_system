import React from 'react';
import ReactECharts from 'echarts-for-react';

const RevenueChart = ({ data }) => {
    if (!data || !data.labels || data.labels.length === 0) {
        return <div>Không có dữ liệu</div>;
    }

    const option = {
        title: {
            text: 'Dự báo Doanh thu',
            textStyle: {
                fontSize: 20,
                fontWeight: 600,
            },
        },
        tooltip: {
            trigger: 'axis',
            axisPointer: {
                type: 'cross',
            },
            formatter: (params) => {
                let result = `<strong>${params[0].axisValue}</strong><br/>`;
                params.forEach((item) => {
                    const value = item.value ? item.value.toLocaleString('vi-VN') : '0';
                    result += `${item.marker} ${item.seriesName}: ${value} VND<br/>`;
                });
                return result;
            },
        },
        legend: {
            data: ['Thực tế', 'Dự đoán', 'Mục tiêu'],
            top: 40,
        },
        grid: {
            left: '3%',
            right: '4%',
            bottom: '3%',
            containLabel: true,
        },
        xAxis: {
            type: 'category',
            data: data.labels,
            boundaryGap: false,
            axisLabel: {
                rotate: 45,
            },
        },
        yAxis: {
            type: 'value',
            axisLabel: {
                formatter: (value) => {
                    return (value / 1000000).toFixed(0) + 'M';
                },
            },
        },
        series: [
            {
                name: 'Thực tế',
                type: 'line',
                data: data.actual,
                lineStyle: {
                    width: 3,
                },
                itemStyle: {
                    color: '#3b82f6',
                },
                smooth: true,
            },
            {
                name: 'Dự đoán',
                type: 'line',
                data: data.predicted,
                lineStyle: {
                    width: 3,
                    type: 'dashed',
                },
                itemStyle: {
                    color: '#10b981',
                },
                smooth: true,
            },
            {
                name: 'Mục tiêu',
                type: 'line',
                data: data.target,
                lineStyle: {
                    width: 2,
                    type: 'dotted',
                },
                itemStyle: {
                    color: '#f59e0b',
                },
                smooth: true,
            },
        ],
    };

    return <ReactECharts option={option} style={{ height: '400px' }} />;
};

export default RevenueChart;
