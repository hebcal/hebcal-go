package hebcal

import (
	"errors"
	"strconv"
)

/**
 * Calculates yahrzeit.
 * `hyear` must be after original `gdate` of death.
 * Returns `undefined` when requested year preceeds or is same as original year.
 *
 * Hebcal uses the algorithm defined in "Calendrical Calculations"
 * by Edward M. Reingold and Nachum Dershowitz.
 *
 * The customary anniversary date of a death is more complicated and depends
 * also on the character of the year in which the first anniversary occurs.
 * There are several cases:
 *
 * * If the date of death is Marcheshvan 30, the anniversary in general depends
 *   on the first anniversary; if that first anniversary was not Marcheshvan 30,
 *   use the day before Kislev 1.
 * * If the date of death is Kislev 30, the anniversary in general again depends
 *   on the first anniversary — if that was not Kislev 30, use the day before
 *   Tevet 1.
 * * If the date of death is Adar II, the anniversary is the same day in the
 *   last month of the Hebrew year (Adar or Adar II).
 * * If the date of death is Adar I 30, the anniversary in a Hebrew year that
 *   is not a leap year (in which Adar only has 29 days) is the last day in
 *   Shevat.
 * * In all other cases, use the normal (that is, same month number) anniversary
 *   of the date of death. [Calendrical Calculations p. 113]
 * @example
 * import {HebrewCalendar} from '@hebcal/core';
 * const dt = new Date(2014, 2, 2); // '2014-03-02' == '30 Adar I 5774'
 * const hd = HebrewCalendar.getYahrzeit(5780, dt); // '30 Sh\'vat 5780'
 * console.log(hd.greg().toLocaleDateString('en-US')); // '2/25/2020'
 * @param {number} hyear Hebrew year
 * @param {Date|HDate} gdate Gregorian or Hebrew date of death
 * @return {HDate} anniversary occurring in hyear
 */
func GetYahrzeit(hyear int, date HDate) (HDate, error) {
	if hyear <= date.Year {
		return HDate{}, errors.New("year " + strconv.Itoa(hyear) + " occurs on or before original date")
	}

	if date.Month == Cheshvan && date.Day == 30 && !LongCheshvan(date.Year+1) {
		// If it's Heshvan 30 it depends on the first anniversary;
		// if that was not Heshvan 30, use the day before Kislev 1.
		date = NewHDateFromRD(HebrewToRD(hyear, Kislev, 1) - 1)
	} else if date.Month == Kislev && date.Day == 30 && ShortKislev(date.Year+1) {
		// If it's Kislev 30 it depends on the first anniversary;
		// if that was not Kislev 30, use the day before Teveth 1.
		date = NewHDateFromRD(HebrewToRD(hyear, Tevet, 1) - 1)
	} else if date.Month == Adar2 {
		// If it's Adar II, use the same day in last month of year (Adar or Adar II).
		date.Month = HMonth(MonthsInHebYear(hyear))
	} else if date.Month == Adar1 && date.Day == 30 && !IsHebLeapYear(hyear) {
		// If it's the 30th in Adar I and year is not a leap year
		// (so Adar has only 29 days), use the last day in Shevat.
		date.Day = 30
		date.Month = Shvat
	}
	// In all other cases, use the normal anniversary of the date of death.

	// advance day to rosh chodesh if needed
	if date.Month == Cheshvan && date.Day == 30 && !LongCheshvan(hyear) {
		date.Month = Kislev
		date.Day = 1
	} else if date.Month == Kislev && date.Day == 30 && ShortKislev(hyear) {
		date.Month = Tevet
		date.Day = 1
	}

	return NewHDate(hyear, date.Month, date.Day), nil
}

/**
 * Calculates a birthday or anniversary (non-yahrzeit).
 * `hyear` must be after original `gdate` of anniversary.
 * Returns `undefined` when requested year preceeds or is same as original year.
 *
 * Hebcal uses the algorithm defined in "Calendrical Calculations"
 * by Edward M. Reingold and Nachum Dershowitz.
 *
 * The birthday of someone born in Adar of an ordinary year or Adar II of
 * a leap year is also always in the last month of the year, be that Adar
 * or Adar II. The birthday in an ordinary year of someone born during the
 * first 29 days of Adar I in a leap year is on the corresponding day of Adar;
 * in a leap year, the birthday occurs in Adar I, as expected.
 *
 * Someone born on the thirtieth day of Marcheshvan, Kislev, or Adar I
 * has his birthday postponed until the first of the following month in
 * years where that day does not occur. [Calendrical Calculations p. 111]
 * @example
 * import {HebrewCalendar} from '@hebcal/core';
 * const dt = new Date(2014, 2, 2); // '2014-03-02' == '30 Adar I 5774'
 * const hd = HebrewCalendar.getBirthdayOrAnniversary(5780, dt); // '1 Nisan 5780'
 * console.log(hd.greg().toLocaleDateString('en-US')); // '3/26/2020'
 * @param {number} hyear Hebrew year
 * @param {Date|HDate} gdate Gregorian or Hebrew date of event
 * @return {HDate} anniversary occurring in `hyear`
 */
func GetBirthdayOrAnniversary(hyear int, date HDate) (HDate, error) {
	if hyear <= date.Year {
		return HDate{}, errors.New("year " + strconv.Itoa(hyear) + " occurs on or before original date")
	}
	isOrigLeap := IsHebLeapYear(date.Year)
	month := date.Month
	day := date.Day
	if (month == Adar1 && !isOrigLeap) || (month == Adar2 && isOrigLeap) {
		month = HMonth(MonthsInHebYear(hyear))
	} else if month == Cheshvan && day == 30 && !LongCheshvan(hyear) {
		month = Kislev
		day = 1
	} else if month == Kislev && day == 30 && ShortKislev(hyear) {
		month = Tevet
		day = 1
	} else if month == Adar1 && day == 30 && isOrigLeap && !IsHebLeapYear(hyear) {
		month = Nisan
		day = 1
	}
	return NewHDate(hyear, month, day), nil
}
