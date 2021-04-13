const settingsSelection = document.getElementById("data-selection")
settingsSelection.addEventListener("change", (event) => {
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
    } else if (event.target.id === "data-new-button-type" || event.target.id === "data-new-button-type-mif") {
        loadSettings();
        loadDetails(null, "second");
    } else if (event.target.id === "data-new-button-type-extended" || event.target.id === "data-new-button-mif-type-extended") {
        loadSettings();
        loadDetails(null, "third");
    }
    if (event.target.id === "data-save-button" || event.target.id === "data-save-button-mif")  {
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
        switch (selection) {
            case "workplaces" : {
                saveWorkplaceMode();
                break
            }
        }
    } else if (event.target.id === "data-new-port-button" || event.target.id === "data-new-port-button-mif") {
        loadDevicePortDetails(null);
    } else if (event.target.id === "data-new-workplace-port-button" || event.target.id === "data-new-workplace-port-button-mif") {
        loadWorkplacePortDetails(null);
    } else if (event.target.id === "workplace-port-delete-button" || event.target.id === "workplace-port-delete-button-mif") {
        deleteWorkplacePort();
    } else if (event.target.id === "port-save-button" || event.target.id === "port-save-button-mif") {
        saveDevicePortDetails();
    } else if (event.target.id === "workplace-port-save-button" || event.target.id === "workplace-port-save-button-mif") {
        saveWorkplacePortDetails();
    } else {
        const tablePortSelectedId = event.target.parentElement.parentElement.parentElement.parentElement.parentElement.id
        let portTable = Metro.getPlugin("#data-port-table", "table");
        let workplacePortTable = Metro.getPlugin("#data-workplace-port-table", "table");
        if (tablePortSelectedId === "data-port-table" && portTable.getSelectedItems().length > 0 && event.target.type === "radio") {
            let selectedItem = portTable.getSelectedItems()[0][0]
            loadDevicePortDetails(selectedItem);
        } else if (tablePortSelectedId === "data-workplace-port-table" && workplacePortTable.getSelectedItems().length > 0 && event.target.type === "radio") {
            let selectedItem = workplacePortTable.getSelectedItems()[0][0]
            loadWorkplacePortDetails(selectedItem);
        }
    }
})

