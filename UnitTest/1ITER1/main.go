package main

import (
	"fmt"
	"time"
)

func main() {
	for i := 0; i < 5; i++ {
		now := time.Now().Format("2006-01-02T15:04:05.000000")
		fmt.Println(now)
	}
}
