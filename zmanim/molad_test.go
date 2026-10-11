package zmanim

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const millisOffset = "2006-01-02T15:04:05.000-07:00"

func fmtZman(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(millisOffset)
}

func TestMoladBasedZmanim(t *testing.T) {
	jerusalem := NewLocation("Jerusalem", "IL", 31.76904, 35.21633, 0, "Asia/Jerusalem")
	// Rapa Iti in Dec 2027 is an extreme case from Shaarei Zmanim: the molad
	// falls on 27 Kislev and the 3-day time on the evening of the 30th.
	rapaIti := NewLocation("Rapa Iti", "PF", -27.61, -144.33, 0, "Pacific/Tahiti")
	tests := []struct {
		loc  *Location
		date string
		want [5]string // molad, 3 days, 7 days, between moldos, 15 days
	}{
		// 29 Cheshvan 5787
		{&jerusalem, "2026-11-09", [5]string{"2026-11-09T22:06:13.504+02:00", "", "", "", ""}},
		// 30 Cheshvan 5787: no molad-based time
		{&jerusalem, "2026-11-10", [5]string{}},
		{&jerusalem, "2026-11-12", [5]string{"", "2026-11-12T22:06:13.504+02:00", "", "", ""}},
		{&jerusalem, "2026-11-16", [5]string{"", "", "2026-11-16T22:06:13.504+02:00", "", ""}},
		{&jerusalem, "2026-11-24", [5]string{"", "", "",
			"2026-11-24T16:28:15.170+02:00", "2026-11-24T22:06:13.504+02:00"}},
		{&jerusalem, "2026-11-25", [5]string{}},
		// 27 Kislev 5788
		{&rapaIti, "2027-12-27", [5]string{"2027-12-27T20:23:00.170-10:00", "", "", "", ""}},
		// 30 Kislev 5788
		{&rapaIti, "2027-12-30", [5]string{"", "2027-12-30T20:23:00.170-10:00", "", "", ""}},
		{&rapaIti, "2028-01-03", [5]string{"", "", "2028-01-03T20:23:00.170-10:00", "", ""}},
		{&rapaIti, "2028-01-11", [5]string{"", "", "",
			"2028-01-11T14:45:01.836-10:00", "2028-01-11T20:23:00.170-10:00"}},
	}
	for _, tc := range tests {
		dt, _ := time.Parse("2006-01-02", tc.date)
		z := New(tc.loc, dt)
		var zero time.Time
		got := [5]string{
			fmtZman(z.ZmanMolad()),
			fmtZman(z.TchilasZmanKidushLevana3Days(zero, zero)),
			fmtZman(z.TchilasZmanKidushLevana7Days(zero, zero)),
			fmtZman(z.SofZmanKidushLevanaBetweenMoldos(zero, zero)),
			fmtZman(z.SofZmanKidushLevana15Days(zero, zero)),
		}
		assert.Equal(t, tc.want, got, "%s %s", tc.loc.Name, tc.date)
	}
}

func TestMoladBasedZmanimDaytimeAdjustment(t *testing.T) {
	jerusalem := NewLocation("Jerusalem", "IL", 31.76904, 35.21633, 0, "Asia/Jerusalem")
	dt, _ := time.Parse("2006-01-02", "2026-11-24")
	z := New(&jerusalem, dt)
	alos, tzeit := z.AlotHaShachar(), z.Tzeit(Tzeit3SmallStars)
	// 16:28 is daytime, so the latest time moves back to that morning's alos
	assert.Equal(t, alos, z.SofZmanKidushLevanaBetweenMoldos(alos, tzeit))
	// 22:06 is already night and is left alone
	assert.Equal(t, "2026-11-24T22:06:13.504+02:00", fmtZman(z.SofZmanKidushLevana15Days(alos, tzeit)))

	// a daytime start (14:45) moves forward to that evening's tzeit
	dt, _ = time.Parse("2006-01-02", "2026-01-21")
	z = New(&jerusalem, dt)
	assert.Equal(t, "2026-01-21T14:45:40.170+02:00",
		fmtZman(z.TchilasZmanKidushLevana3Days(time.Time{}, time.Time{})))
	tzeit = z.Tzeit(Tzeit3SmallStars)
	assert.Equal(t, tzeit, z.TchilasZmanKidushLevana3Days(z.AlotHaShachar(), tzeit))
}
