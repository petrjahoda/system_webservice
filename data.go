package main

import (
	"encoding/json"
	"fmt"
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

type DataPageData struct {
	Version               string
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
	Compacted             string
	DataFilterPlaceholder string
	DateLocale            string
}

type TableWorkplaceSelection struct {
	WorkplaceName      string
	WorkplaceSelection string
}
type TableSelection struct {
	SelectionName  string
	SelectionValue string
	Selection      string
}

type TableDataPageInput struct {
	Data       string
	Workplaces []string
	From       string
	To         string
}

type DataPageOutput struct {
	Result string
}

type TableData struct {
	DataTableSearchTitle    string
	DataTableInfoTitle      string
	DataTableRowsCountTitle string
	TableHeader             []HeaderCell
	TableRows               []TableRow
}
type TableRow struct {
	TableCell []TableCell
}

type TableCell struct {
	CellName string
}

type HeaderCell struct {
	HeaderName string
}

func data(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	timer := time.Now()
	email, _, _ := request.BasicAuth()
	logInfo("DATA", "Sending data page to "+cachedUsersByEmail[email].FirstName+" "+cachedUsersByEmail[email].SecondName)
	var data DataPageData
	data.Version = version
	data.DateLocale = cachedLocales[cachedUsersByEmail[email].Locale]
	data.Company = cachedCompanyName
	data.MenuOverview = getLocale(email, "menu-overview")
	data.MenuWorkplaces = getLocale(email, "menu-workplaces")
	data.MenuCharts = getLocale(email, "menu-charts")
	data.MenuStatistics = getLocale(email, "menu-statistics")
	data.MenuData = getLocale(email, "menu-data")
	data.MenuSettings = getLocale(email, "menu-settings")
	data.DataFilterPlaceholder = getLocale(email, "data-table-search-title")
	data.Compacted = cachedUserSettings[email].menuState
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
		SelectionName:  getLocale(email, "states"),
		SelectionValue: "states",
		Selection:      getSelected(cachedUserSettings[email].dataSelection, "states"),
	})
	if cachedUsersByEmail[email].UserTypeID == 2 {
		logInfo("DATA", "Adding data menu for administrator")
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "users"),
			SelectionValue: "users",
			Selection:      getSelected(cachedUserSettings[email].dataSelection, "users"),
		})
		data.SelectionMenu = append(data.SelectionMenu, TableSelection{
			SelectionName:  getLocale(email, "system-statistics"),
			SelectionValue: "system-statistics",
			Selection:      getSelected(cachedUserSettings[email].dataSelection, "system-statistics"),
		})
	}
	var dataWorkplaces []TableWorkplaceSelection
	for _, workplace := range cachedWorkplacesById {
		dataWorkplaces = append(dataWorkplaces, TableWorkplaceSelection{
			WorkplaceName:      workplace.Name,
			WorkplaceSelection: getWorkplaceSelection(cachedUserSettings[email].selectedWorkplaces, workplace.Name),
		})
	}
	sort.Slice(dataWorkplaces, func(i, j int) bool {
		return dataWorkplaces[i].WorkplaceName < dataWorkplaces[j].WorkplaceName
	})
	data.Workplaces = dataWorkplaces
	tmpl := template.Must(template.ParseFiles("./html/Data.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("DATA", "Date page sent in "+time.Since(timer).String())
}

func getWorkplaceSelection(selectedWorkplaces []string, workplace string) string {
	for _, selectedWorkplace := range selectedWorkplaces {
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

func getTableData(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	timer := time.Now()
	email, _, _ := request.BasicAuth()
	logInfo("DATA", "Sending table data to "+cachedUsersByEmail[email].FirstName+" "+cachedUsersByEmail[email].SecondName)
	var data TableDataPageInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("DATA", "Error parsing data: "+err.Error())
		var responseData DataPageOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("DATA", "Processing data ended")
		return
	}
	logInfo("DATA", "Loading data for "+data.Data+" for "+strconv.Itoa(len(data.Workplaces))+" workplaces")
	loc, err := time.LoadLocation(location)
	if err != nil {
		logError("DATA", "Problem loading timezone, setting Europe/Prague")
		loc, _ = time.LoadLocation("Europe/Prague")
	}
	layout := "2006-01-02T15:04"
	dateFrom, err := time.ParseInLocation(layout, data.From, loc)
	dateFrom = dateFrom.In(time.UTC)
	if err != nil {
		logError("DATA", "Problem parsing date: "+data.From)
		var responseData DataPageOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("DATA", "Processing data ended")
		return
	}
	dateTo, err := time.ParseInLocation(layout, data.To, loc)
	dateTo = dateTo.In(time.UTC)
	if err != nil {
		logError("DATA", "Problem parsing date: "+data.To)
		var responseData DataPageOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("DATA", "Processing data ended")
		return
	}
	logInfo("DATA", "From "+dateFrom.String()+" to "+dateTo.String())
	fmt.Println(data.Data)
	workplaceIds := getWorkplaceIds(data, cachedUsersByEmail[email].Locale)
	updateUserDataSettings(email, data.Data, data)
	logInfo("DATA", "Preprocessing takes "+time.Since(timer).String())
	switch data.Data {
	case "alarms":
		processAlarms(writer, workplaceIds, dateFrom, dateTo, email)
	case "breakdowns":
		processBreakdowns(writer, workplaceIds, dateFrom, dateTo, email)
	case "downtimes":
		processDowntimes(writer, workplaceIds, dateFrom, dateTo, email)
	case "faults":
		processFaults(writer, workplaceIds, dateFrom, dateTo, email)
	case "orders":
		processOrders(writer, workplaceIds, dateFrom, dateTo, email)
	case "packages":
		processPackages(writer, workplaceIds, dateFrom, dateTo, email)
	case "parts":
		processParts(writer, workplaceIds, dateFrom, dateTo, email)
	case "states":
		processStates(writer, workplaceIds, dateFrom, dateTo, email)
	case "users":
		processUsers(writer, workplaceIds, dateFrom, dateTo, email)
	case "system-statistics":
		processSystemStats(writer, dateFrom, dateTo, email)
	}
	logInfo("DATA", "Table data sent in "+time.Since(timer).String())
	return
}

func getWorkplaceIds(data TableDataPageInput, userLocale string) string {
	workplaceNames := `name in ('`
	for _, workplace := range data.Workplaces {
		workplaceNames += workplace + `','`
	}
	workplaceNames = strings.TrimSuffix(workplaceNames, `,'`)
	workplaceNames += ")"
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("DATA", "Problem opening database: "+err.Error())
		return ""
	}
	var workplaces []database.Workplace
	db.Select("id").Where(workplaceNames).Find(&workplaces)
	workplaceIds := `workplace_id in ('`
	for _, workplace := range workplaces {
		workplaceIds += strconv.Itoa(int(workplace.ID)) + `','`
	}
	workplaceIds = strings.TrimSuffix(workplaceIds, `,'`)
	workplaceIds += ")"
	return workplaceIds
}
