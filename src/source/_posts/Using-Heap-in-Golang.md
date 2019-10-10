---
title: Using Heap in Golang
date: 2019-10-10 17:04:28
tags: [golang, heap]
keywords:
description:
---










## Preface

Valid Parentheses[^1]

## Heap internal

## Heap in Golang

### Sort interface

why
```
type Interface interface {
    sort.Interface
    Push(x interface{}) // add x as element Len()
    Pop() interface{}   // remove and return element Len() - 1.
}
```

## Implement a heap

*  小顶堆

max 堆

## General: implement a heap of structs


大顶堆




## Reference

[^1]: [Array Representation Of Binary Heap](https://www.geeksforgeeks.org/array-representation-of-binary-heap/)  