package main

import (
	"encoding/json"
	"github.com/petrjahoda/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

func loadBreakdownStatistics(writer http.ResponseWriter, dateFrom time.Time, dateTo time.Time, email string) {
	timer := time.Now()
	logInfo("STATISTICS-BREAKDOWNS", "Loading breakdowns statistics")
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("STATISTICS-BREAKDOWNS", "Problem opening database: "+err.Error())
		var responseData StatisticsDataOutput
		responseData.Result = "ERR: Problem opening database, " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("STATISTICS-BREAKDOWNS", "Loading breakdowns statistics ended")
		return
	}
	breakdownRecords := downloadBreakdownRecords(email, db, dateTo, dateFrom)
	var responseData StatisticsDataOutput
	breakdownDataByBreakdownSorted := processBreakdownDataByBreakdown(breakdownRecords)
	for _, record := range breakdownDataByBreakdownSorted {
		responseData.SelectionChartData = append(responseData.SelectionChartData, record.name)
		responseData.SelectionChartValue = append(responseData.SelectionChartValue, record.duration.Seconds())
		responseData.SelectionChartText = append(responseData.SelectionChartText, record.duration.Round(time.Second).String())
	}

	breakdownDataByWorkplaceSorted := processBreakdownDataByWorkplace(breakdownRecords)
	for _, data := range breakdownDataByWorkplaceSorted {
		responseData.WorkplaceChartData = append(responseData.WorkplaceChartData, data.name)
		responseData.WorkplaceChartValue = append(responseData.WorkplaceChartValue, data.duration.Seconds())
		responseData.WorkplaceChartText = append(responseData.WorkplaceChartText, data.duration.Round(time.Second).String())
	}

	breakdownDataByDurationSorted := processBreakdownDataByDuration(breakdownRecords)
	for _, data := range breakdownDataByDurationSorted {
		responseData.TimeChartData = append(responseData.TimeChartData, data.name)
		responseData.TimeChartValue = append(responseData.TimeChartValue, data.duration.Seconds())
		responseData.TimeChartText = append(responseData.TimeChartText, data.duration.Round(time.Second).String())
	}

	breakdownDataByUserSorted := processBreakdownDataByUser(breakdownRecords)
	for _, data := range breakdownDataByUserSorted {
		responseData.UsersChartData = append(responseData.UsersChartData, data.name)
		responseData.UsersChartValue = append(responseData.UsersChartValue, data.duration.Seconds())
		responseData.UsersChartText = append(responseData.UsersChartText, data.duration.Round(time.Second).String())
	}

	breakdownDataByStartSorted := processBreakdownDataByDate(err, breakdownRecords, dateFrom, dateTo)
	for _, data := range breakdownDataByStartSorted {
		responseData.DaysChartData = append(responseData.DaysChartData, data.name)
		responseData.DaysChartValue = append(responseData.DaysChartValue, data.duration.Seconds())
		responseData.DaysChartText = append(responseData.DaysChartText, data.duration.Round(time.Second).String())
	}
	usersByEmailSync.RLock()
	responseData.Locale = cachedUsersByEmail[email].Locale
	usersByEmailSync.RUnlock()
	responseData.CalendarChartLocale = getLocale(email, "statistics-breakdowns-calendar")
	responseData.FirstUpperChartLocale = getLocale(email, "statistics-workplaces")
	responseData.SecondUpperChartLocale = getLocale(email, "statistics-selection-breakdowns")
	responseData.ThirdUpperChartLocale = getLocale(email, "statistics-users")
	responseData.FourthUpperChartLocale = getLocale(email, "statistics-duration")
	statesByIdSync.RLock()
	responseData.Color = cachedStatesById[poweroff].Color
	statesByIdSync.RUnlock()
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(responseData)
	logInfo("STATISTICS-BREAKDOWNS", "Breakdowns statistics loaded in "+time.Since(timer).String())
}

func processBreakdownDataByDate(err error, breakdownRecords []database.BreakdownRecord, dateFrom time.Time, dateTo time.Time) []tempData {
	loc, err := time.LoadLocation(cachedLocation)
	breakdownDataByStart := map[string]time.Duration{}
	for _, record := range breakdownRecords {
		if record.DateTimeStart.YearDay() == record.DateTimeEnd.Time.YearDay() {
			breakdownDataByStart[record.DateTimeStart.In(loc).Format("2006-01-02")] = breakdownDataByStart[record.DateTimeStart.In(loc).Format("2006-01-02")] + record.DateTimeEnd.Time.Sub(record.DateTimeStart)
		} else if record.DateTimeEnd.Time.IsZero() {
			record.DateTimeEnd.Time = time.Now()
			tempDate := record.DateTimeStart
			for tempDate.YearDay() != record.DateTimeEnd.Time.YearDay() {
				tempDateNoon := time.Date(tempDate.Year(), tempDate.Month(), tempDate.Day()+1, 0, 0, 0, 0, loc)
				breakdownDataByStart[tempDate.In(loc).Format("2006-01-02")] = breakdownDataByStart[tempDate.In(loc).Format("2006-01-02")] + tempDateNoon.In(loc).Sub(tempDate.In(loc))
				tempDate = time.Date(tempDate.Year(), tempDate.Month(), tempDate.Day()+1, 0, 0, 0, 0, loc)
			}
			breakdownDataByStart[tempDate.In(loc).Format("2006-01-02")] = breakdownDataByStart[tempDate.In(loc).Format("2006-01-02")] + time.Now().In(loc).Sub(tempDate.In(loc))

		} else {
			tempDate := record.DateTimeStart
			for tempDate.YearDay() != record.DateTimeEnd.Time.YearDay() {
				tempDateNoon := time.Date(tempDate.Year(), tempDate.Month(), tempDate.Day()+1, 0, 0, 0, 0, loc)
				breakdownDataByStart[tempDate.In(loc).Format("2006-01-02")] = breakdownDataByStart[tempDate.In(loc).Format("2006-01-02")] + tempDateNoon.In(loc).Sub(tempDate.In(loc))
				tempDate = time.Date(tempDate.Year(), tempDate.Month(), tempDate.Day()+1, 0, 0, 0, 0, loc)
			}
			breakdownDataByStart[tempDate.In(loc).Format("2006-01-02")] = breakdownDataByStart[tempDate.In(loc).Format("2006-01-02")] + record.DateTimeEnd.Time.In(loc).Sub(tempDate.In(loc))
		}

	}
	for dateFrom.In(loc).YearDay() != dateTo.In(loc).YearDay() {
		breakdownDataByStart[dateFrom.In(loc).Format("2006-01-02")] = breakdownDataByStart[dateFrom.In(loc).Format("2006-01-02")] + 0*time.Second
		dateFrom = dateFrom.Add(24 * time.Hour)
	}
	var breakdownDataByStartSorted []tempData
	for date, duration := range breakdownDataByStart {
		breakdownDataByStartSorted = append(breakdownDataByStartSorted, tempData{date, duration})
	}
	sort.Slice(breakdownDataByStartSorted, func(i, j int) bool {
		return breakdownDataByStartSorted[i].name < breakdownDataByStartSorted[j].name
	})
	return breakdownDataByStartSorted
}

func processBreakdownDataByUser(breakdownRecords []database.BreakdownRecord) []tempData {
	breakdownDataByUser := map[string]time.Duration{}
	for _, record := range breakdownRecords {
		if record.DateTimeEnd.Time.IsZero() {
			usersByIdSync.RLock()
			breakdownDataByUser[cachedUsersById[uint(record.UserID)].FirstName+" "+cachedUsersById[uint(record.UserID)].SecondName] = breakdownDataByUser[cachedUsersById[uint(record.UserID)].FirstName+" "+cachedUsersById[uint(record.UserID)].SecondName] + time.Now().Sub(record.DateTimeStart)
			usersByIdSync.RUnlock()
		} else {
			usersByIdSync.RLock()
			breakdownDataByUser[cachedUsersById[uint(record.UserID)].FirstName+" "+cachedUsersById[uint(record.UserID)].SecondName] = breakdownDataByUser[cachedUsersById[uint(record.UserID)].FirstName+" "+cachedUsersById[uint(record.UserID)].SecondName] + record.DateTimeEnd.Time.Sub(record.DateTimeStart)
			usersByIdSync.RUnlock()
		}
	}
	var breakdownDataByUserSorted []tempData
	for userName, duration := range breakdownDataByUser {
		breakdownDataByUserSorted = append(breakdownDataByUserSorted, tempData{userName, duration})
	}
	sort.Slice(breakdownDataByUserSorted, func(i, j int) bool {
		return breakdownDataByUserSorted[i].duration < breakdownDataByUserSorted[j].duration
	})
	return breakdownDataByUserSorted
}

func processBreakdownDataByDuration(breakdownRecords []database.BreakdownRecord) []tempData {
	breakdownDataByStart := map[string]time.Duration{}
	for _, record := range breakdownRecords {
		var duration time.Duration
		if record.DateTimeEnd.Time.IsZero() {
			duration = time.Now().Sub(record.DateTimeStart)
		} else {
			duration = record.DateTimeEnd.Time.Sub(record.DateTimeStart)
		}
		if duration <= 5*time.Minute {
			breakdownDataByStart["0m<5m"] = breakdownDataByStart["0m<5m"] + duration
		} else if duration > 5*time.Minute && duration <= 15*time.Minute {
			breakdownDataByStart["5m<15m"] = breakdownDataByStart["5m<15m"] + duration
		} else if duration > 15*time.Minute && duration <= 30*time.Minute {
			breakdownDataByStart["15m<30m"] = breakdownDataByStart["15m<30m"] + duration
		} else if duration > 30*time.Minute && duration <= 60*time.Minute {
			breakdownDataByStart["30m<60m"] = breakdownDataByStart["30m<60m"] + duration
		} else if duration > 1*time.Hour && duration <= 2*time.Hour {
			breakdownDataByStart["1h<2h"] = breakdownDataByStart["1h<2h"] + duration
		} else if duration > 2*time.Hour && duration <= 4*time.Hour {
			breakdownDataByStart["2h<4h"] = breakdownDataByStart["2h<4h"] + duration
		} else if duration > 4*time.Hour && duration <= 8*time.Hour {
			breakdownDataByStart["4h<8h"] = breakdownDataByStart["4h<8h"] + duration
		} else if duration > 8*time.Hour && duration <= 16*time.Hour {
			breakdownDataByStart["8h<16h"] = breakdownDataByStart["8h<16h"] + duration
		} else if duration > 16*time.Hour && duration <= 32*time.Hour {
			breakdownDataByStart["16h<32h"] = breakdownDataByStart["16h<32h"] + duration
		} else {
			breakdownDataByStart["32h<++"] = breakdownDataByStart["32h<++"] + duration
		}
	}
	var breakdownDataByStartSorted []tempData
	for name, duration := range breakdownDataByStart {
		breakdownDataByStartSorted = append(breakdownDataByStartSorted, tempData{name, duration})
	}
	sort.Slice(breakdownDataByStartSorted, func(i, j int) bool {
		return breakdownDataByStartSorted[i].duration < breakdownDataByStartSorted[j].duration
	})
	return breakdownDataByStartSorted
}

func processBreakdownDataByWorkplace(breakdownRecords []database.BreakdownRecord) []tempData {
	breakdownDataByWorkplace := map[string]time.Duration{}
	for _, breakdownRecord := range breakdownRecords {
		if breakdownRecord.DateTimeEnd.Time.IsZero() {
			workplacesByIdSync.RLock()
			breakdownDataByWorkplace[cachedWorkplacesById[uint(breakdownRecord.WorkplaceID)].Name] = breakdownDataByWorkplace[cachedWorkplacesById[uint(breakdownRecord.WorkplaceID)].Name] + time.Now().Sub(breakdownRecord.DateTimeStart)
			workplacesByIdSync.RUnlock()
		} else {
			workplacesByIdSync.RLock()
			breakdownDataByWorkplace[cachedWorkplacesById[uint(breakdownRecord.WorkplaceID)].Name] = breakdownDataByWorkplace[cachedWorkplacesById[uint(breakdownRecord.WorkplaceID)].Name] + breakdownRecord.DateTimeEnd.Time.Sub(breakdownRecord.DateTimeStart)
			workplacesByIdSync.RUnlock()
		}
	}

	var brwakdownDataByWorkplaceSorted []tempData
	for workplaceName, duration := range breakdownDataByWorkplace {
		brwakdownDataByWorkplaceSorted = append(brwakdownDataByWorkplaceSorted, tempData{workplaceName, duration})
	}
	sort.Slice(brwakdownDataByWorkplaceSorted, func(i, j int) bool {
		return brwakdownDataByWorkplaceSorted[i].duration < brwakdownDataByWorkplaceSorted[j].duration
	})
	return brwakdownDataByWorkplaceSorted
}

func processBreakdownDataByBreakdown(breakdownRecords []database.BreakdownRecord) []tempData {
	breakdownDataByBreakdown := map[string]time.Duration{}
	for _, breakdownRecord := range breakdownRecords {
		if breakdownRecord.DateTimeEnd.Time.IsZero() {
			breakdownByIdSync.RLock()
			breakdownDataByBreakdown[cachedBreakdownsById[uint(breakdownRecord.BreakdownID)].Name] = breakdownDataByBreakdown[cachedBreakdownsById[uint(breakdownRecord.BreakdownID)].Name] + time.Now().Sub(breakdownRecord.DateTimeStart)
			breakdownByIdSync.RUnlock()
		} else {
			breakdownByIdSync.RLock()
			breakdownDataByBreakdown[cachedBreakdownsById[uint(breakdownRecord.BreakdownID)].Name] = breakdownDataByBreakdown[cachedBreakdownsById[uint(breakdownRecord.BreakdownID)].Name] + breakdownRecord.DateTimeEnd.Time.Sub(breakdownRecord.DateTimeStart)
			breakdownByIdSync.RUnlock()
		}
	}
	var breakdownDataByBreakdownsSorted []tempData
	for breakdownName, duration := range breakdownDataByBreakdown {
		breakdownDataByBreakdownsSorted = append(breakdownDataByBreakdownsSorted, tempData{breakdownName, duration})
	}
	sort.Slice(breakdownDataByBreakdownsSorted, func(i, j int) bool {
		return breakdownDataByBreakdownsSorted[i].duration < breakdownDataByBreakdownsSorted[j].duration
	})
	return breakdownDataByBreakdownsSorted

}

func downloadBreakdownRecords(email string, db *gorm.DB, dateTo time.Time, dateFrom time.Time) []database.BreakdownRecord {
	var breakdownRecords []database.BreakdownRecord
	userWebSettingsSync.RLock()
	workplaceNames := cachedUserWebSettings[email]["statistics-selected-workplaces"]
	breakdownNames := cachedUserWebSettings[email]["statistics-selected-types-breakdowns"]
	userNames := cachedUserWebSettings[email]["statistics-selected-users"]
	userWebSettingsSync.RUnlock()
	workplaceIds := getWorkplaceIds(workplaceNames)
	breakdownIds := getBreakdownIds(breakdownNames)
	userIds := getUserIds(userNames)
	if workplaceIds == "" {
		if breakdownNames == "" {
			if userNames == "" {
				db.Where("date_time_start <= ? and date_time_end >= ?", dateTo, dateFrom).Or("date_time_start <= ? and date_time_end is null", dateTo).Or("date_time_start <= ? and date_time_end >= ?", dateFrom, dateTo).Order("date_time_start asc").Find(&breakdownRecords)
			} else {
				db.Where("date_time_start <= ? and date_time_end >= ?", dateTo, dateFrom).Where(userIds).Or("date_time_start <= ? and date_time_end is null", dateTo).Where(userIds).Or("date_time_start <= ? and date_time_end >= ?", dateFrom, dateTo).Where(userIds).Order("date_time_start asc").Find(&breakdownRecords)
			}
		} else {
			if userNames == "" {
				db.Where("date_time_start <= ? and date_time_end >= ?", dateTo, dateFrom).Where(breakdownIds).Or("date_time_start <= ? and date_time_end is null", dateTo).Where(breakdownIds).Or("date_time_start <= ? and date_time_end >= ?", dateFrom, dateTo).Where(breakdownIds).Order("date_time_start asc").Find(&breakdownRecords)
			} else {
				db.Where("date_time_start <= ? and date_time_end >= ?", dateTo, dateFrom).Where(userIds).Where(breakdownIds).Or("date_time_start <= ? and date_time_end is null", dateTo).Where(userIds).Where(breakdownIds).Or("date_time_start <= ? and date_time_end >= ?", dateFrom, dateTo).Where(userIds).Where(breakdownIds).Order("date_time_start asc").Find(&breakdownRecords)
			}
		}
	} else {
		if breakdownNames == "" {
			if userNames == "" {
				db.Where("date_time_start <= ? and date_time_end >= ?", dateTo, dateFrom).Where(workplaceIds).Or("date_time_start <= ? and date_time_end is null", dateTo).Where(workplaceIds).Or("date_time_start <= ? and date_time_end >= ?", dateFrom, dateTo).Where(workplaceIds).Order("date_time_start asc").Find(&breakdownRecords)
			} else {
				db.Where("date_time_start <= ? and date_time_end >= ?", dateTo, dateFrom).Where(userIds).Where(workplaceIds).Or("date_time_start <= ? and date_time_end is null", dateTo).Where(userIds).Where(workplaceIds).Or("date_time_start <= ? and date_time_end >= ?", dateFrom, dateTo).Where(userIds).Where(workplaceIds).Order("date_time_start asc").Find(&breakdownRecords)
			}
		} else {
			if userNames == "" {
				db.Where("date_time_start <= ? and date_time_end >= ?", dateTo, dateFrom).Where(workplaceIds).Where(breakdownIds).Or("date_time_start <= ? and date_time_end is null", dateTo).Where(workplaceIds).Where(breakdownIds).Or("date_time_start <= ? and date_time_end >= ?", dateFrom, dateTo).Where(workplaceIds).Where(breakdownIds).Order("date_time_start asc").Find(&breakdownRecords)
			} else {
				db.Where("date_time_start <= ? and date_time_end >= ?", dateTo, dateFrom).Where(userIds).Where(workplaceIds).Where(breakdownIds).Or("date_time_start <= ? and date_time_end is null", dateTo).Where(userIds).Where(workplaceIds).Where(breakdownIds).Or("date_time_start <= ? and date_time_end >= ?", dateFrom, dateTo).Where(userIds).Where(workplaceIds).Where(breakdownIds).Order("date_time_start asc").Find(&breakdownRecords)
			}
		}
	}
	return breakdownRecords

}

func getBreakdownIds(breakdownNames string) string {
	if len(breakdownNames) == 0 {
		return ""
	}
	breakdownNamesAsArray := strings.Split(breakdownNames, ";")
	breakdownNamesSql := `name in ('`
	for _, user := range breakdownNamesAsArray {
		breakdownNamesSql += user + `','`
	}
	breakdownNamesSql = strings.TrimSuffix(breakdownNamesSql, `,'`)
	breakdownNamesSql += ")"
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("DATA", "Problem opening database: "+err.Error())
		return ""
	}
	var breakdowns []database.Breakdown
	db.Select("id").Where(breakdownNamesSql).Find(&breakdowns)
	breakdownIds := `breakdown_id in ('`
	for _, breakdown := range breakdowns {
		breakdownIds += strconv.Itoa(int(breakdown.ID)) + `','`
	}
	breakdownIds = strings.TrimSuffix(breakdownIds, `,'`)
	breakdownIds += ")"
	return breakdownIds

}
