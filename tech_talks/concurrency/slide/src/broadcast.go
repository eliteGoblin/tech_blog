package main

import (
	"fmt"
	"sync"
)

func main() {
	// DEF_START OMIT
	type Button struct { // <1>
		Clicked *sync.Cond
	}
	button := Button{Clicked: sync.NewCond(&sync.Mutex{})}
	subscribe := func(c *sync.Cond, fn func()) { // <2>
		var goroutineRunning sync.WaitGroup
		goroutineRunning.Add(1)
		go func() {
			goroutineRunning.Done()
			c.L.Lock()
			defer c.L.Unlock()
			c.Wait()
			fn()
		}()
		goroutineRunning.Wait()
	}
	// DEF_END OMIT
	// USE_START OMIT
	var clickRegistered sync.WaitGroup // <3>
	clickRegistered.Add(2)
	subscribe(button.Clicked, func() { // <4>
		fmt.Println("Maximizing window.")
		clickRegistered.Done()
	})
	subscribe(button.Clicked, func() { // <5>
		fmt.Println("Displaying annoying dialogue box!")
		clickRegistered.Done()
	})
	button.Clicked.Broadcast() // <7>
	clickRegistered.Wait()
	// USE_END OMIT
}
