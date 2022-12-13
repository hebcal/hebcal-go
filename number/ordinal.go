package number

import "strconv"

// Ordinal appends English language suffix (1 -> "1st", 2 -> "2nd", etc.)
func Ordinal(n int) string {
	str := strconv.Itoa(n)
	i := n % 100
	if i/10 == 1 {
		return str + "th"
	}
	switch i % 10 {
	case 1:
		return str + "st"
	case 2:
		return str + "nd"
	case 3:
		return str + "rd"
	default:
		return str + "th"
	}
}
