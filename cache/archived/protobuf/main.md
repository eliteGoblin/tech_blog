


#### 思路

*  why protobuf: 基础的tcp unmarshal
    -  没有protobuf如何用tcp实现 unmarshal; 参考 Beej's Guide
    -  仿造nsq 构建tcp server, 序列化 version等
*  protobuf 如何使用
*  grpc对比 golang 内置　rpc(gob) 可以抓包; 或者json rpc;

#### setup

[Install protobuf3 on Ubuntu](https://gist.github.com/sofyanhadia/37787e5ed098c97919b8c593f0ec44d8)

#### raw tcp

*  byte order
    -  与string无关: 
    -  抓包显示tcp传送内容,预期
        +  string 一致
        +  network byte order
        +  float怎么表示? 
    -  7.4. Serialization—How to Pack Data
    ```c++
    htons() host to network short
    htonl() host to network long
    ntohs() network to host short
    ntohl() network to host long
    ```
    -  no similar htons() functions for float types.


#### protobuf

*  marshal && unmarshal完成

tree ./src/frank/

```shell
./src/frank/
├── address.proto
└── proto_defs
    └── address.pb.go
```

Streaming Multiple Messages

If you want to write multiple messages to a single file or stream, it is up to you to keep track of where one message ends and the next begins. The Protocol Buffer wire format is not self-delimiting, so protocol buffer parsers cannot determine where a message ends on their own. The easiest way to solve this problem is to write the size of each message before you write the message itself. When you read the messages back in, you read the size, then read the bytes into a separate buffer, then parse from that buffer. (If you want to avoid copying bytes to a separate buffer, check out the CodedInputStream class (in both C++ and Java) which can be told to limit reads to a certain number of bytes.)

#### gRPC

传送大量数据benchmark

*  tcp 吞吐量
*  gRPC





#### refs

ref0: [A practical guide to protocol buffers (Protobuf) in Go (Golang)](http://www.minaandrawos.com/2014/05/27/practical-guide-protocol-buffers-protobuf-go-golang/)
[Protocol Buffer Basics: Go](https://developers.google.com/protocol-buffers/docs/gotutorial)
[Language Guide (proto3)](https://developers.google.com/protocol-buffers/docs/proto3)
[Install protobuf 3.3 on Ubuntu 16.04](https://gist.github.com/rvegas/e312cb81bbb0b22285bc6238216b709b)  


[Protocol Buffer Basics: Go](https://developers.google.com/protocol-buffers/docs/gotutorial)
[Practical Golang: Using Protobuffs](https://jacobmartins.com/2016/05/24/practical-golang-using-protobuffs/)
[How we use gRPC to build a client/server system in Go](https://medium.com/pantomath/how-we-use-grpc-to-build-a-client-server-system-in-go-dd20045fa1c2)