

1.不支持嵌套锁
2.可以一个goroutine lock，另一个goroutine unlock

golang中的互斥锁定义在src/sync/mutex.go

// A Mutex is a mutual exclusion lock.
// Mutexes can be created as part of other structures;
// the zero value for a Mutex is an unlocked mutex.
//
// A Mutex must not be copied after first use.
type Mutex struct {
    state int32
    sema  uint32
}

const (
    mutexLocked = 1 << iota // mutex is locked
    mutexWoken
    mutexWaiterShift = iota
)


[golang中的锁源码实现：Mutex](http://legendtkl.com/2016/10/23/golang-mutex/)
[Locks Aren't Slow; Lock Contention Is](https://preshing.com/20111118/locks-arent-slow-lock-contention-is/)