package event

import (
	"fmt"

	"github.com/hebcal/hdate"
	"github.com/hebcal/hebcal-go/molad"
)

type moladEvent struct {
	Date      hdate.HDate
	Molad     molad.Molad
	MonthName string
}

func NewMoladEvent(date hdate.HDate, m molad.Molad, monthName string) CalEvent {
	return moladEvent{
		Date:      date,
		Molad:     m,
		MonthName: monthName,
	}
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
