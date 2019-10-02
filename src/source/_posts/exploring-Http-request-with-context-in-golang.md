---
title: Exploring param passing With Golang's HTTP Context Support
date: 2018-04-23 16:59:39
tags:
keywords:
description:
---

{% asset_img gophermega_meitu_1.jpg %}　

#### Preface

想像如下场景：  
　　
某天，满头大汗的客服同事找到你反馈一个用户详情(PersonDetail)的请求失败了，要求你调查．服务架构如下:  

{% asset_img system.jpg %}　

调用顺序解释:  
1.  前端请求通过HTTP到达gateway节点, 目的是获取person_id的PersonDetail, 包括Name,Address,Education信息
2.  gateway节点需要从Name, Address, Education后端节点分别获取相关信息
3.  gateway节点将各个后端节点的response整合成为一个json后，传给前端

在这个过程中，各个后端服务节点都可能出问题导致此次请求失败，该如何查找问题呢？  

这时如果有一个请求对应的一个uniqueId,用此可以将系统中这个request串起来：无论在哪个后端节点，只要在日志中过滤此uniqueId，就能找到对应日志，从而最终解决问题．  

以上机制其实是分布式微服务系统的"刚需"，没有这种机制，几乎不可能进行有效的日志排查．　　

接下来看看我们如何在用golang和HTTP实现的系统中实现上述机制．

<!-- more -->


#### 需求及实现分析

我们需要的其实是一种能在请求端传递uniqueId，并在请求到达的后端能将此uniqueId打印出来的机制:  

考虑到我们使用的是HTTP，用HTTP传递和解析数据的方式都可以实现，如：  

*  url parameter：
    ```shell
    http:my_server_addr/path?unique_id=xxxx
    ```
*  HTTP header
    ```golang
    // HTTP client设置header
    req.Header.Set("My-Request-ID", "xxxx")
    // HTTP server读取Header  
    reqID := req.Header.Get("My-Request-ID")   
    ```
*  cookie，我们本文接下来的例子使用这种方式  
    ```golang
    // 设置请求时使用cookie
    cookie := http.Cookie{Name: "request_id", Value: "xxx", Expires: expiration}
    http.SetCookie(w, &cookie) 
    // 服务端获取cookie
    //   r is a *http.Request
    cookie, _ := r.Cookie("request_id") 
    ```
-  POST json数据，等等  

实际工程中考虑到网络环境的复杂，需要HTTP client对外部HTTP请求能自动实现随时取消，超时取消等功能，这个golang已经有context机制来负责对超时和取消的控制，且在golang 1.7之后已经在HTTP包中集成了context．因此我们本篇文章的目的是：  

*  如何在golang中利用HTTP传递和解析uniqueId
*  如何利用context在进程内传递uniqueId并实现cancel和timeout发起的http request


#### server端实现解析unique request id

##### function wrapper增加函数行为

最简单的实现方式是在每个http handler里面解析HTTP传来的uniqueId，还有一种更好的方式，不用改动已经有的handler，利用函数编程思想，做一个通用的function wrapper来增加解析uniqueId的功能，并利用http的context将uniqueId传入目标函数，实例代码如下：

```golang
func AddContext(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // our own control code here
        ...
        next.ServeHTTP(w, r) 
    })
}
```

###### HTTP request的context传递参数

```golang
func AddContext(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // out extra parse uniqueId code here
        uniqueId := getUniqueIdFromHTTPReq()
        // save uniqueId in context obj
        ctx := context.WithValue(r.Context(), "unique_id", uniqueId)
        // do the real request handle job with uniqueId passed by context
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

##### 完整的例子

以下的代码来自[Simple Golang HTTP Request Context Example](https://gocodecloud.com/blog/2016/11/15/simple-golang-http-request-context-example/)，感谢作者提供的有趣例子：  

*  关键在于StatusPage函数，如果user已经通过cookie设置了Username，则将其通过我们的middleware wrapper：　**AddContext**传入我们的HTTP handler，同样可以用来实现我们开始提出的uniqueId问题．　　

```golang
func main() {
    mux := http.NewServeMux()

    mux.Handle("/", AddContext(http.HandlerFunc(StatusPage)))
    mux.HandleFunc("/login", LoginPage)
    mux.HandleFunc("/logout", LogoutPage)

    log.Println("Start server on port :8085")
    contextedMux := AddContext(mux)
    log.Fatal(http.ListenAndServe(":8085", contextedMux))
}

func StatusPage(w http.ResponseWriter, r *http.Request) {
    if username := r.Context().Value("Username"); username != nil {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Hello " + username.(string) + "\n"))
    } else {
        w.WriteHeader(http.StatusNotFound)
        w.Write([]byte("Not Logged in"))
    }
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
    expiration := time.Now().Add(365 * 24 * time.Hour) // Set to expire in 1 year
    cookie := http.Cookie{Name: "username", Value: "alice_cooper@gmail.com", Expires: expiration}
    http.SetCookie(w, &cookie)
}

func LogoutPage(w http.ResponseWriter, r *http.Request) {
    expiration := time.Now().AddDate(0, 0, -1) //Set to expire in the past
    cookie := http.Cookie{Name: "username", Value: "alice_cooper@gmail.com", Expires: expiration}
    http.SetCookie(w, &cookie)
}

func AddContext(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Println(r.Method, "-", r.RequestURI)
        cookie, _ := r.Cookie("username")
        fmt.Println("Cookie got: ", cookie)
        if cookie != nil {
            ctx := context.WithValue(r.Context(), "Username", cookie.Value)
            next.ServeHTTP(w, r.WithContext(ctx))
        } else {
            next.ServeHTTP(w, r)
        }
    })
}
```

我们可以用curl来发出带cookie的HTTP请求

```shell
curl localhost:8085/ --cookie "username=alice_cooper@gmail.com"
```

#### 实现能自动timeout的HTTP request

当gateway节点请求后端时，本身会作为一个http client，我们如何利用context实现具备timeout行为的HTTP request呢?思路如下:  

*  将阻塞执行的HTTP request函数放入goroutine中，并将结果返回至answerChannel
*  主goroutine用select来判断context的DoneChannel和answerChannel

```golang
func jobWithCancelHandler(w http.ResponseWriter, r * http.Request){
    // 取得request传入的context，增加5秒timeout功能
   ctx,cancel = context.WithTimeout(r.Context(), 5 * time.Second)
   select{
    case <-ctx.Done():
        log.Println(ctx.Err())
        return
    case result:=<-longRunningCalculation(timecost):
        io.WriteString(w,result)
   } 
}
func longRunningCalculation(timeCost int)chan string{
    result:=make(chan string)
    go func (){
        time.Sleep(time.Second*(time.Duration(timeCost)))
        result<-"Done"
    }()
    return result
}
```

#### 总结

本文通过实际工作中排查系统日志需求，引入golang中利用HTTP传递request相关参数的一般方法: 利用context，同时方便的实现了HTTP request的cancel和timeout控制．　　

#### refs

[Simple Golang HTTP Request Context Example](https://gocodecloud.com/blog/2016/11/15/simple-golang-http-request-context-example/)　　

[curl网站开发指南](http://www.ruanyifeng.com/blog/2011/09/curl.html)　　

[使用Golang的Context管理上下文](https://blog.csdn.net/u014029783/article/details/53782864)　　