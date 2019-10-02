Go maps in action
https://blog.golang.org/go-maps-in-action


map是否是reference
[Are maps passed by value or by reference in Go?](https://stackoverflow.com/questions/40680981/are-maps-passed-by-value-or-by-reference-in-go)

答案这句话不对: Map types are reference types, like pointers or slices[1]

但是map地址不一样　https://play.golang.org/

```
package main
import (
    "fmt"
)

func main() {
    mp := make(map[string]string)
    mp2 := mp
    mp3 := mp2
    fmt.Printf("mp %p, mp2 %p mp3 %p\n", &mp, &mp2, &mp3)
    fmt.Printf("mp %v, mp2 %v\n", mp, mp2)
    mp2["test"] = "test"
    fmt.Printf("mp %v, mp2 %v\n", mp, mp2)
}
// 输出
mp 0x1040c128, mp2 0x1040c130 mp3 0x1040c138
mp map[], mp2 map[]
mp map[test:test], mp2 map[test:test]
```

二者地址不一样，但是貌似指向的是同一块地址，一个改了，其他也改了

map分配的内存用完，会发生reallocate么?  如果会，那之前指向map的指针(golang的map变量)是否会失效？