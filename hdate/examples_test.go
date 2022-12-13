package hdate_test

import (
	"fmt"
	"time"

	"github.com/hebcal/hebcal-go/hdate"
)

/****** HDate *****/
func ExampleFromGregorian() {
	hd := hdate.FromGregorian(2008, time.November, 13)
	fmt.Println(hd)
	// Output: 15 Cheshvan 5769
}

func ExampleFromTime() {
	t := time.Date(2008, time.November, 13, 0, 0, 0, 0, time.UTC)
	hd := hdate.FromTime(t)
	fmt.Println(hd)
	// Output: 15 Cheshvan 5769
}

func ExampleFromRD() {
	hd := hdate.FromRD(733359)
	fmt.Println(hd)
	// Output: 15 Cheshvan 5769
}

func ExampleMonthFromName() {
	var month hdate.HMonth
	month, _ = hdate.MonthFromName("Shvat")
	fmt.Printf("%s (%d)\n", month, int(month))
	month, _ = hdate.MonthFromName("cheshvan")
	fmt.Printf("%s (%d)\n", month, int(month))
	month, _ = hdate.MonthFromName("טבת")
	fmt.Printf("%s (%d)\n", month, int(month))
	month, _ = hdate.MonthFromName("Adar 1")
	fmt.Printf("%s (%d)\n", month, int(month))
	month, _ = hdate.MonthFromName("Adar 2")
	fmt.Printf("%s (%d)\n", month, int(month))
	// Output:
	// Sh'vat (11)
	// Cheshvan (8)
	// Tevet (10)
	// Adar I (12)
	// Adar II (13)
}

func ExampleDaysInYear() {
	days := hdate.DaysInYear(5782)
	fmt.Println(days)
	// Output: 384
}

func ExampleDaysInMonth() {
	days := hdate.DaysInMonth(hdate.Kislev, 5783)
	fmt.Println(days)
	// Output: 30
}

func ExampleHebrewToRD() {
	rataDie := hdate.HebrewToRD(5769, hdate.Cheshvan, 15)
	fmt.Println(rataDie)
	// Output: 733359
}

func ExampleHDate_Greg() {
	hd := hdate.New(5765, hdate.Adar2, 22)
	year, month, day := hd.Greg()
	fmt.Println(year, month, day)
	// Output: 2005 April 2
}

func ExampleHDate_Abs() {
	hd := hdate.New(5765, hdate.Adar2, 22)
	rataDie := hd.Abs()
	fmt.Println(rataDie)
	// Output: 732038
}

func ExampleHDate_DaysInMonth() {
	hd := hdate.New(5765, hdate.Adar2, 22)
	days := hd.DaysInMonth()
	fmt.Println(days)
	// Output: 29
}

func ExampleHDate_Weekday() {
	hd := hdate.New(5765, hdate.Adar2, 22)
	dayOfWeek := hd.Weekday()
	fmt.Println(dayOfWeek)
	// Output: Saturday
}

func ExampleHDate_MonthName() {
	hd := hdate.New(5765, hdate.Adar2, 22)
	fmt.Println(hd.MonthName("en"))
	fmt.Println(hd.MonthName("he"))
	// Output:
	// Adar II
	// אַדָר ב׳
}

func ExampleHDate_IsLeapYear() {
	hd := hdate.New(5765, hdate.Adar2, 22)
	leap := hd.IsLeapYear()
	fmt.Println(leap)
	// Output: true
}

func ExampleHDate_Gregorian() {
	hd := hdate.New(5765, hdate.Adar2, 22)
	t := hd.Gregorian()
	fmt.Println(t)
	// Output: 2005-04-02 00:00:00 +0000 UTC
}

func ExampleHDate_Before() {
	orig := hdate.FromGregorian(2014, time.February, 19)
	hd := orig.Before(time.Saturday)
	fmt.Println(hd.Gregorian().Format(time.RFC1123))
	// Output: Sat, 15 Feb 2014 00:00:00 UTC
}

func ExampleHDate_OnOrBefore() {
	orig := hdate.FromGregorian(2014, time.February, 19)
	hd := orig.OnOrBefore(time.Saturday)
	fmt.Println(hd.Gregorian().Format(time.RFC1123))
	orig = hdate.FromGregorian(2014, time.February, 22)
	hd = orig.OnOrBefore(time.Saturday)
	fmt.Println(hd.Gregorian().Format(time.RFC1123))
	orig = hdate.FromGregorian(2014, time.February, 23)
	hd = orig.OnOrBefore(time.Saturday)
	fmt.Println(hd.Gregorian().Format(time.RFC1123))
	// Output:
	// Sat, 15 Feb 2014 00:00:00 UTC
	// Sat, 22 Feb 2014 00:00:00 UTC
	// Sat, 22 Feb 2014 00:00:00 UTC
}

func ExampleHDate_Nearest() {
	orig := hdate.FromGregorian(2014, time.February, 19)
	hd := orig.Nearest(time.Saturday)
	fmt.Println(hd.Gregorian().Format(time.RFC1123))
	orig = hdate.FromGregorian(2014, time.February, 18)
	hd = orig.Nearest(time.Saturday)
	fmt.Println(hd.Gregorian().Format(time.RFC1123))
	// Output:
	// Sat, 22 Feb 2014 00:00:00 UTC
	// Sat, 15 Feb 2014 00:00:00 UTC
}

func ExampleHDate_OnOrAfter() {
	orig := hdate.FromGregorian(2014, time.February, 19)
	hd := orig.OnOrAfter(time.Saturday)
	fmt.Println(hd.Gregorian().Format(time.RFC1123))
	orig = hdate.FromGregorian(2014, time.February, 22)
	hd = orig.OnOrAfter(time.Saturday)
	fmt.Println(hd.Gregorian().Format(time.RFC1123))
	orig = hdate.FromGregorian(2014, time.February, 23)
	hd = orig.OnOrAfter(time.Saturday)
	fmt.Println(hd.Gregorian().Format(time.RFC1123))
	// Output:
	// Sat, 22 Feb 2014 00:00:00 UTC
	// Sat, 22 Feb 2014 00:00:00 UTC
	// Sat, 01 Mar 2014 00:00:00 UTC
}

func ExampleHDate_After() {
	orig := hdate.FromGregorian(2014, time.February, 19)
	hd := orig.After(time.Saturday)
	fmt.Println(hd.Gregorian().Format(time.RFC1123))
	orig = hdate.FromGregorian(2014, time.February, 22)
	hd = orig.After(time.Saturday)
	fmt.Println(hd.Gregorian().Format(time.RFC1123))
	orig = hdate.FromGregorian(2014, time.February, 23)
	hd = orig.After(time.Saturday)
	fmt.Println(hd.Gregorian().Format(time.RFC1123))
	// Output:
	// Sat, 22 Feb 2014 00:00:00 UTC
	// Sat, 01 Mar 2014 00:00:00 UTC
	// Sat, 01 Mar 2014 00:00:00 UTC
}

func ExampleNew() {
	hd := hdate.New(5782, hdate.Nisan, 30)
	fmt.Println(hd)
	// Output: 30 Nisan 5782
}

func ExampleHDate_Next() {
	orig := hdate.New(5782, hdate.Nisan, 30)
	hd := orig.Next()
	fmt.Println(hd)
	// Output: 1 Iyyar 5782
}

func ExampleHDate_Prev() {
	orig := hdate.New(5782, hdate.Tishrei, 1)
	hd := orig.Prev()
	fmt.Println(hd)
	// Output: 29 Elul 5781
}

func ExampleHDate_String() {
	hd := hdate.New(5783, hdate.Elul, 29)
	str := hd.String()
	fmt.Println(str)
	// Output: 29 Elul 5783
}

func ExampleHMonth_String() {
	month := hdate.Adar2
	fmt.Println(month.String())
	// Output: Adar II
}

func ExampleHDate_Month() {
	hd := hdate.New(5769, hdate.Cheshvan, 15)
	fmt.Println(hd.Month())
	// Output: Cheshvan
}

func ExampleHDate_Year() {
	hd := hdate.New(5769, hdate.Cheshvan, 15)
	fmt.Println(hd.Year())
	// Output: 5769
}

func ExampleHDate_Day() {
	hd := hdate.New(5769, hdate.Cheshvan, 15)
	fmt.Println(hd.Day())
	// Output: 15
}

/***** anniversary ******/

func ExampleGetYahrzeit() {
	hd := hdate.FromGregorian(2014, time.March, 2)
	fmt.Println(hd)
	yahrzeit, _ := hdate.GetYahrzeit(5783, hd)
	fmt.Println(yahrzeit)
	// Output:
	// 30 Adar I 5774
	// 30 Sh'vat 5783
}

func ExampleGetBirthdayOrAnniversary() {
	hd := hdate.FromGregorian(2014, time.March, 2)
	fmt.Println(hd)
	yahrzeit, _ := hdate.GetBirthdayOrAnniversary(5783, hd)
	fmt.Println(yahrzeit)
	// Output:
	// 30 Adar I 5774
	// 1 Nisan 5783
}
