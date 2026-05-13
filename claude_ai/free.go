// free.go — аналог утилиты GNU free (информация о памяти)
// Читает данные из /proc/meminfo (Linux).
//
// Аргументы:
//   -h      — справка
//   -m      — вывод в мегабайтах (по умолчанию килобайты)
//   -g      — вывод в гигабайтах
//   -t      — показать итоговую строку (Total)

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// MemInfo хранит поля из /proc/meminfo
type MemInfo struct {
	MemTotal     uint64
	MemFree      uint64
	MemAvailable uint64
	Buffers      uint64
	Cached       uint64
	SwapTotal    uint64
	SwapFree     uint64
}

func main() {
	// Определение флагов
	helpFlag := flag.Bool("h", false, "Показать справку")
	megaBytes := flag.Bool("m", false, "Показать значения в мегабайтах")
	gigaBytes := flag.Bool("g", false, "Показать значения в гигабайтах")
	showTotal := flag.Bool("t", false, "Показать итоговую строку (память + своп)")

	flag.Parse()

	// Вывод справки
	if *helpFlag {
		fmt.Println("Использование: free [ОПЦИИ]")
		fmt.Println("Отображает информацию об использовании памяти системы.")
		fmt.Println()
		fmt.Println("Опции:")
		flag.PrintDefaults()
		os.Exit(0)
	}

	// Нельзя одновременно указывать -m и -g
	if *megaBytes && *gigaBytes {
		fmt.Fprintln(os.Stderr, "Ошибка: нельзя одновременно указывать -m и -g")
		os.Exit(1)
	}

	// Чтение /proc/meminfo
	info, err := readMemInfo()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка чтения /proc/meminfo: %v\n", err)
		os.Exit(1)
	}

	// Определяем делитель и единицу измерения
	var divisor uint64 = 1
	unit := "КБ"
	if *megaBytes {
		divisor = 1024
		unit = "МБ"
	} else if *gigaBytes {
		divisor = 1024 * 1024
		unit = "ГБ"
	}

	// Вспомогательная функция для перевода
	conv := func(kb uint64) uint64 {
		if divisor == 0 {
			return 0
		}
		return kb / divisor
	}

	// Вычисляем "используемую" память: total - free - buffers - cached
	used := uint64(0)
	if info.MemTotal >= info.MemFree+info.Buffers+info.Cached {
		used = info.MemTotal - info.MemFree - info.Buffers - info.Cached
	}
	swapUsed := uint64(0)
	if info.SwapTotal >= info.SwapFree {
		swapUsed = info.SwapTotal - info.SwapFree
	}

	// Заголовок
	fmt.Printf("%-10s %10s %10s %10s %10s %10s\n",
		"", "всего", "используется", "свободно", "буферы/кэш", "доступно")
	fmt.Printf("%-10s %10d %10d %10s %10s %10d  [%s]\n",
		"Память:",
		conv(info.MemTotal),
		conv(used),
		fmt.Sprintf("%d", conv(info.MemFree)),
		fmt.Sprintf("%d", conv(info.Buffers+info.Cached)),
		conv(info.MemAvailable),
		unit,
	)
	fmt.Printf("%-10s %10d %10d %10d\n",
		"Своп:",
		conv(info.SwapTotal),
		conv(swapUsed),
		conv(info.SwapFree),
	)

	// Итоговая строка
	if *showTotal {
		fmt.Printf("%-10s %10d %10d %10d\n",
			"Итого:",
			conv(info.MemTotal+info.SwapTotal),
			conv(used+swapUsed),
			conv(info.MemFree+info.SwapFree),
		)
	}
}

// readMemInfo читает и разбирает /proc/meminfo
func readMemInfo() (*MemInfo, error) {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Используем map для хранения значений из файла
	data := make(map[string]uint64)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		// Ключ без двоеточия
		key := strings.TrimSuffix(parts[0], ":")
		val, err := strconv.ParseUint(parts[1], 10, 64)
		if err != nil {
			continue
		}
		data[key] = val
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &MemInfo{
		MemTotal:     data["MemTotal"],
		MemFree:      data["MemFree"],
		MemAvailable: data["MemAvailable"],
		Buffers:      data["Buffers"],
		Cached:       data["Cached"],
		SwapTotal:    data["SwapTotal"],
		SwapFree:     data["SwapFree"],
	}, nil
}
