package event

import (
	"reflect"
	"testing"

	"github.com/hebcal/hdate"
)

func TestHolidayEventGetCategories(t *testing.T) {
	tests := []struct {
		desc  string
		flags Flags
		chm   int
		want  []string
	}{
		{"Pesach VII", Chag, 0, []string{"holiday", "major"}},
		{"Yom Kippur", Chag | MajorFast, 0, []string{"holiday", "major", "fast"}},
		{"Chanukah: 3 Candles", ChanukahCandles | MinorHoliday, 0, []string{"holiday", "minor"}},
		{"Tzom Gedaliah", MinorFast, 0, []string{"holiday", "fast"}},
		{"Tu BiShvat", MinorHoliday, 0, []string{"holiday", "minor"}},
		{"Yom HaAtzma'ut", ModernHoliday, 0, []string{"holiday", "modern"}},
		{"Shabbat Shekalim", SpecialShabbat, 0, []string{"holiday", "shabbat"}},
		{"Rosh Chodesh Sivan", RoshChodesh, 0, []string{"roshchodesh"}},
		{"Shabbat Mevarchim Chodesh Sivan", ShabbatMevarchim, 0, []string{"mevarchim"}},
		{"Pesach II (CH''M)", ChulOnly, 1, []string{"holiday", "major", "cholhamoed"}},
		// flags alone inconclusive: minor-holiday fallback list vs major default
		{"Lag BaOmer", 0, 0, []string{"holiday", "minor"}},
		{"Erev Purim", Erev, 0, []string{"holiday", "minor"}},
		{"Some Random Holiday", 0, 0, []string{"holiday", "major"}},
	}
	for _, tc := range tests {
		ev := NewHolidayEvent(hdate.New(5784, hdate.Sivan, 1), tc.desc, tc.flags, "")
		ev.CholHaMoedDay = tc.chm
		if got := ev.Categories(); !reflect.DeepEqual(got, tc.want) {
			t.Errorf("%q %b: GetCategories() = %v, want %v", tc.desc, tc.flags, got, tc.want)
		}
	}
}

func TestCategoriesFromFlags(t *testing.T) {
	if got := CategoriesFromFlags(0); !reflect.DeepEqual(got, []string{"unknown"}) {
		t.Errorf("CategoriesFromFlags(0) = %v, want [unknown]", got)
	}
	if got := CategoriesFromFlags(ParshaHashavua); !reflect.DeepEqual(got, []string{"parashat"}) {
		t.Errorf("CategoriesFromFlags(PARSHA_HASHAVUA) = %v", got)
	}
	// the returned slice must be a copy the caller can mutate safely
	a := CategoriesFromFlags(MajorFast)
	a[0] = "mutated"
	b := CategoriesFromFlags(MajorFast)
	if b[0] != "holiday" {
		t.Errorf("CategoriesFromFlags returned a shared slice: %v", b)
	}
}
