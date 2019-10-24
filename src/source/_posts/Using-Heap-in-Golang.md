---
title: Using Heap in Golang
date: 2019-10-10 17:04:28
tags: [golang, heap]
keywords:
description:
---


## Preface

在生活中，常会遇到根据优先级分配资源的问题，简单讲就是插队: 

*  登机: 根据舱位决定登机顺序，头等舱/商务舱先登机，不受到达顺序影响
*  银行: VIP优先获得服务
*  急诊室: 根据病情严重程度决定就诊顺序

等等，在开发中，需要用到优先队列解决排队问题，这是一种非常有用的数据结构。

c++内置优先队列可以直接拿来用，如: 

```c++
std::priority_queue<int> pq;
pq.push(2);
pq.pop();
```

因为Golang中没有泛型支持，实现heap有点tricky，我们来一起看看吧。

<!-- more -->

## 概念分辨

Golang给我们提供了"container/heap"，并没有priority queue package，这两者是什么关系呢？ 在课本中，我们学过用近似完全二叉树实现的binary heap, 那么再加上binary heap, tree，这四者究竟是什么关系呢？

可能很多小伙伴看到这里已经不耐烦的点关闭了: 我就想知道怎么实现priority queue, show me the code! 说这么多没用的干嘛？  

这里涉及到一个我非常认同的学习方法：整体学习法，建立结构性知识[^1]，广泛的建立起知识间的连接，这样才能加深理解，就不必每次用完就搁置，下次用再得从头看起啦。

一图胜千言 [^2]

<div align="center">
{% asset_img concept.jpg %}
</div>

>A heap is a specialized tree-based data structure which is essentially an almost complete tree that satisfies the heap property[^3]

Binary Heap是heap的一种实现，所有的heap都是tree的特殊形式(满足heap property)。

那和我们最终需要的priority queue究竟是什么关系呢？

>The heap is one maximally efficient implementation of an abstract data type called a priority queue[^3]

Heap只是priority queue(以下简称PQ)的一种实现罢了。 PQ作为一个抽象数据类型(ADT)，只要实现满足PQ的方法即可。

在Golang中，也同样内置了heap package来实现PQ。

## ADT: Priority Queue 

PQ最基本的操作是push和pop: 

*  Push: 向PQ插入元素，并自动排序
*  Pop: 从PQ头部弹出元素，并保证剩余元素仍有序

有些情况需要动态改变优先级，需要Update方法

本文的目标便是用heap package来实现PQ。

## Heap Internal

简单对binary heap的实现做下回顾, binary heap其实是complete binary tree, 关于full, complete, perfect binary tree的区别，可以看这里[^4]。

如下便是一个complete binary tree: 

<div align="center">
{% asset_img complete_btree.png %}
</div>

它有如下性质: 

*  从左向右依次排列，可以用数组实现
*  root index从0开始, 任何node i有:
    -  parent index: └(i-1)/2┘
    -  left child: $2*i + 1$
    -  right child: $2*i + 2$

可以看出它的特点: 实现简单，可以找child, parent节点。

array中存储的是heap的level order traversal: 

<div align="center">
{% asset_img binaryheap.png %}
</div>

Binary heap同样遵循heap property: 节点"优先级"都比子树的高，因此每次pop都是取头节点。取决有优先级定义: 数值小优先级高的就是小顶堆，反之则是大顶堆。

这里有一个形象的比喻：叠罗汉是小顶堆[^5], 缅怀这篇文章作者，有趣的Vamei大神: 

<div align="center">
{% asset_img heap_metaphor.jpg %}
</div>

叠罗汉就是一个堆, 图中有三个堆，优先级是人的体重，体重最轻的在顶端，它是小顶堆。

Heap最基本的函数:

* Build: 初始建heap, $O(n)$
* Push: 插入新元素，自动排序 $O(log(n))$
* Pop: 弹出root元素，并使剩余元素扔保持heap property $O(log(n))$

## Heap Support in Golang

Golang在"container/heap" package中实现了heap。由于没有泛型的支持，假设基本元素是int, float, struct的三个堆，如果想复用heap代码，只能是通过interface标准化一些操作，让这些差异的操作由用户实现。

Golang heap.Interface封装了差异操作:

```golang
type Interface interface {
    sort.Interface
    Push(x interface{}) // add x as element Len()
    Pop() interface{}   // remove and return element Len() - 1.
}
```

sort.Interface需要我们提供的slice，实现另一个interface，果然是麻烦啊。

基于heap.Interface, heap package提供了我们上节提到的heap基本函数:

```golang
// 建堆
func Init(h Interface)
// 插入元素
func Push(h Interface, x interface{})
// 弹出root元素
func Pop(h Interface) interface{}
// Update元素(包括优先级)
func Fix(h Interface, i int)
// 删除
func Remove(h Interface, i int) interface{}
```

这里值得注意的是不要混淆heap.Push和自己slice实现的Push函数，自己Push仅仅是slice操作，而heap.Push调用了slice的Push操作，还需要额外操作维护heap property。

## Implement a heap in Golang

我们知道了Golang中使用了heap的套路：并不能直接拿来用，而是需要先实现heap.Interface和其中嵌套的sort.Interface，没有generic，很多地方会比较繁琐。但generic又使语言很快变得复杂难懂，果然trade-off无处不在啊，engineer都是天生纠结的命。

这里先看一下我的实现，采用了interface{}作为元素基本类型，可以同时兼容不同类型的元素，用interface{}来替代generic是常用套路。这里需要另外传入一个predicate，类似STL的自定义比较函数，可以通过传入自定义函数来实现对不同数据类型的比较,包括heap of objects，也能在大，小顶堆之间切换，只要改变判断结果即可。

```golang
type myHeap struct {
    data      []interface{}
    predicate func(x, y interface{}) bool
}

func NewHeap(data []interface{}, predicate func(x, y interface{}) bool) *myHeap {
    hp := &myHeap{
        data:      data,
        predicate: predicate,
    }
    heap.Init(hp)
    return hp
}

func (h myHeap) Len() int {
    return len(h.data)
}

func (h myHeap) Swap(i, j int) {
    h.data[i], h.data[j] = h.data[j], h.data[i]
}

func (h myHeap) Less(i, j int) bool {
    return h.predicate(h.data[i], h.data[j])
}

func (h *myHeap) Push(x interface{}) {
    h.data = append(h.data, x)
}

func (h *myHeap) Pop() interface{} {
    len := len(h.data)
    x := h.data[len-1]
    h.data = h.data[:len-1]
    return x
}
```

有了以上模板可以套用，因为interface{}可以指向任何类型。实现heap的关键就变成了实现predicate 函数:

```golang
func(x, y interface{}) bool
```

当x比y的优先级高时(意味着在数组中x将排在y前面)，返回true。

因此小顶堆，小数获得高优先级，即 return x < y；大顶堆，大数获得高优先级，即 return x > y;

## Use heap 

我们接下来看下如何用我们的模板实现heap：

```golang
package main

import (
    "fmt"
    "container/heap"
)

func main() {
    intHeap := NewHeap([]interface{}{3,1,2},
    func(x, y interface{}) bool {
        return x.(int) < y.(int)
    })
    heap.Init(intHeap)
    fmt.Println(intHeap)
    // [1 3 2]
    heap.Push(intHeap, 0)
    fmt.Println(intHeap)
    // [0 1 2 3]
    for intHeap.Len() > 0 {
      fmt.Printf("%+v ", heap.Pop(intHeap)) //注意别写成 intHeap.Pop()
    }
    // 0 1 2 3 
}
```

在 `heap.Init(intHeap)` 之后，我们构建了heap: 

<div align="center">
{% asset_img heap_init.png %}
</div>

注意到slice中的值正是tree的level traversal结果。同时可看到push和pop均符合预期。push得到一个新的heap: 

<div align="center">
{% asset_img heap_adjust.png %}
</div>

回顾heap的push过程，是将新元素放在slice末尾，同时不断将其上升至root的过程。

## Update element in heap

有时候能动态调整heap中元素的优先级会很有用，heap package提供了fix函数: 

```golang
func Fix(h Interface, i int)
```

i为被调整元素的index, 我们需要先调整slice处于i位置的元素，然后调用此函数以维持heap property: 

```golang
intHeap.data = []interface{}{1,3,2} //初始化heap
heap.Init(intHeap)
intHeap.data[1] = -1  // update 1 -> 9
heap.Fix(intHeap, 1)  // 需要知道update元素的index
fmt.Println(intHeap)  // [-1 1 2]
```

但在实际中，我们并不能知道我们要调整元素的index，因为index是在随着元素的push和pop不断的变化，因此就需要记录在元素内部，同时在swap元素时维护index的最新值，只要我们有元素的指针，就可以实现动态update元素的优先级。

完整的例子见[heap doc](https://golang.org/pkg/container/heap/#example__priorityQueue)，同时上面示例代码在[golang playground](https://play.golang.org/p/NP0x2zhwBi8)


[^1]: Young, S. (n.d.). Learn more study less.
[^2]: [堆和树有什么区别？堆为什么要叫堆，不叫树呢？](https://www.zhihu.com/question/36134980)
[^3]: [Heap (data structure)](https://en.wikipedia.org/wiki/Heap_(data_structure))
[^4]: [完美二叉树, 完全二叉树和完满二叉树](https://www.cnblogs.com/idorax/p/6441043.html)
[^5]: [纸上谈兵: 堆 (heap)](https://www.cnblogs.com/vamei/archive/2013/03/20/2966612.html)
[^6]: [binary heap](https://www.geeksforgeeks.org/binary-heap/)