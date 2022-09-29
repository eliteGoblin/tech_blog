
```
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
```

```
go test -bench BenchmarkAtomic
1 times: got 1
goos: linux
goarch: amd64
BenchmarkAtomic-8       100 times: got 100
10000 times: got 10000
1000000 times: got 1000000
5000000 times: got 5000000
 5000000               393 ns/op
PASS
ok      _/home/frankie/notes/blog/cache/concurrency/src/atomic  2.366s
```

More sophisticated synchronization construction use atomic primitives undercover. Like Mutex and WaitGroup

one of those operations
performs a write, both threads must use atomic
operations

how atomic is implemented? CAS

low-level programming: CPU command: 

大部分操作系统支持CAS，x86指令集上的CAS汇编指令是CMPXCHG