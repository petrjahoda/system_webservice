package main

import (
	"encoding/json"
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
	logInfo("UPDATE", "Updating user settings for "+cachedUsersByEmail[email].FirstName+" "+cachedUsersByEmail[email].SecondName)
	logInfo("UPDATE", "Settings: "+key+", "+value)
	timer := time.Now()
	settings := cachedUserWebSettings[email]
	if settings != nil {
		settings[key] = value
		userSettingsSync.Lock()
		cachedUserWebSettings[email] = settings
		userSettingsSync.Unlock()
		logInfo("UPDATE", "Updated "+key+" to "+value+" in "+time.Since(timer).String())
	}
}

func updateUserWebSettingsFromWeb(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	email, _, _ := request.BasicAuth()
	logInfo("UPDATE", "Updating user web settings for "+cachedUsersByEmail[email].FirstName+" "+cachedUsersByEmail[email].SecondName)
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
	if len(email) == 0 {
		logInfo("UPDATE", "Updating user web settings for parsed email "+cachedUsersByEmail[email].FirstName+" "+cachedUsersByEmail[email].SecondName)
		email = data.Email
	}
	logInfo("UPDATE", "Settings: "+data.Key+", "+data.Value)
	settings := cachedUserWebSettings[email]
	if settings != nil {
		settings[data.Key] = data.Value
		userSettingsSync.Lock()
		cachedUserWebSettings[email] = settings
		userSettingsSync.Unlock()
		logInfo("UPDATE", "Updated "+data.Key+" to "+data.Value+" in "+time.Since(timer).String())
	}
}
