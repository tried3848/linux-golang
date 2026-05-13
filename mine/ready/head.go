package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	qflag := flag.Bool("q", false, "не печатать заголовки с именами файлов")
	vflag := flag.Bool("v", false, "всегда печатать заголовки с именами файлов")
	nflag := flag.Int("n", 10, "количество строк для вывода (по умолчанию 10)")
	hflag := flag.Bool("h", false, "справка о команде head")
	
	flag.Parse()

	if *hflag {
		printHelp()
		return
	}

	files := flag.Args()
	
	// Если файлов нет, читаем из stdin
	if len(files) == 0 {
		headStdin(*nflag)
		return
	}

	if *qflag {
		headQflag(files, *nflag)
		return
	}

	if *vflag {
		headVflag(files, *nflag)
		return
	}

	// Без флагов - показываем заголовки только если файлов несколько
	if len(files) > 1 {
		headVflag(files, *nflag)
	} else {
		headNoFlags(files[0], *nflag)
	}
}

func headNoFlags(filename string, n int) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "head: %s: %v\n", filename, err)
		return
	}
	lines := strings.Split(string(data), "\n")
	// Выводим первые n строк
	for i := 0; i < n && i < len(lines); i++ {
		fmt.Println(lines[i])
	}
}

func headQflag(files []string, n int) {
	for _, filename := range files {
		data, err := os.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "head: %s: %v\n", filename, err)
			continue
		}
		lines := strings.Split(string(data), "\n")
		for i := 0; i < n && i < len(lines); i++ {
			fmt.Println(lines[i])
		}
	}
}

func headVflag(files []string, n int) {
	for i, filename := range files {
		if i > 0 {
			fmt.Println() // пустая строка между файлами
		}
		fmt.Printf("==> %s <==\n", filename)
		
		data, err := os.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "head: %s: %v\n", filename, err)
			continue
		}
		lines := strings.Split(string(data), "\n")
		for j := 0; j < n && j < len(lines); j++ {
			fmt.Println(lines[j])
		}
	}
}

func headStdin(n int) {
	// Чтение из stdin (если данные переданы через pipe)
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// Данные есть в stdin
		data, err := os.ReadFile(os.Stdin.Name())
		if err != nil {
			fmt.Fprintf(os.Stderr, "head: %v\n", err)
			return
		}
		lines := strings.Split(string(data), "\n")
		for i := 0; i < n && i < len(lines); i++ {
			fmt.Println(lines[i])
		}
	}
}

func printHelp() {
	fmt.Println(`Использование: head [флаг] [файл...]
	
Флаги:
  -n число   количество строк для вывода (по умолчанию 10)
  -q         не печатать заголовки с именами файлов
  -v         всегда печатать заголовки с именами файлов
  -h         справка по команде и параметрам

Примеры:
  head file.txt              # первые 10 строк файла
  head -n 20 file.txt        # первые 20 строк файла
  head -q file1.txt file2.txt # без заголовков
  head -v file.txt           # с заголовком для одного файла
  cat file.txt | head        # чтение из stdin`)
}