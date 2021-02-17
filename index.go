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
	user, _ := cachedUsers[email]
	switch user.Locale {
	case "CsCZ":
		{
			menuOverview = cachedLocales[locale].CsCZ
		}
	case "DeDE":
		{
			menuOverview = cachedLocales[locale].DeDE
		}
	case "EnUS":
		{
			menuOverview = cachedLocales[locale].EnUS
		}
	case "EsES":
		{
			menuOverview = cachedLocales[locale].EsES
		}
	case "FrFR":
		{
			menuOverview = cachedLocales[locale].FrFR
		}
	case "ItIT":
		{
			menuOverview = cachedLocales[locale].ItIT
		}
	case "PlPL":
		{
			menuOverview = cachedLocales[locale].PlPL
		}
	case "PtPT":
		{
			menuOverview = cachedLocales[locale].PtPT
		}
	case "SkSK":
		{
			menuOverview = cachedLocales[locale].SkSK
		}
	case "RuRU":
		{
			menuOverview = cachedLocales[locale].RuRU
		}
	default:
		{
			menuOverview = cachedLocales[locale].EnUS
		}
	}
	return menuOverview
}
