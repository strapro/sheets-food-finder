package models

import "strings"

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
