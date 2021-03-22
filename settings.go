package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
	"time"
)

type SettingsPageInput struct {
	Data string
	Name string
	Type string
}

type SettingsPageOutput struct {
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
	UserEmail      string
	UserName       string
}

func settings(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	timer := time.Now()
	email, _, _ := request.BasicAuth()
	logInfo("SETTINGS", "Sending page to "+cachedUsersByEmail[email].FirstName+" "+cachedUsersByEmail[email].SecondName)
	var data SettingsPageOutput
	data.Version = version
	data.Company = cachedCompanyName
	data.UserEmail = email
	data.UserName = cachedUsersByEmail[email].FirstName + " " + cachedUsersByEmail[email].SecondName
	data.MenuOverview = getLocale(email, "menu-overview")
	data.MenuWorkplaces = getLocale(email, "menu-workplaces")
	data.MenuCharts = getLocale(email, "menu-charts")
	data.MenuStatistics = getLocale(email, "menu-statistics")
	data.MenuData = getLocale(email, "menu-data")
	data.MenuSettings = getLocale(email, "menu-settings")
	data.Compacted = cachedUserSettings[email].menuState
	data.DateLocale = cachedLocales[cachedUsersByEmail[email].Locale]
	data.SelectionMenu = append(data.SelectionMenu, TableSelection{
		SelectionName:  getLocale(email, "user-name"),
		SelectionValue: "user",
		Selection:      getSelected(cachedUserSettings[email].settingsSelection, "user-name"),
	})
	if cachedUsersByEmail[email].UserRoleID != 3 {
		logInfo("SETTINGS", "Adding data menu for power user")
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "alarms"),
			SelectionValue: "alarms",
			Selection:      getSelected(cachedUserSettings[email].settingsSelection, "alarms"),
		})
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "breakdowns"),
			SelectionValue: "breakdowns",
			Selection:      getSelected(cachedUserSettings[email].settingsSelection, "breakdowns"),
		})
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "downtimes"),
			SelectionValue: "downtimes",
			Selection:      getSelected(cachedUserSettings[email].settingsSelection, "downtimes"),
		})
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "faults"),
			SelectionValue: "faults",
			Selection:      getSelected(cachedUserSettings[email].settingsSelection, "faults"),
		})
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "operations"),
			SelectionValue: "operations",
			Selection:      getSelected(cachedUserSettings[email].settingsSelection, "operations"),
		})
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "orders"),
			SelectionValue: "orders",
			Selection:      getSelected(cachedUserSettings[email].settingsSelection, "orders"),
		})
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "packages"),
			SelectionValue: "packages",
			Selection:      getSelected(cachedUserSettings[email].settingsSelection, "packages"),
		})
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "parts"),
			SelectionValue: "parts",
			Selection:      getSelected(cachedUserSettings[email].settingsSelection, "parts"),
		})
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "products"),
			SelectionValue: "products",
			Selection:      getSelected(cachedUserSettings[email].settingsSelection, "products"),
		})
	}
	if cachedUsersByEmail[email].UserRoleID == 1 {
		logInfo("SETTINGS", "Adding data menu for administrator")
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "devices"),
			SelectionValue: "devices",
			Selection:      getSelected(cachedUserSettings[email].settingsSelection, "devices"),
		})
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "states"),
			SelectionValue: "states",
			Selection:      getSelected(cachedUserSettings[email].settingsSelection, "states"),
		})
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "users"),
			SelectionValue: "users",
			Selection:      getSelected(cachedUserSettings[email].settingsSelection, "users"),
		})

		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "workplaces"),
			SelectionValue: "workplaces",
			Selection:      getSelected(cachedUserSettings[email].settingsSelection, "workplaces"),
		})
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "workshifts"),
			SelectionValue: "workshifts",
			Selection:      getSelected(cachedUserSettings[email].settingsSelection, "workshifts"),
		})
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "system-settings"),
			SelectionValue: "system-settings",
			Selection:      getSelected(cachedUserSettings[email].settingsSelection, "system-settings"),
		})
	}

	tmpl := template.Must(template.ParseFiles("./html/settings.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("SETTINGS", "Page sent in "+time.Since(timer).String())
}

func loadSettingsData(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	timer := time.Now()
	email, _, _ := request.BasicAuth()
	logInfo("SETTINGS", "Loading settings for "+cachedUsersByEmail[email].FirstName+" "+cachedUsersByEmail[email].SecondName)
	var data SettingsPageInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Loading settings ended")
		return
	}
	logInfo("SETTINGS", "Loading settings for "+data.Data)
	updateUserDataSettings(email, "", data.Data, nil)
	switch data.Data {
	case "alarms":
		loadAlarmsSettings(writer, email)
	case "breakdowns":
		loadBreakdownsSettings(writer, email)
	case "downtimes":
		loadDowntimesSettings(writer, email)
	case "faults":
		loadFaultsSettings(writer, email)
	case "operations":
		loadOperationsSettings(writer, email)
	case "orders":
		loadOrdersSettings(writer, email)
	case "packages":
		loadPackagesSettings(writer, email)
	case "parts":
		loadPartsSettings(writer, email)
	case "products":
		loadProductsSettings(writer, email)
	case "states":
		loadStatesSettings(writer, email)
	case "devices":
		loadDevicesSettings(writer, email)
	case "users":
		loadUsersSettings(writer, email)
	case "workplaces":
		loadWorkplacesSettings(writer, email)
	case "workshifts":
		loadWorkshiftsSettings(writer, email)
	case "system-settings":
		loadSystemSettings(writer, email)
	case "user":
		loadUserSettings(writer, email)
	}
	logInfo("SETTINGS", "Settings loaded in "+time.Since(timer).String())
	return
}

func loadSettingsDetail(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	timer := time.Now()
	email, _, _ := request.BasicAuth()
	logInfo("SETTINGS", "Loading settings detail for "+cachedUsersByEmail[email].FirstName+" "+cachedUsersByEmail[email].SecondName)
	var data SettingsPageInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Loading settings detail ended")
		return
	}

	logInfo("SETTINGS", "Loading details settings for "+data.Data+", "+data.Name)
	switch data.Type {
	case "first":
		{
			switch data.Data {
			case "alarms":
				loadAlarmDetails(data.Name, writer, email)
			case "breakdowns":
				loadBreakdownDetails(data.Name, writer, email)
			case "downtimes":
				loadDowntimeDetails(data.Name, writer, email)
			case "faults":
				loadFaultDetails(data.Name, writer, email)
			case "operations":
				loadOperationDetails(data.Name, writer, email)
			case "orders":
				loadOrderDetails(data.Name, writer, email)
			case "packages":
				loadPackageDetails(data.Name, writer, email)
			case "parts":
				loadPartDetails(data.Name, writer, email)
			case "products":
				loadProductDetails(data.Name, writer, email)
			case "states":
				loadStateDetails(data.Name, writer, email)
			case "devices":
				loadDeviceDetails(data.Name, writer, email)
			case "system-settings":
				loadSystemSettingsDetails(data.Name, writer, email)
			case "users":
				loadUserDetails(data.Name, writer, email)
			case "workplaces":
				loadWorkplaceDetails(data.Name, writer, email)
			case "workshifts":
				loadWorkshiftDetails(data.Name, writer, email)
			}
		}
	case "second":
		{
			switch data.Data {
			case "breakdowns":
				loadBreakdownTypeDetails(data.Name, writer, email)
			case "downtimes":
				loadDowntimeTypeDetails(data.Name, writer, email)
			case "faults":
				loadFaultTypeDetails(data.Name, writer, email)
			case "packages":
				loadPackageTypeDetails(data.Name, writer, email)
			case "users":
				loadUserTypeDetails(data.Name, writer, email)
			case "workplaces":
				loadWorkplaceSectionDetails(data.Name, writer, email)
			}

		}
	case "third":
		{
			switch data.Data {
			case "workplaces":
				loadWorkplaceModeDetails(data.Name, writer, email)
			}
		}

	}

	logInfo("SETTINGS", "Detail settings loaded in "+time.Since(timer).String())
	return
}
