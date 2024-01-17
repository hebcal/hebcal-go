package event

import (
	"strconv"

	"github.com/hebcal/gematriya"
	"github.com/hebcal/hdate"
	"github.com/MaxBGreenberg/hebcal-go/dafyomi"
	"github.com/MaxBGreenberg/hebcal-go/locales"
)

type dafYomiEvent struct {
	Date hdate.HDate
	Daf  dafyomi.Daf
}

func NewDafYomiEvent(hd hdate.HDate, daf dafyomi.Daf) CalEvent {
	return dafYomiEvent{Date: hd, Daf: daf}
}

func (ev dafYomiEvent) GetDate() hdate.HDate {
	return ev.Date
}

func (ev dafYomiEvent) Render(locale string) string {
	name, _ := locales.LookupTranslation(ev.Daf.Name, locale)
	if locale == "he" {
		return name + " דף " + gematriya.Gematriya(ev.Daf.Blatt)
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
