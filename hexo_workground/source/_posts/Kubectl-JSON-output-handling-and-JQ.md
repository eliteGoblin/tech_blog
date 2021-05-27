---
layout: post
title: Kubectl JSON output handling and JQ
date: 2021-05-25 22:17:10
tags: [Kubernetes JQ]
keywords: [Kubernetes kubectl JSON JQ]
---

<img src="https://raw.githubusercontent.com/eliteGoblin/images/master/blog/img/picgo/20210527140611.png" alt="20210527140611" style="width:500px"/>

## WHY

Kubernetes可以看作"Object Store", Object格式: `yaml` or `json`, 两者等价. 

Kubectl的`get`命令读取一系列Objects, 并以JSON格式返回. 这些JSON格式往往很大，有复杂的嵌套，常见的需求: 

*  获取Output的subset, 往往deep nested: 某些fields， 而不是全部JSON
*  Filter out irrelevant objects by fields: e.g.我们只想返回失败的job
*  Sort: sort by creation date, status, etc

相比每次将JSON当做文本，`grep`; 我们需要更贴合JSON的方式. 

本质上JSON output是我们的数据源, 我们需要query language, 就像query relational database一样. 

本文介绍用`JQ`, `Jid`处理`kubectl`输出, 应对常见的情景. 

<!-- more -->

示例代码在[github](https://github.com/eliteGoblin/code_4_blog/kubectl_jq)下载

## Introduction to JQ

`jq`首先是binary tool, 用来query json; 由于强大的query功能，可视为独立的language, 类似`awk`.

先看`jq`提供的基本功能, Demo JSON

```json
{ "store": {
    "book": [ 
      { "category": "reference",
        "author": "Nigel Rees",
        "title": "Sayings of the Century",
        "price": 8.95
      },
      { "category": "fiction",
        "author": "Evelyn Waugh",
        "title": "Sword of Honour",
        "price": 12.99
      },
      { "category": "fiction",
        "author": "Herman Melville",
        "title": "Moby Dick",
        "isbn": "0-553-21311-3",
        "price": 8.99
      },
      { "category": "fiction",
        "author": "J. R. R. Tolkien",
        "title": "The Lord of the Rings",
        "isbn": "0-395-19395-8",
        "price": 22.99
      }
    ],
    "bicycle": {
      "color": "red",
      "price": 19.95
    }
  }
}
```


```sh
# prettify JSON
cat example.json | jq .
# extract fields inside array element
cat example.json | jq '.store.book[1].price'
# extract field of array, 输出multiple elements
cat example.json | jq '.store.book[].price'
# array转化为iterator,即等价于用element多次调用, 不同于输出一个JSON array
cat example.json | jq '.store.book | .[]'
# 构建新的object
cat example.json | jq '.store.book[1] | {bookTitle: .title, bookPrice: .price}'
# 构建array, 将多个elements, construct为一个array
cat example.json | jq '[.store.book | .[]]'
```

[JQ tutorial](https://stedolan.github.io/jq/tutorial/)也提供了一些例子, 值得一读。

## 进一步了解 JQ 

以下主要来自[JQ Manual](https://stedolan.github.io/jq/manual/)  

以一个具体`jq` cmd为例: query k8s jobs, 以创建时间排序, filter, 只留下failed jobs, 并reconstruct object.
```sh
k get jobs -o json --sort-by=.metadata.creationTimestamp | jq -r '.items | .[] | select(.metadata.name | contains("openworld")) | select(.status.active != 1) | select(.status.conditions[0].type == "Failed") | {name: .metadata.name, time: .metadata.creationTimestamp, status: .status}' 
```

可见`jq`的query由一系列`|`组成，类似UNIX的pipe, 连接独立component的input, output. 

每个component, 称为filter, 接收输入，完成输出. 每个filter包含不同的`jq` built-in operator, 实行不同的转换. 

有些filter会产生multiple output(filter A), 当 pipe `filter A`的输出到`filter B`时，会`runs the second filter for each element of the array`; 即不用explicitly写for-loop, 直接pipe两者即可. 

**Every filter has input and output**, 即使常数"123"也是filter: 不管输入是什么，输出"123". JQ里一切皆filter.

### Filters

```sh
.                   # 最简单的filter, echo input
.foo .foo.bar       # 等价于 .foo | .bar, object index
                    #   严格形式为.["foo"], 如果key有特殊字符
.[2]                # array index
.[2:5]              # array slice
.[]                 # array iterator: 将single array转为multiple elements
filterA, filterB    # fan out to 多个filter, combine它们的输出, 
                    #   如echo "{}" | jq '1, 2, 3'输出3个digit, 3行
|                   # 连接两filter, 特殊的filter
(expr)              # group operators, 视expr为expression
```

### Types && Values

> jq supports the same set of datatypes as JSON - numbers, strings, booleans, arrays, objects (which in JSON-speak are hashes with only string keys), and "null".

针对Type转换的主要有`[]`: array construction及`{}`: object construction.

`[]`将multiple input转化为JSON array, 即single object, 如之前提到的

```sh
cat example.json | jq '[.store.book | .[]]'
```

`.[]`是将array打散为multiple output, 因此`[]`和`.[]`为逆过程. 

`{}`reconstruct object: 原始object不能满足需要，我们根据输出，重新构造object, 之前例子: 

```sh
| {name: .metadata.name, time: .metadata.creationTimestamp, status: .status}
```

有一个common case: 我们只想获得input object的some fields as new Object, 可以: 

```
{title: .title | author: .author}
```

可以简写为: `{title, author}`

当我们想递归遍下降历所有的sub objects of input "Big" object: 用: `..`: 

```sh
cat example.json | jq '.. | .price?'
```

此例中，bicycle及book都有price, 我们都列出来; `.bar?`的`?`在没有`bar`index时不至于报错. 

`jq`的基本元素介绍完了, 但实际使用，除了基本原理，了解built-in非常关键，毕竟advance场景, 需要自己实现operator时候不多。

## Built-in operators and functions

列举一些常用的operator和function, 感觉在做functional programming.

### select(boolean_expression)

重要的filter: 一般用来过滤array的elements: 

如果input满足`boolean_expression`, 原样保留; 否则丢弃

boolean_expression可以调用input, 结合其他filter: 

```sh
jq `[1,2,3] | map(select(. >= 2))`
Output: [2, 3]
```

another example: 

```sh
jq '.[] | select(.id == "second")'
Input	[{"id": "first", "val": 1}, {"id": "second", "val": 2}]
Output	{"id": "second", "val": 2}
```

### 四则运算

`+ -`

初看起来很简单, 但功能丰富: 

`a + b`连接两个filter, 将输入传递给both `a`和`b`; 再将两者输出做"add":

*  Numbers: 直接add
*  Arrays:  concat 2 array into 1 bigger array
*  Strings: concat, similar as array
*  Objects: merge 2 objects, 右边覆盖左边

同理有`-`, 但仅支持两边都是`numbers`或`array`

`* / %`

一般仅对两个`numbers`, 也支持`array`, `string`和`objects`, 但效果比较诡异，应该不常用，略过。

### keys

`keys, keys_unsorted`

要求输入为`object`, 输出`keys`为array.

`has(key)`: 
  +  input: object; output: boolean, 是否含有此key.
  +  input: array: array是否有此element

```sh
jq 'map(has("foo"))'
Input	[{"foo": 42}, {}]
Output	[true, false]
```

`in`: if input key in given object, or array; inversed version of `has`

```sh
jq 'map(in([0,1]))'
Input	[2, 0]
Output	[false, true]
```

### map(expr) map_values(expr)

*  input: array; 
*  `expr`: expr or function, 对每个elements或fields应用expr; 
*  output: 应用`expr`后的新array

```sh
jq 'map_values(.+1)'
Input	{"a": 1, "b": 2, "c": 3}
Output	{"a": 2, "b": 3, "c": 4}
```

可以传入其他built-in, 例如获取field type.
```sh
jq 'map(type)'
Input	[0, false, [], {}, null, "hello"]
Output	["number", "boolean", "array", "object", "null", "string"]
```

### filter by data type

arrays, objects, iterables, booleans, numbers, normals, finites, strings, nulls, values, scalars

> These built-ins select only inputs that are arrays, objects, iterables (arrays or objects), booleans, numbers, normal numbers, finite numbers, strings, null, non-null values, and non-iterables, respectively.

只留下match的type
```sh
jq '[].[]|numbers, nulls]'
Input	[[],{},1,"foo",null,true,false]
Output	[1 null]
```

### any all

input: array of boolean values
output: `any`任意一个true为true; `all`需要所有为true.

```sh
jq 'any'
Input	[true, false]
Output	true
```

### array operation

*  `sort`, `sort_by(path_expression)`
```
jq 'sort_by(.foo)'
Input	[{"foo":4, "bar":10}, {"foo":3, "bar":100}, {"foo":2, "bar":1}]
Output	[{"foo":2, "bar":1}, {"foo":3, "bar":100}, {"foo":4, "bar":10}]
```
*  `unique`, `unique_by`
*  `min`, `max`, `min_by`, `max_by`
*  `reverse`

### string operation

一系列的操作: 

*  `contains(s)`
*  `index(s)`, `rindex(s)`
*  `inside`
*  `startWith(str)`, `endWith(str)`
*  `split`, `join`

String interpolation - `\(foo)`

```sh
jq '"The input was \(.), which is one less than \(.+1)"'
Input	42
Output	"The input was 42, which is one less than 43"
```

### length

"length" of values: 支持`string`, `array`, `object`


## 逻辑判断 operator

作用一目了然, 返回true or false: 

*  `==, !=`
```sh
jq '.[] == 1'
Input	[1, 1.0, "1", "banana"]
Output	true
        true
        false
        false
```
*  `>, >=, <=, <`
*  `and` `or` `not`
```sh
jq '[true, false | not]'
Input	null
Output	[false, true]
```

Alternative operator: `a // b`: 若`a`不是`false`或`null`, output`a`; 否则输出`b`; 等于是给定"默认值"

```sh
jq '.foo // 42'
Input	{}
Output	42
```
这些可在[jq play](https://jqplay.org/)上直接看结果.

在local可以直接`jq -n 'expr'`来测试query

## Regular expressions (PCRE)

`jq`支持完整的`re`search: 采用与`php, ruby, sublime`等相同的`re` library: `Oniguruma`

RE操作: 

```
STRING | FILTER( REGEX )
STRING | FILTER( REGEX; FLAGS )
STRING | FILTER( [REGEX] )
STRING | FILTER( [REGEX, FLAGS] )
```

*  STRING为待match string, 作为input
*  FILTER is one of: 
  +  match: 找到match, 输出object
  +  test: Like `match`, but does not return match objects, only `true` or `false`
  +  capture: 保留capture name, 结果存到新object

FLAGS is a string consisting of one of more of the supported flags:

```
g - Global search (find all matches, not just the first)
i - Case insensitive search
m - Multi line mode ('.' will match newlines)
n - Ignore empty matches
p - Both s and m modes are enabled
s - Single line mode ('^' -> '\A', '$' -> '\Z')
l - Find longest possible matches
x - Extended regex format (ignore whitespace and comments)
```

测试一系列的match: 
```sh
jq '.[] | test("a b c # spaces are ignored"; "ix")'
Input	["xabcd", "ABC"]
Output	true
true
```

可以结合`select`来正则匹配: `select(.metadata.name | test("test-"))`

## JQ Output format

特殊的`filter`, 一般在最终输出前做必要的formatting和encoding/decoding, 语法为`@foo`, 有如下几种: 

*  `@text`: 实际调用`tostring`
*  `@json`: serialize to json
*  `@html`: HTML escaping: `<>&'"` map to `&lt;`, `&gt;`, `&amp;`, `&apos;`, `&quot;`.
*  `@uri`: percent-encoding
*  `@csv`, `@tsv`: input必须为`array`, 转为`csv`, `tsv`格式
*  `@sh`: escape for POSIX shell
*  `@base64`: base64 encode; `@base64d`: base64 decode

## JQ, JSON path, Kubectl

[JSON Path standard](https://tools.ietf.org/id/draft-goessner-dispatch-jsonpath-00.html) 提供了不同于`jq`的query syntax, 两者的比较见[JQ doc: For JSONPath users](https://github.com/stedolan/jq/wiki/For-JSONPath-users)  

Kubectl本身支持`JSONPath`, 见[官方doc](https://kubernetes.io/zh/docs/reference/kubectl/jsonpath/), 但正如doc所说:  

> 不支持 JSONPath 正则表达式。如需使用正则表达式进行匹配操作，您可以使用如 jq 之类的工具

```sh
kubectl get pods -o json | jq -r '.items[] | select(.metadata.name | test("test-")).spec.containers[].image'
```

## 用JID visualize field selection

`Jid`使用场景: 生成JSON query statement: 一步步的从一个大json中选取deep nested field, 并能autocomplete.

安装`jid` tool:

```sh
go get -u github.com/simeji/jid/cmd/jid
```

使用`jid`非常简单:
```sh
cat example.json | jid -q | xclip -selection c
```

<img src="https://raw.githubusercontent.com/eliteGoblin/images/master/blog/img/picgo/20210526162922.png" alt="20210526162922" style="width:500px"/>

依次指定field来filter, `tab`自动补全, 最后pipe到clipboard, `ctrl+v`或`xclip -selection c -o`即可得到生成的query: `.store.book[1].price`

## 结论

我们的主要目标是用query `kubectl`的`json` output, 一些建议: 

*  kubectl使用`--sort-by`指定排序规则(只能为`integer`或`string`)
*  如果JSON结果大，复杂，先用`jid`获取`jq`能直接使用的query path
*  将`kubectl`的输出pipe到`jq`中，filter, 转换，得到想要的结果

## Ref

[JQ manual](https://stedolan.github.io/jq/manual/)  
[Kubectl output options](https://gist.github.com/so0k/42313dbb3b547a0f51a547bb968696ba)  

