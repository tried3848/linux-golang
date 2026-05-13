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
	n := 10 // первые 10 строк
	for i := 0; i < n && i < len(lines); i++ {
		fmt.Println(lines[i])
	}
}
