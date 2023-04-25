package generator

import (
	"github.com/subtle-byte/mockigo/cmd/mockigo/config"
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
	err := Generate(&config.InitializedConfig{Config: &config.Config{RootDir: "testdata", MocksDir: "testdata/mocks"}, RootDir: toAbs("testdata"), MocksDir: toAbs("testdata/mocks")}, interfaces)
	require.NoError(t, err)
}
