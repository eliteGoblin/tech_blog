---
title: using protobuf in golang
date: 2018-06-01 11:32:33
tags: [protobuf,golang]
keywords: protobuf golang
description:
---

<div align="center">
{% asset_img serialize.jpg %}
</div>


#### Why, What, How?

做后端的同学会对protobuf时有耳闻，可能也知道它是一种数据交换结构(标准)，那我们通讯时直接用TCP/UDP或者HTTP传输数据就可以了呗，为什么需要用它呢？本文从介绍protobuf解决什么问题开始，如何搭建开发环境，并以两个简单例子说明数据是如何经由protobuf交换．最后以一个TCP通信的例子综合讲述实际工程中protobuf的使用场景，相信实践完本文提供的例子，读者会对在golang开发中如何使用protobuf有初步认识．  

<!-- more -->

##### why we need things like protobuf

用HTTP通信时，可以很方便的以json的形式传送struct对象，但如果我们使用的是如TCP,UDP，甚至串口这些的面向字节流的协议呢？我们需要处理一系列很繁琐的问题:  

*  传输整数时的[byte order][1]问题:  TCP/IP传输时，需要先将字节转为netword byte order(big endian)，接受方收到数据后需要将其转为本机byte order
*  传输浮点数时，除了考虑字节序，还需要考虑浮点数在不同平台表示差异及相互转换
*  传输struct的object，需要指明各字段的类型，占用空间大小等

传输数据这个看起来*简单*的事情其实一点也不简单，有没有一种透明的机制使得我们只需要在发送端输入object，接受方能正确无误的提取到这个object呢？protobuf就是一种被广为接受的解决方案.

##### what is protobuf

全称google protocol buffers  

*  protobuf是一种将object(多数语言会提供struct类型，其对象便是object)marshal成字节流，并从字节流unmarshal成object的一种透明机制（或成为encode/decode
*  支持多平台，多语言间marshal/unmarshal
*  多用在传输数据时，减轻两端业务代码的数据编码/解析负担．
*  二进制编码，高效

##### how protobuf works

0.  protobuf用.proto文件的形式来约定object(protobuf称为message)的数据结构，在下面我们会通过例子讲述如何编写，也可以参考[proto3语法文档][3]，如
    ```protobuf
    message Person {
      string name = 1;
      int32 id = 2;  
      string email = 3;
    }
    ```
1.  protobuf提供编译器protoc，输入.proto文件，自动生成特定语言的数据结构定义程序代码．若是给golang用，还需安装编译器插件protoc-gen-go，生成的文件是\*.pb.go  
2.  有了定义好的数据*.pb.go，我们可以用protobuf的golang library来在代码中实现marshal和unmarshal操作(我们的**终极目标**)

#### 配置本文所需环境

根据上一节描述，我们需要三个东西:  

*  protoc: .proto文件编译器
*  protoc-gen-go:  golang的protoc插件
*  golang protobuf library

我的环境是ubuntu16.04, go1.9.2

##### install protoc

```shell
# Make sure you grab the latest version
curl -OL https://github.com/google/protobuf/releases/download/v3.3.0/protoc-3.3.0-linux-x86_64.zip

# Unzip
unzip protoc-3.3.0-linux-x86_64.zip -d protoc3

# Move protoc to /usr/local/bin/
sudo mv protoc3/bin/* /usr/local/bin/

# Move protoc3/include to /usr/local/include/
sudo mv protoc3/include/* /usr/local/include/

# Optional: change owner
sudo chown $USER /usr/local/bin/protoc
sudo chown -R $USER /usr/local/include/google
```

安装成功的标志:  
*  protoc命令可以执行  
*  /usr/local/include/google/protobuf存在并有文件  

##### install protoc-gen-go

*  设置环境变量$GOBIN
*  执行 go get -u github.com/golang/protobuf/protoc-gen-go

执行成功会发现$GOBIN目录下多了文件protoc-gen-go，命令行可以运行．它是protoc的golang必须插件，生成*.pb.go用的还是使用protoc命令．  

##### install protobuf golang library

golang library主要用到其提供的marshal和unmarshal功能，与protoc-gen-go生成数据结构定义文件*.pb.go必须结合使用才能实现marshal/unmarshal自己的object，缺一不可．  
library本质是我们的代码能import的package，需要安装到业务代码的GOPATH中．api文档见[package proto](https://godoc.org/github.com/golang/protobuf/proto)

安装步骤  

*  设置我们代码的GOPATH
*  go get "github.com/golang/protobuf/proto"

到此为止protobuf的环境就一切准备就绪了(真不容易...)，接下来我们通过两个例子看看protobuf工作过程．　　

#### protobuf by example

代码在[这里](https://github.com/eliteGoblin/Notes/tree/master/cs/Languages/go/protobuf/src/frank)

例子配置步骤:

下载我github目录cs/Languages/go/protobuf下的所有代码到本地，假设下载完毕，地址为: **YOUR_PATH**/protobuf(参见[如何从 GitHub 上下载单个文件夹](https://www.zhihu.com/question/25369412))，执行如下命令:  

```shell
cd YOUR_PATH/protobuf
# 设置GOPATH为此目录
export GOPATH=`pwd`
```

##### 代码结构简介

```shell
cd YOUR_PATH/protobuf
tree ./src/frank/
```
得到如下结构: 

```shell
./src/frank/
├── address.proto          # 简单Person结构一节所用的数据结构定义文件
├── full_addr
│   ├── full_addr.pb.go    # 复杂AddressBook结构自动生成的.pb.go文件
│   └── full_addr_test.go  # 复杂AddressBook结构的演示代码在此
├── full_addr.proto　       # 复杂AddressBook结构的.proto文件
├── main.go                # 简单Person结构演示代码在此
└── proto_defs             # 简单Person结构生成的proto.go文件
    └── address.pb.go
```

##### 简单Person结构

address.proto定义的Person message type:  

```proto
syntax = "proto3";      // 指明版本，有proto2和proto3可选
package proto_defs;     // 指明生成的.pb.go文件package

message Person {
  string name = 1;      // 内置string类型
  int32 id = 2;         // Person唯一Id,内置int32
  string email = 3;
}
```

说明:  

*  field后面数字为Field Numbers:
    -  对于一个message type是唯一的，标识field; 范围[1,2^29 - 1];
    -  1-15和fieldtype需要1byte，其他需要2byte([15,2047])或更多;
    -  最频繁的选[1,15]优化传输
*  protobuf语言规范内置类型有: double, float, int, int32.. 与一般语言很相近，具体参考[官方protobuf 语言规范][3]  

生成我们可用的 address.pb.go文件:  

```shell
cd YOUR_PATH/protobuf/src/frank
protoc -I=. --go_out=./proto_defs address.proto
```
我们的proto_defs/address.pb.go就是这么自动生成，这个文件我们也不应该手动去修改，如何使用呢？在main.go中:  

```golang
person := new(proto_defs.Person)
...
cache, err := proto.Marshal(person)   // object到byte array的marshal
var targObj = new(proto_defs.Person)  
err = proto.Unmarshal(cache, targObj) // byte array到object的unmarshal
...
```

虽然protobuf的准备阶段稍显繁琐，但是在使用时很简洁，和json的marshal/unmarshal很类似． 

这个繁琐也是值得的，有人对比protobuf和json，当传输的是integer时，marshal和unmarshal速度是后者的3倍和8倍，传输float时，速度差异会更明显:  

{% asset_img vsjson.jpg %}

##### 复杂AddressBook结构

关键在于.proto文件的定义方式，使用上仍是简单的marshal/unmarshal，例子来自google go protobuf basic:  

```proto
syntax = "proto3";
package full_addr;
message Person {
  string name = 1;
  int32 id = 2;  
  string email = 3;
  enum PhoneType {
    MOBILE = 0;
    HOME = 1;
    WORK = 2;
  }
  message PhoneNumber {
    string number = 1;
    PhoneType type = 2;
  }
  repeated PhoneNumber phones = 4;
}
// Our address book file is just one of these.
message AddressBook {
  repeated Person people = 1;
}
```

说明:  

*  数据展现的是address book, 分为3层结构
    -  AddressBook由Person组成，是Person数组
    -  Person除了name,id等信息，包含有PhoneNumber数组
*  repeated 代表出现0次和多次，对应slice(数组)
*  enum表示枚举类型，与golang的enum对应，protoc已经方便的为我们生成了enum结构和辅助代码(摘自自动生成的full_addr.pb.go): 
    ```golang
    type Person_PhoneType int32
    const (
        Person_MOBILE Person_PhoneType = 0
        Person_HOME   Person_PhoneType = 1
        Person_WORK   Person_PhoneType = 2
    )
    var Person_PhoneType_name = map[int32]string{
        0: "MOBILE",
        1: "HOME",
        2: "WORK",
    }
    var Person_PhoneType_value = map[string]int32{
        "MOBILE": 0,
        "HOME":   1,
        "WORK":   2,
    }
    ```

full_addr.pb.go生成和测试代码与上一节相似．  

#### 用TCP传输AddressBook

有了protobuf方便的marshal和unmarshal机制，我们可以将其用在tcp通讯中，示意代码:  

client.go
```golang
obj := new(AddressBook)
...  // fill AddressBook object
// marshal AddressBook object to byteArray
byteArray, _ := proto.Marshal(AddressBook)
// write to remote by tcp
conn, __ := net.Dial("tcp", "localhost:8888")
conn.Write(byteArray)
```

server.go
```golang
data := make([]byte, 4096)
n, err := conn.Read(data)
protodata := new(ProtobufTest.AddressBook)
proto.Unmarshal(data[0:n], protodata)
```

完整的代码来自这个[例子](4)

#### Conclusion

*  protobuf以远快于json的速度(一般marshal/unmarshal分别约快3倍，8倍)应用在很多数据传输场景，如TCP/UDP, gRPC等  
*  本文说明了配置protobuf开发环境的详细步骤，并以2个例子说明如何用protobuf序列化/反序列化程序数据结构，这也是protobuf要解决的问题  
*  最后借用一个C/S通信的例子说明protobuf如何在实际通信中应用  

#### refs

[1. Beej's Guide to Network Programming][1]  
[2. A practical guide to protocol buffers (Protobuf) in Go (Golang)][4]  
[3. Protocol Buffer Basics: Go](2)  
[4. Language Guide (proto3)](3)  
[5. Is Protobuf 5x Faster Than JSON?][5]  


[1]: https://beej.us/guide/bgnet/html/single/bgnet.html  
[2]: https://developers.google.com/protocol-buffers/docs/gotutorial  
[3]: https://developers.google.com/protocol-buffers/docs/proto3
[4]: http://www.minaandrawos.com/2014/05/27/practical-guide-protocol-buffers-protobuf-go-golang/
[5]: https://dzone.com/articles/is-protobuf-5x-faster-than-json