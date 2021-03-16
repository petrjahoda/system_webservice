package main

import (
	"github.com/petrjahoda/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strconv"
	"sync"
	"time"
)

var cachedUsersByEmail = map[string]database.User{}
var cachedUsersById = map[uint]database.User{}
var cachedUserSettings = map[string]userSettings{}
var cachedWorkplacesById = map[uint]database.Workplace{}
var cachedWorkplacesByName = map[string]database.Workplace{}
var cachedLocalesByName = map[string]database.Locale{}
var cachedLocales = map[string]string{}
var cachedCompanyName string
var location string
var cachedOrdersById = map[uint]database.Order{}
var cachedOrdersByName = map[string]database.Order{}
var cachedOperationsById = map[uint]database.Operation{}
var cachedProductsById = map[uint]database.Product{}
var cachedProductsByName = map[string]database.Product{}
var cachedWorkplaceModesById = map[uint]database.WorkplaceMode{}
var cachedWorkshiftsById = map[uint]database.Workshift{}
var cachedAlarmsById = map[uint]database.Alarm{}
var cachedBreakdownsById = map[uint]database.Breakdown{}
var cachedBreakdownTypesByName = map[string]database.BreakdownType{}
var cachedBreakdownTypesById = map[uint]database.BreakdownType{}
var cachedDowntimesById = map[uint]database.Downtime{}
var cachedDowntimeTypesByName = map[string]database.DowntimeType{}
var cachedDowntimeTypesById = map[uint]database.DowntimeType{}
var cachedFaultsById = map[uint]database.Fault{}
var cachedFaultTypesByName = map[string]database.FaultType{}
var cachedFaultTypesById = map[uint]database.FaultType{}
var cachedPackagesById = map[uint]database.Package{}
var cachedPackageTypesByName = map[string]database.PackageType{}
var cachedPackageTypesById = map[uint]database.PackageType{}
var cachedPartsById = map[uint]database.Part{}
var cachedStatesById = map[uint]database.State{}
var cachedWorkplaceDevicePorts = map[string][]database.DevicePort{}
var cachedDevicePortsColorsById = map[int]string{}
var cachedUserRolesById = map[uint]database.UserRole{}
var cachedUserRolesByName = map[string]database.UserRole{}
var cachedUserTypesById = map[uint]database.UserType{}
var cachedUserTypesByName = map[string]database.UserType{}

var usersSync sync.RWMutex
var userSettingsSync sync.RWMutex
var localesSync sync.RWMutex
var companyNameSync sync.RWMutex
var workplacesSync sync.RWMutex
var ordersSync sync.RWMutex
var operationsSync sync.RWMutex
var productsSync sync.RWMutex
var workplaceModesSync sync.RWMutex
var workshiftsSync sync.RWMutex
var alarmsSync sync.RWMutex
var breakdownsSync sync.RWMutex
var downtimesSync sync.RWMutex
var faultsSync sync.RWMutex
var packagesSync sync.RWMutex
var partsSync sync.RWMutex
var statesSync sync.RWMutex
var workplaceDevicePortsSync sync.RWMutex
var cachedDevicePortsColorsSync sync.RWMutex

type userSettings struct {
	menuState          string
	sectionStates      []sectionState
	dataSelection      string
	settingsSelection  string
	selectedWorkplaces []string
}
type sectionState struct {
	section string
	state   string
}

func cacheData() {
	for {
		logInfo("CHACHING", "Caching started")
		timer := time.Now()

		db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
		sqlDB, _ := db.DB()
		if err != nil {
			logError("CHACHING", "Problem opening database: "+err.Error())
			return
		}

		cacheUsers(db)

		var workplaces []database.Workplace
		db.Find(&workplaces)
		workplacesSync.Lock()
		for _, workplace := range workplaces {
			cachedWorkplacesById[workplace.ID] = workplace
			cachedWorkplacesByName[workplace.Name] = workplace

		}
		workplacesSync.Unlock()
		logInfo("CHACHING", "Cached "+strconv.Itoa(len(cachedWorkplacesById))+" workplaces")

		var companyName database.Setting
		db.Where("name like 'company'").Find(&companyName)
		companyNameSync.Lock()
		cachedCompanyName = companyName.Value
		companyNameSync.Unlock()
		logInfo("CHACHING", "Cached company name")

		var timezone database.Setting
		db.Where("name=?", "timezone").Find(&timezone)
		companyNameSync.Lock()
		location = timezone.Value
		companyNameSync.Unlock()
		logInfo("CHACHING", "Cached timezone")

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
		logInfo("CHACHING", "Cached "+strconv.Itoa(len(cachedLocalesByName))+" locales")

		var orders []database.Order
		db.Find(&orders)
		ordersSync.Lock()
		for _, order := range orders {
			cachedOrdersById[order.ID] = order
			cachedOrdersByName[order.Name] = order

		}
		ordersSync.Unlock()
		logInfo("CHACHING", "Cached "+strconv.Itoa(len(cachedOrdersById))+" orders")

		var operations []database.Operation
		db.Find(&operations)
		operationsSync.Lock()
		for _, operation := range operations {
			cachedOperationsById[operation.ID] = operation

		}
		operationsSync.Unlock()
		logInfo("CHACHING", "Cached "+strconv.Itoa(len(cachedOperationsById))+" operations")

		cacheProducts(db)

		var workplaceModes []database.WorkplaceMode
		db.Find(&workplaceModes)
		workplaceModesSync.Lock()
		for _, workplaceMode := range workplaceModes {
			cachedWorkplaceModesById[workplaceMode.ID] = workplaceMode

		}
		workplaceModesSync.Unlock()
		logInfo("CHACHING", "Cached "+strconv.Itoa(len(cachedWorkplaceModesById))+" workplace modes")

		cacheWorkshifts(db)
		cacheAlarms(db)
		cacheBreakdowns(db)
		cacheDowntimes(db)
		cacheFaults(db)
		cachePackages(db)
		cacheParts(db)
		cacheStates(db)

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

		}
		workplaceDevicePortsSync.Unlock()
		logInfo("CHACHING", "Cached "+strconv.Itoa(len(cachedWorkplaceDevicePorts))+" workplace deviceports")

		var workplacePorts []database.WorkplacePort
		db.Find(&workplacePorts)
		cachedDevicePortsColorsSync.Lock()
		for _, workplacePort := range workplacePorts {
			cachedDevicePortsColorsById[workplacePort.DevicePortID] = workplacePort.Color
		}
		cachedDevicePortsColorsSync.Unlock()
		logInfo("CHACHING", "Cached "+strconv.Itoa(len(cachedDevicePortsColorsById))+" workplace deviceport colors")

		_ = sqlDB.Close()
		logInfo("CHACHING", "Caching done in "+time.Since(timer).String())
		time.Sleep(1 * time.Minute)
	}
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
	for _, user := range users {
		if len(user.Email) > 0 {
			cachedUsersByEmail[user.Email] = user
			_, userCached := cachedUserSettings[user.Email]
			if !userCached {
				var userSettings userSettings
				cachedUserSettings[user.Email] = userSettings
			}
		}
		cachedUsersById[user.ID] = user
	}
	for _, userType := range userTypes {
		cachedUserTypesById[userType.ID] = userType
		cachedUserTypesByName[userType.Name] = userType
	}
	for _, userRole := range userRoles {
		cachedUserRolesById[userRole.ID] = userRole
		cachedUserRolesByName[userRole.Name] = userRole
	}

	usersSync.Unlock()
	userSettingsSync.Unlock()
	logInfo("CHACHING", "Cached "+strconv.Itoa(len(cachedUsersByEmail))+" users")
	logInfo("CHACHING", "Cached "+strconv.Itoa(len(cachedUserSettings))+" user settings")
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
	logInfo("CHACHING", "Cached "+strconv.Itoa(len(cachedPackagesById))+" packages")
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
	logInfo("CHACHING", "Cached "+strconv.Itoa(len(cachedFaultsById))+" faults")
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
	logInfo("CHACHING", "Cached "+strconv.Itoa(len(cachedDowntimesById))+" downtimes")
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
	logInfo("CHACHING", "Cached "+strconv.Itoa(len(cachedBreakdownsById))+" breakdowns")
}

func cacheWorkshifts(db *gorm.DB) {
	var workshifts []database.Workshift
	db.Find(&workshifts)
	workshiftsSync.Lock()
	for _, workshift := range workshifts {
		cachedWorkshiftsById[workshift.ID] = workshift

	}
	workshiftsSync.Unlock()
	logInfo("CHACHING", "Cached "+strconv.Itoa(len(cachedWorkshiftsById))+" workshifts")
}

func cacheStates(db *gorm.DB) {
	var states []database.State
	db.Find(&states)
	statesSync.Lock()
	for _, state := range states {
		cachedStatesById[state.ID] = state
	}
	statesSync.Unlock()
	logInfo("CHACHING", "Cached "+strconv.Itoa(len(cachedStatesById))+" states")
}

func cacheParts(db *gorm.DB) {
	var parts []database.Part
	db.Find(&parts)
	partsSync.Lock()
	for _, part := range parts {
		cachedPartsById[part.ID] = part

	}
	partsSync.Unlock()
	logInfo("CHACHING", "Cached "+strconv.Itoa(len(cachedPartsById))+" parts")
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
	logInfo("CHACHING", "Cached "+strconv.Itoa(len(cachedProductsById))+" products")
}

func cacheAlarms(db *gorm.DB) {
	var alarms []database.Alarm
	db.Find(&alarms)
	alarmsSync.Lock()
	for _, alarm := range alarms {
		cachedAlarmsById[alarm.ID] = alarm

	}
	alarmsSync.Unlock()
	logInfo("CHACHING", "Cached "+strconv.Itoa(len(cachedAlarmsById))+" alarms")
}
