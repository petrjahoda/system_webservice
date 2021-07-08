package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type StatisticsPageData struct {
	Version               string
	Information           string
	Software              string
	Company               string
	Alarms                string
	MenuOverview          string
	MenuWorkplaces        string
	MenuCharts            string
	MenuStatistics        string
	MenuData              string
	MenuSettings          string
	SelectionMenu         []TableSelection
	Workplaces            []TableWorkplaceSelection
	Types                 []TableTypeSelection
	Users                 []TableUserSelection
	DataFilterPlaceholder string
	DateFrom              string
	DateTo                string
	UserEmail             string
	UserName              string
}

type StatisticsOutput struct {
	Result    string
	Compacted string
}

type StatisticsTypeInput struct {
	Email     string
	Selection string
}

type StatisticsTypeOutput struct {
	Result string
	Data   string
}

type StatisticsDataOutput struct {
	Result                 string
	Compacted              string
	SelectionChartData     []string
	SelectionChartValue    []float64
	SelectionChartText     []string
	WorkplaceChartData     []string
	WorkplaceChartValue    []float64
	WorkplaceChartText     []string
	UsersChartData         []string
	UsersChartValue        []float64
	UsersChartText         []string
	TimeChartData          []string
	TimeChartValue         []float64
	TimeChartText          []string
	DaysChartData          []string
	DaysChartValue         []float64
	DaysChartText          []string
	Locale                 string
	CalendarChartLocale    string
	FirstUpperChartLocale  string
	SecondUpperChartLocale string
	ThirdUpperChartLocale  string
	FourthUpperChartLocale string
}

func statistics(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	timer := time.Now()
	go updatePageCount("statistics")
	email, _, _ := request.BasicAuth()
	go updateWebUserRecord("statistics", email)
	logInfo("STATISTICS", "Sending page to "+email)
	var data StatisticsPageData
	data.Version = version
	companyNameSync.RLock()
	data.Company = cachedCompanyName
	companyNameSync.RUnlock()
	data.MenuOverview = getLocale(email, "menu-overview")
	data.MenuWorkplaces = getLocale(email, "menu-workplaces")
	data.MenuCharts = getLocale(email, "menu-charts")
	data.MenuStatistics = getLocale(email, "menu-statistics")
	data.MenuData = getLocale(email, "menu-data")
	data.MenuSettings = getLocale(email, "menu-settings")
	data.UserEmail = email
	data.UserName = cachedUsersByEmail[email].FirstName + " " + cachedUsersByEmail[email].SecondName
	usersByEmailSync.RLock()
	data.DateFrom = cachedUserWebSettings[email]["statistics-selected-from"]
	data.DateTo = cachedUserWebSettings[email]["statistics-selected-to"]
	selectedType := cachedUserWebSettings[email]["statistics-selected-type"]
	usersByEmailSync.RUnlock()
	//data.SelectionMenu = append(data.SelectionMenu, TableSelection{
	//	SelectionName:  getLocale(email, "overview"),
	//	SelectionValue: "overview",
	//	Selection:      getSelected(selectedType, "overview"),
	//})
	//data.SelectionMenu = append(data.SelectionMenu, TableSelection{
	//	SelectionName:  getLocale(email, "alarms"),
	//	SelectionValue: "alarms",
	//	Selection:      getSelected(selectedType, "alarms"),
	//})
	//data.SelectionMenu = append(data.SelectionMenu, TableSelection{
	//	SelectionName:  getLocale(email, "breakdowns"),
	//	SelectionValue: "breakdowns",
	//	Selection:      getSelected(selectedType, "breakdowns"),
	//})
	data.SelectionMenu = append(data.SelectionMenu, TableSelection{
		SelectionName:  getLocale(email, "downtimes"),
		SelectionValue: "downtimes",
		Selection:      getSelected(selectedType, "downtimes"),
	})
	//data.SelectionMenu = append(data.SelectionMenu, TableSelection{
	//	SelectionName:  getLocale(email, "faults"),
	//	SelectionValue: "faults",
	//	Selection:      getSelected(selectedType, "faults"),
	//})
	//data.SelectionMenu = append(data.SelectionMenu, TableSelection{
	//	SelectionName:  getLocale(email, "orders"),
	//	SelectionValue: "orders",
	//	Selection:      getSelected(selectedType, "orders"),
	//})
	//data.SelectionMenu = append(data.SelectionMenu, TableSelection{
	//	SelectionName:  getLocale(email, "packages"),
	//	SelectionValue: "packages",
	//	Selection:      getSelected(selectedType, "packages"),
	//})
	//data.SelectionMenu = append(data.SelectionMenu, TableSelection{
	//	SelectionName:  getLocale(email, "parts"),
	//	SelectionValue: "parts",
	//	Selection:      getSelected(selectedType, "parts"),
	//})
	usersByEmailSync.RLock()
	userRoleId := cachedUsersByEmail[email].UserRoleID
	usersByEmailSync.RUnlock()

	if userRoleId == 1 {
		logInfo("STATISTICS", "Adding data menu for administrator")
		//data.SelectionMenu = append(data.SelectionMenu, TableSelection{
		//	SelectionName:  getLocale(email, "users"),
		//	SelectionValue: "web-users",
		//	Selection:      getSelected(selectedType, "web-users"),
		//})
		//data.SelectionMenu = append(data.SelectionMenu, TableSelection{
		//	SelectionName:  getLocale(email, "system-statistics"),
		//	SelectionValue: "system-statistics",
		//	Selection:      getSelected(selectedType, "system-statistics"),
		//})
	}
	data.DataFilterPlaceholder = getLocale(email, "data-table-search-title")

	var dataWorkplaces []TableWorkplaceSelection
	workplacesByIdSync.RLock()
	workplacesById := cachedWorkplacesById
	workplacesByIdSync.RUnlock()
	for _, workplace := range workplacesById {
		userWebSettingsSync.RLock()
		selectedWorkplaces := cachedUserWebSettings[email]["statistics-selected-workplaces"]
		userWebSettingsSync.RUnlock()
		dataWorkplaces = append(dataWorkplaces, TableWorkplaceSelection{
			WorkplaceName:      workplace.Name,
			WorkplaceSelection: getWorkplaceWebSelection(selectedWorkplaces, workplace.Name),
		})
	}
	sort.Slice(dataWorkplaces, func(i, j int) bool {
		return dataWorkplaces[i].WorkplaceName < dataWorkplaces[j].WorkplaceName
	})
	data.Workplaces = dataWorkplaces

	var dataUsers []TableUserSelection
	usersByIdSync.RLock()
	usersById := cachedUsersById
	usersByIdSync.RUnlock()
	for _, user := range usersById {
		userWebSettingsSync.RLock()
		selectedUsers := cachedUserWebSettings[email]["statistics-selected-users"]
		userWebSettingsSync.RUnlock()
		dataUsers = append(dataUsers, TableUserSelection{
			UserName:      user.SecondName + " " + user.FirstName + " [" + user.Email + "]",
			UserSelection: getWorkplaceWebSelection(selectedUsers, user.SecondName+" "+user.FirstName+" ["+user.Email+"]"),
		})
	}
	sort.Slice(dataUsers, func(i, j int) bool {
		return dataUsers[i].UserName < dataUsers[j].UserName
	})
	data.Users = dataUsers
	for _, menu := range data.SelectionMenu {
		if menu.Selection == "selected" {
			switch menu.SelectionValue {
			case "downtimes":
				var datasTypes []TableTypeSelection
				downtimesByIdSync.RLock()
				downtimesById := cachedDowntimesById
				downtimesByIdSync.RUnlock()
				for _, downtime := range downtimesById {
					userWebSettingsSync.RLock()
					selectedTypes := cachedUserWebSettings[email]["statistics-selected-types-"+menu.SelectionValue]
					userWebSettingsSync.RUnlock()
					datasTypes = append(datasTypes, TableTypeSelection{
						TypeName:      downtime.Name,
						TypeSelection: getWorkplaceWebSelection(selectedTypes, downtime.Name),
					})
				}
				sort.Slice(datasTypes, func(i, j int) bool {
					return datasTypes[i].TypeName < datasTypes[j].TypeName
				})
				data.Types = datasTypes
			}
			break
		}
	}

	softwareNameSync.RLock()
	data.Software = cachedSoftwareName
	softwareNameSync.RUnlock()
	data.Information = "INF: Page processed in " + time.Since(timer).String()
	tmpl := template.Must(template.ParseFiles("./html/statistics.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("STATISTICS", "Page sent in "+time.Since(timer).String())
}

func loadStatisticsData(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	timer := time.Now()
	go updatePageCount("statistics")
	email, _, _ := request.BasicAuth()
	var data DataPageInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("DATA", "Error parsing data: "+err.Error())
		var responseData StatisticsOutput
		responseData.Result = "ERR: Error parsing data, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("DATA", "Loading statistics ended")
		return
	}
	logInfo("DATA", "Loading statistics for "+data.Data+" for "+strconv.Itoa(len(data.Workplaces))+" workplaces")
	if len(email) == 0 {
		email = data.Email
	}
	logInfo("DATA", "Loading statistics for "+email)
	locationSync.RLock()
	loc, err := time.LoadLocation(cachedLocation)
	locationSync.RUnlock()
	if err != nil {
		logError("DATA", "Problem loading timezone, setting Europe/Prague")
		loc, _ = time.LoadLocation("Europe/Prague")
	}
	layout := "2006-01-02T15:04"
	dateFrom, err := time.ParseInLocation(layout, data.From, loc)
	dateFrom = dateFrom.In(time.UTC)
	if err != nil {
		logError("DATA", "Problem parsing date: "+data.From)
		var responseData StatisticsOutput
		responseData.Result = "ERR: Error parsing data, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("DATA", "Loading statistics ended")
		return
	}
	dateTo, err := time.ParseInLocation(layout, data.To, loc)
	dateTo = dateTo.In(time.UTC)
	if err != nil {
		logError("DATA", "Problem parsing date: "+data.To)
		var responseData StatisticsOutput
		responseData.Result = "ERR: Error parsing data, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("DATA", "Loading statistics ended")
		return
	}
	logInfo("DATA", "From "+dateFrom.String()+" to "+dateTo.String())
	updateUserWebSettings(email, "statistics-selected-type", data.Data)
	updateUserWebSettings(email, "statistics-selected-from", data.From)
	updateUserWebSettings(email, "statistics-selected-to", data.To)
	logInfo("DATA", "Preprocessing takes "+time.Since(timer).String())
	switch data.Data {
	case "overview":
		//loadOverView(writer, workplaceIDs, dateFrom, dateTo, email)
	case "alarms":
		//loadAlarmsStatistics(writer, workplaceIds, dateFrom, dateTo, email)
	case "breakdowns":
		//loadBreakdownStatistics(writer, workplaceIds, dateFrom, dateTo, email)
	case "downtimes":
		loadDowntimesStatistics(writer, dateFrom, dateTo, email)
	case "faults":
		//loadFaultsStatistics(writer, workplaceIds, dateFrom, dateTo, email)
	case "orders":
		//loadOrdersStatistics(writer, workplaceIds, dateFrom, dateTo, email)
	case "packages":
		//loadPackagesStatistics(writer, workplaceIds, dateFrom, dateTo, email)
	case "parts":
		//loadPartsStatistics(writer, workplaceIds, dateFrom, dateTo, email)
	case "users":
		//loadUsersStatistics(writer, workplaceIds, dateFrom, dateTo, email)
	case "system-statistics":
		//loadSystemStatistics(writer, dateFrom, dateTo, email)
	}
	logInfo("DATA", "Statistics loaded in "+time.Since(timer).String())
	return
}

func loadTypesForSelectedStatistics(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	timer := time.Now()
	logInfo("STATISTICS", "Loading statistics types")
	email, _, _ := request.BasicAuth()
	var data StatisticsTypeInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("STATISTICS", "Error parsing data: "+err.Error())
		var responseData StatisticsTypeOutput
		responseData.Result = "ERR: Error parsing data, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("STATISTICS", "Loading statistics ended")
		return
	}
	logInfo("STATISTICS", "Loading statistics types for "+data.Selection)
	if len(email) == 0 {
		email = data.Email
	}
	updateUserWebSettings(email, "statistics-selected-type", data.Selection)
	var dataTypes []TableTypeSelection
	switch data.Selection {
	case "alarms":
		{
			alarmByIdSync.RLock()
			alarmsById := cachedAlarmsById
			alarmByIdSync.RUnlock()
			for _, selection := range alarmsById {
				userWebSettingsSync.RLock()
				selectedTypes := cachedUserWebSettings[email]["statistics-selected-types-"+data.Selection]
				userWebSettingsSync.RUnlock()
				dataTypes = append(dataTypes, TableTypeSelection{
					TypeName:      selection.Name,
					TypeSelection: getWorkplaceWebSelection(selectedTypes, selection.Name),
				})
			}
			sort.Slice(dataTypes, func(i, j int) bool {
				return dataTypes[i].TypeName < dataTypes[j].TypeName
			})
		}
	case "breakdowns":
		{
			breakdownByIdSync.RLock()
			breakdownsById := cachedBreakdownsById
			breakdownByIdSync.RUnlock()
			for _, selection := range breakdownsById {
				userWebSettingsSync.RLock()
				selectedTypes := cachedUserWebSettings[email]["statistics-selected-types-"+data.Selection]
				userWebSettingsSync.RUnlock()
				dataTypes = append(dataTypes, TableTypeSelection{
					TypeName:      selection.Name,
					TypeSelection: getWorkplaceWebSelection(selectedTypes, selection.Name),
				})
			}
			sort.Slice(dataTypes, func(i, j int) bool {
				return dataTypes[i].TypeName < dataTypes[j].TypeName
			})
		}
	case "downtimes":
		{
			downtimesByIdSync.RLock()
			downtimesById := cachedDowntimesById
			downtimesByIdSync.RUnlock()
			for _, selection := range downtimesById {
				userWebSettingsSync.RLock()
				selectedTypes := cachedUserWebSettings[email]["statistics-selected-types-"+data.Selection]
				userWebSettingsSync.RUnlock()
				dataTypes = append(dataTypes, TableTypeSelection{
					TypeName:      selection.Name,
					TypeSelection: getWorkplaceWebSelection(selectedTypes, selection.Name),
				})
			}
			sort.Slice(dataTypes, func(i, j int) bool {
				return dataTypes[i].TypeName < dataTypes[j].TypeName
			})
		}
	case "faults":
		{
			faultsByIdSync.RLock()
			faultsById := cachedFaultsById
			faultsByIdSync.RUnlock()
			for _, selection := range faultsById {
				userWebSettingsSync.RLock()
				selectedTypes := cachedUserWebSettings[email]["statistics-selected-types-"+data.Selection]
				userWebSettingsSync.RUnlock()
				dataTypes = append(dataTypes, TableTypeSelection{
					TypeName:      selection.Name,
					TypeSelection: getWorkplaceWebSelection(selectedTypes, selection.Name),
				})
			}
			sort.Slice(dataTypes, func(i, j int) bool {
				return dataTypes[i].TypeName < dataTypes[j].TypeName
			})
		}
	case "orders":
		{
			ordersByIdSync.RLock()
			ordersById := cachedOrdersById
			ordersByIdSync.RUnlock()
			for _, selection := range ordersById {
				userWebSettingsSync.RLock()
				selectedTypes := cachedUserWebSettings[email]["statistics-selected-types-"+data.Selection]
				userWebSettingsSync.RUnlock()
				dataTypes = append(dataTypes, TableTypeSelection{
					TypeName:      selection.Name,
					TypeSelection: getWorkplaceWebSelection(selectedTypes, selection.Name),
				})
			}
			sort.Slice(dataTypes, func(i, j int) bool {
				return dataTypes[i].TypeName < dataTypes[j].TypeName
			})
		}
	case "packages":
		{
			packagesByIdSync.RLock()
			packagesById := cachedPackagesById
			packagesByIdSync.RUnlock()
			for _, selection := range packagesById {
				userWebSettingsSync.RLock()
				selectedTypes := cachedUserWebSettings[email]["statistics-selected-types-"+data.Selection]
				userWebSettingsSync.RUnlock()
				dataTypes = append(dataTypes, TableTypeSelection{
					TypeName:      selection.Name,
					TypeSelection: getWorkplaceWebSelection(selectedTypes, selection.Name),
				})
			}
			sort.Slice(dataTypes, func(i, j int) bool {
				return dataTypes[i].TypeName < dataTypes[j].TypeName
			})
		}
	}
	var responseData StatisticsPageData
	responseData.Types = dataTypes
	tmpl, err := template.ParseFiles("./html/statistics-selection.html")
	if err != nil {
		logError("SETTINGS", "Problem parsing html file: "+err.Error())
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
	} else {
		_ = tmpl.Execute(writer, responseData)
		logInfo("SETTINGS", "Workplaces updated in "+time.Since(timer).String())
	}
	logInfo("STATISTICS", "Statistics types loaded in "+time.Since(timer).String())
}
