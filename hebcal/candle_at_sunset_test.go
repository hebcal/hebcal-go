package hebcal

import (
	"strings"
	"testing"
	"time"

	"github.com/hebcal/hdate"
	"github.com/hebcal/hebcal-go/zmanim"
	"github.com/stretchr/testify/assert"
)

// candleTimes returns the candle-lighting and Havdalah events between start
// and end, as "MM-DD HH:MM desc" in local time.
func candleTimes(t *testing.T, opts CalOptions, start, end string) string {
	t.Helper()
	s, _ := time.Parse("2006-01-02", start)
	e, _ := time.Parse("2006-01-02", end)
	opts.Start, opts.End = hdate.FromTime(s), hdate.FromTime(e)
	opts.CandleLighting = true
	events, err := HebrewCalendar(&opts)
	if err != nil {
		t.Fatal(err)
	}
	var out []string
	for _, ev := range events {
		te, ok := ev.(TimedEvent)
		if !ok || strings.HasPrefix(te.Desc, "Fast ") || strings.Contains(te.Desc, "Chametz") ||
			te.Desc == "Finish eating chametz" {
			continue
		}
		out = append(out, te.EventTime.Format("01-02 15:04 ")+te.Desc)
	}
	return strings.Join(out, ",")
}

// The expected times are @hebcal/core's with candleLightingMins: 0. The first
// night of Pesach lights at sunset, the second night at tzeit (Havdalah rules),
// and Chanukah on Friday at sunset with Shabbat.
func TestCandleLightingAtSunset(t *testing.T) {
	ny := CalOptions{Location: zmanim.LookupCity("New York"), CandleLightingAtSunset: true}
	assert.Equal(t,
		"03-27 19:15 Candle lighting,03-28 19:58 Havdalah,04-01 19:20 Candle lighting,"+
			"04-02 20:03 Candle lighting,04-03 19:22 Candle lighting,04-04 20:05 Havdalah,"+
			"04-07 19:27 Candle lighting,04-08 20:10 Candle lighting,04-09 20:11 Havdalah,"+
			"04-10 19:30 Candle lighting",
		candleTimes(t, ny, "2026-03-27", "2026-04-10"))
	assert.Equal(t,
		"12-04 16:28 Chanukah: 1 Candle,12-04 16:28 Candle lighting,"+
			"12-05 17:14 Chanukah: 2 Candles,12-05 17:14 Havdalah",
		candleTimes(t, ny, "2026-12-04", "2026-12-05"))

	// The Israeli city custom (40 minutes in Jerusalem) does not apply.
	jlm := CalOptions{Location: zmanim.LookupCity("Jerusalem"), IL: true, CandleLightingAtSunset: true}
	assert.Equal(t,
		"03-27 18:55 Candle lighting,03-28 19:32 Havdalah,04-01 18:58 Candle lighting,"+
			"04-02 19:36 Havdalah",
		candleTimes(t, jlm, "2026-03-27", "2026-04-02"))

	// CandleLightingMins is ignored.
	ny.CandleLightingMins = 40
	assert.True(t, strings.HasPrefix(candleTimes(t, ny, "2026-03-27", "2026-03-27"), "03-27 19:15 "))
}
