// ps.go — аналог утилиты GNU ps (список процессов)
// Читает данные из /proc на Linux.
//
// Аргументы:
//   -h   — справка
//   -a   — показать все процессы (не только текущего пользователя)
//   -u   — показать столбец с именем пользователя
//   -p <pid> — показать информацию только о процессе с указанным PID

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// ProcessInfo хранит информацию об одном процессе
type ProcessInfo struct {
	PID   int
	PPID  int
	State string
	Name  string
	UID   string
	User  string
	RSS   int64 // RSS в кб
}

func main() {
	// Определение флагов
	helpFlag := flag.Bool("h", false, "Показать справку")
	allProcs := flag.Bool("a", false, "Показать все процессы (не только текущего пользователя)")
	showUser := flag.Bool("u", false, "Показать столбец с именем пользователя")
	pidFilter := flag.Int("p", 0, "Показать только процесс с указанным PID (0 = все)")

	flag.Parse()

	// Вывод справки
	if *helpFlag {
		fmt.Println("Использование: ps [ОПЦИИ]")
		fmt.Println("Показывает список процессов.")
		fmt.Println()
		fmt.Println("Опции:")
		flag.PrintDefaults()
		os.Exit(0)
	}

	// Проверка корректности PID
	if *pidFilter < 0 {
		fmt.Fprintln(os.Stderr, "Ошибка: PID не может быть отрицательным")
		os.Exit(1)
	}

	// Получаем UID текущего пользователя (для фильтрации если не -a)
	currentUID := strconv.Itoa(os.Getuid())

	// Читаем список процессов из /proc
	procs, err := readProcesses()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка чтения процессов: %v\n", err)
		os.Exit(1)
	}

	// Сортируем по PID
	sort.Slice(procs, func(i, j int) bool {
		return procs[i].PID < procs[j].PID
	})

	// Печать заголовка
	if *showUser {
		fmt.Printf("%-8s %-5s %-5s %-1s %-10s %s\n", "USER", "PID", "PPID", "S", "RSS(кб)", "COMMAND")
	} else {
		fmt.Printf("%-5s %-5s %-1s %-10s %s\n", "PID", "PPID", "S", "RSS(кб)", "COMMAND")
	}

	// Вывод процессов
	for _, p := range procs {
		// Фильтр по PID
		if *pidFilter > 0 && p.PID != *pidFilter {
			continue
		}
		// Фильтр по пользователю
		if !*allProcs && p.UID != currentUID {
			continue
		}

		if *showUser {
			fmt.Printf("%-8s %-5d %-5d %-1s %-10d %s\n",
				p.User, p.PID, p.PPID, p.State, p.RSS, p.Name)
		} else {
			fmt.Printf("%-5d %-5d %-1s %-10d %s\n",
				p.PID, p.PPID, p.State, p.RSS, p.Name)
		}
	}

	// Если фильтр по PID — предупреждение если не найден
	if *pidFilter > 0 {
		found := false
		for _, p := range procs {
			if p.PID == *pidFilter {
				found = true
				break
			}
		}
		if !found {
			fmt.Fprintf(os.Stderr, "ps: процесс с PID %d не найден\n", *pidFilter)
			os.Exit(1)
		}
	}
}

// readProcesses читает информацию о процессах из /proc
func readProcesses() ([]ProcessInfo, error) {
	// Находим все числовые директории в /proc (каждая = PID)
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil, err
	}

	var procs []ProcessInfo
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		pid, err := strconv.Atoi(e.Name())
		if err != nil {
			continue // не числовая директория
		}

		p, err := readProcStat(pid)
		if err != nil {
			continue // процесс мог завершиться
		}
		procs = append(procs, p)
	}
	return procs, nil
}

// readProcStat читает /proc/<pid>/status и /proc/<pid>/stat
func readProcStat(pid int) (ProcessInfo, error) {
	info := ProcessInfo{PID: pid}

	statusPath := filepath.Join("/proc", strconv.Itoa(pid), "status")
	f, err := os.Open(statusPath)
	if err != nil {
		return info, err
	}
	defer f.Close()

	// Разбираем /proc/<pid>/status с помощью map
	statusFields := make(map[string]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		parts := strings.SplitN(scanner.Text(), ":", 2)
		if len(parts) == 2 {
			statusFields[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	info.Name = statusFields["Name"]
	info.State = strings.Fields(statusFields["State"])[0]

	// Разбираем PPID
	if v, ok := statusFields["PPid"]; ok {
		info.PPID, _ = strconv.Atoi(v)
	}

	// Разбираем UID
	if v, ok := statusFields["Uid"]; ok {
		fields := strings.Fields(v)
		if len(fields) > 0 {
			info.UID = fields[0]
			// Имя пользователя по UID
			u, err := user.LookupId(info.UID)
			if err == nil {
				info.User = u.Username
			} else {
				info.User = info.UID
			}
		}
	}

	// RSS в кб
	if v, ok := statusFields["VmRSS"]; ok {
		fields := strings.Fields(v)
		if len(fields) > 0 {
			info.RSS, _ = strconv.ParseInt(fields[0], 10, 64)
		}
	}

	return info, nil
}
