import React from 'react';
import ReactECharts from 'echarts-for-react';

const OrdersChart = ({ data }) => {
    if (!data || !data.labels || data.labels.length === 0) {
        return <div>Không có dữ liệu</div>;
    }

    const option = {
        title: {
            text: 'Số lượng Đơn hàng',
            textStyle: {
                fontSize: 20,
                fontWeight: 600,
            },
        },
        tooltip: {
            trigger: 'axis',
            axisPointer: {
                type: 'shadow',
            },
            formatter: (params) => {
                const item = params[0];
                return `<strong>${item.axisValue}</strong><br/>${item.marker} ${item.seriesName}: ${item.value} đơn`;
            },
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
            axisLabel: {
                rotate: 45,
            },
        },
        yAxis: {
            type: 'value',
        },
        series: [
            {
                name: 'Số đơn hàng',
                type: 'bar',
                data: data.order_counts,
                itemStyle: {
                    color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                        { offset: 0, color: '#667eea' },
                        { offset: 1, color: '#764ba2' },
                    ]),
                    borderRadius: [8, 8, 0, 0],
                },
                barWidth: '60%',
            },
        ],
    };

    return <ReactECharts option={option} style={{ height: '300px' }} />;
};

export default OrdersChart;
