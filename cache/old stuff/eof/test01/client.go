

/* GetHeadInfo
 */
package main

import (
    "net"
    "os"
    "fmt"
    "time"
)

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintf(os.Stderr, "Usage: %s host:port ", os.Args[0])
        os.Exit(1)
    }
    service := os.Args[1]

    tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
    checkError(err)

    conn, err := net.DialTCP("tcp", nil, tcpAddr)
    checkError(err)

    sContent := []byte("HEAD / HTTP/1.0\r\n\r\n")
    for i := 0; i < len(sContent); i ++ {
        _, err = conn.Write(sContent[i:i+1])
        checkError(err)
        time.Sleep(5 * time.Second)
    }
    
    conn.Close()
    fmt.Println("client conn closed")
    time.Sleep(10 * time.Second)
}


func checkError(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
        os.Exit(1)
    }
}