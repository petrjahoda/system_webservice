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

func savePackageType(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	timer := time.Now()
	logInfo("SETTINGS-PACKAGES", "Saving package type started")
	var data PackageTypeDetailsDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS-PACKAGES", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-PACKAGES", "Saving package type ended")
		return
	}
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-PACKAGES", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-PACKAGES", "Saving package type ended")
		return
	}
	var packageType database.PackageType
	db.Where("id=?", data.Id).Find(&packageType)
	packageType.Name = data.Name
	packageType.Count, _ = strconv.Atoi(data.Count)
	packageType.Note = data.Note
	db.Save(&packageType)
	cachePackages(db)
	logInfo("SETTINGS-PACKAGES", "Package type saved in "+time.Since(timer).String())
}

func savePackage(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	timer := time.Now()
	logInfo("SETTINGS-PACKAGES", "Saving package started")
	var data PackageDetailsDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS-PACKAGES", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-PACKAGES", "Saving package ended")
		return
	}
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-PACKAGES", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-PACKAGES", "Saving package ended")
		return
	}
	var onePackage database.Package
	db.Where("id=?", data.Id).Find(&onePackage)
	onePackage.Name = data.Name
	onePackage.PackageTypeID = int(cachedPackageTypesByName[data.Type].ID)
	onePackage.OrderID = int(cachedOrdersByName[data.Order].ID)
	onePackage.Barcode = data.Barcode
	onePackage.Note = data.Note
	db.Save(&onePackage)
	cachePackages(db)
	logInfo("SETTINGS-PACKAGES", "Package saved in "+time.Since(timer).String())
}

func loadPackageTypeDetails(id string, writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS-PACKAGES", "Loading package type details")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-PACKAGES", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-PACKAGES", "Loading package type details ended")
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
	tmpl := template.Must(template.ParseFiles("./html/settings-detail-package-type.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("SETTINGS-PACKAGES", "Package type details loaded in "+time.Since(timer).String())
}

func loadPackageDetails(id string, writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS-PACKAGES", "Loading package details")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-PACKAGES", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-PACKAGES", "Loading package details ended")
		return
	}
	var onePackage database.Package
	db.Where("id = ?", id).Find(&onePackage)
	var packageTypes []PackageTypeSelection
	for _, packageType := range cachedPackageTypesById {
		if packageType.Name == cachedPackageTypesById[uint(onePackage.PackageTypeID)].Name {
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
	for _, order := range orders {
		if order.Name == cachedOrdersById[uint(onePackage.OrderID)].Name {
			orderSelection = append(orderSelection, OrderSelection{OrderName: order.Name, OrderId: order.ID, OrderSelected: "selected"})
		} else {
			orderSelection = append(orderSelection, OrderSelection{OrderName: order.Name, OrderId: order.ID})
		}
	}
	sort.Slice(orderSelection, func(i, j int) bool {
		return orderSelection[i].OrderName < orderSelection[j].OrderName
	})
	data := PackageDetailsDataOutput{
		PackageName:            onePackage.Name,
		PackageNamePrepend:     getLocale(email, "package-name"),
		PackageTypeName:        cachedPackageTypesById[uint(onePackage.PackageTypeID)].Name,
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
	tmpl := template.Must(template.ParseFiles("./html/settings-detail-package.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("SETTINGS-PACKAGES", "Package details loaded in "+time.Since(timer).String())
}

func loadPackagesSettings(writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS-PACKAGE", "Loading packages settings")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-PACKAGE", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-PACKAGE", "Loading packages settings ended")
		return
	}

	var data PackagesSettingsDataOutput
	data.DataTableSearchTitle = getLocale(email, "data-table-search-title")
	data.DataTableInfoTitle = getLocale(email, "data-table-info-title")
	data.DataTableRowsCountTitle = getLocale(email, "data-table-rows-count-title")

	var records []database.Package
	db.Order("id desc").Find(&records)
	addPackageSettingsTableHeaders(email, &data)
	for _, record := range records {
		addPackageSettingsTableRow(record, &data)
	}

	var typeRecords []database.PackageType
	db.Order("id desc").Find(&typeRecords)
	addPackageSettingsTypeTableHeaders(email, &data)
	for _, record := range typeRecords {
		addPackageSettingsTypeTableRow(record, &data)
	}
	tmpl := template.Must(template.ParseFiles("./html/settings-table-type.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("SETTINGS-PACKAGE", "Packages settings loaded in "+time.Since(timer).String())
}

func addPackageSettingsTypeTableRow(record database.PackageType, data *PackagesSettingsDataOutput) {
	var tableRow TableRowType
	id := TableCellType{CellNameType: strconv.Itoa(int(record.ID))}
	tableRow.TableCellType = append(tableRow.TableCellType, id)
	name := TableCellType{CellNameType: record.Name}
	tableRow.TableCellType = append(tableRow.TableCellType, name)
	data.TableRowsType = append(data.TableRowsType, tableRow)
}

func addPackageSettingsTypeTableHeaders(email string, data *PackagesSettingsDataOutput) {
	id := HeaderCellType{HeaderNameType: "#", HeaderWidthType: "30"}
	data.TableHeaderType = append(data.TableHeaderType, id)
	name := HeaderCellType{HeaderNameType: getLocale(email, "type-name")}
	data.TableHeaderType = append(data.TableHeaderType, name)
}

func addPackageSettingsTableRow(record database.Package, data *PackagesSettingsDataOutput) {
	var tableRow TableRow
	id := TableCell{CellName: strconv.Itoa(int(record.ID))}
	tableRow.TableCell = append(tableRow.TableCell, id)
	name := TableCell{CellName: record.Name}
	tableRow.TableCell = append(tableRow.TableCell, name)
	data.TableRows = append(data.TableRows, tableRow)
}

func addPackageSettingsTableHeaders(email string, data *PackagesSettingsDataOutput) {
	id := HeaderCell{HeaderName: "#", HeaderWidth: "30"}
	data.TableHeader = append(data.TableHeader, id)
	name := HeaderCell{HeaderName: getLocale(email, "package-name")}
	data.TableHeader = append(data.TableHeader, name)
}