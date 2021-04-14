package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/petrjahoda/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"html/template"
	"net/http"
	"sort"
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
	BreakdownDuration          string
	AlarmInformation           string
	AlarmDuration              string
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
	timer := time.Now()
	go updatePageCount("workplaces")
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
	var downtimeRecords []database.DowntimeRecord
	db.Where("date_time_end is null").Find(&downtimeRecords)
	cachedDowntimeRecords := make(map[int]database.DowntimeRecord)
	for _, downtimeRecord := range downtimeRecords {
		cachedDowntimeRecords[downtimeRecord.WorkplaceID] = downtimeRecord
	}
	var orderRecords []database.OrderRecord
	db.Where("date_time_end is null").Find(&orderRecords)
	cachedOrderRecords := make(map[int]database.OrderRecord)
	for _, orderRecord := range orderRecords {
		cachedOrderRecords[orderRecord.WorkplaceID] = orderRecord
	}
	var userRecords []database.UserRecord
	db.Where("date_time_end is null").Find(&userRecords)
	cachedUserRecords := make(map[int]database.UserRecord)
	for _, userRecord := range userRecords {
		cachedUserRecords[userRecord.WorkplaceID] = userRecord
	}
	var breakdownRecords []database.BreakdownRecord
	db.Where("date_time_end is null").Find(&breakdownRecords)
	cachedBreakdownRecords := make(map[int]database.BreakdownRecord)
	for _, breakdownRecord := range breakdownRecords {
		cachedBreakdownRecords[breakdownRecord.WorkplaceID] = breakdownRecord
	}
	var alarmRecords []database.AlarmRecord
	db.Where("date_time_end is null").Find(&alarmRecords)
	cachedAlarmRecords := make(map[int]database.AlarmRecord)
	for _, alarmRecord := range alarmRecords {
		cachedAlarmRecords[alarmRecord.WorkplaceID] = alarmRecord
	}
	var stateRecords []database.StateRecord
	db.Raw("select * from state_records where id in (select distinct max(id) as id from state_records group by workplace_id)").Find(&stateRecords)
	cachedStateRecords := make(map[int]database.StateRecord)
	for _, stateRecord := range stateRecords {
		cachedStateRecords[stateRecord.WorkplaceID] = stateRecord
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
		var tempWorkplaces []database.Workplace
		for _, workplace := range cachedWorkplacesById {
			tempWorkplaces = append(tempWorkplaces, workplace)
		}
		sort.Slice(tempWorkplaces, func(i, j int) bool {
			return tempWorkplaces[i].Name < tempWorkplaces[j].Name
		})

		var pageWorkplaces []Workplace
		for _, workplace := range tempWorkplaces {
			if uint(workplace.WorkplaceSectionID) == workplaceSection.ID {
				var pageWorkplace Workplace
				pageWorkplace.WorkplaceName = workplace.Name
				loc, _ := time.LoadLocation(location)
				data := downloadData(db, time.Date(time.Now().UTC().Year(), time.Now().UTC().Month(), time.Now().UTC().Day(), 0, 0, 0, 0, time.Now().Location()), time.Now().In(loc), workplace.ID, loc, email, 1)

				var totalDuration time.Duration
				for _, duration := range data {
					totalDuration = duration
				}
				startOfToday := time.Date(time.Now().In(loc).Year(), time.Now().In(loc).Month(), time.Now().In(loc).Day(), 0, 0, 0, 0, loc)
				totalTodayDuration := time.Now().In(loc).Sub(startOfToday)
				stateRecord := cachedStateRecords[int(workplace.ID)]
				state := cachedStatesById[uint(stateRecord.StateID)]
				switch stateRecord.StateID {
				case 1:
					pageWorkplace.WorkplaceColor = "background-color: " + state.Color
					pageWorkplace.WorkplaceProductivityColor = "bg-darkGreen"
					pageWorkplace.WorkplaceState = "mif-play"
					pageWorkplace.WorkplaceStateName = getLocale(email, "production")
					pageWorkplace.WorkplaceStateDuration = time.Since(stateRecord.DateTimeStart).Round(time.Second).String()
					pageWorkplace.WorkplaceProductivityToday = strconv.FormatFloat((totalDuration.Seconds()/totalTodayDuration.Seconds())*100, 'f', 1, 64)
					orderRecordId := cachedOrderRecords[int(workplace.ID)].OrderID
					userRecordId := cachedUserRecords[int(workplace.ID)].UserID
					breakdownRecordId := cachedBreakdownRecords[int(workplace.ID)].BreakdownID
					alarmRecordId := cachedAlarmRecords[int(workplace.ID)].AlarmID
					if orderRecordId > 0 {
						pageWorkplace.Information = getLocale(email, "order-name") + ": " + cachedOrdersById[uint(orderRecordId)].Name
						pageWorkplace.OrderDuration = "[" + time.Now().Sub(cachedOrderRecords[int(workplace.ID)].DateTimeStart).Round(time.Second).String() + "]"
						pageWorkplace.UserInformation = getLocale(email, "user-name") + ": " + cachedUsersById[uint(userRecordId)].FirstName + " " + cachedUsersById[uint(userRecordId)].SecondName
					} else if userRecordId > 0 {
						pageWorkplace.Information = getLocale(email, "order-name") + ": -"
						pageWorkplace.UserInformation = getLocale(email, "user-name") + ": " + cachedUsersById[uint(userRecordId)].FirstName + " " + cachedUsersById[uint(userRecordId)].SecondName
					} else {
						pageWorkplace.Information = getLocale(email, "order-name") + ": -"
						pageWorkplace.UserInformation = getLocale(email, "user-name") + ": -"
					}
					if alarmRecordId > 0 {
						pageWorkplace.AlarmInformation = getLocale(email, "alarm-name") + ": " + cachedAlarmsById[uint(alarmRecordId)].Name
						pageWorkplace.AlarmDuration = "[" + time.Now().Sub(cachedAlarmRecords[int(workplace.ID)].DateTimeStart).Round(time.Second).String() + "]"
					} else {
						pageWorkplace.AlarmInformation = getLocale(email, "alarm-name") + ": -"
						pageWorkplace.AlarmDuration = ""
					}
					if breakdownRecordId > 0 {
						pageWorkplace.BreakdownInformation = getLocale(email, "breakdown-name") + ": " + cachedBreakdownsById[uint(breakdownRecordId)].Name
						pageWorkplace.BreakdownDuration = "[" + time.Now().Sub(cachedBreakdownRecords[int(workplace.ID)].DateTimeStart).Round(time.Second).String() + "]"
					} else {
						pageWorkplace.BreakdownInformation = getLocale(email, "breakdown-name") + ": -"
						pageWorkplace.BreakdownDuration = ""
					}
					pageWorkplace.TodayDate = time.Now().Format("02.01.2006")
				case 2:
					pageWorkplace.WorkplaceColor = "background-color: " + state.Color
					pageWorkplace.WorkplaceProductivityColor = "bg-darkOrange"
					pageWorkplace.WorkplaceState = "mif-pause"
					downtimeRecordId := cachedDowntimeRecords[int(workplace.ID)].DowntimeID
					pageWorkplace.WorkplaceStateName = cachedDowntimesById[uint(downtimeRecordId)].Name
					pageWorkplace.WorkplaceStateDuration = time.Since(stateRecord.DateTimeStart).Round(time.Second).String()
					pageWorkplace.WorkplaceProductivityToday = strconv.FormatFloat((totalDuration.Seconds()/totalTodayDuration.Seconds())*100, 'f', 1, 64)
					orderRecordId := cachedOrderRecords[int(workplace.ID)].OrderID
					userRecordId := cachedUserRecords[int(workplace.ID)].UserID
					breakdownRecordId := cachedBreakdownRecords[int(workplace.ID)].BreakdownID
					alarmRecordId := cachedAlarmRecords[int(workplace.ID)].AlarmID
					if orderRecordId > 0 {
						pageWorkplace.Information = getLocale(email, "order-name") + ": " + cachedOrdersById[uint(orderRecordId)].Name
						pageWorkplace.OrderDuration = "[" + time.Now().Sub(cachedOrderRecords[int(workplace.ID)].DateTimeStart).Round(time.Second).String() + "]"
						pageWorkplace.UserInformation = getLocale(email, "user-name") + ": " + cachedUsersById[uint(userRecordId)].FirstName + " " + cachedUsersById[uint(userRecordId)].SecondName
					} else if userRecordId > 0 {
						pageWorkplace.Information = getLocale(email, "order-name") + ": -"
						pageWorkplace.UserInformation = getLocale(email, "user-name") + ": " + cachedUsersById[uint(userRecordId)].FirstName + " " + cachedUsersById[uint(userRecordId)].SecondName
					} else {
						pageWorkplace.Information = getLocale(email, "order-name") + ": -"
						pageWorkplace.UserInformation = getLocale(email, "user-name") + ": -"
					}
					if alarmRecordId > 0 {
						pageWorkplace.AlarmInformation = getLocale(email, "alarm-name") + ": " + cachedAlarmsById[uint(alarmRecordId)].Name
						pageWorkplace.AlarmDuration = "[" + time.Now().Sub(cachedAlarmRecords[int(workplace.ID)].DateTimeStart).Round(time.Second).String() + "]"
					} else {
						pageWorkplace.AlarmInformation = getLocale(email, "alarm-name") + ": -"
						pageWorkplace.AlarmDuration = ""
					}
					if breakdownRecordId > 0 {
						pageWorkplace.BreakdownInformation = getLocale(email, "breakdown-name") + ": " + cachedBreakdownsById[uint(breakdownRecordId)].Name
						pageWorkplace.BreakdownDuration = "[" + time.Now().Sub(cachedBreakdownRecords[int(workplace.ID)].DateTimeStart).Round(time.Second).String() + "]"
					} else {
						pageWorkplace.BreakdownInformation = getLocale(email, "breakdown-name") + ": -"
						pageWorkplace.BreakdownDuration = ""
					}
					pageWorkplace.TodayDate = time.Now().Format("02.01.2006")
				default:
					pageWorkplace.WorkplaceColor = "background-color: " + state.Color
					pageWorkplace.WorkplaceState = "mif-stop"
					pageWorkplace.WorkplaceStateName = getLocale(email, "poweroff")
					pageWorkplace.WorkplaceProductivityColor = "bg-darkRed"
					pageWorkplace.WorkplaceStateDuration = time.Since(stateRecord.DateTimeStart).Round(time.Second).String()
					pageWorkplace.WorkplaceProductivityToday = strconv.FormatFloat((totalDuration.Seconds()/totalTodayDuration.Seconds())*100, 'f', 1, 64)
					pageWorkplace.Information = getLocale(email, "order-name") + ": -"
					pageWorkplace.UserInformation = getLocale(email, "user-name") + ": -"
					breakdownRecordId := cachedBreakdownRecords[int(workplace.ID)].BreakdownID
					alarmRecordId := cachedAlarmRecords[int(workplace.ID)].AlarmID
					if alarmRecordId > 0 {
						pageWorkplace.AlarmInformation = getLocale(email, "alarm-name") + ": " + cachedAlarmsById[uint(alarmRecordId)].Name
						pageWorkplace.AlarmDuration = "[" + time.Now().Sub(cachedAlarmRecords[int(workplace.ID)].DateTimeStart).Round(time.Second).String() + "]"
					} else {
						pageWorkplace.AlarmInformation = getLocale(email, "alarm-name") + ": -"
						pageWorkplace.AlarmDuration = ""
					}
					if breakdownRecordId > 0 {
						pageWorkplace.BreakdownInformation = getLocale(email, "breakdown-name") + ": " + cachedBreakdownsById[uint(breakdownRecordId)].Name
						pageWorkplace.BreakdownDuration = "[" + time.Now().Sub(cachedBreakdownRecords[int(workplace.ID)].DateTimeStart).Round(time.Second).String() + "]"
					} else {
						pageWorkplace.BreakdownInformation = getLocale(email, "breakdown-name") + ": -"
						pageWorkplace.BreakdownDuration = ""
					}
					pageWorkplace.TodayDate = time.Now().Format("02.01.2006")
				}
				pageWorkplaces = append(pageWorkplaces, pageWorkplace)
			}
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
	logInfo("MAIN", "Workplace page sent in "+time.Since(timer).String())
}
