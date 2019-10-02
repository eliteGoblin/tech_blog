---
title: 'leetcode: Two Sum'
date: 2019-07-04 11:19:57
tags: [LeetCode, amazon, misc]
categories:
  - LeetCode
keywords:
description:
---

## Preface

Leetcode的题目断断续续也刷了不少，特别是之前刷了遍soulmachine的leetcode题集[^1]，对题目的分类及解题套路有了一定了解。时间拖的比较久，很多题目现在再看也比较模糊，为了更好的总结和理解，再次记录自己每题的最佳解法，同时产出博客可是对漫漫刷题路的一个激励，借用看到的一句话：

> LC题和刷GRE 填空题一样，都是需要快速得到正反馈的事情。如果长久无法得到反馈，刷题人的情绪就容易变得糟糕[^2]

刷题的目的呢，当然首先是因为自己浓浓的大厂情节啦，其实长远看更重要的是确实从刷题中学到不少，思路也开阔，碰到复杂问题不慌了，代码质量也有肉眼可见的提高。

推荐一些资源: 

*  [Leetcode题目按公司，频率分类](http://206.81.6.248:12306/leetcode/Amazon/algorithm)
*  [LeetCode All in One 题目讲解汇总](https://www.cnblogs.com/grandyang/p/4606334.html) 

所有刷题的代码将会记录在[这个代码仓库](https://github.com/eliteGoblin/sky_ladder)里，现在还在尝试是每刷一道题都记录还是review时候再记录，边刷边感觉吧!

<!-- more -->

## Problem

Two Sum[^3], Given an array of integers, return indices of the two numbers such that they add up to a specific target.

You may assume that each input would have exactly one solution, and you may not use the same element twice.

Example:

Given nums = [2, 7, 11, 15], target = 9,

Because nums[0] + nums[1] = 2 + 7 = 9,
return [0, 1].

## Solution

这道题给定一个数组，要求给出一个数target，输出两个index，使两数之和等于target。  

首先想到暴力，固定一个index，遍历数组其余部分，查看两数之和是否为target，复杂度$O(n^2)$，应该会超时吧。。

而且也不能排序，因为最后要求返回index，一个简洁的办法是： 遍历原数组，得到value-->index的hash，然后遍历数组，只要看target-nums[i]是否在hash中存在就可以啦，注意的点是考虑特殊元素：value为target的一半，可能重复，如何区分是和为两倍的自己还是两个重复元素相加为target的情况：

*  hash中存放的是最后一个value的index
*  遍历数组由前向后遍历，需要判断target-nums[i]的index是否是自己，如果不是返回结果
*  巧妙之处：hash中存放的是最后一个为此value的index，由前向后遍历数组，才能判断是否是两个元素而不是一个元素的两倍

```golang
func twoSum(nums []int, target int) []int {
    mp := make(map[int]int)
    for i, v := range nums {
        mp[v] = i
    }
    for i, v := range nums {
        other := target - v
        if j, ok := mp[other]; ok {
            if i != j {
                return []int{i, j}
            }
        }
    }
    return []int{-1, -1}
}
```

[^1]: [Leetcode题集](https://github.com/soulmachine/leetcode)  
[^2]: [LeetCode-Python 270题刷题心得体会](https://blog.csdn.net/qq_32424059/article/details/89000776)  
[^3]: [1. Two Sum](https://leetcode.com/problems/two-sum/)