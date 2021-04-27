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

func loadUsersTable(writer http.ResponseWriter, workplaceIds string, dateFrom time.Time, dateTo time.Time, email string) {
	timer := time.Now()
	logInfo("DATA-USERS", "Loading users table")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("DATA-USERS", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "ERR: Problem opening database, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("DATA-USERS", "Loading users table ended")
		return
	}
	var userRecords []database.UserRecord
	if workplaceIds == "workplace_id in (')" {
		db.Where("date_time_start <= ? and date_time_end >= ?", dateTo, dateFrom).Or("date_time_start <= ? and date_time_end is null", dateTo).Or("date_time_start <= ? and date_time_end >= ?", dateFrom, dateTo).Order("date_time_start desc").Find(&userRecords)
	} else {
		db.Where("date_time_start <= ? and date_time_end >= ?", dateTo, dateFrom).Where(workplaceIds).Or("date_time_start <= ? and date_time_end is null", dateTo).Where(workplaceIds).Or("date_time_start <= ? and date_time_end >= ?", dateFrom, dateTo).Where(workplaceIds).Order("date_time_start desc").Find(&userRecords)
	}
	var data TableOutput
	data.DataTableSearchTitle = getLocale(email, "data-table-search-title")
	data.DataTableInfoTitle = getLocale(email, "data-table-info-title")
	data.DataTableRowsCountTitle = getLocale(email, "data-table-rows-count-title")
	loc, err := time.LoadLocation(location)
	addUserTableHeaders(email, &data)
	for _, record := range userRecords {
		addUserTableRow(record, &data, db, loc)
	}
	tmpl, err := template.ParseFiles("./html/data-content.html")
	if err != nil {
		logError("SETTINGS", "Problem parsing html file: "+err.Error())
		var responseData TableOutput
		responseData.Result = "ERR: Problem parsing html file: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
	} else {
		data.Result = "INF: Users data processed in " + time.Since(timer).String()
		_ = tmpl.Execute(writer, data)
		logInfo("SETTINGS", "Users data loaded in "+time.Since(timer).String())
	}
}

func addUserTableRow(record database.UserRecord, data *TableOutput, db *gorm.DB, loc *time.Location) {
	var tableRow TableRow
	workplaceNameCell := TableCell{CellName: cachedWorkplacesById[uint(record.WorkplaceID)].Name}
	tableRow.TableCell = append(tableRow.TableCell, workplaceNameCell)
	userStart := TableCell{CellName: record.DateTimeStart.In(loc).Format("2006-01-02 15:04:05")}
	tableRow.TableCell = append(tableRow.TableCell, userStart)
	if record.DateTimeEnd.Time.IsZero() {
		userEnd := TableCell{CellName: time.Now().In(loc).Format("2006-01-02 15:04:05") + " +"}
		tableRow.TableCell = append(tableRow.TableCell, userEnd)
	} else {
		orderEnd := TableCell{CellName: record.DateTimeEnd.Time.In(loc).Format("2006-01-02 15:04:05")}
		tableRow.TableCell = append(tableRow.TableCell, orderEnd)
	}
	userName := TableCell{CellName: cachedUsersById[uint(record.UserID)].FirstName + " " + cachedUsersById[uint(record.UserID)].SecondName}
	tableRow.TableCell = append(tableRow.TableCell, userName)
	var orderRecord database.OrderRecord
	db.Where("id = ?", record.OrderRecordID).Find(&orderRecord)
	orderName := TableCell{CellName: cachedOrdersById[uint(orderRecord.OrderID)].Name}
	tableRow.TableCell = append(tableRow.TableCell, orderName)
	note := TableCell{CellName: record.Note}
	tableRow.TableCell = append(tableRow.TableCell, note)
	data.TableRows = append(data.TableRows, tableRow)
}

func addUserTableHeaders(email string, data *TableOutput) {
	workplaceName := HeaderCell{HeaderName: getLocale(email, "workplace-name")}
	data.TableHeader = append(data.TableHeader, workplaceName)
	userStart := HeaderCell{HeaderName: getLocale(email, "user-start")}
	data.TableHeader = append(data.TableHeader, userStart)
	userEnd := HeaderCell{HeaderName: getLocale(email, "user-end")}
	data.TableHeader = append(data.TableHeader, userEnd)
	userName := HeaderCell{HeaderName: getLocale(email, "user-name")}
	data.TableHeader = append(data.TableHeader, userName)
	orderName := HeaderCell{HeaderName: getLocale(email, "order-name")}
	data.TableHeader = append(data.TableHeader, orderName)
	noteName := HeaderCell{HeaderName: getLocale(email, "note-name")}
	data.TableHeader = append(data.TableHeader, noteName)
}
