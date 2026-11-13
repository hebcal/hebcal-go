package zmanim_test

import (
	"errors"
	"math"
	"testing"
	"time"

	"github.com/hebcal/hebcal-go/zmanim"
	"github.com/stretchr/testify/assert"
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

func TestNewLocationErrors(t *testing.T) {
	_, err := zmanim.NewLocation("x", "US", 90.5, 0, 0, "America/New_York")
	assert.True(t, errors.Is(err, zmanim.ErrLatitude))
	_, err = zmanim.NewLocation("x", "US", 0, -181, 0, "America/New_York")
	assert.True(t, errors.Is(err, zmanim.ErrLongitude))
	_, err = zmanim.NewLocation("x", "US", math.NaN(), 0, 0, "America/New_York")
	assert.True(t, errors.Is(err, zmanim.ErrLatitude))
	loc, err := zmanim.NewLocation("x", "US", 40.7, -74, -5, "America/New_York")
	assert.NoError(t, err)
	assert.Equal(t, 0, loc.Elevation)
	assert.Panics(t, func() { zmanim.MustNewLocation("x", "US", 91, 0, 0, "UTC") })
}

func TestNewErrors(t *testing.T) {
	date := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	_, err := zmanim.New(nil, date)
	assert.Error(t, err)
	_, err = zmanim.New(&zmanim.Location{Latitude: 40, Longitude: -74, TimeZoneID: "No/Such_Zone"}, date)
	assert.True(t, errors.Is(err, zmanim.ErrTimeZone), "%v", err)
	_, err = zmanim.New(&zmanim.Location{Latitude: -91, Longitude: -74, TimeZoneID: "UTC"}, date)
	assert.True(t, errors.Is(err, zmanim.ErrLatitude), "%v", err)
	_, err = zmanim.New(&zmanim.Location{Latitude: 40, Longitude: 181, TimeZoneID: "UTC"}, date)
	assert.True(t, errors.Is(err, zmanim.ErrLongitude), "%v", err)
	assert.Panics(t, func() { zmanim.MustNew(&zmanim.Location{TimeZoneID: "No/Such_Zone"}, date) })
	z, err := zmanim.New(zmanim.LookupCity("Boston"), date)
	assert.NoError(t, err)
	assert.False(t, z.Sunset().IsZero())
}

// AllCities returns a copy: changing it does not change the built-in cities.
func TestAllCitiesIsACopy(t *testing.T) {
	cities := zmanim.AllCities()
	assert.NotEmpty(t, cities)
	name := cities[0].Name
	cities[0].Name = "Atlantis"
	assert.Equal(t, name, zmanim.AllCities()[0].Name)
	assert.Nil(t, zmanim.LookupCity("Atlantis"))
}
