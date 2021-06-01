package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/petrjahoda/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"strings"
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
	Result                  string
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
	Result                   string
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
	Result                      string
}

type WorkplaceDetailsDataOutput struct {
	WorkplaceName             string
	WorkplaceNamePrepend      string
	WorkplaceSectionPrepend   string
	WorkshiftsPrepend         string
	WorkplaceModePrepend      string
	ProductionDowntimePrepend string
	PoweronPoweroffPrepend    string
	CountOkPrepend            string
	CountNokPrepend           string
	WorkplaceCode             string
	WorkplaceCodePrepend      string
	Note                      string
	NotePrepend               string
	CreatedAt                 string
	CreatedAtPrepend          string
	UpdatedAt                 string
	UpdatedAtPrepend          string
	ProductionColorValue      string
	PoweroffColorValue        string
	OkColorValue              string
	NokColorValue             string
	WorkplaceSections         []WorkplaceSectionSelection
	WorkplaceModes            []WorkplaceModeSelection
	ProductionDowntimes       []WorkplacePortSelection
	PoweronPoweroffs          []WorkplacePortSelection
	CountOks                  []WorkplacePortSelection
	CountNoks                 []WorkplacePortSelection
	Workshifts                []WorkshiftSelection
	DataFilterPlaceholder     string
	DataTableSearchTitle      string
	DataTableInfoTitle        string
	DataTableRowsCountTitle   string
	TableHeader               []HeaderCell
	TableRows                 []TableRow
	WorkshiftTableHeader      []WorkshiftHeaderCell
	WorkshiftTableRows        []WorkshiftTableRow
	Result                    string
}

type WorkshiftSelection struct {
	WorkshiftName      string
	WorkshiftSelection string
}

type WorkplacePortSelection struct {
	WorkplacePortName     string
	WorkplacePortId       uint
	WorkplacePortSelected string
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
	PowerOffDuration string
	Note             string
}

type WorkplaceSectionDetailsDataInput struct {
	Id   string
	Name string
	Note string
}

type WorkplaceDetailsDataInput struct {
	Id                          string
	Name                        string
	Section                     string
	ProductionDowntimeSelection string
	ProductionDowntimeColor     string
	PowerOnPowerOffSelection    string
	PowerOnPowerOffColor        string
	WorkShifts                  []string
	CountOkSelection            string
	CountOkColor                string
	CountNokSelection           string
	CountNokColor               string
	Mode                        string
	Code                        string
	Note                        string
}

type WorkplacePortDetailsPageInput struct {
	Data        string
	WorkplaceId string
}

type WorkplacePortDetailsDataInput struct {
	Id            string
	WorkplaceName string
	Name          string
	DevicePortId  string
	StateId       string
	Color         string
	Note          string
}

type WorkplacePortDetailsDataOutput struct {
	WorkplacePortName              string
	WorkplacePortNamePrepend       string
	WorkplacePortDevicePortPrepend string
	WorkplacePortStatePrepend      string
	Color                          string
	Note                           string
	NotePrepend                    string
	CreatedAt                      string
	CreatedAtPrepend               string
	UpdatedAt                      string
	UpdatedAtPrepend               string
	WorkplacePortDevicePorts       []WorkplacePortDevicePortSelection
	WorkplacePortStates            []WorkplacePortStateSelection
	Result                         string
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

func loadWorkplaces(writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS", "Loading workplaces")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS", "Problem opening database: "+err.Error())
		var responseData WorkplacesSettingsDataOutput
		responseData.Result = "ERR: Problem opening database, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Loading workplaces ended with error")
		return
	}
	var records []database.Workplace
	db.Order("id desc").Find(&records)
	var data WorkplacesSettingsDataOutput
	data.DataTableSearchTitle = getLocale(email, "data-table-search-title")
	data.DataTableInfoTitle = getLocale(email, "data-table-info-title")
	data.DataTableRowsCountTitle = getLocale(email, "data-table-rows-count-title")
	addWorkplacesTableHeaders(email, &data)
	for _, record := range records {
		addWorkplacesTableRow(record, &data)
	}
	var typeRecords []database.WorkplaceSection
	db.Order("id desc").Find(&typeRecords)
	addWorkplaceSectionsTableHeaders(email, &data)
	for _, record := range typeRecords {
		addWorkplaceSectionsTableRow(record, &data)
	}
	var extendedRecords []database.WorkplaceMode
	db.Order("id desc").Find(&extendedRecords)
	addWorkplaceModesTableHeaders(email, &data)
	for _, record := range extendedRecords {
		addWorkplaceModesTableRow(record, &data)
	}
	tmpl, err := template.ParseFiles("./html/settings-table-type-extended.html")
	if err != nil {
		logError("SETTINGS", "Problem parsing html file: "+err.Error())
		var responseData FaultSettingsDataOutput
		responseData.Result = "ERR: Problem parsing html file: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
	} else {
		data.Result = "INF: Workplaces processed in " + time.Since(timer).String()
		_ = tmpl.Execute(writer, data)
		logInfo("SETTINGS", "Workplaces loaded in "+time.Since(timer).String())
	}
}

func addWorkplaceModesTableRow(record database.WorkplaceMode, data *WorkplacesSettingsDataOutput) {
	var tableRow TableRowTypeExtended
	id := TableCellTypeExtended{CellNameTypeExtended: strconv.Itoa(int(record.ID))}
	tableRow.TableCellTypeExtended = append(tableRow.TableCellTypeExtended, id)
	name := TableCellTypeExtended{CellNameTypeExtended: record.Name}
	tableRow.TableCellTypeExtended = append(tableRow.TableCellTypeExtended, name)
	data.TableRowsTypeExtended = append(data.TableRowsTypeExtended, tableRow)
}

func addWorkplaceModesTableHeaders(email string, data *WorkplacesSettingsDataOutput) {
	id := HeaderCellTypeExtended{HeaderNameTypeExtended: "#", HeaderWidthTypeExtended: "30"}
	data.TableHeaderTypeExtended = append(data.TableHeaderTypeExtended, id)
	name := HeaderCellTypeExtended{HeaderNameTypeExtended: getLocale(email, "type-name")}
	data.TableHeaderTypeExtended = append(data.TableHeaderTypeExtended, name)
}

func addWorkplaceSectionsTableRow(record database.WorkplaceSection, data *WorkplacesSettingsDataOutput) {
	var tableRow TableRowType
	id := TableCellType{CellNameType: strconv.Itoa(int(record.ID))}
	tableRow.TableCellType = append(tableRow.TableCellType, id)
	name := TableCellType{CellNameType: record.Name}
	tableRow.TableCellType = append(tableRow.TableCellType, name)
	data.TableRowsType = append(data.TableRowsType, tableRow)
}

func addWorkplaceSectionsTableHeaders(email string, data *WorkplacesSettingsDataOutput) {
	id := HeaderCellType{HeaderNameType: "#", HeaderWidthType: "30"}
	data.TableHeaderType = append(data.TableHeaderType, id)
	name := HeaderCellType{HeaderNameType: getLocale(email, "section-name")}
	data.TableHeaderType = append(data.TableHeaderType, name)
}

func addWorkplacesTableRow(record database.Workplace, data *WorkplacesSettingsDataOutput) {
	var tableRow TableRow
	id := TableCell{CellName: strconv.Itoa(int(record.ID))}
	tableRow.TableCell = append(tableRow.TableCell, id)
	name := TableCell{CellName: record.Name}
	tableRow.TableCell = append(tableRow.TableCell, name)
	data.TableRows = append(data.TableRows, tableRow)
}

func addWorkplacesTableHeaders(email string, data *WorkplacesSettingsDataOutput) {
	id := HeaderCell{HeaderName: "#", HeaderWidth: "30"}
	data.TableHeader = append(data.TableHeader, id)
	name := HeaderCell{HeaderName: getLocale(email, "workplace-name")}
	data.TableHeader = append(data.TableHeader, name)
}

func loadWorkplaceSection(id string, writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS", "Loading workplace section")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS", "Problem opening database: "+err.Error())
		var responseData WorkplaceSectionDetailsDataOutput
		responseData.Result = "nok: " + err.Error()

		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-USERS", "Loading workplace section ended with error")
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
	tmpl, err := template.ParseFiles("./html/settings-detail-workplace-section.html")
	if err != nil {
		logError("SETTINGS", "Problem parsing html file: "+err.Error())
		var responseData WorkplaceSectionDetailsDataOutput
		responseData.Result = "ERR: Problem parsing html file: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
	} else {
		data.Result = "INF: Workplace section detail processed in " + time.Since(timer).String()
		_ = tmpl.Execute(writer, data)
		logInfo("SETTINGS", "Workplace section detail loaded in "+time.Since(timer).String())
	}
}

func loadWorkplaceMode(id string, writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS", "Loading workplace mode")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS", "Problem opening database: "+err.Error())
		var responseData WorkplaceModeDetailsDataOutput
		responseData.Result = "ERR: Problem opening database, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-USERS", "Loading workplace mode ended with error")
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
	tmpl, err := template.ParseFiles("./html/settings-detail-workplace-mode.html")
	if err != nil {
		logError("SETTINGS", "Problem parsing html file: "+err.Error())
		var responseData WorkplaceModeDetailsDataOutput
		responseData.Result = "ERR: Problem parsing html file: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
	} else {
		data.Result = "INF: Workplace mode detail processed in " + time.Since(timer).String()
		_ = tmpl.Execute(writer, data)
		logInfo("SETTINGS", "Workplace mode detail loaded in "+time.Since(timer).String())
	}
}

func loadWorkplace(id string, writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS", "Loading workplace")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS", "Problem opening database: "+err.Error())
		var responseData WorkplaceDetailsDataOutput
		responseData.Result = "ERR: Problem opening database, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Loading workplace ended with error")
		return
	}
	var workplace database.Workplace
	db.Where("id = ?", id).Find(&workplace)
	workplaceSections := loadWorkplaceSections(workplace)
	workplaceModes := loadWorkplaceModes(workplace)
	var digitalPorts []database.DevicePort
	db.Where("device_port_type_id = ?", digital).Find(&digitalPorts)
	var analogPorts []database.DevicePort
	db.Where("device_port_type_id = ?", analog).Find(&analogPorts)
	data := WorkplaceDetailsDataOutput{
		WorkplaceName:             workplace.Name,
		WorkplaceNamePrepend:      getLocale(email, "workplace-name"),
		WorkplaceSectionPrepend:   getLocale(email, "section-name"),
		WorkshiftsPrepend:         getLocale(email, "workshifts"),
		WorkplaceModePrepend:      getLocale(email, "type-name"),
		ProductionDowntimePrepend: getLocale(email, "production-downtime"),
		PoweronPoweroffPrepend:    getLocale(email, "poweron-poweroff"),
		CountOkPrepend:            getLocale(email, "good-pieces-name"),
		CountNokPrepend:           getLocale(email, "bad-pieces-name"),
		DataFilterPlaceholder:     getLocale(email, "data-table-search-title"),
		WorkplaceCode:             workplace.Code,
		WorkplaceCodePrepend:      getLocale(email, "code-name"),
		Note:                      workplace.Note,
		NotePrepend:               getLocale(email, "note-name"),
		CreatedAt:                 workplace.CreatedAt.Format("2006-01-02T15:04:05"),
		CreatedAtPrepend:          getLocale(email, "created-at"),
		UpdatedAt:                 workplace.UpdatedAt.Format("2006-01-02T15:04:05"),
		UpdatedAtPrepend:          getLocale(email, "updated-at"),
		WorkplaceSections:         workplaceSections,
		WorkplaceModes:            workplaceModes,
	}
	var workplaceProductionPort database.WorkplacePort
	db.Where("workplace_id = ?", workplace.ID).Where("state_id = ?", production).Find(&workplaceProductionPort)
	var productionDowntimes []WorkplacePortSelection
	for _, port := range digitalPorts {
		if int(port.ID) == workplaceProductionPort.DevicePortID {
			productionDowntimes = append(productionDowntimes, WorkplacePortSelection{WorkplacePortName: cachedDevicesById[uint(port.DeviceID)].Name + " #" + strconv.Itoa(port.PortNumber) + ": " + port.Name + " [" + strconv.Itoa(int(port.ID)) + "]", WorkplacePortId: port.ID, WorkplacePortSelected: "selected"})
			var workplacePort database.WorkplacePort
			db.Where("device_port_id = ?", port.ID).Where("workplace_id = ?", workplace.ID).Find(&workplacePort)
			data.ProductionColorValue = workplacePort.Color
		} else {
			productionDowntimes = append(productionDowntimes, WorkplacePortSelection{WorkplacePortName: cachedDevicesById[uint(port.DeviceID)].Name + " #" + strconv.Itoa(port.PortNumber) + ": " + port.Name + " [" + strconv.Itoa(int(port.ID)) + "]", WorkplacePortId: port.ID})
		}
	}
	data.ProductionDowntimes = productionDowntimes
	var workplacePowerOffPort database.WorkplacePort
	db.Where("workplace_id = ?", workplace.ID).Where("state_id = ?", poweroff).Find(&workplacePowerOffPort)
	var powerOnPowerOffs []WorkplacePortSelection
	for _, port := range analogPorts {
		if int(port.ID) == workplacePowerOffPort.DevicePortID {
			powerOnPowerOffs = append(powerOnPowerOffs, WorkplacePortSelection{WorkplacePortName: cachedDevicesById[uint(port.DeviceID)].Name + " #" + strconv.Itoa(port.PortNumber) + ": " + port.Name + " [" + strconv.Itoa(int(port.ID)) + "]", WorkplacePortId: port.ID, WorkplacePortSelected: "selected"})
			var workplacePort database.WorkplacePort
			db.Where("device_port_id = ?", port.ID).Where("workplace_id = ?", workplace.ID).Find(&workplacePort)
			data.PoweroffColorValue = workplacePort.Color
		} else {
			powerOnPowerOffs = append(powerOnPowerOffs, WorkplacePortSelection{WorkplacePortName: cachedDevicesById[uint(port.DeviceID)].Name + " #" + strconv.Itoa(port.PortNumber) + ": " + port.Name + " [" + strconv.Itoa(int(port.ID)) + "]", WorkplacePortId: port.ID})
		}
	}
	data.PoweronPoweroffs = powerOnPowerOffs
	var countOkPort database.WorkplacePort
	db.Where("workplace_id = ?", workplace.ID).Where("counter_ok is true").Find(&countOkPort)
	var countOks []WorkplacePortSelection
	countOks = append(countOks, WorkplacePortSelection{WorkplacePortName: ""})
	for _, port := range digitalPorts {
		if int(port.ID) == countOkPort.DevicePortID {
			countOks = append(countOks, WorkplacePortSelection{WorkplacePortName: cachedDevicesById[uint(port.DeviceID)].Name + " #" + strconv.Itoa(port.PortNumber) + ": " + port.Name + " [" + strconv.Itoa(int(port.ID)) + "]", WorkplacePortId: port.ID, WorkplacePortSelected: "selected"})
			var workplacePort database.WorkplacePort
			db.Where("device_port_id = ?", port.ID).Where("workplace_id = ?", workplace.ID).Find(&workplacePort)
			data.OkColorValue = workplacePort.Color
		} else {
			countOks = append(countOks, WorkplacePortSelection{WorkplacePortName: cachedDevicesById[uint(port.DeviceID)].Name + " #" + strconv.Itoa(port.PortNumber) + ": " + port.Name + " [" + strconv.Itoa(int(port.ID)) + "]", WorkplacePortId: port.ID})
		}
	}
	data.CountOks = countOks
	var countNokPort database.WorkplacePort
	db.Where("workplace_id = ?", workplace.ID).Where("counter_nok is true").Find(&countNokPort)
	var countNoks []WorkplacePortSelection
	countNoks = append(countNoks, WorkplacePortSelection{WorkplacePortName: ""})
	for _, port := range digitalPorts {
		if int(port.ID) == countNokPort.DevicePortID {
			countNoks = append(countNoks, WorkplacePortSelection{WorkplacePortName: cachedDevicesById[uint(port.DeviceID)].Name + " #" + strconv.Itoa(port.PortNumber) + ": " + port.Name + " [" + strconv.Itoa(int(port.ID)) + "]", WorkplacePortId: port.ID, WorkplacePortSelected: "selected"})
			var workplacePort database.WorkplacePort
			db.Where("device_port_id = ?", port.ID).Where("workplace_id = ?", workplace.ID).Find(&workplacePort)
			data.NokColorValue = workplacePort.Color
		} else {
			countNoks = append(countNoks, WorkplacePortSelection{WorkplacePortName: cachedDevicesById[uint(port.DeviceID)].Name + " #" + strconv.Itoa(port.PortNumber) + ": " + port.Name + " [" + strconv.Itoa(int(port.ID)) + "]", WorkplacePortId: port.ID})
		}
	}
	data.CountNoks = countNoks
	var records []database.WorkplacePort
	db.Where("workplace_id = ?", workplace.ID).Find(&records)
	var workshifts []database.WorkplaceWorkshift
	db.Where("workplace_id = ?", workplace.ID).Order("id desc").Find(&workshifts)
	addWorkplacePortsTableHeaders(email, &data)
	for _, record := range records {
		if record.StateID.Valid || record.CounterNOK || record.CounterOK {
			continue
		} else {
			addWorkplacePortsTableRow(record, &data)
		}
	}
	addWorkplaceWorkshiftsTableHeaders(email, &data)
	for _, workshift := range workshifts {
		addWorkplaceWorkshiftsTableRow(workshift, &data)
	}
	var workshiftSelection []WorkshiftSelection
	for _, workshift := range cachedWorkShiftsById {
		workshiftAdded := false
		for _, workplaceWorkshift := range cachedWorkplaceWorkShiftsById {
			if int(workshift.ID) == workplaceWorkshift.WorkshiftID && int(workplace.ID) == workplaceWorkshift.WorkplaceID {
				workshiftSelection = append(workshiftSelection, WorkshiftSelection{WorkshiftName: workshift.Name + " [" + strconv.Itoa(int(workshift.ID)) + "]", WorkshiftSelection: "selected"})
				workshiftAdded = true
				break
			}
		}
		if !workshiftAdded {
			workshiftSelection = append(workshiftSelection, WorkshiftSelection{WorkshiftName: workshift.Name + " [" + strconv.Itoa(int(workshift.ID)) + "]"})
		}
	}
	sort.Slice(workshiftSelection, func(i, j int) bool {
		return workshiftSelection[i].WorkshiftName < workshiftSelection[j].WorkshiftName
	})
	data.Workshifts = workshiftSelection

	tmpl, err := template.ParseFiles("./html/settings-detail-workplace.html")
	if err != nil {
		logError("SETTINGS", "Problem parsing html file: "+err.Error())
		var responseData WorkplaceDetailsDataOutput
		responseData.Result = "ERR: Problem parsing html file: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
	} else {
		data.Result = "INF: Workplace detail processed in " + time.Since(timer).String()
		_ = tmpl.Execute(writer, data)
		logInfo("SETTINGS", "Workplace detail loaded in "+time.Since(timer).String())
	}
}

func loadWorkplaceModes(workplace database.Workplace) []WorkplaceModeSelection {
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
	return workplaceModes
}

func loadWorkplaceSections(workplace database.Workplace) []WorkplaceSectionSelection {
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
	return workplaceSections
}

func addWorkplaceWorkshiftsTableRow(record database.WorkplaceWorkshift, data *WorkplaceDetailsDataOutput) {
	var tableRow WorkshiftTableRow
	id := WorkshiftTableCell{WorkshiftCellName: strconv.Itoa(int(record.ID))}
	tableRow.WorkshiftTableCell = append(tableRow.WorkshiftTableCell, id)
	name := WorkshiftTableCell{WorkshiftCellName: cachedWorkShiftsById[uint(record.WorkshiftID)].Name}
	tableRow.WorkshiftTableCell = append(tableRow.WorkshiftTableCell, name)
	data.WorkshiftTableRows = append(data.WorkshiftTableRows, tableRow)
}

func addWorkplaceWorkshiftsTableHeaders(email string, data *WorkplaceDetailsDataOutput) {
	id := WorkshiftHeaderCell{WorkshiftHeaderName: "#", WorkshiftHeaderWidth: "30"}
	data.WorkshiftTableHeader = append(data.WorkshiftTableHeader, id)
	name := WorkshiftHeaderCell{WorkshiftHeaderName: getLocale(email, "workshift-name")}
	data.WorkshiftTableHeader = append(data.WorkshiftTableHeader, name)
}

func addWorkplacePortsTableRow(record database.WorkplacePort, data *WorkplaceDetailsDataOutput) {
	var tableRow TableRow
	id := TableCell{CellName: strconv.Itoa(int(record.ID))}
	tableRow.TableCell = append(tableRow.TableCell, id)
	name := TableCell{CellName: record.Name}
	tableRow.TableCell = append(tableRow.TableCell, name)
	data.TableRows = append(data.TableRows, tableRow)
}

func addWorkplacePortsTableHeaders(email string, data *WorkplaceDetailsDataOutput) {
	id := HeaderCell{HeaderName: "#", HeaderWidth: "30"}
	data.TableHeader = append(data.TableHeader, id)
	name := HeaderCell{HeaderName: getLocale(email, "port-name")}
	data.TableHeader = append(data.TableHeader, name)
}

func loadWorkplacePort(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	timer := time.Now()
	logInfo("SETTINGS", "Loading workplace port")
	email, _, _ := request.BasicAuth()
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS", "Problem opening database: "+err.Error())
		var responseData WorkplacePortDetailsDataOutput
		responseData.Result = "ERR: Problem opening database, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Loading workplace port ended with error")
		return
	}
	var data WorkplacePortDetailsPageInput
	err = json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS", "Error parsing data: "+err.Error())
		var responseData WorkplacePortDetailsDataOutput
		responseData.Result = "ERR: Error parsing data, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Loading workplace port ended with error")
		return
	}
	var workplacePort database.WorkplacePort
	db.Where("id=?", data.Data).Find(&workplacePort)
	var devicePorts []database.DevicePort
	db.Raw("select * from device_ports where id not in (select device_port_id from workplace_ports where workplace_id=? and (state_id in (1,3) or counter_ok is true or counter_nok is true))", data.WorkplaceId).Find(&devicePorts)
	var workplacePortDevicePorts []WorkplacePortDevicePortSelection
	for _, port := range devicePorts {
		if port.ID == cachedDevicePortsById[uint(workplacePort.DevicePortID)].ID {
			workplacePortDevicePorts = append(workplacePortDevicePorts, WorkplacePortDevicePortSelection{WorkplacePortDevicePortName: cachedDevicesById[uint(port.DeviceID)].Name + " #" + strconv.Itoa(port.PortNumber) + ": " + port.Name + " [" + strconv.Itoa(int(port.ID)) + "]", WorkplacePortDevicePortId: port.ID, WorkplacePortDevicePortSelected: "selected"})
		} else {
			workplacePortDevicePorts = append(workplacePortDevicePorts, WorkplacePortDevicePortSelection{WorkplacePortDevicePortName: cachedDevicesById[uint(port.DeviceID)].Name + " #" + strconv.Itoa(port.PortNumber) + ": " + port.Name + " [" + strconv.Itoa(int(port.ID)) + "]", WorkplacePortDevicePortId: port.ID})
		}
	}
	sort.Slice(workplacePortDevicePorts, func(i, j int) bool {
		return workplacePortDevicePorts[i].WorkplacePortDevicePortName < workplacePortDevicePorts[j].WorkplacePortDevicePortName
	})
	var states []WorkplacePortStateSelection
	for _, state := range cachedStatesById {
		if state.ID == cachedStatesById[uint(workplacePort.StateID.Int32)].ID {
			states = append(states, WorkplacePortStateSelection{WorkplacePortStateName: state.Name + " [" + strconv.Itoa(int(state.ID)) + "]", WorkplacePortStateId: state.ID, WorkplacePortStateSelected: "selected"})
		} else {
			states = append(states, WorkplacePortStateSelection{WorkplacePortStateName: state.Name + " [" + strconv.Itoa(int(state.ID)) + "]", WorkplacePortStateId: state.ID})
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
		Note:                           workplacePort.Note,
		NotePrepend:                    getLocale(email, "note-name"),
		CreatedAt:                      workplacePort.CreatedAt.Format("2006-01-02T15:04:05"),
		CreatedAtPrepend:               getLocale(email, "created-at"),
		UpdatedAt:                      workplacePort.UpdatedAt.Format("2006-01-02T15:04:05"),
		UpdatedAtPrepend:               getLocale(email, "updated-at"),
		WorkplacePortDevicePorts:       workplacePortDevicePorts,
		WorkplacePortStates:            states,
	}
	tmpl, err := template.ParseFiles("./html/settings-detail-workplace-port.html")
	if err != nil {
		logError("SETTINGS", "Problem parsing html file: "+err.Error())
		var responseData WorkplacePortDetailsDataOutput
		responseData.Result = "ERR: Problem parsing html file: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
	} else {
		dataOut.Result = "INF: Workplace port detail processed in " + time.Since(timer).String()
		_ = tmpl.Execute(writer, dataOut)
		logInfo("SETTINGS", "Workplace port detail loaded in "+time.Since(timer).String())
	}
}

func saveWorkplace(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	timer := time.Now()
	logInfo("SETTINGS", "Saving workplace")
	var data WorkplaceDetailsDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "ERR: Error parsing data, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Saving workplace ended with error")
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
		logInfo("SETTINGS", "Saving workplace ended with error")
		return
	}
	var workplace database.Workplace
	db.Where("id=?", data.Id).Find(&workplace)
	workplace.Name = data.Name
	workplace.WorkplaceModeID = int(cachedWorkplaceModesByName[data.Mode].ID)
	workplace.WorkplaceSectionID = int(cachedWorkplaceSectionsByName[data.Section].ID)
	workplace.Code = data.Code
	workplace.Note = data.Note
	result := db.Save(&workplace)
	var workplacePorts []database.WorkplacePort
	db.Where("workplace_id = ?", workplace.ID).Find(&workplacePorts)
	for _, port := range workplacePorts {
		db.Model(&port).Select("counter_ok", "counter_nok", "state_id").Updates(map[string]interface{}{"counter_ok": false, "counter_nok": false, "state_id": nil})
	}
	if len(data.ProductionDowntimeSelection) > 1 {
		productionDowntimeId := strings.TrimRight(strings.Split(data.ProductionDowntimeSelection, "[")[1], "]")
		var devicePort database.DevicePort
		db.Where("id = ?", productionDowntimeId).Find(&devicePort)
		var workplacePort database.WorkplacePort
		db.Where("device_port_id = ?", productionDowntimeId).Where("workplace_id = ?", workplace.ID).Find(&workplacePort)
		if workplacePort.ID > 0 {
			workplacePort.StateID = sql.NullInt32{Int32: production, Valid: true}
			workplacePort.DevicePortID, _ = strconv.Atoi(productionDowntimeId)
			workplacePort.Name = devicePort.Name
			workplacePort.WorkplaceID = int(workplace.ID)
			workplacePort.Color = data.ProductionDowntimeColor
			result = db.Save(&workplacePort)
			cacheWorkplaces(db)
			logInfo("SETTINGS", "Workplace port "+workplacePort.Name+" saved in "+time.Since(timer).String())
		} else {
			productionDowntimeIdAsInt, _ := strconv.Atoi(productionDowntimeId)
			db.Raw("SELECT * from workplace_ports where workplace_id = ? and device_port_id = ? and deleted_at is not null", int(cachedWorkplacesByName[data.Name].ID), productionDowntimeIdAsInt).Find(&workplacePort)
			workplacePort.StateID = sql.NullInt32{Int32: production, Valid: true}
			workplacePort.DevicePortID = productionDowntimeIdAsInt
			workplacePort.Name = devicePort.Name
			workplacePort.WorkplaceID = int(workplace.ID)
			workplacePort.Color = data.ProductionDowntimeColor
			workplacePort.DeletedAt = gorm.DeletedAt{
				Time:  time.Time{},
				Valid: false,
			}
			result = db.Save(&workplacePort)
			cacheWorkplaces(db)
			cacheWorkplacePorts(db)
			logInfo("SETTINGS", "Workplace port "+workplacePort.Name+" saved in "+time.Since(timer).String())
		}
	}
	if len(data.PowerOnPowerOffSelection) > 1 {
		poweroffDowntimeId := strings.TrimRight(strings.Split(data.PowerOnPowerOffSelection, "[")[1], "]")
		var devicePort database.DevicePort
		db.Where("id = ?", poweroffDowntimeId).Find(&devicePort)
		var workplacePort database.WorkplacePort
		db.Where("device_port_id = ?", poweroffDowntimeId).Where("workplace_id = ?", workplace.ID).Find(&workplacePort)
		if workplacePort.ID > 0 {
			workplacePort.StateID = sql.NullInt32{Int32: poweroff, Valid: true}
			workplacePort.DevicePortID, _ = strconv.Atoi(poweroffDowntimeId)
			workplacePort.Name = devicePort.Name
			workplacePort.WorkplaceID = int(workplace.ID)
			workplacePort.Color = data.PowerOnPowerOffColor
			result = db.Save(&workplacePort)
			cacheWorkplaces(db)
			logInfo("SETTINGS", "Workplace port "+workplacePort.Name+" saved in "+time.Since(timer).String())
		} else {
			poweroffDowntimeIdAsInt, _ := strconv.Atoi(poweroffDowntimeId)
			db.Raw("SELECT * from workplace_ports where workplace_id = ? and device_port_id = ? and deleted_at is not null", int(cachedWorkplacesByName[data.Name].ID), poweroffDowntimeIdAsInt).Find(&workplacePort)
			workplacePort.StateID = sql.NullInt32{Int32: poweroff, Valid: true}
			workplacePort.DevicePortID = poweroffDowntimeIdAsInt
			workplacePort.Name = devicePort.Name
			workplacePort.WorkplaceID = int(workplace.ID)
			workplacePort.Color = data.PowerOnPowerOffColor
			workplacePort.DeletedAt = gorm.DeletedAt{
				Time:  time.Time{},
				Valid: false,
			}
			result = db.Save(&workplacePort)
			cacheWorkplaces(db)
			cacheWorkplacePorts(db)
			logInfo("SETTINGS", "Workplace port "+workplacePort.Name+" saved in "+time.Since(timer).String())
		}
	}
	if len(data.CountOkSelection) > 1 {
		countOkId := strings.TrimRight(strings.Split(data.CountOkSelection, "[")[1], "]")
		var devicePort database.DevicePort
		db.Where("id = ?", countOkId).Find(&devicePort)
		var workplacePort database.WorkplacePort
		db.Where("device_port_id = ?", countOkId).Where("workplace_id = ?", workplace.ID).Find(&workplacePort)
		if workplacePort.ID > 0 {
			fmt.Println("naslo")
			workplacePort.DevicePortID, _ = strconv.Atoi(countOkId)
			workplacePort.Name = devicePort.Name
			workplacePort.WorkplaceID = int(workplace.ID)
			workplacePort.Color = data.CountOkColor
			workplacePort.CounterOK = true
			result = db.Save(&workplacePort)
			cacheWorkplaces(db)
			logInfo("SETTINGS", "Workplace port "+workplacePort.Name+" saved in "+time.Since(timer).String())
		} else {
			fmt.Println("nenaslo")
			countOkIdAsInt, _ := strconv.Atoi(countOkId)
			db.Debug().Raw("SELECT * from workplace_ports where workplace_id = ? and device_port_id = ? and deleted_at is not null", int(cachedWorkplacesByName[data.Name].ID), countOkIdAsInt).Find(&workplacePort)
			workplacePort.DevicePortID = countOkIdAsInt
			workplacePort.Name = devicePort.Name
			workplacePort.WorkplaceID = int(workplace.ID)
			workplacePort.Color = data.CountOkColor
			workplacePort.CounterOK = true
			workplacePort.DeletedAt = gorm.DeletedAt{
				Time:  time.Time{},
				Valid: false,
			}
			result = db.Save(&workplacePort)
			cacheWorkplaces(db)
			cacheWorkplacePorts(db)
			logInfo("SETTINGS", "Workplace port "+workplacePort.Name+" saved in "+time.Since(timer).String())
		}

	}
	if len(data.CountNokSelection) > 1 {
		countNokId := strings.TrimRight(strings.Split(data.CountNokSelection, "[")[1], "]")
		var devicePort database.DevicePort
		db.Where("id = ?", countNokId).Find(&devicePort)
		var workplacePort database.WorkplacePort
		db.Where("device_port_id = ?", countNokId).Where("workplace_id = ?", workplace.ID).Find(&workplacePort)
		if workplacePort.ID > 0 {
			workplacePort.DevicePortID, _ = strconv.Atoi(countNokId)
			workplacePort.Name = devicePort.Name
			workplacePort.WorkplaceID = int(workplace.ID)
			workplacePort.Color = data.CountNokColor
			workplacePort.CounterNOK = true
			result = db.Save(&workplacePort)
			cacheWorkplaces(db)
			logInfo("SETTINGS", "Workplace port "+workplacePort.Name+" saved in "+time.Since(timer).String())
		} else {
			countNokIdAsInt, _ := strconv.Atoi(countNokId)
			db.Raw("SELECT * from workplace_ports where workplace_id = ? and device_port_id = ? and deleted_at is not null", int(cachedWorkplacesByName[data.Name].ID), countNokIdAsInt).Find(&workplacePort)
			workplacePort.DevicePortID = countNokIdAsInt
			workplacePort.Name = devicePort.Name
			workplacePort.WorkplaceID = int(workplace.ID)
			workplacePort.Color = data.CountNokColor
			workplacePort.CounterNOK = true
			workplacePort.DeletedAt = gorm.DeletedAt{
				Time:  time.Time{},
				Valid: false,
			}
			result = db.Save(&workplacePort)
			cacheWorkplaces(db)
			cacheWorkplacePorts(db)
			logInfo("SETTINGS", "Workplace port "+workplacePort.Name+" saved in "+time.Since(timer).String())
		}

	}
	cacheWorkplaces(db)
	cacheWorkplacePorts(db)
	var databaseWorkplaceWorkshifts []database.WorkplaceWorkshift
	db.Where("workplace_id = ?", workplace.ID).Find(&databaseWorkplaceWorkshifts)
	pageWorkshiftIds := make(map[int]string)
	for _, workshift := range data.WorkShifts {
		workshiftAsInt, _ := strconv.Atoi(strings.TrimRight(strings.Split(workshift, "[")[1], "]"))
		pageWorkshiftIds[workshiftAsInt] = workshift
	}

	for _, workplaceWorkshift := range databaseWorkplaceWorkshifts {
		_, found := pageWorkshiftIds[workplaceWorkshift.WorkshiftID]
		if found {
			delete(pageWorkshiftIds, workplaceWorkshift.WorkshiftID)
		} else {
			logInfo("SETTINGS", "Deleting workplace workshift record id "+strconv.Itoa(int(workplaceWorkshift.ID)))
			result = db.Delete(&workplaceWorkshift)
		}
	}
	for _, name := range pageWorkshiftIds {
		logInfo("SETTINGS", "Creating workplace workshift record for "+name+" and "+workplace.Name)
		var workplaceWorkshift database.WorkplaceWorkshift
		workplaceWorkshift.WorkshiftID, _ = strconv.Atoi(strings.TrimRight(strings.Split(name, "[")[1], "]"))
		workplaceWorkshift.WorkplaceID = int(workplace.ID)
		result = db.Save(&workplaceWorkshift)
	}
	cacheWorkShifts(db)
	if result.Error != nil {
		var responseData TableOutput
		responseData.Result = "ERR: Workplace not saved: " + result.Error.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logError("SETTINGS", "Workplace "+workplace.Name+" not saved: "+result.Error.Error())
	} else {
		var responseData TableOutput
		responseData.Result = "INF: Workplace saved in " + time.Since(timer).String()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Workplace "+workplace.Name+" saved in "+time.Since(timer).String())
	}
}

func saveWorkplaceSection(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	timer := time.Now()
	logInfo("SETTINGS", "Saving workplace section")
	var data WorkplaceSectionDetailsDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "ERR: Error parsing data, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Saving workplace section ended with error")
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
		logInfo("SETTINGS", "Saving workplace section ended with error")
		return
	}
	var workplaceSection database.WorkplaceSection
	db.Where("id=?", data.Id).Find(&workplaceSection)
	workplaceSection.Name = data.Name
	workplaceSection.Note = data.Note
	result := db.Save(&workplaceSection)
	cacheWorkplaces(db)
	if result.Error != nil {
		var responseData TableOutput
		responseData.Result = "ERR: Workplace section not saved: " + result.Error.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logError("SETTINGS", "Workplace section "+workplaceSection.Name+" not saved: "+result.Error.Error())
	} else {
		var responseData TableOutput
		responseData.Result = "INF: Workplace section saved in " + time.Since(timer).String()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Workplace section "+workplaceSection.Name+" saved in "+time.Since(timer).String())
	}
}

func saveWorkplaceMode(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	timer := time.Now()
	logInfo("SETTINGS", "Saving workplace mode")
	var data WorkplaceModeDetailsDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "ERR: Error parsing data, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Saving workplace mode ended with error")
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
		logInfo("SETTINGS", "Saving workplace mode ended with error")
		return
	}
	downtimeParsed, err := time.ParseDuration(data.DowntimeDuration)
	if err != nil {
		logError("SETTINGS-PRODUCTS", "Problem parsing downtime duration: "+err.Error())
		downtimeParsed = 0
	}
	poweroffParsed, err := time.ParseDuration(data.PowerOffDuration)
	if err != nil {
		logError("SETTINGS-PRODUCTS", "Problem parsing powerOff duration: "+err.Error())
		poweroffParsed = 0
	}
	var workplaceMode database.WorkplaceMode
	db.Where("id=?", data.Id).Find(&workplaceMode)
	workplaceMode.Name = data.Name
	workplaceMode.PoweroffDuration = poweroffParsed
	workplaceMode.DowntimeDuration = downtimeParsed
	workplaceMode.Note = data.Note
	result := db.Save(&workplaceMode)
	cacheWorkplaces(db)
	if result.Error != nil {
		var responseData TableOutput
		responseData.Result = "ERR: Workplace mode not saved: " + result.Error.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logError("SETTINGS", "Workplace mode "+workplaceMode.Name+" not saved: "+result.Error.Error())
	} else {
		var responseData TableOutput
		responseData.Result = "INF: Workplace mode saved in " + time.Since(timer).String()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Workplace mode "+workplaceMode.Name+" saved in "+time.Since(timer).String())
	}
}

func saveWorkplacePort(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	timer := time.Now()
	logInfo("SETTINGS", "Saving workplace port")
	var data WorkplacePortDetailsDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Saving workplace port ended with error")
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
		logInfo("SETTINGS", "Saving workplace port ended with error")
		return
	}

	var workplacePort database.WorkplacePort
	db.Where("id=?", data.Id).Find(&workplacePort)
	devicePortId, _ := strconv.Atoi(strings.TrimRight(strings.Split(data.DevicePortId, "[")[1], "]"))
	if workplacePort.ID > 0 {
		workplacePort.Name = data.Name
		workplacePort.DevicePortID = devicePortId
		workplacePort.WorkplaceID = int(cachedWorkplacesByName[data.WorkplaceName].ID)
		if len(data.StateId) > 0 {
			stateId, _ := strconv.Atoi(strings.TrimRight(strings.Split(data.DevicePortId, "[")[1], "]"))
			workplacePort.StateID = sql.NullInt32{Int32: int32(stateId), Valid: true}
		}
		workplacePort.Color = data.Color
		workplacePort.Note = data.Note
		db.Save(&workplacePort)
		cacheWorkplaces(db)
		logInfo("SETTINGS", "Workplace port "+workplacePort.Name+" saved in "+time.Since(timer).String())
	} else {
		db.Raw("SELECT * from workplace_ports where name = ? and workplace_id = ? and device_port_id = ? and deleted_at is not null", data.Name, int(cachedWorkplacesByName[data.WorkplaceName].ID), devicePortId).Find(&workplacePort)
		workplacePort.Name = data.Name
		workplacePort.DevicePortID = devicePortId
		workplacePort.WorkplaceID = int(cachedWorkplacesByName[data.WorkplaceName].ID)
		if len(data.StateId) > 0 {
			stateId, _ := strconv.Atoi(strings.TrimRight(strings.Split(data.DevicePortId, "[")[1], "]"))
			workplacePort.StateID = sql.NullInt32{Int32: int32(stateId), Valid: true}
		}
		workplacePort.Color = data.Color
		workplacePort.Note = data.Note
		workplacePort.DeletedAt = gorm.DeletedAt{
			Time:  time.Time{},
			Valid: false,
		}
		db.Save(&workplacePort)
		cacheWorkplaces(db)
		logInfo("SETTINGS", "Workplace port "+workplacePort.Name+" saved in "+time.Since(timer).String())
	}
}

func deleteWorkplacePort(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	timer := time.Now()
	logInfo("SETTINGS", "Deleting workplace port")
	var data WorkplacePortDetailsDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "ERR: Error parsing data, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Deleting workplace port ended with error")
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
		logInfo("SETTINGS", "Deleting workplace port ended with error")
		return
	}

	var workplacePort database.WorkplacePort
	db.Where("id=?", data.Id).Find(&workplacePort)
	db.Delete(&workplacePort)
	cacheWorkplaces(db)
	var responseData TableOutput
	responseData.Result = "INF: Data deleted in " + time.Since(timer).String()
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(responseData)
	logInfo("SETTINGS", "Workplace port "+workplacePort.Name+" deleted in "+time.Since(timer).String())
}
