---
title: 'leetcode: Longest Palindromic Substring'
date: 2019-07-10 13:11:26
hidden: true
tags: [LeetCode, amazon, array]
categories:
  - LeetCode
keywords:
description:
---

## Problem

Given a string s, find the longest palindromic substring in s. You may assume that the maximum length of s is 1000.

Example 1:

Input: "babad"
Output: "bab"
Note: "aba" is also a valid answer.

<!-- more -->

## Analyze

题意是给定字串，找出其中最长的palindrome，即最长左右对称的字串。回忆回文定义：

*  长度为奇数时，从中间字符向两边走，两两对称
*  长度为偶数时，从中间第len/2, len / 2 + 1元素两两相等，并向两边走，字符均两两相等

因此直观做法：遍历字符串，以当前位置为中心(奇数情况）或以当前位置和下一位置为中心(偶数情况)，向两边试探，记录最长的字串

## Solution

以上分析，奇数和偶数试探时可以合并，奇数可以视为偶数的特例，传入的index为相同值，偶数两index相差1

```golang
func longestPalindrome(s string) string {
    if len(s) < 2 {
        return s
    }
    var start, end int
    for i := 0; i < len(s); i++ {
        curS, curE := findLongestPalindrome(s, i, i)
        if curE-curS > end-start {
            start, end = curS, curE
        }
        curS, curE = findLongestPalindrome(s, i, i+1)
        if curE-curS > end-start {
            start, end = curS, curE
        }
    }
    return s[start:end]
}

func findLongestPalindrome(s string, left, right int) (start, end int) {
    for left >= 0 && right < len(s) && s[left] == s[right] {
        left--
        right++
    }
    return left + 1, right
}
```

