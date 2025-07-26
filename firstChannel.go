package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string) // unbuffered channel
	go greet(ch)
	time.Sleep(1 * time.Second)
	fmt.Println(<-ch)
}

func greet(ch chan string) {
	fmt.Println("Greet starting")
	ch <- "Hello"
	fmt.Println("Greet ending")
}

/**
Since channel is unbuffered, sender (greet goroutine) wait for receiver (main goroutine)

Hence output is just:
Greet starting
Hello
*/
