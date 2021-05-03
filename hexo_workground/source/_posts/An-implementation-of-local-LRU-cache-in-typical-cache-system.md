---
title: Tech interview-Implement a in-memory cache
date: 2019-07-11 19:51:45
tags: [golang, tech interview]
keywords:
description:
---

{% asset_img memory_hierarchy.jpg %}　


## 背景

在悉尼找IT工作与在北京找互联网工作有不一样的套路，这边IT中小型公司较多，常常给你做一个具体的assignment，往往是实现一个简化版的系统，大约花2-6小时不等，然后下一轮面试主要基于此问问题。更看重的是系统设计，一般工程实现能力，及coding style, test等细节，不像大厂那样上来几轮算法，看你是不是他们要找的够smart的人。  

<!-- more -->

这种assignment显得更为务实，对于一般公司，能出活最重要，而且在讨论solution的过程中会有很多沟通，why，以及challenge，也很能看出一个人的沟通，团队合作，性格等方面有更深入的了解。个人还比较推崇这种面试风格，一般创业，中小型公司更容易找到做事踏实，适合自己团队的人。

之前在找工作的时候，又接到了一个assignment，实现一个简易版的in-memory cache，算是系统中很经典的一个问题，花了约一天的时间思考和实现，在此和大家分享一下。

## 问题

系统描述如下：

*  Data以key-value形式存放于database中，R/W延时500ms
*  有distributed cache可用，启动时是空的，可以存储key-value，R/W延时100ms。
*  在database中存放的data不会变更(never changes and can be cached forever)
*  database中当前不存在的key永远不会存在

原题可以从这里找到: [^1]

## 分析及设计

*  利用in-memory cache，缩短访问时间，可以直接用golang的map做cache
*  实现LRU，实现cache替换，防止内存溢出，R/W复杂度O(1)，实现思路来自[^2]
*  若从数据库发现数据不存在，给key赋予特殊值，这样能使后续同样的查询请求直接返回，不必击穿至数据库层。

非功能性设计：

*  Clean Architecture，分离business logic及infra detail
*  用Interface隔离细节，实现dependecy injection
*  基于interface, 用gomock来实现unit test
*  不同层cache抽象出同样的interface，用递归实现了多层cache的read-through功能
*  给出了cache library的benchmark

代码见这里: [^3]

有关cache的更新套路，耗子哥的文章非常棒: [^4]


[^1]: [go-test](https://github.com/eliteGoblin/code_4_blog/tree/master/cache_solution/task_src/go-test)
[^2]: [LRU Cache 最近最少使用页面置换缓存器](https://www.cnblogs.com/grandyang/p/4587511.html)
[^3]: [cache_solution](https://github.com/eliteGoblin/code_4_blog/tree/master/cache_solution)
[^4]: [缓存更新的套路](https://coolshell.cn/articles/17416.html)