package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Использование: ./rmdir <имя_директории>")
		return
	}

	dirname := os.Args[1]
	
	// Сначала проверяем существование
	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		fmt.Printf("Директория '%s' не существует\n", dirname)
		return
	}
	
	// Затем удаляем
	err := os.RemoveAll(dirname)
	if err != nil {
		fmt.Println("Ошибка во время удаления:", err)
	} else {
		fmt.Println("Директория удалена", dirname)
	}
}
