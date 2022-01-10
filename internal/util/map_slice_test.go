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

func TestMapSliceWithIndex(t *testing.T) {
	s := MapSliceWithIndex([]int{12, 28}, func(i, n int) string {
		return strconv.Itoa(i) + "_" + strconv.Itoa(n)
	})
	require.Equal(t, []string{"0_12", "1_28"}, s)
}
