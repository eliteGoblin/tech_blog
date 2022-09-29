package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
)

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Println("receive 1 incoming request")
	name, _ := os.Hostname()

	var res int
	for i := 0; i < 10000; i++ {
		for j := 0; j < 1000; j++ {
			res += rand.Int()
		}
	}
	fmt.Fprintf(w, fmt.Sprintf("Hello from HTTP server V2, host on: %s, sum %d", name, res))
}

func main() {

	http.HandleFunc("/", hello)
	log.Fatal(http.ListenAndServe(":80", nil))
}
