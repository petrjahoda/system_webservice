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

type WorkshiftsSettingsDataOutput struct {
	DataTableSearchTitle    string
	DataTableInfoTitle      string
	DataTableRowsCountTitle string
	TableHeader             []HeaderCell
	TableRows               []TableRow
}

type WorkshiftDetailsDataOutput struct {
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
}

type WorkshiftDetailsDataInput struct {
	Id    string
	Name  string
	Start string
	End   string
	Note  string
}

func saveWorkshift(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	timer := time.Now()
	logInfo("SETTINGS-WORKSHIFTS", "Saving workshift started")
	var data WorkshiftDetailsDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("SETTINGS-WORKSHIFTS", "Error parsing data: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-WORKSHIFTS", "Saving workshift ended")
		return
	}
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-WORKSHIFTS", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-WORKSHIFTS", "Saving workshift ended")
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
	db.Save(&shift)
	cacheWorkShifts(db)
	logInfo("SETTINGS-WORKSHIFTS", "Workshift saved in "+time.Since(timer).String())
}

func loadWorkshiftDetails(id string, writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS-WORKSHIFTS", "Loading workshift details")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-WORKSHIFTS", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-WORKSHIFTS", "Loading workshift details ended")
		return
	}
	workshiftId, _ := strconv.Atoi(id)
	workshift := cachedWorkshiftsById[uint(workshiftId)]

	data := WorkshiftDetailsDataOutput{
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
	tmpl := template.Must(template.ParseFiles("./html/settings-detail-workshift.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("SETTINGS-WORKSHIFTS", "Workshift details loaded in "+time.Since(timer).String())
}

func loadWorkshiftsSettings(writer http.ResponseWriter, email string) {
	timer := time.Now()
	logInfo("SETTINGS-WORKSHIFTS", "Loading workshifts settings")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("SETTINGS-WORKSHIFTS", "Problem opening database: "+err.Error())
		var responseData TableOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("SETTINGS-WORKSHIFTS", "Loading workshifts settings ended")
		return
	}
	var records []database.Workshift
	db.Order("id desc").Find(&records)
	var data WorkshiftsSettingsDataOutput
	data.DataTableSearchTitle = getLocale(email, "data-table-search-title")
	data.DataTableInfoTitle = getLocale(email, "data-table-info-title")
	data.DataTableRowsCountTitle = getLocale(email, "data-table-rows-count-title")
	addWorkshiftSettingsTableHeaders(email, &data)
	for _, record := range records {
		addWorkshiftSettingsTableRow(record, &data)
	}
	tmpl := template.Must(template.ParseFiles("./html/settings-table.html"))
	_ = tmpl.Execute(writer, data)
	logInfo("SETTINGS-WORKSHIFTS", "Workshifts settings loaded in "+time.Since(timer).String())
}

func addWorkshiftSettingsTableRow(record database.Workshift, data *WorkshiftsSettingsDataOutput) {
	var tableRow TableRow
	id := TableCell{CellName: strconv.Itoa(int(record.ID))}
	tableRow.TableCell = append(tableRow.TableCell, id)
	name := TableCell{CellName: record.Name}
	tableRow.TableCell = append(tableRow.TableCell, name)
	data.TableRows = append(data.TableRows, tableRow)
}

func addWorkshiftSettingsTableHeaders(email string, data *WorkshiftsSettingsDataOutput) {
	id := HeaderCell{HeaderName: "#", HeaderWidth: "30"}
	data.TableHeader = append(data.TableHeader, id)
	name := HeaderCell{HeaderName: getLocale(email, "workshift-name")}
	data.TableHeader = append(data.TableHeader, name)
}
