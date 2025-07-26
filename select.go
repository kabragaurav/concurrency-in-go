package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(2 * time.Second)
		ch1 <- "Namaste"
	}()

	go func() {
		time.Sleep(1 * time.Second)
		ch2 <- "Namaskar"
	}()

	select {
	case greeting := <-ch1:
		fmt.Println("Receive from ch1", greeting)
	case greeting := <-ch2:
		fmt.Println("Receive from ch2", greeting)
	}
}

/**
Potential output:
Receive from ch2 Namaskar
*/
