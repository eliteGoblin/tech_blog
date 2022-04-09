---
layout: post
title: 'Introduction to Helm -- part 1: why and Helm features'
date: 2022-04-08 01:14:58
tags: [Helm, k8s]
---


# Background

最近在公司参与了create EKS cluster, 来取代自己维护KOPS创建的cluster. 其中重要的一项便是用Helm Chart来manage manifest, 替代之前简陋的K8s manifest template tool: [ktmpl](https://github.com/InQuicker/ktmpl). 特在此总结. 

<!-- more -->

# Why Helm

K8s manifest有templating的需求: 需要维护在不同region, 不同environment(dev, qa, prod, etc)的K8s cluster, 它们share mostly similar configuration, 但是特定参数不一致, 如: 

*  docker image version
*  ASG size
*  Dependency URL, 数据库地址(QA用QA DB, PROD用PROD DB)

等等, you get the idea.

主流的"templating" choice有:

*  Helm: 不仅仅是template engine, 它是全能的k8s package management tool: 即提供config的version control; 同时保存状态信息(在k8s secrets)，如当前安装的状态(Helm release status)
*  Kustomize: 严格意义并不是template engine

## Kustomize: the patch/overlay engine

我司首先采用了Kustomize, 主要是K8s official tool. 

基本原理是定义一个base.yaml, 不同环境在上面patch/rewrite

典型目录如下:

```
├── base
│   ├── kustomization.yaml
│   └── manifests.yaml
├── qa1-au1
│   ├── kustomization.yaml
│   └── patch_hpa_replicas.yaml
├── qa1-us1
│   ├── kustomization.yaml
│   └── patch_hpa_replicas.yaml
```

```
kustomize build qa1-au1
```

![](https://raw.githubusercontent.com/eliteGoblin/images/master/blog/img/picgo/20211116083127.png)

支持[JSON Patch](https://datatracker.ietf.org/doc/html/rfc6902), 针对复杂的patch场景: 

如
```
- op: add
path: /spec/jobTemplate/spec/template/metadata/annotations/iam.amazonaws.com~1role
value: arn:aws:iam::xxx:role/api-verifier
```

## Helm's templating engine

Kustomize是在base.yaml上，不同环境打不同补丁/rewrite， Helm像是选择填空: 

> A template is a form that has placeholders that an automated process will parse to replace them with values. Designed to perform a specific function, it marks the places where you must provide the specifics.

![](https://speedmedia.jfrog.com/08612fe1-9391-4cf3-ac1a-6dd49c36b276/https://media.jfrog.com/wp-content/uploads/2020/07/01160120/Kustomize-Template.png/mxw_1080,f_auto)  

> An overlay is a set of replacement strings. Blocks of text in the original file are entirely replaced with new blocks of text.

![](https://speedmedia.jfrog.com/08612fe1-9391-4cf3-ac1a-6dd49c36b276/https://media.jfrog.com/wp-content/uploads/2020/07/01160206/Kustomize-Overlay.png/mxw_1080,f_auto)  

结论是: 

Helm template更容易阅读(采用Go Template), 功能更强. 适合completely own的代码. 而Kustomize更适合给复杂的第三方k8s manifest做一些小修改，这样避免直接修改第三方manifest.

Helm和Kustomize可以结合使用: 

Helm generate K8s manifest => Kustomize patch => kubectl apply

# Helm as package management tool

Helm code是以Helm Chart形式组织:   `helm create mychart`

```
mychart
├── charts
├── Chart.yaml
├── templates
│   ├── deployment.yaml
│   ├── _helpers.tpl
│   ├── hpa.yaml
│   ├── ingress.yaml
│   ├── NOTES.txt
│   ├── serviceaccount.yaml
│   ├── service.yaml
│   └── tests
│       └── test-connection.yaml
└── values.yaml
```

```s
# 安装Helm Chart
helm install . --generate-name
helm uninstall mychart-1649393824
helm list
# NAME              	NAMESPACE	REVISION	UPDATED                                 	STATUS  	CHART        	APP VERSION
# mychart-1649393824	helm-test	1       	2022-04-08 14:57:05.498323154 +1000 AEST	deployed	mychart-0.1.0	1.16.0
helm get template  mychart-1649393824 # will show k8s manifest
helm get values mychart-1649393824
```

upgrade: 

修改template, 并bump `version` in Chart.yaml, to `0.1.1`, run `helm upgrade mychart .`

变化: 

*  Revision: release revision, 2 (初次安装为1)
*  CHART: 由`mychart-0.1.0`变为`mychart-0.1.1`  

```
NAME   	NAMESPACE	REVISION	UPDATED                                	STATUS  	CHART        	APP VERSION
mychart	helm-test	2       	2022-04-09 21:31:29.27852322 +1000 AEST	deployed	mychart-0.1.1	1.16.0 
```

rollback: `helm rollback mychart 1`: 1 为 Release number.


# Package 和 Repo

Package是把`mychart` folder, 打包为: `mychart-0.1.1.tgz`

Repo是: 

> an HTTP server that houses an index.yaml file and optionally some packaged charts.

常见的有: 

*  Github pages
*  ECR
*  S3

Helm可以add, remove repo, 也可以generate index.

下一篇discuss Helm的template engine.

# Lint and Specification

Lint: `helm lint --strict . --debug`

可定义 `values.schema.json` 规定value的structure

如: 
```
values.yaml         # The default configuration values for this chart
values.schema.json  # OPTIONAL: A JSON Schema for imposing a structure on the values.yaml file
```