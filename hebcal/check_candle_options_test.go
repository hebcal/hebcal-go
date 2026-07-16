package hebcal

import (
	"testing"

	"github.com/hebcal/hebcal-go/zmanim"
	"github.com/stretchr/testify/assert"
)

func TestCheckCandleOptions(t *testing.T) {
	// Look up cities
	boston := zmanim.LookupCity("Boston")
	if boston == nil {
		t.Fatal("Boston not found")
	}
	telAviv := zmanim.LookupCity("Tel Aviv")
	if telAviv == nil {
		t.Fatal("Tel Aviv not found")
	}
	haifa := zmanim.LookupCity("Haifa")
	if haifa == nil {
		t.Fatal("Haifa not found")
	}
	jerusalem := zmanim.LookupCity("Jerusalem")
	if jerusalem == nil {
		t.Fatal("Jerusalem not found")
	}

	tests := []struct {
		name               string
		city               *zmanim.Location
		candleLightingMins int
		wantMins           int
	}{
		// Default offsets (when opts.CandleLightingMins is unspecified/0)
		{"Boston default", boston, 0, -18},
		{"Tel Aviv default", telAviv, 0, -20},
		{"Haifa default", haifa, 0, -30},
		{"Jerusalem default", jerusalem, 0, -40},

		// Explicit offset of 22 minutes
		{"Boston explicit 22", boston, 22, -22},
		{"Tel Aviv explicit 22", telAviv, 22, -22},
		{"Haifa explicit 22", haifa, 22, -22},
		{"Jerusalem explicit 22", jerusalem, 22, -22},

		// Explicit offset of negative 22 minutes (should behave the same as positive 22)
		{"Boston explicit -22", boston, -22, -22},

		// Explicit offset of 18 minutes.
		// Note: since 18 is the default offset, the implementation cannot distinguish
		// between "explicit 18" and "unspecified (defaulting to 18)".
		// Consequently, for Israeli cities, explicit 18 is overridden by city defaults:
		// Tel Aviv gets 20, Haifa gets 30, Jerusalem gets 40.
		{"Boston explicit 18", boston, 18, -18},
		{"Tel Aviv explicit 18", telAviv, 18, -20},
		{"Haifa explicit 18", haifa, 18, -30},
		{"Jerusalem explicit 18", jerusalem, 18, -40},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			opts := &CalOptions{
				CandleLighting:     true,
				Location:           tc.city,
				CandleLightingMins: tc.candleLightingMins,
			}
			err := checkCandleOptions(opts)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantMins, opts.CandleLightingMins)
		})
	}
}

func TestCheckCandleOptions_NoCandleLighting(t *testing.T) {
	opts := &CalOptions{
		CandleLighting:     false,
		CandleLightingMins: 0,
	}
	err := checkCandleOptions(opts)
	assert.NoError(t, err)
	assert.Equal(t, 0, opts.CandleLightingMins)
}

func TestCheckCandleOptions_MissingLocation(t *testing.T) {
	opts := &CalOptions{
		CandleLighting: true,
	}
	err := checkCandleOptions(opts)
	assert.Error(t, err)
	assert.Equal(t, "opts.CandleLighting requires opts.Location", err.Error())
}

func TestCheckCandleOptions_MutuallyExclusiveHavdalah(t *testing.T) {
	boston := zmanim.LookupCity("Boston")
	opts := &CalOptions{
		CandleLighting: true,
		Location:       boston,
		HavdalahMins:   50,
		HavdalahDeg:    8.5,
	}
	err := checkCandleOptions(opts)
	assert.Error(t, err)
	assert.Equal(t, "opts.HavdalahMins and opts.HavdalahDeg are mutually exclusive", err.Error())
}
