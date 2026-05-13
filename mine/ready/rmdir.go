package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	vflag := flag.Bool("v", false, "информация о удалении")
	verflag := flag.Bool("ver", false, "версия команды")
	hflag := flag.Bool("h", false, "справка по команде и параметрам")
	
	flag.Parse()

	if *hflag {
		printHelp()
		return
	}

	if *verflag {
		rmdirVerflag()
		return
	}

	// Получаем имя директории из аргументов после флагов
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Использование: rmdir [флаги] <имя_директории>")
		fmt.Println("Для справки используйте: rmdir -h")
		return
	}
	
	dirname := args[0]
	
	if *vflag {
		rmdirVflag(dirname)
	} else {
		rmdirNoFlag(dirname)
	}
}

func rmdirNoFlag(dirname string) {
	// Сначала проверяем существование
	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		fmt.Printf("rmdir: директория '%s' не существует\n", dirname)
		return
	}
	
	// Затем удаляем
	err := os.RemoveAll(dirname)
	if err != nil {
		fmt.Println("Ошибка во время удаления:", err)
	}
	// В случае успеха ничего не выводим
}

func rmdirVflag(dirname string) {
	// Сначала проверяем существование
	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		fmt.Printf("rmdir: директория '%s' не существует\n", dirname)
		return
	}
	
	err := os.RemoveAll(dirname)
	if err != nil {
		fmt.Println("Ошибка во время удаления:", err)
	} else {
		fmt.Printf("rmdir: удаление директории %s\n", dirname)
	}
}

func rmdirVerflag() {
	fmt.Println(`rmdir (Go-аналог) 0.1
Copyright - нет
Лицензия - нет
Это свободное ПО: используйте и распространяйте`)
}

func printHelp() {
	fmt.Println(`Использование: rmdir [флаги] <имя_директории>
	
Удаляет пустые директории (в данной реализации удаляет любые директории рекурсивно)

Флаги:
  -v    	информация о удалении
  -ver    	версия команды
  -h    	справка по команде и параметрам

Примеры:
  rmdir mydir           # удаляет директорию mydir без вывода
  rmdir -v mydir        # удаляет директорию mydir с выводом информации
  rmdir -ver            # показывает версию программы`)
}