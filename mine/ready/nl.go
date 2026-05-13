package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	// Исправляем флаги под правильные параметры nl
	dflag := flag.String("d", "", "использовать CC как логический разделитель страниц")
	sflag := flag.String("s", "", "добавлять СТРОКУ после номера")
	hflag := flag.Bool("h", false, "справка о команде nl")
	
	flag.Parse()

	if *hflag {
		printHelp()
		return
	}

	files := flag.Args() // получаем файлы после флагов
	
	// Если есть флаг -d
	if *dflag != "" {
		nlwithd(files, *dflag)
		return
	}
	
	// Если есть флаг -s
	if *sflag != "" {
		nlwiths(files, *sflag)
		return
	}

	// Без флагов
	litenl(files)
}

func nlwithd(files []string, delimiter string) {
	if len(files) == 0 {
		files = []string{"-"}
	}
	
	for _, filename := range files {
		var scanner *bufio.Scanner
		
		if filename == "-" {
			scanner = bufio.NewScanner(os.Stdin)
		} else {
			file, err := os.Open(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "nl: %s: %v\n", filename, err)
				continue
			}
			defer file.Close()
			scanner = bufio.NewScanner(file)
		}
		
		lineNum := 1
		for scanner.Scan() {
			line := scanner.Text()
			
			// Проверяем, является ли строка разделителем
			if strings.TrimSpace(line) == delimiter {
				fmt.Println(line) // Выводим строку-разделитель
				lineNum = 1       // Сбрасываем нумерацию для новой секции
				continue
			}
			
			fmt.Printf("%6d\t%s\n", lineNum, line)
			lineNum++
		}
		
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "nl: %s: %v\n", filename, err)
		}
	}
}

func nlwiths(files []string, separator string) {
	if len(files) == 0 {
		files = []string{"-"}
	}
	
	for _, filename := range files {
		var scanner *bufio.Scanner
		
		if filename == "-" {
			scanner = bufio.NewScanner(os.Stdin)
		} else {
			file, err := os.Open(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "nl: %s: %v\n", filename, err)
				continue
			}
			defer file.Close()
			scanner = bufio.NewScanner(file)
		}
		
		lineNum := 1
		for scanner.Scan() {
			line := scanner.Text()
			// Используем свой разделитель вместо табуляции
			fmt.Printf("%6d%s%s\n", lineNum, separator, line)
			lineNum++
		}
		
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "nl: %s: %v\n", filename, err)
		}
	}
}

func litenl(files []string) {
	if len(files) == 0 {
		files = []string{"-"}
	}
	
	for _, filename := range files {
		var scanner *bufio.Scanner
		
		if filename == "-" {
			scanner = bufio.NewScanner(os.Stdin)
		} else {
			file, err := os.Open(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "nl: %s: %v\n", filename, err)
				continue
			}
			defer file.Close()
			scanner = bufio.NewScanner(file)
		}
		
		lineNum := 1
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Printf("%6d\t%s\n", lineNum, line)
			lineNum++
		}
		
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "nl: %s: %v\n", filename, err)
		}
	}
}

func printHelp() {
	fmt.Println(`Использование: nl [флаг] [файл...]
	
Флаги:
  -d СС     использовать CC как логический разделитель страниц
  -s СТРОКА добавлять СТРОКУ после номера
  -h        справка по команде и параметрам

Примеры:
  nl file.txt                    # обычная нумерация строк
  nl -s "-> " file.txt           # свой разделитель после номера
  nl -d "@@" file.txt            # разделитель страниц @@

СС — это два символа, используемые для создания логического разделителя
страниц; при отсутствии второго используется «:». Расширение GNU позволяет
задать больше двух символов, а также задание пустой строки (-d '')
отключает сравнение частей.`)
}