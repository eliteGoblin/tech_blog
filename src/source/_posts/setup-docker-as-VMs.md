---
title: setup docker as 'VMs'
date: 2018-10-14 17:47:29
tags: [docker ssh vm]
keywords:
description:
---

{% asset_img main.png %}

## Preface

项目需要用golang编写一个ssh client, 实现自定义的authorisation policy，log等。为了模拟目标环境，需要搭建多级跳板机的测试环境，要求host之间可以相互ssh。用vagrant，每个host起一个VM太'重'了，在我的旧T430上，环境启动至少需要10s+的时间。且会占据较多的系统资源；于是萌生了用docker实现轻量级VM的想法，在一番试验之后，实现了采用docker-compose和phusion/baseimage的组合来搭建多host的环境，启动耗时3s左右。

<!-- more -->

## What and Why

baseimage是一个经过改良过的ubuntu 18.04LTS，github仓库有6800+的stars。官方的image已经很好了不是吗，为什么要改良呢，几个主要的点

*  多process环境，作者的理念是docker container对应service，而不必非得是single process
*  优化init进程，避免僵尸进程产生
*  enable syslog，重定向至docker logs，避免异常日志被丢弃；
*  SSH server，这个是本文关注的重点
*  cron daemon，允许执行定时任务

可见对比official image，baseimage更像一个"VM"，可以拥有很多系统有用的特性，而且能完美满足自己的目标：搭建轻量，多host测试环境。

可以运行命令体验:  

```
docker run --rm -t -i phusion/baseimage /sbin/my_init -- bash -l
```

如果不需要执行命令，仅想让其像VM一样running，可以

```
docker run --rm -ti phusion/baseimage /sbin/my_init
# 启动之后，可以利用exec获得此container运行的shell
docker exec -ti container_name bash
```

## 实现ssh功能

利用baseimage，仅需要这样定义Dockerfile:  

```Dockerfile
FROM phusion/baseimage
# 避免php警告
RUN sed -i "s/^exit 101$/exit 0/" /usr/sbin/policy-rc.d
# 启动SSH server
RUN rm -f /etc/service/sshd/down
# 假如没有自己的ssh key的话，可以用内置脚本生成测试用key pair
RUN /etc/my_init.d/00_regen_ssh_host_keys.sh
# 必须先执行my_init
CMD ["/sbin/my_init", "--enable-insecure-key"]
```

运行环境并实现SSH

```
docker build -t frank/baseimage .
docker run --rm frank/baseimage
# 因为baseimage不是一次性process, 不会退出，需要重新开启一个terminal
# download and save private key to insecure_key
curl -o insecure_key -fSL https://github.com/phusion/baseimage-docker/raw/master/image/services/sshd/keys/insecure_key
# MUST do this, linux do not allow access to private key too open
chmod 600 insecure_key
# 登录container
ssh root@your_container_ip -i insecure_key
```

## 最终目标环境

利用docker-compose可以很方便的搭建多container的环境，实现一键启动，停止:  

docker-compose.yaml
```yaml
vesion: '3.5'
services:
    jumphost:
        container_name: 'jump_host'
        build 'path to image docker file'
        image 'fsun/jump_env'
        host_name 'jump_host'
        networks:
            mynetwork:
                aliases:
                    - 03.non.jumphost
    prodhost:
        container_name: 'prod_host'
        image 'fsun/jump_env'
        host_name 'prod_host'
        networks:
            mynetwork:
                aliases:
                    - 02.syd.prodhost
networs:
    mynetwork:
        dirver: bridge
```


*  host_name是命令hostname显示的结果
*  aliases用于网络DNS
*  baseimage作者提供了用于测试的ssh private key, 可以直接用来ssh dockercontainer

## Conclusion

通过本文的方法，大家可以在本地搭建由docker container构成的轻量级测试环境，可以相互SSH。再加上docker-compose，可以配置aliases作为环境内部的DNS，非常方便

## Reference

[phusion/baseimage-docker](https://github.com/phusion/baseimage-docker)