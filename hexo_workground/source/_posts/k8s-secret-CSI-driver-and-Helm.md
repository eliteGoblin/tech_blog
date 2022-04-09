---
layout: post
title: k8s secret CSI driver and Helm
date: 2022-03-30 04:05:33
tags:
---

# Background

The Secrets Store CSI Driver secrets-store.csi.k8s.io allows Kubernetes to mount multiple secrets, keys, and certs stored in enterprise-grade external secrets stores into their pods as a volume. Once the Volume is attached, the data in it is mounted into the container's file system.

In Kubernetes, you can use a shared Kubernetes Volume as a simple and efficient way to share data between containers in a Pod.

```
k create cm cm-test --from-file=./files/config.json --from-file=./files/db_password --from-file=./files/mysql
```


进入pod:

```s
root@webserver:/etc/config# ls -l
total 0
lrwxrwxrwx 1 root root 18 Mar 30 05:24 config.json -> ..data/config.json
lrwxrwxrwx 1 root root 18 Mar 30 05:24 db_password -> ..data/db_password
lrwxrwxrwx 1 root root 12 Mar 30 05:24 mysql -> ..data/mysql
```

每个key entry都有一个file, file内容为value.