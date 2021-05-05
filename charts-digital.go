package main

import (
	"encoding/json"
	"github.com/petrjahoda/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

func processDigitalData(writer http.ResponseWriter, workplaceName string, dateFrom time.Time, dateTo time.Time, email string, chartName string) {
	timer := time.Now()
	logInfo("CHARTS-DIGITAL", "Processing digital chart data started")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("CHARTS-DIGITAL", "Problem opening database: "+err.Error())
		var responseData ChartDataPageOutput
		responseData.Result = "ERR: Problem opening database, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("CHARTS-DIGITAL", "Processing data ended")
		return
	}
	var responseData ChartDataPageOutput
	var digitalOutputData []PortData
	allWorkplacePorts := cachedWorkplaceDevicePorts[workplaceName]
	for _, port := range allWorkplacePorts {
		if port.DevicePortTypeID == 1 {
			var digitalData []database.DevicePortDigitalRecord
			db.Select("date_time, data").Where("date_time >= ?", dateFrom).Where("date_time <= ?", dateTo).Where("device_port_id = ?", port.ID).Order("id asc").Find(&digitalData)
			var portData PortData
			portData.PortName = "ID" + strconv.Itoa(int(port.ID)) + ": " + port.Name + " (" + port.Unit + ")"
			portData.PortColor = cachedDevicePortsColorsById[int(port.ID)]
			if len(digitalData) > 0 {
				var initialData Data
				initialData.Time = dateFrom.Unix()
				if digitalData[0].Data == 0 {
					initialData.Value = float32(1)
				} else {
					initialData.Value = float32(0)
				}
				portData.DigitalData = append(portData.DigitalData, initialData)
			}
			lastValue := 0
			for _, data := range digitalData {
				var initialData Data
				initialData.Time = data.DateTime.Unix()
				initialData.Value = float32(data.Data)
				lastValue = data.Data
				portData.DigitalData = append(portData.DigitalData, initialData)
			}
			var initialData Data
			initialData.Time = dateTo.Unix()
			if lastValue == 0 {
				initialData.Value = float32(1)
			} else {
				initialData.Value = float32(0)
			}
			digitalOutputData = append(digitalOutputData, portData)
		}
	}
	responseData.ChartData = digitalOutputData
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
	responseData.Result = "INF: Digital chart data downloaded from database in " + time.Since(timer).String()
	responseData.Type = chartName
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(responseData)
	logInfo("CHARTS-DIGITAL", "Digital chart data processed in "+time.Since(timer).String())
}
