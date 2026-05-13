package main

import (
	"fmt"
	"os"
)

func main() {
	filename := os.Args[1]
	data, _ := os.ReadFile(filename)
	fmt.Print(string(data))
}
