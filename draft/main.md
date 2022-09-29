
# KMS




## delve 调试

*  远程调试
*  调试container
*  集成IDE


## docker cmd vs entrypoint


## How to do load testing

k6s


## Golang containers

*  heap, PQ
*  LRU:
    -  leetcode
*  list
*  bitvector
*  set

[dropbox/godropbox](https://github.com/dropbox/godropbox/blob/master/container/set/set.go)


### HTTPS

目标:
*  Break into pieces
*  基本原理
    -  公钥私钥 alice bob? link?  阮一峰
        +  加密解密实战
        +  数字签名实战
    -  再一篇: PKI，certificate
*  实战:
    -  生成公钥 私钥
*  Tonybai: 
*  HTTPS: C语言HTTPS Client分析
*  load certificate, 或者不安全access


### command flag parsing

#### parsing with flag package

[spf13/pflag](https://github.com/spf13/pflag)
[Golang之使用Flag和Pflag](https://o-my-chenjian.com/2017/09/20/Using-Flag-And-Pflag-With-Golang/)

#### advanced flag: cobra

*  simlife
    -  -v --version show version (generate version from git?)
    -  -N --name name of person to sim
        +  default: John
    -  hello name 
        +  -m 3 : hello time
            *  how about just -m 
        +  --help
    -  shopping where
        +  -w window shopping
*  auto generate command

[Golang之使用cobra](https://o-my-chenjian.com/2017/09/20/Using-Cobra-With-Golang/)

### gRPC

*  embedded rpc in golang
*  gRPC
*  advanced gRPC? proxy? SPOF
*  ideas
    -  分系列文章
        +  gRpc介绍，环境安装，最简单 echo server实现
        +  介绍 unary, bi-directional 原理, 并附上例子
            *  乒乓游戏模拟系统
            *  Login(Person): record your name; if not, cannot play
            *  ListPlayers: list players you can play with
                -  use stream, one name per msg(Person), just for show, same as login use
                -  Person
                    +  id: fsun
                    +  addr
                    +  PingPong Ability
                        *  Attack
                        *  Defense
                        *  每次rand一个，限定范围
                        *  对垒的两人，此次 attack - defense > limit, 则游戏结束，产生赢家
            *  Play(Person)
                -  根据person名字, 调出信息
                -  每次server先发球(attack), client收到，rand一个，予以回应
                -  游戏结束时，打印prompt


第二篇:  gRPC 原理， 类型， 深入例子

# local k8s KIND




