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

type PartsSettingsDataOutput struct {
	DataTableSearchTitle    string
	DataTableInfoTitle      string
	DataTableRowsCountTitle string
	TableHeader             []HeaderCell
	TableRows               []TableRow
}

func loadPartsSettings(writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS-PARTS", "Loading parts settings")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-PARTS", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-PARTS", "Loading parts settings ended")
		return
	}
	var records []database.Part
	db.Order("id desc").Find(&records)
	var data PartsSettingsDataOutput
	data.DataTableSearchTitle = getLocale(email, "data-table-search-title")
	data.DataTableInfoTitle = getLocale(email, "data-table-info-title")
	data.DataTableRowsCountTitle = getLocale(email, "data-table-rows-count-title")
	addPartSettingsTableHeaders(email, &data)
	for _, record := range records {
		addPartSettingsTableRow(record, &data)
	}
	tmpl := template.Must(template.ParseFiles("./html/settings-table.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("SETTINGS-PARTS", "Parts settings loaded in "+time.Since(timer).String())
}

func addPartSettingsTableRow(record database.Part, data *PartsSettingsDataOutput) {
	var tableRow TableRow
	id := TableCell{CellName: strconv.Itoa(int(record.ID))}
	tableRow.TableCell = append(tableRow.TableCell, id)
	name := TableCell{CellName: record.Name}
	tableRow.TableCell = append(tableRow.TableCell, name)
	data.TableRows = append(data.TableRows, tableRow)
}

func addPartSettingsTableHeaders(email string, data *PartsSettingsDataOutput) {
	id := HeaderCell{HeaderName: "#", HeaderWidth: "30"}
	data.TableHeader = append(data.TableHeader, id)
	name := HeaderCell{HeaderName: getLocale(email, "part-name")}
	data.TableHeader = append(data.TableHeader, name)
}
