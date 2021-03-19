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
}

type ProductDetailsDataInput struct {
	Id               string
	Name             string
	Barcode          string
	Cycle            string
	DowntimeDuration string
	Note             string
}

func saveProduct(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	timer := time.Now()
	logInfo("SETTINGS-PRODUCTS", "Saving product started")
	var data ProductDetailsDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS-PRODUCTS", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-PRODUCTS", "Saving product ended")
		return
	}
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-PRODUCTS", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-PRODUCTS", "Saving product ended")
		return
	}

	cycleTime, err := time.ParseDuration(data.Cycle)
	if err != nil {
		logError("SETTINGS-PRODUCTS", "Problem parsing cycle time: "+err.Error())
		cycleTime = 0
	}
	duration, err := time.ParseDuration(data.DowntimeDuration)
	if err != nil {
		logError("SETTINGS-PRODUCTS", "Problem parsing duration: "+err.Error())
		duration = 0
	}
	var product database.Product
	db.Where("id=?", data.Id).Find(&product)
	product.Name = data.Name
	product.Barcode = data.Barcode
	product.CycleTime = int(cycleTime.Seconds())
	product.DownTimeDuration = duration
	product.Note = data.Note
	db.Debug().Save(&product)
	cacheProducts(db)
	logInfo("SETTINGS-PRODUCTS", "Product saved in "+time.Since(timer).String())
}

func loadProductDetails(id string, writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS-PRODUCTS", "Loading product details")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-PRODUCTS", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-PRODUCTS", "Loading product details ended")
		return
	}
	productId, _ := strconv.Atoi(id)
	product := cachedProductsById[uint(productId)]
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
	tmpl := template.Must(template.ParseFiles("./html/settings-detail-product.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("SETTINGS-PRODUCTS", "Product details loaded in "+time.Since(timer).String())
}

func loadProductsSettings(writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS-PRODUCTS", "Loading products settings")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-PRODUCTS", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-PRODUCTS", "Loading products settings ended")
		return
	}
	var records []database.Product
	db.Order("id desc").Find(&records)
	var data ProductsSettingsDataOutput
	data.DataTableSearchTitle = getLocale(email, "data-table-search-title")
	data.DataTableInfoTitle = getLocale(email, "data-table-info-title")
	data.DataTableRowsCountTitle = getLocale(email, "data-table-rows-count-title")
	addProductSettingsTableHeaders(email, &data)
	for _, record := range records {
		addProductSettingsTableRow(record, &data)
	}
	tmpl := template.Must(template.ParseFiles("./html/settings-table.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("SETTINGS-PRODUCTS", "Products settings loaded in "+time.Since(timer).String())
}

func addProductSettingsTableRow(record database.Product, data *ProductsSettingsDataOutput) {
	var tableRow TableRow
	id := TableCell{CellName: strconv.Itoa(int(record.ID))}
	tableRow.TableCell = append(tableRow.TableCell, id)
	name := TableCell{CellName: record.Name}
	tableRow.TableCell = append(tableRow.TableCell, name)
	data.TableRows = append(data.TableRows, tableRow)
}

func addProductSettingsTableHeaders(email string, data *ProductsSettingsDataOutput) {
	id := HeaderCell{HeaderName: "#", HeaderWidth: "30"}
	data.TableHeader = append(data.TableHeader, id)
	name := HeaderCell{HeaderName: getLocale(email, "product-name")}
	data.TableHeader = append(data.TableHeader, name)
}
