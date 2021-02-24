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
	Result       string
	OrderRecords []OrderRecord
}

type OrderRecord struct {
	WorkplaceName     string
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
	TableHeader []HeaderCell
	TableRows   []TableRow
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
	if cachedUsersByEmail[email].UserTypeID == 2 {
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
		logInfo("MAIN", "Processing data ended")
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
		logInfo("MAIN", "Processing data ended")
		return
	}
	dateTo, err := time.Parse(layout, data.To)
	if err != nil {
		logError("MAIN", "Problem parsing date: "+data.To)
		var responseData DataPageOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Processing data ended")
		return
	}
	email, _, _ := request.BasicAuth()
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
			db.Where(workplaceIds).Where("date_time_start >= ?", dateFrom).Where("date_time_end <= ?", dateTo).Or("date_time_end is null").Find(&orderRecords)
			var userRecords []database.UserRecord
			db.Where(workplaceIds).Where("date_time_start >= ?", dateFrom).Where("date_time_end <= ?", dateTo).Or("date_time_end is null").Find(&userRecords)
			var userRecordsByRecordId = map[int]database.UserRecord{}
			for _, record := range userRecords {
				userRecordsByRecordId[record.OrderRecordID] = record
			}
			var orderRecordsOutput []OrderRecord
			for _, record := range orderRecords {
				var orderRecordOutput OrderRecord
				orderRecordOutput.WorkplaceName = cachedWorkplacesById[uint(record.WorkplaceID)].Name
				orderRecordOutput.WorkplaceModeName = cachedWorkplaceModesById[uint(record.WorkplaceModeID)].Name
				orderRecordOutput.WorkshiftName = cachedWorkshiftsById[uint(record.WorkshiftID)].Name
				actualUserId := userRecordsByRecordId[int(record.ID)].UserID
				orderRecordOutput.UserName = cachedUsersById[uint(actualUserId)].FirstName + " " + cachedUsersById[uint(actualUserId)].SecondName
				orderRecordOutput.OrderName = cachedOrdersById[uint(record.OrderID)].Name
				orderRecordOutput.OperationName = cachedOperationsById[uint(record.OperationID)].Name
				orderRecordOutput.AverageCycle = record.AverageCycle
				orderRecordOutput.AverageCycle = record.AverageCycle
				orderRecordOutput.Cavity = record.Cavity
				orderRecordOutput.Ok = record.CountOk
				orderRecordOutput.Nok = record.CountNok
				orderRecordOutput.Note = record.Note
				orderRecordsOutput = append(orderRecordsOutput, orderRecordOutput)
			}
			logInfo("MAIN", "Sending "+strconv.Itoa(len(orderRecordsOutput))+" orders")
			var data TableData

			// TODO: create locales for order table headers
			// TODO: create locales for data-table-rows-count-title	Show entries:	Title for rows steps box
			// TODO: create locales for data-table-search-title	Search:	Title for search input
			// TODO: create locales for data-table-info-title	Showing $1 to $2 of $3 entries	Title for table info block
			// TODO: create locales for data-pagination-prev-title	Prev	Title pagination prev button
			// TODO: create locales for data-pagination-next-title	Next	Title pagination next button
			// TODO: create locales for data-all-records-title	All	Title all records in rows steps block
			// TODO: create locales for data-inspector-title	Inspector	Title for table inspector window

			// TODO: creates headers
			var headerCell HeaderCell
			headerCell.HeaderName = "Workplace Name"
			data.TableHeader = append(data.TableHeader, headerCell)
			var headerCell2 HeaderCell
			headerCell2.HeaderName = "b"
			data.TableHeader = append(data.TableHeader, headerCell2)

			// TODO: loop data and add rows
			//for _, record := range orderRecords {
			//var tableRow1 TableRow
			//var tableRow2 TableRow
			//
			//var tableCellA1 TableCell
			//tableCellA1.CellName = "testA1"
			//
			//var tableCellA2 TableCell
			//tableCellA2.CellName = "testA2"
			//
			//var tableCellB1 TableCell
			//tableCellB1.CellName = "testB1"
			//
			//var tableCellB2 TableCell
			//tableCellB2.CellName = "testB2"
			//
			//tableRow1.TableCell = append(tableRow1.TableCell, tableCellA1)
			//tableRow1.TableCell = append(tableRow1.TableCell, tableCellA2)
			//
			//tableRow2.TableCell = append(tableRow2.TableCell, tableCellB1)
			//tableRow2.TableCell = append(tableRow2.TableCell, tableCellB2)
			//
			//data.TableRows = append(data.TableRows, tableRow1)
			//data.TableRows = append(data.TableRows, tableRow2)
			//}

			tmpl := template.Must(template.ParseFiles("./html/table.html"))
			_ = tmpl.Execute(writer, data)
			logInfo("MAIN", "Home page sent for "+email)
			return

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
