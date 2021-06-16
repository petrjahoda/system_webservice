package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"strings"
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
	data.SelectionMenu = append(data.SelectionMenu, TableSelection{
		SelectionName:  getLocale(email, "user-name"),
		SelectionValue: "username",
		Selection:      getSelected(selectedType, "username"),
	})
	data.SelectionMenu = append(data.SelectionMenu, TableSelection{
		SelectionName:  getLocale(email, "company"),
		SelectionValue: "company",
		Selection:      getSelected(selectedType, "company"),
	})
	// TODO: ADD COMPANY TO LOCALES
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
		SelectionName:  getLocale(email, "states"),
		SelectionValue: "states",
		Selection:      getSelected(selectedType, "states"),
	})
	usersByEmailSync.RLock()
	userRoleId := cachedUsersByEmail[email].UserRoleID
	usersByEmailSync.RUnlock()

	if userRoleId == 1 {
		logInfo("STATISTICS", "Adding data menu for administrator")
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "users"),
			SelectionValue: "users",
			Selection:      getSelected(selectedType, "users"),
		})
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "system-statistics"),
			SelectionValue: "system-statistics",
			Selection:      getSelected(selectedType, "system-statistics"),
		})
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
	logInfo("DATA", "Loading statistics for "+email)
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
	selectedWorkplaces := ""
	for _, workplace := range data.Workplaces {
		selectedWorkplaces += workplace + ";"
	}
	selectedWorkplaces = strings.TrimRight(selectedWorkplaces, ";")
	updateUserWebSettings(email, "statistics-selected-workplaces", selectedWorkplaces)
	//workplaceIds := getWorkplaceIds(data)
	logInfo("DATA", "Preprocessing takes "+time.Since(timer).String())
	switch data.Data {
	case "alarms":
		//loadAlarmsStatistics(writer, workplaceIds, dateFrom, dateTo, email)
	case "breakdowns":
		//loadBreakdownStatistics(writer, workplaceIds, dateFrom, dateTo, email)
	case "downtimes":
		//loadDowntimesStatistics(writer, workplaceIds, dateFrom, dateTo, email)
	case "faults":
		//loadFaultsStatistics(writer, workplaceIds, dateFrom, dateTo, email)
	case "orders":
		//loadOrdersStatistics(writer, workplaceIds, dateFrom, dateTo, email)
	case "packages":
		//loadPackagesStatistics(writer, workplaceIds, dateFrom, dateTo, email)
	case "parts":
		//loadPartsStatistics(writer, workplaceIds, dateFrom, dateTo, email)
	case "states":
		//loadStatesStatistics(writer, workplaceIds, dateFrom, dateTo, email)
	case "users":
		//loadUsersStatistics(writer, workplaceIds, dateFrom, dateTo, email)
	case "system-statistics":
		//loadSystemStatistics(writer, dateFrom, dateTo, email)
	}
	logInfo("DATA", "Statistics loaded in "+time.Since(timer).String())
	return
}
