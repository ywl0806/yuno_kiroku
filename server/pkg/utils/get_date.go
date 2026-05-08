package utils

import "time"

func GetDateFromStr(t string) time.Time {

	tm, _ := time.Parse(time.RFC3339, t)
	return tm
}

func GetFirstDayOfMonth(year int, month int) time.Time {
	return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
}

func GetLastDayOfMonth(year int, month int) time.Time {
	return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local).AddDate(0, 1, 0).Add(-time.Second)
}
