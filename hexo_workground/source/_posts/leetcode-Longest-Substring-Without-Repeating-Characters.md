---
title: 'leetcode: Longest Substring Without Repeating Characters'
date: 2019-07-27 10:55:21
hidden: true
tags: [LeetCode, amazon, array]
categories:
  - LeetCode
keywords:
description:
---


## Problem

Longest Substring Without Repeating Characters[^1], Given a string, find the length of the longest substring without repeating characters.

Example 1:

Input: "abcabcbb"
Output: 3 
Explanation: The answer is "abc", with the length of 3. 

## 分析

要求string不重复字串的长度，可以用map来记录字符出现的位置，但是重复元素的位置如何记录呢？难不成用一个类似multimap的东西么？并不需要，我们用一个prePos来记录本次查询起点前一个位置，这样用i - prePos就得到本次查询子串的长度。

这样我们遍历字串时，首先查看当前字符是否存在map，如果没有，插入map；如果存在，并不一定就是当前字串含重复元素；将其与prePos比较，如果pos <= prePos，说明不与本子串重复；如果pos > prePos，说明检测到重复，这时记录子串长度，并更新prePos与map即可。


## Solution


```golang
func lengthOfLongestSubstring(s string) int {
    sB := []byte(s)
    mp := make(map[byte]int)
    prePos := -1
    var res int
    for i := range s {
        if pos, ok := mp[sB[i]]; ok {
            if pos > prePos {
                prePos = pos
            }
        }
        mp[sB[i]] = i
        res = max(res, i-prePos)
    }
    return res
}
```


[^1]: [Longest Substring Without Repeating Characters](https://leetcode.com/problems/longest-substring-without-repeating-characters/)