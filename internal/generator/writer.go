package generator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

type writer struct {
	buf *bytes.Buffer
	err error
}

func (w *writer) Print(args ...string) {
	for _, arg := range args {
		_, err := w.buf.WriteString(arg)
		if err != nil && w.err == nil {
			w.err = err
		}
	}
}

func (w *writer) Println(args ...string) {
	w.Print(args...)
	_, err := w.buf.WriteString("\n")
	if err != nil && w.err == nil {
		w.err = err
	}
}

func writeToFile(buf []byte, filePath string) error {
	mockDir := filepath.Dir(filePath)
	err := os.MkdirAll(mockDir, 0o777)
	if err != nil {
		return fmt.Errorf("create dirs to mock (mock dir = %s): %w", mockDir, err)
	}

	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("create mock file %q: %w", filePath, err)
	}
	defer f.Close()

	_, err = f.Write(buf)
	if err != nil {
		return fmt.Errorf("write into mock file %q: %w", filePath, err)
	}
	return nil
}
