package utils

import "time"

func GetDateFromStr(t string) time.Time {

	tm, _ := time.Parse(time.RFC3339, t)
	return tm
}
