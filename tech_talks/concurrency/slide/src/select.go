package main

import (
	"fmt"
)

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	close(ch1)
	close(ch2)
	count1, count2 := 0, 0
	for i := 0; i < 100000; i++ {
		select {
		case <-ch1:
			count1++
		case <-ch1:
			count1++
		case <-ch1:
			count1++
		case <-ch2:
			count2++

		}
	}
	fmt.Printf("ch1 is %d, ch2 is %d, rate is %f",
		count1, count2, float32(count1)/float32(count2))
}
