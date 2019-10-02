---
title: A Prod ready SSH client in Golang
date: 2019-07-10 15:49:41
tags: [golang, ssh]
keywords:
description:
---

<div align="center">
{% asset_img main.png %}
</div>


## Preface

Linux下的SSH程序相信大家一定非常熟悉，可以说是登录远程服务器，执行远程命令的首选。有时我们还需要加入一些额外的逻辑来扩展标准SSH的功能，比如同时在一百台机器执行某个命令，或者在允许访问前，检查当前用户的额外权限，或者是否过期等等。

如何做到呢？我们可以利用golang的SSH library来实现一个定制化的SSH程序，然后再基础上加入我们的business logic，就可以啦。用Golang封装良好的SSH library实现一个最基本的SSH client并不难，网上也有很多相关博客，但用起来相比标准的SSH不少问题：

<!-- more -->

*  没有按照SSH方式处理按键: Ctrl+C不退出程序，Ctrl+D退出SSH client
*  SHELL颜色的显示
*  SSH terminal没有resize
*  记录用户所有输入及输出的log: 记录用户的输入输出对关键系统的audit很重要
*  定时退出功能：每次赋予用户的权限只在特定的时间段有效，超过则退出

本文实现的SSH client解决了以上问题，如果想做一款比较顺手的SSH client，可以参考本文展示的代码哦~


## Illustration

本SSH client包含两个基本功能，在远程机器执行命令及创建SSH session。我们

代码在这里[^1]

### 搭建测试环境

首先clone下示例代码
```shell
git clone git@github.com:eliteGoblin/code_4_blog.git
```

用vagrant启动一个虚拟机，作为我们测试的target host: 

```shell
cd golang_ssh_client/test
vagrant up
```

得到虚拟机的ssh配置信息
```
vagrant ssh-config
```
得到输出中的关键信息：

```shell
HostName 127.0.0.1
User vagrant
Port 2222
IdentityFile /home/frankie/git_repo/code_4_blog/golang_ssh_client/test/.vagrant/machines/default/virtualbox/private_key
```

### 运行

我们修改main.go中的const，使你的程序中hardcode的值与vagrant的输出一致: 

```golang
const (
    user    = "vagrant"
    host    = "localhost"
    port    = 2222
    keyPath = "/home/frankie/git_repo/code_4_blog/golang_ssh_client/test/.vagrant/machines/default/virtualbox/private_key"
)
```

然后我们运行

```shell
go build .
./myssh
```

我们就可以得到一个SSH窗口啦: 

<div align="center">
{% asset_img ssh_demo.png %}
</div>

由于我们在demo中执行了远程命令，

```
touch /tmp/test_ssh
```

可以看到命令成功执行，同时在10s后session自动结束。

## Conclusion

这里我们给出了一个的SSH Client，在此基础上我们可以借助Golang丰富的library，实现我们自己的验证，session时间控制等逻辑。

[^1]: [golang_ssh_client](https://github.com/eliteGoblin/code_4_blog/tree/master/golang_ssh_client)