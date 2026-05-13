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
	startFlag := flag.Int("v", 1, "Начальный номер строки")
	widthFlag := flag.Int("w", 6, "Ширина поля номера строки")
	flag.Parse()

	// Обработка флага помощи
	if *helpFlag {
		fmt.Println("Использование: nl [-h] [-v N] [-w N] [файл ...]")
		fmt.Println("  -h     Показать справку")
		fmt.Println("  -v N   Начальный номер строки (по умолчанию 1)")
		fmt.Println("  -w N   Ширина поля номера (по умолчанию 6)")
		os.Exit(0)
	}

	// Проверяем корректность входных значений
	if *widthFlag < 1 || *widthFlag > 20 {
		fmt.Fprintln(os.Stderr, "nl: ширина поля (-w) должна быть от 1 до 20")
		os.Exit(1)
	}
	if *startFlag < 0 {
		fmt.Fprintln(os.Stderr, "nl: начальный номер (-v) не может быть отрицательным")
		os.Exit(1)
	}

	args := flag.Args()
	lineNum := *startFlag

	// Если файлы не указаны — читаем из stdin
	if len(args) == 0 {
		numberLines(os.Stdin, &lineNum, *widthFlag)
		return
	}

	exitCode := 0
	for _, path := range args {
		info, err := os.Stat(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "nl: '%s': нет такого файла\n", path)
			exitCode = 1
			continue
		}
		if info.IsDir() {
			fmt.Fprintf(os.Stderr, "nl: '%s': это директория\n", path)
			exitCode = 1
			continue
		}

		f, err := os.Open(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "nl: не удалось открыть '%s': %v\n", path, err)
			exitCode = 1
			continue
		}
		numberLines(f, &lineNum, *widthFlag)
		f.Close()
	}
	os.Exit(exitCode)
}

// numberLines читает файл и выводит строки с номерами
func numberLines(f *os.File, lineNum *int, width int) {
	// Формат: пробелы + номер + TAB + строка
	format := fmt.Sprintf("%%%dd\t%%s\n", width)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fmt.Printf(format, *lineNum, scanner.Text())
		*lineNum++
	}
}
