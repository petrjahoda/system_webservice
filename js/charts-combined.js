function drawCombinedChart(chartData) {
    let minDate = document.getElementById("from-date").value
    let maxDate = document.getElementById("to-date").value
    startDateAsValue = new Date(minDate) * 1000
    endDateAsValue = new Date(maxDate) * 1000
    let positionInChart
    let locale = getLocaleFrom(chartData);
    let seriesList = [];
    let digitalSampling = "none"
    let sampling = "none"
    moment.locale(locale);
    chartData["ChartData"].forEach((element) => {
        let data = []
        switch (element["PortType"]) {
            case "digital":
                for (const record of element["DigitalData"]) {
                    data.push([new Date(record["Time"] * 1000), record["Value"]]);
                }
                if (element["DigitalData"].length > (8640*2) && flashButton.classList.contains("mif-flash-on")) {
                    digitalSampling = "lttb"
                }
                seriesList.push({
                    name: element["PortName"],
                    color: element["PortColor"],
                    type: 'line',
                    step: 'end',
                    areaStyle: {},
                    symbol: 'none',
                    data: data,
                    sampling: digitalSampling,
                    lineStyle: {
                        width: 1,
                    },
                    emphasis: {
                        focus: 'series'
                    },
                    xAxisIndex: 1,
                    yAxisIndex: 1,
                    cursor: 'default'
                });
                break;
            case "analog":
                if (element["AnalogData"].length > 8640 && flashButton.classList.contains("mif-flash-on")) {
                    sampling = "average"
                }
                for (const record of element["AnalogData"]) {
                    if (record["Value"] === -32768) {
                        data.push([new Date(record["Time"] * 1000), null]);
                    } else {
                        data.push([new Date(record["Time"] * 1000), record["Value"]]);
                    }
                }
                seriesList.push({
                    name: element["PortName"],
                    color: element["PortColor"],
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
                    xAxisIndex: 2,
                    yAxisIndex: 2,
                });
                break;
            default:
                if (element["DigitalData"] !== null) {
                    data.push([new Date(element["DigitalData"][0]["Time"] * 1000), element["DigitalData"][0]["Value"], element["PortName"], new Date(element["DigitalData"][0]["Time"] * 1000), new Date(element["DigitalData"][1]["Time"] * 1000)]);
                    for (let i = 1; i < element["DigitalData"].length; i++) {
                        if (i + 1 === element["DigitalData"].length) {
                            if (element["DigitalData"][i]["Value"] === 1) {
                                data.push([new Date(element["DigitalData"][i]["Time"] * 1000), element["DigitalData"][i]["Value"], element["PortName"], new Date(element["DigitalData"][i]["Time"] * 1000), new Date(element["DigitalData"][i]["Time"] * 1000)]);
                            } else {
                                data.push([new Date(element["DigitalData"][i]["Time"] * 1000), element["DigitalData"][i]["Value"], element["PortName"], new Date(element["DigitalData"][i - 1]["Time"] * 1000), new Date(element["DigitalData"][i]["Time"] * 1000)]);
                            }
                            break
                        }
                        if (element["DigitalData"][i]["Value"] === 1) {
                            data.push([new Date(element["DigitalData"][i]["Time"] * 1000), element["DigitalData"][i]["Value"], element["PortName"], new Date(element["DigitalData"][i]["Time"] * 1000), new Date(element["DigitalData"][i + 1]["Time"] * 1000)]);
                        } else {
                            data.push([new Date(element["DigitalData"][i]["Time"] * 1000), element["DigitalData"][i]["Value"], element["PortName"], new Date(element["DigitalData"][i - 1]["Time"] * 1000), new Date(element["DigitalData"][i]["Time"] * 1000)]);
                        }

                    }
                    seriesList.push({
                        name: element["PortName"],
                        color: element["PortColor"],
                        type: 'line',
                        step: 'end',
                        areaStyle: {},
                        symbol: 'none',
                        sampling: 'none',
                        data: data,
                        lineStyle: {
                            width: 0,
                        },
                        emphasis: {
                            focus: 'series'
                        },
                        xAxisIndex: 0,
                        yAxisIndex: 0,
                    });
                }
                break;
        }
    })
    if (!phoneLinkButton.classList.contains("mif-phonelink-off")) {
        if (chartData["UserData"] !== null) {
            let userData = []
            let color = ""
            for (const element of chartData["UserData"]) {
                color = element["Color"]
                userData.push([new Date(element["FromDate"]), 1, element["Information"], element["FromDate"], element["ToDate"]]);
                userData.push([new Date(element["ToDate"]), 0, element["Information"], element["FromDate"], element["ToDate"]]);
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
                xAxisIndex: 3,
                yAxisIndex: 3,
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
                orderData.push([new Date(element["FromDate"]), 1, element["Information"], element["FromDate"], element["ToDate"]]);
                orderData.push([new Date(element["ToDate"]), 0, element["Information"], element["FromDate"], element["ToDate"]]);
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
                xAxisIndex: 4,
                yAxisIndex: 4,
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
                downtimeData.push([new Date(element["FromDate"]), 1, element["Information"], element["FromDate"], element["ToDate"]]);
                downtimeData.push([new Date(element["ToDate"]), 0, element["Information"], element["FromDate"], element["ToDate"]]);
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
                xAxisIndex: 5,
                yAxisIndex: 5,
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
                breakdownData.push([new Date(element["FromDate"]), 1, element["Information"], element["FromDate"], element["ToDate"]]);
                breakdownData.push([new Date(element["ToDate"]), 0, element["Information"], element["FromDate"], element["ToDate"]]);
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
                xAxisIndex: 6,
                yAxisIndex: 6,
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
                alarmData.push([new Date(element["FromDate"]), 1, element["Information"], element["FromDate"], element["ToDate"]]);
                alarmData.push([new Date(element["ToDate"]), 0, element["Information"], element["FromDate"], element["ToDate"]]);
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
                xAxisIndex: 7,
                yAxisIndex: 7,
                emphasis: {
                    focus: 'series'
                },
            });
        }
    }

    let sliderTopPosition = "90%"
    if (phoneLinkButton.classList.contains("mif-phonelink-off")) {
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
                snap: true,
                type: 'line',
            },
            position: function (point) {
                positionInChart = point[0]
            },
            formatter: function (params) {
                let dateChange = endDateAsValue - startDateAsValue
                let pointerValue = (startDateAsValue + (positionInChart/document.getElementById("chart").offsetWidth)*dateChange) / 1000
                let result = ""
                for (const param of params) {
                    if (pointerValue > +param["value"][3] && pointerValue < +param["value"][4]) {
                        result += '<span style="display:inline-block;margin-right:5px;border-radius:10px;width:9px;height:9px;background-color:' + param["color"] + '"></span>' + param["value"][2] + '&nbsp&nbsp&nbsp<span style="font-size:' + 10 + 'px">' + moment(param["value"][3]).format('LLL') + " - " + moment(param["value"][4]).format("LLL") + '</span>' + "<br>"
                    } else {
                        if (param["seriesIndex"] === 1 && param["value"][1] !== null) {
                            result += '<span style="display:inline-block;margin-right:5px;border-radius:10px;width:9px;height:9px;background-color:' + param["color"] + '"></span>' +param["seriesName"] +": "+ param["value"][1] + "<br>"
                        } else if (param["seriesIndex"] === 0 && param["value"][1] !== null){
                            result += '<span style="display:inline-block;margin-right:5px;border-radius:10px;width:9px;height:9px;background-color:' + param["color"] + '"></span>' +param["seriesName"] +": "+ param["value"][1] + "<br>"
                        }
                    }
                }
                return "<b>" + '<span style="border-bottom: 1px solid;width: 100%;display: block;">' + moment(pointerValue).format('Do MMMM YYYY, h:mm:ss') + "</span></b><br>" + result
            },
        },
        grid: [{
            left: 0,
            right: 0,
            top: '5%',
            height: '3%'
        }, {
            left: 0,
            right: 0,
            top: '9%',
            height: '3%'
        }, {
            left: 0,
            right: 0,
            top: '13%',
            height: '52%'
        }, {
            left: 0,
            right: 0,
            top: '70%',
            height: '3%'
        }, {
            left: 0,
            right: 0,
            top: '74%',
            height: '3%'
        }, {
            left: 0,
            right: 0,
            top: '78%',
            height: '3%'
        }, {
            left: 0,
            right: 0,
            top: '82%',
            height: '3%'
        }, {
            top: '86%',
            height: '3%'
        }],
        legend: {},
        toolbox: {
            right: 0,
            show: true,
            feature: {
                dataZoom: {
                    yAxisIndex: 'none'
                },
                saveAsImage: {
                    type: "png",
                    name: "combined-chart"
                }
            }
        },
        xAxis: [
            {
                type: 'time',
                min: minDate,
                max: maxDate,
                show: false
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
                axisLabel: {
                    formatter: function (value) {
                        return moment(new Date(value)).format('LLL');
                    }
                },
                type: 'time',
                min: minDate,
                max: maxDate,
                show: true,
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
            },
            {
                gridIndex: 6,
                type: 'time',
                min: minDate,
                max: maxDate,
                show: false,
            },
            {
                gridIndex: 7,
                type: 'time',
                min: minDate,
                max: maxDate,
                show: false,
            }
        ],
        yAxis: [
            {
                type: 'value',
                show: false
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
            },
            {
                gridIndex: 6,
                type: 'value',
                show: false,
            },
            {
                gridIndex: 7,
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
                xAxisIndex: [2]
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
                xAxisIndex: [2],
                top: sliderTopPosition
            },
            {
                type: 'inside',
                id: 'terminalInsideChartZoom',
                filterMode: 'none',
                show: false,
                xAxisIndex: [0, 1, 3, 4, 5, 6, 7,8,9]
            }, {
                type: 'slider',
                id: 'terminalSliderChartZoom',
                showDataShadow: false,
                filterMode: 'none',
                show: false,
                xAxisIndex: [0, 1, 3, 4, 5, 6, 7,8,9]
            },
        ],
        series: seriesList,
    };
    option && myChart.setOption(option);
}