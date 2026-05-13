package main

import (
	"fmt"
	"os"
	"flag"
)

func main() {
	
	vflag := flag.Bool("v", false, "сообщение о создании")
	hflag := flag.Bool("h", false, "справка о команде cat")
	pflag := flag.Bool("p", false, "создавать родительские директории при необходимости")
	
	flag.Parse()

	if *hflag {
		printHelp()
		return
	}
	if *vflag {
		mkdirmessage()
		return
	}

	dirs := flag.Args()
	if len(dirs) == 0{
		fmt.Println("Ошибка: Не указано имя директории")
		printHelp()
		return
	}
	
	// Обрабатываем каждую директорию
	for _, dir := range dirs {
		if *pflag {
			mkdirWithParents(dir, *vflag)
		} else {
			simplemkdir(dir, *vflag)
		}
	}
}

func simplemkdir(dirname string, verbose bool){
	err := os.Mkdir(dirname, 0755)
	if err != nil {
		fmt.Println("Ошибка во время создания:", err)
	} else if verbose {
		fmt.Println("Директория создана:", dirname)
	}
}

func mkdirWithParents(dirname string, verbose bool){
	err := os.MkdirAll(dirname, 0755)
	if err != nil {
		fmt.Println("Ошибка во время создания:", err)
	} else if verbose {
		fmt.Println("Директория(ии) создана(ы):", dirname)
	}
}

func mkdirmessage() {
	fmt.Println("Команда для создания директорий")
}

func printHelp() {
	fmt.Println(`Использование: mkdir [флаг] [директория...]
	
Флаги:
  -v    печатать сообщение о каждом созданном каталоге
  -p    создавать родительские директории при необходимости
  -h    справка по команде и параметрам

Примеры:
  mkdir test          - создает директорию test
  mkdir -v test       - создает директорию test с сообщением
  mkdir -p a/b/c      - создает все родительские директории a, a/b, a/b/c
  mkdir -p -v a/b/c   - создает директории с выводом сообщений`)
}