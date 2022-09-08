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
package hebcal

import (
	"strconv"
	"time"
)

type HavdalahEvent struct {
	Date        HDate        // Holiday date of occurrence
	Flags       HolidayFlags // Event flag bitmask
	mins        int
	eventTime   time.Time
	linkedEvent *HolidayEvent
	loc         *HLocation
}

func (ev HavdalahEvent) GetDate() HDate {
	return ev.Date
}

func (ev HavdalahEvent) Render() string {
	prefix := "Havdalah"
	if ev.mins != 0 {
		prefix = "Havdalah (" + strconv.Itoa(ev.mins) + " mins)"
	}
	return prefix + ": " + ev.eventTime.Format(time.Kitchen)
}

func (ev HavdalahEvent) GetFlags() HolidayFlags {
	return ev.Flags
}

func (ev HavdalahEvent) GetEmoji() string {
	return "✨"
}

type CandleLightingEvent struct {
	Date        HDate        // Holiday date of occurrence
	Flags       HolidayFlags // Event flag bitmask
	eventTime   time.Time
	linkedEvent *HolidayEvent
	loc         *HLocation
}

func (ev CandleLightingEvent) GetDate() HDate {
	return ev.Date
}

func (ev CandleLightingEvent) Render() string {
	return "Candle lighting: " + ev.eventTime.Format(time.Kitchen)
}

func (ev CandleLightingEvent) GetFlags() HolidayFlags {
	return ev.Flags
}

func (ev CandleLightingEvent) GetEmoji() string {
	return "🕯️"
}

func makeCandleEvent(hd HDate, opts *CalOptions, ev *HolidayEvent) HEvent {
	havdalahTitle := false
	useHavdalahOffset := false
	dow := hd.Weekday()
	if dow == time.Saturday {
		useHavdalahOffset = true
	}
	mask := LIGHT_CANDLES
	if ev != nil {
		mask = ev.Flags
		if dow != time.Friday {
			if (mask & (LIGHT_CANDLES_TZEIS | CHANUKAH_CANDLES)) != 0 {
				useHavdalahOffset = true
			} else if (mask & YOM_TOV_ENDS) != 0 {
				havdalahTitle = true
				useHavdalahOffset = true
			}
		}
	} else if dow == time.Saturday {
		havdalahTitle = true
		mask = LIGHT_CANDLES_TZEIS
	}
	// if offset is 0 or undefined, we'll use tzeit time
	offset := opts.CandleLightingMins
	if useHavdalahOffset {
		offset = opts.HavdalahMins
	}
	year, month, day := hd.Greg()
	gregDate := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	loc := opts.Location
	zmanim := NewZmanim(loc.Latitude, loc.Longitude, gregDate, loc.TimeZoneId)
	var time time.Time
	if offset != 0 {
		time = zmanim.SunsetOffset(offset)
	} else {
		time = zmanim.Tzeit(opts.HavdalahDeg)
	}
	if time == nilTime {
		return HolidayEvent{} // no sunset
	}
	if havdalahTitle {
		return HavdalahEvent{
			Date:        hd,
			Flags:       mask,
			eventTime:   time,
			mins:        opts.HavdalahMins,
			linkedEvent: ev,
			loc:         loc,
		}
	} else {
		return CandleLightingEvent{
			Date:        hd,
			Flags:       mask,
			eventTime:   time,
			linkedEvent: ev,
			loc:         loc,
		}
	}
}
