---
title: 'leetcode: Best Time to Buy and Sell Stock'
date: 2019-08-11 12:42:40
tags: [LeetCode, amazon, trivial]
categories:
  - LeetCode
keywords:
description:
---

## 题目

Say you have an array for which the ith element is the price of a given stock on day i.

If you were only permitted to complete at most one transaction (i.e., buy one and sell one share of the stock), design an algorithm to find the maximum profit.

Note that you cannot sell a stock before you buy one.

Example 1:

Input: [7,1,5,3,6,4]
Output: 5
Explanation: Buy on day 2 (price = 1) and sell on day 5 (price = 6), profit = 6-1 = 5.　Not 7-1 = 6, as selling price needs to be larger than buying price.

## 分析

给定array代表每天股票价格，而且只能限定买卖一次，求最大收益。直接的思路就是固定一天为买入，遍历余下日子，卖出，求出最大收益。复杂度$$O(n^2)$$

更好的方法是: 每次遍历位置作为sell day而不是buy day，同时在遍历过程中保存当前位置的之前最小value min, 当前股价-min就是收益，记录最大值即可。

## Solution

```golang
func maxProfit(prices []int) int {
    if len(prices) == 0 {
        return 0
    }
    m := 0
    preMin := prices[0]
    for i := 1; i < len(prices); i++ {
        m = max(m, prices[i]-preMin)
        preMin = min(prices[i], preMin)
    }
    return m
}
```