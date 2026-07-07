package hebcal

import (
	"sync"
	"testing"
)

// TestGetHolidaysForYearCache verifies the memoization returns the identical
// cached slice (same backing array) on repeat calls with the same arguments,
// and distinct results for the Israel vs Diaspora schedule.
func TestGetHolidaysForYearCache(t *testing.T) {
	a := GetHolidaysForYear(5784, false)
	b := GetHolidaysForYear(5784, false)
	if len(a) == 0 || len(a) != len(b) || &a[0] != &b[0] {
		t.Fatalf("expected identical cached slice: len(a)=%d len(b)=%d", len(a), len(b))
	}
	// The cached slice is clamped (cap == len) so a caller's append cannot
	// corrupt the shared cache.
	if cap(a) != len(a) {
		t.Errorf("expected cap==len for cached slice, got cap=%d len=%d", cap(a), len(a))
	}
	// Israel and Diaspora schedules are cached separately and differ.
	il := GetHolidaysForYear(5784, true)
	if len(il) == len(a) {
		t.Errorf("expected Israel and Diaspora holiday counts to differ")
	}
}

// TestGetHolidaysForYearConcurrent exercises the RWMutex under -race.
func TestGetHolidaysForYearConcurrent(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			if len(GetHolidaysForYear(5780+(n%8), n%2 == 0)) == 0 {
				t.Errorf("empty holidays for iteration %d", n)
			}
		}(i)
	}
	wg.Wait()
}

func BenchmarkGetHolidaysForYearCached(b *testing.B) {
	GetHolidaysForYear(5784, false) // warm cache
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = GetHolidaysForYear(5784, false)
	}
}

func BenchmarkGetAllHolidaysForYear(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = getAllHolidaysForYear(5784)
	}
}
