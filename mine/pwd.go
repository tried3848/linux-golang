package main
import(
	"fmt"
	"os"
)
func main(){
	dir, err := os.Getwd()

	if err != nil{
		fmt.Printf("Ошибка: %s\n", err)
		return
	}
	fmt.Println(dir)
}
	
