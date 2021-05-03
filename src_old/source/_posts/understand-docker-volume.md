---
title: Understanding Docker Volume
date: 2018-11-07 08:45:05
tags:
keywords:
description:
---

{% asset_img main.jpg %}


## Preface

在Docker的使用过程中，data persist，share一直是比较令人困惑，比如：

*  container内部被改动的文件保存在哪里？
*  container运行结束或者被删除后，如何读取感兴趣的数据？
*  如何在host及container，container之间共享数据？

本文的目标通过介绍docker相关背景知识(union file system, volume)，给出完成常见volume操作的示例，最终达到解答这些问题的目的。

<!-- more -->

## Where does data reside in container

当我们在container中写入数据时，这部分数据保存在哪里呢？

我们首先看一下当前host磁盘容量: 

```console
df -h
文件系统                     容量  已用  可用 已用% 挂载点
udev                         3.8G     0  3.8G    0% /dev
tmpfs                        766M  2.1M  764M    1% /run
/dev/mapper/ubuntu--vg-root  139G   78G   55G   59% 
```

然后启动一个container，写入size为4G的文件: 

```
# ti 交互terminal 
docker run --name frank_ubuntu -ti  ubuntu
# 得到container bash
mkdir /test_data && cd /test_data
dd if=/dev/zero of=4g.bin bs=2G count=2
```

这时再在host运行df -h

```
文件系统                     容量  已用  可用 已用% 挂载点
udev                         3.8G     0  3.8G    0% /dev
tmpfs                        766M  2.1M  764M    1% /run
/dev/mapper/ubuntu--vg-root  139G   82G   51G   62% /
```

正如我们所料，新写入的4g.bin占了host上4G的空间，然后停止container：

```
文件系统                     容量  已用  可用 已用% 挂载点
udev                         3.8G     0  3.8G    0% /dev
tmpfs                        766M  2.1M  764M    1% /run
/dev/mapper/ubuntu--vg-root  139G   82G   51G   62% /
```

可见虽然停止了container，所占用的磁盘空间仍然存在，如何访问已停止的container呢：

```
# i interactive
# a attach STDOUT/ERR 
docker start -ia frank_ubuntu
```

为什么要访问已经停止的container的文件呢？典型的应用场景是：运行在container中的service crash，但是需要看已经stop的container存储的log，如果是用 --rm 启动的container，停止后会自动释放磁盘，因为container被remove了，磁盘容量因此会被回收。

## Docker file system

知道了container写入的文件具有与其相同的生命周期之后，我们来看一看docker的文件系统：

以如下Dockerfile为例：

```
FROM ubuntu:18.04

RUN mkdir /eg_test_1
RUN mkdir /eg_test_2
RUN touch /eg_test_2/hello_from_frank
```

查看输出

```
Sending build context to Docker daemon  4.096kB
Step 1/4 : FROM ubuntu:18.04
 ---> ea4c82dcd15a                              # lvl 1 标准ubuntu
Step 2/4 : RUN mkdir /eg_test_1
 ---> Running in 0f2785c298c8
Removing intermediate container 0f2785c298c8
 ---> fb70e5434d03                              # lvl 2 mkdir /eg_test_1 
Step 3/4 : RUN mkdir /eg_test_2
 ---> Running in dd69984f7b67
Removing intermediate container dd69984f7b67       
 ---> 8df2a5f2290e                              # lvl 3 mkdir /eg_test_2
Step 4/4 : RUN touch /eg_test_2/hello_from_frank
 ---> Running in 867b28469282
Removing intermediate container 867b28469282
 ---> cdea14636b37                              # lvl 4 touch /eg_test_2/hello_from_frank
Successfully built cdea14636b37
```

直观上看到docker构建时是分层的，每一条命令都会改变在前一步成功的基础上增加一层。由Dockerfile build之后的image由一系列的read-only层组成，最上面是一个read-write层，container运行起来之后，对文件系统做的改动，写入的数据等均会保存在最上面的读写层。

> docker image格式：自底向上，由一系列只读层，加上最上面的读写层，称之为Union File System

{% asset_img ufs.png %}

当container被删除后，再重新启动同样的image，会在最上面构建一个全新的读写层，之前container的数据被丢弃。

Union File System 并不能与外界(宿主机，NFS)共享file/directory，以及分离数据与container的生命周期，如何解决？简单的bind机制应运而生。

## Use bind to share

何为bind? 将host的file/directory与container共享，任何一方的修改都对对方立刻可见。bind命令可能大家并不陌生
```
# 将host的/tmp directory mount到container的/host_tmp directory
sd run -ti -v /tmp:/host_tmp --name frank_ubuntu --rm ubuntu
```

运行bind后，查看此container的Mounts选项：

```
docker inspect eager_cray -f '{{json .Mounts}}' | jq
```

注: jq是一个很好用的Command-line JSON processor，安装：

```
sudo apt-get install jq
```

可以看到Type, source, destination

```javascript
[
  {
    "Type": "bind",
    "Source": "/tmp",
    "Destination": "/eg_tmp",
    "Mode": "",
    "RW": true,
    "Propagation": "rprivate"
  }
]
```
另外：官方建议大家使用--mount option，避免-v src:dst 时bind和volume弄混的pitfall，虽然繁琐一点，但是清晰可读大于一切。接下来的例子，将尽量采用建议的做法，因为之前自己也被-v的灵活语法弄得很是混乱，而且对理解概念没有帮助。
```
# docker官方建议使用mount option，更verbose
docker run -it --rm --name frank_ubuntu\
  --mount type=bind,source=/tmp,target=/eg_tmp\
  ubuntu:18.04
```

这样container只要将想persist的数据写入与host bind的mount point，就可以同时实现与host数据共享以及生命周期分离。

这样做非常方便，我们在build image后，运行docker时指定要挂载的host file/folder即可。开发自己docker container管理系统的同学一定熟悉类似的命令：

```
docker run -d -p 9000:9000 -v /var/run/docker.sock:/var/run/docker.sock portainer/portainer
```

这是bind的一个典型用法：将host的docker daemon的unix sock文件，共享给了某个container，这样此container就可以通过读写本container的docker.sock文件，来调用host的docker API，实现host的container环境的管理了。

bind的一些典型用法: 

*  在host与container之间共享配置，例如container bind host的/etc/resolv.conf,实现DNS
*  共享source code及release file，方便container运行

借助bind，实现了数据快速共享及persist，但这是在借助host的文件系统而实现的，有没有更通用的共享方案呢？于是就引出了Volume。

## Persist and share using volume

先看下两段官方描述: 

> volumes are managed by Docker and are isolated from the core functionality of the host machine. A given volume can be mounted into multiple containers simultaneously. When no running container is using a volume, the volume is still available to Docker and is not removed automatically.

> A data volume is a specially-designated directory within one or more containers that bypasses the Union File System. 

我们先直观的感受一下volume: 

创建一个名为 frank_test_vol 的volume

```
docker volume create frank_test_vol
# 查看当前host上volume
docker volume ls
```
可以看到刚刚创建的volume
```
DRIVER              VOLUME NAME
local               frank_test_vol
```

查看volume详细信息: 

```
docker volume inspect frank_test_vol | jq
```

得到

```javascript
[
  {
    "CreatedAt": "2018-11-10T10:33:31+11:00",
    "Driver": "local",
    "Labels": {},
    "Mountpoint": "/var/lib/docker/volumes/frank_test_vol/_data",
    "Name": "frank_test_vol",
    "Options": {},
    "Scope": "local"
  }
]
```

发现Mountpoint在host上的一个路径，这里volume存储在host上:

```
sudo ls -l /var/lib/docker/volumes/frank_test_vol/
```

得到

```
drwxr-xr-x 3 root root  4096 11月 10 10:33 frank_test_vol
```

如何使用这个volume呢：

```
docker run -ti --rm --name frank_ubuntu \
--mount source=frank_test_vol,target=/my_vol \
ubuntu
```

这个命令--mount是在runtime将已经建立的volume mount到container的/my_vol处，而真正存储是在host的/var/lib/docker/volumes/frank_test_vol/

可以查看container的mount信息感受一下与bind的区别：

```
docker inspect -f '{{json .Mounts}}' frank_ubuntu
```

得到

```javascript
[
  {
    "Type": "volume",
    "Name": "frank_test_vol",
    "Source": "/var/lib/docker/volumes/frank_test_vol/_data",
    "Destination": "/my_vol",
    "Driver": "local",
    "Mode": "z",
    "RW": true,
    "Propagation": ""
  }
]
```

可以看到Type是volume而不是之前的bind，同样能看到mount的source和destination信息

目录，这样共享就完成了，而且任何一方对volume的改动都会被对方立刻看到：

```
# in container
touch /my_vol/hello_from_container
# in host
sudo touch /var/lib/docker/volumes/frank_test_vol/_data/hello_from_host
```

查看host/container上的对应目录，均会发现两个文件

```
-rw-r--r-- 1 root root    0 Nov 10 00:15 hello_from_container
-rw-r--r-- 1 root root    0 Nov 10 00:15 hello_from_host
```

可能你会问，假如container运行时想mount到某处，可以image在此处已经有entry，或者说有folder/file存在，会怎么样？ 答案就是会覆盖image的，以运行时挂载的为准。

我们可以看到，volume是独立于container而创建和manage的，container可以在运行时mount volume，而且此volume可以同时mount到多个container

总结一下，volume有如下feature: 

*  跨平台，支持windows和linux
*  Volume本身bypass了Union File System
*  Data volumes can be shared and reused among containers
*  Changes to a data volume are made directly
*  Changes to a data volume will not be included when you update an image
*  Data volumes persist even if the container itself is deleted
*  Volume drivers let you store volumes on remote hosts or cloud providers, to encrypt the contents of volumes, or to add other functionality.
*  New volumes can have their content pre-populated by a container(下面会谈到)

如果container有想独立于image的数据，可以将其存储于volume上，就像独立的u盘，于自己独立，而且可以很方便的与他人分享。

那这样每次我们想使用volume得先在命令行创建volume，再在runtime通过命令行参数绑定此volume，步骤比较繁琐。有时我们需要container运行起来，自动创建其container独有的volume，存放此次运行的log, data。这时匿名volume就会非常有帮助。

volume分为named和anonymous两种，区别是named就像我们之前演示的那样，我们手动创建，并赋予其name，而anonymous是在container运行时pre-populated，也就是自动创建的，并非没有名字，而是一串保证不会重复的随机串作为其名字，既然想自动化，而非命令行手动创建并在docker run时mount，那就需要我们记录在Dockerfile中，这就是Dockerfile的VOLUME命令出现的原因: 

如下Dockerfile
```
FROM ubuntu:18.04

RUN mkdir /my_tmp/
RUN mkdir /my_tmp/eg_test_1
RUN mkdir /my_tmp/eg_test_2
RUN touch /my_tmp/eg_test_2/hello_from_frank

VOLUME /my_tmp

RUN touch /my_tmp/eg_test_2/should_not_in_volume
```
VOLUME命令会干3件事: 

*  创建一个匿名volume
*  拷贝 container **当前** /my_tmp位置的内容到此匿名volume
*  mount volume到container的/my_tmp位置

build并运行
```
docker build --no-cache -t frank/ubuntu .
docker run -ti --rm --name frank_ubuntu frank/ubuntu
```

查看系统volume会发现匿名volume

```
DRIVER              VOLUME NAME
local               16b77fc36729cf3bcd0f37270f5cd4dd13cae7b4026ad1562dfd461104a289a8
```

inspect container, 仅显示Mounts

```
docker inspect -f '{{json .Mounts}}' frank_ubuntu | jq
```

得到

```javascript
[
  {
    "Type": "volume",
    "Name": "16b77fc36729cf3bcd0f37270f5cd4dd13cae7b4026ad1562dfd461104a289a8",
    "Source": "/var/lib/docker/volumes/16b77fc36729cf3bcd0f37270f5cd4dd13cae7b4026ad1562dfd461104a289a8/_data",
    "Destination": "/my_tmp",
    "Driver": "local",
    "Mode": "",
    "RW": true,
    "Propagation": ""
  }
]
```

我们可以看到volume在host上存储的位置，如果查看会发现docker image中建立的文件/文件夹已经被加入到新建立的匿名卷

```
/var/lib/docker/volumes/9cf0b9e60303d021796ee4fa2c661e8287fcfab43257bcd6d7ffe3c9a07717ba/_data/
├── eg_test_1
└── eg_test_2
    └── hello_from_frank
```

但是并没有 /my_tmp/eg_test_2/should_not_in_volume，在我们的Dockerfile中，此文件是在Volume命令之后加入到docker image的，因为docker构建时顺序执行，这时自然不会将还没有建立的文件拷贝到匿名volume中。

用docker命令行--mount也可以达到一样的效果

去掉上面Dockerfile的VOLUME命令: 

```
FROM ubuntu:18.04

RUN mkdir /my_tmp/
RUN mkdir /my_tmp/eg_test_1
RUN mkdir /my_tmp/eg_test_2
RUN touch /my_tmp/eg_test_2/hello_from_frank

```

```
docker run -ti --rm --name frank_ubuntu --mount type=volume,target=/my_tmp
```

和VOLUME效果一样，先创建一个空匿名Volume，如果mount的target不存在，创建一个空folder；如果存在，则把当前存在的内容拷贝至创建的Volume中。

## clean 

最后我们清理一下刚才测试过程中产生的垃圾

```
# 测试环境无所谓，干活时慎重使用!
docker system prune
# 清理没有attach的volume，
docker volume rm `docker volume ls -q -f dangling=true`
# 或者
docker volume prune
```

## Conclusion

本文从container数据持久化说起，引出bind: share host及container folder；之后提到更通用的方案: volume；有两种类型的volume: named及anonymous; 对于anonymous volume，两种使用方式： cmdline及Dockerfile的Volume指令。

*  bind用在简单共享文件，配置，源代码
*  仅想expose, persist container的数据(如程序log)，anonymous volume是首选，最方便的是Dockerfile的VOLUME命令，每次启动container会自动创建匿名volume
*  docker rm --rm 选项会删除随container创建的匿名volume
*  尽量用--mount而非-v选项，-v选项几种支持的模式语法上很模糊，容易混淆。
*  volume是存储的抽象，对container提供了同样的访问接口。不仅支持存放在host的文件系统，也支持NFS等。

常见的Volume命令总结:

```
docker volume create my-vol
docker volume ls
docker volume inspect my-vol
docker run -ti\
  --name frank_ubuntu \
  --mount source=myvol2,target=/app \
  ubuntu:18.04
sd inspect frank_ubuntu -f '{{json .Mounts}}' | jq
docker system prune
docker volume prune
```



## Reference

[Manage data in Docker](https://docs.docker.com/storage/)  
[Understanding Volumes in Docker](https://container-solutions.com/understanding-volumes-docker/)  
[Use bind mounts](https://docs.docker.com/storage/bind-mounts/)  
[Use volumes](https://docs.docker.com/storage/volumes/)  
[Dockerfile reference](https://docs.docker.com/engine/reference/builder/)  
[VOLUME 定义匿名卷](https://yeasy.gitbooks.io/docker_practice/image/dockerfile/volume.html)  