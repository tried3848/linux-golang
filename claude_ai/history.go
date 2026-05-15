package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func getHistoryFile() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "history: не удалось найти домашнюю директорию: %v\n", err)
		os.Exit(1)
	}
	return filepath.Join(home, ".bash_history")
}

func readHistory() ([]string, error) {
	path := getHistoryFile()
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать файл истории: %v", err)
	}
	lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	return lines, nil
}

func writeHistory(lines []string) error {
	path := getHistoryFile()
	content := strings.Join(lines, "\n") + "\n"
	return os.WriteFile(path, []byte(content), 0600)
}

func showHistory(n int) {
	lines, err := readHistory()
	if err != nil {
		fmt.Fprintln(os.Stderr, "history:", err)
		os.Exit(1)
	}

	start := 0
	if n > 0 && n < len(lines) {
		start = len(lines) - n
	}

	for i := start; i < len(lines); i++ {
		fmt.Printf("%5d  %s\n", i+1, lines[i])
	}
}

func clearHistory() {
	err := writeHistory([]string{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "history: не удалось очистить историю: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("История очищена.")
}

func deleteEntry(offset int) {
	lines, err := readHistory()
	if err != nil {
		fmt.Fprintln(os.Stderr, "history:", err)
		os.Exit(1)
	}

	index := offset - 1 // offset начинается с 1
	if offset < 0 {
		index = len(lines) + offset // отрицательные — с конца
	}

	if index < 0 || index >= len(lines) {
		fmt.Fprintf(os.Stderr, "history: %d: смещение за пределами диапазона\n", offset)
		os.Exit(1)
	}

	lines = append(lines[:index], lines[index+1:]...)
	err = writeHistory(lines)
	if err != nil {
		fmt.Fprintf(os.Stderr, "history: не удалось записать историю: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Запись %d удалена.\n", offset)
}

func main() {
	hflag := flag.Bool("h", false, "справка по команде и параметрам")
	cflag := flag.Bool("c", false, "очистить историю")
	doffset := flag.Int("d", 0, "удалить запись по смещению OFFSET")
	flag.Parse()

	if *hflag {
		printHelp()
		return
	}

	if *cflag {
		clearHistory()
		return
	}

	if *doffset != 0 {
		deleteEntry(*doffset)
		return
	}

	// Показать последние N записей или всю историю
	args := flag.Args()
	n := 0
	if len(args) > 0 {
		val, err := strconv.Atoi(args[0])
		if err != nil || val < 0 {
			fmt.Fprintf(os.Stderr, "history: %s: числовой аргумент требуется\n", args[0])
			os.Exit(1)
		}
		n = val
	}

	showHistory(n)
}

func printHelp() {
	fmt.Println(`Использование: history [параметры]
Показывает историю терминала.

Параметры:
  -h            справка по команде и параметрам
  -c            очистить историю (удалить все записи)
  -d OFFSET     удалить запись истории по смещению OFFSET
                (отрицательные значения отсчитываются с конца)
  [n]           показать только последние N записей (опционально)

Примеры:
  history           - показать всю историю
  history -c        - очистить историю
  history -d 5      - удалить 5-ю запись
  history -d -2     - удалить предпоследнюю запись
  history 10        - показать последние 10 записей`)
}
