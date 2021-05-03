---
title: 'leetcode: Add two numbers'
date: 2019-07-08 12:59:42
tags: [LeetCode, amazon, list]
categories:
  - LeetCode
keywords:
description:
---

## Problem

Add two numbers[^1], You are given two non-empty linked lists representing two non-negative integers. The digits are stored in reverse order and each of their nodes contain a single digit. Add the two numbers and return it as a linked list.

You may assume the two numbers do not contain any leading zero, except the number 0 itself.

<!-- more -->

Example:

Input: (2 -> 4 -> 3) + (5 -> 6 -> 4)
Output: 7 -> 0 -> 8
Explanation: 342 + 465 = 807.

## Solution

经典的两数字相加链表版，需要注意的: 

*  两链表为空
*  两链表不等长，剩余链表和进位相加情况
*  链表都走完，仍剩余进位的情况

```golang
func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
    if l1 == nil {
        return l2
    } else if l2 == nil {
        return l1
    }
    carry := 0
    head := &ListNode{
        Next: nil,
    }
    pre := head
    for l1 != nil && l2 != nil {
        value := l1.Val + l2.Val + carry
        pre.Next = &ListNode{
            Val: value % 10,
        }
        pre = pre.Next
        carry = value / 10
        l1 = l1.Next
        l2 = l2.Next
    }
    left := l1
    if left == nil {
        left = l2
    }
    for left != nil {
        value := left.Val + carry
        pre.Next = &ListNode{
            Val: value % 10,
        }
        pre = pre.Next
        carry = value / 10
        left = left.Next
    }
    if carry > 0 {
        pre.Next = &ListNode{
            Val: carry,
        }
    }
    return head.Next
}
```


[^1]: [2. Add Two Numbers](https://leetcode.com/problems/add-two-numbers/)  