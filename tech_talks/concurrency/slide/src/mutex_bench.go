package main

import (
	"fmt"
	"sync"
	"testing"
)

// START OMIT
func BenchmarkMutex(b *testing.B) {
	var target int32
	mutex := sync.Mutex{}
	wg := sync.WaitGroup{}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
			}()
			mutex.Lock()
			target++
			mutex.Unlock()
		}()
	}
	wg.Wait()
	fmt.Printf("%d times: got %d\n", b.N, target)
}
// END OMIT