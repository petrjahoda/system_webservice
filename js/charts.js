
let chartIsLoading = false
const workplaceSelection = document.getElementById("workplace-selection")
workplaceSelection.addEventListener("change", () => {
    if (!chartIsLoading) {
        loadChart();
    }
})

const dataSelection = document.getElementById("data-selection")
dataSelection.addEventListener("change", () => {
    if (!chartIsLoading) {
        loadChart();
    }
})

let chartDom = document.getElementById('chart');
let chartHeight = document.documentElement.clientHeight * 0.9
let chartWidth = document.documentElement.clientWidth
if (!document.getElementById("mainmenu").classList.contains("compacted")) {
    chartWidth = document.documentElement.clientWidth*0.86
}
if (chartHeight < 800) {
    chartHeight = 800;
}
let startDateAsValue = new Date()
let endDateAsValue = new Date()
let borderStart = 50
let borderEnd = chartDom.scrollWidth - 80
console.log(chartDom.scrollWidth)
console.log(chartDom.offsetWidth)
console.log(chartDom.clientWidth)
console.log(document.getElementById("chart-container").clientWidth)
console.log(document.getElementById("chart-container").scrollWidth)
console.log(document.getElementById("chart-container").offsetWidth)
console.log(borderStart)
console.log(borderEnd)
let borderChange = borderEnd - borderStart
let myChart = echarts.init(chartDom, null, {height: chartHeight,width:chartWidth, renderer: 'svg'});
if (document.getElementById('to-date').value === "") {
    let now = new Date();
    now.setMinutes(now.getMinutes() - now.getTimezoneOffset());
    document.getElementById('to-date').value = now.toISOString().slice(0, 16);
    now.setHours(now.getHours() - 24);
    document.getElementById('from-date').value = now.toISOString().slice(0, 16);
}
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

function loadChart() {
    chartIsLoading = true
    document.getElementById("loader").hidden = false
    let flashData = "mif-flash-off"
    let terminalData = "mif-phonelink"
    if (phoneLinkButton.classList.contains("mif-phonelink-off")) {
        terminalData = "mif-phonelink-off"
    }
    if (flashButton.classList.contains("mif-flash-on")) {
        flashData = "mif-flash-on"
    }
    let data = {
        data: document.getElementById("data-selection").value,
        workplace: document.getElementById("workplace-selection").value,
        from: document.getElementById("from-date").value,
        to: document.getElementById("to-date").value,
        flash: flashData,
        terminal: terminalData,
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
            chartIsLoading = false
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
            } else if (result["Type"] === "combined-chart") {
                myChart.clear()
                if (result["ChartData"] !== null) {
                    drawCombinedChart(result)
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
}

dataOkButton.addEventListener("click", () => {
    loadChart();
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
myChart.on('datazoom', function () {
    let zoom = myChart.getOption().dataZoom[0];
    startDateAsValue = zoom.startValue * 1000
    endDateAsValue = zoom.endValue * 1000
    console.log(myChart)
});
