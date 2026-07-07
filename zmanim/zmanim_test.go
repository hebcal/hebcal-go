package zmanim_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hebcal/hebcal-go/zmanim"
	"github.com/stretchr/testify/assert"
)

func TestZmanimChicago(t *testing.T) {
	assert := assert.New(t)
	dt := time.Date(2020, time.June, 5, 12, 0, 0, 0, time.UTC)
	location := zmanim.NewLocation("Chicago", "US", 41.85003, -87.65005, 180, "America/Chicago")
	zman := zmanim.New(&location, dt)
	expected := []string{
		"Thu, 04 Jun 2020 20:21:48 -0500",
		"Fri, 05 Jun 2020 00:49:04 -0500",
		"Fri, 05 Jun 2020 03:25:29 -0500",
		"Fri, 05 Jun 2020 04:03:04 -0500",
		"Fri, 05 Jun 2020 04:12:48 -0500",
		"Fri, 05 Jun 2020 05:16:19 -0500",
		"Fri, 05 Jun 2020 09:02:51 -0500",
		"Fri, 05 Jun 2020 08:26:51 -0500",
		"Fri, 05 Jun 2020 10:18:21 -0500",
		"Fri, 05 Jun 2020 09:54:21 -0500",
		"Fri, 05 Jun 2020 12:49:23 -0500",
		"Fri, 05 Jun 2020 13:27:08 -0500",
		"Fri, 05 Jun 2020 17:13:40 -0500",
		"Fri, 05 Jun 2020 18:48:03 -0500",
		"Fri, 05 Jun 2020 20:22:27 -0500",
		"Fri, 05 Jun 2020 20:50:19 -0500",
		"Fri, 05 Jun 2020 21:03:49 -0500",
		"Fri, 05 Jun 2020 21:13:47 -0500",
	}
	times := makeTestTimes(zman)
	actual := make([]string, 18)
	for idx, t := range times {
		actual[idx] = t.Format(time.RFC1123Z)
	}
	assert.Equal(expected, actual)

	assert.Equal(4530.666666666667, zman.Hour())
	// assert.Equal(2674.500, zman.nightHour())
}

func TestZmanimTelAviv(t *testing.T) {
	assert := assert.New(t)
	dt := time.Date(2021, time.March, 6, 12, 0, 0, 0, time.UTC)
	location := zmanim.NewLocation("Tel Aviv", "IL", 32.08088, 34.78057, 15, "Asia/Jerusalem")
	zman := zmanim.New(&location, dt)
	expected := []string{
		"Fri, 05 Mar 2021 17:41:37 +0200",
		"Fri, 05 Mar 2021 23:51:56 +0200",
		"Sat, 06 Mar 2021 04:50:09 +0200",
		"Sat, 06 Mar 2021 05:11:52 +0200",
		"Sat, 06 Mar 2021 05:18:00 +0200",
		"Sat, 06 Mar 2021 06:02:15 +0200",
		"Sat, 06 Mar 2021 08:57:16 +0200",
		"Sat, 06 Mar 2021 08:21:16 +0200",
		"Sat, 06 Mar 2021 09:55:37 +0200",
		"Sat, 06 Mar 2021 09:31:37 +0200",
		"Sat, 06 Mar 2021 11:52:18 +0200",
		"Sat, 06 Mar 2021 12:21:28 +0200",
		"Sat, 06 Mar 2021 15:16:30 +0200",
		"Sat, 06 Mar 2021 16:29:26 +0200",
		"Sat, 06 Mar 2021 17:42:22 +0200",
		"Sat, 06 Mar 2021 17:58:27 +0200",
		"Sat, 06 Mar 2021 18:11:57 +0200",
		"Sat, 06 Mar 2021 18:18:39 +0200",
	}
	times := makeTestTimes(zman)
	actual := make([]string, 18)
	for idx, t := range times {
		actual[idx] = t.Format(time.RFC1123Z)
	}
	assert.Equal(expected, actual)
}

func makeTestTimes(zman zmanim.Zmanim) []time.Time {
	times := []time.Time{
		zman.GregEve(),
		zman.ChatzotNight(),
		zman.AlotHaShachar(),
		zman.Misheyakir(),
		zman.MisheyakirMachmir(),
		zman.Sunrise(),
		zman.SofZmanShma(),
		zman.SofZmanShmaMGA(),
		zman.SofZmanTfilla(),
		zman.SofZmanTfillaMGA(),
		zman.Chatzot(),
		zman.MinchaGedola(),
		zman.MinchaKetana(),
		zman.PlagHaMincha(),
		zman.Sunset(),
		zman.BeinHashmashos(),
		zman.Tzeit(7.083),
		zman.Tzeit(8.5),
	}
	return times
}

func TestZmanimHelsinki(t *testing.T) {
	assert := assert.New(t)
	location := zmanim.NewLocation("Helsinki", "FI", 60.16952, 24.93545, 26, "Europe/Helsinki")
	dates := []struct {
		yy int
		mm time.Month
		dd int
	}{
		{2020, time.May, 15},
		{2020, time.May, 16},
		{2020, time.May, 22},
		{2020, time.May, 23},
		{2020, time.May, 29},
		{2020, time.May, 30},
		{2020, time.June, 5},
		{2020, time.June, 6},
		{2020, time.July, 31},
		{2020, time.August, 1},
	}
	actual := make([]string, len(dates))
	for idx, date := range dates {
		dt := time.Date(date.yy, date.mm, date.dd, 12, 0, 0, 0, time.UTC)
		zman := zmanim.New(&location, dt)
		var t time.Time
		if dt.Weekday() == time.Friday {
			t = zman.SunsetOffset(-18, true)
		} else {
			t = zman.Tzeit(8.5)
		}
		if t.IsZero() {
			actual[idx] = "undefined"
		} else {
			actual[idx] = t.Format(time.RFC1123Z)
		}
	}
	expected := []string{
		"Fri, 15 May 2020 21:36:00 +0300",
		"Sat, 16 May 2020 23:48:49 +0300",
		"Fri, 22 May 2020 21:52:00 +0300",
		"Sun, 24 May 2020 00:31:13 +0300",
		"Fri, 29 May 2020 22:06:00 +0300",
		"undefined",
		"Fri, 05 Jun 2020 22:18:00 +0300",
		"undefined",
		"Fri, 31 Jul 2020 21:35:00 +0300",
		"Sat, 01 Aug 2020 23:28:10 +0300",
	}
	assert.Equal(expected, actual)
}

func ExampleZmanim_SunsetOffset() {
	dt := time.Date(2020, time.June, 5, 12, 0, 0, 0, time.UTC)
	location := zmanim.NewLocation("Chicago", "US", 41.85003, -87.65005, 180, "America/Chicago")
	zman := zmanim.New(&location, dt)
	fmt.Println(zman.SunsetOffset(-18, true))
	fmt.Println(zman.SunsetOffset(50, true))
	// Output:
	// 2020-06-05 20:04:00 -0500 CDT
	// 2020-06-05 21:12:00 -0500 CDT
}

func TestZmanimAmsterdam(t *testing.T) {
	assert := assert.New(t)
	dt := time.Date(2023, time.May, 29, 12, 0, 0, 0, time.UTC)
	location := zmanim.NewLocation("Amsterdam", "NL", 52.37403, 4.88969, 12, "Europe/Amsterdam")
	zman := zmanim.New(&location, dt)
	alot := zman.AlotHaShachar()
	assert.Equal(alot.IsZero(), true)
	assert.Equal(time.Time{}, alot)
}

// TestZmanimElevation verifies that enabling UseElevation shifts sunrise earlier
// and sunset later (a person at a higher elevation sees past the sea-level
// horizon), while leaving degree-based zmanim unaffected. Reference values for
// Jerusalem (elevation 786m) on 2022-06-21 come from the elevation-aware NOAA
// calculator, matching KosherJava.
func TestZmanimElevation(t *testing.T) {
	assert := assert.New(t)
	dt := time.Date(2022, time.June, 21, 12, 0, 0, 0, time.UTC)
	location := zmanim.NewLocation("Jerusalem", "IL", 31.76904, 35.21633, 786, "Asia/Jerusalem")

	sea := zmanim.New(&location, dt)
	assert.False(sea.UseElevation)
	seaSunrise := sea.Sunrise()
	seaSunset := sea.Sunset()
	seaAlot := sea.AlotHaShachar()

	elev := zmanim.New(&location, dt)
	elev.UseElevation = true
	elevSunrise := elev.Sunrise()
	elevSunset := elev.Sunset()

	// Elevation makes sunrise earlier and sunset later.
	assert.True(elevSunrise.Before(seaSunrise), "elevation sunrise should be earlier")
	assert.True(elevSunset.After(seaSunset), "elevation sunset should be later")

	// Locked-in reference values (Asia/Jerusalem, +03:00).
	assert.Equal("2022-06-21 05:29:18", elevSunrise.Format("2006-01-02 15:04:05"))
	assert.Equal("2022-06-21 19:52:32", elevSunset.Format("2006-01-02 15:04:05"))
	assert.Equal("2022-06-21 19:47:42", seaSunset.Format("2006-01-02 15:04:05"))

	// Degree-based zmanim ignore elevation.
	assert.Equal(seaAlot, elev.AlotHaShachar(), "alot must not change with elevation")
}
