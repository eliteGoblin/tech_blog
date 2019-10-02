---
title: Review Slice Internal When Solving A Golang Puzzle
date: 2018-12-18 15:48:53
tags: [Golang]
keywords:
description:
---

## Preface

在一次日常编程练习中，碰到一个开始觉得颇为奇怪的Golang behaviour，在解决问题过程中，弄清楚了两个Golang比较细节的feature，很有趣的experience，分享给大家。

<!-- more -->

## Peculiar Issue

下面代码，首先将array作为map的key，遍历map的key，将array插入到二维slice中，见[playground](https://play.golang.org/p/AlICQA8tpYX)：
```
mp := make(map[[3]int]bool)
mp[[3]int{1, 2, 3}] = true
mp[[3]int{4, 5, 6}] = true
mp[[3]int{7, 8, 9}] = true
res := make([][]int, 0)
for k := range mp {
    fmt.Printf("key is %+v\n", k[:])
    res = append(res, k[:])
    fmt.Printf("after append: %+v\n", res)
}
```

运行出乎意料，看来后续的append将之前的全部覆盖了，这是怎么回事？

```
key is [1 2 3]
after append: [[1 2 3]]
key is [4 5 6]
after append: [[4 5 6] [4 5 6]]
key is [7 8 9]
after append: [[7 8 9] [7 8 9] [7 8 9]]
```

## Asking a problem

百思不得其解，于是在go-nuts上提问了一发，很快得到回答，不得不说上面真是高手云集，氛围也相当好[^1]。

## Got a answer

输出涉及到Golang的两个feature: loop variable reuse及slice internal，我们逐一分析

### Loop variable reuse

我们loop over map key的代码其实可以展开成为，见[playground](https://play.golang.org/p/PATwkcN-LP0): 

```
a1 := [3]int{1, 2, 3} 
a2 := [3]int{4, 5, 6} 
a3 := [3]int{7, 8, 9} 

s := make([][]int, 0) 
a := a1 
s = append(s, a[:]) 
a = a2 
s = append(s, a[:]) 
a = a3 
s = append(s, a[:]) 

fmt.Println(s) 
```

同样能重现出被后续append覆盖的现象。我们注意被append到\[\]\[\]int的a被reuse，如果换成每次append新\[3\]int问题就不会出现，见[playground](https://play.golang.org/p/UxltLI8hKaj)

```
s := make([][]int, 0) 
s = append(s, a1[:]) 
s = append(s, a2[:]) 
s = append(s, a3[:])
fmt.Println(s)
```

可见是variable reuse引发的问题，but wait: 

```
for k := range mp {
    ...
}
```

这个for怎么看着都像每次新建variable，不是吗？  

还真的不是，参见这个issue，有位高手是这么回复的：

> It is, unfortunately, a very common "gotcha". 

And wtf is a Gotcha? wikipedia说道: 

> Gotcha，在计算机编程领域指在系统、或程序、程序设计语言中一个合法有效的构造但是反直觉、总是容易造成错误，易于使用但其结果是不期望的或者不合逻辑的。 字面上是got you的简写，常用于口语，直译为： “逮着你了”、“捉弄到你了”、“你中计了” 、“骗到你了”

社区已经有人提议改成每次都新建variable而不是reuse了[^2]

但为什么reuse了一个slice，就会造成覆盖的结果呢？这个涉及到slice的内部结构。

### Slice internals

我们看下Slice的内部结构[^3]：

{% asset_img slice.png %}

我们append时：

```
s := make([][]int, 0) 
s = append(s, a[:]) 
```

其实是将slice的header struct的指针append到二维数组，这里\[\]\[\]int非常像c中的指针数组，由于我们复用了slice struct的指针，我们append到\[\]\[\]int的其实是**同样的slice struct的指针**，同时slice struct包含指向同一片数据的指针，这就是为什么会出现这样的输出：

```
after append: [[7 8 9] [7 8 9] [7 8 9]]
```

示意图：

{% asset_img sliceofslice.png %}

[^1]: [Append array to slice of slice behavior puzzles me](https://groups.google.com/forum/#!topic/golang-nuts/GwL3XpcI02Y)
[^2]: [proposal: spec: redefine range loop variables in each iteration](https://github.com/golang/go/issues/20733)  
[^3]: [Go Slices: usage and internals](https://blog.golang.org/go-slices-usage-and-internals)