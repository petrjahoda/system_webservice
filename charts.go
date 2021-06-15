package main

import (
	"encoding/json"
	"fmt"
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
	Flash     string
	Terminal  string
}

type ChartsPageData struct {
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
	SelectionMenu         []ChartSelection
	Workplaces            []ChartWorkplaceSelection
	DataFilterPlaceholder string
	DateLocale            string
	UserEmail             string
	UserName              string
	DateFrom              string
	DateTo                string
	FlashClass            string
	TerminalClass         string
}

type ChartWorkplaceSelection struct {
	WorkplaceName      string
	WorkplaceValue     string
	WorkplaceSelection string
}
type ChartSelection struct {
	SelectionName  string
	SelectionValue string
	Selection      string
}

type ChartDataPageOutput struct {
	Result           string
	Locale           string
	Type             string
	ChartData        []PortData
	OrderData        []TerminalData
	DowntimeData     []TerminalData
	BreakdownData    []TerminalData
	AlarmData        []TerminalData
	UserData         []TerminalData
	OrdersLocale     string
	DowntimesLocale  string
	BreakdownsLocale string
	AlarmsLocale     string
	UsersLocale      string
}

type TerminalData struct {
	Color       string
	FromDate    int64
	ToDate      int64
	Information string
	Note        string
}

type PortData struct {
	PortType    string
	PortName    string
	PortColor   string
	AnalogData  []Data
	DigitalData []Data
}

type Data struct {
	Time  int64
	Value float32
}

func charts(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	timer := time.Now()
	go updatePageCount("charts")
	email, _, _ := request.BasicAuth()
	go updateWebUserRecord("charts", email)
	logInfo("CHARTS", "Sending page to "+cachedUsersByEmail[email].FirstName+" "+cachedUsersByEmail[email].SecondName)
	var data ChartsPageData
	data.Version = version
	localesSync.Lock()
	data.DateLocale = cachedLocales[cachedUsersByEmail[email].Locale]
	localesSync.Unlock()
	data.UserEmail = email
	data.UserName = cachedUsersByEmail[email].FirstName + " " + cachedUsersByEmail[email].SecondName
	companyNameSync.Lock()
	data.Company = cachedCompanyName
	companyNameSync.Unlock()
	data.MenuOverview = getLocale(email, "menu-overview")
	data.MenuWorkplaces = getLocale(email, "menu-workplaces")
	data.MenuCharts = getLocale(email, "menu-charts")
	data.MenuStatistics = getLocale(email, "menu-statistics")
	data.MenuData = getLocale(email, "menu-data")
	data.MenuSettings = getLocale(email, "menu-settings")
	data.DataFilterPlaceholder = getLocale(email, "data-table-search-title")
	data.SelectionMenu = append(data.SelectionMenu, ChartSelection{
		SelectionName:  getLocale(email, "combined-chart"),
		SelectionValue: "combined-chart",
		Selection:      getSelected(cachedUserWebSettings[email]["charts-selected-chart"], "combined-chart"),
	})
	data.SelectionMenu = append(data.SelectionMenu, ChartSelection{
		SelectionName:  getLocale(email, "production-chart"),
		SelectionValue: "production-chart",
		Selection:      getSelected(cachedUserWebSettings[email]["charts-selected-chart"], "production-chart"),
	})
	data.SelectionMenu = append(data.SelectionMenu, ChartSelection{
		SelectionName:  getLocale(email, "analog-data"),
		SelectionValue: "analog-data",
		Selection:      getSelected(cachedUserWebSettings[email]["charts-selected-chart"], "analog-data"),
	})
	data.SelectionMenu = append(data.SelectionMenu, ChartSelection{
		SelectionName:  getLocale(email, "digital-data"),
		SelectionValue: "digital-data",
		Selection:      getSelected(cachedUserWebSettings[email]["charts-selected-chart"], "digital-data"),
	})
	var dataWorkplaces []ChartWorkplaceSelection
	for _, workplace := range cachedWorkplacesById {
		dataWorkplaces = append(dataWorkplaces, ChartWorkplaceSelection{
			WorkplaceName:      workplace.Name,
			WorkplaceValue:     workplace.Name,
			WorkplaceSelection: getSelected(cachedUserWebSettings[email]["charts-selected-workplace"], workplace.Name),
		})
	}
	sort.Slice(dataWorkplaces, func(i, j int) bool {
		return dataWorkplaces[i].WorkplaceName < dataWorkplaces[j].WorkplaceName
	})
	data.Workplaces = dataWorkplaces
	data.DateFrom = cachedUserWebSettings[email]["charts-selected-from"]
	data.DateTo = cachedUserWebSettings[email]["charts-selected-to"]
	fmt.Println("FLASH", cachedUserWebSettings[email]["charts-selected-flash"])
	fmt.Println("TERMINAL", cachedUserWebSettings[email]["charts-selected-terminal"])
	data.FlashClass = "mif-flash-on"
	if len(cachedUserWebSettings[email]["charts-selected-flash"]) > 0 {
		data.FlashClass = cachedUserWebSettings[email]["charts-selected-flash"]
	}
	data.TerminalClass = "mif-phonelink"
	if len(cachedUserWebSettings[email]["charts-selected-terminal"]) > 0 {
		data.TerminalClass = cachedUserWebSettings[email]["charts-selected-terminal"]
	}
	data.Software = cachedSoftwareName
	data.Information = "INF: Page processed in " + time.Since(timer).String()
	tmpl := template.Must(template.ParseFiles("./html/charts.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("CHARTS", "Page sent in "+time.Since(timer).String())
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
	companyNameSync.Lock()
	loc, err := time.LoadLocation(location)
	companyNameSync.Unlock()
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
	updateUserWebSettings(email, "charts-selected-workplace", data.Workplace)
	updateUserWebSettings(email, "charts-selected-chart", data.Data)
	updateUserWebSettings(email, "charts-selected-from", data.From)
	updateUserWebSettings(email, "charts-selected-to", data.To)
	updateUserWebSettings(email, "charts-selected-flash", data.Flash)
	updateUserWebSettings(email, "charts-selected-terminal", data.Terminal)
	switch data.Data {
	case "combined-chart":
		processCombinedChart(writer, data.Workplace, dateFrom, dateTo, email, data.Data)
	case "timeline-chart":
	case "analog-data":
		processAnalogData(writer, data.Workplace, dateFrom, dateTo, email, data.Data)
	case "digital-data":
		processDigitalData(writer, data.Workplace, dateFrom, dateTo, email, data.Data)
	case "production-chart":
		processProductionChart(writer, data.Workplace, dateFrom, dateTo, email, data.Data)
	case "consumption-chart":
	}
	logInfo("CHARTS", "Chart data sent in "+time.Since(timer).String())
	return
}
