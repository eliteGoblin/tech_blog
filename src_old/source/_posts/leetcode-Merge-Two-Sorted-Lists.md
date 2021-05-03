---
title: 'leetcode: Merge Two Sorted Lists'
date: 2019-07-28 18:26:02
tags: [LeetCode, amazon, list]
categories:
  - LeetCode
keywords:
description:
---

## Problem

Merge Two Sorted Lists[^1], Merge two sorted linked lists and return it as a new list. The new list should be made by splicing together the nodes of the first two lists.

## 分析

基础题，相加两链表代表的数字，两链表已经逆序，所以直接从head加起即可，用carry记录当前是否有进位，注意：

*  判断一个已经到达尾部情况
*  两链表加完，仍有carry情况
*  list问题很多用pseudoHead能简化逻辑

## Solution

```golang
func mergeTwoLists(l1 *ListNode, l2 *ListNode) *ListNode {
    if l1 == nil {
        return l2
    }
    if l2 == nil {
        return l1
    }
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


[^1]: [Merge Two Sorted Lists](https://leetcode.com/problems/merge-two-sorted-lists/)