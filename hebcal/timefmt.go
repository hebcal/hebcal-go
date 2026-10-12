package hebcal

import "github.com/hebcal/hebcal-go/event"

// ReformatTimeStr formats a 24-hour (00:00-23:59) time string such as "20:13"
// either as 12-hour time with a suffix (e.g. "8:13pm") or leaves it in 24-hour
// format, according to opts.
//
// suffix is the afternoon marker, such as "p", "pm" or " P.M."; a morning time
// gets "a", "am" or " A.M." instead.
//
// opts.Hour24 forces 24-hour time and opts.Hour12 forces 12-hour time.
// Otherwise the country of opts.Location decides (see event.Uses12HourClock),
// falling back to Israel when opts.IL is set and to the United States when it
// is not. opts may be nil, which means the United States.
func ReformatTimeStr(timeStr, suffix string, opts *CalOptions) string {
	if opts == nil {
		opts = &CalOptions{}
	}
	return event.ReformatTimeStr(timeStr, suffix, clockCountry(opts), hour12Override(opts))
}

// clockCountry is the country code whose clock convention applies to opts.
func clockCountry(opts *CalOptions) string {
	if opts.Location != nil && opts.Location.CountryCode != "" {
		return opts.Location.CountryCode
	}
	if opts.IL {
		return "IL"
	}
	return "US"
}

// hour12Override turns Hour24 and Hour12 into event.ReformatTimeStr's
// optional override: nil when neither is set.
func hour12Override(opts *CalOptions) *bool {
	var b bool
	switch {
	case opts.Hour24:
		b = false
	case opts.Hour12:
		b = true
	default:
		return nil
	}
	return &b
}
