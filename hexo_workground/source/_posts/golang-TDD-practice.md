---
title: golang TDD实践
date: 2018-01-31 19:27:06
tags: [TDD,testing,golang]
keywords: golang testing TDD test-driven
description: 
---

#### Preface 

<div align="center">
{% asset_img review.jpg %}
</div>

如果我们每天的工作就是喝着咖啡，嚼着零食，惬意的写写代码，按时回家。对写好的代码充满信心，不需要时刻担心自己的代码线上是否会出问题，不担心局部小小改动会以意想不到的方式引起故障；让别人一提起你都会说：这家伙的代码，No Problem。相信会让我们的人生幸福指数提高不少，然而现实情况却not even close，让人苦恼。
TDD(Test Driven Development)在实践中被证明是一种很有效的软件设计/开发方法，它不一定会让我们的梦想全部成真，但起码，带来了希望。

本文是作者在golang中学习TDD方法的一个总结：首先概括了TDD要素，然后讨论如何用基于golang的testing框架实现基本的TDD，并给出了例子：如何在解决leetcode问题时应用TDD

<!-- more -->

#### TDD 简介

##### TDD三定律 

采用TDD方法开发软件的流程很有特色，:

1.  在编写不能通过的单元测试前，不能编写生产代码
2.  只编写刚好无法通过的单元测试，不能编译也算不通过
3.  只编写刚好能通过当前失败case的生产代码

##### TDD工作循环

我们的写代码过程其实就是一次次的小循环，测试代码每次只比生产代码早写数分钟. 比如我们需要写一个模块: MakePizza(),对应三定律的具体开发步骤是

1. 　先写单元测试代码
```golang
func TestMakePizza() {
    TestGetOven()
    TestMakeDough()
    TestPreparePizza()
    TestBake()
    TestGetPizza()
}
```
这时测试是FAIL的，因为没有生产代码，编译不通过  
2. 　修改代码使单元测试刚好能PASS，需要补其单元测试代码需要的生产代码的各个函数
```golang
    // TestGetOven 需要
    func TestGetOven() {
        if getOven() != true {
            Error()...
        }
    }
    func getOven() bool {
        return true
    }
```
为使得测试刚好通过，我们需要根据此test case写一些能让测试成功的临时代码
3.  审视代码，是否有重复，是否高内聚，低耦合；重构之，跳到１

##### TDD 总结

*  If it's worth building, it's worth testing
*  TDD本质是design activity：如何写出可测试的代码，一般就意味着高内聚低耦合
*  TDD顺带的好处便是验证代码(单元测试)和丰富的文档（独立的test case可以被认为是模块使用说明）
*  在TDD过程中形成的测试集合可作为模块的requirement和specification，可作为其重要组成部分，但不可能代替全部文档
*  每次一小步：一点单元测试代码，一点代码，在不断迭代中构建强壮的系统
*  不追求完美，依照对系统的重要程度设计测试代码
*  TDD长于细微之处的说明和验证，不擅长处理类似整体设计的大问题，后者一般要配合AMDD方法，即agile model-driven development
*  Test Cases应由开发人员像维护生产代码一样维护，使之尽可能保持整洁，消灭重复

#### golang原生testing框架

##### testing框架要点

Golang内部集成了轻量级的测试框架，要点如下

*  测试代码写在与生产代码一个package内(即同一folder下)，以 **xxx_test.go** 命名
*  测试代码需要 import "testing"，包含两个基础功能测试和benchmark
*  在package 执行 go test 即可自动执行此包下的所有测试

##### 单元测试

*  功能测试: 核心对象是 testing.T，单元测试以函数为基本单元，go test其实实在挨个执行写好的单元测试函数
    -  单元函数命名必须为　TestXxx(t *testing.T)，Test后必须接一个大写字母
    ```golang
    func TestMakePizza(t *testing.T) {...}
    func Testgetoven(t *testing.T) {...}   // 不会执行，应为TestGetOven
    ```
    -  用testing.T的Error和Fail系列方法指示失败
    ```golang
    func TestGetOven(t *testing.T) {
        if oven, err := getOven(); err != nil {
            t.Errorf("err %s", err.Error())
        }
    }
    ```
    -  进行table-driven test:
    ```golang
    var testCases = []struct {
        param1 string
        param2 string
        out string
    } {
        {"11", "12", "out1"},
        ...
    }
    func TestAllCases(t *testing.T) {
        for e := range testCases {
            if e.out != myFunc(e.param1, e,param2) {
                ...
            }
        }
    }
    ```
*  如何精确指定运行的测试case?
    ```golang
    go test -run NameOfTest // 精确指定case name pattern
    go test xxx_test.go     // 运行xxx_test.go的所有测试case
    ```

##### Benchmark

Benchmark 核心对象是testing.B，以函数为基本运行单元
*  函数声名：BenchmarkXxx(*testing.B)
*  运行方式：go test -bench benchFuncName：　benchFuncName可以是regexp
*  bench一般格式    
```golang
func BenchmarkHello(b *testing.B) {
    for i := 0; i < b.N; i++ {
        fmt.Sprintf("hello")
    }
}
```
   执行原理：go会重复执行BenchmarkHello b.N次，b.N会根据运行时测试情况调整
*  RunParallel 来同时运行多个function, 需用 go test -cpu cpu_num来运行
```golang
func BenchmarkTemplateParallel(b *testing.B) {
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            ...
        }
    })
}
```
*  如果Benchmark函数执行需要一个耗时的环境setup函数，此时我们并不像让其计入我们的平均耗时，可以这样：
```golang
func BenchmarkBigLen(b *testing.B) {
    big := NewBig()
    b.ResetTimer()   // 执行完setup后，reset benchmark的timer
    for i := 0; i < b.N; i++ {
        big.Len()
    }
}
```

#### Solving Real World Problem

##### 用TDD解决leetcode　22: generate parentheses

> Given n pairs of parentheses, write a function to generate all combinations of well-formed parentheses.
> For example, given n = 3, a solution set is:
[
  "((()))",
  "(()())",
  "(())()",
  "()(())",
  "()()()"
]

*  首先我们根据题目要求，先写出单元测试函数
    ```golang
    func TestGenerateParenthesis_0(t *testing.T) {
        res := generateParenthesis(0)
        if res == nil || len(res) != 0 ||  {
            t.Errorf("expected empty, got %v", res)
        }
    }
    ```
    我们首先想到测试n == 0情况，预期返回空数组；可见使用了TDD，我们会强制自己在写代码前，考虑极端case
*  修改生产代码，使之PASS
    ```golang
    func generateParenthesis(n int) []string {
        if 0 == n {
            return []string{}
        }
        return nil
    }
    ```
*  查看生产代码：有无可重构之处，有无重复；开始下一个循环，写更多的单元测试，归纳出2个case，用table-driven方式写出：
    ```golang
    import "github.com/stretchr/testify/assert"
    var testCases = []struct {
        in int
        out []string
    } {
        {
            1,
            []string{"()"},
        },
        {
            3,
            []string{
                "((()))",
                "(()())",
                "(())()",
                "()(())",
                "()()()",
            },
        },
    }
    func TestGenerateParenthesis(t *testing.T) {
        if testing.Short() { // 如果指定 go test -short ，则为short模式，此函数会被跳过
            t.Skip()
        } else {
            for _, test := range testCases {
                res := generateParenthesis(test.in)
                // 第三方assert package，提供了方便的assert函数
                assert.Equal(t, len(test.out), len(res), "array size should be same: expected %v, got %v",
                    test.out, res)
                aIsSubsetOfB(t, res, test.out) // 自己根据assert package封装，判断两集合是否有子集关系
                aIsSubsetOfB(t, test.out, res)
            }
        }
    }
    ```

补齐能通过TestGenerateParenthesis函数的代码，我们很有信心能一次AC，最终代码可以在本文最后[附录](#leetcode_22)找到

##### 用benchmark解决性能争论

我们在实际工作中遇到一个小疑惑：我们需要在运行时从别处取值来填充内存的map,需要先clear map，有两种方式：  

删除此前map的所有key，或者重新new一个map  
两者运行时性能孰优孰劣呢？用benchmark一目了然:

```golang
func BenchmarkStatic(b *testing.B) {
    for i := 0; i < b.N; i ++ {
        for k := range mp { // 用删除key的方式clear map
            delete(mp, k)
        }
        getSevenRand()  // 取得7个随机数
        fillMpContent() // 填充此map
    }
}

func BenchmarkDynamic(b *testing.B) {
    for i := 0; i < b.N; i ++ {
        mp = make(map[int]bool) // new map
        getSevenRand()
        fillMpContent()
    }
}
```

benchmark结果:

```
BenchmarkStatic-4        3000000               576 ns/op
BenchmarkDynamic-4       2000000               824 ns/op
```

可见在运行时频繁创建map会耗时较高,完整的benchmark代码也可以在[附录](#benchmark_code)中找到


#### Cool Stuff 

笔者非常喜欢Golang的一大原因便是其提供了一系列非常实用且handy的工具集，make your everyday programming life a lot easier

##### tesing的example构建动态文档

之前提到TDD的一个side effect便是产生出了＂动态＂的文档，这个特殊的文档好处便是不会过时．  
Golang对此提供了方便的example功能，类似test单元函数，但其目的更明确：是文档的一部分，而且可以用godoc工具方便的把example代码嵌入到网页文档中，比如官方string包的[example](https://golang.org/pkg/strings/#pkg-examples)，可以直接在网页文档上点运行，直接能看运行结果，非常方便，以下是要点：

*  类似test，example对应golang的函数，在go test被运行且被验证
*  example函数命名：ExampleXXX()，
    -  example of a package:    Example()
    -  function F :             ExampleF()
    -  a type T :               ExampleT()
    -  method M on type T:      ExampleT_M()
    -  多个example用 ExampleXXX_suffix()
*  需要在函数结尾以注释中用Output关键字指明expected output,不含output的example code不会被执行，但是会编译
*  Output写法:
    -  匹配特定顺序输出:
        ```golang
        func ExampleSalutations() {
            fmt.Println("hello, and")
            fmt.Println("goodbye")
            // Output:
            // hello, and
            // goodbye
        }
        ```
    -  匹配unordered的输出: Unordered output
        ```golang
        func ExamplePerm() {
            for _, value := range Perm(２) {
                fmt.Println(value)
            }
            // Unordered output:
            // 2
            // 0
            // 1
        }
        ```

##### 代码coverage

测试代码coverage能看出大体的测试覆盖情况,在go里可以用很容易的获取

```shell
go test -cover                     // package整体测试覆盖率
go test -coverprofile=coverage.out // 输出测试统计
go tool cover -func=coverage.out   // 函数级别的测试覆盖率
go tool cover -html=coverage.out   // 网页形式生成覆盖率报告
```

网页形式coverage report见下图，可见，Golang除了给出整体的覆盖率，还给出了未被测试代码覆盖的生产代码行
{% asset_img coverage.jpg %}

##### go report

如果有一个工具能直观的看到我们代码的各项指标如fmt, lint, function complexity等，其实就能对代码的质量的出一个大致的了解: 

{% asset_img report.jpg %}

上图是用一个在线工具[goreportcard](https://goreportcard.com)生成的，其代码在github可以找到: [goreportcard project](https://github.com/gojp/goreportcard)

#### 结论

本文从TDD概述开始，介绍了TDD的基本工作方法和意义；本着用Golang实践TDD的目的，介绍了testing框架，并以解决1个leetcode问题的契机展示了如何用TDD方式进行开发．在最后列举了一些有意思的特性和工具.希望对有兴趣了解TDD的朋友一些启发，欢迎讨论及指教．

#### 附录 

#####  leetcode 22 ACed code <h6 id="leetcode_22"></h6>

```golang
func generateParenthesis(n int) []string {
    allValidStrings = make([]string, 0, 256)
    walkGenerateTreeRecusive("", n, 0, 0)
    return allValidStrings
}

var allValidStrings []string

func walkGenerateTreeRecusive(curString string, total int, leftParenthesisCt int, rightParenthesisCt int) {
    if total <= 0 {
        allValidStrings = []string{}
        return
    }
    if leftParenthesisCt == total {
        for i := 0; i < total - rightParenthesisCt; i ++ {
            curString +=  string(')')
        }
        allValidStrings = append(allValidStrings, curString)
    }else {
        walkGenerateTreeRecusive(curString + string('('), total, leftParenthesisCt + 1, rightParenthesisCt)
        if leftParenthesisCt > rightParenthesisCt {
            walkGenerateTreeRecusive(curString + string(')'), total, leftParenthesisCt, rightParenthesisCt + 1)
        }
    }
}
```
#####  benchmark 代码     <h6 id="benchmark_code"></h6>

```golang
package _16_combination_sum_III

import (
    "testing"
    "math/rand"
)

const (
    MAXRAND = 100000
)

var arr []int
func init() {
    arr = make([]int, MAXRAND, MAXRAND)
    for i := 0; i < len(arr); i ++ {
        arr[i] = rand.Int()
    }
}

var mp = make(map[int]bool)
var randIndex = 0

var mpContent [7]int
func getSevenRand() {
    if MAXRAND - randIndex < 7 {
        randIndex = 0
    }
    for i := 0; i < 7; i ++ {
        mpContent[i] =  randIndex + i
    }
    randIndex = randIndex + 7
}

func fillMpContent() {
    for i := 0; i < 7; i ++ {
        mp[mpContent[i]] = true
    }
}

func BenchmarkStatic(b *testing.B) {
    for i := 0; i < b.N; i ++ {
        for k := range mp {
            delete(mp, k)
        }
        getSevenRand()
        fillMpContent()
    }
}

func BenchmarkDynamic(b *testing.B) {
    for i := 0; i < b.N; i ++ {
        mp = make(map[int]bool)
        getSevenRand()
        fillMpContent()
    }
}
```


#### Reference

TDD:  
*Clean Code* Robert C. Martin Chapter 9 单元测试  
[Introduction to Test Driven Development](http://agiledata.org/essays/tdd.html)  
[TestDrivenDevelopment](https://martinfowler.com/bliki/TestDrivenDevelopment.html)  
[The Art of Agile Development: Test-Driven Development](http://www.jamesshore.com/Agile-Book/test_driven_development.html)  
[Test driven development book [closed]](https://stackoverflow.com/questions/797026/test-driven-development-book)  

Golang Testing:  
[Golang basics - writing unit tests](https://blog.alexellis.io/golang-writing-unit-tests/)  
[Testable Examples in Go](https://blog.golang.org/examples)  
[Test-driven development with Go](https://leanpub.com/golang-tdd/read#leanpub-auto-wrapping-up-2)  
[testify package](https://github.com/stretchr/testify#installation)  
[goreportcard](https://github.com/gojp/goreportcard)  
["Dependency Injection" in Golang](http://openmymind.net/Dependency-Injection-In-Go/)  
[The cover story](https://blog.golang.org/cover)  
[How to run test cases in a specified file?](https://stackoverflow.com/questions/16935965/how-to-run-test-cases-in-a-specified-file)  

