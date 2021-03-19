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
	"sort"
	"strconv"
	"strings"
	"time"
)

type DowntimesSettingsDataOutput struct {
	DataTableSearchTitle    string
	DataTableInfoTitle      string
	DataTableRowsCountTitle string
	TableHeader             []HeaderCell
	TableRows               []TableRow
	TableHeaderType         []HeaderCellType
	TableRowsType           []TableRowType
}

type DowntimeDetailsDataOutput struct {
	DowntimeName            string
	DowntimeNamePrepend     string
	DowntimeTypeName        string
	DowntimeTypeNamePrepend string
	Barcode                 string
	BarcodePrepend          string
	Color                   string
	ColorPrepend            string
	Note                    string
	NotePrepend             string
	CreatedAt               string
	CreatedAtPrepend        string
	UpdatedAt               string
	UpdatedAtPrepend        string
	DowntimeTypes           []DowntimeTypeSelection
}

type DowntimeTypeSelection struct {
	DowntimeTypeName     string
	DowntimeTypeId       uint
	DowntimeTypeSelected string
}

type DowntimeTypeDetailsDataOutput struct {
	DowntimeTypeName        string
	DowntimeTypeNamePrepend string
	Note                    string
	NotePrepend             string
	CreatedAt               string
	CreatedAtPrepend        string
	UpdatedAt               string
	UpdatedAtPrepend        string
}

type DowntimeDetailsDataInput struct {
	Id      string
	Name    string
	Type    string
	Barcode string
	Color   string
	Note    string
}

type DowntimeTypeDetailsDataInput struct {
	Id   string
	Name string
	Note string
}

func saveDowntimeType(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	timer := time.Now()
	logInfo("SETTINGS-DOWNTIMES", "Saving downtime type started")
	var data DowntimeTypeDetailsDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS-DOWNTIMES", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-DOWNTIMES", "Saving downtime type ended")
		return
	}
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-DOWNTIMES", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-DOWNTIMES", "Saving downtime ended")
		return
	}
	var downtimeType database.DowntimeType
	db.Where("id=?", data.Id).Find(&downtimeType)
	downtimeType.Name = data.Name
	downtimeType.Note = data.Note
	db.Save(&downtimeType)
	cacheDowntimes(db)
	logInfo("SETTINGS-DOWNTIMES", "Downtime type saved in "+time.Since(timer).String())
}

func saveDowntime(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	timer := time.Now()
	logInfo("SETTINGS-DOWNTIMES", "Saving downtime started")
	var data DowntimeDetailsDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS-DOWNTIMES", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-DOWNTIMES", "Saving downtime ended")
		return
	}
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-DOWNTIMES", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-DOWNTIMES", "Saving downtime ended")
		return
	}
	var downtime database.Downtime
	db.Where("id=?", data.Id).Find(&downtime)
	downtime.Name = data.Name
	downtime.DowntimeTypeID = int(cachedDowntimeTypesByName[data.Type].ID)
	result := strings.TrimRight(data.Color, " none repeat scroll 0% 0% / auto padding-box border-box")
	rgb, err := colors.ParseRGB(result)
	if err != nil {
		logError("SETTINGS-DOWNTIMES", "Problem parsing color: "+err.Error())
	} else {
		downtime.Color = rgb.ToHEX().String()
	}
	downtime.Barcode = data.Barcode
	downtime.Note = data.Note
	db.Save(&downtime)
	cacheDowntimes(db)
	logInfo("SETTINGS-DOWNTIMES", "Downtime saved in "+time.Since(timer).String())
}

func loadDowntimeTypeDetails(id string, writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS-DOWNTIMES", "Loading downtime type details")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-DOWNTIMES", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-DOWNTIMES", "Loading downtime type details ended")
		return
	}
	var downtimeType database.DowntimeType
	db.Where("id = ?", id).Find(&downtimeType)

	data := DowntimeTypeDetailsDataOutput{
		DowntimeTypeName:        downtimeType.Name,
		DowntimeTypeNamePrepend: getLocale(email, "type-name"),
		Note:                    downtimeType.Note,
		NotePrepend:             getLocale(email, "note-name"),
		CreatedAt:               downtimeType.CreatedAt.Format("2006-01-02T15:04:05"),
		CreatedAtPrepend:        getLocale(email, "created-at"),
		UpdatedAt:               downtimeType.UpdatedAt.Format("2006-01-02T15:04:05"),
		UpdatedAtPrepend:        getLocale(email, "updated-at"),
	}
	tmpl := template.Must(template.ParseFiles("./html/settings-detail-downtime-type.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("SETTINGS-DOWNTIMES", "Downtime type details loaded in "+time.Since(timer).String())
}

func loadDowntimeDetails(id string, writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS-DOWNTIMES", "Loading downtime details")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-DOWNTIMES", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-DOWNTIMES", "Loading downtime details ended")
		return
	}
	var downtime database.Downtime
	db.Where("id = ?", id).Find(&downtime)
	var downtypes []DowntimeTypeSelection
	for _, downtimeType := range cachedDowntimeTypesById {
		if downtimeType.Name == cachedDowntimeTypesById[uint(downtime.DowntimeTypeID)].Name {
			downtypes = append(downtypes, DowntimeTypeSelection{DowntimeTypeName: downtimeType.Name, DowntimeTypeId: downtimeType.ID, DowntimeTypeSelected: "selected"})
		} else {
			downtypes = append(downtypes, DowntimeTypeSelection{DowntimeTypeName: downtimeType.Name, DowntimeTypeId: downtimeType.ID})
		}
	}
	sort.Slice(downtypes, func(i, j int) bool {
		return downtypes[i].DowntimeTypeName < downtypes[j].DowntimeTypeName
	})
	data := DowntimeDetailsDataOutput{
		DowntimeName:            downtime.Name,
		DowntimeNamePrepend:     getLocale(email, "downtime-name"),
		DowntimeTypeName:        cachedDowntimeTypesById[uint(downtime.DowntimeTypeID)].Name,
		DowntimeTypeNamePrepend: getLocale(email, "type-name"),
		Barcode:                 downtime.Barcode,
		BarcodePrepend:          getLocale(email, "barcode"),
		Color:                   downtime.Color,
		ColorPrepend:            getLocale(email, "color"),
		Note:                    downtime.Note,
		NotePrepend:             getLocale(email, "note-name"),
		CreatedAt:               downtime.CreatedAt.Format("2006-01-02T15:04:05"),
		CreatedAtPrepend:        getLocale(email, "created-at"),
		UpdatedAt:               downtime.UpdatedAt.Format("2006-01-02T15:04:05"),
		UpdatedAtPrepend:        getLocale(email, "updated-at"),
		DowntimeTypes:           downtypes,
	}
	tmpl := template.Must(template.ParseFiles("./html/settings-detail-downtime.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("SETTINGS-DOWNTIMES", "Downtimes details loaded in "+time.Since(timer).String())
}
func loadDowntimesSettings(writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS-DOWNTIMES", "Loading downtimes settings")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-DOWNTIMES", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-DOWNTIMES", "Loading downtimes settings ended")
		return
	}

	var data DowntimesSettingsDataOutput
	data.DataTableSearchTitle = getLocale(email, "data-table-search-title")
	data.DataTableInfoTitle = getLocale(email, "data-table-info-title")
	data.DataTableRowsCountTitle = getLocale(email, "data-table-rows-count-title")

	var records []database.Downtime
	db.Order("id desc").Find(&records)
	addDowntimeSettingsTableHeaders(email, &data)
	for _, record := range records {
		addDowntimeSettingsTableRow(record, &data)
	}

	var typeRecords []database.DowntimeType
	db.Order("id desc").Find(&typeRecords)
	addDowntimeSettingsTypeTableHeaders(email, &data)
	for _, record := range typeRecords {
		addDowntimeSettingsTypeTableRow(record, &data)
	}

	tmpl := template.Must(template.ParseFiles("./html/settings-table-type.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("SETTINGS-DOWNTIMES", "Downtimes settings loaded in "+time.Since(timer).String())
}

func addDowntimeSettingsTableRow(record database.Downtime, data *DowntimesSettingsDataOutput) {
	var tableRow TableRow
	id := TableCell{CellName: strconv.Itoa(int(record.ID))}
	tableRow.TableCell = append(tableRow.TableCell, id)
	name := TableCell{CellName: record.Name}
	tableRow.TableCell = append(tableRow.TableCell, name)
	data.TableRows = append(data.TableRows, tableRow)
}

func addDowntimeSettingsTableHeaders(email string, data *DowntimesSettingsDataOutput) {
	id := HeaderCell{HeaderName: "#", HeaderWidth: "30"}
	data.TableHeader = append(data.TableHeader, id)
	name := HeaderCell{HeaderName: getLocale(email, "downtime-name")}
	data.TableHeader = append(data.TableHeader, name)
}

func addDowntimeSettingsTypeTableRow(record database.DowntimeType, data *DowntimesSettingsDataOutput) {
	var tableRow TableRowType
	id := TableCellType{CellNameType: strconv.Itoa(int(record.ID))}
	tableRow.TableCellType = append(tableRow.TableCellType, id)
	name := TableCellType{CellNameType: record.Name}
	tableRow.TableCellType = append(tableRow.TableCellType, name)
	data.TableRowsType = append(data.TableRowsType, tableRow)
}

func addDowntimeSettingsTypeTableHeaders(email string, data *DowntimesSettingsDataOutput) {
	id := HeaderCellType{HeaderNameType: "#", HeaderWidthType: "30"}
	data.TableHeaderType = append(data.TableHeaderType, id)
	name := HeaderCellType{HeaderNameType: getLocale(email, "type-name")}
	data.TableHeaderType = append(data.TableHeaderType, name)
}
