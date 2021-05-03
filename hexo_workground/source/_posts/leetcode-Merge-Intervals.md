---
title: 'leetcode: Merge Intervals'
date: 2019-08-04 11:48:08
hidden: true
tags: [LeetCode, amazon, array]
categories:
  - LeetCode
keywords:
description:
---

## 问题

Merge Intervals[^1], Given a collection of intervals, merge all overlapping intervals.

Example 1:

Input: [[1,3],[2,6],[8,10],[15,18]]
Output: [[1,6],[8,10],[15,18]]
Explanation: Since intervals [1,3] and [2,6] overlaps, merge them into [1,6].

## 分析

给定一系列的interval，要求将其merge，每个interval用含有2元素的slice表示。如何merge呢？基础就是两两merge：有overlap的化为一个interval，没有的保持不变，那如何组织merge过程？即如何两两merge使之最终给出结果呢？

可以先按interval的start排序，然后遍历排序的intervals，给定的例子就是排好序的数组。

每个interval与前面的interval尝试merge，如[1,3],[2,6], 如果有overlap，则merge成为[1, 6]。　如果无overlap，如之前merge的结果[1, 6]和[8, 10]，则保留两者，下一个interval再与[8, 10]比较。

因为interval是按start排序的，i1和i2没有overlap，那么i3肯定不会与i1有overlap: 只有i3.start <= i1.end才有，但是i2.start > i1.end，因此i3.start > i1.end

## Solution

```golang
type intervals [][]int

func (a intervals) Less(i, j int) bool {
    return a[i][0] < a[j][0]
}

func (a intervals) Len() int {
    return len(a)
}

func (a intervals) Swap(i, j int) {
    a[i], a[j] = a[j], a[i]
}

func max(a ...int) int {
    m := a[0]
    for i := 1; i < len(a); i++ {
        if a[i] > m {
            m = a[i]
        }
    }
    return m
}

func merge(ins [][]int) [][]int {
    if len(ins) <= 1 {
        return ins
    }
    intervalArr := intervals(ins)
    sort.Sort(intervalArr)
    res := make([][]int, 0)
    res = append(res, intervalArr[0])
    for i := 1; i < len(intervalArr); i++ {
        if intervalArr[i][0] <= res[len(res)-1][1] {
            res[len(res)-1][1] = max(res[len(res)-1][1], intervalArr[i][1])
        } else {
            res = append(res, intervalArr[i])
        }
    }
    return res
}
```

有关golang如何实现自定义的sort可以看这里[^2]。

[^1]: [Merge Intervals](https://leetcode.com/problems/merge-intervals/)
[^2]: [Head First Golang Sort](https://elitegoblin.github.io/2017/09/04/golang-sort-top-down-approach/)