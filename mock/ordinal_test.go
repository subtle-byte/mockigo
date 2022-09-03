package mock

import (
	"fmt"
)

func Example_ordinal() {
	for i := 1; i < 10; i++ {
		fmt.Print(ordinal(i) + " ")
	}
	fmt.Println(ordinal(10))
	for i := 11; i < 20; i++ {
		fmt.Print(ordinal(i) + " ")
	}
	fmt.Println(ordinal(20))
	for i := 21; i < 30; i++ {
		fmt.Print(ordinal(i) + " ")
	}
	fmt.Println(ordinal(30))
	for i := 101; i < 110; i++ {
		fmt.Print(ordinal(i) + " ")
	}
	fmt.Println(ordinal(110))
	// Output:
	// 1st 2nd 3rd 4th 5th 6th 7th 8th 9th 10th
	// 11th 12th 13th 14th 15th 16th 17th 18th 19th 20th
	// 21st 22nd 23rd 24th 25th 26th 27th 28th 29th 30th
	// 101st 102nd 103rd 104th 105th 106th 107th 108th 109th 110th
}
