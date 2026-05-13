package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// FileSignature описывает сигнатуру (magic bytes) для определения типа файла
type FileSignature struct {
	Offset  int
	Magic   []byte
	MimeType string
	Desc    string
}

// Список известных сигнатур файлов
var knownSignatures = []FileSignature{
	{0, []byte{0x7F, 0x45, 0x4C, 0x46}, "application/x-elf", "ELF исполняемый файл"},
	{0, []byte{0xFF, 0xD8, 0xFF}, "image/jpeg", "JPEG изображение"},
	{0, []byte{0x89, 0x50, 0x4E, 0x47}, "image/png", "PNG изображение"},
	{0, []byte{0x47, 0x49, 0x46, 0x38}, "image/gif", "GIF изображение"},
	{0, []byte{0x25, 0x50, 0x44, 0x46}, "application/pdf", "PDF документ"},
	{0, []byte{0x50, 0x4B, 0x03, 0x04}, "application/zip", "ZIP архив"},
	{0, []byte{0x1F, 0x8B}, "application/gzip", "GZIP архив"},
	{0, []byte{0x42, 0x5A, 0x68}, "application/x-bzip2", "BZIP2 архив"},
	{0, []byte{0x75, 0x73, 0x74, 0x61, 0x72}, "application/x-tar", "TAR архив"},
	{0, []byte{0xD0, 0xCF, 0x11, 0xE0}, "application/msword", "Microsoft Office документ"},
	{0, []byte{0x23, 0x21}, "text/x-script", "Скрипт (shebang)"},
	{0, []byte{0xCA, 0xFE, 0xBA, 0xBE}, "application/java", "Java class файл"},
}

func main() {
	// Определение флагов командной строки
	helpFlag := flag.Bool("h", false, "Показать справку")
	mimeFlag := flag.Bool("m", false, "Вывести MIME-тип вместо описания")
	briefFlag := flag.Bool("b", false, "Краткий вывод (без имени файла)")
	flag.Parse()

	// Обработка флага помощи
	if *helpFlag {
		fmt.Println("Использование: file [-h] [-m] [-b] <файл> [файл2 ...]")
		fmt.Println("  -h  Показать справку")
		fmt.Println("  -m  Вывести MIME-тип")
		fmt.Println("  -b  Краткий вывод (без имени файла)")
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "file: не указан файл")
		fmt.Fprintln(os.Stderr, "Используйте 'file -h' для справки")
		os.Exit(1)
	}

	exitCode := 0
	for _, path := range args {
		result, err := detectFile(path, *mimeFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "file: '%s': %v\n", path, err)
			exitCode = 1
			continue
		}

		if *briefFlag {
			fmt.Println(result)
		} else {
			fmt.Printf("%s: %s\n", path, result)
		}
	}
	os.Exit(exitCode)
}

// detectFile определяет тип файла по magic bytes и расширению
func detectFile(path string, mime bool) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("нет такого файла или директории")
	}

	// Специальные типы
	if info.IsDir() {
		if mime {
			return "inode/directory", nil
		}
		return "директория", nil
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return "символическая ссылка", nil
	}

	// Открываем файл и читаем первые байты
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("нет прав на чтение")
	}
	defer f.Close()

	// Читаем до 512 байт для определения типа
	buf := make([]byte, 512)
	n, _ := f.Read(buf)
	buf = buf[:n]

	// Проверяем сигнатуры
	for _, sig := range knownSignatures {
		if sig.Offset+len(sig.Magic) <= len(buf) {
			match := true
			for i, b := range sig.Magic {
				if buf[sig.Offset+i] != b {
					match = false
					break
				}
			}
			if match {
				if mime {
					return sig.MimeType, nil
				}
				return sig.Desc, nil
			}
		}
	}

	// Проверяем, является ли файл текстовым (все байты ASCII или UTF-8 printable)
	isText := true
	for _, b := range buf {
		if b < 0x09 || (b > 0x0D && b < 0x20 && b != 0x1B) {
			isText = false
			break
		}
	}

	if isText {
		// Определяем по расширению
		ext := strings.ToLower(path)
		switch {
		case strings.HasSuffix(ext, ".go"):
			if mime { return "text/x-go", nil }
			return "Go исходный код", nil
		case strings.HasSuffix(ext, ".py"):
			if mime { return "text/x-python", nil }
			return "Python скрипт", nil
		case strings.HasSuffix(ext, ".sh"):
			if mime { return "text/x-shellscript", nil }
			return "Shell скрипт", nil
		case strings.HasSuffix(ext, ".json"):
			if mime { return "application/json", nil }
			return "JSON данные", nil
		case strings.HasSuffix(ext, ".html") || strings.HasSuffix(ext, ".htm"):
			if mime { return "text/html", nil }
			return "HTML документ", nil
		case strings.HasSuffix(ext, ".xml"):
			if mime { return "text/xml", nil }
			return "XML документ", nil
		case strings.HasSuffix(ext, ".csv"):
			if mime { return "text/csv", nil }
			return "CSV данные", nil
		default:
			if mime { return "text/plain", nil }
			return "текстовый файл", nil
		}
	}

	if mime {
		return "application/octet-stream", nil
	}
	return "бинарный файл", nil
}
