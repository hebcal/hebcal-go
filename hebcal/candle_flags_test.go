package hebcal_test

import (
	"testing"
	"time"

	"github.com/hebcal/hdate"
	"github.com/stretchr/testify/assert"

	"github.com/hebcal/hebcal-go/event"
	"github.com/hebcal/hebcal-go/hebcal"
	"github.com/hebcal/hebcal-go/zmanim"
)

// A candle-lighting event carries exactly one of LightCandles (minutes
// before sunset: Shabbat or the start of a chag) and LightCandlesTzeis
// (nightfall / Havdalah time), matching @hebcal/core.
func TestCandleLightingFlags(t *testing.T) {
	cases := []struct {
		date  string
		label string
		time  string
		flags event.Flags
	}{
		{"2024-10-16", "Wed, Erev Sukkot", "17:55", event.ChulOnly | event.Erev | event.LightCandles},
		{"2024-11-08", "plain Friday", "16:25", event.LightCandles},
		{"2020-05-29", "Fri, Shavuot I", "20:01", event.Chag | event.ChulOnly | event.LightCandles},
		{"2020-04-10", "Fri, Yom Tov ends, Shabbat begins", "19:12",
			event.Chag | event.ChulOnly | event.YomTovEnds | event.LightCandles},
		{"2024-10-17", "Thu, Sukkot I", "18:53", event.Chag | event.ChulOnly | event.LightCandlesTzeis},
		{"2021-03-27", "Sat, Erev Pesach", "19:57", event.ChulOnly | event.Erev | event.LightCandlesTzeis},
	}
	loc := zmanim.LookupCity("New York")
	for _, c := range cases {
		t.Run(c.date+" "+c.label, func(t *testing.T) {
			d, _ := time.Parse("2006-01-02", c.date)
			hd := hdate.FromTime(d)
			events, err := hebcal.HebrewCalendar(&hebcal.Options{
				Start: hd, End: hd, CandleLighting: true, Location: loc,
			})
			assert.NoError(t, err)
			var found bool
			for _, ev := range events {
				te, ok := ev.(event.Timed)
				if !ok || te.Desc() != "Candle lighting" {
					continue
				}
				found = true
				assert.Equal(t, c.time, te.Time().Format("15:04"))
				assert.Equal(t, c.flags, te.Flags(), te.Flags().String())
			}
			assert.True(t, found)
		})
	}
}

func TestCandleLightingNeverBothFlags(t *testing.T) {
	both := event.LightCandles | event.LightCandlesTzeis
	for _, city := range []string{"New York", "Jerusalem"} {
		events, err := hebcal.HebrewCalendar(&hebcal.Options{
			Year: 5780, IsHebrewYear: true, NumYears: 12, IL: city == "Jerusalem",
			CandleLighting: true, Location: zmanim.LookupCity(city),
		})
		assert.NoError(t, err)
		for _, ev := range events {
			assert.False(t, ev.Flags().Has(both), ev.Render("en"))
		}
	}
}
