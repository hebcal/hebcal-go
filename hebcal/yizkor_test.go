package hebcal

import (
	"testing"

	"github.com/hebcal/hdate"
	"github.com/hebcal/hebcal-go/event"
	"github.com/stretchr/testify/assert"
)

func TestYizkor(t *testing.T) {
	for _, tc := range []struct {
		il   bool
		want []string
	}{
		{false, []string{"Yom Kippur", "Shmini Atzeret", "Pesach VIII", "Shavuot II"}},
		{true, []string{"Yom Kippur", "Shmini Atzeret", "Pesach VII", "Shavuot"}},
	} {
		events, err := HebrewCalendar(&CalOptions{Year: 5786, IsHebrewYear: true, IL: tc.il,
			NoHolidays: true, Yizkor: true})
		assert.NoError(t, err)
		var got []string
		for _, ev := range events {
			y, ok := ev.(YizkorEvent)
			if !ok {
				t.Fatalf("unexpected event %s", ev.Render("en"))
			}
			assert.Equal(t, "Yizkor", y.Render("en"))
			assert.Equal(t, "יִזְכּוֹר", y.Render("he"))
			assert.Equal(t, event.YIZKOR, y.GetFlags())
			assert.Equal(t, []string{"yizkor"}, y.GetCategories())
			got = append(got, y.LinkedEvent.(event.HolidayEvent).Desc)
		}
		assert.Equal(t, tc.want, got, "il=%v", tc.il)
	}
}

func TestYizkorFollowsItsHoliday(t *testing.T) {
	events, err := HebrewCalendar(&CalOptions{
		Start: hdate.New(5786, hdate.Tishrei, 10), End: hdate.New(5786, hdate.Tishrei, 10),
		Mask: event.MAJOR_FAST | event.YIZKOR,
	})
	assert.NoError(t, err)
	var got []string
	for _, ev := range events {
		got = append(got, ev.Render("en"))
	}
	assert.Equal(t, []string{"Yom Kippur", "Yizkor"}, got)
}
