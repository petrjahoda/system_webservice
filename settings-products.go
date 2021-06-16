package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/petrjahoda/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

type ProductsSettingsDataOutput struct {
	DataTableSearchTitle    string
	DataTableInfoTitle      string
	DataTableRowsCountTitle string
	TableHeader             []HeaderCell
	TableRows               []TableRow
	Result                  string
}

type ProductDetailsDataOutput struct {
	ProductName             string
	ProductNamePrepend      string
	Barcode                 string
	BarcodePrepend          string
	CycleTime               string
	CycleTimePrepend        string
	DowntimeDuration        string
	DowntimeDurationPrepend string
	Note                    string
	NotePrepend             string
	CreatedAt               string
	CreatedAtPrepend        string
	UpdatedAt               string
	UpdatedAtPrepend        string
	Result                  string
}

type ProductDetailsDataInput struct {
	Id               string
	Name             string
	Barcode          string
	Cycle            string
	DowntimeDuration string
	Note             string
}

func loadProducts(writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS", "Loading products")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS", "Problem opening database: "+err.Error())
		var responseData ProductsSettingsDataOutput
		responseData.Result = "ERR: Problem opening database, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Loading products ended with error")
		return
	}
	var records []database.Product
	db.Order("id desc").Find(&records)
	var data ProductsSettingsDataOutput
	data.DataTableSearchTitle = getLocale(email, "data-table-search-title")
	data.DataTableInfoTitle = getLocale(email, "data-table-info-title")
	data.DataTableRowsCountTitle = getLocale(email, "data-table-rows-count-title")
	addProductsTableHeaders(email, &data)
	for _, record := range records {
		addProductsTableRow(record, &data)
	}
	tmpl, err := template.ParseFiles("./html/settings-table.html")
	if err != nil {
		logError("SETTINGS", "Problem parsing html file: "+err.Error())
		var responseData OrdersSettingsDataOutput
		responseData.Result = "ERR: Problem parsing html file: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
	} else {
		data.Result = "INF: Products processed in " + time.Since(timer).String()
		_ = tmpl.Execute(writer, data)
		logInfo("SETTINGS", "Products loaded in "+time.Since(timer).String())
	}
}

func addProductsTableRow(record database.Product, data *ProductsSettingsDataOutput) {
	var tableRow TableRow
	id := TableCell{CellName: strconv.Itoa(int(record.ID))}
	tableRow.TableCell = append(tableRow.TableCell, id)
	name := TableCell{CellName: record.Name}
	tableRow.TableCell = append(tableRow.TableCell, name)
	data.TableRows = append(data.TableRows, tableRow)
}

func addProductsTableHeaders(email string, data *ProductsSettingsDataOutput) {
	id := HeaderCell{HeaderName: "#", HeaderWidth: "30"}
	data.TableHeader = append(data.TableHeader, id)
	name := HeaderCell{HeaderName: getLocale(email, "product-name")}
	data.TableHeader = append(data.TableHeader, name)
}

func loadProduct(id string, writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS", "Loading product")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS", "Problem opening database: "+err.Error())
		var responseData ProductDetailsDataOutput
		responseData.Result = "ERR: Problem opening database, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Loading product ended with error")
		return
	}
	productId, _ := strconv.Atoi(id)
	productsByIdSync.RLock()
	product := cachedProductsById[uint(productId)]
	productsByIdSync.RUnlock()
	data := ProductDetailsDataOutput{
		ProductName:             product.Name,
		ProductNamePrepend:      getLocale(email, "product-name"),
		Barcode:                 product.Barcode,
		BarcodePrepend:          getLocale(email, "barcode"),
		CycleTime:               (time.Duration(product.CycleTime) * time.Second).String(),
		CycleTimePrepend:        getLocale(email, "cycle-name"),
		DowntimeDuration:        product.DownTimeDuration.String(),
		DowntimeDurationPrepend: getLocale(email, "downtime-duration"),
		Note:                    product.Note,
		NotePrepend:             getLocale(email, "note-name"),
		CreatedAt:               product.CreatedAt.Format("2006-01-02T15:04:05"),
		CreatedAtPrepend:        getLocale(email, "created-at"),
		UpdatedAt:               product.UpdatedAt.Format("2006-01-02T15:04:05"),
		UpdatedAtPrepend:        getLocale(email, "updated-at"),
	}
	tmpl, err := template.ParseFiles("./html/settings-detail-product.html")
	if err != nil {
		logError("SETTINGS", "Problem parsing html file: "+err.Error())
		var responseData ProductDetailsDataOutput
		responseData.Result = "ERR: Problem parsing html file: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
	} else {
		data.Result = "INF: Product detail processed in " + time.Since(timer).String()
		_ = tmpl.Execute(writer, data)
		logInfo("SETTINGS", "Product detail loaded in "+time.Since(timer).String())
	}
}

func saveProduct(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	timer := time.Now()
	logInfo("SETTINGS", "Saving product")
	var data ProductDetailsDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "ERR: Error parsing data, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Saving product ended with error")
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
		logInfo("SETTINGS", "Saving product ended with error")
		return
	}

	cycleTime, err := time.ParseDuration(data.Cycle)
	if err != nil {
		logError("SETTINGS", "Problem parsing cycle time: "+err.Error())
		cycleTime = 0
	}
	duration, err := time.ParseDuration(data.DowntimeDuration)
	if err != nil {
		logError("SETTINGS", "Problem parsing duration: "+err.Error())
		duration = 0
	}
	var product database.Product
	db.Where("id=?", data.Id).Find(&product)
	product.Name = data.Name
	product.Barcode = data.Barcode
	product.CycleTime = int(cycleTime.Seconds())
	product.DownTimeDuration = duration
	product.Note = data.Note
	result := db.Save(&product)
	cacheProducts(db)
	if result.Error != nil {
		var responseData TableOutput
		responseData.Result = "ERR: Product not saved: " + result.Error.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logError("SETTINGS", "Product "+product.Name+" not saved: "+result.Error.Error())
	} else {
		var responseData TableOutput
		responseData.Result = "INF: Product saved in " + time.Since(timer).String()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Product "+product.Name+" saved in "+time.Since(timer).String())
	}
}
