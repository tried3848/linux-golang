package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	Pflag := flag.Bool("P", false, "")
	Lflag := flag.Bool("L", false, "")
	hflag := flag.Bool("h", false, "справка о команде pwd")
	flag.Parse()

	if *hflag {
		printHelp()
		return
	}

	if *Pflag {
		pwdPflag()
		return
	}

	if *Lflag {
		pwdLflag()
		return
	}

	// Если нет флагов, работаем как pwd без флагов (по умолчанию как -L)
	pwdnoflag()
}

func pwdnoflag() {
	// По умолчанию pwd ведёт себя как с флагом -L
	pwdLflag()
}

func pwdLflag() {
	// -L - печатает значение переменной окружения $PWD,
	// если она соответствует текущей рабочей директории
	// (учитывает символьные ссылки).
	
	// Получаем значение переменной окружения PWD
	pwdEnv := os.Getenv("PWD")
	
	if pwdEnv == "" {
		// Если переменная PWD не установлена, используем физический путь
		dir, err := os.Getwd()
		if err != nil {
			fmt.Printf("Ошибка: %s\n", err)
			return
		}
		fmt.Println(dir)
		return
	}
	
	// Проверяем, соответствует ли PWD текущей рабочей директории
	// Получаем физический путь текущей директории
	realDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Ошибка: %s\n", err)
		return
	}
	
	// Преобразуем логический путь из PWD в физический через filepath.EvalSymlinks
	logicalDirPhysical, err := filepath.EvalSymlinks(pwdEnv)
	if err != nil {
		// Если не удалось разрешить символьные ссылки, просто выводим PWD
		fmt.Println(pwdEnv)
		return
	}
	
	// Сравниваем физический путь реальной директории с физическим путём из PWD
	if logicalDirPhysical == realDir {
		// Переменная PWD корректна - выводим её (логический путь с учётом ссылок)
		fmt.Println(pwdEnv)
	} else {
		// PWD не соответствует текущей директории - выводим реальный путь
		fmt.Println(realDir)
	}
}

func pwdPflag() {
	// печатает физический путь,
	// разрешая все символьные ссылки
	// (показывает реальное расположение).
	
	// os.Getwd() возвращает очищенный физический путь
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Ошибка: %s\n", err)
		return
	}
	
	// Дополнительно разрешаем символьные ссылки на всякий случай
	physicalPath, err := filepath.EvalSymlinks(dir)
	if err != nil {
		// Если не удалось разрешить, выводим то, что получили от os.Getwd()
		fmt.Println(dir)
		return
	}
	
	fmt.Println(physicalPath)
}

func printHelp() {
	fmt.Println(`Использование: pwd [флаг]
	
Флаги:
  -L    печатает значение переменной окружения $PWD, 
        если она соответствует текущей рабочей директории 
        (учитывает символьные ссылки).
  -P    печатает физический путь, 
        разрешая все символьные ссылки (показывает реальное расположение).
  -h    справка по команде и параметрам`)
}