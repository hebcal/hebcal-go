package event_test

import (
	"testing"

	"github.com/hebcal/hdate"
	"github.com/hebcal/hebcal-go/event"
	"github.com/hebcal/hebcal-go/molad"
	"github.com/stretchr/testify/assert"
)

func TestReformatTimeStr(t *testing.T) {
	tests := []struct {
		timeStr, suffix, cc string
		clock               event.ClockFormat
		want                string
	}{
		{"20:30", "pm", "US", event.ClockAuto, "8:30pm"},
		{"05:19", "pm", "US", event.ClockAuto, "5:19am"},
		{"00:05", "pm", "US", event.ClockAuto, "12:05am"},
		{"12:00", "pm", "US", event.ClockAuto, "12:00pm"},
		{"20:30", "p", "CA", event.ClockAuto, "8:30p"},
		{"05:19", "p", "CA", event.ClockAuto, "5:19a"},
		{"05:19", " P.M.", "US", event.ClockAuto, "5:19 A.M."},
		{"20:30", "pm", "", event.ClockAuto, "8:30pm"},
		{"20:30", "pm", "IL", event.ClockAuto, "20:30"},
		{"05:19", "pm", "GB", event.ClockAuto, "05:19"},
		{"20:30", "pm", "US", event.Clock24, "20:30"},
		{"05:19", "pm", "IL", event.Clock12, "5:19am"},
		{"20:30", "", "US", event.ClockAuto, "8:30"},
		{"05:19", "", "US", event.ClockAuto, "5:19"},
		{"00:05", "", "US", event.ClockAuto, "00:05"},
		{"bogus", "pm", "US", event.ClockAuto, "bogus"},
	}
	for _, tt := range tests {
		got := event.ReformatTimeStr(tt.timeStr, tt.suffix, tt.cc, tt.clock)
		assert.Equal(t, tt.want, got, "%q %q %q", tt.timeStr, tt.suffix, tt.cc)
	}
}

func TestUses12HourClock(t *testing.T) {
	assert.True(t, event.Uses12HourClock("US"))
	assert.True(t, event.Uses12HourClock("AU"))
	assert.False(t, event.Uses12HourClock("IL"))
	assert.False(t, event.Uses12HourClock("GB"))
	assert.False(t, event.Uses12HourClock(""))
}

func TestMoladEventClock(t *testing.T) {
	hd := hdate.New(5786, hdate.Sivan, 29)
	m := molad.New(5786, hdate.Tamuz)
	assert.Contains(t, event.NewMoladEvent(hd, m, "Tammuz", "US", event.ClockAuto).Render("en"),
		"Molad Tammuz: Monday, 6:46am")
	assert.Contains(t, event.NewMoladEvent(hd, m, "Tammuz", "US", event.Clock24).Render("en"),
		"Molad Tammuz: Monday, 6:46 and")
	assert.Contains(t, event.NewMoladEvent(hd, m, "Tammuz", "GB", event.Clock12).Render("en"),
		"Molad Tammuz: Monday, 6:46am")
}
