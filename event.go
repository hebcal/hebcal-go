package hebcal

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hebcal/hebcal-go/dafyomi"
	"github.com/hebcal/hebcal-go/hdate"
	"github.com/hebcal/hebcal-go/locales"
	"github.com/hebcal/hebcal-go/mishnayomi"
	"github.com/hebcal/hebcal-go/sedra"
)

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

type HolidayFlags uint32

const (
	// Chag, yontiff, yom tov
	CHAG HolidayFlags = 1 << iota
	// Light candles 18 minutes before sundown
	LIGHT_CANDLES
	// End of holiday (end of Yom Tov)
	YOM_TOV_ENDS
	// Observed only in the Diaspora (chutz l'aretz)
	CHUL_ONLY
	// Observed only in Israel
	IL_ONLY
	// Light candles in the evening at Tzeit time (3 small stars)
	LIGHT_CANDLES_TZEIS
	// Candle-lighting for Chanukah
	CHANUKAH_CANDLES
	// Rosh Chodesh, beginning of a new Hebrew month
	ROSH_CHODESH
	// Minor fasts like Tzom Tammuz, Ta'anit Esther, ...
	MINOR_FAST
	// Shabbat Shekalim, Zachor, ...
	SPECIAL_SHABBAT
	// Weekly sedrot on Saturdays
	PARSHA_HASHAVUA
	// Daily page of Talmud
	DAF_YOMI
	// Days of the Omer
	OMER_COUNT
	// Yom HaShoah, Yom HaAtzma'ut, ...
	MODERN_HOLIDAY
	// Yom Kippur and Tish'a B'Av
	MAJOR_FAST
	// On the Saturday before Rosh Chodesh
	SHABBAT_MEVARCHIM
	// Molad
	MOLAD
	// Yahrzeit or Hebrew Anniversary
	USER_EVENT
	// Daily Hebrew date ("11th of Sivan, 5780")
	HEBREW_DATE
	// A holiday that's not major, modern, rosh chodesh, or a fast day
	MINOR_HOLIDAY
	// Evening before a major or minor holiday
	EREV
	// Chol haMoed, intermediate days of Pesach or Sukkot
	CHOL_HAMOED
	// Mishna Yomi
	MISHNA_YOMI
	// Yom Kippur Katan, minor day of atonement on the day preceding each Rosh Chodesh
	YOM_KIPPUR_KATAN
	// Zemanim, halachic times of day
	ZMANIM
)

type HEvent interface {
	GetDate() hdate.HDate        // Holiday date of occurrence
	Render(locale string) string // Description (e.g. "Pesach III (CH''M)")
	GetFlags() HolidayFlags      // Event flag bitmask
	GetEmoji() string            // Holiday-specific emoji
	// Returns a simplified (untranslated) description for this event.
	// For example, HolidayEvent supports "Erev Pesach" => "Pesach",
	// and "Sukkot III (CH''M)" => "Sukkot".
	// For many holidays the basename and the event description are
	// the same.
	Basename() string
}

// HolidayEvent represents a built-in holiday like Pesach, Purim or Tu BiShvat
type HolidayEvent struct {
	Date          hdate.HDate  // Holiday date of occurrence
	Desc          string       // Description (e.g. "Pesach III (CH''M)")
	Flags         HolidayFlags // Event flag bitmask
	Emoji         string       // Holiday-specific emoji
	CholHaMoedDay int          // used only for Pesach and Sukkot
	ChanukahDay   int          // used only for Chanukah
}

func (ev HolidayEvent) GetDate() hdate.HDate {
	return ev.Date
}

func (ev HolidayEvent) Render(locale string) string {
	if (ev.Flags & ROSH_CHODESH) != 0 {
		rchStr, _ := locales.LookupTranslation("Rosh Chodesh", locale)
		monthStr, _ := locales.LookupTranslation(ev.Desc[13:], locale)
		return rchStr + " " + monthStr
	}
	str, _ := locales.LookupTranslation(ev.Desc, locale)
	return str
}

func (ev HolidayEvent) GetFlags() HolidayFlags {
	return ev.Flags
}

func (ev HolidayEvent) GetEmoji() string {
	return ev.Emoji
}

var regexes = []*regexp.Regexp{
	regexp.MustCompile(` \d{4}$`),
	regexp.MustCompile(` \(CH''M\)$`),
	regexp.MustCompile(` \(observed\)$`),
	regexp.MustCompile(` \(Hoshana Raba\)$`),
	regexp.MustCompile(` [IV]+$`),
	regexp.MustCompile(`: \d Candles?$`),
	regexp.MustCompile(`: 8th Day$`),
	regexp.MustCompile(`^Erev `),
}

func (ev HolidayEvent) Basename() string {
	str := ev.Desc
	for _, regex := range regexes {
		str = regex.ReplaceAllString(str, "")
	}
	return str
}

func getEnOrdinal(n int) string {
	str := strconv.Itoa(n)
	i := n % 100
	if i/10 == 1 {
		return str + "th"
	}
	switch i % 10 {
	case 1:
		return str + "st"
	case 2:
		return str + "nd"
	case 3:
		return str + "rd"
	default:
		return str + "th"
	}
}

type hebrewDateEvent struct {
	Date hdate.HDate
}

func (ev hebrewDateEvent) GetDate() hdate.HDate {
	return ev.Date
}

func (ev hebrewDateEvent) Render(locale string) string {
	hd := ev.Date
	enMonthName := hd.MonthName("en")
	switch locale {
	case "he":
		return Gematriya(hd.Day) + " " + hd.MonthName("he") + " " + Gematriya(hd.Year)
	case "", "en", "sephardic", "ashkenazi",
		"ashkenazi_litvish", "ashkenazi_poylish", "ashkenazi_standard":
		return getEnOrdinal(hd.Day) + " of " + enMonthName +
			", " + strconv.Itoa(hd.Year)
	case "es":
		monthName, _ := locales.LookupTranslation(enMonthName, locale)
		return strconv.Itoa(hd.Day) + "º " + monthName + " " + strconv.Itoa(hd.Year)

	}
	monthName, _ := locales.LookupTranslation(enMonthName, locale)
	return strconv.Itoa(hd.Day) + " " + monthName + " " + strconv.Itoa(hd.Year)
}

func (ev hebrewDateEvent) GetFlags() HolidayFlags {
	return HEBREW_DATE
}

func (ev hebrewDateEvent) GetEmoji() string {
	return ""
}

func (ev hebrewDateEvent) Basename() string {
	return ev.Date.String()
}

// Represents one of 54 weekly Torah portions, always on a Saturday
type parshaEvent struct {
	Date   hdate.HDate
	Parsha sedra.Parsha
	IL     bool
}

func (ev parshaEvent) GetDate() hdate.HDate {
	return ev.Date
}

func (ev parshaEvent) Render(locale string) string {
	prefix, _ := locales.LookupTranslation("Parashat", locale)
	name, _ := locales.LookupTranslation(ev.Parsha.Name[0], locale)
	if len(ev.Parsha.Name) == 2 {
		p2, _ := locales.LookupTranslation(ev.Parsha.Name[1], locale)
		delim := "-"
		if locale == "he" {
			delim = "־"
		}
		name = name + delim + p2
	}
	return prefix + " " + name
}

func (ev parshaEvent) GetFlags() HolidayFlags {
	return PARSHA_HASHAVUA
}

func (ev parshaEvent) GetEmoji() string {
	return ""
}

func (ev parshaEvent) Basename() string {
	return strings.Join(ev.Parsha.Name, "-")
}

type dafYomiEvent struct {
	Date hdate.HDate
	Daf  dafyomi.DafYomi
}

func (ev dafYomiEvent) GetDate() hdate.HDate {
	return ev.Date
}

func (ev dafYomiEvent) Render(locale string) string {
	name, _ := locales.LookupTranslation(ev.Daf.Name, locale)
	if locale == "he" {
		return name + " דף " + Gematriya(ev.Daf.Blatt)
	}
	return name + " " + strconv.Itoa(ev.Daf.Blatt)
}

func (ev dafYomiEvent) GetFlags() HolidayFlags {
	return DAF_YOMI
}

func (ev dafYomiEvent) GetEmoji() string {
	return ""
}

func (ev dafYomiEvent) Basename() string {
	return ev.Daf.String()
}

type mishnaYomiEvent struct {
	Date   hdate.HDate
	Mishna mishnayomi.MishnaPair
}

func (ev mishnaYomiEvent) GetDate() hdate.HDate {
	return ev.Date
}

func (ev mishnaYomiEvent) Render(locale string) string {
	m1 := ev.Mishna[0]
	m2 := ev.Mishna[1]
	tractate, _ := locales.LookupTranslation(m1.Tractate, locale)
	s := tractate + " " + strconv.Itoa(m1.Chap) + ":" + strconv.Itoa(m1.Verse) + "-"
	sameTractate := m1.Tractate == m2.Tractate
	if !sameTractate {
		tractate, _ := locales.LookupTranslation(m2.Tractate, locale)
		s += tractate + " "
	}
	if sameTractate && m2.Chap == m1.Chap {
		s += strconv.Itoa(m2.Verse)
	} else {
		s += strconv.Itoa(m2.Chap) + ":" + strconv.Itoa(m2.Verse)
	}
	return s
}

func (ev mishnaYomiEvent) GetFlags() HolidayFlags {
	return MISHNA_YOMI
}

func (ev mishnaYomiEvent) GetEmoji() string {
	return ""
}

func (ev mishnaYomiEvent) Basename() string {
	return ev.Mishna.String()
}

type moladEvent struct {
	Date      hdate.HDate
	Molad     hdate.Molad
	MonthName string
}

func (ev moladEvent) GetDate() hdate.HDate {
	return ev.Date
}

func (ev moladEvent) Render(locale string) string {
	return fmt.Sprintf("Molad %s: %s, %d minutes and %d chalakim after %d:00",
		ev.MonthName, ev.Molad.Date.Weekday().String()[0:3],
		ev.Molad.Minutes, ev.Molad.Chalakim, ev.Molad.Hours)
}

func (ev moladEvent) GetFlags() HolidayFlags {
	return MOLAD
}

func (ev moladEvent) GetEmoji() string {
	return ""
}

func (ev moladEvent) Basename() string {
	return "Molad " + ev.MonthName
}
