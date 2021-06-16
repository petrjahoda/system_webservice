package main

import (
	"database/sql"
	"github.com/petrjahoda/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sync"
	"time"
)

var (
	cachedCompanyName string
	companyNameSync   sync.RWMutex
)

var (
	cachedSoftwareName string
	softwareNameSync   sync.RWMutex
)

var (
	cachedLocation string
	locationSync   sync.RWMutex
)

var (
	cachedLocalesByName = map[string]database.Locale{}
	localedByNameSync   sync.RWMutex
)

var (
	cachedLocales = map[string]string{}
	localesSync   sync.RWMutex
)

var (
	cachedAlarmsById = map[uint]database.Alarm{}
	alarmByIdSync    sync.RWMutex
)

var (
	cachedBreakdownsById = map[uint]database.Breakdown{}
	breakdownByIdSync    sync.RWMutex
)

var (
	cachedBreakdownTypesByName = map[string]database.BreakdownType{}
	breakdownTypesByNameSync   sync.RWMutex
)

var (
	cachedBreakdownTypesById = map[uint]database.BreakdownType{}
	breakdownTypesByIdSync   sync.RWMutex
)

var (
	cachedDowntimeTypesByName = map[string]database.DowntimeType{}
	downtimeTypesByNameSync   sync.RWMutex
)

var (
	cachedDowntimeTypesById = map[uint]database.DowntimeType{}
	downtimeTypesByIdSync   sync.RWMutex
)
var (
	cachedDowntimesById = map[uint]database.Downtime{}
	downtimesByIdSync   sync.RWMutex
)

var (
	cachedFaultTypesByName = map[string]database.FaultType{}
	faultTypesByNameSync   sync.RWMutex
)

var (
	cachedFaultsById = map[uint]database.Fault{}
	faultsByIdSync   sync.RWMutex
)

var (
	cachedFaultTypesById = map[uint]database.FaultType{}
	faultTypesByIdSync   sync.RWMutex
)
var (
	cachedOrdersByName = map[string]database.Order{}
	ordersByNameSync   sync.RWMutex
)

var (
	cachedOrdersById = map[uint]database.Order{}
	ordersByIdSync   sync.RWMutex
)

var (
	cachedOperationsById = map[uint]database.Operation{}
	operationsByIdSync   sync.RWMutex
)

var (
	cachedPackageTypesByName = map[string]database.PackageType{}
	packageTypesByNameSync   sync.RWMutex
)

var (
	cachedPackagesById = map[uint]database.Package{}
	packagesByIdSync   sync.RWMutex
)

var (
	cachedPackageTypesById = map[uint]database.PackageType{}
	packageTypesByIdSync   sync.RWMutex
)

var (
	cachedPartsById = map[uint]database.Part{}
	partsByIdSync   sync.RWMutex
)

var (
	cachedProductsByName = map[string]database.Product{}
	productsByNameSync   sync.RWMutex
)

var (
	cachedProductsById = map[uint]database.Product{}
	productsByIdSync   sync.RWMutex
)

var (
	cachedDevicesByName = map[string]database.Device{}
	devicesByNameSync   sync.RWMutex
)

var (
	cachedDevicesById = map[uint]database.Device{}
	devicesByIdSync   sync.RWMutex
)

var (
	cachedDeviceTypesById = map[uint]database.DeviceType{}
	deviceTypesByIdSync   sync.RWMutex
)
var (
	cachedDeviceTypesByName = map[string]database.DeviceType{}
	deviceTypesByNameSync   sync.RWMutex
)
var (
	cachedDevicePortsById = map[uint]database.DevicePort{}
	devicePortsByIdSync   sync.RWMutex
)

var (
	cachedDevicePortTypesById = map[uint]database.DevicePortType{}
	devicePortTypesByIdSync   sync.RWMutex
)

var (
	cachedDevicePortTypesByName = map[string]database.DevicePortType{}
	devicePortTypesByNameSync   sync.RWMutex
)

var (
	cachedDevicePortsColorsById = map[int]string{}
	devicePortsColorsByIdSync   sync.RWMutex
)

var (
	cachedStatesById = map[uint]database.State{}
	statesByIdSync   sync.RWMutex
)

var (
	cachedWorkShiftsById = map[uint]database.Workshift{}
	workShiftsByIdSync   sync.RWMutex
)

var (
	cachedWorkplaceWorkShiftsById = map[uint]database.WorkplaceWorkshift{}
	workplaceWorkShiftsByIdSync   sync.RWMutex
)

var (
	cachedLatestCachedWorkplaceCalendarData = time.Now()
	latestCachedWorkplaceCalendarDataSync   sync.RWMutex
)

var (
	cachedLatestCachedWorkplaceConsumption = time.Now()
	latestCachedWorkplaceConsumptionSync   sync.RWMutex
)

var (
	cachedConsumptionDataByWorkplaceName = map[string]map[string]float32{}
	consumptionDataByWorkplaceNameSync   sync.RWMutex
)

var (
	cachedUserRolesByName = map[string]database.UserRole{}
	userRolesByNameSync   sync.RWMutex
)

var (
	cachedUserRolesById = map[uint]database.UserRole{}
	userRolesByIdSync   sync.RWMutex
)

var (
	cachedUserTypesByName = map[string]database.UserType{}
	userTypesByNameSync   sync.RWMutex
)

var (
	cachedUserTypesById = map[uint]database.UserType{}
	userTypesByIdSync   sync.RWMutex
)

var (
	cachedWorkplaceDevicePorts = map[string][]database.DevicePort{}
	workplaceDevicePortsSync   sync.RWMutex
)
var (
	cachedWorkplacePorts = map[string][]database.WorkplacePort{}
	workplacePortsSync   sync.RWMutex
)

var (
	cachedWorkplaceModesById = map[uint]database.WorkplaceMode{}
	workplaceModesByIdSync   sync.RWMutex
)

var (
	cachedWorkplaceModesByName = map[string]database.WorkplaceMode{}
	workplaceModesByNameSync   sync.RWMutex
)

var (
	cachedWorkplaceSectionsByName = map[string]database.WorkplaceSection{}
	workplaceSectionsByNameSync   sync.RWMutex
)

var (
	cachedWorkplaceSectionsById = map[uint]database.WorkplaceSection{}
	workplaceSectionsByIdSync   sync.RWMutex
)
var (
	cachedWorkplacesProductionRecords = map[string]map[string]time.Duration{}
	workplacesProductionRecordsSync   sync.RWMutex
)

var (
	cachedWorkplacesDowntimeRecords = map[string]map[string]time.Duration{}
	workplacesDowntimeRecordsSync   sync.RWMutex
)

var (
	cachedWorkplacesPoweroffRecords = map[string]map[string]time.Duration{}
	workplacesPoweroffRecordsSync   sync.RWMutex
)

var (
	cachedWorkplacesByName = map[string]database.Workplace{}
	workplacesByNameSync   sync.RWMutex
)

var (
	cachedWorkplacesById = map[uint]database.Workplace{}
	workplacesByIdSync   sync.RWMutex
)

var (
	cachedUsersByEmail = map[string]database.User{}
	usersByEmailSync   sync.RWMutex
)

var (
	cachedUsersById = map[uint]database.User{}
	usersByIdSync   sync.RWMutex
)
var (
	cachedUserWebSettings = map[string]map[string]string{}
	userWebSettingsSync   sync.RWMutex
)

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
	locationSync.RLock()
	loc, _ := time.LoadLocation(cachedLocation)
	locationSync.RUnlock()
	cacheConsumptionData(db, time.Date(time.Now().Add(-720*time.Hour).Year(), time.Now().Add(-720*time.Hour).Month(), time.Now().Add(-720*time.Hour).Day(), 0, 0, 0, 0, loc))
	logInfo("CACHING", "Initial caching done in "+time.Since(timer).String())
}

func cacheProductionData(db *gorm.DB, fromDate time.Time) {
	locationSync.RLock()
	loc, _ := time.LoadLocation(cachedLocation)
	locationSync.RUnlock()
	logInfo("CACHING", "Caching workplace production/downtime/poweroff from: "+fromDate.In(loc).String()+"  back until: "+time.Now().In(loc).String())
	var workplaces []database.Workplace
	db.Find(&workplaces)
	workplacesProductionRecordsSync.RLock()
	workplacesProductionRecords := cachedWorkplacesProductionRecords
	workplacesProductionRecordsSync.RUnlock()
	workplacesDowntimeRecordsSync.RLock()
	workplacesDowntimeRecords := cachedWorkplacesDowntimeRecords
	workplacesDowntimeRecordsSync.RUnlock()
	workplacesPoweroffRecordsSync.RLock()
	workplacesPoweroffRecords := cachedWorkplacesPoweroffRecords
	workplacesPoweroffRecordsSync.RUnlock()
	for _, workplace := range workplaces {
		if workplacesProductionRecords[workplace.Name] == nil {
			workplacesProductionRecords[workplace.Name] = make(map[string]time.Duration)
		}
		if workplacesDowntimeRecords[workplace.Name] == nil {
			workplacesDowntimeRecords[workplace.Name] = make(map[string]time.Duration)
		}
		if workplacesPoweroffRecords[workplace.Name] == nil {
			workplacesPoweroffRecords[workplace.Name] = make(map[string]time.Duration)
		}
		productionRecords := workplacesProductionRecords[workplace.Name]
		downtimeRecords := workplacesDowntimeRecords[workplace.Name]
		poweroffRecords := workplacesPoweroffRecords[workplace.Name]
		tempDate := fromDate
		for tempDate.Before(time.Now()) {
			productionRecords[tempDate.Format("2006-01-02")] = 0
			downtimeRecords[tempDate.Format("2006-01-02")] = 0
			poweroffRecords[tempDate.Format("2006-01-02")] = 0
			tempDate = tempDate.Add(24 * time.Hour)
		}
		var stateRecords []database.StateRecord
		latestCachedWorkplaceCalendarDataSync.Lock()
		cachedLatestCachedWorkplaceCalendarData = time.Now().In(loc)
		latestCachedWorkplaceCalendarDataSync.Unlock()
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
		workplacesProductionRecordsSync.Lock()
		cachedWorkplacesProductionRecords[workplace.Name] = productionRecords
		workplacesProductionRecordsSync.Unlock()
		workplacesDowntimeRecordsSync.Lock()
		cachedWorkplacesDowntimeRecords[workplace.Name] = downtimeRecords
		workplacesDowntimeRecordsSync.Unlock()
		workplacesPoweroffRecordsSync.Lock()
		cachedWorkplacesPoweroffRecords[workplace.Name] = poweroffRecords
		workplacesPoweroffRecordsSync.Unlock()
	}
	logInfo("CACHING", "Production/downtime/poweroff workplace records cached")
}

func cacheConsumptionData(db *gorm.DB, date time.Time) {
	locationSync.RLock()
	loc, _ := time.LoadLocation(cachedLocation)
	locationSync.RUnlock()
	latestCachedWorkplaceConsumptionSync.Lock()
	cachedLatestCachedWorkplaceConsumption = time.Now().In(loc)
	latestCachedWorkplaceConsumptionSync.Unlock()
	var workplaces []database.Workplace
	db.Find(&workplaces)
	for _, workplace := range workplaces {
		consumptionDataByWorkplaceNameSync.Lock()
		if cachedConsumptionDataByWorkplaceName[workplace.Name] == nil {
			cachedConsumptionDataByWorkplaceName[workplace.Name] = make(map[string]float32)
		}
		consumptionDataByWorkplaceNameSync.Unlock()
		consumptionDataByWorkplaceNameSync.RLock()
		tempConsumptionData := cachedConsumptionDataByWorkplaceName[workplace.Name]
		consumptionDataByWorkplaceNameSync.RUnlock()
		tempDate := date
		for tempDate.Before(time.Now()) {
			tempConsumptionData[tempDate.Format("2006-01-02")] = 0
			tempDate = tempDate.Add(24 * time.Hour)
		}
		tempDate = date
		workplacePortsSync.RLock()
		workplacePorts := cachedWorkplacePorts
		workplacePortsSync.RUnlock()
		for _, port := range workplacePorts[workplace.Name] {
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
				consumptionDataByWorkplaceNameSync.Lock()
				cachedConsumptionDataByWorkplaceName[workplace.Name] = tempConsumptionData
				consumptionDataByWorkplaceNameSync.Unlock()
			}
		}
	}
	logInfo("CACHING", "Consumption cached")
}

func cacheWorkplacePorts(db *gorm.DB) {
	var workplacePorts []database.WorkplacePort
	db.Find(&workplacePorts)
	for _, workplacePort := range workplacePorts {
		devicePortsColorsByIdSync.Lock()
		cachedDevicePortsColorsById[workplacePort.DevicePortID] = workplacePort.Color
		devicePortsColorsByIdSync.Unlock()
	}
	logInfo("CACHING", "Workplace port colors cached")
}

func cacheWorkplaceDevicePorts(db *gorm.DB) {
	var allWorkplaces []database.Workplace
	db.Find(&allWorkplaces)

	for _, workplace := range allWorkplaces {
		var allWorkplacePorts []database.WorkplacePort
		var allDevicePorts []database.DevicePort
		db.Where("workplace_id = ?", workplace.ID).Find(&allWorkplacePorts)
		for _, workplacePort := range allWorkplacePorts {
			var devicePort database.DevicePort
			db.Where("id = ?", workplacePort.DevicePortID).Find(&devicePort)
			allDevicePorts = append(allDevicePorts, devicePort)
		}
		workplaceDevicePortsSync.Lock()
		cachedWorkplaceDevicePorts[workplace.Name] = allDevicePorts
		workplaceDevicePortsSync.Unlock()
		workplacePortsSync.Lock()
		cachedWorkplacePorts[workplace.Name] = allWorkplacePorts
		workplacePortsSync.Unlock()
	}

	logInfo("CACHING", "Workplace ports cached")
}

func cacheSystemSettings(db *gorm.DB) {
	var companyName database.Setting
	db.Where("name like 'company'").Find(&companyName)
	companyNameSync.Lock()
	cachedCompanyName = companyName.Value
	companyNameSync.Unlock()
	logInfo("CACHING", "Company name cached: "+companyName.Value)

	var softwareName database.Setting
	db.Where("name like 'software'").Find(&softwareName)
	softwareNameSync.Lock()
	cachedSoftwareName = softwareName.Value
	softwareNameSync.Unlock()
	logInfo("CACHING", "Software name cached: "+softwareName.Value)

	var timezone database.Setting
	db.Where("name=?", "timezone").Find(&timezone)
	locationSync.Lock()
	cachedLocation = timezone.Value
	locationSync.Unlock()
	logInfo("CACHING", "Timezone cached: "+timezone.Value)
}

func cacheLocales(db *gorm.DB) {
	var locales []database.Locale
	db.Find(&locales)
	for _, locale := range locales {
		localedByNameSync.Lock()
		cachedLocalesByName[locale.Name] = locale
		localedByNameSync.Unlock()
	}
	localesSync.Lock()
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
	logInfo("CACHING", "Locales cached")
}

func cacheOrders(db *gorm.DB) {
	var orders []database.Order
	db.Find(&orders)
	for _, order := range orders {
		ordersByIdSync.Lock()
		cachedOrdersById[order.ID] = order
		ordersByIdSync.Unlock()
		ordersByNameSync.Lock()
		cachedOrdersByName[order.Name] = order
		ordersByNameSync.Unlock()
	}
	logInfo("CACHING", "Orders cached")
}

func cacheOperations(db *gorm.DB) {
	var operations []database.Operation
	db.Find(&operations)
	for _, operation := range operations {
		operationsByIdSync.Lock()
		cachedOperationsById[operation.ID] = operation
		operationsByIdSync.Unlock()
	}
	logInfo("CACHING", "Operations cached")
}

func cacheWorkplaces(db *gorm.DB) {
	var workplaceModes []database.WorkplaceMode
	var workplaceSections []database.WorkplaceSection
	var workplaces []database.Workplace
	db.Find(&workplaces)
	db.Find(&workplaceModes)
	db.Find(&workplaceSections)
	for _, workplaceMode := range workplaceModes {
		workplaceModesByIdSync.Lock()
		cachedWorkplaceModesById[workplaceMode.ID] = workplaceMode
		workplaceModesByIdSync.Unlock()
		workplaceModesByNameSync.Lock()
		cachedWorkplaceModesByName[workplaceMode.Name] = workplaceMode
		workplaceModesByNameSync.Unlock()

	}
	for _, workplaceSection := range workplaceSections {
		workplaceSectionsByIdSync.Lock()
		cachedWorkplaceSectionsById[workplaceSection.ID] = workplaceSection
		workplaceSectionsByIdSync.Unlock()
		workplaceSectionsByNameSync.Lock()
		cachedWorkplaceSectionsByName[workplaceSection.Name] = workplaceSection
		workplaceSectionsByNameSync.Unlock()

	}
	for _, workplace := range workplaces {
		workplacesByIdSync.Lock()
		cachedWorkplacesById[workplace.ID] = workplace
		workplacesByIdSync.Unlock()
		workplacesByNameSync.Lock()
		cachedWorkplacesByName[workplace.Name] = workplace
		workplacesByNameSync.Unlock()

	}
	logInfo("CACHING", "Workplaces cached")
	logInfo("CACHING", "Workplace modes cached")
	logInfo("CACHING", "Workplace sections cached")
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
	for _, device := range devices {
		devicesByIdSync.Lock()
		cachedDevicesById[device.ID] = device
		devicesByIdSync.Unlock()
		devicesByNameSync.Lock()
		cachedDevicesByName[device.Name] = device
		devicesByNameSync.Unlock()
	}
	for _, deviceType := range deviceTypes {
		deviceTypesByIdSync.Lock()
		cachedDeviceTypesById[deviceType.ID] = deviceType
		deviceTypesByIdSync.Unlock()
		deviceTypesByNameSync.Lock()
		cachedDeviceTypesByName[deviceType.Name] = deviceType
		deviceTypesByNameSync.Unlock()
	}
	for _, devicePortType := range devicePortTypes {
		devicePortTypesByIdSync.Lock()
		cachedDevicePortTypesById[devicePortType.ID] = devicePortType
		devicePortTypesByIdSync.Unlock()
		devicePortTypesByNameSync.Lock()
		cachedDevicePortTypesByName[devicePortType.Name] = devicePortType
		devicePortTypesByNameSync.Unlock()
	}
	for _, devicePort := range devicePorts {
		devicePortsByIdSync.Lock()
		cachedDevicePortsById[devicePort.ID] = devicePort
		devicePortsByIdSync.Unlock()
	}
	logInfo("CACHING", "Devices cached")
	logInfo("CACHING", "Device ports cached")
	logInfo("CACHING", "Device types cached")
	logInfo("CACHING", "Device port types cached")
}

func cacheUsers(db *gorm.DB) {
	var users []database.User
	var userTypes []database.UserType
	var userRoles []database.UserRole
	db.Find(&users)
	db.Find(&userTypes)
	db.Find(&userRoles)
	for _, user := range users {
		if len(user.Email) > 0 {
			usersByEmailSync.Lock()
			cachedUsersByEmail[user.Email] = user
			usersByEmailSync.Unlock()
			userWebSettingsSync.RLock()
			_, userWebCached := cachedUserWebSettings[user.Email]
			userWebSettingsSync.RUnlock()
			if !userWebCached {
				data := map[string]string{}
				userWebSettingsSync.Lock()
				cachedUserWebSettings[user.Email] = data
				userWebSettingsSync.Unlock()
			}

		}
		usersByIdSync.Lock()
		cachedUsersById[user.ID] = user
		usersByIdSync.Unlock()
	}
	for _, userType := range userTypes {
		userTypesByIdSync.Lock()
		cachedUserTypesById[userType.ID] = userType
		userTypesByIdSync.Unlock()
		userTypesByNameSync.Lock()
		cachedUserTypesByName[userType.Name] = userType
		userTypesByNameSync.Unlock()
	}
	for _, userRole := range userRoles {
		userRolesByIdSync.Lock()
		cachedUserRolesById[userRole.ID] = userRole
		userRolesByIdSync.Unlock()
		userRolesByNameSync.Lock()
		cachedUserRolesByName[userRole.Name] = userRole
		userRolesByNameSync.Unlock()
	}
	logInfo("CACHING", "Users cached")
	logInfo("CACHING", "User settings cached")
}

func cachePackages(db *gorm.DB) {
	var packages []database.Package
	var packageTypes []database.PackageType
	db.Find(&packages)
	db.Find(&packageTypes)
	for _, onePackage := range packages {
		packagesByIdSync.Lock()
		cachedPackagesById[onePackage.ID] = onePackage
		packagesByIdSync.Unlock()

	}
	for _, packageType := range packageTypes {
		packageTypesByIdSync.Lock()
		cachedPackageTypesById[packageType.ID] = packageType
		packageTypesByIdSync.Unlock()
		packageTypesByNameSync.Lock()
		cachedPackageTypesByName[packageType.Name] = packageType
		packageTypesByNameSync.Unlock()
	}
	logInfo("CACHING", "Packages cached")
}

func cacheFaults(db *gorm.DB) {
	var faults []database.Fault
	var faultTypes []database.FaultType
	db.Find(&faults)
	db.Find(&faultTypes)
	for _, fault := range faults {
		faultsByIdSync.Lock()
		cachedFaultsById[fault.ID] = fault
		faultsByIdSync.Unlock()

	}
	for _, faultType := range faultTypes {
		faultTypesByIdSync.Lock()
		cachedFaultTypesById[faultType.ID] = faultType
		faultTypesByIdSync.Unlock()
		faultTypesByNameSync.Lock()
		cachedFaultTypesByName[faultType.Name] = faultType
		faultTypesByNameSync.Unlock()
	}
	logInfo("CACHING", "Faults cached")
}

func cacheDowntimes(db *gorm.DB) {
	var downtimes []database.Downtime
	var downtimeTypes []database.DowntimeType
	db.Find(&downtimes)
	db.Find(&downtimeTypes)
	for _, downtime := range downtimes {
		downtimesByIdSync.Lock()
		cachedDowntimesById[downtime.ID] = downtime
		downtimesByIdSync.Unlock()
	}
	for _, downtimeType := range downtimeTypes {
		downtimeTypesByIdSync.Lock()
		cachedDowntimeTypesById[downtimeType.ID] = downtimeType
		downtimeTypesByIdSync.Unlock()
		downtimeTypesByNameSync.Lock()
		cachedDowntimeTypesByName[downtimeType.Name] = downtimeType
		downtimeTypesByNameSync.Unlock()
	}
	logInfo("CACHING", "Downtimes cached")
}

func cacheBreakdowns(db *gorm.DB) {
	var breakdowns []database.Breakdown
	var breakdownTypes []database.BreakdownType
	db.Find(&breakdowns)
	db.Find(&breakdownTypes)
	for _, breakdown := range breakdowns {
		breakdownByIdSync.Lock()
		cachedBreakdownsById[breakdown.ID] = breakdown
		breakdownByIdSync.Unlock()

	}
	for _, breakdownType := range breakdownTypes {
		breakdownTypesByIdSync.Lock()
		cachedBreakdownTypesById[breakdownType.ID] = breakdownType
		breakdownTypesByIdSync.Unlock()
		breakdownTypesByNameSync.Lock()
		cachedBreakdownTypesByName[breakdownType.Name] = breakdownType
		breakdownTypesByNameSync.Unlock()
	}
	logInfo("CACHING", "Breakdowns cached")
}

func cacheWorkShifts(db *gorm.DB) {
	var workShifts []database.Workshift
	var workplaceWorkShifts []database.WorkplaceWorkshift
	db.Find(&workShifts)
	db.Find(&workplaceWorkShifts)
	for _, workshift := range workShifts {
		workShiftsByIdSync.Lock()
		cachedWorkShiftsById[workshift.ID] = workshift
		workShiftsByIdSync.Unlock()

	}
	for _, workplaceWorkshift := range workplaceWorkShifts {
		workplaceWorkShiftsByIdSync.Lock()
		cachedWorkplaceWorkShiftsById[workplaceWorkshift.ID] = workplaceWorkshift
		workplaceWorkShiftsByIdSync.Unlock()
	}
	logInfo("CACHING", "Workshift cached")
	logInfo("CACHING", "Workplace workshifts cached")
}

func cacheStates(db *gorm.DB) {
	var states []database.State
	db.Find(&states)
	for _, state := range states {
		statesByIdSync.Lock()
		cachedStatesById[state.ID] = state
		statesByIdSync.Unlock()
	}
	logInfo("CACHING", "States cached")
}

func cacheParts(db *gorm.DB) {
	var parts []database.Part
	db.Find(&parts)
	for _, part := range parts {
		partsByIdSync.Lock()
		cachedPartsById[part.ID] = part
		partsByIdSync.Unlock()
	}
	logInfo("CACHING", "Parts cached")
}

func cacheProducts(db *gorm.DB) {
	var products []database.Product
	db.Find(&products)
	for _, product := range products {
		productsByIdSync.Lock()
		cachedProductsById[product.ID] = product
		productsByIdSync.Unlock()
		productsByNameSync.Lock()
		cachedProductsByName[product.Name] = product
		productsByNameSync.Unlock()
	}
	logInfo("CACHING", "Products cached")
}

func cacheAlarms(db *gorm.DB) {
	var alarms []database.Alarm
	db.Find(&alarms)
	for _, alarm := range alarms {
		alarmByIdSync.Lock()
		cachedAlarmsById[alarm.ID] = alarm
		alarmByIdSync.Unlock()
	}
	logInfo("CACHING", "Alarm cached")
}
