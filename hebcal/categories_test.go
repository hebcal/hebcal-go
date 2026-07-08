package hebcal

import (
	"reflect"
	"testing"
	"time"

	"github.com/hebcal/hdate"
	"github.com/hebcal/hebcal-go/event"
)

func TestTimedEventGetCategories(t *testing.T) {
	hd := hdate.New(5784, hdate.Nisan, 14)
	when := time.Date(2024, time.April, 22, 19, 0, 0, 0, time.UTC)
	opts := &CalOptions{}
	tests := []struct {
		desc string
		want []string
	}{
		{"Candle lighting", []string{"candles"}},
		{"Havdalah", []string{"havdalah"}},
		{"Fast begins", []string{"zmanim", "fast"}},
		{"Fast ends", []string{"zmanim", "fast"}},
		{"Finish eating chametz", []string{"zmanim", "achilasChametz"}},
		{"Biur Chametz", []string{"zmanim", "biurChametz"}},
	}
	for _, tc := range tests {
		ev := NewTimedEvent(hd, tc.desc, event.LIGHT_CANDLES, when, 0, nil, opts)
		if got := ev.GetCategories(); !reflect.DeepEqual(got, tc.want) {
			t.Errorf("%q: GetCategories() = %v, want %v", tc.desc, got, tc.want)
		}
	}
}
