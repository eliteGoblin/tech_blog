---
title: Clean Architecure in Go (1) The Introduction
date: 2019-03-29 10:10:21
tags: [Architecure golang]
keywords:
description:
---

## Preface

当我们最初建立一个Golang项目时，很自然问题是：应该怎么layout比较好：即我们怎么应该组织项目代码，在Golang中有什么组织代码的约定"套路"吗？从这个看似简单的问题，表面上是我们如何安排目录结构，其实背后涉及到另一个有趣的问题：我们应该如何architect一个新项目，因为layout其实是project architecture的反映，architecure变化肯定会让layout产生变化。

本文中，我们先看大型golang项目在layout方面的普遍做法，即大型项目的目录结构。基于此进一步深入思考：这种做法能否得出一个清晰的project architecure：来得到更好的项目架构内在指标：低耦合度，可维护性/扩展性，可测性。 如果不能，怎样做才能更好的解决上述软件架构中的pain points? 本文给出的答案是Uncle Bob提出的Clean Architecure。 

最后我们会通过一个简单的case来看我们如何实现Clean Architecture。

<div align="center">
{% asset_img CleanArchitecture.jpg %}
</div>

<!-- more -->

关于这个主题，我在team内部也做了一个分享，见[这里](https://go-talks.appspot.com/github.com/eliteGoblin/Notes/cs/presentations/topics/clean_architecture/slide/clean_arch.slide)

## Golang "Standard" Layout

一般的大型Golang软件的目录结构如下：

<div align="center">
{% asset_img standard_layout.png %}
</div>

这里的Standard其实并不是指golang官方标准，这里总结的layout来自于: [^1] 

在这里各主要目录解释如下

### /cmd directory

项目中build成的binary的入口main.go应该放在此处，如果有多个binary，建立对应的子目录。

此目录下应该实现的是各binary的main module，不应该包含太多代码。作为daemon运行的service的binary也放在此处，这里cmd有些误导，并不专指binary command。

### /internal directory

Go 1.4引入的特性，代码放在internal/下面，不可以被非parent的代码import，也就意味着不能被其他project import，所以一般放置仅限repo内部使用的library。

### /pkg directory

与internal/相反，这里存放可以被其他project复用的library code，应该仔细考虑这里的code是否足够成熟以被其他project复用。

见Prometheus例子[prometheus/prometheus/pkg](https://github.com/prometheus/prometheus/tree/master/pkg)

### /api directory

用来存放api接口的约定，如swagger和protcolbuffer定义文件，参见[OpenShift API](https://github.com/openshift/origin/tree/master/api)

### /configs directory

很好理解，存放配置文件或者配置文件模板，以及诸如consul-template，systemd service file(也可以存放在/init目录下)等等。

### /build directory

一般有两个子目录：

*  build/package directory: cloud (AMI), container (Docker), OS (deb, rpm, pkg) package配置及脚本。
*  build/ci: CI(travis, circle, drone)配置及脚本。

### miscellaneous directory

- /deployments:  system and container orchestration deployment configurations
- /tools: supporting tools
- /scripts: supporting scripts
- /docs: 除了godoc，其他额外的doc放在此处
- /assets: logo, images及其他static files

## Thinking in "Standard Layout"

上述的layout给出了golang project的大致划分和命名的建议，但我们可以看出，大部分的代码还是需要放在/internal下面的，因为一般项目并不会成熟到有很多代码会放在/pkg供别的项目引用。  

对于一个庞大的service来说，我们对于处在internal/目录下的代码如何组织呢？我们怎么layout才能使项目变得更低耦合，可维护/更改及可测试呢？  

显然这个layout除了一些很basic的划分并没有给出上面问题的答案，而这个问题正是layout的**最关键**问题，同时也引出了: layout背后反映的是怎样一种架构？我们如何更好的架构我们的项目？


## Pain points in current golang layout

我们先来总结一下我们日常编程中，遇到的架构/layout的痛点：

* Hard to change, decisions taken too early
* Centered around frameworks/database, business logic is spread everywhere
* Hard to find things needed
* Test: slow, heavy; low test coverage, hard to mock; 
* Circular dependency

其实这些问题绝大多数是由错误的dependency，导致项目不同module之间杂乱的couple在一起：

*  business logic依赖细节如db实现，造成了不能轻易替换db甚至修改db schema
*  高层逻辑依赖细节的另一个坏处就是不容易进行测试，只能启动db实例才可以。
*  过度依赖技术framework，核心逻辑散落在各处，不易寻找。

所以我们的最初问题转换成了：寻找一个低耦合，依赖关系正确的architecure，来解决我们的痛点。

这个属于软件工程的经典话题，我们已经有了很好的指导原则：SOLID。

## SOLID Principle

- Single responsibility principle
- Open–closed principle
- Liskov substitution principle
- Interface segregation principle: 明确细分interface的功能，而不是包含多种功能的 general interface
- Dependency inversion principle

这5个原则中，以下两个和Clean Architecture更为相关：

*  Liskov替换原则：软件模块应该抽象出接口，同样接口的object可以替换而不需要调用者修改代码。
*  Dependency inversion principle: 高层逻辑不应该依赖于细节实现，两者都应该依赖于抽象，即interface。

### Dependency Inversion Principle

原则的官方定义：
> 
* High-level modules should not depend on low-level modules. Both should depend on abstractions (interfaces)
* Abstractions should not depend on details. Details (classes) should depend on abstractions

这个耳熟能详的原则是理解Clean Architeture，即我们接下来讲述架构的关键。

我们先看下没有实践Dependency Inversion Principle的一个典型的架构: 

<div align="center">
{% asset_img dip.png %}
</div>

从架构我们可以看到：控制流(A call B，控制流为A-->B)与依赖流完全一致，导致高层，抽象的模块依赖于低层模块；如果底层模块改变实现细节，可能会引起中层，高层模块的一系列修改。这种架构使模块之间紧紧的耦合在一起。

我们可以利用Dependency Inversion Principle来invert模块的依赖：使底层模块的实现依赖于高层模块定义的接口：

<div align="center">
{% asset_img dip2.png %}
</div>

在上图中我们看到，HL1 package定义了interface，并又ML1 package来实现，使之前HL1依赖于ML1**倒置**过来了，变成ML1依赖HL1，这就是依赖倒置原则。

这个原则的关键在于：我们可以在我们想的地方invert任何依赖，从而修复诸如高层(更抽象的模块)错误依赖底层细节的现状。

## Clean Architecture

讲到这里，郑重引出本文的主角: Clean Architecture(见本文开始的配图)

什么是Clean Architecture呢？是由Uncle Bob老师强力提出的一种软件架构方法，核心是分层架构，约束依赖，博客见这里：[^2]

Clean Architeture结构图的直观感受：

*  分层的结构，每个同心圆代表软件的一层
*  最外是device, UI, db实现等细节层，最内是核心业务逻辑的表现：use case层和entity层
*  越内的层越抽象，是policy；而越外的层越底层，细节，是mechanism；
*  清晰的依赖关系：依赖流只能由外向内，不能相反；细节永远依赖高层，抽象的部分，也即遵循dependency inversion principle

Clean Architecture就先介绍到这里，下一篇会包含: 

*  每层划分的具体实践
*  什么是Use Case
*  分析简单例子来了解如何进行Clean Architecture

Stay tuned，真香警告

<div align="center">
{% asset_img wjz.gif %}
</div>

[^1]: [golang-standards/project-layout](https://github.com/golang-standards/project-layout)
[^2]: [The Clean Architecture](http://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)  





