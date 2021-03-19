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

type DevicesSettingsDataOutput struct {
	DataTableSearchTitle    string
	DataTableInfoTitle      string
	DataTableRowsCountTitle string
	TableHeader             []HeaderCell
	TableRows               []TableRow
}

type DeviceDetailsDataOutput struct {
	DeviceName              string
	DeviceNamePrepend       string
	DeviceTypeNamePrepend   string
	IpAddress               string
	IpAddressPrepend        string
	MacAddress              string
	MacAddressPrepend       string
	DeviceVersion           string
	DeviceVersionPrepend    string
	DeviceSettings          string
	DeviceSettingsPrepend   string
	Note                    string
	NotePrepend             string
	DeviceEnabled           string
	EnabledPrepend          string
	CreatedAt               string
	CreatedAtPrepend        string
	UpdatedAt               string
	UpdatedAtPrepend        string
	DeviceTypes             []DeviceTypeSelection
	DeviceEnabledSelection  []DeviceEnabledSelection
	DataTableSearchTitle    string
	DataTableInfoTitle      string
	DataTableRowsCountTitle string
	TableHeader             []HeaderCell
	TableRows               []TableRow
	PortsHidden             string
}

type DeviceTypeSelection struct {
	DeviceTypeName     string
	DeviceTypeId       uint
	DeviceTypeSelected string
}

type DeviceEnabledSelection struct {
	DeviceEnabled         string
	DeviceEnabledSelected string
}

type DeviceDetailsDataInput struct {
	Id       string
	Name     string
	Type     string
	Ip       string
	Mac      string
	Version  string
	Settings string
	Note     string
	Enabled  string
}

type DevicePortDetailsDataOutput struct {
	DevicePortName                string
	DevicePortNamePrepend         string
	DevicePortTypeNamePrepend     string
	DevicePortFilePosition        string
	DevicePortFilePositionPrepend string
	DevicePortUnit                string
	DevicePortUnitPrepend         string
	PlcDataType                   string
	PlcDataTypePrepend            string
	PlcDataAddress                string
	PlcDataAddressPrepend         string
	Settings                      string
	SettingsPrepend               string
	Note                          string
	NotePrepend                   string
	VirtualEnabledPrepend         string
	CreatedAt                     string
	CreatedAtPrepend              string
	UpdatedAt                     string
	UpdatedAtPrepend              string
	DevicePortTypes               []DevicePortTypeSelection
	PortVirtualSelection          []PortVirtualSelection
}

type DevicePortTypeSelection struct {
	DevicePortTypeName     string
	DevicePortTypeId       uint
	DevicePortTypeSelected string
}

type PortVirtualSelection struct {
	PortVirtual         string
	PortVirtualSelected string
}

type DevicePortDetailsDataInput struct {
	Id             string
	DeviceName     string
	Name           string
	Type           string
	Position       string
	Unit           string
	PlcDataType    string
	PlcDataAddress string
	Settings       string
	Note           string
	Virtual        string
}

type DevicePortDetailsPageInput struct {
	Data string
	Type string
}

func saveDevicePortDetails(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	timer := time.Now()
	logInfo("SETTINGS-DEVICES", "Saving device port settings started")
	var data DevicePortDetailsDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS-DEVICES", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-DEVICES", "Saving device port settings ended")
		return
	}
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-DEVICES", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-DEVICES", "Saving device port settings ended")
		return
	}
	var port database.DevicePort
	db.Where("id=?", data.Id).Find(&port)
	port.Name = data.Name
	port.DevicePortTypeID = int(cachedDevicePortTypesByName[data.Type].ID)
	port.DeviceID = int(cachedDevicesByName[data.DeviceName].ID)
	port.PortNumber, _ = strconv.Atoi(data.Position)
	port.PlcDataType = data.PlcDataType
	port.PlcDataAddress = data.PlcDataAddress
	port.Settings = data.Settings
	port.Unit = data.Unit
	port.Virtual, _ = strconv.ParseBool(data.Virtual)
	port.Note = data.Note
	db.Save(&port)
	cacheDevices(db)
	logInfo("SETTINGS-DEVICES", "Device port settings saved in "+time.Since(timer).String())
}

func loadDevicePortDetail(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	timer := time.Now()
	logInfo("SETTINGS-DEVICES", "Loading device port settings")
	email, _, _ := request.BasicAuth()
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-DEVICES", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-DEVICES", "Loading device port settings ended")
		return
	}
	var data DevicePortDetailsPageInput
	err = json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS-DEVICES", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-DEVICES", "Loading workplace port settings ended")
		return
	}
	var devicePort database.DevicePort
	db.Where("id=?", data.Data).Find(&devicePort)

	var devicePortTypes []DevicePortTypeSelection
	for _, devicePortType := range cachedDevicePortTypesById {
		if devicePortType.Name == cachedDevicePortTypesById[uint(devicePort.DevicePortTypeID)].Name {
			devicePortTypes = append(devicePortTypes, DevicePortTypeSelection{DevicePortTypeName: devicePortType.Name, DevicePortTypeId: devicePortType.ID, DevicePortTypeSelected: "selected"})
		} else {
			devicePortTypes = append(devicePortTypes, DevicePortTypeSelection{DevicePortTypeName: devicePortType.Name, DevicePortTypeId: devicePortType.ID})
		}
	}
	sort.Slice(devicePortTypes, func(i, j int) bool {
		return devicePortTypes[i].DevicePortTypeName < devicePortTypes[j].DevicePortTypeName
	})

	var portVirtualSelection []PortVirtualSelection
	portVirtualSelection = append(portVirtualSelection, PortVirtualSelection{PortVirtual: "true", PortVirtualSelected: checkSelection(devicePort.Virtual, "true")})
	portVirtualSelection = append(portVirtualSelection, PortVirtualSelection{PortVirtual: "false", PortVirtualSelected: checkSelection(devicePort.Virtual, "false")})

	dataOut := DevicePortDetailsDataOutput{
		DevicePortName:                devicePort.Name,
		DevicePortNamePrepend:         getLocale(email, "port-name"),
		DevicePortTypeNamePrepend:     getLocale(email, "port-type-name"),
		DevicePortFilePosition:        strconv.Itoa(devicePort.PortNumber),
		DevicePortFilePositionPrepend: getLocale(email, "port-number"),
		DevicePortUnit:                devicePort.Unit,
		DevicePortUnitPrepend:         getLocale(email, "port-unit"),
		PlcDataType:                   devicePort.PlcDataType,
		PlcDataTypePrepend:            getLocale(email, "plc-data-type"),
		PlcDataAddress:                devicePort.PlcDataAddress,
		PlcDataAddressPrepend:         getLocale(email, "plc-data-address"),
		Settings:                      devicePort.Settings,
		SettingsPrepend:               getLocale(email, "device-settings"),
		Note:                          devicePort.Note,
		NotePrepend:                   getLocale(email, "note-name"),
		VirtualEnabledPrepend:         getLocale(email, "port-virtual"),
		CreatedAt:                     devicePort.CreatedAt.Format("2006-01-02T15:04:05"),
		CreatedAtPrepend:              getLocale(email, "created-at"),
		UpdatedAt:                     devicePort.UpdatedAt.Format("2006-01-02T15:04:05"),
		UpdatedAtPrepend:              getLocale(email, "updated-at"),
		DevicePortTypes:               devicePortTypes,
		PortVirtualSelection:          portVirtualSelection,
	}
	tmpl := template.Must(template.ParseFiles("./html/settings-detail-device-port.html"))
	_ = tmpl.Execute(writer, dataOut)
	logInfo("SETTINGS-DEVICES", "Device port settings loaded in "+time.Since(timer).String())
}

func saveDevice(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	timer := time.Now()
	logInfo("SETTINGS-DEVICES", "Saving device started")
	var data DeviceDetailsDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS-DEVICES", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-DEVICES", "Saving device ended")
		return
	}
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-DEVICES", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-DEVICES", "Saving device ended")
		return
	}
	var device database.Device
	db.Where("id=?", data.Id).Find(&device)
	device.Name = data.Name
	device.DeviceTypeID = int(cachedDeviceTypesByName[data.Type].ID)
	device.IpAddress = data.Ip
	device.MacAddress = data.Mac
	device.Settings = data.Settings
	device.TypeName = data.Version
	device.Note = data.Note
	device.Activated, _ = strconv.ParseBool(data.Enabled)
	db.Save(&device)
	cacheDevices(db)
	logInfo("SETTINGS-DEVICES", "Device saved in "+time.Since(timer).String())
}

func loadDeviceDetails(id string, writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS-DEVICES", "Loading device details")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-DEVICES", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-DEVICES", "Loading device details ended")
		return
	}
	var device database.Device
	db.Where("id = ?", id).Find(&device)

	var deviceTypes []DeviceTypeSelection
	for _, deviceType := range cachedDeviceTypesById {
		if deviceType.Name == cachedDeviceTypesById[uint(device.DeviceTypeID)].Name {
			deviceTypes = append(deviceTypes, DeviceTypeSelection{DeviceTypeName: deviceType.Name, DeviceTypeId: deviceType.ID, DeviceTypeSelected: "selected"})
		} else {
			deviceTypes = append(deviceTypes, DeviceTypeSelection{DeviceTypeName: deviceType.Name, DeviceTypeId: deviceType.ID})
		}
	}
	sort.Slice(deviceTypes, func(i, j int) bool {
		return deviceTypes[i].DeviceTypeName < deviceTypes[j].DeviceTypeName
	})

	var deviceEnabledSelection []DeviceEnabledSelection
	deviceEnabledSelection = append(deviceEnabledSelection, DeviceEnabledSelection{DeviceEnabled: "true", DeviceEnabledSelected: checkSelection(device.Activated, "true")})
	deviceEnabledSelection = append(deviceEnabledSelection, DeviceEnabledSelection{DeviceEnabled: "false", DeviceEnabledSelected: checkSelection(device.Activated, "false")})

	var records []database.DevicePort
	db.Where("device_id = ?", device.ID).Order("id desc").Find(&records)

	data := DeviceDetailsDataOutput{
		DeviceName:             device.Name,
		DeviceNamePrepend:      getLocale(email, "device-name"),
		DeviceTypeNamePrepend:  getLocale(email, "type-name"),
		IpAddress:              device.IpAddress,
		IpAddressPrepend:       getLocale(email, "ip-address"),
		MacAddress:             device.MacAddress,
		MacAddressPrepend:      getLocale(email, "mac-address"),
		DeviceVersion:          device.TypeName,
		DeviceVersionPrepend:   getLocale(email, "device-version"),
		DeviceSettings:         device.Settings,
		DeviceSettingsPrepend:  getLocale(email, "device-settings"),
		Note:                   device.Note,
		NotePrepend:            getLocale(email, "note-name"),
		EnabledPrepend:         getLocale(email, "enabled"),
		CreatedAt:              device.CreatedAt.Format("2006-01-02T15:04:05"),
		CreatedAtPrepend:       getLocale(email, "created-at"),
		UpdatedAt:              device.UpdatedAt.Format("2006-01-02T15:04:05"),
		UpdatedAtPrepend:       getLocale(email, "updated-at"),
		DeviceTypes:            deviceTypes,
		DeviceEnabledSelection: deviceEnabledSelection,
	}
	if device.DeviceTypeID == 2 {
		data.PortsHidden = "hidden"
	}
	addDeviceDetailsTableHeaders(email, &data)
	for _, record := range records {
		addDeviceDetailsTableRow(record, &data)
	}
	tmpl := template.Must(template.ParseFiles("./html/settings-detail-device.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("SETTINGS-DEVICES", "Device details loaded in "+time.Since(timer).String())
}

func addDeviceDetailsTableRow(record database.DevicePort, data *DeviceDetailsDataOutput) {
	var tableRow TableRow
	id := TableCell{CellName: strconv.Itoa(int(record.ID))}
	tableRow.TableCell = append(tableRow.TableCell, id)
	name := TableCell{CellName: record.Name}
	tableRow.TableCell = append(tableRow.TableCell, name)
	data.TableRows = append(data.TableRows, tableRow)
}

func addDeviceDetailsTableHeaders(email string, data *DeviceDetailsDataOutput) {
	id := HeaderCell{HeaderName: "#", HeaderWidth: "30"}
	data.TableHeader = append(data.TableHeader, id)
	name := HeaderCell{HeaderName: getLocale(email, "port-name")}
	data.TableHeader = append(data.TableHeader, name)
}

func loadDevicesSettings(writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS-DEVICES", "Loading devices settings")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-DEVICES", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-DEVICES", "Loading devices settings ended")
		return
	}
	var records []database.Device
	db.Order("id desc").Find(&records)
	var data DevicesSettingsDataOutput
	data.DataTableSearchTitle = getLocale(email, "data-table-search-title")
	data.DataTableInfoTitle = getLocale(email, "data-table-info-title")
	data.DataTableRowsCountTitle = getLocale(email, "data-table-rows-count-title")
	addDeviceSettingsTableHeaders(email, &data)
	for _, record := range records {
		addDeviceSettingsTableRow(record, &data)
	}
	tmpl := template.Must(template.ParseFiles("./html/settings-table.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("SETTINGS-DEVICES", "Devices settings loaded in "+time.Since(timer).String())
}

func addDeviceSettingsTableRow(record database.Device, data *DevicesSettingsDataOutput) {
	var tableRow TableRow
	id := TableCell{CellName: strconv.Itoa(int(record.ID))}
	tableRow.TableCell = append(tableRow.TableCell, id)
	name := TableCell{CellName: record.Name}
	tableRow.TableCell = append(tableRow.TableCell, name)
	data.TableRows = append(data.TableRows, tableRow)
}

func addDeviceSettingsTableHeaders(email string, data *DevicesSettingsDataOutput) {
	id := HeaderCell{HeaderName: "#", HeaderWidth: "30"}
	data.TableHeader = append(data.TableHeader, id)
	name := HeaderCell{HeaderName: getLocale(email, "device-name")}
	data.TableHeader = append(data.TableHeader, name)
}
