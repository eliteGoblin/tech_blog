package main

import (
	"fmt"
	"sync"
	"testing"
)

func BenchmarkNotAtomic(b *testing.B) {
	var target int
	wg := sync.WaitGroup{}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
			}()
			target++
		}()
	}
	wg.Wait()
	fmt.Printf("%d times: got %d\n", b.N, target)
}
