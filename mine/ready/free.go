package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	bflag := flag.Bool("b", false, "информация в байтах") // Исправил v на b
	mflag := flag.Bool("m", false, "информация в мегабайтах")
	hflag := flag.Bool("h", false, "справка по команде и параметрам")

	flag.Parse()

	if *hflag {
		printHelp()
		return
	}

	if *bflag {
		freeinb()
		return
	}

	if *mflag {
		freeinm()
		return
	}

	// По умолчанию (без флагов) - в килобайтах (как оригинальный free)
	freenoflag()
}

func freenoflag() {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	// Парсим значения
	memInfo := parseMemInfo(string(data))

	// Выводим в килобайтах (оригинальный формат)
	fmt.Println("              total        used        free      shared  buff/cache   available")
	fmt.Printf("Mem:     %10d %10d %10d %10d %10d %10d\n",
		memInfo["MemTotal"]/1024, // в KB
		(memInfo["MemTotal"]-memInfo["MemFree"]-memInfo["Buffers"]-memInfo["Cached"])/1024,
		memInfo["MemFree"]/1024,
		memInfo["Shmem"]/1024,
		(memInfo["Buffers"]+memInfo["Cached"])/1024,
		memInfo["MemAvailable"]/1024,
	)
	fmt.Printf("Swap:    %10d %10d %10d %10d %10d %10d\n",
		memInfo["SwapTotal"]/1024,
		(memInfo["SwapTotal"]-memInfo["SwapFree"])/1024,
		memInfo["SwapFree"]/1024,
		0, 0, 0,
	)
}

func freeinb() {
	// show output in bytes
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	memInfo := parseMemInfo(string(data))

	fmt.Println("              total        used        free      shared  buff/cache   available")
	fmt.Printf("Mem:     %10d %10d %10d %10d %10d %10d\n",
		memInfo["MemTotal"], // в байтах
		memInfo["MemTotal"]-memInfo["MemFree"]-memInfo["Buffers"]-memInfo["Cached"],
		memInfo["MemFree"],
		memInfo["Shmem"],
		memInfo["Buffers"]+memInfo["Cached"],
		memInfo["MemAvailable"],
	)
	fmt.Printf("Swap:    %10d %10d %10d\n",
		memInfo["SwapTotal"],
		memInfo["SwapTotal"]-memInfo["SwapFree"],
		memInfo["SwapFree"],
	)
}

func freeinm() {
	// show output in megabytes
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	memInfo := parseMemInfo(string(data))

	// Переводим в мегабайты (делим на 1024*1024)
	fmt.Println("              total        used        free      shared  buff/cache   available")
	fmt.Printf("Mem:     %10d %10d %10d %10d %10d %10d\n",
		memInfo["MemTotal"]/(1024*1024), // в MB
		(memInfo["MemTotal"]-memInfo["MemFree"]-memInfo["Buffers"]-memInfo["Cached"])/(1024*1024),
		memInfo["MemFree"]/(1024*1024),
		memInfo["Shmem"]/(1024*1024),
		(memInfo["Buffers"]+memInfo["Cached"])/(1024*1024),
		memInfo["MemAvailable"]/(1024*1024),
	)
	fmt.Printf("Swap:    %10d %10d %10d\n",
		memInfo["SwapTotal"]/(1024*1024),
		(memInfo["SwapTotal"]-memInfo["SwapFree"])/(1024*1024),
		memInfo["SwapFree"]/(1024*1024),
	)
}

// Вспомогательная функция для парсинга /proc/meminfo
func parseMemInfo(content string) map[string]int64 {
	memInfo := make(map[string]int64)
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			key := strings.TrimSuffix(parts[0], ":")
			value, _ := strconv.ParseInt(parts[1], 10, 64)
			memInfo[key] = value
		}
	}

	return memInfo
}

func printHelp() {
	fmt.Println(`Использование: free [флаг] 
Флаги:
  -b    	показывает в байтах
  -m    	показывает в мегабайтах
  -h    	справка по команде и параметрам
  
Без флагов показывает в килобайтах (KB)
`)
}
