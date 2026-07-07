package zmanim_test

import (
	"math"
	"testing"
	"time"

	"github.com/hebcal/hebcal-go/zmanim"
)

// Reference values were produced by the @hebcal/core JavaScript Zmanim class
// (which uses the same NOAA algorithm) for Jerusalem on 2024-04-26. A 2-second
// tolerance absorbs sub-second rounding differences between @hebcal/noaa and
// noaa-go; the JSON API rounds to the minute, so these are effectively exact.
func TestExtraZmanim(t *testing.T) {
	loc := zmanim.NewLocation("Jerusalem", "IL", 31.76904, 35.21633, 786, "Asia/Jerusalem")
	tz, _ := time.LoadLocation("Asia/Jerusalem")
	dt := time.Date(2024, time.April, 26, 12, 0, 0, 0, time.UTC)

	assertClose := func(t *testing.T, name string, got time.Time, wantISO string) {
		t.Helper()
		want, err := time.ParseInLocation("2006-01-02T15:04:05", wantISO, tz)
		if err != nil {
			t.Fatalf("%s: bad want %q: %v", name, wantISO, err)
		}
		if got.IsZero() {
			t.Errorf("%s: got zero time, want %s", name, wantISO)
			return
		}
		if d := math.Abs(got.Sub(want).Seconds()); d > 2 {
			t.Errorf("%s: got %s want ~%s (off %.0fs)", name, got.Format(time.RFC3339), wantISO, d)
		}
	}

	z := zmanim.New(&loc, dt) // UseElevation defaults to false
	// The MGA "72 minute" sof zman is measured from sea level, so it is the
	// same with or without elevation (see the ze assertions below).
	assertClose(t, "sofZmanShmaMGA", z.SofZmanShmaMGA(), "2024-04-26T08:41:40")
	assertClose(t, "sofZmanTfillaMGA", z.SofZmanTfillaMGA(), "2024-04-26T10:00:08")
	assertClose(t, "seaLevelSunrise", z.SeaLevelSunrise(), "2024-04-26T05:58:14")
	assertClose(t, "seaLevelSunset", z.SeaLevelSunset(), "2024-04-26T19:15:58")
	assertClose(t, "sofZmanShmaMGA16Point1", z.SofZmanShmaMGA16Point1(), "2024-04-26T08:38:56")
	assertClose(t, "sofZmanShmaMGA19Point8", z.SofZmanShmaMGA19Point8(), "2024-04-26T08:28:58")
	assertClose(t, "sofZmanTfillaMGA16Point1", z.SofZmanTfillaMGA16Point1(), "2024-04-26T09:58:22")
	assertClose(t, "sofZmanTfillaMGA19Point8", z.SofZmanTfillaMGA19Point8(), "2024-04-26T09:51:44")
	assertClose(t, "minchaGedolaMGA", z.MinchaGedolaMGA(), "2024-04-26T13:16:20")
	assertClose(t, "minchaKetanaMGA", z.MinchaKetanaMGA(), "2024-04-26T17:11:46")
	assertClose(t, "alosBaalHatanya", z.AlosBaalHatanya(), "2024-04-26T04:36:24")
	assertClose(t, "sofZmanShmaBaalHatanya", z.SofZmanShmaBaalHatanya(), "2024-04-26T09:15:50")
	assertClose(t, "sofZmanTfilaBaalHatanya", z.SofZmanTfilaBaalHatanya(), "2024-04-26T10:22:55")
	assertClose(t, "minchaGedolaBaalHatanya", z.MinchaGedolaBaalHatanya(), "2024-04-26T13:10:39")
	assertClose(t, "minchaKetanaBaalHatanya", z.MinchaKetanaBaalHatanya(), "2024-04-26T16:31:56")
	assertClose(t, "plagHaminchaBaalHatanya", z.PlagHaminchaBaalHatanya(), "2024-04-26T17:55:48")
	assertClose(t, "tzaisBaalHatanya", z.TzaisBaalHatanya(), "2024-04-26T19:41:38")

	// The MGA mincha zmanim use elevation-aware sunrise/sunset.
	ze := zmanim.New(&loc, dt)
	ze.UseElevation = true
	assertClose(t, "minchaGedolaMGA(elev)", ze.MinchaGedolaMGA(), "2024-04-26T13:16:42")
	assertClose(t, "minchaKetanaMGA(elev)", ze.MinchaKetanaMGA(), "2024-04-26T17:14:21")
	// The MGA 72-minute sof zman uses sea-level sunrise/sunset, so elevation
	// must not change it (regression test for the elevation-aware bug).
	assertClose(t, "sofZmanShmaMGA(elev)", ze.SofZmanShmaMGA(), "2024-04-26T08:41:40")
	assertClose(t, "sofZmanTfillaMGA(elev)", ze.SofZmanTfillaMGA(), "2024-04-26T10:00:08")
	// seaLevel variants ignore elevation.
	assertClose(t, "seaLevelSunrise(elev)", ze.SeaLevelSunrise(), "2024-04-26T05:58:14")
	assertClose(t, "seaLevelSunset(elev)", ze.SeaLevelSunset(), "2024-04-26T19:15:58")
}

// TestExtraZmanimPolar confirms the new zmanim return the zero time (rather than
// panicking) where the sun does not reach the required depression.
func TestExtraZmanimPolar(t *testing.T) {
	loc := zmanim.NewLocation("Tromso", "NO", 69.6489, 18.9551, 0, "Europe/Oslo")
	z := zmanim.New(&loc, time.Date(2024, time.June, 21, 12, 0, 0, 0, time.UTC))
	if got := z.AlosBaalHatanya(); !got.IsZero() {
		t.Errorf("expected zero AlosBaalHatanya in polar summer, got %v", got)
	}
	if got := z.SofZmanShmaMGA16Point1(); !got.IsZero() {
		t.Errorf("expected zero SofZmanShmaMGA16Point1 in polar summer, got %v", got)
	}
}
