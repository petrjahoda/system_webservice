let productivityChartDom = document.getElementById('productivity-overview');
let terminalChartDom = document.getElementById('terminal-data-overview');
let calendarChartDom = document.getElementById('calendar-heatmap');
productivityChartDom.style.height = '500px'
terminalChartDom.style.height = '500px'
calendarChartDom.style.height = '250px'
let productivityChart = echarts.init(productivityChartDom);
let calendarChart = echarts.init(calendarChartDom);
let terminalChart = echarts.init(terminalChartDom);
drawProductivityChart();
drawTerminalChart();
drawCalendar();

function drawTerminalChart() {
    let option;

    let labelLeft = {
        position: 'left'
    };
    option = {
        grid: {
            top: 10,
            bottom: 10,
            left: 50,
            right: 20,
        },
        xAxis: {
            scale: true,
            responsive: true,
            type: 'value',
            position: 'top',
            splitLine: {
                lineStyle: {
                    type: 'dashed'
                }
            }
        },
        yAxis: {
            scale: true,
            responsive: true,
            type: 'category',
            axisLine: {show: false},
            axisLabel: {show: false},
            axisTick: {show: false},
            splitLine: {show: false},
            data: ['ten', 'nine', 'eight', 'seven', 'six', 'five', 'four', 'three', 'two', 'one']
        },
        series: [
            {
                type: 'bar',
                label: {
                    show: true,
                    formatter: '{b}'
                },
                data: [
                    {value: 0.01, label: labelLeft},
                    {value: 0.09, label: labelLeft},
                    {value: 0.2, label: labelLeft},
                    {value: 0.44, label: labelLeft},
                    {value: 0.23, label: labelLeft},
                    {value: 0.08, label: labelLeft},
                    {value: 0.17, label: labelLeft},
                    {value: 0.47, label: labelLeft},
                    {value: 1.17, label: labelLeft},
                    {value: 1.36, label: labelLeft},
                ]
            }
        ]
    };

    option && terminalChart.setOption(option);
}

function drawCalendar() {
    let option;
    function getVirtualData(year) {
        year = year || '2021';
        let date = +echarts.number.parseDate(year + '-01-01');
        let end = +echarts.number.parseDate((+year + 1) + '-01-01');
        let dayTime = 3600 * 24 * 1000;
        let data = [];
        for (let time = date; time < end; time += dayTime) {
            data.push([
                echarts.format.formatTime('yyyy-MM-dd', time),
                Math.floor(Math.random() * 100)
            ]);
        }
        return data;
    }
    option = {
        scale: true,
        responsive: true,
        title: {
            top: 30,
            left: 'center',
            text: '2021'
        },
        tooltip: {},
        visualMap: {
            color: 'green',
            min: 0,
            max: 100,
            type: 'piecewise',
            orient: 'horizontal',
            left: 'center',
            top: 65,
            inRange: {
                color: ['#fefefd', '#e8eed0', '#d0dda0', '#b9cc70', '#a1bb40', '#89aa10']
            }
        },
        calendar: {
            top: 120,
            left: 100,
            right: 100,
            cellSize: ['auto', 15],
            range: '2021',
            itemStyle: {
                borderWidth: 0.5
            },
            yearLabel: {show: false},
            dayLabel: {
                firstDay: 0,
                nameMap: ['Po','Ut','St','Ct','Pa','So','Ne']
            },
            monthLabel: {
                nameMap: ['Led','Uno','Bre','Dub','Kve','Cer','Cvc', 'Srp', 'Zar', 'Rij', 'Lis', 'Pro']
            }
        },
        series: {
            type: 'heatmap',
            coordinateSystem: 'calendar',
            data: getVirtualData(2021)
        }
    };

    option && calendarChart.setOption(option);
}



function drawProductivityChart() {
    productivityChart.scale = true
    let option;

    let labelLeft = {
        position: 'insideLeft'
    };
    option = {
        scale: true,
        responsive: true,
        grid: {
            top: 10,
            bottom: 10,
            left: 20,
            right: 20,
        },
        xAxis: {
            scale: true,
            responsive: true,
            type: 'value',
            position: 'top',
            splitLine: {
                lineStyle: {
                    type: 'dashed'
                }
            }
        },
        yAxis: {
            scale: true,
            responsive: true,
            type: 'category',
            axisLine: {show: false},
            axisLabel: {show: false},
            axisTick: {show: false},
            splitLine: {show: false},
            data: ['CNC-12', 'CNC-3', 'eight', 'seven', 'six', 'five', 'four', 'three', 'two', 'one']
        },
        series: [
            {
                color: '#89aa10',
                type: 'bar',
                label: {
                    show: true,
                    formatter: '{b}'
                },
                barWidth: 40,
                data: [
                    {value: 0.01, label: labelLeft, itemStyle: {color: '#fefefd', borderColor: '#e8eed0'}},
                    {value: 0.09, label: labelLeft, itemStyle: {color: '#fefefd', borderColor: '#e8eed0'}},
                    {value: 0.20, label: labelLeft, itemStyle: {color: '#e8eed0', borderColor: '#d0dda0'}},
                    {value: 0.34, label: labelLeft, itemStyle: {color: '#e8eed0', borderColor: '#d0dda0'}},
                    {value: 0.43, label: labelLeft, itemStyle: {color: '#d0dda0', borderColor: '#a1bb40'}},
                    {value: 0.48, label: labelLeft, itemStyle: {color: '#d0dda0', borderColor: '#a1bb40'}},
                    {value: 0.52, label: labelLeft, itemStyle: {color: '#d0dda0', borderColor: '#a1bb40'}},
                    {value: 0.67, label: labelLeft, itemStyle: {color: '#a1bb40', borderColor: '#89aa10'}},
                    {value: 0.87, label: labelLeft, itemStyle: {color: '#89aa10', borderColor: '#89aa10'}},
                    {value: 0.97, label: labelLeft, itemStyle: {color: '#89aa10', borderColor: '#89aa10'}},
                ]
            }
        ]
    };
    option && productivityChart.setOption(option);
}

function resizeCharts() {
    productivityChart.resize()
    calendarChart.resize();
    terminalChart.resize();
}

window.addEventListener('resize', resizeCharts)
