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
package hebcal

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
	var (
		il          bool = opts.IL
		currentYear int  = -1
		// holidaysYear []HolidayEvent
		sedra     Sedra
		beginOmer int
		endOmer   int
		candlesEv HEvent
	)
	events := make([]HEvent, 0, 20)
	for abs := startAbs; abs <= endAbs; abs++ {
		hd := NewHDateFromRD(abs)
		hyear := hd.Year
		if hd.Year != currentYear {
			currentYear = hyear
			// holidaysYear = GetHolidaysForYear(hyear, il)
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
		/*
			const prevEventsLength = evts.length;
			const dow = hd.getDay();
			let candlesEv = undefined;
			const ev = holidaysYear.get(hd.toString()) || [];
			ev.forEach((e) => {
				candlesEv = appendHolidayAndRelated(evts, e, options, candlesEv, dow);
			});
		*/
		if opts.Sedrot && dow == time.Saturday && hyear >= 3762 {
			parsha := sedra.LookupByRD(abs)
			if !parsha.Chag {
				events = append(events, ParshaEvent{Date: hd, Parsha: parsha, IL: il})
			}
		}
		if opts.DafYomi && hyear >= 5684 {
			daf, _ := GetDafYomi(hd)
			events = append(events, DafYomiEvent{Date: hd, Daf: daf})
		}
		if opts.Omer && abs >= beginOmer && abs <= endOmer {
			omerDay := abs - beginOmer + 1
			events = append(events, NewOmerEvent(hd, omerDay))
		}
		/*
			const hmonth = hd.getMonth();
			if (options.molad && dow == SAT && hmonth != ELUL && hd.getDate() >= 23 && hd.getDate() <= 29) {
				const monNext = (hmonth == HDate.monthsInYear(hyear) ? NISAN : hmonth + 1);
				evts.push(new MoladEvent(hd, hyear, monNext));
			}
		*/
		if candlesEv == nil && opts.CandleLighting && (dow == time.Friday || dow == time.Saturday) {
			candlesEv = makeCandleEvent(hd, opts, nil)
			if dow == time.Friday && candlesEv != nil && opts.Sedrot && currentYear >= 3762 {
				// candlesEv.memo = sedra.getString(abs)
			}
		}
		if candlesEv != nil {
			events = append(events, candlesEv)
		}
		if opts.AddHebrewDates || (opts.AddHebrewDatesForEvents && prevEventsLength != len(events)) {
			// events = append(events, HEvent{})
		}
	}
	return events, nil
}

func getStartAndEnd(opts *CalOptions) (int, int, error) {
	if (opts.Start != nil && opts.End == nil) ||
		(opts.Start == nil && opts.End != nil) {
		return 0, 0, errors.New("opts.Start requires opts.End")
	} else if opts.Start != nil && opts.End != nil {
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
			var mlen [13]int
			if IsGregLeapYear(year) {
				mlen = mlenLeap
			} else {
				mlen = mlenStd
			}
			endAbs := startAbs + mlen[opts.Month]
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
		opts.HavdalahDeg = 8.5
	}
	return nil
}
