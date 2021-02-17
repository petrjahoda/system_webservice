package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/petrjahoda/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"html/template"
	"net/http"
	"strings"
	"time"
)

type Workplace struct {
	WorkplaceColor             string
	WorkplaceState             string
	WorkplaceName              string
	WorkplaceStateDuration     string
	WorkplaceProductivityToday string
	WorkplaceProductivityColor string
	Information                string
}

type WorkplaceSection struct {
	SectionName    string
	PanelCompacted string
	Workplaces     []Workplace
}

type WorkplacesData struct {
	Version           string
	Company           string
	Alarms            string
	MenuOverview      string
	MenuWorkplaces    string
	MenuCharts        string
	MenuStatistics    string
	MenuData          string
	MenuSettings      string
	WorkplaceSections []WorkplaceSection
	Compacted         string
}

func workplaces(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	ipAddress := strings.Split(request.RemoteAddr, ":")
	logInfo("MAIN", "Sending home page to "+ipAddress[0])
	email, _, _ := request.BasicAuth()

	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("MAIN", "Problem opening database: "+err.Error())
		return
	}
	var workplaceSections []database.WorkplaceSection
	db.Find(&workplaceSections)

	var sections []WorkplaceSection
	for _, workplaceSection := range workplaceSections {

		var section WorkplaceSection
		section.SectionName = workplaceSection.Name
		section.PanelCompacted = "display:block"
		userSettings := cachedUserSettings[email]
		for _, state := range userSettings.sectionStates {
			if state.section == workplaceSection.Name {
				if state.state != "expand" {
					section.PanelCompacted = "display:none"
				}
			}
		}

		var pageWorkplaces []Workplace
		var workplaces []database.Workplace
		db.Where("workplace_section_id = ?", workplaceSection.ID).Find(&workplaces)
		for _, workplace := range workplaces {
			var pageWorkplace Workplace
			pageWorkplace.WorkplaceName = workplace.Name
			var stateRecord database.StateRecord
			db.Where("workplace_id = ?", workplace.ID).Last(&stateRecord)

			var downtimeRecord database.DowntimeRecord
			db.Where("workplace_id = ?", workplace.ID).Where("date_time_end is null").Last(&downtimeRecord)
			var downtime database.Downtime
			db.Where("id = ?", downtimeRecord.DowntimeID).Find(&downtime)

			var orderRecord database.OrderRecord
			db.Where("workplace_id = ?", workplace.ID).Where("date_time_end is null").Last(&orderRecord)
			var order database.Order
			db.Where("id = ?", orderRecord.OrderID).Find(&order)

			var userRecord database.UserRecord
			db.Where("workplace_id = ?", workplace.ID).Where("date_time_end is null").Last(&userRecord)
			var user database.User
			db.Where("id = ?", userRecord.UserID).Find(&user)

			switch stateRecord.StateID {
			case 1:
				pageWorkplace.WorkplaceColor = "bg-green"
				pageWorkplace.WorkplaceProductivityColor = "bg-darkGreen"
				pageWorkplace.WorkplaceState = "mif-play"
				pageWorkplace.WorkplaceStateDuration = time.Since(stateRecord.DateTimeStart).Round(time.Second).String()
				pageWorkplace.WorkplaceProductivityToday = "23"
				pageWorkplace.Information = order.Name + " " + user.FirstName + " " + user.SecondName
			case 2:
				pageWorkplace.WorkplaceColor = "bg-orange"
				pageWorkplace.WorkplaceProductivityColor = "bg-darkOrange"
				pageWorkplace.WorkplaceState = "mif-pause"
				pageWorkplace.WorkplaceStateDuration = time.Since(stateRecord.DateTimeStart).Round(time.Second).String()
				pageWorkplace.WorkplaceProductivityToday = "23"
				pageWorkplace.Information = downtime.Name + " " + user.FirstName + " " + user.SecondName
			default:
				pageWorkplace.WorkplaceColor = "bg-red"
				pageWorkplace.WorkplaceState = "mif-stop"
				pageWorkplace.WorkplaceProductivityColor = "bg-darkRed"
				pageWorkplace.WorkplaceStateDuration = time.Since(stateRecord.DateTimeStart).Round(time.Second).String()
				pageWorkplace.WorkplaceProductivityToday = "23"
				pageWorkplace.Information = ""
			}

			pageWorkplaces = append(pageWorkplaces, pageWorkplace)
		}
		section.Workplaces = pageWorkplaces
		sections = append(sections, section)
	}

	var data WorkplacesData
	data.Version = version
	data.Company = cachedCompanyName
	data.MenuOverview = getLocale(email, "menu-overview")
	data.MenuWorkplaces = getLocale(email, "menu-workplaces")
	data.MenuCharts = getLocale(email, "menu-charts")
	data.MenuStatistics = getLocale(email, "menu-statistics")
	data.MenuData = getLocale(email, "menu-data")
	data.MenuSettings = getLocale(email, "menu-settings")
	data.WorkplaceSections = sections
	data.Compacted = cachedUserSettings[email].menuState
	tmpl := template.Must(template.ParseFiles("./html/workplaces.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("MAIN", "Home page sent")
}
