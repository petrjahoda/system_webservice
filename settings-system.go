package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/petrjahoda/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

type SystemSettingsDataOutput struct {
	DataTableSearchTitle    string
	DataTableInfoTitle      string
	DataTableRowsCountTitle string
	TableHeader             []HeaderCell
	TableRows               []TableRow
}

type SystemSettingsDetailsDataOutput struct {
	SystemSettingsName         string
	SystemSettingsNamePrepend  string
	SystemSettingsSelection    []SystemSettingsSelection
	SystemSettingsValue        string
	SystemSettingsValuePrepend string
	Enabled                    string
	EnabledPrepend             string
	Note                       string
	NotePrepend                string
	CreatedAt                  string
	CreatedAtPrepend           string
	UpdatedAt                  string
	UpdatedAtPrepend           string
}
type SystemSettingsSelection struct {
	SystemSettingsValue    string
	SystemSettingsSelected string
}

type SystemSettingsDataInput struct {
	Id      string
	Name    string
	Value   string
	Enabled string
	Note    string
}

func loadSystemSettings(writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS", "Loading system settings")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Loading system settings ended with error")
		return
	}
	var data SystemSettingsDataOutput
	data.DataTableSearchTitle = getLocale(email, "data-table-search-title")
	data.DataTableInfoTitle = getLocale(email, "data-table-info-title")
	data.DataTableRowsCountTitle = getLocale(email, "data-table-rows-count-title")
	var records []database.Setting
	db.Order("id desc").Find(&records)
	addSystemSettingsTableHeaders(email, &data)
	for _, record := range records {
		addSystemSettingsTableRow(record, &data)
	}
	tmpl := template.Must(template.ParseFiles("./html/settings-table.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("SETTINGS", "System settings loaded in "+time.Since(timer).String())
}

func addSystemSettingsTableRow(record database.Setting, data *SystemSettingsDataOutput) {
	var tableRow TableRow
	id := TableCell{CellName: strconv.Itoa(int(record.ID))}
	tableRow.TableCell = append(tableRow.TableCell, id)
	name := TableCell{CellName: record.Name}
	tableRow.TableCell = append(tableRow.TableCell, name)
	data.TableRows = append(data.TableRows, tableRow)
}

func addSystemSettingsTableHeaders(email string, data *SystemSettingsDataOutput) {
	id := HeaderCell{HeaderName: "#", HeaderWidth: "30"}
	data.TableHeader = append(data.TableHeader, id)
	name := HeaderCell{HeaderName: getLocale(email, "name")}
	data.TableHeader = append(data.TableHeader, name)
}

func loadSystemSettingsDetails(id string, writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS", "Loading system settings details")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Loading system settings details ended with error")
		return
	}
	var settings database.Setting
	db.Where("id = ?", id).Find(&settings)
	var systemSettingsSelection []SystemSettingsSelection
	systemSettingsSelection = append(systemSettingsSelection, SystemSettingsSelection{SystemSettingsValue: "true", SystemSettingsSelected: checkSelection(settings.Enabled, "true")})
	systemSettingsSelection = append(systemSettingsSelection, SystemSettingsSelection{SystemSettingsValue: "false", SystemSettingsSelected: checkSelection(settings.Enabled, "false")})
	data := SystemSettingsDetailsDataOutput{
		SystemSettingsName:         settings.Name,
		SystemSettingsNamePrepend:  getLocale(email, "name"),
		SystemSettingsValue:        settings.Value,
		SystemSettingsValuePrepend: getLocale(email, "value"),
		Enabled:                    strconv.FormatBool(settings.Enabled),
		EnabledPrepend:             getLocale(email, "enabled"),
		Note:                       settings.Note,
		NotePrepend:                getLocale(email, "note-name"),
		CreatedAt:                  settings.CreatedAt.Format("2006-01-02T15:04:05"),
		CreatedAtPrepend:           getLocale(email, "created-at"),
		UpdatedAt:                  settings.UpdatedAt.Format("2006-01-02T15:04:05"),
		UpdatedAtPrepend:           getLocale(email, "updated-at"),
		SystemSettingsSelection:    systemSettingsSelection,
	}
	tmpl := template.Must(template.ParseFiles("./html/settings-detail-system.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("SETTINGS", "System settings details for "+settings.Name+" loaded in "+time.Since(timer).String())
}

func checkSelection(enabled bool, selection string) string {
	if strconv.FormatBool(enabled) == selection {
		return "selected"
	}
	return ""
}

func saveSystemSettingsDetails(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	timer := time.Now()
	logInfo("SETTINGS", "Saving system settings details")
	var data SystemSettingsDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Saving system settings details ended with error")
		return
	}
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Saving system settings details ended with error")
		return
	}
	var settings database.Setting
	db.Where("id=?", data.Id).Find(&settings)
	settings.Name = data.Name
	settings.Value = data.Value
	settings.Enabled, _ = strconv.ParseBool(data.Enabled)
	settings.Note = data.Note
	db.Save(&settings)
	cacheUsers(db)
	cacheSystemSettings(db)
	logInfo("SETTINGS", "System settings details for "+settings.Name+" saved in "+time.Since(timer).String())
}
