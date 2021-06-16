package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/petrjahoda/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"html/template"
	"net/http"
	"time"
)

type UserSettingsDataOutput struct {
	FirstName            string
	FirstNamePrepend     string
	SecondName           string
	SecondNamePrepend    string
	SelectionName        string
	SelectionNamePrepend string
	Email                string
	EmailPrepend         string
	Password             string
	PasswordPrepend      string
	Phone                string
	PhonePrepend         string
	Locale               string
	LocalePrepend        string
	Note                 string
	NotePrepend          string
	CreatedAt            string
	CreatedAtPrepend     string
	UpdatedAt            string
	UpdatedAtPrepend     string
	Locales              []LocaleSelection
	Result               string
}

func loadUserSettings(writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS", "Loading user settings")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "ERR: Problem opening database, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Loading user settings ended with error")
		return
	}
	var user database.User
	usersByEmailSync.RLock()
	userId := cachedUsersByEmail[email].ID
	usersByEmailSync.RUnlock()
	db.Where("id = ?", userId).Find(&user)
	var locales []LocaleSelection
	locales = append(locales, LocaleSelection{LocaleName: "CsCZ", LocaleSelected: testLocaleForUser(user.Locale, "CsCZ")})
	locales = append(locales, LocaleSelection{LocaleName: "DeDE", LocaleSelected: testLocaleForUser(user.Locale, "DeDE")})
	locales = append(locales, LocaleSelection{LocaleName: "EnUS", LocaleSelected: testLocaleForUser(user.Locale, "EnUS")})
	locales = append(locales, LocaleSelection{LocaleName: "EsES", LocaleSelected: testLocaleForUser(user.Locale, "EsES")})
	locales = append(locales, LocaleSelection{LocaleName: "FrFR", LocaleSelected: testLocaleForUser(user.Locale, "FrFR")})
	locales = append(locales, LocaleSelection{LocaleName: "ItIT", LocaleSelected: testLocaleForUser(user.Locale, "ItIT")})
	locales = append(locales, LocaleSelection{LocaleName: "PlPL", LocaleSelected: testLocaleForUser(user.Locale, "PlPL")})
	locales = append(locales, LocaleSelection{LocaleName: "PtPT", LocaleSelected: testLocaleForUser(user.Locale, "PtPT")})
	locales = append(locales, LocaleSelection{LocaleName: "SkSK", LocaleSelected: testLocaleForUser(user.Locale, "SkSK")})
	locales = append(locales, LocaleSelection{LocaleName: "RuRU", LocaleSelected: testLocaleForUser(user.Locale, "RuRU")})
	data := UserSettingsDataOutput{
		FirstName:         user.FirstName,
		FirstNamePrepend:  getLocale(email, "first-name"),
		SecondName:        user.SecondName,
		SecondNamePrepend: getLocale(email, "second-name"),
		LocalePrepend:     getLocale(email, "locale"),
		Email:             user.Email,
		EmailPrepend:      getLocale(email, "email"),
		Password:          "",
		PasswordPrepend:   getLocale(email, "password"),
		Phone:             user.Phone,
		PhonePrepend:      getLocale(email, "phone"),
		Note:              user.Note,
		NotePrepend:       getLocale(email, "note-name"),
		CreatedAt:         user.CreatedAt.Format("2006-01-02T15:04:05"),
		CreatedAtPrepend:  getLocale(email, "created-at"),
		UpdatedAt:         user.UpdatedAt.Format("2006-01-02T15:04:05"),
		UpdatedAtPrepend:  getLocale(email, "updated-at"),
		Locales:           locales,
	}
	tmpl, err := template.ParseFiles("./html/settings-user.html")
	if err != nil {
		logError("SETTINGS", "Problem parsing html file: "+err.Error())
		var responseData OrdersSettingsDataOutput
		responseData.Result = "ERR: Problem parsing html file: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
	} else {
		data.Result = "INF: User settings processed in " + time.Since(timer).String()
		_ = tmpl.Execute(writer, data)
		logInfo("SETTINGS", "User settings loaded in "+time.Since(timer).String())
	}
}

func saveUserSettings(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	timer := time.Now()
	logInfo("SETTINGS", "Saving user settings started")
	email, _, _ := request.BasicAuth()
	var data UserDetailsDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "ERR: Error parsing data, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Saving user settings ended with error")
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
		logInfo("SETTINGS", "Saving user settings ended with error")
		return
	}
	var user database.User
	usersByEmailSync.RLock()
	userId := cachedUsersByEmail[email].ID
	usersByEmailSync.RUnlock()
	db.Where("id=?", userId).Find(&user)
	user.FirstName = data.FirstName
	user.SecondName = data.SecondName
	user.Email = data.Email
	user.Phone = data.Phone
	user.Note = data.Note
	user.Locale = data.Locale
	if len(data.Password) > 0 {
		user.Password = hashPasswordFromString([]byte(data.Password))
	}
	result := db.Save(&user)
	cacheUsers(db)
	if result.Error != nil {
		var responseData TableOutput
		responseData.Result = "ERR: User settings not saved: " + result.Error.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logError("SETTINGS", "User settings for "+user.FirstName+" "+user.SecondName+" not saved: "+result.Error.Error())
	} else {
		var responseData TableOutput
		responseData.Result = "INF: User settings saved in " + time.Since(timer).String()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "User settings for "+user.FirstName+" "+user.SecondName+" saved in "+time.Since(timer).String())
	}
}
