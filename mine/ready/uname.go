package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	aflag := flag.Bool("v", false, "вся инфа")
	verflag := flag.Bool("ver", false, "версия")
	hflag := flag.Bool("h", false, "справка по команде и параметрам")

	flag.Parse()

	if *hflag {
		printHelp()
		return
	}

	if *verflag {
		unameVflag()
		return
	}

	if *aflag {
		unameAflag()
		return
	}
	unameNOflag()
}

func unameNOflag() {
	unam, err := os.ReadFile("/proc/sys/kernel/ostype")
	if err != nil {
		fmt.Printf("%s\n", err)
	} else {
		fmt.Printf("%s\n", string(unam))
	}
}
func unameVflag() {
	unam, err := os.ReadFile("/proc/sys/kernel/version")
	if err != nil {
		fmt.Printf("%s\n", err)
	} else {
		fmt.Printf("%s\n", string(unam))
	}
}

func unameAflag() {
	ostype, _ := os.ReadFile("/proc/sys/kernel/ostype")
	hostname, _ := os.ReadFile("/proc/sys/kernel/hostname")
	release, _ := os.ReadFile("/proc/sys/kernel/osrelease")
	version, _ := os.ReadFile("/proc/sys/kernel/version")

	// Для domainname (обычно пустой или совпадает с hostname)
	domainname, _ := os.ReadFile("/proc/sys/kernel/domainname")

	// Убираем лишние пробелы и переносы строк
	fmt.Printf("%s %s %s %s %s %s",
		strings.TrimSpace(string(ostype)),
		strings.TrimSpace(string(hostname)),
		strings.TrimSpace(string(domainname)),
		strings.TrimSpace(string(release)),
		strings.TrimSpace(string(version)),
		getMachine(),
	)
}

// Функция для определения архитектуры
func getMachine() string {
	// Пробуем прочитать из /proc/cpuinfo
	cpuinfo, err := os.ReadFile("/proc/cpuinfo")
	if err == nil {
		lines := strings.Split(string(cpuinfo), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "flags") {
				// Определяем архитектуру по флагам CPU
				if strings.Contains(line, " lm ") {
					return "x86_64"
				}
				return "i386"
			}
		}
	}

	// Альтернативный способ через uname -m (вызов внешней команды)
	// Или можно просто вернуть "unknown"
	return "unknown"
}

func printHelp() {
	fmt.Println(`Использование: uname [параметр]
Показывает определенные сведения о системе

Параметры:
  -h            справка по команде и параметрам
  -a			напечатать всю информацию о системе
  -v			напечатать версию ядра`)
}
