---
title: Ports forwarding Basic in K8s
date: 2019-10-25 16:49:19
tags:
keywords:
description:
---

{% asset_img hub.jpg %}

## Preface

在K8s中，我们需要在很多地方映射端口，如

* Container port
* Service port
* Node port

常常会带来疑惑，在这里我们尝试梳理一下有关container, K8s port的方方面面，

<!-- more -->

## 根源: Container Port

K8s作为容器编排系统，是建立在容器技术之上的，因此我们先来看看docker port。

官方Doc:

> By default, when you create a container, it does not publish any of its ports to the outside world.

> To make a port available to services outside of Docker, or to Docker containers which are not connected to the container’s network, use the --publish or -p flag. This creates a firewall rule which maps a container port to a port on the Docker host. Here are some examples

简单的说就是可以通过publish命令来将docker内部port映射到host上，format: 主机端口:容器端口，几个example: 

* -p 8080:80
* -p 192.168.1.100:8080:80
* -p 8080:80/udp
* -p 8080:80/tcp -p 8080:80/udp

可以指定protocol, host ip等。

container ports涉及到的概念有publish, expose，我们来一一看一下


### 测试环境

[arunvelsriram/utils](https://github.com/arunvelsriram/utils)是一个集成了很多有用utils的image，如nc, psql, dig等，调试问题很方便。我们将以其为模板，只留下我们需要的工具，build我们自己的image来测试端口映射。

Dockerfile可以在[这里](https://github.com/eliteGoblin/code_4_blog/tree/master/utils_image)找到。

```
docker build -t elitegoblin/testports .
```

这条命令在我们本地(Ubuntu 18.04)build名为elitegoblin/testports的image。

我们用nc命令(netcat)来模拟监听某个端口的服务，对于调试网络很方便。

### UnPublished Port

Publish意味着将container port映射到host上，只有host admin才有权限，也意味着 **运行时** 进行的操作。

有个问题，没有publish的port，究竟能不能被外界访问呢？我们来测试一下：

我们通过netcat(nc)监听8888端口，但是并不将其publish

```
docker run --name srv --rm elitegoblin/testports nc -l 8888
```

通过`docker inspect srv`命令获得容器ip: 172.17.0.2，在另一个容器中运行curl来看srv的8888端口是否开放

```
docker run --rm -ti elitegoblin/testports bash
# inside container
# curl 172.17.0.2:8888
```

我们可以看到srv端收到HTTP请求，证明没有publish，仍可以被外界访问

```
GET / HTTP/1.1
Host: 172.17.0.2:8888
User-Agent: curl/7.58.0
Accept: */*
```

同样也可以在host上用curl访问到，这表明: 

* Publish是将container port映射到host
* 如果直接访问容器的IP，可以访问到没有publish的port
* Publish不需要先expose, publish unexposed port效果一样

### Expose

docker port的publish简单易懂，关于expose: 

> The EXPOSE instruction does not actually publish the port.  It functions as a type of documentation between the person who builds the image and the person who runs the container, about which ports are intended to be published. 

从文档里我们了解到expose是文档，也就是常常见于Dockerfile中，区别于publish作用于运行时，expose是image author侧的责任，表明此image需要publish哪些ports，但是并不会(也不能)publish，因为image作者并没有host上的任何权限。

> To actually publish the port when running the container, use the -p flag on docker run to publish and map one or more ports, or the -P flag to publish all exposed ports and map them to high-order ports.

Expose并不仅仅是文档，docker可以自动将expose的port map到host的一个random端口。

我们添加EXPOSE directive到我们的dockerfile: 

```
EXPOSE 8888
```

再运行docker
```
docker run --name srv --rm -P elitegoblin/testports nc -l 8888
```

通过`docker port srv`我们可以看到8888被映射到了host:

```
8888/tcp -> 0.0.0.0:32773
```
可以在host上运行curl

```
curl http://localhost:32773
```

同样可以看到HTTP请求到达srv

```
GET / HTTP/1.1
Host: localhost:32773
User-Agent: curl/7.58.0
Accept: */*
```

### Expose but not listen

我们来看一下如果container内部没有使用，但是EXPOSE有无影响，同样在Dockerfile中: 

`EXPOSE 8888`

启动container
```
docker run --name srv --rm -ti -P elitegoblin/testports bash
```
运行命令: `docker port srv`,发现虽然没有listen，仍然进行了映射：

```
8888/tcp -> 0.0.0.0:32775
```

在host上运行`curl http://localhost:32775`，不出意料的得到: `Connection reset by peer`，因为container内部并没有监听。

这时我们在容器内部启动nc监听8888端口，再次在host上curl，请求成功到达srv

```
GET / HTTP/1.1
Host: localhost:32775
User-Agent: curl/7.58.0
Accept: */*
```

说明EXPOSE作用是创建映射规则，并不要求container port被listen。

## Container port in K8s

在K8s环境，我们在deployment的config中常会见到关于container port [^1]

```
ports:
- containerPort: 3306
  name: mysql
```

K8s文档里写道: 

> List of ports to expose from the container. Exposing a port here gives the system additional information about the network connections a container uses, but is primarily informational. Not specifying a port here DOES NOT prevent that port from being exposed.

又一个informational，和EXPOSE极为类似，这个参数就是说明container内部监听了3306端口，有没有它对port是否能被外界访问到并没有影响，完全是为方便阅读者理解。

在K8s内部，其实我们并不关心port是否publish到host上，我们一般是: 

*  container/pod内部监听port
*  通过service访问一组pods
*  service之上完成service port到container port的映射

比如在pod内部我们想访问另一组pods时，对应服务记为example-svc: 

```
curl http://{example-svc-name}:{svc-port}/path/of/my/url
```

这就涉及到另一个转换: service到container port。

## Service Port

一个典型的multi-port service配置如下 [^2]

```yaml
apiVersion: v1
kind: Service
metadata:
  name: my-service
spec:
  selector:
    app: MyApp
  ports:
    - name: http
      protocol: TCP
      port: 80
      targetPort: 9376
    - name: https
      protocol: TCP
      port: 443
      targetPort: 9377
```

它对每个port分别定义了映射规则，这样当pods访问 `http://my-service/xx`时，会被redirect到pods的9376端口；而`http://my-service:443/xx`被redirect到pods的9377端口。

## Node Port

NodePort是一种向外界暴露service时方法：将service映射为host的port，在host上访问此node port会被redirect到内部的service，有点像docker里面的publish port。

这是另一层映射，主机port映射为service的port，我们来修改一下上述多端口的service的定义，增加字段`type: NodePort`: 

my-service.yaml
```yaml
apiVersion: v1
kind: Service
metadata:
  name: my-service
spec:
  selector:
    app: MyApp
  type: NodePort
  ports:
    - name: http
      protocol: TCP
      port: 80
      targetPort: 9376
    - name: https
      protocol: TCP
      port: 443
      targetPort: 9377
```
deploy service: 
```
kubectl apply -f my-service.yaml
```

查看my-service的详细信息 `kubectl describe services my-service`

```
Type:                     NodePort
IP:                       100.64.30.206
Port:                     http  80/TCP
TargetPort:               9376/TCP
NodePort:                 http  30938/TCP
Endpoints:                <none>
Port:                     https  443/TCP
TargetPort:               9377/TCP
NodePort:                 https  32709/TCP
```

我们看到两个三元组，和service定义对应:

* HTTP
    - port: 80
    - target port: 9376
    - node port: 30938
* HTTPS
    - port: 443
    - target port: 9377
    - node port: 32709

可见port映射了两次: node port --> port(service port) --> target port;

## Conclusion

通过分析，我们了解了：

*  在Container内部，EXPOSE有两个作用：指示container listen哪个端口，另外`docker run -P`会自动将EXPOSE的port　publish到host上(随机分配端口); EXPOSE并不控制，也就是对port实际是否能被访问没有作用
*  Deployment的containerPort字段也是informative，可以省略。
*  Service实现了port-->container port的映射
*  Service如果是NodePort，还有node port(host port) --> service port的映射。

有点小绕，画一个图压压惊: 

{% asset_img big_pic.png %}


[^1]: [Should I configure the ports in the Kubernetes deployment?](https://medium.com/faun/should-i-configure-the-ports-in-kubernetes-deployment-c6b3817e495)
[^2]: [Multi-port Services](https://kubernetes.io/docs/concepts/services-networking/service/#multi-port-services)