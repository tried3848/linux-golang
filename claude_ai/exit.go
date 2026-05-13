// exit.go — аналог команды exit оболочки
// Завершает процесс с указанным кодом возврата.
//
// Аргументы:
//   -h        — справка
//   -c <код>  — код завершения (0–255, по умолчанию 0)
//   -v        — вывести сообщение перед завершением
//   -s <msg>  — пользовательское сообщение при завершении

package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// Определение флагов
	helpFlag := flag.Bool("h", false, "Показать справку")
	exitCode := flag.Int("c", 0, "Код завершения (0–255)")
	verbose := flag.Bool("v", false, "Вывести сообщение перед завершением")
	message := flag.String("s", "", "Пользовательское сообщение при завершении")

	flag.Parse()

	// Вывод справки
	if *helpFlag {
		fmt.Println("Использование: exit [ОПЦИИ]")
		fmt.Println("Завершает программу с заданным кодом возврата.")
		fmt.Println()
		fmt.Println("Опции:")
		flag.PrintDefaults()
		fmt.Println()
		fmt.Println("Примеры:")
		fmt.Println("  exit -c 0          # успешное завершение")
		fmt.Println("  exit -c 1 -v       # завершение с кодом 1 и сообщением")
		fmt.Println("  exit -c 42 -s 'Пока!'  # завершение с кодом 42 и текстом 'Пока!'")
		os.Exit(0)
	}

	// Проверка диапазона кода завершения
	if *exitCode < 0 || *exitCode > 255 {
		fmt.Fprintf(os.Stderr, "Ошибка: код завершения должен быть в диапазоне 0–255, получено: %d\n", *exitCode)
		os.Exit(1)
	}

	// Вывод сообщения при необходимости
	if *message != "" {
		fmt.Println(*message)
	} else if *verbose {
		if *exitCode == 0 {
			fmt.Printf("Завершение с кодом %d (успех)\n", *exitCode)
		} else {
			fmt.Printf("Завершение с кодом %d (ошибка)\n", *exitCode)
		}
	}

	// Завершение с указанным кодом
	os.Exit(*exitCode)
}
