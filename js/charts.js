let dt = new Date();
dt = addDate(dt, -1, 'days');
const datePicker = document.getElementById("fromdate")
datePicker.dataset.value = dt.format("%Y-%m-%d")
const dataOkButton = document.getElementById("data-ok-button")

dataOkButton.addEventListener("click", (event) => {
    console.log("getting chart data for " + document.getElementById("data-selection").value)
    console.log("getting chart data for " + document.getElementById("workplace-selection").value)
    let data = {
        data: document.getElementById("data-selection").value,
        workplace: document.getElementById("workplace-selection").value,
        from: document.getElementById("fromdate").value + ";" + document.getElementById("fromtime").value,
        to: document.getElementById("todate").value + ";" + document.getElementById("totime").value
    };
    fetch("/get_chart_data", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            let result = JSON.parse(data);
            if (result["Type"] === "analog-data") {
                drawAnalogChart(result)
            }

        });
    }).catch((error) => {
        console.log(error)
    });
})

function addDate(dt, amount, dateType) {
    switch (dateType) {
        case 'days':
            return dt.setDate(dt.getDate() + amount) && dt;
        case 'weeks':
            return dt.setDate(dt.getDate() + (7 * amount)) && dt;
        case 'months':
            return dt.setMonth(dt.getMonth() + amount) && dt;
        case 'years':
            return dt.setFullYear(dt.getFullYear() + amount) && dt;
    }
}

function drawAnalogChart(chartData) {
    am4core.ready(function () {
        let dataItem;
        am4core.addLicense("CH72348321");
        const chart = am4core.create("chart", am4charts.XYChart);
        let dateAxis = chart.xAxes.push(new am4charts.DateAxis());
        dateAxis.groupData = true;
        dateAxis.groupCount = 17280;
        let valueAxis = chart.yAxes.push(new am4charts.ValueAxis());
        valueAxis.min = 0;

        for (const analogData of chartData["AnalogData"]) {
            console.log(analogData["PortName"] + " started")
            const series = chart.series.push(new am4charts.LineSeries());
            series.dataFields.valueY = "value";
            series.dataFields.dateX = "date";
            series.name = analogData["PortName"];
            let segment = series.segments.template;
            segment.interactionsEnabled = true;
            const data = [];
            for (const oneData of analogData["PortData"]) {
                dataItem = {date: oneData["Time"]*1000, value:oneData["Value"]};
                data.push(dataItem);
            }
            console.log(analogData["PortName"] + " ended with length of " + data.length)
            series.tooltipText = "{valueY}";
            series.data = data;
        }
        chart.legend = new am4charts.Legend();
        chart.legend.position = "right";
        chart.legend.scrollable = true;

        chart.cursor = new am4charts.XYCursor();
        chart.cursor.xAxis = dateAxis;

        let scrollbarX = new am4core.Scrollbar();
        scrollbarX.marginBottom = 20;
        chart.scrollbarX = scrollbarX;
        chart.dateFormatter.language.locale = am4lang_cs_CZ;
    });
}