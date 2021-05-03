---
title: 'Head First Golang Sort'
date: 2017-09-04 19:59:19
tags: [golang]
keywords: golang sort 总结
---


<div align="center">
{% asset_img sort.jpg %}
</div>

#### Preface 

sort作为一项非常基础的需求，对于编程(尤其后端)重要性不言而喻， 本文由日常sort需求出发，分析了go sort编程原理及各需求的go语言的实现. 本文的目标是介绍go中sort工作原理，同时可以作为go sort速查. 文章来源于作者日常编程总结和go语言sort文档. 但是document形式的文章不太利于理解，因此本文采用了head first的方式，先提出了我们最关注的 **如何实现xx功能**，分析go语言的sort支持，给出sort常见需求的最佳实践.

<!-- more -->

#### 常见的sort需求

*  排序ints，floats，strings: [升序](#sort_embedded)，[降序](#sort_desc_order)
*  排序struct array
    -  [固定pattern排序](#sort_static)(即静态:无法运行时修改，修改排序算法需要修改Interface实现)
    -  [不同的struct member排序(动态)](#sort_dynamic)
    -  [多关键字排序(动态)](#sort_multikeys)
*  [stable排序](#sort_stable)

本文的目标之一便是给出上述问题的go最佳实践

#### <a id="sort_basic">golang sort使用原理</a>

本文描述的golang sort针对slice
*  内置类型int，float的slice: golang sort包内置对其支持: 提供api对其进行升序，降序sort
*  struct array的sort: 需要实现Interface，然后调用sort.Sort(&mySortInterface)进行排序，Interface定义：
    ```go
    type Interface interface {
        // Len is the number of elements in the collection.
        Len() int
        // Less reports whether the element with
        // index i should sort before the element with index j.
        Less(i， j int) bool
        // Swap swaps the elements with indexes i and j.
        Swap(i， j int)
    }
    // go基于Interface提供的sort功能
    func IsSorted(data Interface) // 判断排序是否满足Less函数
    func Sort(data Interface)　　　　　// 不稳定排序
    func Stable(data Interface)   // 稳定排序
    ```
    sort.Sort内部主要用quickSort，并调用Interface实现的Less，Swap函数来进行实际的排序.

以上就是golang sort的原理，很简单吧? 基本上来说，我们可以通过:基于我们的struct slice实现Interface的Len，Swap，Less函数来实现排序的目的.但针对一些复杂的需求:如　**动态改变排序算法**，**多关键字排序**等如何优雅的实现呢？ 本文以Interface来指实现了Less，Swap，Len的interface.


#### sort最佳实践

##### sort 内置类型数组 ints，floats，strings <a id="sort_embedded"></a>

对一个具体的内置类型，以int数组为例，golang内置对其提供了与Interface类似的一组两个函数: IsSorted，Sort(内置类型不需要stable)

```golang
func IntsAreSorted(a []int)bool
func Ints(a []int) // 
```
简单的例子:
```golang
arr := []int{2， 1， -3， 1， 0}
sort.Ints(arr)
fmt.Println(arr)
// 结果 [-3 0 1 1 2]
```

sort包还提供float64，string slice的排序实现
```golang
// float64s
func Float64AreSorted(s []float64)bool
func Float64s(s []float64)
// strings
func StringsAreSorted(a []string)bool
func Strings(a []string)
```

<a id="sort_desc_order"></a>
升序排列实现了，降序呢?一句话:用sort.Reverse(Interface)来反转Interface的Less函数，之后变和升序一样，调用sort.Sort函数进行排序:
```golang
// Less returns the opposite of the embedded implementation's Less method.
func (r reverse) Less(i， j int) bool {
    return r.Interface.Less(j， i)
}
```
降序代码示例
```golang
sort.Sort(sort.Reverse(Interface))
```

以Ints为例:
golang内部实现了sort.IntSlice，其实就是用[]int实现了Interface
```go
type IntSlice []int
func (IntSlice)Len()int
func (IntSlice)Less(i， j int)bool
func (IntSlice)Swap(i， j int)
```
示例代码:
```golang
arr := []int{2， 1， -3， 1， 0}
sort.Sort(sort.Reverse(sort.IntSlice(arr)))
fmt.Println(arr)
```

##### sort struct slice 

示例数据集说明: 自定义struct，存储个人基本信息: 名字，年龄，体重(kg)，这是一家三口人:
```golang
type Person struct{
    name string
    age int
    weight float
}
var people = []Person {
    {"Frank"， 30， 70.5}，
    {"Lisha"， 30， 55.6}，
    {"TheBaby"， 0， 5.3}， 
}
```

###### 按固定模式排序 <a id="sort_static"></a>

最常用情形为按照某个member进行排序. 按照之前介绍的[golang排序原理](#sort_basic)，基于struct slice实现Interface，即可以实现按单个member的排序

```golang
type PeopleSorter []Person
func (ps PeopleSorter)Len()int {
    return len(ps)
}
func (ps PeopleSorter)Swap(i， j int) {
    ps[i]， ps[j] = ps[j]， ps[i]
}
func (ps PeopleSorter)Less(i， j int)bool {
    return ps[i].name < ps[j].name
}
// 调用
sort.Sort(PeopleSorter(people))
fmt.Println(people)
// 结果
[{Frank 30 70.5} {Lisha 30 55.6} {TheBaby 0 5.3}]
```

想要multiple key排序，只需要更改Less，如先按照Age，再按照Weight:
```
func (ps PeopleSorter)Less(i， j int)bool {
    if ps[i].age == ps[j].age {
        return ps[i].weight < ps[j].weight
    }
    return ps[i].age < ps[j].age
}
// 结果  [{TheBaby 0 5.3} {Lisha 30 55.6} {Frank 30 70.5}]
```

但是用此实现方式实现如下功能不可行:

1.  同时支持按照Name或Age或Weight排序: 上述实现不能共存于一个Interface
2.  Person所有member全排列排序，按上面所说的方法就得实现6个Less函数... :sob: 

很自然的解决办法就是在Sort之前，动态传入Less函数，有以下两种模式分别解决1， 2问题

###### 按single struct member动态排序 <a id="sort_dynamic"></a>

目标是实现如下效果：　每次调用可以传入不同的conditions
```golang
Sort(people).By(Age)
Sort(people).By(Weight)
... // 可以扩展其他condition
```

实现思路:
> 因为sort.Sort最终实现排序，因此我们需要实现我们自己可以传入Less function的Sort，并在其内部动态构造一个Interface，用其Less接口来调用传入的比较函数， 实现最终排序行为.

1.  接收自己Less函数的Sorter，实现调用内置的by来实现比较
    ```go
    type PeopleSorter struct {
        people []Person
        by func(a， b *Person)bool
    }
    func (ps PeopleSorter)Len()int {
        return len(ps.people)
    }
    func (ps PeopleSorter)Swap(i， j int) {
        ps.people[i]， ps.people[j] = ps.people[j]， ps.people[i]
    }
    func (ps PeopleSorter)Less(i， j int)bool {
        return ps.by(&ps.people[i]， &ps.people[j])
    }
    ```
2.  实现外部调用接口，以便可以这么调用 By(Age).Sort(people)
    ```go
    type By func(a， b *Person)bool的，因为没有带被sort的slice
    func (by By)Sort(people []Person) {
        sort.Sort(
            &PeopleSorter{
                people : people，
                by : by，
        })
    }
    ```
３.  调用演示
    ```go
    Age := func(a， b *Person)bool {
        return a.age < b.age
    }
    Name := func(a， b *Person)bool {
        return a.name < b.name
    }
    By(Age).Sort(people)
    fmt.Println(people)
    // 结果：　[{TheBaby 0 5.3} {Lisha 30 55.6} {Frank 30 70.5}]
    By(Name).Sort(people)
    fmt.Println(people)
    // 结果: [{Frank 30 70.5} {Lisha 30 55.6} {TheBaby 0 5.3}]
    ```

**注**：　由于需要外部传入比较函数，不能直接传 func (i， j int)bool的，因为没有带被sort的slice;因此需要传入func (a， b *Person)bool， Less实现时再转换一下即可

###### 按multiple struct member动态排序 <a id="sort_multikeys"></a>

目标是实现如下效果：　conditions的次序和数量可以任意指定

```go
OrderBy(Age， Name， Weight...).Sort(people)
```

实现思路:

> 以slice形式接收Less函数作为OrderBy的参数，遍历此函数对象数组。OrderBy传入Less函数，Sort传入待排序struct slice，两者组合成为一个Interface，再调用Interface自己实现的Sort函数来最终调用sort.Sort(sorter)来实现排序

1.  最终排序用的Interface: MultiSorter
    ```go
    type LessFunc func (a， b *Person)bool
    type MultiPeopleSorter struct {
        people []Person     // 待排序数组
        lessAr []LessFunc   // less函数array
    }
    func (pms *MultiPeopleSorter)Len()int {
        return len(pms.people)
    }
    func (pms *MultiPeopleSorter)Swap(i， j int) {
        pms.people[i]， pms.people[j] = pms.people[j]， pms.people[i]
    }
    // 多关键字排序的核心实现: 
    //    遍历less函数数组，用当前less比较 a，b
    //    若 a != b ，则返回当前less比较结果
    //    否则继续比较，若全部less无法判定，则返回false
    func (pms *MultiPeopleSorter)Less(i， j int)bool {
        a， b := &pms.people[i]， &pms.people[j]
        for _， less := range pms.lessAr {
            switch {
            case less(a， b):
                return true
            case less(b， a):
                return false
            default:
                // continue
            }
        }
        return false
    }
    ```
2.  OrderBy实现:
    ```go
    // 生成Sorter，传入less函数数组
    func OrderBy(lessArr ...LessFunc) *MultiPeopleSorter {
        return &MultiPeopleSorter{
            lessAr : lessArr，
        }
    }
    ```
3.  Interface自己Sort实现:
    ```go
    func (pms *MultiPeopleSorter)Sort(people []Person) {
        pms.people = people　// 传入待排序数组
        sort.Sort(pms)      // 实际排序行为
    }
    ```

这两个实现思路有点绕，这么理解起来比较容易:
1.  都是实现如下风格的调用: 
    ```go
    By(xxx).Sort(people)
    ```
2.  最终是为了sort.Sort传入Interface， 之前的所有代码都是为这个目的而做准备

##### stable sort <a id="sort_stable"></a>

很简单，只要把sort.Sort换成sort.Stable即可。
我们都知道： quicksort是不稳定排序，而mergesort是稳定排序，查看go源码:

```go
func Sort(data Interface) {
    ...
    quickSort(data， 0， n， maxDepth)
}
func Stable(data Interface) {
    ...
        if m := a + blockSize; m < n {
            symMerge(data， a， m， n)
        }
    ...
}
```

#### 为什么是 < 而不是 <=

这个问题的实质是对collection做search时，find函数中判断两个对象*equal*如何实现：一般可用**相等**和**等价**来实现. 等价实现基于Less函数，当下列条件满足时认为两对象等价:

```go
!Less(a， b) && !Less(b， a)
```

>a不小于b而且b不小于a同时成立则认为两对象等价

出问题的地方在于用**等价**实现，假如Less用<=实现:

```go
func Less(a， b *Person)bool {
    return a.age <= b.age
}
// 下列自身和自身的判断会返回false，即自己和自己不等：
!Less(a， a) && !Less(a， a) 
```
Less(a， a)会返回true，!Less(a， a)会返回false，false && false == false， 表明a和a不等价!

因此推荐在日常写Less函数的时候坚持: <而非<=，可以避免很多麻烦

关于equal和equivalence，在*Effective STL*一书的以下两节讲的很清楚，推荐阅读:

*  Item 19. Understand the difference between equality and equivalence
*  Item 21.Always have comparison functions return false for equal
values.


#### 写在最后

对比其他语言，golang提供的sort不算最方便，比如最简单的按照某一member对slice排序，需要写不少代码，不像STL中直接可以向collection传入lambda.
但本人认为这个其实是因为go不是一门大而全的重型语言，没有泛型的支持， 体现在目前sort接口，更像是go的一个feature而非缺陷。针对日常需求，目前用go还是可以较为便捷的实现。


#### 参考文献

[Package Sort](https://golang.org/pkg/sort/#Interface)　　

Scott Meyers， Effective STL: 50 Specific Ways to Improve Your Use of the Standard Template Library