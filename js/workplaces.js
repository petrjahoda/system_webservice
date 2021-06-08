let timer = 30
let timeleft = timer;
const downloadTimer = setInterval(function () {
    if (timeleft <= 0) {
        location.reload()
    }
    document.getElementById("progress-bar").value = timer - timeleft;
    timeleft -= 1;
}, 1000);

const refreshButton = document.getElementById("data-refresh-button")
refreshButton.addEventListener('click', () => {
    location.reload()
})

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