package sedra

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidDouble(t *testing.T) {
	assert.Equal(t, true, isValidDouble(-26))
	assert.Equal(t, false, isValidDouble(-1))
	assert.Equal(t, false, isValidDouble(26))
}
