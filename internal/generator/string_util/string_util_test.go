package string_util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCamelToSnake(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(CamelToSnake("wordWordWORD2WordWORD"), "word_word_word2_word_word")
	assert.Equal(CamelToSnake("WordWordWORD2WordWORD"), "word_word_word2_word_word")
}

func TestCountMaxConsecutiveUnderscores(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(CountMaxConsecutiveUnderscores("skjfiudf", 0), 0)
	assert.Equal(CountMaxConsecutiveUnderscores("skjf_iudf__", 0), 2)
	assert.Equal(CountMaxConsecutiveUnderscores("skjf___iudf_", 0), 3)
	assert.Equal(CountMaxConsecutiveUnderscores("skjf_iudf_", 5), 5)
}
