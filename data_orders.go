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

func loadOrdersTable(writer http.ResponseWriter, workplaceIds string, dateFrom time.Time, dateTo time.Time, email string) {
	timer := time.Now()
	logInfo("DATA-ORDERS", "Loading orders table")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("DATA-ORDERS", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "ERR: Problem opening database, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("DATA-ORDERS", "Loading orders table ended")
		return
	}
	var orderRecords []database.OrderRecord
	var userRecords []database.UserRecord
	if workplaceIds == "workplace_id in (')" {
		db.Where("date_time_start <= ? and date_time_end >= ?", dateTo, dateFrom).Or("date_time_start <= ? and date_time_end is null", dateTo).Or("date_time_start <= ? and date_time_end >= ?", dateFrom, dateTo).Order("date_time_start desc").Find(&orderRecords)
		db.Where("date_time_start <= ? and date_time_end >= ?", dateTo, dateFrom).Or("date_time_start <= ? and date_time_end is null", dateTo).Or("date_time_start <= ? and date_time_end >= ?", dateFrom, dateTo).Find(&userRecords)
	} else {
		db.Where("date_time_start <= ? and date_time_end >= ?", dateTo, dateFrom).Where(workplaceIds).Or("date_time_start <= ? and date_time_end is null", dateTo).Where(workplaceIds).Or("date_time_start <= ? and date_time_end >= ?", dateFrom, dateTo).Where(workplaceIds).Order("date_time_start desc").Find(&orderRecords)
		db.Where("date_time_start <= ? and date_time_end >= ?", dateTo, dateFrom).Where(workplaceIds).Or("date_time_start <= ? and date_time_end is null", dateTo).Where(workplaceIds).Or("date_time_start <= ? and date_time_end >= ?", dateFrom, dateTo).Where(workplaceIds).Find(&userRecords)
	}
	var userRecordsByRecordId = map[int]database.UserRecord{}
	for _, record := range userRecords {
		userRecordsByRecordId[record.OrderRecordID] = record
	}
	var data TableOutput
	data.Compacted = cachedUserWebSettings[email]["data-selected-size"]
	data.DataTableSearchTitle = getLocale(email, "data-table-search-title")
	data.DataTableInfoTitle = getLocale(email, "data-table-info-title")
	data.DataTableRowsCountTitle = getLocale(email, "data-table-rows-count-title")
	companyNameSync.Lock()
	loc, err := time.LoadLocation(location)
	companyNameSync.Unlock()
	addOrderTableHeaders(email, &data)
	for _, record := range orderRecords {
		addOrderTableRow(record, userRecordsByRecordId, &data, loc)
	}
	tmpl, err := template.ParseFiles("./html/data-content.html")
	if err != nil {
		logError("SETTINGS", "Problem parsing html file: "+err.Error())
		var responseData TableOutput
		responseData.Result = "ERR: Problem parsing html file: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
	} else {
		data.Result = "INF: Orders data processed in " + time.Since(timer).String()
		_ = tmpl.Execute(writer, data)
		logInfo("SETTINGS", "Orders data loaded in "+time.Since(timer).String())
	}
}

func addOrderTableHeaders(email string, data *TableOutput) {
	workplaceName := HeaderCell{HeaderName: getLocale(email, "workplace-name")}
	data.TableHeader = append(data.TableHeader, workplaceName)
	orderStart := HeaderCell{HeaderName: getLocale(email, "order-start")}
	data.TableHeader = append(data.TableHeader, orderStart)
	orderEnd := HeaderCell{HeaderName: getLocale(email, "order-end")}
	data.TableHeader = append(data.TableHeader, orderEnd)
	workplaceModeName := HeaderCell{HeaderName: getLocale(email, "workplacemode-name")}
	data.TableHeader = append(data.TableHeader, workplaceModeName)
	workshiftName := HeaderCell{HeaderName: getLocale(email, "workshift-name")}
	data.TableHeader = append(data.TableHeader, workshiftName)
	userName := HeaderCell{HeaderName: getLocale(email, "user-name")}
	data.TableHeader = append(data.TableHeader, userName)
	orderName := HeaderCell{HeaderName: getLocale(email, "order-name")}
	data.TableHeader = append(data.TableHeader, orderName)
	operationName := HeaderCell{HeaderName: getLocale(email, "operation-name")}
	data.TableHeader = append(data.TableHeader, operationName)
	cycleName := HeaderCell{HeaderName: getLocale(email, "cycle-name")}
	data.TableHeader = append(data.TableHeader, cycleName)
	cavityName := HeaderCell{HeaderName: getLocale(email, "cavity-name")}
	data.TableHeader = append(data.TableHeader, cavityName)
	goodPcsName := HeaderCell{HeaderName: getLocale(email, "good-pieces-name")}
	data.TableHeader = append(data.TableHeader, goodPcsName)
	badPcsName := HeaderCell{HeaderName: getLocale(email, "bad-pieces-name")}
	data.TableHeader = append(data.TableHeader, badPcsName)
	noteName := HeaderCell{HeaderName: getLocale(email, "note-name")}
	data.TableHeader = append(data.TableHeader, noteName)
}

func addOrderTableRow(record database.OrderRecord, userRecordsByRecordId map[int]database.UserRecord, data *TableOutput, loc *time.Location) {
	var tableRow TableRow
	workplaceNameCell := TableCell{CellName: cachedWorkplacesById[uint(record.WorkplaceID)].Name}
	tableRow.TableCell = append(tableRow.TableCell, workplaceNameCell)
	orderStart := TableCell{CellName: record.DateTimeStart.In(loc).Format("2006-01-02 15:04:05")}
	tableRow.TableCell = append(tableRow.TableCell, orderStart)
	if record.DateTimeEnd.Time.IsZero() {
		orderEnd := TableCell{CellName: time.Now().In(loc).Format("2006-01-02 15:04:05") + " +"}
		tableRow.TableCell = append(tableRow.TableCell, orderEnd)
	} else {
		orderEnd := TableCell{CellName: record.DateTimeEnd.Time.In(loc).Format("2006-01-02 15:04:05")}
		tableRow.TableCell = append(tableRow.TableCell, orderEnd)
	}
	workplaceModeNameCell := TableCell{CellName: cachedWorkplaceModesById[uint(record.WorkplaceModeID)].Name}
	tableRow.TableCell = append(tableRow.TableCell, workplaceModeNameCell)
	workshiftName := TableCell{CellName: cachedWorkShiftsById[uint(record.WorkshiftID)].Name}
	tableRow.TableCell = append(tableRow.TableCell, workshiftName)
	actualUserId := userRecordsByRecordId[int(record.ID)].UserID
	userName := TableCell{CellName: cachedUsersById[uint(actualUserId)].FirstName + " " + cachedUsersById[uint(actualUserId)].SecondName}
	tableRow.TableCell = append(tableRow.TableCell, userName)
	orderName := TableCell{CellName: cachedOrdersById[uint(record.OrderID)].Name}
	tableRow.TableCell = append(tableRow.TableCell, orderName)
	operationName := TableCell{CellName: cachedOperationsById[uint(record.OperationID)].Name}
	tableRow.TableCell = append(tableRow.TableCell, operationName)
	averageCycleAsString := strconv.FormatFloat(float64(record.AverageCycle), 'f', 2, 64)
	averageCycle := TableCell{CellName: averageCycleAsString + "s"}
	tableRow.TableCell = append(tableRow.TableCell, averageCycle)
	cavityAsString := strconv.Itoa(record.Cavity)
	cavity := TableCell{CellName: cavityAsString}
	tableRow.TableCell = append(tableRow.TableCell, cavity)
	okAsString := strconv.Itoa(record.CountOk)
	ok := TableCell{CellName: okAsString}
	tableRow.TableCell = append(tableRow.TableCell, ok)
	nokAsString := strconv.Itoa(record.CountNok)
	nok := TableCell{CellName: nokAsString}
	tableRow.TableCell = append(tableRow.TableCell, nok)
	note := TableCell{CellName: record.Note}
	tableRow.TableCell = append(tableRow.TableCell, note)
	data.TableRows = append(data.TableRows, tableRow)
}
