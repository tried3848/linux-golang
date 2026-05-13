package main

import (
	"fmt"
	"io/fs"
	"os"
	"strings"
	"time"
)

func main() {
	// Используем ручной разбор флагов вместо flag пакета
	args := os.Args[1:]
	
	showAll := false
	longFormat := false
	
	// Разбираем флаги
	for _, arg := range args {
		if arg == "-h" {
			printHelp()
			return
		}
		
		// Поддерживаем -a, -l, -la, -al
		if strings.HasPrefix(arg, "-") {
			flags := arg[1:] // убираем дефис
			for _, f := range flags {
				switch f {
				case 'a':
					showAll = true
				case 'l':
					longFormat = true
				default:
					fmt.Printf("Неизвестный флаг: -%c\n", f)
					fmt.Println("Попробуйте 'ls -h' для справки")
					return
				}
			}
		}
	}
	
	// Вызываем соответствующую функцию
	if longFormat && showAll {
		lsWithLongAndAll()
	} else if longFormat {
		lsWithLong()
	} else if showAll {
		lsWithAll()
	} else {
		lsDefault()
	}
}

func readDirectory() []fs.DirEntry {
	entries, err := os.ReadDir(".")
	if err != nil {
		fmt.Println("Ошибка чтения директории:", err)
		os.Exit(1)
	}
	return entries
}

func lsDefault() {
	entries := readDirectory()
	var names []string
	
	for _, entry := range entries {
		if !strings.HasPrefix(entry.Name(), ".") {
			names = append(names, entry.Name())
		}
	}
	
	fmt.Println(strings.Join(names, "  "))
}

func lsWithAll() {
	entries := readDirectory()
	var names []string
	
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	
	fmt.Println(strings.Join(names, "  "))
}

func lsWithLong() {
	entries := readDirectory()
	
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		printLongFormat(entry)
	}
}

func lsWithLongAndAll() {
	entries := readDirectory()
	
	for _, entry := range entries {
		printLongFormat(entry)
	}
}

func printLongFormat(entry fs.DirEntry) {
	info, err := entry.Info()
	if err != nil {
		fmt.Printf("%-20s [ошибка]\n", entry.Name())
		return
	}
	
	// Права доступа
	perms := info.Mode().Perm().String()
	
	// Тип файла
	if info.IsDir() {
		perms = "d" + perms[1:]
	}
	
	// Размер
	size := info.Size()
	
	// Время модификации
	modTime := info.ModTime().Format("Jan 02 15:04")
	
	// Имя
	name := entry.Name()
	
	fmt.Printf("%-10s %8d %12s %s\n", perms, size, modTime, name)
}

func printHelp() {
	fmt.Println(`Использование: ls [флаги]
	
Флаги:
  -a    показывать все файлы, включая скрытые (начинающиеся с .)
  -l    подробный формат вывода (права, размер, дата)
  -la   комбинация -l и -a (подробно о всех файлах)
  -al   то же самое что и -la
  -h    показать эту справку

Примеры:
  ls         # обычный вывод
  ls -a      # все файлы (включая скрытые)
  ls -l      # подробный вывод
  ls -la     # подробный вывод всех файлов
  ls -h      # справка`)
}