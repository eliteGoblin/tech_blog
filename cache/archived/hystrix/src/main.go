package main

import (
  "fmt"
  "github.com/afex/hystrix-go/hystrix"
  "net/http"
  "net"
  "errors"
)

func main() {
  hystrixStreamHandler := hystrix.NewStreamHandler()
  hystrixStreamHandler.Start()
  go http.ListenAndServe(net.JoinHostPort("127.0.0.1", "81"), hystrixStreamHandler)
  hystrix.ConfigureCommand("hello hystrix", hystrix.CommandConfig{
    Timeout:               1000,
    MaxConcurrentRequests: 10,
    ErrorPercentThreshold: 25,
  })
  //for i:= 0; i < 1000000; i ++ {
  //  hystrix.Go("hello hystrix", func() error {
  //    client := &http.Client{}
  //    req, _ := http.NewRequest("GET", "http://localhost:8000", nil)
  //    resp, err := client.Do(req)
  //    fmt.Println("resp got", resp, err)
  //    return nil
  //  }, func(err error) error {
  //      // do this when errors occur
  //      fmt.Println("hello from error ", err)
  //      return nil
  //  })
  //  time.Sleep(time.Millisecond * 100)
  //}

	ch := hystrix.Go("hello hystrix", func() error {
	  //client := &http.Client{}
	  //req, _ := http.NewRequest("GET", "http://localhost:8000", nil)
	  //resp, err := client.Do(req)
	  //fmt.Println("resp got", resp, err)
	  return errors.New("code 500 got")
	}, func(err error) error {
		// do this when errors occur
		fmt.Println("fallback called because ", err)
		return fmt.Errorf("fallback returned error %s", err)
	})

	err := <- ch
	fmt.Println("err before exit", err)
  
  select{}
}
