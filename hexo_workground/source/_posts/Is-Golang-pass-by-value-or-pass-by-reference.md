---
title: Is Golang pass by value or by reference?
date: 2018-08-05 20:35:30
tags: [golang]
keywords:
description:
---


{% asset_img main.jpg %}


#### Preface

在接触到Golang的某一时刻我们肯定会问自己：Golang是值传递还是引用传递？好吧，Golang其实不像C++那样非此即鄙，要么值传递，要么引用传递。我们一起来看看吧

<!-- more -->


#### 问题分析

为什么需要这个问题呢？我们在写函数时，想知道传入的object是否会被函数改变，比如

```
func aYearPassed(ageOfPeople map[string]int) {
    for name:= range ageOfPeople {
        ageOfPeople[name] ++
    }
}
func main(){
    ageOfPeople := map[string]int {
        "Frank" : 12,
        "Lisha" : 11,
    }
    aYearPassed(ageOfPeople)
    fmt.Println(ageOfPeople)
}
```

本例输出 map[Frank:13 Lisha:12]，是否就能说明Golang是引用传递的呢？为时尚早。我们说传递方式其实是针对Golang的具体数据结构，需要逐一分析

#### Slice, Array, String

我们想知道在赋值时发生了什么，就需要看[Slice结构的实现](https://golang.org/src/runtime/slice.go)：

```
type slice struct {
    array unsafe.Pointer
    len   int
    cap   int
}
```

{% asset_img slice.jpg %}

可见slice是struct，存放len, cap信息，数据由array指针指向。形如s[2:4]其实是新生成了struct header: len为2，array指向了其数据地址（旧数据后移2个位置）：

{% asset_img slice_new.jpg %}

而slice赋值给新slc，本质上是浅拷贝：生成新的slice header，复制之前的值，包括数据地址。值得注意的是，slice数据会自动伸缩，可能造成重新分配：导致之前浅拷贝的两个slice最终指向不一样的数据地址

但是array结构没有采用指针指向数据地址的结构，因此array拷贝两者修改独立：

```
slcSrc := []int{1, 2, 3}
fmt.Println(slcSrc)
slcDst := slcSrc
slcDst[1] = 999
fmt.Println(slcDst)
fmt.Println(slcSrc)

primes := [6]int{2, 3, 5, 7, 11, 13}
fmt.Println(primes)
primesChanged := primes
primesChanged[2] = 999
fmt.Println(b)
fmt.Println(primes)
```

从输出结果看来，slice在没有发生数据空间重新分配时，两者共享数据，看起来像是reference；而array是值拷贝，赋值后两者独立变化。

string也类似slice，[代码](https://golang.org/src/runtime/string.go)：

```
type stringStruct struct {
    str unsafe.Pointer
    len int
}
```


#### Map

同理，看[Map实现代码](https://golang.org/src/runtime/hashmap.go)，用hash table实现：类似slice，仍其实是header，有指针指向数据的结构：

```
// A header for a Go map.
type hmap struct {
    // Note: the format of the Hmap is encoded in ../../cmd/internal/gc/reflect.go and
    // ../reflect/type.go. Don't change this structure without also changing that code!
    count     int // # live cells == size of map.  Must be first (used by len() builtin)
    flags     uint8
    B         uint8  // log_2 of # of buckets (can hold up to loadFactor * 2^B items)
    noverflow uint16 // approximate number of overflow buckets; see incrnoverflow for details
    hash0     uint32 // hash seed

    buckets    unsafe.Pointer // array of 2^B Buckets. may be nil if count==0.
    oldbuckets unsafe.Pointer // previous bucket array of half the size, non-nil only when growing
    nevacuate  uintptr        // progress counter for evacuation (buckets less than this have been evacuated)

    extra *mapextra // optional fields
}
```

同样的浅拷贝，同样可能因内存扩展而重新分配内存地址

#### Channel

[实现代码](https://github.com/golang/go/blob/4fc9565ffce91c4299903f7c17a275f0786734a1/src/runtime/chan.go#L17-L29)，核心是hchan结构:

```
type hchan struct {
    qcount   uint           // total data in the queue
    dataqsiz uint           // size of the circular queue
    buf      unsafe.Pointer // points to an array of dataqsiz elements
    elemsize uint16
    closed   uint32
    elemtype *_type // element type
    sendx    uint   // send index
    recvx    uint   // receive index
    recvq    waitq  // list of recv waiters
    sendq    waitq  // list of send waiters
    lock     mutex
}
```

由于channel在make时已经指定了cap，因此不会重新分配数据

#### 结论

*  **No Pass by Ref in Golang**: Golang本质是值传递，但因为数据结构大多都采用header，指针指向数据的方式，很多时候看起来像引用传递
*  Array拷贝的时候就是值拷贝
*  Slice，Map应认为是值传递，尽管看起来像是引用传递，因为可能因内存重新分配造成原结构和被赋值的指向不是同一地址
*  Channel可以认为等同于引用传递，因为数据不可能发生重新分配

到目前为止，Golang的map还有一个["内存泄露的bug"](https://github.com/golang/go/issues/20135)；

#### Reference



[Go Slices: usage and internals](https://blog.golang.org/go-slices-usage-and-internals)  
[Hash table](https://en.wikipedia.org/wiki/Hash_table)  
[runtime: maps do not shrink after elements removal (delete)](https://github.com/golang/go/issues/20135)