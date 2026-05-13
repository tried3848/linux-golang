// history_bang.go — аналог команд !! и !n (повтор команды из истории)
//
// В оболочке bash:
//   !!   — повторить последнюю команду
//   !n   — повторить команду с номером n
//
// Данная утилита читает файл истории (~/.bash_history или указанный),
// выводит нужную команду и сообщает, какую именно нужно выполнить.
//
// Аргументы:
//   -h           — справка
//   -f <файл>    — путь к файлу истории (по умолчанию ~/.bash_history)
//   -n <номер>   — номер команды (0 = последняя, аналог !!)
//   -l           — показать последние 20 строк истории (список)

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// Определение флагов
	helpFlag := flag.Bool("h", false, "Показать справку")
	histFile := flag.String("f", "", "Путь к файлу истории (по умолчанию ~/.bash_history)")
	lineNum := flag.Int("n", 0, "Номер команды в истории (0 = последняя, т.е. !!)")
	listMode := flag.Bool("l", false, "Показать последние 20 команд истории")

	flag.Parse()

	// Вывод справки
	if *helpFlag {
		fmt.Println("Использование: history_bang [ОПЦИИ]")
		fmt.Println("Аналог !! и !n — просмотр и повтор команд из истории bash.")
		fmt.Println()
		fmt.Println("Опции:")
		flag.PrintDefaults()
		fmt.Println()
		fmt.Println("Примеры:")
		fmt.Println("  history_bang -l            # показать последние 20 команд")
		fmt.Println("  history_bang -n 0          # вывести последнюю команду (!!)")
		fmt.Println("  history_bang -n 42         # вывести команду номер 42 (!42)")
		os.Exit(0)
	}

	// Определяем путь к файлу истории
	if *histFile == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Ошибка: не удалось определить домашний каталог")
			os.Exit(1)
		}
		*histFile = filepath.Join(home, ".bash_history")
	}

	// Открываем файл истории
	f, err := os.Open(*histFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка открытия файла истории %q: %v\n", *histFile, err)
		os.Exit(1)
	}
	defer f.Close()

	// Считываем все строки истории в срез (slice)
	var history []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			history = append(history, line)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка чтения файла истории: %v\n", err)
		os.Exit(1)
	}

	total := len(history)
	if total == 0 {
		fmt.Fprintln(os.Stderr, "История команд пуста.")
		os.Exit(1)
	}

	// Режим вывода списка последних 20 команд
	if *listMode {
		start := total - 20
		if start < 0 {
			start = 0
		}
		for i := start; i < total; i++ {
			fmt.Printf("%4d  %s\n", i+1, history[i])
		}
		os.Exit(0)
	}

	// Проверка корректности номера строки
	if *lineNum < 0 {
		fmt.Fprintln(os.Stderr, "Ошибка: номер команды (-n) не может быть отрицательным")
		os.Exit(1)
	}
	if *lineNum > total {
		fmt.Fprintf(os.Stderr, "Ошибка: номер команды %d выходит за пределы истории (всего %d команд)\n", *lineNum, total)
		os.Exit(1)
	}

	// Выбор команды: 0 = последняя (!!)
	var cmd string
	if *lineNum == 0 {
		cmd = history[total-1]
		fmt.Printf("!! -> %s\n", cmd)
	} else {
		cmd = history[*lineNum-1]
		fmt.Printf("!%d -> %s\n", *lineNum, cmd)
	}

	// Подсказка пользователю (реальный запуск запрещён требованиями — no exec)
	fmt.Printf("\nДля выполнения скопируйте команду выше в терминал.\n")
	_ = cmd
}
