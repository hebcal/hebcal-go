package event

import (
	"strconv"

	"github.com/hebcal/gematriya"
	"github.com/hebcal/hdate"
	"github.com/MaxBGreenberg/hebcal-go/dafyomi"
	"github.com/MaxBGreenberg/hebcal-go/locales"
)

type yyomiEvent struct {
	Date hdate.HDate
	Daf  dafyomi.Daf
}

func NewYerushalmiYomiEvent(hd hdate.HDate, daf dafyomi.Daf) CalEvent {
	return yyomiEvent{Date: hd, Daf: daf}
}

func (ev yyomiEvent) GetDate() hdate.HDate {
	return ev.Date
}

func (ev yyomiEvent) Render(locale string) string {
	yerushalmiStr, _ := locales.LookupTranslation("Yerushalmi", locale)
	name, _ := locales.LookupTranslation(ev.Daf.Name, locale)
	if locale == "he" {
		return yerushalmiStr + " " + name + " דף " + gematriya.Gematriya(ev.Daf.Blatt)
	}
	return yerushalmiStr + " " + name + " " + strconv.Itoa(ev.Daf.Blatt)
}

func (ev yyomiEvent) GetFlags() HolidayFlags {
	return YERUSHALMI_YOMI
}

func (ev yyomiEvent) GetEmoji() string {
	return ""
}

func (ev yyomiEvent) Basename() string {
	return ev.Daf.String()
}
