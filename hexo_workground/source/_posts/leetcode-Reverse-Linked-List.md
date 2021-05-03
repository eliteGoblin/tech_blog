---
title: 'leetcode: Reverse Linked List'
date: 2019-08-04 20:19:53
hidden: true
tags: [LeetCode, amazon, list]
categories:
  - LeetCode
keywords:
description:
---

## Problem

Reverse Linked List[^1], Reverse a singly linked list.

Example:

Input: 1->2->3->4->5->NULL
Output: 5->4->3->2->1->NULL

## 分析

又一道list基础题，可以分别用循环和递归解。

循环：create pseudoHead节点，遍历list,每次将node插入pseudoHead的下一个节点，这样最后得到的就是reverse的list。

递归：假设链表 A->B->C，处理A时，先逆转后面链表，得到C->B，但A->next中存储B的指针，因此可以将A挂在B后面，A->next设置为nil，返回C，得到逆转的list。

## Solution

循环
```golang
func reverseListIterative(head *ListNode) *ListNode {
    pseudoHead := &ListNode{}
    var next *ListNode
    for p := head; p != nil; {
        next = p.Next
        p.Next = pseudoHead.Next
        pseudoHead.Next = p
        p = next
    }
    return pseudoHead.Next
}
```

递归
```golang
func reverseListRecursive(head *ListNode) *ListNode {
    if head == nil || head.Next == nil {
        return head
    }
    newHead := reverseListRecursive(head.Next)
    head.Next.Next = head
    head.Next = nil
    return newHead
}
```

[^1]: [Reverse Linked List](https://leetcode.com/problems/reverse-linked-list/)