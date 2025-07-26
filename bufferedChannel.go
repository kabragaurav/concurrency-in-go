package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string) // unbuffered channel
	go aloha(ch)
	time.Sleep(1 * time.Second)
	fmt.Println(<-ch)
}

func aloha(ch chan<- string) {
	fmt.Println("Aloha starting")
	ch <- "Aloha"
	fmt.Println("Aloha ending")
}

/**
Since channel is buffered, sender (greet goroutine) DOES NOT wait for receiver (main goroutine)

Hence output is:
Aloha starting
Aloha ending
Aloha
*/
