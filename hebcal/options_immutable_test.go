package hebcal_test

import (
	"reflect"
	"sync"
	"testing"

	"github.com/hebcal/hdate"
	"github.com/hebcal/hebcal-go/event"
	"github.com/hebcal/hebcal-go/hebcal"
	"github.com/hebcal/hebcal-go/zmanim"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// normalizingOptions returns options that exercise every normalization
// HebrewCalendar performs: an Israeli location (which implies the Israel
// schedule and its own candle-lighting offset), a mask that switches event
// streams on, default Havdalah and Yerushalmi settings, and negative offsets.
func normalizingOptions() hebcal.CalOptions {
	return hebcal.CalOptions{
		Year:             2026,
		Location:         zmanim.LookupCity("Jerusalem"),
		CandleLighting:   true,
		FastStartMins:    -72,
		FastEndDeg:       -6.45,
		TishaBavEndMins:  -30,
		NoRoshChodesh:    true,
		NoModern:         true,
		NoMinorFast:      true,
		NoSpecialShabbat: true,
		YerushalmiYomi:   true,
		DailyLearning:    []string{"dafYomi"},
		Yahrzeits:        []hebcal.UserYahrzeit{},
		UserEvents:       []hebcal.UserEventSpec{{Month: hdate.Nisan, Day: 20, Desc: "Birthday"}},
		Mask: event.RoshChodesh | event.ModernHoliday | event.MinorFast |
			event.SpecialShabbat | event.ParshaHashavua | event.DafYomi |
			event.OmerCount | event.ShabbatMevarchim | event.MishnaYomi |
			event.NachYomi | event.YomKippurKatan | event.BeHaB |
			event.YerushalmiYomi | event.KiddushLevana | event.Yizkor |
			event.LightCandles | event.LightCandlesTzeis | event.YomTovEnds,
	}
}

// HebrewCalendar works on a copy of its options: the caller's struct, and
// the location it points to, are the same after the call as before.
func TestHebrewCalendarDoesNotModifyOptions(t *testing.T) {
	opts := normalizingOptions()
	loc := *opts.Location
	opts.Location = &loc
	want := opts
	wantLoc := loc

	events, err := hebcal.HebrewCalendar(&opts)
	require.NoError(t, err)
	require.NotEmpty(t, events)
	assert.True(t, reflect.DeepEqual(want, opts), "options changed:\n got %+v\nwant %+v", opts, want)
	assert.Equal(t, wantLoc, loc)

	// An early year switches candle-lighting off internally; not for the caller.
	early := hebcal.CalOptions{Year: 1700, CandleLighting: true, Location: zmanim.LookupCity("Boston")}
	wantEarly := early
	_, err = hebcal.HebrewCalendar(&early)
	require.NoError(t, err)
	assert.True(t, reflect.DeepEqual(wantEarly, early))
}

// Events keep the clock format they were generated with, not a reference to
// the options: changing the options afterwards does not change how they render.
func TestEventsDoNotFollowOptions(t *testing.T) {
	opts := hebcal.CalOptions{
		Year: 2026, Month: 4, Location: zmanim.LookupCity("Chicago"),
		CandleLighting: true, Molad: true, ShabbatMevarchim: true,
		KiddushLevanaMaharil: true, SunriseSunset: true,
	}
	events, err := hebcal.HebrewCalendar(&opts)
	require.NoError(t, err)
	render := func() []string {
		out := make([]string, 0, 2*len(events))
		for _, ev := range events {
			out = append(out, ev.Render("en"), event.Memo(ev, "en"))
		}
		return out
	}
	before := render()
	opts.Clock = event.Clock24
	opts.Location = zmanim.LookupCity("London")
	opts.IL = true
	assert.Equal(t, before, render())
}

// One options value can be shared by concurrent calls. Run with -race.
func TestHebrewCalendarConcurrent(t *testing.T) {
	opts := normalizingOptions()
	want, err := hebcal.HebrewCalendar(&opts)
	require.NoError(t, err)
	wantStr := renderAll(want)

	const goroutines = 8
	results := make([][]string, goroutines)
	errs := make([]error, goroutines)
	var wg sync.WaitGroup
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			events, err := hebcal.HebrewCalendar(&opts)
			errs[i] = err
			results[i] = renderAll(events)
		}(i)
	}
	wg.Wait()
	for i := 0; i < goroutines; i++ {
		require.NoError(t, errs[i])
		assert.Equal(t, wantStr, results[i])
	}
}

func renderAll(events []event.CalEvent) []string {
	out := make([]string, len(events))
	for i, ev := range events {
		out[i] = ev.Date().String() + " " + ev.Render("en")
	}
	return out
}

// A nil options pointer means the zero CalOptions.
func TestHebrewCalendarNilOptions(t *testing.T) {
	events, err := hebcal.HebrewCalendar(nil)
	require.NoError(t, err)
	want, err := hebcal.HebrewCalendar(&hebcal.CalOptions{})
	require.NoError(t, err)
	assert.Equal(t, len(want), len(events))
}
