---
title: Security, the big picture - Part I
date: 2018-11-20 17:05:37
tags: [Security golang HTTPS]
keywords:
description:
---

{% asset_img enigma.jpg %}

## Preface

安全是计算机系统中独立又特别部分。尤其是日常编程任务中业务代码实现较多时，往往不会涉及到安全，导致我们对安全仅有模糊的概念，but why should we care?  

安全的重要性不言而喻，尤其是在系统架构层次，当你需要完整的设计一套商用系统时，除了功能实现，性能，高可用等方面的考虑，安全也显得尤为重要。即使是仅实现业务代码，必然会碰到一些场合需要安全相关知识，特别是以SSH，HTTPS为代表的安全技术如此流行的今天，RSA/DSA，SSL/TLS，椭圆曲线，公钥私钥，签名，证书，X509，pem文件，等等等等，一系列的术语名词让人目不暇接，如坠云雾。 零星的关键词搜索带来了更多的术语名词，更多的疑问，使人沮丧。 

作者本人对安全入门的挫折深有体会，很久都搞不清楚SSH，HTTPS工作原理；分不清SSL/SSH区别。但是随着前阵子一个安全相关的项目的完结，基本上以程序员视角，建立起安全的big picture，从而弄清楚了对于安全的很多疑问。再一次深感建立big picture的重要性，见[^1]。因此本篇博客的目的是通过回答一系列的Why，帮助读者建立起安全的框架，知其所以然，是真正理解知识的开始；在big picture建立之后，通过实际操作，编程来体会每一部分，才会在脑中留下知识的深深烙印。这是一劳永逸习得新知识的最佳实践。可是长期来看，最为高效的学习方式。

<!-- more -->

## Ultimate problems security want to solve

简而言之，三防：防偷窥，防篡改，防冒充。如何理解呢，我们假设我们当前处在一个没有任何安全机制的计算机系统中，Alice和Bob通过计算机发message来实现通信，还有一个心怀叵测的用户，Eve，来源于发音相似的,表示窃听的单词：eavesdrop。

### Confidentiality

防偷窥很好理解，没有安全机制意味着所有消息都是以plain text方式发送，如在TCP/IP网络中，Eve可以截获，看到Alice和Bob互相发送的message。如何解决呢？加密(encryption)。何为加密呢？发送方Alice将明文(二进制串)按照某种算法转为无法明白其意思的二进制串(加密，encrypt)；接收方Bob按照与Alice使用算法相逆的算法，将无法解读的二进制串转换为明文(解密，decrypt)。算法可以是通用的，那是否意味着只要采用知道了这种算法就可以破译别人，同样采用此加密算法的message呢，不行的。

要了解上述过程的原理，我们来看一个绝佳的例子，凯撒密码[^2]。理解了这个古老的例子，我们能从其中找到现代密码学基础术语的对应，对理解安全的基本原理非常有帮助。类比是学习的一个非常强大的工具。

{% asset_img caesar.png %}

信息传递对于军队性命攸关，但同时信息又需要保密。直接由口头下达的命令便类似没有安全机制的plain text。凯撒大帝的解决办法是，对命令*加密*。

>  恺撒密码: 明文中的所有字母都在字母表上向後（或向前）按照一个固定数目进行偏移後被替换成密文

这里凯撒密码对应上文：

*  命令是明文
*  按固定数目进行偏移是加密算法
*  逆向的按固定数目进行偏移是解密算法
*  明文按照加密算法，加密得到密文

回到上面的问题，如果别人知道了加密，解密算法，即*按固定数目进行偏移*，是否能破译密文呢？不能，因为还不知道偏移是多少。这个偏移，是密码学中key的前身。即Eve知道了算法，但是我们同时用一个特别的数字作为算法的输入，那么Eve就无法得知如何解密。

基本道理了解了，你已经入门一半了。安全学一个令人望而生畏的部分是其术语，而多个术语代表同一概念的情况让其尤为严重。因此来一波术语对照：

> Plaintext: Decrypted or unencrypted data (it doesn't have to be text only)  
> Cipher: Algorithms for encoding text to make it unreadable to voyeurs  
> Encrypt: Scrambling data to make it unrecognizable  
> Decrypt: Unscrambling data to its original format  
> Key: A complex sequence of alpha-numeric characters, produced by the algorithm, that allows you to scramble and unscramble data  
> Ciphertext: Data that has been encrypted  

在凯撒密码中，key便是偏移的数字，这就是我们耳熟能详的秘钥。

再举一个栗子：

二战时德国科学家，发明了对电报加密的机器，Enigma(见文章开始的配图)：

> 对于潜艇作战，尤其是德国海军的“狼群”战术来说，无线电通讯是潜艇在海上活动获取信息通报情况的最重要手段，而Enigma则是关乎整个无线电通讯安全的核心设备。在1942年之前，装备了Enigma的德国潜艇部队击沉了盟军舰船1000余艘，破译Enigma是战胜德国海军潜艇的关键。
> Enigma采用复式字母替换加密方法，利用键盘、转子、跳线、反射板、显示器进行对称加密/解密。简单讲，按键盘上任意一个字母，该字母会经过键盘-转子-跳线-反射板-显示器而被转换为另一个字母。每转换一次，转子会转动一格。破解Enigma的难度在于不知道当前密钥，包括转子的初始位置和跳线设置。

Enigma机器是用硬件(键盘、转子、跳线等)实现的加密算法，Enigma机器被盗取后，即使通过研究内部结构，弄清楚加密原理，也因无法知道秘钥（当前拨号盘数字组合），而无法破译密文。

是不是只要key(密钥)未知，就一定无法破译呢？当然不是。凯撒密码可以通过穷举26种字母替换组合，并根据语法判断得到明文。Enigma最终也被图灵完美破解。对于现今加密算法而言，只要密钥足够长，破解(以目前计算能力)几乎不可能。

简而言之，加密可以表示如下：

> Encryption: CipherText = Encrypt(PlainText, KeyOfEncrytion)  
> Decryption: PlainText = Decrypt(CipherText, keyOfDecryption)

KeyOfEncrytion和keyOfDecryption是相同的key么？不一定。

key的本质是一串随机的二进制序列，当今安全体系下，有两大类解决方案：

*  对称key
*  非对称key

如何理解？在加密，解密过程中，KeyOfEncrytion和keyOfDecryption是完全一样的，就叫做对称，symmetric。与之相反，加密解密使用的不同的key，称为非对称，asymmetric。

*  非对称涉及到一对key，keyA和keyB
*  非对称key有如下特性：
    -  单个key既可以用来加密，又可以用来解密。
    -  keyA加密后，必须用keyB解密
    -  keyB加密后，必须用keyA解密
    -  也就是采用同种算法情况下，互相能解密对方加密的cipher text，**并不是某个key专用于加密，另一个key专用于解密**
*  非对称key是公钥体系(Public Key Infrastructure)的基础
*  常见的对称加密算法(采用对称key的算法): DES、3DES、AES、Blowfish、IDEA、RC5、RC6
*  常见非对称加密算法: RSA、DSA、ECC
*  算法一定的情况，安全性取决于key的长度。

有了key的话，如何解决Alice和Bob通信不被偷窥的问题呢？

对称加密情况:

*  Alice和Bob都持有相同的key。Alice发送message前，先用事先商量好的加密算法，如3DES。得到ciphertext，发送给Bob
*  Bob接受到，采用3DES算法及同样的key，能解密ciphertext。
*  Eve因为不知道算法，尤其是key而无法破译ciphertext，达到防止被人*偷窥*

非对称加密情况：

*  Alice持有keyA，Bob持有keyB。
*  Alice使用非对称个加密算法，如RSA，以及keyA来加密message
*  Bob收到message，用keyB解密

这么看起来大同小异，为什么这么麻烦，使用同一个key不就完了？  

不能，一个栗子：对称秘钥需要各不相同，假如taobao为了防止客户交易被人窥视，采用对称秘钥加密，来一个人taobao服务器就得存储与此人相同的秘钥，来一亿个人就。。。  

我们将在另外的文章中比较对称和非对称加密各自特点及适用场合。引出重要概念PKI，就是我们总提到的公钥/私钥体系，这是理解SSH，HTTPS技术的基础。

### Integrity

解决了防止被人偷窥的问题，另一个问题是，我们如何保证数据不被人篡改？即到达Bob的数据一定与Alice当时发送的数据**完全一致**呢？虽然数据被加密称为一系列无法解读的二进制串，但数据在传输过程中，仍可能有意或者无意被改变。比如传输时数据出错，或者被恶意用户Eve截获，经过修改后才发送给Bob。  

这里值得一提的是，即使是我们一般意义上能可靠传输TCP，并不能完全保证数据没有发生变化，见[^3]。TCP的做法是采用计算并对比数据的checksum，16-bit ones-complement sum来判断数据是否发生变化，但仍会存在修改过后的数据checksum没变的情况。    

我们采用类似的思路，将原始数据**浓缩**成一串独一无二的二进制数据，称为digest，这个数据就像数据的指纹一样，任何两个数据的digest都不一样。 

Alice在发送数据同时，附上用MD5计算出来的digest，只要Bob在收到数据时，用MD5计算digest，比较Alic附带的digest及Bob计算得到的digest，一致表明数据没有发生改变。

> The minimum checksum you'd need for ensuring flawless data transfer is the MD5 value of the data. Of course anything better than that (SHA-1, SHA-256, SHA-384, SHA-512, Whirlpool, and so on) will work even better, yet MD5 is sufficient

注意digest需以加密为前提，否者恶意用户Eve先篡改Alice的消息，再重新按照MD5计算一个被篡改信息的digest，Bob变无从发现消息被篡改了。

## Conclusion

本文及后续Part II，旨在简明的介绍安全面临的基本问题，以及解决思路。并不涉及具体代码。本文主要讲了如何解决数据传输过程中的偷窥，篡改问题。Part II来解决伪造身份的问题。  

## Stay tuned

以下是接下来本系列文章构想，因目前有更优先级高的事情，以及这部分给我比较大的压力和时间消耗，会搁置一阵自再完成

*  big picture Part II: 防冒充原理, signature。certificate
*  Golang代码实现如何生成对称key，加密，解密，巩固知识。
*  big picture Part III: PKI体系，公钥私钥
*  实际操作公私钥生成，利用公私钥加密解密。
*  安全应用：SSH简介
*  安全应用：SSL/TLS HTTPS简介
*  综合例子，没最终确定：
    -  申请证书等 on AWS，用nginx搭建https静态网站
    -  golang实现HTTPS通信，server/client相互认证


[^1]: 见[如何高效学习](https://keeganlee.me/post/full-stack/20170909)快速学习四步法
[^2]: [凱撒密碼](https://zh.wikipedia.org/wiki/%E5%87%B1%E6%92%92%E5%AF%86%E7%A2%BC)
[^3]: [Can a TCP checksum fail to detect an error?](https://stackoverflow.com/questions/3830206/can-a-tcp-checksum-fail-to-detect-an-error-if-yes-how-is-this-dealt-with)  





