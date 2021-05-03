---
title: 'leetcode: Median of Two Sorted Arrays'
date: 2019-07-12 11:55:29
hidden: true
tags: [LeetCode, amazon, array]
categories:
  - LeetCode
keywords:
description:
---

## 题目

Median of Two Sorted Arrays[^1], There are two sorted arrays nums1 and nums2 of size m and n respectively.

Find the median of the two sorted arrays. The overall run time complexity should be O(log (m+n)).

You may assume nums1 and nums2 cannot be both empty.


## 分析

Hard题目，初看挺懵的，感觉无从下手，而且看了解答也容易久久想不明白。其实不必畏惧，这道题一旦搞懂还是思路还是很清晰的，我们来分析一波：

求两sorted数组的median，最直观的方法便是先merge，然后取median即可，median公式:

$$median = \frac{arr[\frac{len+1}{2} - 1] + arr[\frac{len+2}{2} - 1]}{2.0}$$

也即第(len+1)/2和第(len+2)/2的元素求均值，元素个数为奇数时，其实是一个元素。

这样做的复杂度是O(n)，但要求复杂度O($log(m+n)$)，意味着我们可能需要减治法(reduce-and-conquer)[^2]，即每次丢弃数组的一部分，再剩余数组找答案。如何做呢？

## 减治法

求median可以用上述公式泛化为求两排序数组的Kth元素，我们的问题转化为如何找两排序数组的第K大元素。

考虑一个sorted数组，求第K个元素，我们可以扔掉前$\frac{K}{2}$个，然后再在剩下元素中找第$K-\frac{K}{2}$元素，这就是本题的基本思路：每次扔掉肯定排第K元素前面的$\frac{K}{2}$个元素，reduce问题规模，然后再在剩余集合中conquer，找第$K-\frac{K}{2}$个元素。

我们首先找出两数组A,B的第$\frac{K}{2}$元素，比较，假设A的第$\frac{K}{2}$元素<B的第$\frac{K}{2}$元素，则A的前$\frac{K}{2}$个元素不存在比整体第K个元素大的元素，可以放心的扔掉，比如:

A为[1, 2, 3], B为[2, 3, 4], K=5, K/2=2，A[2] = 2 < B[2] = 3，则A的前2个元素肯定不包含K，可以扔掉，于是在A[3] B[2, 3, 4]中找第3大的元素即可。

考虑特殊情况：

某个元素不够$\frac{K}{2}$该如何呢？这时需要扔掉*另一个*数组的前$\frac{K}{2}$个元素，因为元素多的数组的前$\frac{K}{2}$个元素肯定在最终的前K个元素中，比如:

A为[-999]，B为[2, 3, 4, 5]，K=４，K/2=2, 则B中前两个肯定组成了最终的前4个数，扔掉的B中前两个，则问题变为在A[-999]，B[4, 5]中找第2个元素。而且A若为[999]或者包含一个任意值的数组也一样。

另外若一个数组为空，或K=1的这两种特殊情况处理很直观。

## Solution

对应以上分析，code如下:

```golang
func findMedianSortedArrays(nums1 []int, nums2 []int) float64 {
    midLeft := (len(nums1) + len(nums2) + 1) / 2
    midRight := (len(nums1) + len(nums2) + 2) / 2
    return float64(findKth(nums1, nums2, midLeft)+
        findKth(nums1, nums2, midRight)) / 2.0
}

func findKth(nums1, nums2 []int, k int) int {
    var s, l []int
    if len(nums1) <= len(nums2) {
        s, l = nums1, nums2
    } else {
        s, l = nums2, nums1
    }
    if len(s) == 0 {
        return l[k-1]
    }
    if k == 1 {
        if s[0] <= l[0] {
            return s[0]
        }
        return l[0]
    }
    if len(s) < k/2 {
        return findKth(s, l[k/2:], k-k/2)
    }
    if s[k/2-1] <= l[k/2-1] {
        return findKth(s[k/2:], l, k-k/2)
    }
    return findKth(s, l[k/2:], k-k/2)
}
```


[^1]: [Median of Two Sorted Arrays](https://leetcode.com/problems/median-of-two-sorted-arrays/)
[^2]: [拜托，面试别再问我TopK了！！！](https://yq.aliyun.com/articles/642891)