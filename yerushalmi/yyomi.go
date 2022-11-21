// Hebcal's yerushalmi package calculates the Yerushalmi Yomi, a
// daily regimen of learning the Jerusalem Talmud.
//
// The Yerushalmi Daf Yomi program takes approx. 4.25 years or 51 months.
// Unlike the Daf Yomi Bavli cycle, the Yerushalmi cycle skips both
// Yom Kippur and Tisha B'Av. The page numbers are according to the Vilna
// Edition which is used since 1900.
//
// https://en.wikipedia.org/wiki/Jerusalem_Talmud
package yerushalmi

import (
	"time"

	"github.com/hebcal/hebcal-go/dafyomi"
	"github.com/hebcal/hebcal-go/greg"
	"github.com/hebcal/hebcal-go/hdate"
)

var shas = []dafyomi.Daf{
	{Name: "Berakhot", Blatt: 68},
	{Name: "Peah", Blatt: 37},
	{Name: "Demai", Blatt: 34},
	{Name: "Kilayim", Blatt: 44},
	{Name: "Sheviit", Blatt: 31},
	{Name: "Terumot", Blatt: 59},
	{Name: "Maasrot", Blatt: 26},
	{Name: "Maaser Sheni", Blatt: 33},
	{Name: "Challah", Blatt: 28},
	{Name: "Orlah", Blatt: 20},
	{Name: "Bikkurim", Blatt: 13},
	{Name: "Shabbat", Blatt: 92},
	{Name: "Eruvin", Blatt: 65},
	{Name: "Pesachim", Blatt: 71},
	{Name: "Beitzah", Blatt: 22},
	{Name: "Rosh Hashanah", Blatt: 22},
	{Name: "Yoma", Blatt: 42},
	{Name: "Sukkah", Blatt: 26},
	{Name: "Taanit", Blatt: 26},
	{Name: "Shekalim", Blatt: 33},
	{Name: "Megillah", Blatt: 34},
	{Name: "Chagigah", Blatt: 22},
	{Name: "Moed Katan", Blatt: 19},
	{Name: "Yevamot", Blatt: 85},
	{Name: "Ketubot", Blatt: 72},
	{Name: "Sotah", Blatt: 47},
	{Name: "Nedarim", Blatt: 40},
	{Name: "Nazir", Blatt: 47},
	{Name: "Gittin", Blatt: 54},
	{Name: "Kiddushin", Blatt: 48},
	{Name: "Bava Kamma", Blatt: 44},
	{Name: "Bava Metzia", Blatt: 37},
	{Name: "Bava Batra", Blatt: 34},
	{Name: "Shevuot", Blatt: 44},
	{Name: "Makkot", Blatt: 9},
	{Name: "Sanhedrin", Blatt: 57},
	{Name: "Avodah Zarah", Blatt: 37},
	{Name: "Horayot", Blatt: 19},
	{Name: "Niddah", Blatt: 13},
}

// The number of pages in the Talmud Yerushalmi
const numDapim = 1554

// YerushalmiYomiStartRD is the R.D. date of the first cycle of
// Yerushalmi Yomi.
var YerushalmiYomiStartRD = greg.ToRD(1980, time.February, 2)

// New calculates the Daf Yomi Yerushalmi for given date.
//
// Returns an empty Daf for Yom Kippur and Tisha B'Av.
//
// Panics if the date is before Daf Yomi Yerushalmi cycle began
// (2 February 1980).
func New(hd hdate.HDate) dafyomi.Daf {
	cday := hd.Abs()
	if cday < YerushalmiYomiStartRD {
		panic(hd.String() + " is before Daf Yomi Yerushalmi cycle began")
	}

	// No Daf for Yom Kippur and Tisha B'Av
	if (hd.Month() == hdate.Tishrei && hd.Day() == 10) ||
		(hd.Month() == hdate.Av &&
			((hd.Day() == 9 && hd.Weekday() != time.Saturday) ||
				(hd.Day() == 10 && hd.Weekday() == time.Sunday))) {
		return dafyomi.Daf{}
	}

	nextCycle := YerushalmiYomiStartRD
	prevCycle := YerushalmiYomiStartRD
	for cday >= nextCycle {
		prevCycle = nextCycle
		nextCycle += numDapim
		nextCycle += numSpecialDays(prevCycle, nextCycle)
	}

	total := cday - prevCycle - numSpecialDays(prevCycle, cday)

	for j := 0; j < len(shas); j++ {
		masechet := shas[j]
		if total < masechet.Blatt {
			return dafyomi.Daf{Name: masechet.Name, Blatt: total + 1}
		}
		total -= masechet.Blatt
	}

	panic("Interal error, this code should be unreachable")
}

func numSpecialDays(startAbs, endAbs int) int {
	startYear := hdate.FromRD(startAbs).Year()
	endYear := hdate.FromRD(endAbs).Year()

	specialDays := 0
	for year := startYear; year <= endYear; year++ {
		yk := hdate.New(year, hdate.Tishrei, 10)
		ykAbs := yk.Abs()
		if ykAbs >= startAbs && ykAbs <= endAbs {
			specialDays++
		}
		av9dt := hdate.New(year, hdate.Av, 9)
		if av9dt.Weekday() == time.Saturday {
			av9dt = av9dt.Next()
		}
		av9abs := av9dt.Abs()
		if av9abs >= startAbs && av9abs <= endAbs {
			specialDays++
		}
	}
	return specialDays
}
