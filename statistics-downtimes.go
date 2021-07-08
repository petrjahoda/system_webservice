package main

import (
	"encoding/json"
	"github.com/petrjahoda/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"sort"
	"time"
)

type tempData struct {
	name     string
	duration time.Duration
}

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
	downtimeRecords := downloadDowntimeRecords(email, db, dateTo, dateFrom)
	var responseData StatisticsDataOutput
	downtimeDataByDowntimeSorted := processDataByDowntime(downtimeRecords)
	for _, record := range downtimeDataByDowntimeSorted {
		responseData.SelectionChartData = append(responseData.SelectionChartData, record.name)
		responseData.SelectionChartValue = append(responseData.SelectionChartValue, record.duration.Seconds())
		responseData.SelectionChartText = append(responseData.SelectionChartText, record.duration.Round(time.Second).String())
	}

	downtimeDataByWorkplaceSorted := processDataByWorkplace(downtimeRecords)
	for _, data := range downtimeDataByWorkplaceSorted {
		responseData.WorkplaceChartData = append(responseData.WorkplaceChartData, data.name)
		responseData.WorkplaceChartValue = append(responseData.WorkplaceChartValue, data.duration.Seconds())
		responseData.WorkplaceChartText = append(responseData.WorkplaceChartText, data.duration.Round(time.Second).String())
	}

	downtimeDataByDurationSorted := processDataByDuration(downtimeRecords)
	for _, data := range downtimeDataByDurationSorted {
		responseData.TimeChartData = append(responseData.TimeChartData, data.name)
		responseData.TimeChartValue = append(responseData.TimeChartValue, data.duration.Seconds())
		responseData.TimeChartText = append(responseData.TimeChartText, data.name+": "+data.duration.String())
	}

	downtimeDataByUserSorted := processDataByUser(downtimeRecords)
	for _, data := range downtimeDataByUserSorted {
		responseData.UsersChartData = append(responseData.UsersChartData, data.name)
		responseData.UsersChartValue = append(responseData.UsersChartValue, data.duration.Seconds())
		responseData.UsersChartText = append(responseData.UsersChartText, data.duration.Round(time.Second).String())
	}

	downtimeDataByStartSorted := processDataByDate(err, downtimeRecords, dateFrom, dateTo)
	for _, data := range downtimeDataByStartSorted {
		responseData.DaysChartData = append(responseData.DaysChartData, data.name)
		responseData.DaysChartValue = append(responseData.DaysChartValue, data.duration.Seconds())
		responseData.DaysChartText = append(responseData.DaysChartText, data.duration.Round(time.Second).String())
	}
	usersByEmailSync.RLock()
	responseData.Locale = cachedUsersByEmail[email].Locale
	usersByEmailSync.RUnlock()
	responseData.CalendarChartLocale = getLocale(email, "statistics-downtimes-calendar")
	responseData.FirstUpperChartLocale = getLocale(email, "statistics-workplaces")
	responseData.SecondUpperChartLocale = getLocale(email, "statistics-selection-downtimes")
	responseData.ThirdUpperChartLocale = getLocale(email, "statistics-users")
	responseData.FourthUpperChartLocale = getLocale(email, "statistics-duration")
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(responseData)
	logInfo("STATISTICS-DOWNTIMES", "Downtime statistics loaded in "+time.Since(timer).String())
}

func downloadDowntimeRecords(email string, db *gorm.DB, dateTo time.Time, dateFrom time.Time) []database.DowntimeRecord {
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
				db.Where("date_time_start <= ? and date_time_end >= ?", dateTo, dateFrom).Or("date_time_start <= ? and date_time_end is null", dateTo).Or("date_time_start <= ? and date_time_end >= ?", dateFrom, dateTo).Order("date_time_start asc").Find(&downtimeRecords)
			} else {
				db.Where("date_time_start <= ? and date_time_end >= ?", dateTo, dateFrom).Where(userIds).Or("date_time_start <= ? and date_time_end is null", dateTo).Where(userIds).Or("date_time_start <= ? and date_time_end >= ?", dateFrom, dateTo).Where(userIds).Order("date_time_start asc").Find(&downtimeRecords)
			}
		} else {
			if userNames == "" {
				db.Where("date_time_start <= ? and date_time_end >= ?", dateTo, dateFrom).Where(downtimeIds).Or("date_time_start <= ? and date_time_end is null", dateTo).Where(downtimeIds).Or("date_time_start <= ? and date_time_end >= ?", dateFrom, dateTo).Where(downtimeIds).Order("date_time_start asc").Find(&downtimeRecords)
			} else {
				db.Where("date_time_start <= ? and date_time_end >= ?", dateTo, dateFrom).Where(userIds).Where(downtimeIds).Or("date_time_start <= ? and date_time_end is null", dateTo).Where(userIds).Where(downtimeIds).Or("date_time_start <= ? and date_time_end >= ?", dateFrom, dateTo).Where(userIds).Where(downtimeIds).Order("date_time_start asc").Find(&downtimeRecords)
			}
		}
	} else {
		if downtimeNames == "" {
			if userNames == "" {
				db.Where("date_time_start <= ? and date_time_end >= ?", dateTo, dateFrom).Where(workplaceIds).Or("date_time_start <= ? and date_time_end is null", dateTo).Where(workplaceIds).Or("date_time_start <= ? and date_time_end >= ?", dateFrom, dateTo).Where(workplaceIds).Order("date_time_start asc").Find(&downtimeRecords)
			} else {
				db.Where("date_time_start <= ? and date_time_end >= ?", dateTo, dateFrom).Where(userIds).Where(workplaceIds).Or("date_time_start <= ? and date_time_end is null", dateTo).Where(userIds).Where(workplaceIds).Or("date_time_start <= ? and date_time_end >= ?", dateFrom, dateTo).Where(userIds).Where(workplaceIds).Order("date_time_start asc").Find(&downtimeRecords)
			}
		} else {
			if userNames == "" {
				db.Where("date_time_start <= ? and date_time_end >= ?", dateTo, dateFrom).Where(workplaceIds).Where(downtimeIds).Or("date_time_start <= ? and date_time_end is null", dateTo).Where(workplaceIds).Where(downtimeIds).Or("date_time_start <= ? and date_time_end >= ?", dateFrom, dateTo).Where(workplaceIds).Where(downtimeIds).Order("date_time_start asc").Find(&downtimeRecords)
			} else {
				db.Where("date_time_start <= ? and date_time_end >= ?", dateTo, dateFrom).Where(userIds).Where(workplaceIds).Where(downtimeIds).Or("date_time_start <= ? and date_time_end is null", dateTo).Where(userIds).Where(workplaceIds).Where(downtimeIds).Or("date_time_start <= ? and date_time_end >= ?", dateFrom, dateTo).Where(userIds).Where(workplaceIds).Where(downtimeIds).Order("date_time_start asc").Find(&downtimeRecords)
			}
		}
	}
	return downtimeRecords
}

func processDataByDate(err error, downtimeRecords []database.DowntimeRecord, dateFrom time.Time, dateTo time.Time) []tempData {
	loc, err := time.LoadLocation(cachedLocation)
	downtimeDataByStart := map[string]time.Duration{}
	for _, record := range downtimeRecords {
		if record.DateTimeStart.YearDay() == record.DateTimeEnd.Time.YearDay() {
			downtimeDataByStart[record.DateTimeStart.In(loc).Format("2006-01-02")] = downtimeDataByStart[record.DateTimeStart.In(loc).Format("2006-01-02")] + record.DateTimeEnd.Time.Sub(record.DateTimeStart)
		} else if record.DateTimeEnd.Time.IsZero() {
			record.DateTimeEnd.Time = time.Now()
			tempDate := record.DateTimeStart
			for tempDate.YearDay() != record.DateTimeEnd.Time.YearDay() {
				tempDateNoon := time.Date(tempDate.Year(), tempDate.Month(), tempDate.Day()+1, 0, 0, 0, 0, loc)
				downtimeDataByStart[tempDate.In(loc).Format("2006-01-02")] = downtimeDataByStart[tempDate.In(loc).Format("2006-01-02")] + tempDateNoon.In(loc).Sub(tempDate.In(loc))
				tempDate = time.Date(tempDate.Year(), tempDate.Month(), tempDate.Day()+1, 0, 0, 0, 0, loc)
			}
			downtimeDataByStart[tempDate.In(loc).Format("2006-01-02")] = downtimeDataByStart[tempDate.In(loc).Format("2006-01-02")] + time.Now().In(loc).Sub(tempDate.In(loc))

		} else {
			tempDate := record.DateTimeStart
			for tempDate.YearDay() != record.DateTimeEnd.Time.YearDay() {
				tempDateNoon := time.Date(tempDate.Year(), tempDate.Month(), tempDate.Day()+1, 0, 0, 0, 0, loc)
				downtimeDataByStart[tempDate.In(loc).Format("2006-01-02")] = downtimeDataByStart[tempDate.In(loc).Format("2006-01-02")] + tempDateNoon.In(loc).Sub(tempDate.In(loc))
				tempDate = time.Date(tempDate.Year(), tempDate.Month(), tempDate.Day()+1, 0, 0, 0, 0, loc)
			}
			downtimeDataByStart[tempDate.In(loc).Format("2006-01-02")] = downtimeDataByStart[tempDate.In(loc).Format("2006-01-02")] + record.DateTimeEnd.Time.In(loc).Sub(tempDate.In(loc))
		}

	}
	for dateFrom.In(loc).YearDay() != dateTo.In(loc).YearDay() {
		downtimeDataByStart[dateFrom.In(loc).Format("2006-01-02")] = downtimeDataByStart[dateFrom.In(loc).Format("2006-01-02")] + 0*time.Second
		dateFrom = dateFrom.Add(24 * time.Hour)
	}
	var downtimeDataByStartSorted []tempData
	for date, duration := range downtimeDataByStart {
		downtimeDataByStartSorted = append(downtimeDataByStartSorted, tempData{date, duration})
	}
	sort.Slice(downtimeDataByStartSorted, func(i, j int) bool {
		return downtimeDataByStartSorted[i].name < downtimeDataByStartSorted[j].name
	})
	return downtimeDataByStartSorted
}

func processDataByUser(downtimeRecords []database.DowntimeRecord) []tempData {
	downtimeDataByUser := map[string]time.Duration{}
	for _, record := range downtimeRecords {
		if record.DateTimeEnd.Time.IsZero() {
			usersByIdSync.RLock()
			downtimeDataByUser[cachedUsersById[uint(record.UserID.Int32)].FirstName+" "+cachedUsersById[uint(record.UserID.Int32)].SecondName] = downtimeDataByUser[cachedUsersById[uint(record.UserID.Int32)].FirstName+" "+cachedUsersById[uint(record.UserID.Int32)].SecondName] + time.Now().Sub(record.DateTimeStart)
			usersByIdSync.RUnlock()
		} else {
			usersByIdSync.RLock()
			downtimeDataByUser[cachedUsersById[uint(record.UserID.Int32)].FirstName+" "+cachedUsersById[uint(record.UserID.Int32)].SecondName] = downtimeDataByUser[cachedUsersById[uint(record.UserID.Int32)].FirstName+" "+cachedUsersById[uint(record.UserID.Int32)].SecondName] + record.DateTimeEnd.Time.Sub(record.DateTimeStart)
			usersByIdSync.RUnlock()
		}
	}
	var downtimeDataByUserSorted []tempData
	for userName, duration := range downtimeDataByUser {
		downtimeDataByUserSorted = append(downtimeDataByUserSorted, tempData{userName, duration})
	}
	sort.Slice(downtimeDataByUserSorted, func(i, j int) bool {
		return downtimeDataByUserSorted[i].duration < downtimeDataByUserSorted[j].duration
	})
	return downtimeDataByUserSorted
}

func processDataByDuration(downtimeRecords []database.DowntimeRecord) []tempData {
	downtimeDataByStart := map[string]time.Duration{}
	for _, record := range downtimeRecords {
		var duration time.Duration
		if record.DateTimeEnd.Time.IsZero() {
			duration = time.Now().Sub(record.DateTimeStart)
		} else {
			duration = record.DateTimeEnd.Time.Sub(record.DateTimeStart)
		}
		if duration <= 5*time.Minute {
			downtimeDataByStart["0m<5m"] = downtimeDataByStart["0m<5m"] + duration
		} else if duration > 5*time.Minute && duration <= 15*time.Minute {
			downtimeDataByStart["5m<15m"] = downtimeDataByStart["5m<15m"] + duration
		} else if duration > 15*time.Minute && duration <= 30*time.Minute {
			downtimeDataByStart["15m<30m"] = downtimeDataByStart["15m<30m"] + duration
		} else if duration > 30*time.Minute && duration <= 60*time.Minute {
			downtimeDataByStart["30m<60m"] = downtimeDataByStart["30m<60m"] + duration
		} else if duration > 1*time.Hour && duration <= 2*time.Hour {
			downtimeDataByStart["1h<2h"] = downtimeDataByStart["1h<2h"] + duration
		} else if duration > 2*time.Hour && duration <= 4*time.Hour {
			downtimeDataByStart["2h<4h"] = downtimeDataByStart["2h<4h"] + duration
		} else if duration > 4*time.Hour && duration <= 8*time.Hour {
			downtimeDataByStart["4h<8h"] = downtimeDataByStart["4h<8h"] + duration
		} else if duration > 8*time.Hour && duration <= 16*time.Hour {
			downtimeDataByStart["8h<16h"] = downtimeDataByStart["8h<16h"] + duration
		} else if duration > 16*time.Hour && duration <= 32*time.Hour {
			downtimeDataByStart["16h<32h"] = downtimeDataByStart["16h<32h"] + duration
		} else {
			downtimeDataByStart["32h<++"] = downtimeDataByStart["32h<++"] + duration
		}
	}
	var downtimeDataByStartSorted []tempData
	for name, duration := range downtimeDataByStart {
		downtimeDataByStartSorted = append(downtimeDataByStartSorted, tempData{name, duration})
	}
	sort.Slice(downtimeDataByStartSorted, func(i, j int) bool {
		return downtimeDataByStartSorted[i].duration < downtimeDataByStartSorted[j].duration
	})
	return downtimeDataByStartSorted
}

func processDataByWorkplace(downtimeRecords []database.DowntimeRecord) []tempData {
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

	var downtimeDataByWorkplaceSorted []tempData
	for workplaceName, duration := range downtimeDataByWorkplace {
		downtimeDataByWorkplaceSorted = append(downtimeDataByWorkplaceSorted, tempData{workplaceName, duration})
	}
	sort.Slice(downtimeDataByWorkplaceSorted, func(i, j int) bool {
		return downtimeDataByWorkplaceSorted[i].duration < downtimeDataByWorkplaceSorted[j].duration
	})
	return downtimeDataByWorkplaceSorted
}

func processDataByDowntime(downtimeRecords []database.DowntimeRecord) []tempData {
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
	var downtimeDataByDowntimeSorted []tempData
	for downtimeName, duration := range downtimeDataByDowntime {
		downtimeDataByDowntimeSorted = append(downtimeDataByDowntimeSorted, tempData{downtimeName, duration})
	}
	sort.Slice(downtimeDataByDowntimeSorted, func(i, j int) bool {
		return downtimeDataByDowntimeSorted[i].duration < downtimeDataByDowntimeSorted[j].duration
	})
	return downtimeDataByDowntimeSorted
}
