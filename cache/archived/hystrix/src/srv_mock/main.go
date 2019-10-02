
package main

import (
  "io"
  "fmt"
  "time"
  "net/http"
)

var flag bool

func hello(w http.ResponseWriter, r *http.Request) {
  fmt.Print("hello recved\n")
  if flag == false {
    time.Sleep(10 * time.Second)
  }
  w.WriteHeader(404)
  io.WriteString(w, "Hello !")
}


func main() {
  http.HandleFunc("/", hello)
  go func(){
    time.Sleep(10 * time.Second)
    fmt.Println("server do not sleep")
    flag = true
  }()
  http.ListenAndServe(":8000", nil)
}
