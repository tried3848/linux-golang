package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// FileInfo хранит информацию о файле
type FileInfo struct {
	Name    string
	Size    int64
	Mode    os.FileMode
	ModTime time.Time
	IsDir   bool
}

func main() {
	// Определение флагов командной строки
	helpFlag := flag.Bool("h", false, "Показать справку")
	longFlag := flag.Bool("l", false, "Подробный вывод (права, размер, дата)")
	allFlag := flag.Bool("a", false, "Показать скрытые файлы (начинающиеся с точки)")
	flag.Parse()

	// Обработка флага помощи
	if *helpFlag {
		fmt.Println("Использование: ls [-h] [-l] [-a] [директория]")
		fmt.Println("  -h  Показать справку")
		fmt.Println("  -l  Подробный вывод")
		fmt.Println("  -a  Показать скрытые файлы")
		os.Exit(0)
	}

	// Определяем рабочую директорию
	dir := "."
	if args := flag.Args(); len(args) > 0 {
		dir = args[0]
	}

	// Проверяем корректность пути
	info, err := os.Stat(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ls: невозможно получить доступ к '%s': %v\n", dir, err)
		os.Exit(1)
	}
	if !info.IsDir() {
		fmt.Fprintf(os.Stderr, "ls: '%s' не является директорией\n", dir)
		os.Exit(1)
	}

	// Читаем содержимое директории
	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ls: ошибка чтения директории: %v\n", err)
		os.Exit(1)
	}

	// Собираем список файлов в срез структур FileInfo
	var files []FileInfo
	for _, entry := range entries {
		name := entry.Name()
		// Пропускаем скрытые файлы если флаг -a не установлен
		if !*allFlag && len(name) > 0 && name[0] == '.' {
			continue
		}
		fi, err := entry.Info()
		if err != nil {
			continue
		}
		files = append(files, FileInfo{
			Name:    name,
			Size:    fi.Size(),
			Mode:    fi.Mode(),
			ModTime: fi.ModTime(),
			IsDir:   fi.IsDir(),
		})
	}

	// Сортируем по имени
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name < files[j].Name
	})

	// Выводим результат
	for _, f := range files {
		if *longFlag {
			// Подробный формат: права размер дата имя
			fmt.Printf("%s %8d %s %s\n",
				f.Mode,
				f.Size,
				f.ModTime.Format("Jan 02 15:04"),
				filepath.Base(f.Name),
			)
		} else {
			fmt.Printf("%s  ", f.Name)
		}
	}
	if !*longFlag {
		fmt.Println()
	}
}
