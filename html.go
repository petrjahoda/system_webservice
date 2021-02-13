package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/petrjahoda/database"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"html/template"
	"net/http"
	"strings"
	"time"
)

type HomePageData struct {
	Version        string
	Company        string
	Alarms         string
	MenuOverview   string
	MenuWorkplaces string
	MenuCharts     string
	MenuStatistics string
	MenuData       string
	MenuSettings   string
}

func workplaces(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
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
	tmpl := template.Must(template.ParseFiles("./html/workplaces.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("MAIN", "Home page sent")
}
func charts(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
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
	tmpl := template.Must(template.ParseFiles("./html/charts.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("MAIN", "Home page sent")
}
func statistics(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
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
	tmpl := template.Must(template.ParseFiles("./html/statistics.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("MAIN", "Home page sent")
}
func data(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
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
	tmpl := template.Must(template.ParseFiles("./html/data.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("MAIN", "Home page sent")
}

func settings(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
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
	tmpl := template.Must(template.ParseFiles("./html/settings.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("MAIN", "Home page sent")
}

func updateProgramVersion() {
	logInfo("MAIN", "Writing program version into settings")
	timer := time.Now()
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("MAIN", "Problem opening database: "+err.Error())
		return
	}
	var existingSettings database.Setting
	db.Where("name=?", serviceName).Find(&existingSettings)
	existingSettings.Name = serviceName
	existingSettings.Value = version
	db.Save(&existingSettings)
	logInfo("MAIN", "Program version written into settings in "+time.Since(timer).String())
}

func basicAuth(h httprouter.Handle) httprouter.Handle {
	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		email, password, hasAuth := request.BasicAuth()
		cachedUsersSync.Lock()
		user, userFound := cachedUsers[email]
		cachedUsersSync.Unlock()
		userMatchesPassword := comparePasswords(user.Password, []byte(password))
		if hasAuth && userFound && userMatchesPassword{
			h(writer, request, params)
		} else {
			writer.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			http.Error(writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	}
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		logError("MAIN", "Passwords not matching")
		return false
	}
	return true
}