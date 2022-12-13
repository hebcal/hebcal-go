package hebcal

// Hebcal - A Jewish Calendar Generator
// Copyright (c) 2022 Michael J. Radwin
// Derived from original C version, Copyright (C) 1994-2004 Danny Sadinoff
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
	"errors"
	"math"
	"time"

	"github.com/hebcal/hebcal-go/dafyomi"
	"github.com/hebcal/hebcal-go/greg"
	"github.com/hebcal/hebcal-go/hdate"
	"github.com/hebcal/hebcal-go/mishnayomi"
	"github.com/hebcal/hebcal-go/sedra"
	"github.com/hebcal/hebcal-go/yerushalmi"
	"github.com/hebcal/hebcal-go/zmanim"
)

/*
Calculates holidays and other Hebrew calendar events based on CalOptions.

Each holiday is represented by CalEvent object which includes a date,
a description, flags and optional attributes.
If given no options, returns holidays for the Diaspora for the current Gregorian year.

The date range returned by this function can be controlled by:
  - opts.Year - Gregorian (e.g. 1993) or Hebrew year (e.g. 5749)
  - opts.IsHebrewYear - to interpret year as Hebrew year
  - opts.NumYears - generate calendar for multiple years (default 1)

Alternatively, specify start and end days HDate instances:
  - opts.Start - use specific start date (requires end date)
  - opts.End - use specific end date (requires start date)

Unless opts.NoHolidays == true, default holidays include:
  - Major holidays - Rosh Hashana, Yom Kippur, Pesach, Sukkot, etc.
  - Minor holidays - Purim, Chanukah, Tu BiShvat, Lag BaOmer, etc.
  - Minor fasts - Ta'anit Esther, Tzom Gedaliah, etc. (unless opts.NoMinorFast)
  - Special Shabbatot - Shabbat Shekalim, Zachor, etc. (unless opts.NoSpecialShabbat)
  - Modern Holidays - Yom HaShoah, Yom HaAtzma'ut, etc. (unless opts.NoModern)
  - Rosh Chodesh (unless opts.NoRoshChodesh)

Holiday and Torah reading schedules differ between Israel and the Disapora.
Set opts.IL=true to use the Israeli schedule.

Additional non-default event types can be specified:
  - Parashat HaShavua - weekly Torah Reading on Saturdays (opts.Sedrot)
  - Counting of the Omer (opts.Omer)
  - Babylonian Talmud Daf Yomi (opts.DafYomi)
  - Jerusalem Talmud (Yerushalmi) Yomi (opts.YerushalmiYomi)
  - Mishna Yomi (opts.MishnaYomi)
  - Shabbat Mevarchim HaChodesh on Saturday before Rosh Chodesh (opts.ShabbatMevarchim)
  - Molad announcement on Saturday before Rosh Chodesh (opts.Molad)
  - Yom Kippur Katan (opts.YomKippurKatan)

Candle-lighting and Havdalah times are approximated using latitude and longitude
specified by the HLocation class. The HLocation class contains a small
database of cities with their associated geographic information and time-zone information.
If you ever have any doubts about Hebcal's times, consult your local halachic authority.
If you enter geographic coordinates above the arctic circle or antarctic circle,
the times are guaranteed to be wrong.

To add candle-lighting options, set opts.CandleLighting=true and set
opts.Location to an instance of Location. By default, candle lighting
time is 18 minutes before sundown (40 minutes for Jerusalem) and Havdalah is
calculated according to Tzeit Hakochavim - Nightfall (the point when 3 small stars
are observable in the night time sky with the naked eye). The default Havdalah
option (Tzeit Hakochavim) is calculated when the sun is 8.5° below the horizon.
These defaults can be changed using these options:
  - opts.CandleLightingMins - minutes before sundown to light candles
  - opts.HavdalahMins - minutes after sundown for Havdalah (typical values are 42, 50, or 72).
    Havdalah times are supressed when opts.HavdalahMins=0.
  - opts.HavdalahDeg - degrees for solar depression for Havdalah.
    Default is 8.5 degrees for 3 small stars. Use 7.083 degress for 3 medium-sized stars.
    Havdalah times are supressed when opts.HavdalahDeg=0.

If both opts.CandleLighting=true and opts.Location is specified,
Chanukah candle-lighting times and minor fast start/end times will also be generated.
Chanukah candle-lighting is at dusk (when the sun is 6.0° below the horizon in the evening)
on weekdays, at regular candle-lighting time on Fridays, and at regular Havdalah time on
Saturday night (see above).

Minor fasts begin at Alot HaShachar (sun is 16.1° below the horizon in the morning) and
end when 3 medium-sized stars are observable in the night sky (sun is 7.083° below the horizon
in the evening).

Two options also exist for generating an Event with the Hebrew date:
  - opts.AddHebrewDates - print the Hebrew date for the entire date range
  - opts.AddHebrewDatesForEvents - print the Hebrew date for dates with some events
*/
func HebrewCalendar(opts *CalOptions) ([]CalEvent, error) {
	err := checkCandleOptions(opts)
	if err != nil {
		return nil, err
	}
	if opts.SunriseSunset && opts.Location == nil {
		return nil, errors.New("opts.SunriseSunset requires opts.Location")
	}
	if opts.DailyZmanim && opts.Location == nil {
		return nil, errors.New("opts.DailyZmanim requires opts.Location")
	}
	startAbs, endAbs, err := getStartAndEnd(opts)
	if err != nil {
		return nil, err
	}
	if opts.Location != nil && opts.Location.CountryCode == "IL" {
		opts.IL = true
	}
	opts.Mask = getMaskFromOptions(opts)
	if opts.YerushalmiYomi && opts.YerushalmiEdition == 0 {
		opts.YerushalmiEdition = yerushalmi.Vilna
	}
	beginYerushalmi := yerushalmi.VilnaStartRD
	if opts.YerushalmiEdition == yerushalmi.Schottenstein {
		beginYerushalmi = yerushalmi.SchottensteinStartRD
	}
	var (
		il           bool = opts.IL
		currentYear  int  = -1
		holidaysYear []HolidayEvent
		sedraYear    sedra.Sedra
		beginOmer    int
		endOmer      int
		myIdx        mishnayomi.MishnaYomiIndex
		userEvents   []HolidayEvent
	)
	firstWeekday := time.Weekday(startAbs % 7)
	events := make([]CalEvent, 0, 20)
	for abs := startAbs; abs <= endAbs; abs++ {
		hd := hdate.FromRD(abs)
		hyear := hd.Year()
		if hyear != currentYear {
			currentYear = hyear
			holidaysYear = GetHolidaysForYear(hyear, il)
			if opts.Sedrot || opts.DailySedra {
				sedraYear = sedra.New(hyear, il)
			}
			if opts.Omer {
				beginOmer = hdate.HebrewToRD(hyear, hdate.Nisan, 16)
				endOmer = hdate.HebrewToRD(hyear, hdate.Sivan, 5)
			}
			numUserEvents := len(opts.Yahrzeits) + len(opts.UserEvents)
			if numUserEvents != 0 {
				userEvents = make([]HolidayEvent, 0, numUserEvents)
				for _, yahrzeit := range opts.Yahrzeits {
					origDate := hdate.FromTime(yahrzeit.Date)
					observedDate, err := hdate.GetYahrzeit(currentYear, origDate)
					if err == nil {
						userEvents = append(userEvents, HolidayEvent{
							Date:  observedDate,
							Desc:  yahrzeit.Name,
							Flags: USER_EVENT,
						})
					}
				}
				for _, event := range opts.UserEvents {
					// Watch for ShortKislev and LongCheshvan
					if event.Day <= hdate.DaysInMonth(event.Month, hyear) {
						userEvents = append(userEvents, HolidayEvent{
							Date:  hdate.New(hyear, event.Month, event.Day),
							Desc:  event.Desc,
							Flags: USER_EVENT,
						})
					}
				}
			}
		}
		dow := hd.Weekday()
		prevEventsLength := len(events)
		if opts.SunriseSunset && (!opts.WeeklyAbbreviated || dow == firstWeekday) {
			events = append(events, riseSetEvent{date: hd, opts: opts})
		}
		if opts.DailySedra || (opts.Sedrot && dow == time.Saturday) {
			parsha := sedraYear.LookupByRD(abs)
			if !parsha.Chag {
				events = append(events, parshaEvent{Date: hd, Parsha: parsha, IL: il})
			}
		}
		var candlesEv TimedEvent
		for _, holidayEv := range holidaysYear {
			if hd == holidayEv.Date {
				events, candlesEv = appendHolidayAndRelated(events, candlesEv, holidayEv, opts)
			}
		}
		for _, userEv := range userEvents {
			if abs == userEv.Date.Abs() {
				events = append(events, userEv)
			}
		}
		if !opts.WeeklyAbbreviated || dow == firstWeekday {
			if opts.Omer && abs >= beginOmer && abs <= endOmer {
				omerDay := abs - beginOmer + 1
				events = append(events, NewOmerEvent(hd, omerDay))
			}
			if opts.DafYomi && hyear >= 5684 {
				daf, _ := dafyomi.New(hd)
				events = append(events, dafYomiEvent{Date: hd, Daf: daf})
			}
			if opts.YerushalmiYomi && abs >= beginYerushalmi {
				daf := yerushalmi.New(hd, opts.YerushalmiEdition)
				// daf.Blatt will be 0 to signal no Yerushalmi Yomi on YK and 9Av
				if daf.Blatt != 0 {
					events = append(events, yyomiEvent{Date: hd, Daf: daf})
				}
			}
			if opts.MishnaYomi && abs >= mishnayomi.MishnaYomiStart {
				if len(myIdx) == 0 {
					myIdx = mishnayomi.MakeIndex()
				}
				mishna, _ := myIdx.Lookup(hd)
				events = append(events, mishnaYomiEvent{Date: hd, Mishna: mishna})
			}
			if opts.DailyZmanim {
				zmanEvents := dailyZemanim(hd, opts)
				events = append(events, zmanEvents...)
			}
		}
		if (candlesEv == TimedEvent{}) && opts.CandleLighting && (dow == time.Friday || dow == time.Saturday) {
			candlesEv = makeCandleEvent(hd, opts, nil)
		}
		if (candlesEv != TimedEvent{}) {
			events = append(events, candlesEv)
		}
		if opts.Molad && dow == time.Saturday && hd.Month() != hdate.Elul && hd.Day() >= 23 && hd.Day() <= 29 {
			nextMonthName, nextMonth := nextMonthName(hd.Year(), hd.Month())
			molad := hdate.NewMolad(hd.Year(), nextMonth)
			events = append(events, moladEvent{Date: hd, Molad: molad, MonthName: nextMonthName})
		}
		if (opts.AddHebrewDates && (!opts.WeeklyAbbreviated || dow == firstWeekday)) ||
			((opts.AddHebrewDates || opts.AddHebrewDatesForEvents) && prevEventsLength != len(events)) {
			events = append(events, nil)
			copy(events[prevEventsLength+1:], events[prevEventsLength:])
			events[prevEventsLength] = hebrewDateEvent{Date: hd}
		}
	}
	return events, nil
}

func getStartAndEnd(opts *CalOptions) (int, int, error) {
	if (opts.Start != hdate.HDate{} && opts.End == hdate.HDate{}) ||
		(opts.Start == hdate.HDate{} && opts.End != hdate.HDate{}) {
		return 0, 0, errors.New("opts.Start requires opts.End")
	} else if (opts.Start != hdate.HDate{}) && (opts.End != hdate.HDate{}) {
		return opts.Start.Abs(), opts.End.Abs(), nil
	}
	year := opts.Year
	if year == 0 {
		t := time.Now()
		gy, gm, gd := t.Date()
		if opts.IsHebrewYear {
			today := hdate.FromGregorian(gy, gm, gd)
			year = today.Year()
		} else {
			year = gy
		}
	} else if opts.IsHebrewYear && year < 1 {
		return 0, 0, errors.New("invalid Hebrew year")
	}
	numYears := opts.NumYears
	if numYears == 0 {
		numYears = 1
	}
	if opts.IsHebrewYear {
		startDate := hdate.New(year, hdate.Tishrei, 1)
		// for full Hebrew year, start on Erev Rosh Hashana which
		// is technically in the previous Hebrew year
		// (but conveniently lets us get candle-lighting time for Erev)
		startAbs := startDate.Abs()
		if year > 1 {
			startAbs--
		}
		endDate := hdate.New(year+numYears, hdate.Tishrei, 1)
		endAbs := endDate.Abs() - 1
		return startAbs, endAbs, nil
	} else {
		// disable candle-lighting times for very early dates
		if year < 100 {
			opts.CandleLighting = false
		}
		month := time.January
		if opts.Month != 0 {
			month = opts.Month
		}
		startAbs := greg.ToRD(year, month, 1)
		if opts.Month != 0 {
			endAbs := startAbs + greg.DaysIn(opts.Month, year)
			return startAbs, endAbs - 1, nil
		}
		endAbs := greg.ToRD(year+numYears, time.January, 1)
		return startAbs, endAbs - 1, nil
	}
}

func intAbs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

var israelCityOffset = map[string]int{
	"Jerusalem":       40,
	"Haifa":           30,
	"Zichron Ya'akov": 30,
	"Zichron Ya‘akov": 30,
	"Zichron Ya’akov": 30,
	"Zichron Yaakov":  30,
	"Zikhron Ya'akov": 30,
	"Zikhron Ya'aqov": 30,
	"Zikhron Ya‘akov": 30,
	"Zikhron Ya‘aqov": 30,
	"Zikhron Ya’akov": 30,
	"Zikhron Ya’aqov": 30,
	"Zikhron Yaakov":  30,
}

func checkCandleOptions(opts *CalOptions) error {
	if !opts.CandleLighting {
		return nil
	}
	if opts.Location == nil {
		return errors.New("opts.CandleLighting requires opts.Location")
	}
	if opts.HavdalahMins != 0 && opts.HavdalahDeg != 0.0 {
		return errors.New("opts.HavdalahMins and opts.HavdalahDeg are mutually exclusive")
	}
	min := 18
	if opts.CandleLightingMins != 0 {
		min = opts.CandleLightingMins
	}
	loc := opts.Location
	if loc.CountryCode == "IL" {
		offset := israelCityOffset[loc.Name]
		if offset != 0 && min == 18 {
			min = offset
		}
	}
	opts.CandleLightingMins = -1 * intAbs(min)
	if opts.HavdalahMins != 0 {
		opts.HavdalahMins = intAbs(opts.HavdalahMins)
	} else if opts.HavdalahDeg != 0.0 {
		opts.HavdalahDeg = math.Abs(opts.HavdalahDeg)
	} else {
		opts.HavdalahDeg = zmanim.Tzeit3SmallStars
	}
	return nil
}

const maskLightCandles = LIGHT_CANDLES |
	LIGHT_CANDLES_TZEIS |
	CHANUKAH_CANDLES |
	YOM_TOV_ENDS

func getMaskFromOptions(opts *CalOptions) HolidayFlags {
	if opts.Mask != 0 {
		m := opts.Mask
		if (m & ROSH_CHODESH) != 0 {
			opts.NoRoshChodesh = false
		}
		if (m & MODERN_HOLIDAY) != 0 {
			opts.NoModern = false
		}
		if (m & MINOR_FAST) != 0 {
			opts.NoMinorFast = false
		}
		if (m & SPECIAL_SHABBAT) != 0 {
			opts.NoSpecialShabbat = false
		}
		if (m & PARSHA_HASHAVUA) != 0 {
			opts.Sedrot = true
		}
		if (m & DAF_YOMI) != 0 {
			opts.DafYomi = true
		}
		if (m & OMER_COUNT) != 0 {
			opts.Omer = true
		}
		if (m & SHABBAT_MEVARCHIM) != 0 {
			opts.ShabbatMevarchim = true
		}
		if (m & MISHNA_YOMI) != 0 {
			opts.MishnaYomi = true
		}
		if (m & YOM_KIPPUR_KATAN) != 0 {
			opts.YomKippurKatan = true
		}
		if (m & YERUSHALMI_YOMI) != 0 {
			opts.YerushalmiYomi = true
		}
		return m
	}
	var mask HolidayFlags
	// default opts
	if !opts.NoHolidays {
		mask |= ROSH_CHODESH | YOM_TOV_ENDS | MINOR_FAST |
			SPECIAL_SHABBAT | MODERN_HOLIDAY | MAJOR_FAST |
			MINOR_HOLIDAY | EREV | CHOL_HAMOED |
			LIGHT_CANDLES | LIGHT_CANDLES_TZEIS | CHANUKAH_CANDLES
	}
	if opts.CandleLighting {
		mask |= LIGHT_CANDLES | LIGHT_CANDLES_TZEIS | YOM_TOV_ENDS
	}
	// suppression of defaults
	if opts.NoRoshChodesh {
		mask &= ^ROSH_CHODESH
	}
	if opts.NoModern {
		mask &= ^MODERN_HOLIDAY
	}
	if opts.NoMinorFast {
		mask &= ^MINOR_FAST
	}
	if opts.NoSpecialShabbat {
		mask &= ^SPECIAL_SHABBAT
		mask &= ^SHABBAT_MEVARCHIM
	}
	if opts.IL {
		mask |= IL_ONLY
	} else {
		mask |= CHUL_ONLY
	}
	// non-default opts
	if opts.Sedrot {
		mask |= PARSHA_HASHAVUA
	}
	if opts.DafYomi {
		mask |= DAF_YOMI
	}
	if opts.MishnaYomi {
		mask |= MISHNA_YOMI
	}
	if opts.YerushalmiYomi {
		mask |= YERUSHALMI_YOMI
	}
	if opts.Omer {
		mask |= OMER_COUNT
	}
	if opts.ShabbatMevarchim {
		mask |= SHABBAT_MEVARCHIM
	}
	if opts.YomKippurKatan {
		mask |= YOM_KIPPUR_KATAN
	}
	return mask
}

func appendHolidayAndRelated(events []CalEvent, candlesEv TimedEvent, ev CalEvent, opts *CalOptions) ([]CalEvent, TimedEvent) {
	mask := ev.GetFlags()
	if (!opts.YomKippurKatan && (mask&YOM_KIPPUR_KATAN) != 0) ||
		(opts.NoModern && (mask&MODERN_HOLIDAY) != 0) {
		return events, candlesEv // bail out early
	}
	isMajorFast := (mask & MAJOR_FAST) != 0
	isMinorFast := (mask & MINOR_FAST) != 0
	var startEvent, endEvent TimedEvent
	if opts.CandleLighting && (isMajorFast || isMinorFast) && ev.Render("en") != "Yom Kippur" {
		startEvent, endEvent = makeFastStartEnd(ev, opts)
		if (startEvent != TimedEvent{}) && (isMajorFast || (isMinorFast && !opts.NoMinorFast)) {
			events = append(events, startEvent)
		}
	}
	if (mask & opts.Mask) != 0 {
		if opts.CandleLighting && (mask&maskLightCandles) != 0 {
			if (mask&CHANUKAH_CANDLES) != 0 && !opts.NoHolidays {
				// Replace Chanukah event with a clone that includes candle lighting time.
				// For clarity, allow a "duplicate" candle lighting event to remain for Shabbat
				tmp := makeChanukahCandleLighting(ev.(HolidayEvent), opts)
				if (tmp != TimedEvent{}) {
					ev = tmp
				}
			} else {
				hd := ev.GetDate()
				candlesEv = makeCandleEvent(hd, opts, ev)
			}
		}
		if opts.YomKippurKatan && (mask&YOM_KIPPUR_KATAN) != 0 {
			events = append(events, ev)
		} else if !opts.NoHolidays {
			events = append(events, ev)
		}
	}
	if (endEvent != TimedEvent{}) && (isMajorFast || (isMinorFast && !opts.NoMinorFast)) {
		events = append(events, endEvent)
	}
	return events, candlesEv
}
