package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	filename := os.Args[1]
	data, _ := os.ReadFile(filename)
	content := string(data)
	lines := strings.Count(content, "\n")
	words := len(strings.Fields(content))
	chars := len(content)
	fmt.Println("%s		%s		%s		%s", lines, words, chars, filename)
}
