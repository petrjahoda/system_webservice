const dataSelection = document.getElementById("data-selection")
dataSelection.addEventListener("change", (event) => {
    loadSettings();
    document.getElementById("settings-container-detail").innerHTML = ""
})

const container = document.getElementById("settings-container")
container.addEventListener("click", (event) => {
    const tableSelectedId = event.target.parentElement.parentElement.parentElement.parentElement.parentElement.id
    let table = Metro.getPlugin("#data-table", "table");
    let typeTable = Metro.getPlugin("#type-table", "table");
    if (tableSelectedId === "data-table" && table.getSelectedItems().length > 0 && event.target.type === "radio") {
        let selectedItem = table.getSelectedItems()[0][0]
        loadDetails(selectedItem, false);
    } else if (tableSelectedId === "type-table" && typeTable.getSelectedItems().length > 0 && event.target.type === "radio") {
        let selectedItem = typeTable.getSelectedItems()[0][0]
        loadDetails(selectedItem, true);
    } else if (event.target.id === "data-new-button" || event.target.id === "data-new-button-mif") {
        loadSettings();
        loadDetails(null, false);
    } else if (event.target.id === "data-new-button-type" || event.target.id === "data-new-button-mif-type") {
        loadSettings();
        loadDetails(null, true);
    }
    if (event.target.id === "data-save-button") {
        let selection = document.getElementById("data-selection").value
        switch (selection) {
            case "user" : {
                saveUserSettings();
                break
            }

        }
    }
})

const containerDetail = document.getElementById("settings-container-detail")
containerDetail.addEventListener("click", (event) => {
    if (event.target.id === "data-save-button" || event.target.id === "data-save-button-mif") {
        let selection = document.getElementById("data-selection").value
        console.log("Saving one of " + selection)
        switch (selection) {
            case "alarms" : {
                saveAlarm();
                break
            }
            case "breakdowns" : {
                saveBreakdown();
                break
            }
            case "downtimes" : {
                saveDowntime();
                break
            }
            case "faults" : {
                saveFault();
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
            case "products" : {
                saveProduct();
                break;
            }
            case "parts" : {
                savePart();
                break;
            }
            case "packages" : {
                savePackage();
                break;
            }
            case "states" : {
                saveState();
                break;
            }
            case "users" : {
                saveUser();
                break;
            }
            case "system-settings" : {
                saveSystemSettings();
                break;
            }
            case "workshifts" : {
                saveWorkshift();
                break;
            }
        }
    } else if (event.target.id === "data-type-save-button" || event.target.id === "data-type-save-button-mif") {
        let selection = document.getElementById("data-selection").value
        console.log("Saving type of " + selection)
        switch (selection) {
            case "breakdowns" : {
                saveBreakdownType();
                break
            }
            case "downtimes" : {
                saveDowntimeType();
                break
            }
            case "faults" : {
                saveFaultType();
                break
            }
            case "packages" : {
                savePackageType();
                break;
            }
            case "users" : {
                saveUserType();
                break;
            }
        }
    }
})

function saveSystemSettings() {
    if (document.getElementById("system-settings-name").value.length === 0) {
        document.getElementById("system-settings-name").style.backgroundColor = "#ffcccb"
    } else {
        document.getElementById("system-settings-name").style.backgroundColor = ""
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
                        name: document.getElementById("system-settings-name").value,
                        value: document.getElementById("system-settings-value").value,
                        enabled: document.getElementById("system-settings-selection").value,
                        note: document.getElementById("system-settings-note").value,
                    };
                    fetch("/save_system_settings", {
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

function saveUserSettings() {
    if (document.getElementById("first-name").value.length === 0 || document.getElementById("second-name").value.length === 0) {
        document.getElementById("first-name").style.backgroundColor = "#ffcccb"
        document.getElementById("second-name").style.backgroundColor = "#ffcccb"
    } else {
        document.getElementById("first-name").style.backgroundColor = ""
        document.getElementById("second-name").style.backgroundColor = ""
        if (document.getElementById("data-save-button").classList[1] === "primary") {
            document.getElementById("data-save-button").classList.remove("primary")
            document.getElementById("data-save-button").classList.add("alert")
            document.getElementById("data-save-button-mif").classList.remove("mif-floppy-disk")
            document.getElementById("data-save-button-mif").classList.add("mif-cross")

            setTimeout(function () {
                if (document.getElementById("data-save-button").classList[1] === "alert") {
                    let data = {
                        firstName: document.getElementById("first-name").value,
                        secondName: document.getElementById("second-name").value,
                        email: document.getElementById("email").value,
                        locale: document.getElementById("user-locale-selection").value,
                        password: document.getElementById("password").value,
                        phone: document.getElementById("phone").value,
                        note: document.getElementById("user-note").value,
                    };
                    fetch("/save_user_settings", {
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

function saveUserType() {
    if (document.getElementById("user-type-name").value.length === 0) {
        document.getElementById("user-type-name").style.backgroundColor = "#ffcccb"
    } else {
        document.getElementById("user-type-name").style.backgroundColor = ""
        if (document.getElementById("data-type-save-button").classList[1] === "primary") {
            document.getElementById("data-type-save-button").classList.remove("primary")
            document.getElementById("data-type-save-button").classList.add("alert")
            document.getElementById("data-type-save-button-mif").classList.remove("mif-floppy-disk")
            document.getElementById("data-type-save-button-mif").classList.add("mif-cross")

            setTimeout(function () {
                if (document.getElementById("data-type-save-button").classList[1] === "alert") {
                    let parseId = ""
                    if (Metro.getPlugin("#type-table", "table").getSelectedItems().length > 0) {
                        parseId = Metro.getPlugin("#type-table", "table").getSelectedItems()[0][0]
                    }
                    let data = {
                        id: parseId,
                        name: document.getElementById("user-type-name").value,
                        note: document.getElementById("user-type-note").value,
                    };
                    fetch("/save_user_type", {
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
        } else if (document.getElementById("data-type-save-button").classList[1] === "alert") {
            document.getElementById("data-type-save-button").classList.remove("alert")
            document.getElementById("data-type-save-button").classList.add("primary")
            document.getElementById("data-type-save-button-mif").classList.remove("mif-cross")
            document.getElementById("data-type-save-button-mif").classList.add("mif-floppy-disk")
        }
    }
}


function saveUser() {
    if (document.getElementById("first-name").value.length === 0 || document.getElementById("second-name").value.length === 0) {
        document.getElementById("first-name").style.backgroundColor = "#ffcccb"
        document.getElementById("second-name").style.backgroundColor = "#ffcccb"
    } else {
        document.getElementById("first-name").style.backgroundColor = ""
        document.getElementById("second-name").style.backgroundColor = ""
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
                        firstName: document.getElementById("first-name").value,
                        secondName: document.getElementById("second-name").value,
                        type: document.getElementById("user-type-selection").value,
                        role: document.getElementById("user-role-selection").value,
                        email: document.getElementById("email").value,
                        locale: document.getElementById("user-locale-selection").value,
                        barcode: document.getElementById("barcode").value,
                        password: document.getElementById("password").value,
                        phone: document.getElementById("phone").value,
                        pin: document.getElementById("pin").value,
                        position: document.getElementById("position").value,
                        rfid: document.getElementById("rfid").value,
                        note: document.getElementById("user-note").value,
                    };
                    fetch("/save_user", {
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

function savePackageType() {
    if (document.getElementById("package-type-name").value.length === 0) {
        document.getElementById("package-type-name").style.backgroundColor = "#ffcccb"
    } else {
        document.getElementById("package-type-name").style.backgroundColor = ""
        if (document.getElementById("data-type-save-button").classList[1] === "primary") {
            document.getElementById("data-type-save-button").classList.remove("primary")
            document.getElementById("data-type-save-button").classList.add("alert")
            document.getElementById("data-type-save-button-mif").classList.remove("mif-floppy-disk")
            document.getElementById("data-type-save-button-mif").classList.add("mif-cross")

            setTimeout(function () {
                if (document.getElementById("data-type-save-button").classList[1] === "alert") {
                    let parseId = ""
                    if (Metro.getPlugin("#type-table", "table").getSelectedItems().length > 0) {
                        parseId = Metro.getPlugin("#type-table", "table").getSelectedItems()[0][0]
                    }
                    let data = {
                        id: parseId,
                        name: document.getElementById("package-type-name").value,
                        count: document.getElementById("package-type-count").value,
                        note: document.getElementById("package-type-note").value,
                    };
                    fetch("/save_package_type", {
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
        } else if (document.getElementById("data-type-save-button").classList[1] === "alert") {
            document.getElementById("data-type-save-button").classList.remove("alert")
            document.getElementById("data-type-save-button").classList.add("primary")
            document.getElementById("data-type-save-button-mif").classList.remove("mif-cross")
            document.getElementById("data-type-save-button-mif").classList.add("mif-floppy-disk")
        }
    }
}


function savePackage() {
    if (document.getElementById("package-name").value.length === 0) {
        document.getElementById("package-name").style.backgroundColor = "#ffcccb"
    } else {
        document.getElementById("package-name").style.backgroundColor = ""
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
                        name: document.getElementById("package-name").value,
                        type: document.getElementById("package-type-selection").value,
                        order: document.getElementById("order-selection").value,
                        barcode: document.getElementById("barcode").value,
                        note: document.getElementById("package-note").value,
                    };
                    fetch("/save_package", {
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

function saveFaultType() {
    if (document.getElementById("fault-type-name").value.length === 0) {
        document.getElementById("fault-type-name").style.backgroundColor = "#ffcccb"
    } else {
        document.getElementById("fault-type-name").style.backgroundColor = ""
        if (document.getElementById("data-type-save-button").classList[1] === "primary") {
            document.getElementById("data-type-save-button").classList.remove("primary")
            document.getElementById("data-type-save-button").classList.add("alert")
            document.getElementById("data-type-save-button-mif").classList.remove("mif-floppy-disk")
            document.getElementById("data-type-save-button-mif").classList.add("mif-cross")

            setTimeout(function () {
                if (document.getElementById("data-type-save-button").classList[1] === "alert") {
                    let parseId = ""
                    if (Metro.getPlugin("#type-table", "table").getSelectedItems().length > 0) {
                        parseId = Metro.getPlugin("#type-table", "table").getSelectedItems()[0][0]
                    }
                    let data = {
                        id: parseId,
                        name: document.getElementById("fault-type-name").value,
                        note: document.getElementById("fault-type-note").value,
                    };
                    fetch("/save_fault_type", {
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
        } else if (document.getElementById("data-type-save-button").classList[1] === "alert") {
            document.getElementById("data-type-save-button").classList.remove("alert")
            document.getElementById("data-type-save-button").classList.add("primary")
            document.getElementById("data-type-save-button-mif").classList.remove("mif-cross")
            document.getElementById("data-type-save-button-mif").classList.add("mif-floppy-disk")
        }
    }
}


function saveFault() {
    if (document.getElementById("fault-name").value.length === 0) {
        document.getElementById("fault-name").style.backgroundColor = "#ffcccb"
    } else {
        document.getElementById("fault-name").style.backgroundColor = ""
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
                        name: document.getElementById("fault-name").value,
                        type: document.getElementById("fault-type-selection").value,
                        barcode: document.getElementById("barcode").value,
                        note: document.getElementById("fault-note").value,
                    };
                    fetch("/save_fault", {
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

function saveDowntimeType() {
    if (document.getElementById("downtime-type-name").value.length === 0) {
        document.getElementById("downtime-type-name").style.backgroundColor = "#ffcccb"
    } else {
        document.getElementById("downtime-type-name").style.backgroundColor = ""
        if (document.getElementById("data-type-save-button").classList[1] === "primary") {
            document.getElementById("data-type-save-button").classList.remove("primary")
            document.getElementById("data-type-save-button").classList.add("alert")
            document.getElementById("data-type-save-button-mif").classList.remove("mif-floppy-disk")
            document.getElementById("data-type-save-button-mif").classList.add("mif-cross")

            setTimeout(function () {
                if (document.getElementById("data-type-save-button").classList[1] === "alert") {
                    let parseId = ""
                    if (Metro.getPlugin("#type-table", "table").getSelectedItems().length > 0) {
                        parseId = Metro.getPlugin("#type-table", "table").getSelectedItems()[0][0]
                    }
                    let data = {
                        id: parseId,
                        name: document.getElementById("downtime-type-name").value,
                        note: document.getElementById("downtime-type-note").value,
                    };
                    fetch("/save_downtime_type", {
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
        } else if (document.getElementById("data-type-save-button").classList[1] === "alert") {
            document.getElementById("data-type-save-button").classList.remove("alert")
            document.getElementById("data-type-save-button").classList.add("primary")
            document.getElementById("data-type-save-button-mif").classList.remove("mif-cross")
            document.getElementById("data-type-save-button-mif").classList.add("mif-floppy-disk")
        }
    }
}


function saveDowntime() {
    if (document.getElementById("downtime-name").value.length === 0) {
        document.getElementById("downtime-name").style.backgroundColor = "#ffcccb"
    } else {
        document.getElementById("downtime-name").style.backgroundColor = ""
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
                    let background = ""
                    let colorCursor = document.getElementsByClassName("color-cursor")
                    for (const color of colorCursor) {
                        background = getComputedStyle(color).background
                    }
                    let data = {
                        id: parseId,
                        name: document.getElementById("downtime-name").value,
                        type: document.getElementById("downtime-type-selection").value,
                        barcode: document.getElementById("barcode").value,
                        color: background,
                        note: document.getElementById("downtime-note").value,
                    };
                    fetch("/save_downtime", {
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

function saveBreakdownType() {
    if (document.getElementById("breakdown-type-name").value.length === 0) {
        document.getElementById("breakdown-type-name").style.backgroundColor = "#ffcccb"
    } else {
        document.getElementById("breakdown-type-name").style.backgroundColor = ""
        if (document.getElementById("data-type-save-button").classList[1] === "primary") {
            document.getElementById("data-type-save-button").classList.remove("primary")
            document.getElementById("data-type-save-button").classList.add("alert")
            document.getElementById("data-type-save-button-mif").classList.remove("mif-floppy-disk")
            document.getElementById("data-type-save-button-mif").classList.add("mif-cross")

            setTimeout(function () {
                if (document.getElementById("data-type-save-button").classList[1] === "alert") {
                    let parseId = ""
                    if (Metro.getPlugin("#type-table", "table").getSelectedItems().length > 0) {
                        parseId = Metro.getPlugin("#type-table", "table").getSelectedItems()[0][0]
                    }
                    let data = {
                        id: parseId,
                        name: document.getElementById("breakdown-type-name").value,
                        note: document.getElementById("breakdown-type-note").value,
                    };
                    fetch("/save_breakdown_type", {
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
        } else if (document.getElementById("data-type-save-button").classList[1] === "alert") {
            document.getElementById("data-type-save-button").classList.remove("alert")
            document.getElementById("data-type-save-button").classList.add("primary")
            document.getElementById("data-type-save-button-mif").classList.remove("mif-cross")
            document.getElementById("data-type-save-button-mif").classList.add("mif-floppy-disk")
        }
    }
}

function saveBreakdown() {
    if (document.getElementById("breakdown-name").value.length === 0) {
        document.getElementById("breakdown-name").style.backgroundColor = "#ffcccb"
    } else {
        document.getElementById("breakdown-name").style.backgroundColor = ""
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
                    let background = ""
                    let colorCursor = document.getElementsByClassName("color-cursor")
                    for (const color of colorCursor) {
                        background = getComputedStyle(color).background
                    }
                    let data = {
                        id: parseId,
                        name: document.getElementById("breakdown-name").value,
                        type: document.getElementById("breakdown-type-selection").value,
                        barcode: document.getElementById("barcode").value,
                        color: background,
                        note: document.getElementById("breakdown-note").value,
                    };
                    fetch("/save_breakdown", {
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

function saveWorkshift() {
    let pattern = /^(?:2[0-3]|[01]?[0-9]):[0-5][0-9]:[0-5][0-9]$/
    if (document.getElementById("workshift-name").value.length === 0) {
        document.getElementById("workshift-name").style.backgroundColor = "#ffcccb"
    } else if (!document.getElementById("workshift-start").value.match(pattern)) {
        document.getElementById("workshift-start").style.backgroundColor = "#ffcccb"
    } else if (!document.getElementById("workshift-end").value.match(pattern)) {
        document.getElementById("workshift-end").style.backgroundColor = "#ffcccb"
    } else {
        document.getElementById("workshift-name").style.backgroundColor = ""
        document.getElementById("workshift-start").style.backgroundColor = ""
        document.getElementById("workshift-end").style.backgroundColor = ""
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
                    let background = ""
                    let colorCursor = document.getElementsByClassName("color-cursor")
                    for (const color of colorCursor) {
                        background = getComputedStyle(color).background
                    }
                    let data = {
                        id: parseId,
                        name: document.getElementById("workshift-name").value,
                        start: document.getElementById("workshift-start").value,
                        end: document.getElementById("workshift-end").value,
                        note: document.getElementById("workshift-note").value,
                    };
                    fetch("/save_workshift", {
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

function saveState() {
    if (document.getElementById("state-name").value.length === 0) {
        document.getElementById("state-name").style.backgroundColor = "#ffcccb"
    } else {
        document.getElementById("state-name").style.backgroundColor = ""
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
                    let background = ""
                    let colorCursor = document.getElementsByClassName("color-cursor")
                    for (const color of colorCursor) {
                        background = getComputedStyle(color).background
                    }
                    let data = {
                        id: parseId,
                        name: document.getElementById("state-name").value,
                        color: background,
                        note: document.getElementById("state-note").value,
                    };
                    fetch("/save_state", {
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

function savePart() {
    if (document.getElementById("part-name").value.length === 0) {
        document.getElementById("part-name").style.backgroundColor = "#ffcccb"
    } else {
        document.getElementById("part-name").style.backgroundColor = ""
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
                        name: document.getElementById("part-name").value,
                        barcode: document.getElementById("part-barcode").value,
                        note: document.getElementById("part-note").value,
                    };
                    fetch("/save_part", {
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

function saveProduct() {
    if (document.getElementById("product-name").value.length === 0) {
        document.getElementById("product-name").style.backgroundColor = "#ffcccb"
    } else {
        document.getElementById("product-name").style.backgroundColor = ""
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
                        name: document.getElementById("product-name").value,
                        barcode: document.getElementById("product-barcode").value,
                        cycle: document.getElementById("cycle-time").value,
                        downtimeDuration: document.getElementById("downtime-duration").value,
                        note: document.getElementById("product-note").value,
                    };
                    fetch("/save_product", {
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

function loadDetails(selectedItem, type) {
    let selection = document.getElementById("data-selection").value
    let data = {
        data: selection,
        name: selectedItem,
        type: type
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

