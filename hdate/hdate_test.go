package hdate

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestElapsedDays(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(2110760, elapsedDays(5780))
	assert.Equal(2084447, elapsedDays(5708))
	assert.Equal(1373677, elapsedDays(3762))
	assert.Equal(1340455, elapsedDays(3671))
	assert.Equal(450344, elapsedDays(1234))
	assert.Equal(44563, elapsedDays(123))
	assert.Equal(356, elapsedDays(2))
	assert.Equal(1, elapsedDays(1))
	assert.Equal(2104174, elapsedDays(5762))
	assert.Equal(2104528, elapsedDays(5763))
	assert.Equal(2104913, elapsedDays(5764))
	assert.Equal(2105268, elapsedDays(5765))
	assert.Equal(2105651, elapsedDays(5766))
}

func TestHebrew2RD(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(733359, HebrewToRD(5769, Cheshvan, 15))
	assert.Equal(711262, HebrewToRD(5708, Iyyar, 6))
	assert.Equal(249, HebrewToRD(3762, Tishrei, 1))
	assert.Equal(72, HebrewToRD(3761, Nisan, 1))
	assert.Equal(1, HebrewToRD(3761, Tevet, 18))
	assert.Equal(0, HebrewToRD(3761, Tevet, 17))
	assert.Equal(-1, HebrewToRD(3761, Tevet, 16))
	assert.Equal(-16, HebrewToRD(3761, Tevet, 1))
	assert.Equal(2278650, HebrewToRD(9999, Elul, 29))
	assert.Equal(731840, HebrewToRD(5765, Tishrei, 1))
	assert.Equal(731957, HebrewToRD(5765, Shvat, 1))
	assert.Equal(731987, HebrewToRD(5765, Adar1, 1))
	assert.Equal(732017, HebrewToRD(5765, Adar2, 1))
	assert.Equal(732038, HebrewToRD(5765, Adar2, 22))
	assert.Equal(732046, HebrewToRD(5765, Nisan, 1))
}

func TestRD2Hebrew(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(HDate{Year: 5769, Month: Cheshvan, Day: 15, abs: 733359}, FromRD(733359))
	assert.Equal(HDate{Year: 5708, Month: Iyyar, Day: 6, abs: 711262}, FromRD(711262))
	assert.Equal(HDate{Year: 3762, Month: Tishrei, Day: 1, abs: 249}, FromRD(249))
	assert.Equal(HDate{Year: 3761, Month: Nisan, Day: 1, abs: 72}, FromRD(72))
	assert.Equal(HDate{Year: 3761, Month: Nisan, Day: 8, abs: 79}, FromRD(79))
	assert.Equal(HDate{Year: 3761, Month: Tevet, Day: 18, abs: 1}, FromRD(1))
	assert.Equal(HDate{Year: 3761, Month: Tevet, Day: 17, abs: 0}, FromRD(0))
	assert.Equal(HDate{Year: 3761, Month: Tevet, Day: 16, abs: -1}, FromRD(-1))
	assert.Equal(HDate{Year: 3761, Month: Tevet, Day: 1, abs: -16}, FromRD(-16))
	assert.Equal(HDate{Year: 3761, Month: Kislev, Day: 30, abs: -17}, FromRD(-17))
	assert.Equal(HDate{Year: 9999, Month: Elul, Day: 29, abs: 2278650}, FromRD(2278650))
	assert.Equal(HDate{Year: 5765, Month: Adar2, Day: 22, abs: 732038}, FromRD(732038))
}

func TestMonthNames(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("Nisan", Nisan.String())
	assert.Equal("Tishrei", Tishrei.String())
	assert.Equal("Sh'vat", Shvat.String())
	assert.Equal("Adar I", Adar1.String())
	assert.Equal("Adar II", Adar2.String())
}

func TestMonthNames2(t *testing.T) {
	assert := assert.New(t)
	hd := New(5782, Adar1, 15)
	assert.Equal("Adar I", hd.MonthName("en"))
	hd = New(5783, Adar1, 15)
	assert.Equal("Adar", hd.MonthName("en"))
}

func TestAdar2ResetToAdar1(t *testing.T) {
	assert := assert.New(t)
	hd := New(5782, Adar1, 15)
	assert.Equal(Adar1, hd.Month)
	hd = New(5782, Adar2, 15)
	assert.Equal(Adar2, hd.Month)
	hd = New(5783, Adar1, 15)
	assert.Equal(Adar1, hd.Month)
	hd = New(5783, Adar2, 15)
	assert.Equal(Adar1, hd.Month)
}

func TestMonthFromName(t *testing.T) {
	assert := assert.New(t)
	monthName, _ := MonthFromName("adar")
	assert.Equal(Adar2, monthName)
	monthName, _ = MonthFromName("Adar I")
	assert.Equal(Adar1, monthName)
	monthName, _ = MonthFromName("Adar II")
	assert.Equal(Adar2, monthName)
	monthName, _ = MonthFromName("Adar 2")
	assert.Equal(Adar2, monthName)
	monthName, _ = MonthFromName("Adar 1")
	assert.Equal(Adar1, monthName)
	monthName, _ = MonthFromName("אדר א")
	assert.Equal(Adar1, monthName)
	monthName, _ = MonthFromName("אדר ב")
	assert.Equal(Adar2, monthName)
}

func TestDaysInHebYear(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(385, DaysInYear(5779))
	assert.Equal(355, DaysInYear(5780))
	assert.Equal(353, DaysInYear(5781))
	assert.Equal(384, DaysInYear(5782))
	assert.Equal(355, DaysInYear(5783))
	assert.Equal(383, DaysInYear(5784))
	assert.Equal(355, DaysInYear(5785))
	assert.Equal(354, DaysInYear(5786))
	assert.Equal(385, DaysInYear(5787))
	assert.Equal(355, DaysInYear(5788))
	assert.Equal(354, DaysInYear(5789))
	assert.Equal(383, DaysInYear(3762))
	assert.Equal(354, DaysInYear(3671))
	assert.Equal(353, DaysInYear(1234))
	assert.Equal(355, DaysInYear(123))
	assert.Equal(355, DaysInYear(2))
	assert.Equal(355, DaysInYear(1))

	assert.Equal(353, DaysInYear(5761))
	assert.Equal(354, DaysInYear(5762))
	assert.Equal(385, DaysInYear(5763))
	assert.Equal(355, DaysInYear(5764))
	assert.Equal(383, DaysInYear(5765))
	assert.Equal(354, DaysInYear(5766))
}

func TestDaysInMonth(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(29, DaysInMonth(Iyyar, 5780))
	assert.Equal(30, DaysInMonth(Sivan, 5780))
	assert.Equal(29, DaysInMonth(Cheshvan, 5782))
	assert.Equal(30, DaysInMonth(Cheshvan, 5783))
	assert.Equal(30, DaysInMonth(Kislev, 5783))
	assert.Equal(29, DaysInMonth(Kislev, 5784))

	assert.Equal(30, DaysInMonth(Tishrei, 5765))
	assert.Equal(29, DaysInMonth(Cheshvan, 5765))
	assert.Equal(29, DaysInMonth(Kislev, 5765))
	assert.Equal(29, DaysInMonth(Tevet, 5765))
}

func TestWeekday(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(time.Thursday, New(5769, Cheshvan, 15).Weekday())
	assert.Equal(time.Saturday, New(5708, Iyyar, 6).Weekday())
	assert.Equal(time.Sunday, New(5708, Iyyar, 7).Weekday())
	assert.Equal(time.Thursday, New(3762, Tishrei, 1).Weekday())
	assert.Equal(time.Tuesday, New(3761, Nisan, 1).Weekday())
	assert.Equal(time.Monday, New(3761, Tevet, 18).Weekday())
	assert.Equal(time.Sunday, New(3761, Tevet, 17).Weekday())
	assert.Equal(time.Saturday, New(3761, Tevet, 16).Weekday())
	assert.Equal(time.Friday, New(3761, Tevet, 1).Weekday())
	assert.Equal(time.Tuesday, New(3333, Sivan, 29).Weekday())
	assert.Equal(time.Monday, New(3333, Sivan, 28).Weekday())
	assert.Equal(time.Sunday, New(3333, Sivan, 27).Weekday())
	assert.Equal(time.Saturday, New(3333, Sivan, 26).Weekday())
	assert.Equal(time.Friday, New(3333, Sivan, 25).Weekday())
	assert.Equal(time.Thursday, New(3333, Sivan, 24).Weekday())
	assert.Equal(time.Wednesday, New(3333, Sivan, 23).Weekday())
}

func hd2iso(hd HDate) string {
	year, month, day := hd.Greg()
	d := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return d.Format(time.RFC3339)[:10]
}

func TestBefore(t *testing.T) {
	assert := assert.New(t)
	d := time.Date(2014, time.February, 19, 0, 0, 0, 0, time.UTC)
	hd := FromTime(d)
	assert.Equal("2014-02-15", hd2iso(hd.Before(time.Saturday)))
}

func TestOnOrBefore(t *testing.T) {
	assert := assert.New(t)
	d := time.Date(2014, time.February, 19, 0, 0, 0, 0, time.UTC)
	assert.Equal("2014-02-15", hd2iso(FromTime(d).OnOrBefore(time.Saturday)))
	d = time.Date(2014, time.February, 22, 0, 0, 0, 0, time.UTC)
	assert.Equal("2014-02-22", hd2iso(FromTime(d).OnOrBefore(time.Saturday)))
	d = time.Date(2014, time.February, 23, 0, 0, 0, 0, time.UTC)
	assert.Equal("2014-02-22", hd2iso(FromTime(d).OnOrBefore(time.Saturday)))

}

func TestNearest(t *testing.T) {
	assert := assert.New(t)
	d := time.Date(2014, time.February, 19, 0, 0, 0, 0, time.UTC)
	assert.Equal("2014-02-22", hd2iso(FromTime(d).Nearest(time.Saturday)))
	d = time.Date(2014, time.February, 18, 0, 0, 0, 0, time.UTC)
	assert.Equal("2014-02-15", hd2iso(FromTime(d).Nearest(time.Saturday)))
}

func TestOnOrAfter(t *testing.T) {
	assert := assert.New(t)
	d := time.Date(2014, time.February, 19, 0, 0, 0, 0, time.UTC)
	assert.Equal("2014-02-22", hd2iso(FromTime(d).OnOrAfter(time.Saturday)))
	d = time.Date(2014, time.February, 22, 0, 0, 0, 0, time.UTC)
	assert.Equal("2014-02-22", hd2iso(FromTime(d).OnOrAfter(time.Saturday)))
	d = time.Date(2014, time.February, 23, 0, 0, 0, 0, time.UTC)
	assert.Equal("2014-03-01", hd2iso(FromTime(d).OnOrAfter(time.Saturday)))
}

func TestAfter(t *testing.T) {
	assert := assert.New(t)
	d := time.Date(2014, time.February, 19, 0, 0, 0, 0, time.UTC)
	assert.Equal("2014-02-22", hd2iso(FromTime(d).After(time.Saturday)))
	d = time.Date(2014, time.February, 22, 0, 0, 0, 0, time.UTC)
	assert.Equal("2014-03-01", hd2iso(FromTime(d).After(time.Saturday)))
	d = time.Date(2014, time.February, 23, 0, 0, 0, 0, time.UTC)
	assert.Equal("2014-03-01", hd2iso(FromTime(d).After(time.Saturday)))
}

func TestToString(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("4 Tevet 5511", New(5511, Tevet, 4).String())
	assert.Equal("4 Elul 5782", New(5782, Elul, 4).String())
	assert.Equal("29 Adar II 5749", New(5749, Adar2, 29).String())
}

func TestGreg(t *testing.T) {
	assert := assert.New(t)
	hd := New(5765, Adar2, 22)
	gy, gm, gd := hd.Greg()
	assert.Equal(2005, gy)
	assert.Equal(time.April, gm)
	assert.Equal(2, gd)
}

func ExampleFromGregorian() {
	hd := FromGregorian(2008, time.November, 13)
	fmt.Println(hd)
	// Output: 15 Cheshvan 5769
}

func ExampleFromTime() {
	t := time.Date(2008, time.November, 13, 0, 0, 0, 0, time.UTC)
	hd := FromTime(t)
	fmt.Println(hd)
	// Output: 15 Cheshvan 5769
}

func ExampleFromRD() {
	hd := FromRD(733359)
	fmt.Println(hd)
	// Output: 15 Cheshvan 5769
}

func ExampleMonthFromName() {
	m1, _ := MonthFromName("Shvat")
	m2, _ := MonthFromName("cheshvan")
	m3, _ := MonthFromName("טבת")
	fmt.Printf("%s (%d), %s (%d), %s (%d)", m1, int(m1), m2, int(m2), m3, int(m3))
	// Output: Sh'vat (11), Cheshvan (8), Tevet (10)
}

func ExampleDaysInYear() {
	days := DaysInYear(5782)
	fmt.Println(days)
	// Output: 384
}

func ExampleDaysInMonth() {
	days := DaysInMonth(Kislev, 5783)
	fmt.Println(days)
	// Output: 30
}

func ExampleHebrewToRD() {
	rataDie := HebrewToRD(5769, Cheshvan, 15)
	fmt.Println(rataDie)
	// Output: 733359
}

func ExampleHDate_Greg() {
	hd := New(5765, Adar2, 22)
	year, month, day := hd.Greg()
	fmt.Println(year, month, day)
	// Output: 2005 April 2
}

func ExampleHDate_Abs() {
	hd := New(5765, Adar2, 22)
	rataDie := hd.Abs()
	fmt.Println(rataDie)
	// Output: 732038
}

func ExampleHDate_DaysInMonth() {
	hd := New(5765, Adar2, 22)
	days := hd.DaysInMonth()
	fmt.Println(days)
	// Output: 29
}

func ExampleHDate_Weekday() {
	hd := New(5765, Adar2, 22)
	dayOfWeek := hd.Weekday()
	fmt.Println(dayOfWeek)
	// Output: Saturday
}

func ExampleHDate_MonthName() {
	hd := New(5765, Adar2, 22)
	fmt.Println(hd.MonthName("en"))
	fmt.Println(hd.MonthName("he"))
	// Output:
	// Adar II
	// אַדָר ב׳
}

func ExampleHDate_IsLeapYear() {
	hd := New(5765, Adar2, 22)
	leap := hd.IsLeapYear()
	fmt.Println(leap)
	// Output: true
}

func ExampleHDate_Gregorian() {
	hd := New(5765, Adar2, 22)
	t := hd.Gregorian()
	fmt.Println(t)
	// Output: 2005-04-02 00:00:00 +0000 UTC
}

func ExampleHDate_Before() {
	orig := FromGregorian(2014, time.February, 19)
	hd := orig.Before(time.Saturday)
	fmt.Println(hd.Gregorian().Format(time.RFC1123))
	// Output: Sat, 15 Feb 2014 00:00:00 UTC
}

func ExampleHDate_OnOrBefore() {
	orig := FromGregorian(2014, time.February, 19)
	hd := orig.OnOrBefore(time.Saturday)
	fmt.Println(hd.Gregorian().Format(time.RFC1123))
	orig = FromGregorian(2014, time.February, 22)
	hd = orig.OnOrBefore(time.Saturday)
	fmt.Println(hd.Gregorian().Format(time.RFC1123))
	orig = FromGregorian(2014, time.February, 23)
	hd = orig.OnOrBefore(time.Saturday)
	fmt.Println(hd.Gregorian().Format(time.RFC1123))
	// Output:
	// Sat, 15 Feb 2014 00:00:00 UTC
	// Sat, 22 Feb 2014 00:00:00 UTC
	// Sat, 22 Feb 2014 00:00:00 UTC
}

func ExampleHDate_Nearest() {
	orig := FromGregorian(2014, time.February, 19)
	hd := orig.Nearest(time.Saturday)
	fmt.Println(hd.Gregorian().Format(time.RFC1123))
	orig = FromGregorian(2014, time.February, 18)
	hd = orig.Nearest(time.Saturday)
	fmt.Println(hd.Gregorian().Format(time.RFC1123))
	// Output:
	// Sat, 22 Feb 2014 00:00:00 UTC
	// Sat, 15 Feb 2014 00:00:00 UTC
}

func ExampleHDate_OnOrAfter() {
	orig := FromGregorian(2014, time.February, 19)
	hd := orig.OnOrAfter(time.Saturday)
	fmt.Println(hd.Gregorian().Format(time.RFC1123))
	orig = FromGregorian(2014, time.February, 22)
	hd = orig.OnOrAfter(time.Saturday)
	fmt.Println(hd.Gregorian().Format(time.RFC1123))
	orig = FromGregorian(2014, time.February, 23)
	hd = orig.OnOrAfter(time.Saturday)
	fmt.Println(hd.Gregorian().Format(time.RFC1123))
	// Output:
	// Sat, 22 Feb 2014 00:00:00 UTC
	// Sat, 22 Feb 2014 00:00:00 UTC
	// Sat, 01 Mar 2014 00:00:00 UTC
}

func ExampleHDate_After() {
	orig := FromGregorian(2014, time.February, 19)
	hd := orig.After(time.Saturday)
	fmt.Println(hd.Gregorian().Format(time.RFC1123))
	orig = FromGregorian(2014, time.February, 22)
	hd = orig.After(time.Saturday)
	fmt.Println(hd.Gregorian().Format(time.RFC1123))
	orig = FromGregorian(2014, time.February, 23)
	hd = orig.After(time.Saturday)
	fmt.Println(hd.Gregorian().Format(time.RFC1123))
	// Output:
	// Sat, 22 Feb 2014 00:00:00 UTC
	// Sat, 01 Mar 2014 00:00:00 UTC
	// Sat, 01 Mar 2014 00:00:00 UTC
}

func ExampleNew() {
	hd := New(5782, Nisan, 30)
	fmt.Println(hd)
	// Output: 30 Nisan 5782
}

func ExampleHDate_Next() {
	orig := New(5782, Nisan, 30)
	hd := orig.Next()
	fmt.Println(hd)
	// Output: 1 Iyyar 5782
}

func ExampleHDate_Prev() {
	orig := New(5782, Tishrei, 1)
	hd := orig.Prev()
	fmt.Println(hd)
	// Output: 29 Elul 5781
}

func ExampleHDate_String() {
	hd := New(5783, Elul, 29)
	str := hd.String()
	fmt.Println(str)
	// Output: 29 Elul 5783
}

func ExampleHMonth_String() {
	month := Adar2
	fmt.Println(month.String())
	// Output: Adar II
}
