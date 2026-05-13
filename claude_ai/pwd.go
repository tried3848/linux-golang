package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// Определение флагов командной строки
	helpFlag := flag.Bool("h", false, "Показать справку")
	logicalFlag := flag.Bool("L", false, "Использовать логический путь (с симлинками)")
	physicalFlag := flag.Bool("P", false, "Использовать физический путь (без симлинков)")
	flag.Parse()

	// Обработка флага помощи
	if *helpFlag {
		fmt.Println("Использование: pwd [-h] [-L] [-P]")
		fmt.Println("  -h  Показать справку")
		fmt.Println("  -L  Логический путь (по умолчанию)")
		fmt.Println("  -P  Физический путь (разрешает симлинки)")
		os.Exit(0)
	}

	// Получаем текущую директорию
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "pwd: ошибка получения рабочей директории: %v\n", err)
		os.Exit(1)
	}

	// Если запрошен физический путь — разрешаем симлинки
	if *physicalFlag && !*logicalFlag {
		resolved, err := filepath.EvalSymlinks(dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "pwd: ошибка разрешения симлинков: %v\n", err)
			os.Exit(1)
		}
		dir = resolved
	}

	fmt.Println(dir)
}
