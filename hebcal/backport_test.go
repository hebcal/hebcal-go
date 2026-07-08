package hebcal

import (
	"testing"
	"time"

	"github.com/hebcal/hdate"
	"github.com/hebcal/hebcal-go/event"
	"github.com/hebcal/hebcal-go/zmanim"
)

// shabbatWeekEvents generates the events for the Shabbat of 2026-06-13 (NYC),
// which reads Parashat Sh'lach and blesses the upcoming month of Tammuz.
func shabbatWeekEvents(t *testing.T) []event.CalEvent {
	t.Helper()
	loc := zmanim.LookupCity("New York")
	opts := &CalOptions{
		Start:            hdate.FromGregorian(2026, time.June, 12),
		End:              hdate.FromGregorian(2026, time.June, 13),
		CandleLighting:   true,
		Sedrot:           true,
		ShabbatMevarchim: true,
		Location:         loc,
	}
	events, err := HebrewCalendar(opts)
	if err != nil {
		t.Fatalf("HebrewCalendar: %v", err)
	}
	return events
}

// TestCandleParshaLink verifies that erev-Shabbat candle-lighting links to the
// upcoming week's parsha (backport of the candle->parsha memo).
func TestCandleParshaLink(t *testing.T) {
	var found bool
	for _, ev := range shabbatWeekEvents(t) {
		te, ok := ev.(TimedEvent)
		if !ok || te.Desc != "Candle lighting" {
			continue
		}
		found = true
		if te.LinkedEvent == nil {
			t.Fatal("Candle lighting has no LinkedEvent")
		}
		if got := te.LinkedEvent.Render("en"); got != "Parashat Sh'lach" {
			t.Errorf("candle LinkedEvent.Render = %q, want Parashat Sh'lach", got)
		}
	}
	if !found {
		t.Fatal("no Candle lighting event generated")
	}
}

// TestMevarchimChodeshEvent verifies the upgraded Shabbat Mevarchim event type,
// its renderings, and the molad it carries.
func TestMevarchimChodeshEvent(t *testing.T) {
	var mev event.MevarchimChodeshEvent
	var found bool
	for _, ev := range shabbatWeekEvents(t) {
		if m, ok := ev.(event.MevarchimChodeshEvent); ok {
			mev, found = m, true
		}
	}
	if !found {
		t.Fatal("no MevarchimChodeshEvent generated")
	}
	if mev.MonthName != "Tammuz" {
		t.Errorf("MonthName = %q, want Tammuz", mev.MonthName)
	}
	if got := mev.Render("en"); got != "Shabbat Mevarchim Chodesh Tammuz" {
		t.Errorf("Render = %q", got)
	}
	if got := mev.RenderBrief("en"); got != "Mevarchim Chodesh Tammuz" {
		t.Errorf("RenderBrief = %q", got)
	}
	if got := mev.GetCategories(); len(got) != 1 || got[0] != "mevarchim" {
		t.Errorf("GetCategories = %v", got)
	}
	// molad of Tammuz 5786: Monday, 6h 46m 16 chalakim
	if mev.Molad.Date.Weekday() != time.Monday || mev.Molad.Hours != 6 ||
		mev.Molad.Minutes != 46 || mev.Molad.Chalakim != 16 {
		t.Errorf("Molad = %+v, want Mon 6:46 + 16 chalakim", mev.Molad)
	}
}
