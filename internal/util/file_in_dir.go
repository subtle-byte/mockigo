package util

import (
	"fmt"
	"path/filepath"
)

func FileInDir(dirPath, filePath string) (bool, error) {
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return false, fmt.Errorf("get absolute path of %q: %w", filePath, err)
	}
	absDirPath, err := filepath.Abs(dirPath)
	if err != nil {
		return false, fmt.Errorf("get absolute path of %q: %w", absDirPath, err)
	}
	return filepath.Dir(absFilePath) == absDirPath, nil
}
