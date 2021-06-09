let timer = 60
let timeLeft = timer;
const downloadTimer = setInterval(function () {
    if (timeLeft <= 0) {
        document.getElementById("progress-bar").value = 0
        updateWorkplaces();
    }
    document.getElementById("progress-bar").value = timer - timeLeft;
    timeLeft -= 1;
}, 1000);

const refreshButton = document.getElementById("data-refresh-button")
refreshButton.addEventListener('click', () => {
    updateWorkplaces();
})

function updateWorkplaces() {
    document.getElementById("loader").hidden = false
    let data = {
        email: document.getElementById("user-info").title
    };
    fetch("/update_workplaces", {
        method: "POST",
        body: JSON.stringify(data)

    }).then((response) => {
        response.text().then(function (data) {
            if (data.includes("ERR: ")) {
                JSON.parse(data);
            } else {
                document.getElementById("content-wrapper").innerHTML = data
                timeLeft = timer
                document.getElementById("progress-bar").value = timer - timeLeft;
                document.getElementById("loader").hidden = true
            }
        });
    }).catch(() => {
        document.getElementById("loader").hidden = true
    });
}

function dataCollapse(element) {
    let data = {
        key: element.dataset.titleCaption,
        value: "display:none"
    };
    fetch("/update_user_web_settings_from_web", {
        method: "POST",
        body: JSON.stringify(data)
    }).then(() => {
    }).catch(() => {
    });
}

function dataExpand(element) {
    let data = {
        key: element.dataset.titleCaption,
        value: "display:block"
    };
    fetch("/update_user_web_settings_from_web", {
        method: "POST",
        body: JSON.stringify(data)
    }).then(() => {
    }).catch(() => {
    });
}