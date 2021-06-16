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

func loadPackagesTable(writer http.ResponseWriter, workplaceIds string, dateFrom time.Time, dateTo time.Time, email string) {
	timer := time.Now()
	logInfo("DATA-PACKAGES", "Loading packages table")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("DATA-PACKAGES", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "ERR: Problem opening database, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("DATA-PACKAGES", "Loading packages table ended")
		return
	}
	var packageRecords []database.PackageRecord
	if workplaceIds == "workplace_id in (')" {
		db.Where("date_time >= ?", dateFrom).Where("date_time <= ?", dateTo).Order("date_time desc").Find(&packageRecords)
	} else {
		db.Where("date_time >= ?", dateFrom).Where("date_time <= ?", dateTo).Where(workplaceIds).Order("date_time desc").Find(&packageRecords)
	}
	var data TableOutput
	data.Compacted = cachedUserWebSettings[email]["data-selected-size"]
	data.DataTableSearchTitle = getLocale(email, "data-table-search-title")
	data.DataTableInfoTitle = getLocale(email, "data-table-info-title")
	data.DataTableRowsCountTitle = getLocale(email, "data-table-rows-count-title")
	locationSync.RLock()
	loc, err := time.LoadLocation(cachedLocation)
	locationSync.RUnlock()
	addPackageTableHeaders(email, &data)
	for _, record := range packageRecords {
		addPackageTableRow(record, &data, db, loc)
	}
	tmpl, err := template.ParseFiles("./html/data-content.html")
	if err != nil {
		logError("SETTINGS", "Problem parsing html file: "+err.Error())
		var responseData TableOutput
		responseData.Result = "ERR: Problem parsing html file: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
	} else {
		data.Result = "INF: Packages data processed in " + time.Since(timer).String()
		_ = tmpl.Execute(writer, data)
		logInfo("SETTINGS", "Packages data loaded in "+time.Since(timer).String())
	}
}

func addPackageTableRow(record database.PackageRecord, data *TableOutput, db *gorm.DB, loc *time.Location) {
	var tableRow TableRow
	workplaceNameCell := TableCell{CellName: cachedWorkplacesById[uint(record.WorkplaceID)].Name}
	tableRow.TableCell = append(tableRow.TableCell, workplaceNameCell)
	packageDate := TableCell{CellName: record.DateTime.In(loc).Format("2006-01-02 15:04:05")}
	tableRow.TableCell = append(tableRow.TableCell, packageDate)
	userName := TableCell{CellName: cachedUsersById[uint(record.UserID)].FirstName + " " + cachedUsersById[uint(record.UserID)].SecondName}
	tableRow.TableCell = append(tableRow.TableCell, userName)
	var orderRecord database.OrderRecord
	db.Where("id = ?", record.OrderRecordID).Find(&orderRecord)
	ordersByIdSync.RLock()
	orderName := TableCell{CellName: cachedOrdersById[uint(orderRecord.OrderID)].Name}
	ordersByIdSync.RUnlock()
	tableRow.TableCell = append(tableRow.TableCell, orderName)
	packagesByIdSync.RLock()
	packageName := TableCell{CellName: cachedPackagesById[uint(record.PackageID)].Name}
	packagesByIdSync.RUnlock()
	tableRow.TableCell = append(tableRow.TableCell, packageName)
	countAsString := strconv.Itoa(record.Count)
	count := TableCell{CellName: countAsString}
	tableRow.TableCell = append(tableRow.TableCell, count)
	note := TableCell{CellName: record.Note}
	tableRow.TableCell = append(tableRow.TableCell, note)
	data.TableRows = append(data.TableRows, tableRow)
}

func addPackageTableHeaders(email string, data *TableOutput) {
	workplaceName := HeaderCell{HeaderName: getLocale(email, "workplace-name")}
	data.TableHeader = append(data.TableHeader, workplaceName)
	packageDate := HeaderCell{HeaderName: getLocale(email, "package-date")}
	data.TableHeader = append(data.TableHeader, packageDate)
	userName := HeaderCell{HeaderName: getLocale(email, "user-name")}
	data.TableHeader = append(data.TableHeader, userName)
	orderName := HeaderCell{HeaderName: getLocale(email, "order-name")}
	data.TableHeader = append(data.TableHeader, orderName)
	packageName := HeaderCell{HeaderName: getLocale(email, "package-name")}
	data.TableHeader = append(data.TableHeader, packageName)
	packageCount := HeaderCell{HeaderName: getLocale(email, "package-count")}
	data.TableHeader = append(data.TableHeader, packageCount)
	noteName := HeaderCell{HeaderName: getLocale(email, "note-name")}
	data.TableHeader = append(data.TableHeader, noteName)
}
