let productivityChartDom = document.getElementById('productivity-overview');
let calendarChartDom = document.getElementById('calendar-heatmap');
let terminalDowntimeChartDom = document.getElementById('terminal-downtimes');
let terminalBreakdownChartDom = document.getElementById('terminal-breakdowns');
let terminalAlarmChartDom = document.getElementById('terminal-alarms');
productivityChartDom.hidden = true
calendarChartDom.hidden = true
terminalDowntimeChartDom.hidden = true
terminalBreakdownChartDom.hidden = true
terminalAlarmChartDom.hidden = true
let productivityChart = echarts.init(productivityChartDom);
let calendarChart = echarts.init(calendarChartDom);
let terminalDowntimeChart = echarts.init(terminalDowntimeChartDom);
let terminalBreakdownChart = echarts.init(terminalBreakdownChartDom);
let terminalAlarmChart = echarts.init(terminalAlarmChartDom);

fetch("/load_index_data", {
    method: "POST",
}).then((response) => {
    response.text().then(function (data) {
        let result = JSON.parse(data);
        drawProductivityChart(result);
        drawCalendar(result);
        if (result["TerminalDowntimeNames"] !== null) {
            drawTerminalDowntimeChart(result);
        }
        if (result["TerminalBreakdownNames"] !== null) {
            drawTerminalBreakdownChart(result);
        }
        if (result["TerminalAlarmNames"] !== null) {
            drawTerminalAlarmChart(result);
        }
        resizeCharts()
        productivityChartDom.hidden = false
        terminalDowntimeChartDom.hidden = false
        terminalBreakdownChartDom.hidden = false
        terminalAlarmChartDom.hidden = false
        calendarChartDom.hidden = false
        window.addEventListener('resize', resizeCharts)
    });
}).catch((error) => {
    console.log(error)
});

function drawProductivityChart(data) {
    productivityChartDom.style.height = (data["WorkplaceNames"].length) * 30 + 30 + "px"
    productivityChart.scale = true
    let option;
    option = {
        tooltip: {
            trigger: 'axis',
            axisPointer: {
                type: 'shadow'
            },
            formatter: function (params) {
                return params[0]["name"] +": <b>"+ params[0]["data"].toFixed(1) + "%</b>"
            },
            position: function (point, params, dom, rect, size) {
                return [point[0]-size["contentSize"][0]/2, point[1]];
            }

        },
        title: {
            text: data["ProductivityTodayTitle"]
        },
        scale: true,
        responsive: true,
        grid: {
            top: 30,
            bottom: 0,
            left: 5,
            right: 1,
        },
        xAxis: {
            scale: true,
            min: 0,
            max: 100,
            responsive: true,
            type: 'value',
            position: 'top',
            axisLine: {show: false},
            axisLabel: {show: false},
            axisTick: {show: false},
            splitLine: {
                lineStyle: {
                    type: 'dashed',
                    color: "#e5e5e5"
                }
            },
        },
        yAxis: {
            scale: true,
            responsive: true,
            type: 'category',
            axisLine: {show: false},
            axisLabel: {show: false},
            axisTick: {show: false},
            splitLine: {show: false},

            data: data["WorkplaceNames"]
        },
        series: [
            {
                color: data["TerminalProductionColor"],
                type: 'bar',
                silent: true,
                label: {
                    show: true,
                    // var res = str.split(" ");
                    formatter: '{b}',
                    position: 'insideLeft',
                    fontSize: '16',
                },
                barWidth: 25,
                data: data["WorkplacePercents"],

            }
        ]
    };
    option && productivityChart.setOption(option);
}

function drawTerminalDowntimeChart(data) {
    terminalDowntimeChartDom.style.height = (data["TerminalDowntimeNames"].length) * 30 + 30 + "px"
    console.log(data["Locale"]) //change locale
    let option;
    option = {
        tooltip: {
            trigger: 'axis',
            axisPointer: {
                type: 'shadow'
            },
            formatter: function (params) {
                return params[0]["name"] +": <b>"+ moment.duration(params[0]["data"], "seconds").locale("cs").humanize() + "</b>"
            },
            position: function (point, params, dom, rect, size) {
                return [point[0]-size["contentSize"][0]/2, point[1]];
            }

        },
        title: {
            text: data["DowntimesTitle"]
        },
        grid: {
            top: 30,
            bottom: 0,
            left: 1,
            right: 20,
        },
        xAxis: {
            min: 0,
            scale: true,
            responsive: true,
            type: 'value',
            position: 'bottom',
            axisLine: {show: false},
            axisLabel: {show: false},
            axisTick: {show: false},
            splitLine: {
                lineStyle: {
                    type: 'dashed',
                    color: "#e5e5e5"
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
            data: data["TerminalDowntimeNames"]
        },
        series: [
            {
                color: data["TerminalDowntimeColor"],
                barWidth: 25,
                type: 'bar',
                silent: true,
                label: {
                    show: true,
                    formatter: '{b}',
                    position: 'insideLeft',
                    fontSize: '16',
                },
                data: data["TerminalDowntimeDurations"]
            }
        ]
    };
    option && terminalDowntimeChart.setOption(option);
}

function drawTerminalBreakdownChart(data) {
    terminalBreakdownChartDom.style.height = data["TerminalBreakdownNames"].length * 30 + 30 + "px"
    let option;
    option = {
        tooltip: {
            trigger: 'axis',
            axisPointer: {
                type: 'shadow'
            },
            formatter: function (params) {
                return params[0]["name"] +": <b>"+ moment.duration(params[0]["data"], "seconds").locale("cs").humanize() + "</b>"
            },
            position: function (point, params, dom, rect, size) {
                return [point[0]-size["contentSize"][0]/2, point[1]];
            }

        },
        title: {
            text: data["BreakdownsTitle"]
        },
        grid: {
            top: 30,
            bottom: 0,
            left: 1,
            right: 20,
        },
        xAxis: {
            min: 0,
            scale: true,
            responsive: true,
            type: 'value',
            position: 'top',
            axisLine: {show: false},
            axisLabel: {show: false},
            axisTick: {show: false},
            splitLine: {
                lineStyle: {
                    type: 'dashed',
                    color: "#e5e5e5"
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
            data: data["TerminalBreakdownNames"]
        },
        series: [
            {
                color: data["TerminalBreakdownColor"],
                barWidth: 25,
                type: 'bar',
                silent: true,
                label: {
                    show: true,
                    formatter: '{b}',
                    position: 'insideLeft',
                    fontSize: '16',
                },
                data: data["TerminalBreakdownDurations"]
            }
        ]
    };
    option && terminalBreakdownChart.setOption(option);
}

function drawTerminalAlarmChart(data) {
    terminalAlarmChartDom.style.height = data["TerminalAlarmNames"].length * 30 + 30 + "px"
    let option;
    option = {
        tooltip: {
            trigger: 'axis',
            axisPointer: {
                type: 'shadow'
            },
            formatter: function (params) {
                return params[0]["name"] +": <b>"+ moment.duration(params[0]["data"], "seconds").locale("cs").humanize() + "</b>"
            },
            position: function (point, params, dom, rect, size) {
                return [point[0]-size["contentSize"][0]/2, point[1]];
            }

        },
        title: {
            text: data["AlarmsTitle"]
        },
        grid: {
            top: 30,
            bottom: 0,
            left: 1,
            right: 20,
        },
        xAxis: {
            min: 0,
            scale: true,
            responsive: true,
            type: 'value',
            position: 'top',
            axisLine: {show: false},
            axisLabel: {show: false},
            axisTick: {show: false},
            splitLine: {
                lineStyle: {
                    type: 'dashed',
                    color: "#e5e5e5"
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
            data: data["TerminalAlarmNames"]
        },
        series: [
            {
                color: data["TerminalAlarmColor"],
                barWidth: 25,
                type: 'bar',
                silent: true,
                label: {
                    show: true,
                    formatter: '{b}',
                    position: 'insideLeft',
                    fontSize: '16',
                },
                data: data["TerminalAlarmDurations"]
            }
        ]
    };
    option && terminalAlarmChart.setOption(option);
}

function drawCalendar(data) {
    calendarChartDom.style.height = '250px'
    let option;
    option = {
        scale: true,
        responsive: true,
        tooltip: {
            formatter: function (p) {
                let format = echarts.format.formatTime('dd.MM.yyyy', p.data[0]);
                return format + ': ' + "<b>" + p.data[1] + "%</b>";
            },
            position: function (point, params, dom, rect, size) {
                return [point[0]-size["contentSize"][0]/2, point[1]];
            }
        },
        visualMap: {
            min: 0,
            max: 100,
            calculable: true,
            orient: 'horizontal',
            left: 'center',
            top: '150',
            inRange: {
                color: ['#f3f6e7', '#89aa10']
            },
            formatter: function (value) {
                return +value.toFixed(1) + '%';
            }
        },
        calendar: {
            top: 24,
            left: 30,
            right: 30,
            cellSize: ['auto', 18],
            range: [data["CalendarStart"], data["CalendarEnd"]],
            itemStyle: {
                borderWidth: 0.5
            },
            yearLabel: {show: false},
            dayLabel: {
                firstDay: 0,
                nameMap: data["CalendarDayLabel"]
            },
            monthLabel: {
                nameMap: data["CalendarMonthLabel"]
            }
        },
        series: {
            type: 'heatmap',
            coordinateSystem: 'calendar',
            data: data["CalendarData"],
            emphasis: {
                itemStyle: {
                    borderWidth: 1,
                    borderColor: '#EE6666',
                }
            }
        }
    };
    option && calendarChart.setOption(option);
}

function resizeCharts() {
    setTimeout(function () {
        productivityChart.resize()
        calendarChart.resize();
        terminalDowntimeChart.resize()
        terminalBreakdownChart.resize()
        terminalAlarmChart.resize()
        terminalDowntimeChartDom.hidden = false
        terminalBreakdownChartDom.hidden = false
        terminalAlarmChartDom.hidden = false
        calendarChartDom.hidden = false
        productivityChartDom.hidden = false
    }, 250);
}