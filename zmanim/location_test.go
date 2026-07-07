package zmanim_test

import (
	"testing"

	"github.com/hebcal/hebcal-go/zmanim"
)

func TestNoDuplicateCities(t *testing.T) {
	m := make(map[string]int)
	for idx, city := range zmanim.AllCities() {
		prev, found := m[city.Name]
		if found {
			t.Errorf("Found %s at %d and %d", city.Name, prev, idx)
		}
		m[city.Name] = idx
	}
}

func TestCityElevation(t *testing.T) {
	jerusalem := zmanim.LookupCity("Jerusalem")
	if jerusalem == nil {
		t.Fatal("Jerusalem not found")
	}
	if jerusalem.Elevation != 786 {
		t.Errorf("Jerusalem elevation = %d, want 786", jerusalem.Elevation)
	}
	// Every city should have a non-negative elevation populated.
	for _, city := range zmanim.AllCities() {
		if city.Elevation < 0 {
			t.Errorf("%s has negative elevation %d", city.Name, city.Elevation)
		}
	}
}
