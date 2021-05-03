---
title: 'leetcode: Number of Islands'
date: 2019-07-27 08:54:55
tags: [LeetCode, amazon, dfs]
categories:
  - LeetCode
keywords:
description:
---

## 题目

Number of Islands[^1], Given a 2d grid map of '1's (land) and '0's (water), count the number of islands. An island is surrounded by water and is formed by connecting adjacent lands horizontally or vertically. You may assume all four edges of the grid are all surrounded by water.

## 分析

这道题给定二维数组代表地图，每个格子1是岛屿，0是水，求解岛屿个数；同时指明vert, horizonal可达的才算属于同一个岛屿。

要求岛屿个数，本质是搜索问题：以每个格子为起点，向四周(依题意上下左右四个方向)进行扩散搜索：如果碰到岛屿，则标记为已经搜索过，同时以新岛屿继续扩散式搜索，这样能把此格子所在岛屿全部遍历到。下次遍历如果发现当前格子没有遍历过，则count++，若已经遍历过，跳过。这里用DFS和BFS都可以。

## BFS 解法

```golang
func numIslands(grid [][]byte) int {
    count := 0
    for i := 0; i < len(grid); i++ {
        for j := 0; j < len(grid[0]); j++ {
            if grid[i][j] == '1' {
                count++
                markIslandContainLocation(grid, i, j)
            }
        }
    }
    return count
}

func markIslandContainLocation(grid [][]byte, i, j int) {
    if i < 0 || i >= len(grid) || j < 0 || j >= len(grid[0]) ||
        grid[i][j] != '1' {
        return
    }
    grid[i][j] = '2'
    markIslandContainLocation(grid, i-1, j)
    markIslandContainLocation(grid, i+1, j)
    markIslandContainLocation(grid, i, j-1)
    markIslandContainLocation(grid, i, j+1)
}
```


[^1]: [Number of Islands](https://leetcode.com/problems/number-of-islands/)