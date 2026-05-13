package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// Определение флагов командной строки
	helpFlag := flag.Bool("h", false, "Показать справку")
	printFlag := flag.Bool("p", false, "Вывести итоговый путь после перехода")
	physicalFlag := flag.Bool("P", false, "Разрешить симлинки (физический путь)")
	flag.Parse()

	// Обработка флага помощи
	if *helpFlag {
		fmt.Println("Использование: cd [-h] [-p] [-P] <директория>")
		fmt.Println("  -h  Показать справку")
		fmt.Println("  -p  Напечатать путь после перехода")
		fmt.Println("  -P  Разрешить симлинки")
		fmt.Println()
		fmt.Println("Примечание: данная утилита меняет директорию внутри процесса.")
		fmt.Println("Используйте встроенную команду shell 'cd' для смены директории в терминале.")
		os.Exit(0)
	}

	// Проверяем, что передан аргумент директории
	args := flag.Args()
	if len(args) == 0 {
		// По умолчанию — домашняя директория
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "cd: не удалось определить домашнюю директорию: %v\n", err)
			os.Exit(1)
		}
		args = []string{home}
	}

	target := args[0]

	// Проверяем, что директория существует
	info, err := os.Stat(target)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cd: '%s': нет такого файла или директории\n", target)
		os.Exit(1)
	}
	if !info.IsDir() {
		fmt.Fprintf(os.Stderr, "cd: '%s': не является директорией\n", target)
		os.Exit(1)
	}

	// Меняем директорию процесса
	if err := os.Chdir(target); err != nil {
		fmt.Fprintf(os.Stderr, "cd: не удалось перейти в '%s': %v\n", target, err)
		os.Exit(1)
	}

	// Получаем итоговый путь
	newDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cd: ошибка получения текущей директории: %v\n", err)
		os.Exit(1)
	}

	// Если нужен физический путь
	if *physicalFlag {
		newDir, err = filepath.EvalSymlinks(newDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cd: ошибка разрешения пути: %v\n", err)
			os.Exit(1)
		}
	}

	// Выводим путь если запрошен флаг -p
	if *printFlag {
		fmt.Println(newDir)
	}
}
