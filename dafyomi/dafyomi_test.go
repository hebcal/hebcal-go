package dafyomi_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hebcal/hdate"
	"github.com/MaxBGreenberg/hebcal-go/dafyomi"
	"github.com/stretchr/testify/assert"
)

func TestDafYomi(t *testing.T) {
	assert := assert.New(t)
	daf, _ := dafyomi.New(hdate.FromGregorian(1995, time.December, 17))
	assert.Equal(dafyomi.Daf{Name: "Avodah Zarah", Blatt: 68}, daf)
	daf, _ = dafyomi.New(hdate.FromGregorian(2020, time.June, 18))
	assert.Equal(dafyomi.Daf{Name: "Shabbat", Blatt: 104}, daf)
	daf, _ = dafyomi.New(hdate.FromGregorian(2021, time.March, 23))
	assert.Equal(dafyomi.Daf{Name: "Shekalim", Blatt: 2}, daf)
}

func TestDafYomiEarlyCycles(t *testing.T) {
	assert := assert.New(t)
	daf, _ := dafyomi.New(hdate.FromGregorian(1924, time.December, 7))
	assert.Equal(dafyomi.Daf{Name: "Shekalim", Blatt: 12}, daf)
	daf, _ = dafyomi.New(hdate.FromGregorian(1924, time.December, 8))
	assert.Equal(dafyomi.Daf{Name: "Shekalim", Blatt: 13}, daf)
	daf, _ = dafyomi.New(hdate.FromGregorian(1924, time.December, 9))
	assert.Equal(dafyomi.Daf{Name: "Yoma", Blatt: 2}, daf)
	daf, _ = dafyomi.New(hdate.FromGregorian(1924, time.December, 10))
	assert.Equal(dafyomi.Daf{Name: "Yoma", Blatt: 3}, daf)

	daf, _ = dafyomi.New(hdate.FromGregorian(1961, time.December, 3))
	assert.Equal(dafyomi.Daf{Name: "Shekalim", Blatt: 12}, daf)
	daf, _ = dafyomi.New(hdate.FromGregorian(1961, time.December, 4))
	assert.Equal(dafyomi.Daf{Name: "Shekalim", Blatt: 13}, daf)
	daf, _ = dafyomi.New(hdate.FromGregorian(1961, time.December, 5))
	assert.Equal(dafyomi.Daf{Name: "Yoma", Blatt: 2}, daf)
	daf, _ = dafyomi.New(hdate.FromGregorian(1961, time.December, 6))
	assert.Equal(dafyomi.Daf{Name: "Yoma", Blatt: 3}, daf)

	daf, _ = dafyomi.New(hdate.FromGregorian(1976, time.September, 19))
	assert.Equal(dafyomi.Daf{Name: "Shekalim", Blatt: 12}, daf)
	daf, _ = dafyomi.New(hdate.FromGregorian(1976, time.September, 20))
	assert.Equal(dafyomi.Daf{Name: "Shekalim", Blatt: 13}, daf)
	daf, _ = dafyomi.New(hdate.FromGregorian(1976, time.September, 21))
	assert.Equal(dafyomi.Daf{Name: "Shekalim", Blatt: 14}, daf)
	daf, _ = dafyomi.New(hdate.FromGregorian(1976, time.September, 22))
	assert.Equal(dafyomi.Daf{Name: "Shekalim", Blatt: 15}, daf)
}

func ExampleNew() {
	hd := hdate.FromGregorian(1995, time.December, 17)
	daf, _ := dafyomi.New(hd)
	fmt.Println(daf)
	// Output: Avodah Zarah 68
}
