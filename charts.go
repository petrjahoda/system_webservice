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
	"time"
)

type ChartsDataPageInput struct {
	Data      string
	Workplace string
	From      string
	To        string
}

type ChartsPageData struct {
	Version               string
	Company               string
	Alarms                string
	MenuOverview          string
	MenuWorkplaces        string
	MenuCharts            string
	MenuStatistics        string
	MenuData              string
	MenuSettings          string
	SelectionMenu         []ChartSelection
	Workplaces            []ChartWorkplaceSelection
	Compacted             string
	DataFilterPlaceholder string
	DateLocale            string
}

type ChartWorkplaceSelection struct {
	WorkplaceName      string
	WorkplaceSelection string
}
type ChartSelection struct {
	SelectionName string
	Selection     string
}

type ChartDataPageOutput struct {
	Result     string
	Type       string
	AnalogData []PortData
}

type PortData struct {
	PortName  string
	PortColor string
	PortData  []Data
}

type Data struct {
	Time  int64
	Value float32
}

func charts(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	timer := time.Now()
	email, _, _ := request.BasicAuth()
	logInfo("CHARTS", "Sending data page to "+cachedUsersByEmail[email].FirstName+" "+cachedUsersByEmail[email].SecondName)
	var data ChartsPageData
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

	data.SelectionMenu = append(data.SelectionMenu, ChartSelection{
		SelectionName: getLocale(email, "combined-chart"),
		Selection:     getSelected(cachedUserSettings[email].dataSelection, "combined-chart"),
	})
	data.SelectionMenu = append(data.SelectionMenu, ChartSelection{
		SelectionName: getLocale(email, "timeline-chart"),
		Selection:     getSelected(cachedUserSettings[email].dataSelection, "timeline-chart"),
	})
	data.SelectionMenu = append(data.SelectionMenu, ChartSelection{
		SelectionName: getLocale(email, "analog-data"),
		Selection:     getSelected(cachedUserSettings[email].dataSelection, "analog-data"),
	})
	data.SelectionMenu = append(data.SelectionMenu, ChartSelection{
		SelectionName: getLocale(email, "digital-data"),
		Selection:     getSelected(cachedUserSettings[email].dataSelection, "digital-data"),
	})
	data.SelectionMenu = append(data.SelectionMenu, ChartSelection{
		SelectionName: getLocale(email, "production-chart"),
		Selection:     getSelected(cachedUserSettings[email].dataSelection, "production-chart"),
	})
	data.SelectionMenu = append(data.SelectionMenu, ChartSelection{
		SelectionName: getLocale(email, "consumption-chart"),
		Selection:     getSelected(cachedUserSettings[email].dataSelection, "consumption-chart"),
	})

	var dataWorkplaces []ChartWorkplaceSelection
	for _, workplace := range cachedWorkplacesById {
		dataWorkplaces = append(dataWorkplaces, ChartWorkplaceSelection{
			WorkplaceName:      workplace.Name,
			WorkplaceSelection: getWorkplaceSelection(cachedUserSettings[email].selectedWorkplaces, workplace.Name),
		})
	}
	sort.Slice(dataWorkplaces, func(i, j int) bool {
		return dataWorkplaces[i].WorkplaceName < dataWorkplaces[j].WorkplaceName
	})
	data.Workplaces = dataWorkplaces
	tmpl := template.Must(template.ParseFiles("./html/charts.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("CHARTS", "Charts page sent in "+time.Since(timer).String())
}

func getChartData(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	timer := time.Now()
	email, _, _ := request.BasicAuth()
	logInfo("CHARTS", "Sending chart data to "+cachedUsersByEmail[email].FirstName+" "+cachedUsersByEmail[email].SecondName)
	var data ChartsDataPageInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("CHARTS", "Error parsing data: "+err.Error())
		var responseData DataPageOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("CHARTS", "Processing data ended")
		return
	}
	logInfo("CHARTS", "Loading chart data for "+data.Data+" for "+data.Workplace)
	loc, err := time.LoadLocation(location)
	if err != nil {
		logError("CHARTS", "Problem loading timezone, setting Europe/Prague")
		loc, _ = time.LoadLocation("Europe/Prague")
	}
	layout := "2006-01-02;15:04:05"
	dateFrom, err := time.ParseInLocation(layout, data.From, loc)
	dateFrom = dateFrom.In(time.UTC)
	if err != nil {
		logError("CHARTS", "Problem parsing date: "+data.From)
		var responseData DataPageOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("CHARTS", "Processing data ended")
		return
	}
	dateTo, err := time.ParseInLocation(layout, data.To, loc)
	dateTo = dateTo.In(time.UTC)
	if err != nil {
		logError("CHARTS", "Problem parsing date: "+data.To)
		var responseData DataPageOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("CHARTS", "Processing data ended")
		return
	}
	logInfo("CHARTS", "From "+dateFrom.String()+" to "+dateTo.String())
	logInfo("CHARTS", "Preprocessing takes "+time.Since(timer).String())
	locale := getLocaleForChart(data)
	switch locale.Name {
	case "combined-chart":
		//processAlarms(writer, workplaceIds, dateFrom, dateTo, email)
	case "timeline-chart":
		//processBreakdowns(writer, workplaceIds, dateFrom, dateTo, email)
	case "analog-data":
		processAnalogData(writer, data.Workplace, dateFrom, dateTo, email, locale.Name)
	case "digital-data":
		//processFaults(writer, workplaceIds, dateFrom, dateTo, email)
	case "production-chart":
		//processOrders(writer, workplaceIds, dateFrom, dateTo, email)
	case "consumption-chart":
		//processPackages(writer, workplaceIds, dateFrom, dateTo, email)
	}
	logInfo("CHARTS", "Chart data sent in "+time.Since(timer).String())
	return
}

func getLocaleForChart(data ChartsDataPageInput) database.Locale {
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("DATA", "Problem opening database: "+err.Error())
	}
	var locale database.Locale
	db.Select("name").Where("cs_cz like (?)", data.Data).Or("de_de like (?)", data.Data).Or("en_us like (?)", data.Data).Or("es_es like (?)", data.Data).Or("fr_fr like (?)", data.Data).Or("it_it like (?)", data.Data).Or("pl_pl like (?)", data.Data).Or("pt_pt like (?)", data.Data).Or("sk_sk like (?)", data.Data).Or("ru_ru like (?)", data.Data).Find(&locale)
	return locale
}
