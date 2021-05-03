---
title: Clean Architecure in Go (2)
date: 2019-04-03 13:35:12
tags: [Architecure golang]
keywords:
description:
---

## Preface

接上篇 Clean Architecure in Go (1)，本文将会包含如下内容：

*  Clean Architecture层次划分依据
*  实现Dependency Inversion: Dependency Injection
*  采用Clean Architecture的prject layout示例
*  A simple demo: Clean Architecture实践

<div align="center">
{% asset_img CleanArchitecture.jpg %}
</div>

<!-- more -->

## Analysis of Clean Architecture by layer

### History of Clean Architecture

早在Uncle Bob提出Clean Architecture之前，已经有很多类似的提案：Hexagonal, Onion, DCI, BCE Architecture，所有这些architecture提案解决如下问题，和第一篇文章提到的架构痛点本质上是一回事：

> Independent of Frameworks. The architecture does not depend on the existence of some library of feature laden software. This allows you to use such frameworks as tools, rather than having to cram your system into their limited constraints.  
> Testable. The business rules can be tested without the UI, Database, Web Server, or any other external element.   
> Independent of UI. The UI can change easily, without changing the rest of the system. A Web UI could be replaced with a console UI, for example, without changing the business rules.   
> Independent of Database. You can swap out Oracle or SQL Server, for Mongo, BigTable, CouchDB, or something else. Your business rules are not bound to the database.  
> Independent of any external agency. In fact your business rules simply don’t know anything at all about the outside world.  

Uncle Bob正是在这些架构的基础上，提出了简单易行的Clean Architecture

### Why Layered?

UNIX的架构图和Clean Architecture极为相似，都是layered approach，且外层为细节，内层为核心(kernel)，外层依赖核心的实现，而非相反。  

<div align="center">
{% asset_img unix_arch.png %}
</div>

同理还有TCP/IP协议栈，分层使得高层不必依赖底层实现的细节，只要底层实现了高层需要的接口，无论底层如何实现，高层的逻辑也不必改动。

### The Dependency Rule

正如我们之前提到的，痛点的根源是错误的dependency，Clean Architecture解决之道是：*只允许外层依赖内层*，dependency箭头由外向内，由细节/底层到抽象/高层，绝不允许相反的情况。

每个同心圆代表软件的一层，最外层是具体细节，是mechanism；最内层是最核心的business logic，是policy，是最不容易变化的部分。

我们代码中，有时候会见到我们的业务逻辑模块，输入直接是database的row对应的struct，这就是Anti-pattern，因为database schema是实现细节，不应该影响我们业务逻辑，比如我们即使改动数据库表结构，甚至更换了数据库类型，业务逻辑都不应该有改动。

对我们Golang项目的具体指导意见应该是: 

> Entity和Usecase层不应该有任何依赖，不应该import任何第三方package

### Entity Layer

此层应实现：Use Case中的entity，或者UML图中actor。

用Golang术语来说，可以是核心数据struct及附带的方法，或者一系列的struct和functions。

实现的是最核心，最抽象的商业逻辑，也一般是最不容易改动的部分。可以被其他层直接引用。

### Usecase Layer

系统中所有Usecase实现，控制从entity流入，流出数据数据的过程。此层的改变一般不应引起entity的改动。

我们的设计应该以Usecase为中心，Usecase Diagram绘制出来，Usecase以及entity就都有了。


<div align="center">
{% asset_img ub.png %}
</div>

关于什么是Usecase，请参考[这里](http://www.plainionist.net/Implementing-Clean-Architecture-UseCases/)

Usecase可以，类似面向对象的inheritance及combination，对应Usecase术语是include及extend

### Interface Adapters Layer

此层一般不会有太多代码，主要是起到在配接usecase层及最外层：转换两边的数据结构，使usecase依赖的数据结构与detail的数据结构解耦。

Usecase层及最外detail层因此可以选择最适合各自层实现的数据结构，用本层简单转换一下就行。

一个例子是，之前提到的anti-pattern，最外层数据结构直接反映了database的schema，经由本层转换，转换称为usecase层需要的数据结构。

另一个例子是SQL，本质上我们可以认为其功能也是转换数据：从数据库到内存的golang数据结构，因此SQL应该在本层实现。

Uncle Bob同时提到了MVC也全部应该在此层实现。

### Frameworks and Drivers Layer

总之本层存放各种细节，framework，database，web等初始化。一般来说代码也不会太多。

### A Clean Architecture Layout Example

<div align="center">
{% asset_img clean_arch_structure.png %}
</div>

## Dependency Injection

在前一篇博客介绍Dependency Inversion的部分，我们已经知道invert dependency flow的重要性，dependency injection便是其最常用的实现方法。

例如 A object需要call B，A是高层，B是细节，可以这样: 

```
// pseudo code in golang
// package A follows
package A
type A struct {
    b BInterface
}
type BInterface interface{
    ...
}
// package B follows
// struct B implement BInterface
package B
import "A"
type B struct{
}
b := B{
    ...
}
a := &A{
    b : b, // inject object b into object a
}
a.do() // do function will call a.b 
```

我们可以看到，a的实现并不需要创建object b，由外部模块，一般是main，来初始化b，并且*注入*到a中，从而解除A对B的依赖，实现了dependency inversion。

## Demo: User Registration

千言万语还是得归结于一个concrete sample，本节我们来通过实现一个简单的use case来了解如何应用clean architecture。

这个有趣的例子请参见[Dependency Injection In Go](https://medium.com/full-stack-tips/dependency-injection-in-go-99b09e2cc480)

通过这里例子，我们可以看到一个最简单的用户注册的usecase，如何采用clean architecture的方法论来实现。用户注册本质上可以抽象为：

*  validate新用户信息
*  存储新用户信息

validate可以由我们的business来自己定义，如不允许用163邮箱用户，不允许座机只允许手机号码注册等，这个可以放在entity层实现。但是存储用户信息是细节，有可能存于mysql，或者redis，本地文件甚至只存于本机内存。存储方法应该由最外层来实现。

## Pros and Cons

一图流

<div align="center">
{% asset_img review.png %}
</div>

关于cons： 

*  假如我们逻辑极为简单，也就没必要采用Clean Architecture的分层设计。
*  复用代码会导致不同模块耦合，所以解耦势必会产生一些"duplicate code"，比如usecase层会声明一个与database的schema(即表结构)非常类似的数据结构，看起来两个struct很像duplicate code，其实从Clean Architecture本质上来说，两层实现的模型本就不一样，虽然从代码上看着类似，但也不能复用，避免错误耦合。


## Conclusion

个人对Clean Architecture思想是非常肯定及推崇的，UncleBob在推广它也花了不少功夫，有兴趣可以看他的视频演讲：[这里](https://vimeo.com/43612849)和[这里](https://www.youtube.com/watch?v=Nsjsiz2A9mg)  

然而不得不承认的是，UncleBob提出其已经过了7年，我并没有看到大规模的被采用的迹象，至少顶级的Golang Opensource Project的layout还是比较粗放随意的，不知道是曲高和寡还是当下更流行的小步快跑策略导致人们对传统软件架构理论不屑一顾？个人比较赞同在他处看到的一句话：互联网相对简单的业务逻辑(与传统软件相比)导致我们往往不需要复杂的架构也能解决问题。

## More reading and reference

[Java example](https://github.com/mattia-battiston/clean-architecture-example#why-clean-architecture)  
[Slide of realworld Java clean arch](https://www.slideshare.net/mattiabattiston/real-life-clean-architecture-61242830)  
[.Net Clean Architecture](https://medium.com/@stephanhoekstra/clean-architecture-in-net-8eed6c224c50)  
[The "Implementing Clean Architecture" series](http://www.plainionist.net/Implementing-Clean-Architecture-UseCases)  
[Screaming Architecture](https://blog.cleancoder.com/uncle-bob/2011/09/30/Screaming-Architecture.html)   
[EBI Architecture](https://herbertograca.com/2017/08/24/ebi-architecture/)

