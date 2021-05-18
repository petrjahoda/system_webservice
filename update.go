package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
	"time"
)

type UserSettingsInput struct {
	Settings string
	State    string
}

type UserWorkplacesInput struct {
	Workplaces []string
}

type UserOutput struct {
	Result string
}

func updateUserWorkplaces(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	email, _, _ := request.BasicAuth()
	logInfo("UPDATE", "Updating user workplaces for "+cachedUsersByEmail[email].FirstName+" "+cachedUsersByEmail[email].SecondName)
	timer := time.Now()
	var data UserWorkplacesInput
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
	updateUserDataSettings(email, "", "", data.Workplaces)
	var responseData UserOutput
	responseData.Result = "ok"
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(responseData)
	logInfo("UPDATE", "User workplaces updated in "+time.Since(timer).String())
}

func updateUserSettings(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	email, _, _ := request.BasicAuth()
	logInfo("UPDATE", "Updating user settings for "+cachedUsersByEmail[email].FirstName+" "+cachedUsersByEmail[email].SecondName)
	timer := time.Now()
	var data UserSettingsInput
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
	switch data.Settings {
	case "menu":
		{
			userSettingsSync.Lock()
			userSettings := cachedUserSettings[email]
			if data.State == "false" {
				userSettings.menuState = "compacted js-compact"
			} else {
				userSettings.menuState = ""
			}
			cachedUserSettings[email] = userSettings
			userSettingsSync.Unlock()
		}
	case "section":
		{
			dataSeparated := strings.Split(data.State, "-")
			var sectionState sectionState
			sectionState.section = dataSeparated[1]
			sectionState.state = dataSeparated[0]
			userSettingsSync.Lock()
			userSettings := cachedUserSettings[email]
			for index, state := range userSettings.sectionStates {
				if state.section == sectionState.section {
					userSettings.sectionStates = append(userSettings.sectionStates[0:index], userSettings.sectionStates[index+1:]...)
				}
			}
			userSettings.sectionStates = append(userSettings.sectionStates, sectionState)
			cachedUserSettings[email] = userSettings
			userSettingsSync.Unlock()
		}
	case "compact":
		{
			userSettingsSync.Lock()
			userSettings := cachedUserSettings[email]
			userSettings.compacted = data.State
			cachedUserSettings[email] = userSettings
			userSettingsSync.Unlock()
		}
	}

	logInfo("UPDATE", "User settings updated in "+time.Since(timer).String())
}

func updateUserDataSettings(email string, dataSelection string, settingsSelection string, workplaces []string) {
	logInfo("UPDATE", "Updating user settings for "+cachedUsersByEmail[email].FirstName+" "+cachedUsersByEmail[email].SecondName)
	timer := time.Now()
	userSettingsSync.Lock()
	settings := cachedUserSettings[email]
	if len(dataSelection) > 0 {
		settings.dataSelection = dataSelection
	}
	if len(settingsSelection) > 0 {
		settings.settingsSelection = settingsSelection
	}
	if len(workplaces) > 0 {
		settings.selectedWorkplaces = workplaces
	} else {
		settings.selectedWorkplaces = nil
	}
	cachedUserSettings[email] = settings
	userSettingsSync.Unlock()
	logInfo("UPDATE", "User settings updated in "+time.Since(timer).String())
}
