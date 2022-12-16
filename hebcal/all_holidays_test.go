package hebcal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllHolidaysForYear(t *testing.T) {
	assert.Equal(t, 126, len(getAllHolidaysForYear(5783)))
}
