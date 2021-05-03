---
title: Golang HTTP server中的error propagation
date: 2020-04-04 12:45:01
tags: [golang, HTTP, error]
keywords: error propagation
description: 
---

## Preface 

<div align="center">
{% asset_img chain.jpg %}
</div>

错误处理是程序设计的重要环节，理想的错误处理应为开发人员提供详细的context，并记录日志;它排查问题的重要信息源。并能为用户提供恰到好处的信息: 适当抽象，隐藏程序细节。最好能提供解决建议，或者`reference id`, 协助开发人员缩小排查范围。

特别的, 如何在HTTP server中建立正确的error处理机制, 并能自动生成符合HTTP语义的response呢? 

本文提出一个简单解决方案: 提供error的stack trace, 并能将error转换为HTTP error.

<!-- more -->

## Error propagation

软件模块可能嵌套较多层级，一个error由low-level模块产生，一级级的传递到上层被处理，再反馈给用户，称为error propagation, 像环环相扣的链条. 即error以逆调用链方向传递(称为Wrap, 每层包裹一些context, 再传递出去)，并记录如代码位置, 函数名等关键信息，这样的error stack能还原出错时代码执行路径，对理解问题至关重要.

想想如果没有调用栈, 仅凭一条error message, 且往往是内部库函数产生的error,被很多代码调用, 根本无法理解程序执行路径.

Java有完善的Error handling, 能记录error stack: 

<div align="center">
{% asset_img exception_java.jpg %}
</div>

另一个问题: 在HTTP Server场景，为遵循HTTP标准，还需要考虑出error与HTTP response的转换: 涉及到定制HTTP status code, response header 即response body. 

## Error handling in Golang

Golang自诞生来，因其极简的设计，内部对error propagation并无支持. (直到Go 1.13加入初步支持，远非完美, 前身是`golang.org/x/xerrors`) 但有成熟的三方库可采用，如`github.com/pkg/errors`和`github.com/juju/errors`, 但error propagation思路大同小异.

Error propagation的套路便是由错误源头开始，在向上层传递错误时，层层"包裹", 附加context. 而在错误处理层, 又需要解开包裹，取出错误链条的第一个Error, 即错误源头，从而区分处理不同error。

一个函数调用链以源错误为开端，层层返回给调用层，同时附加调用栈和额外message，这个简单的模型足以处理一般系统的错误。

因此error库至少需要两个函数`Wrap`和`Cause`: `Wrap`用来包装error: 添加context形成新error;`Cause`, 也成为Unwrap则用于取的被包裹的原始错误。

接下来我们来看看如何用`github.com/juju/errors`来实现error propagation并且转换为HTTP response. 

## juju/errors

[juju/errors](https://github.com/juju/errors), github收获千星，提供了简单的error Unwrap逻辑，并内置了一些标准HTTP错误类型。

因为我们想做一些额外定制: 区分是否Server端错误，附加error时HTTP response header等，仅用到`juju/errors`的Wrap和Cause函数, 并没有使用其定义的错误类型，而选择我们自己定义能对应到HTTP标准错误的类型。

在`juju/errors`中

新建error, 替代标准errors函数，新建的error会记录call stack. 

```golang
New(message string) error
Errorf(format string, args ...interface{}) error
```

Wrap error:

```golang
Trace(other error) error
Annotate(other error, message string) error
```

Trace仅Wrap, Annotate会增加额外message. 

```golang
func TestJuju(t *testing.T) {
	err := errors.New("root cause") // line 11
	err = errors.Trace(err) // line 12
	err = errors.Annotatef(err, "wrap with annotate\n") // line 13
	fmt.Printf("%s", errors.ErrorStack(err))
}
```

输出

```shell
errors_propagation/errors/juju_test.go:11: root cause
errors_propagation/errors/juju_test.go:12: 
errors_propagation/errors/juju_test.go:13: wrap with annotate
```

可见Error stack打印出带有函数行数的调用链。

打印Cause: 

```golang
fmt.Printf("Cause is %+v\n", errors.Cause(err))
```

输出root cause: 

```
errors_propagation/errors/juju_test.go:11: root cause
```

可见用errors.New创建的error,带有context信息, 而非`juju/errors`创建的信息则没有, 可以用`errors.Trace`进行wrap. 

测试代码可以在[这里](https://github.com/eliteGoblin/code_4_blog/blob/master/errors_propagation/errors/juju_test.go)找到. 

## error转换为HTTP error response

一般情况是: HTTP server的handler layer接到请求，转入business logic; 层层调用之后，某处报错，错误再返回到handler层，这时需要将错误转换为符合HTTP语义的response. 如: 权限不足返回`403`, 数据未找到返回`404`, 程序崩溃返回`500`. 

有了Error propagation为基础，我们可以提取root cause, 我们需要定义一系列标准HTTP error, 能从error中提取HTTP status code, response header及error message. 

以Status code为例，我们自定义错误需要实现`HTTPStatus() int` interface, 错误处理模块根据interface转换来判断返回的错误是否能提取status code, 同理还有是否是Service Failure等其他功能. 

例如404错误，我们定义: 

```golang
type BadRequest struct {
	message string
	code    string
}

func (br BadRequest) Error() string {
	return br.message
}

func (br BadRequest) Code() string {
	return br.code
}

func (br BadRequest) HTTPStatus() int {
	return http.StatusBadRequest
}

func (br BadRequest) IsServiceFailure() bool {
	return false
}

func NewBadRequest(code, msg string, args ...interface{}) *BadRequest {
	return &BadRequest{
		code:    code,
		message: fmt.Sprintf(msg, args...),
	}
}
```

处理错误时, 利用type conversion来判断是否是我们定义的标准HTTP错误, 既能提取status code.

```golang
if hs, ok := errors.Cause(err).(HasHTTPStatus); ok {
	return hs.HTTPStatus()
}
```

同样，我们可以实现HTTP response header等其他功能.

## Put it together

Error转换为HTTP response机制需要在所有HTTP handler共享, HTTP handler仅需将返回的error传入转换代码，转换模块负责写HTTP status, header, body等. 

一个好办法是function adapter， 将标准的http.Handler添加error返回值，adapt过程中处理error

```golang
func handleExample(w http.ResponseWriter, r *http.Request) error {
    ...
    err = businessLogic()
    if err != nil {
        return errors.Trace(err)
    }
    return nil
}
```

声明adaper, 调用自定义的`WriteError`函数转换error并写HTTP response:

```golang
type WebHandlerFunc func(http.ResponseWriter, *http.Request) error

func AddErrorHandler(handler WebHandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			logrus.Errorf("%+v", errors.ErrorStack(err))
			WriteError(w, err)
		}
	})
}
```

## Example: 模拟read user info的各种可能错误

这里带给大家的例子是: 一个`GET /users?userID=xxx`endpoint, 不同的userID不同的HTTP response, 来模拟真实HTTP server遇到的情况, 404 user not found, 403 forbidden, 500 server error, 正常返回等，可以通过这个精简的例子来看error propgation, 即error是如何转为标准HTTP response的. 

例子代码在[^1]

### 运行demo

```
go build . && ./error_propagation
```

试返回404的user

```
curl -v http://localhost:8090/users\?userID\=notfountindb
```

返回
```
< HTTP/1.1 404 Not Found
< Content-Type: application/json
< Date: Sat, 04 Apr 2020 13:11:40 GMT
< Content-Length: 42
< 
* Connection #0 to host localhost left intact
{"error":"userID: notfountindb not found"}% 
```

同时程序打印error stack: 

```
{
  "level": "error",
  "msg": "userID: notfountindb not found\n/home/frank.sun/git_repo/code_4_blog/errors_propagation/user.go:67: \n/home/frank.sun/git_repo/code_4_blog/errors_propagation/user.go:31: ",
  "time": "2020-04-05T00:11:40+11:00"
}
```

### 代码分析

产生错误时，new我们定义的标准HTTP错误`NewBadRequest`, 同时Wrap: 加入错误产生的行数. 

```golang
func checkUserID(userID string) error {
	if userID == badUserID {
		return errors.Trace(
			se.NewBadRequest("INVALID_USERID", "malformed userID provided: %s", badUserID))
	}
	return nil
}
```

HTTP handler中发现错误, Wrap并返回: 

```golang
func handleUser(w http.ResponseWriter, r *http.Request) error {
    userID := r.URL.Query().Get("userID")
    err := checkUserID(userID)

    if err != nil {
        return errors.Trace(err)
    }
    ...
}
```

, 错误处理层`AddErrorHandler`负责写HTTP response, 核心是调用转换函数, 写Status Code, response header及标准化的error body: 包含`error code`和`error message` 

```golang
func WriteError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	status := se.HTTPStatusCode(err)

	extraHeaders := se.ResponseHeaders(err)
	for k, v := range extraHeaders {
		w.Header().Set(k, v)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(marshalError(err, status))
}
```

## 总结

我们分析error处理的一般思路: error propagation; 及如何将模块error propagate到用户: 转化为HTTP error response. 并结合demo给出了处理框架. 

这套处理思路来自现在任职公司Nearmap, API组(这是一家好公司!); 是经过实际验证，靠谱的解决方案. 同时也解决了之前自己的疑问: 如何在HTTP server中妥善处理error, 在此表示感谢. 


[^1]: [Demo code](https://github.com/eliteGoblin/code_4_blog/blob/master/errors_propagation/)