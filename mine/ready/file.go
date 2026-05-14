package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	verflag := flag.Bool("ver", false, "версия команды")
	hflag := flag.Bool("h", false, "справка по команде и параметрам")
	bflag := flag.Bool("b", false, "не добавлять имена файлов в начало строк вывода")

	flag.Parse()

	if *hflag {
		printHelp()
		return
	}

	if *verflag {
		fileVerflag()
		return
	}

	// Проверяем, передан ли аргумент
	if len(os.Args) < 2 {
		fmt.Println("Ошибка: не указано имя файла")
		fmt.Println("Используйте 'file -h' для справки")
		return
	}

	filename := os.Args[len(os.Args)-1]

	// Проверяем, что filename не является флагом
	if filename == "-b" || filename == "-ver" || filename == "-h" {
		fmt.Println("Ошибка: укажите имя файла")
		return
	}

	info, err := os.Stat(filename)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	mode := info.Mode()

	// Определяем тип файла
	var fileType string
	if mode.IsDir() {
		fileType = "directory"
	} else if mode.IsRegular() {
		fileType = "file"
	} else if mode&os.ModeSymlink != 0 {
		fileType = "symbolic link"
	} else {
		fileType = "other"
	}

	// Выводим результат с учетом флага -b
	if *bflag {
		// Краткий режим - только тип файла
		fmt.Printf("%s\n", fileType)
	} else {
		// Обычный режим - с именем файла
		fmt.Printf("%s: %s\n", filename, fileType)
	}
}

func fileVerflag() {
	fmt.Println(`file (Go-аналог) 0.2
Copyright - нет
Лицензия - нет
Это свободное ПО: используйте и распространяйте`)
}

func printHelp() {
	fmt.Println(`Использование: file [флаги] <имя_файла>
	
Определяет тип файла (директория или обычный файл)

Флаги:
  -b    	не добавлять имена файлов в начало строк вывода
  -ver    	версия команды
  -h    	справка по команде и параметрам

Примеры:
  file document.txt     # выводит: document.txt: file
  file -b document.txt  # выводит: file
  file /home/user       # выводит: /home/user: directory
  file -b /home/user    # выводит: directory
  file -ver             # показывает версию программы
  file -h               # показывает справку`)
}
