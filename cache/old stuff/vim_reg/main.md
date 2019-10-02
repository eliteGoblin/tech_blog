
#### 需求&&目标

*  AND and OR NOT reg combination in VIM
*  稍微扩展，依照vim文档的三个阶段
    -  simple pattern
    -  complex case
    -  full

#### 参考文档

常见需求vim regex

先在Vim实现，几种基本需求， 然后引申至通用reg
[“And” in regular expressions `&&`](http://www.ocpsoft.org/tutorials/regular-expressions/and-in-regex/)
[Regular Expressions: Is there an AND operator?](https://stackoverflow.com/questions/469913/regular-expressions-is-there-an-and-operator)
[Vim Regex : How to search for A AND B NOT C](https://stackoverflow.com/questions/3883985/vim-regex-how-to-search-for-a-and-b-not-c
[Power of g](http://vim.wikia.com/wiki/Power_of_g)
[Fast jump to line that matches a regular expression](http://vim.wikia.com/wiki/Fast_jump_to_line_that_matches_a_regular_expression)
[Vim documentation: pattern](http://vimdoc.sourceforge.net/htmldoc/pattern.html)

\v的使用 简化操作

/\v(.*Frankie)&(.*Test)

Use of "\v" means that in the pattern after it all ASCII characters except
'0'-'9', 'a'-'z', 'A'-'Z' and '_' have a special meaning.  "very magic"

[Vim documentation: pattern](http://vimdoc.sourceforge.net/htmldoc/pattern.html)

或者

目标放在前面

：g方法,以及 :g/pattern/# 显示行号，然后用 : 跳转



需求:

*  grep 
*  vim 定位log
*  程序语言regex


###### 系列构思

*  regex in golang
    -  regex 简介: 范畴，行match
    -  match
        +  简单情况
        +  and, or, not
        +  zero-width assertion
            *   ^$\b
            *   look around: by Jeffrey Jeffson
    -  replace
    
*  regex in vim
*  regex in grep


[我眼里的正则表达式入门教程](http://www.zjmainstay.cn/my-regexp)
[正则表达式30分钟入门教程](http://deerchao.net/tutorials/regex/regex.htm#balancedgroup)