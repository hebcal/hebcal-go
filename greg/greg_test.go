package greg

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGreg2RD(t *testing.T) {
	assert := assert.New(t)
	rataDie, _ := ToRD(1995, time.December, 17)
	assert.Equal(728644, rataDie)
	rataDie, _ = ToRD(1888, time.December, 31)
	assert.Equal(689578, rataDie)
	rataDie, _ = ToRD(2005, time.April, 2)
	assert.Equal(732038, rataDie)
}

func TestGreg2RDEarlyCE(t *testing.T) {
	assert := assert.New(t)
	rataDie, _ := ToRD(88, time.December, 30)
	assert.Equal(32141, rataDie)
	rataDie, _ = ToRD(1, time.January, 1)
	assert.Equal(1, rataDie)
}

func TestGreg2RDNegative(t *testing.T) {
	assert := assert.New(t)
	rataDie, _ := ToRD(-1, time.December, 31)
	assert.Equal(0, rataDie)
	rataDie, _ = ToRD(-1, time.December, 30)
	assert.Equal(-1, rataDie)
	rataDie, _ = ToRD(-1, time.December, 29)
	assert.Equal(-2, rataDie)
	rataDie, _ = ToRD(-1, time.November, 13)
	assert.Equal(-48, rataDie)
	rataDie, _ = ToRD(-1, time.October, 31)
	assert.Equal(-61, rataDie)
	rataDie, _ = ToRD(-1, time.October, 1)
	assert.Equal(-91, rataDie)
	rataDie, _ = ToRD(-1, time.March, 1)
	assert.Equal(-305, rataDie)
	rataDie, _ = ToRD(-1, time.February, 28)
	assert.Equal(-306, rataDie)
	rataDie, _ = ToRD(-1, time.February, 27)
	assert.Equal(-307, rataDie)
	rataDie, _ = ToRD(-1, time.February, 1)
	assert.Equal(-333, rataDie)
	rataDie, _ = ToRD(-1, time.January, 31)
	assert.Equal(-334, rataDie)
	rataDie, _ = ToRD(-1, time.January, 16)
	assert.Equal(-349, rataDie)
	rataDie, _ = ToRD(-1, time.January, 1)
	assert.Equal(-364, rataDie)
}

func TestGreg2RDNegative2(t *testing.T) {
	assert := assert.New(t)
	rataDie, _ := ToRD(-2, time.December, 31)
	assert.Equal(-365, rataDie)
	rataDie, _ = ToRD(-2, time.December, 1)
	assert.Equal(-395, rataDie)
	rataDie, _ = ToRD(-2, time.January, 1)
	assert.Equal(-729, rataDie)
	rataDie, _ = ToRD(-3, time.January, 1)
	assert.Equal(-1094, rataDie)
	rataDie, _ = ToRD(-4, time.January, 1)
	assert.Equal(-1460, rataDie)
	rataDie, _ = ToRD(-100, time.December, 20)
	assert.Equal(-36170, rataDie)
	rataDie, _ = ToRD(-1000, time.June, 15)
	assert.Equal(-365076, rataDie)
}

func TestRD2Greg(t *testing.T) {
	assert := assert.New(t)
	year, month, day := FromRD(737553)
	assert.Equal(2020, year)
	assert.Equal(time.May, month)
	assert.Equal(8, day)
	year2, month2, day2 := FromRD(689578)
	assert.Equal(1888, year2)
	assert.Equal(time.December, month2)
	assert.Equal(31, day2)
	gy, gm, gd := FromRD(732038)
	assert.Equal(2005, gy)
	assert.Equal(time.April, gm)
	assert.Equal(2, gd)
}

func TestRD2Greg88ce(t *testing.T) {
	assert := assert.New(t)
	var year int
	var month time.Month
	var day int
	year, month, day = FromRD(32141)
	assert.Equal(88, year)
	assert.Equal(time.December, month)
	assert.Equal(30, day)
	year, month, day = FromRD(32142)
	assert.Equal(88, year)
	assert.Equal(time.December, month)
	assert.Equal(31, day)
	year, month, day = FromRD(32143)
	assert.Equal(89, year)
	assert.Equal(time.January, month)
	assert.Equal(1, day)
}

func TestRD2Greg1ce(t *testing.T) {
	assert := assert.New(t)
	year, month, day := FromRD(1)
	assert.Equal(1, year)
	assert.Equal(time.January, month)
	assert.Equal(1, day)
}

/*
func TestRD2GregNegative(t *testing.T) {
	assert := assert.New(t)
	var year int
	var month time.Month
	var day int
	year, month, day = RDtoGregorian(0)
	assert.Equal(-1, year)
	assert.Equal(time.December, month)
	assert.Equal(31, day)
	year, month, day = RDtoGregorian(-1)
	assert.Equal(-1, year)
	assert.Equal(time.December, month)
	assert.Equal(30, day)
	year, month, day = RDtoGregorian(-2)
	assert.Equal(-1, year)
	assert.Equal(time.December, month)
	assert.Equal(29, day)
	year, month, day = RDtoGregorian(-48)
	assert.Equal(-1, year)
	assert.Equal(time.November, month)
	assert.Equal(13, day)
	year, month, day = RDtoGregorian(-61)
	assert.Equal(-1, year)
	assert.Equal(time.October, month)
	assert.Equal(31, day)

	year, month, day = RDtoGregorian(-91)
	assert.Equal(-1, year)
	assert.Equal(time.October, month)
	assert.Equal(1, day)
	year, month, day = RDtoGregorian(-305)
	assert.Equal(-1, year)
	assert.Equal(time.March, month)
	assert.Equal(1, day)

	year, month, day = RDtoGregorian(-349)
	assert.Equal(-1, year)
	assert.Equal(time.January, month)
	assert.Equal(16, day)
	year, month, day = RDtoGregorian(-364)
	assert.Equal(-1, year)
	assert.Equal(time.January, month)
	assert.Equal(1, day)
	year, month, day = RDtoGregorian(-365)
	assert.Equal(-2, year)
	assert.Equal(time.December, month)
	assert.Equal(31, day)
	year, month, day = RDtoGregorian(-36535)
	assert.Equal(-100, year)
	assert.Equal(time.December, month)
	assert.Equal(20, day)
}
*/

func ExampleDaysIn() {
	days := DaysIn(time.February, 2004)
	fmt.Println(days)
	// Output: 29
}

func ExampleDateToRD() {
	t := time.Date(2014, time.February, 19, 0, 0, 0, 0, time.UTC)
	rataDie := DateToRD(t)
	fmt.Println(rataDie)
	// Output: 735283
}

func ExampleToRD() {
	rataDie, _ := ToRD(1995, time.December, 17)
	fmt.Println(rataDie)
	// Output: 728644
}

func ExampleFromRD() {
	year, month, day := FromRD(737553)
	fmt.Println(year, month, day)
	// Output: 2020 May 8
}