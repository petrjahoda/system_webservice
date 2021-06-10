
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
document.getElementById('chart').style.height = document.documentElement.clientHeight * 0.9 + 'px'
let startDateAsValue = new Date()
let endDateAsValue = new Date()
let myChart = echarts.init(chartDom, null, {renderer: 'svg'});
if (document.getElementById('to-date').value === "") {
}
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

function loadChart() {
    chartDom.hidden = true
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
                setInterval(function () {
                    myChart.resize()
                    chartDom.hidden = false
                }, 1);
                document.getElementById("loader").hidden = true
            }

        });
    }).catch((error) => {
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

myChart.on('datazoom', function () {
    let zoom = myChart.getOption().dataZoom[0];
    startDateAsValue = zoom.startValue * 1000
    endDateAsValue = zoom.endValue * 1000
});

window.addEventListener('resize', () => {
    myChart.resize()
})
