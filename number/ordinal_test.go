package number_test

import (
	"fmt"

	"github.com/hebcal/hebcal-go/number"
)

func ExampleOrdinal() {
	vals := []int{1, 2, 3, 4, 11, 12, 13, 14, 41, 42, 43, 44}
	for _, val := range vals {
		fmt.Printf("%d -> %s\n", val, number.Ordinal(val))
	}
	// Output:
	// 1 -> 1st
	// 2 -> 2nd
	// 3 -> 3rd
	// 4 -> 4th
	// 11 -> 11th
	// 12 -> 12th
	// 13 -> 13th
	// 14 -> 14th
	// 41 -> 41st
	// 42 -> 42nd
	// 43 -> 43rd
	// 44 -> 44th
}
