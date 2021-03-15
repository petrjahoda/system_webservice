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
}

type PartDetailsDataInput struct {
	Id      string
	Name    string
	Barcode string
	Note    string
}

func savePart(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	timer := time.Now()
	logInfo("SETTINGS-PARTS", "Saving part started")
	var data PartDetailsDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS-PARTS", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-PARTS", "Saving part ended")
		return
	}
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-PARTS", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-PARTS", "Saving part ended")
		return
	}

	var part database.Part
	db.Where("id=?", data.Id).Find(&part)
	part.Name = data.Name
	part.Barcode = data.Barcode
	part.Note = data.Note
	db.Debug().Save(&part)
	cacheParts(db)
	logInfo("SETTINGS-PARTS", "Part saved in "+time.Since(timer).String())
}

func loadPartDetails(id string, writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS-PARTSS", "Loading part details")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-PARTSS", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-PARTSS", "Loading part details ended")
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
	tmpl := template.Must(template.ParseFiles("./html/settings-detail-part.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("SETTINGS-PARTSS", "Part details loaded in "+time.Since(timer).String())
}

func loadPartsSettings(writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS-PARTSS", "Loading parts settings")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-PARTSS", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-PARTSS", "Loading parts settings ended")
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
	logInfo("SETTINGS-PARTSS", "Parts settings loaded in "+time.Since(timer).String())
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
