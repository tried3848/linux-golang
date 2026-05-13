package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// Определение флагов командной строки
	helpFlag := flag.Bool("h", false, "Показать справку")
	parentsFlag := flag.Bool("p", false, "Удалить директорию и пустые родительские директории")
	verboseFlag := flag.Bool("v", false, "Выводить сообщение об удалении каждой директории")
	flag.Parse()

	// Обработка флага помощи
	if *helpFlag {
		fmt.Println("Использование: rmdir [-h] [-p] [-v] <директория> [директория2 ...]")
		fmt.Println("  -h  Показать справку")
		fmt.Println("  -p  Удалить также пустые родительские директории")
		fmt.Println("  -v  Verbose: сообщать об удалении")
		os.Exit(0)
	}

	// Проверяем наличие аргументов
	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "rmdir: не указана директория")
		fmt.Fprintln(os.Stderr, "Используйте 'rmdir -h' для справки")
		os.Exit(1)
	}

	exitCode := 0
	for _, dir := range args {
		if dir == "" {
			fmt.Fprintln(os.Stderr, "rmdir: имя директории не может быть пустым")
			exitCode = 1
			continue
		}

		// Проверяем, что путь существует и является директорией
		info, err := os.Stat(dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "rmdir: не удалось получить информацию о '%s': %v\n", dir, err)
			exitCode = 1
			continue
		}
		if !info.IsDir() {
			fmt.Fprintf(os.Stderr, "rmdir: '%s' не является директорией\n", dir)
			exitCode = 1
			continue
		}

		// Пытаемся удалить директорию (удаляет только пустую)
		if err := os.Remove(dir); err != nil {
			fmt.Fprintf(os.Stderr, "rmdir: не удалось удалить '%s': %v\n", dir, err)
			exitCode = 1
			continue
		}

		if *verboseFlag {
			fmt.Printf("rmdir: удалена директория '%s'\n", dir)
		}

		// Если флаг -p: пытаемся удалить родительские директории
		if *parentsFlag {
			parent := dir
			for {
				// Находим родительскую директорию
				idx := len(parent) - 1
				for idx > 0 && parent[idx] != '/' && parent[idx] != '\\' {
					idx--
				}
				if idx == 0 {
					break
				}
				parent = parent[:idx]
				if parent == "" {
					break
				}
				if err := os.Remove(parent); err != nil {
					// Родительская директория не пустая — прекращаем
					break
				}
				if *verboseFlag {
					fmt.Printf("rmdir: удалена директория '%s'\n", parent)
				}
			}
		}
	}

	os.Exit(exitCode)
}
