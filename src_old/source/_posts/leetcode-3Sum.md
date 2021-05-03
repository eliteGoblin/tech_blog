---
title: 'leetcode: 3Sum'
date: 2019-07-27 10:03:45
tags: [LeetCode, amazon, array]
categories:
  - LeetCode
keywords:
description:
---

## 题目

3Sum[^1], Given an array nums of n integers, are there elements a, b, c in nums such that a + b + c = 0? Find all unique triplets in the array which gives the sum of zero.

Note:

The solution set must not contain duplicate triplets.

## 分析

本题要求返回value的triplets，并不是index，为避免$O(n^3)$的复杂度，可以先排序。

遍历数组，固定第一个数nums[i]，相当于在剩余数组找符合2Sum的问题，即符合$$nums[j] + nums[k] = targer-nums[i]$$的pair即可，用双指针left, right，从两边向中间扫描:，

*  2sum < target，left ++
*  2sum > target, right --
*  2sum == target: left ++, right --

即可找到所有满足条件的pair。


同时注意到target为0，因此只要发现nums[i] > 0即可以退出循环，因为nums[j], nums[k]均>nums[i]>0，因此肯定不满足相加为0的条件。

题目要求不能出现重复的triplets，其实这个本题的一个难点，我们可能会考虑用array: [３]int作为map key来去重(注意golang slice并不能作为map key，因为没有定义slice的equal)，但是[]int包含相同元素但乱序也视为相等，这条思路不通，因此只能我们遍历求结果时，自动跳过重复的三元组，不将其插入到最终结果中即可。

那如何做呢？首先考虑什么情况下会产生重复结果，以下讨论的前提是已经对输入数组排序

### 第一个元素重复

考虑{-1, -1, 0, 1}: 

如果两次都将-1固定为首元素，那会将两个{-1, 0, 1}加入结果集。我们只需要在固定第一个元素时，检查前一个元素是否和当前第一个元素是否相同，相同跳过即可。

### 第二个元素重复

考虑{-1, 0, 0, 1}，同样我们在选第二个元素时，检查是否重复，我们需不需要担心第三个元素重复呢？不需要，前两个元素都没有重复的答案，第三个元素相等代表两组triplet的和不同。

## Solution


```golang
func threeSum(nums []int) [][]int {
    var res [][]int
    if len(nums) < 3 {
        return res
    }
    sort.Ints(nums)
    for i := 0; i < len(nums)-2; i++ {
        if nums[i] > 0 {
            break
        }
        if i > 0 && nums[i] == nums[i-1] {
            continue
        }
        target := -nums[i]
        j, k := i+1, len(nums)-1
        for j < k {
            if j > i+1 && nums[j] == nums[j-1] {
                j++
                continue
            }
            v := nums[j] + nums[k]
            if v > target {
                k--
            } else if v < target {
                j++
            } else {
                res = append(res, []int{nums[i], nums[j], nums[k]})
                j++
                k--
            }
        }
    }
    return res
}
```


[^1]: [3Sum](https://leetcode.com/problems/3sum/)