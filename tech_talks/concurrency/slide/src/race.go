package main

import (
	"fmt"
)

func main() {
	ch := make(chan int, 1)
	go func() {
		select {
			case ch <- 1:
			default:
				fmt.Println("msg lost")
		}
	}()
	select {
	case msg := <-ch:
	  fmt.Printf("got msg %d", msg)
	}
}

