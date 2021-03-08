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

func processBreakdowns(writer http.ResponseWriter, workplaceIds string, dateFrom time.Time, dateTo time.Time, email string) {
	timer := time.Now()
	logInfo("DATA-BREAKDOWNS", "Processing breakdowns started")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("DATA-BREAKDOWNS", "Problem opening database: "+err.Error())
		var responseData DataPageOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("DATA-BREAKDOWNS", "Processing data ended")
		return
	}
	var breakdownRecords []database.BreakdownRecord
	if workplaceIds == "workplace_id in (')" {
		db.Where("date_time_start <= ? and date_time_end >= ?", dateTo, dateFrom).Or("date_time_start <= ? and date_time_end is null", dateTo).Or("date_time_start <= ? and date_time_end >= ?", dateFrom, dateTo).Order("date_time_start desc").Find(&breakdownRecords)
	} else {
		db.Where(workplaceIds).Where("date_time_start <= ? and date_time_end >= ?", dateTo, dateFrom).Or("date_time_start <= ? and date_time_end is null", dateTo).Or("date_time_start <= ? and date_time_end >= ?", dateFrom, dateTo).Order("date_time_start desc").Find(&breakdownRecords)
	}
	var data TableData
	data.DataTableSearchTitle = getLocale(email, "data-table-search-title")
	data.DataTableInfoTitle = getLocale(email, "data-table-info-title")
	data.DataTableRowsCountTitle = getLocale(email, "data-table-rows-count-title")
	addBreakdownTableHeaders(email, &data)
	for _, record := range breakdownRecords {
		addBreakdownTableRow(record, &data)
	}
	tmpl := template.Must(template.ParseFiles("./html/table.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("DATA-BREAKDOWNS", "Breakdowns processed in "+time.Since(timer).String())
}

func addBreakdownTableRow(record database.BreakdownRecord, data *TableData) {
	var tableRow TableRow
	workplaceNameCell := TableCell{CellName: cachedWorkplacesById[uint(record.WorkplaceID)].Name}
	tableRow.TableCell = append(tableRow.TableCell, workplaceNameCell)
	breakdownStart := TableCell{CellName: record.DateTimeStart.Format("2006-01-02 15:04:05")}
	tableRow.TableCell = append(tableRow.TableCell, breakdownStart)
	if record.DateTimeEnd.Time.IsZero() {
		breakdownEnd := TableCell{CellName: time.Now().Format("2006-01-02 15:04:05") + " +"}
		tableRow.TableCell = append(tableRow.TableCell, breakdownEnd)
	} else {
		breakdownEnd := TableCell{CellName: record.DateTimeEnd.Time.Format("2006-01-02 15:04:05")}
		tableRow.TableCell = append(tableRow.TableCell, breakdownEnd)
	}
	userName := TableCell{CellName: cachedUsersById[uint(record.UserID)].FirstName + " " + cachedUsersById[uint(record.UserID)].SecondName}
	tableRow.TableCell = append(tableRow.TableCell, userName)
	breakdownName := TableCell{CellName: cachedBreakdownsById[uint(record.BreakdownID)].Name}
	tableRow.TableCell = append(tableRow.TableCell, breakdownName)
	note := TableCell{CellName: record.Note}
	tableRow.TableCell = append(tableRow.TableCell, note)
	data.TableRows = append(data.TableRows, tableRow)
}

func addBreakdownTableHeaders(email string, data *TableData) {
	workplaceName := HeaderCell{HeaderName: getLocale(email, "workplace-name")}
	data.TableHeader = append(data.TableHeader, workplaceName)
	breakdownStart := HeaderCell{HeaderName: getLocale(email, "breakdown-start")}
	data.TableHeader = append(data.TableHeader, breakdownStart)
	breakdownEnd := HeaderCell{HeaderName: getLocale(email, "breakdown-end")}
	data.TableHeader = append(data.TableHeader, breakdownEnd)
	userName := HeaderCell{HeaderName: getLocale(email, "user-name")}
	data.TableHeader = append(data.TableHeader, userName)
	breakdownName := HeaderCell{HeaderName: getLocale(email, "breakdown-name")}
	data.TableHeader = append(data.TableHeader, breakdownName)
	noteName := HeaderCell{HeaderName: getLocale(email, "note-name")}
	data.TableHeader = append(data.TableHeader, noteName)
}
