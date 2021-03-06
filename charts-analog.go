package main

import (
	"encoding/json"
	"github.com/petrjahoda/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"math"
	"net/http"
	"strconv"
	"time"
)

func processAnalogData(writer http.ResponseWriter, workplaceName string, dateFrom time.Time, dateTo time.Time, email string, chartName string) {
	timer := time.Now()
	logInfo("CHARTS-ANALOG", "Processing analog chart data started")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("CHARTS-ANALOG", "Problem opening database: "+err.Error())
		var responseData ChartDataPageOutput
		responseData.Result = "ERR: Problem opening database, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("CHARTS-ANALOG", "Processing data ended")
		return
	}
	var responseData ChartDataPageOutput
	var analogOutputData []PortData
	workplaceDevicePortsSync.RLock()
	allWorkplacePorts := cachedWorkplaceDevicePorts[workplaceName]
	workplaceDevicePortsSync.RUnlock()
	for _, port := range allWorkplacePorts {
		if port.DevicePortTypeID == analog {
			analogTimeDifference := 20.0
			var analogData []database.DevicePortAnalogRecord
			db.Select("date_time, data").Where("date_time >= ?", dateFrom).Where("date_time <= ?", dateTo).Where("device_port_id = ?", port.ID).Order("date_time").Order("id").Find(&analogData)
			var portData PortData
			portData.PortName = "ID" + strconv.Itoa(int(port.ID)) + ": " + port.Name + " (" + port.Unit + ")"
			devicePortsColorsByIdSync.RLock()
			portColor := cachedDevicePortsColorsById[int(port.ID)]
			devicePortsColorsByIdSync.RUnlock()
			if portColor == "#000000" {
				portData.PortColor = ""
			}
			date := dateFrom
			var initialData Data
			initialData.Time = dateFrom.Unix()
			initialData.Value = math.MinInt16
			portData.AnalogData = append(portData.AnalogData, initialData)
			for _, data := range analogData {
				if data.DateTime.Sub(date).Seconds() > analogTimeDifference {
					initialData.Time = date.Add(1 * time.Second).Unix()
					initialData.Value = math.MinInt16
					portData.AnalogData = append(portData.AnalogData, initialData)
				}
				initialData.Time = data.DateTime.Unix()
				initialData.Value = data.Data
				date = data.DateTime
				portData.AnalogData = append(portData.AnalogData, initialData)
			}
			initialData.Time = dateTo.Unix()
			initialData.Value = math.MinInt16
			portData.AnalogData = append(portData.AnalogData, initialData)
			analogOutputData = append(analogOutputData, portData)
		}
	}
	responseData.ChartData = analogOutputData

	orderData := downloadChartOrderData(db, dateTo, dateFrom, workplaceName)
	responseData.OrderData = orderData
	downtimeData := downloadChartDowntimeData(db, dateTo, dateFrom, workplaceName)
	responseData.DowntimeData = downtimeData
	breakdownData := downloadChartBreakdownData(db, dateTo, dateFrom, workplaceName)
	responseData.BreakdownData = breakdownData
	alarmData := downloadChartAlarmData(db, dateTo, dateFrom, workplaceName)
	responseData.AlarmData = alarmData
	userData := downloadChartUserData(db, dateTo, dateFrom, workplaceName)
	responseData.UserData = userData
	responseData.OrdersLocale = getLocale(email, "orders")
	responseData.DowntimesLocale = getLocale(email, "downtimes")
	responseData.BreakdownsLocale = getLocale(email, "breakdowns")
	responseData.UsersLocale = getLocale(email, "users")
	responseData.AlarmsLocale = getLocale(email, "alarms")
	usersByEmailSync.RLock()
	responseData.Locale = cachedUsersByEmail[email].Locale
	usersByEmailSync.RUnlock()
	responseData.Result = "INF: Analog chart data downloaded from database in " + time.Since(timer).String()
	responseData.Type = chartName
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(responseData)
	logInfo("CHARTS-ANALOG", "Analog chart data processed in "+time.Since(timer).String())
}

func downloadChartUserData(db *gorm.DB, dateTo time.Time, dateFrom time.Time, workplaceName string) []TerminalData {
	var userData []TerminalData
	var userRecords []database.UserRecord
	workplacesByNameSync.RLock()
	workplaceId := cachedWorkplacesByName[workplaceName].ID
	workplacesByNameSync.RUnlock()
	db.Where("date_time_start <= ? and date_time_end >= ?", dateTo, dateFrom).Where("workplace_id = ?", workplaceId).Or("date_time_start <= ? and date_time_end is null", dateTo).Where("workplace_id = ?", workplaceId).Or("date_time_start <= ? and date_time_end >= ?", dateFrom, dateTo).Where("workplace_id = ?", workplaceId).Find(&userRecords)
	for _, record := range userRecords {
		dateTimeStart := record.DateTimeStart.Unix() * 1000
		if record.DateTimeStart.Before(dateFrom) {
			dateTimeStart = dateFrom.Unix() * 1000
		}
		dateTimeEnd := record.DateTimeEnd.Time.Unix() * 1000
		if record.DateTimeEnd.Time.IsZero() {
			dateTimeEnd = dateTo.Unix() * 1000
		}
		oneData := TerminalData{
			Color:    "#2274A5",
			FromDate: dateTimeStart,
			ToDate:   dateTimeEnd,
			Note:     record.Note,
		}
		usersByIdSync.RLock()
		oneData.Information = cachedUsersById[uint(record.UserID)].FirstName + " " + cachedUsersById[uint(record.UserID)].SecondName
		usersByIdSync.RUnlock()
		userData = append(userData, oneData)
	}
	return userData
}

func downloadChartAlarmData(db *gorm.DB, dateTo time.Time, dateFrom time.Time, workplaceName string) []TerminalData {
	var alarmData []TerminalData
	var alarmRecords []database.AlarmRecord
	workplacesByNameSync.RLock()
	workplaceId := cachedWorkplacesByName[workplaceName].ID
	workplacesByNameSync.RUnlock()
	db.Where("date_time_start <= ? and date_time_end >= ?", dateTo, dateFrom).Where("workplace_id = ?", workplaceId).Or("date_time_start <= ? and date_time_end is null", dateTo).Where("workplace_id = ?", workplaceId).Or("date_time_start <= ? and date_time_end >= ?", dateFrom, dateTo).Where("workplace_id = ?", workplaceId).Find(&alarmRecords)
	for _, record := range alarmRecords {
		dateTimeStart := record.DateTimeStart.Unix() * 1000
		if record.DateTimeStart.Before(dateFrom) {
			dateTimeStart = dateFrom.Unix() * 1000
		}
		dateTimeEnd := record.DateTimeEnd.Time.Unix() * 1000
		if record.DateTimeEnd.Time.IsZero() {
			dateTimeEnd = dateTo.Unix() * 1000
		}
		oneData := TerminalData{
			Color:    "grey",
			FromDate: dateTimeStart,
			ToDate:   dateTimeEnd,
		}
		alarmByIdSync.RLock()
		oneData.Information = cachedAlarmsById[uint(record.AlarmID)].Name
		alarmByIdSync.RUnlock()
		alarmData = append(alarmData, oneData)
	}
	return alarmData
}

func downloadChartBreakdownData(db *gorm.DB, dateTo time.Time, dateFrom time.Time, workplaceName string) []TerminalData {
	var breakdownData []TerminalData
	var breakdownRecords []database.BreakdownRecord
	workplacesByNameSync.RLock()
	workplaceId := cachedWorkplacesByName[workplaceName].ID
	workplacesByNameSync.RUnlock()
	db.Where("date_time_start <= ? and date_time_end >= ?", dateTo, dateFrom).Where("workplace_id = ?", workplaceId).Or("date_time_start <= ? and date_time_end is null", dateTo).Where("workplace_id = ?", workplaceId).Or("date_time_start <= ? and date_time_end >= ?", dateFrom, dateTo).Where("workplace_id = ?", workplaceId).Find(&breakdownRecords)
	for _, record := range breakdownRecords {
		dateTimeStart := record.DateTimeStart.Unix() * 1000
		if record.DateTimeStart.Before(dateFrom) {
			dateTimeStart = dateFrom.Unix() * 1000
		}
		dateTimeEnd := record.DateTimeEnd.Time.Unix() * 1000
		if record.DateTimeEnd.Time.IsZero() {
			dateTimeEnd = dateTo.Unix() * 1000
		}
		oneData := TerminalData{
			FromDate: dateTimeStart,
			ToDate:   dateTimeEnd,
			Note:     record.Note,
		}
		statesByIdSync.RLock()
		oneData.Color = cachedStatesById[poweroff].Color
		statesByIdSync.RUnlock()
		breakdownByIdSync.RLock()
		oneData.Information = cachedBreakdownsById[uint(record.BreakdownID)].Name
		breakdownByIdSync.RUnlock()
		breakdownData = append(breakdownData, oneData)
	}
	return breakdownData
}

func downloadChartDowntimeData(db *gorm.DB, dateTo time.Time, dateFrom time.Time, workplaceName string) []TerminalData {
	var downtimeData []TerminalData
	var downtimeRecords []database.DowntimeRecord
	workplacesByNameSync.RLock()
	workplaceId := cachedWorkplacesByName[workplaceName].ID
	workplacesByNameSync.RUnlock()
	db.Where("date_time_start <= ? and date_time_end >= ?", dateTo, dateFrom).Where("workplace_id = ?", workplaceId).Or("date_time_start <= ? and date_time_end is null", dateTo).Where("workplace_id = ?", workplaceId).Or("date_time_start <= ? and date_time_end >= ?", dateFrom, dateTo).Where("workplace_id = ?", workplaceId).Find(&downtimeRecords)
	for _, record := range downtimeRecords {
		dateTimeStart := record.DateTimeStart.Unix() * 1000
		if record.DateTimeStart.Before(dateFrom) {
			dateTimeStart = dateFrom.Unix() * 1000
		}
		dateTimeEnd := record.DateTimeEnd.Time.Unix() * 1000
		if record.DateTimeEnd.Time.IsZero() {
			dateTimeEnd = dateTo.Unix() * 1000
		}
		oneData := TerminalData{
			FromDate: dateTimeStart,
			ToDate:   dateTimeEnd,
			Note:     record.Note,
		}
		statesByIdSync.RLock()
		oneData.Color = cachedStatesById[downtime].Color
		statesByIdSync.RUnlock()
		downtimesByIdSync.RLock()
		oneData.Information = cachedDowntimesById[uint(record.DowntimeID)].Name
		downtimesByIdSync.RUnlock()
		downtimeData = append(downtimeData, oneData)
	}
	return downtimeData
}

func downloadChartOrderData(db *gorm.DB, dateTo time.Time, dateFrom time.Time, workplaceName string) []TerminalData {
	var orderData []TerminalData
	var orderRecords []database.OrderRecord
	workplacesByNameSync.RLock()
	workplaceId := cachedWorkplacesByName[workplaceName].ID
	workplacesByNameSync.RUnlock()
	db.Where("date_time_start <= ? and date_time_end >= ?", dateTo, dateFrom).Where("workplace_id = ?", workplaceId).Or("date_time_start <= ? and date_time_end is null", dateTo).Where("workplace_id = ?", workplaceId).Or("date_time_start <= ? and date_time_end >= ?", dateFrom, dateTo).Where("workplace_id = ?", workplaceId).Order("id").Find(&orderRecords)
	for _, record := range orderRecords {
		ordersByIdSync.RLock()
		productId := int(cachedOrdersById[uint(record.OrderID)].ProductID.Int32)
		ordersByIdSync.RUnlock()
		dateTimeStart := record.DateTimeStart.Unix() * 1000
		if record.DateTimeStart.Before(dateFrom) {
			dateTimeStart = dateFrom.Unix() * 1000
		}
		dateTimeEnd := record.DateTimeEnd.Time.Unix() * 1000
		if record.DateTimeEnd.Time.IsZero() {
			dateTimeEnd = dateTo.Unix() * 1000
		}
		oneData := TerminalData{
			FromDate: dateTimeStart,
			ToDate:   dateTimeEnd,
			Note:     record.Note,
		}
		statesByIdSync.RLock()
		oneData.Color = cachedStatesById[production].Color
		statesByIdSync.RUnlock()
		ordersByIdSync.RLock()
		orderName := cachedOrdersById[uint(record.OrderID)].Name
		ordersByIdSync.RUnlock()
		operationsByIdSync.RLock()
		operationName := cachedOperationsById[uint(record.OperationID)].Name
		operationsByIdSync.RUnlock()
		productsByIdSync.RLock()
		productName := cachedProductsById[uint(productId)].Name
		productsByIdSync.RUnlock()
		oneData.Information = orderName + ", " + operationName + ", " + productName
		orderData = append(orderData, oneData)
	}
	return orderData
}
