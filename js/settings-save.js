function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

function changeDataSaveButtonToWarning() {
    document.getElementById("data-save-button").classList.remove("primary")
    document.getElementById("data-save-button").classList.add("warning")
    document.getElementById("data-save-button-mif").classList.remove("mif-floppy-disk")
    document.getElementById("data-save-button-mif").classList.add("mif-cross")
}

function changeDataTypeSaveButtonToWarning() {
    document.getElementById("data-type-save-button").classList.remove("primary")
    document.getElementById("data-type-save-button").classList.add("warning")
    document.getElementById("data-type-save-button-mif").classList.remove("mif-floppy-disk")
    document.getElementById("data-type-save-button-mif").classList.add("mif-cross")
}

function changeWorkplacePortSaveButtonToWarning() {
    document.getElementById("workplace-port-save-button").classList.remove("primary")
    document.getElementById("workplace-port-save-button").classList.add("warning")
    document.getElementById("workplace-port-save-button-mif").classList.remove("mif-floppy-disk")
    document.getElementById("workplace-port-save-button-mif").classList.add("mif-cross")
}

function changeDevicePortSaveButtonToWarning() {
    document.getElementById("port-save-button").classList.remove("primary")
    document.getElementById("port-save-button").classList.add("warning")
    document.getElementById("port-save-button-mif").classList.remove("mif-floppy-disk")
    document.getElementById("port-save-button-mif").classList.add("mif-cross")
}

function changeDataSaveButtonToPrimary() {
    document.getElementById("data-save-button").classList.remove("warning")
    document.getElementById("data-save-button").classList.add("primary")
    document.getElementById("data-save-button-mif").classList.remove("mif-cross")
    document.getElementById("data-save-button-mif").classList.add("mif-floppy-disk")
}

function changeDataTypeSaveButtonToPrimary() {
    document.getElementById("data-type-save-button").classList.remove("warning")
    document.getElementById("data-type-save-button").classList.add("primary")
    document.getElementById("data-type-save-button-mif").classList.remove("mif-cross")
    document.getElementById("data-type-save-button-mif").classList.add("mif-floppy-disk")
}

function changeDataModeButtonToWarning() {
    document.getElementById("data-mode-save-button").classList.remove("primary")
    document.getElementById("data-mode-save-button").classList.add("warning")
    document.getElementById("data-mode-save-button-mif").classList.remove("mif-floppy-disk")
    document.getElementById("data-mode-save-button-mif").classList.add("mif-cross")
}

function changeDataModeButtonToPrimary() {
    document.getElementById("data-mode-save-button").classList.remove("warning")
    document.getElementById("data-mode-save-button").classList.add("primary")
    document.getElementById("data-mode-save-button-mif").classList.remove("mif-cross")
    document.getElementById("data-mode-save-button-mif").classList.add("mif-floppy-disk")
}

function changeWorkplacePortSaveButtonToPrimary() {
    document.getElementById("workplace-port-save-button").classList.remove("warning")
    document.getElementById("workplace-port-save-button").classList.add("primary")
    document.getElementById("workplace-port-save-button-mif").classList.remove("mif-cross")
    document.getElementById("workplace-port-save-button-mif").classList.add("mif-floppy-disk")
}

function changeDevicePortSaveButtonToPrimary() {
    document.getElementById("port-save-button").classList.remove("warning")
    document.getElementById("port-save-button").classList.add("primary")
    document.getElementById("port-save-button-mif").classList.remove("mif-cross")
    document.getElementById("port-save-button-mif").classList.add("mif-floppy-disk")
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
            changeDataSaveButtonToWarning();
            setTimeout(function () {
                if (document.getElementById("data-save-button").classList[1] === "warning") {
                    let data = {
                        id: sessionStorage.getItem("selected_id"),
                        name: document.getElementById("workshift-name").value,
                        start: document.getElementById("workshift-start").value,
                        end: document.getElementById("workshift-end").value,
                        note: document.getElementById("workshift-note").value,
                    };
                    fetch("/save_workshift", {
                        method: "POST",
                        body: JSON.stringify(data)
                    }).then((response) => {
                        response.text().then(function () {
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
        } else if (document.getElementById("data-save-button").classList[1] === "warning") {
            changeDataSaveButtonToPrimary();
        }
    }
}

function saveAlarm() {
    if (document.getElementById("alarm-name").value.length === 0) {
        document.getElementById("alarm-name").style.backgroundColor = "#ffcccb"
    } else {
        document.getElementById("alarm-name").style.backgroundColor = ""
        if (document.getElementById("data-save-button").classList[1] === "primary") {
            changeDataSaveButtonToWarning();
            setTimeout(function () {
                if (document.getElementById("data-save-button").classList[1] === "warning") {
                    let data = {
                        id: sessionStorage.getItem("selected_id"),
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
                        response.text().then(function () {
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
        } else if (document.getElementById("data-save-button").classList[1] === "warning") {
            changeDataSaveButtonToPrimary();
        }
    }
}

function saveDevice() {
    if (document.getElementById("device-name").value.length === 0) {
        document.getElementById("device-name").style.backgroundColor = "#ffcccb"
    } else {
        document.getElementById("device-name").style.backgroundColor = ""
        if (document.getElementById("data-save-button").classList[1] === "primary") {
            changeDataSaveButtonToWarning();
            setTimeout(function () {
                if (document.getElementById("data-save-button").classList[1] === "warning") {
                    let data = {
                        id: sessionStorage.getItem("selected_id"),
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
                        response.text().then(function () {
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
        } else if (document.getElementById("data-save-button").classList[1] === "warning") {
            changeDataSaveButtonToPrimary();
        }
    }
}

function saveSystemSettings() {
    if (document.getElementById("system-settings-name").value.length === 0) {
        document.getElementById("system-settings-name").style.backgroundColor = "#ffcccb"
    } else {
        document.getElementById("system-settings-name").style.backgroundColor = ""
        if (document.getElementById("data-save-button").classList[1] === "primary") {
            changeDataSaveButtonToWarning();
            setTimeout(function () {
                if (document.getElementById("data-save-button").classList[1] === "warning") {
                    let data = {
                        id: sessionStorage.getItem("selected_id"),
                        name: document.getElementById("system-settings-name").value,
                        value: document.getElementById("system-settings-value").value,
                        enabled: document.getElementById("system-settings-selection").value,
                        note: document.getElementById("system-settings-note").value,
                    };
                    fetch("/save_system_settings", {
                        method: "POST",
                        body: JSON.stringify(data)
                    }).then((response) => {
                        response.text().then(function () {
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
        } else if (document.getElementById("data-save-button").classList[1] === "warning") {
            changeDataSaveButtonToPrimary();
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
            changeDataSaveButtonToWarning();
            setTimeout(function () {
                if (document.getElementById("data-save-button").classList[1] === "warning") {
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
                        response.text().then(function () {
                            document.getElementById("settings-container-detail").innerHTML = ""
                            location.reload()
                        });
                    }).catch((error) => {
                        console.log(error)
                        document.getElementById("settings-container-detail").innerHTML = ""
                        loadSettings();
                    });
                }
            }, 2500);
        } else if (document.getElementById("data-save-button").classList[1] === "warning") {
            changeDataSaveButtonToPrimary();
        }
    }
}

function saveUserType() {
    if (document.getElementById("user-type-name").value.length === 0) {
        document.getElementById("user-type-name").style.backgroundColor = "#ffcccb"
    } else {
        document.getElementById("user-type-name").style.backgroundColor = ""
        if (document.getElementById("data-type-save-button").classList[1] === "primary") {
            changeDataTypeSaveButtonToWarning();
            setTimeout(function () {
                if (document.getElementById("data-type-save-button").classList[1] === "warning") {
                    let data = {
                        id: sessionStorage.getItem("selected_id"),
                        name: document.getElementById("user-type-name").value,
                        note: document.getElementById("user-type-note").value,
                    };
                    fetch("/save_user_type", {
                        method: "POST",
                        body: JSON.stringify(data)
                    }).then((response) => {
                        response.text().then(function () {
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
        } else if (document.getElementById("data-type-save-button").classList[1] === "warning") {
            changeDataTypeSaveButtonToPrimary();
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
            changeDataSaveButtonToWarning();
            setTimeout(function () {
                if (document.getElementById("data-save-button").classList[1] === "warning") {
                    let data = {
                        id: sessionStorage.getItem("selected_id"),
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
                        response.text().then(function () {
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
        } else if (document.getElementById("data-save-button").classList[1] === "warning") {
            changeDataSaveButtonToPrimary()
        }
    }
}

function savePackageType() {
    if (document.getElementById("package-type-name").value.length === 0) {
        document.getElementById("package-type-name").style.backgroundColor = "#ffcccb"
    } else {
        document.getElementById("package-type-name").style.backgroundColor = ""
        if (document.getElementById("data-type-save-button").classList[1] === "primary") {
            changeDataTypeSaveButtonToWarning();
            setTimeout(function () {
                if (document.getElementById("data-type-save-button").classList[1] === "warning") {
                    let data = {
                        id: sessionStorage.getItem("selected_id"),
                        name: document.getElementById("package-type-name").value,
                        count: document.getElementById("package-type-count").value,
                        note: document.getElementById("package-type-note").value,
                    };
                    fetch("/save_package_type", {
                        method: "POST",
                        body: JSON.stringify(data)
                    }).then((response) => {
                        response.text().then(function () {
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
        } else if (document.getElementById("data-type-save-button").classList[1] === "warning") {
            changeDataTypeSaveButtonToPrimary();
        }
    }
}

function savePackage() {
    if (document.getElementById("package-name").value.length === 0) {
        document.getElementById("package-name").style.backgroundColor = "#ffcccb"
    } else {
        document.getElementById("package-name").style.backgroundColor = ""
        if (document.getElementById("data-save-button").classList[1] === "primary") {
            changeDataSaveButtonToWarning();
            setTimeout(function () {
                if (document.getElementById("data-save-button").classList[1] === "warning") {
                    let data = {
                        id: sessionStorage.getItem("selected_id"),
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
                        response.text().then(function () {
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
        } else if (document.getElementById("data-save-button").classList[1] === "warning") {
            changeDataSaveButtonToPrimary();
        }
    }
}

function saveFaultType() {
    if (document.getElementById("fault-type-name").value.length === 0) {
        document.getElementById("fault-type-name").style.backgroundColor = "#ffcccb"
    } else {
        document.getElementById("fault-type-name").style.backgroundColor = ""
        if (document.getElementById("data-type-save-button").classList[1] === "primary") {
            changeDataTypeSaveButtonToWarning();
            setTimeout(function () {
                if (document.getElementById("data-type-save-button").classList[1] === "warning") {
                    let data = {
                        id: sessionStorage.getItem("selected_id"),
                        name: document.getElementById("fault-type-name").value,
                        note: document.getElementById("fault-type-note").value,
                    };
                    fetch("/save_fault_type", {
                        method: "POST",
                        body: JSON.stringify(data)
                    }).then((response) => {
                        response.text().then(function () {
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
        } else if (document.getElementById("data-type-save-button").classList[1] === "warning") {
            changeDataTypeSaveButtonToPrimary();
        }
    }
}

function saveFault() {
    if (document.getElementById("fault-name").value.length === 0) {
        document.getElementById("fault-name").style.backgroundColor = "#ffcccb"
    } else {
        document.getElementById("fault-name").style.backgroundColor = ""
        if (document.getElementById("data-save-button").classList[1] === "primary") {
            changeDataSaveButtonToWarning();
            setTimeout(function () {
                if (document.getElementById("data-save-button").classList[1] === "warning") {
                    let data = {
                        id: sessionStorage.getItem("selected_id"),
                        name: document.getElementById("fault-name").value,
                        type: document.getElementById("fault-type-selection").value,
                        barcode: document.getElementById("barcode").value,
                        note: document.getElementById("fault-note").value,
                    };
                    fetch("/save_fault", {
                        method: "POST",
                        body: JSON.stringify(data)
                    }).then((response) => {
                        response.text().then(function () {
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
        } else if (document.getElementById("data-save-button").classList[1] === "warning") {
            changeDataSaveButtonToPrimary();
        }
    }
}

function saveDowntimeType() {
    if (document.getElementById("downtime-type-name").value.length === 0) {
        document.getElementById("downtime-type-name").style.backgroundColor = "#ffcccb"
    } else {
        document.getElementById("downtime-type-name").style.backgroundColor = ""
        if (document.getElementById("data-type-save-button").classList[1] === "primary") {
            changeDataTypeSaveButtonToWarning();
            setTimeout(function () {
                if (document.getElementById("data-type-save-button").classList[1] === "warning") {
                    let data = {
                        id: sessionStorage.getItem("selected_id"),
                        name: document.getElementById("downtime-type-name").value,
                        note: document.getElementById("downtime-type-note").value,
                    };
                    fetch("/save_downtime_type", {
                        method: "POST",
                        body: JSON.stringify(data)
                    }).then((response) => {
                        response.text().then(function () {
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
        } else if (document.getElementById("data-type-save-button").classList[1] === "warning") {
            changeDataTypeSaveButtonToPrimary();
        }
    }
}

function saveDowntime() {
    if (document.getElementById("downtime-name").value.length === 0) {
        document.getElementById("downtime-name").style.backgroundColor = "#ffcccb"
    } else {
        document.getElementById("downtime-name").style.backgroundColor = ""
        if (document.getElementById("data-save-button").classList[1] === "primary") {
            changeDataSaveButtonToWarning();
            setTimeout(function () {
                if (document.getElementById("data-save-button").classList[1] === "warning") {
                    let data = {
                        id: sessionStorage.getItem("selected_id"),
                        name: document.getElementById("downtime-name").value,
                        type: document.getElementById("downtime-type-selection").value,
                        barcode: document.getElementById("barcode").value,
                        color: document.getElementById("downtime-color").value,
                        note: document.getElementById("downtime-note").value,
                    };
                    fetch("/save_downtime", {
                        method: "POST",
                        body: JSON.stringify(data)
                    }).then((response) => {
                        response.text().then(function () {
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
        } else if (document.getElementById("data-save-button").classList[1] === "warning") {
            changeDataSaveButtonToPrimary()
        }
    }
}

function saveBreakdownType() {
    if (document.getElementById("breakdown-type-name").value.length === 0) {
        document.getElementById("breakdown-type-name").style.backgroundColor = "#ffcccb"
    } else {
        document.getElementById("breakdown-type-name").style.backgroundColor = ""
        if (document.getElementById("data-type-save-button").classList[1] === "primary") {
            changeDataTypeSaveButtonToWarning();
            setTimeout(function () {
                if (document.getElementById("data-type-save-button").classList[1] === "warning") {
                    let data = {
                        id: sessionStorage.getItem("selected_id"),
                        name: document.getElementById("breakdown-type-name").value,
                        note: document.getElementById("breakdown-type-note").value,
                    };
                    fetch("/save_breakdown_type", {
                        method: "POST",
                        body: JSON.stringify(data)
                    }).then((response) => {
                        response.text().then(function () {
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
        } else if (document.getElementById("data-type-save-button").classList[1] === "warning") {
            changeDataTypeSaveButtonToPrimary();
        }
    }
}

function saveBreakdown() {
    if (document.getElementById("breakdown-name").value.length === 0) {
        document.getElementById("breakdown-name").style.backgroundColor = "#ffcccb"
    } else {
        document.getElementById("breakdown-name").style.backgroundColor = ""
        if (document.getElementById("data-save-button").classList[1] === "primary") {
            changeDataSaveButtonToWarning();
            setTimeout(function () {
                if (document.getElementById("data-save-button").classList[1] === "warning") {
                    let data = {
                        id: sessionStorage.getItem("selected_id"),
                        name: document.getElementById("breakdown-name").value,
                        type: document.getElementById("breakdown-type-selection").value,
                        barcode: document.getElementById("barcode").value,
                        color: document.getElementById("breakdown-color").value,
                        note: document.getElementById("breakdown-note").value,
                    };
                    fetch("/save_breakdown", {
                        method: "POST",
                        body: JSON.stringify(data)
                    }).then((response) => {
                        response.text().then(function () {
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
        } else if (document.getElementById("data-save-button").classList[1] === "warning") {
            changeDataSaveButtonToPrimary()
        }
    }
}

function saveState() {
    if (document.getElementById("state-name").value.length === 0) {
        document.getElementById("state-name").style.backgroundColor = "#ffcccb"
    } else {
        document.getElementById("state-name").style.backgroundColor = ""
        if (document.getElementById("data-save-button").classList[1] === "primary") {
            changeDataSaveButtonToWarning();
            setTimeout(function () {
                if (document.getElementById("data-save-button").classList[1] === "warning") {
                    let data = {
                        id: sessionStorage.getItem("selected_id"),
                        name: document.getElementById("state-name").value,
                        color: document.getElementById("state-color").value,
                        note: document.getElementById("state-note").value,
                    };
                    fetch("/save_state", {
                        method: "POST",
                        body: JSON.stringify(data)
                    }).then((response) => {
                        response.text().then(function () {
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
        } else if (document.getElementById("data-save-button").classList[1] === "warning") {
            changeDataSaveButtonToPrimary()
        }
    }
}

function savePart() {
    if (document.getElementById("part-name").value.length === 0) {
        document.getElementById("part-name").style.backgroundColor = "#ffcccb"
    } else {
        document.getElementById("part-name").style.backgroundColor = ""
        if (document.getElementById("data-save-button").classList[1] === "primary") {
            changeDataSaveButtonToWarning();
            setTimeout(function () {
                if (document.getElementById("data-save-button").classList[1] === "warning") {
                    let data = {
                        id: sessionStorage.getItem("selected_id"),
                        name: document.getElementById("part-name").value,
                        barcode: document.getElementById("part-barcode").value,
                        note: document.getElementById("part-note").value,
                    };
                    fetch("/save_part", {
                        method: "POST",
                        body: JSON.stringify(data)
                    }).then((response) => {
                        response.text().then(function () {
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
        } else if (document.getElementById("data-save-button").classList[1] === "warning") {
            changeDataSaveButtonToPrimary()
        }
    }
}

function saveProduct() {
    if (document.getElementById("product-name").value.length === 0) {
        document.getElementById("product-name").style.backgroundColor = "#ffcccb"
    } else {
        document.getElementById("product-name").style.backgroundColor = ""
        if (document.getElementById("data-save-button").classList[1] === "primary") {
            changeDataSaveButtonToWarning();
            setTimeout(function () {
                if (document.getElementById("data-save-button").classList[1] === "warning") {
                    let data = {
                        id: sessionStorage.getItem("selected_id"),
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
                        response.text().then(function () {
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
        } else if (document.getElementById("data-save-button").classList[1] === "warning") {
            changeDataSaveButtonToPrimary()
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
            changeDataSaveButtonToWarning();
            setTimeout(function () {
                if (document.getElementById("data-save-button").classList[1] === "warning") {
                    let data = {
                        id: sessionStorage.getItem("selected_id"),
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
                        response.text().then(function () {
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
        } else if (document.getElementById("data-save-button").classList[1] === "warning") {
            changeDataSaveButtonToPrimary()
        }
    }
}

function saveOperation() {
    if (document.getElementById("operation-name").value.length === 0) {
        document.getElementById("operation-name").style.backgroundColor = "#ffcccb"
    } else {
        document.getElementById("operation-name").style.backgroundColor = ""
        if (document.getElementById("data-save-button").classList[1] === "primary") {
            changeDataSaveButtonToWarning();
            setTimeout(function () {
                if (document.getElementById("data-save-button").classList[1] === "warning") {
                    let data = {
                        id: sessionStorage.getItem("selected_id"),
                        name: document.getElementById("operation-name").value,
                        order: document.getElementById("order-selection").value,
                        barcode: document.getElementById("operation-barcode").value,
                        note: document.getElementById("operation-note").value,
                    };
                    fetch("/save_operation", {
                        method: "POST",
                        body: JSON.stringify(data)
                    }).then((response) => {
                        response.text().then(function () {
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
        } else if (document.getElementById("data-save-button").classList[1] === "warning") {
            changeDataSaveButtonToPrimary()
        }
    }
}

function saveWorkplacePortDetails() {
    if (document.getElementById("workplace-port-name").value.length === 0) {
        document.getElementById("workplace-port-name").style.backgroundColor = "#ffcccb"
    } else {
        document.getElementById("workplace-port-name").style.backgroundColor = ""
        if (document.getElementById("workplace-port-save-button").classList[1] === "primary") {
            changeWorkplacePortSaveButtonToWarning();
            setTimeout(function () {
                if (document.getElementById("workplace-port-save-button").classList[1] === "warning") {
                    let data = {
                        id: sessionStorage.getItem("workplace_port_selected_id"),
                        workplaceName: document.getElementById("workplace-name").value,
                        name: document.getElementById("workplace-port-name").value,
                        devicePortId: document.getElementById("workplace-port-device-port-selection").value,
                        stateId: document.getElementById("workplace-port-state-selection").value,
                        color: document.getElementById("port-color").value,
                        note: document.getElementById("workplace-port-note").value,
                    };
                    fetch("/save_workplace_port_details", {
                        method: "POST",
                        body: JSON.stringify(data)
                    }).then((response) => {
                        response.text().then(function () {
                            document.getElementById("workplace-port-container").innerHTML = ""
                            loadDetails(sessionStorage.getItem("selected_id"), "first");
                        });
                    }).catch((error) => {
                        console.log(error)
                        document.getElementById("workplace-port-container").innerHTML = ""
                        loadDetails(sessionStorage.getItem("selected_id"), "first");
                    });
                }
            }, 2500);
        } else if (document.getElementById("workplace-port-save-button").classList[1] === "warning") {
            changeWorkplacePortSaveButtonToPrimary();
        }
    }
}

function saveDevicePortDetails() {
    if (document.getElementById("device-port-name").value.length === 0) {
        document.getElementById("device-port-name").style.backgroundColor = "#ffcccb"
    } else {
        document.getElementById("device-port-name").style.backgroundColor = ""
        if (document.getElementById("port-save-button").classList[1] === "primary") {
            changeDevicePortSaveButtonToWarning();
            setTimeout(function () {
                if (document.getElementById("port-save-button").classList[1] === "warning") {
                    let data = {
                        id: sessionStorage.getItem("device_port_selected_id"),
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
                        response.text().then(function () {
                            document.getElementById("port-container").innerHTML = ""
                            loadDetails(sessionStorage.getItem("selected_id"), "first");
                        });
                    }).catch((error) => {
                        console.log(error)
                        document.getElementById("port-container").innerHTML = ""
                        loadDetails(sessionStorage.getItem("selected_id"), "first");
                    });
                }
            }, 2500);
        } else if (document.getElementById("port-save-button").classList[1] === "warning") {
            changeDevicePortSaveButtonToPrimary();
        }
    }
}

function saveWorkplace() {
    if (document.getElementById("workplace-name").value.length === 0) {
        document.getElementById("workplace-name").style.backgroundColor = "#ffcccb"
    } else {
        document.getElementById("workplace-name").style.backgroundColor = ""
        if (document.getElementById("data-save-button").classList[1] === "primary") {
            changeDataSaveButtonToWarning()
            setTimeout(function () {
                if (document.getElementById("data-save-button").classList[1] === "warning") {
                    const workshiftElement = document.getElementsByClassName("tag short-tag");
                    let workShiftsFromPage = []
                    for (let index = 0; index < workshiftElement.length; index++) {
                        workShiftsFromPage.push(workshiftElement[index].children[0].innerHTML)
                    }
                    let data = {
                        id: sessionStorage.getItem("selected_id"),
                        name: document.getElementById("workplace-name").value,
                        section: document.getElementById("workplace-section-selection").value,
                        productionDowntimeSelection: document.getElementById("production-downtime-selection").value,
                        productionDowntimeColor: document.getElementById("production-color").value,
                        powerOnPowerOffSelection: document.getElementById("poweron-poweroff-selection").value,
                        powerOnPowerOffColor: document.getElementById("poweroff-color").value,
                        workShifts: workShiftsFromPage,
                        countOkSelection: document.getElementById("count-ok-selection").value,
                        countOkColor: document.getElementById("ok-color").value,
                        countNokSelection: document.getElementById("count-nok-selection").value,
                        countNokColor: document.getElementById("nok-color").value,
                        mode: document.getElementById("workplace-mode-selection").value,
                        code: document.getElementById("code").value,
                        note: document.getElementById("workplace-note").value,
                    };
                    fetch("/save_workplace", {
                        method: "POST",
                        body: JSON.stringify(data)
                    }).then((response) => {
                        response.text().then(function () {
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
        } else if (document.getElementById("data-save-button").classList[1] === "warning") {
            changeDataSaveButtonToPrimary()
        }
    }
}

function saveWorkplaceSection() {
    if (document.getElementById("workplace-section-name").value.length === 0) {
        document.getElementById("workplace-section-name").style.backgroundColor = "#ffcccb"
    } else {
        document.getElementById("workplace-section-name").style.backgroundColor = ""
        if (document.getElementById("data-type-save-button").classList[1] === "primary") {
            changeDataTypeSaveButtonToWarning()
            setTimeout(function () {
                if (document.getElementById("data-type-save-button").classList[1] === "warning") {
                    let data = {
                        id: sessionStorage.getItem("selected_id"),
                        name: document.getElementById("workplace-section-name").value,
                        note: document.getElementById("workplace-section-note").value,
                    };
                    fetch("/save_workplace_section", {
                        method: "POST",
                        body: JSON.stringify(data)
                    }).then((response) => {
                        response.text().then(function () {
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
        } else if (document.getElementById("data-type-save-button").classList[1] === "warning") {
            changeDataTypeSaveButtonToPrimary()
        }
    }
}

function saveWorkplaceMode() {
    if (document.getElementById("workplace-mode-name").value.length === 0) {
        document.getElementById("workplace-mode-name").style.backgroundColor = "#ffcccb"
    } else {
        document.getElementById("workplace-mode-name").style.backgroundColor = ""
        if (document.getElementById("data-mode-save-button").classList[1] === "primary") {
            changeDataModeButtonToWarning();
            setTimeout(function () {
                if (document.getElementById("data-mode-save-button").classList[1] === "warning") {
                    let data = {
                        id: sessionStorage.getItem("selected_id"),
                        name: document.getElementById("workplace-mode-name").value,
                        downtimeDuration: document.getElementById("downtime-duration").value,
                        powerOffDuration: document.getElementById("poweroff-duration").value,
                        note: document.getElementById("workplace-mode-note").value,
                    };
                    fetch("/save_workplace_mode", {
                        method: "POST",
                        body: JSON.stringify(data)
                    }).then((response) => {
                        response.text().then(function () {
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
        } else if (document.getElementById("data-mode-save-button").classList[1] === "warning") {
            changeDataModeButtonToPrimary();
        }
    }
}