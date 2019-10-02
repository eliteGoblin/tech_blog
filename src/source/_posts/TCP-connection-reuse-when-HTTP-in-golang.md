---
title: TCP connection reuse when HTTP in golang
date: 2018-04-11 19:03:59
tags: [golang, HTTP]
keywords: HTTP reuse
description:
---

{% asset_img persistenthttp.jpg %}

#### Preface

实际工作中会碰到这样的情况：

面对频繁的HTTP请求，我们应让它们尽量复用底层的TCP连接，而不是每次HTTP都新建TCP连接，有如下好处：
*  减少了每次TCP握手，释放连接的开销
*  如果TCP使用了TLS，会给每次连接建立带来更大的开销(真的很大哟)　　

那么如何在golang中实现呢？我们首先了解下HTTP关于连接复用的background知识．



<!-- more -->

#### HTTP Background

参见wikipedia:  
> HTTP persistent connection, also called HTTP keep-alive, or HTTP connection reuse, is the idea of using a single TCP connection to send and receive multiple HTTP requests/responses, as opposed to opening a new connection for every single request/response pair

可以理解为单一TCP连接的multiplex．HTTP协议对其的支持

*  HTTP 1.0: client发送的HTTP请求HEADER中需要设置 keepalive
    ```http
    Connection: keep-alive
    ```
    server返回时也会带此header
*  HTTP 1.1: 默认就是reuse, 除非另外指定

那是不是意味着我们的web实现基于HTTP 1.1就自动得到这项好处呢？我们测试一下

#### Golang HTTP reuse之旅

我们如何验证golang实现的HTTP client进行通信时reuse了底层的TCP呢？抓TCP包看是否有多次握手就一目了然了

##### HTTP client 实现

main.go文件(为了清晰去除了err handle代码)
```golang
var (
    httpClient *http.Client
)

const (
    MaxIdleConnections int = 10
    RequestTimeout     int = 5
)

// init HTTPClient
func init() {
    httpClient = createHTTPClient()
}

// createHTTPClient for connection re-use
func createHTTPClient() *http.Client {
    client := &http.Client{
        Transport: &http.Transport{
            MaxIdleConnsPerHost: MaxIdleConnections,
        },
        Timeout: time.Duration(RequestTimeout) * time.Second,
    }
    return client
}

func main() {
    var endPoint string = "http://localhost:8087/"
    for i := 0; i < 3; i++ {
        req, _ := http.NewRequest("GET", endPoint, nil)
        response, _ := httpClient.Do(req)
        // MUST read all response's data
        ioutil.ReadAll(response.Body)
        // Close the connection to reuse it
        response.Body.Close()
    }
}
```

实现要点:  
*  全局唯一http client,可以传入自定义的Transport member, 定义MaxIdleConnections
*  必须全部读取http response的data, 这样TCP就不会残留数据，影响reuse
*  必须调用response.Body.Close()，这也是好习惯：随时清理不需要的资源

client写好，我们需要抓包工具和HTTP server

##### Tcpdump抓包

先用tcpdump抓取lookback上所有的tcp包，保存为p1.pcap文件
```shell
sudo tcpdump -i lo -s 0 -n tcp -w /tmp/p1.pcap
```
再用wireshark打开p1.pcap，指定过滤条件tcp.port == 8087，过滤src或dst port是8087的，就是我们感兴趣的tcp包

{% asset_img wireshark.jpg %}

有了抓包工具，接下来我们进行HTTP测试，看是否如我们预期只有一个TCP握手

##### 第三方HTTP Server测试

为了方便，采用python提供的http module，作为http server:

```shell
python3 -m http.server 8087
go run main.go
```

抓包结果如下:
{% asset_img python_http_not_work.jpg %}

我们注意到还是新建了3次TCP连接，为什么呢？  

进一步观察，我们发现TCP连接每次是由server端关闭的(server先向client发送的FIN)，试了nginx也是一样，似乎两者都处理完请求后主动关闭了TCP连接，我们再试一下server不主动关闭连接时是否会发生我们预期的TCP connection reuse行为

##### 自己实现简单HTTP server，不主动关闭连接

http_echo_server.go  

```golang
func sayHello(w http.ResponseWriter, r *http.Request) {
    fmt.Println("sayHello")
    message := r.URL.Path
    message = strings.TrimPrefix(message, "/")
    message = "Hello " + message
    w.Write([]byte(message))
}
func main() {
    http.HandleFunc("/", sayHello)
    if err := http.ListenAndServe(":8087", nil); err != nil {
        panic(err)
    }
}
```

抓包结果如下:  
{% asset_img reused.jpg %}

终于实现了HTTP connection reuse! 记得我们要求必须:

```golang
response.Body.Close()
```

假如我们不执行这条Close语句，果然没有发生TCP connection reuse，抓包结果:  
{% asset_img not_close_resp.jpg %}

#### 总结

*  HTTP connection reuse是golang的默认行为，但需要保证client端
    -  重用http client， golang已保证其是并发安全
    -  必须读取response数据
    -  比如调用 response.Body.Close()
*  server根据不同实现，也可能主动close TCP connection，影响reuse
*  go-nuts google group是个好地方，golang碰到的疑难杂症建议去上面提问试试

#### refs


[HTTP persistent connection](https://en.wikipedia.org/wiki/HTTP_persistent_connection)  
[Reusing http connections in Golang](https://stackoverflow.com/questions/17948827/reusing-http-connections-in-golang#comment26240189_17953506)  
[Keep-Alive in http requests in golang](https://awmanoj.github.io/tech/2016/12/16/keep-alive-http-requests-in-golang/)   
[TCP connection reuse when doing HTTP](https://groups.google.com/forum/#!topic/golang-nuts/IZ2p1sHkeBQ)  