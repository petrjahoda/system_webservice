package main

import (
	"encoding/json"
	"github.com/petrjahoda/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"html/template"
	"net/http"
	"sort"
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
	AlarmName     string
	WorkplaceName string
	SqlCommand    string
	MessageHeader string
	MessageText   string
	Recipients    string
	Url           string
	Pdf           string
	CreatedAt     string
	UpdatedAt     string
	Workplaces    []string
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
	alarmName := TableCell{CellName: record.Name}
	tableRow.TableCell = append(tableRow.TableCell, alarmName)
	workplaceName := TableCell{CellName: cachedWorkplacesById[uint(record.WorkplaceID)].Name}
	tableRow.TableCell = append(tableRow.TableCell, workplaceName)
	header := TableCell{CellName: record.MessageHeader}
	tableRow.TableCell = append(tableRow.TableCell, header)
	text := TableCell{CellName: record.MessageText}
	tableRow.TableCell = append(tableRow.TableCell, text)
	recipients := TableCell{CellName: record.Recipients}
	tableRow.TableCell = append(tableRow.TableCell, recipients)
	data.TableRows = append(data.TableRows, tableRow)
}

func addAlarmSettingsTableHeaders(email string, data *AlarmsSettingsData) {
	alarmName := HeaderCell{HeaderName: getLocale(email, "alarm-name")}
	data.TableHeader = append(data.TableHeader, alarmName)
	workplaceName := HeaderCell{HeaderName: getLocale(email, "workplace-name")}
	data.TableHeader = append(data.TableHeader, workplaceName)
	header := HeaderCell{HeaderName: getLocale(email, "message-header")}
	data.TableHeader = append(data.TableHeader, header)
	text := HeaderCell{HeaderName: getLocale(email, "message-text")}
	data.TableHeader = append(data.TableHeader, text)
	recipients := HeaderCell{HeaderName: getLocale(email, "recipients")}
	data.TableHeader = append(data.TableHeader, recipients)
}

func processDetailAlarmSettings(name string, writer http.ResponseWriter, email string) {
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
	db.Where("name = ?", name).Find(&alarm)
	var workplaces []string
	for _, workplace := range cachedWorkplacesById {
		workplaces = append(workplaces, workplace.Name)
	}
	sort.Strings(workplaces)
	alarmData := AlarmData{
		AlarmName:     alarm.Name,
		WorkplaceName: cachedWorkplacesById[uint(alarm.WorkplaceID)].Name,
		SqlCommand:    alarm.SqlCommand,
		MessageHeader: alarm.MessageHeader,
		MessageText:   alarm.MessageText,
		Recipients:    alarm.Recipients,
		Url:           alarm.Url,
		Pdf:           alarm.Pdf,
		CreatedAt:     alarm.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     alarm.UpdatedAt.Format("2006-01-02 15:04:05"),
		Workplaces:    workplaces,
	}
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(alarmData)
	logInfo("SETTINGS-ALARMS", "Detail alarm settings processed in "+time.Since(timer).String())
}
