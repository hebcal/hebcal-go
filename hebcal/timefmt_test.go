package hebcal_test

import (
	"testing"

	"github.com/hebcal/hebcal-go/hebcal"
	"github.com/hebcal/hebcal-go/zmanim"
	"github.com/stretchr/testify/assert"
)

func TestReformatTimeStrOptions(t *testing.T) {
	assert := assert.New(t)
	chicago := zmanim.LookupCity("Chicago")
	jerusalem := zmanim.LookupCity("Jerusalem")
	assert.Equal("8:30pm", hebcal.ReformatTimeStr("20:30", "pm", nil))
	assert.Equal("8:30pm", hebcal.ReformatTimeStr("20:30", "pm", &hebcal.CalOptions{}))
	assert.Equal("20:30", hebcal.ReformatTimeStr("20:30", "pm", &hebcal.CalOptions{IL: true}))
	assert.Equal("5:19am", hebcal.ReformatTimeStr("05:19", "pm", &hebcal.CalOptions{Location: chicago}))
	assert.Equal("05:19", hebcal.ReformatTimeStr("05:19", "pm", &hebcal.CalOptions{Location: jerusalem}))
	assert.Equal("05:19", hebcal.ReformatTimeStr("05:19", "pm", &hebcal.CalOptions{Location: chicago, Hour24: true}))
	assert.Equal("5:19am", hebcal.ReformatTimeStr("05:19", "pm", &hebcal.CalOptions{Location: jerusalem, Hour12: true}))
	assert.Equal("05:19", hebcal.ReformatTimeStr("05:19", "pm",
		&hebcal.CalOptions{Location: chicago, Hour12: true, Hour24: true}))
}

// Morning and evening timed events are distinguishable in a 12-hour country,
// and each override applies to TimedEvent.Render().
func TestTimedEventClockFormat(t *testing.T) {
	assert := assert.New(t)
	chicago := zmanim.LookupCity("Chicago")
	jerusalem := zmanim.LookupCity("Jerusalem")
	assert.Equal("Fast begins: 4:42am", findOnDate(t,
		&hebcal.CalOptions{Year: 2022, Month: 4, Location: chicago, CandleLighting: true},
		2022, 4, 15, "Fast begins"))
	assert.Equal("Candle lighting: 7:13pm", findOnDate(t,
		&hebcal.CalOptions{Year: 2022, Month: 4, Location: chicago, CandleLighting: true},
		2022, 4, 15, "Candle lighting"))
	assert.Equal("Fast begins: 04:42", findOnDate(t,
		&hebcal.CalOptions{Year: 2022, Month: 4, Location: chicago, CandleLighting: true, Hour24: true},
		2022, 4, 15, "Fast begins"))
	assert.Equal("Candle lighting: 4:00pm", findOnDate(t,
		&hebcal.CalOptions{Year: 2022, Month: 12, Location: jerusalem, CandleLighting: true, Hour12: true},
		2022, 12, 23, "Candle lighting"))
}
