package event

import (
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/hebcal/gematriya"
	"github.com/hebcal/hdate"
	"github.com/hebcal/locales"
)

type hebrewDateEvent struct {
	Date hdate.HDate
}

func (ev hebrewDateEvent) GetDate() hdate.HDate {
	return ev.Date
}

func (ev hebrewDateEvent) Render(locale string) string {
	return renderHebrewDate(ev.Date, locale, true)
}

// RenderBrief returns the date without its year, e.g. "15th of Cheshvan" or
// "ט״ו חֶשְׁוָן". 1 Tishrei, the first day of the year, keeps it.
func (ev hebrewDateEvent) RenderBrief(locale string) string {
	hd := ev.Date
	showYear := hd.Month() == hdate.Tishrei && hd.Day() == 1
	return renderHebrewDate(hd, locale, showYear)
}

// renderHebrewDate renders hd as "15th of Cheshvan, 5769", "12. Tewet, 5786"
// or "ט״ו חֶשְׁוָן תשס״ט", with or without the year. The month name is
// translated from the locale catalog, and in Latin-script locales its
// apostrophe is typographic ("Sh’vat").
func renderHebrewDate(hd hdate.HDate, locale string, showYear bool) string {
	locale = strings.ToLower(locale)
	day := hd.Day()
	month := monthName(hd, locale)
	switch locale {
	case "he", "he-x-nonikud":
		s := gematriya.Gematriya(day) + " " + month
		if showYear {
			s += " " + gematriya.Gematriya(hd.Year())
		}
		return s
	}
	var nth, of string
	switch {
	case locale == "", locale == "en", locale == "sephardic", strings.HasPrefix(locale, "ashkenazi"):
		nth, of = humanize.Ordinal(day), " of"
	case locale == "es":
		nth = strconv.Itoa(day) + "º"
	default:
		nth = strconv.Itoa(day) + "."
	}
	s := nth + of + " " + strings.ReplaceAll(month, "'", "’")
	if showYear {
		s += ", " + strconv.Itoa(hd.Year())
	}
	return s
}

// monthName returns hd's month translated to locale. The untranslated name
// is spelled "Tamuz", with one m.
func monthName(hd hdate.HDate, locale string) string {
	en := hd.MonthName("en")
	if en == "Tammuz" {
		en = "Tamuz"
	}
	if s, ok := locales.LookupTranslation(en, locale); ok && s != "" {
		return s
	}
	return en
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

func (ev hebrewDateEvent) GetCategories() []string {
	return CategoriesFromFlags(ev.GetFlags())
}

func NewHebrewDateEvent(hd hdate.HDate) CalEvent {
	return hebrewDateEvent{Date: hd}
}
