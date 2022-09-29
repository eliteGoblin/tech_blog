package main

import (
	"fmt"
	"time"
	"sync"
  "runtime"
)

// START OMIT
func main() {
    runtime.GOMAXPROCS(1)
    var wg sync.WaitGroup
    wg.Add(1)
    wg.Add(1)
	var i int

	go func(){
	  defer wg.Done()
	  for{
		i++
	  }
	}()
	go func() {
	  defer wg.Done()
	  for {
	    fmt.Println("got at", time.Now())
	    time.Sleep(time.Millisecond)
	  }
	}()
	wg.Wait()
	
}
// END OMIT