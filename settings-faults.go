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

type FaultSettingsDataOutput struct {
	DataTableSearchTitle    string
	DataTableInfoTitle      string
	DataTableRowsCountTitle string
	TableHeader             []HeaderCell
	TableRows               []TableRow
	TableHeaderType         []HeaderCellType
	TableRowsType           []TableRowType
}

type FaultDetailsDataOutput struct {
	FaultName            string
	FaultNamePrepend     string
	FaultTypeName        string
	FaultTypeNamePrepend string
	Barcode              string
	BarcodePrepend       string
	Color                string
	ColorPrepend         string
	Note                 string
	NotePrepend          string
	CreatedAt            string
	CreatedAtPrepend     string
	UpdatedAt            string
	UpdatedAtPrepend     string
	FaultTypes           []FaultTypeSelection
}

type FaultTypeSelection struct {
	FaultTypeName     string
	FaultTypeId       uint
	FaultTypeSelected string
}

type FaultTypeDetailsDataOutput struct {
	FaultTypeName        string
	FaultTypeNamePrepend string
	Note                 string
	NotePrepend          string
	CreatedAt            string
	CreatedAtPrepend     string
	UpdatedAt            string
	UpdatedAtPrepend     string
}

type FaultDetailsDataInput struct {
	Id      string
	Name    string
	Type    string
	Barcode string
	Note    string
}

type FaultTypeDetailsDataInput struct {
	Id   string
	Name string
	Note string
}

func loadFaults(writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS", "Loading faults")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Loading faults ended with error")
		return
	}
	var data FaultSettingsDataOutput
	data.DataTableSearchTitle = getLocale(email, "data-table-search-title")
	data.DataTableInfoTitle = getLocale(email, "data-table-info-title")
	data.DataTableRowsCountTitle = getLocale(email, "data-table-rows-count-title")
	var records []database.Fault
	db.Order("id desc").Find(&records)
	addFaultsTableHeaders(email, &data)
	for _, record := range records {
		addFaultsTableRow(record, &data)
	}
	var typeRecords []database.FaultType
	db.Order("id desc").Find(&typeRecords)
	addFaultTypesTableHeaders(email, &data)
	for _, record := range typeRecords {
		addFaultTypesTableRow(record, &data)
	}
	tmpl := template.Must(template.ParseFiles("./html/settings-table-type.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("SETTINGS", "Faults loaded in "+time.Since(timer).String())
}

func addFaultsTableRow(record database.Fault, data *FaultSettingsDataOutput) {
	var tableRow TableRow
	id := TableCell{CellName: strconv.Itoa(int(record.ID))}
	tableRow.TableCell = append(tableRow.TableCell, id)
	name := TableCell{CellName: record.Name}
	tableRow.TableCell = append(tableRow.TableCell, name)
	data.TableRows = append(data.TableRows, tableRow)
}

func addFaultsTableHeaders(email string, data *FaultSettingsDataOutput) {
	id := HeaderCell{HeaderName: "#", HeaderWidth: "30"}
	data.TableHeader = append(data.TableHeader, id)
	name := HeaderCell{HeaderName: getLocale(email, "fault-name")}
	data.TableHeader = append(data.TableHeader, name)
}

func addFaultTypesTableRow(record database.FaultType, data *FaultSettingsDataOutput) {
	var tableRow TableRowType
	id := TableCellType{CellNameType: strconv.Itoa(int(record.ID))}
	tableRow.TableCellType = append(tableRow.TableCellType, id)
	name := TableCellType{CellNameType: record.Name}
	tableRow.TableCellType = append(tableRow.TableCellType, name)
	data.TableRowsType = append(data.TableRowsType, tableRow)
}

func addFaultTypesTableHeaders(email string, data *FaultSettingsDataOutput) {
	id := HeaderCellType{HeaderNameType: "#", HeaderWidthType: "30"}
	data.TableHeaderType = append(data.TableHeaderType, id)
	name := HeaderCellType{HeaderNameType: getLocale(email, "type-name")}
	data.TableHeaderType = append(data.TableHeaderType, name)
}

func loadFault(id string, writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS", "Loading fault")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Loading fault ended with error")
		return
	}
	var fault database.Fault
	db.Where("id = ?", id).Find(&fault)
	var faultTypes []FaultTypeSelection
	for _, faultType := range cachedFaultTypesById {
		if faultType.Name == cachedFaultTypesById[uint(fault.FaultTypeID)].Name {
			faultTypes = append(faultTypes, FaultTypeSelection{FaultTypeName: faultType.Name, FaultTypeId: faultType.ID, FaultTypeSelected: "selected"})
		} else {
			faultTypes = append(faultTypes, FaultTypeSelection{FaultTypeName: faultType.Name, FaultTypeId: faultType.ID})
		}
	}
	sort.Slice(faultTypes, func(i, j int) bool {
		return faultTypes[i].FaultTypeName < faultTypes[j].FaultTypeName
	})
	data := FaultDetailsDataOutput{
		FaultName:            fault.Name,
		FaultNamePrepend:     getLocale(email, "fault-name"),
		FaultTypeName:        cachedFaultTypesById[uint(fault.FaultTypeID)].Name,
		FaultTypeNamePrepend: getLocale(email, "type-name"),
		Barcode:              fault.Barcode,
		BarcodePrepend:       getLocale(email, "barcode"),
		Note:                 fault.Note,
		NotePrepend:          getLocale(email, "note-name"),
		CreatedAt:            fault.CreatedAt.Format("2006-01-02T15:04:05"),
		CreatedAtPrepend:     getLocale(email, "created-at"),
		UpdatedAt:            fault.UpdatedAt.Format("2006-01-02T15:04:05"),
		UpdatedAtPrepend:     getLocale(email, "updated-at"),
		FaultTypes:           faultTypes,
	}
	tmpl := template.Must(template.ParseFiles("./html/settings-detail-fault.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("SETTINGS", "Fault "+fault.Name+" loaded in "+time.Since(timer).String())
}

func loadFaultType(id string, writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS", "Loading fault type")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Loading fault type ended with error")
		return
	}
	var faultType database.FaultType
	db.Where("id = ?", id).Find(&faultType)

	data := FaultTypeDetailsDataOutput{
		FaultTypeName:        faultType.Name,
		FaultTypeNamePrepend: getLocale(email, "type-name"),
		Note:                 faultType.Note,
		NotePrepend:          getLocale(email, "note-name"),
		CreatedAt:            faultType.CreatedAt.Format("2006-01-02T15:04:05"),
		CreatedAtPrepend:     getLocale(email, "created-at"),
		UpdatedAt:            faultType.UpdatedAt.Format("2006-01-02T15:04:05"),
		UpdatedAtPrepend:     getLocale(email, "updated-at"),
	}
	tmpl := template.Must(template.ParseFiles("./html/settings-detail-fault-type.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("SETTINGS", "Fault type "+faultType.Name+" loaded in "+time.Since(timer).String())
}

func saveFault(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	timer := time.Now()
	logInfo("SETTINGS", "Saving fault")
	var data FaultDetailsDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Saving fault ended with error")
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
		logInfo("SETTINGS", "Saving fault ended with error")
		return
	}
	var fault database.Fault
	db.Where("id=?", data.Id).Find(&fault)
	fault.Name = data.Name
	fault.FaultTypeID = int(cachedFaultTypesByName[data.Type].ID)
	fault.Barcode = data.Barcode
	fault.Note = data.Note
	db.Save(&fault)
	cacheFaults(db)
	logInfo("SETTINGS", "Fault "+fault.Name+" saved in "+time.Since(timer).String())
}

func saveFaultType(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	timer := time.Now()
	logInfo("SETTINGS", "Saving fault type")
	var data FaultTypeDetailsDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Saving fault type ended with error")
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
		logInfo("SETTINGS", "Saving fault type ended with error")
		return
	}
	var faultType database.FaultType
	db.Where("id=?", data.Id).Find(&faultType)
	faultType.Name = data.Name
	faultType.Note = data.Note
	db.Save(&faultType)
	cacheFaults(db)
	logInfo("SETTINGS", "Fault type "+faultType.Name+" saved in "+time.Since(timer).String())
}
