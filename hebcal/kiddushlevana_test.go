package hebcal

import (
	"strings"
	"testing"

	"github.com/hebcal/hdate"
	"github.com/hebcal/hebcal-go/event"
	"github.com/hebcal/hebcal-go/molad"
	"github.com/hebcal/hebcal-go/zmanim"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func kiddushLevanaMonth(t *testing.T, loc *zmanim.Location, hyear int, hmonth hdate.HMonth) KiddushLevanaEvent {
	t.Helper()
	opts := CalOptions{
		Start:                hdate.New(hyear, hmonth, 1),
		End:                  hdate.New(hyear, hmonth, hdate.DaysInMonth(hmonth, hyear)),
		Location:             loc,
		KiddushLevanaMaharil: true,
		NoHolidays:           true,
	}
	events, err := HebrewCalendar(&opts)
	require.NoError(t, err)
	require.Len(t, events, 1)
	ev, ok := events[0].(KiddushLevanaEvent)
	require.True(t, ok, "got %T", events[0])
	return ev
}

func klSummary(ev KiddushLevanaEvent) string {
	return ev.Date.String() + " | " + ev.Date.Gregorian().Format("Mon Jan 02 2006") +
		" | " + ev.EventTime.Format("15:04")
}

func TestKiddushLevana(t *testing.T) {
	newYork := zmanim.LookupCity("New York")
	jerusalem := zmanim.LookupCity("Jerusalem")
	tests := []struct {
		name   string
		loc    *zmanim.Location
		hyear  int
		hmonth hdate.HMonth
		want   string
	}{
		{"weeknight, before midnight", newYork, 5787, hdate.Cheshvan,
			"14 Cheshvan 5787 | Sun Oct 25 2026 | 21:44"},
		// Maharil's time is Tue Nov 24 2026 9:28am
		{"daytime is moved back to alot hashachar", newYork, 5787, hdate.Kislev,
			"14 Kislev 5787 | Tue Nov 24 2026 | 05:27"},
		{"Motzei Shabbat is not moved", newYork, 5787, hdate.Adar1,
			"13 Adar I 5787 | Sat Feb 20 2027 | 23:40"},
		// Maharil's time is Shabbat Sep 26 2026 9:00am, first day of Sukkot
		{"Shabbat and Yom Tov moves back to candle-lighting", newYork, 5787, hdate.Tishrei,
			"14 Tishrei 5787 | Fri Sep 25 2026 | 18:30"},
		// Maharil's time is Sunday Apr 3 1977, 15 Nisan (Pesach I).
		// @hebcal/core says 18:03; its sunset that day is a minute later.
		{"Shabbat followed by Pesach moves back to Friday", newYork, 5737, hdate.Nisan,
			"13 Nisan 5737 | Fri Apr 01 1977 | 18:02"},
		// Jerusalem candle-lighting is 40 minutes before sunset
		{"Yom Tov in Israel uses Israeli candle-lighting", jerusalem, 5784, hdate.Nisan,
			"14 Nisan 5784 | Mon Apr 22 2024 | 18:33"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ev := kiddushLevanaMonth(t, tt.loc, tt.hyear, tt.hmonth)
			assert.Equal(t, tt.want, klSummary(ev))
		})
	}
}

func TestKiddushLevanaMatchesCandleLighting(t *testing.T) {
	newYork := zmanim.LookupCity("New York")
	ev := kiddushLevanaMonth(t, newYork, 5787, hdate.Tishrei)
	opts := CalOptions{
		Start:          ev.Date,
		End:            ev.Date,
		Location:       newYork,
		CandleLighting: true,
		NoHolidays:     true,
	}
	events, err := HebrewCalendar(&opts)
	require.NoError(t, err)
	require.Len(t, events, 1)
	candles := events[0].(TimedEvent)
	assert.Equal(t, "Candle lighting", candles.Desc)
	assert.Equal(t, candles.EventTime, ev.EventTime)
}

func TestKiddushLevanaMemo(t *testing.T) {
	ev := kiddushLevanaMonth(t, zmanim.LookupCity("New York"), 5787, hdate.Tishrei)
	want := strings.Join([]string{
		"Earliest Kiddush Levana (3 days): Mon, Sep 14, 2026, 2:39pm",
		"Earliest Kiddush Levana (7 days): Fri, Sep 18, 2026, 2:39pm",
		"Latest Kiddush Levana (Maharil): Sat, Sep 26, 2026, 9:00am",
		"Latest Kiddush Levana (15 days): Sat, Sep 26, 2026, 2:38pm",
	}, "\n")
	assert.Equal(t, want, ev.Memo("en"))
	assert.Contains(t, ev.Memo("ashkenazi"), "Latest Kiddush Levanah (Maharil): ")
}

func TestKiddushLevanaEvent(t *testing.T) {
	assert := assert.New(t)
	ev := kiddushLevanaMonth(t, zmanim.LookupCity("New York"), 5787, hdate.Cheshvan)
	assert.Equal("Latest Kiddush Levana: 9:44pm", ev.Render("en"))
	assert.Equal("Latest Kiddush Levanah: 9:44pm", ev.Render("ashkenazi"))
	assert.Equal("סוֹף זְמַן קִדּוּשׁ לְבָנָה: 9:44pm", ev.Render("he"))
	assert.Equal([]string{"zmanim", "kiddushLevana"}, ev.GetCategories())
	assert.Equal("🌔", ev.GetEmoji())
	assert.Equal(event.KIDDUSH_LEVANA, ev.GetFlags())
	assert.Equal("KIDDUSH_LEVANA", ev.GetFlags().String())
	assert.Equal(molad.New(5787, hdate.Cheshvan), ev.Molad)
}

func TestKiddushLevanaOnePerMonth(t *testing.T) {
	opts := CalOptions{
		Year:                 5787,
		IsHebrewYear:         true,
		Location:             zmanim.LookupCity("New York"),
		KiddushLevanaMaharil: true,
		NoHolidays:           true,
	}
	events, err := HebrewCalendar(&opts)
	require.NoError(t, err)
	assert.Len(t, events, 13)
	for _, ev := range events {
		assert.IsType(t, KiddushLevanaEvent{}, ev)
	}
}

func TestKiddushLevanaOnlyInRange(t *testing.T) {
	count := func(start, end int) int {
		opts := CalOptions{
			Start:                hdate.New(5787, hdate.Tishrei, start),
			End:                  hdate.New(5787, hdate.Tishrei, end),
			Location:             zmanim.LookupCity("New York"),
			KiddushLevanaMaharil: true,
			NoHolidays:           true,
		}
		events, err := HebrewCalendar(&opts)
		require.NoError(t, err)
		return len(events)
	}
	assert.Equal(t, 0, count(1, 13))
	assert.Equal(t, 0, count(15, 29))
	assert.Equal(t, 1, count(14, 14))
}

func TestKiddushLevanaMask(t *testing.T) {
	opts := CalOptions{
		Start:    hdate.New(5787, hdate.Cheshvan, 1),
		End:      hdate.New(5787, hdate.Cheshvan, 29),
		Location: zmanim.LookupCity("New York"),
		Mask:     event.KIDDUSH_LEVANA,
	}
	events, err := HebrewCalendar(&opts)
	require.NoError(t, err)
	require.Len(t, events, 1)
	assert.Equal(t, "21:44", events[0].(KiddushLevanaEvent).EventTime.Format("15:04"))
}

func TestKiddushLevanaRequiresLocation(t *testing.T) {
	_, err := HebrewCalendar(&CalOptions{KiddushLevanaMaharil: true})
	assert.EqualError(t, err, "opts.KiddushLevanaMaharil requires opts.Location")
}

func TestKiddushLevanaIgnoresCandleOnlyOptions(t *testing.T) {
	// Without CandleLighting these options do nothing, so conflicting
	// pairs are ignored, as they are when no option needs a location
	conflicting := func(opts CalOptions) *CalOptions {
		opts.HavdalahMins = 50
		opts.HavdalahDeg = 8.5
		opts.FastStartDeg = 16.1
		opts.FastStartMins = 72
		opts.FastEndDeg = 7.083
		opts.FastEndMins = 20
		opts.TishaBavEndDeg = 6.45
		opts.TishaBavEndMins = 50
		return &opts
	}
	newYork := zmanim.LookupCity("New York")
	events, err := HebrewCalendar(conflicting(CalOptions{
		Start:                hdate.New(5787, hdate.Cheshvan, 1),
		End:                  hdate.New(5787, hdate.Cheshvan, 29),
		Location:             newYork,
		KiddushLevanaMaharil: true,
		NoHolidays:           true,
	}))
	require.NoError(t, err)
	require.Len(t, events, 1)
	assert.Equal(t, "21:44", events[0].(KiddushLevanaEvent).EventTime.Format("15:04"))
	_, err = HebrewCalendar(conflicting(CalOptions{
		Year:                 5787,
		IsHebrewYear:         true,
		Location:             newYork,
		KiddushLevanaMaharil: true,
		CandleLighting:       true,
	}))
	assert.EqualError(t, err, "opts.HavdalahMins and opts.HavdalahDeg are mutually exclusive")
}

func TestKiddushLevanaSortedByTime(t *testing.T) {
	boston := zmanim.LookupCity("Boston")
	render := func(start, end hdate.HDate) []string {
		events, err := HebrewCalendar(&CalOptions{
			Start:                start,
			End:                  end,
			Location:             boston,
			CandleLighting:       true,
			KiddushLevanaMaharil: true,
		})
		require.NoError(t, err)
		result := make([]string, len(events))
		for i, ev := range events {
			result[i] = ev.GetDate().String() + " " + ev.Render("en")
		}
		return result
	}
	// late at night, after Havdalah
	assert.Equal(t, []string{
		"12 Adar I 5787 Candle lighting: 5:03pm",
		"13 Adar I 5787 Havdalah: 6:05pm",
		"13 Adar I 5787 Latest Kiddush Levana: 11:40pm",
		"14 Adar I 5787 Purim Katan",
	}, render(hdate.New(5787, hdate.Adar1, 12), hdate.New(5787, hdate.Adar1, 14)))
	// at dawn, before the fast begins
	assert.Equal(t, []string{
		"13 Adar 5780 Erev Purim",
		"13 Adar 5780 Latest Kiddush Levana: 5:42am",
		"13 Adar 5780 Fast begins: 5:43am",
		"13 Adar 5780 Ta'anit Esther",
		"13 Adar 5780 Fast ends: 7:18pm",
	}, render(hdate.New(5780, hdate.Adar1, 13), hdate.New(5780, hdate.Adar1, 13)))
}
