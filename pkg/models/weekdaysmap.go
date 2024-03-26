package models

import (
	"strings"
	"time"
)

var weekDaysMap = map[string]int{
	"ΔΕΥΤΕΡΑ":   1,
	"ΤΡΙΤΗ":     2,
	"ΤΕΤΑΡΤΗ":   3,
	"ΠΕΜΠΤΗ":    4,
	"ΠΑΡΑΣΚΕΥΗ": 5,
	"ΣΑΒΒΑΤΟ":   6,
	"ΚΥΡΙΑΚΗ":   0,
}

func GetWeekDayIndex(weekDay string) int {
	for key := range weekDaysMap {
		if strings.Contains(weekDay, key) {
			return weekDaysMap[key]
		}
	}

	return -1
}

func GetWeekRange() string {
	// Get date of previous Monday
	startDate := time.Now()
	for startDate.Weekday() != time.Monday {
		startDate = startDate.AddDate(0, 0, -1)
	}

	// Get date of next Friday relative to the previous Monday
	endDate := startDate.AddDate(0, 0, 4)

	return startDate.Format("02/01") + "-" + endDate.Format("02/01")
}
