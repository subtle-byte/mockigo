package glob

import (
	"fmt"
	"sort"
	"strings"

	"golang.org/x/exp/maps"
)

func printDirs(dir *Dir, tab string) {
	fmt.Println(tab+"ExcludeInterfaces:", dir.ExcludeInterfaces)
	interfacesExceptions := maps.Keys(dir.InterfacesExceptions)
	sort.Strings(interfacesExceptions)
	fmt.Println(tab + "InterfacesExceptions: [" + strings.Join(interfacesExceptions, ", ") + "]")
	fmt.Println(tab+"ExcludeSubdirs:", dir.ExcludeSubdirs)
	fmt.Println(tab + "Subdirs:")
	subdirsNames := maps.Keys(dir.Subdirs)
	sort.Strings(subdirsNames)
	for _, subdirName := range subdirsNames {
		subdir := dir.Subdirs[subdirName]
		fmt.Println(tab + "├─Name: " + subdirName)
		printDirs(&subdir, tab+"│ ")
	}
}

func ExampleGlob_simple() {
	glob, err := NewGlob([]string{})
	fmt.Println("Glob create error:", err)
	fmt.Println("Glob tree:")
	printDirs(glob.RootDir, "  ")
	// Output:
	// Glob create error: <nil>
	// Glob tree:
	//   ExcludeInterfaces: false
	//   InterfacesExceptions: []
	//   ExcludeSubdirs: false
	//   Subdirs:
}

func ExampleGlob_root() {
	_, err := NewGlob([]string{
		"@A",
	})
	fmt.Println("Glob create error:", err)

	_, err = NewGlob([]string{
		"!",
	})
	fmt.Println("Glob create error:", err)

	glob, err := NewGlob([]string{
		"!.",
		".@A,B",
	})
	fmt.Println("Glob create error:", err)
	fmt.Println("Glob tree:")
	printDirs(glob.RootDir, "  ")
	// Output:
	// Glob create error: rule 1: use . for root dir
	// Glob create error: rule 1: use . for root dir
	// Glob create error: <nil>
	// Glob tree:
	//   ExcludeInterfaces: true
	//   InterfacesExceptions: [A, B]
	//   ExcludeSubdirs: true
	//   Subdirs:
}

func ExampleGlob_tooManyAtSigns() {
	_, err := NewGlob([]string{
		"abc@kj@ie",
	})
	fmt.Println("Glob create error:", err)
	// Output:
	// Glob create error: rule 1: only one @ allowed
}

func ExampleGlob_excludeSome() {
	glob, err := NewGlob([]string{
		"!proto",
		"proto@A,B",
		"proto/good",
		"!proto/good@C",
	})
	fmt.Println("Glob create error:", err)
	fmt.Println("Glob tree:")
	printDirs(glob.RootDir, "  ")
	//Output:
	// Glob create error: <nil>
	// Glob tree:
	//   ExcludeInterfaces: false
	//   InterfacesExceptions: []
	//   ExcludeSubdirs: false
	//   Subdirs:
	//   ├─Name: proto
	//   │ ExcludeInterfaces: true
	//   │ InterfacesExceptions: [A, B]
	//   │ ExcludeSubdirs: true
	//   │ Subdirs:
	//   │ ├─Name: good
	//   │ │ ExcludeInterfaces: false
	//   │ │ InterfacesExceptions: [C]
	//   │ │ ExcludeSubdirs: false
	//   │ │ Subdirs:
}

func ExampleGlob_includeSome() {
	glob, err := NewGlob([]string{
		"!.",
		"internal/service@A,B",
		"internal/server",
		"!internal/server@BadInterface",
	})
	fmt.Println("Glob create error:", err)
	fmt.Println("Glob tree:")
	printDirs(glob.RootDir, "  ")
	// Output:
	// Glob create error: <nil>
	// Glob tree:
	//   ExcludeInterfaces: true
	//   InterfacesExceptions: []
	//   ExcludeSubdirs: true
	//   Subdirs:
	//   ├─Name: internal
	//   │ ExcludeInterfaces: true
	//   │ InterfacesExceptions: []
	//   │ ExcludeSubdirs: true
	//   │ Subdirs:
	//   │ ├─Name: server
	//   │ │ ExcludeInterfaces: false
	//   │ │ InterfacesExceptions: [BadInterface]
	//   │ │ ExcludeSubdirs: false
	//   │ │ Subdirs:
	//   │ ├─Name: service
	//   │ │ ExcludeInterfaces: true
	//   │ │ InterfacesExceptions: [A, B]
	//   │ │ ExcludeSubdirs: true
	//   │ │ Subdirs:
}

func ExampleGlob_revertingRule() {
	glob, err := NewGlob([]string{
		"!abc@A,B,C,D",
		"abc@C,B",
	})
	fmt.Println("Glob create error:", err)
	fmt.Println("Glob tree:")
	printDirs(glob.RootDir, "  ")
	// Output:
	// Glob create error: <nil>
	// Glob tree:
	//   ExcludeInterfaces: false
	//   InterfacesExceptions: []
	//   ExcludeSubdirs: false
	//   Subdirs:
	//   ├─Name: abc
	//   │ ExcludeInterfaces: false
	//   │ InterfacesExceptions: [A, D]
	//   │ ExcludeSubdirs: false
	//   │ Subdirs:
}
