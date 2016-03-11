package gannettApi

import (
	"fmt"
)

/*
	Format a year, month, day, year, hours, minutes, and seconds into a date string
	for querying the Gannett Api

	FormatAsDateSting(2014, 10, 1, 0, 0, 0) == 2014-10-01T00:00:00Z

	For more info
		https://confluence.gannett.com/pages/viewpage.action?title=Search+v4+Recipes&spaceKey=GDPDW#Searchv4Recipes-FilterbyDateRange
*/
func FormatAsDateString(year, month, day, hour, minute, second int) string {

	return fmt.Sprintf("%s-%s-%sT%s:%s:%sZ",
		GetZeroPaddedString(year),
		GetZeroPaddedString(month),
		GetZeroPaddedString(day),
		GetZeroPaddedString(hour),
		GetZeroPaddedString(minute),
		GetZeroPaddedString(second))
}

/*
	Zero pad a number (useful for date strings, e.g. 1 -> 01 )
*/
func GetZeroPaddedString(i int) string {
	if i < 10 {
		return fmt.Sprintf("0%d", i)
	} else {
		return fmt.Sprintf("%d", i)
	}
}
