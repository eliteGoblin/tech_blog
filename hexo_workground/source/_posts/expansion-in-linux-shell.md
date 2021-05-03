---
title: Linux Shell Expansion and Regular Expressions
date: 2018-12-14 14:00:04
tags: [linux RegularExpressions]
keywords:
description:
---

## Preface

我们在初学shell时，经常会纠结于参数分割，是否需要加引号，单引号还是双引号，有时单引号工作，有时又得双引号。稍不注意，就会得到很unexpect的结果。而经常也会感叹于变量引用的灵活：$PATH $(date) \`ls\`，让人眼花缭乱。更别提在我们用grep，sed时，那让人目眩的一大串转义字符了，WTF is that?  

```
# search in vim
/([0-9]\{3\}) [0-9]\{3\}-[0-9]\{4\}
# sed to process 
sed 's/\([0-9]\{2\}\)\/\([0-9]\{2\}\)\/\([0-9]\{4\}\
)$/\3-\1-\2/' distros.txt
```

{% asset_img main.jpg %}　

<!-- more -->

这背后是由于shell对参数做的特殊转换机制起作用: shell expansion，也是一个较为tricky的特性，是造成shell恐惧症的一大源泉。  

本文中，我们将建立起shell expansion的知识框架，同时回答下列问题:  

*  shell expansion分类
*  单引号，双引号作用，分别解决什么问题
*  expansion在grep/egrep中的使用，并解释之前grep及sed正则表达式含义

希望在阅读完本文后，能让你对shell的使用更得心应手。

## Shell Expansion

### 参数分割

我们先看一下shell中对于参数是如何分割的，这是理解shell执行程序的第一步：

> Spaces, tabs, and newlines (linefeed characters) and treats them as delimiters [^1]

就是将空格，tab，换行算作delimter，一视同仁，不算做传入的参数。

### Pathname Expansion

Shell中的匹配路径的wildcard生效的过程，就是在expand wildcard，将其转换为一个或者一类pathname集合，与regular expression尽管有一些相似，是完全不同的机制，wildcard仅用作路径匹配。最常见*, ?, \[\]

```
# 打印/tmp全部文件名称
[me@ubuntu ~]$ echo /tmp/*
```

为什么叫expansion呢？\*字符会替换为/tmp目录下的全部文件名，一个\*字符得到一个集合，厉害吧。


同样~字符，等同于当前用户的home directory，称为Tilde Expansion。

Bash supports the following three simple wildcards:

>\* - Matches any string, including the null string  
>? - Matches any single (one) character.  
>[...] - Matches any one of the enclosed characters.


### Parameter Expansion

非常常见的一种形式，完成变量替换

```
[me@ubuntu ~]$ echo $USER
```

### Command Substitution

同样很常见，用一个命令的输出来expand表达式

```
# 括号形式
[me@ubuntu ~]$ echo $(ls)
[me@ubuntu ~]$ ls -l $(which cp)
# back quote形式
[me@ubuntu ~]$ ls -l `which cp`
```

### Arithmetic Expansion

format为: $((expression))，以四则运算为expansion。
```
[me@ubuntu ~]$ echo $((2 + 2))
[me@ubuntu ~]$ echo $(($((5**2)) * 3))
```

## Brace Expansion

通过规则，枚举等expand成为集合，分为两种：

*  comma separated list of strings，以逗号分割
*  range of integers and strings，以..来标明范围

```
[me@ubuntu ~]$ echo Front-{A,B,C}-Back
Front-A-Back Front-B-Back Front-C-Back
[me@ubuntu ~]$ echo Number_{1..5}
Number_1 Number_2 Number_3 Number_4 Number_5
[me@ubuntu ~]$ echo a{A{1,2},B{3,4}}b
aA1b aA2b aB3b aB4b
```

注意brace内{}，不能包含unqoted whitespace。  

什么是quote，为什么要quote呢？

## Quote

前面提到的各种expansion，由不同的具有特殊意义的char控制: **\$**, **{}**, `等。就像一词多义，这些字符被当做特殊字符，本身与其他字符结合一起有特殊的meaning。我们仅仅想输出这些普通字符怎么办呢：quote，即用引号来消除特殊字符的含义，有两种quote: double quote和single quote，single character quote。

> shell provides a mechanism called quoting to selectively suppress unwanted expansions

### Single Quote

To suppress all expansions，简单粗暴，被单引号quote的字符串的任何char都被当做普通char,失去special meaning。 

### Double Quote

被quote的字符失去special meaning，除了三个:  **\$**, **\\**, **\`**。

### Backslash Quote

和c语言中的转义字符类似，\\字符使接下来的字符失去特殊含义。除了一种情况：\\newline：作用是将newline从input中*消除*，或者说ignore，用在format script时，将一行script拆分为多行，但不会造成副作用。一个有意思的事实是：single quote包围的字符串，无法用backslash使字符串包含'字符。

## Quote in grep command

quote带给人的恐惧尤其表现在当你想在shell中应用正则时，那一长串酸爽的backslash字符，本身正则就容易让人费解，quote will make it worse。  

其实并不难，我们首先理解键入grep后的整个过程： 以 echo "abcc" | grep 'ab.*$' 为例

*  shell fork child process：　grep
*  在argv[1]中传入字符串参数 ab\.\*$
*  grep　以收到的字符串参数为pattern进行匹配

这是最常见的情况，single quote禁止所有expansion，被quote部分为正则字串，我们再试一个命令：

```
[me@ubuntu ~]$ echo "123" | grep '[0-9]{3}'
```

结果并没有匹配。为什么呢？

我们需要先停下来，搞懂一些基本知识：包括shell的metachar，正则metacharacter，知道了正则依赖的特殊字符，才能正确quote，避免常见的pitfall。

## Regular expression metacharacters

正则的metacharacter有：

```
^ $ . [ ] { } - ? * + ( ) | \
```

其实并不完全正确，正则以metacharacter集合不同分为两类：

### Extended Regular Expressions(ERE)

上面提到的metacharacter，也即我们一般最常用的，其实是已经被Extended过的，egrep和grep -E使用的正是ERE。

### POSIX Basic Regular Expressions(BRE)

仅识别 ^ $ . \[ \] * 作为metachar，其他的都当做literal，对比ERE，以下字符被识别为literal 

```
( ) { } ? + |
```

一个有趣，同时也比较容易在使用中造成confusion的事实是，BRE backslash quote了{}, \[\]，不会将其认为literal，而是等同于ERE语法的meta char：

> However (and this is the fun part), the “(”, “)”, “{”, and “}” characters are treated as
metacharacters in BRE if they are escaped with a backslash, whereas with ERE, preceding any metacharacter with a backslash causes it to be treated as a literal.

看一个栗子：

```
[me@ubuntu ~]$ echo "123" | grep '[0-9]{3}' --color=auto
# 无输出，符合预期，{ } 属于ERE
[me@ubuntu ~]$ echo "123" | grep '[0-9]\{3\}' --color=auto
123
# BRE quote了 { }使其被认作ERE
```


## Shell metacharacter

我们在shell运行命令时，先由shell expand parameters，然后再将expand的结果送入程序，因此我们应该避免正则的metacharacter被shell 错误的expand或者解析。尤其注意与正则重叠的字符：

```
|   // pipe
?   // Matching single char
*   // Matching any char
[ ] // Mating, eg: ls -l e[abc].txt
$   // Parameter expansion
( ) // Group Command
{ } // Brace Expansion, Parameter Expansion: ${var}
```

## Common tools's regular expressions

*  支持BRE
    -  grep
    -  vim
    -  sed[^3]
*  支持ERE
    -  grep -E
    -  egrep

## Conclusion

*  核心问题是shell运行命令时，如何将参数传递给命令子进程的。
*  shell的参数由space分割
*  在shell运行程序时，shell会对其metachar进行expand，expand后再将参数传入程序
*  bash的meta characters: 
```
>   Output redirection, (see File Redirection)
>>  Output redirection (append)
<   Input redirection
*   File substitution wildcard; zero or more characters
?   File substitution wildcard; one character
[ ] File substitution wildcard; any character between brackets
`cmd`   Command Substitution
$(cmd)  Command Substitution
|   The Pipe (|)
;   Command sequence, Sequences of Commands
||  OR conditional execution
&&  AND conditional execution
( ) Group commands, Sequences of Commands
&   Run command in the background, Background Processes
#   Comment
$   Expand the value of a variable
\   Prevent or escape interpretation of the next character
<<  Input redirection (see Here Documents)
```
*  正则分为BRE和ERE，各自的meta characters
    -  BRE: ^ $ . \[ \] * 
    -  ERE: ^ $ . \[ \] { } - ? * + ( ) | \
    -  BRE可以用\\将\[\] {}　识别为meta character，与ERE意义相同。
*  用egrep + single quote可以满足最常见的正则应用。vim，sed均采用BRE

## Further reading

[What characters do I need to escape when using sed in a sh script?](https://unix.stackexchange.com/questions/32907/what-characters-do-i-need-to-escape-when-using-sed-in-a-sh-script)  


[^1]: The Linux Command Line Chap7: Seeing The World As The Shell Sees It  
[^2]: 见[BASH Metacharacters and Their Meanings](http://www.angelfire.com/mi/genastorhotz/reality/computers/linux/bashmetachars.html)   
[^3]: [sed](http://users.monash.edu.au/~erict/Resources/sed/)