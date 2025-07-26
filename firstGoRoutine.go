package main

import "fmt"

func main() {
	go hello()
	bye()
}

func hello() {
	fmt.Println("hello")
}

func bye() {
	fmt.Println("bye")
}

/**
Can produce different outputs like:

1.
bye

2.
bye
hello
*/
