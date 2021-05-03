---
title: 'leetcode: Product of Array Except Self'
date: 2019-08-17 09:14:33
tags: [LeetCode, amazon, array]
categories:
  - LeetCode
keywords:
description:
---

## Problem

Given an array nums of n integers where n > 1,  return an array output such that output[i] is equal to the product of all the elements of nums except nums[i].

Example:

Input:  [1,2,3,4]
Output: [24,12,8,6]

## 分析

这道题要求我们不能用除法，既然是求不包含自身的所有元素的积，那试着用乘法：元素左边所有元素之积乘以元素右边所有元素之积就是不含此元素的数组积。方法就是先从左至右遍历一遍数组，求出每个元素左边积组成的数组。同理从右至左遍历数组，求出每个元素右侧积的数组，再将两数组每个index相乘，即得到答案。

过程可以简化，我们不需要额外分配数组，reuse result数组：第一遍将左侧积数组存在res，第二次从右至左遍历时，只用一个变量存储当前右边积即可，同时更新res数组

## Solution

```golang
func productExceptSelf(nums []int) []int {
    res := make([]int, len(nums))
    if len(nums) == 0 {
        return nums
    }
    res[0] = 1
    for i := 1; i < len(nums); i++ {
        res[i] = nums[i-1] * res[i-1]
    }
    product := 1
    for i := len(nums) - 1; i >= 0; i-- {
        res[i] = res[i] * product
        product *= nums[i]
    }
    return res
}
```