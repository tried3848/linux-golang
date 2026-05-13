package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	// Определение флагов командной строки
	helpFlag := flag.Bool("h", false, "Показать справку")
	verboseFlag := flag.Bool("v", false, "Подробный вывод")
	destFlag := flag.String("d", ".", "Директория для извлечения файлов")
	flag.Parse()

	// Обработка флага помощи
	if *helpFlag {
		fmt.Println("Использование: unzip [-h] [-v] [-d <директория>] <архив.zip>")
		fmt.Println("  -h        Показать справку")
		fmt.Println("  -v        Подробный вывод")
		fmt.Println("  -d <dir>  Директория назначения")
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "unzip: не указан zip-архив")
		fmt.Fprintln(os.Stderr, "Используйте 'unzip -h' для справки")
		os.Exit(1)
	}

	archiveName := args[0]

	// Проверяем существование архива
	if _, err := os.Stat(archiveName); err != nil {
		fmt.Fprintf(os.Stderr, "unzip: '%s': нет такого файла\n", archiveName)
		os.Exit(1)
	}

	// Открываем zip-архив
	r, err := zip.OpenReader(archiveName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unzip: не удалось открыть архив '%s': %v\n", archiveName, err)
		os.Exit(1)
	}
	defer r.Close()

	// Создаём директорию назначения
	if err := os.MkdirAll(*destFlag, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "unzip: не удалось создать директорию '%s': %v\n", *destFlag, err)
		os.Exit(1)
	}

	fmt.Printf("Извлечение в: %s\n", *destFlag)

	exitCode := 0
	for _, f := range r.File {
		if err := extractFile(f, *destFlag, *verboseFlag); err != nil {
			fmt.Fprintf(os.Stderr, "unzip: ошибка извлечения '%s': %v\n", f.Name, err)
			exitCode = 1
		}
	}

	os.Exit(exitCode)
}

// extractFile извлекает один файл из архива
func extractFile(f *zip.File, destDir string, verbose bool) error {
	// Формируем путь назначения
	destPath := filepath.Join(destDir, filepath.FromSlash(f.Name))

	if verbose {
		fmt.Printf("  извлечение: %s\n", f.Name)
	}

	// Если это директория — создаём её
	if f.FileInfo().IsDir() {
		return os.MkdirAll(destPath, f.Mode())
	}

	// Создаём родительские директории
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return err
	}

	// Открываем файл внутри архива
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	// Создаём файл на диске
	outFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, rc)
	return err
}
