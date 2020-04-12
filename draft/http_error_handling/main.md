
## 选定一个使用

juju/errors


说清楚思路: 

*  error chain
*  通过cause找出root cause
*  通过type assertion? Is or Interface Assertion?
*  模型: 一个error root cause, 层层trace. root cause需要aware http: 分配http code, header, error code.
*  HTTP
    -  标准化 error:
        ```
        HTTP CODE
        HEADER xxx
        {
            code: ""
            message: ""
        }
        ```
    -  HTTP Code
    ```
    type HasHTTPStatus interface {
        // HTTPStatus returns a HTTP status of an error.
        HTTPStatus() int
    }
    ```
    - HTTP HEADER
    ```
    type HasResponseHeaders interface {
        ResponseHeaders() map[string]string
    }
    ```
    -  区分是否server error: 根据是否有HTTP Status
    -  提供HTTP status:
        ```
        type HasHTTPStatus interface {
            // HTTPStatus returns a HTTP status of an error.
            HTTPStatus() int
        }
        ```
*  嵌入HTTP
    -  gocore/mware/handler.go WebHandler

问题:

*  是否应该每次Trace, 中断会怎么样

用Go 1.13的error 实现

两种Wrap: fmt.Errorf 或者　自定义结构: 实现Unwrap

? Cause, Is Has


No use juju/errors