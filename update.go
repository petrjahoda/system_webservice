package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

type UserWebSettingsInput struct {
	Key   string
	Value string
}

type UserOutput struct {
	Result string
}

func updateUserWebSettings(email string, key string, value string) {
	logInfo("UPDATE", "Updating user settings for "+cachedUsersByEmail[email].FirstName+" "+cachedUsersByEmail[email].SecondName)
	timer := time.Now()
	userSettingsSync.Lock()
	settings := cachedUserWebSettings[email]
	settings[key] = value
	cachedUserWebSettings[email] = settings
	userSettingsSync.Unlock()
	logInfo("UPDATE", "User settings updated in "+time.Since(timer).String())
}

func updateUserWebSettingsFromWeb(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	email, _, _ := request.BasicAuth()
	logInfo("UPDATE", "Updating user settings for "+cachedUsersByEmail[email].FirstName+" "+cachedUsersByEmail[email].SecondName)
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
	userSettingsSync.Lock()
	settings := cachedUserWebSettings[email]
	settings[data.Key] = data.Value
	cachedUserWebSettings[email] = settings
	userSettingsSync.Unlock()
	logInfo("UPDATE", "Updated "+data.Key+" to "+data.Value+" in "+time.Since(timer).String())
}
