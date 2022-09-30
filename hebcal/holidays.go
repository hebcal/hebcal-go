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
	"sort"
	"strconv"
	"time"

	"github.com/hebcal/hebcal-go/hdate"
	"github.com/hebcal/hebcal-go/sedra"
)

type holiday struct {
	mm     hdate.HMonth
	dd     int
	desc   string
	flags  HolidayFlags
	emoji  string
	chmDay int
}

const chanukahEmoji = "🕎"

var staticHolidays = []holiday{
	{mm: hdate.Tishrei, dd: 2, desc: "Rosh Hashana II", flags: CHAG | YOM_TOV_ENDS, emoji: "🍏🍯"},
	{mm: hdate.Tishrei, dd: 9, desc: "Erev Yom Kippur", flags: EREV | LIGHT_CANDLES},
	{mm: hdate.Tishrei, dd: 10, desc: "Yom Kippur", flags: CHAG | MAJOR_FAST | YOM_TOV_ENDS},
	{mm: hdate.Tishrei, dd: 14, desc: "Erev Sukkot", flags: EREV | LIGHT_CANDLES},

	{mm: hdate.Tishrei, dd: 15, desc: "Sukkot I", flags: CHUL_ONLY | CHAG | LIGHT_CANDLES_TZEIS},
	{mm: hdate.Tishrei, dd: 16, desc: "Sukkot II", flags: CHUL_ONLY | CHAG | YOM_TOV_ENDS},
	{mm: hdate.Tishrei, dd: 17, desc: "Sukkot III (CH''M)", flags: CHUL_ONLY | CHOL_HAMOED, chmDay: 1},
	{mm: hdate.Tishrei, dd: 18, desc: "Sukkot IV (CH''M)", flags: CHUL_ONLY | CHOL_HAMOED, chmDay: 2},
	{mm: hdate.Tishrei, dd: 19, desc: "Sukkot V (CH''M)", flags: CHUL_ONLY | CHOL_HAMOED, chmDay: 3},
	{mm: hdate.Tishrei, dd: 20, desc: "Sukkot VI (CH''M)", flags: CHUL_ONLY | CHOL_HAMOED, chmDay: 4},
	{mm: hdate.Tishrei, dd: 22, desc: "Shmini Atzeret",
		flags: CHUL_ONLY | CHAG | LIGHT_CANDLES_TZEIS},
	{mm: hdate.Tishrei, dd: 23, desc: "Simchat Torah",
		flags: CHUL_ONLY | CHAG | YOM_TOV_ENDS},

	{mm: hdate.Tishrei, dd: 15, desc: "Sukkot I", flags: IL_ONLY | CHAG | YOM_TOV_ENDS},
	{mm: hdate.Tishrei, dd: 16, desc: "Sukkot II (CH''M)", flags: IL_ONLY | CHOL_HAMOED, chmDay: 1},
	{mm: hdate.Tishrei, dd: 17, desc: "Sukkot III (CH''M)", flags: IL_ONLY | CHOL_HAMOED, chmDay: 2},
	{mm: hdate.Tishrei, dd: 18, desc: "Sukkot IV (CH''M)", flags: IL_ONLY | CHOL_HAMOED, chmDay: 3},
	{mm: hdate.Tishrei, dd: 19, desc: "Sukkot V (CH''M)", flags: IL_ONLY | CHOL_HAMOED, chmDay: 4},
	{mm: hdate.Tishrei, dd: 20, desc: "Sukkot VI (CH''M)", flags: IL_ONLY | CHOL_HAMOED, chmDay: 5},
	{mm: hdate.Tishrei, dd: 22, desc: "Shmini Atzeret",
		flags: IL_ONLY | CHAG | YOM_TOV_ENDS},

	{mm: hdate.Tishrei, dd: 21, desc: "Sukkot VII (Hoshana Raba)",
		flags: LIGHT_CANDLES | CHOL_HAMOED, chmDay: -1},
	{mm: hdate.Kislev, dd: 24, desc: "Chanukah: 1 Candle",
		flags: EREV | MINOR_HOLIDAY | CHANUKAH_CANDLES, emoji: chanukahEmoji},
	{mm: hdate.Tevet, dd: 10, desc: "Asara B'Tevet", flags: MINOR_FAST},
	{mm: hdate.Shvat, dd: 15, desc: "Tu BiShvat", flags: MINOR_HOLIDAY, emoji: "🌳"},
	{mm: hdate.Adar2, dd: 13, desc: "Erev Purim", flags: EREV | MINOR_HOLIDAY, emoji: "🎭️📜"},
	{mm: hdate.Adar2, dd: 14, desc: "Purim", flags: MINOR_HOLIDAY, emoji: "🎭️📜"},
	{mm: hdate.Nisan, dd: 14, desc: "Erev Pesach", flags: EREV | LIGHT_CANDLES},
	// Pesach Israel
	{mm: hdate.Nisan, dd: 15, desc: "Pesach I",
		flags: IL_ONLY | CHAG | YOM_TOV_ENDS},
	{mm: hdate.Nisan, dd: 16, desc: "Pesach II (CH''M)",
		flags: IL_ONLY | CHOL_HAMOED, chmDay: 1},
	{mm: hdate.Nisan, dd: 17, desc: "Pesach III (CH''M)",
		flags: IL_ONLY | CHOL_HAMOED, chmDay: 2},
	{mm: hdate.Nisan, dd: 18, desc: "Pesach IV (CH''M)",
		flags: IL_ONLY | CHOL_HAMOED, chmDay: 3},
	{mm: hdate.Nisan, dd: 19, desc: "Pesach V (CH''M)",
		flags: IL_ONLY | CHOL_HAMOED, chmDay: 4},
	{mm: hdate.Nisan, dd: 20, desc: "Pesach VI (CH''M)",
		flags: IL_ONLY | CHOL_HAMOED | LIGHT_CANDLES, chmDay: 5},
	{mm: hdate.Nisan, dd: 21, desc: "Pesach VII",
		flags: IL_ONLY | CHAG | YOM_TOV_ENDS},
	// Pesach chutz l'aretz
	{mm: hdate.Nisan, dd: 15, desc: "Pesach I",
		flags: CHUL_ONLY | CHAG | LIGHT_CANDLES_TZEIS},
	{mm: hdate.Nisan, dd: 16, desc: "Pesach II",
		flags: CHUL_ONLY | CHAG | YOM_TOV_ENDS},
	{mm: hdate.Nisan, dd: 17, desc: "Pesach III (CH''M)",
		flags: CHUL_ONLY | CHOL_HAMOED, chmDay: 1},
	{mm: hdate.Nisan, dd: 18, desc: "Pesach IV (CH''M)",
		flags: CHUL_ONLY | CHOL_HAMOED, chmDay: 2},
	{mm: hdate.Nisan, dd: 19, desc: "Pesach V (CH''M)",
		flags: CHUL_ONLY | CHOL_HAMOED, chmDay: 3},
	{mm: hdate.Nisan, dd: 20, desc: "Pesach VI (CH''M)",
		flags: CHUL_ONLY | CHOL_HAMOED | LIGHT_CANDLES, chmDay: 4},
	{mm: hdate.Nisan, dd: 21, desc: "Pesach VII",
		flags: CHUL_ONLY | CHAG | LIGHT_CANDLES_TZEIS},
	{mm: hdate.Nisan, dd: 22, desc: "Pesach VIII",
		flags: CHUL_ONLY | CHAG | YOM_TOV_ENDS},

	{mm: hdate.Iyyar, dd: 14, desc: "Pesach Sheni", flags: MINOR_HOLIDAY},
	{mm: hdate.Iyyar, dd: 18, desc: "Lag BaOmer", flags: MINOR_HOLIDAY, emoji: "🔥"},
	{mm: hdate.Sivan, dd: 5, desc: "Erev Shavuot",
		flags: EREV | LIGHT_CANDLES, emoji: "⛰️🌸"},
	{mm: hdate.Sivan, dd: 6, desc: "Shavuot",
		flags: IL_ONLY | CHAG | YOM_TOV_ENDS, emoji: "⛰️🌸"},
	{mm: hdate.Sivan, dd: 6, desc: "Shavuot I",
		flags: CHUL_ONLY | CHAG | LIGHT_CANDLES_TZEIS, emoji: "⛰️🌸"},
	{mm: hdate.Sivan, dd: 7, desc: "Shavuot II",
		flags: CHUL_ONLY | CHAG | YOM_TOV_ENDS, emoji: "⛰️🌸"},
	{mm: hdate.Av, dd: 15, desc: "Tu B'Av",
		flags: MINOR_HOLIDAY, emoji: "❤️"},
	{mm: hdate.Elul, dd: 1, desc: "Rosh Hashana LaBehemot",
		flags: MINOR_HOLIDAY, emoji: "🐑"},
	{mm: hdate.Elul, dd: 29, desc: "Erev Rosh Hashana",
		flags: EREV | LIGHT_CANDLES, emoji: "🍏🍯"},
}

var staticModernHolidays = []struct {
	firstYear        int // First observed in Hebrew year
	mm               hdate.HMonth
	dd               int
	desc             string
	chul             bool // also display in Chutz L'Aretz
	suppressEmoji    bool
	satPostponeToSun bool
	friPostponeToSun bool
}{
	{firstYear: 5727, mm: hdate.Iyyar, dd: 28, desc: "Yom Yerushalayim",
		chul: true},
	{firstYear: 5737, mm: hdate.Kislev, dd: 6, desc: "Ben-Gurion Day",
		satPostponeToSun: true, friPostponeToSun: true},
	{firstYear: 5750, mm: hdate.Shvat, dd: 30, desc: "Family Day",
		suppressEmoji: true},
	{firstYear: 5758, mm: hdate.Cheshvan, dd: 12, desc: "Yitzhak Rabin Memorial Day"},
	{firstYear: 5764, mm: hdate.Iyyar, dd: 10, desc: "Herzl Day",
		satPostponeToSun: true},
	{firstYear: 5765, mm: hdate.Tamuz, dd: 29, desc: "Jabotinsky Day",
		satPostponeToSun: true},
	{firstYear: 5769, mm: hdate.Cheshvan, dd: 29, desc: "Sigd",
		chul: true, suppressEmoji: true},
	{firstYear: 5777, mm: hdate.Nisan, dd: 10, desc: "Yom HaAliyah",
		chul: true},
	{firstYear: 5777, mm: hdate.Cheshvan, dd: 7, desc: "Yom HaAliyah School Observance"},
}

func tzomGedaliahDate(rh hdate.HDate) hdate.HDate {
	offset := 0
	if rh.Weekday() == time.Thursday {
		offset = 1
	}
	return hdate.New(rh.Year, hdate.Tishrei, 3+offset)
}

func taanitEstherDate(pesach hdate.HDate) hdate.HDate {
	offset := 31
	if pesach.Weekday() == time.Tuesday {
		offset = 33
	}
	return hdate.FromRD(pesach.Abs() - offset)
}

func shushanPurimDate(pesach hdate.HDate) hdate.HDate {
	offset := 29
	if pesach.Weekday() == time.Sunday {
		offset = 28
	}
	return hdate.FromRD(pesach.Abs() - offset)
}

func taanitBechorotDate(pesach hdate.HDate) hdate.HDate {
	if pesach.Prev().Weekday() == time.Saturday {
		return hdate.FromRD(hdate.DayOnOrBefore(time.Thursday, pesach.Abs()))
	} else {
		return hdate.New(pesach.Year, hdate.Nisan, 14)
	}
}

type byDate []HolidayEvent

func (s byDate) Len() int {
	return len(s)
}
func (s byDate) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byDate) Less(i, j int) bool {
	abs1 := s[i].Date.Abs()
	abs2 := s[j].Date.Abs()
	if abs1 == abs2 {
		return s[i].Desc < s[j].Desc
	}
	return abs1 < abs2
}

func nextMonthName(year int, month hdate.HMonth) (string, hdate.HMonth) {
	nextMonth := month + 1
	monthsInYear := hdate.MonthsInYear(year)
	if month == hdate.HMonth(monthsInYear) {
		nextMonth = hdate.Nisan
	}
	nextMonthName := hdate.New(year, nextMonth, 1).MonthName("en")
	return nextMonthName, nextMonth
}

func getAllHolidaysForYear(year int) []HolidayEvent {
	events := make([]HolidayEvent, 0, 120)
	// standard holidays that don't shift based on year
	for _, h := range staticHolidays {
		events = append(events, HolidayEvent{Date: hdate.New(year, h.mm, h.dd),
			Desc: h.desc, Flags: h.flags, Emoji: h.emoji, CholHaMoedDay: h.chmDay})
	}
	// variable holidays
	roshHashana := hdate.New(year, hdate.Tishrei, 1)
	nextRh := hdate.New(year+1, hdate.Tishrei, 1)
	pesach := hdate.New(year, hdate.Nisan, 15)
	pesachAbs := pesach.Abs()
	events = append(events,
		HolidayEvent{
			Date:  roshHashana,
			Desc:  "Rosh Hashana " + strconv.Itoa(year),
			Flags: CHAG | LIGHT_CANDLES_TZEIS,
			Emoji: "🍏🍯"},
		HolidayEvent{
			Date:  hdate.FromRD(hdate.DayOnOrBefore(time.Saturday, 7+roshHashana.Abs())),
			Desc:  "Shabbat Shuva",
			Flags: SPECIAL_SHABBAT,
			Emoji: "🕍"},
		HolidayEvent{
			Date:  tzomGedaliahDate(roshHashana),
			Desc:  "Tzom Gedaliah",
			Flags: MINOR_FAST},
		HolidayEvent{
			Date:  hdate.FromRD(hdate.DayOnOrBefore(time.Saturday, pesachAbs-43)),
			Desc:  "Shabbat Shekalim",
			Flags: SPECIAL_SHABBAT,
			Emoji: "🕍"},
		HolidayEvent{
			Date:  hdate.FromRD(hdate.DayOnOrBefore(time.Saturday, pesachAbs-30)),
			Desc:  "Shabbat Zachor",
			Flags: SPECIAL_SHABBAT,
			Emoji: "🕍"},
		HolidayEvent{
			Date:  taanitEstherDate(pesach),
			Desc:  "Ta'anit Esther",
			Flags: MINOR_FAST},
		HolidayEvent{
			Date:  shushanPurimDate(pesach),
			Desc:  "Shushan Purim",
			Flags: MINOR_HOLIDAY, Emoji: "🎭️📜"},
		HolidayEvent{
			Date:  hdate.FromRD(hdate.DayOnOrBefore(time.Saturday, pesachAbs-14) - 7),
			Desc:  "Shabbat Parah",
			Flags: SPECIAL_SHABBAT,
			Emoji: "🕍"},
		HolidayEvent{
			Date: hdate.FromRD(hdate.DayOnOrBefore(time.Saturday, pesachAbs-14)),
			Desc: "Shabbat HaChodesh", Flags: SPECIAL_SHABBAT,
			Emoji: "🕍"},
		HolidayEvent{
			Date: hdate.FromRD(hdate.DayOnOrBefore(time.Saturday, pesachAbs-1)),
			Desc: "Shabbat HaGadol", Flags: SPECIAL_SHABBAT,
			Emoji: "🕍"},
		HolidayEvent{
			Date:  taanitBechorotDate(pesach),
			Desc:  "Ta'anit Bechorot",
			Flags: MINOR_FAST},
		HolidayEvent{
			Date:  hdate.FromRD(hdate.DayOnOrBefore(time.Saturday, nextRh.Abs()-4)),
			Desc:  "Leil Selichot",
			Flags: MINOR_HOLIDAY,
			Emoji: "🕍"},
	)
	if hdate.IsLeapYear(year) {
		events = append(events, HolidayEvent{
			Date:  hdate.New(year, hdate.Adar1, 14),
			Desc:  "Purim Katan",
			Flags: MINOR_HOLIDAY,
			Emoji: "🎭️"})
		events = append(events, HolidayEvent{
			Date:  hdate.New(year, hdate.Adar1, 15),
			Desc:  "Shushan Purim Katan",
			Flags: MINOR_HOLIDAY,
			Emoji: "🎭️"})
	}
	for candles := 2; candles <= 6; candles++ {
		events = append(events, HolidayEvent{
			Date:        hdate.New(year, hdate.Kislev, 23+candles),
			Desc:        "Chanukah: " + strconv.Itoa(candles) + " Candles",
			Flags:       MINOR_HOLIDAY | CHANUKAH_CANDLES,
			Emoji:       chanukahEmoji,
			ChanukahDay: candles - 1})
	}
	var chanukah7 hdate.HDate
	if hdate.ShortKislev(year) {
		chanukah7 = hdate.New(year, hdate.Tevet, 1)
	} else {
		chanukah7 = hdate.New(year, hdate.Kislev, 30)
	}
	chanukah8 := chanukah7.Next()
	events = append(events,
		HolidayEvent{
			Date:        chanukah7,
			Desc:        "Chanukah: 7 Candles",
			Flags:       MINOR_HOLIDAY | CHANUKAH_CANDLES,
			Emoji:       chanukahEmoji,
			ChanukahDay: 6},
		HolidayEvent{
			Date:        chanukah8,
			Desc:        "Chanukah: 8 Candles",
			Flags:       MINOR_HOLIDAY | CHANUKAH_CANDLES,
			Emoji:       chanukahEmoji,
			ChanukahDay: 7},
		HolidayEvent{
			Date:        chanukah8.Next(),
			Desc:        "Chanukah: 8th Day",
			Flags:       MINOR_HOLIDAY,
			Emoji:       chanukahEmoji,
			ChanukahDay: 8})

	// Tisha BAv and the 3 weeks
	var tamuz17 = hdate.New(year, hdate.Tamuz, 17)
	if tamuz17.Weekday() == time.Saturday {
		tamuz17 = tamuz17.Next()
	}
	events = append(events,
		HolidayEvent{Date: tamuz17, Desc: "Tzom Tammuz", Flags: MINOR_FAST})

	var av9dt = hdate.New(year, hdate.Av, 9)
	var av9title = "Tish'a B'Av"
	if av9dt.Weekday() == time.Saturday {
		av9dt = av9dt.Next()
		av9title += " (observed)"
	}
	var av9abs = av9dt.Abs()
	events = append(events,
		HolidayEvent{
			Date:  hdate.FromRD(hdate.DayOnOrBefore(time.Saturday, av9abs)),
			Desc:  "Shabbat Chazon",
			Flags: SPECIAL_SHABBAT,
			Emoji: "🕍"},
		HolidayEvent{Date: av9dt.Prev(), Desc: "Erev Tish'a B'Av", Flags: EREV | MAJOR_FAST},
		HolidayEvent{Date: av9dt, Desc: av9title, Flags: MAJOR_FAST},
		HolidayEvent{
			Date:  hdate.FromRD(hdate.DayOnOrBefore(time.Saturday, av9abs+7)),
			Desc:  "Shabbat Nachamu",
			Flags: SPECIAL_SHABBAT,
			Emoji: "🕍"})

	// modern holidays
	if year >= 5708 {
		// Yom HaAtzma'ut only celebrated after 1948
		var day int
		if pesach.Weekday() == time.Sunday {
			day = 2
		} else if pesach.Weekday() == time.Saturday {
			day = 3
		} else if year < 5764 {
			day = 4
		} else if pesach.Weekday() == time.Tuesday {
			day = 5
		} else {
			day = 4
		}
		var tmpDate = hdate.New(year, hdate.Iyyar, day)
		events = append(events,
			HolidayEvent{Date: tmpDate, Desc: "Yom HaZikaron", Flags: MODERN_HOLIDAY, Emoji: "🇮🇱"},
			HolidayEvent{Date: tmpDate.Next(), Desc: "Yom HaAtzma'ut", Flags: MODERN_HOLIDAY, Emoji: "🇮🇱"},
		)
	}

	if year >= 5711 {
		// Yom HaShoah first observed in 1951
		var nisan27dt = hdate.New(year, hdate.Nisan, 27)
		/* When the actual date of Yom Hashoah falls on a Friday, the
		 * state of Israel observes Yom Hashoah on the preceding
		 * Thursday. When it falls on a Sunday, Yom Hashoah is observed
		 * on the following Monday.
		 * https://www.ushmm.org/remember/days-of-remembrance/resources/calendar
		 */
		if nisan27dt.Weekday() == time.Friday {
			nisan27dt = nisan27dt.Prev()
		} else if nisan27dt.Weekday() == time.Sunday {
			nisan27dt = nisan27dt.Next()
		}
		events = append(events,
			HolidayEvent{Date: nisan27dt, Desc: "Yom HaShoah", Flags: MODERN_HOLIDAY})
	}

	for _, h := range staticModernHolidays {
		if year >= h.firstYear {
			emoji := "🇮🇱"
			if h.suppressEmoji {
				emoji = ""
			}
			hd := hdate.New(year, h.mm, h.dd)
			if h.friPostponeToSun && hd.Weekday() == time.Friday {
				hd = hd.Next().Next()
			}
			if h.satPostponeToSun && hd.Weekday() == time.Saturday {
				hd = hd.Next()
			}
			flags := MODERN_HOLIDAY
			if !h.chul {
				flags |= IL_ONLY
			}
			events = append(events, HolidayEvent{
				Date:  hd,
				Desc:  h.desc,
				Flags: flags,
				Emoji: emoji})
		}
	}

	// Rosh Chodesh
	monthsInYear := hdate.MonthsInYear(year)
	for i := 1; i <= monthsInYear; i++ {
		isNisan := i == 1
		prevMonthNum := i - 1
		prevMonthYear := year
		if isNisan {
			prevMonthYear = year - 1
			prevMonthNum = hdate.MonthsInYear(prevMonthYear)
		}
		prevMonth := hdate.HMonth(prevMonthNum)
		prevMonthNumDays := hdate.DaysInMonth(prevMonth, prevMonthYear)
		month := hdate.HMonth(i)
		monthName := hdate.New(year, month, 1).MonthName("en")
		desc := "Rosh Chodesh " + monthName
		if prevMonthNumDays == 30 {
			events = append(events,
				HolidayEvent{
					Date:  hdate.New(year, prevMonth, 30),
					Desc:  desc,
					Flags: ROSH_CHODESH,
					Emoji: "🌒"},
				HolidayEvent{
					Date:  hdate.New(year, month, 1),
					Desc:  desc,
					Flags: ROSH_CHODESH,
					Emoji: "🌒"},
			)
		} else if month != hdate.Tishrei {
			events = append(events,
				HolidayEvent{
					Date:  hdate.New(year, month, 1),
					Desc:  desc,
					Flags: ROSH_CHODESH,
					Emoji: "🌒"})
		}

		// Shabbat Mevarchim Chodesh
		if month == hdate.Elul {
			continue
		}
		nextMonthName, _ := nextMonthName(year, month)
		events = append(events,
			HolidayEvent{
				Date:  hdate.New(year, month, 29).OnOrBefore(time.Saturday),
				Desc:  "Shabbat Mevarchim Chodesh " + nextMonthName,
				Flags: SHABBAT_MEVARCHIM,
			})
	}

	// Begin: Yom Kippur Katan
	// start at hdate.Iyyar because one may not fast during hdate.Nisan
	for month := hdate.Iyyar; month <= hdate.HMonth(monthsInYear); month++ {
		nextMonthName, nextMonth := nextMonthName(year, month)
		// Yom Kippur Katan is not observed on the day before Rosh Hashanah.
		// Not observed prior to Rosh Chodesh hdate.Cheshvan because Yom Kippur has just passed.
		// Not observed before Rosh Chodesh hdate.Tevet, because that day is Hanukkah.
		if nextMonth == hdate.Tishrei || nextMonth == hdate.Cheshvan || nextMonth == hdate.Tevet {
			continue
		}
		ykk := hdate.New(year, month, 29)
		dow := ykk.Weekday()
		if dow == time.Friday || dow == time.Saturday {
			ykk = ykk.OnOrBefore(time.Thursday)
		}

		events = append(events,
			HolidayEvent{
				Date:  ykk,
				Desc:  "Yom Kippur Katan " + nextMonthName,
				Flags: MINOR_FAST | YOM_KIPPUR_KATAN})
	}

	sedra := sedra.New(year, false)
	beshalachHd, _ := sedra.FindParshaNum(16)
	events = append(events,
		HolidayEvent{
			Date:  beshalachHd,
			Desc:  "Shabbat Shirah",
			Flags: SPECIAL_SHABBAT,
			Emoji: "🕍"})

	// Birkat Hachamah appears only once every 28 years
	for day := 1; day <= 30; day++ {
		rataDie := hdate.HebrewToRD(year, hdate.Nisan, day)
		elapsed := rataDie + 1373429
		if elapsed%10227 == 172 {
			events = append(events,
				HolidayEvent{
					Date:  hdate.FromRD(rataDie),
					Desc:  "Birkat Hachamah",
					Flags: MINOR_HOLIDAY,
					Emoji: "☀️"})
		}
	}

	sort.Sort(byDate(events))
	return events
}

// Returns a slice of holidays for the year.
// For Israel holiday schedule, specify il=true.
func GetHolidaysForYear(year int, il bool) []HolidayEvent {
	events := getAllHolidaysForYear(year)
	result := make([]HolidayEvent, 0, len(events))
	for _, ev := range events {
		if (il && (ev.Flags&CHUL_ONLY) == 0) || (!il && (ev.Flags&IL_ONLY) == 0) {
			result = append(result, ev)
		}
	}
	return result
}
