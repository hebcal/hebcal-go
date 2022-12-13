package greg_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hebcal/hebcal-go/greg"
	"github.com/stretchr/testify/assert"
)

func TestGreg2RD(t *testing.T) {
	assert := assert.New(t)
	rataDie := greg.ToRD(1995, time.December, 17)
	assert.Equal(int64(728644), rataDie)
	rataDie = greg.ToRD(1888, time.December, 31)
	assert.Equal(int64(689578), rataDie)
	rataDie = greg.ToRD(2005, time.April, 2)
	assert.Equal(int64(732038), rataDie)
}

func TestGreg2RDEarlyCE(t *testing.T) {
	assert := assert.New(t)
	rataDie := greg.ToRD(88, time.December, 30)
	assert.Equal(int64(32141), rataDie)
	rataDie = greg.ToRD(1, time.January, 1)
	assert.Equal(int64(1), rataDie)
}

func TestGreg2RDNegative(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(int64(0), greg.ToRD(0, time.December, 31))
	assert.Equal(int64(-1), greg.ToRD(0, time.December, 30))
	assert.Equal(int64(-2), greg.ToRD(0, time.December, 29))
	assert.Equal(int64(-48), greg.ToRD(0, time.November, 13))
	assert.Equal(int64(-61), greg.ToRD(0, time.October, 31))
	assert.Equal(int64(-91), greg.ToRD(0, time.October, 1))
	assert.Equal(int64(-305), greg.ToRD(0, time.March, 1))
	assert.Equal(int64(-306), greg.ToRD(0, time.February, 29))
	assert.Equal(int64(-307), greg.ToRD(0, time.February, 28))
	assert.Equal(int64(-308), greg.ToRD(0, time.February, 27))
	assert.Equal(int64(-334), greg.ToRD(0, time.February, 1))
	assert.Equal(int64(-335), greg.ToRD(0, time.January, 31))
	assert.Equal(int64(-350), greg.ToRD(0, time.January, 16))
	assert.Equal(int64(-365), greg.ToRD(0, time.January, 1))
	assert.Equal(int64(-366), greg.ToRD(-1, time.December, 31))
}

func TestGreg2RDNegative2(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(int64(-396), greg.ToRD(-1, time.December, 1))
	assert.Equal(int64(-730), greg.ToRD(-1, time.January, 1))
	assert.Equal(int64(-1095), greg.ToRD(-2, time.January, 1))
	assert.Equal(int64(-1460), greg.ToRD(-3, time.January, 1))
	assert.Equal(int64(-36171), greg.ToRD(-99, time.December, 20))
	assert.Equal(int64(-365077), greg.ToRD(-999, time.June, 15))
	assert.Equal(int64(-36536), greg.ToRD(-100, time.December, 20))
}

func TestRD2Greg(t *testing.T) {
	assert := assert.New(t)
	year, month, day := greg.FromRD(737553)
	assert.Equal(2020, year)
	assert.Equal(time.May, month)
	assert.Equal(8, day)
	year2, month2, day2 := greg.FromRD(689578)
	assert.Equal(1888, year2)
	assert.Equal(time.December, month2)
	assert.Equal(31, day2)
	gy, gm, gd := greg.FromRD(732038)
	assert.Equal(2005, gy)
	assert.Equal(time.April, gm)
	assert.Equal(2, gd)
}

func TestRD2Greg88ce(t *testing.T) {
	assert := assert.New(t)
	var year int
	var month time.Month
	var day int
	year, month, day = greg.FromRD(32141)
	assert.Equal(88, year)
	assert.Equal(time.December, month)
	assert.Equal(30, day)
	year, month, day = greg.FromRD(32142)
	assert.Equal(88, year)
	assert.Equal(time.December, month)
	assert.Equal(31, day)
	year, month, day = greg.FromRD(32143)
	assert.Equal(89, year)
	assert.Equal(time.January, month)
	assert.Equal(1, day)
}

func TestRD2Greg1ce(t *testing.T) {
	assert := assert.New(t)
	year, month, day := greg.FromRD(1)
	assert.Equal(1, year)
	assert.Equal(time.January, month)
	assert.Equal(1, day)
}

func TestRD2GregNegative(t *testing.T) {
	assert := assert.New(t)
	var year int
	var month time.Month
	var day int

	year, month, day = greg.FromRD(-730)
	assert.Equal(-1, year)
	assert.Equal(time.January, month)
	assert.Equal(1, day)

	year, month, day = greg.FromRD(-36536)
	assert.Equal(-100, year)
	assert.Equal(time.December, month)
	assert.Equal(20, day)

	year, month, day = greg.FromRD(0)
	assert.Equal(0, year)
	assert.Equal(time.December, month)
	assert.Equal(31, day)
	year, month, day = greg.FromRD(-1)
	assert.Equal(0, year)
	assert.Equal(time.December, month)
	assert.Equal(30, day)
	year, month, day = greg.FromRD(-2)
	assert.Equal(0, year)
	assert.Equal(time.December, month)
	assert.Equal(29, day)
	year, month, day = greg.FromRD(-48)
	assert.Equal(0, year)
	assert.Equal(time.November, month)
	assert.Equal(13, day)
	year, month, day = greg.FromRD(-61)
	assert.Equal(0, year)
	assert.Equal(time.October, month)
	assert.Equal(31, day)

	year, month, day = greg.FromRD(-91)
	assert.Equal(0, year)
	assert.Equal(time.October, month)
	assert.Equal(1, day)

	year, month, day = greg.FromRD(-92)
	assert.Equal(0, year)
	assert.Equal(time.September, month)
	assert.Equal(30, day)

	year, month, day = greg.FromRD(-365)
	assert.Equal(0, year)
	assert.Equal(time.January, month)
	assert.Equal(1, day)

	year, month, day = greg.FromRD(-366)
	assert.Equal(-1, year)
	assert.Equal(time.December, month)
	assert.Equal(31, day)
}

func ExampleDaysIn() {
	days := greg.DaysIn(time.February, 2004)
	fmt.Println(days)
	// Output: 29
}

func ExampleDateToRD() {
	t := time.Date(2014, time.February, 19, 0, 0, 0, 0, time.UTC)
	rataDie := greg.DateToRD(t)
	fmt.Println(rataDie)
	// Output: 735283
}

func ExampleToRD() {
	rataDie := greg.ToRD(1995, time.December, 17)
	fmt.Println(rataDie)
	// Output: 728644
}

func ExampleFromRD() {
	year, month, day := greg.FromRD(737553)
	fmt.Println(year, month, day)
	// Output: 2020 May 8
}
