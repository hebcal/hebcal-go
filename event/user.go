package event

import (
	"github.com/hebcal/hdate"
)

// UserEvent is a personal event: a Yahrzeit, a birthday, an anniversary.
// Construct one with NewUserEvent.
type UserEvent struct {
	date hdate.HDate
	desc string
}

// NewUserEvent constructs a personal event on date. desc is its description,
// which is not translated.
func NewUserEvent(date hdate.HDate, desc string) UserEvent {
	return UserEvent{date: date, desc: desc}
}

// GetDate returns the Hebrew date on which the event occurs.
func (ev UserEvent) GetDate() hdate.HDate {
	return ev.date
}

// Desc returns the description the event was constructed with.
func (ev UserEvent) Desc() string {
	return ev.desc
}

// Render returns the description, whatever the locale.
func (ev UserEvent) Render(locale string) string {
	return ev.desc
}

// URL returns "": a personal event has no page.
func (ev UserEvent) URL() string {
	return ""
}

func (ev UserEvent) GetFlags() Flags {
	return User
}

func (ev UserEvent) GetEmoji() string {
	return ""
}

func (ev UserEvent) Basename() string {
	return ev.desc
}

func (ev UserEvent) GetCategories() []string {
	return CategoriesFromFlags(ev.GetFlags())
}
