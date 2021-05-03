---
title: K8s服务发现机制及DNS
date: 2019-11-14 18:57:43
tags: [kubernetes]
keywords: [kubernetes, service discovery, DNS]
description:
---

{% asset_img coredns.jpeg %}

## 背景

我们知道在K8s中，service用来abstract一组pod：意味着当我们想访问一组pod时，只需要直接访问这个service，service会自动帮我们挑选好一个pod，非常方便。  

比如我们有一个名叫tiles-app的Service: 

```
-   kind: Service
    apiVersion: v1
    metadata:
      name: app-svc
    spec:
      selector:
        component: app
```

在程序中我们可以用service name作为域名访问对应的service

```
http://app-svc/my/path/...
```

Service是K8s的抽象术语，它背后必然是借助经典的Linux/网络技术来实现。理解实现方式能加深我们对K8s系统的了解，而且我们可能也知道在K8s中，有DNS运行，DNS我们很多人对其的认知停留在name到IP地址转换，K8s为什么要用到DNS，这个内部DNS又和我们广域网上的DNS有什么关系呢？

前段时间在组内做了一次关于K8s内部服务发现及DNS的分享，这里做下简单的总结。

<!-- more -->

## 分享目标

理解下列概念:

*  K8s采用DNS作为服务发现的方式，app-svc作为domain name，K8s会给我们返回一个IP地址
*  Service几种发现模式: ClusterIP, Headless, ExternalName，前两者对应Server-side discovery和client-side discovery。
    -  ClusterIP是virtual IP，可以用IP table实现
    -  Headless类似经典DNS: service name返回一系列不同的IP(Pod的IP)
*  利用Dig和Host来调试DNS(镜像可以用 infoblox/dnstools)
*  DNS
    -  DNS是什么(多角度): server-client, distributed system, DNS tree
    -  DNS客户端配置解析: /etc/resolv.conf 
    -  为什么返回的域名最后有个"点": FQDN
    -  了解K8s对DNS的约定: K8s DNS specification
    -  DNS基本概念
        +  Zone
        +  A record
        +  CNAME
        +  Authority(DNS Server)
*  CoreDNS在K8s的基本工作原理: 监听api-server，有新service则将service name及对应的IP插入到CoreDNS记录中，可见K8s是通过古来的DNS技术来实现它炫酷的service discovery。
*  CoreDNS配置解析

以上内容比较多，这也是个把DNS理论和K8s机制串起来的机会，能加深我们对两者的认识。我们学习时应该尽量将孤立的知识点编织成网[^1]

## 成果检验

试着说说一个外部的request如何经过本地DNS cache, DNS server，再从root DNS server递归查询得到远端服务的入口地址(load balancer)，再到如何通过ingress controller，service，最后到达container内部的全过程，如果没什么问题，那恭喜你，你理解了。

示意图

{% asset_img result.png %}

总有一天我会用pyplot/latex画出精致配图的 😅

## Slide

之前的Slide都可以直接online read，这个不知道怎么了，一直提示有错误，目前只能在本地查看了。。

```
git clone git@github.com:eliteGoblin/talks.git
cd coredns/slide/
present .
```

Golang slide的详细方法看下Tony老师的博客吧[^2]

[^1]: [如何高效学习](https://book.douban.com/subject/25783654/)
[^2]: [Golang技术幻灯片的查看方法](https://tonybai.com/2015/08/22/how-to-view-golang-tech-slide/)