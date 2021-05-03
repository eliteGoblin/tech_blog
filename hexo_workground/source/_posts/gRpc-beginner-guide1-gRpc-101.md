---
title: 'gRpc beginner guide1: gRpc 101'
date: 2018-10-03 16:01:05
tags: [golang gRpc]
keywords:
description:
---


{% asset_img gRPC.png %}

## Preface 


I bet you've heard people talking about RPC, gRPC all the time. You may already got questions like: What are they, why we need them in my system? seems like the RESTful API works just fine for us. How I should implement it, hope will not be too complicated.  

In this WHAT/WHY/HOW style blog, I am gonna first illustrate some key about RPC/gRPC, and then show how gRPC works by a simple example, hope you enjoy it. 

<!-- more -->

## What is RPC/gRPC

Some of you may know RPC is stand for *Remote Procedure Call*, it's a mechanism to communicate between processes, which enable programmer to invoke procedures(functions) reside on remote computer, but code just quite as simple as calling local functions.    

That make our life much easier, right? Because in our code, we do not need to worry about the network programming details, just call a local object's method(which is a proxy of remote service, called stub) with some input parameters, and get result back, like calling the local functions!

What is gRpc anyway? It's an implementation of RPC mechanism, or we can say it's an concrete object of RPC class. The initial g may come from Google who invented it.  

## Why we need gRPC

Why we need RPC?  As discussed, RPC made our day easier because it handles a lot of networking communication details, which make access of remote service just like local function call.  I agree with the metaphore that in a microservice system, each service is indispensable part of human body, RPC is like veins/neural connecting them. 

Why we choose gRPC?  gRPC use [protobuf](https://elitegoblin.github.io/2018/06/01/using-protobuf-in-golang/) to transmit data, much faster; And gRPC is quite powerful and versatile: it is built on HTTP/2: so proxies, firewalls and streaming reside in gRPC nature; gRPC integrates seamlessly with ecosystem components like service discovery, name resolver, load balancer, tracing and monitoring. 

## How we use gRpc in Golang

### Setup protobuf environment

As I wrote in [here](https://elitegoblin.github.io/2018/06/01/using-protobuf-in-golang/), you need to set up your protobuf environment first, gRPC use protobuf to transfer data

in summary, you need to:  

```linux
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
# Set up your GOBIN and GOPATH env
export GOBIN=xxx
export GOPATH=xxx
# Install protoc compiler plug-in for golang
go get -u github.com/golang/protobuf/protoc-gen-go
# Put protoc-gen-go in your PATH env
export PATH=$PATH:your_path_of_protoc-gen-go
```

### Download gRPC golang library code

```
cd $GOPATH
go get google.golang.org/grpc
```

### Hello World Example

We are gonna work through official's example which reside in *google.golang.org/grpc/examples/helloworld*:  

```
├── greeter_client
│   └── main.go
├── greeter_server
│   └── main.go
├── helloworld
│   ├── helloworld.pb.go
│   └── helloworld.proto
```

*  greeter_client is a rpc client says hello to remote server
*  greeter_server is a rpc server receive client's simple hello request then give response: just a request/reply way of communicating
*  in helloworld directory
    -  helloworld.proto: rpc protocol definition file, using the proto definition syntanx
    -  helloworld.pb.go: rpc library used by both RPC server and client 

let's run RPC server first:  

```
go run greeter_server/main.go
```

then run client which simply say *hello* to the server:  

```
go run greeter_client/main.go frank
```

we got:  

```
2018/10/03 17:39:36 Greeting: Hello frank
```

We just succeeded in running our first gRPC app!

### Some Illustration 

in helloworld.proto
```
...
// The greeting service definition.
service Greeter {
  // Sends a greeting
  rpc SayHello (HelloRequest) returns (HelloReply) {}
}

// The request message containing the user's name.
message HelloRequest {
  string name = 1;
}

// The response message containing the greetings
message HelloReply {
  string message = 1;
}
```

we can see that only difference with pure protobuf definition is we define a service, which use rpc keyword to indicate methods of a given service, actually, we can see gRPC as communicating with remote service which provide a couple of methods to invoke. 

on the client side, what we just need to do is dial and invoke client.SayHello, on server side, we need to implement the Server interface according to what we need:  

```golang
type server struct{}
// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
    return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}
```

and then register it the service:  

```golang
s := grpc.NewServer()
pb.RegisterGreeterServer(s, &server{})
// Register reflection service on gRPC server.
reflection.Register(s)
```

In RPC world, when you want to provide a service on the server, you need to **register** it on server, so clients can find out it by server address and service name

### Make your change

The helloworld example works directly because the source we got already has gRpc server/client code generated, let's try change it a little bit, add another method for Greeter service:  

```
rpc SayHelloAgain(HelloRequest) returns (HelloReply) {}
```

Basicly, it is quite similiar with SayHello, same input and output. We need to regenerate gRpc code:  

```
cd src/google.golang.org/grpc/examples/
protoc -I helloworld/ helloworld/helloworld.proto --go_out=plugins=grpc:helloworld
```

*plugins=grpc* means use grpc plugin to generate grpc code while processing proto file.

in greeter_server/main.go, add implementation for SayHelloAgain method: 

```
func (s *server) SayHelloAgain(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
    return &pb.HelloReply{Message: "Hello Again" + in.Name}, nil
}
```

in greeter_client/main.go, invoke SayHelloAgain: 

```
rAg, errAg := c.SayHelloAgain(ctx, &pb.HelloRequest{Name: name})
```

run the server, run client with go run greeter_client/main.go frank, you'll see the SayHelloAgain method is correctly invoked:  

```
2018/10/03 17:39:36 Greeting: Hello frank
2018/10/03 17:39:36 Greeting Again: Hello Againfrank
```

the original example could be found [here](https://github.com/grpc/grpc-go/tree/master/examples)

## Conclusion

In this blog we explained what is RPC/gRPC, the reason why we need it: 

*  Easy to use, programmer do not need to go through a lot of network, data encode/decode details
*  Efficient, because it transfers binary data
*  Easy to extend, because it is built on HTTP/2

And we use official's example to show how gRPC works.  

## Reference

[Remote procedure call](https://en.wikipedia.org/wiki/Remote_procedure_call)  
[Why gRPC uses Protobuf and HTTP/2 together?](https://github.com/grpc/grpc/issues/6292)  
[Go Quick Start](https://grpc.io/docs/quickstart/go.html#update-a-grpc-service)  
[What is gRPC?](https://grpc.io/docs/guides/)  
[gRPC Basics - Go](https://grpc.io/docs/tutorials/basic/go.html#generating-client-and-server-code)  