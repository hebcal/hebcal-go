package hebcal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllHolidaysForYear(t *testing.T) {
	// The raw table is pre-filter, so holidays observed on both schedules are
	// listed twice, once CHUL_ONLY and once IL_ONLY. Erev Sukkot and Erev
	// Pesach are among them.
	assert.Equal(t, 135, len(getAllHolidaysForYear(5783)))
}
