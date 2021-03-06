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
	Result                  string
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
	Result               string
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

func loadOperations(writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS", "Loading operations")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS", "Problem opening database: "+err.Error())
		var responseData OperationsSettingsDataOutput
		responseData.Result = "ERR: Problem opening database, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Loading operations ended with error")
		return
	}
	var records []database.Operation
	db.Order("id desc").Find(&records)
	var data OperationsSettingsDataOutput
	data.DataTableSearchTitle = getLocale(email, "data-table-search-title")
	data.DataTableInfoTitle = getLocale(email, "data-table-info-title")
	data.DataTableRowsCountTitle = getLocale(email, "data-table-rows-count-title")
	addOperationsTableHeaders(email, &data)
	for _, record := range records {
		addOperationsTableRow(record, &data)
	}
	tmpl, err := template.ParseFiles("./html/settings-table.html")
	if err != nil {
		logError("SETTINGS", "Problem Parsing html file: "+err.Error())
		var responseData AlarmsSettingsDataOutput
		responseData.Result = "ERR: Problem parsing html file: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
	} else {
		data.Result = "INF: Operations processed in " + time.Since(timer).String()
		_ = tmpl.Execute(writer, data)
		logInfo("SETTINGS", "Operations loaded in "+time.Since(timer).String())
	}
}

func addOperationsTableRow(record database.Operation, data *OperationsSettingsDataOutput) {
	var tableRow TableRow
	id := TableCell{CellName: strconv.Itoa(int(record.ID))}
	tableRow.TableCell = append(tableRow.TableCell, id)
	name := TableCell{CellName: record.Name}
	tableRow.TableCell = append(tableRow.TableCell, name)
	data.TableRows = append(data.TableRows, tableRow)
}

func addOperationsTableHeaders(email string, data *OperationsSettingsDataOutput) {
	id := HeaderCell{HeaderName: "#", HeaderWidth: "30"}
	data.TableHeader = append(data.TableHeader, id)
	name := HeaderCell{HeaderName: getLocale(email, "operation-name")}
	data.TableHeader = append(data.TableHeader, name)
}

func loadOperation(id string, writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS", "Loading operation")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS", "Problem opening database: "+err.Error())
		var responseData OperationDetailsDataOutput
		responseData.Result = "ERR: Problem opening database, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Loading operation ended with error")
		return
	}
	var operation database.Operation
	db.Where("id = ?", id).Find(&operation)
	var orders []OrderSelection
	ordersByIdSync.RLock()
	ordersById := cachedOrdersById
	ordersByIdSync.RUnlock()
	for _, order := range ordersById {
		workplacesByIdSync.RLock()
		cachedOrderName := cachedWorkplacesById[uint(operation.OrderID)].Name
		workplacesByIdSync.RUnlock()
		if order.Name == cachedOrderName {
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
	ordersByIdSync.RLock()
	data.OrderName = cachedOrdersById[uint(operation.OrderID)].Name
	ordersByIdSync.RUnlock()
	tmpl, err := template.ParseFiles("./html/settings-detail-operation.html")
	if err != nil {
		logError("SETTINGS", "Problem parsing html file: "+err.Error())
		var responseData OperationDetailsDataOutput
		responseData.Result = "ERR: Problem parsing html file: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
	} else {
		data.Result = "INF: Operation detail processed in " + time.Since(timer).String()
		_ = tmpl.Execute(writer, data)
		logInfo("SETTINGS", "Operation detail loaded in "+time.Since(timer).String())
	}
}

func saveOperation(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	timer := time.Now()
	logInfo("SETTINGS", "Saving operation")
	var data OperationDetailsDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "ERR: Error parsing data, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Saving operation ended with error")
		return
	}
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "ERR: Problem opening database, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Saving operation ended with error")
		return
	}
	var operation database.Operation
	db.Where("id=?", data.Id).Find(&operation)
	operation.Name = data.Name
	ordersByNameSync.RLock()
	operation.OrderID = int(cachedOrdersByName[data.Order].ID)
	ordersByNameSync.RUnlock()
	operation.Barcode = data.Barcode
	operation.Note = data.Note
	result := db.Save(&operation)
	cacheOperations(db)
	if result.Error != nil {
		var responseData TableOutput
		responseData.Result = "ERR: Operation not saved: " + result.Error.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logError("SETTINGS", "Operation "+operation.Name+" not saved: "+result.Error.Error())
	} else {
		var responseData TableOutput
		responseData.Result = "INF: Operation saved in " + time.Since(timer).String()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Operation "+operation.Name+" saved in "+time.Since(timer).String())
	}
}
