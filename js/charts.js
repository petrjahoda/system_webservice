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
    fetch("/load_chart_data", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            let result = JSON.parse(data);
            updateCharm(result["Result"])
            if (result["Type"] === "analog-data") {
                console.log("BEFORE " + new Date().toISOString())
                drawAnalogChart(result)
                console.log("AFTER  " + new Date().toISOString())
                document.getElementById("loader").hidden = true
            }
        });
    }).catch((error) => {
        updateCharm("ERR: " + error)
        document.getElementById("loader").hidden = true
    });
})

function drawAnalogChart(chartData) {
    const intermediateData = chartData["AnalogData"][0]["PortData"]
    var chartDom = document.getElementById('chart');
    var myChart = echarts.init(chartDom);
    var option;
    var date = [];

    // var data = [Math.random() * 300];
    // var base = +new Date(1968, 9, 3);
    // var oneDay = 10000;
    // for (var i = 1; i < 2000000; i++) {
    //     var jetzt = new Date(base += oneDay);
    //     date.push(jetzt);
    //     data.push(Math.round((Math.random() - 0.5) * 20 + data[i - 1]));
    // }

    var data = [];
    for (const element of intermediateData) {
        date.push(new Date(element["Time"]*1000));
        if (element["Value"] === -32768) {
            data.push(null);
            // data.push(0);
        } else {
            data.push(element["Value"]);
        }
    }
    console.log(data.length.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ","))
    option = {
        tooltip: {
            trigger: 'axis'
        },
        title: {
            left: 'center',
            text: 'ANALOG DATA',
        },
        toolbox: {
            feature: {
                dataZoom: {
                    yAxisIndex: 'none'
                },
                restore: {},
                saveAsImage: {}
            }
        },
        xAxis: {
            type: 'category',
            data: date
        },
        yAxis: {
            type: 'value',
        },
        dataZoom: [{
            type: 'inside',
            start: 0,
            end: 100
        }, {
            start: 0,
            end: 100
        }],
        series: [
            {
                name: 'TEST DATA',
                type: 'line',
                symbol: 'none',
                sampling: 'lttb',
                // itemStyle: {
                //     color: 'rgb(255, 70, 131)'
                // },
                data: data,
                lineStyle: {
                    color: chartData["AnalogData"][0]["Color"],
                    width: 1,
                },
            }
        ]
    };
    console.log("PROCESSED " + new Date().toISOString())
    option && myChart.setOption(option);
    // let x = []
    // let y = []
    // for (const element of intermediateData["PortData"]) {
    //     x.push(new Date(element["Time"]*1000))
    //     if (element["Value"] === -32768) {
    //         y.push(null)
    //     } else {
    //         y.push(element["Value"])
    //     }
    // }
    // let trace1 = {
    //     type: 'scatter',
    //     mode: "lines",
    //
    //     x: x,
    //     y: y,
    //     marker: {
    //         color: 'green',
    //         line: {
    //             width: 2.5
    //         }
    //     },
    //     name: "test1"
    // };
    //
    // let data = [trace1];
    // let layout = {
    //     font: {size: 10, family: 'ProximaNova'},
    //     xaxis: {
    //         rangeslider: {}
    //     },
    //     showlegend: true,
    //     legend: {
    //         orientation: 'h',
    //         yanchor: 'top',
    //         xanchor: 'center',
    //         y: 1,
    //         x: 0.5
    //     }
    // };
    // let config = {responsive: true}
    // Plotly.newPlot('chart', data, layout, config);
}






