package dafyomi

import (
	"fmt"
	"testing"
	"time"

	"github.com/hebcal/hebcal-go/hdate"
	"github.com/stretchr/testify/assert"
)

func TestMishnaYomi(t *testing.T) {
	idx := MakeMishnaYomiIndex()
	hd := hdate.FromGregorian(1995, time.December, 17)
	mishna, err := idx.Lookup(hd)
	assert.Equal(t, nil, err)
	assert.Equal(t, mishna, MishnaPair{
		Mishna{Tractate: "Bava Kamma", Chap: 5, Verse: 7},
		Mishna{Tractate: "Bava Kamma", Chap: 6, Verse: 1},
	})
}

func ExampleMishnaYomiIndex_Lookup() {
	idx := MakeMishnaYomiIndex()
	mishna, _ := idx.Lookup(hdate.FromGregorian(1947, time.May, 20))
	fmt.Println(mishna)
	mishna, _ = idx.Lookup(hdate.FromGregorian(1995, time.December, 17))
	fmt.Println(mishna)
	mishna, _ = idx.Lookup(hdate.FromGregorian(2022, time.August, 1))
	fmt.Println(mishna)
	// Output:
	// Berakhot 1:1-2
	// Bava Kamma 5:7-6:1
	// Terumot 11:3-4
}
