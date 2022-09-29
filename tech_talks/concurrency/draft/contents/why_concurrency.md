## Ideas Come From

*  book


## Why we need concurrency

*  External driving forces:  In real-world systems many things are happening simultaneously and must be addressed “in real-time” by software
    -  Our team works: 
        +  JIRAs seperate works
*  Intenal Forces(ppt gopher 搬砖 snapshot): 
    -  Running in parallel make tasks much quicker
    -  utilize CPU when one task is waiting for IO

[Concepts:  Concurrency](http://sce.uhcl.edu/helm/rationalunifiedprocess/process/workflow/ana_desi/co_cncry.htm#Why%20are%20we%20interested?)
[Concurrency is not parallelism](https://blog.golang.org/concurrency-is-not-parallelism)


### Concurrency vs Parallelism

Concurrency
Programming as the composition of independently executing processes.

Parallelism
Programming as the simultaneous execution of (possibly related) computations.

Concurrency is about dealing with lots of things at once. breaking it into pieces that can be executed independently.

Parallelism is about doing lots of things at once.

Concurrency is about structure, parallelism is about execution.

Concurrency provides a way to structure a solution to solve a problem that may (but not necessarily) be parallelizable. Parallel need concurrency

Concurrency is a property of the code; parallelism is a property of the running
program

## Why Concurrency is hard

### Race Condition

```
var data int
go func() {
    data++
}()
if data == 0 {
    fmt.Printf("the value is %v.\n", data)
}
```

### Atomicity

the atomicity of an operation can change depending on the currently defined scope

is i++ atomic?
```
# is i++ atomic
func BenchmarkAtomic(b *testing.B) {
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
```

```
go test -bench BenchmarkNotAtomic
```

```
➜  atomic git:(master) ✗ go test -bench BenchmarkAtomic
1 times: got 1
goos: linux
goarch: amd64
BenchmarkAtomic-8       100 times: got 95
10000 times: got 8656
1000000 times: got 827856
5000000 times: got 4172625
 5000000               391 ns/op
PASS
ok      _/home/frankie/notes/blog/cache/concurrency/src/atomic  2.356s
```

incrementing a counter includes

*  reading the current value,
*  incrementing the value in memory,
*  writing back the updated value.

is map operation atomic?
```
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
    time.Sleep(5 * time.Second)
}
// fatal error: concurrent map read and map write
```

Is single read/write atomic? 

[Are C++ Reads and Writes of an int Atomic?](https://stackoverflow.com/questions/54188/are-c-reads-and-writes-of-an-int-atomic)
[Cache coherence](https://en.wikipedia.org/wiki/Cache_coherence)

*  So why do we care? Atomicity is important because if something is atomic, implicitly it is safe within concurrent contexts.
*  Most statements in Go is NOT atomic



// Package atomic provides low-level atomic memory primitives
// useful for implementing synchronization algorithms.
//
// These functions require great care to be used correctly.
// Except for special, low-level applications, synchronization is better
// done with channels or the facilities of the sync package.
// Share memory by communicating;
// don't communicate by sharing memory.
//
// The compare-and-swap operation, implemented by the CompareAndSwapT
// functions, is the atomic equivalent of:
//
//  if *addr == old {
//      *addr = new
//      return true
//  }
//  return false
//



## Deadlock, Live lock, starvation

*  Mutual Exclusion
*  Hold and Wait
*  No Preemption
*  Circular Wait

*  Deadlock例子: 通道 
    -  Process1: 拿A等B
    -  Process2: 拿B等A
    -  waiting forever
*  Livelock例子: 两人通道
    ```
    in a pretty dark room: moves to one side to let you pass, but you’ve just done the same. So you move to the other side, but she’s also done the same. Imagine this going on forever, and you understand livelocks
    ```
*  