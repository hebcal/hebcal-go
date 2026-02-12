package sedra_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hebcal/hdate"
	"github.com/hebcal/hebcal-go/locales"
	"github.com/hebcal/hebcal-go/sedra"
	"github.com/stretchr/testify/assert"
)

func TestSedra_Lookup(t *testing.T) {
	assert := assert.New(t)
	sedraYear := sedra.New(5749, false)
	assert.Equal(
		sedra.Parsha{Chag: true},
		sedraYear.Lookup(hdate.FromGregorian(1988, time.October, 1)))
	assert.Equal(
		sedra.Parsha{Name: []string{"Chayei Sara"}, Num: []int{5}, Chag: false},
		sedraYear.Lookup(hdate.FromGregorian(1988, time.November, 5)))
	assert.Equal(
		sedra.Parsha{Name: []string{"Chukat", "Balak"}, Num: []int{39, 40}, Chag: false},
		sedraYear.Lookup(hdate.FromGregorian(1989, time.July, 15)))

	sedraYear = sedra.New(5781, false)
	assert.Equal(
		sedra.Parsha{Name: []string{"Achrei Mot", "Kedoshim"}, Num: []int{29, 30}, Chag: false},
		sedraYear.Lookup(hdate.FromGregorian(2021, time.April, 24)))
	assert.Equal(
		sedra.Parsha{Name: []string{"Bereshit"}, Num: []int{1}, Chag: false},
		sedraYear.Lookup(hdate.FromGregorian(2020, time.October, 17)))
}

func TestSedraAshkenazi(t *testing.T) {
	assert := assert.New(t)
	sedraYear := sedra.New(5749, false)
	want := []string{
		"Vayeilech",
		"Vayeilech",
		"Ha'azinu",
		"",
		"Bereishis",
		"Noach",
		"Lech Lecho",
		"Vayeiro",
		"Chayei Soroh",
		"Toldos",
		"Vayeitzei",
		"Vayishlach",
		"Vayeishev",
		"Mikeitz",
		"Vayigash",
		"Vaichi",
		"Shemos",
		"Vo'eiro",
		"Bo",
		"Beshalach",
		"Yisro",
		"Mishpotim",
		"Terumoh",
		"Tetzaveh",
		"Ki Siso",
		"Vayakheil",
		"Pekudei",
		"Vayikro",
		"Tzav",
		"Shemini",
		"Tazria",
		"Metzoro",
		"",
		"Acharei Mos",
		"Kedoshim",
		"Emor",
		"Behar",
		"Bechukosai",
		"Bemidbar",
		"",
		"Noso",
		"Beha'aloscho",
		"Shelach",
		"Korach",
		"Chukas-Bolok",
		"Pinchos",
		"Matos-Masei",
		"Devorim",
		"Vo'eschanan",
		"Eikev",
		"Re'eih",
		"Shoftim",
		"Ki Seitzei",
		"Ki Sovo",
		"Nitzovim-Vayeilech",
	}
	var got []string
	hd := hdate.New(5749, hdate.Tishrei, 1)
	for hd.Year() == 5749 {
		parsha := sedraYear.Lookup(hd).Render("ashkenazi")
		got = append(got, parsha)
		hd = hd.After(time.Saturday)
	}
	assert.Equal(want, got)
	vb, _ := locales.LookupTranslation("Vezot Haberakhah", "ashkenazi")
	assert.Equal("Vezos Haberochoh", vb)
}

func ExampleSedra_Lookup() {
	sedraYear := sedra.New(5749, false)
	hd := hdate.FromGregorian(1989, time.July, 15)
	parsha := sedraYear.Lookup(hd)
	fmt.Println(parsha)
	// Output: Parashat Chukat-Balak
}

func ExampleSedra_FindParshaNum() {
	sedraYear := sedra.New(5749, false)
	hd, _ := sedraYear.FindParshaNum(16)
	fmt.Println(hd)
	// Output: 15 Sh'vat 5749
}

func TestSedraYearTypes(t *testing.T) {
	years := []int{5701, 5702, 5703, 5708, 5710, 5711, 5713, 5714, 5715, 5717, 5719, 5726, 5734, 5736}
	for _, year := range years {
		diaspora := sedra.New(year, false)
		hd1, _ := diaspora.FindParshaNum(1)
		assert.Equal(t, year, hd1.Year())
		il := sedra.New(year, true)
		hd2, _ := il.FindParshaNum(1)
		assert.Equal(t, year, hd2.Year())
	}
}

func TestSedraEarlyYears(t *testing.T) {
	years := []int{3762, 3761, 3760, 3759, 100, 2, 1}
	for _, year := range years {
		s := sedra.New(year, false)
		assert.Equal(t, year, s.Year)
	}
}
