package main

import (
	"fmt"
	"github.com/subtle-byte/mockigo/cmd/mockigo/config"
	"github.com/subtle-byte/mockigo/internal/dir_walker"
	"github.com/subtle-byte/mockigo/internal/generator"
)

func runMockigo(cnf *config.InitializedConfig) (err error) {
	pkgDirs := map[string]generator.Interfaces{}
	dirWalker := dir_walker.NewWalker()
	dirWalker.Walk(cnf.RootDir, cnf.MocksDir, cnf.WalkGlob.RootDir, func(dirPath string, interfaces generator.Interfaces) {
		pkgDirs[dirPath] = interfaces
	})

	err = generator.Generate(cnf, pkgDirs)
	if err != nil {
		return fmt.Errorf("generate: %w", err)
	}
	return nil
}

func main() {
	cnf, _, err := config.KongConfig()
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	err = runMockigo(cnf)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	return
}
