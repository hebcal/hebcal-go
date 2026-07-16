package event_test

import (
	"testing"

	"github.com/hebcal/hdate"
	"github.com/hebcal/hebcal-go/event"
	"github.com/hebcal/hebcal-go/molad"
	"github.com/hebcal/hebcal-go/sedra"
	"github.com/stretchr/testify/assert"
)

func TestHebrewDateEvent_Extra(t *testing.T) {
	hd := hdate.New(5781, hdate.Sivan, 3)
	ev := event.NewHebrewDateEvent(hd)
	assert.Equal(t, hd, ev.GetDate())
	assert.Equal(t, event.HEBREW_DATE, ev.GetFlags())
	assert.Equal(t, "", ev.GetEmoji())
	assert.Equal(t, []string{"hebdate"}, ev.GetCategories())
}

func TestHolidayEvent_Extra(t *testing.T) {
	hd := hdate.New(5781, hdate.Sivan, 6)
	ev := event.HolidayEvent{
		Date:  hd,
		Desc:  "Shavuot I",
		Flags: event.CHAG,
	}
	assert.Equal(t, hd, ev.GetDate())
	assert.Equal(t, event.CHAG, ev.GetFlags())
	assert.Equal(t, "✡️", ev.GetEmoji())

	ev2 := event.HolidayEvent{
		Date:  hd,
		Desc:  "Shabbat Shekalim",
		Flags: event.SPECIAL_SHABBAT,
	}
	assert.Equal(t, "🕍", ev2.GetEmoji())

	ev3 := event.HolidayEvent{
		Date:  hd,
		Desc:  "Rosh Chodesh Tamuz",
		Flags: event.ROSH_CHODESH,
	}
	assert.Equal(t, "🌒", ev3.GetEmoji())

	ev4 := event.HolidayEvent{
		Date:  hd,
		Desc:  "My Custom Holiday",
		Flags: event.CHAG,
		Emoji: "🎉",
	}
	assert.Equal(t, "🎉", ev4.GetEmoji())

	// Test Categories
	evChol := event.HolidayEvent{
		Date:          hd,
		Desc:          "Pesach III (CH''M)",
		Flags:         event.CHOL_HAMOED,
		CholHaMoedDay: 3,
	}
	assert.Equal(t, []string{"holiday", "major", "cholhamoed"}, evChol.GetCategories())

	evMinor := event.HolidayEvent{
		Date:  hd,
		Desc:  "Lag BaOmer",
		Flags: 0,
	}
	assert.Equal(t, []string{"holiday", "minor"}, evMinor.GetCategories())

	evMajorFall := event.HolidayEvent{
		Date:  hd,
		Desc:  "Some Unknown Major",
		Flags: 0,
	}
	assert.Equal(t, []string{"holiday", "major"}, evMajorFall.GetCategories())

	// Test Render for Rosh Chodesh
	evRoshChodesh := event.HolidayEvent{
		Date:  hd,
		Desc:  "Rosh Chodesh Sh'vat",
		Flags: event.ROSH_CHODESH,
	}
	assert.Equal(t, "Rosh Chodesh Sh'vat", evRoshChodesh.Render("en"))
	assert.Equal(t, "ראש חודש שבט", evRoshChodesh.Render("he-x-nonikud"))

	// Test Render for Shabbat Mevarchim
	evMevarchim := event.HolidayEvent{
		Date:  hd,
		Desc:  "Shabbat Mevarchim Chodesh Kislev",
		Flags: event.SHABBAT_MEVARCHIM,
	}
	assert.Equal(t, "Shabbat Mevarchim Chodesh Kislev", evMevarchim.Render("en"))

	// Test Render for Rosh Hashana Hebrew
	evRoshHashana := event.HolidayEvent{
		Date:  hdate.New(5781, hdate.Tishrei, 1),
		Desc:  "Rosh Hashana 5781",
		Flags: event.CHAG,
	}
	assert.Equal(t, "רֹאשׁ הַשָּׁנָה תשפ״א", evRoshHashana.Render("he"))
	assert.Equal(t, "Rosh Hashana 5781", evRoshHashana.Render("en"))

	// Test Render for Yom Kippur Katan
	evYKK := event.HolidayEvent{
		Date:  hd,
		Desc:  "Yom Kippur Katan Sivan",
		Flags: event.MINOR_FAST | event.YOM_KIPPUR_KATAN,
	}
	assert.Equal(t, "Yom Kippur Katan Sivan", evYKK.Render("en"))

	// Test GetEmoji edge cases
	evMevarchimEmoji := event.HolidayEvent{
		Date:  hd,
		Desc:  "Shabbat Mevarchim Chodesh Kislev",
		Flags: event.SHABBAT_MEVARCHIM,
	}
	assert.Equal(t, "", evMevarchimEmoji.GetEmoji())

	evYKKEmoji := event.HolidayEvent{
		Date:  hd,
		Desc:  "Yom Kippur Katan Sivan",
		Flags: event.YOM_KIPPUR_KATAN | event.MINOR_FAST,
	}
	assert.Equal(t, "", evYKKEmoji.GetEmoji())
}

func TestMevarchimChodeshEvent(t *testing.T) {
	hd := hdate.New(5786, hdate.Sivan, 25)
	m := molad.New(5786, hdate.Tamuz)
	ev := event.NewMevarchimChodeshEvent(hd, "Tammuz", m)

	assert.Equal(t, hd, ev.GetDate())
	assert.Equal(t, event.SHABBAT_MEVARCHIM, ev.GetFlags())
	assert.Equal(t, "", ev.GetEmoji())
	assert.Equal(t, []string{"mevarchim"}, ev.GetCategories())
	assert.Equal(t, "Shabbat Mevarchim Chodesh Tammuz", ev.Basename())
	assert.Equal(t, "Shabbat Mevarchim Chodesh Tammuz", ev.Render("en"))
	assert.Equal(t, "Mevarchim Chodesh Tammuz", ev.RenderBrief("en"))

	// Test Brief render fallback with no spaces
	evBrief := event.NewMevarchimChodeshEvent(hd, "Tammuz", m)
	assert.Equal(t, "Mevarchim Chodesh Tammuz", evBrief.RenderBrief("en"))
}

func TestMoladEvent(t *testing.T) {
	hd := hdate.New(5786, hdate.Sivan, 29)
	m := molad.New(5786, hdate.Tamuz)

	ev := event.NewMoladEvent(hd, m, "Tammuz", "US")
	assert.Equal(t, hd, ev.GetDate())
	assert.Equal(t, event.MOLAD, ev.GetFlags())
	assert.Equal(t, "", ev.GetEmoji())
	assert.Equal(t, []string{"molad"}, ev.GetCategories())
	assert.Equal(t, "Molad Tammuz", ev.Basename())

	// Test renderings for different locales and country codes
	assert.Contains(t, ev.Render("en"), "Molad Tammuz: Monday, 6:46am and 16 chalakim")

	evGB := event.NewMoladEvent(hd, m, "Tammuz", "GB")
	assert.Contains(t, evGB.Render("en"), "Molad Tammuz: Monday, 6:46 and 16 chalakim")

	evDefaultCC := event.NewMoladEvent(hd, m, "Tammuz", "")
	assert.Contains(t, evDefaultCC.Render("en"), "Molad Tammuz: Monday, 6:46am and 16 chalakim")

	// Test hebrew rendering (various hours / periods of day)
	assert.Contains(t, ev.Render("he"), "בֹּקֶר")

	mAfternoon := m
	mAfternoon.Hours = 14
	evAfternoon := event.NewMoladEvent(hd, mAfternoon, "Tammuz", "US")
	assert.Contains(t, evAfternoon.Render("he"), "בַּצׇּהֳרַיִים")
	assert.Contains(t, evAfternoon.Render("he-x-nonikud"), "בצהריים")

	mEvening := m
	mEvening.Hours = 19
	evEvening := event.NewMoladEvent(hd, mEvening, "Tammuz", "US")
	assert.Contains(t, evEvening.Render("he"), "בָּעֶרֶב")

	mNight := m
	mNight.Hours = 23
	evNight := event.NewMoladEvent(hd, mNight, "Tammuz", "US")
	assert.Contains(t, evNight.Render("he"), "בַּלַּ֥יְלָה")

	mNight2 := m
	mNight2.Hours = 2
	evNight2 := event.NewMoladEvent(hd, mNight2, "Tammuz", "US")
	assert.Contains(t, evNight2.Render("he"), "בַּלַּ֥יְלָה")

	// 12-hour clock boundary cases: 0:00 (12:00am) and 12:00 (12:00pm)
	mMidnight := m
	mMidnight.Hours = 0
	mMidnight.Minutes = 0
	evMidnight := event.NewMoladEvent(hd, mMidnight, "Tammuz", "US")
	assert.Contains(t, evMidnight.Render("en"), "12:00am")

	mNoon := m
	mNoon.Hours = 12
	mNoon.Minutes = 0
	evNoon := event.NewMoladEvent(hd, mNoon, "Tammuz", "US")
	assert.Contains(t, evNoon.Render("en"), "12:00pm")

	mPM := m
	mPM.Hours = 13
	mPM.Minutes = 30
	evPM := event.NewMoladEvent(hd, mPM, "Tammuz", "US")
	assert.Contains(t, evPM.Render("en"), "1:30pm")
}

func TestUserEvent(t *testing.T) {
	hd := hdate.New(5780, hdate.Nisan, 15)
	ev := event.UserEvent{
		Date: hd,
		Desc: "My Event",
	}

	assert.Equal(t, hd, ev.GetDate())
	assert.Equal(t, event.USER_EVENT, ev.GetFlags())
	assert.Equal(t, "", ev.GetEmoji())
	assert.Equal(t, "My Event", ev.Basename())
	assert.Equal(t, "My Event", ev.Render("en"))
	assert.Equal(t, []string{"user"}, ev.GetCategories())
}

func TestParshaEvent_Extra(t *testing.T) {
	hd := hdate.New(5780, hdate.Nisan, 15)
	p := sedra.Parsha{
		Name: []string{"Noach"},
		Num:  []int{2},
		Chag: false,
	}
	ev := event.NewParshaEvent(hd, p, false)

	assert.Equal(t, hd, ev.GetDate())
	assert.Equal(t, event.PARSHA_HASHAVUA, ev.GetFlags())
	assert.Equal(t, "", ev.GetEmoji())
	assert.Equal(t, "Noach", ev.Basename())
	assert.Equal(t, []string{"parashat"}, ev.GetCategories())
}

