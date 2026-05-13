package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	aflag := flag.Bool("a", false, "изменить только время доступа")
	dflag := flag.String("d", "", "проанализировать СТРОКУ и использовать её вместо текущего времени")
	hflag := flag.Bool("h", false, "справка о команде touch")

	flag.Parse()

	if *hflag {
		printHelp()
		return
	}

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Использование: ./touch <имя_файла>")
		return
	}
	filename := args[0]

	// Если файла нет - создаём
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			fmt.Println("Ошибка во время создания:", err)
			return
		}
		file.Close()
		fmt.Println("Файл создан")
	}

	if *aflag {
		touchAflag(filename)
	} else if *dflag != "" {
		touchDflag(filename, *dflag)
	} else {
		touchnoflag(filename)
	}
}

func touchnoflag(filename string) {
	now := time.Now()
	err := os.Chtimes(filename, now, now)
	if err != nil {
		fmt.Println("Ошибка при обновлении времени:", err)
	}
}

func touchAflag(filename string) {
	now := time.Now()
	stat, _ := os.Stat(filename)
	mtime := stat.ModTime()
	err := os.Chtimes(filename, now, mtime)
	if err != nil {
		fmt.Println("Ошибка при обновлении времени доступа:", err)
	}
}

func touchDflag(filename string, datestr string) {
	// Парсим дату формата: 20260505051111 -> 2026-05-05 05:11:11
	if len(datestr) == 14 {
		year := datestr[0:4]
		month := datestr[4:6]
		day := datestr[6:8]
		hour := datestr[8:10]
		min := datestr[10:12]
		sec := datestr[12:14]
		
		parsedTime, err := time.Parse("2006-01-02 15:04:05", 
			fmt.Sprintf("%s-%s-%s %s:%s:%s", year, month, day, hour, min, sec))
		if err == nil {
			err = os.Chtimes(filename, parsedTime, parsedTime)
			if err != nil {
				fmt.Println("Ошибка при установке времени:", err)
			}
			return
		}
	}
	
	// Если формат не подошёл
	fmt.Println("Ошибка: используйте формат ГГГГММДДЧЧММСС (например: 20260505051111)")
}

func printHelp() {
	fmt.Println(`Использование: touch [флаг] [файл]
	
Флаги:
  -a    изменить только время доступа
  -d    проанализировать СТРОКУ и использовать её вместо текущего времени
  -h    справка по команде и параметрах
	
Пример:
  ./touch -d 20260505051111 12.txt`)
}