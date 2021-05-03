---
title: 'leetcode: Search in Rotated Sorted Array'
date: 2019-08-17 09:30:04
tags: [LeetCode, amazon, array]
categories:
  - LeetCode
keywords:
description:
---

## Problem

Suppose an array sorted in ascending order is rotated at some pivot unknown to you beforehand.

(i.e., [0,1,2,4,5,6,7] might become [4,5,6,7,0,1,2]).

You are given a target value to search. If found in the array return its index, otherwise return -1.

You may assume no duplicate exists in the array.

Your algorithm's runtime complexity must be in the order of O(log n).

Example 1:

Input: nums = [4,5,6,7,0,1,2], target = 0
Output: 4
Example 2:

Input: nums = [4,5,6,7,0,1,2], target = 3
Output: -1

## 分析

很经典的题，依稀还记得若干年前刚接触此题的困惑，好好的排序数组干嘛非得给rotate一下再搜索，真是闲的。但本题确实比较锻炼递归思维。

题目给定经过一次旋转的数组(可能旋转完还是本身)，所谓旋转就像拿刀将数组一分为二，然后两段的整体顺序反一下。