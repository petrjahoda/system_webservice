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

type BreakdownsSettingsDataOutput struct {
	DataTableSearchTitle    string
	DataTableInfoTitle      string
	DataTableRowsCountTitle string
	TableHeader             []HeaderCell
	TableRows               []TableRow
	TableHeaderType         []HeaderCellType
	TableRowsType           []TableRowType
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

func saveBreakdownType(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	timer := time.Now()
	logInfo("SETTINGS-BREAKDOWNS", "Saving breakdown type started")
	var data BreakdownTypeDetailsDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS-BREAKDOWNS", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-BREAKDOWNS", "Saving breakdown type ended")
		return
	}
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-BREAKDOWNS", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-BREAKDOWNS", "Saving breakdown ended")
		return
	}
	var breakdownType database.BreakdownType
	db.Where("id=?", data.Id).Find(&breakdownType)
	breakdownType.Name = data.Name
	breakdownType.Note = data.Note
	db.Save(&breakdownType)
	cacheBreakdowns(db)
	logInfo("SETTINGS-BREAKDOWNS", "Breakdown type saved in "+time.Since(timer).String())
}

func saveBreakdown(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	timer := time.Now()
	logInfo("SETTINGS-BREAKDOWNS", "Saving breakdown started")
	var data BreakdownDetailsDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS-BREAKDOWNS", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-BREAKDOWNS", "Saving breakdown ended")
		return
	}
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-BREAKDOWNS", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-BREAKDOWNS", "Saving breakdown ended")
		return
	}
	var breakdown database.Breakdown
	db.Where("id=?", data.Id).Find(&breakdown)
	breakdown.Name = data.Name
	breakdown.BreakdownTypeID = int(cachedBreakdownTypesByName[data.Type].ID)
	result := strings.TrimRight(data.Color, " none repeat scroll 0% 0% / auto padding-box border-box")
	rgb, err := colors.ParseRGB(result)
	if err != nil {
		logError("SETTINGS-BREAKDOWNS", "Problem parsing color: "+err.Error())
	} else {
		breakdown.Color = rgb.ToHEX().String()
	}
	breakdown.Barcode = data.Barcode
	breakdown.Note = data.Note
	db.Save(&breakdown)
	cacheBreakdowns(db)
	logInfo("SETTINGS-BREAKDOWNS", "Breakdown saved in "+time.Since(timer).String())
}

func loadBreakdownTypeDetails(id string, writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS-BREAKDOWNS", "Loading breakdown type details")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-BREAKDOWNS", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-BREAKDOWNS", "Loading breakdown type details ended")
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
	tmpl := template.Must(template.ParseFiles("./html/settings-detail-breakdown-type.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("SETTINGS-BREAKDOWNS", "Breakdown type details loaded in "+time.Since(timer).String())
}

func loadBreakdownDetails(id string, writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS-BREAKDOWNS", "Loading breakdown details")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-BREAKDOWNS", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-BREAKDOWNS", "Loading breakdown details ended")
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
	tmpl := template.Must(template.ParseFiles("./html/settings-detail-breakdown.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("SETTINGS-BREAKDOWNS", "Breakdown details loaded in "+time.Since(timer).String())
}

func loadBreakdownsSettings(writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS-BREAKDOWNS", "Loading breakdowns settings")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-BREAKDOWNS", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-BREAKDOWNS", "Loading breakdowns settings ended")
		return
	}

	var data BreakdownsSettingsDataOutput
	data.DataTableSearchTitle = getLocale(email, "data-table-search-title")
	data.DataTableInfoTitle = getLocale(email, "data-table-info-title")
	data.DataTableRowsCountTitle = getLocale(email, "data-table-rows-count-title")

	var records []database.Breakdown
	db.Order("id desc").Find(&records)
	addBreakdownSettingsTableHeaders(email, &data)
	for _, record := range records {
		addBreakdownSettingsTableRow(record, &data)
	}

	var typeRecords []database.BreakdownType
	db.Order("id desc").Find(&typeRecords)
	addBreakdownSettingsTypeTableHeaders(email, &data)
	for _, record := range typeRecords {
		addBreakdownSettingsTypeTableRow(record, &data)
	}

	tmpl := template.Must(template.ParseFiles("./html/settings-table-type.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("SETTINGS-BREAKDOWNS", "Breakdowns settings loaded in "+time.Since(timer).String())
}

func addBreakdownSettingsTableRow(record database.Breakdown, data *BreakdownsSettingsDataOutput) {
	var tableRow TableRow
	id := TableCell{CellName: strconv.Itoa(int(record.ID))}
	tableRow.TableCell = append(tableRow.TableCell, id)
	name := TableCell{CellName: record.Name}
	tableRow.TableCell = append(tableRow.TableCell, name)
	data.TableRows = append(data.TableRows, tableRow)
}
func addBreakdownSettingsTableHeaders(email string, data *BreakdownsSettingsDataOutput) {
	id := HeaderCell{HeaderName: "#", HeaderWidth: "30"}
	data.TableHeader = append(data.TableHeader, id)
	name := HeaderCell{HeaderName: getLocale(email, "breakdown-name")}
	data.TableHeader = append(data.TableHeader, name)
}

func addBreakdownSettingsTypeTableRow(record database.BreakdownType, data *BreakdownsSettingsDataOutput) {
	var tableRow TableRowType
	id := TableCellType{CellNameType: strconv.Itoa(int(record.ID))}
	tableRow.TableCellType = append(tableRow.TableCellType, id)
	name := TableCellType{CellNameType: record.Name}
	tableRow.TableCellType = append(tableRow.TableCellType, name)
	data.TableRowsType = append(data.TableRowsType, tableRow)
}

func addBreakdownSettingsTypeTableHeaders(email string, data *BreakdownsSettingsDataOutput) {
	id := HeaderCellType{HeaderNameType: "#", HeaderWidthType: "30"}
	data.TableHeaderType = append(data.TableHeaderType, id)
	name := HeaderCellType{HeaderNameType: getLocale(email, "type-name")}
	data.TableHeaderType = append(data.TableHeaderType, name)
}
