package atomic

import (
	"fmt"
	"sync"
	"sync/atomic"
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

func BenchmarkAtomic(b *testing.B) {
	var target int32
	wg := sync.WaitGroup{}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
			}()
			atomic.AddInt32(&target, 1)
		}()
	}
	wg.Wait()
	fmt.Printf("%d times: got %d\n", b.N, target)
}
