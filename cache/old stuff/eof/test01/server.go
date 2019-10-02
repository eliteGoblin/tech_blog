/* SimpleEchoServer
 */
package main

import (
    "net"
    "fmt"
    "os"
)

func main() {

    service := ":1201"
    tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
    checkError(err)

    listener, err := net.ListenTCP("tcp", tcpAddr)
    checkError(err)

    for {
        conn, err := listener.Accept()
        if err != nil {
            continue
        }
        handleClient(conn)
        conn.Close() // we're finished
    }
}

func handleClient(conn net.Conn) {
    var buf [256]byte
    for {
        n, err := conn.Read(buf[0:])
        if err != nil {
            fmt.Println("err got ", err)
            return
        }
        fmt.Printf("content got %s len is %d\nil", string(buf[0:]), n)
        _, err2 := conn.Write(buf[0:n])
        if err2 != nil {
            return
        }
    }
}

func checkError(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
        os.Exit(1)
    }
}