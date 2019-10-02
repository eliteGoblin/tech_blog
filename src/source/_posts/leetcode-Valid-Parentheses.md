---
title: 'leetcode: Valid Parentheses'
date: 2019-08-04 20:03:53
tags: [LeetCode, amazon, misc]
categories:
  - LeetCode
keywords:
description:
---

## Problem

Valid Parentheses[^1], iven a string containing just the characters '(', ')', '{', '}', '[' and ']', determine if the input string is valid.

An input string is valid if:

Open brackets must be closed by the same type of brackets.
Open brackets must be closed in the correct order.
Note that an empty string is also considered valid.

## 分析

stack练习的基础题目，如果是只有一种括号，也可以不用显式定义stack,用递归求解：用left, right分别记录剩余左右括号数目，如果left==right==0，则说明是valid，left>right，非法字符，其余碰到match情况让对应括号计数值--即可。==0则为已经用完，这样不必再传入括号个数count。



## Solution

```golang
var mp = map[byte]byte{
    '}': '{',
    ']': '[',
    ')': '(',
}

func isValid(s string) bool {
    stk := NewStack(len(s))
    bytes := []byte(s)
    for _, v := range bytes {
        switch v {
        case '{':
            fallthrough
        case '[':
            fallthrough
        case '(':
            stk.Push(v)
        default:
            top, err := stk.Top()
            if err != nil ||
                mp[v] != top.(byte) {
                return false
            }
            stk.Pop()
        }
    }
    return stk.Len() == 0
}
```

golang没有内置stack，可以用slice实现
```golang
type stack []interface{}
func NewStack(cap int) *stack {
    s := make(stack, 0, cap)
    return &s
}

func (s *stack) Push(v interface{}) {
    *s = append(*s, v)
}

func (s *stack) Pop() error {
    if len(*s) == 0 {
        return errors.New("underflow")
    }
    *s = (*s)[:len(*s)-1]
    return nil
}

func (s *stack) Top() (interface{}, error) {
    if len(*s) == 0 {
        return struct{}{}, errors.New("underflow")
    }
    return (*s)[len(*s)-1], nil
}

func (s *stack) Len() int {
    return len(*s)
}
```



[^1]: [Valid Parentheses](https://leetcode.com/problems/valid-parentheses/)