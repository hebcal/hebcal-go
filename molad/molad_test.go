package molad_test

import (
	"fmt"

	"github.com/hebcal/hebcal-go/hdate"
	"github.com/hebcal/hebcal-go/molad"
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
