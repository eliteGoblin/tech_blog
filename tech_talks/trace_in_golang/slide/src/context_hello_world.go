package main

import (
	"fmt"
	"context"
	"time"
	"errors"
)

// START OMIT
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context){
	  select {
	    case <-ctx.Done():
	       fmt.Printf("context ended before finished")
	    case <-time.After(time.Minute):
	       fmt.Println("function finished")
	  }
	}(ctx)
	err := somethingReturnAError()
	if err != nil {
	  cancel()
	}
	select{}
}
func somethingReturnAError() error {
  return errors.New("error occured")
}
// Output:
context ended before finished
// END OMIT
