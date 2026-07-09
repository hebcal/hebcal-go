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

func TestLoadLocation(t *testing.T) {
	loc, err := zmanim.LoadLocation("America/New_York")
	if err != nil {
		t.Fatalf("LoadLocation: %v", err)
	}
	if loc.String() != "America/New_York" {
		t.Errorf("got %q, want America/New_York", loc.String())
	}
	// A second call must return the identical cached *time.Location.
	loc2, err := zmanim.LoadLocation("America/New_York")
	if err != nil {
		t.Fatalf("LoadLocation (cached): %v", err)
	}
	if loc != loc2 {
		t.Errorf("expected cached LoadLocation to return the same *time.Location pointer")
	}
	// Unknown timezones must still error (and are not cached).
	if _, err := zmanim.LoadLocation("Not/AZone"); err == nil {
		t.Error("expected an error for an invalid timezone")
	}
}
