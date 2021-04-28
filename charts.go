package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
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
	Information           string
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
	UserEmail             string
	UserName              string
}

type ChartWorkplaceSelection struct {
	WorkplaceName      string
	WorkplaceSelection string
}
type ChartSelection struct {
	SelectionName  string
	SelectionValue string
	Selection      string
}

type ChartDataPageOutput struct {
	Result     string
	Locale     string
	Type       string
	AnalogData []PortData
	OrderData  []TerminalData
}

type TerminalData struct {
	Name          string
	Color         string
	FromDate      int64
	ToDate        int64
	DataName      string
	OperationName string
	ProductName   string
	AverageCycle  float32
	CountOk       int
	CountNok      int
	Note          string
}

type PortData struct {
	PortName  string
	PortColor string
	PortData  []Data
}

type Data struct {
	Time  time.Time
	Value float32
}

func charts(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	timer := time.Now()
	go updatePageCount("charts")
	email, _, _ := request.BasicAuth()
	logInfo("CHARTS", "Sending data page to "+cachedUsersByEmail[email].FirstName+" "+cachedUsersByEmail[email].SecondName)
	var data ChartsPageData
	data.Version = version
	data.DateLocale = cachedLocales[cachedUsersByEmail[email].Locale]
	data.UserEmail = email
	data.UserName = cachedUsersByEmail[email].FirstName + " " + cachedUsersByEmail[email].SecondName
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
		SelectionName:  getLocale(email, "combined-chart"),
		SelectionValue: "combined-chart",
		Selection:      getSelected(cachedUserSettings[email].dataSelection, "combined-chart"),
	})
	data.SelectionMenu = append(data.SelectionMenu, ChartSelection{
		SelectionName:  getLocale(email, "timeline-chart"),
		SelectionValue: "timeline-chart",
		Selection:      getSelected(cachedUserSettings[email].dataSelection, "timeline-chart"),
	})
	data.SelectionMenu = append(data.SelectionMenu, ChartSelection{
		SelectionName:  getLocale(email, "analog-data"),
		SelectionValue: "analog-data",
		Selection:      getSelected(cachedUserSettings[email].dataSelection, "analog-data"),
	})
	data.SelectionMenu = append(data.SelectionMenu, ChartSelection{
		SelectionName:  getLocale(email, "digital-data"),
		SelectionValue: "digital-data",
		Selection:      getSelected(cachedUserSettings[email].dataSelection, "digital-data"),
	})
	data.SelectionMenu = append(data.SelectionMenu, ChartSelection{
		SelectionName:  getLocale(email, "production-chart"),
		SelectionValue: "production-chart",
		Selection:      getSelected(cachedUserSettings[email].dataSelection, "production-chart"),
	})
	data.SelectionMenu = append(data.SelectionMenu, ChartSelection{
		SelectionName:  getLocale(email, "consumption-chart"),
		SelectionValue: "consumption-chart",
		Selection:      getSelected(cachedUserSettings[email].dataSelection, "consumption-chart"),
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
	data.Information = "INF: Page processed in " + time.Since(timer).String()
	tmpl := template.Must(template.ParseFiles("./html/charts.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("CHARTS", "Charts page sent in "+time.Since(timer).String())
}

func loadChartData(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	timer := time.Now()
	email, _, _ := request.BasicAuth()
	logInfo("CHARTS", "Sending chart data to "+cachedUsersByEmail[email].FirstName+" "+cachedUsersByEmail[email].SecondName)
	var data ChartsDataPageInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("CHARTS", "Error parsing data: "+err.Error())
		var responseData TableOutput
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
	layout := "2006-01-02T15:04"
	dateFrom, err := time.ParseInLocation(layout, data.From, loc)
	dateFrom = dateFrom.In(time.UTC)
	if err != nil {
		logError("CHARTS", "Problem parsing date: "+data.From)
		var responseData TableOutput
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
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("CHARTS", "Processing data ended")
		return
	}
	logInfo("CHARTS", "From "+dateFrom.String()+" to "+dateTo.String())
	logInfo("CHARTS", "Preprocessing takes "+time.Since(timer).String())
	switch data.Data {
	case "combined-chart":
	case "timeline-chart":
	case "analog-data":
		processAnalogData(writer, data.Workplace, dateFrom, dateTo, email, data.Data)
	case "digital-data":
	case "production-chart":
	case "consumption-chart":
	}
	logInfo("CHARTS", "Chart data sent in "+time.Since(timer).String())
	return
}
