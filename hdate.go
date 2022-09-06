package hebcal

import (
	"errors"
	"math"
	"regexp"
	"strconv"
	s "strings"
	"time"
)

type HMonth int

const (
	/** Nissan / ניסן */
	Nisan HMonth = 1 + iota
	/** Iyyar / אייר */
	Iyyar
	/** Sivan / סיון */
	Sivan
	/** Tamuz (sometimes Tammuz) / תמוז */
	Tamuz
	/** Av / אב */
	Av
	/** Elul / אלול */
	Elul
	/** Tishrei / תִשְׁרֵי */
	Tishrei
	/** Cheshvan / חשון */
	Cheshvan
	/** Kislev / כסלו */
	Kislev
	/** Tevet / טבת */
	Tevet
	/** Sh'vat / שבט */
	Shvat
	/** Adar or Adar Rishon / אדר */
	Adar1
	/** Adar Sheini (only on leap years) / אדר ב׳ */
	Adar2
)

var longMonthNames = []string{
	"",
	"Nisan",
	"Iyyar",
	"Sivan",
	"Tamuz",
	"Av",
	"Elul",
	"Tishrei",
	"Cheshvan",
	"Kislev",
	"Tevet",
	"Sh'vat",
	"Adar I",
	"Adar II",
	"Nisan",
}

// String returns the English name of the month ("Nisan", "Iyyar", ...).
func (m HMonth) String() string {
	if Nisan <= m && m <= Adar2 {
		return longMonthNames[m]
	}
	return "%!HMonth(" + strconv.Itoa(int(m)) + ")"
}

type HDate struct {
	Year  int
	Month HMonth
	Day   int
	abs   int
}

const epoch = -1373428

// Avg year length in the cycle (19 solar years with 235 lunar months)
const avgHebrewYearDays = 365.24682220597794

/**
 * Returns true if Hebrew year is a leap year
 * @param {number} year Hebrew year
 * @return {boolean}
 */
func IsHebLeapYear(year int) bool {
	return (1+year*7)%19 < 7
}

/**
 * Number of months in this Hebrew year (either 12 or 13 depending on leap year)
 * @param {number} year Hebrew year
 * @return {number}
 */
func MonthsInHebYear(year int) int {
	if IsHebLeapYear(year) {
		return 13
	} else {
		return 12
	}
}

/**
 * Number of days in the hebrew YEAR
 * @param {number} year Hebrew year
 * @return {number}
 */
func DaysInHebYear(year int) int {
	return elapsedDays(year+1) - elapsedDays(year)
}

/**
 * true if Cheshvan is long in Hebrew year
 * @param {number} year Hebrew year
 * @return {boolean}
 */
func LongCheshvan(year int) bool {
	return DaysInHebYear(year)%10 == 5
}

/**
 * true if Kislev is short in Hebrew year
 * @param {number} year Hebrew year
 * @return {boolean}
 */
func ShortKislev(year int) bool {
	return DaysInHebYear(year)%10 == 3
}

/**
 * Number of days in Hebrew month in a given year (29 or 30)
 * @param {number} month Hebrew month (e.g. months.TISHREI)
 * @param {number} year Hebrew year
 * @return {number}
 */
func DaysInMonth(month HMonth, year int) int {
	switch month {
	case Iyyar, Tamuz, Elul, Tevet, Adar2:
		return 29
	}
	if (month == Adar1 && !IsHebLeapYear(year)) ||
		(month == Cheshvan && !LongCheshvan(year)) ||
		(month == Kislev) && ShortKislev(year) {
		return 29
	} else {
		return 30
	}
}

var edCache map[int]int = make(map[int]int)

func elapsedDays(year int) int {
	days, ok := edCache[year]
	if ok {
		return days
	}
	days = elapsedDays0(year)
	edCache[year] = days
	return days
}

/**
 * Days from sunday prior to start of Hebrew calendar to mean
 * conjunction of Tishrei in Hebrew YEAR
 * @private
 * @param {number} year Hebrew year
 * @return {number}
 */
func elapsedDays0(year int) int {
	prevYear := year - 1
	mElapsed := 235*(prevYear/19) + // Months in complete 19 year lunar (Metonic) cycles so far
		12*(prevYear%19) + // Regular months in this cycle
		(((prevYear%19)*7 + 1) / 19) // Leap months this cycle

	pElapsed := 204 + 793*(mElapsed%1080)

	hElapsed := 5 +
		12*mElapsed +
		793*(mElapsed/1080) +
		(pElapsed / 1080)

	parts := (pElapsed % 1080) + 1080*(hElapsed%24)

	day := 1 + 29*mElapsed + (hElapsed / 24)

	altDay := day

	if (parts >= 19440) ||
		(((day % 7) == 2) && (parts >= 9924) && !(IsHebLeapYear(year))) ||
		(((day % 7) == 1) && (parts >= 16789) && IsHebLeapYear(prevYear)) {
		altDay = day + 1
	}

	if altDay%7 == 0 || altDay%7 == 3 || altDay%7 == 5 {
		return altDay + 1
	} else {
		return altDay
	}
}

func HebrewToRD(year int, month HMonth, day int) int {
	tempabs := day
	if month < Tishrei {
		monthsInYear := HMonth(MonthsInHebYear(year))
		for m := Tishrei; m <= monthsInYear; m++ {
			tempabs += DaysInMonth(m, year)
		}
		for m := Nisan; m < month; m++ {
			tempabs += DaysInMonth(m, year)
		}
	} else {
		for m := Tishrei; m < month; m++ {
			tempabs += DaysInMonth(m, year)
		}
	}
	return epoch + elapsedDays(year) + tempabs - 1
}

func NewHDate(year int, month HMonth, day int) HDate {
	if month == Adar2 && !IsHebLeapYear(year) {
		month = Adar1
	}
	if month == Adar2+1 {
		month = Nisan
	}
	if month < Nisan || month > Adar2 {
		panic("invalid Hebrew Month " + month.String())
	}
	if day < 1 || day > 30 {
		panic("invalid Hebrew day " + strconv.Itoa(day))
	}
	return HDate{Year: year, Month: month, Day: day}
}

func newYear(year int) int {
	return epoch + elapsedDays(year)
}

func NewHDateFromRD(rataDie int) HDate {
	approx := int(float64(rataDie-epoch) / avgHebrewYearDays)
	year := approx
	for newYear(year) <= rataDie {
		year++
	}
	year--
	var month HMonth
	if rataDie < HebrewToRD(year, Nisan, 1) {
		month = Tishrei
	} else {
		month = Nisan
	}
	for rataDie > HebrewToRD(year, month, DaysInMonth(month, year)) {
		month++
	}
	day := 1 + rataDie - HebrewToRD(year, month, 1)
	return HDate{Year: year, Month: month, Day: day, abs: rataDie}
}

func NewHDateFromGregorian(year int, month time.Month, day int) HDate {
	rataDie, _ := GregorianToRD(year, month, day)
	return NewHDateFromRD(rataDie)
}

func NewHDateFromTime(d time.Time) HDate {
	year, month, day := d.Date()
	rataDie, _ := GregorianToRD(year, month, day)
	return NewHDateFromRD(rataDie)
}

func (hd *HDate) Abs() int {
	if hd.abs == 0 {
		hd.abs = HebrewToRD(hd.Year, hd.Month, hd.Day)
	}
	return hd.abs
}

func (hd HDate) DaysInMonth() int {
	return DaysInMonth(hd.Month, hd.Year)
}

func (hd HDate) Greg() (int, time.Month, int) {
	return RDtoGregorian(hd.Abs())
}

func mod(x, y int) int {
	return x - y*int(math.Floor(float64(x)/float64(y)))
}

// Weekday returns the day of the week specified by hd.
func (hd HDate) Weekday() time.Weekday {
	abs := hd.Abs()
	if abs < 0 {
		dayOfWeek := mod(abs, 7)
		return time.Weekday(dayOfWeek)
	}
	dayOfWeek := abs % 7
	return time.Weekday(dayOfWeek)
}

func (hd HDate) Prev() HDate {
	return NewHDateFromRD(hd.Abs() - 1)
}

func (hd HDate) Next() HDate {
	return NewHDateFromRD(hd.Abs() + 1)
}

func (hd HDate) IsLeapYear() bool {
	return IsHebLeapYear(hd.Year)
}

func (hd HDate) MonthName() string {
	switch hd.Month {
	case Nisan, Iyyar, Sivan, Tamuz, Av, Elul, Tishrei, Cheshvan, Kislev, Tevet, Shvat, Adar2:
		return hd.Month.String()
	case Adar1:
		if hd.IsLeapYear() {
			return "Adar I"
		} else {
			return "Adar"
		}
	default:
		/*NOTREACHED*/
		return ""
	}
}

func (hd HDate) String() string {
	return strconv.Itoa(hd.Day) + " " + hd.MonthName() + " " + strconv.Itoa(hd.Year)
}

/**
 * Converts Hebrew month string name to numeric
 * @param {string} monthName monthName
 * @return {HMonth}
 */
func MonthFromName(monthName string) (HMonth, error) {
	str := s.ToLower(monthName)
	runes := []rune(str)
	/*
		the Hebrew months are unique to their second letter
		N         Nisan  (November?)
		I         Iyyar
		E        Elul
		C        Cheshvan
		K        Kislev
		1        1Adar
		2        2Adar
		Si Sh     Sivan, Shvat
		Ta Ti Te Tamuz, Tishrei, Tevet
		Av Ad    Av, Adar

		אב אד אי אל   אב אדר אייר אלול
		ח            חשון
		ט            טבת
		כ            כסלו
		נ            ניסן
		ס            סיון
		ש            שבט
		תמ תש        תמוז תשרי
	*/
	switch runes[0] {
	case 'n':
	case 'נ':
		if runes[1] == 'o' {
			break /* this catches "november" */
		}
		return Nisan, nil
	case 'i':
		return Iyyar, nil
	case 'e':
		return Elul, nil
	case 'c':
	case 'ח':
		return Cheshvan, nil
	case 'k':
	case 'כ':
		return Kislev, nil
	case 's':
		switch runes[1] {
		case 'i':
			return Sivan, nil
		case 'h':
			return Shvat, nil
		default:
			break
		}
	case 't':
		switch runes[1] {
		case 'a':
			return Tamuz, nil
		case 'i':
			return Tishrei, nil
		case 'e':
			return Tevet, nil
		}
	case 'a':
		switch runes[1] {
		case 'v':
			return Av, nil
		case 'd':
			regex := regexp.MustCompile("(?i)(1|[^i]i|a|א)$")
			if regex.MatchString(monthName) {
				return Adar1, nil
			}
			return Adar2, nil // else assume sheini
		}
	case 'ס':
		return Sivan, nil
	case 'ט':
		return Tevet, nil
	case 'ש':
		return Shvat, nil
	case 'א':
		switch runes[1] {
		case 'ב':
			return Av, nil
		case 'ד':
			regex := regexp.MustCompile("(?i)(1|[^i]i|a|א)$")
			if regex.MatchString(monthName) {
				return Adar1, nil
			}
			return Adar2, nil // else assume sheini
		case 'י':
			return Iyyar, nil
		case 'ל':
			return Elul, nil
		}
	case 'ת':
		switch runes[1] {
		case 'מ':
			return Tamuz, nil
		case 'ש':
			return Tishrei, nil
		}
	}
	return 0, errors.New("unable to parse month name")
}

/**
 * Note: Applying this function to d+6 gives us the DAYNAME on or after an
 * absolute day d. Similarly, applying it to d+3 gives the DAYNAME nearest to
 * absolute date d, applying it to d-1 gives the DAYNAME previous to absolute
 * date d, and applying it to d+7 gives the DAYNAME following absolute date d.
 * @param {time.Weekday} dayOfWeek
 * @param {int} rataDie
 * @return {int}
 */
func dayOnOrBefore(dayOfWeek time.Weekday, rataDie int) int {
	return rataDie - ((rataDie - int(dayOfWeek)) % 7)
}

func onOrBefore(dayOfWeek time.Weekday, rataDie int) HDate {
	return NewHDateFromRD(dayOnOrBefore(dayOfWeek, rataDie))
}

/**
 * Returns an `HDate` representing the a dayNumber before the current date.
 * Sunday=0, Saturday=6
 * @example
 * new HDate(new Date('Wednesday February 19, 2014')).before(6).greg() // Sat Feb 15 2014
 * @param {number} day day of week
 * @return {HDate}
 */
func (hd HDate) Before(dayOfWeek time.Weekday) HDate {
	return onOrBefore(dayOfWeek, hd.Abs()-1)
}

/**
 * Returns an `HDate` representing the a dayNumber on or before the current date.
 * Sunday=0, Saturday=6
 * @example
 * new HDate(new Date('Wednesday February 19, 2014')).onOrBefore(6).greg() // Sat Feb 15 2014
 * new HDate(new Date('Saturday February 22, 2014')).onOrBefore(6).greg() // Sat Feb 22 2014
 * new HDate(new Date('Sunday February 23, 2014')).onOrBefore(6).greg() // Sat Feb 22 2014
 * @param {number} dow day of week
 * @return {HDate}
 */
func (hd HDate) OnOrBefore(dayOfWeek time.Weekday) HDate {
	return onOrBefore(dayOfWeek, hd.Abs())
}

/**
 * Returns an `HDate` representing the nearest dayNumber to the current date
 * Sunday=0, Saturday=6
 * @example
 * new HDate(new Date('Wednesday February 19, 2014')).nearest(6).greg() // Sat Feb 22 2014
 * new HDate(new Date('Tuesday February 18, 2014')).nearest(6).greg() // Sat Feb 15 2014
 * @param {number} dow day of week
 * @return {HDate}
 */
func (hd HDate) Nearest(dayOfWeek time.Weekday) HDate {
	return onOrBefore(dayOfWeek, hd.Abs()+3)
}

/**
 * Returns an `HDate` representing the a dayNumber on or after the current date.
 * Sunday=0, Saturday=6
 * @example
 * new HDate(new Date('Wednesday February 19, 2014')).onOrAfter(6).greg() // Sat Feb 22 2014
 * new HDate(new Date('Saturday February 22, 2014')).onOrAfter(6).greg() // Sat Feb 22 2014
 * new HDate(new Date('Sunday February 23, 2014')).onOrAfter(6).greg() // Sat Mar 01 2014
 * @param {number} dow day of week
 * @return {HDate}
 */
func (hd HDate) OnOrAfter(dayOfWeek time.Weekday) HDate {
	return onOrBefore(dayOfWeek, hd.Abs()+6)
}

/**
 * Returns an `HDate` representing the a dayNumber after the current date.
 * Sunday=0, Saturday=6
 * @example
 * new HDate(new Date('Wednesday February 19, 2014')).after(6).greg() // Sat Feb 22 2014
 * new HDate(new Date('Saturday February 22, 2014')).after(6).greg() // Sat Mar 01 2014
 * new HDate(new Date('Sunday February 23, 2014')).after(6).greg() // Sat Mar 01 2014
 * @param {number} day day of week
 * @return {HDate}
 */
func (hd HDate) After(dayOfWeek time.Weekday) HDate {
	return onOrBefore(dayOfWeek, hd.Abs()+7)
}
