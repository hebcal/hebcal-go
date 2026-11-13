package omer_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hebcal/hdate"
	"github.com/hebcal/hebcal-go/omer"
	"github.com/hebcal/locales"
	"github.com/stretchr/testify/assert"
)

func TestSefira(t *testing.T) {
	assert := assert.New(t)
	omer := omer.MustNewOmerEvent(hdate.New(5770, hdate.Sivan, 2), 46)
	assert.Equal("Eternity within Majesty", omer.Sefira("en"))
	assert.Equal("נֶּֽצַח שֶׁבְּמַלְכוּת", omer.Sefira("he"))
	assert.Equal("Netzach sheb'Malkhut", omer.Sefira("translit"))
}

func ExampleOmerEvent_Sefira() {
	ev := omer.MustNewOmerEvent(hdate.New(5770, hdate.Sivan, 2), 46)
	fmt.Println(ev.Sefira("en"))
	fmt.Println(ev.Sefira("he"))
	fmt.Println(ev.Sefira("translit"))
	// Output:
	// Eternity within Majesty
	// נֶּֽצַח שֶׁבְּמַלְכוּת
	// Netzach sheb'Malkhut
}

// omerEvents returns the 49 Omer events of 5782, indexed by day (index 0 unused).
func omerEvents() []omer.OmerEvent {
	events := make([]omer.OmerEvent, 50)
	startAbs := hdate.New(5782, hdate.Nisan, 16).Abs()
	for i := 1; i <= 49; i++ {
		events[i] = omer.MustNewOmerEvent(hdate.FromRD(startAbs+int64(i)-1), i)
	}
	return events
}

func TestTodayIsEn(t *testing.T) {
	actual := make([]string, 50)
	for i, ev := range omerEvents()[1:] {
		actual[i+1] = ev.TodayIs("en")
	}
	assert.Equal(t, coreTodayIsEn, actual)
}

func TestTodayIsHe(t *testing.T) {
	actual := make([]string, 50)
	for i, ev := range omerEvents()[1:] {
		actual[i+1] = ev.TodayIs("he")
	}
	assert.Equal(t, coreTodayIsHe, actual)
}

func TestTodayIsHeNoNikud(t *testing.T) {
	for i, ev := range omerEvents()[1:] {
		assert.Equal(t, locales.HebrewStripNikkud(coreTodayIsHe[i+1]), ev.TodayIs("he-x-NoNikud"))
	}
}

func TestSefiraHe(t *testing.T) {
	actual := make([]string, 50)
	for i, ev := range omerEvents()[1:] {
		actual[i+1] = ev.Sefira("he")
	}
	assert.Equal(t, coreSefiraHe, actual)
}

func TestSefiraHeNoNikud(t *testing.T) {
	for i, ev := range omerEvents()[1:] {
		assert.Equal(t, locales.HebrewStripNikkud(coreSefiraHe[i+1]), ev.Sefira("he-x-NoNikud"))
	}
}

func TestEmoji(t *testing.T) {
	expected := []string{
		"",
		"①", "②", "③", "④", "⑤", "⑥", "⑦",
		"⑧", "⑨", "⑩", "⑪", "⑫", "⑬", "⑭",
		"⑮", "⑯", "⑰", "⑱", "⑲", "⑳", "㉑",
		"㉒", "㉓", "㉔", "㉕", "㉖", "㉗", "㉘",
		"㉙", "㉚", "㉛", "㉜", "㉝", "㉞", "㉟",
		"㊱", "㊲", "㊳", "㊴", "㊵", "㊶", "㊷",
		"㊸", "㊹", "㊺", "㊻", "㊼", "㊽", "㊾",
	}
	actual := make([]string, 50)
	start := hdate.New(5782, hdate.Nisan, 16)
	startAbs := start.Abs()
	for i := 1; i <= 49; i++ {
		abs := startAbs + int64(i) - 1
		ev := omer.MustNewOmerEvent(hdate.FromRD(abs), i)
		actual[i] = ev.Emoji()
	}
	assert.Equal(t, expected, actual)
}

func ExampleOmerEvent_Emoji() {
	omer := omer.MustNewOmerEvent(hdate.New(5770, hdate.Nisan, 28), 13)
	fmt.Println(omer.Emoji())
	// Output: ⑬
}

func ExampleOmerEvent_TodayIs() {
	omer := omer.MustNewOmerEvent(hdate.New(5770, hdate.Nisan, 28), 13)
	fmt.Println(omer.TodayIs("en"))
	fmt.Println(omer.TodayIs("he"))
	fmt.Println(omer.TodayIs("he-x-NoNikud"))
	// Output:
	// Today is 13 days, which are 1 week and 6 days of the Omer
	// הַיּוֹם שְׁלוֹשָׁה עָשָׂר יוֹם, שֶׁהֵם שָׁבֽוּעַ אֶחָד וְשִׁשָּׁה יָמִים לָעֽוֹמֶר
	// היום שלושה עשר יום, שהם שבוע אחד וששה ימים לעומר
}

func TestRenderBrief(t *testing.T) {
	ev := omer.MustNewOmerEvent(hdate.New(5786, hdate.Iyyar, 7), 22)
	assert.Equal(t, "22nd day of the Omer", ev.Render("en"))
	assert.Equal(t, "Omer day 22", ev.RenderBrief("en"))
	assert.Equal(t, "עוֹמֶר יוֹם 22", ev.RenderBrief("he"))
}

func TestNewOmerEventRange(t *testing.T) {
	hd := hdate.New(5786, hdate.Nisan, 16)
	for _, day := range []int{-1, 0, 50, 1000} {
		_, err := omer.NewOmerEvent(hd, day)
		assert.True(t, errors.Is(err, omer.ErrOmerDay), "day %d", day)
		assert.Panics(t, func() { omer.MustNewOmerEvent(hd, day) })
	}
	for _, day := range []int{1, 7, 49} {
		ev, err := omer.NewOmerEvent(hd, day)
		assert.NoError(t, err)
		assert.Equal(t, day, ev.OmerDay)
		assert.Equal(t, ev, omer.MustNewOmerEvent(hd, day))
	}
}
