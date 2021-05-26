function drawDigitalChart(chartData) {
    let minDate = document.getElementById("from-date").value
    let maxDate = document.getElementById("to-date").value
    startDateAsValue = new Date(minDate) * 1000
    endDateAsValue = new Date(maxDate) * 1000
    let positionInChart
    let locale = getLocaleFrom(chartData);
    moment.locale(locale);
    let seriesList = [];
    let xAxisData = [];
    let yAxisData = [];
    let gridData = [];
    let sampling = "none"
    let position = 0
    for (const digitalData of chartData["ChartData"]) {
        if (digitalData["DigitalData"].length > 8640 && flashButton.classList.contains("mif-flash-on")) {
            sampling = "lttb"
        }
        let data = []
        for (const element of digitalData["DigitalData"]) {
            if (element["Value"] === -32768) {
                data.push([new Date(element["Time"] * 1000), null]);
            } else {
                data.push([new Date(element["Time"] * 1000), element["Value"]]);
            }
        }
        gridData.push({
            left: 50,
            right: 100,
            top: (position*7)+10+"%",
            height: "3%"
        });
        if (position === chartData["ChartData"].length-1) {
            xAxisData.push({
                gridIndex: position,
                type: 'time',
                min: minDate,
                max: maxDate,
                axisLabel: {
                    formatter: function (value) {
                        return moment(new Date(value)).format('LLL');
                    }
                },
                show: true,
            });
        } else {
            xAxisData.push({
                gridIndex: position,
                type: 'time',
                min: minDate,
                max: maxDate,
                show: false,
            });
        }
        yAxisData.push({
            gridIndex: position,
            type: 'value',
            show: false,
        });
        seriesList.push({
            name: digitalData["PortName"],
            color: digitalData["PortColor"],
            areaStyle: {},
            type: 'line',
            symbol: 'none',
            step: 'end',
            data: data,
            sampling: sampling,
            lineStyle: {
                width: 1,
            },
            emphasis: {
                focus: 'series'
            },
            xAxisIndex: position,
            yAxisIndex: position,
        });
        position++
    }
    let initialTerminalPosition = (position*7)+11
    gridData.push({
        left: 50,
        right: 100,
        top: initialTerminalPosition+"%",
        height: "3%"
    });
    xAxisData.push({
        gridIndex: position,
        type: 'time',
        min: minDate,
        max: maxDate,
        show: false,
    });
    yAxisData.push({
        gridIndex: position,
        type: 'value',
        show: false,
    });
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
                xAxisIndex: position,
                yAxisIndex: position,
                emphasis: {
                    focus: 'series'
                },
            });
        }
        gridData.push({
            left: 50,
            right: 100,
            top: initialTerminalPosition+4+"%",
            height: "3%"
        });
        xAxisData.push({
            gridIndex: position+1,
            type: 'time',
            min: minDate,
            max: maxDate,
            show: false,
        });
        yAxisData.push({
            gridIndex: position+1,
            type: 'value',
            show: false,
        });
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
                xAxisIndex: position+1,
                yAxisIndex: position+1,
                emphasis: {
                    focus: 'series'
                },
            });

        }
        gridData.push({
            left: 50,
            right: 100,
            top: initialTerminalPosition+8+"%",
            height: "3%"
        });
        xAxisData.push({
            gridIndex: position+2,
            type: 'time',
            min: minDate,
            max: maxDate,
            show: false,
        });
        yAxisData.push({
            gridIndex: position+2,
            type: 'value',
            show: false,
        });
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
                xAxisIndex: position+2,
                yAxisIndex: position+2,
                emphasis: {
                    focus: 'series'
                },
            });
        }
        gridData.push({
            left: 50,
            right: 100,
            top: initialTerminalPosition+12+"%",
            height: "3%"
        });
        xAxisData.push({
            gridIndex: position+3,
            type: 'time',
            min: minDate,
            max: maxDate,
            show: false,
        });
        yAxisData.push({
            gridIndex: position+3,
            type: 'value',
            show: false,
        });
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
                xAxisIndex: position+3,
                yAxisIndex: position+3,
                emphasis: {
                    focus: 'series'
                },
            });
        }
        gridData.push({
            left: 50,
            right: 100,
            top: initialTerminalPosition+16+"%",
            height: "3%"
        });
        xAxisData.push({
            gridIndex: position+4,
            type: 'time',
            min: minDate,
            max: maxDate,
            show: false,
        });
        yAxisData.push({
            gridIndex: position+4,
            type: 'value',
            show: false,
        });
        console.log(chartData["AlarmData"] !== null)
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
                xAxisIndex: position+4,
                yAxisIndex: position+4,
                emphasis: {
                    focus: 'series'
                },
            });
        }
    }
    let option;
    let sliderTopPosition = initialTerminalPosition+20+"%"
    if (phoneLinkButton.classList.contains("mif-phonelink-off")) {
        sliderTopPosition = initialTerminalPosition+"%"
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
            position: function (point) {
                positionInChart = point[0]
            },
            formatter: function (params) {
                let dateChange = endDateAsValue - startDateAsValue
                let pointerValue = (startDateAsValue + ((positionInChart - borderStart) * (dateChange / borderChange))) / 1000
                let result = ""
                for (const param of params) {
                    if (pointerValue > +param["value"][3] && pointerValue < +param["value"][4]) {
                        result += '<span style="display:inline-block;margin-right:5px;border-radius:10px;width:9px;height:9px;background-color:' + param["color"] + '"></span>' + param["value"][2] + '&nbsp&nbsp&nbsp<span style="font-size:' + 10 + 'px">' + moment(param["value"][3]).format('LLL') + " - " + moment(param["value"][4]).format("LLL") + '</span>' + "<br>"
                    } else {
                        if (param["seriesIndex"] < chartData["ChartData"].length && param["value"][1] !== null) {
                            result += '<span style="display:inline-block;margin-right:5px;border-radius:10px;width:9px;height:9px;background-color:' + param["color"] + '"></span>' + param["seriesName"] + ": " + param["value"][1] + '&nbsp&nbsp&nbsp<span style="font-size:' + 10 + 'px">' + moment(param["value"][0]).format('Do MMMM YYYY, h:mm:ss') + "</span><br>"
                        }
                    }
                }
                return "<b>" + '<span style="border-bottom: 1px solid;width: 100%;display: block;">' + moment(pointerValue).format('Do MMMM YYYY, h:mm:ss') + "</span></b><br>" + result
            },
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
                    name: "digital-data"
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
                top: sliderTopPosition,
            },
            {
                type: 'inside',
                id: 'terminalInsideChartZoom',
                filterMode: 'none',
                show: false,
                xAxisIndex: [1, 2, 3, 4, 5,6]
            }, {
                type: 'slider',
                id: 'terminalSliderChartZoom',
                showDataShadow: false,
                filterMode: 'none',
                show: false,
                xAxisIndex: [1, 2, 3, 4, 5,6]
            },
        ],
        series: seriesList,
    };
    option && myChart.setOption(option);
}