function drawProductionChart(chartData) {
    let minDate = document.getElementById("from-date").value
    let maxDate = document.getElementById("to-date").value
    startDateAsValue = new Date(minDate) * 1000
    endDateAsValue = new Date(maxDate) * 1000
    let positionInChart
    let locale = getLocaleFrom(chartData);
    let seriesList = [];
    let sampling = "none"
    moment.locale(locale);
    for (const analogData of chartData["ChartData"]) {
        if (analogData["DigitalData"] === null) {
            document.getElementById("loader").hidden = true
        }
        if (analogData["DigitalData"].length > 8640 && flashButton.classList.contains("mif-flash-on")) {
            sampling = "average"
        }
        let data = []
        for (const element of analogData["DigitalData"]) {
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
            symbolSize: [2, 2],
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
                xAxisIndex: 5,
                yAxisIndex: 5,
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
                        if (param["seriesIndex"] === 0 && param["value"][1] !== null) {
                            result += '<span style="display:inline-block;margin-right:5px;border-radius:10px;width:9px;height:9px;background-color:' + param["color"] + '"></span>' + param["seriesName"] + ": " + param["value"][1] + "<br>"
                        }
                    }
                }
                return "<b>" + '<span style="border-bottom: 1px solid;width: 100%;display: block;">' + moment(pointerValue).format('Do MMMM YYYY, h:mm:ss') + "</span></b><br>" + result
            },

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
                axisLabel: {
                    show: true,
                    fontSize: 10,
                }
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