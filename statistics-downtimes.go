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

func loadDowntimesStatistics(writer http.ResponseWriter, workplaceIds string, dateFrom time.Time, dateTo time.Time, email string) {
	timer := time.Now()
	logInfo("STATISTICS-DOWNTIMES", "Loading downtimes statistics")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("STATISTICS-DOWNTIMES", "Problem opening database: "+err.Error())
		var responseData StatisticsDataOutput
		responseData.Result = "ERR: Problem opening database, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("STATISTICS-DOWNTIMES", "Loading downtime statistics ended")
		return
	}
	var downtimeRecords []database.DowntimeRecord
	var userRecords []database.UserRecord
	if workplaceIds == "" {
		db.Where("date_time_start <= ? and date_time_end >= ?", dateTo, dateFrom).Or("date_time_start <= ? and date_time_end is null", dateTo).Or("date_time_start <= ? and date_time_end >= ?", dateFrom, dateTo).Order("date_time_start desc").Find(&downtimeRecords)
		db.Where("date_time_start <= ? and date_time_end >= ?", dateTo, dateFrom).Or("date_time_start <= ? and date_time_end is null", dateTo).Or("date_time_start <= ? and date_time_end >= ?", dateFrom, dateTo).Order("date_time_start desc").Find(&userRecords)
	} else {
		db.Where("date_time_start <= ? and date_time_end >= ?", dateTo, dateFrom).Where(workplaceIds).Or("date_time_start <= ? and date_time_end is null", dateTo).Where(workplaceIds).Or("date_time_start <= ? and date_time_end >= ?", dateFrom, dateTo).Where(workplaceIds).Order("date_time_start desc").Find(&downtimeRecords)
		db.Where("date_time_start <= ? and date_time_end >= ?", dateTo, dateFrom).Where(workplaceIds).Or("date_time_start <= ? and date_time_end is null", dateTo).Where(workplaceIds).Or("date_time_start <= ? and date_time_end >= ?", dateFrom, dateTo).Where(workplaceIds).Order("date_time_start desc").Find(&userRecords)
	}
	downtimeDataByDowntime := map[string]time.Duration{}
	downtimeDataByUser := map[uint]time.Duration{}
	for _, downtimeRecord := range downtimeRecords {
		if downtimeRecord.DateTimeEnd.Time.IsZero() {
			downtimesByIdSync.RLock()
			downtimeDataByDowntime[cachedDowntimesById[uint(downtimeRecord.DowntimeID)].Name] = downtimeDataByDowntime[cachedDowntimesById[uint(downtimeRecord.DowntimeID)].Name] + time.Now().Sub(downtimeRecord.DateTimeStart)
			downtimesByIdSync.RUnlock()
		} else {
			downtimesByIdSync.RLock()
			downtimeDataByDowntime[cachedDowntimesById[uint(downtimeRecord.DowntimeID)].Name] = downtimeDataByDowntime[cachedDowntimesById[uint(downtimeRecord.DowntimeID)].Name] + downtimeRecord.DateTimeEnd.Time.Sub(downtimeRecord.DateTimeStart)
			downtimesByIdSync.RUnlock()
		}
		for _, userRecord := range userRecords {
			if userRecord.WorkplaceID == downtimeRecord.WorkplaceID {
				if downtimeRecord.DateTimeEnd.Time.IsZero() && userRecord.DateTimeEnd.Time.IsZero() {
					downtimeDataByUser[uint(userRecord.UserID)] = downtimeDataByUser[uint(userRecord.UserID)] + time.Now().Sub(downtimeRecord.DateTimeStart)
					break
				} else if downtimeRecord.DateTimeStart.Before(userRecord.DateTimeStart) && userRecord.DateTimeStart.Before(downtimeRecord.DateTimeEnd.Time) && !downtimeRecord.DateTimeEnd.Time.IsZero() {
					downtimeDataByUser[uint(userRecord.UserID)] = downtimeDataByUser[uint(userRecord.UserID)] + downtimeRecord.DateTimeEnd.Time.Sub(downtimeRecord.DateTimeStart)
					break
				} else if downtimeRecord.DateTimeStart.Before(userRecord.DateTimeEnd.Time) && userRecord.DateTimeEnd.Time.Before(downtimeRecord.DateTimeEnd.Time) && !downtimeRecord.DateTimeEnd.Time.IsZero() {
					downtimeDataByUser[uint(userRecord.UserID)] = downtimeDataByUser[uint(userRecord.UserID)] + downtimeRecord.DateTimeEnd.Time.Sub(downtimeRecord.DateTimeStart)
					break
				} else if downtimeRecord.DateTimeStart.Before(userRecord.DateTimeStart) && userRecord.DateTimeEnd.Time.Before(downtimeRecord.DateTimeEnd.Time) && !downtimeRecord.DateTimeEnd.Time.IsZero() {
					downtimeDataByUser[uint(userRecord.UserID)] = downtimeDataByUser[uint(userRecord.UserID)] + downtimeRecord.DateTimeEnd.Time.Sub(downtimeRecord.DateTimeStart)
					break
				} else if userRecord.DateTimeStart.Before(downtimeRecord.DateTimeStart) && downtimeRecord.DateTimeEnd.Time.Before(userRecord.DateTimeEnd.Time) && !downtimeRecord.DateTimeEnd.Time.IsZero() {
					downtimeDataByUser[uint(userRecord.UserID)] = downtimeDataByUser[uint(userRecord.UserID)] + downtimeRecord.DateTimeEnd.Time.Sub(downtimeRecord.DateTimeStart)
					break
				}
			}
		}
	}
	var responseData StatisticsDataOutput
	for downtimeName, duration := range downtimeDataByDowntime {
		responseData.DurationChartData = append(responseData.DurationChartData, downtimeName)
		responseData.DurationChartValue = append(responseData.DurationChartValue, duration.Seconds())
		responseData.DurationChartText = append(responseData.DurationChartText, duration.Round(time.Second).String())
	}
	downtimeDataByWorkplace := map[string]time.Duration{}
	for _, downtimeRecord := range downtimeRecords {
		if downtimeRecord.DateTimeEnd.Time.IsZero() {
			workplacesByIdSync.RLock()
			downtimeDataByWorkplace[cachedWorkplacesById[uint(downtimeRecord.DowntimeID)].Name] = downtimeDataByWorkplace[cachedWorkplacesById[uint(downtimeRecord.DowntimeID)].Name] + time.Now().Sub(downtimeRecord.DateTimeStart)
			workplacesByIdSync.RUnlock()
		} else {
			workplacesByIdSync.RLock()
			downtimeDataByWorkplace[cachedWorkplacesById[uint(downtimeRecord.DowntimeID)].Name] = downtimeDataByWorkplace[cachedWorkplacesById[uint(downtimeRecord.DowntimeID)].Name] + downtimeRecord.DateTimeEnd.Time.Sub(downtimeRecord.DateTimeStart)
			workplacesByIdSync.RUnlock()
		}
	}
	for workplaceName, duration := range downtimeDataByWorkplace {
		responseData.WorkplaceChartData = append(responseData.WorkplaceChartData, workplaceName)
		responseData.WorkplaceChartValue = append(responseData.WorkplaceChartValue, duration.Seconds())
		responseData.WorkplaceChartText = append(responseData.WorkplaceChartText, duration.Round(time.Second).String())
	}

	downtimeDataByStart := map[int]time.Duration{}
	for _, record := range downtimeRecords {
		downtimeDataByStart[record.DateTimeStart.Hour()] = downtimeDataByStart[record.DateTimeStart.Hour()] + 1
	}
	for hour, count := range downtimeDataByStart {
		responseData.StartChartData = append(responseData.StartChartData, strconv.Itoa(hour))
		responseData.StartChartValue = append(responseData.StartChartValue, float64(count))
		responseData.StartChartText = append(responseData.StartChartText, strconv.Itoa(hour)+":00 - "+strconv.Itoa(hour+1)+":00")
	}
	for userId, duration := range downtimeDataByUser {
		usersByIdSync.RLock()
		responseData.UsersChartData = append(responseData.UsersChartData, cachedUsersById[userId].FirstName+" "+cachedUsersById[userId].SecondName)
		usersByIdSync.RUnlock()
		responseData.UsersChartValue = append(responseData.UsersChartValue, duration.Seconds())
		responseData.UsersChartText = append(responseData.UsersChartText, duration.Round(time.Second).String())

	}
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(responseData)
	logInfo("STATISTICS", "Downtime statistics loaded in "+time.Since(timer).String())
}
