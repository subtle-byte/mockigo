package dir_walker

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/subtle-byte/mockigo/internal/dir_walker/glob"
	"github.com/subtle-byte/mockigo/internal/generator"
	"github.com/subtle-byte/mockigo/internal/util"
)

type DirEntry struct {
	IsDir bool
	Name  string
}

type Walker struct {
	ReadDir func(dirPath string) ([]DirEntry, error)
}

// ReadDir does the same as os.ReadDir, but does not sort
func ReadDir(dirPath string) ([]DirEntry, error) {
	f, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	dirEntries, err := f.ReadDir(-1)
	if err != nil {
		return nil, err
	}
	return util.MapSlice(dirEntries, func(entry fs.DirEntry) DirEntry {
		return DirEntry{
			IsDir: entry.IsDir(),
			Name:  entry.Name(),
		}
	}), nil
}

func NewWalker() *Walker {
	return &Walker{
		ReadDir: ReadDir,
	}
}

func (w *Walker) Walk(dirPath, forbiddenDirPath string, globDir *glob.Dir, visitor func(dirPath string, interfaces generator.Interfaces)) {
	if dirPath == forbiddenDirPath {
		return
	}
	if globDir.ExcludeInterfaces && globDir.ExcludeSubdirs && len(globDir.InterfacesExceptions) == 0 && len(globDir.Subdirs) == 0 {
		return
	}
	dirEntries, err := w.ReadDir(dirPath)
	if err != nil {
		fmt.Println("ERROR: read dir:", err)
		return
	}
	nonDir := false
	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir {
			subdir, ok := globDir.Subdirs[dirEntry.Name]
			if !ok && globDir.ExcludeSubdirs {
				continue
			}
			subdirPath := filepath.Join(dirPath, dirEntry.Name)
			w.Walk(subdirPath, forbiddenDirPath, &subdir, visitor)
		} else {
			nonDir = true
		}
	}
	if nonDir && !(globDir.ExcludeInterfaces && len(globDir.InterfacesExceptions) == 0) {
		visitor(dirPath, generator.Interfaces{
			IncludeInterfaces:    !globDir.ExcludeInterfaces,
			InterfacesExceptions: globDir.InterfacesExceptions,
		})
	}
}
