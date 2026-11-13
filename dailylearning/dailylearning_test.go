package dailylearning

import (
	"strconv"
	"sync"
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

func (m mockEvent) Date() hdate.HDate           { return m.date }
func (m mockEvent) Render(locale string) string { return m.desc }
func (m mockEvent) Flags() event.Flags          { return event.DailyLearning }
func (m mockEvent) Emoji() string               { return "" }
func (m mockEvent) Basename() string            { return m.desc }
func (m mockEvent) Desc() string                { return m.desc }
func (m mockEvent) URL() string                 { return "" }
func (m mockEvent) Categories() []string        { return []string{"learning"} }

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
	AddCalendarFrom("TestCal", mockFn, startDate)
	assert.True(t, Has("TestCal"))
	assert.True(t, Has("testcal")) // case insensitivity

	// Test Lookup
	hd := hdate.New(5780, hdate.Cheshvan, 1)
	ev := Lookup("TestCal", hd, false)
	assert.NotNil(t, ev)
	assert.Equal(t, "Mock Learning", ev.Render("en"))
	assert.Equal(t, hd, ev.Date())

	// Test Lookup for missing calendar
	assert.Nil(t, Lookup("MissingCal", hd, false))

	// Test Calendars
	AddCalendar("AlphaCal", mockFn)
	cals := Calendars()
	assert.Equal(t, []string{"alphacal", "testcal"}, cals)

	// Test StartDate
	sd, ok := StartDate("TestCal")
	assert.True(t, ok)
	assert.Equal(t, startDate, sd)

	// Test StartDate for calendar registered without start date
	sd2, ok2 := StartDate("AlphaCal")
	assert.False(t, ok2)
	assert.Equal(t, hdate.HDate{}, sd2)

	// Test StartDate for unregistered calendar
	sd3, ok3 := StartDate("Unregistered")
	assert.False(t, ok3)
	assert.Equal(t, hdate.HDate{}, sd3)
}

// The registry may be written while it is read. Run with -race.
func TestRegistryConcurrentUse(t *testing.T) {
	fn := func(hd hdate.HDate, il bool) event.CalEvent {
		return mockEvent{date: hd, desc: "Concurrent"}
	}
	hd := hdate.New(5786, hdate.Tishrei, 1)
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(2)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				AddCalendarFrom("concurrent"+strconv.Itoa(i), fn, hd)
				AddCalendar("concurrentPlain"+strconv.Itoa(i), fn)
			}
		}(i)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				_ = Lookup("concurrent"+strconv.Itoa(i), hd, false)
				_ = Has("concurrent0")
				_ = Calendars()
				_, _ = StartDate("concurrent" + strconv.Itoa(i))
			}
		}(i)
	}
	wg.Wait()
	assert.True(t, Has("Concurrent7"))
	start, ok := StartDate("concurrent7")
	assert.True(t, ok)
	assert.Equal(t, hd, start)
	_, ok = StartDate("concurrentPlain7")
	assert.False(t, ok)
}
