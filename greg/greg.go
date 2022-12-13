// Hebcal's greg package converts between Gregorian dates
// and R.D. (Rata Die) day numbers.
package greg

// Hebcal - A Jewish Calendar Generator
// Copyright (c) 2022 Michael J. Radwin
// Derived from original C version, Copyright (C) 1994-2004 Danny Sadinoff
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

import (
	"math"
	"time"
)

// 1-based month lengths
var monthLen = [13]int{0, 31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

// DaysIn returns the number of days in the Gregorian month.
func DaysIn(m time.Month, year int) int {
	if m == time.February && IsLeapYear(year) {
		return 29
	}
	return monthLen[m]
}

// Returns true if the Gregorian year is a leap year.
func IsLeapYear(year int) bool {
	return (year%4 == 0) && (year%100 != 0 || year%400 == 0)
}

// Converts Gregorian date to absolute R.D. (Rata Die) days.
// Hours, minutes and seconds are ignored
func DateToRD(t time.Time) int64 {
	year, month, day := t.Date()
	abs := ToRD(year, month, day)
	return abs
}

func monthOffset(year int, month time.Month) int {
	if month <= time.February {
		return 0
	} else if IsLeapYear(year) {
		return -1
	} else {
		return -2
	}
}

// Converts Gregorian date to absolute R.D. (Rata Die) days.
//
// Panics if Gregorian year is 0.
func ToRD(year int, month time.Month, day int) int64 {
	py := int64(year - 1)
	abs := 365*py +
		quotient(py, 4) -
		quotient(py, 100) +
		quotient(py, 400) +
		quotient((367*int64(month)-362), 12) +
		int64(monthOffset(year, month)) + int64(day)
	return abs
}

func mod(x, y int64) int64 {
	X := float64(x)
	Y := float64(y)
	return int64(X - Y*math.Floor(X/Y))
}

func quotient(x, y int64) int64 {
	return int64(math.Floor(float64(x) / float64(y)))
}

func yearFromFixed(rataDie int64) int {
	l0 := int64(rataDie) - 1
	n400 := quotient(l0, 146097)
	d1 := mod(l0, 146097)
	n100 := quotient(d1, 36524)
	d2 := mod(d1, 36524)
	n4 := quotient(d2, 1461)
	d3 := mod(d2, 1461)
	n1 := quotient(d3, 365)
	year := 400*n400 + 100*n100 + 4*n4 + n1
	yy := int(year)
	if n100 == 4 || n1 == 4 {
		return yy
	}
	return yy + 1
}

/*
Converts from Rata Die (R.D. number) to Gregorian date.

See the footnote on page 384 of “Calendrical Calculations, Part II:
Three Historical Calendars” by E. M. Reingold,  N. Dershowitz, and S. M.
Clamen, Software--Practice and Experience, Volume 23, Number 4
(April, 1993), pages 383-404 for an explanation.
*/
func FromRD(rataDie int64) (int, time.Month, int) {
	year := yearFromFixed(rataDie)
	var correction int64
	if rataDie < ToRD(year, time.March, 1) {
		correction = 0
	} else if IsLeapYear(year) {
		correction = 1
	} else {
		correction = 2
	}
	priorDays := rataDie - ToRD(year, time.January, 1)
	month := quotient(12*(priorDays+correction)+373, 367)
	day := rataDie - ToRD(year, time.Month(month), 1) + 1
	return year, time.Month(month), int(day)
}
