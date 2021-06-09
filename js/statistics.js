if (document.getElementById('to-date').value === "") {
}
let now = new Date();
now.setMinutes(now.getMinutes() - now.getTimezoneOffset());
document.getElementById('to-date').value = now.toISOString().slice(0, 16);
now.setMonth(now.getMonth() - 1);
document.getElementById('from-date').value = now.toISOString().slice(0, 16);

const dataSelection = document.getElementById("statistics-selection")
dataSelection.addEventListener("change", () => {
    loadData();
})

const dataOkButton = document.getElementById("statistics-ok-button")
dataOkButton.addEventListener("click", () => {
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
        data: document.getElementById("statistics-selection").value,
        workplaces: workplaces,
        from: document.getElementById("from-date").value,
        to: document.getElementById("to-date").value
    };
    fetch("/load_statistics_data", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            if (data.includes("ERR: ")) {
                JSON.parse(data);
            } else {
                // document.getElementById("data-table-container").innerHTML = data
                // if (document.getElementById("data-table").classList.contains("compact")) {
                //     tableButton.classList.remove("mif-menu")
                //     tableButton.classList.add("mif-lines")
                // } else {
                //     tableButton.classList.remove("mif-lines")
                //     tableButton.classList.add("mif-menu")
                // }
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

document.addEventListener('DOMContentLoaded', loadData, false);
