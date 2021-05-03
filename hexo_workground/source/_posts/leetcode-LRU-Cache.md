---
title: 'leetcode: LRU Cache'
date: 2019-07-08 15:20:25
hidden: true
tags: [LeetCode, amazon, list]
categories:
  - LeetCode
keywords:
description:
---

## Problem

LRU Cache[^1], Design and implement a data structure for Least Recently Used (LRU) cache. It should support the following operations: get and put.

get(key) - Get the value (will always be positive) of the key if the key exists in the cache, otherwise return -1.
put(key, value) - Set or insert the value if the key is not already present. When the cache reached its capacity, it should invalidate the least recently used item before inserting a new item.

<!-- more -->

The cache is initialized with a positive capacity.

Follow up:
Could you do both operations in O(1) time complexity?

## Solution

要求O(1)复杂度，第一次接触此题，没有思路，O(1)的Get可以用map实现，但是O(1)的Put，尤其如何是当capacity满时，如何替换least recent use的cache。

思路: 用list来记录元素，most recent element排在head，而least recent used元素排在末尾。这样置换时，直接找到list的tail元素替换即可。用map来记录key-->element的映射，Get/Put时先找map，找到则返回或更新元素值。注意的点: 

*  用container/list，双链表，同时存储head,tail的链接，O(1)时间找head, tail
*  自定义list.Element，为一个struct结构，存储key值。因为置换LRU cache时，需要从map中删除此value对应的key

```golang
type LRUCache struct {
    mp       map[int]*list.Element
    list     *list.List
    capacity int
}

func Constructor(capacity int) LRUCache {
    return LRUCache{
        mp:       make(map[int]*list.Element),
        list:     list.New(),
        capacity: capacity,
    }
}

type element struct {
    key   int
    value int
}

func (this *LRUCache) Get(key int) int {
    if e, ok := this.mp[key]; ok {
        this.list.MoveToFront(e)
        return e.Value.(*element).value
    }
    return -1
}

func (this *LRUCache) Put(key int, value int) {
    if e, ok := this.mp[key]; ok {
        this.list.MoveToFront(e)
        e.Value.(*element).value = value
        return
    }
    if this.list.Len() >= this.capacity {
        tail := this.list.Back()
        key := tail.Value.(*element).key
        delete(this.mp, key)
        this.list.Remove(tail)
    }
    this.list.PushFront(&element{
        key:   key,
        value: value,
    })
    this.mp[key] = this.list.Front()
}
```


[^1]: [146. LRU Cache](https://leetcode.com/problems/lru-cache/)  