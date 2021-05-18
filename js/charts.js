let chartDom = document.getElementById('chart');
let chartHeight = document.documentElement.clientHeight * 0.85
if (chartHeight < 800) {
    chartHeight = 800;
}
let myChart = echarts.init(chartDom, null, {height: chartHeight});
let now = new Date();
now.setMinutes(now.getMinutes() - now.getTimezoneOffset());
document.getElementById('to-date').value = now.toISOString().slice(0, 16);
now.setHours(now.getHours() - 24);
document.getElementById('from-date').value = now.toISOString().slice(0, 16);


const flashButton = document.getElementById("flash-button")
const phonelinkButton = document.getElementById("phonelink-button")
flashButton.addEventListener("click", (event) => {
    if (flashButton.classList.contains("mif-flash-on")) {
        flashButton.classList.remove("mif-flash-on")
        flashButton.classList.add("mif-flash-off")
    } else {
        flashButton.classList.remove("mif-flash-off")
        flashButton.classList.add("mif-flash-on")
    }
})
phonelinkButton.addEventListener("click", (event) => {
    if (phonelinkButton.classList.contains("mif-phonelink-off")) {
        phonelinkButton.classList.remove("mif-phonelink-off")
        phonelinkButton.classList.add("mif-phonelink")
    } else {
        phonelinkButton.classList.remove("mif-phonelink")
        phonelinkButton.classList.add("mif-phonelink-off")
    }
})

phonelinkButton.addEventListener("click", (event) => {
})


const dataOkButton = document.getElementById("data-ok-button")
dataOkButton.addEventListener("click", (event) => {
    document.getElementById("loader").hidden = false
    console.log("getting chart data for " + document.getElementById("data-selection").value)
    console.log("getting chart data for " + document.getElementById("workplace-selection").value)
    let data = {
        data: document.getElementById("data-selection").value,
        workplace: document.getElementById("workplace-selection").value,
        from: document.getElementById("from-date").value,
        to: document.getElementById("to-date").value
    };
    const start = performance.now();
    fetch("/load_chart_data", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        const download = performance.now();
        let difference = download - start
        if (difference < 1000) {
            updateCharm("INF: Chart data downloaded from database in " + difference + "ms")
        } else {
            updateCharm("INF: Chart data downloaded from database in " + difference / 1000 + "s")
        }
        response.text().then(function (data) {
            let result = JSON.parse(data);
            if (result["Type"] === "analog-data") {
                myChart.clear()
                if (result["ChartData"] !== null) {
                    drawAnalogChart(result)
                }
                document.getElementById("loader").hidden = true
            } else if (result["Type"] === "digital-data") {
                myChart.clear()
                if (result["ChartData"] !== null) {
                    drawDigitalChart(result)
                }
                document.getElementById("loader").hidden = true
            }
            const draw = performance.now();
            let difference = draw - download
            if (difference < 1000) {
                updateCharm("INF: Chart data drew in " + difference + "ms")
            } else {
                updateCharm("INF: Chart data drew in " + difference / 1000 + "s")
            }
        });
    }).catch((error) => {
        updateCharm("ERR: " + error)
        document.getElementById("loader").hidden = true
        document.getElementById("loader").style.transform = "none"
    });
})

function drawDigitalChart(chartData) {
    let locale = getLocaleFrom(chartData);
    moment.locale(locale);
    let minDate = document.getElementById("from-date").value
    let maxDate = document.getElementById("to-date").value
    let seriesList = [];
    let xAxisData = [];
    let yAxisData = [];
    let gridData = [];
    let sampling = "none"
    let counter = 1;
    xAxisData.push({
        type: 'time',
        min: minDate,
        max: maxDate,
        show: false
    });
    yAxisData.push({
        type: 'value',
        show: false,
    });
    let initialHeight = "60%"
    if (chartData["ChartData"].length !== 1) {
        initialHeight = (55 / chartData["ChartData"].length) - 1 + "%"
    }
    gridData.push({
        top: '5%',
        bottom: 75,
        left: 50,
        right: 100,
        height: initialHeight,
    });
    let dateAlreadyAdded = false
    for (const digitalData of chartData["ChartData"]) {
        if (counter > 1) {
            console.log(initialHeight)
            console.log((55 / chartData["ChartData"].length) * (counter - 1) + 5)
            gridData.push({
                left: 50,
                right: 100,
                top: ((counter - 1) * (55 / chartData["ChartData"].length) + 5) + "%",
                height: initialHeight
            });

            if (counter === chartData["ChartData"].length) {
                xAxisData.push({
                    axisLabel: {
                        formatter: function (value) {
                            return moment(new Date(value)).format('LLL');
                        }
                    },
                    gridIndex: counter - 1,
                    type: 'time',
                    min: minDate,
                    max: maxDate,
                    show: true,
                });
            } else {
                xAxisData.push({
                    gridIndex: counter - 1,
                    type: 'time',
                    min: minDate,
                    max: maxDate,
                    show: false,
                });
            }
            yAxisData.push({
                gridIndex: counter - 1,
                type: 'value',
                show: false,
            });
        }
        updateCharm("INF: " + digitalData["PortName"] + " with size: " + digitalData["DigitalData"].length)
        if (digitalData["DigitalData"].length > 8640 && flashButton.classList.contains("mif-flash-on")) {
            sampling = "lttb"
        }
        let data = []
        for (const element of digitalData["DigitalData"]) {
            data.push([new Date(element["Time"] * 1000), element["Value"]]);
        }
        dateAlreadyAdded = true
        seriesList.push({
            name: digitalData["PortName"],
            type: 'line',
            step: 'end',
            areaStyle: {},
            symbol: 'none',
            data: data,
            sampling: sampling,
            lineStyle: {
                width: 1,
            },
            emphasis: {
                focus: 'series'
            },
            xAxisIndex: counter - 1,
            yAxisIndex: counter - 1,
        });
        counter++;
    }

    if (!phonelinkButton.classList.contains("mif-phonelink-off")) {
        if (chartData["UserData"] !== null) {
            let data = []
            let color = ""
            for (const element of chartData["UserData"]) {
                color = element["Color"]
                data.push([new Date(element["FromDate"]), 1, element["Information"], moment(new Date(element["FromDate"])).format('LLL'), moment(new Date(element["ToDate"])).format('LLL')]);
                data.push([new Date(element["ToDate"]), 0, element["Information"], moment(new Date(element["FromDate"])).format('LLL'), moment(new Date(element["ToDate"])).format('LLL')]);
            }
            gridData.push({
                left: 50,
                right: 100,
                top: "70%",
                height: "4%"
            });
            xAxisData.push({
                gridIndex: chartData["ChartData"].length,
                type: 'time',
                min: minDate,
                max: maxDate,
                show: false,
            });
            yAxisData.push({
                gridIndex: chartData["ChartData"].length,
                type: 'value',
                show: false,
            });
            seriesList.push({
                name: chartData["UsersLocale"],
                color: color,
                areaStyle: {},
                type: 'line',
                step: 'end',
                symbol: 'none',
                data: data,
                sampling: 'none',
                lineStyle: {
                    width: 0,
                },
                xAxisIndex: chartData["ChartData"].length,
                yAxisIndex: chartData["ChartData"].length,
                emphasis: {
                    focus: 'series'
                },
            });
        }

        if (chartData["OrderData"] !== null) {
            let data = []
            let color = ""
            for (const element of chartData["OrderData"]) {
                color = element["Color"]
                data.push([new Date(element["FromDate"]), 1, element["Information"], moment(new Date(element["FromDate"])).format('LLL'), moment(new Date(element["ToDate"])).format('LLL')]);
                data.push([new Date(element["ToDate"]), 0, element["Information"], moment(new Date(element["FromDate"])).format('LLL'), moment(new Date(element["ToDate"])).format('LLL')]);
            }
            gridData.push({
                left: 50,
                right: 100,
                top: "75%",
                height: "4%"
            });
            xAxisData.push({
                gridIndex: chartData["ChartData"].length + 1,
                type: 'time',
                min: minDate,
                max: maxDate,
                show: false,
            });
            yAxisData.push({
                gridIndex: chartData["ChartData"].length + 1,
                type: 'value',
                show: false,
            });
            seriesList.push({
                name: chartData["OrdersLocale"],
                color: color,
                areaStyle: {},
                type: 'line',
                step: 'end',
                symbol: 'none',
                data: data,
                sampling: 'none',
                lineStyle: {
                    width: 0,
                },
                xAxisIndex: chartData["ChartData"].length + 1,
                yAxisIndex: chartData["ChartData"].length + 1,
                emphasis: {
                    focus: 'series'
                },
            });
        }
        if (chartData["DowntimeData"] !== null) {
            let data = []
            let color = ""
            for (const element of chartData["DowntimeData"]) {
                color = element["Color"]
                data.push([new Date(element["FromDate"]), 1, element["Information"], moment(new Date(element["FromDate"])).format('LLL'), moment(new Date(element["ToDate"])).format('LLL')]);
                data.push([new Date(element["ToDate"]), 0, element["Information"], moment(new Date(element["FromDate"])).format('LLL'), moment(new Date(element["ToDate"])).format('LLL')]);
            }
            gridData.push({
                left: 50,
                right: 100,
                top: "80%",
                height: "4%"
            });
            xAxisData.push({
                gridIndex: chartData["ChartData"].length + 2,
                type: 'time',
                min: minDate,
                max: maxDate,
                show: false,
            });
            yAxisData.push({
                gridIndex: chartData["ChartData"].length + 2,
                type: 'value',
                show: false,
            });
            seriesList.push({
                name: chartData["DowntimesLocale"],
                color: color,
                areaStyle: {},
                type: 'line',
                step: 'end',
                symbol: 'none',
                data: data,
                sampling: 'none',
                lineStyle: {
                    width: 0,
                },
                xAxisIndex: chartData["ChartData"].length + 2,
                yAxisIndex: chartData["ChartData"].length + 2,
                emphasis: {
                    focus: 'series'
                },
            });
        }
        if (chartData["BreakdownData"] !== null) {
            let data = []
            let color = ""
            for (const element of chartData["BreakdownData"]) {
                color = element["Color"]
                data.push([new Date(element["FromDate"]), 1, element["Information"], moment(new Date(element["FromDate"])).format('LLL'), moment(new Date(element["ToDate"])).format('LLL')]);
                data.push([new Date(element["ToDate"]), 0, element["Information"], moment(new Date(element["FromDate"])).format('LLL'), moment(new Date(element["ToDate"])).format('LLL')]);
            }
            gridData.push({
                left: 50,
                right: 100,
                top: "85%",
                height: "4%"
            });
            xAxisData.push({
                gridIndex: chartData["ChartData"].length + 3,
                type: 'time',
                min: minDate,
                max: maxDate,
                show: false,
            });
            yAxisData.push({
                gridIndex: chartData["ChartData"].length + 3,
                type: 'value',
                show: false,
            });
            seriesList.push({
                name: chartData["BreakdownsLocale"],
                color: color,
                areaStyle: {},
                type: 'line',
                step: 'end',
                symbol: 'none',
                data: data,
                sampling: 'none',
                lineStyle: {
                    width: 0,
                },
                xAxisIndex: chartData["ChartData"].length + 3,
                yAxisIndex: chartData["ChartData"].length + 3,
                emphasis: {
                    focus: 'series'
                },
            });
        }
        if (chartData["AlarmData"] !== null) {
            let data = []
            let color = ""
            for (const element of chartData["AlarmData"]) {
                color = element["Color"]
                data.push([new Date(element["FromDate"]), 1, element["Information"], moment(new Date(element["FromDate"])).format('LLL'), moment(new Date(element["ToDate"])).format('LLL')]);
                data.push([new Date(element["ToDate"]), 0, element["Information"], moment(new Date(element["FromDate"])).format('LLL'), moment(new Date(element["ToDate"])).format('LLL')]);
            }
            gridData.push({
                left: 50,
                right: 100,
                top: "90%",
                height: "4%"
            });
            xAxisData.push({
                gridIndex: chartData["ChartData"].length + 4,
                type: 'time',
                min: minDate,
                max: maxDate,
                show: false,
            });
            yAxisData.push({
                gridIndex: chartData["ChartData"].length + 4,
                type: 'value',
                show: false,
            });
            seriesList.push({
                name: chartData["AlarmsLocale"],
                color: color,
                areaStyle: {},
                type: 'line',
                step: 'end',
                symbol: 'none',
                data: data,
                sampling: 'none',
                lineStyle: {
                    width: 0,
                },
                xAxisIndex: chartData["ChartData"].length + 4,
                yAxisIndex: chartData["ChartData"].length + 4,
                emphasis: {
                    focus: 'series'
                },
            });
        }
    }
    let option;
    let sliderTopPosition = "90%"
    if (phonelinkButton.classList.contains("mif-phonelink-off")) {
        sliderTopPosition = "70%"
    }
    option = {
        animation: false,
        textStyle: {
            fontFamily: 'ProximaNova'
        },
        tooltip: {
            trigger: 'axis',
            axisPointer: {
                type: 'line',
            },
            formatter: function (params) {
                let result = ""
                for (const param of params) {
                    let color = param["color"]
                    if (param["axisIndex"] > chartData["ChartData"].length - 1) {
                        result += "<b>" + param["value"][3] + " - " + param["value"][4] + '</b><br><span style="display:inline-block;margin-right:5px;border-radius:10px;width:9px;height:9px;background-color:' + color + '"></span>' + param["value"][2] + "<br><br>"
                    } else {
                        result += "<b>" + moment(new Date(params[0]["axisValue"])).format('Do MMMM YYYY h:mm:ss') + "</b><br>" + '<span style="display:inline-block;margin-right:5px;border-radius:10px;width:9px;height:9px;background-color:' + color + '"></span>' + param["seriesName"] + " [" + param["value"][1] + "]<br><br>"
                    }
                }
                return result.replace(/^\s*<br\s*\/?>|<br\s*\/?>\s*$/g, '')
            },
            position: function (point, params, dom, rect, size) {
                return [point[0] - size["contentSize"][0] / 2, point[1]];
            }
        },
        axisPointer: {
            link: {xAxisIndex: 'all'}
        },
        grid: gridData,
        legend: {},
        toolbox: {
            right: 100,
            show: true,
            feature: {
                dataZoom: {
                    yAxisIndex: 'none'
                },
                saveAsImage: {
                    type: "png",
                    name: "digital"
                }
            }
        },
        yAxis: yAxisData,
        xAxis: xAxisData,
        dataZoom: [{
            type: 'inside',
            realtime: true,
            start: 0,
            end: 100
        }, {
            type: 'slider',
            realtime: true,
            showDataShadow: false,
            start: 0,
            end: 100,
            top: sliderTopPosition,
        }],
        series: seriesList,
    };
    option && myChart.setOption(option);
}


function drawAnalogChart(chartData) {
    let minDate = document.getElementById("from-date").value
    let maxDate = document.getElementById("to-date").value
    let locale = getLocaleFrom(chartData);
    let seriesList = [];
    let sampling = "none"
    moment.locale(locale);
    for (const analogData of chartData["ChartData"]) {
        if (analogData["AnalogData"].length > 8640 && flashButton.classList.contains("mif-flash-on")) {
            sampling = "average"
        }
        let data = []
        for (const element of analogData["AnalogData"]) {
            if (element["Value"] === -32768) {
                data.push([new Date(element["Time"] * 1000), null]);
            } else {
                data.push([new Date(element["Time"] * 1000), element["Value"]]);
            }
        }

        seriesList.push({
            name: analogData["PortName"],
            color: analogData["PortColor"],
            type: 'line',
            symbol: 'none',
            data: data,
            sampling: sampling,
            lineStyle: {
                width: 1,
            },
            emphasis: {
                focus: 'series'
            },
        });
    }
    if (!phonelinkButton.classList.contains("mif-phonelink-off")) {
        if (chartData["UserData"] !== null) {
            let userData = []
            let color = ""
            for (const element of chartData["UserData"]) {
                color = element["Color"]
                userData.push([new Date(element["FromDate"]), 1, element["Information"], moment(new Date(element["FromDate"])).format('LLL'), moment(new Date(element["ToDate"])).format('LLL')]);
                userData.push([new Date(element["ToDate"]), 0, element["Information"], moment(new Date(element["FromDate"])).format('LLL'), moment(new Date(element["ToDate"])).format('LLL')]);
            }
            seriesList.push({
                name: chartData["UsersLocale"],
                color: color,
                areaStyle: {},
                type: 'line',
                step: 'end',
                symbol: 'none',
                data: userData,
                sampling: 'none',
                lineStyle: {
                    width: 0,
                },
                xAxisIndex: 1,
                yAxisIndex: 1,
                emphasis: {
                    focus: 'series'
                },
            });
        }
        if (chartData["OrderData"] !== null) {
            let orderData = []
            let color = ""
            for (const element of chartData["OrderData"]) {
                color = element["Color"]
                orderData.push([new Date(element["FromDate"]), 1, element["Information"], moment(new Date(element["FromDate"])).format('LLL'), moment(new Date(element["ToDate"])).format('LLL')]);
                orderData.push([new Date(element["ToDate"]), 0, element["Information"], moment(new Date(element["FromDate"])).format('LLL'), moment(new Date(element["ToDate"])).format('LLL')]);
            }
            seriesList.push({
                name: chartData["OrdersLocale"],
                color: color,
                areaStyle: {},
                type: 'line',
                step: 'end',
                symbol: 'none',
                data: orderData,
                sampling: 'none',
                lineStyle: {
                    width: 0,
                },
                xAxisIndex: 2,
                yAxisIndex: 2,
                emphasis: {
                    focus: 'series'
                },
            });
        }
        if (chartData["DowntimeData"] !== null) {
            let downtimeData = []
            let color = ""
            for (const element of chartData["DowntimeData"]) {
                color = element["Color"]
                downtimeData.push([new Date(element["FromDate"]), 1, element["Information"], moment(new Date(element["FromDate"])).format('LLL'), moment(new Date(element["ToDate"])).format('LLL')]);
                downtimeData.push([new Date(element["ToDate"]), 0, element["Information"], moment(new Date(element["FromDate"])).format('LLL'), moment(new Date(element["ToDate"])).format('LLL')]);
            }
            seriesList.push({
                name: chartData["DowntimesLocale"],
                color: color,
                areaStyle: {},
                type: 'line',
                step: 'end',
                symbol: 'none',
                data: downtimeData,
                sampling: 'none',
                lineStyle: {
                    width: 0,
                },
                xAxisIndex: 3,
                yAxisIndex: 3,
                emphasis: {
                    focus: 'series'
                },
            });
        }
        if (chartData["BreakdownData"] !== null) {
            let breakdownData = []
            let color = ""
            for (const element of chartData["BreakdownData"]) {
                color = element["Color"]
                breakdownData.push([new Date(element["FromDate"]), 1, element["Information"], moment(new Date(element["FromDate"])).format('LLL'), moment(new Date(element["ToDate"])).format('LLL')]);
                breakdownData.push([new Date(element["ToDate"]), 0, element["Information"], moment(new Date(element["FromDate"])).format('LLL'), moment(new Date(element["ToDate"])).format('LLL')]);
            }
            seriesList.push({
                name: chartData["BreakdownsLocale"],
                color: color,
                areaStyle: {},
                type: 'line',
                step: 'end',
                symbol: 'none',
                data: breakdownData,
                sampling: 'none',
                lineStyle: {
                    width: 0,
                },
                xAxisIndex: 4,
                yAxisIndex: 4,
                emphasis: {
                    focus: 'series'
                },
            });
        }
        if (chartData["AlarmData"] !== null) {
            let alarmData = []
            let color = ""
            for (const element of chartData["AlarmData"]) {
                color = element["Color"]
                alarmData.push([new Date(element["FromDate"]), 1, element["Information"], moment(new Date(element["FromDate"])).format('LLL'), moment(new Date(element["ToDate"])).format('LLL')]);
                alarmData.push([new Date(element["ToDate"]), 0, element["Information"], moment(new Date(element["FromDate"])).format('LLL'), moment(new Date(element["ToDate"])).format('LLL')]);
            }
            seriesList.push({
                name: chartData["AlarmsLocale"],
                color: color,
                areaStyle: {},
                type: 'line',
                step: 'end',
                symbol: 'none',
                data: alarmData,
                sampling: 'none',
                lineStyle: {
                    width: 0,
                },
                xAxisIndex: 5,
                yAxisIndex: 5,
                emphasis: {
                    focus: 'series'
                },
            });
        }
    }

    let sliderTopPosition = "90%"
    if (phonelinkButton.classList.contains("mif-phonelink-off")) {
        sliderTopPosition = "70%"
    }
    let option;
    option = {
        animation: false,
        textStyle: {
            fontFamily: 'ProximaNova'
        },
        tooltip: {
            trigger: 'axis',
            axisPointer: {
                type: 'line',
            },
            formatter: function (params) {
                let result = ""
                for (const param of params) {
                    let color = param["color"]
                    if (param["axisIndex"] > 0) {
                        result += "<b>" + param["value"][3] + " - " + param["value"][4] + '</b><br><span style="display:inline-block;margin-right:5px;border-radius:10px;width:9px;height:9px;background-color:' + color + '"></span>' + param["value"][2] + "<br><br>"
                    } else {
                        result += "<b>" + moment(new Date(params[0]["axisValue"])).format('Do MMMM YYYY h:mm:ss') + "</b><br>" + '<span style="display:inline-block;margin-right:5px;border-radius:10px;width:9px;height:9px;background-color:' + color + '"></span>' + param["seriesName"] + " [" + param["value"][1] + "]<br><br>"
                    }
                }
                return result.replace(/^\s*<br\s*\/?>|<br\s*\/?>\s*$/g, '')
            },
            position: function (point, params, dom, rect, size) {
                return [point[0] - size["contentSize"][0] / 2, point[1]];
            }
        },
        grid: [{
            top: '5%',
            left: 50,
            right: 100,
            height: '60%'
        }, {
            left: 50,
            right: 100,
            top: '70%',
            height: '4%'
        }, {
            left: 50,
            right: 100,
            top: '75%',
            height: '4%'
        }, {
            left: 50,
            right: 100,
            top: '80%',
            height: '4%'
        }, {
            left: 50,
            right: 100,
            top: '85%',
            height: '4%'
        }, {
            left: 50,
            right: 100,
            top: '90%',
            height: '4%'
        }],
        legend: {},
        toolbox: {
            right: 100,
            show: true,
            feature: {
                dataZoom: {
                    yAxisIndex: 'none'
                },
                saveAsImage: {
                    type: "png",
                    name: "analog"
                }
            }
        },
        xAxis: [
            {
                type: 'time',
                axisLabel: {
                    formatter: function (value) {
                        return moment(new Date(value)).format('LLL');
                    }
                },
                min: minDate,
                max: maxDate,
                show: true
            },
            {
                gridIndex: 1,
                type: 'time',
                min: minDate,
                max: maxDate,
                show: false,
            },
            {
                gridIndex: 2,
                type: 'time',
                min: minDate,
                max: maxDate,
                show: false,
            },
            {
                gridIndex: 3,
                type: 'time',
                min: minDate,
                max: maxDate,
                show: false,
            },
            {
                gridIndex: 4,
                type: 'time',
                min: minDate,
                max: maxDate,
                show: false,
            },
            {
                gridIndex: 5,
                type: 'time',
                min: minDate,
                max: maxDate,
                show: false,
            }
        ],
        yAxis: [
            {
                type: 'value',
            },
            {
                gridIndex: 1,
                type: 'value',
                show: false,
            },
            {
                gridIndex: 2,
                type: 'value',
                show: false,
            },
            {
                gridIndex: 3,
                type: 'value',
                show: false,
            },
            {
                gridIndex: 4,
                type: 'value',
                show: false,
            },
            {
                gridIndex: 5,
                type: 'value',
                show: false,
            }
        ],
        axisPointer: {
            link: {xAxisIndex: 'all'}
        },
        dataZoom: [
            {
                type: 'inside',
                id: 'insideChartZoom',
                start: 0,
                end: 100,
                xAxisIndex: [0]
            }, {
                type: 'slider',
                id: 'sliderChartZoom',

                labelFormatter: function (value) {
                    return moment(new Date(value)).format('LLL');
                },
                realtime: true,
                showDataShadow: false,
                start: 0,
                end: 100,
                xAxisIndex: [0],
                top: sliderTopPosition
            },
            {
                type: 'inside',
                id: 'terminalInsideChartZoom',
                filterMode: 'none',
                show: false,
                xAxisIndex: [1, 2, 3, 4, 5]
            }, {
                type: 'slider',
                id: 'terminalSliderChartZoom',
                showDataShadow: false,
                filterMode: 'none',
                show: false,
                xAxisIndex: [1, 2, 3, 4, 5]
            },
        ],
        series: seriesList,
    };
    option && myChart.setOption(option);
}

function getLocaleFrom(chartData) {
    let locale = ""
    switch (chartData["Locale"]) {
        case "CsCZ": {
            locale = "cs";
            break;
        }
        case "DeDE": {
            locale = "de";
            break;
        }
        case "EnUS": {
            locale = "en";
            break;
        }
        case "EsES": {
            locale = "es";
            break;
        }
        case "FrFR": {
            locale = "fr";
            break;
        }
        case "ItIT": {
            locale = "it";
            break;
        }
        case "PlPL": {
            locale = "pl";
            break;
        }
        case "PtPT": {
            locale = "pt";
            break;
        }
        case "SkSK": {
            locale = "sk";
            break;
        }
        case "RuRU": {
            locale = "ru";
            break;
        }
    }
    return locale;
}

window.addEventListener('resize', () => {
    myChart.resize()

})

myChart.on('dataZoom', function (evt) {
    let option = myChart.getOption();
    if (evt["dataZoomId"] !== undefined) {
        myChart.dispatchAction({
            type: 'dataZoom',
            startValue: option.dataZoom[0].startValue,
            endValue: option.dataZoom[0].endValue
        });
    } else if (evt["batch"] !== undefined) {
        myChart.dispatchAction({
            type: 'dataZoom',
            startValue: evt["batch"][0]["startValue"],
            endValue: evt["batch"][0]["endValue"]
        });
    } else if (evt["startValue"] === undefined) {
        myChart.dispatchAction({
            type: 'dataZoom',
            startValue: option.dataZoom[0].startValue,
            endValue: option.dataZoom[0].endValue
        });
    }
});





