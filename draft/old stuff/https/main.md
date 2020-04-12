

#### general
目的 解决最简单case

*  生成自签名证书, 并添加至chrome中
*  nginx接上证书, 打给后端的服务，把证书剥掉
*  后端验证是否收到


问下董昊 openresty


#### 实践

*  实现golang自签名证书
*  nginx剥证书
*  udesk目前production: 


#### self-signed cert

*  failed

#### 非对称加密解密示例


[GO加密解密之RSA](http://blog.studygolang.com/2013/01/go%E5%8A%A0%E5%AF%86%E8%A7%A3%E5%AF%86%E4%B9%8Brsa/)
[Go by Example: Base64 Encoding](https://gobyexample.com/base64-encoding)


#### https

*  据私钥生成公钥

#### 目前进度


*  搞懂非对称加密
*  自签名，实现http client server  InsecureSkipVerify: true 的认证
*  tony白的文章很全面, client认证server的证书
    -  TODO: 对服务端的证书进行校验 该看

#### refs


*  [Go和HTTPS](https://tonybai.com/2015/04/30/go-and-https/)
*  [golang-https-example](https://github.com/jcbsmpsn/golang-https-example)
*  [使用Go实现TLS 服务器和客户端](http://colobu.com/2016/06/07/simple-golang-tls-examples/)