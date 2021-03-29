const menu = document.getElementById("menu")

menu.addEventListener("click", (event) => {
    console.log("changing menu state")
    let data = {
        settings: "menu",
        state: document.getElementById("mainmenu").classList.contains("compacted").toString()
    };
    fetch("/update_user_settings", {
        method: "POST",
        body: JSON.stringify(data)
    }).then(() => {
        setTimeout(function () {
            resizeCharts()
        }, 100);

    }).catch((error) => {
        console.log(error)
    });
})

