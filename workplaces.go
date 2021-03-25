package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/petrjahoda/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Workplace struct {
	WorkplaceColor             string
	UserName                   string
	OrderName                  string
	ProductName                string
	WorkplaceState             string
	WorkplaceStateName         string
	WorkplaceName              string
	WorkplaceStateDuration     string
	WorkplaceProductivityToday string
	WorkplaceProductivityColor string
	Information                string
	OrderDuration              string
	UserInformation            string
	BreakdownInformation       string
	AlarmInformation           string
	TodayDate                  string
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
	UserEmail         string
	UserName          string
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
			var state database.State
			db.Where("id = ?", stateRecord.StateID).Find(&state)

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

			var breakdownRecord database.BreakdownRecord
			db.Where("workplace_id = ?", workplace.ID).Where("date_time_end is null").Last(&breakdownRecord)
			var breakdown database.Breakdown
			db.Where("id = ?", breakdownRecord.BreakdownID).Find(&breakdown)

			var alarmRecord database.AlarmRecord
			db.Where("workplace_id = ?", workplace.ID).Where("date_time_end is null").Last(&alarmRecord)
			var alarm database.Alarm
			db.Where("id = ?", alarmRecord.AlarmID).Find(&alarm)

			productivity := calculateProductivity(db, workplace)

			switch stateRecord.StateID {
			case 1:
				pageWorkplace.WorkplaceColor = "background-color: " + state.Color
				pageWorkplace.WorkplaceProductivityColor = "bg-darkGreen"
				pageWorkplace.WorkplaceState = "mif-play"
				pageWorkplace.WorkplaceStateName = getLocale(email, "production")
				pageWorkplace.WorkplaceStateDuration = time.Since(stateRecord.DateTimeStart).Round(time.Minute).String()
				pageWorkplace.WorkplaceProductivityToday = productivity
				if len(order.Name) > 0 {
					pageWorkplace.Information = getLocale(email, "order-name") + ": " + order.Name
					pageWorkplace.OrderDuration = "[" + time.Now().Sub(orderRecord.DateTimeStart).Round(time.Minute).String() + "]"
					pageWorkplace.UserInformation = getLocale(email, "user-name") + ": " + user.FirstName + " " + user.SecondName
				} else if len(user.FirstName) > 0 {
					pageWorkplace.Information = getLocale(email, "order-name") + ": -"
					pageWorkplace.UserInformation = getLocale(email, "user-name") + ": " + user.FirstName + " " + user.SecondName
				} else {
					pageWorkplace.Information = getLocale(email, "order-name") + ": -"
					pageWorkplace.UserInformation = getLocale(email, "user-name") + ": -"
				}
				if len(alarm.Name) > 0 {
					pageWorkplace.AlarmInformation = getLocale(email, "alarm-name") + ": " + alarm.Name
				} else {
					pageWorkplace.AlarmInformation = getLocale(email, "alarm-name") + ": -"
				}
				if len(breakdown.Name) > 0 {
					pageWorkplace.BreakdownInformation = getLocale(email, "breakdown-name") + ": " + breakdown.Name
				} else {
					pageWorkplace.BreakdownInformation = getLocale(email, "breakdown-name") + ": -"
				}
				pageWorkplace.TodayDate = time.Now().Format("02.01.2006")
			case 2:
				pageWorkplace.WorkplaceColor = "background-color: " + state.Color
				pageWorkplace.WorkplaceProductivityColor = "bg-darkOrange"
				pageWorkplace.WorkplaceState = "mif-pause"
				pageWorkplace.WorkplaceStateName = downtime.Name
				pageWorkplace.WorkplaceStateDuration = time.Since(stateRecord.DateTimeStart).Round(time.Minute).String()
				pageWorkplace.WorkplaceProductivityToday = productivity
				if len(order.Name) > 0 {
					pageWorkplace.Information = getLocale(email, "order-name") + ": " + order.Name
					pageWorkplace.OrderDuration = "[" + time.Now().Sub(orderRecord.DateTimeStart).Round(time.Minute).String() + "]"
					pageWorkplace.UserInformation = getLocale(email, "user-name") + ": " + user.FirstName + " " + user.SecondName
				} else if len(user.FirstName) > 0 {
					pageWorkplace.Information = getLocale(email, "order-name") + ": -"
					pageWorkplace.UserInformation = getLocale(email, "user-name") + ": " + user.FirstName + " " + user.SecondName
				} else {
					pageWorkplace.Information = getLocale(email, "order-name") + ": -"
					pageWorkplace.UserInformation = getLocale(email, "user-name") + ": -"
				}
				if len(alarm.Name) > 0 {
					pageWorkplace.AlarmInformation = getLocale(email, "alarm-name") + ": " + alarm.Name
				} else {
					pageWorkplace.AlarmInformation = getLocale(email, "alarm-name") + ": -"
				}
				if len(breakdown.Name) > 0 {
					pageWorkplace.BreakdownInformation = getLocale(email, "breakdown-name") + ": " + breakdown.Name
				} else {
					pageWorkplace.BreakdownInformation = getLocale(email, "breakdown-name") + ": -"
				}
				pageWorkplace.TodayDate = time.Now().Format("02.01.2006")
			default:
				pageWorkplace.WorkplaceColor = "background-color: " + state.Color
				pageWorkplace.WorkplaceState = "mif-stop"
				pageWorkplace.WorkplaceStateName = getLocale(email, "poweroff")
				pageWorkplace.WorkplaceProductivityColor = "bg-darkRed"
				pageWorkplace.WorkplaceStateDuration = time.Since(stateRecord.DateTimeStart).Round(time.Minute).String()
				pageWorkplace.WorkplaceProductivityToday = productivity
				pageWorkplace.Information = getLocale(email, "order-name") + ": -"
				pageWorkplace.UserInformation = getLocale(email, "user-name") + ": -"
				if len(alarm.Name) > 0 {
					pageWorkplace.AlarmInformation = getLocale(email, "alarm-name") + ": " + alarm.Name
				} else {
					pageWorkplace.AlarmInformation = getLocale(email, "alarm-name") + ": -"
				}
				if len(breakdown.Name) > 0 {
					pageWorkplace.BreakdownInformation = getLocale(email, "breakdown-name") + ": " + breakdown.Name
				} else {
					pageWorkplace.BreakdownInformation = getLocale(email, "breakdown-name") + ": -"
				}
				pageWorkplace.TodayDate = time.Now().Format("02.01.2006")
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
	data.UserEmail = email
	data.UserName = cachedUsersByEmail[email].FirstName + " " + cachedUsersByEmail[email].SecondName
	tmpl := template.Must(template.ParseFiles("./html/workplaces.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("MAIN", "Home page sent")
}

func calculateProductivity(db *gorm.DB, workplace database.Workplace) string {
	noon := time.Date(time.Now().UTC().Year(), time.Now().UTC().Month(), time.Now().UTC().Day(), 0, 0, 0, 0, time.Now().Location())
	var stateRecords []database.StateRecord
	db.Raw("select * from state_records where id >= (select id from state_records where workplace_id=? and date_time_start < ? order by date_time_start desc limit 1) and workplace_id=? order by date_time_start asc", workplace.ID, noon, workplace.ID).Find(&stateRecords)
	firstRecord := true
	var productivitySum time.Duration
	for index, record := range stateRecords {
		if firstRecord {
			if record.StateID == 1 {
				if index+1 > len(stateRecords)-1 {
					productivitySum += time.Now().UTC().Sub(noon)
				} else {
					productivitySum += stateRecords[index+1].DateTimeStart.UTC().Sub(noon)
				}
			}
			firstRecord = false
		} else {
			if record.StateID == 1 {
				if index+1 > len(stateRecords)-1 {
					productivitySum += time.Now().Sub(record.DateTimeStart.UTC())
				} else {
					productivitySum += stateRecords[index+1].DateTimeStart.UTC().Sub(record.DateTimeStart.UTC())
				}
			}
		}
	}
	productivityAsFloat := (productivitySum.Seconds() / time.Now().Sub(noon).Seconds()) * 100
	productivity := strconv.FormatFloat(productivityAsFloat, 'f', 1, 64)
	return productivity
}
