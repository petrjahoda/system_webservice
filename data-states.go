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

func loadStatesTable(writer http.ResponseWriter, workplaceIds string, dateFrom time.Time, dateTo time.Time, email string) {
	timer := time.Now()
	logInfo("DATA-STATES", "Loading states table")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("DATA-STATES", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "ERR: Problem opening database, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("DATA-STATES", "Loading states table ended")
		return
	}
	var orderRecords []database.StateRecord
	if workplaceIds == "workplace_id in (')" {
		db.Where("date_time_start >= ?", dateFrom).Where("date_time_start <= ?", dateTo).Order("date_time_start desc").Find(&orderRecords)
	} else {
		db.Where("date_time_start >= ?", dateFrom).Where("date_time_start <= ?", dateTo).Where(workplaceIds).Order("date_time_start desc").Find(&orderRecords)
	}
	var data TableOutput
	userWebSettingsSync.RLock()
	data.Compacted = cachedUserWebSettings[email]["data-selected-size"]
	userWebSettingsSync.RUnlock()
	data.DataTableSearchTitle = getLocale(email, "data-table-search-title")
	data.DataTableInfoTitle = getLocale(email, "data-table-info-title")
	data.DataTableRowsCountTitle = getLocale(email, "data-table-rows-count-title")
	addStateTableHeaders(email, &data)
	locationSync.RLock()
	loc, err := time.LoadLocation(cachedLocation)
	locationSync.RUnlock()
	for _, record := range orderRecords {
		addStateTableRow(record, loc, &data)
	}
	tmpl, err := template.ParseFiles("./html/data-content.html")
	if err != nil {
		logError("SETTINGS", "Problem parsing html file: "+err.Error())
		var responseData TableOutput
		responseData.Result = "ERR: Problem parsing html file: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
	} else {
		data.Result = "INF: States data processed in " + time.Since(timer).String()
		_ = tmpl.Execute(writer, data)
		logInfo("SETTINGS", "States data loaded in "+time.Since(timer).String())
	}
}

func addStateTableRow(record database.StateRecord, loc *time.Location, data *TableOutput) {
	var tableRow TableRow
	workplacesByIdSync.RLock()
	workplaceNameCell := TableCell{CellName: cachedWorkplacesById[uint(record.WorkplaceID)].Name}
	workplacesByIdSync.RUnlock()
	tableRow.TableCell = append(tableRow.TableCell, workplaceNameCell)
	stateStartDate := TableCell{CellName: record.DateTimeStart.In(loc).Format("2006-01-02 15:04:05")}
	tableRow.TableCell = append(tableRow.TableCell, stateStartDate)
	statesByIdSync.RLock()
	stateName := TableCell{CellName: cachedStatesById[uint(record.StateID)].Name}
	statesByIdSync.RUnlock()
	tableRow.TableCell = append(tableRow.TableCell, stateName)
	note := TableCell{CellName: record.Note}
	tableRow.TableCell = append(tableRow.TableCell, note)
	data.TableRows = append(data.TableRows, tableRow)
}

func addStateTableHeaders(email string, data *TableOutput) {
	workplaceName := HeaderCell{HeaderName: getLocale(email, "workplace-name")}
	data.TableHeader = append(data.TableHeader, workplaceName)
	stateStart := HeaderCell{HeaderName: getLocale(email, "state-start")}
	data.TableHeader = append(data.TableHeader, stateStart)
	stateName := HeaderCell{HeaderName: getLocale(email, "state-name")}
	data.TableHeader = append(data.TableHeader, stateName)
	noteName := HeaderCell{HeaderName: getLocale(email, "note-name")}
	data.TableHeader = append(data.TableHeader, noteName)
}
