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

type StatesSettingsDataOutput struct {
	DataTableSearchTitle    string
	DataTableInfoTitle      string
	DataTableRowsCountTitle string
	TableHeader             []HeaderCell
	TableRows               []TableRow
	Result                  string
}

type StateDetailsDataOutput struct {
	StateName        string
	StateNamePrepend string
	Color            string
	ColorPrepend     string
	Note             string
	NotePrepend      string
	CreatedAt        string
	CreatedAtPrepend string
	UpdatedAt        string
	UpdatedAtPrepend string
	Result           string
}

type StateDetailsDataInput struct {
	Id    string
	Name  string
	Color string
	Note  string
}

func loadStates(writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS", "Loading states")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS", "Problem opening database: "+err.Error())
		var responseData StatesSettingsDataOutput
		responseData.Result = "ERR: Problem opening database, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Loading states ended with error")
		return
	}
	var records []database.State
	db.Order("id desc").Find(&records)
	var data StatesSettingsDataOutput
	data.DataTableSearchTitle = getLocale(email, "data-table-search-title")
	data.DataTableInfoTitle = getLocale(email, "data-table-info-title")
	data.DataTableRowsCountTitle = getLocale(email, "data-table-rows-count-title")
	addStatesTableHeaders(email, &data)
	for _, record := range records {
		addStatesTableRow(record, &data)
	}
	tmpl, err := template.ParseFiles("./html/settings-table.html")
	if err != nil {
		logError("SETTINGS", "Problem parsing html file: "+err.Error())
		var responseData OrdersSettingsDataOutput
		responseData.Result = "ERR: Problem parsing html file: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
	} else {
		data.Result = "INF: States processed in " + time.Since(timer).String()
		_ = tmpl.Execute(writer, data)
		logInfo("SETTINGS", "States loaded in "+time.Since(timer).String())
	}
}

func addStatesTableRow(record database.State, data *StatesSettingsDataOutput) {
	var tableRow TableRow
	id := TableCell{CellName: strconv.Itoa(int(record.ID))}
	tableRow.TableCell = append(tableRow.TableCell, id)
	name := TableCell{CellName: record.Name}
	tableRow.TableCell = append(tableRow.TableCell, name)
	data.TableRows = append(data.TableRows, tableRow)
}

func addStatesTableHeaders(email string, data *StatesSettingsDataOutput) {
	id := HeaderCell{HeaderName: "#", HeaderWidth: "30"}
	data.TableHeader = append(data.TableHeader, id)
	name := HeaderCell{HeaderName: getLocale(email, "state-name")}
	data.TableHeader = append(data.TableHeader, name)
}

func loadState(id string, writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS", "Loading state")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS", "Problem opening database: "+err.Error())
		var responseData StateDetailsDataOutput
		responseData.Result = "ERR: Problem opening database, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Loading state ended with error")
		return
	}
	stateId, _ := strconv.Atoi(id)
	state := cachedStatesById[uint(stateId)]
	data := StateDetailsDataOutput{
		StateName:        state.Name,
		StateNamePrepend: getLocale(email, "state-name"),
		Color:            state.Color,
		ColorPrepend:     getLocale(email, "color"),
		Note:             state.Note,
		NotePrepend:      getLocale(email, "note-name"),
		CreatedAt:        state.CreatedAt.Format("2006-01-02T15:04:05"),
		CreatedAtPrepend: getLocale(email, "created-at"),
		UpdatedAt:        state.UpdatedAt.Format("2006-01-02T15:04:05"),
		UpdatedAtPrepend: getLocale(email, "updated-at"),
	}
	tmpl, err := template.ParseFiles("./html/settings-detail-state.html")
	if err != nil {
		logError("SETTINGS", "Problem parsing html file: "+err.Error())
		var responseData StateDetailsDataOutput
		responseData.Result = "ERR: Problem parsing html file: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
	} else {
		data.Result = "INF: State detail processed in " + time.Since(timer).String()
		_ = tmpl.Execute(writer, data)
		logInfo("SETTINGS", "State detail loaded in "+time.Since(timer).String())
	}
}

func saveState(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	timer := time.Now()
	logInfo("SETTINGS", "Saving state")
	var data StateDetailsDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "ERR: Error parsing data, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Saving state ended with error")
		return
	}
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "ERR: Problem opening database, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Saving state ended with error")
		return
	}
	var state database.State
	db.Where("id=?", data.Id).Find(&state)
	state.Color = data.Color
	state.Name = data.Name
	state.Note = data.Note
	result := db.Save(&state)
	cacheStates(db)
	if result.Error != nil {
		var responseData TableOutput
		responseData.Result = "ERR: State not saved: " + result.Error.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logError("SETTINGS", "State "+state.Name+" not saved: "+result.Error.Error())
	} else {
		var responseData TableOutput
		responseData.Result = "INF: State saved in " + time.Since(timer).String()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "State "+state.Name+" saved in "+time.Since(timer).String())
	}
}
