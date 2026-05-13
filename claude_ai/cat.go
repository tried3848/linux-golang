package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func main() {
	// Определение флагов командной строки
	helpFlag := flag.Bool("h", false, "Показать справку")
	numberFlag := flag.Bool("n", false, "Нумеровать все строки вывода")
	showEndsFlag := flag.Bool("E", false, "Показывать $ в конце каждой строки")
	flag.Parse()

	// Обработка флага помощи
	if *helpFlag {
		fmt.Println("Использование: cat [-h] [-n] [-E] [файл ...]")
		fmt.Println("  -h  Показать справку")
		fmt.Println("  -n  Нумеровать строки")
		fmt.Println("  -E  Показывать символ конца строки $")
		os.Exit(0)
	}

	args := flag.Args()

	// Если файлы не указаны — читаем из stdin
	if len(args) == 0 {
		if err := printFile(os.Stdin, numberFlag, showEndsFlag, new(int)); err != nil {
			fmt.Fprintf(os.Stderr, "cat: ошибка чтения stdin: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Глобальный счётчик строк (сквозная нумерация по всем файлам)
	lineNum := 0
	exitCode := 0

	for _, path := range args {
		// Проверяем существование файла
		info, err := os.Stat(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cat: '%s': нет такого файла или директории\n", path)
			exitCode = 1
			continue
		}
		if info.IsDir() {
			fmt.Fprintf(os.Stderr, "cat: '%s': это директория\n", path)
			exitCode = 1
			continue
		}

		file, err := os.Open(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cat: не удалось открыть '%s': %v\n", path, err)
			exitCode = 1
			continue
		}

		if err := printFile(file, numberFlag, showEndsFlag, &lineNum); err != nil {
			fmt.Fprintf(os.Stderr, "cat: ошибка чтения '%s': %v\n", path, err)
			exitCode = 1
		}
		file.Close()
	}

	os.Exit(exitCode)
}

// printFile читает файл построчно и выводит его содержимое
func printFile(f *os.File, numberFlag *bool, showEndsFlag *bool, lineNum *int) error {
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		*lineNum++

		// Формируем строку вывода
		prefix := ""
		if *numberFlag {
			prefix = fmt.Sprintf("%6d\t", *lineNum)
		}
		suffix := ""
		if *showEndsFlag {
			suffix = "$"
		}
		fmt.Printf("%s%s%s\n", prefix, line, suffix)
	}
	return scanner.Err()
}
