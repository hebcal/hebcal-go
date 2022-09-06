package hebcal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGematriya(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("תשמ״ט", Gematriya(5749))
	assert.Equal("תשע״ד", Gematriya(5774))
	assert.Equal("תש״פ", Gematriya(5780))
	assert.Equal("ג׳", Gematriya(3))
	assert.Equal("י״ד", Gematriya(14))
	assert.Equal("ט״ו", Gematriya(15))
	assert.Equal("ט״ז", Gematriya(16))
	assert.Equal("י״ז", Gematriya(17))
	assert.Equal("כ׳", Gematriya(20))
	assert.Equal("כ״ה", Gematriya(25))
	assert.Equal("ס׳", Gematriya(60))
	assert.Equal("קכ״ג", Gematriya(123))
	assert.Equal("תרי״ג", Gematriya(613))
	assert.Equal("תשמ״ט", Gematriya(5749))
	assert.Equal("ג׳תשס״א", Gematriya(3761))
	assert.Equal("ו׳תשמ״ט", Gematriya(6749))
	assert.Equal("ח׳תשס״ה", Gematriya(8765))
	assert.Equal("כב׳ת״ש", Gematriya(22700))
	assert.Equal("טז׳קכ״ג", Gematriya(16123))
	assert.Equal("א׳קכ״ג", Gematriya(1123))
	assert.Equal("ו׳", Gematriya(6000))
	assert.Equal("ז׳ז׳", Gematriya(7007))
}
