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
	Result       string
	OrderRecords []OrderRecord
}

type OrderRecord struct {
	WorkplaceName     string
	OrderStart        string
	OrderEnd          string
	WorkplaceModeName string
	WorkshiftName     string
	UserName          string
	OrderName         string
	OperationName     string
	AverageCycle      float32
	Cavity            int
	Ok                int
	Nok               int
	Note              string
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
	ipAddress := strings.Split(request.RemoteAddr, ":")
	logInfo("MAIN", "Sending home page to "+ipAddress[0])
	email, _, _ := request.BasicAuth()
	var data DataPageData
	data.Version = version
	data.DateLocale = getLocaleForPickers(email)
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
		logInfo("MAIN", "Adding data menu for administrator")
		data.SelectionMenu = append(data.SelectionMenu, Selection{
			SelectionName: getLocale(email, "users"),
			Selection:     getSelected(cachedUserSettings[email].dataSelection, "users"),
		})
		data.SelectionMenu = append(data.SelectionMenu, Selection{
			SelectionName: getLocale(email, "system-statistics"),
			Selection:     getSelected(cachedUserSettings[email].dataSelection, "system-statistics"),
		})
	}
	for _, workplace := range cachedWorkplacesById {
		data.Workplaces = append(data.Workplaces, WorkplaceSelection{
			WorkplaceName:      workplace.Name,
			WorkplaceSelection: getWorkplaceSelection(cachedUserSettings[email].selectedWorkplaces, workplace.Name),
		})
	}

	tmpl := template.Must(template.ParseFiles("./html/Data.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("MAIN", "Home page sent")
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

func getLocaleForPickers(email string) string {
	var pickerLocale string
	user, _ := cachedUsersByEmail[email]
	switch user.Locale {
	case "CsCZ":
		{
			pickerLocale = "cs-CZ"
		}
	case "DeDE":
		{
			pickerLocale = "de-DE"
		}
	case "EnUS":
		{
			pickerLocale = "en-US"
		}
	case "EsES":
		{
			pickerLocale = "es-MX"
		}
	case "FrFR":
		{
			pickerLocale = "fr-FR"
		}
	case "ItIT":
		{
			pickerLocale = "it-IT"
		}
	case "PlPL":
		{
			pickerLocale = "pl-PL"
		}
	case "PtPT":
		{
			pickerLocale = "pt-BR"
		}
	case "SkSK":
		{
			pickerLocale = "sk-SK"
		}
	case "RuRU":
		{
			pickerLocale = "ru-RU"
		}
	default:
		{
			pickerLocale = "en-US"
		}
	}
	return pickerLocale

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
		logInfo("MAIN", "Processing data ended")
		return
	}
	logInfo("MAIN", "Loading data for "+data.Data)
	logInfo("MAIN", "Number of workplaces selected "+strconv.Itoa(len(data.Workplaces)))

	loc, err := time.LoadLocation(location)
	if err != nil {
		logError("MAIN", "Problem loading timezone, setting Europe/Prague")
		loc, _ = time.LoadLocation("Europe/Prague")
	}
	layout := "2006-01-02;15:04:05"
	dateFrom, err := time.ParseInLocation(layout, data.From, loc)
	dateFrom = dateFrom.In(time.UTC)
	if err != nil {
		logError("MAIN", "Problem parsing date: "+data.From)
		var responseData DataPageOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Processing data ended")
		return
	}
	dateTo, err := time.ParseInLocation(layout, data.To, loc)
	dateTo = dateTo.In(time.UTC)
	if err != nil {
		logError("MAIN", "Problem parsing date: "+data.To)
		var responseData DataPageOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Processing data ended")
		return
	}

	logInfo("MAIN", "From "+dateFrom.String())
	logInfo("MAIN", "To "+dateTo.String())
	email, _, _ := request.BasicAuth()
	workplaceIds, locale := processDataFromDatabase(data, err)

	userSettingsSync.Lock()
	settings := cachedUserSettings[email]
	settings.dataSelection = locale.Name
	settings.selectedWorkplaces = data.Workplaces
	cachedUserSettings[email] = settings
	userSettingsSync.Unlock()

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
			var data TableData

			var headerCell HeaderCell
			headerCell.HeaderName = "a"
			data.TableHeader = append(data.TableHeader, headerCell)
			var headerCell2 HeaderCell
			headerCell2.HeaderName = "b"
			data.TableHeader = append(data.TableHeader, headerCell2)

			var tableRow1 TableRow
			var tableRow2 TableRow

			var tableCellA1 TableCell
			tableCellA1.CellName = "testA1"

			var tableCellA2 TableCell
			tableCellA2.CellName = "testA2"

			var tableCellB1 TableCell
			tableCellB1.CellName = "testB1"

			var tableCellB2 TableCell
			tableCellB2.CellName = "testB2"

			tableRow1.TableCell = append(tableRow1.TableCell, tableCellA1)
			tableRow1.TableCell = append(tableRow1.TableCell, tableCellA2)

			tableRow2.TableCell = append(tableRow2.TableCell, tableCellB1)
			tableRow2.TableCell = append(tableRow2.TableCell, tableCellB2)

			data.TableRows = append(data.TableRows, tableRow1)
			data.TableRows = append(data.TableRows, tableRow2)

			tmpl := template.Must(template.ParseFiles("./html/table.html"))
			_ = tmpl.Execute(writer, data)
			logInfo("MAIN", "Home page sent for "+email)
			return
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
			processOrders(writer, workplaceIds, dateFrom, dateTo, email)
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
	return
}

func processOrders(writer http.ResponseWriter, workplaceIds string, dateFrom time.Time, dateTo time.Time, email string) {
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("MAIN", "Problem opening database: "+err.Error())
		var responseData DataPageOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Processing data ended")
		return
	}
	var orderRecords []database.OrderRecord
	db.Where(workplaceIds).Where("date_time_start >= ?", dateFrom).Where("date_time_start <= ?", dateTo).Order("date_time_start desc").Find(&orderRecords)
	var userRecords []database.UserRecord
	db.Where(workplaceIds).Where("date_time_start >= ?", dateFrom).Where("date_time_start <= ?", dateTo).Find(&userRecords)
	var userRecordsByRecordId = map[int]database.UserRecord{}
	for _, record := range userRecords {
		userRecordsByRecordId[record.OrderRecordID] = record
	}
	var data TableData
	data.DataTableSearchTitle = getLocale(email, "data-table-search-title")
	data.DataTableInfoTitle = getLocale(email, "data-table-info-title")
	data.DataTableRowsCountTitle = getLocale(email, "data-table-rows-count-title")
	addTableHeaders(email, &data)
	for _, record := range orderRecords {
		addTableRow(record, userRecordsByRecordId, &data)
	}
	tmpl := template.Must(template.ParseFiles("./html/table.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("MAIN", "Order data sent for "+email)
}

func addTableRow(record database.OrderRecord, userRecordsByRecordId map[int]database.UserRecord, data *TableData) {
	var tableRow TableRow
	workplaceNameCell := TableCell{CellName: cachedWorkplacesById[uint(record.WorkplaceID)].Name}
	tableRow.TableCell = append(tableRow.TableCell, workplaceNameCell)
	orderStart := TableCell{CellName: record.DateTimeStart.Format("2006-01-02 15:04:05")}
	tableRow.TableCell = append(tableRow.TableCell, orderStart)
	if record.DateTimeEnd.Time.IsZero() {
		orderEnd := TableCell{CellName: time.Now().Format("2006-01-02 15:04:05") + " +"}
		tableRow.TableCell = append(tableRow.TableCell, orderEnd)
	} else {
		orderEnd := TableCell{CellName: record.DateTimeEnd.Time.Format("2006-01-02 15:04:05")}
		tableRow.TableCell = append(tableRow.TableCell, orderEnd)
	}
	workplaceModeNameCell := TableCell{CellName: cachedWorkplaceModesById[uint(record.WorkplaceModeID)].Name}
	tableRow.TableCell = append(tableRow.TableCell, workplaceModeNameCell)
	workshiftName := TableCell{CellName: cachedWorkshiftsById[uint(record.WorkshiftID)].Name}
	tableRow.TableCell = append(tableRow.TableCell, workshiftName)
	actualUserId := userRecordsByRecordId[int(record.ID)].UserID
	userName := TableCell{CellName: cachedUsersById[uint(actualUserId)].FirstName + " " + cachedUsersById[uint(actualUserId)].SecondName}
	tableRow.TableCell = append(tableRow.TableCell, userName)
	orderName := TableCell{CellName: cachedOrdersById[uint(record.OrderID)].Name}
	tableRow.TableCell = append(tableRow.TableCell, orderName)
	operationName := TableCell{CellName: cachedOperationsById[uint(record.OperationID)].Name}
	tableRow.TableCell = append(tableRow.TableCell, operationName)
	averageCycleAsString := strconv.FormatFloat(float64(record.AverageCycle), 'f', 2, 64)
	averageCycle := TableCell{CellName: averageCycleAsString + "s"}
	tableRow.TableCell = append(tableRow.TableCell, averageCycle)
	cavityAsString := strconv.Itoa(record.Cavity)
	cavity := TableCell{CellName: cavityAsString}
	tableRow.TableCell = append(tableRow.TableCell, cavity)
	okAsString := strconv.Itoa(record.CountOk)
	ok := TableCell{CellName: okAsString}
	tableRow.TableCell = append(tableRow.TableCell, ok)
	nokAsString := strconv.Itoa(record.CountNok)
	nok := TableCell{CellName: nokAsString}
	tableRow.TableCell = append(tableRow.TableCell, nok)
	note := TableCell{CellName: record.Note}
	tableRow.TableCell = append(tableRow.TableCell, note)
	data.TableRows = append(data.TableRows, tableRow)
}

func addTableHeaders(email string, data *TableData) {
	workplaceName := HeaderCell{HeaderName: getLocale(email, "workplace-name")}
	data.TableHeader = append(data.TableHeader, workplaceName)
	orderStart := HeaderCell{HeaderName: getLocale(email, "order-start")}
	data.TableHeader = append(data.TableHeader, orderStart)
	orderEnd := HeaderCell{HeaderName: getLocale(email, "order-end")}
	data.TableHeader = append(data.TableHeader, orderEnd)
	workplaceModeName := HeaderCell{HeaderName: getLocale(email, "workplacemode-name")}
	data.TableHeader = append(data.TableHeader, workplaceModeName)
	workshiftName := HeaderCell{HeaderName: getLocale(email, "workshift-name")}
	data.TableHeader = append(data.TableHeader, workshiftName)
	userName := HeaderCell{HeaderName: getLocale(email, "user-name")}
	data.TableHeader = append(data.TableHeader, userName)
	orderName := HeaderCell{HeaderName: getLocale(email, "order-name")}
	data.TableHeader = append(data.TableHeader, orderName)
	operationName := HeaderCell{HeaderName: getLocale(email, "operation-name")}
	data.TableHeader = append(data.TableHeader, operationName)
	cycleName := HeaderCell{HeaderName: getLocale(email, "cycle-name")}
	data.TableHeader = append(data.TableHeader, cycleName)
	cavityName := HeaderCell{HeaderName: getLocale(email, "cavity-name")}
	data.TableHeader = append(data.TableHeader, cavityName)
	goodPcsName := HeaderCell{HeaderName: getLocale(email, "good-pieces-name")}
	data.TableHeader = append(data.TableHeader, goodPcsName)
	badPcsName := HeaderCell{HeaderName: getLocale(email, "bad-pieces-name")}
	data.TableHeader = append(data.TableHeader, badPcsName)
	noteName := HeaderCell{HeaderName: getLocale(email, "note-name")}
	data.TableHeader = append(data.TableHeader, noteName)
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
