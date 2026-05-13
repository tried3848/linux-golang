package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// HistoryEntry хранит одну запись истории
type HistoryEntry struct {
	Number  int
	Command string
}

// historyFile возвращает путь к файлу истории
func historyFile() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".shell_history"
	}
	return filepath.Join(home, ".shell_history")
}

// loadHistory читает историю команд из файла
func loadHistory() []HistoryEntry {
	var entries []HistoryEntry
	f, err := os.Open(historyFile())
	if err != nil {
		return entries // Если файл не существует — пустая история
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	num := 1
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) != "" {
			entries = append(entries, HistoryEntry{Number: num, Command: line})
			num++
		}
	}
	return entries
}

// saveHistory записывает команду в файл истории
func saveHistory(cmd string) error {
	f, err := os.OpenFile(historyFile(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprintln(f, cmd)
	return err
}

func main() {
	// Определение флагов командной строки
	helpFlag := flag.Bool("h", false, "Показать справку")
	clearFlag := flag.Bool("c", false, "Очистить историю команд")
	nFlag := flag.Int("n", 0, "Показать последние N записей")
	flag.Parse()

	// Обработка флага помощи
	if *helpFlag {
		fmt.Println("Использование: history [-h] [-c] [-n N] [команда]")
		fmt.Println("  -h      Показать справку")
		fmt.Println("  -c      Очистить историю")
		fmt.Println("  -n N    Показать последние N команд")
		fmt.Println()
		fmt.Println("Для добавления команды: history <команда>")
		fmt.Println("Файл истории:", historyFile())
		os.Exit(0)
	}

	// Очистка истории
	if *clearFlag {
		if err := os.Remove(historyFile()); err != nil && !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "history: не удалось очистить историю: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("История очищена")
		os.Exit(0)
	}

	// Если переданы аргументы — добавляем команду в историю
	args := flag.Args()
	if len(args) > 0 {
		cmd := strings.Join(args, " ")
		if err := saveHistory(cmd); err != nil {
			fmt.Fprintf(os.Stderr, "history: не удалось сохранить команду: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Добавлено в историю: %s\n", cmd)
		os.Exit(0)
	}

	// Загружаем и выводим историю
	entries := loadHistory()

	if len(entries) == 0 {
		fmt.Println("История пуста")
		os.Exit(0)
	}

	// Если указано -n, выводим только последние N записей
	if *nFlag > 0 {
		if *nFlag > len(entries) {
			*nFlag = len(entries)
		}
		entries = entries[len(entries)-*nFlag:]
	}

	for _, e := range entries {
		fmt.Printf("%5d  %s\n", e.Number, e.Command)
	}
}

// Вспомогательная функция для !! и !n: получить команду из истории
func getHistoryCommand(n int) (string, error) {
	entries := loadHistory()
	if len(entries) == 0 {
		return "", fmt.Errorf("история пуста")
	}
	if n == 0 {
		// !! — последняя команда
		return entries[len(entries)-1].Command, nil
	}
	// !n — команда с номером n
	for _, e := range entries {
		if e.Number == n {
			return e.Command, nil
		}
	}
	return "", fmt.Errorf("команда %d не найдена в истории", n)
}

// Функция для получения строки в виде числа
func parseNum(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("недопустимое число: %s", s)
	}
	return n, nil
}
