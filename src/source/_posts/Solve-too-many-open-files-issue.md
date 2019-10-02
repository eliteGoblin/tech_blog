---
title: 由Too many open files问题看system limit
date: 2018-02-23 17:49:00
tags: [linux limit]
description:
---

<div align="center">
{% asset_img toomanyfilesopen.jpg %}
</div>

#### Preface

本文的缘起是作者写的golang服务，某天发现日志报错: socket: Too many open files, 定位到出问题的行:  
```golang
import "http"
client = &http.Client{
  Transport: tr,
  Timeout:   time.Duration(5 * time.Second),
}
response, err := client.Do(req) // 这里报错!
```
系统初始化时，会发起约1500余次http请求，看错误提示应该是http client在open socket时打开文件过多，超出系统限制, 但是通过linux命令ulimit将支持最多打开的文件数限制到65535, 仍是会报此错误．为什么限制远大于程序需要仍会发生这种错误? 如何修复这种看似毫无头绪的错误呢? 我们一步一步来.    

<!-- more -->

#### Ulimit命令

出现这种和超出系统资源限制的提示时，首先想到Linux的ulimit命令: 我们首先　man ulimit 
>ULIMIT(3)                                                                              Linux Programmer's Manual                                                                             ULIMIT(3)
NAME
       ulimit - get and set user limits
SYNOPSIS
       #include <ulimit.h>
       long ulimit(int cmd, long newlimit);
DESCRIPTION
       Warning: This routine is obsolete.  Use getrlimit(2), setrlimit(2), and sysconf(3) instead.  For the shell command ulimit(), see bash(1).

可见，这里ulimit是指Linux的API,并不是我们的shell. 我们输入 man bash, 找到ulimit部分:  
> ulimit [-HSTabcdefilmnpqrstuvx [limit]]
              Provides  control  over  the  resources  available to the shell and to processes started by it, on systems that allow such control.  The -H and -S options specify that the hard or soft
              limit is set for the given resource.  A hard limit cannot be increased by a non-root user once it is set; a soft limit may be increased up to the value of the hard limit.   If  neither
              -H  nor  -S  is specified, both the soft and hard limits are set.  
              ...

这里有几个重点需要注意的:  

*  ulimit设置的是shell的resource limit, 和通过shell启动的process的.  
*  limit有hard limit和soft limit之分, 通过 H或S来指定
    -  non-root user不能increase system limit
    -  soft limit不能设置为超过hard limit

为什么要有soft, hard limit来共同作用呢? 此limit是指什么的limit, 是系统的么?　从man说明的我们仍有疑惑:    
>ulimit -n The maximum number of open file descriptors  

#### getrlimit and setrlimit Functions

查阅了APUE的 7.11 getrlimit and setrlimit Functions, 我们知道:  

*  对于process，有一系列的limits, 其中一些可以用函数getrlimit/setrlimit来读取,设置
*  通过rlimit struct来和Linux API传递limit数据
  ```c
  struct rlimit {
    rlim_t rlim_cur; /* soft limit: current limit */
    rlim_t rlim_max; /* hard limit: maximum value for rlim_cur */
  };
  ```

*  soft limit用来设置当前session/process的limit, 而hard limit是为soft limit设置的资源上限, 两者的关系
    -  process只能设置soft limit不大于hard limit
    -  process只能decrease hard limit但必须不小于soft limit, 即必须时刻保证: soft_limit <= hard_limit
    -  只有superuser process能increase hard limit
*  设置的limit影响calling process, 并会被 **children process继承**, 这也是为什么shell调用ulimit设置，后续在此shell运行的process会采用设置的值

我们的问题所关心的参数RLIMIT_NOFILE:  

> RLIMIT_NOFILE The maximum number of open files per process. Changing
this limit affects the value returned by the sysconf function
for its _SC_OPEN_MAX argument  

可见此次错误仅和此process有关，与系统已经打开的文件句柄数无关  

#### golang 打印 process limit

既然Linux提供了get/set syslimit的功能, 在golang里我们可以直接调用syscall来查看在出错前的RLIMIT_NOFILE被设置为多少来确认一下是否此错误为我们所想是此limit导致的: 

```golang
import "syscall"
var limit syscall.Rlimit
syscall.Getrlimit(syscall.RLIMIT_NOFILE, &limit)
fmt.Println("limit got %+v", limit)
// Output: 1024
```

Bingo! 确认我们的程序在执行时RLIMIT_NOFILE并没有按我们的预期设置为65535, 导致了问题  

#### 问题的解决

最终确认问题是由于我们使用了systemd启动此程序，systemd默认此程序的RLIMIT_NOFILE为1024, 需要修改对应的service文件:  

```shell
# udesk_cti.service
[Unit]
...
[Service]
LimitNOFILE=65535
[Install]
...
```
修改后，执行 sudo systemctl daemon-reload, 问题解决

ulimit修改的是本session的limit, 如果想在system-wide修改, 需要修改 /etc/security/limits.conf, 加入配置.  

插入配置的格式是: 
```shell
<domain>  <type>  <item>  <value>
```
如果想设置nofile默认为65535, 在文件中加入如下两行即可  
```shell
* hard nofile 65535
* soft nofile 65535
```


#### 后记&&结论

其实当时被这个问题困扰了一阵子，真正的解决是在一个同事发现不经由systemd启动的服务没有此问题才开始怀疑systemd的配置的. 之后看了APUE的system limit一节才全部弄懂全部细节，而非像文中描述的抽丝剥茧式一步一步推导解决的．另一个启发是: 当软件层次比较复杂时，从最接近问题之处(即错误现场)来追寻: 如直接打印本process看到的limit是否是自己预期的65535, 能很快定位问题．

#### 参考文献

*Advanced Programming in the UNIX® Environment* aka APUE  
[Systemd services and resource limits](https://fredrikaverpil.github.io/2016/04/27/systemd-and-resource-limits/)  
[Package syscall doc](https://golang.org/pkg/syscall/#Getrlimit)  
