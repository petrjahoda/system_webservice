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
	Id   string
	Type string
}

type SettingsPageOutput struct {
	Version        string
	Information    string
	Software       string
	Company        string
	Alarms         string
	MenuOverview   string
	MenuWorkplaces string
	MenuCharts     string
	MenuStatistics string
	MenuData       string
	MenuSettings   string
	SelectionMenu  []TableSelection
	DateLocale     string
	UserEmail      string
	UserName       string
}

func settings(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	timer := time.Now()
	go updatePageCount("settings")
	email, _, _ := request.BasicAuth()
	go updateWebUserRecord("settings", email)
	logInfo("SETTINGS", "Sending page to "+email)
	var data SettingsPageOutput
	data.Version = version
	companyNameSync.RLock()
	data.Company = cachedCompanyName
	companyNameSync.RUnlock()
	data.UserEmail = email
	usersByEmailSync.RLock()
	data.UserName = cachedUsersByEmail[email].FirstName + " " + cachedUsersByEmail[email].SecondName
	usersByEmailSync.RUnlock()
	data.MenuOverview = getLocale(email, "menu-overview")
	data.MenuWorkplaces = getLocale(email, "menu-workplaces")
	data.MenuCharts = getLocale(email, "menu-charts")
	data.MenuStatistics = getLocale(email, "menu-statistics")
	data.MenuData = getLocale(email, "menu-data")
	data.MenuSettings = getLocale(email, "menu-settings")
	usersByEmailSync.RLock()
	userLocale := cachedUsersByEmail[email].Locale
	usersByEmailSync.RUnlock()
	localesSync.RLock()
	data.DateLocale = cachedLocales[userLocale]
	localesSync.RUnlock()
	userWebSettingsSync.RLock()
	selectedType := cachedUserWebSettings[email]["settings-selected-type"]
	userWebSettingsSync.RUnlock()
	data.SelectionMenu = append(data.SelectionMenu, TableSelection{
		SelectionName:  getLocale(email, "user-name"),
		SelectionValue: "user",
		Selection:      getSelected(selectedType, "user-name"),
	})
	usersByEmailSync.RLock()
	userRoleId := cachedUsersByEmail[email].UserRoleID
	usersByEmailSync.RUnlock()
	if userRoleId != user {
		logInfo("SETTINGS", "Adding data menu for power user")
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "alarms"),
			SelectionValue: "alarms",
			Selection:      getSelected(selectedType, "alarms"),
		})
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "breakdowns"),
			SelectionValue: "breakdowns",
			Selection:      getSelected(selectedType, "breakdowns"),
		})
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "downtimes"),
			SelectionValue: "downtimes",
			Selection:      getSelected(selectedType, "downtimes"),
		})
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "faults"),
			SelectionValue: "faults",
			Selection:      getSelected(selectedType, "faults"),
		})
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "operations"),
			SelectionValue: "operations",
			Selection:      getSelected(selectedType, "operations"),
		})
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "orders"),
			SelectionValue: "orders",
			Selection:      getSelected(selectedType, "orders"),
		})
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "packages"),
			SelectionValue: "packages",
			Selection:      getSelected(selectedType, "packages"),
		})
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "parts"),
			SelectionValue: "parts",
			Selection:      getSelected(selectedType, "parts"),
		})
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "products"),
			SelectionValue: "products",
			Selection:      getSelected(selectedType, "products"),
		})
	}
	if userRoleId == administrator {
		logInfo("SETTINGS", "Adding data menu for administrator")
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "devices"),
			SelectionValue: "devices",
			Selection:      getSelected(selectedType, "devices"),
		})
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "states"),
			SelectionValue: "states",
			Selection:      getSelected(selectedType, "states"),
		})
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "users"),
			SelectionValue: "users",
			Selection:      getSelected(selectedType, "users"),
		})

		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "workplaces"),
			SelectionValue: "workplaces",
			Selection:      getSelected(selectedType, "workplaces"),
		})
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "workshifts"),
			SelectionValue: "workshifts",
			Selection:      getSelected(selectedType, "workshifts"),
		})
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "system-settings"),
			SelectionValue: "system-settings",
			Selection:      getSelected(selectedType, "system-settings"),
		})
	}
	softwareNameSync.RLock()
	data.Software = cachedSoftwareName
	softwareNameSync.RUnlock()
	data.Information = "INF: Page processed in " + time.Since(timer).String()
	tmpl := template.Must(template.ParseFiles("./html/settings.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("SETTINGS", "Page sent in "+time.Since(timer).String())
}

func loadSettingsData(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	timer := time.Now()
	go updatePageCount("settings")
	email, _, _ := request.BasicAuth()
	logInfo("SETTINGS", "Loading settings for "+email)
	var data SettingsPageInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "ERR: Error parsing data, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Loading settings ended with error")
		return
	}
	logInfo("SETTINGS", "Loading settings for "+data.Data)
	updateUserWebSettings(email, "settings-selected-type", data.Data)
	switch data.Data {
	case "alarms":
		loadAlarms(writer, email)
	case "breakdowns":
		loadBreakdowns(writer, email)
	case "downtimes":
		loadDowntimes(writer, email)
	case "faults":
		loadFaults(writer, email)
	case "operations":
		loadOperations(writer, email)
	case "orders":
		loadOrders(writer, email)
	case "packages":
		loadPackages(writer, email)
	case "parts":
		loadParts(writer, email)
	case "products":
		loadProducts(writer, email)
	case "states":
		loadStates(writer, email)
	case "devices":
		loadDevices(writer, email)
	case "users":
		loadUsers(writer, email)
	case "workplaces":
		loadWorkplaces(writer, email)
	case "workshifts":
		loadWorkShifts(writer, email)
	case "system-settings":
		loadSystemSettings(writer, email)
	case "user":
		loadUserSettings(writer, email)
	}
	logInfo("SETTINGS", "Settings loaded in "+time.Since(timer).String())
	return
}

func loadSettingsDetail(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	timer := time.Now()
	email, _, _ := request.BasicAuth()
	logInfo("SETTINGS", "Loading settings detail for "+email)
	var data SettingsPageInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "ERR: Error parsing data, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Loading settings detail ended with error")
		return
	}
	logInfo("SETTINGS", "Loading details settings for "+data.Data+", "+data.Id)
	switch data.Type {
	case "first":
		{
			switch data.Data {
			case "alarms":
				loadAlarm(data.Id, writer, email)
			case "breakdowns":
				loadBreakdown(data.Id, writer, email)
			case "downtimes":
				loadDowntime(data.Id, writer, email)
			case "faults":
				loadFault(data.Id, writer, email)
			case "operations":
				loadOperation(data.Id, writer, email)
			case "orders":
				loadOrder(data.Id, writer, email)
			case "packages":
				loadPackage(data.Id, writer, email)
			case "parts":
				loadPart(data.Id, writer, email)
			case "products":
				loadProduct(data.Id, writer, email)
			case "states":
				loadState(data.Id, writer, email)
			case "devices":
				loadDevice(data.Id, writer, email)
			case "system-settings":
				loadSystemSettingsDetails(data.Id, writer, email)
			case "users":
				loadUser(data.Id, writer, email)
			case "workplaces":
				loadWorkplace(data.Id, writer, email)
			case "workshifts":
				loadWorkshift(data.Id, writer, email)
			}
		}
	case "second":
		{
			switch data.Data {
			case "breakdowns":
				loadBreakdownTypes(data.Id, writer, email)
			case "downtimes":
				loadDowntimeType(data.Id, writer, email)
			case "faults":
				loadFaultType(data.Id, writer, email)
			case "packages":
				loadPackageType(data.Id, writer, email)
			case "users":
				loadUserType(data.Id, writer, email)
			case "workplaces":
				loadWorkplaceSection(data.Id, writer, email)
			}

		}
	case "third":
		{
			switch data.Data {
			case "workplaces":
				loadWorkplaceMode(data.Id, writer, email)
			}
		}

	}
	logInfo("SETTINGS", "Detail settings loaded in "+time.Since(timer).String())
	return
}
