package event

import (
	"strings"
	"testing"
	"time"

	"github.com/hebcal/hdate"
	"github.com/hebcal/hebcal-go/sedra"
)

func TestHolidayEventURL(t *testing.T) {
	tests := []struct {
		name string
		ev   HolidayEvent
		want string
	}{
		{
			name: "ordinary holiday uses the Gregorian year",
			ev:   NewHolidayEvent(hdate.New(5787, hdate.Tishrei, 15), "Sukkot I", Chag, ""),
			want: "https://www.hebcal.com/holidays/sukkot-2026",
		},
		{
			name: "apostrophes are dropped from the slug",
			ev:   NewHolidayEvent(hdate.New(5786, hdate.Iyyar, 5), "Yom HaAtzma'ut", ModernHoliday, ""),
			want: "https://www.hebcal.com/holidays/yom-haatzmaut-2026",
		},
		{
			name: "Israel-only observance carries i=on",
			ev:   NewHolidayEvent(hdate.New(5786, hdate.Iyyar, 5), "Yom HaAtzma'ut", ModernHoliday|ILOnly, ""),
			want: "https://www.hebcal.com/holidays/yom-haatzmaut-2026?i=on",
		},
		{
			name: "Yom Kippur Katan has no page",
			ev:   NewHolidayEvent(hdate.New(5780, hdate.Shvat, 29), "Yom Kippur Katan Adar", MinorFast, ""),
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ev.URL(); got != tt.want {
				t.Errorf("URL() = %q, want %q", got, tt.want)
			}
		})
	}
}

// Asara B'Tevet can fall twice in one Gregorian year, so its URL is
// disambiguated with a full date rather than just the year.
func TestAsaraBTevetURLUsesFullDate(t *testing.T) {
	ev := NewHolidayEvent(hdate.New(5787, hdate.Tevet, 10), "Asara B'Tevet", MinorFast, "")
	want := "https://www.hebcal.com/holidays/asara-btevet-20261220"
	if got := ev.URL(); got != want {
		t.Errorf("URL() = %q, want %q", got, want)
	}
}

// Chanukah can start in December and end in January; the whole festival is
// filed under the year it began, so the January days point back a year.
func TestChanukahJanuaryDaysUsePreviousYear(t *testing.T) {
	// 5784: Chanukah runs 7 Dec 2023 - 15 Dec 2023, entirely in December.
	dec := NewHolidayEvent(hdate.New(5784, hdate.Kislev, 25), "Chanukah: 1 Candle", MinorHoliday, "")
	if got, want := dec.URL(), "https://www.hebcal.com/holidays/chanukah-2023"; got != want {
		t.Errorf("December day: URL() = %q, want %q", got, want)
	}
	// 5785: the last days fall in January 2025 but still belong to 2024.
	jan := NewHolidayEvent(hdate.New(5785, hdate.Tevet, 2), "Chanukah: 8th Day", MinorHoliday, "")
	gy, gm, _ := jan.Date().ProlepticGreg()
	if gm != 1 {
		t.Fatalf("test precondition: expected a January date, got %d-%02d", gy, int(gm))
	}
	if got, want := jan.URL(), "https://www.hebcal.com/holidays/chanukah-2024"; got != want {
		t.Errorf("January day: URL() = %q, want %q", got, want)
	}
}

// In a leap year both Adars carry the description "Rosh Chodesh Adar", but
// hebcal.com files them under separate pages. A two-day Rosh Chodesh starts on
// the 30th of the preceding month, which must not be mistaken for that month.
func TestLeapYearAdarURLsAreDistinct(t *testing.T) {
	if !hdate.IsLeapYear(5787) {
		t.Fatal("test precondition: 5787 should be a leap year")
	}
	cases := []struct {
		name string
		ev   HolidayEvent
		want string
	}{
		{"Adar I, first day (30 Sh'vat)", NewHolidayEvent(hdate.New(5787, hdate.Shvat, 30), "Rosh Chodesh Adar", RoshChodesh, ""), "rosh-chodesh-adar-i"},
		{"Adar I, second day", NewHolidayEvent(hdate.New(5787, hdate.Adar1, 1), "Rosh Chodesh Adar", RoshChodesh, ""), "rosh-chodesh-adar-i"},
		{"Adar II, first day (30 Adar I)", NewHolidayEvent(hdate.New(5787, hdate.Adar1, 30), "Rosh Chodesh Adar", RoshChodesh, ""), "rosh-chodesh-adar-ii"},
		{"Adar II, second day", NewHolidayEvent(hdate.New(5787, hdate.Adar2, 1), "Rosh Chodesh Adar", RoshChodesh, ""), "rosh-chodesh-adar-ii"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.ev.URL(); !strings.Contains(got, "/"+c.want+"-") {
				t.Errorf("URL() = %q, want it to contain %q", got, c.want)
			}
		})
	}
}

// A common year has a single Adar and keeps the undecorated slug, even though
// hdate numbers that month Adar1.
func TestCommonYearAdarHasNoNumeral(t *testing.T) {
	if hdate.IsLeapYear(5780) {
		t.Fatal("test precondition: 5780 should be a common year")
	}
	ev := NewHolidayEvent(hdate.New(5780, hdate.Adar1, 1), "Rosh Chodesh Adar", RoshChodesh, "")
	if got, want := ev.URL(), "https://www.hebcal.com/holidays/rosh-chodesh-adar-2020"; got != want {
		t.Errorf("URL() = %q, want %q", got, want)
	}
}

// hebcal.com spells this month with one m in Rosh Chodesh paths, but the fast
// keeps the two-m spelling.
func TestTammuzSlugs(t *testing.T) {
	rc := NewHolidayEvent(hdate.New(5787, hdate.Tamuz, 1), "Rosh Chodesh Tammuz", RoshChodesh, "")
	if !strings.Contains(rc.URL(), "rosh-chodesh-tamuz-") {
		t.Errorf("Rosh Chodesh: URL() = %q, want rosh-chodesh-tamuz", rc.URL())
	}
	fast := NewHolidayEvent(hdate.New(5787, hdate.Tamuz, 17), "Tzom Tammuz", MinorFast, "")
	if !strings.Contains(fast.URL(), "tzom-tammuz-") {
		t.Errorf("fast: URL() = %q, want tzom-tammuz", fast.URL())
	}
}

// Years the website does not cover produce no link rather than one that 404s.
func TestOutOfRangeYearsHaveNoURL(t *testing.T) {
	for _, hy := range []int{3000, 8000} {
		ev := NewHolidayEvent(hdate.New(hy, hdate.Tishrei, 15), "Sukkot I", Chag, "")
		gy, _, _ := ev.Date().ProlepticGreg()
		if gy >= 100 && gy <= 2999 {
			continue // still in range; nothing to assert
		}
		if got := ev.URL(); got != "" {
			t.Errorf("Hebrew year %d (Gregorian %d): URL() = %q, want empty", hy, gy, got)
		}
	}
}

// parshaFor builds a Parsha value from portion names, for URL tests that do
// not need the real reading schedule.
func parshaFor(names ...string) sedra.Parsha {
	return sedra.Parsha{Name: names}
}

func TestParshaEventURL(t *testing.T) {
	// 8 Aug 2026 is a Saturday; Parashat Eikev is read that day.
	hd := hdate.FromGregorian(2026, time.August, 8)
	ev := NewParshaEvent(hd, parshaFor("Eikev"), false)
	want := "https://www.hebcal.com/sedrot/eikev-20260808"
	if got := URL(ev); got != want {
		t.Errorf("URL() = %q, want %q", got, want)
	}
	il := NewParshaEvent(hd, parshaFor("Eikev"), true)
	if got := URL(il); got != want+"?i=on" {
		t.Errorf("Israel URL() = %q, want %q", got, want+"?i=on")
	}
}

// A doubled portion keeps both names in the slug.
func TestDoubledParshaURL(t *testing.T) {
	hd := hdate.FromGregorian(2026, time.March, 14)
	ev := NewParshaEvent(hd, parshaFor("Vayakhel", "Pekudei"), false)
	if got := URL(ev); !strings.Contains(got, "/sedrot/vayakhel-pekudei-") {
		t.Errorf("URL() = %q, want it to contain /sedrot/vayakhel-pekudei-", got)
	}
}

func TestHebrewDateEventURL(t *testing.T) {
	ev := NewHebrewDateEvent(hdate.New(5786, hdate.Nisan, 15))
	got := URL(ev)
	for _, want := range []string{"hd=15", "hm=Nisan", "hy=5786", "h2g=1"} {
		if !strings.Contains(got, want) {
			t.Errorf("URL() = %q, want it to contain %q", got, want)
		}
	}
}

// URL returns "" for events that implement CalEvent but not URLer, rather
// than panicking on the type assertion.
func TestURLHelperHandlesEventsWithoutURLs(t *testing.T) {
	ev := UserEvent{}
	if got := URL(ev); got != "" {
		t.Errorf("URL() = %q, want empty", got)
	}
}
