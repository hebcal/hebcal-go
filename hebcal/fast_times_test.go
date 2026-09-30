package hebcal

import (
	"strings"
	"testing"

	"github.com/hebcal/hebcal-go/zmanim"
	"github.com/stretchr/testify/assert"
)

// fastTimes returns the "Fast begins" and "Fast ends" times of a 2026
// calendar, as "MM-DD begins|ends HH:MM".
func fastTimes(t *testing.T, city string, opts CalOptions) string {
	t.Helper()
	opts.Year = 2026
	opts.Location = zmanim.LookupCity(city)
	opts.CandleLighting = true
	events, err := HebrewCalendar(&opts)
	if err != nil {
		t.Fatal(err)
	}
	var out []string
	for _, ev := range events {
		te, ok := ev.(TimedEvent)
		if !ok || !strings.HasPrefix(te.Desc, "Fast ") {
			continue
		}
		out = append(out, te.EventTime.Format("01-02 ")+strings.TrimPrefix(te.Desc, "Fast ")+te.EventTime.Format(" 15:04"))
	}
	return strings.Join(out, ",")
}

// The expected times are @hebcal/core 6.11's for the same calendars. The
// fasts are Ta'anit Esther (03-02), Ta'anit Bechorot (04-01, no end), Tzom
// Tammuz (07-02), Tish'a B'Av (07-22/23), Tzom Gedaliah (09-14) and Asara
// B'Tevet (12-20).
func TestFastStartEndOptions(t *testing.T) {
	cases := []struct {
		name string
		city string
		opts CalOptions
		want string
	}{
		{"diaspora defaults", "New York", CalOptions{},
			"03-02 begins 05:08,03-02 ends 18:22,04-01 begins 05:16,07-02 begins 03:41,07-02 ends 21:11," +
				"07-22 begins 20:21,07-23 ends 20:54,09-14 begins 05:13,09-14 ends 19:40,12-20 begins 05:48,12-20 ends 17:09"},
		// minor fasts end 15 minutes after sunset in Israel; Tish'a B'Av does not
		{"israel defaults", "Jerusalem", CalOptions{IL: true},
			"03-02 begins 04:53,03-02 ends 17:53,04-01 begins 05:15,07-02 begins 04:10,07-02 ends 20:04," +
				"07-22 begins 19:43,07-23 ends 20:11,09-14 begins 05:09,09-14 ends 19:02,12-20 begins 05:17,12-20 ends 16:54"},
		{"FastStartDeg", "New York", CalOptions{FastStartDeg: 19.8},
			"03-02 begins 04:48,03-02 ends 18:22,04-01 begins 04:55,07-02 begins 03:07,07-02 ends 21:11," +
				"07-22 begins 20:21,07-23 ends 20:54,09-14 begins 04:52,09-14 ends 19:40,12-20 begins 05:27,12-20 ends 17:09"},
		{"FastStartMins", "New York", CalOptions{FastStartMins: 72},
			"03-02 begins 05:16,03-02 ends 18:22,04-01 begins 05:27,07-02 begins 04:17,07-02 ends 21:11," +
				"07-22 begins 20:21,07-23 ends 20:54,09-14 begins 05:23,09-14 ends 19:40,12-20 begins 06:04,12-20 ends 17:09"},
		{"FastEndDeg in Israel", "Jerusalem", CalOptions{IL: true, FastEndDeg: 8.5},
			"03-02 begins 04:53,03-02 ends 18:14,04-01 begins 05:15,07-02 begins 04:10,07-02 ends 20:31," +
				"07-22 begins 19:43,07-23 ends 20:11,09-14 begins 05:09,09-14 ends 19:23,12-20 begins 05:17,12-20 ends 17:19"},
		{"FastEndMins", "New York", CalOptions{FastEndMins: 45},
			"03-02 begins 05:08,03-02 ends 18:34,04-01 begins 05:16,07-02 begins 03:41,07-02 ends 21:16," +
				"07-22 begins 20:21,07-23 ends 20:54,09-14 begins 05:13,09-14 ends 19:52,12-20 begins 05:48,12-20 ends 17:16"},
		{"TishaBavEndDeg", "New York", CalOptions{TishaBavEndDeg: 8.5},
			"03-02 begins 05:08,03-02 ends 18:22,04-01 begins 05:16,07-02 begins 03:41,07-02 ends 21:11," +
				"07-22 begins 20:21,07-23 ends 21:08,09-14 begins 05:13,09-14 ends 19:40,12-20 begins 05:48,12-20 ends 17:09"},
		{"TishaBavEndMins in Israel", "Jerusalem", CalOptions{IL: true, TishaBavEndMins: 50},
			"03-02 begins 04:53,03-02 ends 17:53,04-01 begins 05:15,07-02 begins 04:10,07-02 ends 20:04," +
				"07-22 begins 19:43,07-23 ends 20:32,09-14 begins 05:09,09-14 ends 19:02,12-20 begins 05:17,12-20 ends 16:54"},
		// negative values are taken as positive
		{"negative FastEndMins", "New York", CalOptions{FastEndMins: -45},
			"03-02 begins 05:08,03-02 ends 18:34,04-01 begins 05:16,07-02 begins 03:41,07-02 ends 21:16," +
				"07-22 begins 20:21,07-23 ends 20:54,09-14 begins 05:13,09-14 ends 19:52,12-20 begins 05:48,12-20 ends 17:16"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.want, fastTimes(t, c.city, c.opts))
		})
	}
}

func TestFastOptionsAreMutuallyExclusive(t *testing.T) {
	for _, opts := range []CalOptions{
		{FastStartDeg: 19.8, FastStartMins: 72},
		{FastEndDeg: 8.5, FastEndMins: 45},
		{TishaBavEndDeg: 8.5, TishaBavEndMins: 50},
	} {
		opts.Year = 2026
		opts.Location = zmanim.LookupCity("New York")
		opts.CandleLighting = true
		_, err := HebrewCalendar(&opts)
		assert.Error(t, err)
	}
}
