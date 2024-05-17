package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<html><body><h1>Hello, World2!</h1></body></html>")
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Server listening on port 8081...")
	http.ListenAndServe("localhost:8081", nil)
}
