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
var cachedLocalesByName = map[string]database.Locale{}
var cachedCompanyName string
var location string
var cachedOrdersById = map[uint]database.Order{}
var cachedOperationsById = map[uint]database.Operation{}
var cachedWorkplaceModesById = map[uint]database.WorkplaceMode{}
var cachedWorkshiftsById = map[uint]database.Workshift{}

var usersSync sync.RWMutex
var userSettingsSync sync.RWMutex
var localesSync sync.RWMutex
var companyNameSync sync.RWMutex
var workplacesSync sync.RWMutex
var ordersSync sync.RWMutex
var operationsSync sync.RWMutex
var workplaceModesSync sync.RWMutex
var workshiftsSync sync.RWMutex

type userSettings struct {
	menuState          string
	sectionStates      []sectionState
	dataSelection      string
	selectedWorkplaces []string
}
type sectionState struct {
	section string
	state   string
}

func cacheData() {
	for {
		db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
		sqlDB, _ := db.DB()
		if err != nil {
			logError("MAIN", "Problem opening database: "+err.Error())
			return
		}
		var users []database.User
		db.Find(&users)
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
		usersSync.Unlock()
		userSettingsSync.Unlock()
		logInfo("MAIN", "Cached "+strconv.Itoa(len(cachedUsersByEmail))+" users")
		logInfo("MAIN", "Cached "+strconv.Itoa(len(cachedUserSettings))+" user settings")

		var workplaces []database.Workplace
		db.Find(&workplaces)
		workplacesSync.Lock()
		for _, workplace := range workplaces {
			cachedWorkplacesById[workplace.ID] = workplace

		}
		workplacesSync.Unlock()
		logInfo("MAIN", "Cached "+strconv.Itoa(len(cachedWorkplacesById))+" workplaces")

		var companyName database.Setting
		db.Where("name like 'company'").Find(&companyName)
		companyNameSync.Lock()
		cachedCompanyName = companyName.Value
		companyNameSync.Unlock()
		logInfo("MAIN", "Cached company name")

		logInfo("MAIN", "Reading timezone from database")

		var timezone database.Setting
		db.Where("name=?", "timezone").Find(&timezone)
		companyNameSync.Lock()
		location = timezone.Value
		companyNameSync.Unlock()
		logInfo("MAIN", "Cached timezone")

		var locales []database.Locale
		db.Find(&locales)
		localesSync.Lock()
		for _, locale := range locales {
			cachedLocalesByName[locale.Name] = locale
		}
		localesSync.Unlock()
		logInfo("MAIN", "Cached "+strconv.Itoa(len(cachedLocalesByName))+" locales")

		var orders []database.Order
		db.Find(&orders)
		ordersSync.Lock()
		for _, order := range orders {
			cachedOrdersById[order.ID] = order

		}
		ordersSync.Unlock()
		logInfo("MAIN", "Cached "+strconv.Itoa(len(cachedOrdersById))+" orders")

		var operations []database.Operation
		db.Find(&operations)
		operationsSync.Lock()
		for _, operation := range operations {
			cachedOperationsById[operation.ID] = operation

		}
		operationsSync.Unlock()
		logInfo("MAIN", "Cached "+strconv.Itoa(len(cachedOperationsById))+" operations")

		var workplaceModes []database.WorkplaceMode
		db.Find(&workplaceModes)
		workplaceModesSync.Lock()
		for _, workplaceMode := range workplaceModes {
			cachedWorkplaceModesById[workplaceMode.ID] = workplaceMode

		}
		workplaceModesSync.Unlock()
		logInfo("MAIN", "Cached "+strconv.Itoa(len(cachedWorkplaceModesById))+" workplace modes")

		var workshifts []database.Workshift
		db.Find(&workshifts)
		workshiftsSync.Lock()
		for _, workshift := range workshifts {
			cachedWorkshiftsById[workshift.ID] = workshift

		}
		workshiftsSync.Unlock()
		logInfo("MAIN", "Cached "+strconv.Itoa(len(cachedWorkshiftsById))+" workshifts")

		_ = sqlDB.Close()
		time.Sleep(1 * time.Minute)
	}
}
