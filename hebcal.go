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
	"strings"
	"time"
)

/*
Calculates holidays and other Hebrew calendar events based on CalOptions.

Each holiday is represented by an HEvent object which includes a date,
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
  - Daf Yomi (opts.DafYomi)
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

Lastly, translation and transliteration of event titles is controlled by
opts.Locale and the Locale API.
This package supports three locales by default:
  - en - default, Sephardic transliterations (e.g. "Shabbat")
  - ashkenazi - Ashkenazi transliterations (e.g. "Shabbos")
  - he - Hebrew (e.g. "שַׁבָּת")

Additional locales (such as ru or fr) are supported by the
https://github.com/hebcal/hebcal-locales hebcal-locales-go package
*/
func HebrewCalendar(opts *CalOptions) ([]HEvent, error) {
	err := checkCandleOptions(opts)
	if err != nil {
		return nil, err
	}
	startAbs, endAbs, err := getStartAndEnd(opts)
	if err != nil {
		return nil, err
	}
	if opts.Location != nil && opts.Location.CountryCode == "IL" {
		opts.IL = true
	}
	opts.Mask = getMaskFromOptions(opts)
	var (
		il           bool = opts.IL
		currentYear  int  = -1
		holidaysYear []HolidayEvent
		sedra        Sedra
		beginOmer    int
		endOmer      int
	)
	events := make([]HEvent, 0, 20)
	for abs := startAbs; abs <= endAbs; abs++ {
		hd := NewHDateFromRD(abs)
		hyear := hd.Year
		if hd.Year != currentYear {
			currentYear = hyear
			holidaysYear = GetHolidaysForYear(hyear, il)
			if opts.Sedrot && currentYear >= 3762 {
				sedra = NewSedra(hyear, il)
			}
			if opts.Omer {
				beginOmer = HebrewToRD(hyear, Nisan, 16)
				endOmer = HebrewToRD(hyear, Sivan, 5)
			}
		}
		dow := hd.Weekday()
		prevEventsLength := len(events)
		var candlesEv TimedEvent
		for _, holidayEv := range holidaysYear {
			if hd == holidayEv.Date {
				events, candlesEv = appendHolidayAndRelated(events, candlesEv, holidayEv, opts)
			}
		}
		if opts.Sedrot && dow == time.Saturday && hyear >= 3762 {
			parsha := sedra.LookupByRD(abs)
			if !parsha.Chag {
				events = append(events, parshaEvent{Date: hd, Parsha: parsha, IL: il})
			}
		}
		if opts.DafYomi && hyear >= 5684 {
			daf, _ := GetDafYomi(hd)
			events = append(events, dafYomiEvent{Date: hd, Daf: daf})
		}
		if opts.Omer && abs >= beginOmer && abs <= endOmer {
			omerDay := abs - beginOmer + 1
			events = append(events, newOmerEvent(hd, omerDay))
		}
		/*
			const hmonth = hd.getMonth();
			if (options.molad && dow == SAT && hmonth != ELUL && hd.getDate() >= 23 && hd.getDate() <= 29) {
				const monNext = (hmonth == HDate.monthsInYear(hyear) ? NISAN : hmonth + 1);
				evts.push(new MoladEvent(hd, hyear, monNext));
			}
		*/
		if (candlesEv == TimedEvent{}) && opts.CandleLighting && (dow == time.Friday || dow == time.Saturday) {
			candlesEv = makeCandleEvent(hd, opts, nil)
		}
		if (candlesEv != TimedEvent{}) {
			events = append(events, candlesEv)
		}
		if opts.AddHebrewDates || (opts.AddHebrewDatesForEvents && prevEventsLength != len(events)) {
			events = append(events, hebrewDateEvent{Date: hd})
		}
	}
	return events, nil
}

func getStartAndEnd(opts *CalOptions) (int, int, error) {
	if (opts.Start != HDate{} && opts.End == HDate{}) ||
		(opts.Start == HDate{} && opts.End != HDate{}) {
		return 0, 0, errors.New("opts.Start requires opts.End")
	} else if (opts.Start != HDate{}) && (opts.End != HDate{}) {
		return opts.Start.Abs(), opts.End.Abs(), nil
	}
	year := opts.Year
	if year == 0 {
		t := time.Now()
		gy, gm, gd := t.Date()
		if opts.IsHebrewYear {
			today := NewHDateFromGregorian(gy, gm, gd)
			year = today.Year
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
		startDate := NewHDate(year, Tishrei, 1)
		// for full Hebrew year, start on Erev Rosh Hashana which
		// is technically in the previous Hebrew year
		// (but conveniently lets us get candle-lighting time for Erev)
		startAbs := startDate.Abs() - 1
		endDate := NewHDate(year+numYears, Tishrei, 1)
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
		startAbs, _ := GregorianToRD(year, month, 1)
		if opts.Month != 0 {
			endAbs := startAbs + DaysIn(opts.Month, year)
			return startAbs, endAbs - 1, nil
		}
		endAbs, _ := GregorianToRD(year+numYears, time.January, 1)
		return startAbs, endAbs - 1, nil
	}
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
	if loc.CountryCode == "IL" && strings.HasPrefix(loc.Name, "Jerusalem") && min == 18 {
		min = 40
	}
	opts.CandleLightingMins = -1 * intAbs(min)
	if opts.HavdalahMins != 0 {
		opts.HavdalahMins = intAbs(opts.HavdalahMins)
	} else if opts.HavdalahDeg != 0.0 {
		opts.HavdalahDeg = math.Abs(opts.HavdalahDeg)
	} else {
		opts.HavdalahDeg = Tzeit3SmallStars
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

func appendHolidayAndRelated(events []HEvent, candlesEv TimedEvent, ev HEvent, opts *CalOptions) ([]HEvent, TimedEvent) {
	mask := ev.GetFlags()
	if !opts.YomKippurKatan && (mask&YOM_KIPPUR_KATAN) != 0 {
		return events, candlesEv // bail out early
	}
	loc := opts.Location
	isMajorFast := (mask & MAJOR_FAST) != 0
	isMinorFast := (mask & MINOR_FAST) != 0
	var startEvent, endEvent TimedEvent
	if opts.CandleLighting && (isMajorFast || isMinorFast) && ev.Render() != "Yom Kippur" {
		startEvent, endEvent = makeFastStartEnd(ev, loc)
		if (startEvent != TimedEvent{}) && (isMajorFast || (isMinorFast && !opts.NoMinorFast)) {
			events = append(events, startEvent)
		}
	}
	if (mask & opts.Mask) != 0 {
		if opts.CandleLighting && (mask&maskLightCandles) != 0 {
			if (mask&CHANUKAH_CANDLES) != 0 && !opts.NoHolidays {
				// Replace Chanukah event with a clone that includes candle lighting time.
				// For clarity, allow a "duplicate" candle lighting event to remain for Shabbat
				ev = makeChanukahCandleLighting(ev.(HolidayEvent), opts)
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

type hebrewDateEvent struct {
	Date HDate
}

func (ev hebrewDateEvent) GetDate() HDate {
	return ev.Date
}

func (ev hebrewDateEvent) Render() string {
	return ev.Date.String()
}

func (ev hebrewDateEvent) GetFlags() HolidayFlags {
	return HEBREW_DATE
}

func (ev hebrewDateEvent) GetEmoji() string {
	return ""
}
