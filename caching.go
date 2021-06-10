package main

import (
	"database/sql"
	"github.com/petrjahoda/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strconv"
	"sync"
	"time"
)

var cachedCompanyName string
var location string
var cachedLocales = map[string]string{}

var cachedUserWebSettings = map[string]map[string]string{}
var cachedBreakdownTypesByName = map[string]database.BreakdownType{}
var cachedDevicesByName = map[string]database.Device{}
var cachedDevicePortsByName = map[string]database.DevicePort{}
var cachedWorkplaceDevicePorts = map[string][]database.DevicePort{}
var cachedWorkplacePorts = map[string][]database.WorkplacePort{}
var cachedDeviceTypesByName = map[string]database.DeviceType{}
var cachedDevicePortTypesByName = map[string]database.DevicePortType{}
var cachedDowntimeTypesByName = map[string]database.DowntimeType{}
var cachedFaultTypesByName = map[string]database.FaultType{}
var cachedLocalesByName = map[string]database.Locale{}
var cachedOrdersByName = map[string]database.Order{}
var cachedPackageTypesByName = map[string]database.PackageType{}
var cachedProductsByName = map[string]database.Product{}
var cachedUsersByEmail = map[string]database.User{}
var cachedUserRolesByName = map[string]database.UserRole{}
var cachedUserTypesByName = map[string]database.UserType{}
var cachedWorkplacesByName = map[string]database.Workplace{}
var cachedWorkplaceModesByName = map[string]database.WorkplaceMode{}
var cachedWorkplaceSectionsByName = map[string]database.WorkplaceSection{}
var cachedWorkplacesProductionRecords = map[string]map[string]time.Duration{}
var cachedWorkplacesDowntimeRecords = map[string]map[string]time.Duration{}
var cachedWorkplacesPoweroffRecords = map[string]map[string]time.Duration{}
var latestCachedWorkplaceCalendarData = time.Now()
var latestCachedWorkplaceConsumption = time.Now()
var cachedDevicePortsColorsById = map[int]string{}
var cachedAlarmsById = map[uint]database.Alarm{}
var cachedBreakdownsById = map[uint]database.Breakdown{}
var cachedBreakdownTypesById = map[uint]database.BreakdownType{}
var cachedDevicesById = map[uint]database.Device{}
var cachedDevicePortsById = map[uint]database.DevicePort{}
var cachedDevicePortTypesById = map[uint]database.DevicePortType{}
var cachedDeviceTypesById = map[uint]database.DeviceType{}
var cachedDowntimesById = map[uint]database.Downtime{}
var cachedDowntimeTypesById = map[uint]database.DowntimeType{}
var cachedFaultsById = map[uint]database.Fault{}
var cachedFaultTypesById = map[uint]database.FaultType{}
var cachedOperationsById = map[uint]database.Operation{}
var cachedOrdersById = map[uint]database.Order{}
var cachedPackagesById = map[uint]database.Package{}
var cachedPackageTypesById = map[uint]database.PackageType{}
var cachedPartsById = map[uint]database.Part{}
var cachedProductsById = map[uint]database.Product{}
var cachedStatesById = map[uint]database.State{}
var cachedUsersById = map[uint]database.User{}
var cachedUserRolesById = map[uint]database.UserRole{}
var cachedUserTypesById = map[uint]database.UserType{}
var cachedWorkplacesById = map[uint]database.Workplace{}
var cachedWorkplaceModesById = map[uint]database.WorkplaceMode{}
var cachedWorkplaceSectionsById = map[uint]database.WorkplaceSection{}
var cachedWorkShiftsById = map[uint]database.Workshift{}
var cachedWorkplaceWorkShiftsById = map[uint]database.WorkplaceWorkshift{}
var cachedConsumptionDataByWorkplaceName = map[string]map[string]float32{}

var workplacesRecords sync.RWMutex
var alarmsSync sync.RWMutex
var breakdownsSync sync.RWMutex
var companyNameSync sync.RWMutex
var devicesSync sync.RWMutex
var devicePortColors sync.RWMutex
var downtimesSync sync.RWMutex
var faultsSync sync.RWMutex
var localesSync sync.RWMutex
var operationsSync sync.RWMutex
var ordersSync sync.RWMutex
var packagesSync sync.RWMutex
var partsSync sync.RWMutex
var productsSync sync.RWMutex
var statesSync sync.RWMutex
var usersSync sync.RWMutex
var userSettingsSync sync.RWMutex
var workplacesSync sync.RWMutex
var workplaceDevicePortsSync sync.RWMutex
var workShiftsSync sync.RWMutex
var consumptionDataSync sync.RWMutex

func cacheData() {
	logInfo("CACHING", "Caching started")
	timer := time.Now()
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("CACHING", "Problem opening database: "+err.Error())
		return
	}

	cacheSystemSettings(db)
	cacheProductionData(db, time.Now().AddDate(-1, 0, 0))
	cacheUsers(db)
	cacheDevices(db)
	cacheLocales(db)
	cacheOrders(db)
	cacheOperations(db)
	cacheProducts(db)
	cacheWorkplaces(db)
	cacheWorkShifts(db)
	cacheAlarms(db)
	cacheBreakdowns(db)
	cacheDowntimes(db)
	cacheFaults(db)
	cachePackages(db)
	cacheParts(db)
	cacheStates(db)
	cacheWorkplaceDevicePorts(db)
	cacheWorkplacePorts(db)
	companyNameSync.Lock()
	loc, _ := time.LoadLocation(location)
	companyNameSync.Unlock()
	cacheConsumptionData(db, time.Date(time.Now().Add(-720*time.Hour).Year(), time.Now().Add(-720*time.Hour).Month(), time.Now().Add(-720*time.Hour).Day(), 0, 0, 0, 0, loc))
	logInfo("CACHING", "Initial caching done in "+time.Since(timer).String())
}

func cacheProductionData(db *gorm.DB, fromDate time.Time) {
	companyNameSync.Lock()
	loc, _ := time.LoadLocation(location)
	companyNameSync.Unlock()
	logInfo("CACHING", "Caching workplace production/downtime/poweroff from: "+fromDate.In(loc).String()+"  back until: "+time.Now().In(loc).String())
	var workplaces []database.Workplace
	db.Find(&workplaces)
	for _, workplace := range workplaces {
		if cachedWorkplacesProductionRecords[workplace.Name] == nil {
			cachedWorkplacesProductionRecords[workplace.Name] = make(map[string]time.Duration)
		}
		if cachedWorkplacesDowntimeRecords[workplace.Name] == nil {
			cachedWorkplacesDowntimeRecords[workplace.Name] = make(map[string]time.Duration)
		}
		if cachedWorkplacesPoweroffRecords[workplace.Name] == nil {
			cachedWorkplacesPoweroffRecords[workplace.Name] = make(map[string]time.Duration)
		}
		productionRecords := cachedWorkplacesProductionRecords[workplace.Name]
		downtimeRecords := cachedWorkplacesDowntimeRecords[workplace.Name]
		poweroffRecords := cachedWorkplacesPoweroffRecords[workplace.Name]
		tempDate := fromDate
		for tempDate.Before(time.Now()) {
			productionRecords[tempDate.Format("2006-01-02")] = 0
			downtimeRecords[tempDate.Format("2006-01-02")] = 0
			poweroffRecords[tempDate.Format("2006-01-02")] = 0
			tempDate = tempDate.Add(24 * time.Hour)
		}
		var stateRecords []database.StateRecord
		latestCachedWorkplaceCalendarData = time.Now().In(loc)
		db.Select("state_id, date_time_start").Where("date_time_start >= ?", fromDate).Where("date_time_start <= ?", time.Now().In(loc)).Where("workplace_id = ?", workplace.ID).Order("date_time_start asc").Order("id asc").Find(&stateRecords)
		for index, record := range stateRecords {
			nextDate := time.Now().In(loc)
			if index < len(stateRecords)-1 {
				nextDate = stateRecords[index+1].DateTimeStart
			}
			if record.StateID == production {
				if record.DateTimeStart.In(loc).Day() == nextDate.In(loc).Day() {
					productionRecords[record.DateTimeStart.In(loc).Format("2006-01-02")] += nextDate.In(loc).Sub(record.DateTimeStart.In(loc))
				} else {
					endOfRecordDay := time.Date(record.DateTimeStart.In(loc).Year(), record.DateTimeStart.In(loc).Month(), record.DateTimeStart.In(loc).Day()+1, 0, 0, 0, 0, loc)
					for record.DateTimeStart.In(loc).Before(nextDate.In(loc)) {
						productionRecords[record.DateTimeStart.In(loc).Format("2006-01-02")] += endOfRecordDay.In(loc).Sub(record.DateTimeStart.In(loc))
						record.DateTimeStart = endOfRecordDay.In(loc)
						endOfRecordDay = time.Date(record.DateTimeStart.In(loc).Year(), record.DateTimeStart.In(loc).Month(), record.DateTimeStart.In(loc).Day()+1, 0, 0, 0, 0, loc)
					}
					endOfRecordDay = endOfRecordDay.In(loc).Add(-24 * time.Hour)
					productionRecords[nextDate.In(loc).Format("2006-01-02")] += nextDate.In(loc).Sub(endOfRecordDay.In(loc))
				}
			}
			if record.StateID == downtime {
				if record.DateTimeStart.In(loc).Day() == nextDate.In(loc).Day() {
					downtimeRecords[record.DateTimeStart.In(loc).Format("2006-01-02")] += nextDate.In(loc).Sub(record.DateTimeStart.In(loc))
				} else {
					endOfRecordDay := time.Date(record.DateTimeStart.In(loc).Year(), record.DateTimeStart.In(loc).Month(), record.DateTimeStart.In(loc).Day()+1, 0, 0, 0, 0, loc)
					for record.DateTimeStart.In(loc).Before(nextDate.In(loc)) {
						downtimeRecords[record.DateTimeStart.In(loc).Format("2006-01-02")] += endOfRecordDay.In(loc).Sub(record.DateTimeStart.In(loc))
						record.DateTimeStart = endOfRecordDay.In(loc)
						endOfRecordDay = time.Date(record.DateTimeStart.In(loc).Year(), record.DateTimeStart.In(loc).Month(), record.DateTimeStart.In(loc).Day()+1, 0, 0, 0, 0, loc)
					}
					endOfRecordDay = endOfRecordDay.In(loc).Add(-24 * time.Hour)
					downtimeRecords[nextDate.In(loc).Format("2006-01-02")] += nextDate.In(loc).Sub(endOfRecordDay.In(loc))
				}
			}
		}
		today := time.Now().In(loc).Format("2006-01-02")
		for date, productionDuration := range productionRecords {
			downtimeDuration := downtimeRecords[date]
			if date == today {
				todaysDuration := time.Now().In(loc).Sub(time.Date(time.Now().In(loc).Year(), time.Now().In(loc).Month(), time.Now().In(loc).Day(), 0, 0, 0, 0, loc))
				poweroffRecords[date] = time.Duration(len(workplaces))*todaysDuration - downtimeDuration - productionDuration
				continue
			}
			poweroffRecords[date] = time.Duration(len(workplaces)*24)*time.Hour - downtimeDuration - productionDuration
		}
		workplacesRecords.Lock()
		cachedWorkplacesProductionRecords[workplace.Name] = productionRecords
		cachedWorkplacesDowntimeRecords[workplace.Name] = downtimeRecords
		cachedWorkplacesPoweroffRecords[workplace.Name] = poweroffRecords
		workplacesRecords.Unlock()
	}
	logInfo("CACHING", "Cached production/downtime/poweroff data with length of  "+strconv.Itoa(len(cachedWorkplacesProductionRecords))+" workplaces")
}

func cacheConsumptionData(db *gorm.DB, date time.Time) {
	companyNameSync.Lock()
	loc, _ := time.LoadLocation(location)
	companyNameSync.Unlock()
	latestCachedWorkplaceConsumption = time.Now().In(loc)
	consumptionDataSync.Lock()
	var workplaces []database.Workplace
	db.Find(&workplaces)
	for _, workplace := range workplaces {
		if cachedConsumptionDataByWorkplaceName[workplace.Name] == nil {
			cachedConsumptionDataByWorkplaceName[workplace.Name] = make(map[string]float32)
		}
		tempConsumptionData := cachedConsumptionDataByWorkplaceName[workplace.Name]
		tempDate := date
		for tempDate.Before(time.Now()) {
			tempConsumptionData[tempDate.Format("2006-01-02")] = 0
			tempDate = tempDate.Add(24 * time.Hour)
		}
		tempDate = date
		for _, port := range cachedWorkplacePorts[workplace.Name] {
			if port.StateID.Int32 == poweroff {
				for tempDate.Before(time.Now()) {
					toDate := tempDate.Add(24 * time.Hour)
					var result sql.NullFloat64
					db.Raw("select ((sum(Data)/count(id))*230*(count(id)/360))/1000 from device_port_analog_records where device_port_id=? and date_time >= ? and date_time <= ?", port.DevicePortID, tempDate.In(loc), toDate.In(loc)).Scan(&result)
					if result.Valid {
						tempConsumptionData[tempDate.In(loc).Format("2006-01-02")] = float32(result.Float64)
					} else {
						tempConsumptionData[tempDate.In(loc).Format("2006-01-02")] = 0.0
					}
					tempDate = tempDate.Add(24 * time.Hour)
				}
				cachedConsumptionDataByWorkplaceName[workplace.Name] = tempConsumptionData
			}
		}
	}
	consumptionDataSync.Unlock()
	logInfo("CACHING", "Cached consumption data for "+strconv.Itoa(len(cachedConsumptionDataByWorkplaceName))+" workplaces")
}

func cacheWorkplacePorts(db *gorm.DB) {
	var workplacePorts []database.WorkplacePort
	db.Find(&workplacePorts)
	devicePortColors.Lock()
	for _, workplacePort := range workplacePorts {
		cachedDevicePortsColorsById[workplacePort.DevicePortID] = workplacePort.Color
	}
	devicePortColors.Unlock()
	logInfo("CACHING", "Cached "+strconv.Itoa(len(cachedDevicePortsColorsById))+" workplace port colors")
}

func cacheWorkplaceDevicePorts(db *gorm.DB) {
	var allWorkplaces []database.Workplace
	db.Find(&allWorkplaces)
	workplaceDevicePortsSync.Lock()
	for _, workplace := range allWorkplaces {
		var allWorkplacePorts []database.WorkplacePort
		var allDevicePorts []database.DevicePort
		db.Where("workplace_id = ?", workplace.ID).Find(&allWorkplacePorts)
		for _, workplacePort := range allWorkplacePorts {
			var devicePort database.DevicePort
			db.Where("id = ?", workplacePort.DevicePortID).Find(&devicePort)
			allDevicePorts = append(allDevicePorts, devicePort)
		}
		cachedWorkplaceDevicePorts[workplace.Name] = allDevicePorts
		cachedWorkplacePorts[workplace.Name] = allWorkplacePorts

	}
	workplaceDevicePortsSync.Unlock()
	logInfo("CACHING", "Cached "+strconv.Itoa(len(cachedWorkplaceDevicePorts))+" workplace ports")
}

func cacheSystemSettings(db *gorm.DB) {
	var companyName database.Setting
	db.Where("name like 'company'").Find(&companyName)
	companyNameSync.Lock()
	cachedCompanyName = companyName.Value
	companyNameSync.Unlock()
	logInfo("CACHING", "Cached company name")

	var timezone database.Setting
	db.Where("name=?", "timezone").Find(&timezone)
	companyNameSync.Lock()
	location = timezone.Value
	companyNameSync.Unlock()
	logInfo("CACHING", "Cached timezone")
}

func cacheLocales(db *gorm.DB) {
	var locales []database.Locale
	db.Find(&locales)
	localesSync.Lock()
	for _, locale := range locales {
		cachedLocalesByName[locale.Name] = locale
	}
	cachedLocales["CsCZ"] = "cs-CZ"
	cachedLocales["DeDE"] = "de-DE"
	cachedLocales["EnUS"] = "en-US"
	cachedLocales["EsES"] = "es-MX"
	cachedLocales["FrFR"] = "fr-FR"
	cachedLocales["ItIT"] = "it-IT"
	cachedLocales["PlPL"] = "pl-PL"
	cachedLocales["PtPT"] = "pt-BR"
	cachedLocales["SkSK"] = "sk-SK"
	cachedLocales["RuRU"] = "ru-RU"
	cachedLocales["EnUS"] = "en-US"
	localesSync.Unlock()
	logInfo("CACHING", "Cached "+strconv.Itoa(len(cachedLocalesByName))+" locales")
}

func cacheOrders(db *gorm.DB) {
	var orders []database.Order
	db.Find(&orders)
	ordersSync.Lock()
	for _, order := range orders {
		cachedOrdersById[order.ID] = order
		cachedOrdersByName[order.Name] = order
	}
	ordersSync.Unlock()
	logInfo("CACHING", "Cached "+strconv.Itoa(len(cachedOrdersById))+" orders")
}

func cacheOperations(db *gorm.DB) {
	var operations []database.Operation
	db.Find(&operations)
	operationsSync.Lock()
	for _, operation := range operations {
		cachedOperationsById[operation.ID] = operation
	}
	operationsSync.Unlock()
	logInfo("CACHING", "Cached "+strconv.Itoa(len(cachedOperationsById))+" operations")
}

func cacheWorkplaces(db *gorm.DB) {
	var workplaceModes []database.WorkplaceMode
	var workplaceSections []database.WorkplaceSection
	var workplaces []database.Workplace
	db.Find(&workplaces)
	db.Find(&workplaceModes)
	db.Find(&workplaceSections)
	workplacesSync.Lock()
	for _, workplaceMode := range workplaceModes {
		cachedWorkplaceModesById[workplaceMode.ID] = workplaceMode
		cachedWorkplaceModesByName[workplaceMode.Name] = workplaceMode

	}
	for _, workplaceSection := range workplaceSections {
		cachedWorkplaceSectionsById[workplaceSection.ID] = workplaceSection
		cachedWorkplaceSectionsByName[workplaceSection.Name] = workplaceSection

	}
	for _, workplace := range workplaces {
		cachedWorkplacesById[workplace.ID] = workplace
		cachedWorkplacesByName[workplace.Name] = workplace

	}
	workplacesSync.Unlock()
	logInfo("CACHING", "Cached "+strconv.Itoa(len(cachedWorkplacesById))+" workplaces")
	logInfo("CACHING", "Cached "+strconv.Itoa(len(cachedWorkplaceModesById))+" workplace modes")
	logInfo("CACHING", "Cached "+strconv.Itoa(len(cachedWorkplaceSectionsById))+" workplace sections")
}

func cacheDevices(db *gorm.DB) {
	var devices []database.Device
	var deviceTypes []database.DeviceType
	var devicePortTypes []database.DevicePortType
	var devicePorts []database.DevicePort
	db.Find(&devices)
	db.Find(&deviceTypes)
	db.Find(&devicePortTypes)
	db.Find(&devicePorts)
	devicesSync.Lock()
	cachedDevicesById = map[uint]database.Device{}
	cachedDevicesByName = map[string]database.Device{}
	cachedDevicePortsById = map[uint]database.DevicePort{}
	cachedDevicePortsByName = map[string]database.DevicePort{}
	cachedDeviceTypesById = map[uint]database.DeviceType{}
	cachedDeviceTypesByName = map[string]database.DeviceType{}
	cachedDevicePortTypesById = map[uint]database.DevicePortType{}
	cachedDevicePortTypesByName = map[string]database.DevicePortType{}
	for _, device := range devices {
		cachedDevicesById[device.ID] = device
		cachedDevicesByName[device.Name] = device
	}
	for _, deviceType := range deviceTypes {
		cachedDeviceTypesById[deviceType.ID] = deviceType
		cachedDeviceTypesByName[deviceType.Name] = deviceType
	}
	for _, devicePortType := range devicePortTypes {
		cachedDevicePortTypesById[devicePortType.ID] = devicePortType
		cachedDevicePortTypesByName[devicePortType.Name] = devicePortType
	}
	for _, devicePort := range devicePorts {
		cachedDevicePortsById[devicePort.ID] = devicePort
		cachedDevicePortsByName[devicePort.Name] = devicePort
	}

	devicesSync.Unlock()
	logInfo("CACHING", "Cached "+strconv.Itoa(len(cachedDevicesById))+" devices")
	logInfo("CACHING", "Cached "+strconv.Itoa(len(cachedDevicePortsById))+" device ports")
	logInfo("CACHING", "Cached "+strconv.Itoa(len(cachedDeviceTypesById))+" device types")
	logInfo("CACHING", "Cached "+strconv.Itoa(len(cachedDevicePortTypesById))+" device port types")
}

func cacheUsers(db *gorm.DB) {
	var users []database.User
	var userTypes []database.UserType
	var userRoles []database.UserRole
	db.Find(&users)
	db.Find(&userTypes)
	db.Find(&userRoles)
	usersSync.Lock()
	userSettingsSync.Lock()
	cachedUsersByEmail = map[string]database.User{}
	cachedUsersById = map[uint]database.User{}
	cachedUserWebSettings = map[string]map[string]string{}
	for _, user := range users {
		if len(user.Email) > 0 {
			cachedUsersByEmail[user.Email] = user
			_, userWebCached := cachedUserWebSettings[user.Email]
			if !userWebCached {
				data := map[string]string{}
				cachedUserWebSettings[user.Email] = data
			}

		}
		cachedUsersById[user.ID] = user
	}
	cachedUserTypesById = map[uint]database.UserType{}
	cachedUserTypesByName = map[string]database.UserType{}
	for _, userType := range userTypes {
		cachedUserTypesById[userType.ID] = userType
		cachedUserTypesByName[userType.Name] = userType
	}
	cachedUserRolesById = map[uint]database.UserRole{}
	cachedUserRolesByName = map[string]database.UserRole{}
	for _, userRole := range userRoles {
		cachedUserRolesById[userRole.ID] = userRole
		cachedUserRolesByName[userRole.Name] = userRole
	}
	usersSync.Unlock()
	userSettingsSync.Unlock()
	logInfo("CACHING", "Cached "+strconv.Itoa(len(cachedUsersByEmail))+" users")
	logInfo("CACHING", "Cached "+strconv.Itoa(len(cachedUserWebSettings))+" user settings")
}

func cachePackages(db *gorm.DB) {
	var packages []database.Package
	var packageTypes []database.PackageType
	db.Find(&packages)
	db.Find(&packageTypes)
	packagesSync.Lock()
	for _, onePackage := range packages {
		cachedPackagesById[onePackage.ID] = onePackage

	}
	for _, packageType := range packageTypes {
		cachedPackageTypesById[packageType.ID] = packageType
		cachedPackageTypesByName[packageType.Name] = packageType
	}
	packagesSync.Unlock()
	logInfo("CACHING", "Cached "+strconv.Itoa(len(cachedPackagesById))+" packages")
}

func cacheFaults(db *gorm.DB) {
	var faults []database.Fault
	var faultTypes []database.FaultType
	db.Find(&faults)
	db.Find(&faultTypes)
	faultsSync.Lock()
	for _, fault := range faults {
		cachedFaultsById[fault.ID] = fault

	}
	for _, faultType := range faultTypes {
		cachedFaultTypesById[faultType.ID] = faultType
		cachedFaultTypesByName[faultType.Name] = faultType
	}
	faultsSync.Unlock()
	logInfo("CACHING", "Cached "+strconv.Itoa(len(cachedFaultsById))+" faults")
}

func cacheDowntimes(db *gorm.DB) {
	var downtimes []database.Downtime
	var downtimeTypes []database.DowntimeType
	db.Find(&downtimes)
	db.Find(&downtimeTypes)
	downtimesSync.Lock()
	for _, downtime := range downtimes {
		cachedDowntimesById[downtime.ID] = downtime
	}
	for _, downtimeType := range downtimeTypes {
		cachedDowntimeTypesById[downtimeType.ID] = downtimeType
		cachedDowntimeTypesByName[downtimeType.Name] = downtimeType
	}
	downtimesSync.Unlock()
	logInfo("CACHING", "Cached "+strconv.Itoa(len(cachedDowntimesById))+" downtimes")
}

func cacheBreakdowns(db *gorm.DB) {
	var breakdowns []database.Breakdown
	var breakdownTypes []database.BreakdownType
	db.Find(&breakdowns)
	db.Find(&breakdownTypes)
	breakdownsSync.Lock()
	for _, breakdown := range breakdowns {
		cachedBreakdownsById[breakdown.ID] = breakdown

	}
	for _, breakdownType := range breakdownTypes {
		cachedBreakdownTypesById[breakdownType.ID] = breakdownType
		cachedBreakdownTypesByName[breakdownType.Name] = breakdownType
	}
	breakdownsSync.Unlock()
	logInfo("CACHING", "Cached "+strconv.Itoa(len(cachedBreakdownsById))+" breakdowns")
}

func cacheWorkShifts(db *gorm.DB) {
	var workShifts []database.Workshift
	var workplaceWorkShifts []database.WorkplaceWorkshift
	db.Find(&workShifts)
	db.Find(&workplaceWorkShifts)
	workShiftsSync.Lock()
	cachedWorkShiftsById = map[uint]database.Workshift{}
	cachedWorkplaceWorkShiftsById = map[uint]database.WorkplaceWorkshift{}
	for _, workshift := range workShifts {
		cachedWorkShiftsById[workshift.ID] = workshift

	}
	for _, workplaceWorkshift := range workplaceWorkShifts {
		cachedWorkplaceWorkShiftsById[workplaceWorkshift.ID] = workplaceWorkshift
	}
	workShiftsSync.Unlock()
	logInfo("CACHING", "Cached "+strconv.Itoa(len(cachedWorkShiftsById))+" workShifts")
	logInfo("CACHING", "Cached "+strconv.Itoa(len(cachedWorkplaceWorkShiftsById))+" workplace workShifts")
}

func cacheStates(db *gorm.DB) {
	var states []database.State
	db.Find(&states)
	statesSync.Lock()
	for _, state := range states {
		cachedStatesById[state.ID] = state
	}
	statesSync.Unlock()
	logInfo("CACHING", "Cached "+strconv.Itoa(len(cachedStatesById))+" states")
}

func cacheParts(db *gorm.DB) {
	var parts []database.Part
	db.Find(&parts)
	partsSync.Lock()
	for _, part := range parts {
		cachedPartsById[part.ID] = part

	}
	partsSync.Unlock()
	logInfo("CACHING", "Cached "+strconv.Itoa(len(cachedPartsById))+" parts")
}

func cacheProducts(db *gorm.DB) {
	var products []database.Product
	db.Find(&products)
	productsSync.Lock()
	for _, product := range products {
		cachedProductsById[product.ID] = product
		cachedProductsByName[product.Name] = product

	}
	productsSync.Unlock()
	logInfo("CACHING", "Cached "+strconv.Itoa(len(cachedProductsById))+" products")
}

func cacheAlarms(db *gorm.DB) {
	var alarms []database.Alarm
	db.Find(&alarms)
	alarmsSync.Lock()
	for _, alarm := range alarms {
		cachedAlarmsById[alarm.ID] = alarm
	}
	alarmsSync.Unlock()
	logInfo("CACHING", "Cached "+strconv.Itoa(len(cachedAlarmsById))+" alarms")
}
