---
title: '一题两解:一道leetcode问题的DP和Greedy解法'
date: 2017-09-01 16:12:32
tags: [algorithms,leetcode]
---

<div align="center">
{% asset_img greedy.jpg %}
</div>

#### Preface

DP(Dynamic Programming)和Greedy算法，区别于具体的算法，是algorithm paradigm，可称为算法思想，被用来求Optimization类问题，这类问题用一般方法解决起来颇有难度，而且常常复杂度太高。DP和Greedy会让不熟悉的同学有所畏惧。

DP和Greedy都需要专题性的练习，本文由一道具体的题引出DP和Greedy解法，读者可以基于这个具体的问题来了解两种算法思想的应用，以及对比两者对于解决同一道题目的不同，加深对两者应用场景的认识

#### 为什么要学习DP和Greedy?

*  在技术面试中如果能展示DP和Greedy求最优解一般都会给自己带来加分，这也是本人学习DP和Greedy的直接动力。
*  很多常用算法都采用了DP和Greedy思想，参见[这里](#algo_dp_greedy_algoindex)，熟悉DP和Greedy思想更容易理解这类算法

<!-- more -->

#### 引子：一道leetcode题目

问题: 给出一个pair: (a,b), a < b的数组，如(1,2) (3, 4) (0,6) (5,7) (8,9) (5,9)

求pair构成的最长序列: p1, p2, p3... pi　
要求：　前后两个pair A，B: B.first > A.Second
本例的最长序列为：(1,2) (3, 4) (5,7) (8,9)  

原题在这里: [Maximum Length of Pair Chain](https://leetcode.com/problems/maximum-length-of-pair-chain/description/)

#### DP初步及解法

##### DP简介

What is DP and Why DP ?　　
DP是对递归解法的实现的优化，以经典的Fibonacci numbers问题为例，递归解法: 

```c
int fib(int n)
{
   if ( n <= 1 )
      return n;
   return fib(n-1) + fib(n-2);
}
```

递归解法包含很多次重复计算, 如fib(0)执行3次，fib(1) 5次

```
                         fib(5)
                     /             \
               fib(4)                fib(3)
             /      \                /     \
         fib(3)      fib(2)         fib(2)    fib(1)
        /     \        /    \       /    \
  fib(2)   fib(1)  fib(1) fib(0) fib(1) fib(0)
  /    \
fib(1) fib(0)
```

求fib(n)问题，随着n增大重复计算也会急剧增长，为了解决重复计算问题，很自然的思路是把中间计算的fib(m)存入表中，用查表来代替重复计算，如何查表有两种具体的technique:

*  执行Top-down recursive，需要fib(m)值时，先查表: 表中有就不必计算，这就是DP算法的Memoization解法

> memoization is an optimization technique used primarily to speed up computer programs by storing the results of expensive function calls and returning the cached result when the same inputs occur again -- Wikipedia

* 　第二种: 最终目的求fib(n)，因为fib(n) = fib(n-1) + fib(n-2)，需要先求较小的值。可以直接从最小值开始计算fib(0)，fib(1)... 最终得出fib(n)，这就是DP算法的tablular解法

**可见DP就是采用recursive思想，并利用Memoization或Tabular technique来优化求中间值过程，避免重复计算**

经典算法教程归纳的DP两个特征和上面提到的等价: 

*  Optimal Substructure：　问题最优解可籍由其子问题的最优解来获得，暗示了递归算法，找到了递归关系就基本解决了问题
*  Overlapping subproblems:　直接递归时子问题存在重叠，即存在重复计算问题，需要用Memoization或tabular technique来解决

##### DP解法

目的是找到与子问题的递归关系，关键在于子问题的界定：

排序的pair的更利于我们处理，可以按照first element排序，排序后的pair为:

```
(0,6) (1,2) (3, 4) (5,7) (5,9) (8,9) 
```

本题是求6个pair的最长连接序列，递归的想法便是: 已知k个pair的最长序列(0< k < 6)，如何求6个pair的最长序列？

目标问题用f(6)表示，解法为: 

> 用第6个pair(8,9)分别尝试于前5个子问题的最优解进行配接，即检查(8,9)和f(1) ... f(5)的最长序列是否能组成更长的序列，并取其中的最大值为f(6)的解，这就把f(6)成功的转化为对子问题的求解

公式

```
f(n) = max {
     v(i) = {
        f(i) + 1 // if pair(n)与f(i)的最长sequence可以组成更长的sequence
        f(i)     // otherwise
    }
} i = [1,n-1]
```

代码因而也比较直观：

```go
func findLongestChain(pairs [][]int) int {
    if len(pairs) == 0 {
        return 0
    }
    sortPairsArrByFirstElement(pairs)       // 根据first element来排序pairs数组
    dp := make([]int, len(pairs))           // 存放子问题的最优解
    ele2OfDPRes := make([]int, len(pairs))　　// 存放子问题的最大序列的second element,以判断某pair是否可以与子问题组成更长的序列
    dp[0] = 1                               // 初始化：　f(1)的最长sequence是自己，len为1
    ele2OfDPRes[0] = pairs[0][1]
    max := 1
    for i := 1; i < len(pairs); i++ {
        dp[i] = 1
        ele2OfDPRes[i] = pairs[i][1]
        for j := 0; j < i; j++ {
            len := 0
            if pairs[i][0] > ele2OfDPRes[j] { // 新pair可以与f(i)组成更长的sequence
                len = dp[j] + 1
                if len > dp[i] {
                    dp[i] = len
                    ele2OfDPRes[i] = pairs[i][1]
                }
            } else {                          // 新pair无法与f(i)组成更长的sequence，最长仍为f(i)的sequence
                len = dp[j]
                if len > dp[i] {
                    dp[i] = len
                    ele2OfDPRes[i] = pairs[j][1]
                }
            }
        }
        if dp[i] > max {
            max = dp[i]
        }
    }
    return max
}
```


#### Greedy初步及解法

##### Greedy初步

如果说DP本质是用递归求解问题，是Top-down的，Greedy的本质便是由小及大，每个子问题做出当前看起来*最好*的选择，到最后得到的便是问题最优解。

> Greedy is an algorithmic paradigm that builds up a solution piece by piece, always choosing the next piece that offers the most obvious and immediate benefit.

但对于某些问题，每步的局部最优并不意味着得到的solution是全局最优的，因此Greedy算法最核心的办法是能证明一个问题能用局部最优法一步步求出的结果也是全局最优，举个例子：

*  问题1: 小偷，来到一家商店，店里有金沙，铜沙，铜沙，但他的背包只能装50kg的物品，那么该如何选择使带走的东西价值最大呢？ 
    这个问题可以用Greedy算法解决：　首先装金沙，如果装完包内有空余，装银沙，否则装铜沙。即每步开始装当前单价最贵的物品，按此思路装满就是最优解。用反证法可以证明此解法正确。
*  问题２：小偷，来到另一家商店，这次店里是金砖，银砖，铜砖，同样背包限重50kg, 那么这次他还能用贪婪算法成功的带走最大价值呢？
    答案是不行的，一个反例：　金砖0.1kg, 银砖50kg, 铜砖10kg, 50kg银砖价值>0.1kg金砖，因此最优解是只装银砖，而非先装金砖

上面列举的问题是经典的背包问题，分别是fractional knapsack和0-1 knapsack，前者可以用Greedy算法解决

##### Greedy解法

本题的关键便是找到Greedy算法，并证明用Greedy算法求出的solution是最优的

还是先对pairs排序，这次以second element为key进行升序排列：

```
(1,2) (3, 4) (0,6) (5,7) (8,9) (5,9)
```

Greedy算法便是: 
*  第一个元素肯定入选最长sequence
*  依次按顺序遍历pair，每次选second element最小的，假设与之前得到的sequence匹配(此pair的first element > sequence最后pair的second element)，则此pair在sequence上，否则找下一个

证明: 反证法
*  第一条：假设第一个元素A不属于最长sequence，假设是pair B，但是有B.second > A.second，因此A也满足最长sequence
*  第二条：已经得到部分最长sequence，当前满足sequence的最小pair为A，假设最终sequence选B没选A，同上反证也可以证明A满足sequence，而且A的second element更小，比B更优。

证明了Greedy可行，实现很直观:

```go
func findLongestChain(pairs [][]int) int {
    if len(pairs) <= 1 {
        return len(pairs)
    }
    sortPairsArrBySecElement(pairs)
    rightMost := pairs[0][1]
    lenPairs := 1
    for i := 1; i < len(pairs); i++ {
        if pairs[i][0] > rightMost {
            lenPairs++
            rightMost = pairs[i][1]
        }
    }
    return lenPairs
}
```

可见，若Greedy算法可行，则Greedy一般而言是最优解法，时间，空间，代码实现复杂度都不大

#### 为什么Greedy和DP对pair数组排序方式不一样

其实是DP按first，second排序都可以，因为DP对元素是如何被排序的没有依赖，只要能够知道当前pair是否和某个特定子问题的sequence构成更长的sequence即可。
而Greedy只能按照second排序，因为如下按pair.first排序的例子：
```
(1, 10000) (2, 3) (5, 6)
```
很明显不能选(1, 10000)作为最长sequence节点

#### <h4 id=algo_dp_greedy_algoindex>DP和Greedy应用举例</h>

*  DP
    -  string algorithms
        +   longest common subsequence
        +   longest increasing subsequence
        +   longest common substring
        +   Levenshtein distance 
    -  0-1 knapsack
    -  Bellman–Ford algorithm
    -  Floyd's all-pairs shortest path algorithm
    -  travelling salesman problem
*  Greedy
    -  Kruskal’s Minimum Spanning Tree (MST)
    -  Prim’s Minimum Spanning Tree
    -  Dijkstra’s Shortest Path
    -  Huffman Coding
    -  Activity Selection problem(本题变种，基本一样)
    -  fractional knapsack

[较全的DP应用场景见这里](#https://en.wikipedia.org/wiki/Dynamic_programming#Algorithms_that_use_dynamic_programming#Algorithms_that_use_dynamic_programming)

#### 结语

本文是笔者学习DP和Greedy算法的一个引子，今后会安排时间依着Sedgewick的*Algorithms*一书以及geekforgeek的题集系列继续深入学习
在接近完成本文时，看到了知乎上的问题，提到DP的本质是**状态转移方程**，颇有道理，之后随着对DP的理解加深会持续更正本文，或新开文章讨论。

学习是一个循序渐进的过程。一点心得分享给大家：

> 克服了对*“难题”*的畏惧，也就解决了难题的大半

#### 参考文献

[什么是动态规划？动态规划的意义是什么？](https://www.zhihu.com/question/23995189)  
[Dynamic Programming | Set 20 (Maximum Length Chain of Pairs)](http://www.geeksforgeeks.org/dynamic-programming-set-20-maximum-length-chain-of-pairs/)  
[Greedy Algorithms | Set 1 (Activity Selection Problem)](http://www.geeksforgeeks.org/?p=18528)