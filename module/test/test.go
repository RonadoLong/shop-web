package main

import (
	"time"
	"fmt"
	"strconv"
)

func main()  {

	t := time.Now()
	year, month, day := t.Date()
	thisMonthFirstDay := time.Date(year, month, day+1, 0, 0, 0, 0, t.Location())
	fmt.Println(thisMonthFirstDay)

	fmt.Println(thisMonthFirstDay.Sub(t).Seconds())



	i := float64(123.12)

	string := strconv.FormatFloat(i, 'E', -1, 64)
	int,_:=strconv.Atoi(string)

	fmt.Println(int)
}