package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
	"strings"
	"time"
)

type SettingsPageData struct {
	Version        string
	Company        string
	Alarms         string
	MenuOverview   string
	MenuWorkplaces string
	MenuCharts     string
	MenuStatistics string
	MenuData       string
	MenuSettings   string
	Compacted      string
	SelectionMenu  []TableSelection
	DateLocale     string
}

type SettingsDataPageInput struct {
	Data string
	Name string
}

func settings(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	ipAddress := strings.Split(request.RemoteAddr, ":")
	logInfo("SETTINGS", "Sending home page to "+ipAddress[0])
	email, _, _ := request.BasicAuth()
	var data SettingsPageData
	data.Version = version
	data.Company = cachedCompanyName
	data.MenuOverview = getLocale(email, "menu-overview")
	data.MenuWorkplaces = getLocale(email, "menu-workplaces")
	data.MenuCharts = getLocale(email, "menu-charts")
	data.MenuStatistics = getLocale(email, "menu-statistics")
	data.MenuData = getLocale(email, "menu-data")
	data.MenuSettings = getLocale(email, "menu-settings")
	data.Compacted = cachedUserSettings[email].menuState
	data.DateLocale = cachedLocales[cachedUsersByEmail[email].Locale]
	data.SelectionMenu = append(data.SelectionMenu, TableSelection{
		SelectionName:  getLocale(email, "alarms"),
		SelectionValue: "alarms",
		Selection:      getSelected(cachedUserSettings[email].dataSelection, "alarms"),
	})
	data.SelectionMenu = append(data.SelectionMenu, TableSelection{
		SelectionName:  getLocale(email, "breakdowns"),
		SelectionValue: "breakdowns",
		Selection:      getSelected(cachedUserSettings[email].dataSelection, "breakdowns"),
	})
	data.SelectionMenu = append(data.SelectionMenu, TableSelection{
		SelectionName:  getLocale(email, "downtimes"),
		SelectionValue: "downtimes",
		Selection:      getSelected(cachedUserSettings[email].dataSelection, "downtimes"),
	})
	data.SelectionMenu = append(data.SelectionMenu, TableSelection{
		SelectionName:  getLocale(email, "faults"),
		SelectionValue: "faults",
		Selection:      getSelected(cachedUserSettings[email].dataSelection, "faults"),
	})
	data.SelectionMenu = append(data.SelectionMenu, TableSelection{
		SelectionName:  getLocale(email, "operations"),
		SelectionValue: "operations",
		Selection:      getSelected(cachedUserSettings[email].dataSelection, "operations"),
	})
	data.SelectionMenu = append(data.SelectionMenu, TableSelection{
		SelectionName:  getLocale(email, "orders"),
		SelectionValue: "orders",
		Selection:      getSelected(cachedUserSettings[email].dataSelection, "orders"),
	})
	data.SelectionMenu = append(data.SelectionMenu, TableSelection{
		SelectionName:  getLocale(email, "packages"),
		SelectionValue: "packages",
		Selection:      getSelected(cachedUserSettings[email].dataSelection, "packages"),
	})
	data.SelectionMenu = append(data.SelectionMenu, TableSelection{
		SelectionName:  getLocale(email, "parts"),
		SelectionValue: "parts",
		Selection:      getSelected(cachedUserSettings[email].dataSelection, "parts"),
	})
	data.SelectionMenu = append(data.SelectionMenu, TableSelection{
		SelectionName:  getLocale(email, "products"),
		SelectionValue: "products",
		Selection:      getSelected(cachedUserSettings[email].dataSelection, "products"),
	})
	data.SelectionMenu = append(data.SelectionMenu, TableSelection{
		SelectionName:  getLocale(email, "states"),
		SelectionValue: "states",
		Selection:      getSelected(cachedUserSettings[email].dataSelection, "states"),
	})
	if cachedUsersByEmail[email].UserTypeID == 2 {
		logInfo("SETTINGS", "Adding data menu for administrator")
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "devices"),
			SelectionValue: "devices",
			Selection:      getSelected(cachedUserSettings[email].dataSelection, "devices"),
		})
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "system-settings"),
			SelectionValue: "system-settings",
			Selection:      getSelected(cachedUserSettings[email].dataSelection, "system-settings"),
		})
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "users"),
			SelectionValue: "users",
			Selection:      getSelected(cachedUserSettings[email].dataSelection, "users"),
		})
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "workplaces"),
			SelectionValue: "workplaces",
			Selection:      getSelected(cachedUserSettings[email].dataSelection, "workplaces"),
		})
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "workshifts"),
			SelectionValue: "workshifts",
			Selection:      getSelected(cachedUserSettings[email].dataSelection, "workshifts"),
		})
	}
	tmpl := template.Must(template.ParseFiles("./html/settings.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("SETTINGS", "Home page sent")
}

func getSettingsData(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	timer := time.Now()
	email, _, _ := request.BasicAuth()
	logInfo("SETTINGS", "Sending settings data to "+cachedUsersByEmail[email].FirstName+" "+cachedUsersByEmail[email].SecondName)
	var data SettingsDataPageInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS", "Error parsing data: "+err.Error())
		var responseData DataPageOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Processing data ended")
		return
	}
	logInfo("SETTINGS", "Loading settings for "+data.Data)
	updateUserDataSettings(email, data.Data, nil)
	switch data.Data {
	case "alarms":
		processAlarmsSettings(writer, email)
	case "breakdowns":
		//processBreakdownsSettings(writer, email)
	case "downtimes":
		//processDowntimesSettings(writer, email)
	case "faults":
		//processFaultsSettings(writer, email)
	case "operations":
		//processOperationsSettings(writer, email)
	case "orders":
		//processOrdersSettings(writer, email)
	case "packages":
		//processPackagesSettings(writer, email)
	case "parts":
		//processPartsSettings(writer, email)
	case "products":
		//processProductsSettings(writer, email)
	case "states":
		//processStatesSettings(writer, email)
	case "devices":
		//processDevicesSettings(writer, email)
	case "system-settings":
		//processSystemSettings(writer, email)
	case "users":
		//processUsersSettings(writer, email)
	case "workplace":
		//processWorkplacesSettings(writer, email)
	case "workshifts":
		//processWorkshiftsSettings(writer, email)
	}
	logInfo("SETTINGS", "Settings data sent in "+time.Since(timer).String())
	return
}

func getDetailSettings(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	timer := time.Now()
	email, _, _ := request.BasicAuth()
	logInfo("SETTINGS", "Sending settings data to "+cachedUsersByEmail[email].FirstName+" "+cachedUsersByEmail[email].SecondName)
	var data SettingsDataPageInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS", "Error parsing data: "+err.Error())
		var responseData DataPageOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Processing data ended")
		return
	}
	logInfo("SETTINGS", "Loading details settings for "+data.Data+", "+data.Name)
	switch data.Data {
	case "alarms":
		processDetailAlarmSettings(data.Name, writer, email)
	case "breakdowns":
		//processBreakdownsSettings(writer, email)
	case "downtimes":
		//processDowntimesSettings(writer, email)
	case "faults":
		//processFaultsSettings(writer, email)
	case "operations":
		//processOperationsSettings(writer, email)
	case "orders":
		//processOrdersSettings(writer, email)
	case "packages":
		//processPackagesSettings(writer, email)
	case "parts":
		//processPartsSettings(writer, email)
	case "products":
		//processProductsSettings(writer, email)
	case "states":
		//processStatesSettings(writer, email)
	case "devices":
		//processDevicesSettings(writer, email)
	case "system-settings":
		//processSystemSettings(writer, email)
	case "users":
		//processUsersSettings(writer, email)
	case "workplace":
		//processWorkplacesSettings(writer, email)
	case "workshifts":
		//processWorkshiftsSettings(writer, email)
	}
	logInfo("SETTINGS", "Detail settings sent in "+time.Since(timer).String())
	return
}
