package main

import (
	"fmt"
	"time"
	"flag"
	"os"
)

func main() {
	dflag := flag.Bool("d", false, "время из строки")
	rflag := flag.Bool("r", false, "время в соответсвии")
	hflag := flag.Bool("h", false, "справка о команде")
	flag.Parse()
	mskLocation, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		fmt.Println(err)
		return
	}

	if *hflag {
		printHelp()
		return
	}

	if *rflag {
		now := time.Now().In(mskLocation)
		fmt.Println(now.Format(time.RFC1123Z))
		return
	}

	if *dflag{
		if flag.NArg() < 1 {
			fmt.Fprintln(os.Stderr, "Ошибка: нужна строка с датой")
			os.Exit(1)
		}
		timeStr := flag.Arg(0)
		mineTime(timeStr, mskLocation)
		return
	}
	printCurrentTime(mskLocation)
}

func printCurrentTime(location *time.Location) {
	now := time.Now().In(location)
	fmt.Printf("%d %s %d %02d:%02d:%02d %s\n", 
		now.Day(), now.Month(), now.Year(), 
		now.Hour(), now.Minute(), now.Second(), 
		location)
}

func mineTime(timeStr string, location *time.Location) {
	var year, month, day, hour, minute, second int
	
	// Простой разбор формата "2006-01-02 15:04:05"
	_, err := fmt.Sscanf(timeStr, "%d-%d-%d %d:%d:%d", &year, &month, &day, &hour, &minute, &second)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка: не удалось распарсить дату '%s'\n", timeStr)
		fmt.Fprintln(os.Stderr, "Используйте формат: ГГГГ-ММ-ДД ЧЧ:ММ:СС")
		os.Exit(1)
	}
	
	// // Создаем время
	// parsedTime := time.Date(year, time.Month(month), day, hour, minute, second, 0, location)
	
	fulltime := fmt.Sprintf("%02d:%02d:%02d", hour, minute, second)
	
	fmt.Printf("%d %s %d %s %s\n", day, time.Month(month), year, fulltime, location)
}


func printHelp() {
	fmt.Println(`Использование: date [флаг] []
	
Флаги:
  -d    показывать не текущее время, а время, описанное заданной СТРОКОЙ
  -R 	выводит дату в соответсвии с RFC-2822
  -h    справка по команде и параметрах`)
}
