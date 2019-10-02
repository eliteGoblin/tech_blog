---
title: which namespace does RBAC take effect in k8s
date: 2019-09-14 13:42:34
tags: [kubernetes, RBAC]
keywords:
description:
---

<div align="center">
{% asset_img rbac.png %}
</div>

## Preface

最近工作转向k8s，自己比较满意，之前很久就有兴趣，也看完了一本kubernetes up and running, 但并没有进一步深入研究。这次参与公司的k8s cluster升级项目，收获良多，算是正式入门了，看来深入的掌握知识果然还是得learning by doing，压力动力具备，目标导向，才能高效的学习。后续也会退出k8s, istio系列blog。

回归正题，在升级到k8s v1.13时，需要开启RBAC模块提升安全性，对每个service account(以下简称sa)赋予least priviledge。这需要建立特定的role/cluster role，并将其bind到service account。这就涉及到授权作用于哪个namespace的问题(namespace是k8s的一个机制：可以将k8s集群进一步分成子集群，起到隔离的效果)。

<!-- more -->

## 问题背景

K8s通过RESTFul API来实现declarative definition. 也就是k8s其概念都会对应为内部的Object, 我们需要的操作就会转化为对Object的CRUD操作。

RBAC作用的原理是request到达k8s api server后，authenticate为user/group或者sa, 以sa为例，检查其bind到哪个role/clusterrole上(bind是通过rolebinding/clusterrolebinding完成的，体现了k8s的松耦合设计)。

service account object需要指定namespace, rolebinding也有，同样role也有namespace，三者的namespace如何共同决定authorization起作用的namespace呢？ 感觉比较乱，其实不然，我们来实际验证一把。

## 分析&&验证

### 实验环境&&工具介绍

我们创建三个不同的namespace， 创建ns.yaml: 
```
---
apiVersion: v1
kind: Namespace
metadata:
  name: ns-sa
---
apiVersion: v1
kind: Namespace
metadata:
  name: ns-role
---
apiVersion: v1
kind: Namespace
metadata:
  name: ns-rolebinding

```
并创建
```
kubectl apply -f ./ns.yaml
```

k8s很方便的内置了auth subcommand来验证某个sa是否有特定权限, 假设我们查看name为test的sa是否在kube-system有list pods的权限: 

```
kubectl auth can-i list pods --namespace kube-system --as system:serviceaccount:kube-system:test
```

*  sa必须以system:serviceaccount:{namespace}:{sa name}的形式提供
*  --as表示当前user act as指定的sa，需要当前用户有impersonate权限，类似AWS的assume role。

接下来我们分别验证sa, role/clusterrole, rolebinding/clusterrolebinding在不同namespace组合下实际授权效果。

### service account并不决定作用namespace

标题是我们的结论，我们实际验证: 

建立test-sa.yaml
```
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: test-sa
  namespace: ns-sa
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  namespace: ns-role
  name: test-role
rules:
- apiGroups: [""] # "" indicates the core API group
  resources: ["pods"]
  verbs: ["list"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: test-rolebinding
  namespace: ns-role
subjects:
- kind: ServiceAccount
  name: test-sa
  namespace: ns-sa
roleRef:
  kind: Role #this must be Role or ClusterRole
  name: test-role # this must match the name of the Role or ClusterRole you wish to bind to
  apiGroup: rbac.authorization.k8s.io
```

并应用
```
kubectl apply -f test-sa.yaml
```
在namespace ns-sa下，test-sa并没有权限，但是在ns-role下就有:
```
kubectl auth can-i list pods --namespace ns-role --as system:serviceaccount:ns-sa:test-sa
# yes
kubectl auth can-i list pods --namespace ns-sa --as system:serviceaccount:ns-sa:test-sa
# no
```

可见sa的namespace只决定sa object建立在哪个namespace,并不能决定授权起作用的namespace。
cleanup:

```
kubectl delete -f test-sa.yaml
```

在证明了sa不影响后，clusterrole+clusterrolebinding的作用范围也很直观: 两者都是cluster wide，也必然在cluster wide范围生效。比较tricky的是role+clusterrolebinding和clusterrole+rolebinding。



### role不支持clusterrolebinding

下面的role_clusterrolebinding.yaml的validation失效: 

role_clusterbinding.yaml
```
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: test-sa
  namespace: ns-sa
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  namespace: ns-role
  name: test-role
rules:
- apiGroups: [""] # "" indicates the core API group
  resources: ["pods"]
  verbs: ["list"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: test-cluserrolebinding
subjects:
- kind: ServiceAccount
  name: test-sa
  namespace: ns-sa
roleRef:
  kind: Role #this must be Role or ClusterRole
  name: test-role # this must match the name of the Role or ClusterRole you wish to bind to
  apiGroup: rbac.authorization.k8s.io
```

```
kubectl apply -f role_clusterbinding.yaml
```

error msg: 

```
serviceaccount/test-sa created
role.rbac.authorization.k8s.io/test-role created
The ClusterRoleBinding "test-cluserrolebinding" is invalid: roleRef.kind: Unsupported value: "Role": supported values: "ClusterRole"
```

可见role并不能被clusterrolebinding绑定到sa。


cleanup:

```
kubectl delete -f role_clusterbinding.yaml
```

### role与rolebinding必须处于同一namespace


在之前模板中，我们发现rolebinding object并不能指定role的namespace(在roleRef字段), 如果两者分属于不同的namespace，运行时会提示找不到对应的role:

```
metadata:
  name: test-rolebinding
  namespace: ns-role
subjects:
- kind: ServiceAccount
  name: test-sa
  namespace: ns-sa
roleRef:
  kind: Role #this must be Role or ClusterRole
  name: test-role # this must match the name of the Role or ClusterRole you wish to bind to
  apiGroup: rbac.authorization.k8s.io
```

这其实表明rolebinding必须和role在同一namespace下。


### clusterrole + rolebinding, 用rolebinding的namespace

比较好理解，也是常用的pattern: clusterrole定义了模板，被不同namespace绑定，将priviledge限制到某个namespace。

clusterrole_rolebinding.yaml
```
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: test-sa
  namespace: ns-sa
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: test-clusterrole
rules:
- apiGroups: [""] # "" indicates the core API group
  resources: ["pods"]
  verbs: ["list"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: test-cluserrolebinding
subjects:
- kind: ServiceAccount
  name: test-sa
  namespace: ns-sa
roleRef:
  kind: ClusterRole #this must be Role or ClusterRole
  name: test-clusterrole # this must match the name of the Role or ClusterRole you wish to bind to
  apiGroup: rbac.authorization.k8s.io
```

创建RBAC: 
```
kubectl apply -f clusterrole_rolebinding.yaml
```

测试

```
kubectl auth can-i list pods --namespace kube-system --as system:serviceaccount:ns-sa:test-sa
# no
kubectl auth can-i list pods --namespace ns-role --as system:serviceaccount:ns-sa:test-sa
# yes
```

## 结论

经过分情况分析，我们发现: 

*  sa的namespace不影响授权作用namespace
*  role+rolebinding，两者必须处于同一namespace，仅作用于指定namespace下
*  clusterrole+clusterrolebinding，作用于cluster wide
*  clusterrole+rolebinding作用于rolebinding的namespace

再次分析可以得出： 只要看rolebinding/clusterrolebinding的namespace即可，同时值得注意的是不指定namespace kubectl会默认为当前namespace(如rolebinding的metadata:namespace)，容易造成unexpected behaviour。

## Reference


[Understanding Kubernetes RBAC](https://rancher.com/understanding-kubernetes-rbac/)  
[kubectl的用户认证授权使用kubeconfig或者token进行权限认证](https://jimmysong.io/posts/kubectl-user-authentication-authorization/)  
[Kubernetes namespace default service account](https://stackoverflow.com/questions/52995962/kubernetes-namespace-default-service-account)  
[How to: RBAC best practices and workarounds
Introduction](http://docs.heptio.com/content/tutorials/rbac.html)