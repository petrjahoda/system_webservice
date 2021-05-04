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
	allWorkplacePorts := cachedWorkplaceDevicePorts[workplaceName]
	for _, port := range allWorkplacePorts {
		if port.DevicePortTypeID == 2 {
			var analogData []database.DevicePortAnalogRecord
			db.Select("date_time, data").Where("date_time >= ?", dateFrom).Where("date_time <= ?", dateTo).Where("device_port_id = ?", port.ID).Order("id asc").Find(&analogData)
			var portData PortData
			portData.PortName = "ID" + strconv.Itoa(int(port.ID)) + ": " + port.Name + " (" + port.Unit + ")"
			portData.PortColor = cachedDevicePortsColorsById[int(port.ID)]
			date := dateFrom
			for _, data := range analogData {
				for data.DateTime.Sub(date).Seconds() > 20 {
					var initialData Data
					initialData.Time = date.Unix()
					initialData.Value = math.MinInt16
					portData.PortData = append(portData.PortData, initialData)
					date = date.Add(10 * time.Second)
				}
				var initialData Data
				initialData.Time = data.DateTime.Unix()
				initialData.Value = data.Data
				date = data.DateTime
				portData.PortData = append(portData.PortData, initialData)
			}
			for dateTo.Sub(date).Seconds() > 20 {
				var initialData Data
				initialData.Time = date.Unix()
				initialData.Value = math.MinInt16
				portData.PortData = append(portData.PortData, initialData)
				date = date.Add(10 * time.Second)

			}
			analogOutputData = append(analogOutputData, portData)
		}
	}
	responseData.AnalogData = analogOutputData
	var terminalData []TerminalData
	var orderRecords []database.OrderRecord
	db.Where("date_time_start <= ? and date_time_end >= ?", dateTo, dateFrom).Where("workplace_id = ?", cachedWorkplacesByName[workplaceName].ID).Or("date_time_start <= ? and date_time_end is null", dateTo).Where("workplace_id = ?", cachedWorkplacesByName[workplaceName].ID).Or("date_time_start <= ? and date_time_end >= ?", dateFrom, dateTo).Where("workplace_id = ?", cachedWorkplacesByName[workplaceName].ID).Find(&orderRecords)
	for _, record := range orderRecords {
		productId := int(cachedOrdersById[uint(record.OrderID)].ProductID.Int32)
		dateTimeEnd := time.Now().Unix() * 1000
		if !record.DateTimeEnd.Time.IsZero() {
			dateTimeEnd = record.DateTimeEnd.Time.Unix() * 1000
		}
		oneOrderData := TerminalData{
			Name:          "production",
			Color:         "green",
			FromDate:      record.DateTimeStart.Unix() * 1000,
			ToDate:        dateTimeEnd,
			DataName:      cachedOrdersById[uint(record.OrderID)].Name,
			OperationName: cachedOperationsById[uint(record.OperationID)].Name,
			ProductName:   cachedProductsById[uint(productId)].Name,
			AverageCycle:  record.AverageCycle,
			CountOk:       record.CountOk,
			CountNok:      record.CountNok,
			Note:          record.Note,
		}
		terminalData = append(terminalData, oneOrderData)
	}
	var downtimeRecords []database.DowntimeRecord
	db.Where("date_time_start <= ? and date_time_end >= ?", dateTo, dateFrom).Where("workplace_id = ?", cachedWorkplacesByName[workplaceName].ID).Or("date_time_start <= ? and date_time_end is null", dateTo).Where("workplace_id = ?", cachedWorkplacesByName[workplaceName].ID).Or("date_time_start <= ? and date_time_end >= ?", dateFrom, dateTo).Where("workplace_id = ?", cachedWorkplacesByName[workplaceName].ID).Find(&downtimeRecords)
	for _, record := range downtimeRecords {
		dateTimeEnd := time.Now().Unix() * 1000
		if !record.DateTimeEnd.Time.IsZero() {
			dateTimeEnd = record.DateTimeEnd.Time.Unix() * 1000
		}
		oneOrderData := TerminalData{
			Name:     "downtime",
			Color:    "orange",
			FromDate: record.DateTimeStart.Unix() * 1000,
			ToDate:   dateTimeEnd,
			DataName: cachedDowntimesById[uint(record.DowntimeID)].Name,
			Note:     record.Note,
		}
		terminalData = append(terminalData, oneOrderData)
	}
	responseData.OrderData = terminalData
	responseData.Locale = cachedUsersByEmail[email].Locale
	responseData.Result = "INF: Analog chart data downloaded from database in " + time.Since(timer).String()
	responseData.Type = chartName
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(responseData)
	logInfo("CHARTS-ANALOG", "Analog chart data processed in "+time.Since(timer).String())
}
