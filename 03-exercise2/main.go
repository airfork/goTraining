package main

import (
	"fmt"
	"sort"
)

func main() {
	type people []string
	studygroup := people{"Zeno", "Al", "Cody", "Sam"}

	otherpeople := []string{"Zeno", "Al", "Cody", "Sam"}

	n := []int{10, 20, 24, 12, 100, 300, 1000, 1230, 233, 40}
	sort.Strings(studygroup)
	sort.Strings(otherpeople)
	sort.Ints(n)

	fmt.Println("Sorted people slice:", studygroup)
	fmt.Println("Sorted other slice:", otherpeople)
	fmt.Println("Sorted int slice:", n)
}
