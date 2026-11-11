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

func calendar(t *testing.T, opts *hebcal.CalOptions) []event.CalEvent {
	t.Helper()
	events, err := hebcal.HebrewCalendar(opts)
	assert.NoError(t, err)
	return events
}

func find(events []event.CalEvent, desc string) event.CalEvent {
	for _, ev := range events {
		if event.Desc(ev) == desc {
			return ev
		}
	}
	return nil
}

func TestTimedEventEmoji(t *testing.T) {
	events := calendar(t, &hebcal.CalOptions{
		Year: 5787, IsHebrewYear: true,
		CandleLighting: true,
		Location:       zmanim.LookupCity("Boston"),
	})
	var candles, havdalah, yomTovHavdalah, erevYomTov int
	for _, ev := range events {
		te, ok := hebcal.AsTimedEvent(ev)
		if !ok {
			continue
		}
		dow := te.Date.Weekday()
		switch te.Desc {
		case "Candle lighting":
			candles++
			assert.Equal(t, "🕯️", ev.GetEmoji(), te.Date.String())
			if dow != time.Friday {
				erevYomTov++
			}
		case "Havdalah":
			havdalah++
			assert.Equal(t, "✨", ev.GetEmoji(), te.Date.String())
			if dow != time.Saturday {
				yomTovHavdalah++
			}
		}
	}
	assert.Greater(t, candles, 50)
	assert.Greater(t, havdalah, 50)
	assert.Greater(t, erevYomTov, 0)
	assert.Greater(t, yomTovHavdalah, 0)
}

func TestChanukahCandleLightingIsTheHoliday(t *testing.T) {
	events := calendar(t, &hebcal.CalOptions{
		Start:          hdate.New(5783, hdate.Kislev, 25),
		End:            hdate.New(5783, hdate.Kislev, 29), // Friday 6th candle
		CandleLighting: true,
		Location:       zmanim.LookupCity("Jerusalem"),
	})
	for _, desc := range []string{"Chanukah: 2 Candles", "Chanukah: 6 Candles"} {
		ev := find(events, desc)
		if assert.NotNil(t, ev, desc) {
			assert.Equal(t, "Chanukah", ev.Basename())
			assert.Equal(t, []string{"holiday", "minor"}, ev.GetCategories())
			assert.Equal(t, "https://www.hebcal.com/holidays/chanukah-2022", event.URL(ev))
		}
	}
}

func TestCandleLightingMemo(t *testing.T) {
	opts := func(start, end hdate.HDate) *hebcal.CalOptions {
		return &hebcal.CalOptions{Start: start, End: end, CandleLighting: true,
			Sedrot: true, Location: zmanim.LookupCity("Boston")}
	}
	memos := func(events []event.CalEvent, locale string) []string {
		var out []string
		for _, ev := range events {
			if te, ok := ev.(hebcal.TimedEvent); ok && te.Desc == "Candle lighting" {
				out = append(out, event.Memo(ev, locale))
			}
		}
		return out
	}
	// Friday 23 Nisan 5787: the coming Shabbat's portion
	evs := calendar(t, opts(hdate.New(5787, hdate.Nisan, 23), hdate.New(5787, hdate.Nisan, 23)))
	assert.Equal(t, []string{"Parashat Achrei Mot"}, memos(evs, "en"))
	assert.Equal(t, []string{"פָּרָשַׁת אַחֲרֵי מוֹת"}, memos(evs, "he"))
	// Friday 18 Tishrei 5786, Chol HaMoed Sukkot: the festival reading
	evs = calendar(t, opts(hdate.New(5786, hdate.Tishrei, 18), hdate.New(5786, hdate.Tishrei, 18)))
	assert.Equal(t, []string{"Sukkot Shabbat Chol ha-Moed"}, memos(evs, "en"))
	// Candle-lighting before Yom Tov has no memo of its own
	evs = calendar(t, opts(hdate.New(5787, hdate.Nisan, 14), hdate.New(5787, hdate.Nisan, 14)))
	assert.Equal(t, []string{""}, memos(evs, "en"))
}

func TestMevarchimAndYKKMemo(t *testing.T) {
	evs := calendar(t, &hebcal.CalOptions{
		Start:            hdate.New(5787, hdate.Cheshvan, 1),
		End:              hdate.New(5787, hdate.Cheshvan, 30),
		ShabbatMevarchim: true,
		YomKippurKatan:   true,
		NoHolidays:       true,
		Location:         zmanim.LookupCity("London"),
	})
	mev := find(evs, "Shabbat Mevarchim Chodesh Kislev")
	if assert.NotNil(t, mev) {
		assert.Equal(t, "Molad Kislev: Monday, 22:27 and 3 chalakim", event.Memo(mev, "en"))
		assert.Equal(t, "Molad Kislav: Lundi, 22:27 et 3 'halakim", event.Memo(mev, "fr"))
	}
	ykk := find(evs, "Yom Kippur Katan Kislev")
	if assert.NotNil(t, ykk) {
		assert.Equal(t, "Minor Day of Atonement on the day preceeding Rosh Chodesh Kislev", event.Memo(ykk, "en"))
	}
}

func TestIsraelCandleMins(t *testing.T) {
	friday := hdate.New(5787, hdate.Cheshvan, 19)
	candles := func(loc zmanim.Location) string {
		evs := calendar(t, &hebcal.CalOptions{Start: friday, End: friday, CandleLighting: true, Location: &loc})
		return evs[0].Render("en")
	}
	jlm := *zmanim.LookupCity("Jerusalem")
	want := candles(jlm)
	byName := jlm
	byName.Name = "Jerusalem, Israel"
	assert.Equal(t, want, candles(byName))
	byID := jlm
	byID.Name = "Yerushalayim"
	byID.GeoNameID = 281184
	assert.Equal(t, want, candles(byID))
	other := jlm
	other.Name = "Somewhere, Israel"
	assert.NotEqual(t, want, candles(other))
}

func TestDescAndBasename(t *testing.T) {
	evs := calendar(t, &hebcal.CalOptions{
		Start: hdate.New(5786, hdate.Nisan, 16), End: hdate.New(5786, hdate.Nisan, 22),
		Omer: true, Sedrot: true, Molad: true,
	})
	var got []string
	for _, ev := range evs {
		got = append(got, event.Desc(ev)+" | "+ev.Basename())
	}
	assert.Contains(t, got, "Omer 1 | Omer 1")
	assert.Contains(t, got, "Pesach II | Pesach")
	m := calendar(t, &hebcal.CalOptions{Start: hdate.New(5787, hdate.Cheshvan, 27), End: hdate.New(5787, hdate.Cheshvan, 27), Molad: true, NoHolidays: true})
	if assert.Len(t, m, 1) {
		assert.Equal(t, "Molad Kislev 5787", event.Desc(m[0]))
	}
}

func TestYomKippurKatanMemoSmartApostrophe(t *testing.T) {
	evs, err := hebcal.HebrewCalendar(&hebcal.CalOptions{
		Start:          hdate.New(5786, hdate.Tevet, 1),
		End:            hdate.New(5786, hdate.Shvat, 30),
		YomKippurKatan: true,
		NoHolidays:     true,
	})
	assert.NoError(t, err)
	ykk := find(evs, "Yom Kippur Katan Sh'vat")
	if assert.NotNil(t, ykk) {
		assert.Equal(t, "Minor Day of Atonement on the day preceeding Rosh Chodesh Sh’vat", event.Memo(ykk, "en"))
	}
}
