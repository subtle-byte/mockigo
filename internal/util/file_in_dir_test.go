package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSameDir(t *testing.T) {
	in, err := FileInDir("..", "file.go")
	assert.NoError(t, err)
	assert.Equal(t, false, in)

	in, err = FileInDir(".", "./file.go")
	assert.NoError(t, err)
	assert.Equal(t, true, in)
}
