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
        if (digitalData["DigitalData"] === null) {
            document.getElementById("loader").hidden = true
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
    if (!phoneLinkButton.classList.contains("mif-phonelink-off")) {
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
                height: "3%"
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
                top: "74%",
                height: "3%"
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
                top: "78%",
                height: "3%"
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
                top: "82%",
                height: "3%"
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
                top: "86%",
                height: "3%"
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
    if (phoneLinkButton.classList.contains("mif-phonelink-off")) {
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