package dir_walker

import (
	"fmt"
	"sort"
)

func ExampleReadDir_good() {
	entries, err := ReadDir("testdata")
	// sort the issues after they come out so that the Output error is equivalent, irrelevant of system
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name > entries[j].Name
	})
	fmt.Println("read dir error:", err)
	fmt.Printf("%+v\n", entries)
	//Output:
	// read dir error: <nil>
	// [{IsDir:true Name:sub} {IsDir:false Name:1.txt}]
}

func ExampleReadDir_bad() {
	_, err := ReadDir("dir_that_does_not_exist")
	fmt.Println("read dir error:", err)
	//Output:
	// read dir error: open dir_that_does_not_exist: no such file or directory
}
