package main

import (
	"database/sql"
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

type OrdersSettingsDataOutput struct {
	DataTableSearchTitle    string
	DataTableInfoTitle      string
	DataTableRowsCountTitle string
	TableHeader             []HeaderCell
	TableRows               []TableRow
	Result                  string
}

type OrderDetailsDataOutput struct {
	OrderName              string
	OrderNamePrepend       string
	ProductName            string
	ProductNamePrepend     string
	WorkplaceName          string
	WorkplaceNamePrepend   string
	Barcode                string
	BarcodePrepend         string
	DateTimeRequest        string
	DateTimeRequestPrepend string
	CountRequest           string
	CountRequestPrepend    string
	Cavity                 string
	CavityPrepend          string
	Note                   string
	NotePrepend            string
	CreatedAt              string
	CreatedAtPrepend       string
	UpdatedAt              string
	UpdatedAtPrepend       string
	Products               []ProductSelection
	Workplaces             []WorkplaceSelection
	Result                 string
}

type ProductSelection struct {
	ProductName     string
	ProductId       uint
	ProductSelected string
}

type OrderDetailsDataInput struct {
	Id              string
	Name            string
	Product         string
	Workplace       string
	Barcode         string
	CountRequest    string
	DateTimeRequest string
	Cavity          string
	Note            string
}

func loadOrders(writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS", "Loading orders")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS", "Problem opening database: "+err.Error())
		var responseData OrdersSettingsDataOutput
		responseData.Result = "ERR: Problem opening database, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Loading orders ended with error")
		return
	}
	var records []database.Order
	db.Order("id desc").Find(&records)
	var data OrdersSettingsDataOutput
	data.DataTableSearchTitle = getLocale(email, "data-table-search-title")
	data.DataTableInfoTitle = getLocale(email, "data-table-info-title")
	data.DataTableRowsCountTitle = getLocale(email, "data-table-rows-count-title")
	addOrdersTableHeaders(email, &data)
	for _, record := range records {
		addOrdersTableRow(record, &data)
	}
	tmpl, err := template.ParseFiles("./html/settings-table.html")
	if err != nil {
		logError("SETTINGS", "Problem parsing html file: "+err.Error())
		var responseData OrdersSettingsDataOutput
		responseData.Result = "ERR: Problem parsing html file: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
	} else {
		data.Result = "INF: Orders processed in " + time.Since(timer).String()
		_ = tmpl.Execute(writer, data)
		logInfo("SETTINGS", "Orders loaded in "+time.Since(timer).String())
	}
}

func addOrdersTableRow(record database.Order, data *OrdersSettingsDataOutput) {
	var tableRow TableRow
	id := TableCell{CellName: strconv.Itoa(int(record.ID))}
	tableRow.TableCell = append(tableRow.TableCell, id)
	name := TableCell{CellName: record.Name}
	tableRow.TableCell = append(tableRow.TableCell, name)
	data.TableRows = append(data.TableRows, tableRow)
}

func addOrdersTableHeaders(email string, data *OrdersSettingsDataOutput) {
	id := HeaderCell{HeaderName: "#", HeaderWidth: "30"}
	data.TableHeader = append(data.TableHeader, id)
	name := HeaderCell{HeaderName: getLocale(email, "order-name")}
	data.TableHeader = append(data.TableHeader, name)
}

func loadOrder(id string, writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS", "Loading order")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS", "Problem opening database: "+err.Error())
		var responseData OrderDetailsDataOutput
		responseData.Result = "ERR: Problem opening database, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Loading order ended with error")
		return
	}
	var order database.Order
	db.Where("id = ?", id).Find(&order)
	var workplaces []WorkplaceSelection
	workplacesByIdSync.RLock()
	workplacesById := cachedWorkplacesById
	workplacesByIdSync.RUnlock()
	for _, workplace := range workplacesById {
		if workplace.Name == workplacesById[uint(order.WorkplaceID.Int32)].Name {
			workplaces = append(workplaces, WorkplaceSelection{WorkplaceName: workplace.Name, WorkplaceId: workplace.ID, WorkplaceSelected: "selected"})
		} else {
			workplaces = append(workplaces, WorkplaceSelection{WorkplaceName: workplace.Name, WorkplaceId: workplace.ID})
		}
	}
	sort.Slice(workplaces, func(i, j int) bool {
		return workplaces[i].WorkplaceName < workplaces[j].WorkplaceName
	})
	var products []ProductSelection
	productsByIdSync.RLock()
	productsById := cachedProductsById
	productsByIdSync.RUnlock()
	for _, product := range productsById {
		if product.Name == productsById[uint(order.ProductID.Int32)].Name {
			products = append(products, ProductSelection{ProductName: product.Name, ProductId: product.ID, ProductSelected: "selected"})
		} else {
			products = append(products, ProductSelection{ProductName: product.Name, ProductId: product.ID})
		}
	}
	sort.Slice(products, func(i, j int) bool {
		return products[i].ProductName < products[j].ProductName
	})
	requiredDate := order.CreatedAt.Format("2006-01-02T15:04:05")
	if !order.DateTimeRequest.Time.IsZero() {
		requiredDate = order.DateTimeRequest.Time.Format("2006-01-02T15:04:05")
	}
	data := OrderDetailsDataOutput{
		OrderName:              order.Name,
		OrderNamePrepend:       getLocale(email, "order-name"),
		ProductNamePrepend:     getLocale(email, "product-name"),
		WorkplaceNamePrepend:   getLocale(email, "workplace-name"),
		DateTimeRequest:        requiredDate,
		DateTimeRequestPrepend: getLocale(email, "date-requested"),
		CountRequest:           strconv.Itoa(order.CountRequest),
		CountRequestPrepend:    getLocale(email, "count-requested"),
		Cavity:                 strconv.Itoa(order.Cavity),
		CavityPrepend:          getLocale(email, "cavity-name"),
		Barcode:                order.Barcode,
		BarcodePrepend:         getLocale(email, "barcode"),
		Note:                   order.Note,
		NotePrepend:            getLocale(email, "note-name"),
		CreatedAt:              order.CreatedAt.Format("2006-01-02T15:04:05"),
		CreatedAtPrepend:       getLocale(email, "created-at"),
		UpdatedAt:              order.UpdatedAt.Format("2006-01-02T15:04:05"),
		UpdatedAtPrepend:       getLocale(email, "updated-at"),
		Products:               products,
		Workplaces:             workplaces,
	}
	workplacesByIdSync.RLock()
	data.WorkplaceName = cachedWorkplacesById[uint(order.WorkplaceID.Int32)].Name
	workplacesByIdSync.RUnlock()
	productsByIdSync.RLock()
	data.ProductName = cachedProductsById[uint(order.ProductID.Int32)].Name
	productsByIdSync.RUnlock()
	tmpl, err := template.ParseFiles("./html/settings-detail-order.html")
	if err != nil {
		logError("SETTINGS", "Problem parsing html file: "+err.Error())
		var responseData OrderDetailsDataOutput
		responseData.Result = "ERR: Problem parsing html file: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
	} else {
		data.Result = "INF: Order detail processed in " + time.Since(timer).String()
		_ = tmpl.Execute(writer, data)
		logInfo("SETTINGS", "Order detail loaded in "+time.Since(timer).String())
	}
}

func saveOrder(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	timer := time.Now()
	logInfo("SETTINGS", "Saving order")
	var data OrderDetailsDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "ERR: Error parsing data, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Saving order ended with error")
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
		logInfo("SETTINGS", "Saving order ended with error")
		return
	}
	var order database.Order
	db.Where("id=?", data.Id).Find(&order)
	order.Name = data.Name
	if len(data.Product) > 0 {
		productsByNameSync.RLock()
		order.ProductID = sql.NullInt32{Int32: int32(cachedProductsByName[data.Product].ID), Valid: true}
		productsByNameSync.RUnlock()
	}
	if len(data.Workplace) > 0 {
		workplacesByNameSync.RLock()
		order.WorkplaceID = sql.NullInt32{Int32: int32(cachedWorkplacesByName[data.Workplace].ID), Valid: true}
		workplacesByNameSync.RUnlock()
	}
	timeParsed, _ := time.Parse("2006-01-02T15:04:05", data.DateTimeRequest)
	order.DateTimeRequest = sql.NullTime{
		Time:  timeParsed,
		Valid: true,
	}
	order.CountRequest, _ = strconv.Atoi(data.CountRequest)
	order.Cavity, _ = strconv.Atoi(data.Cavity)
	order.Note = data.Note
	result := db.Save(&order)
	cacheOrders(db)
	if result.Error != nil {
		var responseData TableOutput
		responseData.Result = "ERR: Order not saved: " + result.Error.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logError("SETTINGS", "Order "+order.Name+" not saved: "+result.Error.Error())
	} else {
		var responseData TableOutput
		responseData.Result = "INF: Order saved in " + time.Since(timer).String()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Order "+order.Name+" saved in "+time.Since(timer).String())
	}
}
