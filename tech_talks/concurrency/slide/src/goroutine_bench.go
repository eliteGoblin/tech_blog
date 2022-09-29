package main

import (
	"sync"
	"testing"
)

var wg sync.WaitGroup
var begin = make(chan struct{})
var c = make(chan struct{})
var token struct{}

func BenchmarkContextSwitch(b *testing.B) {
	sender := func() {
		defer wg.Done()
		<-begin // <1>
		for i := 0; i < b.N; i++ {
			c <- token // <2>
		}
	}
	receiver := func() {
		defer wg.Done()
		<-begin // <1>
		for i := 0; i < b.N; i++ {
			<-c // <3>
		}
	}
	wg.Add(2)
	go sender()
	go receiver()
	b.StartTimer() // <4>
	close(begin)   // <5>
	wg.Wait()
}
