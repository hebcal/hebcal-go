package hebcal

import (
	"sort"
	"strconv"
	"time"
)

type HolidayFlags uint32

const (
	/* Chag, yontiff, yom tov */
	CHAG HolidayFlags = 1 << iota
	/* Light candles 18 minutes before sundown */
	LIGHT_CANDLES
	/* End of holiday (end of Yom Tov)  */
	YOM_TOV_ENDS
	/* Observed only in the Diaspora (chutz l'aretz)  */
	CHUL_ONLY
	/* Observed only in Israel */
	IL_ONLY
	/* Light candles in the evening at Tzeit time (3 small stars) */
	LIGHT_CANDLES_TZEIS
	/* Candle-lighting for Chanukah */
	CHANUKAH_CANDLES
	/* Rosh Chodesh, beginning of a new Hebrew month */
	ROSH_CHODESH
	/* Minor fasts like Tzom Tammuz, Ta'anit Esther, ... */
	MINOR_FAST
	/* Shabbat Shekalim, Zachor, ... */
	SPECIAL_SHABBAT
	/* Weekly sedrot on Saturdays */
	PARSHA_HASHAVUA
	/* Daily page of Talmud */
	DAF_YOMI
	/* Days of the Omer */
	OMER_COUNT
	/* Yom HaShoah, Yom HaAtzma'ut, ... */
	MODERN_HOLIDAY
	/* Yom Kippur and Tish'a B'Av */
	MAJOR_FAST
	/* On the Saturday before Rosh Chodesh */
	SHABBAT_MEVARCHIM
	/* Molad */
	MOLAD
	/* Yahrzeit or Hebrew Anniversary */
	USER_EVENT
	/* Daily Hebrew date ("11th of Sivan, 5780") */
	HEBREW_DATE
	/* A holiday that's not major, modern, rosh chodesh, or a fast day */
	MINOR_HOLIDAY
	/* Evening before a major or minor holiday */
	EREV
	/* Chol haMoed, intermediate days of Pesach or Sukkot */
	CHOL_HAMOED
	/* Mishna Yomi */
	MISHNA_YOMI
	/* Yom Kippur Katan, minor day of atonement on the day preceeding each Rosh Chodesh */
	YOM_KIPPUR_KATAN
)

type HEvent struct {
	Date          HDate
	Desc          string
	Flags         HolidayFlags
	Emoji         string
	CholHaMoedDay int
	ChanukahDay   int
}

type holiday struct {
	mm     HMonth
	dd     int
	desc   string
	flags  HolidayFlags
	emoji  string
	chmDay int
}

const chanukahEmoji = "🕎"

var staticHolidays = []holiday{
	{mm: Tishrei, dd: 2, desc: "Rosh Hashana II", flags: CHAG | YOM_TOV_ENDS, emoji: "🍏🍯"},
	{mm: Tishrei, dd: 9, desc: "Erev Yom Kippur", flags: EREV | LIGHT_CANDLES},
	{mm: Tishrei, dd: 10, desc: "Yom Kippur", flags: CHAG | MAJOR_FAST | YOM_TOV_ENDS},
	{mm: Tishrei, dd: 14, desc: "Erev Sukkot", flags: EREV | LIGHT_CANDLES},

	{mm: Tishrei, dd: 15, desc: "Sukkot I", flags: CHUL_ONLY | CHAG | LIGHT_CANDLES_TZEIS},
	{mm: Tishrei, dd: 16, desc: "Sukkot II", flags: CHUL_ONLY | CHAG | YOM_TOV_ENDS},
	{mm: Tishrei, dd: 17, desc: "Sukkot III (CH''M)", flags: CHUL_ONLY | CHOL_HAMOED, chmDay: 1},
	{mm: Tishrei, dd: 18, desc: "Sukkot IV (CH''M)", flags: CHUL_ONLY | CHOL_HAMOED, chmDay: 2},
	{mm: Tishrei, dd: 19, desc: "Sukkot V (CH''M)", flags: CHUL_ONLY | CHOL_HAMOED, chmDay: 3},
	{mm: Tishrei, dd: 20, desc: "Sukkot VI (CH''M)", flags: CHUL_ONLY | CHOL_HAMOED, chmDay: 4},
	{mm: Tishrei, dd: 22, desc: "Shmini Atzeret",
		flags: CHUL_ONLY | CHAG | LIGHT_CANDLES_TZEIS},
	{mm: Tishrei, dd: 23, desc: "Simchat Torah",
		flags: CHUL_ONLY | CHAG | YOM_TOV_ENDS},

	{mm: Tishrei, dd: 15, desc: "Sukkot I", flags: IL_ONLY | CHAG | YOM_TOV_ENDS},
	{mm: Tishrei, dd: 16, desc: "Sukkot II (CH''M)", flags: IL_ONLY | CHOL_HAMOED, chmDay: 1},
	{mm: Tishrei, dd: 17, desc: "Sukkot III (CH''M)", flags: IL_ONLY | CHOL_HAMOED, chmDay: 2},
	{mm: Tishrei, dd: 18, desc: "Sukkot IV (CH''M)", flags: IL_ONLY | CHOL_HAMOED, chmDay: 3},
	{mm: Tishrei, dd: 19, desc: "Sukkot V (CH''M)", flags: IL_ONLY | CHOL_HAMOED, chmDay: 4},
	{mm: Tishrei, dd: 20, desc: "Sukkot VI (CH''M)", flags: IL_ONLY | CHOL_HAMOED, chmDay: 5},
	{mm: Tishrei, dd: 22, desc: "Shmini Atzeret",
		flags: IL_ONLY | CHAG | YOM_TOV_ENDS},

	{mm: Tishrei, dd: 21, desc: "Sukkot VII (Hoshana Raba)",
		flags: LIGHT_CANDLES | CHOL_HAMOED, chmDay: -1},
	{mm: Kislev, dd: 24, desc: "Chanukah: 1 Candle",
		flags: EREV | MINOR_HOLIDAY | CHANUKAH_CANDLES, emoji: chanukahEmoji},
	{mm: Tevet, dd: 10, desc: "Asara B'Tevet", flags: MINOR_FAST},
	{mm: Shvat, dd: 15, desc: "Tu BiShvat", flags: MINOR_HOLIDAY, emoji: "🌳"},
	{mm: Adar2, dd: 13, desc: "Erev Purim", flags: EREV | MINOR_HOLIDAY, emoji: "🎭️📜"},
	{mm: Adar2, dd: 14, desc: "Purim", flags: MINOR_HOLIDAY, emoji: "🎭️📜"},
	{mm: Nisan, dd: 14, desc: "Erev Pesach", flags: EREV | LIGHT_CANDLES},
	// Pesach Israel
	{mm: Nisan, dd: 15, desc: "Pesach I",
		flags: IL_ONLY | CHAG | YOM_TOV_ENDS},
	{mm: Nisan, dd: 16, desc: "Pesach II (CH''M)",
		flags: IL_ONLY | CHOL_HAMOED, chmDay: 1},
	{mm: Nisan, dd: 17, desc: "Pesach III (CH''M)",
		flags: IL_ONLY | CHOL_HAMOED, chmDay: 2},
	{mm: Nisan, dd: 18, desc: "Pesach IV (CH''M)",
		flags: IL_ONLY | CHOL_HAMOED, chmDay: 3},
	{mm: Nisan, dd: 19, desc: "Pesach V (CH''M)",
		flags: IL_ONLY | CHOL_HAMOED, chmDay: 4},
	{mm: Nisan, dd: 20, desc: "Pesach VI (CH''M)",
		flags: IL_ONLY | CHOL_HAMOED | LIGHT_CANDLES, chmDay: 5},
	{mm: Nisan, dd: 21, desc: "Pesach VII",
		flags: IL_ONLY | CHAG | YOM_TOV_ENDS},
	// Pesach chutz l'aretz
	{mm: Nisan, dd: 15, desc: "Pesach I",
		flags: CHUL_ONLY | CHAG | LIGHT_CANDLES_TZEIS},
	{mm: Nisan, dd: 16, desc: "Pesach II",
		flags: CHUL_ONLY | CHAG | YOM_TOV_ENDS},
	{mm: Nisan, dd: 17, desc: "Pesach III (CH''M)",
		flags: CHUL_ONLY | CHOL_HAMOED, chmDay: 1},
	{mm: Nisan, dd: 18, desc: "Pesach IV (CH''M)",
		flags: CHUL_ONLY | CHOL_HAMOED, chmDay: 2},
	{mm: Nisan, dd: 19, desc: "Pesach V (CH''M)",
		flags: CHUL_ONLY | CHOL_HAMOED, chmDay: 3},
	{mm: Nisan, dd: 20, desc: "Pesach VI (CH''M)",
		flags: CHUL_ONLY | CHOL_HAMOED | LIGHT_CANDLES, chmDay: 4},
	{mm: Nisan, dd: 21, desc: "Pesach VII",
		flags: CHUL_ONLY | CHAG | LIGHT_CANDLES_TZEIS},
	{mm: Nisan, dd: 22, desc: "Pesach VIII",
		flags: CHUL_ONLY | CHAG | YOM_TOV_ENDS},

	{mm: Iyyar, dd: 14, desc: "Pesach Sheni", flags: MINOR_HOLIDAY},
	{mm: Iyyar, dd: 18, desc: "Lag BaOmer", flags: MINOR_HOLIDAY, emoji: "🔥"},
	{mm: Sivan, dd: 5, desc: "Erev Shavuot",
		flags: EREV | LIGHT_CANDLES, emoji: "⛰️🌸"},
	{mm: Sivan, dd: 6, desc: "Shavuot",
		flags: IL_ONLY | CHAG | YOM_TOV_ENDS, emoji: "⛰️🌸"},
	{mm: Sivan, dd: 6, desc: "Shavuot I",
		flags: CHUL_ONLY | CHAG | LIGHT_CANDLES_TZEIS, emoji: "⛰️🌸"},
	{mm: Sivan, dd: 7, desc: "Shavuot II",
		flags: CHUL_ONLY | CHAG | YOM_TOV_ENDS, emoji: "⛰️🌸"},
	{mm: Av, dd: 15, desc: "Tu B'Av",
		flags: MINOR_HOLIDAY, emoji: "❤️"},
	{mm: Elul, dd: 1, desc: "Rosh Hashana LaBehemot",
		flags: MINOR_HOLIDAY, emoji: "🐑"},
	{mm: Elul, dd: 29, desc: "Erev Rosh Hashana",
		flags: EREV | LIGHT_CANDLES, emoji: "🍏🍯"},
}

func tzomGedaliahDate(rh HDate) HDate {
	offset := 0
	if rh.Weekday() == time.Thursday {
		offset = 1
	}
	return NewHDate(rh.Year, Tishrei, 3+offset)
}

func taanitEstherDate(pesach HDate) HDate {
	offset := 31
	if pesach.Weekday() == time.Tuesday {
		offset = 33
	}
	return NewHDateFromRD(pesach.Abs() - offset)
}

func shushanPurimDate(pesach HDate) HDate {
	offset := 29
	if pesach.Weekday() == time.Sunday {
		offset = 28
	}
	return NewHDateFromRD(pesach.Abs() - offset)
}

func taanitBechorotDate(pesach HDate) HDate {
	if pesach.Prev().Weekday() == time.Saturday {
		return NewHDateFromRD(dayOnOrBefore(time.Thursday, pesach.Abs()))
	} else {
		return NewHDate(pesach.Year, Nisan, 14)
	}
}

type byDate []HEvent

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

func nextMonthName(year int, month HMonth) (string, HMonth) {
	nextMonth := month + 1
	monthsInYear := MonthsInHebYear(year)
	if month == HMonth(monthsInYear) {
		nextMonth = Nisan
	}
	nextMonthName := NewHDate(year, nextMonth, 1).MonthName()
	return nextMonthName, nextMonth
}

func GetAllHolidaysForYear(year int) []HEvent {
	events := make([]HEvent, 0, 120)
	// standard holidays that don't shift based on year
	for _, h := range staticHolidays {
		events = append(events, HEvent{Date: NewHDate(year, h.mm, h.dd),
			Desc: h.desc, Flags: h.flags, Emoji: h.emoji, CholHaMoedDay: h.chmDay})
	}
	// variable holidays
	roshHashana := NewHDate(year, Tishrei, 1)
	nextRh := NewHDate(year+1, Tishrei, 1)
	pesach := NewHDate(year, Nisan, 15)
	pesachAbs := pesach.Abs()
	events = append(events,
		HEvent{
			Date:  roshHashana,
			Desc:  "Rosh Hashana " + strconv.Itoa(year),
			Flags: CHAG | LIGHT_CANDLES_TZEIS,
			Emoji: "🍏🍯"},
		HEvent{
			Date:  NewHDateFromRD(dayOnOrBefore(time.Saturday, 7+roshHashana.Abs())),
			Desc:  "Shabbat Shuva",
			Flags: SPECIAL_SHABBAT,
			Emoji: "🕍"},
		HEvent{
			Date:  tzomGedaliahDate(roshHashana),
			Desc:  "Tzom Gedaliah",
			Flags: MINOR_FAST},
		HEvent{
			Date:  NewHDateFromRD(dayOnOrBefore(time.Saturday, pesachAbs-43)),
			Desc:  "Shabbat Shekalim",
			Flags: SPECIAL_SHABBAT,
			Emoji: "🕍"},
		HEvent{
			Date:  NewHDateFromRD(dayOnOrBefore(time.Saturday, pesachAbs-30)),
			Desc:  "Shabbat Zachor",
			Flags: SPECIAL_SHABBAT,
			Emoji: "🕍"},
		HEvent{
			Date:  taanitEstherDate(pesach),
			Desc:  "Ta'anit Esther",
			Flags: MINOR_FAST},
		HEvent{
			Date:  shushanPurimDate(pesach),
			Desc:  "Shushan Purim",
			Flags: MINOR_HOLIDAY, Emoji: "🎭️📜"},
		HEvent{
			Date:  NewHDateFromRD(dayOnOrBefore(time.Saturday, pesachAbs-14) - 7),
			Desc:  "Shabbat Parah",
			Flags: SPECIAL_SHABBAT,
			Emoji: "🕍"},
		HEvent{
			Date: NewHDateFromRD(dayOnOrBefore(time.Saturday, pesachAbs-14)),
			Desc: "Shabbat HaChodesh", Flags: SPECIAL_SHABBAT,
			Emoji: "🕍"},
		HEvent{
			Date: NewHDateFromRD(dayOnOrBefore(time.Saturday, pesachAbs-1)),
			Desc: "Shabbat HaGadol", Flags: SPECIAL_SHABBAT,
			Emoji: "🕍"},
		HEvent{
			Date:  taanitBechorotDate(pesach),
			Desc:  "Ta'anit Bechorot",
			Flags: MINOR_FAST},
		HEvent{
			Date:  NewHDateFromRD(dayOnOrBefore(time.Saturday, nextRh.Abs()-4)),
			Desc:  "Leil Selichot",
			Flags: MINOR_HOLIDAY,
			Emoji: "🕍"},
	)
	if IsHebLeapYear(year) {
		events = append(events, HEvent{
			Date:  NewHDate(year, Adar1, 14),
			Desc:  "Purim Katan",
			Flags: MINOR_HOLIDAY,
			Emoji: "🎭️"})
	}
	for candles := 2; candles <= 6; candles++ {
		events = append(events, HEvent{
			Date:        NewHDate(year, Kislev, 23+candles),
			Desc:        "Chanukah: " + strconv.Itoa(candles) + " Candles",
			Flags:       MINOR_HOLIDAY | CHANUKAH_CANDLES,
			Emoji:       chanukahEmoji,
			ChanukahDay: candles - 1})
	}
	var chanukah7 HDate
	if ShortKislev(year) {
		chanukah7 = NewHDate(year, Tevet, 1)
	} else {
		chanukah7 = NewHDate(year, Kislev, 30)
	}
	chanukah8 := chanukah7.Next()
	events = append(events,
		HEvent{
			Date:        chanukah7,
			Desc:        "Chanukah: 7 Candles",
			Flags:       MINOR_HOLIDAY | CHANUKAH_CANDLES,
			Emoji:       chanukahEmoji,
			ChanukahDay: 6},
		HEvent{
			Date:        chanukah8,
			Desc:        "Chanukah: 8 Candles",
			Flags:       MINOR_HOLIDAY | CHANUKAH_CANDLES,
			Emoji:       chanukahEmoji,
			ChanukahDay: 7},
		HEvent{
			Date:        chanukah8.Next(),
			Desc:        "Chanukah: 8th Day",
			Flags:       MINOR_HOLIDAY,
			Emoji:       chanukahEmoji,
			ChanukahDay: 8})

	// Tisha BAv and the 3 weeks
	var tamuz17 = NewHDate(year, Tamuz, 17)
	if tamuz17.Weekday() == time.Saturday {
		tamuz17 = tamuz17.Next()
	}
	events = append(events,
		HEvent{Date: tamuz17, Desc: "Tzom Tammuz", Flags: MINOR_FAST})

	var av9dt = NewHDate(year, Av, 9)
	var av9title = "Tish'a B'Av"
	if av9dt.Weekday() == time.Saturday {
		av9dt = av9dt.Next()
		av9title += " (observed)"
	}
	var av9abs = av9dt.Abs()
	events = append(events,
		HEvent{
			Date:  NewHDateFromRD(dayOnOrBefore(time.Saturday, av9abs)),
			Desc:  "Shabbat Chazon",
			Flags: SPECIAL_SHABBAT,
			Emoji: "🕍"},
		HEvent{Date: av9dt.Prev(), Desc: "Erev Tish'a B'Av", Flags: EREV | MAJOR_FAST},
		HEvent{Date: av9dt, Desc: av9title, Flags: MAJOR_FAST},
		HEvent{
			Date:  NewHDateFromRD(dayOnOrBefore(time.Saturday, av9abs+7)),
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
		var tmpDate = NewHDate(year, Iyyar, day)
		events = append(events,
			HEvent{Date: tmpDate, Desc: "Yom HaZikaron", Flags: MODERN_HOLIDAY, Emoji: "🇮🇱"},
			HEvent{Date: tmpDate.Next(), Desc: "Yom HaAtzma'ut", Flags: MODERN_HOLIDAY, Emoji: "🇮🇱"},
		)
	}

	if year >= 5711 {
		// Yom HaShoah first observed in 1951
		var nisan27dt = NewHDate(year, Nisan, 27)
		/* When the actual date of Yom Hashoah falls on a Friday, the
		 * state of Israel observes Yom Hashoah on the preceding
		 * Thursday. When it falls on a Sunday, Yom Hashoah is observed
		 * on the following Monday.
		 * http://www.ushmm.org/remembrance/dor/calendar/
		 */
		if nisan27dt.Weekday() == time.Friday {
			nisan27dt = nisan27dt.Prev()
		} else if nisan27dt.Weekday() == time.Sunday {
			nisan27dt = nisan27dt.Next()
		}
		events = append(events,
			HEvent{Date: nisan27dt, Desc: "Yom HaShoah", Flags: MODERN_HOLIDAY})
	}

	if year >= 5727 {
		// Yom Yerushalayim only celebrated after 1967
		events = append(events,
			HEvent{
				Date:  NewHDate(year, Iyyar, 28),
				Desc:  "Yom Yerushalayim",
				Flags: MODERN_HOLIDAY,
				Emoji: "🇮🇱"})
	}

	if year >= 5769 {
		events = append(events,
			HEvent{
				Date:  NewHDate(year, Cheshvan, 29),
				Desc:  "Sigd",
				Flags: MODERN_HOLIDAY})
	}

	if year >= 5777 {
		events = append(events,
			HEvent{
				Date:  NewHDate(year, Cheshvan, 7),
				Desc:  "Yom HaAliyah School Observance",
				Flags: MODERN_HOLIDAY,
				Emoji: "🇮🇱"},
			HEvent{
				Date:  NewHDate(year, Nisan, 10),
				Desc:  "Yom HaAliyah",
				Flags: MODERN_HOLIDAY,
				Emoji: "🇮🇱"},
		)
	}

	// Rosh Chodesh
	monthsInYear := MonthsInHebYear(year)
	for i := 1; i <= monthsInYear; i++ {
		isNisan := i == 1
		prevMonthNum := i - 1
		prevMonthYear := year
		if isNisan {
			prevMonthYear = year - 1
			prevMonthNum = MonthsInHebYear(prevMonthYear)
		}
		prevMonth := HMonth(prevMonthNum)
		prevMonthNumDays := DaysInMonth(prevMonth, prevMonthYear)
		month := HMonth(i)
		monthName := NewHDate(year, month, 1).MonthName()
		desc := "Rosh Chodesh " + monthName
		if prevMonthNumDays == 30 {
			events = append(events,
				HEvent{
					Date:  NewHDate(year, prevMonth, 30),
					Desc:  desc,
					Flags: ROSH_CHODESH,
					Emoji: "🌒"},
				HEvent{
					Date:  NewHDate(year, month, 1),
					Desc:  desc,
					Flags: ROSH_CHODESH,
					Emoji: "🌒"},
			)
		} else if month != Tishrei {
			events = append(events,
				HEvent{
					Date:  NewHDate(year, month, 1),
					Desc:  desc,
					Flags: ROSH_CHODESH,
					Emoji: "🌒"})
		}

		// Shabbat Mevarchim Chodesh
		if month == Elul {
			continue
		}
		nextMonthName, _ := nextMonthName(year, month)
		events = append(events,
			HEvent{
				Date:  NewHDate(year, month, 29).OnOrBefore(time.Saturday),
				Desc:  "Shabbat Mevarchim Chodesh " + nextMonthName,
				Flags: SHABBAT_MEVARCHIM,
			})
	}

	// Begin: Yom Kippur Katan
	// start at Iyyar because one may not fast during Nisan
	for month := Iyyar; month <= HMonth(monthsInYear); month++ {
		nextMonthName, nextMonth := nextMonthName(year, month)
		// Yom Kippur Katan is not observed on the day before Rosh Hashanah.
		// Not observed prior to Rosh Chodesh Cheshvan because Yom Kippur has just passed.
		// Not observed before Rosh Chodesh Tevet, because that day is Hanukkah.
		if nextMonth == Tishrei || nextMonth == Cheshvan || nextMonth == Tevet {
			continue
		}
		ykk := NewHDate(year, month, 29)
		dow := ykk.Weekday()
		if dow == time.Friday || dow == time.Saturday {
			ykk = ykk.OnOrBefore(time.Thursday)
		}

		events = append(events,
			HEvent{
				Date:  ykk,
				Desc:  "Yom Kippur Katan " + nextMonthName,
				Flags: MINOR_FAST | YOM_KIPPUR_KATAN})
	}

	sedra := NewSedra(year, false)
	beshalachHd, _ := sedra.FindParshaNum(16)
	events = append(events,
		HEvent{
			Date:  beshalachHd,
			Desc:  "Shabbat Shirah",
			Flags: SPECIAL_SHABBAT,
			Emoji: "🕍"})

	sort.Sort(byDate(events))
	return events
}

func GetHolidaysForYear(year int, il bool) []HEvent {
	events := GetAllHolidaysForYear(year)
	result := make([]HEvent, 0, len(events))
	for _, ev := range events {
		if (il && (ev.Flags&CHUL_ONLY) == 0) || (!il && (ev.Flags&IL_ONLY) == 0) {
			result = append(result, ev)
		}
	}
	return result
}
