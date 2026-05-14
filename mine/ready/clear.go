package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func main() {

	xflag := flag.Bool("v", false, "только видимое")
	verflag := flag.Bool("ver", false, "версия команды")
	hflag := flag.Bool("h", false, "справка по команде и параметрам")

	flag.Parse()

	if *hflag {
		printHelp()
		return
	}

	if *verflag {
		clearVerflag()
		return
	}

	if *xflag {
		clearxflag()
		return
	}
	cleanoflag()
}

func cleanoflag() {
	// Полная очистка экрана и буфера прокрутки
	clearScreen(true)
}

func clearxflag() {
	// Очистка только видимой области, сохранение истории прокрутки
	clearScreen(false)
}

func clearScreen(clearScrollback bool) {
	switch runtime.GOOS {
	case "windows":
		// Windows: используем cmd команды
		clearWindows(clearScrollback)
	default:
		// Linux, macOS, Unix: используем ANSI escape-последовательности
		clearUnix(clearScrollback)
	}
}

func clearWindows(clearScrollback bool) {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()

	// Windows не поддерживает отдельную очистку буфера прокрутки через простые команды
	if !clearScrollback {
		// На Windows можно дополнительно очистить только видимую область
		// но это сложнее, поэтому просто выводим сообщение
		fmt.Print("\033[2J\033[H") // ANSI последовательности для Windows 10+
	}
}

func clearUnix(clearScrollback bool) {
	if clearScrollback {
		// Очистка экрана и буфера прокрутки
		// \033[2J - очистить весь экран
		// \033[3J - очистить буфер прокрутки (поддерживается не везде)
		// \033[H - переместить курсор в верхний левый угол
		fmt.Print("\033[2J\033[3J\033[H")
	} else {
		// Только видимая область
		fmt.Print("\033[2J\033[H")
	}
}

func clearVerflag() {
	fmt.Println(`clear (Go-аналог) 0.2
Аналог утилиты clear на Go
Copyright - свободное ПО
Лицензия - MIT

Особенности:
- Кроссплатформенная поддержка (Windows, Linux, macOS)
- Использует ANSI escape-последовательности
- Параметр -v сохраняет историю прокрутки
- Параметр -h показывает справку`)
}

func printHelp() {
	fmt.Println(`Использование: clear [параметры]

Очищает экран терминала.

Параметры:
  -h            справка по команде и параметрам
  -ver          вывести версию программы
  -v            не пытаться очистить буфер прокрутки (только видимая область)

Примеры:
  clear         # полная очистка экрана и истории
  clear -v      # очистить только видимую область
  clear -h      # показать справку
  clear -ver    # показать версию

Примечания:
  - Без параметров очищается весь экран и буфер прокрутки
  - С параметром -v очищается только видимая область, историю можно вернуть
  - На некоторых терминалах очистка буфера прокрутки может не работать
  - Для Windows требует Windows 10+ с поддержкой ANSI или использует cmd`)
}

// Дополнительная функция: очистка с использованием последовательностей
func clearWithSequence(sequence string) {
	fmt.Print(sequence)
}

// Альтернативная реализация: через внешнюю команду
func clearWithExternalCommand() error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cls")
	default:
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	return cmd.Run()
}

// Дополнительная функция: получение размеров терминала
func getTerminalSize() (width, height int) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return 80, 24 // значения по умолчанию
	}

	// Парсим "rows cols"
	parts := strings.Split(strings.TrimSpace(string(out)), " ")
	if len(parts) == 2 {
		height, _ = strconv.Atoi(parts[0])
		width, _ = strconv.Atoi(parts[1])
	}
	return
}

// Функция для "умной" очистки - запоминает позицию курсора
func clearAndSavePosition() {
	fmt.Print("\033[s")        // сохранить позицию курсора
	fmt.Print("\033[2J\033[H") // очистить экран
	fmt.Print("\033[u")        // восстановить позицию курсора
}
