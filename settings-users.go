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

type UsersSettingsDataOutput struct {
	DataTableSearchTitle    string
	DataTableInfoTitle      string
	DataTableRowsCountTitle string
	TableHeader             []HeaderCell
	TableRows               []TableRow
	TableHeaderType         []HeaderCellType
	TableRowsType           []TableRowType
}

func loadUsersSettings(writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS-USERS", "Loading users settings")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-USERS", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-USERS", "Loading users settings ended")
		return
	}

	var data UsersSettingsDataOutput
	data.DataTableSearchTitle = getLocale(email, "data-table-search-title")
	data.DataTableInfoTitle = getLocale(email, "data-table-info-title")
	data.DataTableRowsCountTitle = getLocale(email, "data-table-rows-count-title")

	var records []database.User
	db.Order("id desc").Find(&records)
	addUserSettingsTableHeaders(email, &data)
	for _, record := range records {
		addUserSettingsTableRow(record, &data)
	}

	var typeRecords []database.UserType
	db.Order("id desc").Find(&typeRecords)
	addUserSettingsTypeTableHeaders(email, &data)
	for _, record := range typeRecords {
		addUserSettingsTypeTableRow(record, &data)
	}

	tmpl := template.Must(template.ParseFiles("./html/settings-table-type.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("SETTINGS-USERS", "Users settings loaded in "+time.Since(timer).String())
}

func addUserSettingsTableRow(record database.User, data *UsersSettingsDataOutput) {
	var tableRow TableRow
	id := TableCell{CellName: strconv.Itoa(int(record.ID))}
	tableRow.TableCell = append(tableRow.TableCell, id)
	name := TableCell{CellName: record.FirstName + " " + record.SecondName}
	tableRow.TableCell = append(tableRow.TableCell, name)
	data.TableRows = append(data.TableRows, tableRow)
}

func addUserSettingsTableHeaders(email string, data *UsersSettingsDataOutput) {
	id := HeaderCell{HeaderName: "#", HeaderWidth: "30"}
	data.TableHeader = append(data.TableHeader, id)
	name := HeaderCell{HeaderName: getLocale(email, "user-name")}
	data.TableHeader = append(data.TableHeader, name)
}

func addUserSettingsTypeTableRow(record database.UserType, data *UsersSettingsDataOutput) {
	var tableRow TableRowType
	id := TableCellType{CellNameType: strconv.Itoa(int(record.ID))}
	tableRow.TableCellType = append(tableRow.TableCellType, id)
	name := TableCellType{CellNameType: record.Name}
	tableRow.TableCellType = append(tableRow.TableCellType, name)
	data.TableRowsType = append(data.TableRowsType, tableRow)
}

func addUserSettingsTypeTableHeaders(email string, data *UsersSettingsDataOutput) {
	id := HeaderCellType{HeaderNameType: "#", HeaderWidthType: "30"}
	data.TableHeaderType = append(data.TableHeaderType, id)
	name := HeaderCellType{HeaderNameType: getLocale(email, "type-name")}
	data.TableHeaderType = append(data.TableHeaderType, name)
}
