let now = new Date();
now.setMinutes(now.getMinutes() - now.getTimezoneOffset());
document.getElementById('to-date').value = now.toISOString().slice(0,16);
now.setHours(now.getHours() - 24);
document.getElementById('from-date').value = now.toISOString().slice(0,16);

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
        dateAxis.groupCount = 3840;
        let valueAxis = chart.yAxes.push(new am4charts.ValueAxis());
        valueAxis.min = 0;
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
                if (oneData["Value"] === -32768) {
                    dataItem = {date: oneData["Time"]*1000};
                } else {
                    dataItem = {date: oneData["Time"]*1000, value:oneData["Value"]};
                    if (oneData["Value"] < valueAxis.min) {
                        valueAxis.min = oneData["Value"];
                    }
                }
                data.push(dataItem);
            }
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
        switch (chartData["Locale"]) {
            case "CsCZ": chart.dateFormatter.language.locale = am4lang_cs_CZ;break;
            case "DeDE": chart.dateFormatter.language.locale = am4lang_de_DE;break;
            case "EsES": chart.dateFormatter.language.locale = am4lang_es_ES;break;
            case "FrFR": chart.dateFormatter.language.locale = am4lang_fr_FR;break;
            case "ItIT": chart.dateFormatter.language.locale = am4lang_it_IT;break;
            case "PlPL": chart.dateFormatter.language.locale = am4lang_pl_PL;break;
            case "PtPT": chart.dateFormatter.language.locale = am4lang_pt_PT;break;
            case "SkSK": chart.dateFormatter.language.locale = am4lang_cs_CZ;break;
            case "RuRU": chart.dateFormatter.language.locale = am4lang_sr_RS;break;
        }
    });
    am4core.options.autoDispose = true;
}