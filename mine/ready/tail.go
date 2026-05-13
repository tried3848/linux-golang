package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	qflag := flag.Bool("q", false, "не печатать заголовки с именами файлов")
	vflag := flag.Bool("v", false, "всегда печатать заголовки с именами файлов")
	nflag := flag.Int("n", 10, "количество строк для вывода (по умолчанию 10)")
	hflag := flag.Bool("h", false, "справка о команде tail")

	flag.Parse()
	if *hflag {
		printHelp()
		return
	}
	
	files := flag.Args()
	
	// Если файлов нет, читаем из stdin
	if len(files) == 0 {
		tailStdin(*nflag)
		return
	}
	
	if *qflag {
		tailQflag(files, *nflag)
		return
	}
	if *vflag {
		tailVflag(files, *nflag)
		return
	}
	
	// Без флагов - показываем заголовки только если файлов несколько
	if len(files) > 1 {
		tailVflag(files, *nflag)
	} else {
		tailNoFlags(files[0], *nflag)
	}
}

// Без флагов (один файл)
func tailNoFlags(filename string, n int) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "tail: %s: %v\n", filename, err)
		return
	}
	lines := strings.Split(string(data), "\n")
	
	// Удаляем последнюю пустую строку если есть
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	
	// Выводим последние n строк
	start := len(lines) - n
	if start < 0 {
		start = 0
	}
	for i := start; i < len(lines); i++ {
		fmt.Println(lines[i])
	}
}

// С флагом -v (всегда печатать заголовки)
func tailVflag(files []string, n int) {
	for i, filename := range files {
		if i > 0 {
			fmt.Println() // пустая строка между файлами
		}
		fmt.Printf("==> %s <==\n", filename)
		
		data, err := os.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "tail: %s: %v\n", filename, err)
			continue
		}
		lines := strings.Split(string(data), "\n")
		
		// Удаляем последнюю пустую строку если есть
		if len(lines) > 0 && lines[len(lines)-1] == "" {
			lines = lines[:len(lines)-1]
		}
		
		// Выводим последние n строк
		start := len(lines) - n
		if start < 0 {
			start = 0
		}
		for j := start; j < len(lines); j++ {
			fmt.Println(lines[j])
		}
	}
}

// С флагом -q (не печатать заголовки)
func tailQflag(files []string, n int) {
	for _, filename := range files {
		data, err := os.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "tail: %s: %v\n", filename, err)
			continue
		}
		lines := strings.Split(string(data), "\n")
		
		// Удаляем последнюю пустую строку если есть
		if len(lines) > 0 && lines[len(lines)-1] == "" {
			lines = lines[:len(lines)-1]
		}
		
		// Выводим последние n строк
		start := len(lines) - n
		if start < 0 {
			start = 0
		}
		for j := start; j < len(lines); j++ {
			fmt.Println(lines[j])
		}
	}
}

// Чтение из stdin
func tailStdin(n int) {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// Данные есть в stdin
		data, err := os.ReadFile(os.Stdin.Name())
		if err != nil {
			fmt.Fprintf(os.Stderr, "tail: %v\n", err)
			return
		}
		lines := strings.Split(string(data), "\n")
		
		// Удаляем последнюю пустую строку если есть
		if len(lines) > 0 && lines[len(lines)-1] == "" {
			lines = lines[:len(lines)-1]
		}
		
		// Выводим последние n строк
		start := len(lines) - n
		if start < 0 {
			start = 0
		}
		for i := start; i < len(lines); i++ {
			fmt.Println(lines[i])
		}
	}
}

func printHelp() {
	fmt.Println(`Использование: tail [флаг] [файл...]
	
Флаги:
  -n число   количество строк для вывода (по умолчанию 10)
  -q         не печатать заголовки с именами файлов
  -v         всегда печатать заголовки с именами файлов
  -h         справка по команде и параметрам

Примеры:
  tail file.txt              # последние 10 строк файла
  tail -n 20 file.txt        # последние 20 строк файла
  tail -q file1.txt file2.txt # без заголовков
  tail -v file.txt           # с заголовком для одного файла
  cat file.txt | tail        # чтение из stdin`)
}