// arch.go — аналог утилиты GNU arch (вывод архитектуры процессора)
// Использует пакет runtime для определения архитектуры без вызова внешних команд.
//
// Аргументы:
//   -h  — справка
//   -o  — также вывести название ОС
//   -v  — подробный вывод (arch + os + версия Go runtime)

package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
)

// SysInfo хранит информацию о системе
type SysInfo struct {
	Arch      string // архитектура процессора (GOARCH)
	OS        string // операционная система (GOOS)
	GoVersion string // версия Go runtime
	NumCPU    int    // количество логических процессоров
}

func main() {
	// Определение флагов
	helpFlag := flag.Bool("h", false, "Показать справку")
	showOS := flag.Bool("o", false, "Показать также название операционной системы")
	verbose := flag.Bool("v", false, "Подробный вывод (архитектура, ОС, версия Go, CPU)")

	flag.Parse()

	// Вывод справки
	if *helpFlag {
		fmt.Println("Использование: arch [ОПЦИИ]")
		fmt.Println("Выводит архитектуру аппаратной платформы.")
		fmt.Println()
		fmt.Println("Опции:")
		flag.PrintDefaults()
		fmt.Println()
		fmt.Println("Возможные значения архитектуры: amd64, arm64, 386, arm, mips и др.")
		os.Exit(0)
	}

	// Собираем информацию о системе
	info := SysInfo{
		Arch:      runtime.GOARCH,
		OS:        runtime.GOOS,
		GoVersion: runtime.Version(),
		NumCPU:    runtime.NumCPU(),
	}

	// Режим подробного вывода
	if *verbose {
		fmt.Printf("Архитектура : %s\n", info.Arch)
		fmt.Printf("ОС          : %s\n", info.OS)
		fmt.Printf("Go runtime  : %s\n", info.GoVersion)
		fmt.Printf("CPU (логич.): %d\n", info.NumCPU)
		os.Exit(0)
	}

	// Режим с выводом ОС
	if *showOS {
		fmt.Printf("%s %s\n", info.Arch, info.OS)
		os.Exit(0)
	}

	// Стандартный вывод — только архитектура
	fmt.Println(info.Arch)
}
