package event_test

import (
	"fmt"
	"math/bits"
	"testing"

	"github.com/hebcal/hebcal-go/event"
	"github.com/stretchr/testify/assert"
)

func TestHolidayFlagsHas(t *testing.T) {
	assert := assert.New(t)
	f := event.Chag | event.LightCandles
	assert.True(f.Has(event.Chag))
	assert.True(f.Has(event.Chag | event.LightCandles))
	assert.False(f.Has(event.Chag | event.Erev))
	assert.False(f.Has(event.Erev))
	assert.True(f.Has(0))
}

func TestHolidayFlagsHasAny(t *testing.T) {
	assert := assert.New(t)
	f := event.Chag | event.LightCandles
	assert.True(f.HasAny(event.Chag))
	assert.True(f.HasAny(event.Chag | event.Erev))
	assert.False(f.HasAny(event.Erev | event.MinorFast))
	assert.False(f.HasAny(0))
}

func TestHolidayFlagsWithWithout(t *testing.T) {
	assert := assert.New(t)
	var f event.Flags
	f = f.With(event.Chag | event.Erev)
	assert.Equal(event.Chag|event.Erev, f)
	f = f.Without(event.Chag)
	assert.Equal(event.Erev, f)
	f = f.Without(event.Molad) // clearing an unset flag is a no-op
	assert.Equal(event.Erev, f)
}

func TestFlagsString(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("0", event.Flags(0).String())
	assert.Equal("Chag", event.Chag.String())
	assert.Equal("Chag|LightCandles", (event.LightCandles | event.Chag).String())
	assert.Equal("DailyLearning", event.DailyLearning.String())
	assert.Equal("MinorFast|BeHaB", (event.MinorFast | event.BeHaB).String())
	assert.Equal("Erev|0x80000000", (event.Erev | 1<<31).String())
	assert.Equal("MajorFast", fmt.Sprintf("%v", event.MajorFast))
}

// Every flag through Yizkor (currently the last one) must have a name.
func TestFlagsStringCoversAllFlags(t *testing.T) {
	last := bits.TrailingZeros32(uint32(event.Yizkor))
	for i := 0; i <= last; i++ {
		s := (event.Flags(1) << i).String()
		assert.NotContains(t, s, "0x", "bit %d has no name", i)
	}
}
