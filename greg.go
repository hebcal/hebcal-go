package hebcal

import (
	"errors"
	"time"
)

// 1-based month lengths
var mlenStd = [13]int{0, 31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
var mlenLeap = [13]int{0, 31, 29, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

func IsGregLeapYear(year int) bool {
	return (year%4 == 0) && (year%100 != 0 || year%400 == 0)
}

func GregorianToRD(year int, month time.Month, day int) (int, error) {
	if year == 0 {
		return 0, errors.New("invalid Gregorian year")
	}
	var monthOffset int
	if month <= time.February {
		monthOffset = 0
	} else if IsGregLeapYear(year) {
		monthOffset = -1
	} else {
		monthOffset = -2
	}
	dayOfYear := ((367*int(month) - 362) / 12) + monthOffset + day
	var prevYear int
	if year >= 1 {
		prevYear = year - 1
	} else {
		prevYear = year + 1
		var yearLen int
		if IsGregLeapYear(year) {
			yearLen = 366
		} else {
			yearLen = 365
		}
		dayOfYear = -1 * (yearLen - dayOfYear)
	}
	rataDie := 365*prevYear + /* + days in prior years */
		(prevYear / 4) - /* + Julian Leap years */
		(prevYear / 100) + /* - century years */
		(prevYear / 400) + /* + Gregorian leap years */
		dayOfYear /* days this year */
	return rataDie, nil
}

/* like math.intAbs() but for ints */
func intAbs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

/*
func negativeRDtoGregorian(rataDie int, year int, d3 int) (int, time.Month, int) {
	mar1, _ := GregorianToRD(year, time.March, 1)
	var correction int
	if mar1 <= rataDie {
		correction = 0
	} else if IsLeapYear(year) {
		correction = 1
	} else {
		correction = 2
	}
	jan1, _ := GregorianToRD(year, time.January, 1)
	priorDays := rataDie - jan1
	month := (12*(priorDays+correction) + 373) / 367
	month1st, _ := GregorianToRD(year, time.Month(month), 1)
	day := rataDie - month1st + 1
	fmt.Printf("rataDie=%d, year=%d, jan1=%d, mar1=%d, priorDays=%d, correction=%d, month=%d, month1st=%d, day=%d, d3=%d\n", rataDie, year, jan1, mar1, priorDays, correction, month, month1st, day, d3)
	return year, time.Month(month), day
}
*/

func RDtoGregorian(rataDie int) (int, time.Month, int) {
	d0 := rataDie - 1
	n400 := d0 / 146097
	d1 := d0 % 146097
	n100 := d1 / 36524
	d2 := d1 % 36524
	n4 := d2 / 1461
	d3 := d2 % 1461
	n1 := d3 / 365
	year := 400*n400 + 100*n100 + 4*n4 + n1
	if n100 == 4 || n1 == 4 {
		return year, time.December, 31
	}
	// if rataDie <= 0 {
	// 	return negativeRDtoGregorian(rataDie, year-1, d3)
	// }
	year++
	var day = (d3 % 365) + 1
	var month = 1
	var mlen [13]int
	if IsGregLeapYear(year) {
		mlen = mlenLeap
	} else {
		mlen = mlenStd
	}
	for numDays := mlen[month]; numDays < day; numDays = mlen[month] {
		day -= numDays
		month++
	}
	return year, time.Month(month), day
}