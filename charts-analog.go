package main

import (
	"encoding/json"
	"fmt"
	"github.com/petrjahoda/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
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
		responseData.Result = "nok: " + err.Error()
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
			db.Where("date_time >= ?", dateFrom).Where("date_time <= ?", dateTo).Where("device_port_id = ?", port.ID).Order("id asc").Find(&analogData)
			var portData PortData
			portData.PortName = port.Name
			date := dateFrom
			for _, data := range analogData {
				for data.DateTime.Sub(date).Seconds() > 20 {
					var initialData Data
					initialData.Time = date.Unix()
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
				portData.PortData = append(portData.PortData, initialData)
				date = date.Add(10 * time.Second)

			}
			fmt.Println(len(portData.PortData))
			analogOutputData = append(analogOutputData, portData)
		}
	}
	responseData.AnalogData = analogOutputData
	responseData.Result = "ok"
	responseData.Type = chartName
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(responseData)
	logInfo("CHARTS-ANALOG", "Analog chart data processed in "+time.Since(timer).String())
}
