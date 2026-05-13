package main

import (
	"archive/tar"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	// Определение флагов командной строки
	helpFlag := flag.Bool("h", false, "Показать справку")
	createFlag := flag.Bool("c", false, "Создать архив")
	extractFlag := flag.Bool("x", false, "Извлечь архив")
	gzipFlag := flag.Bool("z", false, "Использовать gzip сжатие")
	verboseFlag := flag.Bool("v", false, "Подробный вывод")
	fileFlag := flag.String("f", "", "Имя архива")
	flag.Parse()

	// Обработка флага помощи
	if *helpFlag {
		fmt.Println("Использование: tar [-h] [-c|-x] [-z] [-v] -f <архив> [файлы...]")
		fmt.Println("  -h  Показать справку")
		fmt.Println("  -c  Создать архив")
		fmt.Println("  -x  Извлечь архив")
		fmt.Println("  -z  Использовать gzip сжатие (.tar.gz)")
		fmt.Println("  -v  Подробный вывод")
		fmt.Println("  -f  Имя файла архива (обязательно)")
		os.Exit(0)
	}

	// Проверяем обязательный аргумент -f
	if *fileFlag == "" {
		fmt.Fprintln(os.Stderr, "tar: необходимо указать имя архива (-f)")
		os.Exit(1)
	}

	// Проверяем, что указана операция
	if !*createFlag && !*extractFlag {
		fmt.Fprintln(os.Stderr, "tar: необходимо указать операцию: -c (создать) или -x (извлечь)")
		os.Exit(1)
	}

	if *createFlag && *extractFlag {
		fmt.Fprintln(os.Stderr, "tar: нельзя одновременно использовать -c и -x")
		os.Exit(1)
	}

	if *createFlag {
		// Создание архива
		args := flag.Args()
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "tar: не указаны файлы для архивирования")
			os.Exit(1)
		}
		if err := createTar(*fileFlag, args, *gzipFlag, *verboseFlag); err != nil {
			fmt.Fprintf(os.Stderr, "tar: ошибка создания архива: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Извлечение архива
		if err := extractTar(*fileFlag, *gzipFlag, *verboseFlag); err != nil {
			fmt.Fprintf(os.Stderr, "tar: ошибка извлечения архива: %v\n", err)
			os.Exit(1)
		}
	}
}

// createTar создаёт tar-архив из указанных файлов
func createTar(archive string, sources []string, useGzip bool, verbose bool) error {
	f, err := os.Create(archive)
	if err != nil {
		return err
	}
	defer f.Close()

	var tw *tar.Writer
	if useGzip {
		// Используем gzip-сжатие
		gw := gzip.NewWriter(f)
		defer gw.Close()
		tw = tar.NewWriter(gw)
	} else {
		tw = tar.NewWriter(f)
	}
	defer tw.Close()

	for _, src := range sources {
		if err := addToTar(tw, src, verbose); err != nil {
			return fmt.Errorf("ошибка добавления '%s': %w", src, err)
		}
	}
	return nil
}

// addToTar рекурсивно добавляет файлы в tar-архив
func addToTar(tw *tar.Writer, path string, verbose bool) error {
	return filepath.Walk(path, func(file string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Создаём заголовок
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Name = filepath.ToSlash(file)

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if verbose {
			fmt.Printf("  добавление: %s\n", file)
		}

		// Копируем содержимое файла (директории не копируем)
		if !info.IsDir() {
			f, err := os.Open(file)
			if err != nil {
				return err
			}
			defer f.Close()
			if _, err := io.Copy(tw, f); err != nil {
				return err
			}
		}
		return nil
	})
}

// extractTar извлекает tar-архив
func extractTar(archive string, useGzip bool, verbose bool) error {
	f, err := os.Open(archive)
	if err != nil {
		return err
	}
	defer f.Close()

	var tr *tar.Reader
	if useGzip {
		gr, err := gzip.NewReader(f)
		if err != nil {
			return fmt.Errorf("ошибка gzip: %w", err)
		}
		defer gr.Close()
		tr = tar.NewReader(gr)
	} else {
		tr = tar.NewReader(f)
	}

	// Итерируемся по записям архива
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break // Конец архива
		}
		if err != nil {
			return err
		}

		destPath := filepath.FromSlash(header.Name)

		if verbose {
			fmt.Printf("  извлечение: %s\n", header.Name)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(destPath, os.FileMode(header.Mode)); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
				return err
			}
			outFile, err := os.Create(destPath)
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
		}
	}
	return nil
}
