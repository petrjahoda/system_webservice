package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/petrjahoda/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type DataPageInput struct {
	Data       string
	Workplaces []string
	From       string
	To         string
}
type DataPageOutput struct {
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
	DateLocale            string
	UserEmail             string
	UserName              string
	DateFrom              string
	DateTo                string
}

type TableSelection struct {
	SelectionName  string
	SelectionValue string
	Selection      string
}

type TableWorkplaceSelection struct {
	WorkplaceName      string
	WorkplaceSelection string
}

func data(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	timer := time.Now()
	go updatePageCount("data")
	email, _, _ := request.BasicAuth()
	go updateWebUserRecord("data", email)
	logInfo("DATA", "Sending page to "+email)
	var data DataPageOutput
	data.Version = version
	usersByEmailSync.RLock()
	userLocale := cachedUsersByEmail[email].Locale
	usersByEmailSync.RUnlock()
	localesSync.RLock()
	data.DateLocale = cachedLocales[userLocale]
	localesSync.RUnlock()
	data.UserEmail = email
	usersByEmailSync.RLock()
	data.UserName = cachedUsersByEmail[email].FirstName + " " + cachedUsersByEmail[email].SecondName
	usersByEmailSync.RUnlock()
	userWebSettingsSync.RLock()
	data.DateFrom = cachedUserWebSettings[email]["data-selected-from"]
	data.DateTo = cachedUserWebSettings[email]["data-selected-to"]
	userWebSettingsSync.RUnlock()
	companyNameSync.RLock()
	data.Company = cachedCompanyName
	companyNameSync.RUnlock()
	data.MenuOverview = getLocale(email, "menu-overview")
	data.MenuWorkplaces = getLocale(email, "menu-workplaces")
	data.MenuCharts = getLocale(email, "menu-charts")
	data.MenuStatistics = getLocale(email, "menu-statistics")
	data.MenuData = getLocale(email, "menu-data")
	data.MenuSettings = getLocale(email, "menu-settings")
	data.DataFilterPlaceholder = getLocale(email, "data-table-search-title")
	userWebSettingsSync.RLock()
	selectedType := cachedUserWebSettings[email]["data-selected-type"]
	userWebSettingsSync.RUnlock()
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
		logInfo("DATA", "Adding data menu for administrator")
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
	var dataWorkplaces []TableWorkplaceSelection
	workplacesByIdSync.RLock()
	workplacesById := cachedWorkplacesById
	workplacesByIdSync.RUnlock()
	for _, workplace := range workplacesById {
		userWebSettingsSync.RLock()
		selectedWorkplace := cachedUserWebSettings[email]["data-selected-workplaces"]
		userWebSettingsSync.RUnlock()
		dataWorkplaces = append(dataWorkplaces, TableWorkplaceSelection{
			WorkplaceName:      workplace.Name,
			WorkplaceSelection: getWorkplaceWebSelection(selectedWorkplace, workplace.Name),
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
	tmpl := template.Must(template.ParseFiles("./html/data.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("DATA", "Page sent in "+time.Since(timer).String())
}

func getWorkplaceWebSelection(selectedWorkplaces string, workplace string) string {
	for _, selectedWorkplace := range strings.Split(selectedWorkplaces, ";") {
		if selectedWorkplace == workplace {
			return "selected"
		}
	}
	return ""
}

func getSelected(selection string, menu string) string {
	if selection == menu {
		return "selected"
	}
	return ""
}

func loadTableData(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	timer := time.Now()
	go updatePageCount("data")
	email, _, _ := request.BasicAuth()
	logInfo("DATA", "Loading table for "+email)
	var data DataPageInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("DATA", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "ERR: Error parsing data, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("DATA", "Loading table ended")
		return
	}
	logInfo("DATA", "Loading data for "+data.Data+" for "+strconv.Itoa(len(data.Workplaces))+" workplaces")
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
		var responseData TableOutput
		responseData.Result = "ERR: Error parsing data, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("DATA", "Loading table ended")
		return
	}
	dateTo, err := time.ParseInLocation(layout, data.To, loc)
	dateTo = dateTo.In(time.UTC)
	if err != nil {
		logError("DATA", "Problem parsing date: "+data.To)
		var responseData TableOutput
		responseData.Result = "ERR: Error parsing data, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("DATA", "Loading table ended")
		return
	}
	logInfo("DATA", "From "+dateFrom.String()+" to "+dateTo.String())
	updateUserWebSettings(email, "data-selected-type", data.Data)
	updateUserWebSettings(email, "data-selected-from", data.From)
	updateUserWebSettings(email, "data-selected-to", data.To)
	userWebSettingsSync.RLock()
	workplaceNames := cachedUserWebSettings[email]["data-selected-workplaces"]
	userWebSettingsSync.RUnlock()
	workplaceIds := getWorkplaceIds(workplaceNames)

	logInfo("DATA", "Preprocessing takes "+time.Since(timer).String())
	switch data.Data {
	case "alarms":
		loadAlarmsTable(writer, workplaceIds, dateFrom, dateTo, email)
	case "breakdowns":
		loadBreakdownTable(writer, workplaceIds, dateFrom, dateTo, email)
	case "downtimes":
		loadDowntimesTable(writer, workplaceIds, dateFrom, dateTo, email)
	case "faults":
		loadFaultsTable(writer, workplaceIds, dateFrom, dateTo, email)
	case "orders":
		loadOrdersTable(writer, workplaceIds, dateFrom, dateTo, email)
	case "packages":
		loadPackagesTable(writer, workplaceIds, dateFrom, dateTo, email)
	case "parts":
		loadPartsTable(writer, workplaceIds, dateFrom, dateTo, email)
	case "states":
		loadStatesTable(writer, workplaceIds, dateFrom, dateTo, email)
	case "users":
		loadUsersTable(writer, workplaceIds, dateFrom, dateTo, email)
	case "system-statistics":
		loadSystemStatsTable(writer, dateFrom, dateTo, email)
	}
	logInfo("DATA", "Table loaded in "+time.Since(timer).String())
	return
}

func getWorkplaceIds(workplaceNames string) string {
	if len(workplaceNames) == 0 {
		return ""
	}
	workplaceNamesAsArray := strings.Split(workplaceNames, ";")
	workplaceNamesSql := `name in ('`
	for _, workplace := range workplaceNamesAsArray {
		workplaceNamesSql += workplace + `','`
	}
	workplaceNamesSql = strings.TrimSuffix(workplaceNamesSql, `,'`)
	workplaceNamesSql += ")"
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("DATA", "Problem opening database: "+err.Error())
		return ""
	}
	var workplaces []database.Workplace
	db.Select("id").Where(workplaceNamesSql).Find(&workplaces)
	workplaceIds := `workplace_id in ('`
	for _, workplace := range workplaces {
		workplaceIds += strconv.Itoa(int(workplace.ID)) + `','`
	}
	workplaceIds = strings.TrimSuffix(workplaceIds, `,'`)
	workplaceIds += ")"
	return workplaceIds
}
