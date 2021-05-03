---
title: 'leetcode: Maximum Subarray'
date: 2019-07-28 18:34:03
hidden: true
tags: [LeetCode, amazon, array]
categories:
  - LeetCode
keywords:
description:
---

## Problem

Maximum Subarray[^1], Given an integer array nums, find the contiguous subarray (containing at least one number) which has the largest sum and return its sum.

## 分析

比如给定例子: {-2, 1, 3} 要求连续数组的最大和，如何求呢？

先遍历数组，似乎可以先丢弃负数，直接寻找相加得正数的子数组，但针对{2, -1, 3}，因为a[0]+a[1]大于0，三个数构成最大和，而非仅3或2。进一步分析，每次遍历元素时，目的判断当前元素是否与之前元素形成连续子数组还是应该新起子数组：

*  当前元素与之前子数组和 < 0，则说明对后面和积累不利，负数不应该继续累加到后面元素，只会要子数组和最小，则中断之前的子数组，新起数组
*  当前元素与之前子数组 > 0，则正数可以对后面的累加有正的贡献，当前元素采纳为子数组，继续之前子数组

因为题意要求至少选择一个，如果全部为负数的数组，结果为绝对值最小元素，为简化逻辑，每次比较：

*  若当前元素+前面子数组和 < 当前元素，重新开始子数组
*  否则之前子数组累加当前元素值，视当前元素为其一部分

## Solution

```golang
const (
    uintMax = ^uint(0)
    intMax  = int(uintMax >> 1)
    intMin  = -intMax - 1
)

func maxSubArray(nums []int) int {
    res := intMin
    curSum := 0
    for _, v := range nums {
        curSum = max(curSum+v, v)
        res = max(res, curSum)
    }
    return res
}
```


[^1]: [Maximum Subarray](https://leetcode.com/problems/maximum-subarray/)