package main

import (
	"encoding/json"
	"github.com/jinzhu/now"
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

type IndexPageData struct {
	Version               string
	Information           string
	Company               string
	Alarms                string
	MenuOverview          string
	MenuWorkplaces        string
	MenuCharts            string
	MenuStatistics        string
	MenuData              string
	MenuSettings          string
	UserEmail             string
	UserName              string
	Workplaces            []IndexWorkplaceSelection
	DataFilterPlaceholder string
	Software              string
}

type IndexWorkplaceSelection struct {
	WorkplaceName      string
	WorkplaceSelection string
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
	ProductivityYearTitle      string
	OverviewMonthTitle         string
	ConsumptionMonthTitle      string
	DowntimesTitle             string
	BreakdownsTitle            string
	AlarmsTitle                string
	CalendarDayLabel           []string
	CalendarMonthLabel         []string
	CalendarData               [][]string
	ConsumptionData            []string
	MonthDataDays              []string
	MonthDataProduction        []string
	MonthDataDowntime          []string
	MonthDataPoweroff          []string
	CalendarStart              string
	CalendarEnd                string
	Locale                     string
	ProductionLocale           string
	DowntimeLocale             string
	PoweroffLocale             string
	Result                     string
}

type IndexDataWorkplace struct {
	Name  string
	Value float64
}

type IndexPageInput struct {
	Email string
}

func loadIndexData(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	timer := time.Now()
	email, _, _ := request.BasicAuth()
	if len(email) == 0 {
		var data IndexPageInput
		err := json.NewDecoder(request.Body).Decode(&data)
		if err != nil {
			logError("SETTINGS", "Problem parsing email: "+err.Error())
			var responseData TableOutput
			responseData.Result = "ERR: Problem parsing email, " + err.Error()
			writer.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(writer).Encode(responseData)
			logInfo("SETTINGS", "Loading downtime ended with error")
			return
		}
		email = data.Email
	}
	logInfo("INDEX", "Loading index data for "+cachedUsersByEmail[email].FirstName+" "+cachedUsersByEmail[email].SecondName)
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("INDEX", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "ERR: Problem opening database, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("INDEX", "Loading alarms table ended")
		return
	}
	data := IndexData{}
	companyNameSync.Lock()
	loc, err := time.LoadLocation(location)
	companyNameSync.Unlock()
	cacheProductionData(db, time.Date(latestCachedWorkplaceCalendarData.Year(), latestCachedWorkplaceCalendarData.Month(), latestCachedWorkplaceCalendarData.Day(), 0, 0, 0, 0, loc))
	cacheConsumptionData(db, time.Date(latestCachedWorkplaceConsumption.Year(), latestCachedWorkplaceConsumption.Month(), latestCachedWorkplaceConsumption.Day(), 0, 0, 0, 0, loc))
	workplaceNames, workplacePercents := downloadProductionData(loc, email)
	terminalDowntimeNames, terminalDowntimeValues := downloadTerminalDowntimeData(db, email)
	terminalBreakdownNames, terminalBreakdownValues := downloadTerminalBreakdownData(db, email)
	terminalAlarmNames, terminalAlarmValues := downloadTerminalAlarmData(db, email)
	calendarData := downloadCalendarData(email, loc)
	monthDataDays, monthDataProduction, monthDataDowntime, monthDataPoweroff := downloadIndexChartData(email, loc)
	consumptionData := downloadConsumptionData(email)
	timeBack := time.Now().In(loc).AddDate(0, -11, 0)
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
	data.ProductivityTodayTitle = getLocale(email, "production-today")
	data.ProductivityYearTitle = getLocale(email, "production-last-year")
	data.OverviewMonthTitle = getLocale(email, "overview-last-month")
	data.DowntimesTitle = getLocale(email, "downtimes")
	data.BreakdownsTitle = getLocale(email, "breakdowns")
	data.AlarmsTitle = getLocale(email, "alarms")
	data.CalendarDayLabel = strings.Split(getLocale(email, "day-names"), ",")
	data.CalendarMonthLabel = strings.Split(getLocale(email, "month-names"), ",")
	data.CalendarData = calendarData
	data.ConsumptionData = consumptionData
	data.MonthDataDays = monthDataDays
	data.MonthDataProduction = monthDataProduction
	data.MonthDataDowntime = monthDataDowntime
	data.MonthDataPoweroff = monthDataPoweroff
	data.CalendarStart = time.Date(timeBack.Year(), timeBack.Month(), 1, 0, 0, 0, 0, loc).Format("2006-01-02")
	data.CalendarEnd = now.EndOfMonth().Format("2006-01-02")
	data.Locale = cachedUsersByEmail[email].Locale
	data.ProductionLocale = getLocale(email, "production")
	data.DowntimeLocale = getLocale(email, "downtime")
	data.PoweroffLocale = getLocale(email, "poweroff")
	data.OverviewMonthTitle = getLocale(email, "month-overview")
	data.ConsumptionMonthTitle = getLocale(email, "month-consumption-overview")
	data.Result = "INF: Data processed in " + time.Since(timer).String()
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(data)
	logInfo("INDEX", "Index data sent in "+time.Since(timer).String())
}

func downloadIndexChartData(email string, loc *time.Location) ([]string, []string, []string, []string) {
	var monthDataDays []string
	var monthDataProduction []string
	var monthDataDowntime []string
	var monthDataPoweroff []string
	var workplaces []database.Workplace
	if len(cachedUserWebSettings[email]["index-selected-workplaces"]) == 0 {
		for _, workplace := range cachedWorkplacesByName {
			workplaces = append(workplaces, workplace)
		}
	} else {
		for _, workplace := range strings.Split(cachedUserWebSettings[email]["index-selected-workplaces"], ";") {
			workplaces = append(workplaces, cachedWorkplacesByName[workplace])
		}
	}
	var cachedProductionData = map[string]time.Duration{}
	var cachedDowntimeData = map[string]time.Duration{}
	var cachedPoweroffData = map[string]time.Duration{}
	for _, workplace := range workplaces {
		workplacesRecords.Lock()
		productionData := cachedWorkplacesProductionRecords[workplace.Name]
		workplacesRecords.Unlock()
		for date, duration := range productionData {
			cachedProductionData[date] = cachedProductionData[date] + duration
		}
		workplacesRecords.Lock()
		downtimeData := cachedWorkplacesDowntimeRecords[workplace.Name]
		workplacesRecords.Unlock()
		for date, duration := range downtimeData {
			cachedDowntimeData[date] = cachedDowntimeData[date] + duration
		}
		workplacesRecords.Lock()
		poweroffData := cachedWorkplacesPoweroffRecords[workplace.Name]
		workplacesRecords.Unlock()
		for date, duration := range poweroffData {
			cachedPoweroffData[date] = cachedPoweroffData[date] + duration
		}
	}
	todayNoon := time.Date(time.Now().UTC().Year(), time.Now().UTC().Month(), time.Now().UTC().Day(), 0, 0, 0, 0, loc).In(loc)
	for i := 30; i >= 0; i-- {
		if i == 0 {
			totalTodayDuration := time.Now().In(loc).Sub(todayNoon).Seconds()
			dayWorkplaceDuration := time.Duration(len(workplaces)*int(totalTodayDuration)) * time.Second
			day := time.Now().Add(time.Duration(-24*i) * time.Hour)
			monthDataDays = append(monthDataDays, day.Format("2006-01-02"))
			productionPercentage := (cachedProductionData[day.Format("2006-01-02")].Seconds() / dayWorkplaceDuration.Seconds()) * 100
			productionPercentageAsString := strconv.FormatFloat(productionPercentage, 'f', 1, 64)
			monthDataProduction = append(monthDataProduction, productionPercentageAsString)
			downtimePercentage := (cachedDowntimeData[day.Format("2006-01-02")].Seconds() / dayWorkplaceDuration.Seconds()) * 100
			downtimePercentageAsString := strconv.FormatFloat(downtimePercentage, 'f', 1, 64)
			monthDataDowntime = append(monthDataDowntime, downtimePercentageAsString)
			poweroffPercentage := strconv.FormatFloat(100.0-productionPercentage-downtimePercentage, 'f', 1, 64)
			monthDataPoweroff = append(monthDataPoweroff, poweroffPercentage)
		} else {
			day := time.Now().Add(time.Duration(-24*i) * time.Hour)
			totalTodayDuration := time.Duration(len(workplaces)*86400) * time.Second
			monthDataDays = append(monthDataDays, day.Format("2006-01-02"))
			productionPercentage := (cachedProductionData[day.Format("2006-01-02")].Seconds() / totalTodayDuration.Seconds()) * 100
			productionPercentageAsString := strconv.FormatFloat(productionPercentage, 'f', 1, 64)
			monthDataProduction = append(monthDataProduction, productionPercentageAsString)
			downtimePercentage := (cachedDowntimeData[day.Format("2006-01-02")].Seconds() / totalTodayDuration.Seconds()) * 100
			downtimePercentageAsString := strconv.FormatFloat(downtimePercentage, 'f', 1, 64)
			monthDataDowntime = append(monthDataDowntime, downtimePercentageAsString)
			poweroffPercentage := strconv.FormatFloat(100.0-productionPercentage-downtimePercentage, 'f', 1, 64)
			monthDataPoweroff = append(monthDataPoweroff, poweroffPercentage)
		}
	}
	return monthDataDays, monthDataProduction, monthDataDowntime, monthDataPoweroff
}

func downloadCalendarData(email string, loc *time.Location) [][]string {
	var workplaces []database.Workplace
	if len(cachedUserWebSettings[email]["index-selected-workplaces"]) == 0 {
		for _, workplace := range cachedWorkplacesByName {
			workplaces = append(workplaces, workplace)
		}
	} else {
		for _, workplace := range strings.Split(cachedUserWebSettings[email]["index-selected-workplaces"], ";") {
			workplaces = append(workplaces, cachedWorkplacesByName[workplace])
		}
	}
	var calendarData [][]string
	todayNoon := time.Date(time.Now().UTC().Year(), time.Now().UTC().Month(), time.Now().UTC().Day(), 0, 0, 0, 0, loc).In(loc)
	totalTodayDuration := time.Now().In(loc).Sub(todayNoon).Seconds()
	var cachedCalendarData = map[string]time.Duration{}
	for _, workplace := range workplaces {
		workplacesRecords.Lock()
		data := cachedWorkplacesProductionRecords[workplace.Name]
		workplacesRecords.Unlock()
		for date, duration := range data {
			cachedCalendarData[date] = cachedCalendarData[date] + duration
		}
	}
	for date, duration := range cachedCalendarData {
		if date == todayNoon.Format("2006-01-02") {
			dayWorkplaceDuration := time.Duration(len(workplaces)*int(totalTodayDuration)) * time.Second
			percentage := strconv.FormatFloat((duration.Seconds()/dayWorkplaceDuration.Seconds())*100, 'f', 1, 64)
			calendarData = append(calendarData, []string{date, percentage})
		} else {
			dayWorkplaceDuration := time.Duration(len(workplaces)*86400) * time.Second
			percentage := strconv.FormatFloat((duration.Seconds()/dayWorkplaceDuration.Seconds())*100, 'f', 1, 64)
			calendarData = append(calendarData, []string{date, percentage})
		}
	}
	return calendarData
}

func downloadConsumptionData(email string) []string {
	var workplaces []database.Workplace
	if len(cachedUserWebSettings[email]["index-selected-workplaces"]) == 0 {
		for _, workplace := range cachedWorkplacesByName {
			workplaces = append(workplaces, workplace)
		}
	} else {
		for _, workplace := range strings.Split(cachedUserWebSettings[email]["index-selected-workplaces"], ";") {
			workplaces = append(workplaces, cachedWorkplacesByName[workplace])
		}
	}
	var cachedConsumptionData = map[string]float32{}
	for _, workplace := range workplaces {
		consumptionDataSync.Lock()
		data := cachedConsumptionDataByWorkplaceName[workplace.Name]
		consumptionDataSync.Unlock()
		for date, value := range data {
			cachedConsumptionData[date] = cachedConsumptionData[date] + value
		}
	}

	var consumptionData []string
	for i := 30; i >= 0; i-- {
		day := time.Now().Add(time.Duration(-24*i) * time.Hour)
		consumptionData = append(consumptionData, strconv.FormatFloat(float64(cachedConsumptionData[day.Format("2006-01-02")]), 'f', 1, 64))
	}
	return consumptionData
}

func downloadTerminalAlarmData(db *gorm.DB, email string) ([]string, []float64) {
	var alarmRecords []database.AlarmRecord
	if len(cachedUserWebSettings[email]["index-selected-workplaces"]) > 0 {
		workplaceIds := `workplace_id in ('`
		for _, workplace := range strings.Split(cachedUserWebSettings[email]["index-selected-workplaces"], ";") {
			workplaceIds += strconv.Itoa(int(cachedWorkplacesByName[workplace].ID)) + `','`
		}
		workplaceIds = strings.TrimSuffix(workplaceIds, `,'`)
		workplaceIds += ")"
		db.Where("date_time_end is null").Where(workplaceIds).Find(&alarmRecords)
	} else {
		db.Where("date_time_end is null").Find(&alarmRecords)
	}
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
		terminalAlarmNames = append(terminalAlarmNames, workplace.Name)
		terminalAlarmValues = append(terminalAlarmValues, workplace.Value)
	}
	return terminalAlarmNames, terminalAlarmValues
}

func downloadTerminalBreakdownData(db *gorm.DB, email string) ([]string, []float64) {
	var breakdownRecords []database.BreakdownRecord
	if len(cachedUserWebSettings[email]["index-selected-workplaces"]) > 0 {
		workplaceIds := `workplace_id in ('`
		for _, workplace := range strings.Split(cachedUserWebSettings[email]["index-selected-workplaces"], ";") {
			workplaceIds += strconv.Itoa(int(cachedWorkplacesByName[workplace].ID)) + `','`
		}
		workplaceIds = strings.TrimSuffix(workplaceIds, `,'`)
		workplaceIds += ")"
		db.Where("date_time_end is null").Where(workplaceIds).Find(&breakdownRecords)
	} else {
		db.Where("date_time_end is null").Find(&breakdownRecords)
	}
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
		terminalBreakdownNames = append(terminalBreakdownNames, workplace.Name)
		terminalBreakdownValues = append(terminalBreakdownValues, workplace.Value)
	}
	return terminalBreakdownNames, terminalBreakdownValues
}

func downloadTerminalDowntimeData(db *gorm.DB, email string) ([]string, []float64) {
	var downtimeRecords []database.DowntimeRecord
	if len(cachedUserWebSettings[email]["index-selected-workplaces"]) > 0 {
		workplaceIds := `workplace_id in ('`
		for _, workplace := range strings.Split(cachedUserWebSettings[email]["index-selected-workplaces"], ";") {
			workplaceIds += strconv.Itoa(int(cachedWorkplacesByName[workplace].ID)) + `','`
		}
		workplaceIds = strings.TrimSuffix(workplaceIds, `,'`)
		workplaceIds += ")"
		db.Where("date_time_end is null").Where(workplaceIds).Find(&downtimeRecords)
	} else {
		db.Where("date_time_end is null").Find(&downtimeRecords)
	}

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
		terminalDowntimeNames = append(terminalDowntimeNames, workplace.Name)
		terminalDowntimeValues = append(terminalDowntimeValues, workplace.Value)
	}
	return terminalDowntimeNames, terminalDowntimeValues
}

func downloadProductionData(loc *time.Location, email string) ([]string, []float64) {
	var workplaces []database.Workplace
	if len(cachedUserWebSettings[email]["index-selected-workplaces"]) == 0 {
		for _, workplace := range cachedWorkplacesByName {
			workplaces = append(workplaces, workplace)
		}
	} else {
		for _, workplace := range strings.Split(cachedUserWebSettings[email]["index-selected-workplaces"], ";") {
			workplaces = append(workplaces, cachedWorkplacesByName[workplace])
		}
	}
	var workplaceNames []string
	var workplacePercents []float64
	var indexDataWorkplaces []IndexDataWorkplace
	todayNoon := time.Date(time.Now().UTC().Year(), time.Now().UTC().Month(), time.Now().UTC().Day(), 0, 0, 0, 0, loc).In(loc)
	totalTodayDuration := time.Now().In(loc).Sub(todayNoon)
	for _, workplace := range workplaces {
		workplacesRecords.Lock()
		data := cachedWorkplacesProductionRecords[workplace.Name]
		workplacesRecords.Unlock()
		var indexDataWorkplace IndexDataWorkplace
		indexDataWorkplace.Name = workplace.Name
		indexDataWorkplace.Value = (data[todayNoon.Format("2006-01-02")].Seconds() / totalTodayDuration.Seconds()) * 100
		indexDataWorkplaces = append(indexDataWorkplaces, indexDataWorkplace)
	}

	sort.Slice(indexDataWorkplaces, func(i, j int) bool {
		return indexDataWorkplaces[i].Value < indexDataWorkplaces[j].Value
	})
	for _, workplace := range indexDataWorkplaces {
		workplaceNames = append(workplaceNames, workplace.Name)
		workplacePercents = append(workplacePercents, workplace.Value)
	}
	return workplaceNames, workplacePercents
}

func index(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	timer := time.Now()
	go updatePageCount("index")
	email, _, _ := request.BasicAuth()
	logInfo("INDEX", "Sending home page to "+cachedUsersByEmail[email].FirstName+" "+cachedUsersByEmail[email].SecondName)
	var data IndexPageData
	data.Version = version
	companyNameSync.Lock()
	data.Company = cachedCompanyName
	companyNameSync.Unlock()
	data.MenuOverview = getLocale(email, "menu-overview")
	data.MenuWorkplaces = getLocale(email, "menu-workplaces")
	data.MenuCharts = getLocale(email, "menu-charts")
	data.MenuStatistics = getLocale(email, "menu-statistics")
	data.MenuData = getLocale(email, "menu-data")
	data.MenuSettings = getLocale(email, "menu-settings")
	data.UserEmail = email
	data.UserName = cachedUsersByEmail[email].FirstName + " " + cachedUsersByEmail[email].SecondName
	var dataWorkplaces []IndexWorkplaceSelection
	for _, workplace := range cachedWorkplacesById {
		dataWorkplaces = append(dataWorkplaces, IndexWorkplaceSelection{
			WorkplaceName:      workplace.Name,
			WorkplaceSelection: getWorkplaceWebSelection(cachedUserWebSettings[email]["index-selected-workplaces"], workplace.Name),
		})
	}
	sort.Slice(dataWorkplaces, func(i, j int) bool {
		return dataWorkplaces[i].WorkplaceName < dataWorkplaces[j].WorkplaceName
	})
	data.Workplaces = dataWorkplaces
	data.DataFilterPlaceholder = getLocale(email, "data-table-search-title")
	data.Software = cachedSoftwareName
	data.Information = "INF: Page processed in " + time.Since(timer).String()
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
