// clear.go — аналог утилиты GNU clear (очистка экрана терминала)
// Отправляет ANSI escape-последовательности для очистки экрана.
//
// Аргументы:
//   -h  — справка
//   -s  — только прокрутить экран вниз (без ANSI-кодов, совместимо с любым терминалом)
//   -x  — не сбрасывать позицию прокрутки (только очистить видимую область)

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// Константы ANSI escape-кодов
const (
	// Полная очистка экрана и перемещение курсора в начало
	ansiClearFull = "\033[2J\033[H"
	// Очистка только видимой области без сброса прокрутки
	ansiClearView = "\033[H\033[J"
)

func main() {
	// Определение флагов
	helpFlag := flag.Bool("h", false, "Показать справку")
	scroll := flag.Bool("s", false, "Очистить прокруткой (без ANSI-кодов, совместимый режим)")
	noReset := flag.Bool("x", false, "Не сбрасывать прокрутку — очистить только видимую область")

	flag.Parse()

	// Вывод справки
	if *helpFlag {
		fmt.Println("Использование: clear [ОПЦИИ]")
		fmt.Println("Очищает экран терминала.")
		fmt.Println()
		fmt.Println("Опции:")
		flag.PrintDefaults()
		os.Exit(0)
	}

	// Несовместимые флаги
	if *scroll && *noReset {
		fmt.Fprintln(os.Stderr, "Ошибка: флаги -s и -x несовместимы")
		os.Exit(1)
	}

	if *scroll {
		// Совместимый режим: просто вывести несколько пустых строк
		fmt.Print(strings.Repeat("\n", 40))
	} else if *noReset {
		// Очистить только видимую область терминала
		fmt.Print(ansiClearView)
	} else {
		// Полная очистка с перемещением курсора в верхний левый угол
		fmt.Print(ansiClearFull)
	}
}
