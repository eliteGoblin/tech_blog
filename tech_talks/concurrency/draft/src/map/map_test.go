package mapdemo

import (
	"testing"
	"time"
)

func TestConcurrentMap(t *testing.T) {
	m := make(map[int]int)
	go func() {
		for {
			_ = m[1]
		}
	}()
	go func() {
		for {
			m[2] = 2
		}
	}()
	time.Sleep(20 * time.Second)
}
