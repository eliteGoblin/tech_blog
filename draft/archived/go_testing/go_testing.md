
*  primarily a design activity with future correctness as a side effect
*  文件结构: 放在同一package下
*  go test
    -  在当前目录查找所有 *_test.go , -v 打印详细信息
    -  最小单元是单个test函数
        +  TestXxx(*testing.T):必须以Test开头,且后后一个单词首字符必须大写，Testnumber 是不会执行的
        +  用Error, Fail 指示测试失败
        +  go test -cover 来显示coverage
            ```
            go test -coverprofile=coverage.out 
            go tool cover -func=coverage.out // 显示函数级的
            go tool cover -html=coverage.out // 网页方式展示
            ```
    -  在当前目录递归执行go test: [这里](https://github.com/stretchr/gorc)
    -  [TableDrivenTests](https://github.com/golang/go/wiki/TableDrivenTests)
    -  benchmark:
        +  命名: BenchmarkXxx(*testing.B)
        +  go test -bench regexp
        +  bench一般格式
            ```
            func BenchmarkHello(b *testing.B) {
                for i := 0; i < b.N; i++ {
                    fmt.Sprintf("hello")
                }
            }
            ```
        +  RunParallel 来同时运行多个function, 需用 go test -cpu cpu_num来运行
            ```
            func BenchmarkTemplateParallel(b *testing.B) {
                b.RunParallel(func(pb *testing.PB) {
                    for pb.Next() {
                        ...
                    }
                })
            }
            ```
    -  examples:
        +  是用于Docmentation目的go code;
            *  go test被运行且被验证
            *  可在godoc web page点击run按钮运行
            *  API中含有可运行的code的好处时不会过时
            *  不含output的example code不会被执行，但是会编译
        +  用法:
            *  匹配特定顺序输出:
                ```
                func ExampleSalutations() {
                    fmt.Println("hello, and")
                    fmt.Println("goodbye")
                    // Output:
                    // hello, and
                    // goodbye
                }
                ```
                -  函数申明: ExampleXXX()
                    +  example of a package:    Example()
                    +  function F :             ExampleF()
                    +  a type T :               ExampleT()
                    +  method M on type T:      ExampleT_M()
                    +  多个example用 ExampleXXX_suffix()
            *  匹配unordered的输出: Unordered output
                ```
                func ExamplePerm() {
                    for _, value := range Perm(4) {
                        fmt.Println(value)
                    }
                    // Unordered output:
                    // 4
                    // 2
                    // 1
                    // 3
                    // 0
                }
                ```
*  执行test
    -  在package目录下: go test
    -  在外部： 指定package路径: go test /path_to_package/package_name
    -  支持short模式，外部用 -test.short flag 内部用testing.Short()获取
        +  一般short模式，可以调用t.Skip()来跳过某些case


示例: leetcode 771, 22

代码
```
func numJewelsInStones(J string, S string) int {
    ret := 0
    mp := make(map[rune]bool)
    for _, jewel := range J {
        mp[jewel] = true
    }
    for _, stone := range S {
        if _, ok := mp[stone]; ok {
            ret ++
        }
    }
    return ret
}
```
测试
```
func TestZeroValue(t *testing.T) {
    num := numJewelsInStones("", "")
    if num != 0 {
        t.Errorf("Jewery %s, Stone %s, expected %d, got %d", "", "", 0, num)
    }
}

var numJewelsInStonesTests =  []struct {
    jewels string
    stones string
    result int
} {
    {"", "", 0},
    {"aA", "aAAbbbb", 3},
    {"z", "ZZ", 0},
}


func TestNumJewelsInStones(t *testing.T) {
    for _, e := range numJewelsInStonesTests {
        if e.result != numJewelsInStones(e.jewels, e.stones) {
            t.Errorf("%+v: expected %d, got %d", e, e.result, numJewelsInStones(e.jewels, e.stones))
        }
    }
}
```


###### refs


[Golang basics - writing unit tests](https://blog.alexellis.io/golang-writing-unit-tests/)
[Testable Examples in Go](https://blog.golang.org/examples)
[Test-driven development with Go](https://leanpub.com/golang-tdd/read#leanpub-auto-wrapping-up-2)
[testify package](https://github.com/stretchr/testify#installation)
[goreportcard](https://github.com/gojp/goreportcard)
["Dependency Injection" in Golang](http://openmymind.net/Dependency-Injection-In-Go/)
[The cover story](https://blog.golang.org/cover)