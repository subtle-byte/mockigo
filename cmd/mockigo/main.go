package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/subtle-byte/mockigo/internal/dir_walker"
	"github.com/subtle-byte/mockigo/internal/dir_walker/glob"
	"github.com/subtle-byte/mockigo/internal/generator"
	"gopkg.in/yaml.v2"
)

type Config struct {
	RootDir  string `yaml:"root-dir"`
	MocksDir string `yaml:"mocks-dir"`
	Walk     []string
}

func runMockigo(rootDir, mocksDir string, walkRules []string) (err error) {
	rootDir, err = filepath.Abs(rootDir)
	if err != nil {
		return fmt.Errorf("get abs path for root dir: %w", err)
	}
	mocksDir, err = filepath.Abs(mocksDir)
	if err != nil {
		return fmt.Errorf("get abs path for mocks dir: %w", err)
	}

	dirGlob, err := glob.NewGlob(walkRules)
	if err != nil {
		return fmt.Errorf("check the walk rules: %w", err)
	}

	pkgDirs := map[string]generator.Interfaces{}
	dirWalker := dir_walker.NewWalker()
	dirWalker.Walk(rootDir, mocksDir, dirGlob.RootDir, func(dirPath string, interfaces generator.Interfaces) {
		pkgDirs[dirPath] = interfaces
	})

	err = generator.Generate(rootDir, pkgDirs, mocksDir)
	if err != nil {
		return fmt.Errorf("generate: %w", err)
	}
	return nil
}

func main() {
	var rootDir, mocksDir, walkRulesArg string
	flag.StringVar(&rootDir, "root-dir", "internal", "Directory from which the recursive search for go files begins.")
	flag.StringVar(&mocksDir, "mocks-dir", "internal/mocks", "Directory where are the generated mocks saved.")
	flag.StringVar(&walkRulesArg, "walk", "", "Rules of walking finding interfaces")
	// flag.StringVar(&cpuProfileFile, "cpu-profile-file", "", "")
	flag.Parse()

	confFile, err := os.Open("./mockigo.yaml")
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("ERROR: open mockigo.yaml:", err)
		return
	}
	if err == nil {
		defer confFile.Close()
	}
	conf := Config{}
	yaml.NewDecoder(confFile).Decode(&conf)

	var walkRules []string
	if walkRulesArg != "" {
		walkRules = strings.Split(walkRulesArg, ";")
	} else if len(conf.Walk) > 0 {
		walkRules = conf.Walk
	}

	if rootDir == "internal" && conf.RootDir != "" {
		rootDir = conf.RootDir
	}
	if mocksDir == "internal/mocks" && conf.MocksDir != "" {
		mocksDir = conf.MocksDir
	}

	err = runMockigo(rootDir, mocksDir, walkRules)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
}
