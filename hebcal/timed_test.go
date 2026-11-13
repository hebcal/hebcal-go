package hebcal_test

import (
	"testing"

	"github.com/hebcal/hebcal-go/event"
	"github.com/hebcal/hebcal-go/hebcal"
	"github.com/hebcal/hebcal-go/zmanim"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Every event with a clock time satisfies event.Timed, Kiddush Levana
// included, and its kind agrees with its untranslated description.
func TestTimedInterfaceAndKind(t *testing.T) {
	opts := hebcal.CalOptions{
		Year: 2026, Location: zmanim.LookupCity("Chicago"), CandleLighting: true,
		KiddushLevanaMaharil: true, Sedrot: true, Yizkor: true, DailyZmanim: true,
	}
	events, err := hebcal.HebrewCalendar(&opts)
	require.NoError(t, err)

	wantDesc := map[hebcal.TimedKind]string{
		hebcal.CandleLighting: "Candle lighting",
		hebcal.Havdalah:       "Havdalah",
		hebcal.FastBegins:     "Fast begins",
		hebcal.FastEnds:       "Fast ends",
		hebcal.FinishChametz:  "Finish eating chametz",
		hebcal.BiurChametz:    "Biur Chametz",
		hebcal.KiddushLevana:  "Latest Kiddush Levana",
	}
	seen := map[hebcal.TimedKind]int{}
	for _, ev := range events {
		timed, ok := ev.(event.Timed)
		kinded, hasKind := ev.(interface{ Kind() hebcal.TimedKind })
		assert.Equal(t, ok, hasKind, ev.Desc())
		if !ok {
			_, isKind := hebcal.TimedKindOf(ev)
			assert.False(t, isKind, ev.Desc())
			continue
		}
		assert.False(t, timed.Time().IsZero(), ev.Desc())
		kind := kinded.Kind()
		k, ok := hebcal.TimedKindOf(ev)
		assert.True(t, ok)
		assert.Equal(t, kind, k)
		seen[kind]++
		if want, ok := wantDesc[kind]; ok {
			assert.Equal(t, want, ev.Desc())
		}
		switch kind {
		case hebcal.ChanukahCandles:
			assert.Contains(t, ev.Desc(), "Chanukah: ")
			assert.True(t, ev.Flags().Has(event.ChanukahCandles))
			assert.Equal(t, "Chanukah", ev.Basename())
		case hebcal.Zman:
			assert.Equal(t, event.Zmanim, ev.Flags())
		case hebcal.KiddushLevana:
			_, isKL := ev.(hebcal.KiddushLevanaEvent)
			assert.True(t, isKL)
		}
		// the deprecated helper finds the same events
		te, ok := hebcal.AsTimedEvent(ev)
		require.True(t, ok)
		assert.Equal(t, timed.Time(), te.Time())
		assert.Equal(t, ev.Desc(), te.Desc())
	}
	for kind := hebcal.CandleLighting; kind <= hebcal.Zman; kind++ {
		assert.NotZero(t, seen[kind], kind.String())
	}
	assert.Zero(t, seen[hebcal.TimedOther])
}

// Linked events are reached through event.Linked, whatever their type.
func TestLinkedInterface(t *testing.T) {
	opts := hebcal.CalOptions{
		Year: 2026, Location: zmanim.LookupCity("Chicago"), CandleLighting: true,
		Sedrot: true, Yizkor: true,
	}
	events, err := hebcal.HebrewCalendar(&opts)
	require.NoError(t, err)
	linked := map[string]string{}
	for _, ev := range events {
		if l, ok := ev.(event.Linked); ok && l.LinkedEvent() != nil {
			y, m, d := ev.Date().Greg()
			if y == 2026 && (m == 9 || m == 10 || m == 7) {
				linked[ev.Date().String()+" "+ev.Desc()] = l.LinkedEvent().Desc()
			}
			_ = d
		}
	}
	assert.Equal(t, "Yom Kippur", linked["10 Tishrei 5787 Yizkor"])
	assert.Equal(t, "Erev Yom Kippur", linked["9 Tishrei 5787 Candle lighting"])
	assert.Equal(t, "Yom Kippur", linked["10 Tishrei 5787 Havdalah"])
	assert.Equal(t, "Tish'a B'Av", linked["9 Av 5786 Fast ends"])
	assert.Equal(t, "Erev Sukkot", linked["14 Tishrei 5787 Candle lighting"])
	assert.Equal(t, "Parashat Pinchas", linked["18 Tammuz 5786 Candle lighting"])
}

func TestTimedKindString(t *testing.T) {
	assert.Equal(t, "CandleLighting", hebcal.CandleLighting.String())
	assert.Equal(t, "Zman", hebcal.Zman.String())
	assert.Equal(t, "TimedOther", hebcal.TimedKind(0).String())
	assert.Equal(t, "TimedKind(99)", hebcal.TimedKind(99).String())
}
