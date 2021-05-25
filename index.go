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

func loadIndexData(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	timer := time.Now()
	email, _, _ := request.BasicAuth()
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
	loc, err := time.LoadLocation(location)
	workplaceNames, workplacePercents := downloadProductionData(db, loc, email)
	terminalDowntimeNames, terminalDowntimeValues := downloadTerminalDowntimeData(db, email)
	terminalBreakdownNames, terminalBreakdownValues := downloadTerminalBreakdownData(db, email)
	terminalAlarmNames, terminalAlarmValues := downloadTerminalAlarmData(db, email)
	productionData, downtimeData, poweroffData := downloadIndexData(db, loc, email)
	monthDataDays, monthDataProduction, monthDataDowntime, monthDataPoweroff := processProductionData(productionData, loc, downtimeData, poweroffData)
	consumptionData := downloadConsumptionData(db, loc, email)
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
	data.CalendarData = productionData
	data.ConsumptionData = consumptionData
	data.MonthDataDays = monthDataDays
	data.MonthDataProduction = monthDataProduction
	data.MonthDataDowntime = monthDataDowntime
	data.MonthDataPoweroff = monthDataPoweroff
	data.CalendarStart = time.Date(timeBack.Year(), timeBack.Month(), 1, 0, 0, 0, 0, loc).Format("2006-01-02")
	data.CalendarEnd = now.EndOfMonth().In(loc).Format("2006-01-02")
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

func processProductionData(productionData [][]string, loc *time.Location, downtimeData [][]string, poweroffData [][]string) ([]string, []string, []string, []string) {
	var monthDataDays []string
	var monthDataProduction []string
	var monthDataDowntime []string
	var monthDataPoweroff []string
	for index, data := range productionData {
		layout := "2006-01-02"
		date, err := time.Parse(layout, data[0])
		if err != nil {
			logError("INDEX", "Problem parsing date "+err.Error()+" "+data[0])
		}
		if time.Now().In(loc).Sub(date).Hours() < 720 {
			monthDataDays = append(monthDataDays, data[0])
			monthDataProduction = append(monthDataProduction, data[1])
			monthDataDowntime = append(monthDataDowntime, downtimeData[index][1])
			monthDataPoweroff = append(monthDataPoweroff, poweroffData[index][1])
		}
	}
	return monthDataDays, monthDataProduction, monthDataDowntime, monthDataPoweroff
}

func downloadConsumptionData(db *gorm.DB, loc *time.Location, email string) []string {
	consumptionDataSync.Lock()
	todayNoon := time.Date(consumptionDataLastFullDateTime.UTC().Year(), consumptionDataLastFullDateTime.UTC().Month(), consumptionDataLastFullDateTime.UTC().Day(), 0, 0, 0, 0, loc).In(loc)
	for _, workplace := range cachedWorkplacesById {
		tempConsumptionData := cachedConsumptionDataByWorkplaceName[workplace.Name]
		for _, port := range cachedWorkplacePorts[workplace.Name] {
			if port.StateID.Int32 == 3 {
				var devicePortAnalogRecords []database.DevicePortAnalogRecord
				db.Where("device_port_id = ?", port.DevicePortID).Where("date_time >= ?", todayNoon).Find(&devicePortAnalogRecords)
				tempConsumptionData[todayNoon.Format("2006-01-02")] = 0
				for _, record := range devicePortAnalogRecords {
					tempConsumptionData[record.DateTime.Format("2006-01-02")] += record.Data
				}
				cachedConsumptionDataByWorkplaceName[workplace.Name] = tempConsumptionData
			}
		}
	}
	consumptionDataLastFullDateTime = time.Now().In(loc)
	consumptionDataSync.Unlock()
	cachedConsumptionData := make(map[string]float32)
	initialDate := time.Now().In(loc).AddDate(0, -1, 0)
	for initialDate.Format("2006-01-02") != time.Now().In(loc).Format("2006-01-02") {
		cachedConsumptionData[initialDate.Format("2006-01-02")] = 0.0
		initialDate = initialDate.Add(24 * time.Hour)
	}
	if len(cachedUserWebSettings[email]["index-selected-workplaces"]) > 0 {
		for _, workplace := range strings.Split(cachedUserWebSettings[email]["index-selected-workplaces"], ";") {
			for date, consumption := range cachedConsumptionDataByWorkplaceName[workplace] {
				cachedConsumptionData[date] += consumption
			}
		}
	} else {
		for _, data := range cachedConsumptionDataByWorkplaceName {
			for date, consumption := range data {
				cachedConsumptionData[date] += consumption
			}
		}
	}
	var temporaryData [][]string
	for key, value := range cachedConsumptionData {
		dateTimeParsed, _ := time.Parse("2006-01-02", key)
		if time.Now().In(loc).Sub(dateTimeParsed).Hours() <= 720 {
			consumption := strconv.FormatFloat(float64(value*230/1000), 'f', 1, 64)
			temporaryData = append(temporaryData, []string{key, consumption})
		}
	}
	sort.Slice(temporaryData[:], func(i, j int) bool {
		return temporaryData[i][0] < temporaryData[j][0]
	})
	var consumptionData []string
	for _, consumption := range temporaryData {
		consumptionData = append(consumptionData, consumption[1])
	}
	return consumptionData
}

func downloadIndexData(db *gorm.DB, loc *time.Location, email string) ([][]string, [][]string, [][]string) {
	productionRecords, downtimeRecords, poweroffRecords := downloadYearData(db, time.Now().In(loc).AddDate(-1, 0, 0), time.Now().In(loc), loc, email)
	var productionData [][]string
	var downtimeData [][]string
	var poweroffData [][]string
	for key, value := range productionRecords {
		productionPercentage := strconv.FormatFloat(value.Seconds()*100/(value.Seconds()+downtimeRecords[key].Seconds()+poweroffRecords[key].Seconds()), 'f', 1, 64)
		downtimePercentage := strconv.FormatFloat(downtimeRecords[key].Seconds()*100/(value.Seconds()+downtimeRecords[key].Seconds()+poweroffRecords[key].Seconds()), 'f', 1, 64)
		poweroffPercentage := strconv.FormatFloat(poweroffRecords[key].Seconds()*100/(value.Seconds()+downtimeRecords[key].Seconds()+poweroffRecords[key].Seconds()), 'f', 1, 64)
		productionData = append(productionData, []string{key, productionPercentage})
		downtimeData = append(downtimeData, []string{key, downtimePercentage})
		poweroffData = append(poweroffData, []string{key, poweroffPercentage})
	}
	sort.Slice(productionData[:], func(i, j int) bool {
		return productionData[i][0] < productionData[j][0]
	})
	sort.Slice(downtimeData[:], func(i, j int) bool {
		return downtimeData[i][0] < downtimeData[j][0]
	})
	sort.Slice(poweroffData[:], func(i, j int) bool {
		return poweroffData[i][0] < poweroffData[j][0]
	})
	layout := "2006-01-02"
	initialDate, err := time.Parse(layout, productionData[0][0])
	if err != nil {
		logError("INDEX", "Problem parsing date "+err.Error()+" "+productionData[0][0])
	}
	var datesToAdd []string
	for _, data := range productionData {
		actualDate, err := time.Parse(layout, data[0])
		if err != nil {
			logError("INDEX", "Problem parsing date "+err.Error()+" "+data[0])
		}
		if actualDate != initialDate {
			for actualDate != initialDate {
				datesToAdd = append(datesToAdd, initialDate.Format("2006-01-02"))
				initialDate = initialDate.Add(24 * time.Hour)
			}
			initialDate = initialDate.Add(24 * time.Hour)
		} else {
			initialDate = initialDate.Add(24 * time.Hour)
		}
	}
	for _, dateToAdd := range datesToAdd {
		productionData = append(productionData, []string{dateToAdd, "0.0"})
		downtimeData = append(downtimeData, []string{dateToAdd, "0.0"})
		poweroffData = append(poweroffData, []string{dateToAdd, "100.0"})
	}
	sort.Slice(productionData[:], func(i, j int) bool {
		return productionData[i][0] < productionData[j][0]
	})
	sort.Slice(downtimeData[:], func(i, j int) bool {
		return downtimeData[i][0] < downtimeData[j][0]
	})
	sort.Slice(poweroffData[:], func(i, j int) bool {
		return poweroffData[i][0] < poweroffData[j][0]
	})
	return productionData, downtimeData, poweroffData
}

func downloadYearData(db *gorm.DB, fromDate time.Time, toDate time.Time, loc *time.Location, email string) (map[string]time.Duration, map[string]time.Duration, map[string]time.Duration) {
	productionRecords := make(map[string]time.Duration)
	downtimeRecords := make(map[string]time.Duration)
	poweroffRecords := make(map[string]time.Duration)
	var stateRecords []database.StateRecord
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
	for _, workplace := range workplaces {
		db.Select("state_id, date_time_start").Where("date_time_start >= ?", fromDate).Where("date_time_start <= ?", toDate).Where("workplace_id = ?", workplace.ID).Order("id asc").Find(&stateRecords)
		for index, record := range stateRecords {
			nextDate := time.Now().In(loc)
			if index < len(stateRecords)-1 {
				nextDate = stateRecords[index+1].DateTimeStart
			}
			if record.StateID == 1 {
				if record.DateTimeStart.In(loc).Day() == nextDate.In(loc).Day() {
					productionRecords[record.DateTimeStart.In(loc).Format("2006-01-02")] += nextDate.In(loc).Sub(record.DateTimeStart.In(loc))
				} else {
					endOfRecordDay := time.Date(record.DateTimeStart.In(loc).Year(), record.DateTimeStart.In(loc).Month(), record.DateTimeStart.In(loc).Day()+1, 0, 0, 0, 0, loc)
					for record.DateTimeStart.In(loc).Before(nextDate.In(loc)) {
						productionRecords[record.DateTimeStart.In(loc).Format("2006-01-02")] += endOfRecordDay.In(loc).Sub(record.DateTimeStart.In(loc))
						record.DateTimeStart = endOfRecordDay.In(loc)
						endOfRecordDay = time.Date(record.DateTimeStart.In(loc).Year(), record.DateTimeStart.In(loc).Month(), record.DateTimeStart.In(loc).Day()+1, 0, 0, 0, 0, loc)
					}
					endOfRecordDay = endOfRecordDay.In(loc).Add(-24 * time.Hour)
					productionRecords[nextDate.In(loc).Format("2006-01-02")] += nextDate.In(loc).Sub(endOfRecordDay.In(loc))
				}
			}
			if record.StateID == 2 {
				if record.DateTimeStart.In(loc).Day() == nextDate.In(loc).Day() {
					downtimeRecords[record.DateTimeStart.In(loc).Format("2006-01-02")] += nextDate.In(loc).Sub(record.DateTimeStart.In(loc))
				} else {
					endOfRecordDay := time.Date(record.DateTimeStart.In(loc).Year(), record.DateTimeStart.In(loc).Month(), record.DateTimeStart.In(loc).Day()+1, 0, 0, 0, 0, loc)
					for record.DateTimeStart.In(loc).Before(nextDate.In(loc)) {
						downtimeRecords[record.DateTimeStart.In(loc).Format("2006-01-02")] += endOfRecordDay.In(loc).Sub(record.DateTimeStart.In(loc))
						record.DateTimeStart = endOfRecordDay.In(loc)
						endOfRecordDay = time.Date(record.DateTimeStart.In(loc).Year(), record.DateTimeStart.In(loc).Month(), record.DateTimeStart.In(loc).Day()+1, 0, 0, 0, 0, loc)
					}
					endOfRecordDay = endOfRecordDay.In(loc).Add(-24 * time.Hour)
					downtimeRecords[nextDate.In(loc).Format("2006-01-02")] += nextDate.In(loc).Sub(endOfRecordDay.In(loc))
				}
			}
		}
	}
	today := time.Now().In(loc).Format("2006-01-02")
	for date, productionDuration := range productionRecords {
		downtimeDuration := downtimeRecords[date]
		if date == today {
			todaysDuration := time.Now().In(loc).Sub(time.Date(time.Now().In(loc).Year(), time.Now().In(loc).Month(), time.Now().In(loc).Day(), 0, 0, 0, 0, loc))
			poweroffRecords[date] = time.Duration(len(workplaces))*todaysDuration - downtimeDuration - productionDuration
			continue
		}
		poweroffRecords[date] = time.Duration(len(workplaces)*24)*time.Hour - downtimeDuration - productionDuration
	}

	return productionRecords, downtimeRecords, poweroffRecords
}

func downloadData(db *gorm.DB, fromDate time.Time, toDate time.Time, workplaceId uint, loc *time.Location, email string, stateId int) map[string]time.Duration {
	stateRecordsAsMap := make(map[string]time.Duration)
	var stateRecords []database.StateRecord
	if workplaceId == 0 {
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
		for _, workplace := range workplaces {
			db.Select("state_id, date_time_start").Where("date_time_start >= ?", fromDate).Where("date_time_start <= ?", toDate).Where("workplace_id = ?", workplace.ID).Order("id asc").Find(&stateRecords)
			for index, record := range stateRecords {
				nextDate := time.Now().In(loc)
				if index < len(stateRecords)-1 {
					nextDate = stateRecords[index+1].DateTimeStart
				}
				if record.StateID == stateId {
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
			if record.StateID == stateId {
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

func downloadProductionData(db *gorm.DB, loc *time.Location, email string) ([]string, []float64) {
	var workplaceNames []string
	var workplacePercents []float64
	var indexDataWorkplaces []IndexDataWorkplace
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
	for _, workplace := range workplaces {
		data := downloadData(db, time.Date(time.Now().UTC().Year(), time.Now().UTC().Month(), time.Now().UTC().Day(), 0, 0, 0, 0, time.Now().Location()), time.Now().In(loc), workplace.ID, loc, email, 1)
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
	data.Compacted = cachedUserWebSettings[email]["menu"]
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
