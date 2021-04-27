function changeWorkplacePortDeleteButtonToWarning() {
    document.getElementById("workplace-port-delete-button").classList.remove("alert")
    document.getElementById("workplace-port-delete-button").classList.add("warning")
    document.getElementById("workplace-port-delete-button-mif").classList.remove("mif-floppy-disk")
    document.getElementById("workplace-port-delete-button-mif").classList.add("mif-cross")
}

function changeWorkplacePortDeleteButtonToPrimary() {
    document.getElementById("workplace-port-delete-button").classList.remove("warning")
    document.getElementById("workplace-port-delete-button").classList.add("alert")
    document.getElementById("workplace-port-delete-button-mif").classList.remove("mif-cross")
    document.getElementById("workplace-port-delete-button-mif").classList.add("mif-minus")
}

function deleteWorkplacePort() {
    if (document.getElementById("workplace-port-delete-button").classList[1] === "alert") {
        changeWorkplacePortDeleteButtonToWarning();
        setTimeout(function () {
            if (document.getElementById("workplace-port-delete-button").classList[1] === "warning") {
                let data = {
                    id: sessionStorage.getItem("workplace_port_selected_id"),
                };
                fetch("/delete_workplace_port", {
                    method: "POST",
                    body: JSON.stringify(data)
                }).then((response) => {
                    response.text().then(function (data) {
                        let result = JSON.parse(data);
                        updateCharm(result["Result"])
                        document.getElementById("workplace-port-container").innerHTML = ""
                        loadDetails(sessionStorage.getItem("selected_id"), "first");
                    });
                }).catch((error) => {
                    updateCharm("ERR: " + error)
                    document.getElementById("workplace-port-container").innerHTML = ""
                    loadDetails(sessionStorage.getItem("selected_id"), "first");
                });
            }
        }, 2500);
    } else if (document.getElementById("workplace-port-delete-button").classList[1] === "warning") {
        changeWorkplacePortDeleteButtonToPrimary();
    }
}