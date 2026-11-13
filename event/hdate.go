package event

import (
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/hebcal/gematriya"
	"github.com/hebcal/hdate"
	"github.com/hebcal/locales"
)

// HebrewDateEvent is a day's Hebrew date as an event of its own, e.g. "15th
// of Cheshvan, 5769".
type HebrewDateEvent struct {
	date hdate.HDate
}

// GetDate returns the Hebrew date the event names.
func (ev HebrewDateEvent) GetDate() hdate.HDate {
	return ev.date
}

func (ev HebrewDateEvent) Render(locale string) string {
	return renderHebrewDate(ev.date, locale, true)
}

// RenderBrief returns the date without its year, e.g. "15th of Cheshvan" or
// "ט״ו חֶשְׁוָן". 1 Tishrei, the first day of the year, keeps it.
func (ev HebrewDateEvent) RenderBrief(locale string) string {
	hd := ev.date
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
	s := nth + of + " " + SmartApostrophe(month)
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

func (ev HebrewDateEvent) GetFlags() HolidayFlags {
	return HEBREW_DATE
}

func (ev HebrewDateEvent) GetEmoji() string {
	return ""
}

func (ev HebrewDateEvent) Basename() string {
	return ev.date.String()
}

func (ev HebrewDateEvent) GetCategories() []string {
	return CategoriesFromFlags(ev.GetFlags())
}

// NewHebrewDateEvent constructs the Hebrew date event for hd.
func NewHebrewDateEvent(hd hdate.HDate) HebrewDateEvent {
	return HebrewDateEvent{date: hd}
}
