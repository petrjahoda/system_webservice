package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/petrjahoda/database"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type UsersSettingsDataOutput struct {
	DataTableSearchTitle    string
	DataTableInfoTitle      string
	DataTableRowsCountTitle string
	TableHeader             []HeaderCell
	TableRows               []TableRow
	TableHeaderType         []HeaderCellType
	TableRowsType           []TableRowType
}

type UserTypeDetailsDataOutput struct {
	UserTypeName        string
	UserTypeNamePrepend string
	Note                string
	NotePrepend         string
	CreatedAt           string
	CreatedAtPrepend    string
	UpdatedAt           string
	UpdatedAtPrepend    string
}

type UserDetailsDataOutput struct {
	FirstName            string
	FirstNamePrepend     string
	SecondName           string
	SecondNamePrepend    string
	UserTypeName         string
	UserTypeNamePrepend  string
	UserRoleName         string
	UserRoleNamePrepend  string
	SelectionName        string
	SelectionNamePrepend string
	Barcode              string
	BarcodePrepend       string
	Email                string
	EmailPrepend         string
	Password             string
	PasswordPrepend      string
	Phone                string
	PhonePrepend         string
	PIN                  string
	PINPrepend           string
	Rfid                 string
	RfidPrepend          string
	Position             string
	PositionPrepend      string
	Locale               string
	LocalePrepend        string
	Note                 string
	NotePrepend          string
	CreatedAt            string
	CreatedAtPrepend     string
	UpdatedAt            string
	UpdatedAtPrepend     string
	UserTypes            []UserTypeSelection
	UserRoles            []UserRoleSelection
	Locales              []LocaleSelection
}

type UserTypeSelection struct {
	UserTypeName     string
	UserTypeId       uint
	UserTypeSelected string
}

type UserRoleSelection struct {
	UserRoleName     string
	UserRoleId       uint
	UserRoleSelected string
}

type LocaleSelection struct {
	LocaleName     string
	LocaleSelected string
}

type UserDetailsDataInput struct {
	Id         string
	FirstName  string
	SecondName string
	Type       string
	Role       string
	Locale     string
	Barcode    string
	Password   string
	Email      string
	Phone      string
	Pin        string
	Position   string
	Rfid       string
	Note       string
}

type UserTypeDetailsDataInput struct {
	Id   string
	Name string
	Note string
}

func loadUsers(writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS", "Loading users")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Loading users ended with error")
		return
	}
	var data UsersSettingsDataOutput
	data.DataTableSearchTitle = getLocale(email, "data-table-search-title")
	data.DataTableInfoTitle = getLocale(email, "data-table-info-title")
	data.DataTableRowsCountTitle = getLocale(email, "data-table-rows-count-title")
	var records []database.User
	db.Order("id desc").Find(&records)
	addUsersTableHeaders(email, &data)
	for _, record := range records {
		addUsersTableRow(record, &data)
	}
	var typeRecords []database.UserType
	db.Order("id desc").Find(&typeRecords)
	addUserTypesTableHeaders(email, &data)
	for _, record := range typeRecords {
		addUserTypesTableRow(record, &data)
	}
	tmpl := template.Must(template.ParseFiles("./html/settings-table-type.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("SETTINGS", "Users loaded in "+time.Since(timer).String())
}

func addUsersTableRow(record database.User, data *UsersSettingsDataOutput) {
	var tableRow TableRow
	id := TableCell{CellName: strconv.Itoa(int(record.ID))}
	tableRow.TableCell = append(tableRow.TableCell, id)
	name := TableCell{CellName: record.FirstName + " " + record.SecondName}
	tableRow.TableCell = append(tableRow.TableCell, name)
	data.TableRows = append(data.TableRows, tableRow)
}

func addUsersTableHeaders(email string, data *UsersSettingsDataOutput) {
	id := HeaderCell{HeaderName: "#", HeaderWidth: "30"}
	data.TableHeader = append(data.TableHeader, id)
	name := HeaderCell{HeaderName: getLocale(email, "user-name")}
	data.TableHeader = append(data.TableHeader, name)
}

func addUserTypesTableRow(record database.UserType, data *UsersSettingsDataOutput) {
	var tableRow TableRowType
	id := TableCellType{CellNameType: strconv.Itoa(int(record.ID))}
	tableRow.TableCellType = append(tableRow.TableCellType, id)
	name := TableCellType{CellNameType: record.Name}
	tableRow.TableCellType = append(tableRow.TableCellType, name)
	data.TableRowsType = append(data.TableRowsType, tableRow)
}

func addUserTypesTableHeaders(email string, data *UsersSettingsDataOutput) {
	id := HeaderCellType{HeaderNameType: "#", HeaderWidthType: "30"}
	data.TableHeaderType = append(data.TableHeaderType, id)
	name := HeaderCellType{HeaderNameType: getLocale(email, "type-name")}
	data.TableHeaderType = append(data.TableHeaderType, name)
}

func loadUser(id string, writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS", "Loading user")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Loading user ended with error")
		return
	}
	var user database.User
	db.Where("id = ?", id).Find(&user)
	var userTypes []UserTypeSelection
	for _, userType := range cachedUserTypesById {
		if userType.Name == cachedUserTypesById[uint(user.UserTypeID)].Name {
			userTypes = append(userTypes, UserTypeSelection{UserTypeName: userType.Name, UserTypeId: userType.ID, UserTypeSelected: "selected"})
		} else {
			userTypes = append(userTypes, UserTypeSelection{UserTypeName: userType.Name, UserTypeId: userType.ID})
		}
	}
	sort.Slice(userTypes, func(i, j int) bool {
		return userTypes[i].UserTypeName < userTypes[j].UserTypeName
	})
	var userRoles []UserRoleSelection
	for _, userRole := range cachedUserRolesById {
		if userRole.Name == cachedUserRolesById[uint(user.UserRoleID)].Name {
			userRoles = append(userRoles, UserRoleSelection{UserRoleName: userRole.Name, UserRoleId: userRole.ID, UserRoleSelected: "selected"})
		} else {
			userRoles = append(userRoles, UserRoleSelection{UserRoleName: userRole.Name, UserRoleId: userRole.ID})
		}
	}
	sort.Slice(userTypes, func(i, j int) bool {
		return userTypes[i].UserTypeName < userTypes[j].UserTypeName
	})
	var locales []LocaleSelection
	locales = append(locales, LocaleSelection{LocaleName: "CsCZ", LocaleSelected: testLocaleForUser(user.Locale, "CsCZ")})
	locales = append(locales, LocaleSelection{LocaleName: "DeDE", LocaleSelected: testLocaleForUser(user.Locale, "DeDE")})
	locales = append(locales, LocaleSelection{LocaleName: "EnUS", LocaleSelected: testLocaleForUser(user.Locale, "EnUS")})
	locales = append(locales, LocaleSelection{LocaleName: "EsES", LocaleSelected: testLocaleForUser(user.Locale, "EsES")})
	locales = append(locales, LocaleSelection{LocaleName: "FrFR", LocaleSelected: testLocaleForUser(user.Locale, "FrFR")})
	locales = append(locales, LocaleSelection{LocaleName: "ItIT", LocaleSelected: testLocaleForUser(user.Locale, "ItIT")})
	locales = append(locales, LocaleSelection{LocaleName: "PlPL", LocaleSelected: testLocaleForUser(user.Locale, "PlPL")})
	locales = append(locales, LocaleSelection{LocaleName: "PtPT", LocaleSelected: testLocaleForUser(user.Locale, "PtPT")})
	locales = append(locales, LocaleSelection{LocaleName: "SkSK", LocaleSelected: testLocaleForUser(user.Locale, "SkSK")})
	locales = append(locales, LocaleSelection{LocaleName: "RuRU", LocaleSelected: testLocaleForUser(user.Locale, "RuRU")})
	data := UserDetailsDataOutput{
		FirstName:           user.FirstName,
		FirstNamePrepend:    getLocale(email, "first-name"),
		SecondName:          user.SecondName,
		SecondNamePrepend:   getLocale(email, "second-name"),
		UserRoleName:        cachedUserRolesById[uint(user.UserRoleID)].Name,
		UserRoleNamePrepend: getLocale(email, "role-name"),
		UserTypeName:        cachedUserTypesById[uint(user.UserTypeID)].Name,
		UserTypeNamePrepend: getLocale(email, "type-name"),
		LocalePrepend:       getLocale(email, "locale"),
		Barcode:             user.Barcode,
		BarcodePrepend:      getLocale(email, "barcode"),
		Email:               user.Email,
		EmailPrepend:        getLocale(email, "email"),
		Password:            "",
		PasswordPrepend:     getLocale(email, "password"),
		Phone:               user.Phone,
		PhonePrepend:        getLocale(email, "phone"),
		PIN:                 user.Pin,
		PINPrepend:          getLocale(email, "pin"),
		Position:            user.Position,
		PositionPrepend:     getLocale(email, "position"),
		Rfid:                user.Rfid,
		RfidPrepend:         getLocale(email, "rfid"),
		Note:                user.Note,
		NotePrepend:         getLocale(email, "note-name"),
		CreatedAt:           user.CreatedAt.Format("2006-01-02T15:04:05"),
		CreatedAtPrepend:    getLocale(email, "created-at"),
		UpdatedAt:           user.UpdatedAt.Format("2006-01-02T15:04:05"),
		UpdatedAtPrepend:    getLocale(email, "updated-at"),
		UserTypes:           userTypes,
		UserRoles:           userRoles,
		Locales:             locales,
	}
	tmpl := template.Must(template.ParseFiles("./html/settings-detail-user.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("SETTINGS", "User "+user.FirstName+" "+user.SecondName+" loaded in "+time.Since(timer).String())
}

func loadUserType(id string, writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS", "Loading user type")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Loading user type ended with error")
		return
	}
	var userType database.UserType
	db.Where("id = ?", id).Find(&userType)
	data := UserTypeDetailsDataOutput{
		UserTypeName:        userType.Name,
		UserTypeNamePrepend: getLocale(email, "type-name"),
		Note:                userType.Note,
		NotePrepend:         getLocale(email, "note-name"),
		CreatedAt:           userType.CreatedAt.Format("2006-01-02T15:04:05"),
		CreatedAtPrepend:    getLocale(email, "created-at"),
		UpdatedAt:           userType.UpdatedAt.Format("2006-01-02T15:04:05"),
		UpdatedAtPrepend:    getLocale(email, "updated-at"),
	}
	tmpl := template.Must(template.ParseFiles("./html/settings-detail-user-type.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("SETTINGS", "User type "+userType.Name+" loaded in "+time.Since(timer).String())
}

func saveUser(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	timer := time.Now()
	logInfo("SETTINGS", "Saving user")
	var data UserDetailsDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Saving user ended with error")
		return
	}
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Saving user ended with error")
		return
	}
	var user database.User
	db.Where("id=?", data.Id).Find(&user)
	user.FirstName = data.FirstName
	user.SecondName = data.SecondName
	user.UserTypeID = int(cachedUserTypesByName[data.Type].ID)
	user.UserRoleID = int(cachedUserRolesByName[data.Role].ID)
	user.Barcode = data.Barcode
	user.Email = data.Email
	user.Phone = data.Phone
	user.Pin = data.Pin
	user.Rfid = data.Rfid
	user.Note = data.Note
	user.Locale = data.Locale
	if len(data.Password) > 0 {
		user.Password = hashPasswordFromString([]byte(data.Password))
	}
	db.Save(&user)
	cacheUsers(db)
	logInfo("SETTINGS", "User "+user.FirstName+" "+user.SecondName+" saved in "+time.Since(timer).String())
}

func saveUserType(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	timer := time.Now()
	logInfo("SETTINGS", "Saving user type")
	var data UserTypeDetailsDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Saving user type ended with error")
		return
	}
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Saving user type ended with error")
		return
	}
	var userType database.UserType
	db.Where("id=?", data.Id).Find(&userType)
	userType.Name = data.Name
	userType.Note = data.Note
	db.Save(&userType)
	cacheUsers(db)
	logInfo("SETTINGS", "User type "+userType.Name+" saved in "+time.Since(timer).String())
}

func testLocaleForUser(userLocale string, locale string) string {
	if userLocale == locale {
		return "selected"
	}
	return ""
}

func hashPasswordFromString(pwd []byte) string {
	logInfo("SETTINGS", "Hashing password")
	timer := time.Now()
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		logError("SETTINGS", "Cannot hash password: "+err.Error())
		return ""
	}
	logInfo("SETTINGS", "Password hashed in  "+time.Since(timer).String())
	return string(hash)
}
