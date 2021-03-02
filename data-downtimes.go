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

func processDowntimes(writer http.ResponseWriter, workplaceIds string, dateFrom time.Time, dateTo time.Time, email string) {
	timer := time.Now()
	logInfo("DATA-DOWNTIMES", "Processing downtimes started")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("DATA-DOWNTIMES", "Problem opening database: "+err.Error())
		var responseData DataPageOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("DATA-DOWNTIMES", "Processing data ended")
		return
	}
	var downtimeRecords []database.DowntimeRecord
	var userRecords []database.UserRecord
	if workplaceIds == "workplace_id in (')" {
		db.Where("date_time_start >= ?", dateFrom).Where("date_time_start <= ?", dateTo).Order("date_time_start desc").Find(&downtimeRecords)
		db.Where("date_time_start >= ?", dateFrom).Where("date_time_start <= ?", dateTo).Find(&userRecords)
	} else {
		db.Where(workplaceIds).Where("date_time_start >= ?", dateFrom).Where("date_time_start <= ?", dateTo).Order("date_time_start desc").Find(&downtimeRecords)
		db.Where(workplaceIds).Where("date_time_start >= ?", dateFrom).Where("date_time_start <= ?", dateTo).Find(&userRecords)
	}
	var userRecordsByRecordId = map[int]database.UserRecord{}
	for _, record := range userRecords {
		userRecordsByRecordId[record.OrderRecordID] = record
	}
	var data TableData
	data.DataTableSearchTitle = getLocale(email, "data-table-search-title")
	data.DataTableInfoTitle = getLocale(email, "data-table-info-title")
	data.DataTableRowsCountTitle = getLocale(email, "data-table-rows-count-title")
	addDowntimeTableHeaders(email, &data)
	for _, record := range downtimeRecords {
		addDowntimeTableRow(record, userRecordsByRecordId, &data)
	}
	tmpl := template.Must(template.ParseFiles("./html/table.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("DATA-DOWNTIMES", "Downtimes processed in "+time.Since(timer).String())
}

func addDowntimeTableRow(record database.DowntimeRecord, userRecordsByRecordId map[int]database.UserRecord, data *TableData) {
	var tableRow TableRow
	workplaceNameCell := TableCell{CellName: cachedWorkplacesById[uint(record.WorkplaceID)].Name}
	tableRow.TableCell = append(tableRow.TableCell, workplaceNameCell)
	orderStart := TableCell{CellName: record.DateTimeStart.Format("2006-01-02 15:04:05")}
	tableRow.TableCell = append(tableRow.TableCell, orderStart)
	if record.DateTimeEnd.Time.IsZero() {
		orderEnd := TableCell{CellName: time.Now().Format("2006-01-02 15:04:05") + " +"}
		tableRow.TableCell = append(tableRow.TableCell, orderEnd)
	} else {
		orderEnd := TableCell{CellName: record.DateTimeEnd.Time.Format("2006-01-02 15:04:05")}
		tableRow.TableCell = append(tableRow.TableCell, orderEnd)
	}
	actualUserId := userRecordsByRecordId[int(record.ID)].UserID
	userName := TableCell{CellName: cachedUsersById[uint(actualUserId)].FirstName + " " + cachedUsersById[uint(actualUserId)].SecondName}
	tableRow.TableCell = append(tableRow.TableCell, userName)
	downtimeName := TableCell{CellName: cachedDowntimesById[uint(record.DowntimeID)].Name}
	tableRow.TableCell = append(tableRow.TableCell, downtimeName)
	note := TableCell{CellName: record.Note}
	tableRow.TableCell = append(tableRow.TableCell, note)
	data.TableRows = append(data.TableRows, tableRow)
}

func addDowntimeTableHeaders(email string, data *TableData) {
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
