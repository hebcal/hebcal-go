package hebcal_test

import (
	"testing"
	"time"

	"github.com/hebcal/hdate"
	"github.com/hebcal/hebcal-go/event"
	"github.com/hebcal/hebcal-go/hebcal"
	"github.com/hebcal/hebcal-go/zmanim"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Before 14 September 1752 the civil date of a Hebrew date depends on
// NoJulian, but zmanim are always computed for the proleptic Gregorian date,
// so a timed event falls on its own day and its time does not depend on
// NoJulian.
func TestZmanimAreProlepticBefore1752(t *testing.T) {
	timesFor := func(noJulian bool) []time.Time {
		// Year < 1753 turns candle lighting off, so give an explicit range.
		opts := hebcal.Options{
			Start:    hdate.New(5460, hdate.Tishrei, 1),
			End:      hdate.New(5461, hdate.Elul, 29),
			Location: zmanim.LookupCity("Jerusalem"), IL: true,
			CandleLighting: true, KiddushLevanaMaharil: true,
			NoJulian: noJulian,
		}
		events, err := hebcal.HebrewCalendar(&opts)
		require.NoError(t, err)
		var times []time.Time
		for _, ev := range events {
			timed, ok := ev.(event.Timed)
			if !ok {
				continue
			}
			tm := timed.Time()
			gy, gm, gd := ev.Date().ProlepticGreg()
			y, m, d := tm.Date()
			assert.Equal(t, [3]int{gy, int(gm), gd}, [3]int{y, int(m), d}, ev.Desc())
			assert.Equal(t, ev.Date().Weekday(), tm.Weekday(), ev.Desc())
			times = append(times, tm)
		}
		require.NotEmpty(t, times)
		return times
	}
	julian := timesFor(false)
	proleptic := timesFor(true)
	assert.Equal(t, proleptic, julian)
}

func TestIsAssurBemelachaBefore1752(t *testing.T) {
	jer := zmanim.LookupCity("Jerusalem")
	loc, err := time.LoadLocation(jer.TimeZoneID)
	require.NoError(t, err)
	// Saturday 3 July 1700 (proleptic Gregorian), mid-afternoon.
	tm := time.Date(1700, time.July, 3, 14, 0, 0, 0, loc)
	require.Equal(t, time.Saturday, tm.Weekday())
	assur, err := hebcal.IsAssurBemelacha(tm, jer, true, false)
	require.NoError(t, err)
	assert.True(t, assur)
	// Wednesday 7 July 1700, mid-afternoon.
	assur, err = hebcal.IsAssurBemelacha(tm.AddDate(0, 0, 4), jer, true, false)
	require.NoError(t, err)
	assert.False(t, assur)
}
