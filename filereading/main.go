package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Please only pass in one file as a parameter")
		os.Exit(1)
	}

	filename, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(2)
	}
	io.Copy(os.Stdout, filename)
}
