package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

type UserWebSettingsInput struct {
	Email string
	Key   string
	Value string
}

type UserOutput struct {
	Result string
}

func updateUserWebSettings(email string, key string, value string) {
	logInfo("UPDATE", "Updating user settings for "+email)
	logInfo("UPDATE", "Settings: "+key+", "+value)
	timer := time.Now()
	userWebSettingsSync.RLock()
	settings := cachedUserWebSettings[email]
	userWebSettingsSync.RUnlock()
	if settings != nil {
		settings[key] = value
		userWebSettingsSync.Lock()
		cachedUserWebSettings[email] = settings
		userWebSettingsSync.Unlock()
		logInfo("UPDATE", "Updated "+key+" to "+value+" in "+time.Since(timer).String())
	}
}

func updateUserWebSettingsFromWeb(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	email, _, _ := request.BasicAuth()
	logInfo("UPDATE", "Updating user web settings for "+email)
	timer := time.Now()
	var data UserWebSettingsInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("UPDATE", "Error parsing data: "+err.Error())
		var responseData UserOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("UPDATE", "Loading settings ended with error")
		return
	}
	fmt.Println(data.Email)
	if len(email) == 0 {
		email = data.Email
		logInfo("UPDATE", "Updating user web settings for parsed email "+email)
	}
	logInfo("UPDATE", "Settings: "+data.Key+", "+data.Value)
	userWebSettingsSync.RLock()
	settings := cachedUserWebSettings[email]
	userWebSettingsSync.RUnlock()
	if settings != nil {
		settings[data.Key] = data.Value
		userWebSettingsSync.Lock()
		cachedUserWebSettings[email] = settings
		userWebSettingsSync.Unlock()
		logInfo("UPDATE", "Updated "+data.Key+" to "+data.Value+" in "+time.Since(timer).String())
	}
}
