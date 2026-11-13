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
	assert.Equal(t, event.HebrewDate, ev.GetFlags())
	assert.Equal(t, "", ev.GetEmoji())
	assert.Equal(t, []string{"hebdate"}, ev.GetCategories())
}

func TestHolidayEvent_Extra(t *testing.T) {
	hd := hdate.New(5781, hdate.Sivan, 6)
	ev := event.NewHolidayEvent(
		hd,
		"Shavuot I",
		event.Chag,
		"")
	assert.Equal(t, hd, ev.GetDate())
	assert.Equal(t, event.Chag, ev.GetFlags())
	assert.Equal(t, "✡️", ev.GetEmoji())

	ev2 := event.NewHolidayEvent(
		hd,
		"Shabbat Shekalim",
		event.SpecialShabbat,
		"")
	assert.Equal(t, "🕍", ev2.GetEmoji())

	ev3 := event.NewHolidayEvent(
		hd,
		"Rosh Chodesh Tamuz",
		event.RoshChodesh,
		"")
	assert.Equal(t, "🌒", ev3.GetEmoji())

	ev4 := event.NewHolidayEvent(
		hd,
		"My Custom Holiday",
		event.Chag,
		"🎉")
	assert.Equal(t, "🎉", ev4.GetEmoji())

	// Test Categories
	evChol := event.NewHolidayEvent(hd, "Pesach III (CH''M)", event.CholHamoed, "")
	evChol.CholHaMoedDay = 3
	assert.Equal(t, []string{"holiday", "major", "cholhamoed"}, evChol.GetCategories())

	evMinor := event.NewHolidayEvent(
		hd,
		"Lag BaOmer",
		0,
		"")
	assert.Equal(t, []string{"holiday", "minor"}, evMinor.GetCategories())

	evMajorFall := event.NewHolidayEvent(
		hd,
		"Some Unknown Major",
		0,
		"")
	assert.Equal(t, []string{"holiday", "major"}, evMajorFall.GetCategories())

	// Test Render for Rosh Chodesh
	evRoshChodesh := event.NewHolidayEvent(
		hd,
		"Rosh Chodesh Sh'vat",
		event.RoshChodesh,
		"")
	assert.Equal(t, "Rosh Chodesh Sh’vat", evRoshChodesh.Render("en"))
	assert.Equal(t, "ראש חודש שבט", evRoshChodesh.Render("he-x-nonikud"))

	// Test Render for Shabbat Mevarchim
	evMevarchim := event.NewHolidayEvent(
		hd,
		"Shabbat Mevarchim Chodesh Kislev",
		event.ShabbatMevarchim,
		"")
	assert.Equal(t, "Shabbat Mevarchim Chodesh Kislev", evMevarchim.Render("en"))

	// Test Render for Rosh Hashana Hebrew
	evRoshHashana := event.NewHolidayEvent(
		hdate.New(5781, hdate.Tishrei, 1),
		"Rosh Hashana 5781",
		event.Chag,
		"")
	assert.Equal(t, "רֹאשׁ הַשָּׁנָה תשפ״א", evRoshHashana.Render("he"))
	assert.Equal(t, "Rosh Hashana 5781", evRoshHashana.Render("en"))

	// Test Render for Yom Kippur Katan
	evYKK := event.NewHolidayEvent(
		hd,
		"Yom Kippur Katan Sivan",
		event.MinorFast|event.YomKippurKatan,
		"")
	assert.Equal(t, "Yom Kippur Katan Sivan", evYKK.Render("en"))

	// Test GetEmoji edge cases
	evMevarchimEmoji := event.NewHolidayEvent(
		hd,
		"Shabbat Mevarchim Chodesh Kislev",
		event.ShabbatMevarchim,
		"")
	assert.Equal(t, "", evMevarchimEmoji.GetEmoji())

	evYKKEmoji := event.NewHolidayEvent(
		hd,
		"Yom Kippur Katan Sivan",
		event.YomKippurKatan|event.MinorFast,
		"")
	assert.Equal(t, "", evYKKEmoji.GetEmoji())
}

func TestMevarchimChodeshEvent(t *testing.T) {
	hd := hdate.New(5786, hdate.Sivan, 25)
	m := molad.New(5786, hdate.Tamuz)
	ev := event.NewMevarchimChodeshEvent(hd, "Tammuz", m)

	assert.Equal(t, hd, ev.GetDate())
	assert.Equal(t, event.ShabbatMevarchim, ev.GetFlags())
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

	ev := event.NewMoladEvent(hd, m, "Tammuz", "US", event.ClockAuto)
	assert.Equal(t, hd, ev.GetDate())
	assert.Equal(t, event.Molad, ev.GetFlags())
	assert.Equal(t, "", ev.GetEmoji())
	assert.Equal(t, []string{"molad"}, ev.GetCategories())
	assert.Equal(t, "Molad Tammuz 5786", ev.Basename())

	// Test renderings for different locales and country codes
	assert.Contains(t, ev.Render("en"), "Molad Tammuz: Monday, 6:46am and 16 chalakim")

	evGB := event.NewMoladEvent(hd, m, "Tammuz", "GB", event.ClockAuto)
	assert.Contains(t, evGB.Render("en"), "Molad Tammuz: Monday, 6:46 and 16 chalakim")

	evDefaultCC := event.NewMoladEvent(hd, m, "Tammuz", "", event.ClockAuto)
	assert.Contains(t, evDefaultCC.Render("en"), "Molad Tammuz: Monday, 6:46am and 16 chalakim")

	// Test hebrew rendering (various hours / periods of day)
	assert.Contains(t, ev.Render("he"), "בֹּקֶר")

	mAfternoon := m
	mAfternoon.Hours = 14
	evAfternoon := event.NewMoladEvent(hd, mAfternoon, "Tammuz", "US", event.ClockAuto)
	assert.Contains(t, evAfternoon.Render("he"), "בַּצׇּהֳרַיִים")
	assert.Contains(t, evAfternoon.Render("he-x-nonikud"), "בצהריים")

	mEvening := m
	mEvening.Hours = 19
	evEvening := event.NewMoladEvent(hd, mEvening, "Tammuz", "US", event.ClockAuto)
	assert.Contains(t, evEvening.Render("he"), "בָּעֶרֶב")

	mNight := m
	mNight.Hours = 23
	evNight := event.NewMoladEvent(hd, mNight, "Tammuz", "US", event.ClockAuto)
	assert.Contains(t, evNight.Render("he"), "בַּלַּ֥יְלָה")

	mNight2 := m
	mNight2.Hours = 2
	evNight2 := event.NewMoladEvent(hd, mNight2, "Tammuz", "US", event.ClockAuto)
	assert.Contains(t, evNight2.Render("he"), "בַּלַּ֥יְלָה")

	// 12-hour clock boundary cases: 0:00 (12:00am) and 12:00 (12:00pm)
	mMidnight := m
	mMidnight.Hours = 0
	mMidnight.Minutes = 0
	evMidnight := event.NewMoladEvent(hd, mMidnight, "Tammuz", "US", event.ClockAuto)
	assert.Contains(t, evMidnight.Render("en"), "12:00am")

	mNoon := m
	mNoon.Hours = 12
	mNoon.Minutes = 0
	evNoon := event.NewMoladEvent(hd, mNoon, "Tammuz", "US", event.ClockAuto)
	assert.Contains(t, evNoon.Render("en"), "12:00pm")

	mPM := m
	mPM.Hours = 13
	mPM.Minutes = 30
	evPM := event.NewMoladEvent(hd, mPM, "Tammuz", "US", event.ClockAuto)
	assert.Contains(t, evPM.Render("en"), "1:30pm")
}

func TestUserEvent(t *testing.T) {
	hd := hdate.New(5780, hdate.Nisan, 15)
	ev := event.NewUserEvent(
		hd,
		"My Event")

	assert.Equal(t, hd, ev.GetDate())
	assert.Equal(t, event.User, ev.GetFlags())
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
	assert.Equal(t, event.ParshaHashavua, ev.GetFlags())
	assert.Equal(t, "", ev.GetEmoji())
	assert.Equal(t, "Noach", ev.Basename())
	assert.Equal(t, []string{"parashat"}, ev.GetCategories())
}

func TestSmartApostrophe(t *testing.T) {
	assert.Equal(t, "Sh’vat", event.SmartApostrophe("Sh'vat"))
	assert.Equal(t, "Tish’a B’Av", event.SmartApostrophe("Tish'a B'Av"))
	assert.Equal(t, "Sivan", event.SmartApostrophe("Sivan"))
}
