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

type FaultsSettingsDataOutput struct {
	DataTableSearchTitle    string
	DataTableInfoTitle      string
	DataTableRowsCountTitle string
	TableHeader             []HeaderCell
	TableRows               []TableRow
	TableHeaderType         []HeaderCellType
	TableRowsType           []TableRowType
}

func loadFaultsSettings(writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS-FAULTS", "Loading faults settings")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-FAULTS", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-FAULTS", "Loading faults settings ended")
		return
	}
	var data FaultsSettingsDataOutput
	data.DataTableSearchTitle = getLocale(email, "data-table-search-title")
	data.DataTableInfoTitle = getLocale(email, "data-table-info-title")
	data.DataTableRowsCountTitle = getLocale(email, "data-table-rows-count-title")

	var records []database.Fault
	db.Order("id desc").Find(&records)
	addFaultSettingsTableHeaders(email, &data)
	for _, record := range records {
		addFaultSettingsTableRow(record, &data)
	}

	var typeRecords []database.FaultType
	db.Order("id desc").Find(&typeRecords)
	addFaultSettingsTypeTableHeaders(email, &data)
	for _, record := range typeRecords {
		addFaultSettingsTypeTableRow(record, &data)
	}

	tmpl := template.Must(template.ParseFiles("./html/settings-table-type.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("SETTINGS-FAULTS", "Faults settings loaded in "+time.Since(timer).String())
}

func addFaultSettingsTableRow(record database.Fault, data *FaultsSettingsDataOutput) {
	var tableRow TableRow
	id := TableCell{CellName: strconv.Itoa(int(record.ID))}
	tableRow.TableCell = append(tableRow.TableCell, id)
	name := TableCell{CellName: record.Name}
	tableRow.TableCell = append(tableRow.TableCell, name)
	data.TableRows = append(data.TableRows, tableRow)
}

func addFaultSettingsTableHeaders(email string, data *FaultsSettingsDataOutput) {
	id := HeaderCell{HeaderName: "#", HeaderWidth: "30"}
	data.TableHeader = append(data.TableHeader, id)
	name := HeaderCell{HeaderName: getLocale(email, "fault-name")}
	data.TableHeader = append(data.TableHeader, name)
}

func addFaultSettingsTypeTableRow(record database.FaultType, data *FaultsSettingsDataOutput) {
	var tableRow TableRowType
	id := TableCellType{CellNameType: strconv.Itoa(int(record.ID))}
	tableRow.TableCellType = append(tableRow.TableCellType, id)
	name := TableCellType{CellNameType: record.Name}
	tableRow.TableCellType = append(tableRow.TableCellType, name)
	data.TableRowsType = append(data.TableRowsType, tableRow)
}

func addFaultSettingsTypeTableHeaders(email string, data *FaultsSettingsDataOutput) {
	id := HeaderCellType{HeaderNameType: "#", HeaderWidthType: "30"}
	data.TableHeaderType = append(data.TableHeaderType, id)
	name := HeaderCellType{HeaderNameType: getLocale(email, "type-name")}
	data.TableHeaderType = append(data.TableHeaderType, name)
}
