package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func psnoflag() {
	procDir, err := os.Open("/proc")
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}
	defer procDir.Close()

	entries, err := procDir.Readdirnames(0)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	fmt.Printf("%-8s %-3s %-8s %-20s\n", "PID", "S", "CPU", "COMMAND")
	fmt.Printf("%-8s %-3s %-8s %-20s\n", "---", "-", "---", "-------")

	for _, entry := range entries {
		pid, err := strconv.Atoi(entry)
		if err != nil {
			continue
		}

		comm, state, cpuTotal := getProcessBasicInfo(pid)
		if comm == "" {
			continue
		}

		fmt.Printf("%-8d %-3s %-8d %-20s\n", pid, state, cpuTotal, comm)
	}
}

func psWithFlagF() {
	procDir, err := os.Open("/proc")
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}
	defer procDir.Close()

	entries, err := procDir.Readdirnames(0)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	// Расширенная шапка для -f
	fmt.Printf("%-10s %-8s %-8s %-4s %-8s %-5s %-8s %s\n",
		"UID", "PID", "PPID", "C", "STIME", "TTY", "TIME", "CMD")
	fmt.Printf("%-10s %-8s %-8s %-4s %-8s %-5s %-8s %s\n",
		"---", "---", "----", "-", "-----", "---", "----", "---")

	bootTime := getBootTime()

	for _, entry := range entries {
		pid, err := strconv.Atoi(entry)
		if err != nil {
			continue
		}

		printProcessInfo(pid, bootTime)
	}
}

// Новая функция: вывод информации для конкретного PID в формате -f
func psWithFlagP(targetPid int) {
	// Проверяем, существует ли процесс
	statPath := fmt.Sprintf("/proc/%d/stat", targetPid)
	if _, err := os.Stat(statPath); os.IsNotExist(err) {
		fmt.Printf("Ошибка: процесс с PID %d не найден\n", targetPid)
		return
	}

	// Выводим шапку
	fmt.Printf("%-10s %-8s %-8s %-4s %-8s %-5s %-8s %s\n",
		"UID", "PID", "PPID", "C", "STIME", "TTY", "TIME", "CMD")
	fmt.Printf("%-10s %-8s %-8s %-4s %-8s %-5s %-8s %s\n",
		"---", "---", "----", "-", "-----", "---", "----", "---")

	bootTime := getBootTime()
	printProcessInfo(targetPid, bootTime)
}

// Вспомогательная функция для получения базовой информации о процессе
func getProcessBasicInfo(pid int) (comm string, state string, cpuTotal int64) {
	statPath := fmt.Sprintf("/proc/%d/stat", pid)
	file, err := os.Open(statPath)
	if err != nil {
		return "", "", 0
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	line := scanner.Text()

	openParen := strings.Index(line, "(")
	closeParen := strings.LastIndex(line, ")")

	if openParen == -1 || closeParen == -1 {
		return "", "", 0
	}

	comm = line[openParen+1 : closeParen]
	afterComm := strings.TrimSpace(line[closeParen+1:])
	fields := strings.Fields(afterComm)

	if len(fields) < 14 {
		return "", "", 0
	}

	state = fields[0]
	utime, _ := strconv.ParseInt(fields[12], 10, 64)
	stime, _ := strconv.ParseInt(fields[13], 10, 64)
	cpuTotal = utime + stime

	return comm, state, cpuTotal
}

// Вспомогательная функция для печати информации о процессе (для -f и -p)
func printProcessInfo(pid int, bootTime int64) {
	// Читаем stat файл
	statPath := fmt.Sprintf("/proc/%d/stat", pid)
	statFile, err := os.Open(statPath)
	if err != nil {
		return
	}
	defer statFile.Close()

	scanner := bufio.NewScanner(statFile)
	scanner.Scan()
	statLine := scanner.Text()

	// Парсим stat
	openParen := strings.Index(statLine, "(")
	closeParen := strings.LastIndex(statLine, ")")

	if openParen == -1 || closeParen == -1 {
		return
	}

	comm := statLine[openParen+1 : closeParen]
	afterComm := strings.TrimSpace(statLine[closeParen+1:])
	fields := strings.Fields(afterComm)

	if len(fields) < 22 {
		return
	}

	ppid := fields[1] // PPID
	utime, _ := strconv.ParseInt(fields[12], 10, 64)
	stime, _ := strconv.ParseInt(fields[13], 10, 64)
	starttime, _ := strconv.ParseInt(fields[19], 10, 64)

	// Считаем CPU usage (C) - упрощённо
	cpuUsage := (utime + stime) / 100

	// Считаем STIME (время запуска)
	stimeStr := formatStartTime(starttime, bootTime)

	// Получаем UID
	uid := getUID(pid)

	// Время CPU
	cpuTime := fmt.Sprintf("%02d:%02d:%02d",
		(utime+stime)/3600,
		((utime+stime)%3600)/60,
		(utime+stime)%60)

	// TTY (упрощённо)
	tty := "?"

	fmt.Printf("%-10s %-8d %-8s %-4d %-8s %-5s %-8s %s\n",
		uid, pid, ppid, cpuUsage, stimeStr, tty, cpuTime, comm)
}

func getBootTime() int64 {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return 0
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "btime") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				btime, _ := strconv.ParseInt(fields[1], 10, 64)
				return btime
			}
		}
	}
	return 0
}

func formatStartTime(startticks int64, bootTime int64) string {
	clockTick := int64(100)
	startSeconds := bootTime + (startticks / clockTick)
	t := time.Unix(startSeconds, 0)
	return t.Format("15:04")
}

func getUID(pid int) string {
	statusPath := fmt.Sprintf("/proc/%d/status", pid)
	file, err := os.Open(statusPath)
	if err != nil {
		return "?"
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Uid:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				return fields[1]
			}
		}
	}
	return "?"
}

func printHelp() {
	fmt.Println(`Использование: ps [параметры]
Показывает запушенные процессы

Параметры:
  -h            справка по команде и параметрам
  -f           расширяет вывод, добавляя дополнительную информацию: UID (пользователь, запустивший процесс), 
                                PPID (идентификатор родительского процесса), 
                                C (использование процессора), STIME (время запуска процесса) и другие столбцы. 
  -p            позволяет получить информацию о процессе по его PID.`)
}

func main() {
	// Парсим аргументы командной строки
	if len(os.Args) < 2 {
		psnoflag()
		return
	}

	switch os.Args[1] {
	case "-h":
		printHelp()
	case "-f":
		psWithFlagF()
	case "-p":
		if len(os.Args) < 3 {
			fmt.Println("Ошибка: после -p необходимо указать PID")
			fmt.Println("Пример: ps -p 1234")
			return
		}
		pid, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Printf("Ошибка: PID должен быть числом, получено '%s'\n", os.Args[2])
			return
		}
		psWithFlagP(pid)
	default:
		fmt.Printf("Неизвестный параметр: %s\n", os.Args[1])
		fmt.Println("Используйте -h для получения справки")
	}
}
