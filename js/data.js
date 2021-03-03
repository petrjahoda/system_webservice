
let dt = new Date();
dt = addDate(dt, -1, 'months');
const datePicker = document.getElementById("fromdate")
datePicker.dataset.value = dt.format("%Y-%m-%d")
const dataOkButton = document.getElementById("data-ok-button")

dataOkButton.addEventListener("click", (event) => {
    console.log("getting data for " + document.getElementById("data-selection").value)
    const workplacesElement = document.getElementsByClassName("tag short-tag");
    let workplaces = []
    for (let index = 0; index < workplacesElement.length; index++) {
        workplaces.push(workplacesElement[index].children[0].innerHTML)
    }
    let data = {
        data: document.getElementById("data-selection").value,
        workplaces: workplaces,
        from: document.getElementById("fromdate").value + ";" + document.getElementById("fromtime").value,
        to: document.getElementById("todate").value + ";" + document.getElementById("totime").value
    };
    fetch("/get_table_data", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            document.getElementById("data-table-container").innerHTML = data
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