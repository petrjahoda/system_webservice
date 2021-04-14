let timeleft = 60;
const downloadTimer = setInterval(function () {
    if (timeleft <= 0) {
        clearInterval(downloadTimer);
        loadIndexData();
    }
    document.getElementById("progress-bar").value = 60 - timeleft;
    timeleft -= 1;
}, 1000);

let productivityChartDom = document.getElementById('productivity-overview');
let calendarChartDom = document.getElementById('calendar-heatmap');
let terminalDowntimeChartDom = document.getElementById('terminal-downtimes');
let terminalBreakdownChartDom = document.getElementById('terminal-breakdowns');
let terminalAlarmChartDom = document.getElementById('terminal-alarms');
let consumptionChartDom = document.getElementById('consumption-chart');
let daysDom = document.getElementById('days-chart');
let refreshButton = document.getElementById("data-refresh-button");

let productivityChart = echarts.init(productivityChartDom);
let calendarChart = echarts.init(calendarChartDom);
let terminalDowntimeChart = echarts.init(terminalDowntimeChartDom);
let terminalBreakdownChart = echarts.init(terminalBreakdownChartDom);
let terminalAlarmChart = echarts.init(terminalAlarmChartDom);
let consumptionChart = echarts.init(consumptionChartDom);
let daysChart = echarts.init(daysDom)

loadIndexData();

window.addEventListener('resize', () => location.reload())

refreshButton.addEventListener('click', () => {
    const workplacesElement = document.getElementsByClassName("tag short-tag");
    let workplaces = []
    for (let index = 0; index < workplacesElement.length; index++) {
        workplaces.push(workplacesElement[index].children[0].innerHTML)
    }
    let data = {
        workplaces: workplaces,
    };
    fetch("/update_user_workplaces", {
        method: "POST",
        body: JSON.stringify(data)
    }).then(() => {
        loadIndexData();
        timeleft=60
    }).catch((error) => {
        console.log(error)
    });
})

function loadIndexData() {
    fetch("/load_index_data", {
        method: "POST",
    }).then((response) => {
        response.text().then(function (data) {
            let result = JSON.parse(data);
            drawProductivityChart(result);
            productivityChart.resize()
            drawCalendar(result);
            calendarChart.resize();
            drawDaysChart(result);
            daysChart.resize()
            drawConsumptionChart(result);
            consumptionChart.resize()
            drawTerminalDowntimeChart(result);
            terminalDowntimeChart.resize()
            drawTerminalBreakdownChart(result);
            terminalBreakdownChart.resize()
            drawTerminalAlarmChart(result);
            terminalAlarmChart.resize()
        });
    }).catch((error) => {
        console.log(error)
    });
}

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
                return params[0]["name"] + ": <b>" + params[0]["data"].toFixed(1) + "%</b>"
            },
            position: function (point, params, dom, rect, size) {
                return [point[0] - size["contentSize"][0] / 2, point[1]];
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
    if (data["TerminalDowntimeNames"] === null) {
        terminalDowntimeChartDom.style.height = "0px"
    } else {
        terminalDowntimeChartDom.style.height = (data["TerminalDowntimeNames"].length) * 30 + 30 + "px"
    }
    let option;
    option = {
        tooltip: {
            trigger: 'axis',
            axisPointer: {
                type: 'shadow'
            },
            formatter: function (params) {
                return params[0]["name"] + ": <b>" + moment.duration(params[0]["data"], "seconds").locale("cs").humanize() + "</b>"
            },
            position: function (point, params, dom, rect, size) {
                return [point[0] - size["contentSize"][0] / 2, point[1]];
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
    if (data["TerminalDowntimeNames"] === null) {
        document.getElementById("terminal-breakdowns").style.marginTop = "0px"
    } else {
        document.getElementById("terminal-breakdowns").style.marginTop = "30px"
    }
    if (data["TerminalBreakdownNames"] === null) {
        terminalBreakdownChartDom.style.height = "0px"
    } else {
        terminalBreakdownChartDom.style.height = data["TerminalBreakdownNames"].length * 30 + 30 + "px"
    }
    let option;
    option = {
        tooltip: {
            trigger: 'axis',
            axisPointer: {
                type: 'shadow'
            },
            formatter: function (params) {
                return params[0]["name"] + ": <b>" + moment.duration(params[0]["data"], "seconds").locale("cs").humanize() + "</b>"
            },
            position: function (point, params, dom, rect, size) {
                return [point[0] - size["contentSize"][0] / 2, point[1]];
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
    if (data["TerminalDowntimeNames"] === null && data["TerminalBreakdownNames"] === null) {
        document.getElementById("terminal-alarms").style.marginTop = "0px"
    } else {
        document.getElementById("terminal-alarms").style.marginTop = "30px"
    }
    if (data["TerminalAlarmNames"] === null) {
        terminalAlarmChartDom.style.height = "0px"
    } else {
        terminalAlarmChartDom.style.height = data["TerminalAlarmNames"].length * 30 + 30 + "px"
    }
    let option;
    option = {
        tooltip: {
            trigger: 'axis',
            axisPointer: {
                type: 'shadow'
            },
            formatter: function (params) {
                return params[0]["name"] + ": <b>" + moment.duration(params[0]["data"], "seconds").locale("cs").humanize() + "</b>"
            },
            position: function (point, params, dom, rect, size) {
                return [point[0] - size["contentSize"][0] / 2, point[1]];
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
    let actualWidth = parseInt(window.getComputedStyle(document.getElementById("inside")).width)
    let option;
    option = {
        title: {
            text: data["ProductivityYearTitle"],
            x: 'center'
        },
        scale: true,
        responsive: true,
        tooltip: {
            formatter: function (p) {
                let format = echarts.format.formatTime('dd.MM.yyyy', p.data[0]);
                return format + ': ' + "<b>" + p.data[1] + "%</b>";
            },
            position: function (point, params, dom, rect, size) {
                return [point[0] - size["contentSize"][0] / 2, point[1]];
            }
        },
        visualMap: {
            min: 0,
            max: 100,
            calculable: true,
            orient: 'horizontal',
            left: 'center',
            top: 72 + actualWidth / 53 * 6,
            inRange: {
                color: ['#f3f6e7', '#89aa10']
            },
            formatter: function (value) {
                return +value.toFixed(1) + '%';
            }
        },
        calendar: {
            top: 50,
            left: 30,
            right: 30,
            cellSize: ['auto', actualWidth / 53],
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

function drawDaysChart(data) {
    daysDom.style.height = "250px"
    let option;

    option = {
        legend: {
            bottom: "0"
        },
        title: {
            text: data["OverviewMonthTitle"],
            x: 'center'
        },
        tooltip: {
            trigger: 'axis',
            axisPointer: {
                type: 'none'
            },
            formatter: function (params) {
                return "<b>"+echarts.format.formatTime('dd.MM.yyyy', params[0]["axisValue"]) +"</b><br>"+params[2]["seriesName"] +": <b>" + params[2]["data"]+"%</b><br>"+ params[1]["seriesName"] +": <b>" + params[1]["data"]+"%</b><br>"+ params[0]["seriesName"] +": <b>" + params[0]["data"]+"%</b><br>"
            },
            position: function (point, params, dom, rect, size) {
                return [point[0] - size["contentSize"][0] / 2, point[1]];
            }
        },
        xAxis: {
            data: data['MonthDataDays'],
            axisLine: {onZero: true},
            splitLine: {show: false},
            splitArea: {show: false}
        },
        yAxis: {
            show:false,
            axisLine: {onZero: true},
            splitLine: {show: false},
            splitArea: {show: false},
            max: 100,
        },
        grid: {
            top: 30,
            bottom: 50,
            left: 25,
            right: 25,
        },
        series: [
            {
                name: data['PoweroffLocale'],
                type: 'bar',
                stack: 'one',
                data: data['MonthDataPoweroff'],
                color: data['TerminalBreakdownColor']
            },

            {
                name: data['DowntimeLocale'],
                type: 'bar',
                stack: 'one',
                data: data['MonthDataDowntime'],
                color: data["TerminalDowntimeColor"]
            },
            {
                name: data['ProductionLocale'],
                type: 'bar',
                stack: 'one',
                data: data['MonthDataProduction'],
                color: data['TerminalProductionColor']
            },

        ]
    };
    option && daysChart.setOption(option);
}

function drawConsumptionChart(data) {
    consumptionChartDom.style.height = "300px"
    let emphasisStyle = {
        itemStyle: {
            borderWidth: 1,
            borderColor: '#EE6666',
        }
    };
    option = {
        xAxis: {
            type: 'category',
            data: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun','Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun','Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun','Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun', 'Sat', 'Sun']
        },
        yAxis: {
            type: 'value',
            show: false,
            axisLine: {onZero: true},
            splitLine: {show: false},
            splitArea: {show: false},
        },
        grid: {
            top: 30,
            bottom: 30,
            left: 25,
            right: 25,
        },
        tooltip: {
        },
        series: [{
            data: [820, 932, 901, 934, 1290, 1330, 1320,820, 932, 901, 934, 1290, 1330, 1320,820, 932, 901, 934, 1290, 1330, 1320,820, 932, 901, 934, 1290, 1330, 1320, 123, 321],
            type: 'bar',
            symbol: 'none',
            emphasis: emphasisStyle,
            color: data["TerminalBreakdownColor"],
        }]
    };
    option && consumptionChart.setOption(option);
}
