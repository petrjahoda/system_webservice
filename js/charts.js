let chartDom = document.getElementById('chart');
let myChart = echarts.init(chartDom);
let now = new Date();
now.setMinutes(now.getMinutes() - now.getTimezoneOffset());
document.getElementById('to-date').value = now.toISOString().slice(0, 16);
now.setHours(now.getHours() - 24);
document.getElementById('from-date').value = now.toISOString().slice(0, 16);
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
            updateCharm("INF: Analog chart data downloaded from database in " + difference + "ms")
        } else {
            updateCharm("INF: Analog chart data downloaded from database in " + difference / 1000 + "s")
        }
        document.getElementById("loader").style.transform = "rotateY(180deg)"
        response.text().then(function (data) {
            let result = JSON.parse(data);
            if (result["Type"] === "analog-data") {
                myChart.clear()
                if (result["AnalogData"] !== null) {
                    drawAnalogChart(result)
                }
                document.getElementById("loader").style.transform = "none"
                document.getElementById("loader").hidden = true
            }
            const draw = performance.now();
            let difference = draw - download
            if (difference < 1000) {
                updateCharm("INF: Analog chart data drew in " + difference + "ms")
            } else {
                updateCharm("INF: Analog chart data drew in " + difference / 1000 + "s")
            }
        });
    }).catch((error) => {
        updateCharm("ERR: " + error)
        document.getElementById("loader").hidden = true
        document.getElementById("loader").style.transform = "none"
    });
})

function drawAnalogChart(chartData) {
    let locale = getLocaleFrom(chartData);
    let seriesList = [];
    let date = []
    let sampling = "none"
    let dateAlreadyAdded = false
    moment.locale(locale);
    for (const analogData of chartData["AnalogData"]) {
        updateCharm("INF: " + analogData["PortName"] + " with size: " + analogData["PortData"].length)
        if (analogData["PortData"].length > 8640) {
            sampling = "average"
        }
        let data = []
        for (const element of analogData["PortData"]) {
            if (!dateAlreadyAdded) {
                date.push(moment(new Date(element["Time"] * 1000)).format('Do MMM YYYY, h:mm:ss'));
            }
            if (element["Value"] === -32768) {
                data.push(null);
            } else {
                data.push(element["Value"]);
            }
        }
        dateAlreadyAdded = true
        seriesList.push({
            name: analogData["PortName"],
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
    let option;
    option = {
        animation: false,
        textStyle: {
            fontFamily: 'ProximaNova'
        },
        tooltip: {
            trigger: 'axis',
            axisPointer: {
                type: 'cross',
                snap: true,

            },
            position: function (point, params, dom, rect, size) {
                return [point[0] - size["contentSize"][0] / 2, point[1]];
            }
        },
        grid: {
            top: 50,
            bottom: 75,
            left: 50,
            right: 100,
        },
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
        yAxis: {
            type: 'value',
        },
        xAxis: {
            data: date,
            axisLabel: {}
        },
        dataZoom: [{
            type: 'inside',
            realtime: true,
            start: 0,
            end: 100
        }, {
            type: 'slider',
            realtime: true,
            start: 0,
            end: 100,
        }],
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





