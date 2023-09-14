package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetCaseForFirstLetter(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(SetCaseForFirstLetter("", true), "")
	assert.Equal(SetCaseForFirstLetter("", false), "")
	assert.Equal(SetCaseForFirstLetter("game", true), "Game")
	assert.Equal(SetCaseForFirstLetter("Game", true), "Game")
	assert.Equal(SetCaseForFirstLetter("game", false), "game")
	assert.Equal(SetCaseForFirstLetter("Game", false), "game")
}

func TestCountMaxConsecutiveUnderscores(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(CountMaxConsecutiveUnderscores("skjfiudf", 0), 0)
	assert.Equal(CountMaxConsecutiveUnderscores("skjf_iudf__", 0), 2)
	assert.Equal(CountMaxConsecutiveUnderscores("skjf___iudf_", 0), 3)
	assert.Equal(CountMaxConsecutiveUnderscores("skjf_iudf_", 5), 5)
}
