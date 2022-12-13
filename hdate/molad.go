package hdate

type Molad struct {
	Date     HDate
	Hours    int
	Minutes  int
	Chalakim int
}

// Borrowed from com.kosherjava.zmanim.hebrewcalendar.JewishDate

// Days from the beginning of Sunday till molad BaHaRaD.
// Calculated as 1 day, 5 hours and 204 chalakim = (24 + 5) * 1080 + 204 = 31524
const chalakimMoladTohu int64 = 31524

const chalakimPerMinute = 18
const chalakimPerHour = 1080
const chalakimPerDay int64 = 25920 // 24 * 1080

// The number of chalakim in an average Jewish month. A month has 29 days, 12 hours and 793
// chalakim (44 minutes and 3.3 seconds) for a total of 765,433 chalakim
const chalakimPerMonth int64 = 765433 // (29 * 24 + 12) * 1080 + 793

// Converts the NISSAN based constants used by this class to numeric month
// starting from TISHREI. This is required for Molad claculations.
func getJewishMonthOfYear(year int, month HMonth) int {
	offset := 0
	if IsLeapYear(year) {
		offset = 1
	}
	return (int(month)+(offset+5))%(offset+12) + 1
}

// Returns the number of chalakim (parts - 1080 to the hour) from the original hypothetical Molad Tohu to the year
// and month passed in.
func getChalakimSinceMoladTohu(year int, month HMonth) int64 {
	// Jewish lunar month = 29 days, 12 hours and 793 chalakim
	// chalakim since Molad Tohu BeHaRaD - 1 day, 5 hours and 204 chalakim
	monthOfYear := getJewishMonthOfYear(year, month)
	prevYear := int64(year) - 1
	monthsElapsed := (235 * (prevYear / 19)) + // Months in complete 19 year lunar (Metonic) cycles so far
		(12 * (prevYear % 19)) + // Regular months in this cycle
		((7*(prevYear%19) + 1) / 19) + // Leap months this cycle
		int64(monthOfYear-1) // add elapsed months till the start of the molad of the month
	// return chalakim prior to BeHaRaD + number of chalakim since
	return chalakimMoladTohu + (chalakimPerMonth * monthsElapsed)
}

func moladToAbsDate(chalakim int64) int64 {
	return (chalakim / chalakimPerDay) + epoch - 1
}

// Sets the molad time (hours minutes and chalakim) based on the number of chalakim since the start of the day.
func (m *Molad) setMoladTime(chalakim int) {
	adjustedChalakim := chalakim
	m.Hours = adjustedChalakim / chalakimPerHour
	adjustedChalakim = adjustedChalakim - (m.Hours * chalakimPerHour)
	m.Minutes = adjustedChalakim / chalakimPerMinute
	m.Chalakim = (adjustedChalakim - m.Minutes*chalakimPerMinute)
}

func NewMolad(year int, month HMonth) Molad {
	chalakim := getChalakimSinceMoladTohu(year, month)
	hd := FromRD(moladToAbsDate(chalakim))
	conjunctionDay := chalakim / chalakimPerDay
	conjunctionParts := chalakim - conjunctionDay*chalakimPerDay
	molad := Molad{}
	molad.setMoladTime(int(conjunctionParts))
	if molad.Hours >= 6 {
		hd = hd.Next()
	}
	molad.Date = hd
	molad.Hours = (molad.Hours + 18) % 24
	return molad
}

// This is Dershowitz & Reingold's floating point version here for reference
func molad(year int, month HMonth) float64 {
	yy := year
	if month < Tishrei {
		yy++
	}
	monthsElapsed := int(month-Tishrei) + (235*yy-234)/19
	return float64(epoch) - 876.0/25920.0 + float64(monthsElapsed)*(29.5+(793.0/25920.0))
}
