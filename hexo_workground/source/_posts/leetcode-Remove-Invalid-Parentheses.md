---
title: 'leetcode: Remove Invalid Parentheses'
date: 2019-08-11 12:15:22
hidden: true
tags: [LeetCode, amazon, string, bfs]
categories:
  - LeetCode
keywords:
description:
---

## 题目

Remove the minimum number of invalid parentheses in order to make the input string valid. Return all possible results.

Note: The input string may contain letters other than the parentheses ( and ).

Example 1:

Input: "()())()"
Output: ["()()()", "(())()"]

## 分析

题目给定含有一系列左右括号的字串，可以删除字符，要求给出remove最少的方案，相同remove个数都返回。

递归，即DFS是一个比较自然的思路，但本题要求删除最少，DFS适合求解有无此方案，不适合求最少。可行思路可以是：每次遍历所有位置，remove一个，看是否合法；如果找到，则停止搜索，返回答案。如果找不到，则以每次remove一个字符的到的字串为基础，再remove一个，直到找到。相当于固定了字串的长度，按长度递减的顺序搜索。

因此可以用BFS解，queue中存储待查找的字串，遍历完当前queue后，如果全部invalid，则将remove一个字符的所有字符存入当前queue。举例: 

1.  队列当前有长度为3的两个字串: A: )() B: ()( 
2.  A为invalid, remove一个字符，并存入队列, 队列变为 B:()( A1: () A2: )) A3: )(

当发现valid string后，停止erase。判断为当前queue的所有剩余string即可。但是当前queue中存在长度为n和n-1两种字串，后者是一些字串erase一个元素，会不会错误的把n-1长度的合法字串也存入结果呢？不会的，合法字串的左右括号和必须为偶数，两字串长度相差为１并不会都是valid。

## Solution

```
func removeInvalidParentheses(s string) []string {
    res := make([]string, 0)
    visited := make(map[string]bool)
    visited[s] = true
    lst := list.New()
    lst.PushBack(s)
    found := false

    for lst.Len() > 0 {
        e := lst.Front()
        lst.Remove(e)
        cur := e.Value.(string)
        if isValidParentheses(cur) {
            res = append(res, cur)
            found = true
        }
        if found {
            continue
        }
        for i := range cur {
            if cur[i] != '(' && cur[i] != ')' {
                continue
            }
            eraseOne := cur[:i] + cur[i+1:] // 之前写成 eraseOne := s[:i] + s[i+1:], 极难发现!!!
            if _, ok := visited[eraseOne]; !ok {
                lst.PushBack(eraseOne)
                visited[eraseOne] = true
            }
        }
    }
    if len(res) == 0 {
        res = []string{""}
    }
    return res
}

func isValidParentheses(s string) bool {
    count := 0
    for i := range s {
        if s[i] == '(' {
            count++
        } else if s[i] == ')' {
            count--
            if count < 0 {
                return false
            }
        }
    }
    return count == 0
}
```

