// rm.go — аналог утилиты GNU rm (удаление файлов и каталогов)
//
// Аргументы:
//   -h  — справка
//   -r  — рекурсивное удаление каталогов
//   -f  — принудительное удаление без запроса (игнорировать несуществующие)
//   -i  — интерактивный режим: запрашивать подтверждение перед каждым удалением

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	// Определение флагов
	helpFlag := flag.Bool("h", false, "Показать справку")
	recursive := flag.Bool("r", false, "Рекурсивно удалять каталоги и их содержимое")
	force := flag.Bool("f", false, "Не сообщать об ошибке если файл не существует, не запрашивать подтверждение")
	interactive := flag.Bool("i", false, "Запрашивать подтверждение перед каждым удалением")

	flag.Parse()

	// Вывод справки
	if *helpFlag {
		fmt.Println("Использование: rm [ОПЦИИ] ФАЙЛ [ФАЙЛ...]")
		fmt.Println("Удаляет файлы или каталоги.")
		fmt.Println()
		fmt.Println("Опции:")
		flag.PrintDefaults()
		fmt.Println()
		fmt.Println("ВНИМАНИЕ: удалённые файлы восстановить невозможно!")
		os.Exit(0)
	}

	// Проверка наличия аргументов
	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "Ошибка: необходимо указать хотя бы один файл или каталог")
		fmt.Fprintln(os.Stderr, "Используйте -h для справки")
		os.Exit(1)
	}

	// Нельзя одновременно задать -f и -i (они противоречат друг другу)
	if *force && *interactive {
		fmt.Fprintln(os.Stderr, "Ошибка: флаги -f и -i несовместимы")
		os.Exit(1)
	}

	stdin := bufio.NewReader(os.Stdin)
	exitCode := 0

	// Обрабатываем каждый переданный путь
	for _, path := range flag.Args() {
		err := removeEntry(path, *recursive, *force, *interactive, stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "rm: %s: %v\n", path, err)
			exitCode = 1
		}
	}

	os.Exit(exitCode)
}

// removeEntry удаляет файл или каталог по указанному пути
func removeEntry(path string, recursive, force, interactive bool, stdin *bufio.Reader) error {
	info, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			if force {
				return nil // при -f не сообщаем об отсутствии файла
			}
			return fmt.Errorf("нет такого файла или каталога")
		}
		return err
	}

	// Если это каталог — требуется флаг -r
	if info.IsDir() {
		if !recursive {
			return fmt.Errorf("это каталог (используйте -r для рекурсивного удаления)")
		}
		// Интерактивный запрос на удаление каталога
		if interactive {
			if !confirm(fmt.Sprintf("Удалить каталог '%s' рекурсивно?", path), stdin) {
				return nil
			}
		}
		return os.RemoveAll(path)
	}

	// Файл или символьная ссылка
	if interactive {
		if !confirm(fmt.Sprintf("Удалить файл '%s'?", path), stdin) {
			return nil
		}
	}

	return os.Remove(path)
}

// confirm запрашивает у пользователя подтверждение (y/n)
func confirm(prompt string, reader *bufio.Reader) bool {
	fmt.Printf("%s [y/N] ", prompt)
	resp, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	resp = strings.TrimSpace(strings.ToLower(resp))
	return resp == "y" || resp == "yes" || resp == "д" || resp == "да"
}
