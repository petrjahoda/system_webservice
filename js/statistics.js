let calendarChartDom = document.getElementById('calendar-chart');
let calendarChart = echarts.init(calendarChartDom, null, {renderer: 'svg'});
let firstUpperChartDom = document.getElementById('first-upper-chart');
let firstUpperChart = echarts.init(firstUpperChartDom, null, {renderer: 'svg'});
let secondUpperChartDom = document.getElementById('seconds-upper-chart');
let secondUpperChart = echarts.init(secondUpperChartDom, null, {renderer: 'svg'});
let thirdUpperChartDom = document.getElementById('third-upper-chart');
let thirdUpperChart = echarts.init(thirdUpperChartDom, null, {renderer: 'svg'});
let fourthUpperChartDom = document.getElementById('fourth-upper-chart');
let fourthUpperChart = echarts.init(fourthUpperChartDom, null, {renderer: 'svg'});

if (document.getElementById('to-date').value === "") {
}
let now = new Date();
now.setMinutes(now.getMinutes() - now.getTimezoneOffset());
document.getElementById('to-date').value = now.toISOString().slice(0, 16);
now.setMonth(now.getMonth() - 1);
document.getElementById('from-date').value = now.toISOString().slice(0, 16);

const dataSelection = document.getElementById("statistics-selection")
dataSelection.addEventListener("change", () => {
    let data = {
        email: document.getElementById("user-info").title,
        selection: document.getElementById("statistics-selection").value,
    };
    fetch("/load_types_for_selected_statistics", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            document.getElementById("types").innerHTML = data
        });
        loadStatisticsChartData()
    }).catch(() => {
    });
})

const statisticsRefreshButton = document.getElementById("statistics-refresh-button")
statisticsRefreshButton.addEventListener("click", () => {
    document.getElementById("loader").hidden = false
    statisticsRefreshButton.disabled = true
    const workplacesElement = document.getElementById("workplaces").getElementsByClassName("tag short-tag");
    let workplaces = ""
    for (let index = 0; index < workplacesElement.length; index++) {
        workplaces += workplacesElement[index].children[0].innerHTML + ";"
    }
    workplaces = workplaces.slice(0, -1)
    let data = {
        email: document.getElementById("user-info").title,
        key: "statistics-selected-workplaces",
        value: workplaces,
    };
    fetch("/update_user_web_settings_from_web", {
        method: "POST",
        body: JSON.stringify(data)
    }).then(() => {
        let typesElement = document.getElementById("types").getElementsByClassName("tag short-tag")
        let types = ""
        for (let index = 0; index < typesElement.length; index++) {
            types += typesElement[index].children[0].innerHTML + ";"
        }
        types = types.slice(0, -1)
        let data = {
            email: document.getElementById("user-info").title,
            key: "statistics-selected-types-" + document.getElementById("statistics-selection").value,
            value: types,
        };
        fetch("/update_user_web_settings_from_web", {
            method: "POST",
            body: JSON.stringify(data)
        }).then(() => {
            let usersElement = document.getElementById("users").getElementsByClassName("tag short-tag")
            let users = ""
            for (let index = 0; index < usersElement.length; index++) {
                users += usersElement[index].children[0].innerHTML + ";"
            }
            users = users.slice(0, -1)
            let data = {
                email: document.getElementById("user-info").title,
                key: "statistics-selected-users",
                value: users,
            };
            fetch("/update_user_web_settings_from_web", {
                method: "POST",
                body: JSON.stringify(data)
            }).then(() => {
                loadStatisticsChartData();
                document.getElementById("loader").hidden = true
                statisticsRefreshButton.disabled = false
            }).catch(() => {
                document.getElementById("loader").hidden = true
                statisticsRefreshButton.disabled = false
            });
        }).catch(() => {
            document.getElementById("loader").hidden = true
            statisticsRefreshButton.disabled = false
        });
    }).catch(() => {
        document.getElementById("loader").hidden = true
        statisticsRefreshButton.disabled = false
    });
})


function drawCalendarChart(data) {
    let locale = getLocaleFrom(data);
    moment.locale(locale);
    calendarChartDom.style.height = "250px"
    let emphasisStyle = {
        itemStyle: {
            borderWidth: 1,
            borderColor: '#EE6666',
        }
    };
    let option = {
        title: {
            text: data["CalendarChartLocale"],
            x: 'center'
        },
        textStyle: {
            fontFamily: 'ProximaNova'
        },
        xAxis: {
            type: 'category',
            data: data["DaysChartData"]
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
            right: 50,
        },
        tooltip: {
            trigger: 'axis',
            axisPointer: {
                type: 'shadow'
            },
            formatter: function (params) {
                return "<b>" + moment(new Date(params[0]["axisValue"])).format('LL') + "</b><br>" + '<span style="display:inline-block;margin-right:5px;width:9px;height:9px;background-color:' + params[0]["color"] + '"></span>' + convertDurationToString(params[0]["value"])
            },
            position: function (point, params, dom, rect, size) {
                return [point[0] - size["contentSize"][0] / 2, point[1]];
            }
        },
        series: [{
            data: data["DaysChartValue"],
            type: 'bar',
            symbol: 'none',
            emphasis: emphasisStyle,
        }]
    };
    option && calendarChart.setOption(option);
}

function loadStatisticsChartData() {
    document.getElementById("loader").hidden = false
    document.getElementById("res")
    let data = {
        data: document.getElementById("statistics-selection").value,
        from: document.getElementById("from-date").value,
        to: document.getElementById("to-date").value,
        email: document.getElementById("user-info").title
    };
    fetch("/load_statistics_data", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            try {
                let result = JSON.parse(data);
                console.log(result)
                drawFirstUpperChart(result)
                drawSecondUpperChart(result)
                drawThirdUpperChart(result)
                drawFourthUpperChart(result)
                drawCalendarChart(result)
                calendarChart.resize()
                firstUpperChart.resize()
                secondUpperChart.resize()
                thirdUpperChart.resize()
                fourthUpperChart.resize()

            } catch {
            }
            document.getElementById("loader").hidden = true
        });
    }).catch(() => {
        document.getElementById("loader").hidden = true
    });
}

function drawFourthUpperChart(data) {
    if (data["TimeChartData"] === null) {
        fourthUpperChartDom.style.height = "0px"
    } else {
        fourthUpperChartDom.style.height = (data["TimeChartData"].length) * 30 + 30 + "px"
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
                return params[0]["name"] + ": <b>" + data["TimeChartText"][params[0]["dataIndex"]] + "</b>"
            },
            position: function (point, params, dom, rect, size) {
                return [point[0] - size["contentSize"][0] / 2, point[1]];
            }

        },
        title: {
            text: data["FourthUpperChartLocale"]
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
            data: data["TimeChartData"]
        },
        series: [
            {
                // color: data["TerminalDowntimeColor"],
                barWidth: 25,
                type: 'bar',
                silent: true,
                label: {
                    show: true,
                    formatter: '{b}',
                    position: 'insideLeft',
                    fontSize: '16',
                },
                data: data["TimeChartValue"]
            }
        ]
    };
    option && fourthUpperChart.setOption(option);
}

function drawThirdUpperChart(data) {
    if (data["UsersChartData"] === null) {
        thirdUpperChartDom.style.height = "0px"
    } else {
        thirdUpperChartDom.style.height = (data["UsersChartData"].length) * 30 + 30 + "px"
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
                return params[0]["name"] + ": <b>" + data["UsersChartText"][params[0]["dataIndex"]] + "</b>"
            },
            position: function (point, params, dom, rect, size) {
                return [point[0] - size["contentSize"][0] / 2, point[1]];
            }

        },
        title: {
            text: data["ThirdUpperChartLocale"]
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
            data: data["UsersChartData"]
        },
        series: [
            {
                // color: data["TerminalDowntimeColor"],
                barWidth: 25,
                type: 'bar',
                silent: true,
                label: {
                    show: true,
                    formatter: '{b}',
                    position: 'insideLeft',
                    fontSize: '16',
                },
                data: data["UsersChartValue"]
            }
        ]
    };
    option && thirdUpperChart.setOption(option);
}

function drawSecondUpperChart(data) {
    if (data["SelectionChartData"] === null) {
        secondUpperChartDom.style.height = "0px"
    } else {
        secondUpperChartDom.style.height = (data["SelectionChartData"].length) * 30 + 30 + "px"
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
                return params[0]["name"] + ": <b>" + data["SelectionChartText"][params[0]["dataIndex"]] + "</b>"
            },
            position: function (point, params, dom, rect, size) {
                return [point[0] - size["contentSize"][0] / 2, point[1]];
            }

        },
        title: {
            text: data["SecondUpperChartLocale"]
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
            data: data["SelectionChartData"]
        },
        series: [
            {
                // color: data["TerminalDowntimeColor"],
                barWidth: 25,
                type: 'bar',
                silent: true,
                label: {
                    show: true,
                    formatter: '{b}',
                    position: 'insideLeft',
                    fontSize: '16',
                },
                data: data["SelectionChartValue"]
            }
        ]
    };
    option && secondUpperChart.setOption(option);
}

function drawFirstUpperChart(data) {
    if (data["WorkplaceChartData"] === null) {
        firstUpperChartDom.style.height = "0px"
    } else {
        firstUpperChartDom.style.height = (data["WorkplaceChartData"].length) * 30 + 30 + "px"
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
                return params[0]["name"] + ": <b>" + data["WorkplaceChartText"][params[0]["dataIndex"]] + "</b>"
            },
            position: function (point, params, dom, rect, size) {
                return [point[0] - size["contentSize"][0] / 2, point[1]];
            }

        },
        title: {
            text: data["FirstUpperChartLocale"]
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
            data: data["WorkplaceChartData"]
        },
        series: [
            {
                // color: data["TerminalDowntimeColor"],
                barWidth: 25,
                type: 'bar',
                silent: true,
                label: {
                    show: true,
                    formatter: '{b}',
                    position: 'insideLeft',
                    fontSize: '16',
                },
                data: data["WorkplaceChartValue"]
            }
        ]
    };
    option && firstUpperChart.setOption(option);
}

function addDate(dt, amount, dateType) {
    switch (dateType) {
        case 'days':
            return dt.setDate(dt.getDate() + amount) && dt;
        case 'weeks':
            return dt.setDate(dt.getDate() + (7 * amount)) && dt;
        case 'months':
            return dt.setMonth(dt.getMonth() + amount) && dt;
        case 'years':
            return dt.setFullYear(dt.getFullYear() + amount) && dt;
    }
}

window.addEventListener('resize', () => {
    firstUpperChart.resize()
    secondUpperChart.resize()
    thirdUpperChart.resize()
    fourthUpperChart.resize()
    calendarChart.resize()
})

function convertDurationToString(value) {
    const sec = parseInt(value, 10); // convert value to number if it's string
    let hours = Math.floor(sec / 3600); // get hours
    let minutes = Math.floor((sec - (hours * 3600)) / 60); // get minutes
    let seconds = sec - (hours * 3600) - (minutes * 60); //  get seconds
    if (hours < 10) {
        hours = "0" + hours;
    }
    if (minutes < 10) {
        minutes = "0" + minutes;
    }
    if (seconds < 10) {
        seconds = "0" + seconds;
    }
    return hours + 'h' + minutes + 'h' + seconds + 's'; // Return is HH : MM : SS
}