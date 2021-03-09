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

func processFaults(writer http.ResponseWriter, workplaceIds string, dateFrom time.Time, dateTo time.Time, email string) {
	timer := time.Now()
	logInfo("DATA-FAULTS", "Processing faults started")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("DATA-FAULTS", "Problem opening database: "+err.Error())
		var responseData DataPageOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("DATA-FAULTS", "Processing data ended")
		return
	}
	var orderRecords []database.FaultRecord
	if workplaceIds == "workplace_id in (')" {
		db.Where("date_time >= ?", dateFrom).Where("date_time <= ?", dateTo).Order("date_time desc").Find(&orderRecords)
	} else {
		db.Where("date_time >= ?", dateFrom).Where("date_time <= ?", dateTo).Where(workplaceIds).Order("date_time desc").Find(&orderRecords)
	}
	var data TableData
	data.DataTableSearchTitle = getLocale(email, "data-table-search-title")
	data.DataTableInfoTitle = getLocale(email, "data-table-info-title")
	data.DataTableRowsCountTitle = getLocale(email, "data-table-rows-count-title")
	addFaultTableHeaders(email, &data)
	for _, record := range orderRecords {
		addFaultTableRow(record, &data, db)
	}
	tmpl := template.Must(template.ParseFiles("./html/data-content.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("DATA-FAULTS", "Faults processed in "+time.Since(timer).String())
}

func addFaultTableRow(record database.FaultRecord, data *TableData, db *gorm.DB) {
	var tableRow TableRow
	workplaceNameCell := TableCell{CellName: cachedWorkplacesById[uint(record.WorkplaceID)].Name}
	tableRow.TableCell = append(tableRow.TableCell, workplaceNameCell)
	faultDate := TableCell{CellName: record.DateTime.Format("2006-01-02 15:04:05")}
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

func addFaultTableHeaders(email string, data *TableData) {
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
