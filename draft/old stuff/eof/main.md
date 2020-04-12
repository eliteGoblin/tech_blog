
#### 基本问题

*  TCP发送端拆包实践 test01
    +  每次1byte, 发送20次，　server端收到多少次?  expect: <= 20; expect > 1循环;
    +  结果: 符合预期: client每次发送1byte, server端每次收到1byte 
    

*  读取nbyte, 什么时候返回? 
    -  < n: 会返回吗？
        +  无EOF
            *  预期: 不会
        +  包含EOF: 如expect 10byte, 收到5byte+EOF: 
            *  folder: test02
            *  server read 256, 但是client send 32bytes; server的read仍会返回，应该是tcp包收到就返回，看read返回值的到读取的个数. 
    -  >= n: 返回: 多次取



#### 思路

*  EOF in file
    -  二进制查看文件
*  EOF in TCP(HTTP response)



#### general

*  UNIX mac: CTRL + D
*  eof == -1
    ```c++
    #include <stdio.h>
    int main()
    {
        char ch;
        ch = getchar();
        printf("Input Char Is :%d",ch);
        return -1;
    }
    ```

#### EOF in golang

```golang
package "io"
// EOF is the error returned by Read when no more input is available.
// Functions should return EOF only to signal a graceful end of input.
// If the EOF occurs unexpectedly in a structured data stream,
// the appropriate error is either ErrUnexpectedEOF or some other error
// giving more detail.
var EOF = errors.New("EOF")
```

##### 接受和EOF关系

*  conn.Close()之前server端不会收到EOF,没办法手动