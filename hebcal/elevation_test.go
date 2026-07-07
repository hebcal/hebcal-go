package hebcal_test

import (
	"strings"
	"testing"

	"github.com/hebcal/hebcal-go/hebcal"
	"github.com/hebcal/hebcal-go/zmanim"
	"github.com/stretchr/testify/assert"
)

// findOnDate returns the first event rendered in English on the given Gregorian
// date whose description starts with prefix, or "".
func findOnDate(t *testing.T, opts *hebcal.CalOptions, y, m, d int, prefix string) string {
	t.Helper()
	events, err := hebcal.HebrewCalendar(opts)
	assert.NoError(t, err)
	for _, ev := range events {
		gy, gm, gd := ev.GetDate().Greg()
		if gy == y && int(gm) == m && gd == d {
			r := ev.Render("en")
			if strings.HasPrefix(r, prefix) {
				return r
			}
		}
	}
	return ""
}

// TestUseElevationCalOption verifies that CalOptions.UseElevation flows through
// to candle-lighting and daily zmanim: sunrise/sunset-based times shift with the
// location's elevation, while degree-based zmanim do not. Jerusalem has an
// elevation of 786m.
func TestUseElevationCalOption(t *testing.T) {
	assert := assert.New(t)
	loc := zmanim.LookupCity("Jerusalem")
	newOpts := func(useElevation bool) *hebcal.CalOptions {
		return &hebcal.CalOptions{
			Year: 2022, Month: 4, Location: loc,
			CandleLighting: true, DailyZmanim: true, UseElevation: useElevation,
		}
	}
	y, m, d := 2022, 4, 8 // a Friday in Jerusalem

	// Candle lighting (sunset-based) is later at elevation.
	assert.Equal("Candle lighting: 6:23", findOnDate(t, newOpts(false), y, m, d, "Candle lighting"))
	assert.Equal("Candle lighting: 6:27", findOnDate(t, newOpts(true), y, m, d, "Candle lighting"))

	// Sunrise (sunset-based) is earlier at elevation.
	assert.Equal("Sunrise: 6:19", findOnDate(t, newOpts(false), y, m, d, "Sunrise"))
	assert.Equal("Sunrise: 6:14", findOnDate(t, newOpts(true), y, m, d, "Sunrise"))

	// Alot HaShachar (degree-based) is unaffected by elevation.
	assert.Equal("Alot HaShachar: 5:05", findOnDate(t, newOpts(false), y, m, d, "Alot HaShachar"))
	assert.Equal("Alot HaShachar: 5:05", findOnDate(t, newOpts(true), y, m, d, "Alot HaShachar"))
}

// TestUseElevationDefaultFalse confirms the default (unset) matches the
// non-elevation calculation.
func TestUseElevationDefaultFalse(t *testing.T) {
	assert := assert.New(t)
	loc := zmanim.LookupCity("Jerusalem")
	dflt := &hebcal.CalOptions{Year: 2022, Month: 4, Location: loc, CandleLighting: true}
	explicit := &hebcal.CalOptions{Year: 2022, Month: 4, Location: loc, CandleLighting: true, UseElevation: false}
	assert.Equal(
		findOnDate(t, dflt, 2022, 4, 8, "Candle lighting"),
		findOnDate(t, explicit, 2022, 4, 8, "Candle lighting"),
	)
}
