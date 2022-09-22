package main

import (
	"fmt"
	"time"

	"os"

	"github.com/hebcal/hebcal-go"
	"github.com/hebcal/hebcal-go/hdate"
	"github.com/hebcal/hebcal-go/locales"
	getopt "github.com/pborman/getopt/v2"
)

func main() {
	opt := getopt.New()
	var (
		help = opt.BoolLong("help", 0, "print this help text")
		/*inFileName*/ _ = opt.StringLong("infile", 'I', "", "Get non-yahrtzeit Hebrew user events from specified file. The format is: mmm dd string, Where mmm is a Hebrew month name", "INFILE")
		/*today_sw*/ _ = opt.BoolLong("today", 't', "Only output for today's date")
		/*noGreg_sw*/ _ = opt.BoolLong("today-brief", 'T', "Print today's pertinent information")
		/*yahrtzeitFileName*/ _ = opt.StringLong("yahrtzeit", 'Y', "", "Get yahrtzeit dates from specified file. The format is: mm dd yyyy string. The first three fields specify a *Gregorian* date.", "YAHRTZEIT")
		/*ashkenazi_sw*/ _ = opt.BoolLong("ashkenazi", 'a', "Use Ashkenazi Hebrew transliterations")
		/*euroDates_sw*/ _ = opt.BoolLong("euro-dates", 'e', "Output 'European' dates -- DD.MM.YYYY")
		/*twentyFourHour_sw*/ _ = opt.BoolLong("24hour", 'E', "Output 24-hour times (e.g. 18:37 instead of 6:37)")
		/*iso8601dates_sw */ _ = opt.BoolLong("iso-8601", 'g', "Output ISO 8601 dates -- YYYY-MM-DD")
		lang                   = opt.StringLong("lang", 0, "en", "Use LANG titles", "LANG")
		/*printMolad_sw*/ _ = opt.BoolLong("molad", 'M', "Print the molad on Shabbat Mevorchim")
		/*printSunriseSunset_sw*/ _ = opt.BoolLong("sunrise-and-sunset", 'O', "Output sunrise and sunset times every day")
		/*tabs_sw*/ _ = opt.BoolLong("tabs", 'r', "Tab delineated format")
		/*sedraAllWeek_sw*/ _ = opt.BoolLong("daily-sedra", 'S', "Print sedrah of the week on all calendar days")
		version_sw            = opt.BoolLong("version", 0, "Show version number")
		/*weekday_sw*/ _ = opt.BoolLong("weekday", 'w', "Add day of the week")
		/*abbrev_sw*/ _ = opt.BoolLong("abbreviated", 'W', "Weekly view. Omer, dafyomi, and non-date-specific zemanim are shown once a week, on the day which corresponds to the first day in the range.")
		/*yearDigits_sw*/ _ = opt.BoolLong("year-abbrev", 'y', "Print only last two digits of year")
		cityNameArg         = opt.StringLong("city", 'C', "New York", "City for candle-lighting", "CITY")
		tzid                = opt.StringLong("timezone", 'z', "America/New_York", "Use specified timezone, overriding the -C (localize to city) switch", "TIMEZONE")
		/*utf8_hebrew_sw*/ _ = 1
		/*latitude*/ _ = 0.0
		/*longitude*/ _ = 0.0
		/*zemanim_sw*/ _ = opt.StringLong("zmanim", 'Z', "Print zemanim (experimental)")
	)

	var latitude float64
	var longitude float64
	opt.FlagLong(&latitude, "latitude", 'l', "Set the latitude for solar calculations", "LATITUDE")
	opt.FlagLong(&longitude, "longitude", 'L', "Set the longitude for solar calculations", "LONGITUDE")

	calOptions := hebcal.CalOptions{}
	opt.FlagLong(&calOptions.CandleLighting,
		"candlelighting", 'c', "Print candlelighting times")
	opt.FlagLong(&calOptions.AddHebrewDates,
		"add-hebrew-dates", 'd', "Print the Hebrew date for the entire date range")
	opt.FlagLong(&calOptions.AddHebrewDatesForEvents, "add-hebrew-dates-for-events", 'D', "Print the Hebrew date for dates with some event")

	opt.FlagLong(&calOptions.IsHebrewYear,
		"hebrew-date", 'H', "Use Hebrew date ranges - only needed when e.g. hebcal -H 5373")

	opt.FlagLong(&calOptions.DafYomi,
		"daf-yomi", 'F', "Output the Daf Yomi for the entire date range")
	opt.FlagLong(&calOptions.NoHolidays,
		"no-holidays", 'h', "Suppress default holidays")
	opt.FlagLong(&calOptions.NoRoshChodesh,
		"no-rosh-chodesh", 'x', "Suppress Rosh Chodesh")

	opt.FlagLong(&calOptions.IL,
		"israeli", 'i', "Israeli holiday and sedra schedule")
	opt.FlagLong(&calOptions.NoModern,
		"no-modern", 0, "Suppress modern holidays")
	opt.FlagLong(&calOptions.Omer,
		"omer", 'o', "Add days of the Omer")
	opt.FlagLong(&calOptions.Sedrot,
		"sedrot", 's', "Add weekly sedrot on Saturday")

	calOptions.CandleLightingMins = 18
	opt.FlagLong(&calOptions.CandleLightingMins,
		"candle-mins", 'b', "Set candle-lighting to occur this many minutes before sundown", "MINUTES")

	opt.FlagLong(&calOptions.HavdalahMins,
		"havdalah-mins", 'm', "Set Havdalah to occur this many minutes after sundown", "MINUTES")
	opt.FlagLong(&calOptions.HavdalahDeg,
		"havdalah-deg", 0, "Set Havdalah to occur this many degrees below the horizon", "DEGREES")

	calOptions.NumYears = 1
	opt.FlagLong(&calOptions.NumYears,
		"years", 0, "Generate events for N years (default 1)")

	if err := opt.Getopt(os.Args, nil); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if *help {
		opt.PrintUsage(os.Stderr)
		os.Exit(0)
	}
	if *version_sw {
		fmt.Println("foo")
		os.Exit(0)
	}

	if lang != nil {
		found := false
		for _, a := range locales.AllLocales {
			if a == *lang {
				found = true
				break
			}
		}
		if !found {
			fmt.Fprintf(os.Stderr, "Unknown lang '%s'; using default\n", *lang)
			*lang = "en"
		}
	}

	if cityNameArg != nil {
		city := hebcal.LookupCity(*cityNameArg)
		if (city == hebcal.HLocation{}) {
			fmt.Fprintf(os.Stderr, "unknown city: %s. Use a nearby city or geographic coordinates.\n", *cityNameArg)
			os.Exit(1)
		}
		calOptions.Location = &city
		calOptions.CandleLighting = true
	}

	if tzid != nil {
		_, err := time.LoadLocation(*tzid)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		calOptions.Location.TimeZoneId = *tzid
	}

	// Get the remaining positional parameters
	args := opt.Args()
	fmt.Println(args)
	events, err := hebcal.HebrewCalendar(&calOptions)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, ev := range events {
		desc := ev.Render(*lang)
		fmt.Printf("%s %s\n", hd2iso(ev.GetDate()), desc)
	}

}

func hd2iso(hd hdate.HDate) string {
	year, month, day := hd.Greg()
	d := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return d.Format(time.RFC3339)[:10]
}
