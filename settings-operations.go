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
	"time"
)

type OperationsSettingsDataOutput struct {
	DataTableSearchTitle    string
	DataTableInfoTitle      string
	DataTableRowsCountTitle string
	TableHeader             []HeaderCell
	TableRows               []TableRow
}

type OperationDetailsDataOutput struct {
	OperationName        string
	OperationNamePrepend string
	OrderName            string
	OrderNamePrepend     string
	Barcode              string
	BarcodePrepend       string
	Note                 string
	NotePrepend          string
	CreatedAt            string
	CreatedAtPrepend     string
	UpdatedAt            string
	UpdatedAtPrepend     string
	Orders               []OrderSelection
}

type OrderSelection struct {
	OrderName     string
	OrderId       uint
	OrderSelected string
}

type OperationDetailsDataInput struct {
	Id      string
	Name    string
	Order   string
	Barcode string
	Note    string
	Url     string
	Pdf     string
}

func loadOperationsSettings(writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS-OPERATIONS", "Loading operations settings")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-OPERATIONS", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-OPERATIONS", "Loading operations settings ended")
		return
	}
	var records []database.Operation
	db.Order("id desc").Find(&records)
	var data OperationsSettingsDataOutput
	data.DataTableSearchTitle = getLocale(email, "data-table-search-title")
	data.DataTableInfoTitle = getLocale(email, "data-table-info-title")
	data.DataTableRowsCountTitle = getLocale(email, "data-table-rows-count-title")
	addOperationSettingsTableHeaders(email, &data)
	for _, record := range records {
		addOperationSettingsTableRow(record, &data)
	}
	tmpl := template.Must(template.ParseFiles("./html/settings-table.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("SETTINGS-OPERATIONS", "Operations settings loaded in "+time.Since(timer).String())
}

func addOperationSettingsTableRow(record database.Operation, data *OperationsSettingsDataOutput) {
	var tableRow TableRow
	id := TableCell{CellName: strconv.Itoa(int(record.ID))}
	tableRow.TableCell = append(tableRow.TableCell, id)
	name := TableCell{CellName: record.Name}
	tableRow.TableCell = append(tableRow.TableCell, name)
	data.TableRows = append(data.TableRows, tableRow)
}

func addOperationSettingsTableHeaders(email string, data *OperationsSettingsDataOutput) {
	id := HeaderCell{HeaderName: "#", HeaderWidth: "30"}
	data.TableHeader = append(data.TableHeader, id)
	name := HeaderCell{HeaderName: getLocale(email, "operation-name")}
	data.TableHeader = append(data.TableHeader, name)
}

func loadOperationDetails(id string, writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS-OPERATIONS", "Loading operation details")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-OPERATIONS", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-OPERATIONS", "Loading operation details ended")
		return
	}
	var operation database.Operation
	db.Where("id = ?", id).Find(&operation)
	var orders []OrderSelection
	for _, order := range cachedOrdersById {
		if order.Name == cachedWorkplacesById[uint(operation.OrderID)].Name {
			orders = append(orders, OrderSelection{OrderName: order.Name, OrderId: order.ID, OrderSelected: "selected"})
		} else {
			orders = append(orders, OrderSelection{OrderName: order.Name, OrderId: order.ID})
		}
	}
	sort.Slice(orders, func(i, j int) bool {
		return orders[i].OrderName < orders[j].OrderName
	})

	data := OperationDetailsDataOutput{
		OperationName:        operation.Name,
		OperationNamePrepend: getLocale(email, "operation-name"),
		OrderName:            cachedOrdersById[uint(operation.OrderID)].Name,
		OrderNamePrepend:     getLocale(email, "order-name"),
		Barcode:              operation.Barcode,
		BarcodePrepend:       getLocale(email, "barcode"),
		Note:                 operation.Note,
		NotePrepend:          getLocale(email, "note-name"),
		CreatedAt:            operation.CreatedAt.Format("2006-01-02T15:04:05"),
		CreatedAtPrepend:     getLocale(email, "created-at"),
		UpdatedAt:            operation.UpdatedAt.Format("2006-01-02T15:04:05"),
		UpdatedAtPrepend:     getLocale(email, "updated-at"),
		Orders:               orders,
	}
	tmpl := template.Must(template.ParseFiles("./html/settings-detail-operation.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("SETTINGS-OPERATIONS", "Operation details loaded in "+time.Since(timer).String())
}

func saveOperation(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	timer := time.Now()
	logInfo("SETTINGS-OPERATIONS", "Saving operation started")
	var data OperationDetailsDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS-OPERATIONS", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-OPERATIONS", "Saving operation ended")
		return
	}
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-OPERATIONS", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-OPERATIONS", "Saving operation ended")
		return
	}
	var operation database.Operation
	db.Where("id=?", data.Id).Find(&operation)
	operation.Name = data.Name
	operation.OrderID = int(cachedOrdersByName[data.Order].ID)
	operation.Barcode = data.Barcode
	operation.Note = data.Note
	db.Save(&operation)
	cacheAlarms(db)
	logInfo("SETTINGS-OPERATIONS", "Operation saved in "+time.Since(timer).String())
}
