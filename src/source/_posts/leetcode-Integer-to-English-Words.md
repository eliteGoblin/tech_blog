---
title: 'leetcode: Integer to English Words'
date: 2019-08-04 19:35:48
tags: [LeetCode, amazon, trivial]
categories:
  - LeetCode
keywords:
description:
---

## Problem

Integer to English Words[^1], Convert a non-negative integer to its english words representation. Given input is guaranteed to be less than 231 - 1.Example 1:

Input: 123
Output: "One Hundred Twenty Three"

## 分析

总体思路是依照英文读数字的习惯，每三位分析：超过百位的数字加上特定单位，thousand, million, billion。32位整数最多不超过billion。实现过程善用map可以大大简化代码，另外注意0及trim多余的空格。


## Solution

```golang
func numberToWords(num int) string {
    res := convertHundred(num % 1000)
    v := []string{"Thousand", "Million", "Billion"}
    for i := 0; i < 3; i++ {
        num /= 1000
        if num%1000 > 0 {
            res = convertHundred(num%1000) + " " + v[i] + " " + res
        }
    }
    if res == "" {
        return "Zero"
    }
    return strings.TrimSpace(res)
}

func convertHundred(num int) string {
    under20 := []string{
        "", "One", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten", "Eleven", "Twelve", "Thirteen", "Fourteen", "Fifteen", "Sixteen", "Seventeen", "Eighteen", "Nineteen",
    }
    decimal := []string{
        "", "", "Twenty", "Thirty", "Forty", "Fifty", "Sixty", "Seventy", "Eighty", "Ninety",
    }
    res := ""
    a, b, c := num/100, num%100, num%10
    if b < 20 {
        res = under20[b]
    } else {
        res = decimal[b/10]
        if c > 0 {
            res += " " + under20[c]
        }
    }
    if a > 0 {
        res = under20[a] + " Hundred " + res
    }
    return strings.TrimSpace(res)
}
```

[^1]: [Integer to English Words](https://leetcode.com/problems/integer-to-english-words/)