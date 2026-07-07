package hebcal_test

import (
	"testing"
	"time"

	"github.com/hebcal/hdate"
	"github.com/hebcal/hebcal-go/hebcal"
	"github.com/hebcal/hebcal-go/zmanim"
	"github.com/stretchr/testify/assert"
)

func atLoc(t *testing.T, s, tzid string) time.Time {
	t.Helper()
	loc, err := time.LoadLocation(tzid)
	assert.NoError(t, err)
	tm, err := time.ParseInLocation("2006-01-02T15:04:05", s, loc)
	assert.NoError(t, err)
	return tm
}

// Reference truth values come from @hebcal/core's isAssurBemlacha.
func TestIsAssurBemlacha(t *testing.T) {
	assert := assert.New(t)
	jer := zmanim.NewLocation("Jerusalem", "IL", 31.76904, 35.21633, 786, "Asia/Jerusalem")
	ny := zmanim.NewLocation("New York", "US", 40.71427, -74.00597, 0, "America/New_York")

	tests := []struct {
		loc  *zmanim.Location
		il   bool
		when string
		tzid string
		want bool
	}{
		{&jer, true, "2024-04-26T18:00:00", "Asia/Jerusalem", false},   // Fri before sunset (~19:16)
		{&jer, true, "2024-04-26T19:30:00", "Asia/Jerusalem", true},    // Fri after sunset
		{&jer, true, "2024-04-27T19:00:00", "Asia/Jerusalem", true},    // Shabbat before tzais
		{&jer, true, "2024-04-27T20:30:00", "Asia/Jerusalem", false},   // Shabbat after tzais
		{&ny, false, "2025-06-21T17:08:10", "America/New_York", true},  // Shabbat afternoon
		{&ny, false, "2024-04-24T21:00:00", "America/New_York", false}, // Wednesday
	}
	for _, tc := range tests {
		got, err := hebcal.IsAssurBemlacha(atLoc(t, tc.when, tc.tzid), tc.loc, tc.il, false)
		assert.NoError(err, tc.when)
		assert.Equal(tc.want, got, "%s %s", tc.loc.Name, tc.when)
	}
}

func TestIsAssurBemlachaPolarError(t *testing.T) {
	assert := assert.New(t)
	tromso := zmanim.NewLocation("Tromso", "NO", 69.6489, 18.9551, 0, "Europe/Oslo")
	// Midnight sun: sunset does not occur, so an error is returned.
	_, err := hebcal.IsAssurBemlacha(atLoc(t, "2024-06-21T12:00:00", "Europe/Oslo"), &tromso, false, false)
	assert.Error(err)
}

func TestGetHolidaysOnDate(t *testing.T) {
	assert := assert.New(t)
	// 15 Nisan 5784 = Pesach I (2024-04-23).
	hd := hdate.New(5784, hdate.Nisan, 15)
	events := hebcal.GetHolidaysOnDate(hd, false)
	found := false
	for _, ev := range events {
		if ev.Desc == "Pesach I" {
			found = true
		}
	}
	assert.True(found, "expected Pesach I on 15 Nisan 5784")

	// A random weekday with no holiday.
	none := hebcal.GetHolidaysOnDate(hdate.New(5784, hdate.Iyyar, 3), false)
	assert.Empty(none)
}
