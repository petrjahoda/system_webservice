let now = new Date();
now.setMinutes(now.getMinutes() - now.getTimezoneOffset());
document.getElementById('to-date').value = now.toISOString().slice(0, 16);
now.setMonth(now.getMonth() - 1);
document.getElementById('from-date').value = now.toISOString().slice(0, 16);


const tableButton = document.getElementById("table-button")
tableButton.addEventListener("click", (event) => {
    if (tableButton.classList.contains("mif-menu")) {
        tableButton.classList.remove("mif-menu")
        tableButton.classList.add("mif-lines")
        document.getElementById("data-table").classList.add("compact")
    } else {
        tableButton.classList.remove("mif-lines")
        tableButton.classList.add("mif-menu")
        document.getElementById("data-table").classList.remove("compact")
    }
})

const dataSelection = document.getElementById("data-selection")
dataSelection.addEventListener("change", (event) => {
    loadData();
})

const dataOkButton = document.getElementById("data-ok-button")
dataOkButton.addEventListener("click", (event) => {
    loadData();
})

function loadData() {
    document.getElementById("loader").hidden = false
    const workplacesElement = document.getElementsByClassName("tag short-tag");
    let workplaces = []
    for (let index = 0; index < workplacesElement.length; index++) {
        workplaces.push(workplacesElement[index].children[0].innerHTML)
    }
    let data = {
        data: document.getElementById("data-selection").value,
        workplaces: workplaces,
        from: document.getElementById("from-date").value,
        to: document.getElementById("to-date").value
    };
    fetch("/load_table_data", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            if (data.includes("ERR: ")) {
                let result = JSON.parse(data);
                updateCharm(result["Result"])
            } else {
                document.getElementById("data-table-container").innerHTML = data
                updateCharm(document.getElementById("hidden-data-information").innerText)
            }
            document.getElementById("loader").hidden = true
        });
    }).catch((error) => {
        console.log(error)
        document.getElementById("loader").hidden = true
    });
}

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

function load() {
    const workplacesElement = document.getElementsByClassName("tag short-tag");
    let workplaces = []
    for (let index = 0; index < workplacesElement.length; index++) {
        workplaces.push(workplacesElement[index].children[0].innerHTML)
    }
    let data = {
        data: document.getElementById("data-selection").value,
        workplaces: workplaces,
        from: document.getElementById("from-date").value,
        to: document.getElementById("to-date").value
    };
    fetch("/load_table_data", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            document.getElementById("data-table-container").innerHTML = data
        });
    }).catch((error) => {
        console.log(error)
    });
}

document.addEventListener('DOMContentLoaded', load, false);