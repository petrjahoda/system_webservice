package main

import (
	"encoding/json"
	"fmt"
	"github.com/petrjahoda/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

func loadDowntimesStatistics(writer http.ResponseWriter, dateFrom time.Time, dateTo time.Time, email string) {
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

	userWebSettingsSync.RLock()
	workplaceNames := cachedUserWebSettings[email]["statistics-selected-workplaces"]
	downtimeNames := cachedUserWebSettings[email]["statistics-selected-types-downtimes"]
	userNames := cachedUserWebSettings[email]["statistics-selected-users"]
	userWebSettingsSync.RUnlock()
	workplaceIds := getWorkplaceIds(workplaceNames)
	downtimeIds := getDowntimeIds(downtimeNames)

	userIds := getUserIds(userNames)
	if workplaceIds == "" {
		if downtimeNames == "" {
			if userNames == "" {
				db.Where("date_time_start <= ? and date_time_end >= ?", dateTo, dateFrom).Or("date_time_start <= ? and date_time_end is null", dateTo).Or("date_time_start <= ? and date_time_end >= ?", dateFrom, dateTo).Order("date_time_start desc").Find(&downtimeRecords)
			} else {
				db.Where("date_time_start <= ? and date_time_end >= ?", dateTo, dateFrom).Where(userIds).Or("date_time_start <= ? and date_time_end is null", dateTo).Where(userIds).Or("date_time_start <= ? and date_time_end >= ?", dateFrom, dateTo).Where(userIds).Order("date_time_start desc").Find(&downtimeRecords)
			}
		} else {
			if userNames == "" {
				db.Where("date_time_start <= ? and date_time_end >= ?", dateTo, dateFrom).Where(downtimeIds).Or("date_time_start <= ? and date_time_end is null", dateTo).Where(downtimeIds).Or("date_time_start <= ? and date_time_end >= ?", dateFrom, dateTo).Where(downtimeIds).Order("date_time_start desc").Find(&downtimeRecords)
			} else {
				db.Where("date_time_start <= ? and date_time_end >= ?", dateTo, dateFrom).Where(userIds).Where(downtimeIds).Or("date_time_start <= ? and date_time_end is null", dateTo).Where(userIds).Where(downtimeIds).Or("date_time_start <= ? and date_time_end >= ?", dateFrom, dateTo).Where(userIds).Where(downtimeIds).Order("date_time_start desc").Find(&downtimeRecords)
			}
		}
	} else {
		if downtimeNames == "" {
			if userNames == "" {
				db.Where("date_time_start <= ? and date_time_end >= ?", dateTo, dateFrom).Where(workplaceIds).Or("date_time_start <= ? and date_time_end is null", dateTo).Where(workplaceIds).Or("date_time_start <= ? and date_time_end >= ?", dateFrom, dateTo).Where(workplaceIds).Order("date_time_start desc").Find(&downtimeRecords)
			} else {
				db.Where("date_time_start <= ? and date_time_end >= ?", dateTo, dateFrom).Where(userIds).Where(workplaceIds).Or("date_time_start <= ? and date_time_end is null", dateTo).Where(userIds).Where(workplaceIds).Or("date_time_start <= ? and date_time_end >= ?", dateFrom, dateTo).Where(userIds).Where(workplaceIds).Order("date_time_start desc").Find(&downtimeRecords)
			}
		} else {
			if userNames == "" {
				db.Where("date_time_start <= ? and date_time_end >= ?", dateTo, dateFrom).Where(workplaceIds).Where(downtimeIds).Or("date_time_start <= ? and date_time_end is null", dateTo).Where(workplaceIds).Where(downtimeIds).Or("date_time_start <= ? and date_time_end >= ?", dateFrom, dateTo).Where(workplaceIds).Where(downtimeIds).Order("date_time_start desc").Find(&downtimeRecords)
			} else {
				db.Where("date_time_start <= ? and date_time_end >= ?", dateTo, dateFrom).Where(userIds).Where(workplaceIds).Where(downtimeIds).Or("date_time_start <= ? and date_time_end is null", dateTo).Where(userIds).Where(workplaceIds).Where(downtimeIds).Or("date_time_start <= ? and date_time_end >= ?", dateFrom, dateTo).Where(userIds).Where(workplaceIds).Where(downtimeIds).Order("date_time_start desc").Find(&downtimeRecords)
			}
		}
	}
	downtimeDataByDowntime := map[string]time.Duration{}
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
	}
	var responseData StatisticsDataOutput
	for downtimeName, duration := range downtimeDataByDowntime {
		responseData.SelectionChartData = append(responseData.SelectionChartData, downtimeName)
		responseData.SelectionChartValue = append(responseData.SelectionChartValue, duration.Seconds())
		responseData.SelectionChartText = append(responseData.SelectionChartText, duration.Round(time.Second).String())
	}
	downtimeDataByWorkplace := map[string]time.Duration{}
	for _, downtimeRecord := range downtimeRecords {
		if downtimeRecord.DateTimeEnd.Time.IsZero() {
			workplacesByIdSync.RLock()
			downtimeDataByWorkplace[cachedWorkplacesById[uint(downtimeRecord.WorkplaceID)].Name] = downtimeDataByWorkplace[cachedWorkplacesById[uint(downtimeRecord.WorkplaceID)].Name] + time.Now().Sub(downtimeRecord.DateTimeStart)
			workplacesByIdSync.RUnlock()
		} else {
			workplacesByIdSync.RLock()
			downtimeDataByWorkplace[cachedWorkplacesById[uint(downtimeRecord.WorkplaceID)].Name] = downtimeDataByWorkplace[cachedWorkplacesById[uint(downtimeRecord.WorkplaceID)].Name] + downtimeRecord.DateTimeEnd.Time.Sub(downtimeRecord.DateTimeStart)
			workplacesByIdSync.RUnlock()
		}
	}
	for workplaceName, duration := range downtimeDataByWorkplace {
		responseData.WorkplaceChartData = append(responseData.WorkplaceChartData, workplaceName)
		responseData.WorkplaceChartValue = append(responseData.WorkplaceChartValue, duration.Seconds())
		responseData.WorkplaceChartText = append(responseData.WorkplaceChartText, duration.Round(time.Second).String())
	}

	downtimeDataByStart := map[int]int{}
	for _, record := range downtimeRecords {
		if record.DateTimeEnd.Valid {
			downtimeDataByStart[int(record.DateTimeEnd.Time.Sub(record.DateTimeStart).Round(1*time.Hour).Hours())] = downtimeDataByStart[int(record.DateTimeEnd.Time.Sub(record.DateTimeStart).Round(1*time.Hour).Hours())] + 1
		}

	}
	for duration, count := range downtimeDataByStart {
		responseData.TimeChartData = append(responseData.TimeChartData, strconv.Itoa(duration))
		responseData.TimeChartValue = append(responseData.TimeChartValue, float64(count))
		responseData.TimeChartText = append(responseData.TimeChartText, "+"+strconv.Itoa(duration)+"m: "+strconv.Itoa(count)+"x")
	}
	downtimeDataByUser := map[uint]time.Duration{}
	for _, record := range downtimeRecords {
		if record.DateTimeEnd.Time.IsZero() {
			downtimeDataByUser[uint(record.UserID.Int32)] = downtimeDataByUser[uint(record.UserID.Int32)] + time.Now().Sub(record.DateTimeStart)
		} else {
			downtimeDataByUser[uint(record.UserID.Int32)] = downtimeDataByUser[uint(record.UserID.Int32)] + record.DateTimeEnd.Time.Sub(record.DateTimeStart)
		}

	}
	for userId, duration := range downtimeDataByUser {
		fmt.Println(userId, duration)
		usersByIdSync.RLock()
		if userId == 0 {
			responseData.UsersChartData = append(responseData.UsersChartData, "-")
		} else {
			responseData.UsersChartData = append(responseData.UsersChartData, cachedUsersById[userId].FirstName+" "+cachedUsersById[userId].SecondName)
		}
		usersByIdSync.RUnlock()
		responseData.UsersChartValue = append(responseData.UsersChartValue, duration.Seconds())
		responseData.UsersChartText = append(responseData.UsersChartText, duration.Round(time.Second).String())

	}
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(responseData)
	logInfo("STATISTICS", "Downtime statistics loaded in "+time.Since(timer).String())
}
