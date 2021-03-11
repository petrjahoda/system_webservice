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

func updateUserSettings(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	email, _, _ := request.BasicAuth()
	logInfo("UPDATE", "Updating user settings for "+cachedUsersByEmail[email].FirstName+" "+cachedUsersByEmail[email].SecondName)
	timer := time.Now()
	var data UserSettingsInput
	_ = json.NewDecoder(request.Body).Decode(&data)
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
	}
	cachedUserSettings[email] = settings
	userSettingsSync.Unlock()
	logInfo("UPDATE", "User settings updated in "+time.Since(timer).String())
}
