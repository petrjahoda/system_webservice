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

func processPackages(writer http.ResponseWriter, workplaceIds string, dateFrom time.Time, dateTo time.Time, email string) {
	timer := time.Now()
	logInfo("DATA-PACKAGES", "Processing packages started")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("DATA-PACKAGES", "Problem opening database: "+err.Error())
		var responseData DataPageOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("DATA-PACKAGES", "Processing data ended")
		return
	}
	var packageRecords []database.PackageRecord
	if workplaceIds == "workplace_id in (')" {
		db.Where("date_time >= ?", dateFrom).Where("date_time <= ?", dateTo).Order("date_time desc").Find(&packageRecords)
	} else {
		db.Where(workplaceIds).Where("date_time >= ?", dateFrom).Where("date_time <= ?", dateTo).Order("date_time desc").Find(&packageRecords)
	}
	var data TableData
	data.DataTableSearchTitle = getLocale(email, "data-table-search-title")
	data.DataTableInfoTitle = getLocale(email, "data-table-info-title")
	data.DataTableRowsCountTitle = getLocale(email, "data-table-rows-count-title")
	addPackageTableHeaders(email, &data)
	for _, record := range packageRecords {
		addPackageTableRow(record, &data, db)
	}
	tmpl := template.Must(template.ParseFiles("./html/table.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("DATA-PACKAGES", "Packages processed in "+time.Since(timer).String())
}

func addPackageTableRow(record database.PackageRecord, data *TableData, db *gorm.DB) {
	var tableRow TableRow
	workplaceNameCell := TableCell{CellName: cachedWorkplacesById[uint(record.WorkplaceID)].Name}
	tableRow.TableCell = append(tableRow.TableCell, workplaceNameCell)
	packageDate := TableCell{CellName: record.DateTime.Format("2006-01-02 15:04:05")}
	tableRow.TableCell = append(tableRow.TableCell, packageDate)
	userName := TableCell{CellName: cachedUsersById[uint(record.UserID)].FirstName + " " + cachedUsersById[uint(record.UserID)].SecondName}
	tableRow.TableCell = append(tableRow.TableCell, userName)
	var orderRecord database.OrderRecord
	db.Where("id = ?", record.OrderRecordID).Find(&orderRecord)
	orderName := TableCell{CellName: cachedOrdersById[uint(orderRecord.OrderID)].Name}
	tableRow.TableCell = append(tableRow.TableCell, orderName)
	packageName := TableCell{CellName: cachedPackagesById[uint(record.PackageID)].Name}
	tableRow.TableCell = append(tableRow.TableCell, packageName)
	countAsString := strconv.Itoa(record.Count)
	count := TableCell{CellName: countAsString}
	tableRow.TableCell = append(tableRow.TableCell, count)
	note := TableCell{CellName: record.Note}
	tableRow.TableCell = append(tableRow.TableCell, note)
	data.TableRows = append(data.TableRows, tableRow)
}

func addPackageTableHeaders(email string, data *TableData) {
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
