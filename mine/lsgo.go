package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	entries, err := os.ReadDir(".")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	
	// Собираем имена в слайс
	names := make([]string, len(entries))
	for i, entry := range entries {
		names[i] = entry.Name()
	}
	
	// Выводим через запятую
	fmt.Println(strings.Join(names, ", "))
}
