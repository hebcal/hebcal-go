package sedra

import (
	"fmt"
	"testing"
	"time"

	"github.com/hebcal/hebcal-go/hdate"
	"github.com/stretchr/testify/assert"
)

func TestSedra_Lookup(t *testing.T) {
	assert := assert.New(t)
	sedra := New(5749, false)
	assert.Equal(
		Parsha{Chag: true},
		sedra.Lookup(hdate.FromGregorian(1988, time.October, 1)))
	assert.Equal(
		Parsha{Name: []string{"Chayei Sara"}, Num: []int{5}, Chag: false},
		sedra.Lookup(hdate.FromGregorian(1988, time.November, 5)))
	assert.Equal(
		Parsha{Name: []string{"Chukat", "Balak"}, Num: []int{39, 40}, Chag: false},
		sedra.Lookup(hdate.FromGregorian(1989, time.July, 15)))

	sedra = New(5781, false)
	assert.Equal(
		Parsha{Name: []string{"Achrei Mot", "Kedoshim"}, Num: []int{29, 30}, Chag: false},
		sedra.Lookup(hdate.FromGregorian(2021, time.April, 24)))
	assert.Equal(
		Parsha{Name: []string{"Bereshit"}, Num: []int{1}, Chag: false},
		sedra.Lookup(hdate.FromGregorian(2020, time.October, 17)))
}

func ExampleSedra_Lookup() {
	sedra := New(5749, false)
	hd := hdate.FromGregorian(1989, time.July, 15)
	parsha := sedra.Lookup(hd)
	fmt.Println(parsha)
	// Output: Parashat Chukat-Balak
}

func ExampleSedra_FindParshaNum() {
	sedra := New(5749, false)
	hd, _ := sedra.FindParshaNum(16)
	fmt.Println(hd)
	// Output: 15 Sh'vat 5749
}

func TestSedraYearTypes(t *testing.T) {
	years := []int{5701, 5702, 5703, 5708, 5710, 5711, 5713, 5714, 5715, 5717, 5719, 5726, 5734, 5736}
	for _, year := range years {
		diaspora := New(year, false)
		hd1, _ := diaspora.FindParshaNum(1)
		assert.Equal(t, year, hd1.Year())
		il := New(year, true)
		hd2, _ := il.FindParshaNum(1)
		assert.Equal(t, year, hd2.Year())
	}
}

func TestIsValidDouble(t *testing.T) {
	assert.Equal(t, true, isValidDouble(-26))
	assert.Equal(t, false, isValidDouble(-1))
	assert.Equal(t, false, isValidDouble(26))
}

func TestSedraEarlyYears(t *testing.T) {
	years := []int{3762, 3761, 3760, 3759, 100, 2, 1}
	for _, year := range years {
		fmt.Println(year)
		s := New(year, false)
		assert.Equal(t, year, s.Year)
	}
}
