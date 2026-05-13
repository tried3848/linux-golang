package main
import(
	"fmt"
	"os"
)

func main(){

	filename := os.Args[1]
	info, err := os.Stat(filename) 
	if filename == ""{
		fmt.Printf("Установить название файла")
		return
	}
	if err != nil {
		fmt.Println("Ошибка: %s", err)
	}
	mode := info.Mode()
	if mode.IsDir(){
		fmt.Printf("%s:directory\n",filename)
	} else if mode.IsRegular(){
		fmt.Printf("%s:file\n",filename)
	}
}
