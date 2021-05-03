---
title: 'leetcode: Trapping Rain Water'
date: 2019-07-27 09:32:08
hidden: true
tags: [LeetCode, amazon, dfs]
categories:
  - LeetCode
keywords:
description:
---

## 题目

Trapping Rain Water[^1], Given n non-negative integers representing an elevation map where the width of each bar is 1, compute how much water it is able to trap after raining.

<div align="center">
{% asset_img rainwatertrap.png %}
</div>

## 分析

挺有意思的一道题，给定一个数组代表每个位置的高度，求能容纳多少雨水，看着很高级的样子，其实需要我们找出容器中所有的"坑"，并依次计算每个坑的面积，累加即可。这个思路转换是本题的关键，我们只用关心当前位置是否属于坑的一部分，当前位置会对大坑贡献多少面积即可，这样简化了问题，而不必陷入复杂思路：先描绘出不规则图形，再想办法计算此图形面积。

直观的思路就是对于每个位置，判断是否是"坑"：即分别向左右扫描，找出左右最高的点，算作坑的边缘，取两端较小值作为坑壁高，记为H，当前位置能容纳的水面积为 H - height[i]，累加所有位置存水量即可得到总储水量。

## Solution

可以从扫描两遍，得到两个数组：maxHeightFromLeft, maxHeightFromRight，分别记录对当前位置来说，左右边缘，然后再遍历一遍累加存水量。共三遍，而且需要两个$O(n)$空间，小技巧可以简化：

*  只需要额外的一个数组记录maxHeightFromLeft，第二次遍历时用一个变量记录来自右边的最大值，就可以计算出当前位置的存水量
*  同样在右侧开始遍历时，在得出右侧边缘高度后，就开始计算存水量：即两边缘的差值，省去最后一次遍历


```golang
func trap(height []int) int {
    res := 0
    dp := make([]int, len(height))
    maxValue := 0
    for i := 0; i < len(dp);i ++ {
        dp[i] = maxValue
        maxValue = max(maxValue, height[i])
    }
    maxValue = 0
    for i := len(dp) - 1; i >= 0; i -- {
        dp[i] = min(dp[i], maxValue)
        maxValue = max(maxValue, height[i])
        if dp[i] > height[i] {
            res += dp[i] - height[i]
        }
    }
    return res
}
```

[^1]: [Trapping Rain Water](https://leetcode.com/problems/trapping-rain-water/)