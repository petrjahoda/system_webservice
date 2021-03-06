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

type PackagesSettingsDataOutput struct {
	DataTableSearchTitle    string
	DataTableInfoTitle      string
	DataTableRowsCountTitle string
	TableHeader             []HeaderCell
	TableRows               []TableRow
	TableHeaderType         []HeaderCellType
	TableRowsType           []TableRowType
	Result                  string
}

type PackageTypeDetailsDataOutput struct {
	PackageTypeName        string
	PackageTypeNamePrepend string
	Count                  string
	CountPrepend           string
	Note                   string
	NotePrepend            string
	CreatedAt              string
	CreatedAtPrepend       string
	UpdatedAt              string
	UpdatedAtPrepend       string
	Result                 string
}

type PackageDetailsDataOutput struct {
	PackageName            string
	PackageNamePrepend     string
	PackageTypeName        string
	PackageTypeNamePrepend string
	OrderNamePrepend       string
	Barcode                string
	BarcodePrepend         string
	Note                   string
	NotePrepend            string
	CreatedAt              string
	CreatedAtPrepend       string
	UpdatedAt              string
	UpdatedAtPrepend       string
	PackageTypes           []PackageTypeSelection
	Orders                 []OrderSelection
	Result                 string
}

type PackageTypeSelection struct {
	PackageTypeName     string
	PackageTypeId       uint
	PackageTypeSelected string
}

type PackageDetailsDataInput struct {
	Id      string
	Name    string
	Type    string
	Order   string
	Barcode string
	Note    string
}

type PackageTypeDetailsDataInput struct {
	Id    string
	Name  string
	Count string
	Note  string
}

func loadPackages(writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS", "Loading packages")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS", "Problem opening database: "+err.Error())
		var responseData PackagesSettingsDataOutput
		responseData.Result = "ERR: Problem opening database, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Loading packages ended with error")
		return
	}
	var data PackagesSettingsDataOutput
	data.DataTableSearchTitle = getLocale(email, "data-table-search-title")
	data.DataTableInfoTitle = getLocale(email, "data-table-info-title")
	data.DataTableRowsCountTitle = getLocale(email, "data-table-rows-count-title")
	var records []database.Package
	db.Order("id desc").Find(&records)
	addPackagesTableHeaders(email, &data)
	for _, record := range records {
		addPackagesTableRow(record, &data)
	}
	var typeRecords []database.PackageType
	db.Order("id desc").Find(&typeRecords)
	addPackageTypesTableHeaders(email, &data)
	for _, record := range typeRecords {
		addPackageTypesTableRow(record, &data)
	}
	tmpl, err := template.ParseFiles("./html/settings-table-type.html")
	if err != nil {
		logError("SETTINGS", "Problem parsing html file: "+err.Error())
		var responseData FaultSettingsDataOutput
		responseData.Result = "ERR: Problem parsing html file: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
	} else {
		data.Result = "INF: Packages processed in " + time.Since(timer).String()
		_ = tmpl.Execute(writer, data)
		logInfo("SETTINGS", "Packages loaded in "+time.Since(timer).String())
	}
}

func addPackageTypesTableRow(record database.PackageType, data *PackagesSettingsDataOutput) {
	var tableRow TableRowType
	id := TableCellType{CellNameType: strconv.Itoa(int(record.ID))}
	tableRow.TableCellType = append(tableRow.TableCellType, id)
	name := TableCellType{CellNameType: record.Name}
	tableRow.TableCellType = append(tableRow.TableCellType, name)
	data.TableRowsType = append(data.TableRowsType, tableRow)
}

func addPackageTypesTableHeaders(email string, data *PackagesSettingsDataOutput) {
	id := HeaderCellType{HeaderNameType: "#", HeaderWidthType: "30"}
	data.TableHeaderType = append(data.TableHeaderType, id)
	name := HeaderCellType{HeaderNameType: getLocale(email, "type-name")}
	data.TableHeaderType = append(data.TableHeaderType, name)
}

func addPackagesTableRow(record database.Package, data *PackagesSettingsDataOutput) {
	var tableRow TableRow
	id := TableCell{CellName: strconv.Itoa(int(record.ID))}
	tableRow.TableCell = append(tableRow.TableCell, id)
	name := TableCell{CellName: record.Name}
	tableRow.TableCell = append(tableRow.TableCell, name)
	data.TableRows = append(data.TableRows, tableRow)
}

func addPackagesTableHeaders(email string, data *PackagesSettingsDataOutput) {
	id := HeaderCell{HeaderName: "#", HeaderWidth: "30"}
	data.TableHeader = append(data.TableHeader, id)
	name := HeaderCell{HeaderName: getLocale(email, "package-name")}
	data.TableHeader = append(data.TableHeader, name)
}

func loadPackage(id string, writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS", "Loading package")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS", "Problem opening database: "+err.Error())
		var responseData PackageDetailsDataOutput
		responseData.Result = "ERR: Problem opening database, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Loading package ended with error")
		return
	}
	var onePackage database.Package
	db.Where("id = ?", id).Find(&onePackage)
	var packageTypes []PackageTypeSelection
	packageTypesByIdSync.RLock()
	packageTypesById := cachedPackageTypesById
	packageTypesByIdSync.RUnlock()
	for _, packageType := range packageTypesById {
		if packageType.Name == packageTypesById[uint(onePackage.PackageTypeID)].Name {
			packageTypes = append(packageTypes, PackageTypeSelection{PackageTypeName: packageType.Name, PackageTypeId: packageType.ID, PackageTypeSelected: "selected"})
		} else {
			packageTypes = append(packageTypes, PackageTypeSelection{PackageTypeName: packageType.Name, PackageTypeId: packageType.ID})
		}
	}
	sort.Slice(packageTypes, func(i, j int) bool {
		return packageTypes[i].PackageTypeName < packageTypes[j].PackageTypeName
	})
	var orders []database.Order
	db.Find(&orders)
	var orderSelection []OrderSelection
	ordersByIdSync.RLock()
	ordersById := cachedOrdersById
	ordersByIdSync.RUnlock()
	for _, order := range orders {
		if order.Name == ordersById[uint(onePackage.OrderID)].Name {
			orderSelection = append(orderSelection, OrderSelection{OrderName: order.Name, OrderId: order.ID, OrderSelected: "selected"})
		} else {
			orderSelection = append(orderSelection, OrderSelection{OrderName: order.Name, OrderId: order.ID})
		}
	}
	sort.Slice(orderSelection, func(i, j int) bool {
		return orderSelection[i].OrderName < orderSelection[j].OrderName
	})
	data := PackageDetailsDataOutput{
		PackageName:        onePackage.Name,
		PackageNamePrepend: getLocale(email, "package-name"),

		PackageTypeNamePrepend: getLocale(email, "type-name"),
		OrderNamePrepend:       getLocale(email, "order-name"),
		Barcode:                onePackage.Barcode,
		BarcodePrepend:         getLocale(email, "barcode"),
		Note:                   onePackage.Note,
		NotePrepend:            getLocale(email, "note-name"),
		CreatedAt:              onePackage.CreatedAt.Format("2006-01-02T15:04:05"),
		CreatedAtPrepend:       getLocale(email, "created-at"),
		UpdatedAt:              onePackage.UpdatedAt.Format("2006-01-02T15:04:05"),
		UpdatedAtPrepend:       getLocale(email, "updated-at"),
		PackageTypes:           packageTypes,
		Orders:                 orderSelection,
	}
	packageTypesByIdSync.RLock()
	data.PackageTypeName = cachedPackageTypesById[uint(onePackage.PackageTypeID)].Name
	packageTypesByIdSync.RUnlock()
	tmpl, err := template.ParseFiles("./html/settings-detail-package.html")
	if err != nil {
		logError("SETTINGS", "Problem parsing html file: "+err.Error())
		var responseData PackageDetailsDataOutput
		responseData.Result = "ERR: Problem parsing html file: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
	} else {
		data.Result = "INF: Package detail processed in " + time.Since(timer).String()
		_ = tmpl.Execute(writer, data)
		logInfo("SETTINGS", "Package detail loaded in "+time.Since(timer).String())
	}
}

func loadPackageType(id string, writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS", "Loading package type")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS", "Problem opening database: "+err.Error())
		var responseData PackageTypeDetailsDataOutput
		responseData.Result = "ERR: Problem opening database, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Loading package type ended with error")
		return
	}
	var packageType database.PackageType
	db.Where("id = ?", id).Find(&packageType)
	data := PackageTypeDetailsDataOutput{
		PackageTypeName:        packageType.Name,
		PackageTypeNamePrepend: getLocale(email, "type-name"),
		Count:                  strconv.Itoa(packageType.Count),
		CountPrepend:           getLocale(email, "count-requested"),
		Note:                   packageType.Note,
		NotePrepend:            getLocale(email, "note-name"),
		CreatedAt:              packageType.CreatedAt.Format("2006-01-02T15:04:05"),
		CreatedAtPrepend:       getLocale(email, "created-at"),
		UpdatedAt:              packageType.UpdatedAt.Format("2006-01-02T15:04:05"),
		UpdatedAtPrepend:       getLocale(email, "updated-at"),
	}
	tmpl, err := template.ParseFiles("./html/settings-detail-package-type.html")
	if err != nil {
		logError("SETTINGS", "Problem parsing html file: "+err.Error())
		var responseData PackageTypeDetailsDataOutput
		responseData.Result = "ERR: Problem parsing html file: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
	} else {
		data.Result = "INF: Package type detail processed in " + time.Since(timer).String()
		_ = tmpl.Execute(writer, data)
		logInfo("SETTINGS", "Package type detail loaded in "+time.Since(timer).String())
	}
}

func savePackage(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	timer := time.Now()
	logInfo("SETTINGS", "Saving package")
	var data PackageDetailsDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "ERR: Error parsing data, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Saving package ended with error")
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
		logInfo("SETTINGS", "Saving package ended with error")
		return
	}
	var onePackage database.Package
	db.Where("id=?", data.Id).Find(&onePackage)
	onePackage.Name = data.Name
	packageTypesByNameSync.RLock()
	onePackage.PackageTypeID = int(cachedPackageTypesByName[data.Type].ID)
	packageTypesByNameSync.RUnlock()
	ordersByNameSync.RLock()
	onePackage.OrderID = int(cachedOrdersByName[data.Order].ID)
	ordersByNameSync.RUnlock()
	onePackage.Barcode = data.Barcode
	onePackage.Note = data.Note
	result := db.Save(&onePackage)
	cachePackages(db)
	if result.Error != nil {
		var responseData TableOutput
		responseData.Result = "ERR: Package not saved: " + result.Error.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logError("SETTINGS", "Package "+onePackage.Name+" not saved: "+result.Error.Error())
	} else {
		var responseData TableOutput
		responseData.Result = "INF: Package saved in " + time.Since(timer).String()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Package "+onePackage.Name+" saved in "+time.Since(timer).String())
	}
}

func savePackageType(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	timer := time.Now()
	logInfo("SETTINGS", "Saving package type")
	var data PackageTypeDetailsDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "ERR: Error parsing data, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Saving package type ended with error")
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
		logInfo("SETTINGS", "Saving package type ended with error")
		return
	}
	var packageType database.PackageType
	db.Where("id=?", data.Id).Find(&packageType)
	packageType.Name = data.Name
	packageType.Count, _ = strconv.Atoi(data.Count)
	packageType.Note = data.Note
	result := db.Save(&packageType)
	cachePackages(db)
	if result.Error != nil {
		var responseData TableOutput
		responseData.Result = "ERR: Package type not saved: " + result.Error.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logError("SETTINGS", "Package type "+packageType.Name+" not saved: "+result.Error.Error())
	} else {
		var responseData TableOutput
		responseData.Result = "INF: Package type saved in " + time.Since(timer).String()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Package type "+packageType.Name+" saved in "+time.Since(timer).String())
	}
}
