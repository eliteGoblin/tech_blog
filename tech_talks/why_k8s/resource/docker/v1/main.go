package main

import (
	"fmt"
	"net/http"
	"os"
)

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Println("receive 1 incoming request")
	name, _ := os.Hostname()
	fmt.Fprintf(w, fmt.Sprintf("Hello from HTTP server V1, host on: %s", name))
}

func main() {

	http.HandleFunc("/", hello)

	http.ListenAndServe(":80", nil)
}
