---
title: 漫谈Url Encode
date: 2017-10-19 16:23:01
tags: [http]
keywords: http url
---

<div align="center">
{% asset_img url.gif %}
</div>

#### 什么是Url Encode以及Why

我们日常都要接触各种各样的url，比如我们在京东搜索 *氢氟酸*，浏览器会生成如下url，请求给京东web服务器:

```url
https://search.jd.com/Search?keyword=氢氟酸&enc=utf-8
```

<!-- more -->

此url由如下部分组成：

*  search.jd.com是HTTP请求的地址，会被DNS服务器解析成为ip地址
*  /Search是请求的path，不同的请求path对应不同的web后端处理逻辑，这里我们调用的/Search表明我们请求的是搜索服务
*  ?keyword=氢氟酸&enc=utf-8 **?** 后面是query parameter，**&** 是多个parameter的分隔符，指明搜索关键字和编码方式

目前为止，url很直观，但当我们搜索 **氢氟 酸**(多一个空格)和C&K时，url变成如下的样子：
```
https://search.jd.com/Search?keyword=氢氟%20酸&enc=utf-8
https://search.jd.com/Search?keyword=C%26K&enc=utf-8
```

看起来是空格和 **&** 被转为 **%xx** 发送出去了，为什么需要这种转换? 

因为url存在特殊字符，比如&，它起到分割不同param的作用，但是我们想搜索的keyword带有&的话，比如这样：keyword=C&K，url在解析的时候就会错误的认为keyword=C，因此需要有种方法把对于url有特殊含义的字符(如&)转换，这就是Url Encode的作用


#### 如何Url Encode


那么url中到底哪些字符需要被转义(URL Encode)呢？如何进行？

> URLs can only be sent over the Internet using the ASCII character-set.

也就是说最终传输的url仅仅会包含ASCII码指定的字符

url中包含的字符分为: 

*  reserved character: 
    ```
    !   *   '   (   )   ;   :   @   &   =   +   $   ,   /   ?   #   [   ]
    ```
*  unreserved character
    ```
    0-9a-zA-A-_.~
    ```
*  Other characters
    如空格 **"** 等

[RFC 1738](http://www.ietf.org/rfc/rfc1738.txt) 做了规定
> "...Only alphanumerics [0-9a-zA-Z], the special characters "$-_.+!*'()," [not including the quotes - ed], and reserved characters used for their reserved purposes may be used unencoded within a URL."  
> "只有字母和数字[0-9a-zA-Z]、一些特殊符号"\$-_.+!*'(),"[不包括双引号]、以及某些保留字，才可以不经过编码直接用于URL。"


当我们需要在url中传送url保留字符时，就需要进行url encode了:

Reserved characters after url encoding 

| ! | # | $ | & | ' | ( | ) | * | + | , | / | : | ; | = | ? | @ | [ |  ] |
| -- | -- | -- | -- | -- | -- | -- | -- | -- | -- | -- | -- | -- | -- | -- |  -- | -- | -- |
| %21 | %23 | %24 | %26 | %27 | %28 | %29 | %2A | %2B | %2C | %2F | %3A | %3B | %3D | %3F |  %40 | %5B | %5D |  

可见，编码的方式是用 **%** 加ASCII hex码完成的，因此url encode又成为percent-encoding，
url-encoding也能对unreserved字符进行编码，如a被编码为%61但是不建议这样做，没必要而且可能会有潜在的兼容性问题。除了reserverd字符和unreserved字符，其他字符必须进行url编码:  

| newline | space | " | % |  < | > | \ | ^ | ` | { | } | ｜　|
| -- | -- | -- | -- | -- | -- |-- | -- | -- |-- | -- | -- |
| %0A| %20 | %22 | %25 | %3C | %3E | %5C | %5E | %60 | %7B | %7D | %7C |


#### url装载非Ascii码信息

比如在京东搜索**台灯**，url如何装载汉字信息呢？答案仍是url encoding，汉字被编码为一系列的含　**%** 的ASCII串

有些url为什么会包含汉字呢?
{% asset_img chars.jpg %}
其实这只是浏览器显示效果，打开chrome控制台，可以看到汉字是被url encode过的
{% asset_img header.jpg %}

#### golang中实现Url Encode

利用net/url包可以很方便的进行url编码操作  


把string进行url encode，这样就能将其安全的作为url使用

```
import "http/url"
safeUrl := url.QueryEscape(myQueryString)
```

如果想生成带parameter的url呢? [stackoverflow上一个回答](https://stackoverflow.com/questions/13820280/encode-decode-urls)提供的代码很方便:

```golang
var Url *url.URL
// 构建url base
Url, err := url.Parse("http://www.example.com")
if err != nil {
    panic("boom")
}

// Path可以直接传入reserved chars
Url.Path += "/some/path/or/other_with_funny_characters?_or_not/"
parameters := url.Values{}
// 加入query部分，即生成?后面的url
parameters.Add("hello", "42")
parameters.Add("hello", "54")
parameters.Add("vegetable", "potato#fresh") // #为reserved字符
Url.RawQuery = parameters.Encode()

fmt.Printf("Encoded URL is %q\n", Url.String())
```

输出如下

```
http://www.example.com/some/path/or/other_with_funny_characters%3F_or_not/?hello=42&hello=54&vegetable=potato%23fresh
```

#### 参考文献　

[关于URL编码](http://www.ruanyifeng.com/blog/2010/02/url_encoding.html)  
[Percent-encoding](https://en.wikipedia.org/wiki/Percent-encoding)  
[HTML URL Encoding Reference](https://www.w3schools.com/TagS/ref_urlencode.asp)
[Encode / decode URLs](https://stackoverflow.com/questions/13820280/encode-decode-urls)
