---
title: 'Setuid in Linux File Permission'
date: 2018-11-12 08:39:55
tags: [setuid, linux]
keywords:
description:
---

{% asset_img main.jpg %}

## Preface

在linux下，有时我们的进程需要临时提升权限，比如临时在file system中当前用户不具备access之处创建文件，修改记录，应该如何做到呢？linux提供了setuid/setgid机制，我们经常会用到的sudo就是其中的一个例子，接下来让我们看看相关的linux知识，并分析几个实际的例子加深理解。

<!-- more -->

### Big Picture

Linux的权限系统简洁且优雅，不同module像齿轮一样紧密配合，主要组成部分:  

*  File Access Permission Bit
*  User/Group
*  User/Group的运行时切换

File Access Permission Bit 和 User/Group机制组成了我们最熟悉的权限系统，File Access Permission Bit规定了不同user对此文件的access，User/Group提供了不同User/Group的标识，两者一起决定着用户可否对特定文件执行某个操作。

### File Access Permission Bit

我们对linux的file access permission bit一定不陌生，我们来review一下: 

```
ls -l
-rw-rw-r--   1 frankie normal_users   2355 10月  4 08:28 _config.yml
drwxrwxr-x 329 frankie normal_users  12K 10月  3 16:55 node_modules
```


<div align="center">
{% asset_img permission.png %}
</div>

file 同时还存有userId, groupId, 这里例子分别是frankie及normal_users组。基于特定的userId/groupId，permission bit才能根据permission bit，得出该用户/用户组有无对file的特定权限。


### User/Group

File Permission Bit解决了: 给定用户对本文件有什么权限的问题，User/Group机制正是提供了当前用户是谁的问题。

Linux系统的每个用户都对应有userId(也称为uid)，每个userId都属于某个group(group有唯一groupId，也称为gid)。每个user是否只能属于1个group呢，不是的，可以同时属于多个group，称为supplemenray group

> A user on Linux belongs to a primary group, which is specified in the /etc/passwd file, and can be assigned to multiple supplementary groups, which are specific in the /etc/group file

> Primary group – Specifies a group that the operating system assigns to files that are created by the user. Each user must belong to a primary group.
> Secondary groups – Specifies one or more groups to which a user also belongs. Users can belong to up to 15 secondary groups.

利用id命令查看当前用户uid, gid等信息
```
uid=1000(frankie) gid=1000(frankie) 组=1000(frankie),4(adm),24(cdrom),27(sudo),30(dip),46(plugdev),113(lpadmin),128(sambashare),999(docker)
```

可见frankie用户，除了primary group为normal_users外，还属于sudo, docker等group。supplemenray group与primary的区别除了create file默认为primary group外，用作permission bit时效果是一样的：

假设用户frankie的由supplementary group之一是docker，文件access bit如下:  

```
srw-rw---- 1 root docker 0 11月 12 08:05 /run/docker.sock
```

则frankie对 /run/docker.sock文件有 rw- 权限，因为frankie属于docker用户组(由supplementary group指定)


### UserId/GroupId in processes

前面说到User/Group结合File Permission Bit决定最终File Access Permission，也就是是否具有某权限。我们谈到用户是否有权限时，其实指的是用户的 **进程** 是否对某文件有rwx的权限，例如我们在编辑某文件，提示permission denied，其实是用户启动的VIM，对文件不具备写权限。  

那进程也需要有UserId/GroupId，系统才能由此结合File Permission Bit来最终判定进程的permission。那么用户启动的进程，UserId/GroupId是否就是与当前用户的完全一样呢？我们需要了解一下Linux中process内部UserId/GroupId实现机制。 对于GroupId，其实现完全与UserId相同，下面仅讨论UserId。 

Process 关联有两组UserId:  

*  RealUserId
*  EffectiveUserId

> The real user ID and real group ID identify who we really are. These two fields are taken from our entry in the password file when we log in
> The effective user ID, effective group ID, and supplementary group IDs
determine our file access permissions

一般这两个UserId是一样的，例如我们执行VIM程序时，RealUserId是我们自己，与当前用户登录名相同。EffectiveUserId也同样是当前用户，根据file permission bit，结合当前用户是否是file的owner，及所属Group决定对文件的access permission。那为什么除了RealUserId外还需要一个EffectiveUserId呢？  

因为有时我们需要process以file owner的身份执行，而非当前用户，这样就实现了"相当于以另一个用户的身份执行此程序"，这就是setuid bit机制。

> When the setuid bit is used, the behavior described above it's modified so that when an executable is launched, it does not run with the privileges of the user who launched it, but with that of the file owner instead. For example, if an executable has the setuid bit set on it, and it's owned by root, when launched by a normal user, it will run with root privileges. 

*  setuid bit仅作用于file，对folder无效
*  setuid需要file owner的execute bit开启，否则无效

回到之前的问题，正是由于process有时需要以另外用户的身份运行，而不总是与当前登录的用户相同，因此需要额外的EffectiveUserId来区分：RealUserId标识真正登录的用户，EffectiveUserId正如其名，标识当前作用于File Permission Bit的UserId，就是说检查Linux检查某进程的Access Permission，看的是EffectiveUserId，而不是RealUserId，后者仅起到记录此process是被哪个用户启动的作用。


## How to use setuid bit

setuid bit是一个特殊的file permission bit, 我们通过一个简单golang程序程序，演示setuid的作用，[play ground](https://play.golang.org/p/SCPkBcNeu2_T)

test_setuid.go
```
package main

import (
    "fmt"
    "os/user"
    "syscall"

)

// getRealUserID is to get the real user id of current process
func GetRealUserID() int {
    return syscall.Getuid()
}

// GetEffectiveUserID is to get the effective user id of current process
func GetEffectiveUserID() int {
    return syscall.Geteuid()
}

func GetUserNameFromId(userId int) (name string, err error) {
    u, err := user.LookupId(fmt.Sprintf("%d", userId))
    if err != nil {
        return "", err
    }
    return u.Username, nil
}

func main() {
  fmt.Println("This is demo program showing realuser and effective user")
  realUserId  := GetRealUserID()
  realUserName, _ := GetUserNameFromId(realUserId)
  effectiveUserId  := GetEffectiveUserID()
  effectiveUserName, _ := GetUserNameFromId(effectiveUserId)
  fmt.Printf("RealUserId: %d RealUserName %s\n", realUserId, realUserName)
  fmt.Printf("EffectiveUserId: %d EffectiveUserName %s\n", effectiveUserId, effectiveUserName)
}
```

此demo通过直接调用syscall，打印出程序的RealUser和EffectiveUser。

编译称为二进制: 

```
go build -o test_setuid test_setuid.go
```

我们先查看一下我本机当前用户，uid和gid为1000, name均为frankie

```
id
uid=1000(frankie) gid=1000(frankie) 组=1000(frankie), ...
```

我们运行一下，得到

```
RealUserId: 1000 RealUserName frankie
EffectiveUserId: 1000 EffectiveUserName frankie
```

符合预期：这也是最常见的情况，RealUserId和EffectiveUserId相同，都是当前用户。

我们开启setuid bit

```
chmod u+s test_setuid
```

在zsh下，以红色高亮标识开启此bit的binary 

<div align="center">
{% asset_img after_setuid.png %}
</div>

之前还是绿颜色的，而且注意execution bit由x变成s  

怎么测试以另一个用户运行test_setuid呢？正好系统里有root用户，我们先切换到root: 

```
sudo su -
```

得到root的prompt，再次运行test_setuid，得到

```
RealUserId: 0 RealUserName root
EffectiveUserId: 1000 EffectiveUserName frankie
```

看到变化了吗？RealUserId指示真正的用户变成了root，但是起作用的UserId(effectiveUserId)，是test_setuid，也就是file的owner，这也意味着，file access权限从root降低到普通用户frankie了，不能为所欲为了。这样实现了以另一个user的身份运行process，从结果来说，permission可能升高，也可能降低，以目标effectiveUser为准。  

看到这里，我们对如何使用setuid机制已经明了，就是简单的开启file的一个特殊bit即可，但说到底，这样，实现在不同user的权限之间切换，有何实际意义呢？我们来看几个实际的例子

## Usage Cases

### passwd

passwd用来更改当前用户的密码，直接运行，在其提示下先输入当前用户的密码，就可以修改新密码了。存储密码的文件位于 /etc/shadow，Linux将用户的密码存储于此处，passwd程序通过修改此文件的内容来达到修改密码的目的，查看此文件的permission bit：

```
-rw-r----- 1 root shadow 1.5K 11月 13 08:11 /etc/shadow
```

发现只有root用户才可以对其进行write操作，而我们运行passwd时，是以当前用户运行的，并没有sudo或者需要输入root密码之类的操作，实现了普通用户的进行修改只有root才能写的操作，怎么做到的呢？答案就是setuid，我们看下passwd程序： 

```
ll /usr/bin/passwd
# 得到
-rwsr-xr-x 1 root root 59K 1月  26  2018 /usr/bin/passwd
```

owner的s表明是开启了setuid bit，而其owner正是root，这样才会在运行时以root运行，修改/etc/shadow文件。

### sudo

sudo我们一定不陌生，最常用来执行只有root才有权限的操作，事实上，sudo可以以任何其他user身份运行process，不一定非得是root 

```
sudo -l        # List available commands.
sudo command   # Run command as root.
sudo -u user command  # Run command as user.
```

查看sudo程序的属性: 

```
ll /usr/bin/sudo
-rwsr-xr-x 1 root root 146K 1月  18  2018 /usr/bin/sudo
```

同样开启了setuid bit，owner为root，这样才会执行root专有的操作，sudo工作原理如下:  

*  Read and parse /etc/sudoers, look up the invoking user and its permissions,
*  Ask the invoking user for a password (this is usually the user's password, but can also be the target user's password or skipped as with NOPASSWD)
*  Create a child process in which it calls setuid() to change to the target user execute a shell or the command given as arguemnts in this child

Linux中，child process会继承parent process的RealUserId及EffectiveUserId。 这时子进程本该继承父进程sudo的RealUserId(普通用户)及EffectiveUserId(root)，但sudo会修改其RealUserId及EffectiveUserId为sudo命令行指定的目标用户。  

注意只有root才可以更改自己及子进程的EffectiveUserId并同时修改RealUserId，除此之外process不允许修改RealUserId。这样sudo就实现了完全以另一个user运行某个进程，完全的意思是RealUserId也会改变。

```
sudo ./test_setuid
# output
RealUserId: 0 RealUserName root
EffectiveUserId: 1000 EffectiveUserName frankie
```

可见RealUserId已经被改为root，但是EffectiveUserId由于setuid bit的关系，仍是普通用户frankie。

## Why saved-effective-id is useful

其实process除了RealUserId及EffectiveUserId之外，还有saved set-user-ID，为什么需要这个呢？

> Having a saved user id allows you to drop your privileges (by switching the effective uid to the real one) and then regain them (by switching the effective uid to the saved one) only when needed.

什么意思呢？有时进程需要的User转换需要多次往复转换：假设某程序在执行完root操作后，需要转换为普通用户，从root转换为普通用户是允许的，但是随后还会需要从普通用户转换为root，但这时setuid bit不起作用了，因为这不是首次运行，进程已经完全以普通用户身份运行了，不可能转换为root，怎么办呢？root转换为普通用户之前，会自动将当前的EffectiveUserId(即root)存储到saved-effective-id中，以保证随后可以转换回来。

APUE关于进程如何在RealUserId，EffectiveUserId，saved-effective-id的作用这方面讲解非常精辟，感兴趣的同学可以看 chap 8.11，这里放一张图大家感受一下： 

<div align="center">
{% asset_img summary.png %}
</div>


## Regarding Group

GroupId有与UserId完全对应的feature: 


*  setgid bit V.S. setuid
*  相对于process有的Real/Effective/Saved-SetUserId，process有Real/Effective/Saved-SetGroupId

因此上面针对User的讨论，可以完全对应到Group，设置setgid bit:

```
chmod g+s targetFile
```


## Setuid pitfalls

Setuid在很多情况，尤其设计到权限转换之处非常有用，但是如果不谨慎，及容易造成严重的安全漏洞，所以很多情况会被系统禁止，下面是作者在一个项目中踩过的坑：

*  bash/bash script无法设置setuid bit: Linux ignores the setuid¹ bit on all interpreted executables (i.e. executables starting with a #! line)
*  ssh不支持setuid。一个例子：parent process开启setuid，在子进程fork&exec child ssh，ssh并不会继承父进程的EffectiveUserId，因为ssh程序一旦发现EffectiveUserId与RealUserId不一致，便会将EffectiveUserId重新设置为RealUserId, 这样是有security concern。
*  无法设置一些重要环境变量，如LD_LIBRARY_PATH:  while executing setuid programs, $LD_LIBRARY_PATH is ignored，同样是security concern，防止因替换系统动态库，普通用户取得root权限。


## Conclusion

*  进程有三组user/group id，是进程动态改变权限的基础
    -  RealUserId: 用户登录名
    -  EffectiveUserId: 当前作用于file permission bit的用户Id
    -  Saved setuid: 用作权限提升及下降。
*  EffectiveUserId与file permission bit结合，决定进程对于特定文件的permission access
*  setuid通过修改进程EffectiveUserId为进程二进制的owner来实现以另一用户权限运行。
*  setuid仅作用于file，对folder无效
*  setuid机制可以很方便的改变程序的权限，但是同样也很容易称为系统安全漏洞

## Reference

Chap4.4, Chap8.11 Stevens, R., 2013. Advanced Programming in the UNIX Environment. 3rd ed. United States: Pearson Education, Inc.  
[GID, current, primary, supplementary, effective and real group IDs?](https://unix.stackexchange.com/questions/18198/gid-current-primary-supplementary-effective-and-real-group-ids)  
[How do the internals of sudo work](https://unix.stackexchange.com/questions/80344/how-do-the-internals-of-sudo-work)  
[System Administration Guide: Basic Administration](https://docs.oracle.com/cd/E19253-01/817-1985/userconcept-35906/index.html)  
[How to use special permissions: the setuid, setgid and sticky bits](https://linuxconfig.org/how-to-use-special-permissions-the-setuid-setgid-and-sticky-bits)  
[Switching user using sudo](http://researchhubs.com/post/computing/linux-cmd/sudo-command.html)  
[How does sudo really work?](https://unix.stackexchange.com/questions/126914/how-does-sudo-really-work)  
[Allow setuid on shell scripts](https://unix.stackexchange.com/questions/364/allow-setuid-on-shell-scripts)  
[c - Program can't load after setting the setuid bit on](https://www.google.com.au/search?q=setuid+LD_LIBRARY_PATH+ignored&oq=setuid+LD_LIBRARY_PATH+ignored&aqs=chrome..69i57.7027j0j1&client=ubuntu&sourceid=chrome&ie=UTF-8)

