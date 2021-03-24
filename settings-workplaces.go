package main

import (
	"database/sql"
	"encoding/json"
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

func deleteWorkplacePort(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	timer := time.Now()
	logInfo("SETTINGS-WORKPLACES", "Deleting workplace port started")
	var data WorkplacePortDetailsDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS-WORKPLACES", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-WORKPLACES", "Deleting workplace port ended")
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
		logInfo("SETTINGS-WORKPLACES", "Deleting workplace port ended")
		return
	}

	var workplacePort database.WorkplacePort
	db.Where("id=?", data.Id).Find(&workplacePort)
	db.Delete(&workplacePort)
	cacheWorkplaces(db)
	logInfo("SETTINGS-WORKPLACES", "Workplace port deleted in "+time.Since(timer).String())
}

func saveWorkplacePortDetails(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	timer := time.Now()
	logInfo("SETTINGS-WORKPLACES", "Saving workplace port started")
	var data WorkplacePortDetailsDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS-WORKPLACES", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-WORKPLACES", "Saving workplace port ended")
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
		logInfo("SETTINGS-WORKPLACES", "Saving workplace port ended")
		return
	}

	var workplacePort database.WorkplacePort
	db.Where("id=?", data.Id).Find(&workplacePort)
	workplacePort.Name = data.Name
	workplacePort.DevicePortID, _ = strconv.Atoi(strings.TrimRight(strings.Split(data.DevicePortId, "[")[1], "]"))
	workplacePort.WorkplaceID = int(cachedWorkplacesByName[data.WorkplaceName].ID)
	if len(data.StateId) > 0 {
		stateId, _ := strconv.Atoi(strings.TrimRight(strings.Split(data.DevicePortId, "[")[1], "]"))
		workplacePort.StateID = sql.NullInt32{Int32: int32(stateId), Valid: true}
	}
	workplacePort.Color = data.Color
	workplacePort.Note = data.Note
	db.Save(&workplacePort)
	cacheWorkplaces(db)
	logInfo("SETTINGS-WORKPLACES", "Workplace port saved in "+time.Since(timer).String())
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
	db.Save(&workplace)

	var workplacePorts []database.WorkplacePort
	db.Where("workplace_id = ?", workplace.ID).Find(&workplacePorts)
	for _, port := range workplacePorts {
		db.Model(&port).Select("counter_ok", "counter_nok", "state_id").Updates(map[string]interface{}{"counter_ok": false, "counter_nok": false, "state_id": nil})
	}
	if len(data.ProductionDowntimeSelection) > 1 {
		productionDowntimeId := strings.TrimRight(strings.Split(data.ProductionDowntimeSelection, "[")[1], "]")
		var devicePort database.DevicePort
		db.Where("id = ?", productionDowntimeId).Find(&devicePort)
		var port database.WorkplacePort
		db.Where("device_port_id = ?", productionDowntimeId).Where("workplace_id = ?", workplace.ID).Find(&port)
		port.StateID = sql.NullInt32{Int32: 1, Valid: true}
		port.DevicePortID, _ = strconv.Atoi(productionDowntimeId)
		port.Name = devicePort.Name
		port.WorkplaceID = int(workplace.ID)
		port.Color = data.ProductionDowntimeColor
		db.Save(&port)
	}
	if len(data.PowerOnPowerOffSelection) > 1 {
		poweroffDowntimeId := strings.TrimRight(strings.Split(data.PowerOnPowerOffSelection, "[")[1], "]")
		var devicePort database.DevicePort
		db.Where("id = ?", poweroffDowntimeId).Find(&devicePort)
		var port database.WorkplacePort
		db.Where("device_port_id = ?", poweroffDowntimeId).Where("workplace_id = ?", workplace.ID).Find(&port)
		port.StateID = sql.NullInt32{Int32: 3, Valid: true}
		port.DevicePortID, _ = strconv.Atoi(poweroffDowntimeId)
		port.Name = devicePort.Name
		port.WorkplaceID = int(workplace.ID)
		port.Color = data.PowerOnPowerOffColor
		db.Save(&port)
	}
	if len(data.CountOkSelection) > 1 {
		countOkId := strings.TrimRight(strings.Split(data.CountOkSelection, "[")[1], "]")
		var devicePort database.DevicePort
		db.Where("id = ?", countOkId).Find(&devicePort)
		var port database.WorkplacePort
		db.Where("device_port_id = ?", countOkId).Where("workplace_id = ?", workplace.ID).Find(&port)
		port.DevicePortID, _ = strconv.Atoi(countOkId)
		port.Name = devicePort.Name
		port.WorkplaceID = int(workplace.ID)
		port.Color = data.CountOkColor
		port.CounterOK = true
		db.Save(&port)
	}
	if len(data.CountNokSelection) > 1 {
		countNokId := strings.TrimRight(strings.Split(data.CountNokSelection, "[")[1], "]")
		var devicePort database.DevicePort
		db.Where("id = ?", countNokId).Find(&devicePort)
		var port database.WorkplacePort
		db.Where("device_port_id = ?", countNokId).Where("workplace_id = ?", workplace.ID).Find(&port)
		port.DevicePortID, _ = strconv.Atoi(countNokId)
		port.Name = devicePort.Name
		port.WorkplaceID = int(workplace.ID)
		port.Color = data.CountNokColor
		port.CounterNOK = true
		db.Save(&port)
	}
	cacheWorkplaces(db)

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
			logInfo("SETTINGS-WORKPLACES", "Deleting workplace workshift record id "+strconv.Itoa(int(workplaceWorkshift.ID)))
			db.Delete(&workplaceWorkshift)
		}
	}
	for _, name := range pageWorkshiftIds {
		logInfo("SETTINGS-WORKPLACES", "Creating workplace workshift record for "+name+" and "+workplace.Name)
		var workplaceWorkshift database.WorkplaceWorkshift
		workplaceWorkshift.WorkshiftID, _ = strconv.Atoi(strings.TrimRight(strings.Split(name, "[")[1], "]"))
		workplaceWorkshift.WorkplaceID = int(workplace.ID)
		db.Save(&workplaceWorkshift)
	}
	cacheWorkShifts(db)
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
	poweroffParsed, err := time.ParseDuration(data.PowerOffDuration)
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

	var digitalPorts []database.DevicePort
	db.Where("device_port_type_id = 1").Find(&digitalPorts)
	var analogPorts []database.DevicePort
	db.Where("device_port_type_id = 2").Find(&analogPorts)

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
	db.Where("workplace_id = ?", workplace.ID).Where("state_id = 1").Find(&workplaceProductionPort)
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

	var workplacePoweroffPort database.WorkplacePort
	db.Where("workplace_id = ?", workplace.ID).Where("state_id = 3").Find(&workplacePoweroffPort)
	var poweronPoweroffs []WorkplacePortSelection
	for _, port := range analogPorts {
		if int(port.ID) == workplacePoweroffPort.DevicePortID {
			poweronPoweroffs = append(poweronPoweroffs, WorkplacePortSelection{WorkplacePortName: cachedDevicesById[uint(port.DeviceID)].Name + " #" + strconv.Itoa(port.PortNumber) + ": " + port.Name + " [" + strconv.Itoa(int(port.ID)) + "]", WorkplacePortId: port.ID, WorkplacePortSelected: "selected"})
			var workplacePort database.WorkplacePort
			db.Where("device_port_id = ?", port.ID).Where("workplace_id = ?", workplace.ID).Find(&workplacePort)
			data.PoweroffColorValue = workplacePort.Color
		} else {
			poweronPoweroffs = append(poweronPoweroffs, WorkplacePortSelection{WorkplacePortName: cachedDevicesById[uint(port.DeviceID)].Name + " #" + strconv.Itoa(port.PortNumber) + ": " + port.Name + " [" + strconv.Itoa(int(port.ID)) + "]", WorkplacePortId: port.ID})
		}
	}
	data.PoweronPoweroffs = poweronPoweroffs

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
	addWorkplacePortDetailsTableHeaders(email, &data)
	for _, record := range records {
		if record.StateID.Valid || record.CounterNOK || record.CounterOK {
			continue
		} else {
			addWorkplacePortDetailsTableRow(record, &data)
		}
	}
	addWorkplaceWorkshiftDetailsTableHeaders(email, &data)
	for _, workshift := range workshifts {
		addWorkplaceWorkshiftDetailsTableRow(workshift, &data)
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

	tmpl := template.Must(template.ParseFiles("./html/settings-detail-workplace.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("SETTINGS-WORKPLACES", "Workplace details loaded in "+time.Since(timer).String())
}

func addWorkplaceWorkshiftDetailsTableRow(record database.WorkplaceWorkshift, data *WorkplaceDetailsDataOutput) {
	var tableRow WorkshiftTableRow
	id := WorkshiftTableCell{WorkshiftCellName: strconv.Itoa(int(record.ID))}
	tableRow.WorkshiftTableCell = append(tableRow.WorkshiftTableCell, id)
	name := WorkshiftTableCell{WorkshiftCellName: cachedWorkShiftsById[uint(record.WorkshiftID)].Name}
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