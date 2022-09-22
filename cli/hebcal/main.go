package main

import (
	"fmt"
	"strconv"
	"time"

	"os"

	"github.com/hebcal/hebcal-go"
	"github.com/hebcal/hebcal-go/greg"
	"github.com/hebcal/hebcal-go/hdate"
	"github.com/hebcal/hebcal-go/locales"
	getopt "github.com/pborman/getopt/v2"
)

type RangeType int

const (
	YEAR RangeType = 0 + iota
	MONTH
	DAY
	TODAY
)

var defaultCity = "New York"
var defaultLocation = hebcal.HLocation{}
var userLocation = hebcal.HLocation{}
var calOptions hebcal.CalOptions = hebcal.CalOptions{
	Location: &defaultLocation,
}
var lang = "en"
var theYear = 0
var theGregMonth time.Month = 0
var theHebMonth hdate.HMonth = 0
var theDay = 0
var yearDirty = false
var rangeType = YEAR

func handleArgs() {
	opt := getopt.New()
	var (
		help = opt.BoolLong("help", 0, "print this help text")
		/*inFileName*/ _ = opt.StringLong("infile", 'I', "", "Get non-yahrtzeit Hebrew user events from specified file. The format is: mmm dd string, Where mmm is a Hebrew month name", "INFILE")
		today_sw         = opt.BoolLong("today", 't', "Only output for today's date")
		noGreg_sw        = opt.BoolLong("today-brief", 'T', "Print today's pertinent information")
		/*yahrtzeitFileName*/ _ = opt.StringLong("yahrtzeit", 'Y', "", "Get yahrtzeit dates from specified file. The format is: mm dd yyyy string. The first three fields specify a *Gregorian* date.", "YAHRTZEIT")
		ashkenazi_sw            = opt.BoolLong("ashkenazi", 'a', "Use Ashkenazi Hebrew transliterations")
		/*euroDates_sw*/ _ = opt.BoolLong("euro-dates", 'e', "Output 'European' dates -- DD.MM.YYYY")
		/*twentyFourHour_sw*/ _ = opt.BoolLong("24hour", 'E', "Output 24-hour times (e.g. 18:37 instead of 6:37)")
		/*iso8601dates_sw */ _ = opt.BoolLong("iso-8601", 'g', "Output ISO 8601 dates -- YYYY-MM-DD")
		/*printMolad_sw*/ _ = opt.BoolLong("molad", 'M', "Print the molad on Shabbat Mevorchim")
		/*printSunriseSunset_sw*/ _ = opt.BoolLong("sunrise-and-sunset", 'O', "Output sunrise and sunset times every day")
		/*tabs_sw*/ _ = opt.BoolLong("tabs", 'r', "Tab delineated format")
		/*sedraAllWeek_sw*/ _ = opt.BoolLong("daily-sedra", 'S', "Print sedrah of the week on all calendar days")
		version_sw            = opt.BoolLong("version", 0, "Show version number")
		/*weekday_sw*/ _ = opt.BoolLong("weekday", 'w', "Add day of the week")
		/*abbrev_sw*/ _ = opt.BoolLong("abbreviated", 'W', "Weekly view. Omer, dafyomi, and non-date-specific zemanim are shown once a week, on the day which corresponds to the first day in the range.")
		/*yearDigits_sw*/ _ = opt.BoolLong("year-abbrev", 'y', "Print only last two digits of year")
		cityNameArg         = opt.StringLong("city", 'C', "", "City for candle-lighting", "CITY")
		latitudeStr         = opt.StringLong("latitude", 'l', "", "Set the latitude for solar calculations", "LATITUDE")
		longitudeStr        = opt.StringLong("longitude", 'L', "", "Set the longitude for solar calculations", "LONGITUDE")
		tzid                = opt.StringLong("timezone", 'z', "", "Use specified timezone, overriding the -C (localize to city) switch", "TIMEZONE")
		utf8_hebrew_sw      = opt.BoolLong("", '8', "Use UTF-8 Hebrew")
		/*zemanim_sw*/ _ = opt.StringLong("zmanim", 'Z', "Print zemanim (experimental)")
	)

	opt.FlagLong(&lang, "lang", 0, "Use LANG titles", "LANG")

	opt.FlagLong(&calOptions.CandleLighting,
		"candlelighting", 'c', "Print candlelighting times")
	opt.FlagLong(&calOptions.AddHebrewDates,
		"add-hebrew-dates", 'd', "Print the Hebrew date for the entire date range")
	opt.FlagLong(&calOptions.AddHebrewDatesForEvents, "add-hebrew-dates-for-events", 'D', "Print the Hebrew date for dates with some event")

	opt.FlagLong(&calOptions.IsHebrewYear,
		"hebrew-date", 'H', "Use Hebrew date ranges - only needed when e.g. hebcal -H 5373")

	opt.FlagLong(&calOptions.DafYomi,
		"daf-yomi", 'F', "Output the Daf Yomi for the entire date range")
	opt.FlagLong(&calOptions.MishnaYomi,
		"mishna-yomi", 0, "Output the Mishna Yomi for the entire date range")

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

	if *ashkenazi_sw && *utf8_hebrew_sw {
		fmt.Fprintf(os.Stderr, "Cannot specify both options -a and -8\n")
		os.Exit(1)
	} else if *ashkenazi_sw {
		lang = "ashkenazi"
	} else if *utf8_hebrew_sw {
		lang = "he"
	}
	checkLang()

	if calOptions.CandleLighting && (cityNameArg == nil || *cityNameArg == "") {
		cityNameArg = &defaultCity
	}

	if cityNameArg != nil && *cityNameArg != "" {
		city := hebcal.LookupCity(*cityNameArg)
		if city == defaultLocation {
			fmt.Fprintf(os.Stderr, "unknown city: %s. Use a nearby city or geographic coordinates.\n", *cityNameArg)
			os.Exit(1)
		}
		calOptions.Location = &city
		calOptions.CandleLighting = true
	}

	hasLat := false
	if latitudeStr != nil && *latitudeStr != "" {
		latdeg := 0
		latmin := 0
		n, err := fmt.Sscanf(*latitudeStr, "%d,%d", &latdeg, &latmin)
		if err != nil || n != 2 {
			fmt.Fprintf(os.Stderr, "unable to read latitude argument: %s\n", *latitudeStr)
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		if (intAbs(latdeg) > 90) || latmin > 60 || latmin < 0 {
			fmt.Fprintf(os.Stderr, "Error, latitude argument out of range: %s\n", *latitudeStr)
			os.Exit(1)
		}
		latmin = intAbs(latmin)
		if latdeg < 0 {
			latmin = -latmin
		}
		userLocation.Latitude = float64(latdeg) + (float64(latmin) / 60.0)
		hasLat = true
	}

	hasLong := false
	if longitudeStr != nil && *longitudeStr != "" {
		longdeg := 0
		longmin := 0
		n, err := fmt.Sscanf(*longitudeStr, "%d,%d", &longdeg, &longmin)
		if err != nil || n != 2 {
			fmt.Fprintf(os.Stderr, "unable to read longitude argument: %s\n", *longitudeStr)
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		if (intAbs(longdeg) > 180) || longmin > 60 || longmin < 0 {
			fmt.Fprintf(os.Stderr, "Error, longitude argument out of range: %s\n", *longitudeStr)
			os.Exit(1)
		}
		longmin = intAbs(longmin)
		if longdeg < 0 {
			longmin = -longmin
		}
		userLocation.Longitude = float64(-1*longdeg) + (float64(longmin) / -60.0)
		hasLong = true
	}

	if hasLat && hasLong {
		if tzid == nil || *tzid == "" {
			fmt.Fprintf(os.Stderr, "Error, latitude and longitude requires -z/--timezone\n")
			os.Exit(1)
		}
		_, err := time.LoadLocation(*tzid)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		userLocation.TimeZoneId = *tzid
		if calOptions.IL {
			userLocation.CountryCode = "IL"
		}
		calOptions.Location = &userLocation
		calOptions.CandleLighting = true
	}

	if calOptions.CandleLighting && calOptions.HavdalahDeg == 0.0 && calOptions.HavdalahMins == 0 {
		calOptions.HavdalahMins = 72
	}

	if *noGreg_sw {
		*today_sw = true
	}

	gregTodayYY, gregTodayMM, gregTodayDD := time.Now().Date()

	if *today_sw {
		calOptions.AddHebrewDates = true
		rangeType = TODAY
		theGregMonth = gregTodayMM /* year and month specified */
		theDay = gregTodayDD       /* printc theDay of theMonth */
		yearDirty = true
		calOptions.Omer = true
		calOptions.IsHebrewYear = false
	}

	// Get the remaining positional parameters
	args := opt.Args()

	switch len(args) {
	case 0:
		if calOptions.IsHebrewYear {
			hd := hdate.FromGregorian(gregTodayYY, gregTodayMM, gregTodayDD)
			theYear = hd.Year
		} else {
			theYear = gregTodayYY
		}
	case 1:
		yy, err := strconv.Atoi(args[0])
		if err == nil {
			theYear = yy     /* just year specified */
			yearDirty = true /* print whole year */
		} else {
			switch args[0] {
			case "help":
				opt.PrintUsage(os.Stderr)
				os.Exit(0)
			case "info":
				fmt.Println("info - To Be Implemented")
				os.Exit(0)
			case "cities":
				fmt.Println("cities - To Be Implemented")
				os.Exit(0)
			case "copying":
				fmt.Println("copying - To Be Implemented")
				os.Exit(0)
			case "warranty":
				fmt.Println("warranty - To Be Implemented")
				os.Exit(0)
			default:
				fmt.Fprintf(os.Stderr, "unrecognized command '%s'\n", args[0])
				os.Exit(1)
			}
		}
	case 2:
		yy, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		theYear = yy
		parseGregOrHebMonth(&calOptions, theYear, args[0], &theGregMonth, &theHebMonth)
		yearDirty = true
		rangeType = MONTH
	case 3:
		dd, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		theDay = dd
		yy, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		theYear = yy
		parseGregOrHebMonth(&calOptions, theYear, args[0], &theGregMonth, &theHebMonth)
		yearDirty = true
		rangeType = DAY
	default:
		opt.PrintUsage(os.Stderr)
		os.Exit(1)
	}

	if calOptions.NumYears != 1 && rangeType != YEAR {
		fmt.Fprintf(os.Stderr, "Sorry, --years option works only with entire-year calendars")
		os.Exit(1)
	}
}

func checkLang() {
	if lang != "en" {
		found := false
		for _, a := range locales.AllLocales {
			if a == lang {
				found = true
				break
			}
		}
		if !found {
			fmt.Fprintf(os.Stderr, "Unknown lang '%s'; using default\n", lang)
			lang = "en"
		}
	}
}

func parseGregOrHebMonth(calOptions *hebcal.CalOptions, theYear int, arg string, gregMonth *time.Month, hebMonth *hdate.HMonth) {
	mm, err := strconv.Atoi(arg)
	if err == nil {
		if calOptions.IsHebrewYear {
			fmt.Fprintf(os.Stderr, "Don't use numbers to specify Hebrew months.\n")
			os.Exit(1)
		}
		*gregMonth = time.Month(mm) /* gregorian month */
	} else {
		hm, err := hdate.MonthFromName(arg)
		if err == nil {
			*hebMonth = hm
			calOptions.IsHebrewYear = true /* automagically turn it on */
			if hm == hdate.Adar2 && !hdate.IsLeapYear(theYear) {
				*hebMonth = hdate.Adar1 /* silently fix this mistake */
			}
		} else {
			fmt.Fprintf(os.Stderr, "Unknown Hebrew month: %s.\n", arg)
			os.Exit(1)
		}
	}
}

func main() {
	handleArgs()
	if theYear < 1 || (calOptions.IsHebrewYear && theYear < 3761) {
		fmt.Fprintf(os.Stderr, "Sorry, hebcal can only handle dates in the common era.\n")
		os.Exit(1)
	}
	switch rangeType {
	case TODAY:
		calOptions.AddHebrewDates = true
		calOptions.Start = hdate.FromGregorian(theYear, theGregMonth, theDay)
		calOptions.End = calOptions.Start
	case DAY:
		calOptions.AddHebrewDates = true
		if calOptions.IsHebrewYear {
			calOptions.Start = hdate.New(theYear, theHebMonth, theDay)
		} else {
			calOptions.Start = hdate.FromGregorian(theYear, theGregMonth, theDay)
		}
		calOptions.End = calOptions.Start
	case MONTH:
		if calOptions.IsHebrewYear {
			calOptions.Start = hdate.New(theYear, theHebMonth, 1)
			calOptions.End = hdate.New(theYear, theHebMonth, calOptions.Start.DaysInMonth())
		} else {
			calOptions.Start = hdate.FromGregorian(theYear, theGregMonth, 1)
			calOptions.End = hdate.FromGregorian(theYear, theGregMonth, greg.DaysIn(theGregMonth, theYear))
		}
	case YEAR:
		calOptions.Year = theYear
	default:
		panic("Oh, NO! internal error #17q!")
	}

	events, err := hebcal.HebrewCalendar(&calOptions)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, ev := range events {
		desc := ev.Render(lang)
		fmt.Printf("%s %s\n", hd2iso(ev.GetDate()), desc)
	}
}

func hd2iso(hd hdate.HDate) string {
	year, month, day := hd.Greg()
	d := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return d.Format(time.RFC3339)[:10]
}

func intAbs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
