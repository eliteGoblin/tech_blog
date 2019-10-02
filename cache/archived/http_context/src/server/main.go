package main

import (
	"fmt"
	"net/http"
	"utils"
)

func sayHello(w http.ResponseWriter, r *http.Request) {
	fmt.Println("final server")
	fmt.Println(utils.RequestIDFromContext(r.Context()))
	message := "Final server"
	w.Write([]byte(message))
}
func main() {
	http.HandleFunc("/", sayHello)
	if err := http.ListenAndServe(":8089", nil); err != nil {
		panic(err)
	}
}
