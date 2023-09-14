package util

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMapSlice(t *testing.T) {
	s := MapSlice([]int{12, 28}, func(n int) string {
		return strconv.Itoa(n)
	})
	require.Equal(t, []string{"12", "28"}, s)
}

func TestSliceToSet(t *testing.T) {
	set := SliceToSet([]int{12, 28})
	require.Equal(t, map[int]struct{}{12: struct{}{}, 28: struct{}{}}, set)
}
