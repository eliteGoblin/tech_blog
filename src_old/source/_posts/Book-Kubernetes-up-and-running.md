---
title: 'Book: Kubernetes up and running'
date: 2019-10-05 16:39:08
tags: [k8s]
keywords: Kubernetes
description:
---


<div align="center">
{% asset_img k8s_up_running.jpg %}
</div>

## Preface

Kubernetes(K8s)，　被成为云原生操作系统，可谓是当下最火的技术之一， 同时也是容器编排的实际标准。Devops文化的流行，又对其传播起到至关重要的推动作用。  

由于最近在做搭建K8s集群的任务，也算是对K8s的使用正式入了门。实话说，K8s的学习曲线并不算平坦，概念繁多，其出色的松耦合设计在一定程度上又加剧了理解的复杂性。众多的技术术语往往搞得初学者一头雾水，面对分门别类的官方reference又不知让人从何下手。  

不过官方reference从来都不是初学者的好帮手，好比初学英语，我们也不能靠一本字典迅速取得突破。这时一本好的教材就显得尤为关键，这里我推荐的一本书便是 _Kubernetes up and running_， 属于O'REILY非常经典的up and running丛书，继承了系列的简约明快的风格，配合书中一步步可重现步骤(书提供配套代码)，快速帮初学者建立起基本概念，以及通过实践加深对知识点的理解。而且此书的作者之一是Kelsey Hitower老师，鼎鼎大名的K8s布道者，也是[Kubernetes the hard way](https://github.com/kelseyhightower/kubernetes-the-hard-way)的作者(目前已收获近18000 stars)。  

<!-- more -->

## The book

在书中全面介绍了K8s的核心概念: pods, deployments, service, label/annotations, ingress, replicateset, daemonset, jobs, configmap/secret, rbac等等。 而且本书的最新第二版也可以在[Azure官网免费获得](https://azure.microsoft.com/en-us/resources/kubernetes-up-and-running/)。　

以下是我制作的此书的mindmap，我觉得最终能建立起整个知识体系，明白各个component解决的问题以及是如何相互作用的，才算真正入门即开始掌握K8s。

<embed id="embed" src="k8s_big_picture.svg" type="image/svg+xml">

mindmup在线分享图也可以在[这里找到](https://atlas.mindmup.com/2019/10/1db2cef0e74811e98c02d176d0347468/kubernetes/index.html)

