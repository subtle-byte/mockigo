package path_trie

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	assert := assert.New(t)
	trie := New()
	trie.LoadPath(strings.Split("github.com/aa/aa/cc", "/"), 6)
	trie.LoadPath(strings.Split("github.com/aa/cc", "/"), 2)
	trie.LoadPath(strings.Split("github.com/aa", "/"), nil)
	reducedPaths := trie.ReducedPaths()
	assert.Len(reducedPaths, 3)
	assert.Contains(reducedPaths, ReducedPath{
		Rest: strings.Split("aa/aa/cc", "/"), Meta: 6,
	})
	assert.Contains(reducedPaths, ReducedPath{
		Rest: strings.Split("github.com/aa/cc", "/"), Meta: 2,
	})
	assert.Contains(reducedPaths, ReducedPath{
		Rest: strings.Split("aa", "/"), Meta: nil,
	})
}
