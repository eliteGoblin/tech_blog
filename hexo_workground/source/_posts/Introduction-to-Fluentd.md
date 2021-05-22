---
title: Fluentd and Graylog logging solution in k8s
subtitle: This is subtitle
date: 2021-05-04 17:14:21
tags: [logging, Fluentd]
---

<img src="https://raw.githubusercontent.com/eliteGoblin/images/master/blog/img/picgo/20210504171645.png" alt="20210504171645" style="width:500px"/>

## Why 

搭建于k8s的微服务系统，log aggregation更是一项基本需求: pod随时可能在Node见迁移，找不到原始log文件. 

K8s官方并没有给出具体logging solution, 只给出了[logging architecture](https://kubernetes.io/docs/concepts/cluster-administration/logging/)  

这里介绍一套我们组在用的solution: Fluentd+Graylog, 配置简单(相比EFK), 且提供对应的golang logging library实现. 

即本文解决问题: 收集golang backend产生的日志，并aggregate到graylog.

<!-- more -->

## 本Logging solution的high level workflow

1.  app pod产生日志，由k8s写入node的`/var/log/containers` folder下
2.  Fluentd作为daemonset, 运行于k8s cluster, parse log under `/var/log/containers`
3.  Fluentd获取log作为input，经过内部workflow: filter, match, output, 最终将日志发送给Graylog(通过TCP 协议)
4.  Graylog 从listen的TCP port, 获取log消息，建立ES索引，并能从Graylog dashboard看到对应logging message.

本主题将分解为系列博客: 

1.  [Introduction to Fluentd](https://blog.franksun.org/2021/05/04/Introduction-to-Fluentd/)  
2.  [Setup Graylog in K8s and Istio](https://blog.franksun.org/2021/05/18/Setup-Graylog-in-K8s-and-Istio/)  
3.  [Structure logging in Golang](https://blog.franksun.org/2021/05/19/Structure-logging-in-Golang/)  

本文代码测试环境: 
*  GKE 1.17
*  Istio: 1.7
*  Fluentd: 1.12.3

## Introduction to Fluentd

<img src="https://raw.githubusercontent.com/eliteGoblin/images/master/blog/img/picgo/20210504175255.png" alt="20210504175255" style="width:500px"/>  

从官方架构图可看出，Fluentd作为适配input(主要是各种logging message), output(多种多样: 各种持久化数据库，monitoring system)的中间件, 提供可定制化的: filter, routing, parse, buffer, format等; 最终实现input, output解耦. 

Fluentd由Ruby实现，生态里包含很多开源[Plugin](https://www.fluentd.org/plugins/all). 

Fluentd接受到的data, 化为内部的event

## A Fluentd event's lifecycle

### 本地测试环境

[docker-compose-http](https://github.com/eliteGoblin/code_4_blog/blob/master/fluentd_graylog/fluentd_1/docker-compose-fluentd-http.yaml)  

```yaml
version: '3.2'
services:
  # debug network util
  tool:
    image: praqma/network-multitool
    networks:
      - graylog
  fluentd:
    image: fluent/fluentd:v1.12-1
    volumes:
      - ./conf:/fluentd/etc/
    ports:
      - "8888:8888"
      - "24224:24224"
      - "24224:24224/udp"
networks:
  graylog:
    driver: bridge
```

Fluentd conf: 
```conf
<source>
  @type http
  port 8888
  bind 0.0.0.0
</source>

<match test.cycle>
  @type stdout
</match>
```

*  配置HTTP 输入(plugin `in_http`): 根据不同路径赋予label
*  Match label， 并输出到stdout(plugin `out_stdout`)
*  Fluentd将不同Plugins组装成message processing pipeline

### A HTTP event's life cycle

任何Fluentd的输入都被转化为一系列events， event结构: 

*  `tag`: 表明events的origin, 用于后续message routing
*  `time`: event到达的timestamp
*  `record`: actual message payload

整个lifecycle一般分为如下几步: 
0.  Setup: config Fluentd, 设置好pipeline
1.  Input: 配置input plugin, 获取输入; 不同的输入有不同的`tag`
2.  Filter: filter和rewrite events
3.  Matches: match label, 并通过output plugin, output

这里通过HTTP event来阐述life cycle, input一个event via HTTP:

```bash
curl -i -X POST -d 'json={"action":"login","user":2}' http://localhost:8888/test.cycle
```

Fluentd输出: 

```sh
fluentd_1  | 2021-05-17 22:47:54.452096657 +0000 test.cycle: {"action":"login","user":2}
```

*  tag: test.cycle, from URL path
*  `record`: 实际输入

### Input from log file

`sample_log/source_photos.log` 几条取自app的json log: 

通过input plugin: `in_tail`: 

```conf
<source>
  @type tail
  path /var/log/*.log
  pos_file /tmp/app.log.pos
  tag api.app.*
  format json
  read_from_head true
</source>

<match api.app.**>
  @type stdout
</match>
```

Source为HTTP输入， 并赋予tag: `api.app`, 同样match此label, 输出到stdout:

```bash
fluentd_1  | 1970-01-01 00:33:41.000000000 +0000 api.app.var.log.source_photos.log: {"elevation":0,"latitude":-27.4748854,"longitude":153.0279261,"message":"Metadata Request Query Parameters","request-id":"bb5988fe-39df-4dea-b820-acd828368c3c","severity":"info","surveyID":"a8f9130c-b1be-11ea-b13b-f32156d4454a","url":"https://api-qa.nearmapdev.com/photos/v2/surveyresources/100-a8f9130c-b1be-11ea-b13b-f32156d4454a/photos/ground/153.0279261,-27.4748854/metadata.json?apikey=xxx"}
fluentd_1  | 1970-01-01 00:33:41.000000000 +0000 api.app.var.log.source_photos.log: {"message":"Preview Request Query Parameters","request":{"ID":"d5a0b536-b1be-11ea-8a8b-bfd11cfa842d","SurveyID":"a8f9130c-b1be-11ea-b13b-f32156d4454a","ImageType":1,"X":0,"Y":0,"Z":0,"Width":250,"Height":250,"Background":{"R":0,"G":0,"B":0,"A":0},"Photo":null,"SkipExistingPreviews":false,"GenerationID":""},"request-id":"4bb710e8-1fb6-4a89-9c01-9458608853c0","severity":"info","url":"https://api-qa.nearmapdev.com/photos/v2/surveyresources/100-a8f9130c-b1be-11ea-b13b-f32156d4454a/photos/d5a0b536-b1be-11ea-8a8b-bfd11cfa842d/preview/250x250.jpg?apikey=xxx"}
fluentd_1  | 1970-01-01 00:33:41.000000000 +0000 api.app.var.log.source_photos.log: {"error":"Failed to SET redis value: context canceled","message":"Failed to add value to cache","request-id":"ce7513cf-aafe-4825-b174-cfb3216ad144","severity":"error","url":"https://api-qa.nearmapdev.com/photos/v2/surveyresources/100-a294c788-f986-11e7-82f2-67f998099448/photos/70fba45c-f987-11e7-aacc-9be7b8095f12/0/0/0/512x512.jpg?apiKEY=xxx"}
```

### Labels to jump in message pipeline

一般message会根据top-down依次走完config定义好的各步骤, 如果想"跳跃"pipeline, 需要Fluentd的`labels`机制. 

```
<source>
  @type http
  bind 0.0.0.0
  port 8888
  @label @STAGING
</source>

<filter test.cycle>
  @type grep
  <exclude>
    key action
    pattern ^login$
  </exclude>
</filter>

<label @STAGING>
  <filter test.cycle>
    @type grep
    <exclude>
      key action
      pattern ^logout$
    </exclude>
  </filter>

  <match test.cycle>
    @type stdout
  </match>
</label>
```

*  在Source中指定label
*  Source完毕直达label section, 跳过filter步骤

<img src="https://raw.githubusercontent.com/eliteGoblin/images/master/blog/img/picgo/20210518082921.png" alt="20210518082921" style="width:500px"/> 


可用[calyptia](https://config.calyptia.com/#/visualizer)来visualize Fluentd/FluentBit's config


### Docker logging driver

Application容器化之后，程序输出到stdout, stderr的内容会被docker logging driver捕获并处理. 假设fluentd container name为`fluentd_1_fluentd_1`

`docker logs fluentd_1_fluentd_1`, 同样可以看到输出到stdout的logs, 参见[View logs for a container or service](https://docs.docker.com/config/containers/logging/)  

Docker logging driver支持多种模式， 见[Configure logging drivers](https://docs.docker.com/config/containers/logging/configure/)

Kubernetes下，一般默认config为[json-file](https://docs.docker.com/config/containers/logging/json-file/), 即docker logging driver捕获STDOUT, STDERR输入，写入到json文件: 

*  One line per log, 默认不支持multiline log
*  Json file logging driver写入format会对捕获的message进一步封装: 无论log的格式是什么，将其作为string, 每行转化为一个json, 包含3 fields: `log`, `stream`, `time`.
```javascript
{
  "log": "Log line is here\n",
  "stream": "stdout",
  "time": "2019-01-01T11:11:11.111111111Z"
}
``` 

*  `json-file`并不指期待输入的log format为json, 而是将每行内容转化为json
*  `docker info --format '{{.LoggingDriver}}'`用来查看docker配置的logging driver.

### Send logs to remote backend

Fluentd作为logging的关键一环，最终目的是将log发送给remote logging backend, 常见的有: 

*  EFK stack(Elasticsearch + Fluentd + Kibana): Fluentd发送给Elastic Search, Kibana作为UI.
*  Graylog: 本身提供UI, 后端采用ES及mongodb.

官方文档提供了将Log直接送到EFK的container, 见[这里](https://docs.fluentd.org/container-deployment/docker-compose)

## Resource

[deployment on k8s](https://docs.fluentd.org/container-deployment/kubernetes)  
[fluentd-kubernetes-daemonset](https://github.com/fluent/fluentd-kubernetes-daemonset)  
[建好的Graylog docker](https://github.com/fluent/fluentd-kubernetes-daemonset/blob/master/docker-image/v1.12/debian-graylog/Dockerfile)  