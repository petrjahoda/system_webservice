package main

import (
	"github.com/petrjahoda/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strconv"
	"sync"
	"time"
)

var cachedUsers = map[string]database.User{}
var cachedLocales = map[string]database.Locale{}
var cachedCompanyName string
var cachedUsersSync sync.Mutex
var cachedLocalesSync sync.Mutex
var cachedCompanyNameSync sync.Mutex

func cacheUsers() {
	for {
		db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
		sqlDB, _ := db.DB()
		if err != nil {
			logError("MAIN", "Problem opening database: "+err.Error())
			return
		}
		var users []database.User
		db.Where("email is not null").Find(&users)
		cachedUsersSync.Lock()
		for _, user := range users {
			cachedUsers[user.Email] = user
		}
		cachedUsersSync.Unlock()
		logInfo("MAIN", "Cached "+strconv.Itoa(len(cachedUsers))+" users")
		_ = sqlDB.Close()
		time.Sleep(1 * time.Minute)
	}
}
func cacheCompanyName() {
	for {
		db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
		sqlDB, _ := db.DB()
		if err != nil {
			logError("MAIN", "Problem opening database: "+err.Error())
			return
		}
		var companyName database.Setting
		db.Where("name like 'company'").Find(&companyName)
		cachedCompanyNameSync.Lock()
		cachedCompanyName = companyName.Value
		cachedCompanyNameSync.Unlock()
		_ = sqlDB.Close()
		time.Sleep(1 * time.Hour)
	}
}

func cacheLocales() {
	for {
		db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
		sqlDB, _ := db.DB()
		if err != nil {
			logError("MAIN", "Problem opening database: "+err.Error())
			return
		}
		var locales []database.Locale
		db.Find(&locales)
		cachedLocalesSync.Lock()
		for _, locale := range locales {
			cachedLocales[locale.Name] = locale
		}
		cachedLocalesSync.Unlock()
		logInfo("MAIN", "Cached "+strconv.Itoa(len(cachedLocales))+" locales")
		_ = sqlDB.Close()
		time.Sleep(10 * time.Minute)
	}
}
