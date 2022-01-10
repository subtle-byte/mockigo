package glob

import (
	"fmt"
	"strings"
)

type DirName = string

type Dir struct {
	ExcludeSubdirs bool
	Subdirs        map[DirName]Dir

	ExcludeInterfaces    bool
	InterfacesExceptions map[string]struct{}
}

func newDir() Dir {
	return Dir{
		Subdirs:              map[string]Dir{},
		InterfacesExceptions: map[string]struct{}{},
	}
}

type Glob struct {
	RootDir *Dir
}

func addDir(parentDir *Dir, excludingRule bool, ruleDirs []string, ruleInterfaces []string) {
	if len(ruleDirs) != 0 {
		subdirName := ruleDirs[0]
		subdir, ok := parentDir.Subdirs[subdirName]
		if !ok {
			subdir = newDir()
			subdir.ExcludeSubdirs = parentDir.ExcludeSubdirs
			subdir.ExcludeInterfaces = parentDir.ExcludeInterfaces
		}
		addDir(&subdir, excludingRule, ruleDirs[1:], ruleInterfaces)
		parentDir.Subdirs[subdirName] = subdir
		return
	}
	if len(ruleInterfaces) == 0 {
		*parentDir = newDir()
		parentDir.ExcludeSubdirs = excludingRule
		parentDir.ExcludeInterfaces = excludingRule
		return
	}
	if parentDir.ExcludeInterfaces == excludingRule {
		for _, ruleInterface := range ruleInterfaces {
			delete(parentDir.InterfacesExceptions, ruleInterface)
		}
	} else {
		for _, ruleInterface := range ruleInterfaces {
			parentDir.InterfacesExceptions[ruleInterface] = struct{}{}
		}
	}
}

func NewGlob(rules []string) (*Glob, error) {
	rootDir := newDir()
	for i, rule := range rules {
		rule := strings.TrimSpace(rule)
		dirsAndInterfaces := strings.Split(rule, "@")
		var dirsStr string
		var interfacesStr string
		if len(dirsAndInterfaces) == 1 {
			dirsStr = dirsAndInterfaces[0]
		} else if len(dirsAndInterfaces) == 2 {
			dirsStr = dirsAndInterfaces[0]
			interfacesStr = dirsAndInterfaces[1]
		} else {
			return nil, fmt.Errorf("rule %v: only one @ allowed", i+1)
		}
		var excluding = false
		if strings.HasPrefix(dirsStr, "!") {
			excluding = true
			dirsStr = strings.TrimPrefix(dirsStr, "!")
		}
		if dirsStr == "" {
			return nil, fmt.Errorf("rule %v: use . for root dir", i+1)
		}
		var dirs []string
		if dirsStr != "." {
			dirs = strings.Split(dirsStr, "/")
		}
		var interfaces []string
		if interfacesStr != "" {
			interfaces = strings.Split(interfacesStr, ",")
		}
		addDir(&rootDir, excluding, dirs, interfaces)
	}
	return &Glob{&rootDir}, nil
}
