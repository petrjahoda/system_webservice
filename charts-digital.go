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
	workplaceDevicePortsSync.RLock()
	allWorkplacePorts := cachedWorkplaceDevicePorts[workplaceName]
	workplaceDevicePortsSync.RUnlock()
	for _, port := range allWorkplacePorts {
		if port.DevicePortTypeID == digital {
			var digitalData []database.DevicePortDigitalRecord
			db.Select("date_time, data").Where("date_time >= ?", dateFrom).Where("date_time <= ?", dateTo).Where("device_port_id = ?", port.ID).Order("date_time").Order("id").Find(&digitalData)
			var portData PortData
			portData.PortName = "ID" + strconv.Itoa(int(port.ID)) + ": " + port.Name + " (" + port.Unit + ")"
			devicePortsColorsByIdSync.RLock()
			portColor := cachedDevicePortsColorsById[int(port.ID)]
			devicePortsColorsByIdSync.RUnlock()
			if portColor == "#000000" {
				portData.PortColor = ""
			}
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
				initialData.Value = float32(0)
			} else {
				initialData.Value = float32(1)
			}
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
	usersByEmailSync.RLock()
	responseData.Locale = cachedUsersByEmail[email].Locale
	usersByEmailSync.RUnlock()
	responseData.Result = "INF: Digital chart data downloaded from database in " + time.Since(timer).String()
	responseData.Type = chartName
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(responseData)
	logInfo("CHARTS-DIGITAL", "Digital chart data processed in "+time.Since(timer).String())
}
