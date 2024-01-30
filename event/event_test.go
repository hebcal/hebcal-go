package event_test

import (
	"testing"

	"github.com/hebcal/hdate"
	"github.com/MaxBGreenberg/hebcal-go/event"
	"github.com/MaxBGreenberg/hebcal-go/sedra"
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
