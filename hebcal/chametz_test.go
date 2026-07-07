package hebcal_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hebcal/hebcal-go/hebcal"
	"github.com/hebcal/hebcal-go/zmanim"
	"github.com/stretchr/testify/assert"
)

// chametzEvents renders every event whose (locale) description mentions chametz,
// prefixed with its weekday, Gregorian date, and emoji.
func chametzEvents(t *testing.T, gregYear int, locale string) []string {
	t.Helper()
	loc := zmanim.LookupCity("Chicago")
	opts := &hebcal.CalOptions{Year: gregYear, Location: loc, CandleLighting: true}
	events, err := hebcal.HebrewCalendar(opts)
	assert.NoError(t, err)
	var out []string
	for _, ev := range events {
		r := ev.Render(locale)
		if strings.Contains(strings.ToLower(r), "chametz") || strings.Contains(r, "חָמֵץ") {
			hd := ev.GetDate()
			y, m, d := hd.Greg()
			dow := hd.Weekday().String()[:3]
			out = append(out, fmt.Sprintf("%s %04d-%02d-%02d %s %s", dow, y, m, d, ev.GetEmoji(), r))
		}
	}
	return out
}

// Normal year: Erev Pesach 5782 is Friday 2022-04-15. Both sof zman achilat
// chametz and sof zman biur chametz are emitted that morning, biur one halachic
// hour after achilat.
func TestErevPesachChametz(t *testing.T) {
	assert := assert.New(t)
	expected := []string{
		"Fri 2022-04-15 🍞 Finish eating chametz: 10:37",
		"Fri 2022-04-15 🔥 Biur Chametz: 11:44",
	}
	assert.Equal(expected, chametzEvents(t, 2022, "en"))
}

// Erev Pesach on Shabbat: 14 Nisan 5781 is Saturday 2021-03-27. Chametz cannot
// be burned on Shabbat, so Biur Chametz is moved to Friday 2021-03-26, while sof
// zman achilat chametz stays on Shabbat.
func TestErevPesachChametzOnShabbat(t *testing.T) {
	assert := assert.New(t)
	expected := []string{
		"Fri 2021-03-26 🔥 Biur Chametz: 11:54",
		"Sat 2021-03-27 🍞 Finish eating chametz: 10:51",
	}
	assert.Equal(expected, chametzEvents(t, 2021, "en"))
}

// The events carry their Hebrew translations.
func TestErevPesachChametzHebrew(t *testing.T) {
	assert := assert.New(t)
	expected := []string{
		"Fri 2022-04-15 🍞 סוֹף זְמַן אֲכִילַת חָמֵץ: 10:37",
		"Fri 2022-04-15 🔥 בִּעוּר חָמֵץ: 11:44",
	}
	assert.Equal(expected, chametzEvents(t, 2022, "he"))
}

// Without candle-lighting, no chametz events are emitted.
func TestErevPesachChametzRequiresCandleLighting(t *testing.T) {
	assert := assert.New(t)
	loc := zmanim.LookupCity("Chicago")
	opts := &hebcal.CalOptions{Year: 2022, Location: loc}
	events, err := hebcal.HebrewCalendar(opts)
	assert.NoError(err)
	for _, ev := range events {
		assert.NotContains(strings.ToLower(ev.Render("en")), "chametz")
	}
}
