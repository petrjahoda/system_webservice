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