package dailylearning

import (
	"testing"

	"github.com/hebcal/hdate"
	"github.com/hebcal/hebcal-go/event"
	"github.com/stretchr/testify/assert"
)

// mockEvent is a simple implementation of event.CalEvent for testing
type mockEvent struct {
	date hdate.HDate
	desc string
}

func (m mockEvent) GetDate() hdate.HDate        { return m.date }
func (m mockEvent) Render(locale string) string { return m.desc }
func (m mockEvent) GetFlags() event.HolidayFlags { return event.DAILY_LEARNING }
func (m mockEvent) GetEmoji() string            { return "" }
func (m mockEvent) Basename() string            { return m.desc }
func (m mockEvent) GetCategories() []string     { return []string{"learning"} }

func TestDailyLearningRegistry(t *testing.T) {
	// Backup and reset calendars map
	oldCalendars := calendars
	defer func() { calendars = oldCalendars }()
	calendars = make(map[string]calendar)

	// Define a simple calendar function
	mockFn := func(hd hdate.HDate, il bool) event.CalEvent {
		return mockEvent{date: hd, desc: "Mock Learning"}
	}

	startDate := hdate.New(5780, hdate.Tishrei, 1)

	// Test Has before addition
	assert.False(t, Has("TestCal"))

	// Test AddCalendar with start date
	AddCalendar("TestCal", mockFn, startDate)
	assert.True(t, Has("TestCal"))
	assert.True(t, Has("testcal")) // case insensitivity

	// Test Lookup
	hd := hdate.New(5780, hdate.Cheshvan, 1)
	ev := Lookup("TestCal", hd, false)
	assert.NotNil(t, ev)
	assert.Equal(t, "Mock Learning", ev.Render("en"))
	assert.Equal(t, hd, ev.GetDate())

	// Test Lookup for missing calendar
	assert.Nil(t, Lookup("MissingCal", hd, false))

	// Test GetCalendars
	AddCalendar("AlphaCal", mockFn)
	cals := GetCalendars()
	assert.Equal(t, []string{"alphacal", "testcal"}, cals)

	// Test GetStartDate
	sd, ok := GetStartDate("TestCal")
	assert.True(t, ok)
	assert.Equal(t, startDate, sd)

	// Test GetStartDate for calendar registered without start date
	sd2, ok2 := GetStartDate("AlphaCal")
	assert.False(t, ok2)
	assert.Equal(t, hdate.HDate{}, sd2)

	// Test GetStartDate for unregistered calendar
	sd3, ok3 := GetStartDate("Unregistered")
	assert.False(t, ok3)
	assert.Equal(t, hdate.HDate{}, sd3)
}
