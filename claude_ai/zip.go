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
	recursiveFlag := flag.Bool("r", false, "Рекурсивно добавить директории")
	flag.Parse()

	// Обработка флага помощи
	if *helpFlag {
		fmt.Println("Использование: zip [-h] [-v] [-r] <архив.zip> <файл/директория> [...]")
		fmt.Println("  -h  Показать справку")
		fmt.Println("  -v  Подробный вывод")
		fmt.Println("  -r  Рекурсивно добавить директории")
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "zip: необходимо указать имя архива и хотя бы один файл")
		fmt.Fprintln(os.Stderr, "Используйте 'zip -h' для справки")
		os.Exit(1)
	}

	archiveName := args[0]
	sources := args[1:]

	// Создаём zip-файл
	zipFile, err := os.Create(archiveName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "zip: не удалось создать архив '%s': %v\n", archiveName, err)
		os.Exit(1)
	}
	defer zipFile.Close()

	// Создаём zip writer
	zw := zip.NewWriter(zipFile)
	defer zw.Close()

	exitCode := 0
	for _, src := range sources {
		info, err := os.Stat(src)
		if err != nil {
			fmt.Fprintf(os.Stderr, "zip: '%s': нет такого файла\n", src)
			exitCode = 1
			continue
		}

		if info.IsDir() {
			if !*recursiveFlag {
				fmt.Fprintf(os.Stderr, "zip: '%s': директория (используйте -r для рекурсии)\n", src)
				exitCode = 1
				continue
			}
			// Рекурсивно добавляем директорию
			if err := addDirToZip(zw, src, *verboseFlag); err != nil {
				fmt.Fprintf(os.Stderr, "zip: ошибка добавления директории '%s': %v\n", src, err)
				exitCode = 1
			}
		} else {
			// Добавляем отдельный файл
			if err := addFileToZip(zw, src, *verboseFlag); err != nil {
				fmt.Fprintf(os.Stderr, "zip: ошибка добавления файла '%s': %v\n", src, err)
				exitCode = 1
			}
		}
	}

	os.Exit(exitCode)
}

// addFileToZip добавляет один файл в zip-архив
func addFileToZip(zw *zip.Writer, path string, verbose bool) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return err
	}

	// Создаём заголовок файла в архиве
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = filepath.ToSlash(path)
	header.Method = zip.Deflate

	writer, err := zw.CreateHeader(header)
	if err != nil {
		return err
	}

	if verbose {
		fmt.Printf("  добавление: %s\n", path)
	}

	_, err = io.Copy(writer, f)
	return err
}

// addDirToZip рекурсивно добавляет директорию в zip-архив
func addDirToZip(zw *zip.Writer, dir string, verbose bool) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		return addFileToZip(zw, path, verbose)
	})
}
