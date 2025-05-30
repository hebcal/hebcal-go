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
	assert.Equal(t, "3 Sziván 5781", ev.Render("hu"))
	assert.Equal(t, "ג׳ סיון תשפ״א", ev.Render("he-x-NoNikud"))
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
	assert.Equal(t, "פָּרָשַׁת מַּטּוֹת־מַסְעֵי", ev.Render("he"))
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
	assert.Equal(t, "יום כפור קטן סיון", ev.Render("he-x-NoNikud"))
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
