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
	go updatePageCount("settings")
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
	data.Compacted = cachedUserWebSettings[email]["menu"]
	data.DateLocale = cachedLocales[cachedUsersByEmail[email].Locale]
	data.SelectionMenu = append(data.SelectionMenu, TableSelection{
		SelectionName:  getLocale(email, "user-name"),
		SelectionValue: "user",
		Selection:      getSelected(cachedUserWebSettings[email]["settings-selected-type"], "user-name"),
	})
	if cachedUsersByEmail[email].UserRoleID != 3 {
		logInfo("SETTINGS", "Adding data menu for power user")
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "alarms"),
			SelectionValue: "alarms",
			Selection:      getSelected(cachedUserWebSettings[email]["settings-selected-type"], "alarms"),
		})
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "breakdowns"),
			SelectionValue: "breakdowns",
			Selection:      getSelected(cachedUserWebSettings[email]["settings-selected-type"], "breakdowns"),
		})
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "downtimes"),
			SelectionValue: "downtimes",
			Selection:      getSelected(cachedUserWebSettings[email]["settings-selected-type"], "downtimes"),
		})
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "faults"),
			SelectionValue: "faults",
			Selection:      getSelected(cachedUserWebSettings[email]["settings-selected-type"], "faults"),
		})
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "operations"),
			SelectionValue: "operations",
			Selection:      getSelected(cachedUserWebSettings[email]["settings-selected-type"], "operations"),
		})
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "orders"),
			SelectionValue: "orders",
			Selection:      getSelected(cachedUserWebSettings[email]["settings-selected-type"], "orders"),
		})
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "packages"),
			SelectionValue: "packages",
			Selection:      getSelected(cachedUserWebSettings[email]["settings-selected-type"], "packages"),
		})
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "parts"),
			SelectionValue: "parts",
			Selection:      getSelected(cachedUserWebSettings[email]["settings-selected-type"], "parts"),
		})
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "products"),
			SelectionValue: "products",
			Selection:      getSelected(cachedUserWebSettings[email]["settings-selected-type"], "products"),
		})
	}
	if cachedUsersByEmail[email].UserRoleID == 1 {
		logInfo("SETTINGS", "Adding data menu for administrator")
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "devices"),
			SelectionValue: "devices",
			Selection:      getSelected(cachedUserWebSettings[email]["settings-selected-type"], "devices"),
		})
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "states"),
			SelectionValue: "states",
			Selection:      getSelected(cachedUserWebSettings[email]["settings-selected-type"], "states"),
		})
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "users"),
			SelectionValue: "users",
			Selection:      getSelected(cachedUserWebSettings[email]["settings-selected-type"], "users"),
		})

		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "workplaces"),
			SelectionValue: "workplaces",
			Selection:      getSelected(cachedUserWebSettings[email]["settings-selected-type"], "workplaces"),
		})
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "workshifts"),
			SelectionValue: "workshifts",
			Selection:      getSelected(cachedUserWebSettings[email]["settings-selected-type"], "workshifts"),
		})
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "system-settings"),
			SelectionValue: "system-settings",
			Selection:      getSelected(cachedUserWebSettings[email]["settings-selected-type"], "system-settings"),
		})
	}
	data.Information = "INF: Page processed in " + time.Since(timer).String()
	tmpl := template.Must(template.ParseFiles("./html/settings.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("SETTINGS", "Page sent in "+time.Since(timer).String())
}

func loadSettingsData(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	timer := time.Now()
	go updatePageCount("settings")
	email, _, _ := request.BasicAuth()
	logInfo("SETTINGS", "Loading settings for "+cachedUsersByEmail[email].FirstName+" "+cachedUsersByEmail[email].SecondName)
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
	logInfo("SETTINGS", "Loading settings detail for "+cachedUsersByEmail[email].FirstName+" "+cachedUsersByEmail[email].SecondName)
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
