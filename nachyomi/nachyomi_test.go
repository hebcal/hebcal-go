package nachyomi_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/hebcal/hdate"
	"github.com/MaxBGreenberg/hebcal-go/dafyomi"
	"github.com/MaxBGreenberg/hebcal-go/nachyomi"
	"github.com/stretchr/testify/assert"
)

func TestNachYomi(t *testing.T) {
	idx := nachyomi.MakeIndex()
	chap, err := idx.Lookup(hdate.FromGregorian(2022, time.April, 5))
	assert.Equal(t, nil, err)
	assert.Equal(t, chap, dafyomi.Daf{Name: "I Samuel", Blatt: 31})
	chap, err = idx.Lookup(hdate.FromGregorian(2022, time.October, 25))
	assert.Equal(t, nil, err)
	assert.Equal(t, chap, dafyomi.Daf{Name: "Ezekiel", Blatt: 14})
}

func TestNachYomiBefore(t *testing.T) {
	idx := nachyomi.MakeIndex()
	hd := hdate.FromGregorian(1995, time.December, 17)
	chap, err := idx.Lookup(hd)
	assert.Equal(t, errors.New("before Nach Yomi cycle began"), err)
	assert.Equal(t, chap, dafyomi.Daf{Name: "", Blatt: 0})
}

func ExampleNachYomiIndex_Lookup() {
	idx := nachyomi.MakeIndex()
	chapter, _ := idx.Lookup(hdate.FromGregorian(2022, time.August, 1))
	fmt.Println(chapter)
	// Output: Isaiah 47
}
