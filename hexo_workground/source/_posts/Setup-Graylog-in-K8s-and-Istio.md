---
layout: post
title: Setup Graylog in K8s and Istio
subtitle: This is subtitle
date: 2021-05-18 21:04:48
tags: [logging, Fluentd, k8s, Istio]
---

## Background

[系列文章]()第2篇. 本文讲述如何在K8s下搭建Graylog, 作为demo环境; 并配置Istio ingress使其internet accessible.

主要内容:  
*  Graylog基础
*  GKE创建K8s+Istion
*  Deploy Graylog, 及dependency: mongo和elasticsearch
*  运行demo app, 确认日志到达Graylog

<!-- more -->

## Introduction to Graylog

[Graylog](https://www.graylog.org/), 开源的日志一体化解决方案，其实是 Graylog stack, 包含: 

*  Graylog本身: UI, config management
*  MongoDB: 存储Graylog metadata(非log data)
*  ElasticSearch: 存储log, searching engine

### Graylog's Input

Very flexible, 支持常见的format: 

*  Ingest syslog
*  Ingest journald
*  Ingest Raw/Plaintext
*  Ingest GELF: TCP, UDP, HTTP, Kafka
*  Ingest from files
*  Ingest JSON path from HTTP API
*  AWS logs


Graylog可以config多个input; 每个Input独立accept message; 

如GELF TCP, 需config `bind addr`, `port`, `tls`等; Config完成后, Graylog会listen此addr+port. 

<img src="https://raw.githubusercontent.com/eliteGoblin/images/master/blog/img/picgo/20210519080009.png" alt="20210519080009" style="width:500px"/>  

### Graylog Streams

Graylog存储时支持logging message分离, 分为不同的stream(每个stream ElasticSearch有独立的index); 

即分治法, 将logging message按自定义规则group为不同的逻辑分组，如 HTTP500, HTTP200;  搜索时只搜索HTTP500, 没必要每次全局搜索. 

Note:  
*  Graylog可以config多个stream, stream间独立
*  每个incoming message都会根据routing rule, route到特定的stream
*  一个Message可以被route到多个stream

e.g following message: 

```
message: INSERT failed (out of disk space)
level: 3 (error)
source: database-host-1

message: Added user 'foo'.
level: 6 (informational)
source: database-host-2

message: smtp ERR: remote closed the connection
level: 3 (error)
source: application-x
```

只想看DB error, create a stream, rule为: 

*  Field level must be greater than `4`
*  Field source must match regular expression `^database-host-\d+`

Graylog stream matching其实是为message add field: array type `streams`, 存储stream的ID. 后续ElasticSearch可以根据此field建立索引. 

<img src="https://raw.githubusercontent.com/eliteGoblin/images/master/blog/img/picgo/20210519081504.png" alt="20210519081504" style="width:700px"/>  

见[Graylog doc: Streams](https://docs.graylog.org/en/4.0/pages/streams.html)

### Demo version

*  Kubernetes: 1.19
*  Graylog: 3.0
*  MongoDB: 3
*  ElasticSearch: 6.7.2

All demo code can be downloaded in [github](https://github.com/eliteGoblin/code_4_blog/tree/master/fluentd_graylog/fluentd_2)

### Graylog in Docker

[官方文档](https://docs.graylog.org/en/4.0/pages/installation/docker.html)提供了Graylog的docker-compose file: 

```yml
version: '3'
services:
  # MongoDB: https://hub.docker.com/_/mongo/
  mongo:
    image: mongo:4.2
    networks:
      - graylog
  # Elasticsearch: https://www.elastic.co/guide/en/elasticsearch/reference/7.10/docker.html
  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch-oss:7.10.2
    environment:
      - http.host=0.0.0.0
      - transport.host=localhost
      - network.host=0.0.0.0
      - "ES_JAVA_OPTS=-Xms512m -Xmx512m"
    ulimits:
      memlock:
        soft: -1
        hard: -1
    deploy:
      resources:
        limits:
          memory: 1g
    networks:
      - graylog
  # Graylog: https://hub.docker.com/r/graylog/graylog/
  graylog:
    image: graylog/graylog:4.0
    environment:
      # CHANGE ME (must be at least 16 characters)!
      - GRAYLOG_PASSWORD_SECRET=somepasswordpepper
      # Password: admin
      - GRAYLOG_ROOT_PASSWORD_SHA2=8c6976e5b5410415bde908bd4dee15dfb167a9c873fc4bb8a81f6f2ab448a918
      - GRAYLOG_HTTP_EXTERNAL_URI=http://127.0.0.1:9000/
    entrypoint: /usr/bin/tini -- wait-for-it elasticsearch:9200 --  /docker-entrypoint.sh
    networks:
      - graylog
    restart: always
    depends_on:
      - mongo
      - elasticsearch
    ports:
      # Graylog web interface and REST API
      - 9000:9000
      # Syslog TCP
      - 1514:1514
      # Syslog UDP
      - 1514:1514/udp
      # GELF TCP
      - 12201:12201
      # GELF UDP
      - 12201:12201/udp
networks:
  graylog:
    driver: bridge
```

Note: 

*  本地放访问URL需要与`GRAYLOG_HTTP_EXTERNAL_URI`设置的一致, `localhost`无法访问`127.0.0.1`
*  初始ID/Pass: admin/admin

## Create Kubernetes on GKE

```sh
# create cluster, v1.17
gcloud container clusters create ${CLUSTER_NAME} --cluster-version=1.17 --num-nodes=3
gcloud container clusters get-credentials ${CLUSTER_NAME}
# install istio 1.7
istioctl install --set profile=demo
kubectl create namespace istioinaction
kubectl config set-context $(kubectl config current-context) --namespace=istioinaction
# rename current context
kubectl ctx istio_test=.
# get ingress's IP addr
export URL=$(kubectl -n istio-system get svc istio-ingressgateway -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
```

Note: 

*  需要安装gcloud, Istioctl 1.7, kubectl的两个插件: ctx, ns
*  Istio默认安装好后会spinup LB, 作为cluster的ingress point; 我们需要记录下来之后config graylog
*  此脚本将当前context设置为新建立的cluster, 并命名为

## Deploy Graylog in K8s

首先create demo namespace: `k create ns graylog-demo`

Deploy Graylog dependent component, 代码在[这里](https://github.com/eliteGoblin/code_4_blog/tree/master/fluentd_graylog/fluentd_2)下载: 


```
k apply -f mongo-deploy.yaml
k apply -f es-deploy.yaml
```

修改`graylog-deploy.yaml`, 设置`GRAYLOG_HTTP_EXTERNAL_URI`为ingress-controller的IP addr, 即上一步得到的`$URL`, e.g: 

```yaml
- name: GRAYLOG_HTTP_EXTERNAL_URI
          value: http://34.116.94.91/
```

Deploy graylog: `k apply -f graylog-deploy.yaml`


确认deploy正常: `k get pods -w`

```sh
NAME                              READY   STATUS    RESTARTS   AGE
es-deploy-86c7dfcb7b-7684m        1/1     Running   0          7m17s
es-deploy-86c7dfcb7b-g4ks5        1/1     Running   0          7m17s
graylog-deploy-6866cc494d-sc4tm   1/1     Running   0          4s
mongo-deploy-5864d85d5b-cx7jt     1/1     Running   0          7m42s
```

## Config Istio Ingress

Istio默认安装提供了ingress controller, 我们需要配置route, 将HTTP request引入到Graylog中: 

```yaml
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: myk-ingress-gateway
  namespace: graylog-demo
spec:
  selector:
    istio: ingressgateway # use istio default controller
  servers:
  - port:
      number: 80
      name: http
      protocol: HTTP
    hosts:
    - "*"
---
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: graylog-virtualservice
  namespace: graylog-demo
spec:
  hosts:
  - "*"
  gateways:
  - myk-ingress-gateway
  http:
  - route:
    - destination:
        host: graylog3
        port:
          number: 80
```

*  Gateway用来设置ingress-controller的listen port: listen `80`端口
*  VirtualService设置listen port的routing rule: route到K8s service `graylog3` (k8s service 对应内部DNS name)

回顾`graylog3`定义: 

```yml
apiVersion: v1
kind: Service
metadata:
  name: graylog3
  namespace: graylog-demo
spec:
  selector:
    service: graylog-deploy
  ports:
  - name: http-dashboard
    port: 80
    targetPort: 9000
  - name: tcp-input
    port: 12201
    targetPort: 12201
```

Port `80` route到Pod 9000端口，即Graylog的Web UI.

我们访问`http://34.116.94.91`即可看到Graylog dashboard, ID/Pass: `admin/admin`

<img src="https://raw.githubusercontent.com/eliteGoblin/images/master/blog/img/picgo/20210519094139.png" alt="20210519094139" style="width:500px"/>

到此，我们Graylog demo deploy完成; 工作环境需要给Graylog stack附加持久化存储，即PVC. 

## Structured log and GELF

一般建议log输出为JSON, 并采用structure log: 即将关键信息分离到各个fields, 而不是混合输出为同一个fields

**不可取**:

```go
log.Errorf("requestID %s failed with HTTP code %d", requestID, httpCode)
```

应该: 
```
log.WithField("requestID", %s).
    WithField("HTTP code", %d).
    Error("request failed")
```

而GELF是Graylog建议的log fields"约定", 遵循`structured log`, 标准化一些fields, 用来替代之前流行的syslog标准: 

*  GELF message是JSON string
*  GELF内置data types, log需要遵循data type约定，否则Graylog parse时会报错
*  Mandatory fields:
  +  version: type `string (UTF-8)`, GELF spec version, e.g `1.1`
  +  host: type `string (UTF-8)`, name of the host, source or application
  +  short_message: `string (UTF-8)`, short descriptive message
*  Optional GELF fields
  +  full_message: 
  +  timestamp:
  +  level: type `number`, standard syslog levels, DEFAULT 1 (ALERT)
  +  _[additional field] : Other custom fields
    +  type `string` or `number`, 程序自定义的fields
    +  Log library需要给fields附加prefix`_`

Example GELF message payload: 

```json
{
  "version": "1.1",
  "host": "example.org",
  "short_message": "A short message that helps you identify what is going on",
  "full_message": "Backtrace here\n\nmore stuff",
  "timestamp": 1385053862.3072,
  "level": 1,
  "_user_id": 9001,
  "_some_info": "foo",
  "_some_env_var": "bar"
}
```

Syslog的`severity` level: 

<img src="https://raw.githubusercontent.com/eliteGoblin/images/master/blog/img/picgo/20210520085811.png" alt="20210520085811" style="width:500px"/>  

本文的Demo设置Graylog 接受GELF TCP, 即用TCP协议传输log message, 可以直接测试 via TCP message:

```sh
echo -n -e '{ "version": "1.1", "host": "example.org", "short_message": "A short message", "level": 5, "_some_info": "foo" }'"\0" | nc -w0 graylog.example.com 12201
```

见 [官方doc: GELF](https://docs.graylog.org/en/3.3/pages/gelf.html)  

## Send logs in k8s cluster to graylog

### Config Graylog Input

Setup 两个Input: HTTP 和 TCP

GELF HTTP, listen 12201

<img src="https://raw.githubusercontent.com/eliteGoblin/images/master/blog/img/picgo/20210520083416.png" alt="20210520083416" style="width:500px"/>  

### Send log directly via cronjob

我们设置cronjob来通过HTTP向Graylog发送日志, using GELF format

```yaml
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: curl-cron-job
spec:
  schedule: "* * * * *"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: curl-job
            image: alpine:3.9.4
            args:
            - /bin/sh
            - -c
            - apk add curl -y; while true; do curl -XPOST http://graylog3:12201/gelf -p0 -d '{"short_message":"Hello there, i am your corn job ;)", "host":"alpine-k8s.org", "facility":"test", "_foo":"bar"}';sleep 1s; done
          restartPolicy: OnFailure
```

Create cronjob: `k apply -f log_generate_cronjob.yaml`

循环向Graylog发送log: `curl -XPOST http://graylog3:12201/gelf`, inside cluster

```json
{
  "short_message": "Hello there, i am your corn job ;)",
  "host": "alpine-k8s.org",
  "facility": "test",
  "_foo": "bar"
}
```

Demo log message包含了GELF mandatory fields: `short_message`和`host`(似乎不加`version`也被接受)

可以从Graylog dashboard看到log message: 

<img src="https://raw.githubusercontent.com/eliteGoblin/images/master/blog/img/picgo/20210520084552.png" alt="20210520084552" style="width:500px"/>  

### Config Fluentd to send log message

直接向Graylog发送log message成功，接下来我们加入Fluentd: 

*  Container产生log到stdout, 被docker logging driver捕获，写入到host的`/var/log/containers`folder
*  Fluentd监视`/var/log`, 发送日志到Graylog, via GELF TCP, port 12201

首先Clean up: 

```sh
k delete -f log_generate_cronjob.yaml
```

删除Graylog `GELF HTTP` input, 并Create `GELF TCP` input

Deploy Fluentd Daemonset: 关键config

```yaml
containers:
- name: fluentd
  image: fluent/fluentd-kubernetes-daemonset:v1-debian-graylog
  imagePullPolicy: IfNotPresent
  env:
    - name:  FLUENT_GRAYLOG_HOST
      value: "graylog3.graylog-demo.svc.cluster.local"
    - name:  FLUENT_GRAYLOG_PORT
      value: "12201"
    - name:  FLUENT_GRAYLOG_PROTOCOL
      value: "tcp"
  volumeMounts:
        - name: varlog
          mountPath: /var/log
        - name: varlibdockercontainers
          mountPath: /var/lib/docker/containers
          readOnly: true
```

令Fluentd向Graylog server发送log message, via TCP port `12201`. 这里采用默认Fluentd配置. 

Create daemonset: `k apply -f fluentd_daemonset.yaml`

观察Graylog, 发现有新的log: 

对比我们发送的logging message: 

```json
{
  "time": "2021-04-30 00:21:24.383 +00:00",
  "message": "frank debug",
  "severity": "info",
  "level": 6
}
```

<img src="https://raw.githubusercontent.com/eliteGoblin/images/master/blog/img/picgo/20210520092645.png" alt="20210520092645" style="width:500px"/>

一些有趣的现象: 

* 多了一些fields, 添加docker, k8s metadata: 
  +  `docker`
  +  `kubernetes`
  +  `source`
  +  `stream`
  +  `tag`: tag added by Fluentd
*  同时发现我们的JSON log message, 作为string存在`message` fields, expected, 因为Docker logging driver作了转换: 从STDOUT获取的每一行，无论format, 都作为string; 

我们需要能从Docker logging driver记录的log `unpack` 我们的fields, 因此需要进一步配置Fluentd. 

同时可以SSH登录Node, 验证logging file内容，确实是被`Docker logging driver`统一处理过: 

```
sudo tail -f /var/log/containers/graylog-deploy-6866cc494d-mbvmx_graylog-demo_graylog3-25783c1d79760b91a6c6d0650524aa6631d02fde25b6c7a5fa63691d79339afe.log
```

```json
{"log":"{\"time\": \"2021-04-30 00:21:24.383 +00:00\", \"message\": \"frank debug\", severity: \"info\", level: 6}\n","stream":"stdout","time":"2021-05-19T23:21:00.604887883Z"}
{"log":"{\"time\": \"2021-04-30 00:21:24.383 +00:00\", \"message\": \"frank debug\", severity: \"info\", level: 6}\n","stream":"stdout","time":"2021-05-19T23:21:05.606245351Z"}
```

## 直接利用GKE的logging solution

GKE默认采用`fluent-bit`作为logging collector, 并提供logging dashboard, 见[Customizing Cloud Logging logs for Google Kubernetes Engine with Fluentd](https://cloud.google.com/architecture/customizing-stackdriver-logs-fluentd)  

## Reference

[Graylog With Kubernetes in GKE](https://dzone.com/articles/graylog-with-kubernetes-in-gke)   