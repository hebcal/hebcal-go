package hebcal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllHolidaysForYear(t *testing.T) {
	assert.Equal(t, 127, len(getAllHolidaysForYear(5783)))
}
