package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1) // need to wait for one goroutine before exiting main goroutine
	go hey(&wg)
	wg.Wait()
	seeyou()
}

func hey(wg *sync.WaitGroup) {
	defer wg.Done() // equivalent to wg.Add(-1)
	fmt.Println("hey")
}

func seeyou() {
	fmt.Println("seeyou")
}

/**
Deterministically produces output:
hey
seeyou
*/
