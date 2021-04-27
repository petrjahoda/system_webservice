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

type PartsSettingsDataOutput struct {
	DataTableSearchTitle    string
	DataTableInfoTitle      string
	DataTableRowsCountTitle string
	TableHeader             []HeaderCell
	TableRows               []TableRow
	Result                  string
}

type PartDetailsDataOutput struct {
	PartName         string
	PartNamePrepend  string
	Barcode          string
	BarcodePrepend   string
	Note             string
	NotePrepend      string
	CreatedAt        string
	CreatedAtPrepend string
	UpdatedAt        string
	UpdatedAtPrepend string
	Result           string
}

type PartDetailsDataInput struct {
	Id      string
	Name    string
	Barcode string
	Note    string
}

func loadParts(writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS", "Loading parts")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS", "Problem opening database: "+err.Error())
		var responseData PartsSettingsDataOutput
		responseData.Result = "ERR: Problem opening database, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Loading parts ended with error")
		return
	}
	var records []database.Part
	db.Order("id desc").Find(&records)
	var data PartsSettingsDataOutput
	data.DataTableSearchTitle = getLocale(email, "data-table-search-title")
	data.DataTableInfoTitle = getLocale(email, "data-table-info-title")
	data.DataTableRowsCountTitle = getLocale(email, "data-table-rows-count-title")
	addPartsTableHeaders(email, &data)
	for _, record := range records {
		addPartsTableRow(record, &data)
	}
	tmpl, err := template.ParseFiles("./html/settings-table.html")
	if err != nil {
		logError("SETTINGS", "Problem parsing html file: "+err.Error())
		var responseData OrdersSettingsDataOutput
		responseData.Result = "ERR: Problem parsing html file: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
	} else {
		data.Result = "INF: Parts processed in " + time.Since(timer).String()
		_ = tmpl.Execute(writer, data)
		logInfo("SETTINGS", "Parts loaded in "+time.Since(timer).String())
	}
}

func addPartsTableRow(record database.Part, data *PartsSettingsDataOutput) {
	var tableRow TableRow
	id := TableCell{CellName: strconv.Itoa(int(record.ID))}
	tableRow.TableCell = append(tableRow.TableCell, id)
	name := TableCell{CellName: record.Name}
	tableRow.TableCell = append(tableRow.TableCell, name)
	data.TableRows = append(data.TableRows, tableRow)
}

func addPartsTableHeaders(email string, data *PartsSettingsDataOutput) {
	id := HeaderCell{HeaderName: "#", HeaderWidth: "30"}
	data.TableHeader = append(data.TableHeader, id)
	name := HeaderCell{HeaderName: getLocale(email, "part-name")}
	data.TableHeader = append(data.TableHeader, name)
}

func loadPart(id string, writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS", "Loading part")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS", "Problem opening database: "+err.Error())
		var responseData PartDetailsDataOutput
		responseData.Result = "ERR: Problem opening database, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Loading part ended with error")
		return
	}
	partId, _ := strconv.Atoi(id)
	part := cachedPartsById[uint(partId)]
	data := PartDetailsDataOutput{
		PartName:         part.Name,
		PartNamePrepend:  getLocale(email, "part-name"),
		Barcode:          part.Barcode,
		BarcodePrepend:   getLocale(email, "barcode"),
		Note:             part.Note,
		NotePrepend:      getLocale(email, "note-name"),
		CreatedAt:        part.CreatedAt.Format("2006-01-02T15:04:05"),
		CreatedAtPrepend: getLocale(email, "created-at"),
		UpdatedAt:        part.UpdatedAt.Format("2006-01-02T15:04:05"),
		UpdatedAtPrepend: getLocale(email, "updated-at"),
	}
	tmpl, err := template.ParseFiles("./html/settings-detail-part.html")
	if err != nil {
		logError("SETTINGS", "Problem parsing html file: "+err.Error())
		var responseData PartDetailsDataOutput
		responseData.Result = "ERR: Problem parsing html file: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
	} else {
		data.Result = "INF: Part detail processed in " + time.Since(timer).String()
		_ = tmpl.Execute(writer, data)
		logInfo("SETTINGS", "Part detail loaded in "+time.Since(timer).String())
	}
}

func savePart(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	timer := time.Now()
	logInfo("SETTINGS", "Saving part")
	var data PartDetailsDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "ERR: Error parsing data, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Saving part ended with error")
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
		logInfo("SETTINGS", "Saving part ended with error")
		return
	}

	var part database.Part
	db.Where("id=?", data.Id).Find(&part)
	part.Name = data.Name
	part.Barcode = data.Barcode
	part.Note = data.Note
	result := db.Save(&part)
	cacheParts(db)
	if result.Error != nil {
		var responseData TableOutput
		responseData.Result = "ERR: Part not saved: " + result.Error.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logError("SETTINGS", "Part "+part.Name+" not saved: "+result.Error.Error())
	} else {
		var responseData TableOutput
		responseData.Result = "INF: Part saved in " + time.Since(timer).String()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Part "+part.Name+" saved in "+time.Since(timer).String())
	}
}
