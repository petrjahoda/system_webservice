package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
)

type UserSettingsInput struct {
	Settings string
	State    string
}

func updateUserSettings(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	logInfo("MAIN", "Parsing data")
	var data UserSettingsInput
	_ = json.NewDecoder(request.Body).Decode(&data)
	email, _, _ := request.BasicAuth()
	switch data.Settings {
	case "menu":
		{
			cachedUserSettingsSync.Lock()
			userSettings := cachedUserSettings[email]
			if data.State == "false" {
				userSettings.menuState = "compacted js-compact"
			} else {
				userSettings.menuState = ""
			}
			cachedUserSettings[email] = userSettings
			cachedUserSettingsSync.Unlock()
		}
	case "section":
		{
			dataSeparated := strings.Split(data.State, "-")
			var sectionState sectionState
			sectionState.section = dataSeparated[1]
			sectionState.state = dataSeparated[0]
			cachedUserSettingsSync.Lock()
			userSettings := cachedUserSettings[email]
			for index, state := range userSettings.sectionStates {
				if state.section == sectionState.section {
					userSettings.sectionStates = append(userSettings.sectionStates[0:index], userSettings.sectionStates[index+1:]...)
				}
			}
			userSettings.sectionStates = append(userSettings.sectionStates, sectionState)
			cachedUserSettings[email] = userSettings
			cachedUserSettingsSync.Unlock()
		}
	}
}
