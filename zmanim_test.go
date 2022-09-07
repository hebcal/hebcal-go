package hebcal

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestZmanimChicago(t *testing.T) {
	assert := assert.New(t)
	latitude := 41.85003
	longitude := -87.65005
	dt := time.Date(2020, time.June, 5, 12, 0, 0, 0, time.UTC)
	zmanim := NewZmanim(latitude, longitude, dt, "America/Chicago")
	expected := []string{
		"Thu, 04 Jun 2020 20:21:34 -0500",
		"Fri, 05 Jun 2020 00:49:01 -0500",
		"Fri, 05 Jun 2020 03:25:34 -0500",
		"Fri, 05 Jun 2020 04:03:11 -0500",
		"Fri, 05 Jun 2020 04:12:55 -0500",
		"Fri, 05 Jun 2020 04:42:37 -0500",
		"Fri, 05 Jun 2020 05:16:28 -0500",
		"Fri, 05 Jun 2020 09:02:54 -0500",
		"Fri, 05 Jun 2020 08:26:30 -0500",
		"Fri, 05 Jun 2020 10:18:23 -0500",
		"Fri, 05 Jun 2020 09:54:00 -0500",
		"Fri, 05 Jun 2020 12:49:21 -0500",
		"Fri, 05 Jun 2020 13:27:05 -0500",
		"Fri, 05 Jun 2020 17:13:32 -0500",
		"Fri, 05 Jun 2020 18:47:53 -0500",
		"Fri, 05 Jun 2020 20:22:15 -0500",
		"Fri, 05 Jun 2020 20:56:06 -0500",
		"Fri, 05 Jun 2020 21:13:28 -0500",
	}
	times := []time.Time{
		zmanim.GregEve(),
		zmanim.ChatzotNight(),
		zmanim.AlotHaShachar(),
		zmanim.Misheyakir(),
		zmanim.MisheyakirMachmir(),
		zmanim.Dawn(),
		zmanim.Sunrise(),
		zmanim.SofZmanShma(),
		zmanim.SofZmanShmaMGA(),
		zmanim.SofZmanTfilla(),
		zmanim.SofZmanTfillaMGA(),
		zmanim.Chatzot(),
		zmanim.MinchaGedola(),
		zmanim.MinchaKetana(),
		zmanim.PlagHaMincha(),
		zmanim.Sunset(),
		zmanim.Dusk(),
		zmanim.Tzeit(8.5),
	}
	actual := make([]string, 18)
	for idx, t := range times {
		actual[idx] = t.Format(time.RFC1123Z)
	}
	assert.Equal(expected, actual)

	assert.Equal(4528916, zmanim.hour())
	assert.Equal(75, zmanim.hourMins())
	assert.Equal(2674500, zmanim.nightHour())
	assert.Equal(44, zmanim.nightHourMins())
}

func TestZmanimTelAviv(t *testing.T) {
	assert := assert.New(t)
	latitude := 32.08088
	longitude := 34.78057
	dt := time.Date(2021, time.March, 6, 12, 0, 0, 0, time.UTC)
	zmanim := NewZmanim(latitude, longitude, dt, "Asia/Jerusalem")
	expected := []string{
		"Fri, 05 Mar 2021 17:41:21 +0200",
		"Fri, 05 Mar 2021 23:51:55 +0200",
		"Sat, 06 Mar 2021 04:50:19 +0200",
		"Sat, 06 Mar 2021 05:12:02 +0200",
		"Sat, 06 Mar 2021 05:18:11 +0200",
		"Sat, 06 Mar 2021 05:38:01 +0200",
		"Sat, 06 Mar 2021 06:02:30 +0200",
		"Sat, 06 Mar 2021 08:57:23 +0200",
		"Sat, 06 Mar 2021 08:21:00 +0200",
		"Sat, 06 Mar 2021 09:55:41 +0200",
		"Sat, 06 Mar 2021 09:31:20 +0200",
		"Sat, 06 Mar 2021 11:52:17 +0200",
		"Sat, 06 Mar 2021 12:21:26 +0200",
		"Sat, 06 Mar 2021 15:16:20 +0200",
		"Sat, 06 Mar 2021 16:29:12 +0200",
		"Sat, 06 Mar 2021 17:42:05 +0200",
		"Sat, 06 Mar 2021 18:06:34 +0200",
		"Sat, 06 Mar 2021 18:18:23 +0200",
	}
	times := []time.Time{
		zmanim.GregEve(),
		zmanim.ChatzotNight(),
		zmanim.AlotHaShachar(),
		zmanim.Misheyakir(),
		zmanim.MisheyakirMachmir(),
		zmanim.Dawn(),
		zmanim.Sunrise(),
		zmanim.SofZmanShma(),
		zmanim.SofZmanShmaMGA(),
		zmanim.SofZmanTfilla(),
		zmanim.SofZmanTfillaMGA(),
		zmanim.Chatzot(),
		zmanim.MinchaGedola(),
		zmanim.MinchaKetana(),
		zmanim.PlagHaMincha(),
		zmanim.Sunset(),
		zmanim.Dusk(),
		zmanim.Tzeit(8.5),
	}
	actual := make([]string, 18)
	for idx, t := range times {
		actual[idx] = t.Format(time.RFC1123Z)
	}
	assert.Equal(expected, actual)
}

func TestZmanimHelsinki(t *testing.T) {
	assert := assert.New(t)
	latitude := 60.16952
	longitude := 24.93545
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
		zmanim := NewZmanim(latitude, longitude, dt, "Europe/Helsinki")
		var t time.Time
		if dt.Weekday() == time.Friday {
			t = zmanim.SunsetOffset(-18)
		} else {
			t = zmanim.Tzeit(8.5)
		}
		if (t == time.Time{}) {
			actual[idx] = "undefined"
		} else {
			actual[idx] = t.Format(time.RFC1123Z)
		}
	}
	expected := []string{
		"Fri, 15 May 2020 21:34:00 +0300",
		"Sat, 16 May 2020 23:45:23 +0300",
		"Fri, 22 May 2020 21:51:00 +0300",
		"Sun, 24 May 2020 00:25:46 +0300",
		"Fri, 29 May 2020 22:05:00 +0300",
		"undefined",
		"Fri, 05 Jun 2020 22:17:00 +0300",
		"undefined",
		"Fri, 31 Jul 2020 21:36:00 +0300",
		"Sat, 01 Aug 2020 23:31:14 +0300",
	}
	assert.Equal(expected, actual)
}
