const dataSelection = document.getElementById("data-selection")
dataSelection.addEventListener("change", (event) => {
    loadSettings();
})

const container = document.getElementById("settings-container")
container.addEventListener("click", (event) => {
    let table = Metro.getPlugin("#data-table", "table");
    if (table.getSelectedItems().length > 0 && event.target.type === "radio") {
        let selectedItem = table.getSelectedItems()[0][0]
        loadDetails(selectedItem);
    } else if (event.target.id === "data-new-button" || event.target.id === "data-new-button-mif") {
        loadSettings();
        loadDetails();
    }
})

const containerDetail = document.getElementById("settings-container-detail")
containerDetail.addEventListener("click", (event) => {
    if (event.target.id === "data-save-button" || event.target.id === "data-save-button-mif") {
        if (document.getElementById("alarm-name").value.length === 0) {
            for (const element of document.getElementsByClassName("important")) {
                element.style = "border:1px solid red"
            }
        } else {
            for (const element of document.getElementsByClassName("important")) {
                element.style = ""
            }
            if (document.getElementById("data-save-button").classList[1] === "primary") {
                document.getElementById("data-save-button").classList.remove("primary")
                document.getElementById("data-save-button").classList.add("alert")
                document.getElementById("data-save-button-mif").classList.remove("mif-floppy-disk")
                document.getElementById("data-save-button-mif").classList.add("mif-cross")
                let selection = document.getElementById("data-selection").value
                console.log("Saving: " + selection)
                setTimeout(function () {
                    if (document.getElementById("data-save-button").classList[1] === "alert") {
                        switch (selection) {
                            case "alarms" : {
                                saveAlarm();
                                break
                            }
                            default: {
                                console.log("default")
                            }
                        }
                    }
                }, 2500);
            } else if (document.getElementById("data-save-button").classList[1] === "alert") {
                document.getElementById("data-save-button").classList.remove("alert")
                document.getElementById("data-save-button").classList.add("primary")
                document.getElementById("data-save-button-mif").classList.remove("mif-cross")
                document.getElementById("data-save-button-mif").classList.add("mif-floppy-disk")
            }
        }
    }
})


function saveAlarm() {
    let parseId = ""
    if (Metro.getPlugin("#data-table", "table").getSelectedItems().length > 0) {
        parseId = Metro.getPlugin("#data-table", "table").getSelectedItems()[0][0]
    }
    let data = {
        id: parseId,
        name: document.getElementById("alarm-name").value,
        workplace: document.getElementById("workplace-selection").value,
        sql: document.getElementById("sql-command").value,
        header: document.getElementById("message-header").value,
        text: document.getElementById("message-text").value,
        recipients: document.getElementById("recipients").value,
        url: document.getElementById("url").value,
        pdf: document.getElementById("pdf").value,
    };
    fetch("/save_detail_alarm", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            document.getElementById("settings-container-detail").innerHTML = ""
            loadSettings();
        });
    }).catch((error) => {
        console.log(error)
        document.getElementById("settings-container-detail").innerHTML = ""
        loadSettings();
    });
}

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

function loadSettings() {
    console.log("getting settings for " + document.getElementById("data-selection").value)
    let data = {
        data: document.getElementById("data-selection").value,
    };
    fetch("/get_settings_data", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            document.getElementById("settings-container").innerHTML = data
        });
    }).catch((error) => {
        console.log(error)
    });
}

function loadDetails(selectedItem) {
    let selection = document.getElementById("data-selection").value
    let data = {
        data: selection,
        name: selectedItem
    };
    fetch("/get_detail_settings", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            document.getElementById("settings-container-detail").innerHTML = data
        });
    }).catch((error) => {
        console.log(error)
    });
}