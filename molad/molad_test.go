package molad_test

import (
	"fmt"
	"testing"

	"github.com/hebcal/hdate"
	"github.com/hebcal/hebcal-go/event"
	"github.com/hebcal/hebcal-go/molad"
	"github.com/stretchr/testify/assert"
)

func ExampleNew() {
	month := hdate.Iyyar
	molad := molad.New(5783, month)
	dayOfWeek := molad.Date.Weekday().String()
	fmt.Printf("Molad %s: %s, %d minutes and %d chalakim after %d:00",
		month.String(), dayOfWeek,
		molad.Minutes, molad.Chalakim, molad.Hours)
	// Output: Molad Iyyar: Thursday, 8 minutes and 13 chalakim after 14:00
}

func TestMoladEvent_Render(t *testing.T) {
	month := hdate.Iyyar
	molad := molad.New(5783, month)
	hd := hdate.New(5783, hdate.Nisan, 24)
	ev := event.NewMoladEvent(hd, molad, month.String(), "US")
	assert.Equal(t, "Molad Iyyar: Thursday, 2:08pm and 13 chalakim", ev.Render("en"))
	assert.Equal(t, "מוֹלָד הָלְּבָנָה אִיָּיר יִהְיֶה בַּיּוֹם חֲמִישִׁי בשָׁבוּעַ, בְּשָׁעָה 14 בַּצׇּהֳרַיִים, ו-8 דַּקּוֹת ו-13 חֲלָקִים", ev.Render("he"))
	assert.Equal(t, "מולד הלבנה אייר יהיה ביום חמישי בשבוע, בשעה 14 בצהריים, ו-8 דקות ו-13 חלקים", ev.Render("he-x-NoNikud"))
	// A 24-hour country (e.g. IL) omits the am/pm suffix.
	evIL := event.NewMoladEvent(hd, molad, month.String(), "IL")
	assert.Equal(t, "Molad Iyyar: Thursday, 14:08 and 13 chalakim", evIL.Render("en"))
}
