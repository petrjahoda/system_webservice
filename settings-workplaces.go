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

type WorkplacesSettingsDataOutput struct {
	DataTableSearchTitle    string
	DataTableInfoTitle      string
	DataTableRowsCountTitle string
	TableHeader             []HeaderCell
	TableRows               []TableRow
	TableHeaderType         []HeaderCellType
	TableRowsType           []TableRowType
	TableHeaderTypeExtended []HeaderCellTypeExtended
	TableRowsTypeExtended   []TableRowTypeExtended
}

type WorkplaceModeDetailsDataOutput struct {
	WorkplaceModeName        string
	WorkplaceModeNamePrepend string
	DowntimeDuration         string
	DowntimeDurationPrepend  string
	PoweroffDuration         string
	PoweroffDurationPrepend  string
	Note                     string
	NotePrepend              string
	CreatedAt                string
	CreatedAtPrepend         string
	UpdatedAt                string
	UpdatedAtPrepend         string
}

type WorkplaceSectionDetailsDataOutput struct {
	WorkplaceSectionName        string
	WorkplaceSectionNamePrepend string
	Note                        string
	NotePrepend                 string
	CreatedAt                   string
	CreatedAtPrepend            string
	UpdatedAt                   string
	UpdatedAtPrepend            string
}

type WorkplaceDetailsDataOutput struct {
	WorkplaceName           string
	WorkplaceNamePrepend    string
	WorkplaceSectionPrepend string
	WorkplaceModePrepend    string
	WorkplaceCode           string
	WorkplaceCodePrepend    string
	Note                    string
	NotePrepend             string
	CreatedAt               string
	CreatedAtPrepend        string
	UpdatedAt               string
	UpdatedAtPrepend        string
	WorkplaceSections       []WorkplaceSectionSelection
	WorkplaceModes          []WorkplaceModeSelection
	DataTableSearchTitle    string
	DataTableInfoTitle      string
	DataTableRowsCountTitle string
	TableHeader             []HeaderCell
	TableRows               []TableRow
	WorkshiftTableHeader    []WorkshiftHeaderCell
	WorkshiftTableRows      []WorkshiftTableRow
}

type WorkplaceSectionSelection struct {
	WorkplaceSectionName     string
	WorkplaceSectionId       uint
	WorkplaceSectionSelected string
}

type WorkplaceModeSelection struct {
	WorkplaceModeName     string
	WorkplaceModeId       uint
	WorkplaceModeSelected string
}

type WorkplaceModeDetailsDataInput struct {
	Id               string
	Name             string
	DowntimeDuration string
	PoweroffDuration string
	Note             string
}

type WorkplaceSectionDetailsDataInput struct {
	Id   string
	Name string
	Note string
}

type WorkplaceDetailsDataInput struct {
	Id      string
	Name    string
	Section string
	Mode    string
	Code    string
	Note    string
}

type WorkplacePortDetailsPageInput struct {
	Data string
	Type string
}

type WorkplacePortDetailsDataOutput struct {
	WorkplacePortName              string
	WorkplacePortNamePrepend       string
	WorkplacePortDevicePortPrepend string
	WorkplacePortStatePrepend      string
	Color                          string
	CounterOkEnabledPrepend        string
	CounterNokEnabledPrepend       string
	HighValue                      float32
	HighValuePrepend               string
	LowValue                       float32
	LowValuePrepend                string
	Note                           string
	NotePrepend                    string
	CreatedAt                      string
	CreatedAtPrepend               string
	UpdatedAt                      string
	UpdatedAtPrepend               string
	WorkplacePortDevicePorts       []WorkplacePortDevicePortSelection
	WorkplacePortStates            []WorkplacePortStateSelection
	CounterOkSelection             []CounterOkSelection
	CounterNokSelection            []CounterNokSelection
}

type WorkplacePortDevicePortSelection struct {
	WorkplacePortDevicePortName     string
	WorkplacePortDevicePortId       uint
	WorkplacePortDevicePortSelected string
}

type WorkplacePortStateSelection struct {
	WorkplacePortStateName     string
	WorkplacePortStateId       uint
	WorkplacePortStateSelected string
}

type CounterOkSelection struct {
	CounterOk         string
	CounterOkSelected string
}

type CounterNokSelection struct {
	CounterNok         string
	CounterNokSelected string
}

func loadWorkplacePortDetail(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	timer := time.Now()
	logInfo("SETTINGS-WORKPLACES", "Loading workplace port settings")
	email, _, _ := request.BasicAuth()
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-WORKPLACES", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-WORKPLACES", "Loading workplace port settings ended")
		return
	}
	var data WorkplacePortDetailsPageInput
	err = json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS-WORKPLACES", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-WORKPLACES", "Loading workplace port ended")
		return
	}
	var workplacePort database.WorkplacePort
	db.Where("id=?", data.Data).Find(&workplacePort)

	var counterOkSelections []CounterOkSelection
	counterOkSelections = append(counterOkSelections, CounterOkSelection{CounterOk: "true", CounterOkSelected: checkSelection(workplacePort.CounterOK, "true")})
	counterOkSelections = append(counterOkSelections, CounterOkSelection{CounterOk: "false", CounterOkSelected: checkSelection(workplacePort.CounterOK, "false")})

	var counterNokSelections []CounterNokSelection
	counterNokSelections = append(counterNokSelections, CounterNokSelection{CounterNok: "true", CounterNokSelected: checkSelection(workplacePort.CounterOK, "true")})
	counterNokSelections = append(counterNokSelections, CounterNokSelection{CounterNok: "false", CounterNokSelected: checkSelection(workplacePort.CounterOK, "false")})

	var workplacePortDevicePorts []WorkplacePortDevicePortSelection
	for _, devicePort := range cachedDevicePortsById {
		if devicePort.ID == cachedDevicePortsById[uint(workplacePort.DevicePortID)].ID {
			workplacePortDevicePorts = append(workplacePortDevicePorts, WorkplacePortDevicePortSelection{WorkplacePortDevicePortName: cachedDevicesById[uint(devicePort.DeviceID)].Name + ": [" + strconv.Itoa(devicePort.PortNumber) + "]" + devicePort.Name, WorkplacePortDevicePortId: devicePort.ID, WorkplacePortDevicePortSelected: "selected"})
		} else {
			workplacePortDevicePorts = append(workplacePortDevicePorts, WorkplacePortDevicePortSelection{WorkplacePortDevicePortName: cachedDevicesById[uint(devicePort.DeviceID)].Name + ": [" + strconv.Itoa(devicePort.PortNumber) + "]" + devicePort.Name, WorkplacePortDevicePortId: devicePort.ID})
		}
	}
	sort.Slice(workplacePortDevicePorts, func(i, j int) bool {
		return workplacePortDevicePorts[i].WorkplacePortDevicePortName < workplacePortDevicePorts[j].WorkplacePortDevicePortName
	})

	var states []WorkplacePortStateSelection
	for _, state := range cachedStatesById {
		if state.ID == cachedStatesById[uint(workplacePort.StateID)].ID {
			states = append(states, WorkplacePortStateSelection{WorkplacePortStateName: state.Name, WorkplacePortStateId: state.ID, WorkplacePortStateSelected: "selected"})
		} else {
			states = append(states, WorkplacePortStateSelection{WorkplacePortStateName: state.Name, WorkplacePortStateId: state.ID})
		}
	}
	sort.Slice(states, func(i, j int) bool {
		return states[i].WorkplacePortStateName < states[j].WorkplacePortStateName
	})

	dataOut := WorkplacePortDetailsDataOutput{
		WorkplacePortName:              workplacePort.Name,
		WorkplacePortNamePrepend:       getLocale(email, "port-name"),
		WorkplacePortDevicePortPrepend: getLocale(email, "device-name"),
		WorkplacePortStatePrepend:      getLocale(email, "state-name"),
		Color:                          workplacePort.Color,
		CounterOkEnabledPrepend:        getLocale(email, "counter-ok"),
		CounterNokEnabledPrepend:       getLocale(email, "counter-nok"),
		HighValue:                      workplacePort.HighValue,
		HighValuePrepend:               getLocale(email, "high-value"),
		LowValue:                       workplacePort.LowValue,
		LowValuePrepend:                getLocale(email, "low-value"),
		Note:                           workplacePort.Note,
		NotePrepend:                    getLocale(email, "note-name"),
		CreatedAt:                      workplacePort.CreatedAt.Format("2006-01-02T15:04:05"),
		CreatedAtPrepend:               getLocale(email, "created-at"),
		UpdatedAt:                      workplacePort.UpdatedAt.Format("2006-01-02T15:04:05"),
		UpdatedAtPrepend:               getLocale(email, "updated-at"),
		CounterOkSelection:             counterOkSelections,
		CounterNokSelection:            counterNokSelections,
		WorkplacePortDevicePorts:       workplacePortDevicePorts,
		WorkplacePortStates:            states,
	}
	tmpl := template.Must(template.ParseFiles("./html/settings-detail-workplace-port.html"))
	_ = tmpl.Execute(writer, dataOut)
	logInfo("SETTINGS-WORKPLACES", "Workplace port settings loaded in "+time.Since(timer).String())
}

func saveWorkplace(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	timer := time.Now()
	logInfo("SETTINGS-WORKPLACES", "Saving workplace started")
	var data WorkplaceDetailsDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS-WORKPLACES", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-WORKPLACES", "Saving workplace ended")
		return
	}
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-WORKPLACES", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-WORKPLACES", "Saving workplace ended")
		return
	}

	var workplace database.Workplace
	db.Where("id=?", data.Id).Find(&workplace)
	workplace.Name = data.Name
	workplace.WorkplaceModeID = int(cachedWorkplaceModesByName[data.Mode].ID)
	workplace.WorkplaceSectionID = int(cachedWorkplaceSectionsByName[data.Section].ID)
	workplace.Code = data.Code
	workplace.Note = data.Note
	db.Debug().Save(&workplace)
	cacheWorkplaces(db)
	logInfo("SETTINGS-WORKPLACES", "Workplace saved in "+time.Since(timer).String())
}

func saveWorkplaceSection(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	timer := time.Now()
	logInfo("SETTINGS-WORKPLACES", "Saving workplace section started")
	var data WorkplaceSectionDetailsDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS-WORKPLACES", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-WORKPLACES", "Saving workplace section ended")
		return
	}
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-WORKPLACES", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-WORKPLACES", "Saving workplace section ended")
		return
	}

	var workplaceSection database.WorkplaceSection
	db.Where("id=?", data.Id).Find(&workplaceSection)
	workplaceSection.Name = data.Name
	workplaceSection.Note = data.Note
	db.Save(&workplaceSection)
	cacheWorkplaces(db)
	logInfo("SETTINGS-WORKPLACES", "Workplace section saved in "+time.Since(timer).String())
}

func saveWorkplaceMode(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	timer := time.Now()
	logInfo("SETTINGS-WORKPLACES", "Saving workplace mode started")
	var data WorkplaceModeDetailsDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS-WORKPLACES", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-WORKPLACES", "Saving workplace mode ended")
		return
	}
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-WORKPLACES", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-WORKPLACES", "Saving workplace mode ended")
		return
	}
	downtimeParsed, err := time.ParseDuration(data.DowntimeDuration)
	if err != nil {
		logError("SETTINGS-PRODUCTS", "Problem parsing downtime duration: "+err.Error())
		downtimeParsed = 0
	}
	poweroffParsed, err := time.ParseDuration(data.PoweroffDuration)
	if err != nil {
		logError("SETTINGS-PRODUCTS", "Problem parsing poweroff duration: "+err.Error())
		poweroffParsed = 0
	}
	var workplaceMode database.WorkplaceMode
	db.Where("id=?", data.Id).Find(&workplaceMode)
	workplaceMode.Name = data.Name
	workplaceMode.PoweroffDuration = poweroffParsed
	workplaceMode.DowntimeDuration = downtimeParsed
	workplaceMode.Note = data.Note
	db.Save(&workplaceMode)
	cacheWorkplaces(db)
	logInfo("SETTINGS-WORKPLACES", "Workplace mode saved in "+time.Since(timer).String())
}

func loadWorkplaceDetails(id string, writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS-WORKPLACES", "Loading workplace details")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-WORKPLACES", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-WORKPLACES", "Loading workplace details ended")
		return
	}
	var workplace database.Workplace
	db.Where("id = ?", id).Find(&workplace)

	var workplaceSections []WorkplaceSectionSelection
	for _, workplaceSection := range cachedWorkplaceSectionsById {
		if workplaceSection.Name == cachedWorkplaceSectionsById[uint(workplace.WorkplaceSectionID)].Name {
			workplaceSections = append(workplaceSections, WorkplaceSectionSelection{WorkplaceSectionName: workplaceSection.Name, WorkplaceSectionId: workplaceSection.ID, WorkplaceSectionSelected: "selected"})
		} else {
			workplaceSections = append(workplaceSections, WorkplaceSectionSelection{WorkplaceSectionName: workplaceSection.Name, WorkplaceSectionId: workplaceSection.ID})
		}
	}
	sort.Slice(workplaceSections, func(i, j int) bool {
		return workplaceSections[i].WorkplaceSectionName < workplaceSections[j].WorkplaceSectionName
	})

	var workplaceModes []WorkplaceModeSelection
	for _, workplaceMode := range cachedWorkplaceModesById {
		if workplaceMode.Name == cachedWorkplaceModesById[uint(workplace.WorkplaceModeID)].Name {
			workplaceModes = append(workplaceModes, WorkplaceModeSelection{WorkplaceModeName: workplaceMode.Name, WorkplaceModeId: workplaceMode.ID, WorkplaceModeSelected: "selected"})
		} else {
			workplaceModes = append(workplaceModes, WorkplaceModeSelection{WorkplaceModeName: workplaceMode.Name, WorkplaceModeId: workplaceMode.ID})
		}
	}
	sort.Slice(workplaceModes, func(i, j int) bool {
		return workplaceModes[i].WorkplaceModeName < workplaceModes[j].WorkplaceModeName
	})

	var records []database.WorkplacePort
	db.Where("workplace_id = ?", workplace.ID).Order("id desc").Find(&records)
	var workshifts []database.WorkplaceWorkshift
	db.Where("workplace_id = ?", workplace.ID).Order("id desc").Find(&workshifts)

	data := WorkplaceDetailsDataOutput{
		WorkplaceName:           workplace.Name,
		WorkplaceNamePrepend:    getLocale(email, "workplace-name"),
		WorkplaceSectionPrepend: getLocale(email, "section-name"),
		WorkplaceModePrepend:    getLocale(email, "type-name"),
		WorkplaceCode:           workplace.Code,
		WorkplaceCodePrepend:    getLocale(email, "code-name"),
		Note:                    workplace.Note,
		NotePrepend:             getLocale(email, "note-name"),
		CreatedAt:               workplace.CreatedAt.Format("2006-01-02T15:04:05"),
		CreatedAtPrepend:        getLocale(email, "created-at"),
		UpdatedAt:               workplace.UpdatedAt.Format("2006-01-02T15:04:05"),
		UpdatedAtPrepend:        getLocale(email, "updated-at"),
		WorkplaceSections:       workplaceSections,
		WorkplaceModes:          workplaceModes,
	}
	addWorkplacePortDetailsTableHeaders(email, &data)
	for _, record := range records {
		addWorkplacePortDetailsTableRow(record, &data)
	}

	addWorkplaceWorkshiftDetailsTableHeaders(email, &data)

	for _, workshift := range workshifts {
		addWorkplaceWorkshiftDetailsTableRow(workshift, &data)
	}
	tmpl := template.Must(template.ParseFiles("./html/settings-detail-workplace.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("SETTINGS-WORKPLACES", "Workplace details loaded in "+time.Since(timer).String())
}

func addWorkplaceWorkshiftDetailsTableRow(record database.WorkplaceWorkshift, data *WorkplaceDetailsDataOutput) {
	var tableRow WorkshiftTableRow
	id := WorkshiftTableCell{WorkshiftCellName: strconv.Itoa(int(record.ID))}
	tableRow.WorkshiftTableCell = append(tableRow.WorkshiftTableCell, id)
	name := WorkshiftTableCell{WorkshiftCellName: cachedWorkshiftsById[uint(record.WorkshiftID)].Name}
	tableRow.WorkshiftTableCell = append(tableRow.WorkshiftTableCell, name)
	data.WorkshiftTableRows = append(data.WorkshiftTableRows, tableRow)
}

func addWorkplaceWorkshiftDetailsTableHeaders(email string, data *WorkplaceDetailsDataOutput) {
	id := WorkshiftHeaderCell{WorkshiftHeaderName: "#", WorkshiftHeaderWidth: "30"}
	data.WorkshiftTableHeader = append(data.WorkshiftTableHeader, id)
	name := WorkshiftHeaderCell{WorkshiftHeaderName: getLocale(email, "workshift-name")}
	data.WorkshiftTableHeader = append(data.WorkshiftTableHeader, name)
}

func addWorkplacePortDetailsTableRow(record database.WorkplacePort, data *WorkplaceDetailsDataOutput) {
	var tableRow TableRow
	id := TableCell{CellName: strconv.Itoa(int(record.ID))}
	tableRow.TableCell = append(tableRow.TableCell, id)
	name := TableCell{CellName: record.Name}
	tableRow.TableCell = append(tableRow.TableCell, name)
	data.TableRows = append(data.TableRows, tableRow)
}

func addWorkplacePortDetailsTableHeaders(email string, data *WorkplaceDetailsDataOutput) {
	id := HeaderCell{HeaderName: "#", HeaderWidth: "30"}
	data.TableHeader = append(data.TableHeader, id)
	name := HeaderCell{HeaderName: getLocale(email, "port-name")}
	data.TableHeader = append(data.TableHeader, name)
}

func loadWorkplaceSectionDetails(id string, writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS-WORKPLACES", "Loading workplace section details")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-WORKPLACES", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-USERS", "Loading workplace section details ended")
		return
	}
	var workplaceSection database.WorkplaceSection
	db.Where("id = ?", id).Find(&workplaceSection)

	data := WorkplaceSectionDetailsDataOutput{
		WorkplaceSectionName:        workplaceSection.Name,
		WorkplaceSectionNamePrepend: getLocale(email, "section-name"),
		Note:                        workplaceSection.Note,
		NotePrepend:                 getLocale(email, "note-name"),
		CreatedAt:                   workplaceSection.CreatedAt.Format("2006-01-02T15:04:05"),
		CreatedAtPrepend:            getLocale(email, "created-at"),
		UpdatedAt:                   workplaceSection.UpdatedAt.Format("2006-01-02T15:04:05"),
		UpdatedAtPrepend:            getLocale(email, "updated-at"),
	}
	tmpl := template.Must(template.ParseFiles("./html/settings-detail-workplace-section.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("SETTINGS-WORKPLACES", "Workplace mode section loaded in "+time.Since(timer).String())
}

func loadWorkplaceModeDetails(id string, writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS-WORKPLACES", "Loading workplace mode details")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-WORKPLACES", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-USERS", "Loading workplace mode details ended")
		return
	}
	var workplaceMode database.WorkplaceMode
	db.Where("id = ?", id).Find(&workplaceMode)

	data := WorkplaceModeDetailsDataOutput{
		WorkplaceModeName:        workplaceMode.Name,
		WorkplaceModeNamePrepend: getLocale(email, "type-name"),
		DowntimeDuration:         workplaceMode.DowntimeDuration.String(),
		DowntimeDurationPrepend:  getLocale(email, "downtime-duration"),
		PoweroffDuration:         workplaceMode.PoweroffDuration.String(),
		PoweroffDurationPrepend:  getLocale(email, "poweroff-duration"),
		Note:                     workplaceMode.Note,
		NotePrepend:              getLocale(email, "note-name"),
		CreatedAt:                workplaceMode.CreatedAt.Format("2006-01-02T15:04:05"),
		CreatedAtPrepend:         getLocale(email, "created-at"),
		UpdatedAt:                workplaceMode.UpdatedAt.Format("2006-01-02T15:04:05"),
		UpdatedAtPrepend:         getLocale(email, "updated-at"),
	}
	tmpl := template.Must(template.ParseFiles("./html/settings-detail-workplace-mode.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("SETTINGS-WORKPLACES", "Workplace mode details loaded in "+time.Since(timer).String())
}

func loadWorkplacesSettings(writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS-WORKPLACES", "Loading workplaces settings")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-WORKPLACES", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-WORKPLACES", "Loading workplaces settings ended")
		return
	}
	var records []database.Workplace
	db.Order("id desc").Find(&records)
	var data WorkplacesSettingsDataOutput
	data.DataTableSearchTitle = getLocale(email, "data-table-search-title")
	data.DataTableInfoTitle = getLocale(email, "data-table-info-title")
	data.DataTableRowsCountTitle = getLocale(email, "data-table-rows-count-title")

	addWorkplaceSettingsTableHeaders(email, &data)
	for _, record := range records {
		addWorkplaceSettingsTableRow(record, &data)
	}

	var typeRecords []database.WorkplaceSection
	db.Order("id desc").Find(&typeRecords)
	addWorkplaceSectionSettingsTypeTableHeaders(email, &data)
	for _, record := range typeRecords {
		addWorkplaceSectionSettingsTypeTableRow(record, &data)
	}

	var extendedRecords []database.WorkplaceMode
	db.Order("id desc").Find(&extendedRecords)
	addWorkplaceModeSettingsTypeTableHeaders(email, &data)
	for _, record := range extendedRecords {
		addWorkplaceModeSettingsTypeTableRow(record, &data)
	}

	tmpl := template.Must(template.ParseFiles("./html/settings-table-type-extended.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("SETTINGS-WORKPLACES", "Workplaces settings loaded in "+time.Since(timer).String())
}

func addWorkplaceModeSettingsTypeTableRow(record database.WorkplaceMode, data *WorkplacesSettingsDataOutput) {
	var tableRow TableRowTypeExtended
	id := TableCellTypeExtended{CellNameTypeExtended: strconv.Itoa(int(record.ID))}
	tableRow.TableCellTypeExtended = append(tableRow.TableCellTypeExtended, id)
	name := TableCellTypeExtended{CellNameTypeExtended: record.Name}
	tableRow.TableCellTypeExtended = append(tableRow.TableCellTypeExtended, name)
	data.TableRowsTypeExtended = append(data.TableRowsTypeExtended, tableRow)
}

func addWorkplaceModeSettingsTypeTableHeaders(email string, data *WorkplacesSettingsDataOutput) {
	id := HeaderCellTypeExtended{HeaderNameTypeExtended: "#", HeaderWidthTypeExtended: "30"}
	data.TableHeaderTypeExtended = append(data.TableHeaderTypeExtended, id)
	name := HeaderCellTypeExtended{HeaderNameTypeExtended: getLocale(email, "type-name")}
	data.TableHeaderTypeExtended = append(data.TableHeaderTypeExtended, name)
}

func addWorkplaceSectionSettingsTypeTableRow(record database.WorkplaceSection, data *WorkplacesSettingsDataOutput) {
	var tableRow TableRowType
	id := TableCellType{CellNameType: strconv.Itoa(int(record.ID))}
	tableRow.TableCellType = append(tableRow.TableCellType, id)
	name := TableCellType{CellNameType: record.Name}
	tableRow.TableCellType = append(tableRow.TableCellType, name)
	data.TableRowsType = append(data.TableRowsType, tableRow)
}

func addWorkplaceSectionSettingsTypeTableHeaders(email string, data *WorkplacesSettingsDataOutput) {
	id := HeaderCellType{HeaderNameType: "#", HeaderWidthType: "30"}
	data.TableHeaderType = append(data.TableHeaderType, id)
	name := HeaderCellType{HeaderNameType: getLocale(email, "section-name")}
	data.TableHeaderType = append(data.TableHeaderType, name)
}

func addWorkplaceSettingsTableRow(record database.Workplace, data *WorkplacesSettingsDataOutput) {
	var tableRow TableRow
	id := TableCell{CellName: strconv.Itoa(int(record.ID))}
	tableRow.TableCell = append(tableRow.TableCell, id)
	name := TableCell{CellName: record.Name}
	tableRow.TableCell = append(tableRow.TableCell, name)
	data.TableRows = append(data.TableRows, tableRow)
}

func addWorkplaceSettingsTableHeaders(email string, data *WorkplacesSettingsDataOutput) {
	id := HeaderCell{HeaderName: "#", HeaderWidth: "30"}
	data.TableHeader = append(data.TableHeader, id)
	name := HeaderCell{HeaderName: getLocale(email, "workplace-name")}
	data.TableHeader = append(data.TableHeader, name)
}
