package hebcal

import (
	"testing"
	"time"

	"github.com/hebcal/hdate"
	"github.com/hebcal/hebcal-go/dailylearning"
	"github.com/hebcal/hebcal-go/event"
	"github.com/hebcal/hebcal-go/zmanim"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func renderDayByTimeOfDay(t *testing.T, year int, month time.Month, day int) []string {
	dailylearning.AddCalendar("testTimeOfDay", func(hd hdate.HDate, il bool) event.CalEvent {
		return event.HolidayEvent{Date: hd, Desc: "Daily learning", Flags: event.DAILY_LEARNING}
	})
	toronto := zmanim.NewLocation("Toronto", "CA", 43.70011, -79.4163, 0, "America/Toronto")
	hd := hdate.FromGregorian(year, month, day)
	events, err := HebrewCalendar(&CalOptions{
		Start:                hd,
		End:                  hd,
		Location:             &toronto,
		CandleLighting:       true,
		HavdalahMins:         42,
		Omer:                 true,
		Sedrot:               true,
		KiddushLevanaMaharil: true,
		DailyLearning:        []string{"testTimeOfDay"},
	})
	require.NoError(t, err)
	result := make([]string, len(events))
	for i, ev := range events {
		result[i] = ev.Render("en")
	}
	return result
}

func TestTimeOfDayLearningAfterMinorFast(t *testing.T) {
	assert.Equal(t, []string{
		"Fast begins: 5:54am",
		"Tzom Gedaliah",
		"Daily learning",
		"Fast ends: 7:28pm",
	}, renderDayByTimeOfDay(t, 2027, time.October, 4))
}

func TestTimeOfDayErevAfterLearning(t *testing.T) {
	assert.Equal(t, []string{
		"Latest Kiddush Levana: 6:06am",
		"Daily learning",
		"Erev Sukkot",
		"Candle lighting: 6:16pm",
	}, renderDayByTimeOfDay(t, 2027, time.October, 15))
}

func TestTimeOfDayErevPesach(t *testing.T) {
	assert.Equal(t, []string{
		"Fast begins: 5:12am",
		"Ta'anit Bechorot",
		"Finish eating chametz: 11:07am",
		"Biur Chametz: 12:13pm",
		"Daily learning",
		"Erev Pesach",
		"Candle lighting: 7:37pm",
	}, renderDayByTimeOfDay(t, 2028, time.April, 10))
}

func TestTimeOfDayOmer(t *testing.T) {
	assert.Equal(t, []string{
		"49th day of the Omer",
		"Daily learning",
		"Erev Shavuot",
		"Candle lighting: 8:33pm",
	}, renderDayByTimeOfDay(t, 2028, time.May, 30))
}

func TestTimeOfDayErevTishaBav(t *testing.T) {
	assert.Equal(t, []string{
		"Parashat Devarim",
		"Shabbat Chazon",
		"Daily learning",
		"Erev Tish'a B'Av",
		"Fast begins: 8:34pm",
		"Havdalah (42 min): 9:16pm",
	}, renderDayByTimeOfDay(t, 2022, time.August, 6))
}

func TestTimeOfDayExported(t *testing.T) {
	hd := hdate.New(5787, hdate.Tishrei, 14)
	assert.Equal(t, TimeOfDayLearning,
		TimeOfDay(event.HolidayEvent{Date: hd, Desc: "Psalms 1", Flags: event.DAILY_LEARNING}))
	assert.Equal(t, TimeOfDayEvening,
		TimeOfDay(event.HolidayEvent{Date: hd, Desc: "Erev Sukkot", Flags: event.EREV | event.LIGHT_CANDLES}))
	assert.Equal(t, TimeOfDayDay, TimeOfDay(event.UserEvent{Date: hd, Desc: "Birthday"}))
}
