package main

import (
	"fmt"
	"time"
)

func main() {
	mskLocation, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		fmt.Println(err)
		return
	}
	now := time.Now().In(mskLocation)

	day := now.Day()
	month := now.Month()
	year := now.Year()
	hour := now.Hour()
	minute := now.Minute()
	second := now.Second()
	location := now.Location()

	fulltime := fmt.Sprintf("%02d:%02d:%02d", hour, minute, second)

	fmt.Println(day, month, year, fulltime, location)
}
