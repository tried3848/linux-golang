package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Использование: ./mkdir <имя_директории>")
		return
	}

	dirname := os.Args[1]
	
	
	
	// Затем добавляем
	err := os.Mkdir(dirname, 0755)
	if err != nil {
		fmt.Println("Ошибка во время добавления:", err)
	} else {
		fmt.Println("Директория добавлена", dirname)
	}
}
