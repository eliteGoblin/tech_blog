


Share memory by communicating, don’t communicate by shar‐
ing memory.”
That said, Go does provide traditional locking mechanisms in the sync package. Most
locking issues can be solved using either channels or traditional locks


So which should you use?
Use whichever is most expressive and/or most simple.

inherently more composable than memory access
synchronization primitives.

p47 pic

[Frequently Asked Questions (FAQ)](https://golang.org/doc/faq)
> Regarding mutexes, the sync package implements them, but we hope Go programming
style will encourage people to try higher-level techniques. In particular, consider struc‐
turing your program so that only one goroutine at a time is ever responsible for a par‐
ticular piece of data.


*  gorotine: a few KB
```
package main

import (
    "fmt"
    "runtime"
    "sync"
)

func main() {
    memConsumed := func() uint64 {
        runtime.GC()
        var s runtime.MemStats
        runtime.ReadMemStats(&s)
        return s.Sys
    }

    var c <-chan interface{}
    var wg sync.WaitGroup
    noop := func() { wg.Done(); <-c } // <1>

    const numGoroutines = 1e4 // <2>
    wg.Add(numGoroutines)
    before := memConsumed() // <3>
    for i := numGoroutines; i > 0; i-- {
        go noop()
    }
    wg.Wait()
    after := memConsumed() // <4>
    fmt.Printf("%.3fkb", float64(after-before)/numGoroutines/1000)
}
```
