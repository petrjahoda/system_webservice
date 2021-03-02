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

func processStates(writer http.ResponseWriter, workplaceIds string, dateFrom time.Time, dateTo time.Time, email string) {
	timer := time.Now()
	logInfo("DATA-STATES", "Processing states started")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("DATA-STATES", "Problem opening database: "+err.Error())
		var responseData DataPageOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("DATA-STATES", "Processing data ended")
		return
	}
	var orderRecords []database.StateRecord
	if workplaceIds == "workplace_id in (')" {
		db.Where("date_time_start >= ?", dateFrom).Where("date_time_start <= ?", dateTo).Order("date_time_start desc").Find(&orderRecords)
	} else {
		db.Where(workplaceIds).Where("date_time_start >= ?", dateFrom).Where("date_time_start <= ?", dateTo).Order("date_time_start desc").Find(&orderRecords)
	}
	var data TableData
	data.DataTableSearchTitle = getLocale(email, "data-table-search-title")
	data.DataTableInfoTitle = getLocale(email, "data-table-info-title")
	data.DataTableRowsCountTitle = getLocale(email, "data-table-rows-count-title")
	addStateTableHeaders(email, &data)
	for _, record := range orderRecords {
		addStateTableRow(record, &data)
	}
	tmpl := template.Must(template.ParseFiles("./html/table.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("DATA-STATES", "States processed in "+time.Since(timer).String())
}

func addStateTableRow(record database.StateRecord, data *TableData) {
	var tableRow TableRow
	workplaceNameCell := TableCell{CellName: cachedWorkplacesById[uint(record.WorkplaceID)].Name}
	tableRow.TableCell = append(tableRow.TableCell, workplaceNameCell)
	stateStartDate := TableCell{CellName: record.DateTimeStart.Format("2006-01-02 15:04:05")}
	tableRow.TableCell = append(tableRow.TableCell, stateStartDate)
	stateName := TableCell{CellName: cachedStatesById[uint(record.StateID)].Name}
	tableRow.TableCell = append(tableRow.TableCell, stateName)
	note := TableCell{CellName: record.Note}
	tableRow.TableCell = append(tableRow.TableCell, note)
	data.TableRows = append(data.TableRows, tableRow)
}

func addStateTableHeaders(email string, data *TableData) {
	workplaceName := HeaderCell{HeaderName: getLocale(email, "workplace-name")}
	data.TableHeader = append(data.TableHeader, workplaceName)
	stateStart := HeaderCell{HeaderName: getLocale(email, "state-start")}
	data.TableHeader = append(data.TableHeader, stateStart)
	stateName := HeaderCell{HeaderName: getLocale(email, "state-name")}
	data.TableHeader = append(data.TableHeader, stateName)
	noteName := HeaderCell{HeaderName: getLocale(email, "note-name")}
	data.TableHeader = append(data.TableHeader, noteName)
}
