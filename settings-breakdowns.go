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

type BreakdownsSettingsDataOutput struct {
	DataTableSearchTitle    string
	DataTableInfoTitle      string
	DataTableRowsCountTitle string
	TableHeader             []HeaderCell
	TableRows               []TableRow
	TableHeaderType         []HeaderCellType
	TableRowsType           []TableRowType
	Result                  string
}

type BreakdownDetailsDataOutput struct {
	BreakdownName            string
	BreakdownNamePrepend     string
	BreakdownTypeName        string
	BreakdownTypeNamePrepend string
	Barcode                  string
	BarcodePrepend           string
	Color                    string
	ColorPrepend             string
	Note                     string
	NotePrepend              string
	CreatedAt                string
	CreatedAtPrepend         string
	UpdatedAt                string
	UpdatedAtPrepend         string
	BreakdownTypes           []BreakdownTypeSelection
	Result                   string
}

type BreakdownTypeSelection struct {
	BreakdownTypeName     string
	BreakdownTypeId       uint
	BreakdownTypeSelected string
}

type BreakdownTypeDetailsDataOutput struct {
	BreakdownTypeName        string
	BreakdownTypeNamePrepend string
	Note                     string
	NotePrepend              string
	CreatedAt                string
	CreatedAtPrepend         string
	UpdatedAt                string
	UpdatedAtPrepend         string
	Result                   string
}

type BreakdownDetailsDataInput struct {
	Id      string
	Name    string
	Type    string
	Barcode string
	Color   string
	Note    string
}

type BreakdownTypeDetailsDataInput struct {
	Id   string
	Name string
	Note string
}

func loadBreakdowns(writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS", "Loading breakdowns")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS", "Problem opening database: "+err.Error())
		var responseData BreakdownsSettingsDataOutput
		responseData.Result = "ERR: Problem opening database, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Loading breakdowns ended with error")
		return
	}
	var data BreakdownsSettingsDataOutput
	data.DataTableSearchTitle = getLocale(email, "data-table-search-title")
	data.DataTableInfoTitle = getLocale(email, "data-table-info-title")
	data.DataTableRowsCountTitle = getLocale(email, "data-table-rows-count-title")
	var records []database.Breakdown
	db.Order("id desc").Find(&records)
	addBreakdownsTableHeaders(email, &data)
	for _, record := range records {
		addBreakdownsTableRow(record, &data)
	}
	var typeRecords []database.BreakdownType
	db.Order("id desc").Find(&typeRecords)
	addBreakdownTypesTableHeaders(email, &data)
	for _, record := range typeRecords {
		addBreakdownTypesTableRow(record, &data)
	}
	data.Result = "INF: Breakdowns processed in " + time.Since(timer).String()

	tmpl, err := template.ParseFiles("./html/settings-table-type.html")
	if err != nil {
		logError("SETTINGS", "Problem parsing html file: "+err.Error())
		var responseData AlarmsSettingsDataOutput
		responseData.Result = "ERR: Problem parsing html file: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
	} else {
		data.Result = "INF: Breakdowns processed in " + time.Since(timer).String()
		_ = tmpl.Execute(writer, data)
		logInfo("SETTINGS", "Breakdowns loaded in "+time.Since(timer).String())
	}
}

func addBreakdownsTableRow(record database.Breakdown, data *BreakdownsSettingsDataOutput) {
	var tableRow TableRow
	id := TableCell{CellName: strconv.Itoa(int(record.ID))}
	tableRow.TableCell = append(tableRow.TableCell, id)
	name := TableCell{CellName: record.Name}
	tableRow.TableCell = append(tableRow.TableCell, name)
	data.TableRows = append(data.TableRows, tableRow)
}
func addBreakdownsTableHeaders(email string, data *BreakdownsSettingsDataOutput) {
	id := HeaderCell{HeaderName: "#", HeaderWidth: "30"}
	data.TableHeader = append(data.TableHeader, id)
	name := HeaderCell{HeaderName: getLocale(email, "breakdown-name")}
	data.TableHeader = append(data.TableHeader, name)
}

func addBreakdownTypesTableRow(record database.BreakdownType, data *BreakdownsSettingsDataOutput) {
	var tableRow TableRowType
	id := TableCellType{CellNameType: strconv.Itoa(int(record.ID))}
	tableRow.TableCellType = append(tableRow.TableCellType, id)
	name := TableCellType{CellNameType: record.Name}
	tableRow.TableCellType = append(tableRow.TableCellType, name)
	data.TableRowsType = append(data.TableRowsType, tableRow)
}

func addBreakdownTypesTableHeaders(email string, data *BreakdownsSettingsDataOutput) {
	id := HeaderCellType{HeaderNameType: "#", HeaderWidthType: "30"}
	data.TableHeaderType = append(data.TableHeaderType, id)
	name := HeaderCellType{HeaderNameType: getLocale(email, "type-name")}
	data.TableHeaderType = append(data.TableHeaderType, name)
}

func loadBreakdownTypes(id string, writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS", "Loading breakdown types")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS", "Problem opening database: "+err.Error())
		var responseData BreakdownTypeDetailsDataOutput
		responseData.Result = "ERR: Problem opening database, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Loading breakdown types ended with error")
		return
	}
	var breakdownType database.BreakdownType
	db.Where("id = ?", id).Find(&breakdownType)
	data := BreakdownTypeDetailsDataOutput{
		BreakdownTypeName:        breakdownType.Name,
		BreakdownTypeNamePrepend: getLocale(email, "type-name"),
		Note:                     breakdownType.Note,
		NotePrepend:              getLocale(email, "note-name"),
		CreatedAt:                breakdownType.CreatedAt.Format("2006-01-02T15:04:05"),
		CreatedAtPrepend:         getLocale(email, "created-at"),
		UpdatedAt:                breakdownType.UpdatedAt.Format("2006-01-02T15:04:05"),
		UpdatedAtPrepend:         getLocale(email, "updated-at"),
	}
	tmpl, err := template.ParseFiles("./html/settings-detail-breakdown-type.html")
	if err != nil {
		logError("SETTINGS", "Problem parsing html file: "+err.Error())
		var responseData BreakdownTypeDetailsDataOutput
		responseData.Result = "ERR: Problem parsing html file: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
	} else {
		data.Result = "INF: Breakdown type detail processed in " + time.Since(timer).String()
		_ = tmpl.Execute(writer, data)
		logInfo("SETTINGS", "Breakdown type detail loaded in "+time.Since(timer).String())
	}
}

func loadBreakdown(id string, writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS", "Loading breakdown")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS", "Problem opening database: "+err.Error())
		var responseData BreakdownDetailsDataOutput
		responseData.Result = "ERR: Problem opening database, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Loading breakdown ended with error")
		return
	}
	var breakdown database.Breakdown
	db.Where("id = ?", id).Find(&breakdown)
	var breakdownTypes []BreakdownTypeSelection
	for _, breakdownType := range cachedBreakdownTypesById {
		if breakdownType.Name == cachedBreakdownTypesById[uint(breakdown.BreakdownTypeID)].Name {
			breakdownTypes = append(breakdownTypes, BreakdownTypeSelection{BreakdownTypeName: breakdownType.Name, BreakdownTypeId: breakdownType.ID, BreakdownTypeSelected: "selected"})
		} else {
			breakdownTypes = append(breakdownTypes, BreakdownTypeSelection{BreakdownTypeName: breakdownType.Name, BreakdownTypeId: breakdownType.ID})
		}
	}
	sort.Slice(breakdownTypes, func(i, j int) bool {
		return breakdownTypes[i].BreakdownTypeName < breakdownTypes[j].BreakdownTypeName
	})
	data := BreakdownDetailsDataOutput{
		BreakdownName:            breakdown.Name,
		BreakdownNamePrepend:     getLocale(email, "breakdown-name"),
		BreakdownTypeName:        cachedBreakdownTypesById[uint(breakdown.BreakdownTypeID)].Name,
		BreakdownTypeNamePrepend: getLocale(email, "type-name"),
		Barcode:                  breakdown.Barcode,
		BarcodePrepend:           getLocale(email, "barcode"),
		Color:                    breakdown.Color,
		ColorPrepend:             getLocale(email, "color"),
		Note:                     breakdown.Note,
		NotePrepend:              getLocale(email, "note-name"),
		CreatedAt:                breakdown.CreatedAt.Format("2006-01-02T15:04:05"),
		CreatedAtPrepend:         getLocale(email, "created-at"),
		UpdatedAt:                breakdown.UpdatedAt.Format("2006-01-02T15:04:05"),
		UpdatedAtPrepend:         getLocale(email, "updated-at"),
		BreakdownTypes:           breakdownTypes,
	}
	tmpl, err := template.ParseFiles("./html/settings-detail-breakdown.html")
	if err != nil {
		logError("SETTINGS", "Problem parsing html file: "+err.Error())
		var responseData BreakdownDetailsDataOutput
		responseData.Result = "ERR: Problem parsing html file: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
	} else {
		data.Result = "INF: Breakdown detail processed in " + time.Since(timer).String()
		_ = tmpl.Execute(writer, data)
		logInfo("SETTINGS", "Breakdown detail loaded in "+time.Since(timer).String())
	}
}

func saveBreakdownType(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	timer := time.Now()
	logInfo("SETTINGS", "Saving breakdown type started")
	var data BreakdownTypeDetailsDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "ERR: Error parsing data, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Saving breakdown type ended with error")
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
		logInfo("SETTINGS", "Saving breakdown type ended with error")
		return
	}
	var breakdownType database.BreakdownType
	db.Where("id=?", data.Id).Find(&breakdownType)
	breakdownType.Name = data.Name
	breakdownType.Note = data.Note
	result := db.Save(&breakdownType)
	cacheBreakdowns(db)
	if result.Error != nil {
		var responseData TableOutput
		responseData.Result = "ERR: Breakdown type not saved: " + result.Error.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logError("SETTINGS", "Breakdown type "+breakdownType.Name+" not saved: "+result.Error.Error())
	} else {
		var responseData TableOutput
		responseData.Result = "INF: Breakdown type saved in " + time.Since(timer).String()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Breakdown type "+breakdownType.Name+" saved in "+time.Since(timer).String())
	}
}

func saveBreakdown(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	timer := time.Now()
	logInfo("SETTINGS", "Saving breakdown started")
	var data BreakdownDetailsDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "ERR: Error parsing data, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Saving breakdown ended with error")
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
		logInfo("SETTINGS", "Saving breakdown ended with error")
		return
	}
	var breakdown database.Breakdown
	db.Where("id=?", data.Id).Find(&breakdown)
	breakdown.Name = data.Name
	breakdown.BreakdownTypeID = int(cachedBreakdownTypesByName[data.Type].ID)
	breakdown.Color = data.Color
	breakdown.Barcode = data.Barcode
	breakdown.Note = data.Note
	result := db.Save(&breakdown)
	cacheBreakdowns(db)
	if result.Error != nil {
		var responseData TableOutput
		responseData.Result = "ERR: Breakdown not saved: " + result.Error.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logError("SETTINGS", "Breakdown "+breakdown.Name+" not saved: "+result.Error.Error())
	} else {
		var responseData TableOutput
		responseData.Result = "INF: Breakdown saved in " + time.Since(timer).String()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Breakdown "+breakdown.Name+" saved in "+time.Since(timer).String())
	}
}
