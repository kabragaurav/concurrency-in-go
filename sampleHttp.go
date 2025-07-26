package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/hello", Hello)
	fmt.Println("Serving on port 8080")
	http.ListenAndServe(":8080", nil)
}

func Hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

/***
Fire GET http://localhost:8080/hello to get response
*/
