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
	Company               string
	Alarms                string
	MenuOverview          string
	MenuWorkplaces        string
	MenuCharts            string
	MenuStatistics        string
	MenuData              string
	MenuSettings          string
	Compacted             string
	UserEmail             string
	UserName              string
	Workplaces            []IndexWorkplaceSelection
	DataFilterPlaceholder string
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
	DowntimesTitle             string
	BreakdownsTitle            string
	AlarmsTitle                string
	CalendarDayLabel           []string
	CalendarMonthLabel         []string
	CalendarData               [][]string
	CalendarStart              string
	CalendarEnd                string
	Locale                     string
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
	loc, err := time.LoadLocation(location)
	workplaceNames, workplacePercents := downloadProductionData(db, loc, email)
	terminalDowntimeNames, terminalDowntimeValues := downloadTerminalDowntimeData(db, email)
	terminalBreakdownNames, terminalBreakdownValues := downloadTerminalBreakdownData(db, email)
	terminalAlarmNames, terminalAlarmValues := downloadTerminalAlarmData(db, email)
	calendarData := downloadCalendarData(db, loc, email)
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
	data.ProductivityTodayTitle = "Productivity Today" // getlocale
	data.DowntimesTitle = getLocale(email, "downtimes")
	data.BreakdownsTitle = getLocale(email, "breakdowns")
	data.AlarmsTitle = getLocale(email, "alarms")
	data.CalendarDayLabel = []string{"Po", "Út", "St", "Čt", "Pá", "So", "Ne"}                                             //get locale
	data.CalendarMonthLabel = []string{"Led", "Úno", "Bře", "Dub", "Kvě", "Čer", "Čvc", "Srp", "Zář", "Říj", "Lis", "Pro"} // get locale
	data.CalendarData = calendarData
	data.CalendarStart = time.Now().In(loc).AddDate(0, -11, 0).Format("2006-01-02")
	data.CalendarEnd = now.EndOfMonth().In(loc).Format("2006-01-02")
	data.Locale = cachedUsersByEmail[email].Locale
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(data)
	logInfo("INDEX", "Index data sent in "+time.Since(timer).String())
}

func downloadCalendarData(db *gorm.DB, loc *time.Location, email string) [][]string {
	stateRecordsAsMap := downloadData(db, time.Now().In(loc).AddDate(-1, 0, 0), time.Now().In(loc), 0, loc, email)
	var calendarData [][]string
	for key, value := range stateRecordsAsMap {
		if key != time.Now().In(loc).Format("2006-01-02") {
			totalDayDuration := 24 * time.Hour
			percentage := ""
			if len(cachedUserSettings[email].selectedWorkplaces) == 0 {
				percentage = strconv.FormatFloat(value.Seconds()/(totalDayDuration.Seconds()*float64(len(cachedWorkplacesById)))*100, 'f', 1, 64)
			} else {
				percentage = strconv.FormatFloat(value.Seconds()/(totalDayDuration.Seconds()*float64(len(cachedUserSettings[email].selectedWorkplaces)))*100, 'f', 1, 64)
			}
			calendarData = append(calendarData, []string{key, percentage})
		} else {
			startOfToday := time.Date(time.Now().In(loc).Year(), time.Now().In(loc).Month(), time.Now().In(loc).Day(), 0, 0, 0, 0, loc)
			totalTodayDuration := time.Now().In(loc).Sub(startOfToday)
			percentage := ""
			if len(cachedUserSettings[email].selectedWorkplaces) == 0 {
				percentage = strconv.FormatFloat(value.Seconds()/(totalTodayDuration.Seconds()*float64(len(cachedWorkplacesById)))*100, 'f', 1, 64)
			} else {
				percentage = strconv.FormatFloat(value.Seconds()/(totalTodayDuration.Seconds()*float64(len(cachedUserSettings[email].selectedWorkplaces)))*100, 'f', 1, 64)
			}
			calendarData = append(calendarData, []string{key, percentage})
		}
	}
	return calendarData
}

func downloadData(db *gorm.DB, fromDate time.Time, toDate time.Time, workplaceId uint, loc *time.Location, email string) map[string]time.Duration {
	stateRecordsAsMap := make(map[string]time.Duration)
	var stateRecords []database.StateRecord
	if workplaceId == 0 {
		var workplaces []database.Workplace
		if len(cachedUserSettings[email].selectedWorkplaces) == 0 {
			for _, workplace := range cachedWorkplacesByName {
				workplaces = append(workplaces, workplace)
			}
		} else {
			for _, workplace := range cachedUserSettings[email].selectedWorkplaces {
				workplaces = append(workplaces, cachedWorkplacesByName[workplace])
			}
		}
		for _, workplace := range workplaces {
			db.Select("state_id, date_time_start").Where("date_time_start >= ?", fromDate).Where("date_time_start <= ?", toDate).Where("workplace_id = ?", workplace.ID).Order("id asc").Find(&stateRecords)
			for index, record := range stateRecords {
				nextDate := time.Now().In(loc)
				if index < len(stateRecords)-1 {
					nextDate = stateRecords[index+1].DateTimeStart
				}
				if record.StateID == 1 {
					if record.DateTimeStart.In(loc).Day() == nextDate.In(loc).Day() {
						stateRecordsAsMap[record.DateTimeStart.In(loc).Format("2006-01-02")] += nextDate.In(loc).Sub(record.DateTimeStart.In(loc))
					} else {
						endOfRecordDay := time.Date(record.DateTimeStart.In(loc).Year(), record.DateTimeStart.In(loc).Month(), record.DateTimeStart.In(loc).Day()+1, 0, 0, 0, 0, loc)
						for record.DateTimeStart.In(loc).Before(nextDate.In(loc)) {
							stateRecordsAsMap[record.DateTimeStart.In(loc).Format("2006-01-02")] += endOfRecordDay.In(loc).Sub(record.DateTimeStart.In(loc))
							record.DateTimeStart = endOfRecordDay.In(loc)
							endOfRecordDay = time.Date(record.DateTimeStart.In(loc).Year(), record.DateTimeStart.In(loc).Month(), record.DateTimeStart.In(loc).Day()+1, 0, 0, 0, 0, loc)
						}
						endOfRecordDay = endOfRecordDay.In(loc).Add(-24 * time.Hour)
						stateRecordsAsMap[nextDate.In(loc).Format("2006-01-02")] += nextDate.In(loc).Sub(endOfRecordDay.In(loc))
					}
				}
			}
		}

	} else {
		db.Select("state_id, date_time_start").Where("date_time_start >= ?", fromDate).Where("date_time_start <= ?", toDate).Where("workplace_id = ?", workplaceId).Order("id asc").Find(&stateRecords)
		for index, record := range stateRecords {
			nextDate := time.Now().In(loc)
			if index < len(stateRecords)-1 {
				nextDate = stateRecords[index+1].DateTimeStart
			}
			if record.StateID == 1 {
				if record.DateTimeStart.In(loc).Day() == nextDate.In(loc).Day() {
					stateRecordsAsMap[record.DateTimeStart.In(loc).Format("2006-01-02")] += nextDate.In(loc).Sub(record.DateTimeStart.In(loc))
				} else {

					endOfRecordDay := time.Date(record.DateTimeStart.In(loc).Year(), record.DateTimeStart.In(loc).Month(), record.DateTimeStart.In(loc).Day()+1, 0, 0, 0, 0, loc)
					for record.DateTimeStart.In(loc).Before(nextDate.In(loc)) {
						stateRecordsAsMap[record.DateTimeStart.In(loc).Format("2006-01-02")] += endOfRecordDay.In(loc).Sub(record.DateTimeStart.In(loc))
						record.DateTimeStart = endOfRecordDay.In(loc)
						endOfRecordDay = time.Date(record.DateTimeStart.In(loc).Year(), record.DateTimeStart.In(loc).Month(), record.DateTimeStart.In(loc).Day()+1, 0, 0, 0, 0, loc)
					}
					endOfRecordDay = endOfRecordDay.In(loc).Add(-24 * time.Hour)
					stateRecordsAsMap[nextDate.In(loc).Format("2006-01-02")] += nextDate.In(loc).Sub(endOfRecordDay.In(loc))
				}
			}
		}
	}
	return stateRecordsAsMap
}

func downloadTerminalAlarmData(db *gorm.DB, email string) ([]string, []float64) {
	var alarmRecords []database.AlarmRecord
	if len(cachedUserSettings[email].selectedWorkplaces) > 0 {
		workplaceIds := `workplace_id in ('`
		for _, workplace := range cachedUserSettings[email].selectedWorkplaces {
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
	if len(cachedUserSettings[email].selectedWorkplaces) > 0 {
		workplaceIds := `workplace_id in ('`
		for _, workplace := range cachedUserSettings[email].selectedWorkplaces {
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
	if len(cachedUserSettings[email].selectedWorkplaces) > 0 {
		workplaceIds := `workplace_id in ('`
		for _, workplace := range cachedUserSettings[email].selectedWorkplaces {
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

func downloadProductionData(db *gorm.DB, loc *time.Location, email string) ([]string, []float64) {
	var workplaceNames []string
	var workplacePercents []float64
	var indexDataWorkplaces []IndexDataWorkplace
	var workplaces []database.Workplace
	if len(cachedUserSettings[email].selectedWorkplaces) == 0 {
		for _, workplace := range cachedWorkplacesByName {
			workplaces = append(workplaces, workplace)
		}
	} else {
		for _, workplace := range cachedUserSettings[email].selectedWorkplaces {
			workplaces = append(workplaces, cachedWorkplacesByName[workplace])
		}
	}
	for _, workplace := range workplaces {
		data := downloadData(db, time.Date(time.Now().UTC().Year(), time.Now().UTC().Month(), time.Now().UTC().Day(), 0, 0, 0, 0, time.Now().Location()), time.Now().In(loc), workplace.ID, loc, email)
		var totalDuration time.Duration
		for _, duration := range data {
			totalDuration = duration
		}
		startOfToday := time.Date(time.Now().In(loc).Year(), time.Now().In(loc).Month(), time.Now().In(loc).Day(), 0, 0, 0, 0, loc)
		totalTodayDuration := time.Now().In(loc).Sub(startOfToday)
		var indexDataWorkplace IndexDataWorkplace
		indexDataWorkplace.Name = workplace.Name
		indexDataWorkplace.Value = (totalDuration.Seconds() / totalTodayDuration.Seconds()) * 100
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
	var dataWorkplaces []IndexWorkplaceSelection
	for _, workplace := range cachedWorkplacesById {
		dataWorkplaces = append(dataWorkplaces, IndexWorkplaceSelection{
			WorkplaceName:      workplace.Name,
			WorkplaceSelection: getWorkplaceSelection(cachedUserSettings[email].selectedWorkplaces, workplace.Name),
		})
	}
	sort.Slice(dataWorkplaces, func(i, j int) bool {
		return dataWorkplaces[i].WorkplaceName < dataWorkplaces[j].WorkplaceName
	})
	data.Workplaces = dataWorkplaces
	data.DataFilterPlaceholder = getLocale(email, "data-table-search-title")
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
