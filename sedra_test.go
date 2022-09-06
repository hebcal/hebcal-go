package hebcal

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSedra(t *testing.T) {
	assert := assert.New(t)
	sedra := NewSedra(5749, false)
	assert.Equal(
		Parsha{Chag: true},
		sedra.Lookup(NewHDateFromGregorian(1988, time.October, 1)))
	assert.Equal(
		Parsha{Name: []string{"Chayei Sara"}, Num: []int{5}, Chag: false},
		sedra.Lookup(NewHDateFromGregorian(1988, time.November, 5)))
	assert.Equal(
		Parsha{Name: []string{"Chukat", "Balak"}, Num: []int{39, 40}, Chag: false},
		sedra.Lookup(NewHDateFromGregorian(1989, time.July, 15)))

	sedra = NewSedra(5781, false)
	assert.Equal(
		Parsha{Name: []string{"Achrei Mot", "Kedoshim"}, Num: []int{29, 30}, Chag: false},
		sedra.Lookup(NewHDateFromGregorian(2021, time.April, 24)))
	assert.Equal(
		Parsha{Name: []string{"Bereshit"}, Num: []int{1}, Chag: false},
		sedra.Lookup(NewHDateFromGregorian(2020, time.October, 17)))
}
