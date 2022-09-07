package hebcal

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
	assert.Equal(HDate{Year: 5769, Month: Cheshvan, Day: 15, abs: 733359}, NewHDateFromRD(733359))
	assert.Equal(HDate{Year: 5708, Month: Iyyar, Day: 6, abs: 711262}, NewHDateFromRD(711262))
	assert.Equal(HDate{Year: 3762, Month: Tishrei, Day: 1, abs: 249}, NewHDateFromRD(249))
	assert.Equal(HDate{Year: 3761, Month: Nisan, Day: 1, abs: 72}, NewHDateFromRD(72))
	assert.Equal(HDate{Year: 3761, Month: Nisan, Day: 8, abs: 79}, NewHDateFromRD(79))
	assert.Equal(HDate{Year: 3761, Month: Tevet, Day: 18, abs: 1}, NewHDateFromRD(1))
	assert.Equal(HDate{Year: 3761, Month: Tevet, Day: 17, abs: 0}, NewHDateFromRD(0))
	assert.Equal(HDate{Year: 3761, Month: Tevet, Day: 16, abs: -1}, NewHDateFromRD(-1))
	assert.Equal(HDate{Year: 3761, Month: Tevet, Day: 1, abs: -16}, NewHDateFromRD(-16))
	assert.Equal(HDate{Year: 3761, Month: Kislev, Day: 30, abs: -17}, NewHDateFromRD(-17))
	assert.Equal(HDate{Year: 9999, Month: Elul, Day: 29, abs: 2278650}, NewHDateFromRD(2278650))
	assert.Equal(HDate{Year: 5765, Month: Adar2, Day: 22, abs: 732038}, NewHDateFromRD(732038))
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
	hd := NewHDate(5782, Adar1, 15)
	assert.Equal("Adar I", hd.MonthName())
	hd = NewHDate(5783, Adar1, 15)
	assert.Equal("Adar", hd.MonthName())
}

func TestAdar2ResetToAdar1(t *testing.T) {
	assert := assert.New(t)
	hd := NewHDate(5782, Adar1, 15)
	assert.Equal(Adar1, hd.Month)
	hd = NewHDate(5782, Adar2, 15)
	assert.Equal(Adar2, hd.Month)
	hd = NewHDate(5783, Adar1, 15)
	assert.Equal(Adar1, hd.Month)
	hd = NewHDate(5783, Adar2, 15)
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
	assert.Equal(385, DaysInHebYear(5779))
	assert.Equal(355, DaysInHebYear(5780))
	assert.Equal(353, DaysInHebYear(5781))
	assert.Equal(384, DaysInHebYear(5782))
	assert.Equal(355, DaysInHebYear(5783))
	assert.Equal(383, DaysInHebYear(5784))
	assert.Equal(355, DaysInHebYear(5785))
	assert.Equal(354, DaysInHebYear(5786))
	assert.Equal(385, DaysInHebYear(5787))
	assert.Equal(355, DaysInHebYear(5788))
	assert.Equal(354, DaysInHebYear(5789))
	assert.Equal(383, DaysInHebYear(3762))
	assert.Equal(354, DaysInHebYear(3671))
	assert.Equal(353, DaysInHebYear(1234))
	assert.Equal(355, DaysInHebYear(123))
	assert.Equal(355, DaysInHebYear(2))
	assert.Equal(355, DaysInHebYear(1))

	assert.Equal(353, DaysInHebYear(5761))
	assert.Equal(354, DaysInHebYear(5762))
	assert.Equal(385, DaysInHebYear(5763))
	assert.Equal(355, DaysInHebYear(5764))
	assert.Equal(383, DaysInHebYear(5765))
	assert.Equal(354, DaysInHebYear(5766))
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
	assert.Equal(time.Thursday, NewHDate(5769, Cheshvan, 15).Weekday())
	assert.Equal(time.Saturday, NewHDate(5708, Iyyar, 6).Weekday())
	assert.Equal(time.Sunday, NewHDate(5708, Iyyar, 7).Weekday())
	assert.Equal(time.Thursday, NewHDate(3762, Tishrei, 1).Weekday())
	assert.Equal(time.Tuesday, NewHDate(3761, Nisan, 1).Weekday())
	assert.Equal(time.Monday, NewHDate(3761, Tevet, 18).Weekday())
	assert.Equal(time.Sunday, NewHDate(3761, Tevet, 17).Weekday())
	assert.Equal(time.Saturday, NewHDate(3761, Tevet, 16).Weekday())
	assert.Equal(time.Friday, NewHDate(3761, Tevet, 1).Weekday())
	assert.Equal(time.Tuesday, NewHDate(3333, Sivan, 29).Weekday())
	assert.Equal(time.Monday, NewHDate(3333, Sivan, 28).Weekday())
	assert.Equal(time.Sunday, NewHDate(3333, Sivan, 27).Weekday())
	assert.Equal(time.Saturday, NewHDate(3333, Sivan, 26).Weekday())
	assert.Equal(time.Friday, NewHDate(3333, Sivan, 25).Weekday())
	assert.Equal(time.Thursday, NewHDate(3333, Sivan, 24).Weekday())
	assert.Equal(time.Wednesday, NewHDate(3333, Sivan, 23).Weekday())
}

func hd2iso(hd HDate) string {
	year, month, day := hd.Greg()
	d := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return d.Format(time.RFC3339)[:10]
}

func TestBefore(t *testing.T) {
	assert := assert.New(t)
	d := time.Date(2014, time.February, 19, 0, 0, 0, 0, time.UTC)
	hd := NewHDateFromTime(d)
	assert.Equal("2014-02-15", hd2iso(hd.Before(time.Saturday)))
}

func TestOnOrBefore(t *testing.T) {
	assert := assert.New(t)
	d := time.Date(2014, time.February, 19, 0, 0, 0, 0, time.UTC)
	assert.Equal("2014-02-15", hd2iso(NewHDateFromTime(d).OnOrBefore(time.Saturday)))
	d = time.Date(2014, time.February, 22, 0, 0, 0, 0, time.UTC)
	assert.Equal("2014-02-22", hd2iso(NewHDateFromTime(d).OnOrBefore(time.Saturday)))
	d = time.Date(2014, time.February, 23, 0, 0, 0, 0, time.UTC)
	assert.Equal("2014-02-22", hd2iso(NewHDateFromTime(d).OnOrBefore(time.Saturday)))

}

func TestNearest(t *testing.T) {
	assert := assert.New(t)
	d := time.Date(2014, time.February, 19, 0, 0, 0, 0, time.UTC)
	assert.Equal("2014-02-22", hd2iso(NewHDateFromTime(d).Nearest(time.Saturday)))
	d = time.Date(2014, time.February, 18, 0, 0, 0, 0, time.UTC)
	assert.Equal("2014-02-15", hd2iso(NewHDateFromTime(d).Nearest(time.Saturday)))
}

func TestOnOrAfter(t *testing.T) {
	assert := assert.New(t)
	d := time.Date(2014, time.February, 19, 0, 0, 0, 0, time.UTC)
	assert.Equal("2014-02-22", hd2iso(NewHDateFromTime(d).OnOrAfter(time.Saturday)))
	d = time.Date(2014, time.February, 22, 0, 0, 0, 0, time.UTC)
	assert.Equal("2014-02-22", hd2iso(NewHDateFromTime(d).OnOrAfter(time.Saturday)))
	d = time.Date(2014, time.February, 23, 0, 0, 0, 0, time.UTC)
	assert.Equal("2014-03-01", hd2iso(NewHDateFromTime(d).OnOrAfter(time.Saturday)))
}

func TestAfter(t *testing.T) {
	assert := assert.New(t)
	d := time.Date(2014, time.February, 19, 0, 0, 0, 0, time.UTC)
	assert.Equal("2014-02-22", hd2iso(NewHDateFromTime(d).After(time.Saturday)))
	d = time.Date(2014, time.February, 22, 0, 0, 0, 0, time.UTC)
	assert.Equal("2014-03-01", hd2iso(NewHDateFromTime(d).After(time.Saturday)))
	d = time.Date(2014, time.February, 23, 0, 0, 0, 0, time.UTC)
	assert.Equal("2014-03-01", hd2iso(NewHDateFromTime(d).After(time.Saturday)))
}

func TestToString(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("4 Tevet 5511", NewHDate(5511, Tevet, 4).String())
	assert.Equal("4 Elul 5782", NewHDate(5782, Elul, 4).String())
	assert.Equal("29 Adar II 5749", NewHDate(5749, Adar2, 29).String())
}

func TestGreg(t *testing.T) {
	assert := assert.New(t)
	hd := NewHDate(5765, Adar2, 22)
	gy, gm, gd := hd.Greg()
	assert.Equal(2005, gy)
	assert.Equal(time.April, gm)
	assert.Equal(2, gd)
}

func ExampleNewHDateFromTime() {
	d := time.Date(2008, time.November, 13, 0, 0, 0, 0, time.UTC)
	fmt.Println(NewHDateFromTime(d).String())
	// Output: 15 Cheshvan 5769
}

func ExampleNewHDateFromRD() {
	fmt.Println(NewHDateFromRD(733359).String())
	// Output: 15 Cheshvan 5769
}

func ExampleMonthFromName() {
	m1, _ := MonthFromName("Shvat")
	m2, _ := MonthFromName("cheshvan")
	m3, _ := MonthFromName("טבת")
	fmt.Printf("%s (%d), %s (%d), %s (%d)", m1, int(m1), m2, int(m2), m3, int(m3))
	// Output: Sh'vat (11), Cheshvan (8), Tevet (10)
}
