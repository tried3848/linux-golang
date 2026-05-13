package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// Определение флагов командной строки
	helpFlag := flag.Bool("h", false, "Показать справку")
	parentsFlag := flag.Bool("p", false, "Создать родительские директории при необходимости")
	verboseFlag := flag.Bool("v", false, "Выводить сообщение о каждой созданной директории")
	flag.Parse()

	// Обработка флага помощи
	if *helpFlag {
		fmt.Println("Использование: mkdir [-h] [-p] [-v] <директория> [директория2 ...]")
		fmt.Println("  -h  Показать справку")
		fmt.Println("  -p  Создать родительские директории")
		fmt.Println("  -v  Verbose: сообщать о каждом создании")
		os.Exit(0)
	}

	// Проверяем наличие аргументов
	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "mkdir: не указана директория")
		fmt.Fprintln(os.Stderr, "Используйте 'mkdir -h' для справки")
		os.Exit(1)
	}

	// Создаём каждую директорию из аргументов
	exitCode := 0
	for _, dir := range args {
		// Проверяем, что имя директории не пустое
		if dir == "" {
			fmt.Fprintln(os.Stderr, "mkdir: имя директории не может быть пустым")
			exitCode = 1
			continue
		}

		var err error
		if *parentsFlag {
			// Создаём всю цепочку директорий
			err = os.MkdirAll(dir, 0755)
		} else {
			// Создаём только одну директорию
			err = os.Mkdir(dir, 0755)
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "mkdir: невозможно создать директорию '%s': %v\n", dir, err)
			exitCode = 1
			continue
		}

		// Подробный вывод
		if *verboseFlag {
			fmt.Printf("mkdir: создана директория '%s'\n", dir)
		}
	}

	os.Exit(exitCode)
}
