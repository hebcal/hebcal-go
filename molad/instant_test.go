package molad_test

import (
	"testing"
	"time"

	"github.com/hebcal/hdate"
	"github.com/hebcal/hebcal-go/molad"
	"github.com/stretchr/testify/assert"
)

const millisUTC = "2006-01-02T15:04:05.000Z"

func TestMoladInstant(t *testing.T) {
	tests := []struct {
		year  int
		month hdate.HMonth
		want  [5]string // molad, 3 days, 7 days, between moldos, 15 days
	}{
		{5784, hdate.Nisan, [5]string{"2024-04-08T20:36:26.837Z", "2024-04-11T20:36:26.837Z",
			"2024-04-15T20:36:26.837Z", "2024-04-23T14:58:28.503Z", "2024-04-23T20:36:26.837Z"}},
		// on 27 Tevet
		{5541, hdate.Shvat, [5]string{"1781-01-24T19:57:20.170Z", "1781-01-27T19:57:20.170Z",
			"1781-01-31T19:57:20.170Z", "1781-02-08T14:19:21.836Z", "1781-02-08T19:57:20.170Z"}},
		// before the Gregorian reform: the instant is proleptic Gregorian
		{5000, hdate.Nisan, [5]string{"1240-03-31T16:29:43.504Z", "1240-04-03T16:29:43.504Z",
			"1240-04-07T16:29:43.504Z", "1240-04-15T10:51:45.170Z", "1240-04-15T16:29:43.504Z"}},
	}
	for _, tc := range tests {
		m := molad.New(tc.year, tc.month)
		got := [5]string{}
		for i, ti := range []time.Time{m.Instant(), m.TchilasZmanKidushLevana3Days(),
			m.TchilasZmanKidushLevana7Days(), m.SofZmanKidushLevanaBetweenMoldos(),
			m.SofZmanKidushLevana15Days()} {
			got[i] = ti.Format(millisUTC)
		}
		assert.Equal(t, tc.want, got, "molad %d %v", tc.year, tc.month)
	}
}
