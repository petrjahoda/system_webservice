package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/petrjahoda/database"
	"gopkg.in/go-playground/colors.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type StatesSettingsDataOutput struct {
	DataTableSearchTitle    string
	DataTableInfoTitle      string
	DataTableRowsCountTitle string
	TableHeader             []HeaderCell
	TableRows               []TableRow
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
}

type StateDetailsDataInput struct {
	Id    string
	Name  string
	Color string
	Note  string
}

func saveState(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	timer := time.Now()
	logInfo("SETTINGS-STATES", "Saving state started")
	var data StateDetailsDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS-STATES", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-STATES", "Saving state ended")
		return
	}
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-STATES", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-STATES", "Saving state ended")
		return
	}
	var state database.State
	db.Where("id=?", data.Id).Find(&state)
	result := strings.TrimRight(data.Color, " none repeat scroll 0% 0% / auto padding-box border-box")
	rgb, err := colors.ParseRGB(result)
	if err != nil {
		logError("SETTINGS-STATES", "Problem parsing color: "+err.Error())
	} else {
		state.Color = rgb.ToHEX().String()
	}
	state.Name = data.Name
	state.Note = data.Note
	db.Save(&state)
	cacheStates(db)
	logInfo("SETTINGS-STATES", "State saved in "+time.Since(timer).String())
}

func loadStateDetails(id string, writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS-STATES", "Loading state details")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-STATES", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-STATES", "Loading state details ended")
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
	tmpl := template.Must(template.ParseFiles("./html/settings-detail-state.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("SETTINGS-STATES", "State details loaded in "+time.Since(timer).String())
}

func loadStatesSettings(writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS-STATES", "Loading states settings")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-STATES", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-STATES", "Loading states settings ended")
		return
	}
	var records []database.State
	db.Order("id desc").Find(&records)
	var data StatesSettingsDataOutput
	data.DataTableSearchTitle = getLocale(email, "data-table-search-title")
	data.DataTableInfoTitle = getLocale(email, "data-table-info-title")
	data.DataTableRowsCountTitle = getLocale(email, "data-table-rows-count-title")
	addStateSettingsTableHeaders(email, &data)
	for _, record := range records {
		addStateSettingsTableRow(record, &data)
	}
	tmpl := template.Must(template.ParseFiles("./html/settings-table.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("SETTINGS-STATES", "States settings loaded in "+time.Since(timer).String())
}

func addStateSettingsTableRow(record database.State, data *StatesSettingsDataOutput) {
	var tableRow TableRow
	id := TableCell{CellName: strconv.Itoa(int(record.ID))}
	tableRow.TableCell = append(tableRow.TableCell, id)
	name := TableCell{CellName: record.Name}
	tableRow.TableCell = append(tableRow.TableCell, name)
	data.TableRows = append(data.TableRows, tableRow)
}

func addStateSettingsTableHeaders(email string, data *StatesSettingsDataOutput) {
	id := HeaderCell{HeaderName: "#", HeaderWidth: "30"}
	data.TableHeader = append(data.TableHeader, id)
	name := HeaderCell{HeaderName: getLocale(email, "state-name")}
	data.TableHeader = append(data.TableHeader, name)
}
