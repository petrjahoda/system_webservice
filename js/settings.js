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
    let typeTableExtended = Metro.getPlugin("#type-table-extended", "table");
    if (tableSelectedId === "data-table" && table.getSelectedItems().length > 0 && event.target.type === "radio") {
        let selectedItem = table.getSelectedItems()[0][0]
        loadDetails(selectedItem, "first");
    } else if (tableSelectedId === "type-table" && typeTable.getSelectedItems().length > 0 && event.target.type === "radio") {
        let selectedItem = typeTable.getSelectedItems()[0][0]
        loadDetails(selectedItem, "second");
    } else if (tableSelectedId === "type-table-extended" && typeTableExtended.getSelectedItems().length > 0 && event.target.type === "radio") {
        let selectedItem = typeTableExtended.getSelectedItems()[0][0]
        loadDetails(selectedItem, "third");
    } else if (event.target.id === "data-new-button" || event.target.id === "data-new-button-mif") {
        loadSettings();
        loadDetails(null, "first");
    } else if (event.target.id === "data-new-button-type" || event.target.id === "data-new-button-mif-type") {
        loadSettings();
        loadDetails(null, "second");
    } else if (event.target.id === "data-new-button-type-extended" || event.target.id === "data-new-button-mif-type-extended") {
        loadSettings();
        loadDetails(null, "third");
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
            case "devices" : {
                saveDevice();
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
            case "workplaces": {
                saveWorkplace();
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
            case "workplaces" : {
                saveWorkplaceSection()
                break;
            }
        }
    } else if (event.target.id === "data-mode-save-button" || event.target.id === "data-mode-save-button-mif") {
        let selection = document.getElementById("data-selection").value
        console.log("Saving mode of " + selection)
        switch (selection) {
            case "workplaces" : {
                saveWorkplaceMode();
                break
            }
        }
    } else if (event.target.id === "data-new-port-button" || event.target.id === "data-new-port-button-mif") {
        loadDevicePortDetails(null, "device");
    } else if (event.target.id === "data-new-workplace-port-button" || event.target.id === "data-new-workplace-port-button-mif") {
        loadWorkplacePortDetails(null, "workplace");
    } else if (event.target.id === "data-delete-workplace-port-button" || event.target.id === "data-delete-workplace-port-button-mif") {
        deleteWorkplacePortFromWorkplace(null, "workplace");
    } else if (event.target.id === "data-new-workshift-button" || event.target.id === "data-new-workshift-button-mif") {
        loadWorkshiftDetails(null, "workshift");
    } else if (event.target.id === "data-delete-workshift-button" || event.target.id === "data-delete-workshift-button-mif") {
        deleteWorkshiftFromWorkplace(null, "workshift");
    } else if (event.target.id === "port-save-button" || event.target.id === "port-save-button-mif") {
        saveDevicePortDetails();
    } else if (event.target.id === "workplace-port-save-button" || event.target.id === "workplace-port-save-button-mif") {
        saveWorkplacePortDetails();
    } else if (event.target.id === "workshift-save-button" || event.target.id === "workshift-save-button-mif") {
        saveWorkshiftDetails();
    } else {
        const tablePortSelectedId = event.target.parentElement.parentElement.parentElement.parentElement.parentElement.id
        let portTable = Metro.getPlugin("#data-port-table", "table");
        let workplacePortTable = Metro.getPlugin("#data-workplace-port-table", "table");
        let workshiftTable = Metro.getPlugin("#data-workshift-table", "table");
        if (tablePortSelectedId === "data-port-table" && portTable.getSelectedItems().length > 0 && event.target.type === "radio") {
            let selectedItem = portTable.getSelectedItems()[0][0]
            loadDevicePortDetails(selectedItem, "device");
        } else if (tablePortSelectedId === "data-workplace-port-table" && workplacePortTable.getSelectedItems().length > 0 && event.target.type === "radio") {
            let selectedItem = workplacePortTable.getSelectedItems()[0][0]
            loadWorkplacePortDetails(selectedItem, "workplace");
        } else if (tablePortSelectedId === "data-workshift-table" && workshiftTable.getSelectedItems().length > 0 && event.target.type === "radio") {
            let selectedItem = workshiftTable.getSelectedItems()[0][0]
            loadWorkshiftDetails(selectedItem, "workplace");
        }

    }
})


function saveWorkplace() {
    if (document.getElementById("workplace-name").value.length === 0) {
        document.getElementById("workplace-name").style.backgroundColor = "#ffcccb"
    } else {
        document.getElementById("workplace-name").style.backgroundColor = ""
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
                        name: document.getElementById("workplace-name").value,
                        section: document.getElementById("workplace-section-selection").value,
                        mode: document.getElementById("workplace-mode-selection").value,
                        code: document.getElementById("code").value,
                        note: document.getElementById("workplace-note").value,
                    };
                    fetch("/save_workplace", {
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


function saveWorkplaceSection() {
    if (document.getElementById("workplace-section-name").value.length === 0) {
        document.getElementById("workplace-section-name").style.backgroundColor = "#ffcccb"
    } else {
        document.getElementById("workplace-section-name").style.backgroundColor = ""
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
                        name: document.getElementById("workplace-section-name").value,
                        note: document.getElementById("workplace-section-note").value,
                    };
                    fetch("/save_workplace_section", {
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

function saveWorkplaceMode() {
    if (document.getElementById("workplace-mode-name").value.length === 0) {
        document.getElementById("workplace-mode-name").style.backgroundColor = "#ffcccb"
    } else {
        document.getElementById("workplace-mode-name").style.backgroundColor = ""
        if (document.getElementById("data-mode-save-button").classList[1] === "primary") {
            document.getElementById("data-mode-save-button").classList.remove("primary")
            document.getElementById("data-mode-save-button").classList.add("alert")
            document.getElementById("data-mode-save-button-mif").classList.remove("mif-floppy-disk")
            document.getElementById("data-mode-save-button-mif").classList.add("mif-cross")

            setTimeout(function () {
                if (document.getElementById("data-mode-save-button").classList[1] === "alert") {
                    let parseId = ""
                    if (Metro.getPlugin("#type-table-extended", "table").getSelectedItems().length > 0) {
                        parseId = Metro.getPlugin("#type-table-extended", "table").getSelectedItems()[0][0]
                    }
                    let data = {
                        id: parseId,
                        name: document.getElementById("workplace-mode-name").value,
                        downtimeDuration: document.getElementById("downtime-duration").value,
                        poweroffDuration: document.getElementById("poweroff-duration").value,
                        note: document.getElementById("workplace-mode-note").value,
                    };
                    fetch("/save_workplace_mode", {
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
        } else if (document.getElementById("data-mode-save-button").classList[1] === "alert") {
            document.getElementById("data-mode-save-button").classList.remove("alert")
            document.getElementById("data-mode-save-button").classList.add("primary")
            document.getElementById("data-mode-save-button-mif").classList.remove("mif-cross")
            document.getElementById("data-mode-save-button-mif").classList.add("mif-floppy-disk")
        }
    }
}

function saveWorkplacePortDetails() {
    if (document.getElementById("workplace-port-name").value.length === 0) {
        document.getElementById("workplace-port-name").style.backgroundColor = "#ffcccb"
    } else {
        document.getElementById("workplace-port-name").style.backgroundColor = ""
        if (document.getElementById("workplace-port-save-button").classList[1] === "primary") {
            document.getElementById("workplace-port-save-button").classList.remove("primary")
            document.getElementById("workplace-port-save-button").classList.add("alert")
            document.getElementById("workplace-port-save-button-mif").classList.remove("mif-floppy-disk")
            document.getElementById("workplace-port-save-button-mif").classList.add("mif-cross")

            setTimeout(function () {
                if (document.getElementById("workplace-port-save-button").classList[1] === "alert") {
                    let parseId = ""
                    if (Metro.getPlugin("#data-workplace-port-table", "table").getSelectedItems().length > 0) {
                        parseId = Metro.getPlugin("#data-workplace-port-table", "table").getSelectedItems()[0][0]
                    }
                    let background = ""
                    let colorCursor = document.getElementsByClassName("color-cursor")
                    for (const color of colorCursor) {
                        background = getComputedStyle(color).background
                    }
                    let data = {
                        id: parseId,
                        workplaceName: document.getElementById("workplace-name").value,
                        name: document.getElementById("workplace-port-name").value,
                        devicePortId: document.getElementById("workplace-port-device-port-selection").value,
                        stateId: document.getElementById("workplace-port-state-selection").value,
                        color: background,
                        counterOk: document.getElementById("workplace-port-counter-ok-selection").value,
                        counterNok: document.getElementById("workplace-port-counter-nok-selection").value,
                        highValue: document.getElementById("workplace-port-high-value").value,
                        lowValue: document.getElementById("workplace-port-low-value").value,
                        note: document.getElementById("workplace-port-note").value,
                    };
                    fetch("/save_workplace_port_details", {
                        method: "POST",
                        body: JSON.stringify(data)
                    }).then((response) => {
                        response.text().then(function (data) {
                            document.getElementById("workplace-port-container").innerHTML = ""
                            let table = Metro.getPlugin("#data-table", "table");
                            let selectedItem = table.getSelectedItems()[0][0]
                            loadDetails(selectedItem, "first");
                        });
                    }).catch((error) => {
                        console.log(error)
                        document.getElementById("workplace-port-container").innerHTML = ""
                        let table = Metro.getPlugin("#data-table", "table");
                        let selectedItem = table.getSelectedItems()[0][0]
                        loadDetails(selectedItem, "first");
                    });
                }
            }, 2500);
        } else if (document.getElementById("workplace-port-save-button").classList[1] === "alert") {
            document.getElementById("workplace-port-save-button").classList.remove("alert")
            document.getElementById("workplace-port-save-button").classList.add("primary")
            document.getElementById("workplace-port-save-button-mif").classList.remove("mif-cross")
            document.getElementById("workplace-port-save-button-mif").classList.add("mif-floppy-disk")
        }
    }
}

function saveDevicePortDetails() {
    if (document.getElementById("device-port-name").value.length === 0) {
        document.getElementById("device-port-name").style.backgroundColor = "#ffcccb"
    } else {
        document.getElementById("device-port-name").style.backgroundColor = ""
        if (document.getElementById("port-save-button").classList[1] === "primary") {
            document.getElementById("port-save-button").classList.remove("primary")
            document.getElementById("port-save-button").classList.add("alert")
            document.getElementById("port-save-button-mif").classList.remove("mif-floppy-disk")
            document.getElementById("port-save-button-mif").classList.add("mif-cross")

            setTimeout(function () {
                if (document.getElementById("port-save-button").classList[1] === "alert") {
                    let parseId = ""
                    if (Metro.getPlugin("#data-port-table", "table").getSelectedItems().length > 0) {
                        parseId = Metro.getPlugin("#data-port-table", "table").getSelectedItems()[0][0]
                    }
                    let data = {
                        id: parseId,
                        deviceName: document.getElementById("device-name").value,
                        name: document.getElementById("device-port-name").value,
                        type: document.getElementById("device-port-type-selection").value,
                        position: document.getElementById("device-port-file-position").value,
                        unit: document.getElementById("device-port-unit").value,
                        plcDataType: document.getElementById("device-port-plc-data-type").value,
                        plcDataAddress: document.getElementById("device-port-plc-data-address").value,
                        settings: document.getElementById("device-port-settings").value,
                        note: document.getElementById("device-port-note").value,
                        virtual: document.getElementById("device-port-virtual-selection").value,
                    };
                    fetch("/save_device_port_details", {
                        method: "POST",
                        body: JSON.stringify(data)
                    }).then((response) => {
                        response.text().then(function (data) {
                            document.getElementById("port-container").innerHTML = ""
                            let table = Metro.getPlugin("#data-table", "table");
                            let selectedItem = table.getSelectedItems()[0][0]
                            loadDetails(selectedItem, "first");
                        });
                    }).catch((error) => {
                        console.log(error)
                        document.getElementById("port-container").innerHTML = ""
                        let table = Metro.getPlugin("#data-table", "table");
                        let selectedItem = table.getSelectedItems()[0][0]
                        loadDetails(selectedItem, "first");
                    });
                }
            }, 2500);
        } else if (document.getElementById("port-save-button").classList[1] === "alert") {
            document.getElementById("port-save-button").classList.remove("alert")
            document.getElementById("port-save-button").classList.add("primary")
            document.getElementById("port-save-button-mif").classList.remove("mif-cross")
            document.getElementById("port-save-button-mif").classList.add("mif-floppy-disk")
        }
    }
}

function loadWorkplacePortDetails(selectedPort, type) {
    console.log(type + ": loading workplace port details for " + selectedPort)
    let data = {
        data: selectedPort,
        type: type
    };
    fetch("/load_workplace_port_detail", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            document.getElementById("workplace-port-container").innerHTML = data
            setTimeout(function () {
                document.getElementById("data-delete-workplace-port-button").hidden = false
                document.getElementById("workplace-port-container").scrollIntoView();
            }, 100);

        });
    }).catch((error) => {
        console.log(error)
    });
}

function loadDevicePortDetails(selectedPort, type) {
    console.log(type + ": loading device port details for " + selectedPort)
    let data = {
        data: selectedPort,
        type: type
    };
    fetch("/load_device_port_detail", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            document.getElementById("port-container").innerHTML = data
            setTimeout(function () {
                document.getElementById("port-container").scrollIntoView();
            }, 100);

        });
    }).catch((error) => {
        console.log(error)
    });
}

function saveDevice() {
    if (document.getElementById("device-name").value.length === 0) {
        document.getElementById("device-name").style.backgroundColor = "#ffcccb"
    } else {
        document.getElementById("device-name").style.backgroundColor = ""
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
                        name: document.getElementById("device-name").value,
                        type: document.getElementById("device-type-selection").value,
                        ip: document.getElementById("ip-address").value,
                        mac: document.getElementById("mac-address").value,
                        version: document.getElementById("device-version-name").value,
                        settings: document.getElementById("device-settings").value,
                        note: document.getElementById("device-note").value,
                        enabled: document.getElementById("device-enabled-selection").value
                    };
                    fetch("/save_device", {
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
    console.log("loading")
    let selection = document.getElementById("data-selection").value
    let data = {
        data: selection,
        name: selectedItem,
        type: type
    };
    console.log(data)
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

