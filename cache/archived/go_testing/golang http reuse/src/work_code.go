package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

func dialOnce() {
	http.DefaultTransport.(*http.Transport).Dial = dial

	for i := 0; i < 3; i++ {
		fmt.Println("loop ", i)
		r, err := http.Get("http://localhost:5995/doSomething")
		if err != nil {
			log.Fatal(err)
		}
		r.Body.Close()
	}
}

func dial(netw, addr string) (net.Conn, error) {
	fmt.Printf("dial %s %s\n", netw, addr)
	return net.Dial(netw, addr)
}
