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
	f := event.CHAG | event.LIGHT_CANDLES
	assert.True(f.Has(event.CHAG))
	assert.True(f.Has(event.CHAG | event.LIGHT_CANDLES))
	assert.False(f.Has(event.CHAG | event.EREV))
	assert.False(f.Has(event.EREV))
	assert.True(f.Has(0))
}

func TestHolidayFlagsHasAny(t *testing.T) {
	assert := assert.New(t)
	f := event.CHAG | event.LIGHT_CANDLES
	assert.True(f.HasAny(event.CHAG))
	assert.True(f.HasAny(event.CHAG | event.EREV))
	assert.False(f.HasAny(event.EREV | event.MINOR_FAST))
	assert.False(f.HasAny(0))
}

func TestHolidayFlagsWithWithout(t *testing.T) {
	assert := assert.New(t)
	var f event.HolidayFlags
	f = f.With(event.CHAG | event.EREV)
	assert.Equal(event.CHAG|event.EREV, f)
	f = f.Without(event.CHAG)
	assert.Equal(event.EREV, f)
	f = f.Without(event.MOLAD) // clearing an unset flag is a no-op
	assert.Equal(event.EREV, f)
}

func TestHolidayFlagsString(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("0", event.HolidayFlags(0).String())
	assert.Equal("CHAG", event.CHAG.String())
	assert.Equal("CHAG|LIGHT_CANDLES", (event.LIGHT_CANDLES | event.CHAG).String())
	assert.Equal("DAILY_LEARNING", event.DAILY_LEARNING.String())
	assert.Equal("EREV|0x80000000", (event.EREV | 1<<31).String())
	assert.Equal("MAJOR_FAST", fmt.Sprintf("%v", event.MAJOR_FAST))
}

// Every flag through DAILY_LEARNING (currently the last one) must have a name.
func TestHolidayFlagsStringCoversAllFlags(t *testing.T) {
	last := bits.TrailingZeros32(uint32(event.DAILY_LEARNING))
	for i := 0; i <= last; i++ {
		s := (event.HolidayFlags(1) << i).String()
		assert.NotContains(t, s, "0x", "bit %d has no name", i)
	}
}
