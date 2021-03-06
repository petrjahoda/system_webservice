let now = new Date();
now.setMinutes(now.getMinutes() - now.getTimezoneOffset());
document.getElementById('to-date').value = now.toISOString().slice(0, 16);
now.setMonth(now.getMonth() - 1);
document.getElementById('from-date').value = now.toISOString().slice(0, 16);

const tableButton = document.getElementById("table-button")
tableButton.addEventListener("click", () => {
    let data
    if (tableButton.classList.contains("mif-menu")) {
        tableButton.classList.remove("mif-menu")
        tableButton.classList.add("mif-lines")
        document.getElementById("data-table").classList.add("compact")
        data = {
            email: document.getElementById("user-info").title,
            key: "data-selected-size",
            value: "compact"
        };

    } else {
        tableButton.classList.remove("mif-lines")
        tableButton.classList.add("mif-menu")
        document.getElementById("data-table").classList.remove("compact")
        data = {
            email: document.getElementById("user-info").title,
            key: "data-selected-size",
            value: ""
        };

    }
    fetch("/update_user_web_settings_from_web", {
        method: "POST",
        body: JSON.stringify(data)
    }).then(() => {
    }).catch(() => {
    });
})

const dataSelection = document.getElementById("data-selection")
dataSelection.addEventListener("change", () => {
    loadData();
})

const dataOkButton = document.getElementById("data-ok-button")
dataOkButton.addEventListener("click", async () => {
    const workplacesElement = document.getElementsByClassName("tag short-tag");
    let workplaces = ""
    for (let index = 0; index < workplacesElement.length; index++) {
        workplaces += workplacesElement[index].children[0].innerHTML + ";"
    }
    workplaces = workplaces.slice(0, -1)
    let data = {
        email: document.getElementById("user-info").title,
        key: "data-selected-workplaces",
        value: workplaces,
    };
    fetch("/update_user_web_settings_from_web", {
        method: "POST",
        body: JSON.stringify(data)
    }).then(() => {
        loadData();
    }).catch(() => {
    });
})

function loadData() {
    dataOkButton.disabled = true
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
                JSON.parse(data);
            } else {
                document.getElementById("data-table-container").innerHTML = data
                if (document.getElementById("data-table").classList.contains("compact")) {
                    tableButton.classList.remove("mif-menu")
                    tableButton.classList.add("mif-lines")
                } else {
                    tableButton.classList.remove("mif-lines")
                    tableButton.classList.add("mif-menu")
                }
            }
            document.getElementById("loader").hidden = true
            dataOkButton.disabled = false
        });
    }).catch(() => {
        document.getElementById("loader").hidden = true
        dataOkButton.disabled = false
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
