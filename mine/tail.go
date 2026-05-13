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
	n := 10
	for i := len(lines) - n; i < len(lines); i++ {
		if i >= 0 {
			fmt.Println(lines[i])
		}
	}

}
