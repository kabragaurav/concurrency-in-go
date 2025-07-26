package main

import (
	"fmt"
	"time"
)

func main() {
	go hi()
	// NEVER use Sleep for aliveness in prod
	time.Sleep(1 * time.Second)
	tata()
}

func hi() {
	fmt.Println("hi")
}

func tata() {
	fmt.Println("tata")
}

/**
In general, produces:
hi
tata
*/
