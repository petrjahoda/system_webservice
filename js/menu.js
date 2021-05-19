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
        location.reload()
    }).catch((error) => {
        console.log(error)
    });
})

let charms = document.getElementById("right-charms")
if (sessionStorage.getItem("charm") === "open") {
    charms.classList.add("open")
} else {
    charms.classList.remove("open")
}

let infoButton = document.getElementById("info-button")
infoButton.addEventListener('click', () => {
    if (sessionStorage.getItem("charm") === "open") {
        sessionStorage.setItem("charm", "close")
        charms.classList.remove("open")
    } else {
        charms.classList.add("open")
        sessionStorage.setItem("charm", "open")
    }
})

function updateCharm(text) {
    charms.innerHTML = text + "<br>" + charms.innerHTML
}

const logout = document.getElementById("logout-button")
logout.addEventListener('click', () => {
    let request = new XMLHttpRequest();
    request.open("get", "/rest/login", false, "a", "false");
    request.send();
    window.location.replace("/");
    document.execCommand("ClearAuthenticationCache")
    document.execCommand('ClearAuthenticationCache', false);
    window.location= ("http://log:out@localhost:82/")
})
