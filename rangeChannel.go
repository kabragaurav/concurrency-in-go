package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string)
	greetings := []string{"Hello", "Hi", "Hey", "Hola", "Aloha"}
	go sendToChannel(ch, greetings)
	time.Sleep(2 * time.Second)
	for greeting := range ch {
		fmt.Println("Receive from channel", greeting)
	}
}

func sendToChannel(ch chan string, greetings []string) {
	for _, greeting := range greetings {
		ch <- greeting
	}
	close(ch)
}
