package event_test

import (
	"testing"

	"github.com/hebcal/hdate"
	"github.com/hebcal/hebcal-go/event"
	"github.com/hebcal/hebcal-go/molad"
	"github.com/stretchr/testify/assert"
)

func TestReformatTimeStr(t *testing.T) {
	yes, no := true, false
	tests := []struct {
		timeStr, suffix, cc string
		hour12              *bool
		want                string
	}{
		{"20:30", "pm", "US", nil, "8:30pm"},
		{"05:19", "pm", "US", nil, "5:19am"},
		{"00:05", "pm", "US", nil, "12:05am"},
		{"12:00", "pm", "US", nil, "12:00pm"},
		{"20:30", "p", "CA", nil, "8:30p"},
		{"05:19", "p", "CA", nil, "5:19a"},
		{"05:19", " P.M.", "US", nil, "5:19 A.M."},
		{"20:30", "pm", "", nil, "8:30pm"},
		{"20:30", "pm", "IL", nil, "20:30"},
		{"05:19", "pm", "GB", nil, "05:19"},
		{"20:30", "pm", "US", &no, "20:30"},
		{"05:19", "pm", "IL", &yes, "5:19am"},
		{"20:30", "", "US", nil, "8:30"},
		{"05:19", "", "US", nil, "5:19"},
		{"00:05", "", "US", nil, "00:05"},
		{"bogus", "pm", "US", nil, "bogus"},
	}
	for _, tt := range tests {
		got := event.ReformatTimeStr(tt.timeStr, tt.suffix, tt.cc, tt.hour12)
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

func TestMoladEventHour12(t *testing.T) {
	hd := hdate.New(5786, hdate.Sivan, 29)
	m := molad.New(5786, hdate.Tamuz)
	yes, no := true, false
	assert.Contains(t, event.NewMoladEventHour12(hd, m, "Tammuz", "US", nil).Render("en"),
		"Molad Tammuz: Monday, 6:46am")
	assert.Contains(t, event.NewMoladEventHour12(hd, m, "Tammuz", "US", &no).Render("en"),
		"Molad Tammuz: Monday, 6:46 and")
	assert.Contains(t, event.NewMoladEventHour12(hd, m, "Tammuz", "GB", &yes).Render("en"),
		"Molad Tammuz: Monday, 6:46am")
}
