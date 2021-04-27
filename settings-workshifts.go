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
	"strings"
	"time"
)

type WorkShiftsSettingsDataOutput struct {
	DataTableSearchTitle    string
	DataTableInfoTitle      string
	DataTableRowsCountTitle string
	TableHeader             []HeaderCell
	TableRows               []TableRow
	Result                  string
}

type WorkShiftDetailsDataOutput struct {
	WorkshiftName        string
	WorkshiftNamePrepend string
	Start                string
	StartPrepend         string
	End                  string
	EndPrepend           string
	Note                 string
	NotePrepend          string
	CreatedAt            string
	CreatedAtPrepend     string
	UpdatedAt            string
	UpdatedAtPrepend     string
	Result               string
}

type WorkShiftDetailsDataInput struct {
	Id    string
	Name  string
	Start string
	End   string
	Note  string
}

func loadWorkShifts(writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS", "Loading workshifts")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS", "Problem opening database: "+err.Error())
		var responseData WorkShiftsSettingsDataOutput
		responseData.Result = "ERR: Problem opening database, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Loading workshifts ended with error")
		return
	}
	var records []database.Workshift
	db.Order("id desc").Find(&records)
	var data WorkShiftsSettingsDataOutput
	data.DataTableSearchTitle = getLocale(email, "data-table-search-title")
	data.DataTableInfoTitle = getLocale(email, "data-table-info-title")
	data.DataTableRowsCountTitle = getLocale(email, "data-table-rows-count-title")
	addWorkShiftsTableHeaders(email, &data)
	for _, record := range records {
		addWorkShiftsTableRow(record, &data)
	}
	tmpl, err := template.ParseFiles("./html/settings-table.html")
	if err != nil {
		logError("SETTINGS", "Problem parsing html file: "+err.Error())
		var responseData OrdersSettingsDataOutput
		responseData.Result = "ERR: Problem parsing html file: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
	} else {
		data.Result = "INF: Workshifts processed in " + time.Since(timer).String()
		_ = tmpl.Execute(writer, data)
		logInfo("SETTINGS", "Workshifts loaded in "+time.Since(timer).String())
	}
}

func addWorkShiftsTableRow(record database.Workshift, data *WorkShiftsSettingsDataOutput) {
	var tableRow TableRow
	id := TableCell{CellName: strconv.Itoa(int(record.ID))}
	tableRow.TableCell = append(tableRow.TableCell, id)
	name := TableCell{CellName: record.Name}
	tableRow.TableCell = append(tableRow.TableCell, name)
	data.TableRows = append(data.TableRows, tableRow)
}

func addWorkShiftsTableHeaders(email string, data *WorkShiftsSettingsDataOutput) {
	id := HeaderCell{HeaderName: "#", HeaderWidth: "30"}
	data.TableHeader = append(data.TableHeader, id)
	name := HeaderCell{HeaderName: getLocale(email, "workshift-name")}
	data.TableHeader = append(data.TableHeader, name)
}

func loadWorkshift(id string, writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS", "Loading workshift")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS", "Problem opening database: "+err.Error())
		var responseData WorkShiftDetailsDataOutput
		responseData.Result = "ERR: Problem opening database, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Loading workshift ended with error")
		return
	}
	workshiftId, _ := strconv.Atoi(id)
	workshift := cachedWorkShiftsById[uint(workshiftId)]
	data := WorkShiftDetailsDataOutput{
		WorkshiftName:        workshift.Name,
		WorkshiftNamePrepend: getLocale(email, "workshift-name"),
		Start:                workshift.WorkshiftStart.In(time.UTC).Format("15:04:05"),
		StartPrepend:         getLocale(email, "workshift-start"),
		End:                  workshift.WorkshiftEnd.In(time.UTC).Format("15:04:05"),
		EndPrepend:           getLocale(email, "workshift-end"),
		Note:                 workshift.Note,
		NotePrepend:          getLocale(email, "note-name"),
		CreatedAt:            workshift.CreatedAt.Format("2006-01-02T15:04:05"),
		CreatedAtPrepend:     getLocale(email, "created-at"),
		UpdatedAt:            workshift.UpdatedAt.Format("2006-01-02T15:04:05"),
		UpdatedAtPrepend:     getLocale(email, "updated-at"),
	}
	tmpl, err := template.ParseFiles("./html/settings-detail-workshift.html")
	if err != nil {
		logError("SETTINGS", "Problem parsing html file: "+err.Error())
		var responseData WorkShiftDetailsDataOutput
		responseData.Result = "ERR: Problem parsing html file: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
	} else {
		data.Result = "INF: Workshift detail processed in " + time.Since(timer).String()
		_ = tmpl.Execute(writer, data)
		logInfo("SETTINGS", "Workshift detail loaded in "+time.Since(timer).String())
	}
}

func saveWorkshift(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	timer := time.Now()
	logInfo("SETTINGS", "Saving workshift")
	var data WorkShiftDetailsDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "ERR: Error parsing data, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Saving workshift ended with error")
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
		logInfo("SETTINGS", "Saving workshift ended with error")
		return

	}
	workshiftStart := strings.Split(data.Start, ":")
	workshiftEnd := strings.Split(data.End, ":")
	workshiftStartHour, _ := strconv.Atoi(workshiftStart[0])
	workshiftStartMinute, _ := strconv.Atoi(workshiftStart[1])
	workshiftStartSecond, _ := strconv.Atoi(workshiftStart[2])
	workshiftEndHour, _ := strconv.Atoi(workshiftEnd[0])
	workshiftEndMinute, _ := strconv.Atoi(workshiftEnd[1])
	workshiftEndSecond, _ := strconv.Atoi(workshiftEnd[2])
	var shift database.Workshift
	db.Where("id=?", data.Id).Find(&shift)
	shift.Name = data.Name
	shift.WorkshiftStart = time.Date(2000, 1, 1, workshiftStartHour, workshiftStartMinute, workshiftStartSecond, 0, time.UTC)
	shift.WorkshiftEnd = time.Date(2000, 1, 1, workshiftEndHour, workshiftEndMinute, workshiftEndSecond, 0, time.UTC)
	result := db.Save(&shift)
	cacheWorkShifts(db)
	if result.Error != nil {
		var responseData TableOutput
		responseData.Result = "ERR: Workshift not saved: " + result.Error.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logError("SETTINGS", "Workshift "+shift.Name+" not saved: "+result.Error.Error())
	} else {
		var responseData TableOutput
		responseData.Result = "INF: Workshift saved in " + time.Since(timer).String()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS", "Workshift "+shift.Name+" saved in "+time.Since(timer).String())
	}
}
