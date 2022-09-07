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

func makeCandleEvent(hd HDate, loc *HLocation, ev *HolidayEvent) HEvent {
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
	// const offset = useHavdalahOffset ? options.havdalahMins : options.candleLightingMins;
	var offset int = -18
	if useHavdalahOffset {
		offset = 42
	}
	year, month, day := hd.Greg()
	gregDate := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	zmanim := NewZmanim(loc.Latitude, loc.Longitude, gregDate, loc.TimeZoneId)
	var time time.Time
	if offset != 0 {
		time = zmanim.SunsetOffset(offset)
	} else {
		time = zmanim.Tzeit(8.5)
	}
	if time == nilTime {
		return HolidayEvent{} // no sunset
	}
	if havdalahTitle {
		return HavdalahEvent{
			Date:        hd,
			Flags:       mask,
			eventTime:   time,
			mins:        42,
			linkedEvent: ev,
		}
	} else {
		return CandleLightingEvent{
			Date:        hd,
			Flags:       mask,
			eventTime:   time,
			linkedEvent: ev,
		}
	}
}
