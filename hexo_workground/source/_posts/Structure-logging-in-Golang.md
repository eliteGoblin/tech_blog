---
layout: post
title: Structure logging in Golang
date: 2021-05-20 05:20:05
tags:
---

## 背景

接[Logging系列文章](https://blog.franksun.org/2021/05/04/Introduction-to-Fluentd/), Structure logging是当下best practice. 如何在Golang实现呢? 

本文给出solution, 取自当前我们组使用的logging library, 实现以下需求: 

*  Structure logging in JSON, 1 log per line
*  输出logging到stderr
*  不同level: Debug, Info, Warn, Error
*  结合context, 传递带额外信息的logger: e.g 同一HTTP request, 附带同一`request_id`

<!-- more -->

## 标准化Log message

我们这里规定Application JSON log必须有的字段, 这里不直接采用GELF格式，解耦考虑. 

*  time: log 产生时间, `RFC3339`
*  severity: `debug`, `info`, `warning`, `error`, `fatal`, `panic`
*  message: message body, 不应该embedded额外信息
*  其他Custom fields不限

基于[sirupsen/logrus](https://github.com/sirupsen/logrus)  
借用[joonix/log](https://github.com/joonix/log)已经实现的fluentd format; (其实并没有Fluent format, 是针对GKE要求的几个标准fields)

所以Log process pipeline: 

*  App向STDERR输出JSON格式的log, 1 line per log, 包含3个标准fields ^^
*  Docker logging driver(json-file将其写入node的`/var/log/containers`目录，将JSON log以json string的方式封装进入另一个JSON: 每行还是对应一个log message
*  Fluentd 采用`in_tail` plugin, 不停读取日志文件，并从json string中提取原始消息，并根据config: 
    +  转换field name
    +  新增docker/k8s metadata: 如pod id, source, tag等
*  Fluetd最后采用`out_gelf`插件，将Fluentd event转为GELF格式，发送给Graylog.

## App Log library

前面提到，是对[logrus](https://github.com/sirupsen/logrus)的简单封装  

通过Demo看基本功能: 

初始化，并打印log: 

```go
lg := log.NewFluentd(false)
lg.Info("hello log")
```

```json
{"message":"hello log","severity":"info","time":"2021-05-23T07:03:21+10:00"}
```

`log.NewFluentd` create logger, JSON format, 带有3个标准fields: `time`, `severity`, `message`; 通过`FluentdFormatter`实现

通过`logger.WithField`, 附加更多custom fields, 实现`structure logging`


通过`Context`传递fields实现scoped log: 即多条logging message属于同一scope, 如都有相同的`request_id`. 

```golang
// 生成logger object
flog := lg.WithField("request_id", "fake-request-id")
ctx := context.Background()
// embed logger至context
logCtx := log.NewContext(ctx, flog)
// 从context中提取logger
logFromCtx, _ := log.FromContext(logCtx)
logFromCtx.Info("hello")
```

输出带有`request_id`的log: 

```json
{"message":"hello","request_id":"fake-request-id","severity":"info","time":"2021-05-23T07:03:21+10:00"}
```

## Other loggers

`Logrus`目前处于archived mode,即没有新feature,仅限urgent bug fix; 其他选择: 

*  [zerolog](https://github.com/rs/zerolog)  
*  [zap](https://github.com/uber-go/zap)  

最好在自己application定义logging interface, 解耦app logging及logging implementation.

A workable Fluentd config in kubernetes: [fluentd config](https://github.com/eliteGoblin/code_4_blog/blob/master/fluentd_graylog/fluentd_nm_config/logging/docker/fluentd/data/fluent.conf)
  
  
