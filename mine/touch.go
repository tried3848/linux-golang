package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Использование: ./touch <имя_директории>")
		return
	}

	filename := os.Args[1]
	// Создать пустой файл
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Ошибка во время создания:", err)
	} else {
		fmt.Println("Файл создан", err)
	}
	file.Close()

	// Обновить время изменения
	now := time.Now()
	err = os.Chtimes("file.txt", now, now)

}
