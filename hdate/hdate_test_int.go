package hdate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestElapsedDays(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(2110760, elapsedDays(5780))
	assert.Equal(2084447, elapsedDays(5708))
	assert.Equal(1373677, elapsedDays(3762))
	assert.Equal(1340455, elapsedDays(3671))
	assert.Equal(450344, elapsedDays(1234))
	assert.Equal(44563, elapsedDays(123))
	assert.Equal(356, elapsedDays(2))
	assert.Equal(1, elapsedDays(1))
	assert.Equal(2104174, elapsedDays(5762))
	assert.Equal(2104528, elapsedDays(5763))
	assert.Equal(2104913, elapsedDays(5764))
	assert.Equal(2105268, elapsedDays(5765))
	assert.Equal(2105651, elapsedDays(5766))
}
