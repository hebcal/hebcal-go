package hdate_test

import (
	"fmt"

	"github.com/hebcal/hebcal-go/hdate"
)

func ExampleNewMolad() {
	molad := hdate.NewMolad(5783, hdate.Iyyar)
	dayOfWeek := molad.Date.Weekday().String()[0:3]
	fmt.Printf("Molad Iyyar: %s, %d minutes and %d chalakim after %d:00",
		dayOfWeek, molad.Minutes, molad.Chalakim, molad.Hours)
	// Output: Molad Iyyar: Thu, 8 minutes and 13 chalakim after 14:00
}
