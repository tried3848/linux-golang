// touch.go — аналог утилиты GNU touch
// Создаёт файл если не существует, или обновляет время доступа/изменения.
//
// Аргументы:
//   -h         — справка
//   -c         — не создавать файл если он не существует
//   -t <штамп> — установить время в формате YYYYMMDDHHMMSS
//   -a         — обновить только время доступа (atime)

package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	// Определение флагов командной строки
	helpFlag := flag.Bool("h", false, "Показать справку")
	noCreate := flag.Bool("c", false, "Не создавать файл, если он не существует")
	timeStamp := flag.String("t", "", "Установить время в формате YYYYMMDDHHMMSS")
	atimeOnly := flag.Bool("a", false, "Обновить только время доступа (atime)")

	flag.Parse()

	// Вывод справки
	if *helpFlag {
		fmt.Println("Использование: touch [ОПЦИИ] ФАЙЛ [ФАЙЛ...]")
		fmt.Println("Создаёт файлы или обновляет метки времени.")
		fmt.Println()
		fmt.Println("Опции:")
		flag.PrintDefaults()
		fmt.Println()
		fmt.Println("Пример формата времени: 20240315120000 (15 марта 2024, 12:00:00)")
		os.Exit(0)
	}

	// Проверка наличия хотя бы одного файла
	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "Ошибка: необходимо указать хотя бы один файл")
		fmt.Fprintln(os.Stderr, "Используйте -h для справки")
		os.Exit(1)
	}

	// Разбор метки времени, если задана
	var targetTime time.Time
	if *timeStamp != "" {
		// Допустимые форматы
		formats := []string{
			"200601021504",     // YYYYMMDDHHMM
			"20060102150405",   // YYYYMMDDHHMMSS
		}
		var parseErr error
		parsed := false
		for _, fmt_ := range formats {
			targetTime, parseErr = time.ParseInLocation(fmt_, *timeStamp, time.Local)
			if parseErr == nil {
				parsed = true
				break
			}
		}
		if !parsed {
			fmt.Fprintf(os.Stderr, "Ошибка: неверный формат времени %q. Используйте YYYYMMDDHHMM или YYYYMMDDHHMMSS\n", *timeStamp)
			os.Exit(1)
		}
		// Проверка разумности даты
		if targetTime.Year() < 1970 || targetTime.Year() > 2100 {
			fmt.Fprintln(os.Stderr, "Ошибка: год должен быть в диапазоне 1970–2100")
			os.Exit(1)
		}
	} else {
		targetTime = time.Now()
	}

	// Обработка каждого файла
	exitCode := 0
	for _, filename := range flag.Args() {
		err := touchFile(filename, targetTime, *noCreate, *atimeOnly)
		if err != nil {
			fmt.Fprintf(os.Stderr, "touch: %s: %v\n", filename, err)
			exitCode = 1
		}
	}
	os.Exit(exitCode)
}

// touchFile — создаёт файл или обновляет его метки времени
func touchFile(path string, t time.Time, noCreate bool, atimeOnly bool) error {
	// Проверяем существование файла
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		if noCreate {
			// Флаг -c: файл не создаём
			return nil
		}
		// Создаём пустой файл
		f, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("не удалось создать файл: %w", err)
		}
		f.Close()
	} else if err != nil {
		return fmt.Errorf("ошибка доступа: %w", err)
	}

	// Устанавливаем метки времени
	// Если -a: обновляем только atime, mtime оставляем текущим файла
	if atimeOnly {
		info, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("stat: %w", err)
		}
		// mtime сохраняем прежним
		return os.Chtimes(path, t, info.ModTime())
	}

	// По умолчанию: обновляем оба времени
	return os.Chtimes(path, t, t)
}
