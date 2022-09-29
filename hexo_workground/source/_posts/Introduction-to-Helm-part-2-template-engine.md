---
layout: post
title: 'Introduction to Helm -- part 2: template engine'
date: 2022-04-09 11:55:39
tags: [Helm, k8s]
---


# [Background](#background)

上篇介绍了why Helm 和package management feature, 本篇讨论如何使用Helm template engine.

<!-- more -->

# [Debug](#debug)

local try and error: 即时generate当前template:  

```s
# cd mychart
helm template . --debug 2>/dev/null
```

Go tempalte提供`printf`打印object, 如打印`Values` object.

```yaml
{{ printf "%#v" $.Values }}
{{ printf "%t" .Values.favorite }}
{{ list 1 2 3 | toStrings }}
{{ dict 1 2.66 | toJson }}
{{ index (list 1 2 3) 2}}
```

# [Go template](#go-template)

##  Objects

Templates are executed by applying them to a data structure. 

Text Template + Go Object = 结果

The current object is represented as `.`, e.g 若object是`string`: 

则render:  `{{ . }}`

若当前object有`Name` field, 为string, 则: `{{ .Name }}`.  


关于Object传递到template: 

```go
type Student struct {
    Name string
}
s := Student{"Satish"}
tmpl, __ := tmpl.Parse("Hello {{.Name}}!")
tmpl.Execute(os.Stdout, s) // 在上面的template, object s 就是 "dot"
```

## Cursor 和 Context

When execution begins, `$` is set to the data argument passed to Execute; "Dot"代表当前cursor. 


> Execution of the template walks the structure and sets the cursor, represented by a period '.' and called "dot", to the value at the current location in the structure as execution proceeds.

Cursor即当前object/context,会变化: 如range, with 会设置不同的object, 到其内层template中. 

`range`, `with` 包裹的template: 无法access外层Object,需要定义变量: 即这两者会让scope缩小(由输入的Object `$` 变为某个sub Object)


## Variables

```go
type Person struct {
        Name   string
        Emails []string
}
const tmpl = `{{$name := .Name}}
{{range .Emails}}
    Name is {{$name}}, email is {{.}}
{{end}}
`
func main() {
        person := Person{
                Name:   "Satish",
                Emails: []string{
                    "satish@rubylearning.org", 
                    "satishtalim@gmail.com"
                },
        }

        t := template.New("Person template")
        t, _ := t.Parse(tmpl)
        err = t.Execute(os.Stdout, person)
}
```

传入的Object为`email` string, 无法access外层Object, 如Person's name: 需要定义变量.

`with`也类似: 以下输出都是"output"

```yaml
{{with "output"}}{{printf "%q" .}}{{end}}
	# A with action using dot.

{{with $x := "output" | printf "%q"}}{{$x}}{{end}}
	# A with action that creates and uses a variable.

{{with $x := "output"}}{{printf "%q" $x}}{{end}}
	# A with action that uses a variable in another action.

{{with $x := "output"}}{{$x | printf "%q"}}{{end}}
	# The same, but pipelined.
```

## Named template

Tempalte内还可定义named template: 
```yaml
{{define "T1"}}ONE{{end}}
{{define "T2"}}TWO{{end}}
{{define "T3"}}{{template "T1"}} {{template "T2"}}{{end}}
{{template "T3"}}
# output: ONE TWO
```

## Pipeline

类似Linux的pipeline, 将多个component chain在一起: 

> A pipeline may be "chained" by separating a sequence of commands with pipeline characters '|'. In a chained pipeline, the result of each command is passed as the last argument of the following command. The output of the final command in the pipeline is the value of the pipeline.

A command is: 
*  a simple value (argument) or 
*  a function or 
*  method call

当前cmd的last input是上一个cmd的output

```yaml
# A function call. The printf parameters are identical to fmt.Printf from Go.
{{printf "%q" "output"}}
# A function call whose final argument comes from the previous command.
{{"output" | printf "%q"}}
# A parenthesized argument.
{{printf "%q" (print "out" "put")}}
# A more elaborate call.
{{"put" | printf "%s%s" "out" | printf "%q"}}
# A longer chain.
{{"output" | printf "%s" | printf "%q"}}
```

## Flow Control

两类: 

control structures: 

*  if/else
*  range
*  with

同时还有`action`: 

*  define: declares a new named template inside of your template
*  template: imports a named template

注: `include`是function, 可用在pipeline中，因此一般用`include`而不是`template`

The `include` function allows you to bring in another template, and then pass the results to other template functions.

> “Because template is an action, and not a function, there is no way to pass the output of a template call to other functions; the data is simply inserted inline.”  

### `if/else`: 

```yaml
{{ if PIPELINE }}
  # Do something
{{ else if OTHER PIPELINE }}
  # Do something else
{{ else }}
  # Default case
{{ end }}
```

关于`false` condition: 
*  boolean false
*  0
*  nil
*  empty string, slice, map, etc

### `with` 

用来modify scope的:  
```yaml
{{ with PIPELINE }}
  # restricted scope
{{ end }}
```

例子见下面: 
### range

values: 
```yaml
favorite:
  drink: coffee
  food: pizza
pizzaToppings:
  - mushrooms
  - cheese
  - peppers
  - onions
```

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ .Release.Name }}-configmap
data:
  myvalue: "Hello World"
  {{- with .Values.favorite }}
  drink: {{ .drink | default "tea" | quote }}
  food: {{ .food | upper | quote }}
  {{- end }}
  toppings: |-
    {{- range .Values.pizzaToppings }}
    - {{ . | title | quote }}
    {{- end }}    
```

注意`{{ . | title | quote }}`, scope在`range`之后发生了改变: 变为list element.

### Whitespace   

`if/else`和`range`很常见, 但Go template engine会产生些unexpected whitespace(包含space和换行): 

> When the template engine runs, it removes the contents inside of {{ and }}, but it leaves the remaining whitespace exactly as is.

```yaml
food: "cake"
{{ if eq .Values.favorite.drink "coffee" }}
mug: "true"
{{ end }}
```

会生成额外的换行(由`if`而来): 

```yaml
food: "PIZZA"

mug: "true"
```

解决方法是: 

*  `{{- `(with the dash and space added):  whitespace should be chomped left: 消除左whitespace直到遇到非whitespace字符:
*  ` -}}`: 消除右whitespace直到非whitespace字符

效果示意: 

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ .Release.Name }}-configmap
data:
  myvalue: "Hello World"
  drink: {{ .Values.favorite.drink | default "tea" | quote }}
  food: {{ .Values.favorite.food | upper | quote }}*
**{{- if eq .Values.favorite.drink "coffee" }}
  mug: "true"*
**{{- end }}
```

*  `*`代表删除的whitespace!
*  `if` 行没有产生实际内容(仅仅控制内容产生), `{{- `消除其whitespace
*  `end`行也同样没有内容，不消除会多出一个空行. 

e.g: 

```yaml
drink: {{ .Values.favorite.drink | default "tea" | quote }}
food: {{ .Values.favorite.food | upper | quote }}
{{ if eq .Values.favorite.drink "coffee" }}
mug: "true"
{{ end }}
another: "123"
```

values为:

```yaml
favorite:
  drink: coffee
  food: pizza
```

输出: 

```yaml
drink: "coffee"
food: "PIZZA"

mug: "true"

another: "123"
```

> Be careful! Newlines are whitespace!

还可以用 `ident` 来控制空格, 如: 

```yaml
{{ include "api-service-tpl.rollout.env.common" . | indent 12 }}
```

[Helm doc: Flow Control](https://helm.sh/docs/chart_template_guide/control_structures/)  里讲解的很清楚.  


## Data type and Frequent Functions

Variables in templates are [typed](https://helm.sh/docs/chart_template_guide/data_types/)(因为Go是typed):  

*  string
*  bool
*  int
*  float64
*  byte array
*  object/struct
*  list: immutable
*  dict: a string-keyed map (`map[string]interface{}`), where the value is one of the previous types

function建立在data type基石上, 理解一门语言的data type至关重要. 

Helm template function有两部分:  

*  Go Template Function
*  [sprig function for Go Template](http://masterminds.github.io/sprig/)

全部Helm function在[Template Function List](https://helm.sh/docs/chart_template_guide/function_list/)  

一些常见函数:  
```yaml
# if .Bar empty/non-exist, set to "foo"
default "foo" .Bar 
# takes a list of values and returns the first non-empty one
coalesce 0 1 2 # 1
# create a list, output to type conversion function
{{ list 1 2 3 | toStrings }} # 1 2 3
# Create a dict: "name1" "value1" "name2" "value2" "name3" "value 3"
{{ dict 1 2.66 | toJson }} # {"1":2.66}
# To get the nth element of a list, multi-dimensional: index $mylist i j k
{{ index (list 1 2 3) 2}} # 3
# merge $dest $source1 $source2, Merge two or more dictionaries into one, first get precedence
{{ merge (dict "name" "Frank") (dict "name" "Joe") (dict "name" "Jacob") (dict "addr" "bondi") }} # map[addr:bondi name:Frank]
# mustMerge will return an error in case of unsuccessful merge
```

下一篇谈Helm template相关的advanced features.
# Reference

[Go text/template](https://pkg.go.dev/text/template)  
[Helm from basics to advanced — part II](https://banzaicloud.com/blog/creating-helm-charts-part-2/)   
[Values Files](https://helm.sh/docs/chart_template_guide/values_files/)  
