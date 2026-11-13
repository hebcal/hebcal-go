package hebcal

// Hebcal - A Jewish Calendar Generator
// Copyright (c) 2022 Michael J. Radwin
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

import (
	"fmt"
	"strings"
	"time"

	"github.com/hebcal/hdate"

	"github.com/hebcal/hebcal-go/event"
	"github.com/hebcal/hebcal-go/zmanim"
	"github.com/hebcal/locales"
)

// TimedEvent is used for Candle-lighting, Havdalah, and fast start/end
type TimedEvent struct {
	event.HolidayEvent
	EventTime    time.Time
	LinkedEvent  event.CalEvent
	sunsetOffset int
	clock        clock
	// chagReading is the untranslated name of the festival reading that
	// displaces the weekly portion, for erev-Shabbat candle-lighting in a
	// festival week (see Memo).
	chagReading string
}

// AsTimedEvent returns the TimedEvent ev is or embeds: a TimedEvent itself,
// or a KiddushLevanaEvent.
func AsTimedEvent(ev event.CalEvent) (TimedEvent, bool) {
	if t, ok := ev.(interface{ timedEvent() TimedEvent }); ok {
		return t.timedEvent(), true
	}
	return TimedEvent{}, false
}

func (ev TimedEvent) timedEvent() TimedEvent { return ev }

// isChanukah reports whether ev is the candle-lighting of a night of
// Chanukah, which stands for the holiday itself.
func (ev TimedEvent) isChanukah() bool {
	return ev.Flags.Has(event.CHANUKAH_CANDLES) && ev.LinkedEvent != nil
}

// Memo returns, for erev-Shabbat candle-lighting, the coming Shabbat's Torah
// portion ("Parashat Vayera"), or the festival reading that displaces it
// ("Pesach I"). Other timed events have no memo of their own.
func (ev TimedEvent) Memo(locale string) string {
	if ev.Desc != "Candle lighting" {
		return ""
	}
	if ev.chagReading != "" {
		s, _ := locales.LookupTranslation(ev.chagReading, locale)
		return s
	}
	if ev.LinkedEvent != nil && ev.LinkedEvent.GetFlags() == event.PARSHA_HASHAVUA {
		return ev.LinkedEvent.Render(locale)
	}
	return ""
}

// GetDesc returns the untranslated description, e.g. "Candle lighting" or
// "Chanukah: 3 Candles".
func (ev TimedEvent) GetDesc() string {
	return ev.Desc
}

// NewTimedEvent constructs a timed event. opts chooses the clock format its
// time is rendered in (see ReformatTimeStr) and may be nil; the event keeps
// that format, not a reference to opts.
func NewTimedEvent(hd hdate.HDate, desc string, flags event.HolidayFlags, t time.Time,
	sunsetOffset int, linkedEvent event.CalEvent, opts *CalOptions) TimedEvent {
	return newTimedEvent(hd, desc, flags, t, sunsetOffset, linkedEvent, clockFor(opts))
}

func newTimedEvent(hd hdate.HDate, desc string, flags event.HolidayFlags, t time.Time,
	sunsetOffset int, linkedEvent event.CalEvent, clock clock) TimedEvent {
	if t.IsZero() {
		return TimedEvent{}
	}
	var emoji string
	switch desc {
	case "Candle lighting":
		emoji = "🕯️"
	case "Havdalah":
		emoji = "✨"
	}
	// Round to the nearest minute, matching @hebcal/core's TimedEvent
	// (this.eventTime = Zmanim.roundTime(eventTime)). SunsetOffset already
	// yields a whole minute, so in practice this only affects the degree-based
	// tzeit times (havdalah, fast begin/end, chametz). The ZMANIM daily-times
	// feature has no @hebcal/core equivalent and keeps its unrounded times.
	if flags != event.ZMANIM {
		t = roundToMinute(t)
	}
	return TimedEvent{
		HolidayEvent: event.HolidayEvent{
			Date:  hd,
			Desc:  desc,
			Flags: flags,
			Emoji: emoji,
		},
		EventTime:    t,
		LinkedEvent:  linkedEvent,
		clock:        clock,
		sunsetOffset: sunsetOffset,
	}
}

// roundToMinute rounds t to the nearest minute (>= 30 seconds rounds up),
// matching @hebcal/core Zmanim.roundTime.
func roundToMinute(t time.Time) time.Time {
	if t.IsZero() {
		return t
	}
	if t.Second() == 0 && t.Nanosecond() == 0 {
		return t
	}
	year, month, day := t.Date()
	hour, min, sec := t.Clock()
	if sec >= 30 {
		min++
	}
	return time.Date(year, month, day, hour, min, 0, 0, t.Location())
}

func (ev TimedEvent) GetDate() hdate.HDate {
	return ev.Date
}

func (ev TimedEvent) Render(locale string) string {
	timeStr := ev.clock.formatTime(ev.EventTime)
	return fmt.Sprintf("%s: %s", ev.RenderBrief(locale), timeStr)
}

// RenderBrief returns the translated description without its time, e.g.
// "Candle lighting" or "Havdalah (50 min)".
func (ev TimedEvent) RenderBrief(locale string) string {
	desc, _ := locales.LookupTranslation(ev.Desc, locale)
	if ev.Desc == "Havdalah" && ev.sunsetOffset != 0 {
		minStr, _ := locales.LookupTranslation("min", locale)
		desc = fmt.Sprintf("%s (%d %s)", desc, ev.sunsetOffset, minStr)
	}
	return desc
}

func (ev TimedEvent) GetFlags() event.HolidayFlags {
	return ev.Flags
}

func (ev TimedEvent) GetEmoji() string {
	return ev.Emoji
}

// Basename returns the description, or for Chanukah candle-lighting the
// holiday's basename ("Chanukah").
func (ev TimedEvent) Basename() string {
	if ev.isChanukah() {
		return ev.LinkedEvent.Basename()
	}
	return ev.Desc
}

// GetCategories returns the category and sub-categories for a timed event,
// keyed by its description, matching TimedEvent.getCategories() in
// @hebcal/core.
//
// Chanukah candle-lighting takes the holiday's categories.
func (ev TimedEvent) GetCategories() []string {
	if ev.isChanukah() {
		return ev.LinkedEvent.GetCategories()
	}
	switch ev.Desc {
	case "Candle lighting":
		return []string{"candles"}
	case "Havdalah":
		return []string{"havdalah"}
	case "Fast begins", "Fast ends":
		return []string{"zmanim", "fast"}
	case "Finish eating chametz":
		return []string{"zmanim", "achilasChametz"}
	case "Biur Chametz":
		return []string{"zmanim", "biurChametz"}
	}
	return []string{"unknown"}
}

// newZmanim builds a Zmanim for the Gregorian date of hd at the options'
// Location, honoring opts.UseElevation so that sunrise/sunset-based times
// account for the location's elevation when requested. Degree-based zmanim are
// never affected by elevation.
func newZmanim(hd hdate.HDate, opts *calConfig) zmanim.Zmanim {
	year, month, day := hd.Greg()
	gregDate := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	z := zmanim.New(opts.Location, gregDate)
	z.UseElevation = opts.UseElevation
	return z
}

func makeCandleEvent(hd hdate.HDate, opts *calConfig, ev event.CalEvent) TimedEvent {
	havdalahTitle := false
	useHavdalahOffset := false
	dow := hd.Weekday()
	if dow == time.Saturday {
		useHavdalahOffset = true
	}
	flags := event.LIGHT_CANDLES
	if ev != nil {
		flags = ev.GetFlags()
		if dow != time.Friday {
			if flags.HasAny(event.LIGHT_CANDLES_TZEIS | event.CHANUKAH_CANDLES) {
				useHavdalahOffset = true
			} else if flags.Has(event.YOM_TOV_ENDS) {
				havdalahTitle = true
				useHavdalahOffset = true
			}
		}
	} else if dow == time.Saturday {
		havdalahTitle = true
		flags = event.LIGHT_CANDLES_TZEIS
	}
	// A zero Havdalah offset means tzeit. A zero candle-lighting offset
	// means sunset itself (CandleLightingAtSunset).
	offset := opts.CandleLightingMins
	if useHavdalahOffset {
		offset = opts.HavdalahMins
	}
	z := newZmanim(hd, opts)
	var eventTime time.Time
	if useHavdalahOffset && offset == 0 {
		eventTime = z.Tzeit(opts.HavdalahDeg)
	} else {
		eventTime = z.SunsetOffset(offset, true)
	}
	if eventTime.IsZero() {
		return TimedEvent{} // no sunset
	}
	desc := "Candle lighting"
	if havdalahTitle {
		desc = "Havdalah"
	} else if !flags.Has(event.CHANUKAH_CANDLES) {
		// Exactly one of the two candle-lighting flags, chosen by how the
		// time was computed: LIGHT_CANDLES for minutes before sunset
		// (Shabbat or the start of a chag), LIGHT_CANDLES_TZEIS for
		// nightfall / Havdalah time (the second night of Yom Tov, or a chag
		// that begins Saturday night). Chanukah keeps its own flags: its
		// Friday and Saturday events reuse this event as the Chanukah event.
		if useHavdalahOffset {
			flags = flags.Without(event.LIGHT_CANDLES).With(event.LIGHT_CANDLES_TZEIS)
		} else {
			flags = flags.Without(event.LIGHT_CANDLES_TZEIS).With(event.LIGHT_CANDLES)
		}
	}
	return newTimedEvent(hd, desc, flags, eventTime, offset, ev, opts.clock)
}

func makeChanukahCandleLighting(ev event.HolidayEvent, opts *calConfig) TimedEvent {
	hd := ev.Date
	dow := hd.Weekday()
	if dow == time.Friday || dow == time.Saturday {
		timedEv := makeCandleEvent(hd, opts, ev)
		if (timedEv == TimedEvent{}) {
			return timedEv // no sunset
		}
		timedEv.Desc = ev.Desc
		timedEv.Emoji = ev.Emoji
		timedEv.ChanukahDay = ev.ChanukahDay
		return timedEv
	}
	z := newZmanim(hd, opts)
	candleLightingTime := z.BeinHashmashos()
	if candleLightingTime.IsZero() {
		return TimedEvent{} // no sunset
	}
	return TimedEvent{
		HolidayEvent: ev,
		// Round to the nearest minute, as NewTimedEvent does and as
		// @hebcal/core's TimedChanukahEvent does with Zmanim.roundTime.
		EventTime:   roundToMinute(candleLightingTime),
		LinkedEvent: ev,
		clock:       opts.clock,
	}
}

func makeFastStartEnd(ev event.CalEvent, opts *calConfig) (TimedEvent, TimedEvent) {
	hd := ev.GetDate()
	z := newZmanim(hd, opts)
	desc := event.Desc(ev)
	flags := ev.GetFlags()
	var startEvent, endEvent TimedEvent
	if desc == "Erev Tish'a B'Av" {
		sunset := z.Sunset()
		startEvent = newTimedEvent(hd, "Fast begins", flags, sunset, 0, ev, opts.clock)
	} else if strings.HasPrefix(desc, "Tish'a B'Av") {
		endEvent = newTimedEvent(hd, "Fast ends", flags, fastEndTime(&z, true, opts), 0, ev, opts.clock)
	} else {
		startEvent = newTimedEvent(hd, "Fast begins", flags, fastStartTime(&z, opts), 0, ev, opts.clock)
		if hd.Weekday() != time.Friday && !(hd.Day() == 14 && hd.Month() == hdate.Nisan) {
			endEvent = newTimedEvent(hd, "Fast ends", flags, fastEndTime(&z, false, opts), 0, ev, opts.clock)
		}
	}
	return startEvent, endEvent
}

const (
	// alot16Point1 is Alot HaShachar, the default start of a minor fast.
	alot16Point1 = 16.1
	// tzeitTucazinsky is the default end of Tish'a B'Av, 6.45 degrees as
	// calculated by Rabbi Yechiel Michel Tucazinsky.
	tzeitTucazinsky = 6.45
	// minorFastEndMinutesIL is the default end of a minor fast in Israel, 15
	// minutes after sunset (Rabbi Deblitzky's practice).
	minorFastEndMinutesIL = 15
)

// fastStartTime is when a minor fast begins: opts.FastStartMins before
// sunrise, when the sun is opts.FastStartDeg below the horizon, or at Alot
// HaShachar.
func fastStartTime(z *zmanim.Zmanim, opts *calConfig) time.Time {
	if opts.FastStartMins != 0 {
		return z.SunriseOffset(-opts.FastStartMins, true)
	}
	deg := opts.FastStartDeg
	if deg == 0 {
		deg = alot16Point1
	}
	return z.TimeAtAngle(deg, true)
}

// fastEndTime is when a fast ends. Tish'a B'Av ends opts.TishaBavEndMins
// after sunset, at tzeit opts.TishaBavEndDeg, or at tzeit 6.45 degrees, and
// ignores opts.IL. A minor fast ends opts.FastEndMins after sunset, at tzeit
// opts.FastEndDeg, or else 15 minutes after sunset in Israel and at tzeit
// 7.083 degrees elsewhere.
func fastEndTime(z *zmanim.Zmanim, isTishaBav bool, opts *calConfig) time.Time {
	if isTishaBav {
		if opts.TishaBavEndMins != 0 {
			return z.SunsetOffset(opts.TishaBavEndMins, true)
		}
		deg := opts.TishaBavEndDeg
		if deg == 0 {
			deg = tzeitTucazinsky
		}
		return z.Tzeit(deg)
	}
	if opts.FastEndMins != 0 {
		return z.SunsetOffset(opts.FastEndMins, true)
	}
	if opts.FastEndDeg != 0 {
		return z.Tzeit(opts.FastEndDeg)
	}
	if opts.IL {
		return z.SunsetOffset(minorFastEndMinutesIL, true)
	}
	return z.Tzeit(zmanim.Tzeit3MediumStars)
}

// makeErevPesachChametz returns the "Finish eating chametz" (sof zman achilat
// chametz) and "Biur Chametz" (sof zman biur chametz) events for Erev Pesach.
// It is modeled on hebcal-es6's makeErevPesachChametzEvents.
//
// Sof zman achilat chametz is the same as sof zman tefilla (Gra, sunrise plus 4
// halachic hours); sof zman biur chametz is one halachic hour later. When Erev
// Pesach falls on Shabbat, chametz cannot be burned on Shabbat, so the Biur
// Chametz event is emitted on the Friday before instead (see makeBiurChametz)
// and is omitted here.
func makeErevPesachChametz(ev event.CalEvent, opts *calConfig) []event.CalEvent {
	hd := ev.GetDate()
	z := newZmanim(hd, opts)
	events := make([]event.CalEvent, 0, 2)
	achilas := z.SofZmanAchilasChametz()
	if achilas.IsZero() {
		return events
	}
	achilasEv := newTimedEvent(hd, "Finish eating chametz", 0, achilas, 0, nil, opts.clock)
	achilasEv.Emoji = "🍞"
	events = append(events, achilasEv)
	if hd.Weekday() != time.Saturday {
		if biurEv := makeBiurChametz(hd, opts); (biurEv != TimedEvent{}) {
			events = append(events, biurEv)
		}
	}
	return events
}

// makeBiurChametz returns the "Biur Chametz" (sof zman biur chametz) event for
// the given day, or an empty TimedEvent if the sun does not rise/set.
func makeBiurChametz(hd hdate.HDate, opts *calConfig) TimedEvent {
	z := newZmanim(hd, opts)
	biur := z.SofZmanBiurChametz()
	if biur.IsZero() {
		return TimedEvent{}
	}
	biurEv := newTimedEvent(hd, "Biur Chametz", 0, biur, 0, nil, opts.clock)
	biurEv.Emoji = "🔥"
	return biurEv
}

type riseSetEvent struct {
	date    hdate.HDate
	sunrise time.Time
	sunset  time.Time
	clock   clock
}

func newRiseSetEvent(hd hdate.HDate, opts *calConfig) riseSetEvent {
	z := newZmanim(hd, opts)
	return riseSetEvent{date: hd, sunrise: z.Sunrise(), sunset: z.Sunset(), clock: opts.clock}
}

func (ev riseSetEvent) GetDate() hdate.HDate {
	return ev.date
}

func (ev riseSetEvent) Render(locale string) string {
	riseStr := ev.clock.formatTime(ev.sunrise)
	setStr := ev.clock.formatTime(ev.sunset)
	return fmt.Sprintf("Sunrise: %s; Sunset %s", riseStr, setStr)
}

func (ev riseSetEvent) GetFlags() event.HolidayFlags {
	return event.ZMANIM
}

func (ev riseSetEvent) GetEmoji() string {
	return ""
}

func (ev riseSetEvent) Basename() string {
	return ev.Render("en")
}

func (ev riseSetEvent) GetCategories() []string {
	return []string{"zmanim"}
}

func dailyZemanim(date hdate.HDate, opts *calConfig) []event.CalEvent {
	z := newZmanim(date, opts)
	times := []struct {
		desc string
		t    time.Time
	}{
		{"Alot HaShachar", z.AlotHaShachar()},
		{"Misheyakir", z.Misheyakir()},
		{"Misheyakir Machmir", z.MisheyakirMachmir()},
		{"Sunrise", z.Sunrise()},
		{"Kriat Shema, sof zeman (MGA)", z.SofZmanShmaMGA()},
		{"Kriat Shema, sof zeman (GRA)", z.SofZmanShma()},
		{"Tefilah, sof zeman (MGA)", z.SofZmanTfillaMGA()},
		{"Tefilah, sof zeman (GRA)", z.SofZmanTfilla()},
		{"Chatzot HaYom", z.Chatzot()},
		{"Mincha Gedolah", z.MinchaGedola()},
		{"Mincha Ketanah", z.MinchaKetana()},
		{"Plag HaMincha", z.PlagHaMincha()},
		{"Sunset", z.Sunset()},
		{"Bein HaShemashot", z.BeinHashmashos()},
		{"Tzeit HaKochavim", z.Tzeit(zmanim.Tzeit3SmallStars)},
	}
	events := make([]event.CalEvent, 0, len(times))
	for _, zman := range times {
		if !zman.t.IsZero() {
			ev := newTimedEvent(date, zman.desc, event.ZMANIM, zman.t, 0, nil, opts.clock)
			events = append(events, ev)
		}
	}
	return events
}

// URL returns the Chanukah page for Chanukah candle lighting, and "" for
// every other timed event.
//
// The override matters because TimedEvent embeds event.HolidayEvent and would
// otherwise inherit its URL(), producing links like /holidays/havdalah-2026
// for a description that is not a holiday name.
//
// Only Chanukah is treated as linkable. Candle lighting before Pesach also
// carries a LinkedEvent, but a clock time is not itself an observance with a
// page, and following every link would attach a holiday URL to each Friday
// evening in the calendar.
func (ev TimedEvent) URL() string {
	if !ev.Flags.Has(event.CHANUKAH_CANDLES) || ev.LinkedEvent == nil {
		return ""
	}
	return event.URL(ev.LinkedEvent)
}
