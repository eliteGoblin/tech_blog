
##### 关键问题

*  shell 是否转义
    -  ' 和 "的区别
    -  转义的范围和原理
*  grep egrep 等区别



[When grep “\\” XXFile I got “Trailing Backslash”](https://stackoverflow.com/questions/20342464/when-grep-xxfile-i-got-trailing-backslash)

[grep “+” operator does not work](https://askubuntu.com/questions/293148/grep-operator-does-not-work)


lookahead asserstion:
查找不符合格式的call_id
?! 必须
```
cat udesk_cti.log | grep -P '"call_id\\":\\"(?![0-9a-zA-Z]+)'
grep -P 'request\]((hupall (?!normal_clearing lin_uuid [a-z0-9]{8}-))|(uuid_[a-z]{3,9} (?![a-z0-9]{8}-)))' udesk_cti.log | less
```


[What is the difference between `grep`, `egrep`, and `fgrep`?](https://unix.stackexchange.com/questions/17949/what-is-the-difference-between-grep-egrep-and-fgrep)

[正则表达式的零宽度先行断言（Lookahead）和后行断言（Lookbehind](https://leongfeng.github.io/2017/03/10/regex-java-assertions/)