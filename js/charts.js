let now = new Date();
now.setMinutes(now.getMinutes() - now.getTimezoneOffset());
document.getElementById('to-date').value = now.toISOString().slice(0, 16);
now.setHours(now.getHours() - 24);
document.getElementById('from-date').value = now.toISOString().slice(0, 16);

const dataOkButton = document.getElementById("data-ok-button")

dataOkButton.addEventListener("click", (event) => {
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
        am4core.addLicense("CH72348321");
        const chart = am4core.create("chart", am4charts.XYChart);
        chart.height = am4core.percent(70);
        let dateAxis = chart.xAxes.push(new am4charts.DateAxis());
        dateAxis.groupData = true;
        dateAxis.groupCount = 8640;
        let valueAxis = chart.yAxes.push(new am4charts.ValueAxis());
        valueAxis.min = 0;
        let start = 0;
        let end = 0
        for (const analogData of chartData["AnalogData"]) {
            const series = chart.series.push(new am4charts.LineSeries());
            series.connect = false;
            series.dataFields.valueY = "value";
            series.dataFields.dateX = "date";
            series.name = analogData["PortName"];
            series.stroke = am4core.color(analogData["PortColor"]);
            let segment = series.segments.template;
            segment.interactionsEnabled = true;
            const data = [];
            for (const oneData of analogData["PortData"]) {
                let dataItem;
                if (oneData["Value"] === -32768) {
                    dataItem = {date: oneData["Time"] * 1000};
                } else {
                    dataItem = {date: oneData["Time"] * 1000, value: oneData["Value"]};
                    if (oneData["Value"] < valueAxis.min) {
                        valueAxis.min = oneData["Value"];
                    }
                }
                if (start === 0) {
                    start = oneData["Time"]
                }
                end = oneData["Time"]
                data.push(dataItem);
            }
            series.tooltipText = "{valueY}";
            series.data = data;
        }
        chart.height=am4core.percent(100)
        chart.legend = new am4charts.Legend();
        chart.legend.position = "top";
        chart.legend.scrollable = true;
        chart.cursor = new am4charts.XYCursor();
        chart.cursor.xAxis = dateAxis;
        chart.leftAxesContainer.width = 50;
        chart.rightAxesContainer.width = 50;

        switch (chartData["Locale"]) {
            case "CsCZ":
                chart.dateFormatter.language.locale = am4lang_cs_CZ;
                break;
            case "DeDE":
                chart.dateFormatter.language.locale = am4lang_de_DE;
                break;
            case "EsES":
                chart.dateFormatter.language.locale = am4lang_es_ES;
                break;
            case "FrFR":
                chart.dateFormatter.language.locale = am4lang_fr_FR;
                break;
            case "ItIT":
                chart.dateFormatter.language.locale = am4lang_it_IT;
                break;
            case "PlPL":
                chart.dateFormatter.language.locale = am4lang_pl_PL;
                break;
            case "PtPT":
                chart.dateFormatter.language.locale = am4lang_pt_PT;
                break;
            case "SkSK":
                chart.dateFormatter.language.locale = am4lang_cs_CZ;
                break;
            case "RuRU":
                chart.dateFormatter.language.locale = am4lang_sr_RS;
                break;
        }


        let terminalData = am4core.create("order-chart", am4charts.XYChart);
        terminalData.dateFormatter.dateFormat = "yyyy-MM-dd hh:mm:ss";
        terminalData.data = chartData["OrderData"]
        let terminalCategoryAxis = terminalData.yAxes.push(new am4charts.CategoryAxis());
        terminalCategoryAxis.dataFields.category = "Name";
        terminalCategoryAxis.renderer.grid.template.location = 0;
        terminalCategoryAxis.renderer.labels.template.disabled = true;
        terminalData.leftAxesContainer.width = 50;
        terminalData.rightAxesContainer.width = 50;

        let terminalDateAxis = terminalData.xAxes.push(new am4charts.DateAxis());
        terminalDateAxis.min = start*1000
        terminalDateAxis.max = end*1000
        terminalDateAxis.baseInterval = { count: 1, timeUnit: "second" };
        let series1 = terminalData.series.push(new am4charts.ColumnSeries());
        series1.columns.template.tooltipText = "{DataName}: {openDateX} - {dateX}";
        series1.dataFields.openDateX = "FromDate";
        series1.dataFields.dateX = "ToDate";
        series1.dataFields.categoryY = "Name";
        series1.columns.template.propertyFields.fill = "Color"; // get color from data
        series1.columns.template.propertyFields.stroke = "Color";
        series1.columns.template.strokeOpacity = 1;

        let cellSize = 10;
        terminalData.events.on("datavalidated", function(ev) {
            let chart = ev.target;
            let categoryAxis = chart.yAxes.getIndex(0);
            let adjustHeight = chart.data.length * cellSize - categoryAxis.pixelHeight;
            let targetHeight = chart.pixelHeight + adjustHeight;
            terminalData.svgContainer.htmlElement.style.height = targetHeight + "px";
        });
        let scrollbarX = new am4core.Scrollbar();
        scrollbarX.marginBottom = 20;
        scrollbarX.updateWhileMoving = false
        chart.scrollbarX = scrollbarX


        dateAxis.events.on("selectionextremeschanged", dateAxisChanged);
        function dateAxisChanged(ev) {
            terminalDateAxis.min = dateAxis._minZoomed;
            terminalDateAxis.max = dateAxis._maxZoomed;
        }
    });


    am4core.options.autoDispose = true;
}