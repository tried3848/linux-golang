package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	args := os.Args[1:]
	
	if len(args) >= 1 && args[0] == "-h" {
		printHelp()
		return
	}
	
	exitcom(args)
}

func exitcom(args []string) {
	var exitCode int = 0
	
	// Если есть аргументы, пытаемся получить код выхода
	if len(args) > 0 {
		code, err := strconv.Atoi(args[0])
		if err != nil {
			// Если аргумент не число, выводим ошибку и выходим с кодом 1
			fmt.Fprintf(os.Stderr, "exit: %s: числовой аргумент требуется\n", args[0])
			exitCode = 1
		} else {
			exitCode = code
		}
	} else {
		// Если аргумент не указан, кодом выхода будет статус последней команды
		// В этом примере используем 0, но в реальной оболочке нужно хранить последний код
		exitCode = 0
	}
	
	os.Exit(exitCode)
}

func printHelp() {
	fmt.Println(`Использование exit [n]
    Выход из командного процессора.
    
    Закрывает командный процессор с состоянием n. Если n не указан,
    состоянием выхода будет состояние последней выполненной команды.
    
    Аргументы:
    n    Код возврата (целое число от 0 до 255)`)
}