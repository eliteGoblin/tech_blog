---
title: "Golang BDD入门: Ginkgo和Gomega实现(2)"
date: 2020-04-11 11:45:15
tags: [Golang, testing]
keywords: [Golang, BDD, testing]
description:
---

![ginkgo.jpg](https://raw.githubusercontent.com/eliteGoblin/images/master/blog/img/picgo/ginkgo.jpg)

## Preface

接[上篇](https://www.franksun.cn/2019/11/28/BDD-1/), BDD解决的问题及与TDD的区别; 以及`Given-When-Then`语法. 

继续看如何用ginkgo框架实现Behaviour Test.

<!-- more -->

## 初识Ginkgo

ginkgo, [github](https://github.com/onsi/ginkgo)上3800+ stars, 是最流行的BDD框架. 

回顾上篇的购物车例子，先直观感受ginkgo测试代码: 
 
>初始状态购物车为空(Given)

> 当添加1个商品A时(When)
    购物车商品列表个数为1 (Then)
    购物车商品列表不重复商品个数为1 (Then)
    购物车总价显示A的价格 (Then)


> 当添加1个商品A, 2个B时(When)
    购物车商品个数为3 (Then)
    购物车商品列表不重复商品个数为2 (Then)
    购物车总价显示价格A+2*B(Then)

```golang
var _ = Describe("GinkoCart", func() {
	var (
		cart *ginko_cart.Cart
		err error
	)
	Context(`Start with empty cart`, func(){
		BeforeEach(func() {
			cart = ginko_cart.NewCart()
		})
		When(`Add One A item to cart`, func(){
			BeforeEach(func() {
				err = cart.AddItem(ginko_cart.Item{
					Name: "A",
					Price: 3.99,
					Qty  : 1,
				})
			})
			It(`Should no error`, func() {
				Expect(err).To(BeNil())
			})
			It(`Should display items count as 1`, func() {
				Expect(cart.TotalItems()).To(Equal(1))
			})
			It(`Should display items count as 1`, func() {
				Expect(cart.TotalUniqueItems()).To(Equal(1))
			})
			It(`Should display items total price as A`, func() {
				Expect(cart.TotalPrice()).To(Equal(3.99))
			})
		})
		When(`Add One A and Two B item to cart`, func(){
			BeforeEach(func() {
				err = cart.AddItem(ginko_cart.Item{
					Name: "A",
					Price: 3.99,
					Qty  : 1,
				})
				Expect(err).To(BeNil())
				err = cart.AddItem(ginko_cart.Item{
					Name: "B",
					Price: 12.99,
					Qty  : 2,
				})
				Expect(err).To(BeNil())
			})
			It(`Should display items count as 3`, func() {
				Expect(cart.TotalItems()).To(Equal(3))
			})
			It(`Should display unique items count as 2`, func() {
				Expect(cart.TotalUniqueItems()).To(Equal(2))
			})
			It(`Should display items total price as A+B`, func() {
				Expect(cart.TotalPrice()).To(Equal(3.99+2*12.99))
			})
		})
	})
})
```

Test code可读性很高: key word和annotation可以与自然语言一一转换. 

ginkgo与Javascript BDD框架高度相似, 可见两者关键词几乎完全一致: 

```javascript
describe("A spec", function() {
  it("is just a function, so it can contain any code", function() {
    var foo = 0;
    foo += 1;

    expect(foo).toEqual(1);
  });

  it("can have more than one expectation", function() {
    var foo = 0;
    foo += 1;

    expect(foo).toEqual(1);
    expect(true).toEqual(true);
  });
});
```

Ginkgo的[document](https://onsi.github.io/ginkgo/)是很好的start up guide.

## 第一个Ginkgo例子

Ginkgo依托于golang原生testing框架, 即可用`go test ./...`运行, 也可通过ginkgo binary(安装`go install github.com/onsi/ginkgo`). 封装了ginko测试框架的各种feature, 实际中我用的很少, 仅用来初始化测试代码. 

本节通过简单购物车例子了解如何写BDD测试代码, 完整的例子代码在[github](https://github.com/eliteGoblin/code_4_blog/tree/master/ginkgo_cart)

### 初始化

首先进入待测试package:

```bash
cd code_4_blog/ginkgo_cart
```

初始化
```
ginkgo bootstrap
```

生成以suite_test.go文件, 将ginko嵌入testing, 用`go test ./...`可运行Ginkgo测试代码. 

以上生成了新test suite, 接下来向suite添加测试specs, 生成ginkgo_cart package测试文件

```bash
ginkgo generate ginkgo_cart
```

### 运行

生成`ginko_cart_test.go`, 注意测试文件在`ginko_cart_test`package, 需import package`ginko_cart`. 目的是: BDD层级高于Unit test, 不应了解package内部实现, 测试package外部接口即可. 

编写测试代码，运行: `go test ./...`, 

![20200412152437.png](https://raw.githubusercontent.com/eliteGoblin/images/master/blog/img/picgo/20200412152437.png)


## Ginkgo关键词

Ginkgo测试代码骨架由一系列关键词关联的闭包组成, 常用key word有

*  Describe/Context/When: 测试逻辑块
*  BeforeEach/AfterEach/JustBeforeEach/JustAfterEach: 初始化测试用例块
*  It: 单一Spec, 测试case

Key word的声明均为传入body函数, 如Describe
```golang
Describe(text string, body func()) bool
```

Sample代码片段: 以分析执行顺序 

```golang
var _ = Describe("Nest Test Demo", func() {
	Context("MyTest level1", func() {
		BeforeEach(func() {
			fmt.Println("beforeEach level 1")
		})
		It("spec 3-1 in level1", func(){
			fmt.Println("sepc on level 1")
		})
		Context("MyTest level2", func() {
			BeforeEach(func() {
				fmt.Println("beforeEach level 2")
			})
			Context("MyTest level3", func() {
				BeforeEach(func() {
					fmt.Println("beforeEach level 3")
				})
				It("spec 3-1 in level3", func() {
					fmt.Println("A simple spec in level 3")
				})
				It("3-2 in level3", func() {
					fmt.Println("A simple spec in level 3")
				})
			})
		})
	})
})
```

### Describe, Context, When

三者被称为Container: 对Ginkgo均属同类节, 仅名称不一样. 

一般Describe用于最顶层: 描述完整的测试场景; 包含Context/When, Context/When本身可以嵌套包含下级Context/When.

Describe, Context, When组织成Tree结构: Describe是root, Context和When是普通TreeNode. 

三者可以包含的节点，除了自身，还包括其他Key word节点: BeforeEach, JustBeforeEach, It. 

测试代码逻辑应包裹在BeforeEach, JustBeforeEach, It中，不应直接在Container node实现. 

### It

Ginko执行以It的基本单元: 以定义的顺序执行(It数即为Ginkgo中的Spec数). 示例定义三个It node, 处于不同层次. 执行顺序为: `It 1-1`, `It 3-1`, `It 3-2`. 

It一般包含Assertion逻辑: `Exect(...)`, 即最终的测试结果和预期的比较. 

测试执行逻辑实现于BeforeEach, JustBeforeEach中.   

### BeforeEach, JustBeforeEach

`BeforeEach`声明于Container节点内部, container node每个child执行前都会执行`BeforeEach`. 一般用来Setup test env: 声明测试用变量, 初始化

`JustBeforeEach`很类似, 区别是永远执行于`BeforeEach`之后: 等从root到It node所有`BeforeEach`执行完;才再从root到It node执行所有`JustBeforeEach`; 一般实现测试执行逻辑: 如request HTTP, 添加商品到购物车. 总之是得出输出，以便`It` node与expect比较.

### Demo code 分析

示例各种节点内部组成为Tree:

![20200412194916.png](https://raw.githubusercontent.com/eliteGoblin/images/master/blog/img/picgo/20200412194916.png)

运行示例得到输出: 

```
beforeEach level 1
sepc 1-1 on level 1
•beforeEach level 1
beforeEach level 2
beforeEach level 3
Spec 3-1 in level 3
•beforeEach level 1
beforeEach level 2
beforeEach level 3
Spec 3-2 in level 3
```

可见: 

*  是以各`It` node定义顺序执行
*  每个`It`执行前，走了从root到`It`的path: 顺序执行各context node的`BeforeEach`函数

为什么是层次结构呢? `BeforeEach`实现本层Context environment setup, 本层测试逻辑出现分支: 有了Context子节点, 次层的`BeforeEach`定制次层的environment, 并再次分支: 再继续延伸出子Context...

## It 与Matcher

购物车demo中: 其中一个It

```golang
Expect(cart.TotalItems()).To(Equal(3))
```
这种自然语言风格的assertion是由Ginkgo配套的Gomega实现的: expect返回封装了测试输出值的Assertion:

```golang
func Expect(actual interface{}, extra ...interface{}) Assertion
```

Assertion是interface, 简化版本(为语义通顺，还包含几个类似function):

```golang
type Assertion interface {
	To(matcher types.GomegaMatcher, optionalDescription ...interface{}) bool
	ToNot(matcher types.GomegaMatcher, optionalDescription ...interface{}) bool
}
```

`To`接收`GomegaMatcher`, 其封装了Expect value: Equal调用了Ginkgo的EqualMatcher. 

```golang
func Equal(expected interface{}) types.GomegaMatcher {
	return &matchers.EqualMatcher{
		Expected: expected,
	}
}
```

加上Assertion封装了实际value, 两者的比较可得出结论.而`ToNot`是`To`的相反情况. 

如果想比较自定义的复杂类型: 可实现GomegaMatcher:

```golang
type GomegaMatcher interface {
	Match(actual interface{}) (success bool, err error)
	FailureMessage(actual interface{}) (message string)
	NegatedFailureMessage(actual interface{}) (message string)
}
```

## 其他常用Feature

Focus:

仅执行特定Node及之下的It: 在keyword之前加`F`: `FContext`, `FIt`, 但会使`go test`fail(返回 1), CI集成Ginkgo需注意.

Pending

与Focus相反: 不执行特定Node及之下的It.  在keyword之前加`X`.但默认不会使`go test` fail(若想让其fail, 加 --failOnPending)

Skip:

根据代码runtime结果决定是否跳过某It(Pending是编译时): 

```
It("spec 1-1 in level1", func(){
    if somecondition {
        Skip("special condition wasn't met")
    }
    fmt.Println("sepc 1-1 on level 1")
})
```

Skip仅能置于It之下，否则会Panic.

Eventually

测试异步逻辑: 如发送请求到队列, 需持续polling. 在Gomega实现: 

```golang
Eventually(func() []int {
    return thing.SliceImMonitoring
}, TIMEOUT, POLLING_INTERVAL).Should(HaveLen(2))
```

TIMTOUT为总超时时间, 默认１s;POLLING_INTERVAL为每次polling间隔, 默认10ms.

Ginkgo还支持benchmark及run in parallel, 可参考[Ginkgo doc](https://onsi.github.io/ginkgo/#parallel-specs)

祝大家BDD愉快！