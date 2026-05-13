package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

func main() {
	// Определяем флаги
	cflag := flag.Bool("c", false, "напечатать количество байт")
	nflag := flag.Bool("n", false, "напечатать количество символов")
	hflag := flag.Bool("h", false, "справка о команде wc")

	flag.Parse()

	if *hflag {
		printHelp()
		return
	}
	
	// Получаем имя файла из аргументов
	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Ошибка: не указан файл")
		printHelp()
		return
	}
	
	filename := args[0]
	
	// Вызываем соответствующую функцию в зависимости от флага
	if *cflag {
		wcCflag(filename)
	} else if *nflag {
		wcNflag(filename)
	} else {
		wcnoflag(filename)
	}
}

func wcnoflag(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка чтения файла %s: %v\n", filename, err)
		return
	}
	content := string(data)
	lines := strings.Count(content, "\n")
	words := len(strings.Fields(content))
	bytes := len(content)
	// Используем Printf с форматированием, а не Println
	fmt.Printf("%d\t%d\t%d\t%s\n", lines, words, bytes, filename)
}

func wcCflag(filename string) {
	// -c напечатать количество байт
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка чтения файла %s: %v\n", filename, err)
		return
	}
	bytes := len(data)
	fmt.Printf("%d\t%s\n", bytes, filename)
}

func wcNflag(filename string) {
	// -n напечатать количество символов (c учетом Unicode)
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка чтения файла %s: %v\n", filename, err)
		return
	}
	content := string(data)
	// Подсчитываем количество рун (символов Unicode)
	chars := utf8.RuneCountInString(content)
	fmt.Printf("%d\t%s\n", chars, filename)
}

func printHelp() {
	fmt.Println(`Использование: wc [флаг] [файл]
	
Флаги:
  -c    напечатать количество байт
  -n    напечатать количество символов
  -h    справка по команде и параметрам

Примеры:
  wc file.txt         # выводит строки, слова, байты и имя файла
  wc -c file.txt      # выводит только количество байт
  wc -n file.txt      # выводит только количество символов
  wc -h               # выводит справку`)
}