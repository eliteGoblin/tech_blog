---
title: using gomock in TDD
date: 2018-05-18 14:27:20
tags: [TDD,gomock,golang]
keywords: gomock TDD test golang  
---

<div align="center">
{% asset_img jack_mock2.jpg %}
</div>

#### Preface 

在[golang TDD实践](https://elitegoblin.github.io/2018/01/31/golang%20TDD%E5%AE%9E%E8%B7%B5/)一文中我们领略了TDD method的优美与强大．引入TDD可以极大的提升代码的内部质量和设计水平，从而提升我们的自信，减少恐惧．但如果我们的功能代码依赖外部服务结果时该如何测试呢？假如我们的模块依赖数据库读取的数据，那我们在进行TDD时必须先搭建好DB及初始化好数据么？这很多时候会是一个的负担，阻止我们进行TDD．最好能有种机制能根据我们的test case轻量的模拟外部服务行为，这就是TDD里mock object的意义．

<!-- more -->

#### mock机制及gomock简介

##### mock机制

我们所做的一切目的都是测试我们自己代码逻辑，例如代码依赖数据库返回的数据时，可以用mock机制来方便的为test case生成其依赖的测试数据.  

如何mock外部服务呢？关键在于将外部服务建模成为object, 将其public method抽象成为interface, 从而能decouple我们的代码和外部服务;这样就能用mock object来代替真的外部服务object,从而测试我们代码正确性.  

**mock object是interface和测试数据的二合一：** 经过interface的decouple,我们的代码变为依赖interface的output，我们create mock object的过程，其实就是建立test case的过程，指定test case的input及output, 供测试框架的测试代码调用．

##### gomock basic

[gomock](https://github.com/golang/mock)是golang自带的mock framework，那它具体如何使用呢？大致步骤:  

*  前提: install gomock
*  代码抽象外部依赖服务为interface, 放入单独文件如service_handler.go
*  用安装好的mockgen程序， 传入service_handler.go, 根据interface生成mock obj的实现，如service_mock.go
*  利用mock object代码编写unit test function, **在每个test case**设置好mock object的*input --> output* 对，作为测试case, 并将设置好的mock object嵌入自己代码
*  运行测试，利用mock object的预先配置好的输入和输出，检查自己的程序时候符合预期．

##### install gomock

gomock分为两部分:  

*  golang library code: gomock library
    ```shell
    # set GOPATH to folder which contains src, pkg, bin...
    export GOPATH=...
    go get github.com/golang/mock/gomock
    ```
*  gomock generator : 安装mockgen程序，根据代码的interface生成mock object的代码(类似protobuf自动生成代码过程)
    ```shell
    # 在前一步GOPATH和go get完成的基础上
    # go install基于gomock library代码
    go install github.com/golang/mock/mockgen
    ```

#### mock a redis-server

我们以一个例子来介绍如何使用gomock：　　

我们演示代码依赖以存储于Redis的数据(personId，personName)对，均为普通string.
如果我们想测试，需要mock redis(确切的说是mock redis client的interface，通过此interface获取redis-server上的数据)，完整的代码见[这里](https://github.com/eliteGoblin/Notes/tree/master/cs/Languages/go/gomock/src/frank)

##### sample工程介绍

作用：实现简单的调用redis interface，获取value，并打印

*  project目录结构
```shell           
./src/frank/
├── main.go
├── redis_helper
│   └── redis_helper.go       # redis interface
├── redis_helper_mock
│   └── redis_helper_mock.go  # mockgen生成的mock object代码，实现了interface
└── report
    ├── gen_report.go         # 调用interface的业务代码:打印redis value
    └── gen_report_test.go    # 在此用mock object实现测试代码
```
*  redis interface定义
```golang
package redis_helper
// 待mock的interface
type RedisHelper interface {
    ConnectToRedisServer(connStr string) error // 连接至redisServer
    GetKey(keyName string) (string, error)     // 获取key value
    SetKey(keyName string) error               // set key by keyName
    GetPersonsNameMatchPrefix(keyPrefixName string) ([]string, error) // prefix匹配所有person name
} 
```
*   用户代码嵌入实现interface的object，因此可以将实际访问redis的object替换为mock object. 

```golang
type Reporter struct {
    redisHelper redis_helper.RedisHelper
}
func (selfPtr *Reporter)ShowPersonName(personId string) (personName string, err error) {
    return selfPtr.redisHelper.GetKey(personId)
}
```
*  如下命令生成mock object的RedisHelper interface代码，输出到redis_helper_mock.go

> mockgen -destination=redis_helper_mock/redis_helper_mock.go -package=redis_helper_mock frank/redis_helper RedisHelper 

*  destination: 输出文件
*  package: 生成mock代码的package
*  frank/redis_helper: interface所在目录
*  RedisHelper: 待mock的interface的name，一般希望不同interface存放在不同mock文件中，每次指定一个interface来mock

我们的**最终目的**就是做到不依赖redis服务而编写ShowPersonName的单元测试代码

##### basic case using mock object

gen_report_test.go
```golang
func TestShowPersonName(t *testing.T) {
    mockCtrl := gomock.NewController(t) 
    defer mockCtrl.Finish()
    mockRedisHelper := redis_helper_mock.NewMockRedisHelper(mockCtrl)
    mockRedisHelper.EXPECT().GetKey("person_1").Return("frank", nil)

    reportObj := report.Report{
        RedisHelper : mockRedisHelper, // 用mock object初始化业务代码
    }
    name, err := reportObj.RedisHelper.GetKey("person_1")
    assert.True(t, name == "frank")    // EXPECT已经设置好，"person_1"返回"frank"
}
```
要点:  
*  必须先NewController，在测试函数执行最后调Finish()
*  用.EXEPCT().Method(inputs...).Return(outputs...)方式来指定(inputs, outputs)pair
*  用EXPECT预设的case必须在Finish()之前满足调用数量一致，调用顺序和EXPECT顺序可以不等，可用Times()函数制订某method调用次数.  
    ```golang
    mockRedisHelper.EXPECT().GetKey("person_2").Return("lisha", nil)
    mockRedisHelper.EXPECT().GetKey("person_1").Return("frank", nil)

    name, err := mockRedisHelper.GetKey("person_2")
    assert.True(t, name == "lisha") // 正确
    // 最终会报错，因为没有调用GetKey("person_1")　
    ```


##### mock order

有时我们需要验证的操作之前包含前后顺序关系，比如redis connection应先init再使用，可以用After()来实现

```golang
mockCtrl := gomock.NewController(t)
defer mockCtrl.Finish()

mockRedisHelper := redis_helper_mock.NewMockRedisHelper(mockCtrl)
callFirst := mockRedisHelper.EXPECT().ConnectToRedisServer("localhost:6704").Return(nil).Times(1)
mockRedisHelper.EXPECT().GetKey("person_1").Return("frank", nil).After(callFirst)
mockRedisHelper.EXPECT().GetKey("person_2").Return("lisha", nil).After(callFirst)

mockRedisHelper.ConnectToRedisServer("localhost:6704")
name, err := mockRedisHelper.GetKey("person_2")
assert.True(t, name == "lisha" && err == nil)
name, err = mockRedisHelper.GetKey("person_1")
assert.True(t, name == "frank" && err == nil)
```

如果先调用GetKey后调用ConnectToRedisServer会报错．After适用于两两调用间的先后关系，如想指定一系列操作的先后关系，用gomock.InOrder：

```golang
gomock.InOrder(
    mockDoer.EXPECT().DoSomething(1, "first this"),
    mockDoer.EXPECT().DoSomething(2, "then this"),
    mockDoer.EXPECT().DoSomething(3, "then this"),
    mockDoer.EXPECT().DoSomething(4, "finally this"),
)
```

##### mock type match

前文提到的mock object配置(inputs, outputs)对，当需要不需要精确匹配inputs时，gomock提供模糊匹配inputs的机制:match

```golang
mockRedisHelper.EXPECT().ConnectToRedisServer(gomock.Any()).Return(nil).Times(1)
```
当不关心传入的值时，匹配任意inputs可用gomock.Any()，其他内置的还有: 

> 
gomock.Eq(x): uses reflection to match values that are DeepEqual to x  
gomock.Nil(): matches nil  
gomock.Not(m): (where m is a Matcher) matches values not matched by the matcher m  
gomock.Not(x): (where x is not a Matcher) matches values not DeepEqual to x  

**注** Times(1)指定ConnectToRedisServer只能被调用一次，否则test失败

##### 其他高级特性

*  gomock支持自定义match行为，需要实现gomock.Match interface
*  可以给mockgen生成的mock object加入自定义行为: mock object其实只是实现了interface，本身并没有任何功能，一般通过EXPECT配置特定输入输出，但也可以用Do()为其添加职责(可参考生成的mock object代码，动态增加职责，用到了decorator模式)
    ```golang
    mockDoer.EXPECT().
    DoSomething(gomock.Any(), gomock.Any()).
    Return(nil).
    Do(func(x int, y string) {
        if x > len(y) {
            t.Fail()
        }
        // 其他更复杂的自定义行为 ... 
    })
    ```

#### conclusion

*  mock的目的是使单元测试代码不以来外部服务部署，进一步促进了代码与外部依赖的解偶．
*  通过mock object的测试后，仍需要有途径保证访问真实外部服务代码的正确性．
*  mock其实是interface实现与test case的二合一，可以在测试代码中替换真实interface实现
*  gomock基本test case通过一系列的EXPECT来指定，调用次数需匹配设定，顺序没有要求．
*  gomock能实现诸如调用次数，调用顺序，自定义匹配Matche，自定义mock object行为等高级特性．

#### reference

[Testing with GoMock: A Tutorial](https://blog.codecentric.de/en/2017/08/gomock-tutorial/)