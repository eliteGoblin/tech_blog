---
title: "Golang BDD入门: Ginkgo和Gomega实现(1)"
date: 2019-11-28 16:41:32
tags: [Golang, testing]
keywords: [Golang, BDD, testing]
description:
---


{% asset_img Ginkgo.png %}　


## Preface

在之前文章讲述过如何在Golang中进行TDD，相信大家也听说过BDD(Behavior Drive Development)， 和TDD仅差第一个单词: 由Test换成Behavior， 这个代表什么？ BDD和TDD看着很类似，都是test case first，然后实现可以通过满足test case的代码，那它和TDD关系如何，区别在哪？在Golang中怎么进行BDD实践呢？我们一一道来。

<!-- more -->

## BDD 简介

### BDD vs TDD

从它的名字Behavior Drive Development我们可以看出，这是一种指导Development方法，它的作用便是通过BDD，更容易构建出符合我们预期的软件。所谓我们的预期便是一个个的test case，到这里还是感觉和TDD很像，其实不然。

在TDD中，关注的核心点是function，即认为程序最基本单元是function, 其test case可认为等同与Unit test, TDD与unit test的区别是TDD强调测试和开发结合而成的工作流: 写test case --> 写代码 --> 通过测试，继续写更多测试，下一次循环。

而BDD相比TDD关注更高层的行为，而不是函数级的行为。也就是在BDD中，不会强调函数的功能正确，太底层了，这是unit test应该做的事情。BDD关注的是user story，即用户在特定场景， 与软件交互发生的行为，这个Behavior指的是高层模块(而非基本单元，函数)的行为。

如何区分BDD和TDD呢？可以这么理解： **TDD是给程序员的**，用来验证开发者的最基本模块的功能: 在什么输入，应该产生什么输出，保证实现的边界，健全性。而BDD，由于其test case描述的是更高级模块的行为，脱离了具体的实现，容易用自然语言来描述，**BDD是给产品经理的**，告诉他们系统的行为，请看下面的例子。

###  BDD示例: 购物车系统

我们在实现购物车系统，可以先和产品经理确定其行为，这些精细化的需求可以直接转化为test case: 

初始状态购物车为空(Given)
*  当添加1个商品A时(When)
  +  购物车商品列表显示A，个数为1 (Then)
  +  购物车商品类型列表显示A
  +  购物车总价显示A的价格
*  当添加2个商品A时(When)
  +  购物车商品总列表显示A，个数为2 (Then)
  +  购物车商品类型列表显示A
  +  购物车总价显示A的价格 x 2
...

这里test case完全不涉及怎么实现，怎么实现是程序员的事情， 需要与功能要求解耦。有了test case, 便不太可能出现沟通的理解歧义，同时BDD鼓励程序员实现过程中，不断验证当前实现是否满足约定，直到通过。

### BDD 在测试架构中的位置

可见BDD测试因为面向更高层的功能，不强调系统边界， 健全性，凭BDD本身无法保证系统的各种问题，因此BDD一般要辅以TDD，即底层用单元测试保证代码/实现的质量，在高层/抽象层确保系统的行为符合预期。我个人倾向金字塔模型，不同种类的测试分层，底层的test case应该最多，越高级/抽象的测试，越少：

{% asset_img pyramid.png %}

BDD应处于components一层。

### BDD 更多思考

摘录一些关于BDD回答[^1]，加深理解: 

> BDD，不是跟TDD一个层级的，B是说代码的行为，或许比单元测试高那么一点点吧，主要是跟ATDD（接收测试驱动开发）、SBE（实例化需求）等实践一并提及的，因为他们都是对应到传统测试理论里面，高于单元和模块测试，从功能测试、集成到系统、性能等这些高级别测试的范围。所以说，TDD、BDD根本不是一个层面的东西，解决的是不同的问题。但BDD、ATDD、SBE基本上都认为TDD是基础，也即是说，他们主张做BDD、ATDD、SBE必做TDD。但反之则未必。

> BDD的核心价值是体现在正确的对系统行为进行设计，所以它并非一种行之有效的测试方法。它强调的是系统最终的实现与用户期望的行为是一致的、验证代码实现是否符合设计目标。但是它本身并不强调对系统功能、性能以及边界值等的健全性做保证，无法像完整的测试一样发现系统的各种问题。但BDD倡导的用简洁的自然语言描述系统行为的理念，可以明确的根据设计产生测试，并保障测试用例的质量。


### GWT 语法

BDD的第一步便是约定test case，需要贴近自然语言，方便非开发人员理解。一般采用Give-When-Then的形式来组织test case，称为GWT语法。

对照前一节提到的购物车test case：Given表示测试初始条件，When表示与系统交互action的发生，Then描述系统的行为(结果)。


只要对test case达成共识，并用GWT描述，下一步便可以开始BDD实践了。


## BDD in golang

在实现时，我们需要将GWT组织test case记录翻译为测试代码，并运行测试，通过一系列的assertion来检查实现是否符合test case的预期。我们当然可以直接通过golang的内置testing来实现，这里我们采用Ginkgo+gomega来组织test case，让我们的测试代码读起来非常接近自然语言。

为什么要用两者结合呢？ginkgo实现了test case的组织(以便让其读起来像自然语言)，并加入了其他一些便利功能: 初始化，后续处理，异步等等。而Gomega设计的目的便是与ginkgo一起工作，实现易读的assertion(ginkgo中成为match)功能，其官方也提到: 

> Gomega is ginkgo's preferred matcher library


后篇blog我们来看如何用Ginkgo和Gomega进行BDD test。




[^1]: [TDD 与 BDD 仅仅是语言描述上的区别么？](https://www.zhihu.com/question/20161970)