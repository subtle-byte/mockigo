package generator

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerator(t *testing.T) {
	toAbs := func(rel string) string {
		abs, err := filepath.Abs(rel)
		require.NoError(t, err)
		return abs
	}
	interfaces := map[string]Interfaces{
		toAbs("testdata"): {
			IncludeInterfaces:    true,
			InterfacesExceptions: map[string]struct{}{"Filtered": {}},
		},
	}
	err := Generate(toAbs("testdata"), interfaces, toAbs("testdata/mocks"))
	require.NoError(t, err)
}
