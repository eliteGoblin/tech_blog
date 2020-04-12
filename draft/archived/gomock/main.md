

#### general

*  expect无顺序, 想有顺序
    -  局部顺序: after
    -  指定全部顺序: gomock.InOrder
*  有expect一定得调用，不是optioal的
*  expect指定mock object的全部调用case, 调用的没有被expect也是fail



```shell
mockgen -destination=redis_helper_mock/redis_helper_mock.go -package=redis_helper_mock frank/redis_helper RedisHelper
```


#### 


