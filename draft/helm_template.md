
# Point

*  基本template, 结构, 如何调试
*  pipeline, 常见function
*  处理space, - 和 ident

*  include unexported template

# advanced

*  subchart如何传值: by example
*  library chart 

*  Helm provider
*  ArgoCD using Helm



# Helm

Before we jump into creating common code, lets do a quick review of some relevant Helm concepts. A named template (sometimes called a partial or a subtemplate) is simply a template defined inside of a file, and given a name. In the templates/ directory, any file that begins with an underscore(_) is not expected to output a Kubernetes manifest file. So by convention, helper templates and partials are placed in a _*.tpl or _*.yaml files.


Play with Helm: 

```
helm template . --debug 2>/dev/null
```

# Spaces

By default, all text between actions is copied verbatim when the template is executed. For example, the string " items are made of " in the example above appears on standard output when the program is run.

However, to aid in formatting template source code, if an action's left delimiter (by default "{{") is followed immediately by a minus sign and white space, all trailing white space is trimmed from the immediately preceding text. Similarly, if the right delimiter ("}}") is preceded by white space and a minus sign, all leading white space is trimmed from the immediately following text. In these trim markers, the white space must be present: "{{- 3}}" is like "{{3}}" but trims the immediately preceding text, while "{{-3}}" parses as an action containing the number -3.

# Function

index
Returns the result of indexing its first argument by the
following arguments. Thus "index x 1 2 3" is, in Go syntax,
x[1][2][3]. Each indexed item must be a map, slice, or array.

```
default DEFAULT_VALUE GIVEN_VALUE
```


[Helm from basics to advanced](https://banzaicloud.com/blog/creating-helm-charts-part-2/) 两篇都不错

With if/else scope: [Flow Control](https://helm.sh/docs/chart_template_guide/control_structures/): In this section, we'll talk about if, with, and range. The others are covered in the "Named Templates" section later in this guide.


Whitespace 消除: 

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

# Scope, cursor

From https://pkg.go.dev/text/template#section-directories
 Execution of the template walks the structure and sets the cursor, represented by a period '.' and called "dot", to the value at the current location in the structure as execution proceeds


> Or, we can use $ for accessing the object Release.Name from the parent scope. $ is mapped to the root scope

*  $ 指 root scope; 如 `.Values`指root的Values object
*  . 指当前scope, 默认是root scope, 可通过with, range 修改

range 改变inner scope: 
```
{{- range .Values.pizzaToppings }}
- {{ . | title | quote }}
{{- end }}  
```

>  Each time through the loop, . is set to the current pizza topping. That is, the first time, . is set to mushrooms. The second iteration it is set to cheese, and so on.






# Pipeline

> A pipeline may be "chained" by separating a sequence of commands with pipeline characters '|'. In a chained pipeline, the result of each command is passed as the last argument of the following command. The output of the final command in the pipeline is the value of the pipeline.