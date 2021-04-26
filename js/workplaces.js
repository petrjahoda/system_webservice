const refreshButton = document.getElementById("data-refresh-button")
refreshButton.addEventListener('click', () => {
    location.reload()
})

let timeleft = 60;
const downloadTimer = setInterval(function () {
    if (timeleft <= 0) {
        clearInterval(downloadTimer);
        location.reload()
    }
    document.getElementById("progress-bar").value = 60 - timeleft;
    timeleft -= 1;
}, 1000);


function dataCollapse(element) {
    console.log(element.dataset.titleCaption + " collapsed")
    let data = {
        settings: "section",
        state: "collapse-" + element.dataset.titleCaption
    };
    fetch("/update_user_settings", {
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
        settings: "section",
        state: "expand-" + element.dataset.titleCaption
    };
    fetch("/update_user_settings", {
        method: "POST",
        body: JSON.stringify(data)
    }).then(() => {
    }).catch((error) => {
        console.log(error)
    });
}