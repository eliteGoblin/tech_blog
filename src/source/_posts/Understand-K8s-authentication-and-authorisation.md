---
title: Understand K8s authentication and authorisation
date: 2019-10-05 13:02:26
tags: [kubernetes, RBAC]
keywords: [authentication, authorisation]
description:
---

## Preface

在做应用开发时，有时会被权限的问题搞得很烦，不得不花大量的时间处理一系列的 _access denied_，但由于缺乏对权限系统的了解，往往陷入 _keep trying until it works_ 的境地，影响开发效率，解决过程中并没有习得什么系统性知识。  

不可否认的一点，安全性对商业系统至关重要，作为程序员，根据实际需要，或多或少需要了解一些相关知识和术语，如HTTPS, certificate, CA, X509, OAuth, IAM不等。

在K8s中也不例外，尤其是当我们自己搭建K8s集群，部署自己服务到K8s时，需要处理一系列权限问题，就需要我们对K8s授权系统有全面了解。

最近在小组内做了一次关于K8s RBAC的分享，系统的梳理了下K8s内部的authentication和authorisation，希望内容对大家有帮助。

<!-- more -->

## Authentication and Authorisation

成熟的权限系统都可以划分为： Authentication和Authorisation，分清两者的区别是理解权限系统的关键。

Authentication: 又称“验证”、“鉴权”, 回答用户是谁的问题。就像我们的身份证或passport，证明我们是谁。  


Authorisation: 又称授权，回答特定用户能做什么的问题。权限一般由一系列的规则组成，每个规则可以抽象为： user　verb resource 三元组。

Linux的权限系统是一个典型例子，比如 /etc/hosts文件：

```
-rw-r--r--  1 root root     222 Jul 12 14:44 hosts
```

只允许root user对其 read和write，这里user是root, verb是read, write, resource是/etc/hosts文件。

K8s内部，也通过一系列的规则定义各个用户的权限。Authorisation是通过RBAC模块完成的。

根据安全系统的**least priviledge**原则，一般默认不应授予用户任何权限，用户需要的每个权限需要显式指明。

## Slide

这里Authentication and Authorisation其实仅针对api-server本身，我们通过kubectl, dashboard访问K8s cluster，其实背后是直接向api-server发起HTTPS request，api-server的这里Authentication因此建立在HTTPS之上。

在[此次分享](https://go-talks.appspot.com/github.com/eliteGoblin/Notes/cs/presentations/topics/Kubernetes_RBAC/slide/main.slide#1)中，我们解决/回答如下问题:  

*  理解K8s的authentication和authorization的big picture，这是我非常推崇的学习方法，掌握了某个领域的big picture(最好画出思维导图)，能分解成为一系列的模块，并知道这些模块分别解决什么问题，如何结合在一些作为整体运行的，相比通过解决问题式的学习方法，更容易让我们获得自信，且面对未知问题知道从何做起。
*  RBAC作为K8s的授权模块，理解其关键概念：service account, ole/clusterrole, rolebinding/clusterrolebinding。
*  如何赋予自己的服务以K8s权限：如读取kube-system namespace下的secrets。


