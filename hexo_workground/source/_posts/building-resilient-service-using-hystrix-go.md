---
title: building resilient microservice using hystrix-go
date: 2018-08-03 23:09:35
tags: [golang, hystrix-go]
keywords: golang circuit-breaker
description:
---

<div align="center">
{% asset_img hystrix.png %}
</div>

#### Preface 

在微服务的搭建过程中，时常会听到熔断，降级。什么意思呢？以我们家里用的保险丝为例：其实它是一种circuit breaker，当电流超过阈值时，会断路，家电对电力资源的请求无法被满足，以保护电网安全。在微服务的世界里，由于存在非常多的跨服务调用，当某个繁忙的服务因网络或运行出现异常时，会影响依赖其的服务：经常是开始响应变慢，错误率上升；由于调用者一般都有重试机制，发现超时或者出错会重试几次，进一步加剧了已经出问题节点的负担，使之恶化导致完全停止服务；若调用者不限定timeout, 会因为过长的等待时间影响用户体验;更严重的是，其会积压大量用户请求，耗尽内存，很可能因为一个外部不太重要的服务，而造成关键服务停止，进而引发连锁反应，像雪崩一样，整个微服务系统因此崩溃；这个现象也称之为级联失败。

<!-- more -->

如何解决呢？思路也很简单：借鉴断路器思想：当依赖的服务出现问题时，切断与其联系，期望其能恢复，这就是熔断；当外部服务熔断时，采取一些补偿措施，使得自己接受的整个请求不至于失败：如视频网站个性化推荐系统熔断后，可以转向获取当天最热视频top 10推荐给用户，不至于因为推荐系统超时没有返回，导致整个用户页面显示500错误，这就是服务降级。熔断被归纳为[circuit breaker pattern](https://martinfowler.com/bliki/CircuitBreaker.html)，除了最基本的切断线路，还有自动恢复(circuit close)机制：当预设阈值超出，则将断路器置为open状态，切断与依赖服务关系，所有请求立刻返回失败；过一定时间，若检测到服务恢复，则断路器置为close状态。

我们接下来讲述如何用Hystrix-go实现服务熔断和降级。


#### What is Hystrix-go

Hystrix是NetFlix开发，非常popular的避免级联故障的库，用java实现；

> In a distributed environment, inevitably some of the many service dependencies will fail. Hystrix is a library that helps you control the interactions between these distributed services by adding latency tolerance and fault tolerance logic. Hystrix does this by isolating points of access between the services, stopping cascading failures across them, and providing fallback options, all of which improve your system’s overall resiliency.

Hystrix-go是由个人开发的轻量模拟Hystrix功能的Golang library，简单，好用。

#### Hello world Hystrix-go 

需要先导入包 

```
go get "github.com/afex/hystrix-go/hystrix"
```

最简单的hystrix-go code:

```golang
hystrix.Go("command name", func() error {
    // normal path code
    return nil
}, func(err error) error {
    // do this when errors occur 
    return nil
})
```
跑起来了，不是么？hystrix.Go有三个参数

*  command name为此调用的名字
*  hystrix.Go 第一个函数代表正常执行功能
*  第二个函数为fallback，出现错误时自动执行

这段代码体现了Hystrix-go最核心的功能：为正常执行的功能指定出错时的fallback函数。听起来像降级，但是熔断在哪里？

#### Hystrix-go 配置


##### fallback阈值配置

Hystrix-go可以通过下列代码配置：

```golang
hystrix.ConfigureCommand("unique command name", hystrix.CommandConfig{
    Timeout:               1000,
    MaxConcurrentRequests: 100,
})
```

这里配置了针对单次请求何种情况会自动调用fallback: 

*  请求超过Timeout时间没有响应
*  最大并发超过限定

意味着单次请求如果超时或者超过并发限制，会直接导致fallback函数调用，产生一次error；  

如果想定义自己的错误，比如收到500： 

```golang
ch := hystrix.Go("hello hystrix", func() error {
  if resp.Code == 500 {
    return errors.New("code 500 got")         // 在这里返回错误，fallback会被调用
  }
  return nil
}, func(err error) error {
    // do this when errors occur
    fmt.Println("fallback called because ", err)
    return fmt.Errorf("fallback returned error %s", err)
})
err := <- ch
fmt.Println("error got after hystrix.Go: ", err)
```

产生如下输出:

```
fallback called because  code 500 got
err before exit fallback failed with 'fallback returned error code 500 got'. run error was 'code 500 got'
```

可以看出error传递路径

##### 熔断即恢复阈值配置

产生的error频率和下列配置项决定了何时熔断：

*  RequestVolumeThreshold
*  SleepWindow
*  ErrorPercentThreshold

配置项意义参见代码settings.go：

```golang
var (
    // DefaultTimeout is how long to wait for command to complete, in milliseconds
    DefaultTimeout = 1000
    // DefaultMaxConcurrent is how many commands of the same type can run at the same time
    DefaultMaxConcurrent = 10
    // DefaultVolumeThreshold is the minimum number of requests needed before a circuit can be tripped due to health
    DefaultVolumeThreshold = 20
    // DefaultSleepWindow is how long, in milliseconds, to wait after a circuit opens before testing for recovery
    DefaultSleepWindow = 5000
    // DefaultErrorPercentThreshold causes circuits to open once the rolling measure of errors exceeds this percent of requests
    DefaultErrorPercentThreshold = 50
)
```

当错误率ErrorPercentThreshold达到阈值，即发生熔断；这时所有的请求并不会真正发送，其fallback都是立刻返回，错误原因：circuit open；这时我们可以把降级的逻辑放在fallback函数中执行。SleepWindow和RequestVolumeThreshold定义了断路器如何恢复为closed状态。


#### 同步 V.S 异步

hystrix.Go内部是启动goroutine执行请求，返回chan error指示错误信息，调用代码：

```golang
output := make(chan bool, 1)
errors := hystrix.Go("my_command", func() error {
    // talk to other services
    output <- true
    return nil
}, nil)

select {
case out := <-output:
    // success
case err := <-errors:
    // failure
}
```

也支持同步调用，一次调用只返回一个错误

```golang
err := hystrix.Do("my_command", func() error {
    // talk to other services
    return nil
}, nil)
```

#### Hystrix-go总结

*  避免client过长等待，超出并发限制；同时和自定义错误一起可以触发fallback，熔断机制
*  设定错误超过阈值时断路，避免向运行出问题服务进一步增加其负担
*  抽取出的断路器成为非常有价值的监控点
*  断路器除了自动熔断，恢复；一般还应支持人工干预
*  熔断后的fallback函数实现降级逻辑


#### reference

[CircuitBreaker](https://martinfowler.com/bliki/CircuitBreaker.html)  
[Hystrix wiki](https://github.com/Netflix/Hystrix/wiki)  
[Fault Tolerance in Go](http://thediscoblog.com/blog/2015/02/07/fault-tolerance-in-go/)  
[afex/hystrix-go](https://github.com/afex/hystrix-go)  


