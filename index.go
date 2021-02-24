package main

import (
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
	"strings"
)

type IndexPageData struct {
	Version        string
	Company        string
	Alarms         string
	MenuOverview   string
	MenuWorkplaces string
	MenuCharts     string
	MenuStatistics string
	MenuData       string
	MenuSettings   string
	Compacted      string
}

func index(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	ipAddress := strings.Split(request.RemoteAddr, ":")
	logInfo("MAIN", "Sending home page to "+ipAddress[0])
	email, _, _ := request.BasicAuth()
	var data IndexPageData
	data.Version = version
	data.Company = cachedCompanyName
	data.MenuOverview = getLocale(email, "menu-overview")
	data.MenuWorkplaces = getLocale(email, "menu-workplaces")
	data.MenuCharts = getLocale(email, "menu-charts")
	data.MenuStatistics = getLocale(email, "menu-statistics")
	data.MenuData = getLocale(email, "menu-data")
	data.MenuSettings = getLocale(email, "menu-settings")
	data.Compacted = cachedUserSettings[email].menuState
	tmpl := template.Must(template.ParseFiles("./html/index.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("MAIN", "Home page sent")
}

func getLocale(email string, locale string) string {
	var menuOverview string
	user, _ := cachedUsersByEmail[email]
	switch user.Locale {
	case "CsCZ":
		{
			menuOverview = cachedLocalesByName[locale].CsCZ
		}
	case "DeDE":
		{
			menuOverview = cachedLocalesByName[locale].DeDE
		}
	case "EnUS":
		{
			menuOverview = cachedLocalesByName[locale].EnUS
		}
	case "EsES":
		{
			menuOverview = cachedLocalesByName[locale].EsES
		}
	case "FrFR":
		{
			menuOverview = cachedLocalesByName[locale].FrFR
		}
	case "ItIT":
		{
			menuOverview = cachedLocalesByName[locale].ItIT
		}
	case "PlPL":
		{
			menuOverview = cachedLocalesByName[locale].PlPL
		}
	case "PtPT":
		{
			menuOverview = cachedLocalesByName[locale].PtPT
		}
	case "SkSK":
		{
			menuOverview = cachedLocalesByName[locale].SkSK
		}
	case "RuRU":
		{
			menuOverview = cachedLocalesByName[locale].RuRU
		}
	default:
		{
			menuOverview = cachedLocalesByName[locale].EnUS
		}
	}
	return menuOverview
}
