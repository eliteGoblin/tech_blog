---
title: 'leetcode: Merge k Sorted Lists'
date: 2019-08-04 12:11:50
tags: [LeetCode, amazon, list]
categories:
  - LeetCode
keywords:
description:
---

## Problem

Merge k Sorted Lists[^1], Merge k sorted linked lists and return it as one sorted list. Analyze and describe its complexity.

Example:

Input:
[
  1->4->5,
  1->3->4,
  2->6
]
Output: 1->1->2->3->4->4->5->6

## 分析

是之前merge sorted list的延伸，之前的题目要求merge两个list,现在k个可以用k-1次两两merge来实现，有没有更快的算法呢？可以用merge sort的分治思路，要merge K个list,那将K个list分成两组，分别递归merge2个K/2个的list, 再将两者结果做一次merge。

为什么这样的复杂度低呢？以4个sorted list为例，每个list元素个数为n。假设两两merge，过程如下:

1. merge l1, l2; 遍历2n个node，得到l12
2. merge结果和l3, 遍历3n个node，得到l123
3. merge结果和l4, 遍历4n个node

共计遍历2n+3n+4n=10n个node

而用分治法：

1. merge l1, l2得到l12;merge l3, l4得到l34；共计遍历4n个node
2. merge l12和l34，共计遍历4n个node

分治法总计: 4n + 4n = 8n

## Solution

```golang
func mergeKLists(lists []*ListNode) *ListNode {
    if len(lists) == 0 {
        return nil
    }
    if len(lists) == 1 {
        return lists[0]
    }
    mid := len(lists) / 2
    l1 := mergeKLists(lists[:mid])
    l2 := mergeKLists(lists[mid:])
    pseudoHead := &ListNode{}
    pre := pseudoHead
    for l1 != nil && l2 != nil {
        if l1.Val <= l2.Val {
            pre.Next = l1
            pre = l1
            l1 = l1.Next
        } else {
            pre.Next = l2
            pre = l2
            l2 = l2.Next
        }
    }
    if l1 != nil {
        pre.Next = l1
    } else {
        pre.Next = l2
    }
    return pseudoHead.Next
}
```

[^1]: [Merge k Sorted Lists](https://leetcode.com/problems/merge-k-sorted-lists/)