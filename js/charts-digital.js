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
    let counter = 1;
    let showAxis = false
    if (chartData["ChartData"].length === 1) {
        showAxis = true
    }
    xAxisData.push({
        axisLabel: {
            formatter: function (value) {
                return moment(new Date(value)).format('LLL');
            }
        },
        type: 'time',
        min: minDate,
        max: maxDate,
        show: showAxis
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
        let color = ""
        for (const element of digitalData["DigitalData"]) {
            color = element["Color"]
            data.push([new Date(element["Time"] * 1000), element["Value"]]);
        }
        dateAlreadyAdded = true
        seriesList.push({
            name: digitalData["PortName"],
            color: color,
            areaStyle: {},
            type: 'line',
            step: 'end',
            symbol: 'none',
            data: data,
            sampling: sampling,
            lineStyle: {
                width: 0,
            },
            xAxisIndex: counter - 1,
            yAxisIndex: counter - 1,
            emphasis: {
                focus: 'series'
            },
        });
        counter++;
    }
    if (!phoneLinkButton.classList.contains("mif-phonelink-off")) {
        if (chartData["UserData"] !== null) {
            let userData = []
            let color = ""
            for (const element of chartData["UserData"]) {
                color = element["Color"]
                userData.push([new Date(element["FromDate"]), 1, element["Information"], element["FromDate"], element["ToDate"]]);
                userData.push([new Date(element["ToDate"]), 0, element["Information"], element["FromDate"], element["ToDate"]]);
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
                data: userData,
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
            let orderData = []
            let color = ""
            for (const element of chartData["OrderData"]) {
                color = element["Color"]
                orderData.push([new Date(element["FromDate"]), 1, element["Information"], element["FromDate"], element["ToDate"]]);
                orderData.push([new Date(element["ToDate"]), 0, element["Information"], element["FromDate"], element["ToDate"]]);
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
                data: orderData,
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
            let downtimeData = []
            let color = ""
            for (const element of chartData["DowntimeData"]) {
                color = element["Color"]
                downtimeData.push([new Date(element["FromDate"]), 1, element["Information"], element["FromDate"], element["ToDate"]]);
                downtimeData.push([new Date(element["ToDate"]), 0, element["Information"], element["FromDate"], element["ToDate"]]);
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
                data: downtimeData,
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
            let breakdownData = []
            let color = ""
            for (const element of chartData["BreakdownData"]) {
                color = element["Color"]
                breakdownData.push([new Date(element["FromDate"]), 1, element["Information"], element["FromDate"], element["ToDate"]]);
                breakdownData.push([new Date(element["ToDate"]), 0, element["Information"], element["FromDate"], element["ToDate"]]);
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
                data: breakdownData,
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
            let alarmData = []
            let color = ""
            for (const element of chartData["AlarmData"]) {
                color = element["Color"]
                alarmData.push([new Date(element["FromDate"]), 1, element["Information"], element["FromDate"], element["ToDate"]]);
                alarmData.push([new Date(element["ToDate"]), 0, element["Information"], element["FromDate"], element["ToDate"]]);
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
                data: alarmData,
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
                            result += '<span style="display:inline-block;margin-right:5px;border-radius:10px;width:9px;height:9px;background-color:' + param["color"] + '"></span>' + param["seriesName"] + ": " + param["value"][1] + "<br>"
                        }
                    }
                }
                return "<b>" + moment(pointerValue).format('Do MMMM YYYY, h:mm:ss') + "</b><br>" + result
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