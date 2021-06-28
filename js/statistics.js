if (document.getElementById('to-date').value === "") {
}
let now = new Date();
now.setMinutes(now.getMinutes() - now.getTimezoneOffset());
document.getElementById('to-date').value = now.toISOString().slice(0, 16);
now.setMonth(now.getMonth() - 1);
document.getElementById('from-date').value = now.toISOString().slice(0, 16);

const dataSelection = document.getElementById("statistics-selection")
dataSelection.addEventListener("change", () => {
    let data = {
        email: document.getElementById("user-info").title,
        selection: document.getElementById("statistics-selection").value,
    };
    fetch("/load_types_for_selected_statistics", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            document.getElementById("types").innerHTML = data
        });
        loadStatisticsChartData()
    }).catch(() => {
    });
})

const dataOkButton = document.getElementById("statistics-ok-button")
dataOkButton.addEventListener("click", () => {
    const workplacesElement = document.getElementById("workplaces").getElementsByClassName("tag short-tag");
    let workplaces = ""
    for (let index = 0; index < workplacesElement.length; index++) {
        workplaces += workplacesElement[index].children[0].innerHTML + ";"
    }
    workplaces = workplaces.slice(0, -1)
    let data = {
        email: document.getElementById("user-info").title,
        key: "statistics-selected-workplaces",
        value: workplaces,
    };
    fetch("/update_user_web_settings_from_web", {
        method: "POST",
        body: JSON.stringify(data)
    }).then(() => {
        let typesElement = document.getElementById("types").getElementsByClassName("tag short-tag")
        let types = ""
        for (let index = 0; index < typesElement.length; index++) {
            types += typesElement[index].children[0].innerHTML + ";"
        }
        types = types.slice(0, -1)
        let data = {
            email: document.getElementById("user-info").title,
            key: "statistics-selected-types-" + document.getElementById("statistics-selection").value,
            value: types,
        };
        fetch("/update_user_web_settings_from_web", {
            method: "POST",
            body: JSON.stringify(data)
        }).then(() => {
            let usersElement = document.getElementById("users").getElementsByClassName("tag short-tag")
            let users = ""
            for (let index = 0; index < usersElement.length; index++) {
                users += usersElement[index].children[0].innerHTML + ";"
            }
            users = users.slice(0, -1)
            let data = {
                email: document.getElementById("user-info").title,
                key: "statistics-selected-users",
                value: users,
            };
            fetch("/update_user_web_settings_from_web", {
                method: "POST",
                body: JSON.stringify(data)
            }).then(() => {
                loadStatisticsChartData();
                document.getElementById("loader").hidden = true
            }).catch(() => {
                document.getElementById("loader").hidden = true
            });
        }).catch(() => {
            document.getElementById("loader").hidden = true
        });
    }).catch(() => {
        document.getElementById("loader").hidden = true
    });
})

function loadStatisticsChartData() {
    document.getElementById("loader").hidden = false
    let data = {
        data: document.getElementById("statistics-selection").value,
        from: document.getElementById("from-date").value,
        to: document.getElementById("to-date").value
    };
    fetch("/load_statistics_data", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            try {
                let result = JSON.parse(data);
                document.getElementById("lower-content").innerText = JSON.stringify(result, null, "\t")
            } catch {
            }
            document.getElementById("loader").hidden = true
        });
    }).catch(() => {
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
