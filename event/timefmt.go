package event

import (
	"fmt"
	"strconv"
	"strings"
)

// hour12Countries lists the ISO country codes that conventionally use a
// 12-hour clock, matching the hour12cc table in @hebcal/core reformatTimeStr.
var hour12Countries = map[string]bool{
	"US": true, "CA": true, "BR": true, "AU": true, "NZ": true, "DO": true,
	"PR": true, "GR": true, "IN": true, "KR": true, "NP": true, "ZA": true,
}

// Uses12HourClock reports whether the country with ISO 3166 code cc (e.g.
// "US", "GB", "IL") conventionally uses a 12-hour clock.
func Uses12HourClock(cc string) bool {
	return hour12Countries[cc]
}

// ReformatTimeStr formats a 24-hour (00:00-23:59) time string such as "20:13"
// either as 12-hour time with a suffix (e.g. "8:13pm") or leaves it in 24-hour
// format, matching @hebcal/core reformatTimeStr.
//
// suffix is the afternoon marker, such as "p", "pm" or " P.M." (include a
// leading space if you want one); for a morning time its "p" or "P" is
// replaced by "a" or "A", so "pm" becomes "am".
//
// cc is the ISO 3166 country code of the location. When hour12 is nil, the
// 12-hour form is used only for countries that conventionally use a 12-hour
// clock (see Uses12HourClock). A non-nil hour12 overrides the country: true
// always yields 12-hour time and false always yields 24-hour time.
//
// Pass cc "" to mean the United States; a caller with an Israel schedule
// should pass "IL".
func ReformatTimeStr(timeStr, suffix, cc string, hour12 *bool) string {
	if hour12 != nil && !*hour12 {
		return timeStr
	}
	if cc == "" {
		cc = "US"
	}
	if hour12 == nil && !hour12Countries[cc] {
		return timeStr
	}
	hm := strings.SplitN(timeStr, ":", 2)
	if len(hm) != 2 {
		return timeStr
	}
	hour, err := strconv.Atoi(hm[0])
	if err != nil {
		return timeStr
	}
	hourStr := strconv.Itoa(hour)
	if hour < 12 && suffix != "" {
		suffix = strings.Replace(suffix, "p", "a", 1)
		suffix = strings.Replace(suffix, "P", "A", 1)
		if hour == 0 {
			hourStr = "12"
		}
	} else if hour > 12 {
		hourStr = strconv.Itoa(hour % 12)
	} else if hour == 0 {
		hourStr = "00"
	}
	return fmt.Sprintf("%s:%s%s", hourStr, hm[1], suffix)
}
