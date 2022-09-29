package main

import (
	"context"
	"fmt"
)
func getRequestIdFromContext(ctx context.Context)string {
		if v := ctx.Value("requestId"); v != nil {
			return v.(string)
		}
		return ""
}

func doSomething(ctx context.Context) {
  fmt.Println("enter doSomething")
  fmt.Printf("requestId: %s\n", getRequestIdFromContext(ctx))
}

func main() {
	ctx := context.WithValue(context.Background(), "requestId", "uuid-1234")
	doSomething(ctx)
	// Output:
	// requestId: uuid-1234
}