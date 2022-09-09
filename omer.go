package hebcal

// Hebcal - A Jewish Calendar Generator
// Copyright (c) 2022 Michael J. Radwin
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

import "strconv"

type omerEvent struct {
	Date            HDate
	OmerDay         int
	WeekNumber      int
	DaysWithinWeeks int
}

func NewOmerEvent(hd HDate, omerDay int) omerEvent {
	if omerDay < 1 || omerDay > 49 {
		panic("invalid omerDay")
	}
	week := ((omerDay - 1) / 7) + 1
	days := (omerDay % 7)
	if days == 0 {
		days = 7
	}
	return omerEvent{Date: hd, OmerDay: omerDay, WeekNumber: week, DaysWithinWeeks: days}
}

func (ev omerEvent) GetDate() HDate {
	return ev.Date
}

func (ev omerEvent) Render() string {
	return strconv.Itoa(ev.OmerDay) + " day of the Omer"
}

func (ev omerEvent) GetFlags() HolidayFlags {
	return OMER_COUNT
}

func (ev omerEvent) GetEmoji() string {
	return ""
}
