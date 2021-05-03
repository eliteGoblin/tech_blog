---
title: Full-text search with Golang and ElasticSearch
date: 2018-09-12 16:17:39
tags: [golang ElasticSearch]
keywords:
description:
---

## Preface

搜索，对于很多项目来说非常重要，甚至是核心功能。常见的一种搜索形式是: 输入关键字组合，给出全系统中匹配或相似的1到多条信息(通常称之为1个document)。也称此为full-text search 。  

企业FAQ系统，就可以用此思路搭建：用户输入问题，搜索出系统中和用户问题最相似的问题，然后给出已记录的标准答案；相当于计算机"*理解*"用户的疑问。实现时，既可以手工实现：用分词算法切割用户输入的自然语言，得到关键词数组，用bm25算法，根据词频计算用户问题和系统FAQ问题的相似度，排序之后给出推荐；或者采用开源解决方案，ElasticSearch就是其中被广为采纳的一个。  

之前在悉尼找工作时，正好一家media公司给出的编程题便是用Golang搭建文本搜索系统，当时研究了一下ElasticSearch，给出了解决方案；于是就把这个有趣的经历记录下来。

<div align="center">
{% asset_img ElasticsearchDocker.png %}
</div>

<!-- more -->

## 目标

*  实现article的简单"*增改查*"系统，article的json格式: 
    ```golang
    type Article struct {
        Id string       `json:"id"`
        Title string    `json:"title"`
        Date string     `json:"date"`
        Body string     `json:"body"`
        Tags []string   `json:"tags"`
        CreatedTime string `json:"created_time"`
    }
    ```
*  按ID查找article: 
    ```
    GET /articles/{article_id}
    ```
*  Upsert(Update+Insert)操作
    ```
    POST /articles
    ```
*  指定日期，tag查询articles
    ```golang
    GET /tags/{tag_name}/{date}
    // 返回
    type SearchByTagResult struct {
        Tag string              `json:"tag"`
        Count int               `json:"count"`
        Articles []string       `json:"articles"`
        RelatedTags []string    `json:"related_tags"`
    }
    ```
    +  Tag: 查询的Tag
    +  Count: 命中的Article数目
    +  Articles: 命中的article id数组
    +  RelatedTags: 命中的所有article的所有tag的集合(不包括查询的tag)

## 题目思考

article的增改实现很trival，真正需要考虑的是查询如何实现；

不考虑第三方组件:  
*  需要按照tag, 和date两个维度过滤article
*  用map等实现自己的索引，得到需要的结果集

得出的解决方案，扩展性差：如果想增加过滤条件，如create_time等，需要较大的索引结构改动，于是想到用ElasticSearch来为article建立索引，用其query功能得到article结果集合，然后提取题目要求的结果信息。

完整的项目代码在这里: [article_api](https://github.com/eliteGoblin/code_4_blog/tree/master/article_api)


## ElasticSearch方案实现

### 环境搭建

最portable的方式就是将其放入容器，这也是题目要求的，在reviewer的机器上，以最少的步骤跑起来。ES和golang app各占一个container, 用docker-compose(简单容器编排工具)，实现container的启动及网络互连:  

docker-compose.yml
```yml
version: '3.5'
services:
  search_api:
    container_name: 'article_api'
    build: './fairfax'
    restart: 'on-failure'
    ports:
      - '8080:8080'
    depends_on:
      - elasticsearch
  elasticsearch:
    container_name: 'elasticsearch'
    image: 'docker.elastic.co/elasticsearch/elasticsearch:6.2.3'
    ports:
      - '9200:9200'
```

说明:  
*  完整运行环境由service组成：search_api为golang app, 另一个是ElasticsSearch
*  golang app没有预先build好的image，build指定包含其Dockerfile的目录
*  两个container处于同一网段，通过<service_name>:port方式访问别的service

### ElasticSearch Client Library

首先考虑如何在Golang中访问ES，采用[olivere/elastic](https://github.com/olivere/elastic)作为client library, star 2900+。

连接ElasticSearch:  

```golang
import "github.com/olivere/elastic"
client, err := elastic.NewClient(
                elastic.SetURL("http://elasticsearch:9200"),
                elastic.SetSniff(false),
            )
```

### 建立空索引

```golang
_, err = client.CreateIndex("article").Do(context.Background())
```

### Upsert Document

```golang
// Article为json格式，代表document，见代码定义
func (selfPtr *ESDb)AddDocToDb(article *core.Article) error {
    _, err := selfPtr.esClient.Index().
        Index("article").
        Type("doc").
        Id(article.Id).
        BodyJson(article).
        Refresh("wait_for").
        Do(context.Background())
    return err
}
```

### 查找Document

```golang
func (selfPtr *ESDb)SearchByTag(tag string, date string) (res *SearchByTagResult, err error){
    query := elastic.NewBoolQuery()
    q1 := elastic.NewTermsQuery("tags", tag)
    q2 := elastic.NewMatchQuery("date", date)
    query.Must(q1, q2)

    searchResult, err := selfPtr.esClient.Search().
        Index("article").                   // search in index "article"
        Query(query).                               // specify the query
        Pretty(true).                       // pretty print request and response JSON
        Sort("created_time.keyword", false).    // sort by "id" field, ascending
        Do(context.Background())                    // execute
    ...
}
```

*  用多个子query合并成一个bool query的方式实现"与"的操作
*  Sort排序

### 提取query result

```golang
// core.Article是自己定义的document的golang struct
var esItem core.Article
for _, item := range searchResult.Each(reflect.TypeOf(esItem)) {
    if t, ok := item.(core.Article); ok {
        return &t, nil // 返回匹配到的第一条Document
    }
}
```

## 结论

本文搭建了Golang+ElasticSearch环境，并实现了简单的增改查功能(主要代码来自于[reference2][How to Build a Search Service with Go and Elasticsearch])。ElasticSearch由于其强大的搜索及集群方案已经普遍用在全文搜索，FAQ，及日志分析的ELK Stack等解决方案中，值得进一步深入学习。

## Reference

[Full-text search](https://en.wikipedia.org/wiki/Full-text_search)  
[How to Build a Search Service with Go and Elasticsearch](https://outcrawl.com/go-elastic-search-service/)  
[TF-IDF与余弦相似性的应用（一）：自动提取关键词](http://www.ruanyifeng.com/blog/2013/03/tf-idf.html)  
[TF-IDF与余弦相似性的应用（二）：找出相似文章](http://www.ruanyifeng.com/blog/2013/03/cosine_similarity.html)  