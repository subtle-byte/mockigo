package dir_walker

import (
	"fmt"
	"sort"

	"github.com/subtle-byte/mockigo/internal/dir_walker/glob"
	"github.com/subtle-byte/mockigo/internal/generator"
	"golang.org/x/exp/maps"
)

func File(name string) DirEntry {
	return DirEntry{
		IsDir: false,
		Name:  name,
	}
}

func Dir(name string) DirEntry {
	return DirEntry{
		IsDir: true,
		Name:  name,
	}
}

var virtualFS = map[string][]DirEntry{
	".":                           {Dir("internal")},
	"internal":                    {Dir("service"), Dir("generated"), Dir("generated2"), Dir("mocks"), Dir("empty")},
	"internal/service":            {File("some.go"), Dir("subservice")},
	"internal/service/subservice": {File("some2.go")},
	"internal/generated":          {Dir("good"), Dir("bad"), File("some7.go")},
	"internal/generated/good":     {File("some3.go")},
	"internal/generated/bad":      {File("some4.go")},
	"internal/generated2":         {File("some8.go")},
	"internal/mocks":              {File("some5.go"), File("some6")},
	"internal/empty":              {},
}

func fakeReadDir(dirPath string) ([]DirEntry, error) {
	fmt.Printf("read dir %s\n", dirPath)
	e, ok := virtualFS[dirPath]
	if !ok {
		return nil, fmt.Errorf("%q not ok", dirPath)
	}
	return e, nil
}

func visit(dirPath string, interfaces generator.Interfaces) {
	includeInterfaces := "all"
	if !interfaces.IncludeInterfaces {
		includeInterfaces = "non"
	}
	exceptions := maps.Keys(interfaces.InterfacesExceptions)
	sort.Strings(exceptions)
	fmt.Printf("visit %s: interfaces: %s but %v\n", dirPath, includeInterfaces, exceptions)
}

func ExampleWalk_walkAll() {
	glob, err := glob.NewGlob([]string{})
	fmt.Println("glob create error:", err)

	walker := Walker{ReadDir: fakeReadDir}
	walker.Walk(".", "something-invalid", glob.RootDir, visit)

	//Output:
	// glob create error: <nil>
	// read dir .
	// read dir internal
	// read dir internal/service
	// read dir internal/service/subservice
	// visit internal/service/subservice: interfaces: all but []
	// visit internal/service: interfaces: all but []
	// read dir internal/generated
	// read dir internal/generated/good
	// visit internal/generated/good: interfaces: all but []
	// read dir internal/generated/bad
	// visit internal/generated/bad: interfaces: all but []
	// visit internal/generated: interfaces: all but []
	// read dir internal/generated2
	// visit internal/generated2: interfaces: all but []
	// read dir internal/mocks
	// visit internal/mocks: interfaces: all but []
	// read dir internal/empty
}

func ExampleWalk_walkSome() {
	glob, err := glob.NewGlob([]string{
		"!internal/service@BadInterface",
		"!internal/service/subservice",
		"internal/service/subservice@GoodInterface",
		"!internal/generated",
		"internal/generated/good",
		"!internal/generated2",
	})
	fmt.Println("glob create error:", err)

	walker := Walker{ReadDir: fakeReadDir}
	walker.Walk(".", "internal/mocks", glob.RootDir, visit)

	//Output:
	// glob create error: <nil>
	// read dir .
	// read dir internal
	// read dir internal/service
	// read dir internal/service/subservice
	// visit internal/service/subservice: interfaces: non but [GoodInterface]
	// visit internal/service: interfaces: all but [BadInterface]
	// read dir internal/generated
	// read dir internal/generated/good
	// visit internal/generated/good: interfaces: all but []
	// read dir internal/empty
}

func ExampleWalk_walkBad() {
	glob, err := glob.NewGlob([]string{})
	fmt.Println("glob create error:", err)

	walker := Walker{ReadDir: fakeReadDir}
	walker.Walk("dir_that_does_not_exist", "internal/mocks", glob.RootDir, visit)

	//Output:
	// glob create error: <nil>
	// read dir dir_that_does_not_exist
	// ERROR: read dir: "dir_that_does_not_exist" not ok
}
