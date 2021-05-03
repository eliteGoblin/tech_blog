---
title: 'leetcode: Reverse Integer'
date: 2019-08-11 13:07:58
hidden: true
tags: [LeetCode, amazon, trivial]
categories:
  - LeetCode
keywords:
description:
---

## Problem

Given a 32-bit signed integer, reverse digits of an integer.

Example 1:

Input: 123
Output: 321

## 分析

逆转整数，基本思路是不断得对10取余数，取得最后一个digit，然后再将 res * 10 + digit，得到reverse后的数。但tricky部分是会溢出，一个简单的方法是将其放入int64，逆转，再转换类型即可。

## Solution

```golang
const (
    uint32Max = ^uint32(0)
    int32Max  = int32(uint32Max >> 1)
    int32Min  = -int32Max - 1
)
func reverse(x int) int {
    var res int64
    for x != 0 {
        res *= 10
        res += int64(x) % 10
        x /= 10
    }
    if res > int64(int32Max) || res < int64(int32Min){
        return 0
    }
    return int(res)
}
``` 