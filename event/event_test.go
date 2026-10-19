package event_test

import (
	"testing"

	"github.com/hebcal/hdate"
	"github.com/hebcal/hebcal-go/event"
	"github.com/hebcal/hebcal-go/sedra"
	"github.com/stretchr/testify/assert"
)

func TestHolidayEvent_Basename(t *testing.T) {
	ev := event.HolidayEvent{Desc: "Sukkot III (CH''M)"}
	assert.Equal(t, "Sukkot", ev.Basename())
	ev = event.HolidayEvent{Desc: "Chanukah: 1 Candle"}
	assert.Equal(t, "Chanukah", ev.Basename())
	ev = event.HolidayEvent{Desc: "Chanukah: 7 Candles"}
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
	ykk := event.HolidayEvent{Date: hdate.New(5786, hdate.Kislev, 29), Desc: "Yom Kippur Katan Tevet",
		Flags: event.MINOR_FAST | event.YOM_KIPPUR_KATAN}
	assert.Equal(t, "Yom Kippur Katan Tevet", ykk.Render("en"))
	assert.Equal(t, "Yom Kippur Katan", event.RenderBrief(ykk, "en"))
	assert.Equal(t, "יוֹם כִּפּוּר קָטָן", event.RenderBrief(ykk, "he"))
	// an event with no brief form renders in full
	purim := event.HolidayEvent{Date: hdate.New(5786, hdate.Adar1, 14), Desc: "Purim", Flags: event.MINOR_HOLIDAY}
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
	ev := event.HolidayEvent{
		Date:  hd,
		Desc:  "Yom Kippur Katan Sivan",
		Flags: event.MINOR_FAST | event.YOM_KIPPUR_KATAN}
	assert.Equal(t, "Yom Kippur Katan Sivan", ev.Basename())
	assert.Equal(t, "Yom Kippur Katan Sivan", ev.Render("en"))
	assert.Equal(t, "יוֹם כִּפּוּר קָטָן סִיוָן", ev.Render("he"))
	assert.Equal(t, "День Раскаяния Катан Сиван", ev.Render("ru"))
	assert.Equal(t, "יום כיפור קטן סיון", ev.Render("he-x-NoNikud"))
}

func TestRoshHashanaLocale(t *testing.T) {
	hd := hdate.New(5771, hdate.Tishrei, 1)
	ev := event.HolidayEvent{
		Date:  hd,
		Desc:  "Rosh Hashana 5771",
		Flags: event.CHAG | event.LIGHT_CANDLES_TZEIS,
		Emoji: "🍏🍯"}
	assert.Equal(t, "Rosh Hashana", ev.Basename())
	assert.Equal(t, "Rosh Hashana 5771", ev.Render("en"))
	assert.Equal(t, "רֹאשׁ הַשָּׁנָה תשע״א", ev.Render("he"))
	assert.Equal(t, "ראש השנה תשע״א", ev.Render("he-x-NoNikud"))
	assert.Equal(t, "Рош-А-Шана 5771", ev.Render("ru"))
}
