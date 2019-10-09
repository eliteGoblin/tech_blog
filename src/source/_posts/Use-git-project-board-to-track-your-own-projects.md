---
title: Use git project board to track your own projects
date: 2019-10-09 12:56:34
tags: [tools]
keywords:
description:
---

{% asset_img kanban_main.jpeg %}

## Preface

在日常工作学习中，经常会冒出很多想法记到todo list中，如抽时间看的文章，研究的开源项目，想要写的blog。之前会都记录在个人笔记里(markdown)，笔记会同步到github。  

用文字的方式来记录，管理todo list，有很多drawback: 

*  一个大文件来记录事项和包含的细节，杂乱无章
*  不可视化，无法看到过去完成的条目，缺少完成任务的正面反馈(要知道完成一项具体任务带来的满足感符合人类天性)
*  无法方便的附加截图等。

Github提供了project board,　实现了kanban board的功能，用来管理自己的todo绰绰有余。

设置好之后，现在每天一到办公室就会先看自己的project board，之前杂乱的todo list变成可视化条目，完成小任务带来的成就感也变成我一天中期待的事。而且在使用过程中，也感到自己speed变得smooth起来，而且可视化使得制定和实现周目标变得更容易。同时其很好的集成了github的一些功能，实现了基本的自动化，让我们就来看看它是如何工作的吧。

<!-- more -->

## Kanban board

我们先来看看什么是Kanban board: 

> A kanban board is an agile project management tool designed to help visualize work, limit work-in-progress, and maximize efficiency (or flow). Kanban boards use cards, columns, and continuous improvement to help technology and service teams commit to the right amount of work, and get it done!

直观上看，kanban board分为很多列，每列记录一个个card, 对应一项具体任务(如JIRA ticket)。

{% asset_img kanban.png %}

列可以是:  ready, in progress, done代表card的不同生命周期，也可以根据自己需要创建别的列。比如我加了一列：Committed this week代表本周的计划。 当card状态改变时，将其拖到对应列即可。

比如我的[project](https://github.com/users/eliteGoblin/projects/2): 

{% asset_img my_board.png %}

这样自己的TODO list就方便的建立起来，随时可以访问。同时本周的计划一目了然，而且也能看到之前完成的任务，获得成就感。

突然之间变得agile起来了呢，想想都有点小激动 😁

## Reference

[about-project-boards](https://help.github.com/en/articles/about-project-boards)  
[Markdown玩转Emoji](https://www.jianshu.com/p/e66c9a26a5d5)