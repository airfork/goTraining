package main

import "fmt"

func main() {
	foo(1, 2)
	foo(1, 2, 3)
	aSlice := []int{1, 2, 3, 4}
	foo(aSlice...)
	foo()

	//Example of func expression and double return in func
	// half := func(num int) (int, bool) {
	// 	return num / 2, num%2 == 0
	// }
}

//varadic function, with example of range loop
func foo(nums ...int) {
	for _, v := range nums {
		fmt.Println(v)
	}
}
