function loadSettings() {
    let data = {
        data: document.getElementById("data-selection").value,
    };
    fetch("/load_settings_data", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            if (data.includes("ERR: ")) {
                JSON.parse(data);
            } else {
                document.getElementById("settings-container").innerHTML = data
            }
        });
    }).catch(() => {
    });
}

function loadDetails(selectedItem, type) {
    let selection = document.getElementById("data-selection").value
    let data = {
        data: selection,
        id: selectedItem,
        type: type
    };
    sessionStorage.setItem("selected_id", selectedItem)
    fetch("/load_settings_detail", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            if (data.includes("ERR: ")) {
                JSON.parse(data);
            } else {
                document.getElementById("settings-container-detail").innerHTML = data
            }
        });
    }).catch(() => {
    });
}

function loadWorkplacePortDetails(selectedPort) {
    sessionStorage.setItem("workplace_port_selected_id", selectedPort)
    let data = {
        data: selectedPort,
        workplaceId: sessionStorage.getItem("selected_id")
    };
    fetch("/load_workplace_port_detail", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            if (data.includes("ERR: ")) {
                JSON.parse(data);
            } else {
                document.getElementById("workplace-port-container").innerHTML = data
                setTimeout(function () {
                    document.getElementById("workplace-port-container").scrollIntoView();
                }, 100);
                if (document.getElementById("workplace-port-created-at").value === "0001-01-01T00:00:00") {
                    document.getElementById("workplace-port-delete-button").hidden = true
                }
            }
        });
    }).catch(() => {
    });
}

function loadDevicePortDetails(selectedPort) {
    sessionStorage.setItem("device_port_selected_id", selectedPort)
    let data = {
        data: selectedPort,
        deviceId: sessionStorage.getItem("selected_id"),
    };
    fetch("/load_device_port_detail", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            if (data.includes("ERR: ")) {
                JSON.parse(data);
            } else {
                document.getElementById("port-container").innerHTML = data
                setTimeout(function () {
                    document.getElementById("port-container").scrollIntoView();
                }, 100);
            }
        });
    }).catch(() => {
    });
}