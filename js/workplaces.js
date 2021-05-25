let timeleft = 60;
const downloadTimer = setInterval(function () {
    if (timeleft <= 0) {
        clearInterval(downloadTimer);
        updateWorkplaces();
    }
    document.getElementById("progress-bar").value = 60 - timeleft;
    timeleft -= 1;
}, 1000);

const refreshButton = document.getElementById("data-refresh-button")
refreshButton.addEventListener('click', () => {
    updateWorkplaces();
})

function updateWorkplaces() {
    document.getElementById("loader").hidden = false
    fetch("/update_workplaces", {
        method: "POST",
    }).then((response) => {
        response.text().then(function (data) {
            if (data.includes("ERR: ")) {
                let result = JSON.parse(data);
                updateCharm(result["Result"])
            } else {
                document.getElementById("content-wrapper").innerHTML = data
                updateCharm(document.getElementById("hidden-information").innerText)
                timeleft = 60
                document.getElementById("progress-bar").value = 60 - timeleft;
                document.getElementById("loader").hidden = true
            }
        });
    }).catch((error) => {
        updateCharm("ERR: " + error)
        document.getElementById("loader").hidden = true
    });
}

function dataCollapse(element) {
    console.log(element.dataset.titleCaption + " collapsed")
    let data = {
        key: element.dataset.titleCaption,
        value: "display:none"
    };
    fetch("/update_user_web_settings_from_web", {
        method: "POST",
        body: JSON.stringify(data)
    }).then(() => {
    }).catch((error) => {
        console.log(error)
    });
}

function dataExpand(element) {
    console.log(element.dataset.titleCaption + " expanded")
    let data = {
        key: element.dataset.titleCaption,
        value: "display:block"
    };
    fetch("/update_user_web_settings_from_web", {
        method: "POST",
        body: JSON.stringify(data)
    }).then(() => {
    }).catch((error) => {
        console.log(error)
    });
}