package hebcal_test

import (
	"errors"
	"testing"

	"github.com/hebcal/hdate"
	"github.com/hebcal/hebcal-go/hebcal"
	"github.com/hebcal/hebcal-go/zmanim"
	"github.com/stretchr/testify/assert"
)

// HebrewCalendar's option errors are exported values: each matches itself
// and its kind with errors.Is, and keeps its message.
func TestOptionErrors(t *testing.T) {
	chicago := zmanim.LookupCity("Chicago")
	tests := []struct {
		name string
		opts hebcal.CalOptions
		want error
		kind error
		msg  string
	}{
		{"candles", hebcal.CalOptions{CandleLighting: true}, hebcal.ErrCandleLightingLocation,
			hebcal.ErrLocationRequired, "opts.CandleLighting requires opts.Location"},
		{"kiddush levana", hebcal.CalOptions{KiddushLevanaMaharil: true}, hebcal.ErrKiddushLevanaLocation,
			hebcal.ErrLocationRequired, "opts.KiddushLevanaMaharil requires opts.Location"},
		{"sunrise", hebcal.CalOptions{SunriseSunset: true}, hebcal.ErrSunriseSunsetLocation,
			hebcal.ErrLocationRequired, "opts.SunriseSunset requires opts.Location"},
		{"zmanim", hebcal.CalOptions{DailyZmanim: true}, hebcal.ErrDailyZmanimLocation,
			hebcal.ErrLocationRequired, "opts.DailyZmanim requires opts.Location"},
		{"havdalah", hebcal.CalOptions{CandleLighting: true, Location: chicago, HavdalahMins: 50, HavdalahDeg: 8.5},
			hebcal.ErrHavdalahMinsAndDeg, hebcal.ErrMutuallyExclusive,
			"opts.HavdalahMins and opts.HavdalahDeg are mutually exclusive"},
		{"fast start", hebcal.CalOptions{CandleLighting: true, Location: chicago, FastStartMins: 72, FastStartDeg: 16.1},
			hebcal.ErrFastStartMinsAndDeg, hebcal.ErrMutuallyExclusive,
			"opts.FastStartMins and opts.FastStartDeg are mutually exclusive"},
		{"fast end", hebcal.CalOptions{CandleLighting: true, Location: chicago, FastEndMins: 20, FastEndDeg: 7.083},
			hebcal.ErrFastEndMinsAndDeg, hebcal.ErrMutuallyExclusive,
			"opts.FastEndMins and opts.FastEndDeg are mutually exclusive"},
		{"tisha bav", hebcal.CalOptions{CandleLighting: true, Location: chicago, TishaBavEndMins: 20, TishaBavEndDeg: 6.45},
			hebcal.ErrTishaBavEndMinsAndDeg, hebcal.ErrMutuallyExclusive,
			"opts.TishaBavEndMins and opts.TishaBavEndDeg are mutually exclusive"},
		{"start without end", hebcal.CalOptions{Start: hdate.New(5786, hdate.Tishrei, 1)},
			hebcal.ErrStartRequiresEnd, nil, "opts.Start requires opts.End"},
		{"hebrew year", hebcal.CalOptions{Year: -5, IsHebrewYear: true},
			hebcal.ErrInvalidHebrewYear, nil, "invalid Hebrew year"},
	}
	for _, tt := range tests {
		opts := tt.opts
		_, err := hebcal.HebrewCalendar(&opts)
		if assert.Error(t, err, tt.name) {
			assert.True(t, errors.Is(err, tt.want), tt.name)
			assert.Equal(t, tt.msg, err.Error(), tt.name)
			if tt.kind != nil {
				assert.True(t, errors.Is(err, tt.kind), tt.name)
			}
			assert.False(t, errors.Is(err, zmanim.ErrTimeZone), tt.name)
		}
	}
	assert.False(t, errors.Is(hebcal.ErrCandleLightingLocation, hebcal.ErrMutuallyExclusive))
	assert.False(t, errors.Is(hebcal.ErrHavdalahMinsAndDeg, hebcal.ErrLocationRequired))
}

// A location that cannot be used is an error, not a panic, and only when
// the calendar needs times from it.
func TestUnusableLocation(t *testing.T) {
	badZone := &zmanim.Location{Name: "Nowhere", Latitude: 40, Longitude: -74, TimeZoneID: "Mars/Olympus_Mons"}
	_, err := hebcal.HebrewCalendar(&hebcal.CalOptions{Year: 2026, CandleLighting: true, Location: badZone})
	assert.True(t, errors.Is(err, zmanim.ErrTimeZone), "%v", err)

	badLat := &zmanim.Location{Name: "Nowhere", Latitude: 140, Longitude: -74, TimeZoneID: "America/New_York"}
	_, err = hebcal.HebrewCalendar(&hebcal.CalOptions{Year: 2026, SunriseSunset: true, Location: badLat})
	assert.True(t, errors.Is(err, zmanim.ErrLatitude), "%v", err)

	// With no times asked for, the location only picks the clock format.
	_, err = hebcal.HebrewCalendar(&hebcal.CalOptions{Year: 2026, Location: badZone})
	assert.NoError(t, err)

	_, err = hebcal.IsAssurBemelacha(hdate.New(5786, hdate.Tishrei, 1).Gregorian(), badZone, false, false)
	assert.True(t, errors.Is(err, zmanim.ErrTimeZone), "%v", err)
}
