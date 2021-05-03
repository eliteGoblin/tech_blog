---
title: "初探Go Module: 理论篇"
date: 2019-11-19 17:50:19
tags: [golang]
keywords: [golang, module, gomodule]
description:
---

{% asset_img modules.png %}

## Preface

Package management是现代语言中非常重要的一环，我们也许已经开始用Go module管理package dependency，但在Go1.11之前，Golang官方并没有成熟的依赖关系系统，只有`go get`可以使用。  
在已经有了比较流行的解决方案：govendor, godep情况下，是什么使得Golang认为有必要推出module呢？module和这些"准官方"的solution比起来又有什么优势呢？

探索这些问题能帮助我们更好的理解Go Module，同时不仅仅应用于Golang，了解Module的机制其实让我们有机会了解一类通用问题的解决思路: dependency management(如Ubuntu中的apt-get)，体会为什么如Russ Cox所说: least version selection会是解决依赖问题的最优方案。

这篇文章也是对Russ Cox的Gomodule系列文章[^1]阅读的一个总结，大家感兴趣的话可以看看，非常深刻，会让你对Go Module的理解更上一个level。这篇文章只谈Go module的背景，理念，解决的问题。

<!-- more -->

## 为什么要有package management

问题咋一看很显然：我们当然希望只需要在code中import一系列package，就可以直接使用，而不必像最早C/C++那样，需要在系统中手动安装library，还要注意include的本地路径问题，现在package management能做到import一个类似URL的全局package name，它来帮我们下载好，省我们很多事。

听起来简单愉快，那有什么问题呢？问题就出在兼容性上。

我们直觉倾向用最新的东西: less bug, more features，一般情况下也是成立的。但麻烦的是breaking change，我们依赖A package，A升级了，修改了其中一个接口，造成我们编译失败。

同时试想一下我们在build一个非常复杂的系统，依赖关系也变得非常复杂，上千个依赖形成一张大网。升级一个package可能造成大面积影响: 我们升级了A, 而B依赖A，B也得升级，同时升级有可能break依赖B的C，D，牵一发而动全身。

## Package management is hard

Package management面对的是一系列复杂的问题，这也小伙伴么看到各种各样的apt-get依赖问题头很大的原因。

其中一个问题: Dependency hell: 

> The dependency issue arises around shared packages or libraries on which several other packages have dependencies but where they depend on different and incompatible versions of the shared packages. If the shared package or library can only be installed in a single version, the user may need to address the problem by obtaining newer or older versions of the dependent packages. This, in turn, may break other dependencies and push the problem to another set of packages

Russ Cox也在他的paper提到了本质是NP Complete问题[^2]，不一定总有解决方法。

{% asset_img version-sat.svg %}

> A needs B and C; B needs D version 1, not 2; and C needs D version 2, not 1. In this case, assuming it's not possible to choose both versions of D, there is no way to build A

这个问题也被称为"钻石型依赖"，我们将在后面解释Go Module如何解决它。

Package management 面对的核心问题便是依赖问题。

## Package management的基石: Versioning

正如标题所说，我们需要给每个package一个version，为什么呢？

我们先来看一个反例: `go get`，这里专指go 1.11之前的go get，不支持version，例子来源Russ Cox的presentation[^3]，非常值得一看。

假设我们依赖Library D，在之前我们已经`go get`了一份D到本地硬盘，但当时D的版本是1.0。

Library C1.8依赖D 1.4引入的接口，即使我们运行`go get`，它也不会下载新版，因为本地已经有D了，(go get 不认识version)，这时我们依赖　**太老**。

我们想解决这个问题:　升级所有package到最新版，`go get -u C`，但这时D的最新版1.6已经发布，但D的作者引入breaking change，与C1.8不兼容，还是会编译失败，这时我们依赖 **太新**。

{% asset_img goget_problem.png %}

没有Version，我们注定会陷入一系列的太老，太新的问题。

### Semver

这个奇怪的词是Sematic Version的缩写，其实我们很多人已经在用它了: 

{% asset_img semver.png %}

其实就是说我们用形如`1.12.3`标记一个版本，从遵循semver语义的版本号可以看出有无大改动，breaking change等。

详细定义请见其网站[^4]

## 灵魂拷问: 我们究竟需要什么

每个依赖都指明了版本，很好，在此基础上我们可以精确的指明依赖的版本。我们再回忆一下我们真正想要的是什么？

我们在构建软件，需要依赖。依赖需要正确的版本。什么是正确的版本呢？经过依赖库作者测试过得：我们想用依赖库作者用的一模一样的依赖库，说得很绕，如何理解？

假设我们需要D 1.3.0，D作者发布时，依赖C 2.5.4，因此我们也需要指明我们依赖C 2.5.4(间接依赖)。同时代码中我们不想指明package版本(太繁琐，而且升级依赖版本会造成修改代码)，因此代码中依赖C，我们的代码只能依赖同一个C。

同时我们不希望package management自动为我们升级到最新版本，因为升级可能会break。只有在确定需要的情况下才升级。

这种尽量让自己的依赖版本贴合依赖库作者发布时的情况，称为高保真(high fidelity) build，正是我们需要的。

## Module与其他package management的对比

Module已经是Golang的一部分，以下是它对比其他系统的优点：

*  上节提到的高保真构建(high fidelity build): 采用能满足需求的oldest version，因此需要一个文件记录需要的oldest version(go.mod file)，简单。
*  Upgrade平稳：因为只要满足oldest version，升级倾向保守，比如: 我们代码3个依赖库A, B, C都依赖了D 1.3.1, 只有A，B，C其中一个需要升级D才会升级，否则一直沿用D 1.3，即使一个很新的版本1.20.1已经发布。 
*  赋予库的使用者更多控制权，不允许库的作者将依赖pin到一个old version，库使用者可以在自己代码中指明依赖，这样避免了一些难以处理的情况：如我们依赖A 1.5，B2.3，B2.3指明其必须有依赖A0.1，而我们代码无法使用A1.5，build 失败。Go module module允许作者overwrite库的依赖版本。

### 兼容性

go modules 设计时天生支持其它package management(因为后有的module)，但其它系统一般不支持go module.

> Module-aware builds can also consume requirement information not just from go.mod files but also from all known pre-existing version metadata files in the Go ecosystem: GLOCKFILE, Godeps/Godeps.json, Gopkg.lock, dependencies.tsv, glide.lock, vendor.conf, vendor.yml, vendor/manifest, and vendor/vendor.json.

Go module会从其它package management的配置中读取依赖关系。

## Go module的feature

### 温和的upgrade: minimum version selection

我们来看一下Go moduel支持的操作，以及是其算法，能帮助我们更好的理解其操作。

Go module的版本选择核心算法只有百余行，如Russ Cox所说：

> It's a simple, clean algorithm.

算法支持4个基本操作，构建list，升级一个package，升级全部package，降级一个package。修改一个package可能造成依赖其的package改动，我们这里要尽量让改动最小化。

*  构建build list: 先加入target的直接依赖到一个list，再遍历list，加入间接依赖，如果发现有依赖同一package，但是不同版本，用两者较新版本。
*  升级一个package为指定版本: 构建升级前build list，再将升级版本依赖的list append到之前list，合并相同package，采用较新版本。
*  升级所有package到最新版本: 构建build list，确保每个package采用了最新版本。
*  降级一个package: 降级一个package，再确保所有依赖其的package不再require降级前的版本，必要的话降级这些package。

同时Module支持exclude和replace特定package，注意这两个操作只作用于本地，意味着如果你在写一个library，别人用你的library时，你的exclude和replace会失效，这也是之前提到的：赋予库的使用者更多控制权，特殊操作不应该由library决定，而应该由使用者决定。

更多细节参见Russ Cox的文章。[^5]

### 解决钻石型依赖: semantic import versioning

Go module规定: 

> If an old package and a new package have the same import path, the new package must be backwards compatible with the old package.

Known as Import Compatibility Rule.

结合Semver，这表示如果Major version没变的话，两者就认为是兼容的。Go module需要Semver来实现它的semantic import versioning。

这也意味着Go module用不同的import path来代表breaking change，所以两个不兼容的版本对于go来说，本质上是两个不同的package，因此这两个package可以独立选择版本。

假设我们在写一个library，当我们做出不兼容改动，我们需要升级Major版本号。

{% asset_img semver_import.png %}

正如下段话强调的: 

> Given a version number MAJOR.MINOR.PATCH, increment the: MAJOR version when you make incompatible API changes, MINOR version when you add functionality in a backwards compatible manner, and PATCH version when you make backwards compatible bug fixes.

但Major version是0, 1时，不需要加在path上，两者视为兼容版本。[^6]

将package不兼容的两个版本视为两个package，从code的import path上区分，解决了diamond problem。

## Module是万能的吗？

或者说，Go Module能解决任何情况的dependecy问题吗？当然不是，有些不兼容的情况需要两者的作者协作来解决。

软件的核心在于cooperation，Gomoduel作为一个工具只能解决部分问题。

{% asset_img cooperation.png %}

上述截图同样出自[^3]，时间上从26:26开始。

当Break发生时，Go module提倡使用者和library作者的共同effort来fix，比如: library C1.8依赖D1.4，当D升级到D1.6时，与C1.8不兼容。这时C和D的作者应该Cooperation: 

*  C作者回滚升级，仍指定依赖 >= D1.4
*  D作者发布D1.7，修复与C1.8的依赖问题。
*  C作者发布C1.9，指定依赖 >= D1.7

反例就是我们在dep和一些Linux的包关系系统中的约束: D >= 1.4; D < D1.6; 这会让依赖管理变得极其复杂，而且没有提倡cooperation，碰到坑就绕开，导致坑越来越多。

碰到Module无法自动处理的依赖问题时，需要找到不兼容的根源，然后不兼容的两者共同fix这个问题，同时Module在其中的角色是: 升级时尽量减少改变影响的范围，保守的采取满足requirement的oldest version，当不兼容被依赖包作者修复后，更新依赖。


## 结语

Go module已经从最开始发布的略显稚嫩，一片质疑中变为了成熟的Golang官方包管理工具，同时又在后续版本集成了proxy和checksum，解决依赖包被删除，不可用，及被篡改问题。

Gomodule已经被Hugo, Kubernetes(and it's client), Prometheus(and it's client), aws-go-client, Traefik, Syncthing, frp, etcd, nsq, caddy等一系列项目使用，我们小组最近也在计划从之前的govendor迁移到go module，因为涉及到CI和private repo的集成，因此需要一些额外工作，不过应该不太多。

在此推荐大家在新项目中直接用gomodule，除了目前vscode因为采用gopls对go module支持不够好，Goland无缝使用，Vim听说也没什么大问题。

[^1]: [Go & Versioning](https://research.swtch.com/vgo)
[^2]: [Version SAT](https://research.swtch.com/version-sat)
[^3]: [Go with Versions](https://www.youtube.com/watch?v=F8nrpe0XWRg)
[^4]: [Semantic Versioning 2.0.0](https://semver.org/)
[^5]: [Minimal Version Selection](https://research.swtch.com/vgo-mvs)
[^6]: [FAQ: Why are major versions v0, v1 omitted from import paths?](https://github.com/golang/go/issues/24301)