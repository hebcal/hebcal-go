package zmanim

import "testing"

func TestNoDuplicateCities(t *testing.T) {
	m := make(map[string]int)
	for idx, city := range AllCities() {
		prev, found := m[city.Name]
		if found {
			t.Errorf("Found %s at %d and %d", city.Name, prev, idx)
		}
		m[city.Name] = idx
	}
}
