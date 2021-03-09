const dataOkButton = document.getElementById("data-ok-button")
dataOkButton.addEventListener("click", (event) => {
    console.log("getting settings for " + document.getElementById("data-selection").value)
    let data = {
        data: document.getElementById("data-selection").value,
    };
    fetch("/get_settings_data", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            document.getElementById("data-save-button").classList.add("disabled")
            document.getElementById("settings-container").innerHTML = data
        });
    }).catch((error) => {
        console.log(error)
    });
})

// const container = document.getElementById("settings-container")
// container.addEventListener("click", (event) => {
//     let table = Metro.getPlugin("#data-table", "table");
//     if (table.getSelectedItems().length > 0 && event.target.type === "radio") {
//         let selectedItem = table.getSelectedItems()[0][0]
//         let selection = document.getElementById("data-selection").value
//         switch (selection) {
//             case "alarms" : {
//                 document.getElementById("data-save-button").classList.remove("disabled")
//                 console.log("getting data for " + selectedItem)
//                 let data = {
//                     data: "alarms",
//                     name: selectedItem
//                 };
//                 fetch("/get_detail_settings", {
//                     method: "POST",
//                     body: JSON.stringify(data)
//                 }).then((response) => {
//                     response.text().then(function (data) {
//                         let result = JSON.parse(data);
//                         document.getElementById("alarm-name").value = result["AlarmName"]
//                         document.getElementById("sql-command").value = result["SqlCommand"]
//                         document.getElementById("message-header").value = result["MessageHeader"]
//                         document.getElementById("message-text").value = result["MessageText"]
//                         document.getElementById("recipients").value = result["Recipients"]
//                         document.getElementById("url").value = result["Url"]
//                         document.getElementById("pdf").value = result["Pdf"]
//                         document.getElementById("created-at").value = result["CreatedAt"]
//                         document.getElementById("updated-at").value = result["UpdatedAt"]
//
//                         var select = Metro.getPlugin("#workplace-selection", 'select');
//                         select.data({
//                             "myvalue": "mytext",
//                             "myvalue2": "mytext2"
//                         });
//
//                         selection = document.getElementById("workplace-selection")
//                         for (const selectionElement of selection) {
//                             console.log(selectionElement)
//
//                         }
//                         for (const selectionElement of selection) {
//                             console.log(selectionElement)
//                             if (selectionElement.value === "myvalue2") {
//                                 selectionElement.setAttribute("selected", "selected")
//                             }
//                         }
//
//                         // let workplaces = document.getElementById("workplace-selection")
//                         // for (const workplace of result["Workplaces"]) {
//                         //     let option = document.createElement("option");
//                         //     option.value = workplace;
//                         //     option.text = workplace;
//                         //     if (workplace === result["WorkplaceName"]) {
//                         //         console.log("match")
//                         //         option.selected = true
//                         //     }
//                         //     workplaces.add(option)
//                         // }
//                     });
//                 }).catch((error) => {
//                     console.log(error)
//                 });
//             }
//         }
//     }
// })


