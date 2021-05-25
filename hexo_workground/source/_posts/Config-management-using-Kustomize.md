---
layout: post
title: Config management using Kustomize
date: 2021-05-24 22:36:30
tags: [Kubernetes, Kustomize]
---

<img src="https://raw.githubusercontent.com/eliteGoblin/images/master/blog/img/picgo/20210525113423.png" alt="20210525113423" style="width:500px"/>  

## WHY

我们通过一系列的yaml来Config application deployed on Kubernetes. 

我们的每个application往往需要一组K8s resource, 即一组`yaml` files.

问题, 对于Kubernetes: 

*  我们如何高效管理不同application的manifest: "分治法", 不同application不同的folder, 之下是其group of manifest.
*  不同环境的manifest, 大体相似，细微config差异，如何reuse manifest, 但区别不同的环境?

<!-- more -->

在Kustomize之前，已经有很多tool来管理不同环境的manifest, 可以看出对其的需求之强, 参见[Config management design doc](https://github.com/kubernetes/community/blob/master/contributors/design-proposals/architecture/declarative-application-management.md)

在Kustomize之前，我们组用[ktmpl](https://github.com/jimmycuadra/ktmpl)来apply manifest到不同环境; 机制是` Parameterization Templates`: 所有的环境共享一致的template, template含有不同的parameters, 不同环境定制这些parameters.

Kustomize是Kunernetes官方的config management tool, 解决第二个问题. 机制是对将可复用的manifest定义在base文件中，不同环境对其做`patch`: 即base file + patch(in different env)组成最终apply的manifest.

## Config组织结构: 

*  `base`存放common manifest
*  `overlays`存放不同环境的patches

```
├── base
│   ├── configMap.yaml
│   ├── deployment.yaml
│   ├── kustomization.yaml
│   └── service.yaml
└── overlays
    ├── production
    │   ├── deployment.yaml
    │   └── kustomization.yaml
    └── staging
        ├── kustomization.yaml
        └── map.yaml
```

可以看出: 

*  `overlays`下每个folder代表一个环境: `staging`, `prod`
*  每个folder下都有`kustomization.yaml`: 为Kustomize的配置文件.

## Kustomize feature

通过file`kustomization.yaml`来配置:  

*  Name prefix: 不同overlays给k8s resource以不同的name prefix
*  Common label, annotation
*  ConfigMap, secret generator: 生成ConfigMap, Secret, with hash suffix; 并自动reference.
*  Diff of overlays: 以`diff`输出的形式，查看`overlays`对`base`的patch.
*  Image tag: 配置image's Tag
*  Namespace: 配置resource的namespace


## Hello world

代码见[这里](https://github.com/eliteGoblin/code_4_blog/tree/master/kustomize/helloworld), 示例来自[Demo: hello world with variants](https://github.com/kubernetes-sigs/kustomize/tree/master/examples/helloWorld)

`helloword` folder存放demo app的manifest, `tree base`

```
base
├── configMap.yaml
├── deployment.yaml
├── kustomization.yaml
└── service.yaml
```

`kustomize build base`会将`base`下所有manifest混合在一起输出, 一般这么用: `kustomize build base | k apply -f -`

加入`staging`和`prod`的patches, in `overlays` folder, 变成: 

```
├── base
│   ├── configMap.yaml
│   ├── deployment.yaml
│   ├── kustomization.yaml
│   └── service.yaml
└── overlays
    ├── production
    │   ├── deployment.yaml
    │   └── kustomization.yaml
    └── staging
        ├── kustomization.yaml
        └── map.yaml
```

### Base's Kustomization

`base/kustomization.yaml`config common labels for all resources: 

```yaml
commonLabels:
  app: hello
```

我们修改为`app: my-hello`, 验证是否所有`base`下resource都发生了改变: 

```sh
sed -i.bak 's/app: hello/app: my-hello/' \
    base/kustomization.yaml
```

通过`kustomize build base`可以看到效果: 

<img src="https://raw.githubusercontent.com/eliteGoblin/images/master/blog/img/picgo/20210525104303.png" alt="20210525104303" style="width:500px"/>  

`resources`列出`base`folder包含的resource. 

### Overlays intro

以`staging`的kustomization.yaml为例: 

```yaml
namePrefix: staging-
commonLabels:
  variant: staging
  org: acmeCorporation
commonAnnotations:
  note: Hello, I am staging!
bases:
- ../../base
patchesStrategicMerge:
- map.yaml
```

*  `namePrefix`表明以`staging-`prefix各resource name
*  `commonLabels`和`commonAnnotations`类似上述`base`的case
*  `bases`: 指明base manifest
*  `patchesStrategicMerge`: 以merge的形式update相应的resource: 

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: the-map # (1)
data:
  altGreeting: "Have a pineapple!"
  enableRisky: "true"
```

*  `(1)`指明了update target, name必须match.

### Diff

直接看下`staging`和`production`的不同(同理`base`与具体环境): 

```sh
diff \
  <(kustomize build overlays/staging) \
  <(kustomize build overlays/production) |\
  more
```

```diff
3,4c3,4
<   altGreeting: Have a pineapple!
<   enableRisky: "true"
---
>   altGreeting: Good Morning!
>   enableRisky: "false"
8c8
<     note: Hello, I am staging!
---
>     note: Hello, I am production!
12,13c12,13
<     variant: staging
<   name: staging-the-map
---
>     variant: production
>   name: production-the-map
...(truncate)
```

## Patches总结

Kustomize支持两种patch: 

*  `patchesStrategicMerge`: 给出`partial object`, 对原object做`Upsert`, [design doc: Strategic Merge Patch](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-api-machinery/strategic-merge-patch.md)  
*  `JSON patch`: 根据[JSON patch 标准](https://datatracker.ietf.org/doc/html/rfc6902), 对JSON object施加标准operations, 达到update的效果. 比merge更为imperative, operation有: 
    +  add
    +  remove
    +  replace
    +  move
    +  copy

### JSON patch

示例见[这里](https://github.com/eliteGoblin/code_4_blog/tree/master/kustomize/json_patch)

`json_patch`为demo folder, 下面`kustomization.yaml`: 

```yaml
resources:
- ingress.yaml
patches:
- path: ingress_patch.json
  target:
    group: networking.k8s.io
    version: v1beta1
    kind: Ingress
    name: my-ingress
```

*  `resources`: 待patch的manifest
*  `patches`: JSON patch文件，格式也为JSON
    +  `path`: file path
    +  `target`: 用来filter待update的resource, 可以一次[patch多个resources](https://github.com/kubernetes-sigs/kustomize/blob/master/examples/patchMultipleObjects.md)  


Verify patch diff: 

```sh
diff \                     
  <(cat json_patch/ingress.yaml) \
  <(kustomize build json_patch) |\       
  more
```

```diff
7c7
<   - host: foo.bar.com
---
>   - host: foo.bar.io
10,11c10
<       - path: /
<         backend:
---
>       - backend:
13,15c12,17
<           servicePort: 8888
<       - path: /api
<         backend:
---
>           servicePort: 80
>         path: /
>       - backend:
>           servicePort: 7700
>         path: /healthz
>       - backend:
18,19c20,21
<       - path: /test
<         backend:
---
>         path: /api
>       - backend:
21a24
>         path: /test
```

Note:  
*  JSON patch并不要求object为JSON, 也可为yaml, 见 `ingress_patch.yaml`为example.
*  `patches`之前版本也写为`patchesJson6902`

相关Kustomize文档: 

*  [JSON Patching](https://github.com/kubernetes-sigs/kustomize/blob/master/examples/jsonpatch.md)  
*  [Patching multiple resources at once](https://github.com/kubernetes-sigs/kustomize/blob/master/examples/patchMultipleObjects.md)  


## Reference

[简洁的Intro slide](https://speakerdeck.com/spesnova/introduction-to-kustomize)  
[JSON Patching](https://github.com/kubernetes-sigs/kustomize/blob/master/examples/jsonpatch.md)  
[Patching multiple resources at once](https://github.com/kubernetes-sigs/kustomize/blob/master/examples/patchMultipleObjects.md)  
