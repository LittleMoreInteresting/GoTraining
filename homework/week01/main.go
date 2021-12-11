package main

import (
	"fmt"
	"time"
)

var initTime = time.Now().AddDate(-1, -1, -1)

func main() {

	fmt.Printf("%v\n", time.Since(initTime))
	time.Sleep(2 * time.Second)
	fmt.Printf("%v\n", initTime)
	time.Sleep(2 * time.Second)
	fmt.Printf("%v\n", time.Since(initTime))
	fmt.Printf("%v\n", initTime)
	for i := 0; i < 2; i++ {
		fmt.Printf("%v\n", i)
	}
}
