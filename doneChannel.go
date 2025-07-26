package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Starting main")
	ch := make(chan struct{})
	go doSomeProcessing(ch)
	<-ch
	fmt.Println("Completing main")
}

func doSomeProcessing(ch chan struct{}) {
	fmt.Println("Starting doSomeProcessing")
	time.Sleep(2 * time.Second) // simulate some work
	fmt.Println("Finished doSomeProcessing")
	close(ch)
}

/**
Output:
Starting main
Starting doSomeProcessing
Finished doSomeProcessing
Completing main
*/
