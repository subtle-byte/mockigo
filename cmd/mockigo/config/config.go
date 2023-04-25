package config

import (
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/alecthomas/kong-yaml"
	"github.com/subtle-byte/mockigo/internal/dir_walker/glob"
	"path/filepath"
)

type Config struct {
	RootDir    string   `help:"where to start walking and looking for interfaces" default:"internal"`
	MocksDir   string   `help:"Where to place mocks in a directory(only when not in-dir)" default:"internal/mocks"`
	PkgPrefix  string   `help:"prefix for mock packages(only when in-dir)"`
	PkgPostfix string   `help:"postfix for mock packages(only when in-dir)" default:"_mock"`
	InDir      bool     `help:"whether the generation should happen in the pkg dir[ignores mocksDir]" default:"false"`
	Walk       []string `help:"walking rules for ignoring file(s)"`
}

type InitializedConfig struct {
	*Config
	RootDir  string `yaml:"root-dir"`
	MocksDir string `yaml:"mocks-dir"`
	WalkGlob *glob.Glob
}

type Cli struct {
	Config
}

func KongConfig() (*InitializedConfig, *kong.Context, error) {
	cli := &Cli{}
	ctx := kong.Parse(cli, kong.Configuration(kongyaml.Loader, "./mockigo.yaml"))
	cnf, err := cli.Init()
	return cnf, ctx, err
}

func (cnf *Config) Init() (*InitializedConfig, error) {
	retVal := &InitializedConfig{Config: cnf}
	var err error

	retVal.RootDir, err = filepath.Abs(cnf.RootDir)
	if err != nil {
		return nil, fmt.Errorf("get abs path for root dir: %w", err)
	}
	retVal.MocksDir, err = filepath.Abs(cnf.MocksDir)
	if err != nil {
		return nil, fmt.Errorf("get abs path for mocks dir: %w", err)
	}

	retVal.WalkGlob, err = glob.NewGlob(cnf.Walk)
	if err != nil {
		return nil, fmt.Errorf("check the walk rules: %w", err)
	}
	return retVal, nil
}
