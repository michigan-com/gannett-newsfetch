package lib

import "time"

func SameTime(date1, date2 time.Time) bool {
	date1 = date1.UTC()
	date2 = date2.UTC()
	return date1.Year() == date2.Year() &&
		date1.Month() == date2.Month() &&
		date1.Day() == date2.Day() &&
		date1.Hour() == date2.Hour() &&
		date1.Minute() == date2.Minute() &&
		date1.Second() == date2.Second()
}

/*
	Given a string date, return the date. If anything goes wrong, return time.Now()
*/
func GannettDateStringToDate(dateString string) time.Time {
	// https://golang.org/src/time/format.go#L64
	// Idk, a regular date string wasnt working, cause why would it
	date, err := time.Parse(time.RFC3339Nano, dateString)
	if err != nil {
		return time.Now()
	}
	return date.UTC()
}
