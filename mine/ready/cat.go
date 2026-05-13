package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func main() {
	bflag := flag.Bool("b", false, "нумеровать не пустые строки")
	nflag := flag.Bool("n", false, "нумеровать все строки")
	hflag := flag.Bool("h", false, "справка о команде cat")
	
	flag.Parse()

	if *hflag {
		printHelp()
		return
	}

	if *bflag && *nflag {
		fmt.Println("Нельзя использовать -b и -n вместе")
		return
	}

	// Обрабатываем каждый файл - этот цикл должен быть ВНУТРИ main
	for _, filename := range flag.Args() {
		if err := processFile(filename, *bflag, *nflag); err != nil {
			fmt.Fprintf(os.Stderr, "ошибка при чтении %s: %v\n", filename, err)
			os.Exit(1)
		}
	}
}

func processFile(filename string, bFlag, nFlag bool) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 1

	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case nFlag:
			// Нумеруем все строки
			fmt.Printf("%6d  %s\n", lineNum, line)
			lineNum++

		case bFlag:
			// Нумеруем только непустые строки
			if line != "" {
				fmt.Printf("%6d  %s\n", lineNum, line)
				lineNum++
			} else {
				fmt.Println(line)
			}

		default:
			// Без нумерации
			fmt.Println(line)
		}
	}

	return scanner.Err()
}

func printHelp() {
	fmt.Println(`Использование: cat [флаг] [файлы...]
	
Флаги:
  -b    нумеровать не пустые строки
  -n    нумеровать все строки
  -h    справка по команде и параметрам`)
}