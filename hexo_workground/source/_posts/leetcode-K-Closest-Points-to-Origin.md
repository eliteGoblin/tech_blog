---
title: 'leetcode: K Closest Points to Origin'
date: 2019-07-21 13:14:46
hidden: true
tags: [LeetCode, amazon, heap, sort]
categories:
  - LeetCode
keywords:
description:
---

## Problem

K Closest Points to Origin[^1], We have a list of points on the plane.  Find the K closest points to the origin (0, 0).

(Here, the distance between two points on a plane is the Euclidean distance.)

You may return the answer in any order.  The answer is guaranteed to be unique (except for the order that it is in.)

<!-- more -->

## 分析

给定一个坐标数组，要求返回欧几里得距离最小的K个数，本质是返回数组的TopK问题，只不过排序规则为几何距离。

TopK属于面试经典问题，给定n个元素的array，返回前k个元素，有如下解法: 

*  全部排序: 对数组整体排序，返回前K个元素；$O(nlogn)$
*  部分排序：bubble sort K次，$O(kn)$
*  heap: 构建K容量的max-heap，每次新元素替换heap的root,并fix: $O(logK(n))$
*  partition: 以第K大元素作partition，$O(n)$
*  bucket sort[^2], 用map统计数组中不同值出现的次数; 建立桶：len(array)个不同slot，每次slot存储记录出现次数为slot index的元素array。从后向前遍历桶，得到前K个元素

这篇博客对TopK问题分析的很好: [^3]

## Solution

这里先贴下全部排序法和heap法

### 全部排序

非常straigtforward，关键在于Golang的自定义数组排序，可参考我之前的博客[^4]

```golang
// sort with qsort
type pointsSlice [][]int

func (points pointsSlice) Len() int {
    return len(points)
}

func (points pointsSlice) Less(i, j int) bool {
    distI := distSquare(points[i][0], points[i][1])
    distJ := distSquare(points[j][0], points[j][1])
    return distI < distJ
}

func (points pointsSlice) Swap(i, j int) {
    points[i], points[j] = points[j], points[i]
}

func distSquare(x, y int) int {
    return x*x + y*y
}
func kClosestQSort(points [][]int, K int) [][]int {
    pointsSort := pointsSlice(points)
    sort.Sort(pointsSort)
    return [][]int(pointsSort)[:K]
}
```

### heap

思路如前分析，关键在于如何实现Golang的heap，实现一个slice based的heap, 需要:

*  实现sort.Interface接口
*  实现Append，Pop借口 

比较繁琐，并不能像C++，Java一行代码创建heap/priority queue，后续会有专门文章介绍如何在Golang中实现heap。

```golang

import (
    "container/heap"
)

type maxHeap [][]int

func (heap maxHeap) Len() int {
    return len(heap)
}

func (heap maxHeap) Less(i, j int) bool {
    distI := distSquare(heap[i][1], heap[i][1])
    distJ := distSquare(heap[j][0], heap[j][1])
    return distI >= distJ // max heap
}

func (heap maxHeap) Swap(i, j int) {
    heap[i], heap[j] = heap[j], heap[i]
}

func (heap *maxHeap) Push(v interface{}) {
    *heap = append(*heap, v.([]int))
}

func (heap *maxHeap) Pop() interface{} {
    old := *heap
    n := len(old)
    x := old[n-1]
    *heap = (*heap)[:n-1]
    return x
}

func kClosestHeap(points [][]int, K int) [][]int {
    if len(points) <= K {
        return points
    }
    hp := maxHeap(points[:K])
    heap.Init(&hp) // 以slice为基础建立heap
    for i := K; i < len(points); i++ {
        top := [][]int(points)[0]
        distTop := distSquare(top[0], top[1])
        cur := points[i]
        distCur := distSquare(cur[0], cur[1])
        if distCur < distTop { // 发现更小的元素，替换heap当前最大值
            [][]int(points)[0] = points[i]
            heap.Fix(&hp, 0) // 由于新插入node，重新调整heap
        }
    }
    return [][]int(hp)
}
```

[^1]: [973. K Closest Points to Origin](https://leetcode.com/problems/k-closest-points-to-origin/)  
[^2]: [347_前K个高频元素](https://www.cnblogs.com/xugenpeng/p/9950007.html)  
[^3]: [拜托，面试别再问我TopK了！！！](https://yq.aliyun.com/articles/642891)
[^4]: [Head First Golang Sort](https://elitegoblin.github.io/2017/09/04/golang-sort-top-down-approach/)