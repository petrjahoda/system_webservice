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

type AlarmsSettingsData struct {
	AlarmName               string
	WorkplaceName           string
	SqlCommand              string
	MessageHeader           string
	MessageText             string
	Recipients              string
	Url                     string
	Pdf                     string
	CreatedAt               string
	UpdatedAt               string
	DataTableSearchTitle    string
	DataTableInfoTitle      string
	DataTableRowsCountTitle string
	TableHeader             []HeaderCell
	TableRows               []TableRow
}

type AlarmData struct {
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
}

type WorkplaceSelection struct {
	WorkplaceName     string
	WorkplaceId       uint
	WorkplaceSelected string
}

type SaveDetailAlarm struct {
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

func saveDetailAlarm(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	timer := time.Now()
	logInfo("SETTINGS-ALARMS", "Saving alarm started")
	var data SaveDetailAlarm
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS-ALARMS", "Error parsing data: "+err.Error())
		var responseData DataPageOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-ALARMS", "Saving alarm ended")
		return
	}
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-ALARMS", "Problem opening database: "+err.Error())
		var responseData DataPageOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-ALARMS", "Saving alarm ended")
		return
	}
	var alarm database.Alarm
	db.Where("id=?", data.Id).Find(&alarm)
	alarm.Name = data.Name
	alarm.WorkplaceID = int(cachedWorkplacesByName[data.Workplace].ID)
	alarm.SqlCommand = data.Sql
	alarm.MessageHeader = data.Header
	alarm.MessageText = data.Text
	alarm.Recipients = data.Recipients
	alarm.Url = data.Url
	alarm.Pdf = data.Pdf
	db.Save(&alarm)
	logInfo("SETTINGS-ALARMS", "Alarm saved in "+time.Since(timer).String())
}

func processAlarmsSettings(writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS-ALARMS", "Processing alarms settings started")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-ALARMS", "Problem opening database: "+err.Error())
		var responseData DataPageOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-ALARMS", "Processing alarms settings ended")
		return
	}
	var alarms []database.Alarm
	db.Find(&alarms)

	var data AlarmsSettingsData
	data.DataTableSearchTitle = getLocale(email, "data-table-search-title")
	data.DataTableInfoTitle = getLocale(email, "data-table-info-title")
	data.DataTableRowsCountTitle = getLocale(email, "data-table-rows-count-title")
	data.AlarmName = getLocale(email, "alarm-name")
	data.WorkplaceName = getLocale(email, "workplace-name")
	data.SqlCommand = getLocale(email, "sql-command")
	data.MessageHeader = getLocale(email, "message-header")
	data.MessageText = getLocale(email, "message-text")
	data.Recipients = getLocale(email, "recipients")
	data.Url = getLocale(email, "url")
	data.Pdf = getLocale(email, "pdf")
	data.CreatedAt = getLocale(email, "created-at")
	data.UpdatedAt = getLocale(email, "updated-at")

	addAlarmSettingsTableHeaders(email, &data)
	for _, record := range alarms {
		addAlarmSettingsTableRow(record, &data)
	}

	tmpl := template.Must(template.ParseFiles("./html/settings-alarms.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("SETTINGS-ALARMS", "Alarms settings processed in "+time.Since(timer).String())
}

func addAlarmSettingsTableRow(record database.Alarm, data *AlarmsSettingsData) {
	var tableRow TableRow
	id := TableCell{CellName: strconv.Itoa(int(record.ID))}
	tableRow.TableCell = append(tableRow.TableCell, id)
	alarmName := TableCell{CellName: record.Name}
	tableRow.TableCell = append(tableRow.TableCell, alarmName)
	data.TableRows = append(data.TableRows, tableRow)
}

func addAlarmSettingsTableHeaders(email string, data *AlarmsSettingsData) {
	id := HeaderCell{HeaderName: "#"}
	data.TableHeader = append(data.TableHeader, id)
	alarmName := HeaderCell{HeaderName: getLocale(email, "alarm-name")}
	data.TableHeader = append(data.TableHeader, alarmName)
}

func processDetailAlarmSettings(id string, writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS-ALARMS", "Processing detail alarm settings started")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-ALARMS", "Problem opening database: "+err.Error())
		var responseData DataPageOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-ALARMS", "Processing detail alarm settings ended")
		return
	}
	var alarm database.Alarm
	db.Where("id = ?", id).Find(&alarm)
	var workplaces []WorkplaceSelection
	for _, workplace := range cachedWorkplacesById {
		if workplace.Name == cachedWorkplacesById[uint(alarm.WorkplaceID)].Name {
			workplaces = append(workplaces, WorkplaceSelection{WorkplaceName: workplace.Name, WorkplaceId: workplace.ID, WorkplaceSelected: "selected"})
		} else {
			workplaces = append(workplaces, WorkplaceSelection{WorkplaceName: workplace.Name, WorkplaceId: workplace.ID})
		}
	}

	sort.Slice(workplaces, func(i, j int) bool {
		return workplaces[i].WorkplaceName < workplaces[j].WorkplaceName
	})
	data := AlarmData{
		AlarmName:            alarm.Name,
		AlarmNamePrepend:     getLocale(email, "alarm-name"),
		WorkplaceName:        cachedWorkplacesById[uint(alarm.WorkplaceID)].Name,
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
		CreatedAt:            alarm.CreatedAt.Format("2006-01-02 15:04:05"),
		CreatedAtPrepend:     getLocale(email, "created-at"),
		UpdatedAt:            alarm.UpdatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAtPrepend:     getLocale(email, "updated-at"),
		Workplaces:           workplaces,
	}
	tmpl := template.Must(template.ParseFiles("./html/settings-alarms-details.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("SETTINGS-ALARMS", "Detail alarm settings processed in "+time.Since(timer).String())
}
