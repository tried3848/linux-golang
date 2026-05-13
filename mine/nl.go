package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	filename := os.Args[1]
	data, _ := os.ReadFile(filename)
	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		fmt.Printf("%6d  %s\n", i+1, line)
	}
}
