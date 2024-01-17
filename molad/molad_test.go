package molad_test

import (
	"fmt"
	"testing"

	"github.com/hebcal/hdate"
	"github.com/MaxBGreenberg/hebcal-go/event"
	"github.com/MaxBGreenberg/hebcal-go/molad"
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
	ev := event.NewMoladEvent(hd, molad, month.String())
	assert.Equal(t, "Molad Iyyar: Thu, 8 minutes and 13 chalakim after 14:00", ev.Render("en"))
	assert.Equal(t, "מוֹלָד הָלְּבָנָה אִיָיר יִהְיֶה בַּיּוֹם חֲמִישִׁי בשָׁבוּעַ, בְּשָׁעָה 14 בַּצׇּהֳרַיִים, ו-8 דַּקּוֹת ו-13 חֲלָקִים", ev.Render("he"))
}
