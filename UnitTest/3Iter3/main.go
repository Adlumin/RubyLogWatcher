package main

import (
	"fmt"
	"os"
)

func write(loglines []string) {
	file, err := os.Create("production.log")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	for i := range loglines {
		fmt.Fprintf(file, " %s\n", loglines[i])
	}
	fmt.Println("production.log created")
}

func main() {

	loglines := make([]string, 0, 10000)

}
