const dataSelection = document.getElementById("data-selection")
dataSelection.addEventListener("change", (event) => {
    loadSettings();
    document.getElementById("settings-container-detail").innerHTML = ""
})

const container = document.getElementById("settings-container")
container.addEventListener("click", (event) => {
    let table = Metro.getPlugin("#data-table", "table");
    if (table.getSelectedItems().length > 0 && event.target.type === "radio") {
        let selectedItem = table.getSelectedItems()[0][0]
        console.log("Loading details for " + selectedItem)
        loadDetails(selectedItem);
    } else if (event.target.id === "data-new-button" || event.target.id === "data-new-button-mif") {
        loadSettings();
        loadDetails();
    }
})

const containerDetail = document.getElementById("settings-container-detail")


containerDetail.addEventListener("click", (event) => {
    if (event.target.id === "data-save-button" || event.target.id === "data-save-button-mif") {
        let selection = document.getElementById("data-selection").value
        console.log(selection)
        switch (selection) {
            case "alarms" : {
                saveAlarm();
                break
            }
            case "operations" : {
                saveOperation();
                break
            }
            case "orders" : {
                saveOrder();
                break;
            }
        }
    }
})

function saveOrder() {
    if (document.getElementById("order-name").value.length === 0) {
        document.getElementById("order-name").style.backgroundColor = "#ffcccb"
    } else if (!Date.parse(document.getElementById("date-time-request").value)) {
        document.getElementById("date-time-request").style.backgroundColor = "#ffcccb"
    } else {
        document.getElementById("order-name").style.backgroundColor = ""
        document.getElementById("date-time-request").style.backgroundColor = ""
        if (document.getElementById("data-save-button").classList[1] === "primary") {
            document.getElementById("data-save-button").classList.remove("primary")
            document.getElementById("data-save-button").classList.add("alert")
            document.getElementById("data-save-button-mif").classList.remove("mif-floppy-disk")
            document.getElementById("data-save-button-mif").classList.add("mif-cross")

            setTimeout(function () {
                if (document.getElementById("data-save-button").classList[1] === "alert") {
                    let parseId = ""
                    if (Metro.getPlugin("#data-table", "table").getSelectedItems().length > 0) {
                        parseId = Metro.getPlugin("#data-table", "table").getSelectedItems()[0][0]
                    }
                    let data = {
                        id: parseId,
                        name: document.getElementById("order-name").value,
                        product: document.getElementById("product-selection").value,
                        workplace: document.getElementById("workplace-selection").value,
                        countRequest: document.getElementById("count-request").value,
                        dateTimeRequest: document.getElementById("date-time-request").value,
                        cavity: document.getElementById("cavity").value,
                        barcode: document.getElementById("order-barcode").value,
                        note: document.getElementById("order-note").value,
                    };
                    fetch("/save_order", {
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
            }, 2500);
        } else if (document.getElementById("data-save-button").classList[1] === "alert") {
            document.getElementById("data-save-button").classList.remove("alert")
            document.getElementById("data-save-button").classList.add("primary")
            document.getElementById("data-save-button-mif").classList.remove("mif-cross")
            document.getElementById("data-save-button-mif").classList.add("mif-floppy-disk")
        }
    }
}


function saveOperation() {
    if (document.getElementById("operation-name").value.length === 0) {
        document.getElementById("operation-name").style.backgroundColor = "#ffcccb"
    } else {
        document.getElementById("operation-name").style.backgroundColor = ""
        if (document.getElementById("data-save-button").classList[1] === "primary") {
            document.getElementById("data-save-button").classList.remove("primary")
            document.getElementById("data-save-button").classList.add("alert")
            document.getElementById("data-save-button-mif").classList.remove("mif-floppy-disk")
            document.getElementById("data-save-button-mif").classList.add("mif-cross")

            setTimeout(function () {
                if (document.getElementById("data-save-button").classList[1] === "alert") {
                    let parseId = ""
                    if (Metro.getPlugin("#data-table", "table").getSelectedItems().length > 0) {
                        parseId = Metro.getPlugin("#data-table", "table").getSelectedItems()[0][0]
                    }
                    let data = {
                        id: parseId,
                        name: document.getElementById("operation-name").value,
                        order: document.getElementById("order-selection").value,
                        barcode: document.getElementById("operation-barcode").value,
                        note: document.getElementById("operation-note").value,
                    };
                    fetch("/save_operation", {
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
            }, 2500);
        } else if (document.getElementById("data-save-button").classList[1] === "alert") {
            document.getElementById("data-save-button").classList.remove("alert")
            document.getElementById("data-save-button").classList.add("primary")
            document.getElementById("data-save-button-mif").classList.remove("mif-cross")
            document.getElementById("data-save-button-mif").classList.add("mif-floppy-disk")
        }
    }
}

function saveAlarm() {
    if (document.getElementById("alarm-name").value.length === 0) {
        document.getElementById("alarm-name").style.backgroundColor = "#ffcccb"
    } else {
        document.getElementById("alarm-name").style.backgroundColor = ""
        if (document.getElementById("data-save-button").classList[1] === "primary") {
            document.getElementById("data-save-button").classList.remove("primary")
            document.getElementById("data-save-button").classList.add("alert")
            document.getElementById("data-save-button-mif").classList.remove("mif-floppy-disk")
            document.getElementById("data-save-button-mif").classList.add("mif-cross")

            setTimeout(function () {
                if (document.getElementById("data-save-button").classList[1] === "alert") {
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
                    fetch("/save_alarm", {
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
            }, 2500);
        } else if (document.getElementById("data-save-button").classList[1] === "alert") {
            document.getElementById("data-save-button").classList.remove("alert")
            document.getElementById("data-save-button").classList.add("primary")
            document.getElementById("data-save-button-mif").classList.remove("mif-cross")
            document.getElementById("data-save-button-mif").classList.add("mif-floppy-disk")
        }
    }
}

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

function loadSettings() {
    let data = {
        data: document.getElementById("data-selection").value,
    };
    fetch("/load_settings_data", {
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
    fetch("/load_settings_detail", {
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