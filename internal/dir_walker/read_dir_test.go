package dir_walker

import "fmt"

func ExampleReadDir_good() {
	entries, err := ReadDir("testdata")
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
