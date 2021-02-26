package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/petrjahoda/database"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func updateProgramVersion() {
	logInfo("SYSTEM", "Writing program version into settings")
	timer := time.Now()
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SYSTEM", "Problem opening database: "+err.Error())
		return
	}
	var existingSettings database.Setting
	db.Where("name=?", serviceName).Find(&existingSettings)
	existingSettings.Name = serviceName
	existingSettings.Value = version
	db.Save(&existingSettings)
	logInfo("SYSTEM", "Program version written into settings in "+time.Since(timer).String())
}

func basicAuth(h httprouter.Handle) httprouter.Handle {
	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		email, password, hasAuth := request.BasicAuth()
		usersSync.Lock()
		user, userFound := cachedUsersByEmail[email]
		usersSync.Unlock()
		userMatchesPassword := comparePasswords(user.Password, []byte(password))
		if hasAuth && userFound && userMatchesPassword {
			h(writer, request, params)
		} else {
			writer.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			http.Error(writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	}
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		logError("SYSTEM", "Passwords not matching")
		return false
	}
	return true
}
