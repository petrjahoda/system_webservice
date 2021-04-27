package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/petrjahoda/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"html/template"
	"net/http"
	"sort"
	"strconv"
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
	Result                  string
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
	Result                  string
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
	Result                  string
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

func loadDowntimes(writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS", "Loading downtimes")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS", "Problem opening database: "+err.Error())
		var responseData DowntimesSettingsDataOutput
		responseData.Result = "ERR: Problem opening database, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Loading downtimes ended with error")
		return
	}
	var data DowntimesSettingsDataOutput
	data.DataTableSearchTitle = getLocale(email, "data-table-search-title")
	data.DataTableInfoTitle = getLocale(email, "data-table-info-title")
	data.DataTableRowsCountTitle = getLocale(email, "data-table-rows-count-title")
	var records []database.Downtime
	db.Order("id desc").Find(&records)
	addDowntimesTableHeaders(email, &data)
	for _, record := range records {
		addDowntimesTableRow(record, &data)
	}
	var typeRecords []database.DowntimeType
	db.Order("id desc").Find(&typeRecords)
	addDowntimeTypesTableHeaders(email, &data)
	for _, record := range typeRecords {
		addDowntimeTypesTableRow(record, &data)
	}
	tmpl, err := template.ParseFiles("./html/settings-table-type.html")
	if err != nil {
		logError("SETTINGS", "Problem parsing html file: "+err.Error())
		var responseData DowntimesSettingsDataOutput
		responseData.Result = "ERR: Problem parsing html file: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
	} else {
		data.Result = "INF: Downtimes processed in " + time.Since(timer).String()
		_ = tmpl.Execute(writer, data)
		logInfo("SETTINGS", "Downtimes loaded in "+time.Since(timer).String())
	}
}

func addDowntimesTableRow(record database.Downtime, data *DowntimesSettingsDataOutput) {
	var tableRow TableRow
	id := TableCell{CellName: strconv.Itoa(int(record.ID))}
	tableRow.TableCell = append(tableRow.TableCell, id)
	name := TableCell{CellName: record.Name}
	tableRow.TableCell = append(tableRow.TableCell, name)
	data.TableRows = append(data.TableRows, tableRow)
}

func addDowntimesTableHeaders(email string, data *DowntimesSettingsDataOutput) {
	id := HeaderCell{HeaderName: "#", HeaderWidth: "30"}
	data.TableHeader = append(data.TableHeader, id)
	name := HeaderCell{HeaderName: getLocale(email, "downtime-name")}
	data.TableHeader = append(data.TableHeader, name)
}

func addDowntimeTypesTableRow(record database.DowntimeType, data *DowntimesSettingsDataOutput) {
	var tableRow TableRowType
	id := TableCellType{CellNameType: strconv.Itoa(int(record.ID))}
	tableRow.TableCellType = append(tableRow.TableCellType, id)
	name := TableCellType{CellNameType: record.Name}
	tableRow.TableCellType = append(tableRow.TableCellType, name)
	data.TableRowsType = append(data.TableRowsType, tableRow)
}

func addDowntimeTypesTableHeaders(email string, data *DowntimesSettingsDataOutput) {
	id := HeaderCellType{HeaderNameType: "#", HeaderWidthType: "30"}
	data.TableHeaderType = append(data.TableHeaderType, id)
	name := HeaderCellType{HeaderNameType: getLocale(email, "type-name")}
	data.TableHeaderType = append(data.TableHeaderType, name)
}

func loadDowntime(id string, writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS", "Loading downtime")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS", "Problem opening database: "+err.Error())
		var responseData DowntimeDetailsDataOutput
		responseData.Result = "ERR: Problem opening database, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Loading downtime ended with error")
		return
	}
	var downtime database.Downtime
	db.Where("id = ?", id).Find(&downtime)
	var downTypes []DowntimeTypeSelection
	for _, downtimeType := range cachedDowntimeTypesById {
		if downtimeType.Name == cachedDowntimeTypesById[uint(downtime.DowntimeTypeID)].Name {
			downTypes = append(downTypes, DowntimeTypeSelection{DowntimeTypeName: downtimeType.Name, DowntimeTypeId: downtimeType.ID, DowntimeTypeSelected: "selected"})
		} else {
			downTypes = append(downTypes, DowntimeTypeSelection{DowntimeTypeName: downtimeType.Name, DowntimeTypeId: downtimeType.ID})
		}
	}
	sort.Slice(downTypes, func(i, j int) bool {
		return downTypes[i].DowntimeTypeName < downTypes[j].DowntimeTypeName
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
		DowntimeTypes:           downTypes,
	}
	tmpl, err := template.ParseFiles("./html/settings-detail-downtime.html")
	if err != nil {
		logError("SETTINGS", "Problem parsing html file: "+err.Error())
		var responseData DowntimeDetailsDataOutput
		responseData.Result = "ERR: Problem parsing html file: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
	} else {
		data.Result = "INF: Downtime detail processed in " + time.Since(timer).String()
		_ = tmpl.Execute(writer, data)
		logInfo("SETTINGS", "Downtime detail loaded in "+time.Since(timer).String())
	}
}

func loadDowntimeType(id string, writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS", "Loading downtime type")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS", "Problem opening database: "+err.Error())
		var responseData DowntimeTypeDetailsDataOutput
		responseData.Result = "ERR: Problem opening database, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Loading downtime type ended with error")
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
	tmpl, err := template.ParseFiles("./html/settings-detail-downtime-type.html")
	if err != nil {
		logError("SETTINGS", "Problem parsing html file: "+err.Error())
		var responseData DowntimeTypeDetailsDataOutput
		responseData.Result = "ERR: Problem parsing html file: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
	} else {
		data.Result = "INF: Downtime type detail processed in " + time.Since(timer).String()
		_ = tmpl.Execute(writer, data)
		logInfo("SETTINGS", "Downtime type detail loaded in "+time.Since(timer).String())
	}
}

func saveDowntime(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	timer := time.Now()
	logInfo("SETTINGS", "Saving downtime")
	var data DowntimeDetailsDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "ERR: Error parsing data, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Saving downtime ended with error")
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
		logInfo("SETTINGS", "Saving downtime ended with error")
		return
	}
	var downtime database.Downtime
	db.Where("id=?", data.Id).Find(&downtime)
	downtime.Name = data.Name
	downtime.DowntimeTypeID = int(cachedDowntimeTypesByName[data.Type].ID)
	downtime.Color = data.Color
	downtime.Barcode = data.Barcode
	downtime.Note = data.Note
	result := db.Save(&downtime)
	cacheDowntimes(db)
	if result.Error != nil {
		var responseData TableOutput
		responseData.Result = "ERR: Downtime not saved: " + result.Error.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logError("SETTINGS", "Downtime "+downtime.Name+" not saved: "+result.Error.Error())
	} else {
		var responseData TableOutput
		responseData.Result = "INF: Downtime saved in " + time.Since(timer).String()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Downtime "+downtime.Name+" saved in "+time.Since(timer).String())
	}
}

func saveDowntimeType(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	timer := time.Now()
	logInfo("SETTINGS", "Saving downtime type")
	var data DowntimeTypeDetailsDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "ERR: Error parsing data, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Saving downtime type ended with error")
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
		logInfo("SETTINGS", "Saving downtime ended with error")
		return
	}
	var downtimeType database.DowntimeType
	db.Where("id=?", data.Id).Find(&downtimeType)
	downtimeType.Name = data.Name
	downtimeType.Note = data.Note
	result := db.Save(&downtimeType)
	cacheDowntimes(db)
	if result.Error != nil {
		var responseData TableOutput
		responseData.Result = "ERR: Downtime Type not saved: " + result.Error.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logError("SETTINGS", "Downtime Type "+downtimeType.Name+" not saved: "+result.Error.Error())
	} else {
		var responseData TableOutput
		responseData.Result = "INF: Downtime Type saved in " + time.Since(timer).String()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Downtime Type "+downtimeType.Name+" saved in "+time.Since(timer).String())
	}
}
