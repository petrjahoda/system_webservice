package main

import (
	"encoding/json"
	"github.com/petrjahoda/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"html/template"
	"net/http"
	"time"
)

func loadAlarmsTable(writer http.ResponseWriter, workplaceIds string, dateFrom time.Time, dateTo time.Time, email string) {
	timer := time.Now()
	logInfo("DATA-ALARMS", "Loading alarms table")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("DATA-ALARMS", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "ERR: Problem opening database, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("DATA-ALARMS", "Loading alarms table ended")
		return
	}
	var alarmRecords []database.AlarmRecord
	if workplaceIds == "workplace_id in (')" {
		db.Where("date_time_start <= ? and date_time_end >= ?", dateTo, dateFrom).Or("date_time_start <= ? and date_time_end is null", dateTo).Or("date_time_start <= ? and date_time_end >= ?", dateFrom, dateTo).Order("date_time_start desc").Find(&alarmRecords)
	} else {
		db.Where("date_time_start <= ? and date_time_end >= ?", dateTo, dateFrom).Where(workplaceIds).Or("date_time_start <= ? and date_time_end is null", dateTo).Where(workplaceIds).Or("date_time_start <= ? and date_time_end >= ?", dateFrom, dateTo).Where(workplaceIds).Order("date_time_start desc").Find(&alarmRecords)
	}
	var data TableOutput
	data.Compacted = cachedUserWebSettings[email]["data-selected-size"]
	data.DataTableSearchTitle = getLocale(email, "data-table-search-title")
	data.DataTableInfoTitle = getLocale(email, "data-table-info-title")
	data.DataTableRowsCountTitle = getLocale(email, "data-table-rows-count-title")
	companyNameSync.Lock()
	loc, err := time.LoadLocation(location)
	companyNameSync.Unlock()
	addAlarmTableHeaders(email, &data)
	for _, record := range alarmRecords {
		addAlarmTableRow(record, &data, loc)
	}
	tmpl, err := template.ParseFiles("./html/data-content.html")
	if err != nil {
		logError("SETTINGS", "Problem parsing html file: "+err.Error())
		var responseData TableOutput
		responseData.Result = "ERR: Problem parsing html file: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
	} else {
		data.Result = "INF: Alarms data processed in " + time.Since(timer).String()
		_ = tmpl.Execute(writer, data)
		logInfo("SETTINGS", "Alarms data loaded in "+time.Since(timer).String())
	}
}

func addAlarmTableRow(record database.AlarmRecord, data *TableOutput, loc *time.Location) {
	var tableRow TableRow
	workplaceNameCell := TableCell{CellName: cachedWorkplacesById[uint(record.WorkplaceID)].Name}
	tableRow.TableCell = append(tableRow.TableCell, workplaceNameCell)
	alarmStart := TableCell{CellName: record.DateTimeStart.In(loc).Format("2006-01-02 15:04:05")}
	tableRow.TableCell = append(tableRow.TableCell, alarmStart)
	if record.DateTimeEnd.Time.IsZero() {
		alarmEnd := TableCell{CellName: time.Now().In(loc).Format("2006-01-02 15:04:05") + " +"}
		tableRow.TableCell = append(tableRow.TableCell, alarmEnd)
	} else {
		alarmEnd := TableCell{CellName: record.DateTimeEnd.Time.In(loc).Format("2006-01-02 15:04:05")}
		tableRow.TableCell = append(tableRow.TableCell, alarmEnd)
	}
	alarmName := TableCell{CellName: cachedAlarmsById[uint(record.AlarmID)].Name}
	tableRow.TableCell = append(tableRow.TableCell, alarmName)
	if record.DateTimeProcessed.Time.IsZero() {
		alarmProcessed := TableCell{CellName: time.Now().In(loc).Format("2006-01-02 15:04:05") + " +"}
		tableRow.TableCell = append(tableRow.TableCell, alarmProcessed)
	} else {
		alarmProcessed := TableCell{CellName: record.DateTimeEnd.Time.In(loc).Format("2006-01-02 15:04:05")}
		tableRow.TableCell = append(tableRow.TableCell, alarmProcessed)
	}
	data.TableRows = append(data.TableRows, tableRow)
}

func addAlarmTableHeaders(email string, data *TableOutput) {
	workplaceName := HeaderCell{HeaderName: getLocale(email, "workplace-name")}
	data.TableHeader = append(data.TableHeader, workplaceName)
	alarmStart := HeaderCell{HeaderName: getLocale(email, "alarm-start")}
	data.TableHeader = append(data.TableHeader, alarmStart)
	alarmEnd := HeaderCell{HeaderName: getLocale(email, "alarm-end")}
	data.TableHeader = append(data.TableHeader, alarmEnd)
	alarmName := HeaderCell{HeaderName: getLocale(email, "alarm-name")}
	data.TableHeader = append(data.TableHeader, alarmName)
	alarmProcessed := HeaderCell{HeaderName: getLocale(email, "alarm-processed")}
	data.TableHeader = append(data.TableHeader, alarmProcessed)
}
