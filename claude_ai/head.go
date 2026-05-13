// head.go — аналог утилиты GNU head (вывод первых строк файла)
//
// Аргументы:
//   -h        — справка
//   -n <число> — количество выводимых строк (по умолчанию 10)
//   -c <число> — количество выводимых байт (вместо строк)
//   -q        — не выводить заголовок с именем файла

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	// Определение флагов
	helpFlag := flag.Bool("h", false, "Показать справку")
	numLines := flag.Int("n", 10, "Количество строк для вывода")
	numBytes := flag.Int("c", 0, "Количество байт для вывода (0 = не используется)")
	quiet := flag.Bool("q", false, "Не выводить заголовки с именами файлов")

	flag.Parse()

	// Вывод справки
	if *helpFlag {
		fmt.Println("Использование: head [ОПЦИИ] [ФАЙЛ...]")
		fmt.Println("Выводит первые строки (или байты) каждого файла.")
		fmt.Println()
		fmt.Println("Опции:")
		flag.PrintDefaults()
		os.Exit(0)
	}

	// Проверка корректности аргументов
	if *numLines <= 0 {
		fmt.Fprintln(os.Stderr, "Ошибка: число строк (-n) должно быть положительным")
		os.Exit(1)
	}
	if *numBytes < 0 {
		fmt.Fprintln(os.Stderr, "Ошибка: число байт (-c) не может быть отрицательным")
		os.Exit(1)
	}
	// -n и -c несовместимы
	if *numBytes > 0 && *numLines != 10 {
		fmt.Fprintln(os.Stderr, "Ошибка: нельзя одновременно использовать -n и -c")
		os.Exit(1)
	}

	exitCode := 0
	files := flag.Args()

	// Если файлы не указаны — читаем из stdin
	if len(files) == 0 {
		err := printHead(os.Stdin, "", *numLines, *numBytes, false)
		if err != nil {
			fmt.Fprintf(os.Stderr, "head: ошибка чтения stdin: %v\n", err)
			exitCode = 1
		}
	} else {
		// Заголовок показываем если файлов несколько или не задан -q
		showHeader := len(files) > 1 && !*quiet
		for _, filename := range files {
			f, err := os.Open(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "head: %s: %v\n", filename, err)
				exitCode = 1
				continue
			}
			if showHeader {
				fmt.Printf("==> %s <==\n", filename)
			}
			err = printHead(f, filename, *numLines, *numBytes, false)
			f.Close()
			if err != nil {
				fmt.Fprintf(os.Stderr, "head: %s: %v\n", filename, err)
				exitCode = 1
			}
			if showHeader {
				fmt.Println()
			}
		}
	}

	os.Exit(exitCode)
}

// printHead выводит первые n строк или b байт из reader
func printHead(r io.Reader, name string, n, b int, _ bool) error {
	if b > 0 {
		// Режим байтов: читаем ровно b байт
		buf := make([]byte, b)
		read, err := io.ReadAtLeast(r, buf, 1)
		if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
			return err
		}
		_, writeErr := os.Stdout.Write(buf[:read])
		return writeErr
	}

	// Режим строк: выводим первые n строк
	scanner := bufio.NewScanner(r)
	count := 0
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		count++
		if count >= n {
			break
		}
	}
	return scanner.Err()
}
