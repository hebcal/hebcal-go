package hebcal

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDafYomi(t *testing.T) {
	assert := assert.New(t)
	daf, _ := GetDafYomi(NewHDateFromGregorian(1995, time.December, 17))
	assert.Equal(DafYomi{Name: "Avodah Zarah", Blatt: 68}, daf)
	daf, _ = GetDafYomi(NewHDateFromGregorian(2020, time.June, 18))
	assert.Equal(DafYomi{Name: "Shabbat", Blatt: 104}, daf)
	daf, _ = GetDafYomi(NewHDateFromGregorian(2021, time.March, 23))
	assert.Equal(DafYomi{Name: "Shekalim", Blatt: 2}, daf)

}

func TestDafYomiEarlyCycles(t *testing.T) {
	assert := assert.New(t)
	daf, _ := GetDafYomi(NewHDateFromGregorian(1924, time.December, 7))
	assert.Equal(DafYomi{Name: "Shekalim", Blatt: 12}, daf)
	daf, _ = GetDafYomi(NewHDateFromGregorian(1924, time.December, 8))
	assert.Equal(DafYomi{Name: "Shekalim", Blatt: 13}, daf)
	daf, _ = GetDafYomi(NewHDateFromGregorian(1924, time.December, 9))
	assert.Equal(DafYomi{Name: "Yoma", Blatt: 2}, daf)
	daf, _ = GetDafYomi(NewHDateFromGregorian(1924, time.December, 10))
	assert.Equal(DafYomi{Name: "Yoma", Blatt: 3}, daf)

	daf, _ = GetDafYomi(NewHDateFromGregorian(1961, time.December, 3))
	assert.Equal(DafYomi{Name: "Shekalim", Blatt: 12}, daf)
	daf, _ = GetDafYomi(NewHDateFromGregorian(1961, time.December, 4))
	assert.Equal(DafYomi{Name: "Shekalim", Blatt: 13}, daf)
	daf, _ = GetDafYomi(NewHDateFromGregorian(1961, time.December, 5))
	assert.Equal(DafYomi{Name: "Yoma", Blatt: 2}, daf)
	daf, _ = GetDafYomi(NewHDateFromGregorian(1961, time.December, 6))
	assert.Equal(DafYomi{Name: "Yoma", Blatt: 3}, daf)

	daf, _ = GetDafYomi(NewHDateFromGregorian(1976, time.September, 19))
	assert.Equal(DafYomi{Name: "Shekalim", Blatt: 12}, daf)
	daf, _ = GetDafYomi(NewHDateFromGregorian(1976, time.September, 20))
	assert.Equal(DafYomi{Name: "Shekalim", Blatt: 13}, daf)
	daf, _ = GetDafYomi(NewHDateFromGregorian(1976, time.September, 21))
	assert.Equal(DafYomi{Name: "Shekalim", Blatt: 14}, daf)
	daf, _ = GetDafYomi(NewHDateFromGregorian(1976, time.September, 22))
	assert.Equal(DafYomi{Name: "Shekalim", Blatt: 15}, daf)
}

func ExampleGetDafYomi() {
	hd := NewHDateFromGregorian(1995, time.December, 17)
	daf, _ := GetDafYomi(hd)
	fmt.Println(daf)
	// Output: Avodah Zarah 68
}
