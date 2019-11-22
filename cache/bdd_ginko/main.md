


beforeeach: 

*  其存在context/describe的所有it都会执行
*  试做所有it都一样，不管处于nest第几位置，都会运行beforeeach的statement
*  多级beforeeach是链式关系, parent beforeeach先执行 
    ```
    beforeeach
      ...
        beforeeach
    ```
*  Describe blocks to describe the individual behaviors of your code
*  Context blocks to exercise those behaviors under different circumstances
*  Describe包含多个Context
*  BeforeEach and AfterEach blocks run for each It block
*  JustBeforeEach blocks are guaranteed to be run after all the BeforeEach blocks have run and just before the It block has run

## 区别

While TDD focuses on the technical, or implementation details, BDD focuses on visible, behavioral details. 

Another way to think about it is that TDD focuses on ensuring that the individual parts of the application do what they should, while BDD focuses on ensuring that these parts work together as expected.

## Ginko

As with popular BDD frameworks in other languages, Ginkgo allows you to group tests in Describe and Context container blocks. Ginkgo provides the It and Specify blocks which can hold your assertions. It also comes with handy structural utilities such as BeforeEach, AfterEach, BeforeSuite, AfterSuite and JustBeforeEach that allow you to separate test configuration from test creation, and improve code reuse.

Ginkgo comes with support for writing asynchronous tests. This makes testing code that use goroutines and/or channels as easy as testing synchronous code.


may only call It from within a Describe, Context or When

*  Describe/Context/When: Describe, Context and When blocks are functionally
//equivalent.  The difference is purely semantic
*  It/Specify: //It blocks contain your test code and assertions.  You CANNOT nest any other Ginkgo blocks
//within an It block.
*  Measure: Measure blocks run the passed in body function repeatedly (determined by the samples argument) and accumulate metrics provided to the Benchmarker by the body function.
*  You can use BeforeEach to set up state for your specs. You use It to specify a single spec.
In order to share state between a BeforeEach and an It you use closure variables, typically defined at the top of the most relevant Describe or Context container.
*  In general, the only code within a container block should be an It block or a BeforeEach/JustBeforeEach/JustAfterEach/AfterEach block, or closure variable declarations. It is generally a mistake to make an assertion in a container block.
*  It is also a mistake to initialize a closure variable in a container block. If one of your Its mutates that variable, subsequent Its will receive the mutated value. This is a case of test pollution and can be hard to track down. Always initialize your variables in BeforeEach blocks.
*  Pending Specs: P, X: By default, Ginkgo will print out a description for each skipped spec. You can suppress this by setting the --noisySkippings=false flag.
*  
*  By Doc: 附加test Addition Doc
*  ? beforeeach justbeforeeach, 嵌套是如何工作的？ 
    -  beforeeach下面的嵌套共运行一次还是每个leaf都会运行。

## To next blog

*  Think a example



## Next 

next step, write blog series: BDD test with Ginko & Omega

*  Write 101 blog: basic example: need a good example for it
*  Before each, just, focus, pending, async, CI
*  BDD and more (cucumber also have go BDD library)


*  Before each, just, after
*  async
*  maybe parallel? 至少弄清是干什么的
*  CI: `ginkgo -r --randomizeAllSpecs --randomizeSuites --failOnPending --cover --trace --race --progress`

## Target

BDD结合

## TODO

*  实践文章:  Getting Started with BDD in Go Using Ginkgo
*  


## Reference

[TDD 与 BDD 仅仅是语言描述上的区别么？](https://www.zhihu.com/question/20161970)
[Getting Started with BDD in Go Using Ginkgo](https://semaphoreci.com/community/tutorials/getting-started-with-bdd-in-go-using-ginkgo)
[行为驱动开发（BDD）你准备好了吗？](https://blog.csdn.net/chancein007/article/details/77543494)