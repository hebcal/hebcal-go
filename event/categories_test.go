package event

import (
	"reflect"
	"testing"

	"github.com/hebcal/hdate"
)

func TestHolidayEventGetCategories(t *testing.T) {
	tests := []struct {
		desc  string
		flags HolidayFlags
		chm   int
		want  []string
	}{
		{"Pesach VII", CHAG, 0, []string{"holiday", "major"}},
		{"Yom Kippur", CHAG | MAJOR_FAST, 0, []string{"holiday", "major", "fast"}},
		{"Chanukah: 3 Candles", CHANUKAH_CANDLES | MINOR_HOLIDAY, 0, []string{"holiday", "minor"}},
		{"Tzom Gedaliah", MINOR_FAST, 0, []string{"holiday", "fast"}},
		{"Tu BiShvat", MINOR_HOLIDAY, 0, []string{"holiday", "minor"}},
		{"Yom HaAtzma'ut", MODERN_HOLIDAY, 0, []string{"holiday", "modern"}},
		{"Shabbat Shekalim", SPECIAL_SHABBAT, 0, []string{"holiday", "shabbat"}},
		{"Rosh Chodesh Sivan", ROSH_CHODESH, 0, []string{"roshchodesh"}},
		{"Shabbat Mevarchim Chodesh Sivan", SHABBAT_MEVARCHIM, 0, []string{"mevarchim"}},
		{"Pesach II (CH''M)", CHUL_ONLY, 1, []string{"holiday", "major", "cholhamoed"}},
		// flags alone inconclusive: minor-holiday fallback list vs major default
		{"Lag BaOmer", 0, 0, []string{"holiday", "minor"}},
		{"Erev Purim", EREV, 0, []string{"holiday", "minor"}},
		{"Some Random Holiday", 0, 0, []string{"holiday", "major"}},
	}
	for _, tc := range tests {
		ev := HolidayEvent{Date: hdate.New(5784, hdate.Sivan, 1), Desc: tc.desc, Flags: tc.flags, CholHaMoedDay: tc.chm}
		if got := ev.GetCategories(); !reflect.DeepEqual(got, tc.want) {
			t.Errorf("%q %b: GetCategories() = %v, want %v", tc.desc, tc.flags, got, tc.want)
		}
	}
}

func TestCategoriesFromFlags(t *testing.T) {
	if got := CategoriesFromFlags(0); !reflect.DeepEqual(got, []string{"unknown"}) {
		t.Errorf("CategoriesFromFlags(0) = %v, want [unknown]", got)
	}
	if got := CategoriesFromFlags(PARSHA_HASHAVUA); !reflect.DeepEqual(got, []string{"parashat"}) {
		t.Errorf("CategoriesFromFlags(PARSHA_HASHAVUA) = %v", got)
	}
	// the returned slice must be a copy the caller can mutate safely
	a := CategoriesFromFlags(MAJOR_FAST)
	a[0] = "mutated"
	b := CategoriesFromFlags(MAJOR_FAST)
	if b[0] != "holiday" {
		t.Errorf("CategoriesFromFlags returned a shared slice: %v", b)
	}
}
