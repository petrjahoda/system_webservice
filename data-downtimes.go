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

func loadDowntimesTable(writer http.ResponseWriter, workplaceIds string, dateFrom time.Time, dateTo time.Time, email string) {
	timer := time.Now()
	logInfo("DATA-DOWNTIMES", "Loading downtimes table")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("DATA-DOWNTIMES", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "ERR: Problem opening database, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("DATA-DOWNTIMES", "Loading downtimes table ended")
		return
	}
	var downtimeRecords []database.DowntimeRecord
	if workplaceIds == "workplace_id in (')" {
		db.Where("date_time_start <= ? and date_time_end >= ?", dateTo, dateFrom).Or("date_time_start <= ? and date_time_end is null", dateTo).Or("date_time_start <= ? and date_time_end >= ?", dateFrom, dateTo).Order("date_time_start desc").Find(&downtimeRecords)
	} else {
		db.Where("date_time_start <= ? and date_time_end >= ?", dateTo, dateFrom).Where(workplaceIds).Or("date_time_start <= ? and date_time_end is null", dateTo).Where(workplaceIds).Or("date_time_start <= ? and date_time_end >= ?", dateFrom, dateTo).Where(workplaceIds).Order("date_time_start desc").Find(&downtimeRecords)
	}
	var data TableOutput
	userWebSettingsSync.RLock()
	data.Compacted = cachedUserWebSettings[email]["data-selected-size"]
	userWebSettingsSync.RUnlock()
	data.DataTableSearchTitle = getLocale(email, "data-table-search-title")
	data.DataTableInfoTitle = getLocale(email, "data-table-info-title")
	data.DataTableRowsCountTitle = getLocale(email, "data-table-rows-count-title")
	locationSync.RLock()
	loc, err := time.LoadLocation(cachedLocation)
	locationSync.RUnlock()
	addDowntimeTableHeaders(email, &data)
	for _, record := range downtimeRecords {
		addDowntimeTableRow(record, &data, loc)
	}
	tmpl, err := template.ParseFiles("./html/data-content.html")
	if err != nil {
		logError("SETTINGS", "Problem parsing html file: "+err.Error())
		var responseData TableOutput
		responseData.Result = "ERR: Problem parsing html file: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
	} else {
		data.Result = "INF: Downtimes data processed in " + time.Since(timer).String()
		_ = tmpl.Execute(writer, data)
		logInfo("SETTINGS", "Downtimes data loaded in "+time.Since(timer).String())
	}
}

func addDowntimeTableRow(downtimeRecord database.DowntimeRecord, data *TableOutput, loc *time.Location) {
	var tableRow TableRow
	workplacesByIdSync.RLock()
	workplaceNameCell := TableCell{CellName: cachedWorkplacesById[uint(downtimeRecord.WorkplaceID)].Name}
	workplacesByIdSync.RUnlock()
	tableRow.TableCell = append(tableRow.TableCell, workplaceNameCell)
	orderStart := TableCell{CellName: downtimeRecord.DateTimeStart.In(loc).Format("2006-01-02 15:04:05")}
	tableRow.TableCell = append(tableRow.TableCell, orderStart)
	if downtimeRecord.DateTimeEnd.Time.IsZero() {
		orderEnd := TableCell{CellName: time.Now().In(loc).Format("2006-01-02 15:04:05") + " +"}
		tableRow.TableCell = append(tableRow.TableCell, orderEnd)
	} else {
		orderEnd := TableCell{CellName: downtimeRecord.DateTimeEnd.Time.In(loc).Format("2006-01-02 15:04:05")}
		tableRow.TableCell = append(tableRow.TableCell, orderEnd)
	}
	actualUserId := downtimeRecord.UserID.Int32
	usersByIdSync.RLock()
	userName := TableCell{CellName: cachedUsersById[uint(actualUserId)].FirstName + " " + cachedUsersById[uint(actualUserId)].SecondName}
	usersByIdSync.RUnlock()
	tableRow.TableCell = append(tableRow.TableCell, userName)
	downtimesByIdSync.RLock()
	downtimeName := TableCell{CellName: cachedDowntimesById[uint(downtimeRecord.DowntimeID)].Name}
	downtimesByIdSync.RUnlock()
	tableRow.TableCell = append(tableRow.TableCell, downtimeName)
	note := TableCell{CellName: downtimeRecord.Note}
	tableRow.TableCell = append(tableRow.TableCell, note)
	data.TableRows = append(data.TableRows, tableRow)
}

func addDowntimeTableHeaders(email string, data *TableOutput) {
	workplaceName := HeaderCell{HeaderName: getLocale(email, "workplace-name")}
	data.TableHeader = append(data.TableHeader, workplaceName)
	downtimeStart := HeaderCell{HeaderName: getLocale(email, "downtime-start")}
	data.TableHeader = append(data.TableHeader, downtimeStart)
	downtimeEnd := HeaderCell{HeaderName: getLocale(email, "downtime-end")}
	data.TableHeader = append(data.TableHeader, downtimeEnd)
	userName := HeaderCell{HeaderName: getLocale(email, "user-name")}
	data.TableHeader = append(data.TableHeader, userName)
	downtimeName := HeaderCell{HeaderName: getLocale(email, "downtime-name")}
	data.TableHeader = append(data.TableHeader, downtimeName)
	noteName := HeaderCell{HeaderName: getLocale(email, "note-name")}
	data.TableHeader = append(data.TableHeader, noteName)
}
