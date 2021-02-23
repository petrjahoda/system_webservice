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
	"strconv"
	"strings"
	"time"
)

type DataPageData struct {
	Version        string
	Company        string
	Alarms         string
	MenuOverview   string
	MenuWorkplaces string
	MenuCharts     string
	MenuStatistics string
	MenuData       string
	MenuSettings   string
	SelectionMenu  []string
	Workplaces     []string
	Compacted      string
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

func data(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	ipAddress := strings.Split(request.RemoteAddr, ":")
	logInfo("MAIN", "Sending home page to "+ipAddress[0])
	email, _, _ := request.BasicAuth()
	var data DataPageData
	data.Version = version
	data.Company = cachedCompanyName
	data.MenuOverview = getLocale(email, "menu-overview")
	data.MenuWorkplaces = getLocale(email, "menu-workplaces")
	data.MenuCharts = getLocale(email, "menu-charts")
	data.MenuStatistics = getLocale(email, "menu-statistics")
	data.MenuData = getLocale(email, "menu-data")
	data.MenuSettings = getLocale(email, "menu-settings")
	data.Compacted = cachedUserSettings[email].menuState
	data.SelectionMenu = append(data.SelectionMenu, getLocale(email, "alarms"))
	data.SelectionMenu = append(data.SelectionMenu, getLocale(email, "breakdowns"))
	data.SelectionMenu = append(data.SelectionMenu, getLocale(email, "downtimes"))
	data.SelectionMenu = append(data.SelectionMenu, getLocale(email, "faults"))
	data.SelectionMenu = append(data.SelectionMenu, getLocale(email, "orders"))
	data.SelectionMenu = append(data.SelectionMenu, getLocale(email, "packages"))
	data.SelectionMenu = append(data.SelectionMenu, getLocale(email, "parts"))
	data.SelectionMenu = append(data.SelectionMenu, getLocale(email, "states"))
	if cachedUsers[email].UserTypeID == 2 {
		logInfo("MAIN", "Adding data menu for administrator")
		data.SelectionMenu = append(data.SelectionMenu, getLocale(email, "users"))
		data.SelectionMenu = append(data.SelectionMenu, getLocale(email, "system-statistics"))
	}

	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("MAIN", "Problem opening database: "+err.Error())
		return
	}
	var workplaces []database.Workplace
	db.Select("name").Find(&workplaces)
	for _, workplace := range workplaces {
		data.Workplaces = append(data.Workplaces, workplace.Name)
	}

	tmpl := template.Must(template.ParseFiles("./html/Data.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("MAIN", "Home page sent")
}

func getData(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var data DataPageInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("MAIN", "Error parsing data: "+err.Error())
		var responseData DataPageOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Processing chart data ended")
		return
	}
	logInfo("MAIN", "Loading data for "+data.Data)
	logInfo("MAIN", "Number of workplaces selected "+strconv.Itoa(len(data.Workplaces)))
	logInfo("MAIN", "From "+data.From)
	logInfo("MAIN", "To "+data.To)
	layout := "2006-01-02;15:04:05"
	dateFrom, err := time.Parse(layout, data.From)
	if err != nil {
		logError("MAIN", "Problem parsing date: "+data.From)
		var responseData DataPageOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Processing chart data ended")
		return
	}
	dateTo, err := time.Parse(layout, data.To)
	if err != nil {
		logError("MAIN", "Problem parsing date: "+data.To)
		var responseData DataPageOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Processing chart data ended")
		return
	}

	workplaceIds, locale := processDataFromDatabase(data, err)
	switch locale.Name {
	case "alarms":
		{
			db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
			sqlDB, _ := db.DB()
			defer sqlDB.Close()
			if err != nil {
				logError("MAIN", "Problem opening database: "+err.Error())

			}
			var records []database.AlarmRecord
			db.Where(workplaceIds).Where("date_time_start >= ?", dateFrom).Where("date_time_end <= ?", dateTo).Or("date_time_end is null").Find(&records)
			fmt.Println(len(records))
			email, _, _ := request.BasicAuth()
			fmt.Println(email)
		}
	case "breakdowns":
		{
			db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
			sqlDB, _ := db.DB()
			defer sqlDB.Close()
			if err != nil {
				logError("MAIN", "Problem opening database: "+err.Error())

			}
			var records []database.BreakdownRecord
			db.Where(workplaceIds).Where("date_time_start >= ?", dateFrom).Where("date_time_end <= ?", dateTo).Or("date_time_end is null").Find(&records)
			fmt.Println(len(records))
			email, _, _ := request.BasicAuth()
			fmt.Println(email)
		}
	case "downtimes":
		{
			db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
			sqlDB, _ := db.DB()
			defer sqlDB.Close()
			if err != nil {
				logError("MAIN", "Problem opening database: "+err.Error())

			}
			var records []database.DowntimeRecord
			db.Where(workplaceIds).Where("date_time_start >= ?", dateFrom).Where("date_time_end <= ?", dateTo).Or("date_time_end is null").Find(&records)
			fmt.Println(len(records))
			email, _, _ := request.BasicAuth()
			fmt.Println(email)
		}
	case "faults":
		{
			db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
			sqlDB, _ := db.DB()
			defer sqlDB.Close()
			if err != nil {
				logError("MAIN", "Problem opening database: "+err.Error())

			}
			var records []database.FaultRecord
			db.Where(workplaceIds).Where("date_time_start >= ?", dateFrom).Where("date_time_end <= ?", dateTo).Or("date_time_end is null").Find(&records)
			fmt.Println(len(records))
			email, _, _ := request.BasicAuth()
			fmt.Println(email)
		}
	case "orders":
		{
			db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
			sqlDB, _ := db.DB()
			defer sqlDB.Close()
			if err != nil {
				logError("MAIN", "Problem opening database: "+err.Error())

			}
			var records []database.OrderRecord
			db.Where(workplaceIds).Where("date_time_start >= ?", dateFrom).Where("date_time_end <= ?", dateTo).Or("date_time_end is null").Find(&records)
			fmt.Println(len(records))
			email, _, _ := request.BasicAuth()
			fmt.Println(email)
			// TODO: cache users by id
			// TODO: cache workplaces by id
			// TODO: cache orders by id
			// TODO: cache operations by id
			// TODO: cache workplace mode by id
			// TODO: cache workshifts by id
			// TODO: download user records and make map by order_record_id
			// TODO: create array of outcome orders
			// TODO: loop order records, for every order create new outcome order
			// TODO: 	get order name
			// TODO: 	get operation name
			// TODO: 	get workplace name
			// TODO: 	get workplace mode name
			// TODO: 	get workshift name
			// TODO: 	set average cycle
			// TODO: 	get user_id from map and get user name from cached users
			// TODO: 	set cavity
			// TODO: 	set ok
			// TODO: 	set nok
			// TODO: 	set note
			// TODO: add outcome order to array of outcome orders

		}
	case "packages":
		{
			db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
			sqlDB, _ := db.DB()
			defer sqlDB.Close()
			if err != nil {
				logError("MAIN", "Problem opening database: "+err.Error())

			}
			var records []database.PackageRecord
			db.Where(workplaceIds).Where("date_time_start >= ?", dateFrom).Where("date_time_end <= ?", dateTo).Or("date_time_end is null").Find(&records)
			fmt.Println(len(records))
			email, _, _ := request.BasicAuth()
			fmt.Println(email)
		}
	case "parts":
		{
			db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
			sqlDB, _ := db.DB()
			defer sqlDB.Close()
			if err != nil {
				logError("MAIN", "Problem opening database: "+err.Error())

			}
			var records []database.PartRecord
			db.Where(workplaceIds).Where("date_time_start >= ?", dateFrom).Where("date_time_end <= ?", dateTo).Or("date_time_end is null").Find(&records)
			fmt.Println(len(records))
			email, _, _ := request.BasicAuth()
			fmt.Println(email)
		}
	case "states":
		{
			db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
			sqlDB, _ := db.DB()
			defer sqlDB.Close()
			if err != nil {
				logError("MAIN", "Problem opening database: "+err.Error())

			}
			var records []database.StateRecord
			db.Where(workplaceIds).Where("date_time_start >= ?", dateFrom).Where("date_time_end <= ?", dateTo).Or("date_time_end is null").Find(&records)
			fmt.Println(len(records))
			email, _, _ := request.BasicAuth()
			fmt.Println(email)
		}
	case "users":
		{
			db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
			sqlDB, _ := db.DB()
			defer sqlDB.Close()
			if err != nil {
				logError("MAIN", "Problem opening database: "+err.Error())

			}
			var records []database.UserRecord
			db.Where(workplaceIds).Where("date_time_start >= ?", dateFrom).Where("date_time_end <= ?", dateTo).Or("date_time_end is null").Find(&records)
			fmt.Println(len(records))
			email, _, _ := request.BasicAuth()
			fmt.Println(email)
		}
	case "system-statistics":
		{
			db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
			sqlDB, _ := db.DB()
			defer sqlDB.Close()
			if err != nil {
				logError("MAIN", "Problem opening database: "+err.Error())

			}
			var records []database.SystemRecord
			db.Where(workplaceIds).Where("date_time_start >= ?", dateFrom).Where("date_time_end <= ?", dateTo).Or("date_time_end is null").Find(&records)
			fmt.Println(len(records))
			email, _, _ := request.BasicAuth()
			fmt.Println(email)
		}
	}
}

func processDataFromDatabase(data DataPageInput, err error) (string, database.Locale) {
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
		logError("MAIN", "Problem opening database: "+err.Error())
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
