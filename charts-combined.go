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

func processCombinedChart(writer http.ResponseWriter, workplaceName string, dateFrom time.Time, dateTo time.Time, email string, chartName string) {
	timer := time.Now()
	logInfo("CHARTS-COMBINED", "Processing combined chart data started")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("CHARTS-COMBINED", "Problem opening database: "+err.Error())
		var responseData ChartDataPageOutput
		responseData.Result = "ERR: Problem opening database, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("CHARTS-COMBINED", "Processing data ended")
		return
	}
	var responseData ChartDataPageOutput
	var combinedOutputData []PortData
	workplacePortsSync.RLock()
	allWorkplacePorts := cachedWorkplacePorts[workplaceName]
	workplacePortsSync.RUnlock()
	for _, port := range allWorkplacePorts {
		if port.StateID.Int32 == production {
			var digitalData []database.DevicePortDigitalRecord
			db.Select("date_time, data").Where("date_time >= ?", dateFrom).Where("date_time <= ?", dateTo).Where("device_port_id = ?", port.DevicePortID).Order("date_time").Order("id").Find(&digitalData)
			var portData PortData
			portData.PortType = "digital"
			devicePortsByIdSync.RLock()
			workplaceDevicePort := cachedDevicePortsById[uint(port.DevicePortID)]
			devicePortsByIdSync.RUnlock()
			portData.PortName = "ID" + strconv.Itoa(int(workplaceDevicePort.ID)) + ": " + workplaceDevicePort.Name + " (" + workplaceDevicePort.Unit + ")"
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
			combinedOutputData = append(combinedOutputData, portData)
		}
		if port.StateID.Int32 == poweroff {
			analogTimeDifference := 20.0
			var analogData []database.DevicePortAnalogRecord
			db.Select("date_time, data").Where("date_time >= ?", dateFrom).Where("date_time <= ?", dateTo).Where("device_port_id = ?", port.ID).Order("date_time").Order("id").Find(&analogData)
			var portData PortData
			portData.PortType = "analog"
			devicePortsByIdSync.RLock()
			workplaceDevicePort := cachedDevicePortsById[uint(port.DevicePortID)]
			devicePortsByIdSync.RUnlock()
			portData.PortName = "ID" + strconv.Itoa(int(workplaceDevicePort.ID)) + ": " + workplaceDevicePort.Name + " (" + workplaceDevicePort.Unit + ")"
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
			combinedOutputData = append(combinedOutputData, portData)
		}
	}
	workplacesByNameSync.RLock()
	workplaceId := cachedWorkplacesByName[workplaceName].ID
	workplacesByNameSync.RUnlock()
	var initialStateRecord database.StateRecord
	db.Select("id, date_time_start, state_id").Where("workplace_id = ?", workplaceId).Where("date_time_start < ?", dateFrom).Last(&initialStateRecord)
	var allStateRecords []database.StateRecord
	db.Select("date_time_start, state_id").Where("workplace_id = ?", workplaceId).Where("date_time_start < ?", dateTo).Where("id >= ?", initialStateRecord.ID).Order("date_time_start").Order("id").Find(&allStateRecords)
	var productionStateData PortData
	var downtimeStateData PortData
	var poweroffStateData PortData
	productionStateData.PortType = "state"
	downtimeStateData.PortType = "state"
	poweroffStateData.PortType = "state"
	productionStateData.PortName = getLocale(email, "production")
	statesByIdSync.RLock()
	productionStateData.PortColor = cachedStatesById[production].Color
	downtimeStateData.PortName = getLocale(email, "downtime")
	downtimeStateData.PortColor = cachedStatesById[downtime].Color
	poweroffStateData.PortName = getLocale(email, "poweroff")
	poweroffStateData.PortColor = cachedStatesById[poweroff].Color
	statesByIdSync.RUnlock()
	startLoop := true
	actualState := 0
	for _, record := range allStateRecords {
		if startLoop {
			var data Data
			data.Time = dateFrom.Unix()
			data.Value = 1.0
			if record.StateID == 1 {
				productionStateData.DigitalData = append(productionStateData.DigitalData, data)
				actualState = production
			} else if record.StateID == 2 {
				downtimeStateData.DigitalData = append(downtimeStateData.DigitalData, data)
				actualState = downtime
			} else {
				poweroffStateData.DigitalData = append(poweroffStateData.DigitalData, data)
				actualState = poweroff
			}
			startLoop = false
			continue
		}
		var data Data
		data.Time = record.DateTimeStart.Unix()
		data.Value = 0.0
		if actualState == production {
			productionStateData.DigitalData = append(productionStateData.DigitalData, data)
		} else if actualState == downtime {
			downtimeStateData.DigitalData = append(downtimeStateData.DigitalData, data)
		} else {
			poweroffStateData.DigitalData = append(poweroffStateData.DigitalData, data)
		}
		data.Time = record.DateTimeStart.Unix()
		data.Value = 1.0
		if record.StateID == production {
			productionStateData.DigitalData = append(productionStateData.DigitalData, data)
			actualState = production
		} else if record.StateID == downtime {
			downtimeStateData.DigitalData = append(downtimeStateData.DigitalData, data)
			actualState = downtime
		} else {
			poweroffStateData.DigitalData = append(poweroffStateData.DigitalData, data)
			actualState = poweroff
		}
	}

	var data Data
	data.Time = dateTo.Unix()
	data.Value = 0.0
	if actualState == production {
		productionStateData.DigitalData = append(productionStateData.DigitalData, data)
	} else if actualState == downtime {
		downtimeStateData.DigitalData = append(downtimeStateData.DigitalData, data)
	} else {
		poweroffStateData.DigitalData = append(poweroffStateData.DigitalData, data)
	}
	combinedOutputData = append(combinedOutputData, productionStateData)
	combinedOutputData = append(combinedOutputData, downtimeStateData)
	combinedOutputData = append(combinedOutputData, poweroffStateData)
	responseData.ChartData = combinedOutputData
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
	responseData.Result = "INF: Combined chart data downloaded from database in " + time.Since(timer).String()
	responseData.Type = chartName
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(responseData)
	logInfo("CHARTS-COMBINED", "Combined chart data processed in "+time.Since(timer).String())
}
