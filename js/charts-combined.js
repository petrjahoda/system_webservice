function drawCombinedChart(chartData) {
    let minDate = document.getElementById("from-date").value
    let maxDate = document.getElementById("to-date").value
    let locale = getLocaleFrom(chartData);
    let seriesList = [];
    let sampling = "none"
    moment.locale(locale);
    chartData["ChartData"].forEach((element, index) => {
        console.log(element)
        let data = []
        switch (element["PortType"]) {
            case "digital":

                for (const record of element["DigitalData"]) {
                    data.push([new Date(record["Time"] * 1000), record["Value"]]);
                }
                seriesList.push({
                    name: element["PortName"],
                    color: element["PortColor"],
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
                    xAxisIndex: 1,
                    yAxisIndex: 1,
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
                    for (const record of element["DigitalData"]) {
                        data.push([new Date(record["Time"] * 1000), record["Value"]]);
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
    let positionInChart

    let borderStart = 50
    let borderEnd = chartDom.scrollWidth-100
    let borderChange = borderEnd-borderStart
    let startDateAsValue = new Date(minDate) * 1000
    let endDateAsValue = new Date(maxDate) * 1000
    let dateChange = endDateAsValue - startDateAsValue
    option = {
        animation: false,
        textStyle: {
            fontFamily: 'ProximaNova'
        },
        tooltip: {
            trigger: 'axis',
            axisPointer: {
                type: 'none',
            },
            position: function (point, params, dom, rect, size) {
                positionInChart = point[0]
                return [point[0] - size["contentSize"][0] / 2, point[1]];
            },
            // formatter: function (params, point) {
            //     let result = ""
            //     let doIt = true
            //     let pointerDate = new Date((startDateAsValue + ((positionInChart - borderStart) * (dateChange / borderChange))) / 1000)
            //     for (const series of seriesList) {
            //         if (series["xAxisIndex"] === 0 && doIt) {
            //             console.log(series["data"].filter(word => word[1] === 1))
            //             for (let i = 0; i < series["data"].length; i++) {
            //                 if (series["data"][i][0]<pointerDate && series["data"][i][1] === 1&&series["data"][i+1][0]>new Date((startDateAsValue + ((positionInChart - borderStart) * (dateChange / borderChange))) / 1000) ) {
            //                     result += series["name"]
            //                     doIt = false
            //                     break
            //                 }
            //             }
            //
            //         }
            //
            //     }
            //     // return "<b>" + moment(new Date((startDateAsValue + ((positionInChart - borderStart) * (dateChange / borderChange)))/1000).format('LLL') + "</b><br>" + result)
            //     return "<b>" + new Date((startDateAsValue + ((positionInChart - borderStart) * (dateChange / borderChange))) / 1000) + "</b><br>" + result
            // },

        },
        grid: [{
            left: 50,
            right: 100,
            top: '5%',
            height: '3%'
        }, {
            left: 50,
            right: 100,
            top: '9%',
            height: '3%'
        }, {
            left: 50,
            right: 100,
            top: '13%',
            height: '52%'
        }, {
            left: 50,
            right: 100,
            top: '70%',
            height: '3%'
        }, {
            left: 50,
            right: 100,
            top: '74%',
            height: '3%'
        }, {
            left: 50,
            right: 100,
            top: '78%',
            height: '3%'
        }, {
            left: 50,
            right: 100,
            top: '82%',
            height: '3%'
        }, {
            left: 50,
            right: 100,
            top: '86%',
            height: '3%'
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
                xAxisIndex: [0, 1, 3, 4, 5, 6, 7]
            }, {
                type: 'slider',
                id: 'terminalSliderChartZoom',
                showDataShadow: false,
                filterMode: 'none',
                show: false,
                xAxisIndex: [0, 1, 3, 4, 5, 6, 7]
            },
        ],
        series: seriesList,
    };
    option && myChart.setOption(option);
}