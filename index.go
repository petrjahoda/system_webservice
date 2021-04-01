package main

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/now"
	"github.com/julienschmidt/httprouter"
	"github.com/petrjahoda/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"time"
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
	UserEmail      string
	UserName       string
}

type IndexData struct {
	WorkplaceNames             []string
	WorkplacePercents          []float64
	TerminalProductionColor    string
	TerminalDowntimeNames      []string
	TerminalDowntimeDurations  []float64
	TerminalDowntimeColor      string
	TerminalBreakdownNames     []string
	TerminalBreakdownDurations []float64
	TerminalBreakdownColor     string
	TerminalAlarmNames         []string
	TerminalAlarmDurations     []float64
	TerminalAlarmColor         string
	ProductivityTodayTitle     string
	DowntimesTitle             string
	BreakdownsTitle            string
	AlarmsTitle                string
	CalendarDayLabel           []string
	CalendarMonthLabel         []string
	CalendarData               [][]string
	CalendarStart              string
	CalendarEnd                string
}

type IndexDataWorkplace struct {
	Name  string
	Value float64
}

func loadIndexData(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	timer := time.Now()
	email, _, _ := request.BasicAuth()
	logInfo("INDEX", "Loading index data for "+cachedUsersByEmail[email].FirstName+" "+cachedUsersByEmail[email].SecondName)
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("DATA-ALARMS", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("DATA-ALARMS", "Loading alarms table ended")
		return
	}
	data := IndexData{}
	workplaceNames, workplacePercents := downloadProductionData(db)
	terminalDowntimeNames, terminalDowntimeValues := downloadTerminalDowntimeData(db)
	terminalBreakdownNames, terminalBreakdownValues := downloadTerminalBreakdownData(db)
	terminalAlarmNames, terminalAlarmValues := downloadTerminalAlarmData(db)
	calendarData := downloadCalendarData(db)
	loc, err := time.LoadLocation(location)
	data.WorkplaceNames = workplaceNames
	data.WorkplacePercents = workplacePercents
	data.TerminalDowntimeNames = terminalDowntimeNames
	data.TerminalDowntimeDurations = terminalDowntimeValues
	data.TerminalBreakdownNames = terminalBreakdownNames
	data.TerminalBreakdownDurations = terminalBreakdownValues
	data.TerminalAlarmNames = terminalAlarmNames
	data.TerminalAlarmDurations = terminalAlarmValues
	data.TerminalDowntimeColor = cachedStatesById[2].Color
	data.TerminalProductionColor = cachedStatesById[1].Color
	data.TerminalBreakdownColor = cachedStatesById[3].Color
	data.TerminalAlarmColor = "grey"
	data.ProductivityTodayTitle = "Productivity Today"
	data.DowntimesTitle = "Downtimes"
	data.BreakdownsTitle = "Breakdowns"
	data.AlarmsTitle = "Alarms"
	data.CalendarDayLabel = []string{"Po", "Út", "St", "Čt", "Pá", "So", "Ne"}
	data.CalendarMonthLabel = []string{"Led", "Úno", "Bře", "Dub", "Kvě", "Čer", "Čvc", "Srp", "Zář", "Říj", "Lis", "Pro"}
	data.CalendarData = calendarData
	data.CalendarStart = time.Now().In(loc).AddDate(0, -11, 0).Format("2006-01-02")
	data.CalendarEnd = now.EndOfMonth().In(loc).Format("2006-01-02")
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(data)
	logInfo("INDEX", "Index data sent in "+time.Since(timer).String())
}

func downloadCalendarData(db *gorm.DB) [][]string {
	var stateRecords []database.StateRecord
	db.Where("date_time_start >= ?", time.Now().AddDate(-1, 0, 0)).Where("state_id = 1").Order("date_time_start asc").Find(&stateRecords)
	stateRecordsAsMap := make(map[string]float64)
	for _, record := range stateRecords {
		if record.DateTimeStart.Day() == record.DateTimeEnd.Time.Day() {
			if record.DateTimeEnd.Time.IsZero() {
				record.DateTimeEnd.Time = time.Now().UTC()
			}
			stateRecordsAsMap[record.DateTimeStart.Format("2006-01-02")] += record.DateTimeEnd.Time.Sub(record.DateTimeStart).Seconds()
		} else {
			if record.DateTimeEnd.Time.IsZero() {
				record.DateTimeEnd.Time = time.Now().UTC()
			}
			actualDay := record.DateTimeStart
			endOfActualDay := time.Date(record.DateTimeStart.Year(), record.DateTimeStart.Month(), record.DateTimeStart.Day()+1, 0, 0, 0, 0, time.UTC)
			for actualDay.Day() < record.DateTimeEnd.Time.Day() {
				stateRecordsAsMap[record.DateTimeStart.Format("2006-01-02")] += endOfActualDay.Sub(actualDay).Seconds()
				actualDay = endOfActualDay
				endOfActualDay = endOfActualDay.Add(24 * time.Hour)
			}
			stateRecordsAsMap[record.DateTimeStart.Format("2006-01-02")] += record.DateTimeEnd.Time.Sub(endOfActualDay).Seconds()
		}
	}
	var calendarData [][]string
	for key, value := range stateRecordsAsMap {
		percentage := strconv.FormatFloat(value/86400, 'f', 1, 64)
		calendarData = append(calendarData, []string{key, percentage})
	}
	return calendarData
}

func downloadTerminalAlarmData(db *gorm.DB) ([]string, []float64) {
	var alarmRecords []database.AlarmRecord
	db.Where("date_time_end is null").Limit(10).Find(&alarmRecords)

	var indexDataWorkplaces []IndexDataWorkplace
	for _, alarmRecord := range alarmRecords {
		var indexDataWorkplace IndexDataWorkplace
		indexDataWorkplace.Name = cachedWorkplacesById[uint(alarmRecord.WorkplaceID)].Name + ": " + cachedAlarmsById[uint(alarmRecord.AlarmID)].Name
		indexDataWorkplace.Value = time.Since(alarmRecord.DateTimeStart).Seconds()
		indexDataWorkplaces = append(indexDataWorkplaces, indexDataWorkplace)
	}
	sort.Slice(indexDataWorkplaces, func(i, j int) bool {
		return indexDataWorkplaces[i].Value < indexDataWorkplaces[j].Value
	})
	var terminalAlarmNames []string
	var terminalAlarmValues []float64
	for _, workplace := range indexDataWorkplaces {
		percentage := time.Duration(workplace.Value) * time.Second
		terminalAlarmNames = append(terminalAlarmNames, workplace.Name+": "+percentage.String())
		terminalAlarmValues = append(terminalAlarmValues, workplace.Value)
	}
	return terminalAlarmNames, terminalAlarmValues
}

func downloadTerminalBreakdownData(db *gorm.DB) ([]string, []float64) {
	var breakdownRecords []database.BreakdownRecord
	db.Where("date_time_end is null").Limit(10).Find(&breakdownRecords)

	var indexDataWorkplaces []IndexDataWorkplace
	for _, breakdownRecord := range breakdownRecords {
		var indexDataWorkplace IndexDataWorkplace
		indexDataWorkplace.Name = cachedWorkplacesById[uint(breakdownRecord.WorkplaceID)].Name + ": " + cachedBreakdownsById[uint(breakdownRecord.BreakdownID)].Name
		indexDataWorkplace.Value = time.Since(breakdownRecord.DateTimeStart).Seconds()
		indexDataWorkplaces = append(indexDataWorkplaces, indexDataWorkplace)
	}
	sort.Slice(indexDataWorkplaces, func(i, j int) bool {
		return indexDataWorkplaces[i].Value < indexDataWorkplaces[j].Value
	})
	var terminalBreakdownNames []string
	var terminalBreakdownValues []float64
	for _, workplace := range indexDataWorkplaces {
		percentage := time.Duration(workplace.Value) * time.Second
		terminalBreakdownNames = append(terminalBreakdownNames, workplace.Name+": "+percentage.String())
		terminalBreakdownValues = append(terminalBreakdownValues, workplace.Value)
	}
	return terminalBreakdownNames, terminalBreakdownValues
}

func downloadTerminalDowntimeData(db *gorm.DB) ([]string, []float64) {
	var downtimeRecords []database.DowntimeRecord
	db.Where("date_time_end is null").Limit(10).Find(&downtimeRecords)

	var indexDataWorkplaces []IndexDataWorkplace
	for _, downtimeRecord := range downtimeRecords {
		var indexDataWorkplace IndexDataWorkplace
		indexDataWorkplace.Name = cachedWorkplacesById[uint(downtimeRecord.WorkplaceID)].Name + ": " + cachedDowntimesById[uint(downtimeRecord.DowntimeID)].Name
		indexDataWorkplace.Value = time.Since(downtimeRecord.DateTimeStart).Seconds()
		indexDataWorkplaces = append(indexDataWorkplaces, indexDataWorkplace)
	}
	sort.Slice(indexDataWorkplaces, func(i, j int) bool {
		return indexDataWorkplaces[i].Value < indexDataWorkplaces[j].Value
	})
	var terminalDowntimeNames []string
	var terminalDowntimeValues []float64
	for _, workplace := range indexDataWorkplaces {
		percentage := time.Duration(workplace.Value) * time.Second
		terminalDowntimeNames = append(terminalDowntimeNames, workplace.Name+": "+percentage.String())
		terminalDowntimeValues = append(terminalDowntimeValues, workplace.Value)
	}
	return terminalDowntimeNames, terminalDowntimeValues
}

func downloadProductionData(db *gorm.DB) ([]string, []float64) {
	var workplaceNames []string
	var workplacePercents []float64
	var indexDataWorkplaces []IndexDataWorkplace
	for _, workplace := range cachedWorkplacesByName {
		productivity := calculateProductivity(db, workplace)
		var indexDataWorkplace IndexDataWorkplace
		indexDataWorkplace.Name = workplace.Name
		indexDataWorkplace.Value, _ = strconv.ParseFloat(productivity, 32)
		indexDataWorkplaces = append(indexDataWorkplaces, indexDataWorkplace)
	}
	sort.Slice(indexDataWorkplaces, func(i, j int) bool {
		return indexDataWorkplaces[i].Value < indexDataWorkplaces[j].Value
	})
	for _, workplace := range indexDataWorkplaces {
		percentage := fmt.Sprintf("%.1f", workplace.Value)
		workplaceNames = append(workplaceNames, workplace.Name+": "+percentage+"%")
		workplacePercents = append(workplacePercents, workplace.Value)
	}
	return workplaceNames, workplacePercents
}

func index(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	timer := time.Now()
	email, _, _ := request.BasicAuth()
	logInfo("INDEX", "Sending home page to "+cachedUsersByEmail[email].FirstName+" "+cachedUsersByEmail[email].SecondName)
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
	data.UserEmail = email
	data.UserName = cachedUsersByEmail[email].FirstName + " " + cachedUsersByEmail[email].SecondName
	tmpl := template.Must(template.ParseFiles("./html/index.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("INDEX", "Index page sent in "+time.Since(timer).String())
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
