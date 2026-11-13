package event_test

import (
	"testing"

	"github.com/hebcal/hdate"
	"github.com/hebcal/hebcal-go/event"
	"github.com/hebcal/hebcal-go/sedra"
	"github.com/stretchr/testify/assert"
)

func TestHolidayEvent_Basename(t *testing.T) {
	ev := event.NewHolidayEvent(hdate.HDate{}, "Sukkot III (CH''M)", 0, "")
	assert.Equal(t, "Sukkot", ev.Basename())
	ev = event.NewHolidayEvent(hdate.HDate{}, "Chanukah: 1 Candle", 0, "")
	assert.Equal(t, "Chanukah", ev.Basename())
	ev = event.NewHolidayEvent(hdate.HDate{}, "Chanukah: 7 Candles", 0, "")
	assert.Equal(t, "Chanukah", ev.Basename())
}

func TestHebrewDateEvent(t *testing.T) {
	hd := hdate.New(5781, hdate.Sivan, 3)
	ev := event.NewHebrewDateEvent(hd)
	assert.Equal(t, "3 Sivan 5781", ev.Basename())
	assert.Equal(t, "3rd of Sivan, 5781", ev.Render("en"))
	assert.Equal(t, "ג׳ סִיוָן תשפ״א", ev.Render("he"))
	assert.Equal(t, "3. Sziván, 5781", ev.Render("hu"))
	assert.Equal(t, "3º Sivan, 5781", ev.Render("es"))
	assert.Equal(t, "ג׳ סיון תשפ״א", ev.Render("he-x-NoNikud"))
	assert.Equal(t, "3rd of Sivan", event.RenderBrief(ev, "en"))
	assert.Equal(t, "ג׳ סִיוָן", event.RenderBrief(ev, "he"))
	assert.Equal(t, "3. Сиван", event.RenderBrief(ev, "ru"))
	// the first day of the year keeps its year
	ev = event.NewHebrewDateEvent(hdate.New(5787, hdate.Tishrei, 1))
	assert.Equal(t, "1st of Tishrei, 5787", event.RenderBrief(ev, "en"))
	// typographic apostrophe, one-m Tamuz, catalog month names
	assert.Equal(t, "1st of Sh’vat, 5786", event.NewHebrewDateEvent(hdate.New(5786, hdate.Shvat, 1)).Render("en"))
	assert.Equal(t, "17th of Tamuz", event.RenderBrief(event.NewHebrewDateEvent(hdate.New(5786, hdate.Tamuz, 17)), "en"))
	assert.Equal(t, "10th of Teves", event.RenderBrief(event.NewHebrewDateEvent(hdate.New(5786, hdate.Tevet, 10)), "ashkenazi"))
}

func TestRenderBrief(t *testing.T) {
	ykk := event.NewHolidayEvent(
		hdate.New(5786, hdate.Kislev, 29),
		"Yom Kippur Katan Tevet",
		event.MinorFast|event.YomKippurKatan,
		"")
	assert.Equal(t, "Yom Kippur Katan Tevet", ykk.Render("en"))
	assert.Equal(t, "Yom Kippur Katan", event.RenderBrief(ykk, "en"))
	assert.Equal(t, "יוֹם כִּפּוּר קָטָן", event.RenderBrief(ykk, "he"))
	// an event with no brief form renders in full
	purim := event.NewHolidayEvent(hdate.New(5786, hdate.Adar1, 14), "Purim", event.MinorHoliday, "")
	assert.Equal(t, "Purim", event.RenderBrief(purim, "en"))
	parsha := event.NewParshaEvent(hdate.New(5786, hdate.Cheshvan, 5),
		sedra.Parsha{Name: []string{"Noach"}, Num: []int{2}}, false)
	assert.Equal(t, "Parashat Noach", event.RenderBrief(parsha, "en"))
}

func TestParshaEvent_Render(t *testing.T) {
	parsha := sedra.Parsha{
		Name: []string{"Matot", "Masei"},
		Num:  []int{42, 43},
		Chag: false,
	}
	hd := hdate.New(5783, hdate.Tamuz, 26)
	ev := event.NewParshaEvent(hd, parsha, false)
	assert.Equal(t, "Parashat Matot-Masei", ev.Render("en"))
	assert.Equal(t, "Глава Матот-Масей", ev.Render("ru"))
	assert.Equal(t, "פָּרָשַׁת מַטּוֹת־מַסְעֵי", ev.Render("he"))
}

func TestYomKippurKatanEvent(t *testing.T) {
	hd := hdate.New(5771, hdate.Tevet, 29)
	ev := event.NewHolidayEvent(
		hd,
		"Yom Kippur Katan Sivan",
		event.MinorFast|event.YomKippurKatan,
		"")
	assert.Equal(t, "Yom Kippur Katan Sivan", ev.Basename())
	assert.Equal(t, "Yom Kippur Katan Sivan", ev.Render("en"))
	assert.Equal(t, "יוֹם כִּפּוּר קָטָן סִיוָן", ev.Render("he"))
	assert.Equal(t, "День Раскаяния Катан Сиван", ev.Render("ru"))
	assert.Equal(t, "יום כיפור קטן סיון", ev.Render("he-x-NoNikud"))
}

func TestRoshHashanaLocale(t *testing.T) {
	hd := hdate.New(5771, hdate.Tishrei, 1)
	ev := event.NewHolidayEvent(
		hd,
		"Rosh Hashana 5771",
		event.Chag|event.LightCandlesTzeis,
		"🍏🍯")
	assert.Equal(t, "Rosh Hashana", ev.Basename())
	assert.Equal(t, "Rosh Hashana 5771", ev.Render("en"))
	assert.Equal(t, "רֹאשׁ הַשָּׁנָה תשע״א", ev.Render("he"))
	assert.Equal(t, "ראש השנה תשע״א", ev.Render("he-x-NoNikud"))
	assert.Equal(t, "Рош-А-Шана 5771", ev.Render("ru"))
}

// A caller can build a holiday whose description is shorter than its flags
// suggest; rendering it must not slice past the end of the string.
func TestRenderShortDescriptions(t *testing.T) {
	hd := hdate.New(5786, hdate.Shvat, 1)
	for _, tc := range []struct {
		desc  string
		flags event.Flags
	}{
		{"RC", event.RoshChodesh},
		{"", event.RoshChodesh},
		{"YKK", event.YomKippurKatan | event.MinorFast},
		{"", event.YomKippurKatan},
	} {
		ev := event.NewHolidayEvent(hd, tc.desc, tc.flags, "")
		assert.NotPanics(t, func() {
			_ = ev.Render("en")
			_ = ev.Render("he")
			_ = ev.RenderBrief("en")
			_ = ev.Memo("en")
			_ = ev.URL()
			_ = ev.Basename()
		}, "%q", tc.desc)
	}
	rc := event.NewHolidayEvent(hd, "Rosh Chodesh Sh'vat", event.RoshChodesh, "")
	assert.Equal(t, "Rosh Chodesh Sh’vat", rc.Render("en"))
}
