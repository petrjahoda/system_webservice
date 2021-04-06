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

func loadSystemStatsTable(writer http.ResponseWriter, dateFrom time.Time, dateTo time.Time, email string) {
	timer := time.Now()
	logInfo("DATA-SYSTEM-STATS", "Loading system stats table")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("DATA-SYSTEM-STATS", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("DATA-SYSTEM-STATS", "Loading system stats table ended")
		return
	}
	var systemRecords []database.SystemRecord
	db.Debug().Where("created_at >= ?", dateFrom).Where("created_at <= ?", dateTo).Order("created_at desc").Find(&systemRecords)
	var data TableOutput
	data.DataTableSearchTitle = getLocale(email, "data-table-search-title")
	data.DataTableInfoTitle = getLocale(email, "data-table-info-title")
	data.DataTableRowsCountTitle = getLocale(email, "data-table-rows-count-title")
	loc, err := time.LoadLocation(location)
	addSystemTableHeaders(email, &data)
	for _, record := range systemRecords {
		addSystemTableRow(record, &data, loc)
	}
	tmpl := template.Must(template.ParseFiles("./html/data-content.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("DATA-SYSTEM-STATS", "System stats table loaded in "+time.Since(timer).String())
}

func addSystemTableRow(record database.SystemRecord, data *TableOutput, loc *time.Location) {
	var tableRow TableRow
	systemDate := TableCell{CellName: record.CreatedAt.In(loc).Format("2006-01-02 15:04:05")}
	tableRow.TableCell = append(tableRow.TableCell, systemDate)

	databaseSizeAsString := strconv.FormatFloat(float64(record.DatabaseSizeInMegaBytes), 'f', 0, 64)
	databaseSize := TableCell{CellName: databaseSizeAsString + " MB"}
	tableRow.TableCell = append(tableRow.TableCell, databaseSize)

	databaseGrowthAsString := strconv.FormatFloat(float64(record.DatabaseGrowthInMegaBytes), 'f', 0, 64)
	databaseGrowth := TableCell{CellName: databaseGrowthAsString + " MB"}
	tableRow.TableCell = append(tableRow.TableCell, databaseGrowth)

	discFreeSizeAsString := strconv.FormatFloat(float64(record.DiscFreeSizeInMegaBytes), 'f', 0, 64)
	discFreeSize := TableCell{CellName: discFreeSizeAsString + " MB"}
	tableRow.TableCell = append(tableRow.TableCell, discFreeSize)

	estimatedFreeSizeAsString := strconv.FormatFloat(float64(record.EstimatedDiscFreeSizeInDays), 'f', 0, 64)
	estimatedFreeSize := TableCell{CellName: estimatedFreeSizeAsString}
	tableRow.TableCell = append(tableRow.TableCell, estimatedFreeSize)

	data.TableRows = append(data.TableRows, tableRow)
}

func addSystemTableHeaders(email string, data *TableOutput) {
	faultDate := HeaderCell{HeaderName: getLocale(email, "system-date")}
	data.TableHeader = append(data.TableHeader, faultDate)
	databaseSize := HeaderCell{HeaderName: getLocale(email, "database-size")}
	data.TableHeader = append(data.TableHeader, databaseSize)
	databaseGrowth := HeaderCell{HeaderName: getLocale(email, "database-growth")}
	data.TableHeader = append(data.TableHeader, databaseGrowth)
	discFreeSpace := HeaderCell{HeaderName: getLocale(email, "disc-free-space")}
	data.TableHeader = append(data.TableHeader, discFreeSpace)
	estimatedFreeSpace := HeaderCell{HeaderName: getLocale(email, "estimated-free-space")}
	data.TableHeader = append(data.TableHeader, estimatedFreeSpace)
}
