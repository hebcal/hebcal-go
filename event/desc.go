package event

// Descer is implemented by events that can report their untranslated
// description, the stable identifier @hebcal/core calls getDesc(): "Erev
// Pesach", "Parashat Bereshit", "Omer 7", "29 Elul 5786", "Molad Cheshvan
// 5787". Unlike Basename, it tells apart the days of one holiday.
//
// Like URLer, it is a separate interface rather than a method on CalEvent,
// which is implemented outside this package; use the Desc helper.
type Descer interface {
	GetDesc() string
}

// Desc returns ev's untranslated description, or its Basename when ev does
// not implement Descer.
func Desc(ev CalEvent) string {
	if d, ok := ev.(Descer); ok {
		return d.GetDesc()
	}
	return ev.Basename()
}

// Memoer is implemented by events that carry a longer note of their own: the
// molad announced on Shabbat Mevarchim, the meaning of Yom Kippur Katan, the
// Torah portion of the Shabbat that erev-Shabbat candle-lighting begins, the
// Kiddush Levana times, a daily-learning reading's details.
type Memoer interface {
	// Memo returns the note in locale, or "" if the event has none.
	Memo(locale string) string
}

// Memo returns ev's memo in locale, or "" if it has none.
func Memo(ev CalEvent, locale string) string {
	if m, ok := ev.(Memoer); ok {
		return m.Memo(locale)
	}
	return ""
}

func (ev HolidayEvent) GetDesc() string          { return ev.Desc }
func (ev UserEvent) GetDesc() string             { return ev.Desc }
func (ev HebrewDateEvent) GetDesc() string       { return ev.date.String() }
func (ev MevarchimChodeshEvent) GetDesc() string { return ev.Basename() }

// GetDesc returns "Parashat " and the portion's untranslated name, e.g.
// "Parashat Vayakhel-Pekudei".
func (ev ParshaEvent) GetDesc() string { return "Parashat " + ev.Basename() }

// Memo returns the meaning of Yom Kippur Katan, untranslated, e.g. "Minor Day
// of Atonement on the day preceeding Rosh Chodesh Sivan" (sic). The month
// name carries a typographic apostrophe, as in a rendered title ("Sh’vat").
// Other holidays have no memo of their own.
func (ev HolidayEvent) Memo(locale string) string {
	if !ev.Flags.Has(YOM_KIPPUR_KATAN) || len(ev.Desc) <= len(yomKippurKatanPrefix) {
		return ""
	}
	return "Minor Day of Atonement on the day preceeding Rosh Chodesh " +
		SmartApostrophe(ev.Desc[len(yomKippurKatanPrefix):])
}

const yomKippurKatanPrefix = "Yom Kippur Katan "
