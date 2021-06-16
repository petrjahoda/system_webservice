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

type AlarmsSettingsDataOutput struct {
	DataTableSearchTitle    string
	DataTableInfoTitle      string
	DataTableRowsCountTitle string
	TableHeader             []HeaderCell
	TableRows               []TableRow
	Result                  string
}

type AlarmDetailsDataInput struct {
	Id         string
	Name       string
	Workplace  string
	Sql        string
	Header     string
	Text       string
	Recipients string
	Url        string
	Pdf        string
}

type AlarmDetailsDataOutput struct {
	AlarmName            string
	AlarmNamePrepend     string
	WorkplaceName        string
	WorkplaceNamePrepend string
	SqlCommand           string
	SqlCommandPrepend    string
	MessageHeader        string
	MessageHeaderPrepend string
	MessageText          string
	MessageTextPrepend   string
	Recipients           string
	RecipientsPrepend    string
	Url                  string
	UrlPrepend           string
	Pdf                  string
	PdfPrepend           string
	CreatedAt            string
	CreatedAtPrepend     string
	UpdatedAt            string
	UpdatedAtPrepend     string
	Workplaces           []WorkplaceSelection
	Result               string
}

type WorkplaceSelection struct {
	WorkplaceName     string
	WorkplaceId       uint
	WorkplaceSelected string
}

func loadAlarms(writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS", "Loading alarms")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS", "Problem opening database: "+err.Error())
		var responseData AlarmsSettingsDataOutput
		responseData.Result = "ERR: Problem opening database, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Loading alarms ended with error")
		return
	}
	var data AlarmsSettingsDataOutput
	data.DataTableSearchTitle = getLocale(email, "data-table-search-title")
	data.DataTableInfoTitle = getLocale(email, "data-table-info-title")
	data.DataTableRowsCountTitle = getLocale(email, "data-table-rows-count-title")
	var records []database.Alarm

	db.Order("id desc").Find(&records)
	addAlarmSettingsTableHeaders(email, &data)
	for _, record := range records {
		addAlarmSettingsTableRow(record, &data)
	}
	tmpl, err := template.ParseFiles("./html/settings-table.html")
	if err != nil {
		logError("SETTINGS", "Problem parsing html file: "+err.Error())
		var responseData AlarmsSettingsDataOutput
		responseData.Result = "ERR: Problem parsing html file: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
	} else {
		data.Result = "INF: Alarms processed in " + time.Since(timer).String()
		_ = tmpl.Execute(writer, data)
		logInfo("SETTINGS", "Alarms loaded in "+time.Since(timer).String())
	}
}

func addAlarmSettingsTableRow(record database.Alarm, data *AlarmsSettingsDataOutput) {
	var tableRow TableRow
	id := TableCell{CellName: strconv.Itoa(int(record.ID))}
	tableRow.TableCell = append(tableRow.TableCell, id)
	name := TableCell{CellName: record.Name}
	tableRow.TableCell = append(tableRow.TableCell, name)
	data.TableRows = append(data.TableRows, tableRow)
}

func addAlarmSettingsTableHeaders(email string, data *AlarmsSettingsDataOutput) {
	id := HeaderCell{HeaderName: "#", HeaderWidth: "30"}
	data.TableHeader = append(data.TableHeader, id)
	name := HeaderCell{HeaderName: getLocale(email, "alarm-name")}
	data.TableHeader = append(data.TableHeader, name)
}

func loadAlarm(id string, writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS", "Loading alarm")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS", "Problem opening database: "+err.Error())
		var responseData AlarmDetailsDataOutput
		responseData.Result = "ERR: Problem opening database, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Loading alarm ended with error")
		return
	}
	var alarm database.Alarm
	db.Where("id = ?", id).Find(&alarm)
	var workplaces []WorkplaceSelection
	workplacesByIdSync.RLock()
	workplacesById := cachedWorkplacesById
	workplacesByIdSync.RUnlock()
	for _, workplace := range workplacesById {
		if workplace.Name == workplacesById[uint(alarm.WorkplaceID)].Name {
			workplaces = append(workplaces, WorkplaceSelection{WorkplaceName: workplace.Name, WorkplaceId: workplace.ID, WorkplaceSelected: "selected"})
		} else {
			workplaces = append(workplaces, WorkplaceSelection{WorkplaceName: workplace.Name, WorkplaceId: workplace.ID})
		}
	}
	sort.Slice(workplaces, func(i, j int) bool {
		return workplaces[i].WorkplaceName < workplaces[j].WorkplaceName
	})
	data := AlarmDetailsDataOutput{
		AlarmName:        alarm.Name,
		AlarmNamePrepend: getLocale(email, "alarm-name"),

		WorkplaceNamePrepend: getLocale(email, "workplace-name"),
		SqlCommand:           alarm.SqlCommand,
		SqlCommandPrepend:    getLocale(email, "sql-command"),
		MessageHeader:        alarm.MessageHeader,
		MessageHeaderPrepend: getLocale(email, "message-header"),
		MessageText:          alarm.MessageText,
		MessageTextPrepend:   getLocale(email, "message-text"),
		Recipients:           alarm.Recipients,
		RecipientsPrepend:    getLocale(email, "recipients"),
		Url:                  alarm.Url,
		UrlPrepend:           getLocale(email, "url"),
		Pdf:                  alarm.Pdf,
		PdfPrepend:           getLocale(email, "pdf"),
		CreatedAt:            alarm.CreatedAt.Format("2006-01-02T15:04:05"),
		CreatedAtPrepend:     getLocale(email, "created-at"),
		UpdatedAt:            alarm.UpdatedAt.Format("2006-01-02T15:04:05"),
		UpdatedAtPrepend:     getLocale(email, "updated-at"),
		Workplaces:           workplaces,
	}
	workplacesByIdSync.RLock()
	data.WorkplaceName = cachedWorkplacesById[uint(alarm.WorkplaceID)].Name
	workplacesByIdSync.RUnlock()

	tmpl, err := template.ParseFiles("./html/settings-detail-alarm.html")
	if err != nil {
		logError("SETTINGS", "Problem parsing html file: "+err.Error())
		var responseData AlarmsSettingsDataOutput
		responseData.Result = "ERR: Problem parsing html file: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
	} else {
		data.Result = "INF: Alarm detail processed in " + time.Since(timer).String()
		_ = tmpl.Execute(writer, data)
		logInfo("SETTINGS", "Alarm detail loaded in "+time.Since(timer).String())
	}
}

func saveAlarm(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	timer := time.Now()
	logInfo("SETTINGS", "Saving alarm")
	var data AlarmDetailsDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "ERR: Error parsing data, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Saving alarm ended with error")
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
		logInfo("SETTINGS", "Saving alarm ended with error")
		return
	}
	var alarm database.Alarm
	db.Where("id=?", data.Id).Find(&alarm)
	alarm.Name = data.Name
	workplacesByNameSync.RLock()
	alarm.WorkplaceID = int(cachedWorkplacesByName[data.Workplace].ID)
	workplacesByNameSync.RUnlock()
	alarm.SqlCommand = data.Sql
	alarm.MessageHeader = data.Header
	alarm.MessageText = data.Text
	alarm.Recipients = data.Recipients
	alarm.Url = data.Url
	alarm.Pdf = data.Pdf
	result := db.Save(&alarm)
	cacheAlarms(db)
	if result.Error != nil {
		var responseData TableOutput
		responseData.Result = "ERR: Alarm not saved: " + result.Error.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logError("SETTINGS", "Alarm "+alarm.Name+" not saved: "+result.Error.Error())
	} else {
		var responseData TableOutput
		responseData.Result = "INF: Alarm saved in " + time.Since(timer).String()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Alarm "+alarm.Name+" saved in "+time.Since(timer).String())
	}
}
