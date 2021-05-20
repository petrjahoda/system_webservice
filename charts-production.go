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

func processProductionChart(writer http.ResponseWriter, workplaceName string, dateFrom time.Time, dateTo time.Time, email string, chartName string) {
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
	allWorkplacePorts := cachedWorkplacePorts[workplaceName]
	for _, port := range allWorkplacePorts {
		if port.StateID.Int32 == 1 {
			var digitalData []database.DevicePortDigitalRecord
			db.Select("date_time, data").Where("date_time >= ?", dateFrom).Where("date_time <= ?", dateTo).Where("device_port_id = ?", port.DevicePortID).Order("id asc").Find(&digitalData)
			var portData PortData
			portData.PortName = "ID" + strconv.Itoa(int(port.ID)) + ": " + port.Name
			portData.PortColor = cachedDevicePortsColorsById[int(port.ID)]
			initialCounter := 1
			for _, data := range digitalData {
				if data.Data == 1 {
					var initialData Data
					initialData.Time = data.DateTime.Unix()
					initialData.Value = float32(initialCounter)
					initialCounter++
					portData.DigitalData = append(portData.DigitalData, initialData)
				}
			}
			initialCounter--
			var initialData Data
			initialData.Time = dateTo.Unix()
			initialData.Value = float32(initialCounter)
			portData.DigitalData = append(portData.DigitalData, initialData)
			digitalOutputData = append(digitalOutputData, portData)
		}
	}
	responseData.ChartData = digitalOutputData
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
	responseData.Locale = cachedUsersByEmail[email].Locale
	responseData.Result = "INF: Digital chart data downloaded from database in " + time.Since(timer).String()
	responseData.Type = chartName
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(responseData)
	logInfo("CHARTS-DIGITAL", "Digital chart data processed in "+time.Since(timer).String())
}
