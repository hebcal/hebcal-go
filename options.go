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
package hebcal

import "time"

// CalOptions are used by HebrewCalendar() to generate a slice of events
type CalOptions struct {
	/* latitude/longitude/tzid used for candle-lighting */
	Location *HLocation
	/* Gregorian or Hebrew year */
	Year int
	/* to interpret year as Hebrew year */
	IsHebrewYear bool
	/* Gregorian month (to filter results to a single month) */
	Month time.Month
	/* generate calendar for multiple years (default 1) */
	NumYears int
	/* use specific start date (requires end date) */
	Start *HDate
	/* use specific end date (requires start date) */
	End *HDate
	/* calculate candle-lighting and havdalah times */
	CandleLighting bool
	/* minutes before sundown to light candles (default 18) */
	CandleLightingMins int
	/*
	 * minutes after sundown for Havdalah (typical values are 42, 50, or 72).
	 * If `undefined` (the default), calculate Havdalah according to Tzeit Hakochavim -
	 * Nightfall (the point when 3 small stars are observable in the night time sky with
	 * the naked eye). If `0`, Havdalah times are supressed.
	 */
	HavdalahMins int
	/*
	 * degrees for solar depression for Havdalah.
	 * Default is 8.5 degrees for 3 small stars.
	 * Use 7.083 degress for 3 medium-sized stars.
	 * Havdalah times are supressed when `havdalahDeg=0`.
	 */
	HavdalahDeg float64
	/* calculate parashah hashavua on Saturdays */
	Sedrot bool
	/* Israeli holiday and sedra schedule */
	IL bool
	/* suppress minor fasts */
	NoMinorFast bool
	/* suppress modern holidays */
	NoModern bool
	/* suppress Rosh Chodesh & Shabbat Mevarchim */
	NoRoshChodesh    bool
	ShabbatMevarchim bool
	/* suppress Special Shabbat */
	NoSpecialShabbat bool
	/* suppress regular holidays */
	NoHolidays bool
	/* include Daf Yomi */
	DafYomi bool
	/* include Mishna Yomi */
	MishnaYomi bool
	/* include Days of the Omer */
	Omer bool
	/* include event announcing the molad */
	Molad bool
	/*
	 * translate event titles according to a locale
	 * Default value is `en`, also built-in are `he` and `ashkenazi`.
	 * Additional locales (such as `ru` or `fr`) are provided by the
	 * {@link https://github.com/hebcal/hebcal-locales @hebcal/locales} package
	 */
	Locale string
	/* print the Hebrew date for the entire date range */
	AddHebrewDates bool
	/* print the Hebrew date for dates with some events */
	AddHebrewDatesForEvents bool
	/* use bitmask from `flags` to filter events */
	Mask int
	/*
	 * include Yom Kippur Katan (default `false`).
	 * יוֹם כִּפּוּר קָטָן is a minor day of atonement occurring monthly on the day preceeding each Rosh Chodesh.
	 * Yom Kippur Katan is omitted in Elul (on the day before Rosh Hashanah),
	 * Tishrei (Yom Kippur has just passed), Kislev (due to Chanukah)
	 * and Nisan (fasting not permitted during Nisan).
	 * When Rosh Chodesh occurs on Shabbat or Sunday, Yom Kippur Katan is observed on the preceding Thursday.
	 * @see {@link https://en.wikipedia.org/wiki/Yom_Kippur_Katan#Practices Wikipedia Yom Kippur Katan practices}
	 */
	YomKippurKatan bool
	/*
	 * Whether to use 12-hour time (as opposed to 24-hour time).
	 * Possible values are `true` and `false` the default is locale dependent.
	 */
	Hour12 bool
}
