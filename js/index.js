let timer = 60
let timeLeft = timer;
const downloadTimer = setInterval(function () {
    if (timeLeft <= 0) {
        const workplacesElement = document.getElementsByClassName("tag short-tag");
        let workplaces = ""
        for (let index = 0; index < workplacesElement.length; index++) {
            workplaces += workplacesElement[index].children[0].innerHTML + ";"
        }
        workplaces = workplaces.slice(0, -1)
        let data = {
            key: "index-selected-workplaces",
            value: workplaces,
        };
        console.log(workplaces)
        fetch("/update_user_web_settings_from_web", {
            method: "POST",
            body: JSON.stringify(data)
        }).then(() => {
            loadIndexData();
            timeLeft = timer
        }).catch(() => {
        });
    }
    document.getElementById("progress-bar").value = timer - timeLeft;
    timeLeft -= 1;
}, 1000);

let productivityChartDom = document.getElementById('productivity-overview');
let calendarChartDom = document.getElementById('calendar-heatmap');
let terminalDowntimeChartDom = document.getElementById('terminal-downtimes');
let terminalBreakdownChartDom = document.getElementById('terminal-breakdowns');
let terminalAlarmChartDom = document.getElementById('terminal-alarms');
let consumptionChartDom = document.getElementById('consumption-chart');
consumptionChartDom.style.marginTop = "30px"
let daysDom = document.getElementById('days-chart');
let refreshButton = document.getElementById("data-refresh-button");

let productivityChart = echarts.init(productivityChartDom, null, {renderer: 'svg'});
let calendarChart = echarts.init(calendarChartDom, null, {renderer: 'svg'});
let terminalDowntimeChart = echarts.init(terminalDowntimeChartDom, null, {renderer: 'svg'});
let terminalBreakdownChart = echarts.init(terminalBreakdownChartDom, null, {renderer: 'svg'});
let terminalAlarmChart = echarts.init(terminalAlarmChartDom, null, {renderer: 'svg'});
let consumptionChart = echarts.init(consumptionChartDom, null, {renderer: 'svg'});
let daysChart = echarts.init(daysDom, null, {renderer: 'svg'})

window.addEventListener('resize', () => {
    daysChart.resize()
    consumptionChart.resize()
    terminalAlarmChart.resize()
    terminalDowntimeChart.resize()
    terminalBreakdownChart.resize()
    calendarChart.resize()
    productivityChart.resize()
})
refreshButton.addEventListener('click', () => {
    const workplacesElement = document.getElementsByClassName("tag short-tag");
    let workplaces = ""
    for (let index = 0; index < workplacesElement.length; index++) {
        workplaces += workplacesElement[index].children[0].innerHTML + ";"
    }
    workplaces = workplaces.slice(0, -1)
    let data = {
        key: "index-selected-workplaces",
        value: workplaces,
    };
    console.log(workplaces)
    fetch("/update_user_web_settings_from_web", {
        method: "POST",
        body: JSON.stringify(data)
    }).then(() => {
        loadIndexData();
        timeLeft = timer
    }).catch(() => {
    });
})

function loadIndexData() {
    console.log("LOADING INDEX DATA")
    document.getElementById("loader").hidden = false
    let data = {
        email: document.getElementById("user-info").title
    };
    fetch("/load_index_data", {
        method: "POST",
        body: JSON.stringify(data)

    }).then((response) => {
        response.text().then(function (data) {
            let result = JSON.parse(data);
            drawProductivityChart(result);
            drawCalendar(result);
            drawDaysChart(result);
            drawConsumptionChart(result);
            drawTerminalDowntimeChart(result);
            drawTerminalBreakdownChart(result);
            drawTerminalAlarmChart(result);
            document.getElementById("loader").hidden = true
            productivityChart.resize()
            calendarChart.resize();
            daysChart.resize()
            consumptionChart.resize()
            terminalDowntimeChart.resize()
            terminalBreakdownChart.resize()
            terminalAlarmChart.resize()
        });
    }).catch(() => {
        document.getElementById("loader").hidden = true
    });
}

function drawProductivityChart(data) {
    productivityChartDom.style.height = (data["WorkplaceNames"].length) * 30 + 30 + "px"
    productivityChart.scale = true
    let option;
    option = {
        textStyle: {
            fontFamily: 'ProximaNova'
        },
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
            left: 10,
            right: 20,
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
        textStyle: {
            fontFamily: 'ProximaNova'
        },
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
            left: 0,
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
        textStyle: {
            fontFamily: 'ProximaNova'
        },
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
            left: 0,
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
        textStyle: {
            fontFamily: 'ProximaNova'
        },
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
            left: 0,
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
    let locale = getLocaleFrom(data);
    moment.locale(locale);
    calendarChartDom.style.height = '250px'
    let actualWidth = (window.innerWidth - document.getElementById("navview-menu").offsetWidth) * 0.61
    if (!document.getElementById("mainmenu").classList.contains("compacted")) {
        actualWidth = (window.innerWidth - document.getElementById("navview-menu").offsetWidth) * 0.54
    }
    let option;
    option = {
        textStyle: {
            fontFamily: 'ProximaNova'
        },
        title: {
            text: data["ProductivityYearTitle"],
            x: 'center'
        },
        scale: true,
        responsive: true,
        tooltip: {
            formatter: function (param) {
                return "<b>" + moment(new Date(param["value"][0])).format('LL') + "</b><br>" + '<span style="display:inline-block;margin-right:5px;width:9px;height:9px;background-color:' + param["color"] + '"></span>' + param["value"][1] + "%"
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
                firstDay: 1,
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
    console.log(data["ConsumptionData"])
    console.log(data["MonthDataDays"])
    console.log(data["ConsumptionData"].slice(-data["MonthDataDays"].length))
    daysDom.style.height = "250px"
    let option;
    option = {
        textStyle: {
            fontFamily: 'ProximaNova'
        },
        legend: {
            bottom: "0",
            data: [data['ProductionLocale'], data['DowntimeLocale'], data['PoweroffLocale']]
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
                let result = ""
                for (const param of params) {
                    result = '<span style="display:inline-block;margin-right:5px;width:9px;height:9px;background-color:' + param["color"] + '"></span>' + "<b>" + "</b>" + param["value"] + "%<br>" + result
                }
                result = "<b>" + moment(new Date(params[0]["axisValue"])).format('LL') + "</b><br>" + result
                return result
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
            show: false,
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
    let locale = getLocaleFrom(data);
    moment.locale(locale);
    consumptionChartDom.style.height = "250px"
    let emphasisStyle = {
        itemStyle: {
            borderWidth: 1,
            borderColor: '#EE6666',
        }
    };
    let option = {
        title: {
            text: data["ConsumptionMonthTitle"],
            x: 'center'
        },
        textStyle: {
            fontFamily: 'ProximaNova'
        },
        xAxis: {
            type: 'category',
            data: data["MonthDataDays"]
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
            bottom: 50,
            left: 25,
            right: 25,
        },
        tooltip: {
            trigger: 'axis',
            axisPointer: {
                type: 'shadow'
            },
            formatter: function (params) {
                return "<b>" + moment(new Date(params[0]["axisValue"])).format('LL') + "</b><br>" + '<span style="display:inline-block;margin-right:5px;width:9px;height:9px;background-color:' + params[0]["color"] + '"></span>' + params[0]["value"] + " kWh"
            },
            position: function (point, params, dom, rect, size) {
                return [point[0] - size["contentSize"][0] / 2, point[1]];
            }
        },
        series: [{
            data: data["ConsumptionData"].slice(-data["MonthDataDays"].length),
            type: 'bar',
            symbol: 'none',
            emphasis: emphasisStyle,
            color: data["TerminalBreakdownColor"],
        }]
    };
    option && consumptionChart.setOption(option);
}

loadIndexData();