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
	targets := Targets{
		Include:    true,
		Exceptions: map[string]struct{}{"Filtered": {}},
	}
	err := Generate(Config{
		TargetPkgDirPath: toAbs("testdata"),
		Targets:          targets,
		OutFilePath:      toAbs("testdata/mocks.go"),
		OutPkgName: func(inspectedPkgName string) string {
			return inspectedPkgName + "_test"
		},
		OutPublic: true,
		GoGenCmd:  "some command",
	})
	require.NoError(t, err)
}
