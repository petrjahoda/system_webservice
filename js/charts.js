let chartDom = document.getElementById('chart');
let chartHeight = document.documentElement.clientHeight * 0.9
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
flashButton.addEventListener("click", () => {
    if (flashButton.classList.contains("mif-flash-on")) {
        flashButton.classList.remove("mif-flash-on")
        flashButton.classList.add("mif-flash-off")
    } else {
        flashButton.classList.remove("mif-flash-off")
        flashButton.classList.add("mif-flash-on")
    }
})

const phoneLinkButton = document.getElementById("phonelink-button")
phoneLinkButton.addEventListener("click", () => {
    if (phoneLinkButton.classList.contains("mif-phonelink-off")) {
        phoneLinkButton.classList.remove("mif-phonelink-off")
        phoneLinkButton.classList.add("mif-phonelink")
    } else {
        phoneLinkButton.classList.remove("mif-phonelink")
        phoneLinkButton.classList.add("mif-phonelink-off")
    }
})

const dataOkButton = document.getElementById("data-ok-button")
dataOkButton.addEventListener("click", () => {
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
            } else if (result["Type"] === "production-chart") {
                myChart.clear()
                if (result["ChartData"] !== null) {
                    drawProductionChart(result)
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

window.addEventListener('resize', () => {
    myChart.resize()
})