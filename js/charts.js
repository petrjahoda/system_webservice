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
                drawAnalogChart(result)
                document.getElementById("loader").hidden = true
            }
        });
    }).catch((error) => {
        updateCharm("ERR: " + error)
        document.getElementById("loader").hidden = true
    });
})

function drawAnalogChart(chartData) {
    const intermediateData = chartData["AnalogData"][0]
    let x = []
    let y = []
    for (const element of intermediateData["PortData"]) {
        x.push(element["Time"])
        if (element["Value"] === -32768) {
            y.push(null)
        } else {
            y.push(element["Value"])
        }
    }
    let trace1 = {
        type: 'scatter',
        mode: "lines",

        x: x,
        y: y,
        marker: {
            color: 'green',
            line: {
                width: 2.5
            }
        },
        name: "test1"
    };

    let data = [trace1];
    let layout = {
        font: {size: 10, family: 'ProximaNova'},
        xaxis: {
            rangeslider: {}
        },
        showlegend: true,
        legend: {
            orientation: 'h',
            yanchor: 'top',
            xanchor: 'center',
            y: 1,
            x: 0.5
        }
    };
    let config = {responsive: true}
    Plotly.newPlot('chart', data, layout, config);
}




