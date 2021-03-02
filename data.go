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
	SelectionMenu         []Selection
	Workplaces            []WorkplaceSelection
	Compacted             string
	DataFilterPlaceholder string
	DateLocale            string
}

type WorkplaceSelection struct {
	WorkplaceName      string
	WorkplaceSelection string
}
type Selection struct {
	SelectionName string
	Selection     string
}

type DataPageInput struct {
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
	data.SelectionMenu = append(data.SelectionMenu, Selection{
		SelectionName: getLocale(email, "alarms"),
		Selection:     getSelected(cachedUserSettings[email].dataSelection, "alarms"),
	})
	data.SelectionMenu = append(data.SelectionMenu, Selection{
		SelectionName: getLocale(email, "breakdowns"),
		Selection:     getSelected(cachedUserSettings[email].dataSelection, "breakdowns"),
	})
	data.SelectionMenu = append(data.SelectionMenu, Selection{
		SelectionName: getLocale(email, "downtimes"),
		Selection:     getSelected(cachedUserSettings[email].dataSelection, "downtimes"),
	})
	data.SelectionMenu = append(data.SelectionMenu, Selection{
		SelectionName: getLocale(email, "faults"),
		Selection:     getSelected(cachedUserSettings[email].dataSelection, "faults"),
	})

	data.SelectionMenu = append(data.SelectionMenu, Selection{
		SelectionName: getLocale(email, "orders"),
		Selection:     getSelected(cachedUserSettings[email].dataSelection, "orders"),
	})
	data.SelectionMenu = append(data.SelectionMenu, Selection{
		SelectionName: getLocale(email, "packages"),
		Selection:     getSelected(cachedUserSettings[email].dataSelection, "packages"),
	})
	data.SelectionMenu = append(data.SelectionMenu, Selection{
		SelectionName: getLocale(email, "parts"),
		Selection:     getSelected(cachedUserSettings[email].dataSelection, "parts"),
	})
	data.SelectionMenu = append(data.SelectionMenu, Selection{
		SelectionName: getLocale(email, "states"),
		Selection:     getSelected(cachedUserSettings[email].dataSelection, "states"),
	})
	if cachedUsersByEmail[email].UserTypeID == 2 {
		logInfo("DATA", "Adding data menu for administrator")
		data.SelectionMenu = append(data.SelectionMenu, Selection{
			SelectionName: getLocale(email, "users"),
			Selection:     getSelected(cachedUserSettings[email].dataSelection, "users"),
		})
		data.SelectionMenu = append(data.SelectionMenu, Selection{
			SelectionName: getLocale(email, "system-statistics"),
			Selection:     getSelected(cachedUserSettings[email].dataSelection, "system-statistics"),
		})
	}
	var dataWorkplaces []WorkplaceSelection
	for _, workplace := range cachedWorkplacesById {
		dataWorkplaces = append(dataWorkplaces, WorkplaceSelection{
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
	var data DataPageInput
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
	layout := "2006-01-02;15:04:05"
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
	workplaceIds, locale := getWorkplaceIds(data, err)
	updateUserDataSettings(email, locale, data)
	logInfo("DATA", "Preprocessing takes "+time.Since(timer).String())
	switch locale.Name {
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

func getWorkplaceIds(data DataPageInput, err error) (string, database.Locale) {
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
		return "", database.Locale{}
	}
	var workplaces []database.Workplace
	db.Select("id").Where(workplaceNames).Find(&workplaces)
	workplaceIds := `workplace_id in ('`
	for _, workplace := range workplaces {
		workplaceIds += strconv.Itoa(int(workplace.ID)) + `','`
	}
	workplaceIds = strings.TrimSuffix(workplaceIds, `,'`)
	workplaceIds += ")"
	var locale database.Locale
	db.Select("name").Where("cs_cz like (?)", data.Data).Or("de_de like (?)", data.Data).Or("en_us like (?)", data.Data).Or("es_es like (?)", data.Data).Or("fr_fr like (?)", data.Data).Or("it_it like (?)", data.Data).Or("pl_pl like (?)", data.Data).Or("pt_pt like (?)", data.Data).Or("sk_sk like (?)", data.Data).Or("ru_ru like (?)", data.Data).Find(&locale)
	return workplaceIds, locale
}
