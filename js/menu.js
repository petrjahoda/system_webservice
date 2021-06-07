let charms = document.getElementById("right-charms")
if (sessionStorage.getItem("charm") === "open") {
    charms.classList.add("open")
} else {
    charms.classList.remove("open")
}

const menu = document.getElementById("menu")
menu.addEventListener("click", (event) => {
    let compacted = "compacted js-compact"
    if (document.getElementById("mainmenu").classList.contains("compacted")) {
        compacted = ""
    }
    let data = {
        key: "menu",
        value: compacted
    };
    fetch("/update_user_web_settings_from_web", {
        method: "POST",
        body: JSON.stringify(data)
    }).then(() => {
        location.reload()
    }).catch((error) => {
        console.log(error)
    });
})

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
    let userAgentString = navigator.userAgent;
    let chromeAgent = userAgentString.indexOf("Chrome") > -1;
    let IExplorerAgent = userAgentString.indexOf("MSIE") > -1 || userAgentString.indexOf("rv:") > -1;
    let firefoxAgent = userAgentString.indexOf("Firefox") > -1;
    let safariAgent = userAgentString.indexOf("Safari") > -1;
    if ((chromeAgent) && (safariAgent)) {
        safariAgent = false;
    }
    let operaAgent = userAgentString.indexOf("OP") > -1;
    if ((chromeAgent) && (operaAgent)) {
        chromeAgent = false;
    }
    if (safariAgent) {
        let request = new XMLHttpRequest();
        request.open("get", "/rest/login", false, "a", "false");
        request.send();
        document.execCommand("ClearAuthenticationCache")
        document.execCommand('ClearAuthenticationCache', true);
        window.location.replace(window.location.href);
    } else {
        $.ajax({
            type: "GET",
            url: window.location.href,
            dataType: 'json',
            async: true,
            username: "nobody",
            password: "nothing",
            data: '{ "comment" }'
        })
    }
})

function getLocaleFrom(chartData) {
    let locale = ""
    switch (chartData["Locale"]) {
        case "CsCZ": {
            locale = "cs";
            break;
        }
        case "DeDE": {
            locale = "de";
            break;
        }
        case "EnUS": {
            locale = "en";
            break;
        }
        case "EsES": {
            locale = "es";
            break;
        }
        case "FrFR": {
            locale = "fr";
            break;
        }
        case "ItIT": {
            locale = "it";
            break;
        }
        case "PlPL": {
            locale = "pl";
            break;
        }
        case "PtPT": {
            locale = "pt";
            break;
        }
        case "SkSK": {
            locale = "sk";
            break;
        }
        case "RuRU": {
            locale = "ru";
            break;
        }
    }
    return locale;
}
