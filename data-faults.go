package main

import (
	"encoding/json"
	"github.com/petrjahoda/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

func loadFaultsTable(writer http.ResponseWriter, workplaceIds string, dateFrom time.Time, dateTo time.Time, email string) {
	timer := time.Now()
	logInfo("DATA-FAULTS", "Loading faults table")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("DATA-FAULTS", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "ERR: Problem opening database, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("DATA-FAULTS", "Loading faults table ended")
		return
	}
	var orderRecords []database.FaultRecord
	if workplaceIds == "workplace_id in (')" {
		db.Where("date_time >= ?", dateFrom).Where("date_time <= ?", dateTo).Order("date_time desc").Find(&orderRecords)
	} else {
		db.Where("date_time >= ?", dateFrom).Where("date_time <= ?", dateTo).Where(workplaceIds).Order("date_time desc").Find(&orderRecords)
	}
	var data TableOutput
	data.DataTableSearchTitle = getLocale(email, "data-table-search-title")
	data.DataTableInfoTitle = getLocale(email, "data-table-info-title")
	data.DataTableRowsCountTitle = getLocale(email, "data-table-rows-count-title")
	loc, err := time.LoadLocation(location)
	addFaultTableHeaders(email, &data)
	for _, record := range orderRecords {
		addFaultTableRow(record, &data, db, loc)
	}
	tmpl, err := template.ParseFiles("./html/data-content.html")
	if err != nil {
		logError("SETTINGS", "Problem parsing html file: "+err.Error())
		var responseData TableOutput
		responseData.Result = "ERR: Problem parsing html file: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
	} else {
		data.Result = "INF: Faults data processed in " + time.Since(timer).String()
		_ = tmpl.Execute(writer, data)
		logInfo("SETTINGS", "Faults data loaded in "+time.Since(timer).String())
	}
}

func addFaultTableRow(record database.FaultRecord, data *TableOutput, db *gorm.DB, loc *time.Location) {
	var tableRow TableRow
	workplaceNameCell := TableCell{CellName: cachedWorkplacesById[uint(record.WorkplaceID)].Name}
	tableRow.TableCell = append(tableRow.TableCell, workplaceNameCell)
	faultDate := TableCell{CellName: record.DateTime.In(loc).Format("2006-01-02 15:04:05")}
	tableRow.TableCell = append(tableRow.TableCell, faultDate)
	userName := TableCell{CellName: cachedUsersById[uint(record.UserID)].FirstName + " " + cachedUsersById[uint(record.UserID)].SecondName}
	tableRow.TableCell = append(tableRow.TableCell, userName)
	var orderRecord database.OrderRecord
	db.Where("id = ?", record.OrderRecordID).Find(&orderRecord)
	orderName := TableCell{CellName: cachedOrdersById[uint(orderRecord.OrderID)].Name}
	tableRow.TableCell = append(tableRow.TableCell, orderName)
	faultName := TableCell{CellName: cachedFaultsById[uint(record.FaultID)].Name}
	tableRow.TableCell = append(tableRow.TableCell, faultName)
	countAsString := strconv.Itoa(record.Count)
	count := TableCell{CellName: countAsString}
	tableRow.TableCell = append(tableRow.TableCell, count)
	note := TableCell{CellName: record.Note}
	tableRow.TableCell = append(tableRow.TableCell, note)
	data.TableRows = append(data.TableRows, tableRow)
}

func addFaultTableHeaders(email string, data *TableOutput) {
	workplaceName := HeaderCell{HeaderName: getLocale(email, "workplace-name")}
	data.TableHeader = append(data.TableHeader, workplaceName)
	faultDate := HeaderCell{HeaderName: getLocale(email, "fault-date")}
	data.TableHeader = append(data.TableHeader, faultDate)
	userName := HeaderCell{HeaderName: getLocale(email, "user-name")}
	data.TableHeader = append(data.TableHeader, userName)
	orderName := HeaderCell{HeaderName: getLocale(email, "order-name")}
	data.TableHeader = append(data.TableHeader, orderName)
	faultName := HeaderCell{HeaderName: getLocale(email, "fault-name")}
	data.TableHeader = append(data.TableHeader, faultName)
	faultCount := HeaderCell{HeaderName: getLocale(email, "fault-count")}
	data.TableHeader = append(data.TableHeader, faultCount)
	noteName := HeaderCell{HeaderName: getLocale(email, "note-name")}
	data.TableHeader = append(data.TableHeader, noteName)
}
