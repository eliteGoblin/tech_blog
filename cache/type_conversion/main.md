被unsigned 和signed转换搞得很烦

```
// -1 overflow uint
// fmt.Println(-1, uint(-1))
a := -1
fmt.Println(uint(a))
```