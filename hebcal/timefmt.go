package hebcal

import "github.com/hebcal/hebcal-go/event"

// ReformatTimeStr formats a 24-hour (00:00-23:59) time string such as "20:13"
// either as 12-hour time with a suffix (e.g. "8:13pm") or leaves it in 24-hour
// format, according to opts.
//
// suffix is the afternoon marker, such as "p", "pm" or " P.M."; a morning time
// gets "a", "am" or " A.M." instead.
//
// opts.Clock forces 12-hour time (event.Clock12) or 24-hour time
// (event.Clock24). Otherwise the country of opts.Location decides (see event.Uses12HourClock),
// falling back to Israel when opts.IL is set and to the United States when it
// is not. opts may be nil, which means the United States.
func ReformatTimeStr(timeStr, suffix string, opts *CalOptions) string {
	if opts == nil {
		opts = &CalOptions{}
	}
	return event.ReformatTimeStr(timeStr, suffix, clockCountry(opts), opts.Clock)
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
