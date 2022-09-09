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

import (
	"fmt"
	"strings"
	"time"
)

/*
type TimedEvent2 interface {
	GetDate() HDate         // Holiday date of occurrence
	Render() string         // Description (e.g. "Pesach III (CH''M)")
	GetFlags() HolidayFlags // Event flag bitmask
	GetEmoji() string       // Holiday-specific emoji
	GetTime() time.Time
	GetLinkedEvent() HEvent
	GetLocation() *HLocation
}
*/

// TimedEvent is used for Candle-lighting, Havdalah, and fast start/end
type TimedEvent struct {
	HolidayEvent
	eventTime    time.Time
	sunsetOffset int
	loc          *HLocation
	linkedEvent  HEvent
}

func NewTimedEvent(hd HDate, desc string, flags HolidayFlags, t time.Time, sunsetOffset int, linkedEvent HEvent, loc *HLocation) TimedEvent {
	if t == nilTime {
		return TimedEvent{}
	}
	var emoji string
	switch flags {
	case LIGHT_CANDLES, LIGHT_CANDLES_TZEIS:
		emoji = "🕯️"
	case YOM_TOV_ENDS:
		emoji = "✨"
	case CHANUKAH_CANDLES:
		emoji = chanukahEmoji
	}
	return TimedEvent{
		HolidayEvent: HolidayEvent{
			Date:  hd,
			Desc:  desc,
			Flags: flags,
			Emoji: emoji,
		},
		eventTime:    t,
		linkedEvent:  linkedEvent,
		loc:          loc,
		sunsetOffset: sunsetOffset,
	}
}

func (ev TimedEvent) GetDate() HDate {
	return ev.Date
}

func (ev TimedEvent) Render(locale string) string {
	desc := ev.Desc
	if desc == "Havdalah" && ev.sunsetOffset != 0 {
		desc = fmt.Sprintf("Havdalah (%d min)", ev.sunsetOffset)
	}
	timeStr := ev.eventTime.Format(time.Kitchen)
	return fmt.Sprintf("%s: %s", desc, timeStr[0:len(timeStr)-2])
}

func (ev TimedEvent) GetFlags() HolidayFlags {
	return ev.Flags
}

func (ev TimedEvent) GetEmoji() string {
	return ev.Emoji
}

func makeCandleEvent(hd HDate, opts *CalOptions, ev HEvent) TimedEvent {
	havdalahTitle := false
	useHavdalahOffset := false
	dow := hd.Weekday()
	if dow == time.Saturday {
		useHavdalahOffset = true
	}
	flags := LIGHT_CANDLES
	if ev != nil {
		flags = ev.GetFlags()
		if dow != time.Friday {
			if (flags & (LIGHT_CANDLES_TZEIS | CHANUKAH_CANDLES)) != 0 {
				useHavdalahOffset = true
			} else if (flags & YOM_TOV_ENDS) != 0 {
				havdalahTitle = true
				useHavdalahOffset = true
			}
		}
	} else if dow == time.Saturday {
		havdalahTitle = true
		flags = LIGHT_CANDLES_TZEIS
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
	var eventTime time.Time
	if offset != 0 {
		eventTime = zmanim.SunsetOffset(offset)
	} else {
		eventTime = zmanim.Tzeit(opts.HavdalahDeg)
	}
	if eventTime == nilTime {
		return TimedEvent{} // no sunset
	}
	desc := "Candle lighting"
	if havdalahTitle {
		desc = "Havdalah"
	}
	return NewTimedEvent(hd, desc, flags, eventTime, offset, ev, loc)
}

func makeChanukahCandleLighting(ev HolidayEvent, opts *CalOptions) TimedEvent {
	hd := ev.Date
	dow := hd.Weekday()
	if dow == time.Friday || dow == time.Saturday {
		timedEv := makeCandleEvent(hd, opts, ev)
		timedEv.Desc = ev.Desc
		timedEv.ChanukahDay = ev.ChanukahDay
		return timedEv
	}
	loc := opts.Location
	year, month, day := hd.Greg()
	gregDate := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	zmanim := NewZmanim(loc.Latitude, loc.Longitude, gregDate, loc.TimeZoneId)
	candleLightingTime := zmanim.Dusk()
	return TimedEvent{
		HolidayEvent: ev,
		eventTime:    candleLightingTime,
		linkedEvent:  &ev,
		loc:          loc,
	}
}

// Tzais (nightfall) based on the opinion of the Geonim calculated at
// the sun's position at 8.5° below the western horizon.
// https://kosherjava.com/zmanim/docs/api/com/kosherjava/zmanim/ComplexZmanimCalendar.html#getTzaisGeonim8Point5Degrees()
const Tzeit3SmallStars = 8.5

// Tzais (nightfall) based on the opinion of the
// Geonim calculated as 30 minutes after sunset during the equinox
// (on March 16, about 4 days before the astronomical equinox, the day that
// a solar hour is 60 minutes) in Yerushalayim.
// https://kosherjava.com/zmanim/docs/api/com/kosherjava/zmanim/ComplexZmanimCalendar.html#getTzaisGeonim7Point083Degrees()
const Tzeit3MediumStars = 7.083

func makeFastStartEnd(ev HEvent, loc *HLocation) (TimedEvent, TimedEvent) {
	year, month, day := ev.GetDate().Greg()
	gregDate := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	zmanim := NewZmanim(loc.Latitude, loc.Longitude, gregDate, loc.TimeZoneId)
	hd := ev.GetDate()
	desc := ev.Render("en")
	flags := ev.GetFlags()
	var startEvent, endEvent TimedEvent
	if desc == "Erev Tish'a B'Av" {
		sunset := zmanim.Sunset()
		startEvent = NewTimedEvent(hd, "Fast begins", flags, sunset, 0, ev, loc)
	} else if strings.HasPrefix(desc, "Tish'a B'Av") {
		tzeit := zmanim.Tzeit(Tzeit3MediumStars)
		endEvent = NewTimedEvent(hd, "Fast ends", flags, tzeit, 0, ev, loc)
	} else {
		dawn := zmanim.AlotHaShachar()
		startEvent = NewTimedEvent(hd, "Fast begins", flags, dawn, 0, ev, loc)
		if hd.Weekday() != time.Friday && !(hd.Day == 14 && hd.Month == Nisan) {
			tzeit := zmanim.Tzeit(Tzeit3MediumStars)
			endEvent = NewTimedEvent(hd, "Fast ends", flags, tzeit, 0, ev, loc)
		}
	}
	return startEvent, endEvent
}
