package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	// Читаем информацию о памяти из /proc/meminfo
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	// Выводим нужные поля
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.Contains(line, "MemTotal:") ||
			strings.Contains(line, "MemFree:") ||
			strings.Contains(line, "MemAvailable:") ||
			strings.Contains(line, "Buffers:") ||
			strings.Contains(line, "Cached:") ||
			strings.Contains(line, "SwapTotal:") ||
			strings.Contains(line, "SwapFree:") {
			fmt.Println(line)
		}
	}
}
