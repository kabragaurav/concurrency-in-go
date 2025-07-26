package main

import (
	"fmt"
	"sync"
)

var once sync.Once
var instance *Singleton

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			obj := createInstance()
			fmt.Printf("%p\n", obj)
		}()
	}

	wg.Wait()
}

type Singleton struct {
}

func createInstance() *Singleton {
	once.Do(func() {
		instance = &Singleton{}
	})
	return instance
}

/**
Possible Output:
0x1023703a0
0x1023703a0
0x1023703a0
0x1023703a0
0x1023703a0
*/
